package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	p2p "github.com/SadlifeCaw/DisysyMandatory2/p2p"
	"google.golang.org/grpc"
)

type server struct {
	p2p.UnimplementedNodeServer
}

//different states of a Node
const (
	OUT        = "OUT"
	REQUESTING = "REQUESTING"
	CRITICAL   = "CRITICAL"
)

var TIME_TO_WAIT = 10 * time.Second

var status = OUT
var ownPort string
var GrantedPermissionToEnter = false
var counter = 0

var knownPortsSlice = make([]string, 0)
var requestQueue = make([]string, 0)

var mu sync.Mutex

func main() {
	var file, _ = os.Open(os.Args[1])

	var reader = bufio.NewReader(file)
	var line, _ = reader.ReadString('\n')

	//read given ports
	line = strings.TrimSuffix(line, "\r\n")
	knownPortsSlice = strings.Split(line, ";")
	log.Printf("Other nodes exist at Ports: %s, %s", knownPortsSlice[0], knownPortsSlice[1])

	//trim new line on last port
	knownPortsSlice[1] = strings.TrimSuffix(knownPortsSlice[1], "\n")

	//read own port
	line, _ = reader.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")
	ownPort = line

	//set up logging to appropiate file
	logPath := GetLoggerFileName(ownPort)

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("-----------------------------------")
	log.Println("Logging to custom file on port", ownPort)

	//begin actual implementation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ListenToPort()
	go CheckQueue(ctx)

	//give other Nodes enough time to start up
	log.Printf("Waiting 10 seconds before beginning")
	time.Sleep(TIME_TO_WAIT)

	//keep trying to enter the critical section
	//should maybe be tied to some listener for change in acceptance counter
	for {
		time.Sleep(time.Second * 3)
		attemptEnter(ctx)
		time.Sleep(time.Millisecond * 100)
		leave(ctx)
	}
}

func ListenToPort() {
	lis, err := net.Listen("tcp", ownPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	p2p.RegisterNodeServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func GetLoggerFileName(port string) (filename string) {
	switch port {

	case ":2000":
		return "log_node1.txt"
	case ":2001":
		return "log_node2.txt"
	case ":2002":
		return "log_node3.txt"

	default:
		log.Printf("Error, illegal port while settign up Logger")
		return "" //why? because go :(
	}
}

//queue implementation, could be abstracted to a struct/it's own file in the p2p folder
func Enqueue(queue []string, port string) []string {
	return append(queue, port)
}

func Dequeue(queue []string) (string, []string) {
	if len(queue) != 0 {
		poppedElement := queue[0]
		queue = queue[1:]
		return poppedElement, queue
	} else {
		return "", queue
	}
}

//contionusly see if other Nodes have given acceptance for another Node to enter the critical section
func CheckQueue(ctx context.Context) {
	for {
		//only loop if this Node has NOT been granted permission from the other Nodes
		if !GrantedPermissionToEnter {
			var NextPortToEnter string

			mu.Lock()
			NextPortToEnter, requestQueue = Dequeue(requestQueue)
			mu.Unlock()

			//if no Node in queue, skip iteration
			if NextPortToEnter == "" {
				continue
			}

			NextPortToEnter = fmt.Sprintf("localhost%s", NextPortToEnter)

			conn, _ := grpc.Dial(NextPortToEnter, grpc.WithInsecure(), grpc.WithBlock())
			NextPortNode := p2p.NewNodeClient(conn)

			mu.Lock()
			GrantedPermissionToEnter = true
			mu.Unlock()

			//Increase this Node acceptance counter
			//needs to happen once for every known Node before being GrantedPermissionToEnter in the critical section
			NextPortNode.Accept(ctx, &p2p.AcceptMessage{})
		}
	}
}

//tell all known Nodes that this Node wants to enter
func attemptEnter(ctx context.Context) {
	status = REQUESTING

	for i := 0; i < len(knownPortsSlice); i++ {
		conn, err := grpc.Dial(knownPortsSlice[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		Node := p2p.NewNodeClient(conn)

		//enqueue this Node in other Node's queue
		Node.Request(ctx, &p2p.RequestMessage{Message: ownPort, Node: ownPort})

		//close connection to other Node
		conn.Close()
	}

	//check if this Node has been accepted by every other Node
	for {
		if counter == len(knownPortsSlice) {
			mu.Lock()
			counter = 0 //reset acceptance counter
			mu.Unlock()

			status = CRITICAL
			log.Printf("Node with port %s : ENTER", ownPort)

			for i := 0; i < len(knownPortsSlice); i++ {
				//create connection to other Node
				conn, _ := grpc.Dial(knownPortsSlice[i], grpc.WithInsecure(), grpc.WithBlock())
				Node := p2p.NewNodeClient(conn)

				//tell other Nodes this Node has entered the critical section
				Node.EnterCritical(ctx, &p2p.EnterMessage{Message: ownPort})

				//close connection to other Node
				conn.Close()
			}
			//break loop
			break
		}
	}
}

func leave(ctx context.Context) {
	status = OUT
	log.Printf("Node with port %s : LEAVE", ownPort)

	for i := 0; i < len(knownPortsSlice); i++ {

		//create connection to other Node
		conn, _ := grpc.Dial(knownPortsSlice[i], grpc.WithInsecure(), grpc.WithBlock())
		Node := p2p.NewNodeClient(conn)

		//inform other Node that this Node is leaving the critical section
		Node.LeaveCritical(ctx, &p2p.LeaveMessage{Message: ownPort})

		//close connection to other Node
		conn.Close()
	}
}

//when requesting to enter the critical section, enqueue this Node
func (s *server) Request(ctx context.Context, in *p2p.RequestMessage) (*p2p.RequestReply, error) {
	mu.Lock()
	requestQueue = Enqueue(requestQueue, in.Node)
	mu.Unlock()

	return &p2p.RequestReply{}, nil
}

func (s *server) Accept(ctx context.Context, in *p2p.AcceptMessage) (*p2p.AcceptReply, error) {
	mu.Lock()
	counter++
	mu.Unlock()

	return &p2p.AcceptReply{}, nil
}

func (s *server) LeaveCritical(ctx context.Context, in *p2p.LeaveMessage) (*p2p.LeaveReply, error) {
	log.Printf("Node with port %s : LEAVE", in.Message)
	mu.Lock()
	GrantedPermissionToEnter = false
	mu.Unlock()
	return &p2p.LeaveReply{}, nil
}

func (s *server) EnterCritical(ctx context.Context, in *p2p.EnterMessage) (*p2p.EmptyReply, error) {
	log.Printf("Node with port %s : ENTER", in.Message)
	return &p2p.EmptyReply{}, nil
}
