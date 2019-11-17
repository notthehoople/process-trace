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

	// Read in the Input

	bucketValueMap := make(map[int]int)
	inputArray := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
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

			inputArray = append(inputArray, readText)
			//		} else {
			//			break
		}
	}
	fmt.Println("timeCounter is:", timeCounter)
	fmt.Println(bucketValueMap)
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
