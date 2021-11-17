# DisysyMandatory2

## panic: open log_node1.txt: Access is denied.

If you get this error, you need to make sure the files aren't protected by your system, as Anti Virus programs sometimes likes blocking these kinds of things.

## Running the program

The program is designed to be used with 3 nodes exactly.

To run it, open 3 terminals and in each of them do

```
cd Node/
```

Then run the go file with an input file.
Below is an example with using node1.txt as input (use a different .txt file for each node, 1 to 3)

Be advised: It is important to be fast enough, as the program only waits for 10 seconds before it starts running. If all nodes are not present at this time, the program might not run correctly.

```
go run . ../node1.txt
```

**Important**: After a Node has been created by the command above, you have 10 seconds to activate every other node. If all Nodes are not initialized when the program properly starts, it will not function correctly.
If this is too little time, you can change the variable "TIME_TO_WAIT" in Node/Node.go
## Logging

Each Node logs its operations to its relevant log file, which can be found in /Node
A quick overview:

- log_node1.txt logs from Port 2000

- log_node2.txt logs from Port 2001

- log_node3.txt logs from Port 2002
