package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter Mail ID to check domain records:")
	fmt.Println("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Unable to scan the input: %v\n", err)
	}
}

func checkDomain(domain string) {
	var (
		hasMX, hasSPF, hasDMARC bool
		spfRecord, dmarcRecord  string
	)

	// Check MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up MX records: %v\n", err)
	} else if len(mxRecords) > 0 {
		hasMX = true
	}

	// Check TXT records for SPF
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records: %v\n", err)
	} else {
		for _, record := range txtRecords {
			if strings.HasPrefix(record, "v=spf1") {
				hasSPF = true
				spfRecord = record
				break
			}
		}
	}

	// Check TXT records for DMARC
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up DMARC records: %v\n", err)
	} else {
		for _, record := range dmarcRecords {
			if strings.HasPrefix(record, "v=DMARC1") {
				hasDMARC = true
				dmarcRecord = record
				break
			}
		}
	}

	fmt.Printf(
		"domain=%v, hasMX=%v, hasSPF=%v, spfRecord=%q, hasDMARC=%v, dmarcRecord=%q\n",
		domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord,
	)
}
