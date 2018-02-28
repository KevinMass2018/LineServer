# LineServer
The lineserver code to serve lines from a text file when client send request

How to build:

To build Linux based executables:
    ./build.sh linux   
To build MacBook based executables:
    ./build.sh  mac
    
How to run:
To get command arguments helper:
    ./run.sh -h
    
Usage of ./main:
  -filename string
    	complete filename includes path, a string (default "test_file.txt")
  -p int
    	port number, an integer (default 54321)

For example, if the server want to serving the file /users/peter/temp/file1.txt with the TCP port 12345, the running command is:
    ./run.sh  -p=12345 -file="/users/peter/temp/file1.txt"
    
The design considerations:

When the server starts, it preprocess the inpute file. Build up a line-mapping file to store the offset info for each line in the input file.  The offset info is stored as a 8-byte BigEndian binary number. The offset location of line k is (k-1)*8 in the line-mapping file.

The server then listen to the port, 12345, as above example. 

When client send "GET nnnn" request, the server process it, calculate the location (nnnn-1)*8, then open the line-mapping file, seek to that location, and find the 8-byte binary number, and retrieve this data as offset-file, which is the offset of that line in the original input file.   Then it Seek to this offset-file, and read the corrsponding line and send the line content back to client.

The server also handles client request like "SHUTDOWN", which will shutdown the server; "QUIT", which will disconnect this client, and stop the thread for this client.  In case client send garbage request, SERVER processes it, send back an "ERR" msg, and continue. 


Test Cases:

Manual test includes:  
   building a client to send "GET nnnn", "SHUTDOWN", "QUIT" and also error input like: random commands "GET 0", "GET <huge>", "abcde;gkj;  a'dgkj;h",
    Have multiple clients from same machine as server sending the above commands
    Have multiple clients from same and different machines as server sending the above commands to server
    
Automatic test includes:
    Building a multiple thread to automatically send "GET nnnn" concurrently, continuously,  on small file, and 1GB file, and 10GB file
    

Summary:

So far the code work pretty well in manual test cases.  

Further work and testing need to be done for concurrent/continuously sending client request from 100+ clients, and performance data will be updated to this repo when it is ready.

with 10GB file,  the system performs close to 1GB file, only about 5% latency difference in a limited test, could be due to other factors. 

Used package includes: bufio, os, net, strconv, encoding/binary, fmt, flag, time, strings, github.com/Arafatk/glot

During this work, the following WebPage is consulted:  
    Golang TCP socket programming,  
    Each of the above package website,  
    Golang mutiple thread processing, 
    Golang control channel programing, 
    How to build golang executables for different target,  
    How to pass command auguments to bash scripts,
    How to plot graph in golang
    And many others,...

Totally spend two weeks in after hours for this work



/**********************************/

client.go is the separate code written to test lineserver performance.  The way to start client is:

go run client.go <line-server-ip>:<line-server-port>

E.g.: go run client.go 10.20.30.40:12345

10.20.30.40 is the lineserver IP reachable by client.  12345 is the TCP port line server is listening to
