package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"time"
)

func main() {
	numCPUs := runtime.NumCPU()
	fmt.Printf("Number of CPUs: %d\n", numCPUs)

	// Set the number of calculations to perform
	numCalculationsPtr := flag.Int("i", 60000, "number of calculations to perform")

	// Parse the command-line arguments
	flag.Parse()

	// Use the value of the -i flag as numCalculations
	numCalculations := *numCalculationsPtr

	// Start the timer
	startTime := time.Now()

	// Create a channel to collect results
	result := make(chan *big.Int, numCalculations)

	// Perform calculations in parallel using goroutines
	for i := 0; i < numCalculations; i++ {
		go fibonacci(int64(i), result)
	}

	// Wait for all goroutines to finish and collect results
	for i := 0; i < numCalculations; i++ {
		<-result
	}

	// Calculate the elapsed time
	elapsedTime := time.Since(startTime)

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		return
	}

	filename := fmt.Sprintf("cpu_stress_results_%s.txt", hostname)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening or creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s - Completed %d calculations in %s\n", currentTime, numCalculations, elapsedTime))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	defer file.Close()

	fmt.Printf("Results appended to %s\n", filename)

}

func fibonacci(n int64, result chan *big.Int) {
	if n <= 1 {
		result <- big.NewInt(n)
		return
	}
	a := big.NewInt(0)
	b := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		a, b = b, a.Add(a, b)
	}
	result <- b
}
