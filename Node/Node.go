package main

import (
	"bufio"
	"context"
	"fmt"
	gh "github.com/SadlifeCaw/DisysyMandatory2/p2p"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type server struct {
	gh.UnimplementedNodeServer
}

const (
	OUT = "OUT"
	REQUESTING = "REQUESTING"
	CRITICAL = "CRITICAL"
)

var (
	status = OUT
	port string
	allowed = false
	ports = make([]string, 0)
	counter = 0
	requestQueue = make ([]string, 0)
	mu sync.Mutex
)

func main(){
	var file, _ = os.Open(os.Args[1])
	var reader = bufio.NewReader(file)
	var line, _ = reader.ReadString('\n')
	line = strings.TrimSuffix(line, "\r\n")
	ports = strings.Split(line, ";")
	log.Printf("Ports: %s, %s", ports[0], ports[1])
	line, _ = reader.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")
	port = line

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go listen()
	go queueHandler(ctx)

	var input string

	fmt.Scan(&input)
	time.Sleep(time.Second * 10)

	for {
		time.Sleep(time.Millisecond * 1) //Try with 10
		attemptEnter(ctx)
		time.Sleep(time.Second) //Try with 10
		leave(ctx)
	}
}

func listen(){
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	gh.RegisterNodeServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func enqueue(queue []string, str string) []string {
	return append(queue, str)
}

func dequeue(queue []string) (string, []string){
	if len(queue) != 0 {
		dequeued := queue[0]
		queue = queue[1:]
		return dequeued, queue
	} else {
		return "", queue
	}
}

func queueHandler(ctx context.Context) {
	for {
		if !allowed {
			var allowPort string
			mu.Lock()
			allowPort, requestQueue = dequeue(requestQueue)
			mu.Unlock()
			if allowPort == "" {
				continue
			}
			allowPort = fmt.Sprintf("localhost%s", allowPort)
			conn, _ := grpc.Dial(allowPort, grpc.WithInsecure(), grpc.WithBlock())
			c := gh.NewNodeClient(conn)
			mu.Lock()
			allowed = true
			mu.Unlock()
			c.Accept(ctx, &gh.AcceptMessage{})
		}
	}
}

func attemptEnter(ctx context.Context){
	status = REQUESTING
	for i := 0; i < len(ports); i++ {
		conn, _ := grpc.Dial(ports[i], grpc.WithInsecure())
		c := gh.NewNodeClient(conn)
		c.Request(ctx, &gh.RequestMessage{Node: port})
	}

	for {
		if counter == len(ports){
			mu.Lock()
			counter = 0
			mu.Unlock()
			status = CRITICAL
			log.Printf("Node with port %s has entered the critical state", port)
			for  i := 0; i < len(ports); i++{
				conn, _ := grpc.Dial(ports[i], grpc.WithInsecure(), grpc.WithBlock())
				c := gh.NewNodeClient(conn)
				c.EnterCritical(ctx, &gh.EnterMessage{Message: port})
				conn.Close()
			}
			break
		}
	}
}

func leave(ctx context.Context) {
	status = OUT
	log.Printf("Node with port %s has left the critical state", port)
	for i := 0; i < len(ports); i++ {
		conn, _ := grpc.Dial(ports[i], grpc.WithInsecure(), grpc.WithBlock())
		c := gh.NewNodeClient(conn)
		c.LeaveCritical(ctx, &gh.LeaveMessage{Message: port})
	}
}

func (s *server) Request(ctx context.Context, in *gh.RequestMessage) (*gh.RequestReply, error) {
	mu.Lock()
	requestQueue = enqueue(requestQueue, in.Node)
	mu.Unlock()
	return &gh.RequestReply{}, nil
}

func (s *server) Accept(ctx context.Context, in *gh.AcceptMessage) (*gh.AcceptReply, error){
	mu.Lock()
	counter++
	mu.Unlock()
	return &gh.AcceptReply{}, nil
}

func (s *server) LeaveCritical(ctx context.Context, in *gh.LeaveMessage) (*gh.LeaveReply, error) {
	log.Printf("Node with port %s has left the critical state", in.Message)
	mu.Lock()
	allowed = false
	mu.Unlock()
	return &gh.LeaveReply{}, nil
}

func (s *server) EnterCritical(ctx context.Context, in *gh.EnterMessage) (*gh.EmptyReply, error) {
	log.Printf("Node with port %s, has entered", in.Message)
	return &gh.EmptyReply{}, nil
}