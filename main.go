package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
`)
}

func main() {
	printBanner()

	if len(os.Args) != 3 {
		fmt.Println("Usage: subuniq <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	seen := make(map[string]bool)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lower := strings.ToLower(line)
		if !seen[lower] && lower != "" {
			seen[lower] = true
			outputFile.WriteString(line + "\n")
		}
	}

	fmt.Println("\n Done! Unique subdomains saved.")
}
