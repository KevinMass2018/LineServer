package main

import (
       "fmt"
       "bufio"
       "os"
       "encoding/binary"
)

/* line offset is stored as a 8-byte number in the line_map file */
const OFFSET_BYTE_SIZE = 8

/* 
 * File name to store the line-offset mapping info .
 * This system generated file locates in the same directory as executable
*/
const line_map_file = "line_map.txt"  


/*
 * The total lines of the to-be-served file
 */

var total_line_num uint64

/*
 * Preprocessing the to-be-served file and generate a line mapping file.
 * The starting offset of each line in the to-be-served file is calculated   
 * and stored as binary BigEndian number in a fixed location, i.e, (line-1)*8, 
 * in the mapping file.
 *
 * This way, when client request a particular line, the server can open this
 * mapping file, go to location: (line-1)*8, and retreive the offset. Then open the 
 * to-be-served file, and Seek to this offset, and retrieve the line, and 
 * send back to the client
 */


func line_svr_preprocessor(filename string) int {


     fmt.Println("Text File Preprocessor  Starts -----    ")

     /* Open the to-be-served file. Assume this file pre-exists */

     file, err := os.Open(filename)
     if err != nil {
     	checkError(err)
	return SERVER_FAIL
     }
     defer file.Close()

     /* Create the line map file */
     var mapFile *os.File
     mapFile, err = os.Create(line_map_file)
     if mapFile == nil || err != nil {
        checkError(err)
     	return SERVER_FAIL
     }


     /* Enable a file scanner */
     sc := bufio.NewScanner(file)

     offset     := uint64(0)
     line_num   := uint64(0)   

     /*
      * Scan the original file line by line
      */
     for sc.Scan() {

	 linefile := sc.Bytes()
	 linelen := len(linefile) + 1 // Add the "\n" character

	 /* 
	  * Binary encoding the offset number into a 8-byte array
	  * with BigEndian byte-order
	  */
	 offset_byte := make([]byte, OFFSET_BYTE_SIZE)
	 binary.BigEndian.PutUint64(offset_byte, offset)

         /*
          * Write offset value at the fixed position: line_number * 8
          */
	 _, err = mapFile.WriteAt(offset_byte, int64(line_num*OFFSET_BYTE_SIZE))

	 if err != nil {
     	     checkError(err)
             return SERVER_FAIL
     	 }

	 /* Increment the line number and the offset number */
	 offset   += (uint64)(linelen)
         line_num += (uint64)(1)

    }
    
    /* Save the total line number */
    total_line_num = line_num

    fmt.Println("Text File Preprocessor  Ends -----    ")
    return SERVER_SUCCESS
}
