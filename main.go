package main

import (
       "fmt"
       "net"
       "flag"
       "strconv"
       "log"
       "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
     fmt.Println("/*******************************/")
     fmt.Println("            Line Server          ")
     fmt.Println("/*******************************/")

     filename  :=  flag.String("filename", "test_file.txt", 
     	      	  	      "complete filename includes path, a string")

     log_file  := flag.String("serverlogfile", "serverlogfile.txt",
                              "server log file name, a string") 

     portnum  :=  flag.Int("p", 54321, "port number, an integer")
    
     flag.Parse()

     /*
      * preprocessing the file and store the line offset data into a file
      */
     ret := line_svr_preprocessor(*filename)

     if ret != SERVER_SUCCESS {
     	return
     }     	

     /* Create server log file and log server output to this file */
     log.SetOutput(&lumberjack.Logger{
	Filename:   *log_file,
    	MaxSize:    1, // megabytes
    	MaxBackups: 5,
    	MaxAge:     30, //days
    	Compress:   false, // disabled by default
     })


     service := ":"+strconv.Itoa(*portnum)
     
     tcpAddr, err := net.ResolveTCPAddr("tcp4", service)

     if err != nil {
         checkError(err)
         return 
     }

     /*
      * Listen to the particular addr:port
      */
     listener, err := net.ListenTCP("tcp4", tcpAddr)
     if err != nil {
         checkError(err)
         return 
     }

     /*
      * Start a new thread to handle SHUTDOWN client request 
      */

     req_chan := make( chan string)
     go handle_shutdown_request(req_chan) 
     
     /*
      * Start a for loop for to support concurrent client requests
      */
  
     for {

     	 conn, err := listener.Accept()
	 defer conn.Close()

	 if err != nil {
	    checkError(err)
	    continue
	 }

	 /*
	  * Start a thread to handle request from this client
	  */
	 go handleClientRequest(conn, &req_chan, *filename)

     }
}       