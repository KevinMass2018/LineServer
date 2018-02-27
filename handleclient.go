package main

import (
       "os"
       "fmt"
       "net"
       "strings"
       "strconv"
       "bufio"
       "encoding/binary"
)

/*
 * Return values when processing the client request
 */
const  CLIENT_QUIT     =  0
const  CLIENT_INVALID  =  1
const  CLIENT_SHUTDOWN =  2
const  CLIENT_VALID    =  3
const  SERVER_FAIL     =  4
const  SERVER_SUCCESS  =  5

func handleClientRequest (conn net.Conn, req_chan  *chan string, filename string) {

     defer conn.Close()

     for {

     	 buf := make([]byte, 1024)  

         /* Read the received data into buf */
         _, err := conn.Read(buf[0:])

	 if (err != nil) {

	     checkError(err)
	     fmt.Println("S <= ERR")
             conn.Write([]byte("ERR"))
	     return
	 }

         result := string(buf[0:])

         /*
          *Split the received string into slice to get the line number
          */
         result_arr := strings.Fields(result)

         fmt.Println("C =>",result)

	 
	 ret := processRequest(conn, result_arr, req_chan, filename)	 	 	 
	 switch  ret  {
	     case CLIENT_QUIT :  // client QUIT

              	 fmt.Println("Client sent QUIT")
              	 fmt.Println("S <=  Connection Closed")
              	 conn.Write([]byte("Connection Closed"))    
		 return

	     case CLIENT_INVALID  :  // client invalid input
	     	 fmt.Println("S <= ERR")
              	 conn.Write([]byte("ERR"))
	         continue

	     case CLIENT_SHUTDOWN :  // client send shut
              	 /*
               	  * in case client sent SHUTDOWN request
	       	  * send it to the control channel and handle
               	  * it in a separate thread
               	  */
              	 fmt.Println("Client sent SHUTDOWN")
              	 fmt.Println("S <= Server Shutdown")
              	 conn.Write([]byte("Server SHUTDOWN"))

              	 *req_chan <- "SHUTDOWN"
	      	 return	 

	     case CLIENT_VALID   :
	     	  continue		 
		  
	     case SERVER_SUCCESS : 	  		  
		  continue

             case SERVER_FAIL    :  // continue the thread for now
                 fmt.Println("S <= ERR")
                 conn.Write([]byte("ERR"))
                 continue

	     default :
	         fmt.Println("S <= ERR")
              	 conn.Write([]byte("ERR"))
                 continue	 
 	 }

    }
}

/*
 * This thread handle specific client request, like SHUTDOWN server
 */

func handle_shutdown_request (req_chan chan string) {

     for {
     	 req := <-req_chan
	 
	 if strings.ToUpper(req) == "SHUTDOWN"  {

	    os.Exit(0)
	 }
     
     }
}

func processRequest (conn net.Conn, result_arr []string, 
     		     req_chan  *chan string, filename string) int {

     /*
      * client request length has to be 1 2, or 3
      * GET nnnn or SHUTDOWN, QUIT  + \r\n
      */

     req_len := len(result_arr)      

     if req_len > 3 || req_len <1   {
     	return CLIENT_INVALID
     }


     switch strings.ToUpper(result_arr[0]) {

	 case "SHUTDOWN" :
	      return CLIENT_SHUTDOWN

	 case "QUIT"     :
	      return	CLIENT_QUIT	       
		       
	 case  "GET"     :
	      return processGetRequest(conn, result_arr, filename)		  	

	 default :   // invalid keyword
	      return CLIENT_INVALID     
     }	     
}



func  processGetRequest  (conn net.Conn,  result_arr []string,
      			  filename string) int {


         /*
          * By protocol, the result_arr[1] stores the line number
          */
         line,err := strconv.Atoi(result_arr[1])

         /*
	  * Bail out in case input wrong
          */
         if  err !=nil {
             return CLIENT_INVALID
         }

         linecontent, ret := getContentbyLine(uint64(line), filename)

	 if ret != SERVER_SUCCESS {
	    return ret
	 }

	 /*
	  * Sent the line content back to client
	  */
         fmt.Println("S <= OK")
         fmt.Println("S <=",linecontent)

         var arr []byte
         arr = []byte(linecontent)
         _, err = conn.Write(arr)

      	 if err != nil {
            checkError(err)
            return SERVER_FAIL
      	 }

	 return CLIENT_VALID
}


/*
 * Input is the line number from client request, and the file name to be served
 * This function returns the content of the line as a string, and a status code
 */

func  getContentbyLine (orig_line uint64, filename string) (string, int) {

      /*
       * The line is stored in zero based system in line map file.
       * The client request is using 1 based system. Hence the conversion.
       */

      if orig_line < 1 || orig_line > total_line_num {
         return "ERR", CLIENT_INVALID
      }

      newLine := (uint64)(orig_line-1)

      /*
       * The location to store offset info of the corresponding line
       * is (line * 8) in the line_map file
       */

      offset_map := (int64)(newLine*OFFSET_BYTE_SIZE)

      mapfile, err := os.Open(line_map_file)

      if mapfile == nil || err != nil {
         checkError(err)
         return "ERR", SERVER_FAIL
      }

      defer mapfile.Close()

      /*
       * Move to the offset of the file to read the content
       */
      _,err  = mapfile.Seek(offset_map, 0)

      if err != nil {
         checkError(err)
         return "ERR", SERVER_FAIL
      }

      rd := bufio.NewReader(mapfile)

      if rd == nil {
         fmt.Println("Map file reader can not be initialized")
         return "ERR", SERVER_FAIL
      }

      /*
       * Get the offset value from the lineMap file
       * By design this is a 8 bytes value.
       */
      var offset_byte []byte
      var tmp         byte

      for i:=0; i < OFFSET_BYTE_SIZE; i++ {

          tmp, err = rd.ReadByte()

          if err != nil {
             checkError(err)
             return "ERR", SERVER_FAIL
          }
          offset_byte = append(offset_byte, tmp)
      }

      /* 
       * Convert the byte array into a uint64 number
       * with BigEndian byte-order
       */

      offset := binary.BigEndian.Uint64(offset_byte[0:])

      /*
       * Now open the to-be-served file
       */

      orig_file, err := os.Open(filename)
      if err != nil {
         checkError(err)
         return "ERR", SERVER_FAIL
      }

      defer orig_file.Close()

      _, err = orig_file.Seek(int64(offset), 0)

      if err != nil {
         checkError(err)
         return "ERR", SERVER_FAIL
      }


      rd_orig := bufio.NewReader(orig_file)

      if rd_orig == nil {
         fmt.Println("Original file reader can not be initialized")
         return "ERR", SERVER_FAIL
      }

      linecontent, err := rd_orig.ReadString('\n')

      if err != nil {
         checkError(err)
         return "ERR", SERVER_FAIL
      }


      return linecontent, SERVER_SUCCESS

}


/*
 * Error handling routine
 */
func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Encounting error: %s\n", err.Error())
    }
}

