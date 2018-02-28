package main

import (
       "net"
       "os"
       "fmt"
       "time"
       "strconv"
       "github.com/Arafatk/glot"

)


/******************************************************
 *
 * main function of the line server client module     *
 *
 ******************************************************/

const  X_RANGE = 5000    // control how many lines to request
const  Y_RANGE = 1500
const  X_SHIFT = 0       // control starting from which line
const  CLIENT_COUNT = 50 // control simulate how many clients
       

func main() {


     fmt.Println("/*******************************/")
     fmt.Println("    Line Server  Client          ")
     fmt.Println("/*******************************/")



     /*
      * Sanity check the command line argument number is valid
      */

     if len(os.Args) != 2 {
        fmt.Println ("0")
	os.Exit(1)
     }

     server := os.Args[1]

     for  {
       
        for i:=0; i<CLIENT_COUNT ; i++ {
            oneClientRequest(server)

	    /* Wait 10ms before starting another thread */
	    duration, _ := time.ParseDuration("10ms")
            time.Sleep(duration)
	}    
    }

}


func oneClientRequest (server string) {

     /*
      * Dial to the server
      */
     conn,err := net.Dial("tcp", server)


     if (err != nil ) {
        fmt.Println ("1")
     	return
     }
     
     /*
      * bail out if connection can not be established
      */
     if conn == nil {

     	fmt.Println("2")
     	return
     }

    defer conn.Close()


     sum := int64(0)
     line_idx := int64(1) + X_SHIFT
     var x,y  []float64 
     
     for {
        
	 text := "GET " + strconv.Itoa(int(line_idx))+" \n"    	 

	 fmt.Println (text)	
	 start := time.Now()
	 
	 conn.Write([]byte(text))
	
	 var buf [8196]byte 
	 conn.Read(buf[0:])

	 elapsed := time.Since(start)

         fmt.Println("Message Received, elapsed:", elapsed)
         fmt.Println(string(buf[0:]))

	 sum += int64(elapsed)

	 fmt.Println(line_idx)	 

	 x = append(x, float64(line_idx))
	 y = append(y, float64(elapsed)/1000)

	 duration, _ := time.ParseDuration("1ms")

	 
	 time.Sleep(duration)

        line_idx = (line_idx + 1)


	 if line_idx > X_RANGE + X_SHIFT   {
	    break
	 }


     }
          
     avg := int(sum/(X_RANGE))
     latency_plot (x, y, avg)

}


func  latency_plot(x, y []float64, avg int) {

     dimensions := 2
     // The dimensions supported by the plot
     persist := true
     debug := false
     plot, _ := glot.NewPlot(dimensions, persist, debug)

     avg = int(avg/1000)

     fmt.Println("average", avg)

     pointGroupName := "Average latency per request " + strconv.Itoa(avg) + " us"
     style := "lines"

     points := [][]float64{x, y}
        // Adding a point group
        plot.AddPointGroup(pointGroupName, style, points)
        // A plot type used to make points/ curves and customize and save them as an image.
        plot.SetTitle("Latency: 1000 req/s, Single thread, 10GB")
        // Optional: Setting the title of the plot
        plot.SetXLabel("line-number")
        plot.SetYLabel("latency(us)")
        // Optional: Setting label for X and Y axis
        plot.SetXrange(X_SHIFT, X_SHIFT+X_RANGE)
        plot.SetYrange(0, Y_RANGE)
        // Optional: Setting axis ranges
        plot.SavePlot("plots/plot.png")
}


