package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func processInput(debug bool) {
	var lowerBucket, upperBucket, cpuLatency, readResult, timeCounter, cpuLatencyVal int
	var lowerBucketString, lowBucketValString string
	var connectionCounter int

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

	timeCounter = 0

	// Setup the connection to Graphite for delivering our data
	connectGraphite, err := net.Dial("tcp", "localhost:2003")
	if err != nil {
		fmt.Printf("Failed to connect")
		return
	}

	// We want the service to run forever, so loop forever
	for {
		// Loop through the data blocks we receive each second
		for timeCounter < 16 {
			// Scans a line from Stdin(Console)
			scanner.Scan()
			// Holds the string that scanned
			readText := scanner.Text()
			if len(readText) != 0 {
				if strings.HasPrefix(readText, "@usecs:") {
					// We've found the @usecs: line, which separates each new block of output. Let's count it as it's a new "second", then loop for data
					timeCounter++
				} else {
					// First let's look the special case where data is something like [1]  2 |@@@@   |
					readResult, _ = fmt.Sscanf(readText, "[%d] %d", &lowerBucket, &cpuLatency)
					if readResult == 2 {
						// Both terms have matched in the fmt.Sscanf, so we've found our special case
						lowerBucketString = strconv.Itoa(lowerBucket)
					} else {
						// Didn't find the special case, so let's look for the other types of data we see
						switch strings.Count(readText, "K") {
						case 0:
							// Our data has no Ks, so will be something like [4, 8)    5 |@@@   |
							readResult, _ = fmt.Sscanf(readText, "[%d, %d) %d", &lowerBucket, &upperBucket, &cpuLatency)
							lowerBucketString = strconv.Itoa(lowerBucket)
						case 1:
							// Our data has 1 K in it, so will be something like [512, 1K)   5 |@@   |
							readResult, _ = fmt.Sscanf(readText, "[%d, %dK) %d", &lowerBucket, &upperBucket, &cpuLatency)
							lowerBucketString = strconv.Itoa(lowerBucket)
						case 2:
							// Our data has 2 Ks in it, so will be something like [1K, 2K)   5 |@@   |
							readResult, _ = fmt.Sscanf(readText, "[%dK, %dK) %d", &lowerBucket, &upperBucket, &cpuLatency)
							lowerBucketString = strconv.Itoa(lowerBucket) + "K"
						}
					}

					bucketValueMap[lowerBucketString] += cpuLatency
				}
			}
		}

		// When we get here, we've read and processed 15 different inputs so it's time to summarise and send the details to Graphite.
		// Our output is in the format <hostname>.cpu-lat.<lowbucket> <value> <unix-epoch-time> <lowbucket>

		timeNow := time.Now()
		unixEpoch := timeNow.Unix()

		for lowBucketValString, cpuLatencyVal = range bucketValueMap {
			if debug {
				fmt.Printf("%s.cpu-lat.%s %d %d\n", hostName, lowBucketValString, cpuLatencyVal, unixEpoch)
			}
			fmt.Fprintf(connectGraphite, "%s.cpu-lat.%s %d %d\n", hostName, lowBucketValString, cpuLatencyVal, unixEpoch)
		}

		// Now we've delivered our data, clear the map and let's go again
		for k := range bucketValueMap {
			bucketValueMap[k] = 0
		}

		// Set timeCounter to 1 since we've already seen the next usecs line
		timeCounter = 1

		// Play a little nicer with network ports. We don't want to open and close a connection to Graphite every 15 seconds
		// We also don't want to open a port and leave it open forever, in case there's a firewall on the network path that objects
		// After 15 minutes (60 times round the 15 second loop) we'll close our connection and open a fresh one
		connectionCounter++
		if connectionCounter > 60 {
			connectGraphite.Close()
			connectGraphite, err = net.Dial("tcp", "localhost:2003")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to connect when reconnecting\n")
				break
			}
			connectionCounter = 0
		}
	}
}

// Main routine
func main() {
	var debug bool = false

	flag.BoolVar(&debug, "debug", false, "Turn debug on")

	flag.Parse()

	processInput(debug)
}
