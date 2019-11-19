package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func processInput(debug bool) {
	var readResult int
	var lowerBucket, upperBucket, cpuLatency int
	var timeCounter int
	var lowerBucketString string
	var lowBucketValString string
	var cpuLatencyVal int

	// Setup at process start

	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	if debug {
		fmt.Println("Hostname:", hostName)
	}

	bucketValueMap := make(map[string]int)
	scanner := bufio.NewScanner(os.Stdin)

	// We want the service to run forever, so loop forever
	for {

		for timeCounter = 0; timeCounter < 15; {
			// Scans a line from Stdin(Console)
			scanner.Scan()
			// Holds the string that scanned
			readText := scanner.Text()
			if len(readText) != 0 {
				if strings.HasPrefix(readText, "@usecs:") {

					//					if debug {
					fmt.Println("Found usecs line")
					//					}

					// BUG HERE. WE FIND A NEW usecs THEN WE PROCESS THE RESULTS WE'VE SEEN. WE'RE ALWAYS OUT OF ORDER AT THIS POINT
					// SEE output FILE
					timeCounter++
				} else {
					fmt.Println("Found a line to process")
					// Need to test for an empty line, or badly formatted line and deal with that
					readResult, _ = fmt.Sscanf(readText, "[%d] %d", &lowerBucket, &cpuLatency)

					if debug {
						fmt.Println("readResult on readText:", readResult, readText)
					}

					if readResult == 2 {
						lowerBucketString = strconv.Itoa(lowerBucket)
					} else {

						switch strings.Count(readText, "K") {
						case 0:
							if debug {
								fmt.Println("No Ks")
							}
							readResult, _ = fmt.Sscanf(readText, "[%d, %d) %d", &lowerBucket, &upperBucket, &cpuLatency)
							lowerBucketString = strconv.Itoa(lowerBucket)
						case 1:
							if debug {
								fmt.Println("1 Ks")
							}
							readResult, _ = fmt.Sscanf(readText, "[%d, %dK) %d", &lowerBucket, &upperBucket, &cpuLatency)
							lowerBucketString = strconv.Itoa(lowerBucket)
						case 2:
							if debug {
								fmt.Println("2 Ks")
							}
							readResult, _ = fmt.Sscanf(readText, "[%dK, %dK) %d", &lowerBucket, &upperBucket, &cpuLatency)
							lowerBucketString = strconv.Itoa(lowerBucket) + "K"
						}
						//readResult, _ = fmt.Sscanf(readText, "[%d, %d) %d", &lowerBucket, &upperBucket, &cpuLatency)
						//readResult, _ = fmt.Sscanf(readText, "[%s, %s) %d", &lowerBucketString, &upperBucketString, &cpuLatency)

					}

					if debug {
						fmt.Println("readResult:", readResult)
						fmt.Println("lowerBucketString:", lowerBucketString)
						fmt.Println("cpuLatency:", cpuLatency)
					}

					bucketValueMap[lowerBucketString] += cpuLatency
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

		for lowBucketValString, cpuLatencyVal = range bucketValueMap {
			fmt.Printf("%s.cpu-lat.%s %d %d\n", hostName, lowBucketValString, cpuLatencyVal, unixEpoch)
		}

		// Send data
		// Check number of times we've sent data to Graphite
		// If it's around an hour, close the network connection and reopen it
		//now := time.Now()
		//secs := now.Unix()

		// Now we've delivered our data, clear the map and let's go again
		for k := range bucketValueMap {
			bucketValueMap[k] = 0
		}
		fmt.Println("timeCounter is:", timeCounter)
		fmt.Println(bucketValueMap)
	}

}

// Main routine
func main() {
	var debug bool = false

	flag.BoolVar(&debug, "debug", false, "Turn debug on")

	flag.Parse()

	if debug {
		fmt.Println("[debug] Calling main routine")
	}

	processInput(debug)

	if debug {
		fmt.Println("[debug] Returned from main routine")
	}
}
