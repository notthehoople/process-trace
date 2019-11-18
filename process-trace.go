package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func processInput(debug bool) {
	var readResult, readResultShort int = 0, 0
	var lowerBucket, upperBucket, cpuLatency int
	var timeCounter int = 0

	// Setup at process start

	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	if debug {
		fmt.Println("Hostname:", hostName)
	}

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
				if strings.HasPrefix(readText, "@usecs:") {

					if debug {
						fmt.Println("Found usecs line")
					}

					timeCounter++
				} else {
					readResult, _ = fmt.Sscanf(readText, "[%d, %d) %d", &lowerBucket, &upperBucket, &cpuLatency)
					fmt.Println("readResult:", readResult)

					// If we have read lines of the correct format, parse them and retain the cpuLatency
					if readResult > 0 {
						if readResult == 1 {
							readResultShort, _ = fmt.Sscanf(readText, "[%d] %d", &lowerBucket, &cpuLatency)
							fmt.Println("readResultShort:", readResultShort)
							if debug {
								fmt.Println("Line Read: ", readText)
								fmt.Println("lowerBucket: ", lowerBucket)
								fmt.Println("cpuLatency: ", cpuLatency)
							}
						} else {
							if debug {
								fmt.Println("Line Read: ", readText)
								fmt.Println("lowerBucket: ", lowerBucket)
								fmt.Println("upperBucket: ", upperBucket)
								fmt.Println("cpuLatency: ", cpuLatency)
							}
						}

						bucketValueMap[lowerBucket] += cpuLatency
					} else {
						if debug {
							fmt.Println("FAILED: ", readText)
						}
					}
				}
			}
		}

		// When we get here, we've read and processed 15 different inputs so it's time to summarise
		if debug {
			fmt.Println("timeCounter is:", timeCounter)
			fmt.Println(bucketValueMap)
		}

		// Send the details to Graphite. Our output is in the format <hostname>.cpu-lat.<lowbucket> <value> <unix-epoch-time> <lowbucket>

		timeNow := time.Now()
		unixEpoch := timeNow.Unix()

		for lowBucketVal, cpuLatencyVal := range bucketValueMap {
			fmt.Printf("%s.cpu-lat.%d %d %d\n", hostName, lowBucketVal, cpuLatencyVal, unixEpoch)
		}

		// Send data
		// Check number of times we've sent data to Graphite
		// If it's around an hour, close the network connection and reopen it
		//now := time.Now()
		//secs := now.Unix()

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
