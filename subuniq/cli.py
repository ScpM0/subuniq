from pathlib import Path
import argparse

def subuniq(input_file: str, output_file: str):
    try:
        input_path = Path(input_file).resolve()
        output_path = Path(output_file).resolve()

        with input_path.open('r') as f:
            subdomains = f.read().splitlines()

        unique_subdomains = sorted(set(sub.strip().lower() for sub in subdomains if sub.strip()))

        with output_path.open('w') as f:
            for subdomain in unique_subdomains:
                f.write(subdomain + '\n')

        print(f"\nSubUniq: {len(unique_subdomains)} unique subdomains saved to '{output_path}'.")

    except FileNotFoundError:
        print(f"\nSubUniq: File '{input_file}' not found.")
    except Exception as e:
        print(f"\nSubUniq: Error - {e}")

def main():
    parser = argparse.ArgumentParser(description='SubUniq: Remove duplicate subdomains.')
    parser.add_argument('input_file', help='Input file with subdomains')
    parser.add_argument('output_file', help='Output file for unique subdomains')

    args = parser.parse_args()
    subuniq(args.input_file, args.output_file)
