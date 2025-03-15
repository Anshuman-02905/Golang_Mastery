# Domain Email Security Checker

## Overview
This Go application checks the email security configurations of a domain by verifying:
- **MX Records** (Mail Exchange)
- **SPF Records** (Sender Policy Framework)
- **DMARC Records** (Domain-based Message Authentication, Reporting & Conformance)

## Features
- Scans a given domain for its email security records.
- Identifies whether MX, SPF, and DMARC records are present.
- Displays SPF and DMARC records if available.

## Prerequisites
- Go (1.18 or later)
- Internet connection (for DNS queries)

## Installation
1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd <repository-folder>
   ```
2. Build the application:
   ```sh
   go build -o domain-checker
   ```

## Usage
Run the executable and enter domain names one at a time:
```sh
go run main.go
```
Then input a domain (e.g., `example.com`) and press **Enter** to check its records.

## Example Output
```
Enter Mail ID to check domain records:
domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord
example.com, true, true, "v=spf1 include:_spf.example.com ~all", true, "v=DMARC1; p=none; rua=mailto:dmarc-reports@example.com"
```

## Error Handling
- If DNS queries fail, errors are logged but the program continues running.
- Invalid or unreachable domains will return errors without crashing the application.

## License
This project is open-source under the MIT License.

## Author
Anshuman Mandal
