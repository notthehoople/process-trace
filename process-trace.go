package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	//"strings"
)

func processInput(debug bool) {
	var readResult int = 0
	var lowerBucket, upperBucket, cpuLatency int
	var timeCounter int = 0

	// Setup at process start

	bucketValueMap := make(map[int]int)
	scanner := bufio.NewScanner(os.Stdin)

	// We want the service to run forever, so loop forever
	for {

		for timeCounter = 0; timeCounter < 16; {
			// Scans a line from Stdin(Console)
			scanner.Scan()
			// Holds the string that scanned
			readText := scanner.Text()
			if len(readText) != 0 {
				if readText == "@usecs:" {
					fmt.Println("We found usecs!")
					timeCounter++
				} else {
					readResult, _ = fmt.Sscanf(readText, "[%d, %d) %d", &lowerBucket, &upperBucket, &cpuLatency)
					if readResult > 0 {
						fmt.Println("READ: ", readText)
						fmt.Println("lowerBucket: ", lowerBucket)
						fmt.Println("upperBucket: ", upperBucket)
						fmt.Println("cpuLatency: ", cpuLatency)

						bucketValueMap[lowerBucket] += cpuLatency
					} else {
						fmt.Println("FAILED: ", readText)
					}
				}
			}
		}

		// When we get here, we've read and processed 15 different inputs so it's time to summarise
		fmt.Println("timeCounter is:", timeCounter)
		fmt.Println(bucketValueMap)

		// Send the details to Graphite
		// Send data
		// Check number of times we've sent data to Graphite
		// If it's around an hour, close the network connection and reopen it

		// Now we've delivered our data, clear the map and let's go again
		for k := range bucketValueMap {
			delete(bucketValueMap, k)
		}
		timeCounter = 0
	}

}

// Main routine
func main() {
	var debug bool = false

	flag.BoolVar(&debug, "debug", false, "Turn debug on")

	flag.Parse()

	if debug {
		fmt.Println("[debug] Calling main routine\n")
	}

	processInput(debug)

	if debug {
		fmt.Println("[debug] Returned from main routine\n")
	}
}
