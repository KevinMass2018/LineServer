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
