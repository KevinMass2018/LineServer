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

This is output:    
Usage of ./main:
  -filename string
    	complete filename includes path, a string (default "test_file.txt")
  -p int
    	port number, an integer (default 54321)
  -serverlogfile string
    		 server log file name, a string (default "serverlogfile.txt")

For example, if the server want to serving the file /users/peter/temp/file1.txt with the TCP port 12345, and the log file in local directory as log.txt,  the running command is:
    ./run.sh  -p=12345 -file="/users/peter/temp/file1.txt" -serverlogfile="log.txt"
    
The design considerations:

When the server starts, it preprocesses the inpute file, and builds up a line-mapping file to store the offset info for each line of the input file.  The offset info is stored as a 8-byte BigEndian binary number. The offset location of line k is (k-1)*8 in the line-mapping file.

The server then listen to the port, 12345, as above example. 

When client sends "GET nnnn" request, the server processes it, calculates the location (nnnn-1)*8, then open the line-mapping file, seek to that location, and find the 8-byte binary number, and retrieve this data as offset-file, which is the offset of that line in the original file.   Then it opens the original file, Seeks to this offset-file, and read the corrsponding line and sends the line content back to client.

The server also handles client request like "SHUTDOWN", which will shutdown the server; "QUIT", which will disconnect this client, and stop the thread for this client.  In case client sends garbage request, SERVER processes it, sends back an "ERR" msg, and continue. 


Test Cases:

Manual test includes:  
   building a client to send "GET nnnn", "SHUTDOWN", "QUIT" and also error input like: random commands "GET 0", "GET <huge>", "abcde;gkj;  a'dgkj;h",
    Have multiple clients from same machine as server sending the above commands
    Have multiple clients from same and different machines as server sending the above commands to server
    Tested the server running in MacBook, and RHE Linux release 6.9.     


Automatic test includes

    During this test, the server is running in single-core MacBook, running OS X Yosemite Version 10.10.5.    
    2.2GHz Intel Core i7, 16GB RAM


    Building a client code which initializes multiple threads, and each thread automatically send "GET nnnn" concurrently, continuously,  on small 1GB file, and 10GB file
    

Summary:

   The server works pretty well in manual test cases.  

   The test results for multiple clients, multiple requests per seconds against 1GB/10GB file are captured in the following table.
   In this table, the data is the round-trip time,  measured in client side, from client sends a request, until it receives the response in different combinations:



                                      100 request/second                          1000request/second
                                      
     single client (1GB file )               885us                                          594us
     single client (10GB file)               658us                                          613us
     10-client     (1GB file)                1906us                                         4707us 
     10-client     (10GB file)               2002us                                         4686us
     



From the above table, with this algorithm, the file size only make difference in preprocessing.  It doesn't make too much difference in performance when line server serving client request (The single client and 100 request/second is an exception, with 1Gb 885us, and 10GB 658us, could be due to test environement at that moment). 

For single clients, increasing client request rate from 100 to 1000 doesn't make too much difference, probably because server CPU is not a limitation in either case.  Maybe increasing the request rate to 10000 req/s will reach the CPU limit and expose the difference?

For multiple clients case, increasing client request rate from 100 to 1000 make big difference in terms of round trip time. It's more than doubled. As you can see from the data, the server performance is downgraded with higher client request frequency.



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

Spent about 20 hours in this work.  The following is the time break down:

    4 hours learning golang

    4 hours design and implement the preprocessor, including went down a wrong road then with help from team get back on track

    2 hours implement basic client/server and TCP socket code

    2 hours make the server support concurrent clients requests, and continuous client requests

    2 hours corner cases test, boundary cases test , garbage input test, error handling and bug fixes

    2 hours test client code implemenation, data collection and graph plotting

    2 hours performance and scaling test, data analysis and bug fix 

    1 hour build scripting on different target

    1 hour publishing code and results to GitHub Repository



/**********************************/

client.go is the separate code written to test lineserver performance.  The way to start client is:

go run client.go <line-server-ip>:<line-server-port>

E.g.: go run client.go 10.20.30.40:12345

10.20.30.40 is the lineserver IP reachable by client.  12345 is the TCP port line server is listening to
