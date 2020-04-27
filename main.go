package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var sFile string // For storing source file
	var dFile string // For storing destination file (.csv)

	// Set flags
	flag.StringVar(&sFile, "f", "", "Filename of hosts file.")
	flag.StringVar(&dFile, "o", "results.csv", "Name of output file.")

	// Parse flags
	flag.Parse()

	// Check to see if a valid flag was passed
	if sFile == "" {
		// Display usage and exit
		flag.Usage()
		os.Exit(1)
	}

	// Open the provide source file
	f, err := os.Open(sFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Close file before main returns
	defer f.Close()

	// This is to store all rows
	var allRows [][]string
	// This is to set the header row
	allRows = append(allRows, []string{"Hostname", "Resolved IP"})

	// This creates a scanner that will scan each line of the provided source file
	s := bufio.NewScanner(f)
	for s.Scan() {
		fmt.Println("Attempting to resolve:", s.Text())
		result, err := net.LookupHost(s.Text())
		if err != nil {
			allRows = append(allRows, []string{s.Text(), err.Error()})
			continue
		}
		fmt.Printf("'%v' resolved to -> %v\n", s.Text(), result[0])
		allRows = append(allRows, []string{s.Text(), result[0]})
	}

	// Create or append destination file
	fo, err := os.OpenFile(dFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fo.Close()

	// Write all rows to csv
	w := csv.NewWriter(fo)
	fmt.Println("Attempting to write all results..")
	err = w.WriteAll(allRows)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Done - results saved to:", dFile)
}
