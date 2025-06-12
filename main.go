package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "sort"
    "strings"
)

func main() {
    inputFile := flag.String("i", "", "Input file path (required)")
    outputFile := flag.String("o", "", "Output file path (required)")
    flag.Parse()

    if *inputFile == "" || *outputFile == "" {
        fmt.Println("Usage: subuniq -i input.txt -o output.txt")
        os.Exit(1)
    }

    subs, err := readLines(*inputFile)
    if err != nil {
        log.Fatalf("Error reading input file: %v", err)
    }

    uniqueSubs := uniqueStrings(subs)

    sort.Strings(uniqueSubs)

    err = writeLines(uniqueSubs, *outputFile)
    if err != nil {
        log.Fatalf("Error writing output file: %v", err)
    }

    fmt.Printf("SubUniq: Saved %d unique subdomains to '%s'.\n", len(uniqueSubs), *outputFile)
}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            lines = append(lines, strings.ToLower(line))
        }
    }
    return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    for _, line := range lines {
        _, err := file.WriteString(line + "\n")
        if err != nil {
            return err
        }
    }
    return nil
}

func uniqueStrings(input []string) []string {
    uniqueMap := make(map[string]bool)
    for _, str := range input {
        uniqueMap[str] = true
    }
    var result []string
    for k := range uniqueMap {
        result = append(result, k)
    }
    return result
}
