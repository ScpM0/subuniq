# SubUniq

**SubUniq** is a simple command-line tool written in Go that removes duplicate subdomains from a text file and saves the unique subdomains to another file.

---

## Features

- Removes duplicate subdomains from your list.
- Outputs sorted (alphabetically) unique subdomains.
- Lightweight and fast thanks to Go.
- Easy to use from the terminal.
- Case-insensitive deduplication  
- Ignores empty lines  
- Clean and simple CLI interface  
- Outputs results preserving original casing  
- Handles large files efficiently (future enhancement)


---

## Requirements

- Linux, macOS, or Windows.
- Go installed (only if you want to build/install from source).
- Or use precompiled binaries if available.

---

## Installation

### 1. Install via Go

Run:

```bash
go install github.com/ScpM0/subuniq@latest
```
```bash
sudo mv ~/go/bin/subuniq /usr/local/bin/
```

# Usage
```bash
subuniq -i input_file.txt -o output_file.txt
```
