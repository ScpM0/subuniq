package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const version = "V1.2.1"

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
	fmt.Println("Usage: subuniq -i input.txt -o output.txt [-ignore sub1,sub2] [-format plain|json|csv] [-filter ".gov.eg"] [-valid]")
}

func isValidSubdomain(sub string) bool {
	pattern := `^(?:[a-zA-Z0-9_-]+\.)+[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(sub)
}

func main() {
	inputPath := flag.String("i", "", "Input file path (required)")
	outputPath := flag.String("o", "", "Output file path (required)")
	ignore := flag.String("ignore", "", "Comma separated substrings to ignore")
	format := flag.String("format", "plain", "Output format: plain, json, csv")
	filter := flag.String("filter", "", "Only include subdomains containing this substring")
	validate := flag.Bool("valid", false, "Only include valid subdomains")

	flag.Parse()

	if len(os.Args) == 1 {
		printBanner()
		printUsage()
		os.Exit(0)
	}

	if *inputPath == "" || *outputPath == "" {
		printUsage()
		os.Exit(1)
	}

	ignoreList := []string{}
	if *ignore != "" {
		for _, v := range strings.Split(*ignore, ",") {
			ignoreList = append(ignoreList, strings.TrimSpace(strings.ToLower(v)))
		}
	}

	inputFile, err := os.Open(*inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	seen := make(map[string]bool)
	mu := sync.Mutex{}
	totalLines := 0
	progressStep := 1000

	for scanner.Scan() {
		totalLines++
		line := strings.TrimSpace(scanner.Text())
		lower := strings.ToLower(line)
		if lower == "" {
			continue
		}

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

		if *filter != "" && !strings.Contains(lower, *filter) {
			continue
		}

		if *validate && !isValidSubdomain(lower) {
			continue
		}

		mu.Lock()
		if !seen[lower] {
			seen[lower] = true
		}
		mu.Unlock()

		if totalLines%progressStep == 0 {
			fmt.Printf("Processed %d lines...\n", totalLines)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	uniqueSubs := make([]string, 0, len(seen))
	for sub := range seen {
		uniqueSubs = append(uniqueSubs, sub)
	}
	sort.Strings(uniqueSubs)

	outputFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

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

	fmt.Println("\nDone! Unique subdomains saved.")
	fmt.Printf("Total input lines: %d\n", totalLines)
	fmt.Printf("Unique subdomains: %d\n", len(uniqueSubs))
	if len(ignoreList) > 0 {
		fmt.Printf("Ignored substrings: %v\n", ignoreList)
	}
	if *filter != "" {
		fmt.Printf("Filtered by: %s\n", *filter)
	}
	if *validate {
		fmt.Println("Only valid subdomains included")
	}
}
