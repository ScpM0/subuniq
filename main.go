package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

const version = "V1.2.0"

// printBanner prints the tool banner and version
func printBanner() {
	fmt.Println(`
   _____         _      _    _         _        
  / ____|       | |    | |  | |       (_)       
 | (___   _   _ | |__  | |  | | _ __   _   __ _ 
  \___ \ | | | || '_ \ | |  | || '_ \ | | / _' |
  ____) || |_| || |_) || |__| || | | || || (_| |
 |_____/  \__,_||_.__/  \____/ |_| |_||_| \__, |
                                             | |
                                             |_|                                          
 SUBUNIQ - Subdomain Deduplication Tool by Mostafa
  Version: ` + version + `
`)
}

// printUsage prints the usage information
func printUsage() {
	fmt.Println("Usage: subuniq -i input.txt -o output.txt [-ignore sub1,sub2] [-format plain|json|csv]")
}

func main() {
	// Define flags
	inputPath := flag.String("i", "", "Input file path (required)")
	outputPath := flag.String("o", "", "Output file path (required)")
	ignore := flag.String("ignore", "", "Comma separated substrings to ignore")
	format := flag.String("format", "plain", "Output format: plain, json, csv")

	flag.Parse()

	// If no flags provided, just print banner and usage and exit
	if len(os.Args) == 1 {
		printBanner()
		printUsage()
		os.Exit(0)
	}

	// If required flags are missing, print usage and exit
	if *inputPath == "" || *outputPath == "" {
		printUsage()
		os.Exit(1)
	}

	// Prepare ignore list
	ignoreList := []string{}
	if *ignore != "" {
		for _, v := range strings.Split(*ignore, ",") {
			ignoreList = append(ignoreList, strings.TrimSpace(strings.ToLower(v)))
		}
	}

	// Open input file
	inputFile, err := os.Open(*inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	linesChan := make(chan string)

	// Goroutine to read lines from input file
	go func() {
		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			linesChan <- scanner.Text()
		}
		close(linesChan)

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			os.Exit(1)
		}
	}()

	// Map to keep track of unique subdomains
	seen := make(map[string]bool)
	mu := sync.Mutex{}

	totalLines := 0
	for line := range linesChan {
		totalLines++
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)
		if lower == "" {
			continue
		}

		// Check if line contains any ignored substrings
		ignored := false
		for _, ign := range ignoreList {
			if strings.Contains(lower, ign) {
				ignored = true
				break
			}
		}
		if ignored {
			continue
		}

		// Add unique subdomain
		mu.Lock()
		if !seen[lower] {
			seen[lower] = true
		}
		mu.Unlock()
	}

	// Convert map keys to a slice and sort them
	uniqueSubs := make([]string, 0, len(seen))
	for sub := range seen {
		uniqueSubs = append(uniqueSubs, sub)
	}
	sort.Strings(uniqueSubs)

	// Open output file
	outputFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Write output in the requested format
	switch *format {
	case "plain":
		for _, line := range uniqueSubs {
			_, _ = outputFile.WriteString(line + "\n")
		}
	case "json":
		jsonEncoder := json.NewEncoder(outputFile)
		err := jsonEncoder.Encode(uniqueSubs)
		if err != nil {
			fmt.Printf("Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
	case "csv":
		writer := csv.NewWriter(outputFile)
		for _, line := range uniqueSubs {
			err := writer.Write([]string{line})
			if err != nil {
				fmt.Printf("Error writing CSV: %v\n", err)
				os.Exit(1)
			}
		}
		writer.Flush()
	default:
		fmt.Println("Unsupported format. Use plain, json, or csv.")
		os.Exit(1)
	}

	// Print summary
	fmt.Println("\nDone! Unique subdomains saved.")
	fmt.Printf("Total input lines: %d\n", totalLines)
	fmt.Printf("Unique subdomains: %d\n", len(uniqueSubs))
	if len(ignoreList) > 0 {
		fmt.Printf("Ignored substrings: %v\n", ignoreList)
	}
}
