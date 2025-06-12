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
    _____           _      _   _         _
   / ____|         | |    | | | |       (_)
  | (___   _   _  | |__  | | | | _ __   _   __ _
   \___ \ | | | | | '_ \ | | | || '_ \ | | / _' |
   ____) || |_| | | |_) || |__| || | | || || (_| |
  |_____/   \__,_||_.__/   \____/ |_| |_||_| \__, |
                                             | |
                                             |_|
 SUBUNIQ - Subdomain Deduplication Tool by Mostafa
  Version: ` + version + `
`)
}

// printUsage prints the usage information
func printUsage() {

	fmt.Println("Usage: subuniq -i input.txt -o output.txt [-ignore sub1,sub2] [-format plain|json|csv] [-filter \".gov.eg\"] [-valid]")
}

func isValidSubdomain(sub string) bool {
	
	pattern := `^(?:[a-zA-Z0-9_-]+\.)+[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(sub)
}

func main() {
	// Define command-line flags
	inputPath := flag.String("i", "", "Input file path (required)")
	outputPath := flag.String("o", "", "Output file path (required)")
	ignore := flag.String("ignore", "", "Comma separated substrings to ignore")
	format := flag.String("format", "plain", "Output format: plain, json, csv")
	filter := flag.String("filter", "", "Only include subdomains containing this substring")
	validate := flag.Bool("valid", false, "Only include valid subdomains")

	// Parse command-line arguments into the defined flags
	flag.Parse()

	// If no arguments are provided, print the banner and usage, then exit.
	if len(os.Args) == 1 {
		printBanner()
		printUsage()
		os.Exit(0)
	}

	// Ensure required flags are provided
	if *inputPath == "" || *outputPath == "" {
		printUsage() // Print usage if required arguments are missing
		os.Exit(1)
	}

	// Process ignore list: split by comma, trim spaces, and convert to lowercase
	ignoreList := []string{}
	if *ignore != "" {
		for _, v := range strings.Split(*ignore, ",") {
			ignoreList = append(ignoreList, strings.TrimSpace(strings.ToLower(v)))
		}
	}

	// Open the input file
	inputFile, err := os.Open(*inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close() // Ensure the input file is closed when main exits

	// Initialize data structures for deduplication and concurrency control
	scanner := bufio.NewScanner(inputFile) // Scanner to read file line by line
	seen := make(map[string]bool)          // Map to store unique subdomains (case-insensitive)
	mu := sync.Mutex{}                     // Mutex for concurrent map access (though not strictly needed for single-threaded scan)
	totalLines := 0                        // Counter for total lines processed
	progressStep := 1000                   // How often to print progress updates

	// Iterate through each line of the input file
	for scanner.Scan() {
		totalLines++
		line := strings.TrimSpace(scanner.Text()) // Read and trim whitespace from the line
		lower := strings.ToLower(line)            // Convert line to lowercase for case-insensitive processing
		if lower == "" {
			continue // Skip empty lines
		}

		// Check if the subdomain should be ignored
		ignored := false
		for _, ign := range ignoreList {
			if strings.Contains(lower, ign) {
				ignored = true
				break
			}
		}
		if ignored {
			continue // Skip ignored subdomains
		}

		// Check if the subdomain matches the filter
		if *filter != "" && !strings.Contains(lower, *filter) {
			continue // Skip subdomains that don't match the filter
		}

		// Validate subdomain format if the -valid flag is set
		if *validate && !isValidSubdomain(lower) {
			continue // Skip invalid subdomains
		}

		// Add subdomain to the 'seen' map if it's unique
		mu.Lock()           // Lock to prevent race conditions when modifying map
		if !seen[lower] {
			seen[lower] = true
		}
		mu.Unlock()         // Unlock after map modification

		// Print progress update
		if totalLines%progressStep == 0 {
			fmt.Printf("Processed %d lines...\n", totalLines)
		}
	}

	// Check for any errors during scanning the input file
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Extract unique subdomains from the map into a slice
	uniqueSubs := make([]string, 0, len(seen))
	for sub := range seen {
		uniqueSubs = append(uniqueSubs, sub)
	}
	sort.Strings(uniqueSubs) // Sort the unique subdomains alphabetically

	// Create the output file
	outputFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close() // Ensure the output file is closed

	// Write unique subdomains to the output file based on the specified format
	switch *format {
	case "plain":
		for _, line := range uniqueSubs {
			_, _ = outputFile.WriteString(line + "\n")
		}
	case "json":
		jsonEncoder := json.NewEncoder(outputFile)
		// Encode the slice of unique subdomains directly to JSON
		err := jsonEncoder.Encode(uniqueSubs)
		if err != nil {
			fmt.Printf("Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
	case "csv":
		writer := csv.NewWriter(outputFile)
		for _, line := range uniqueSubs {
			err := writer.Write([]string{line}) // Write each subdomain as a single-column row
			if err != nil {
				fmt.Printf("Error writing CSV: %v\n", err)
				os.Exit(1)
			}
		}
		writer.Flush() // Ensure all buffered CSV data is written to the file
	default:
		fmt.Println("Unsupported format. Use plain, json, or csv.")
		os.Exit(1)
	}

	// Print final statistics
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

