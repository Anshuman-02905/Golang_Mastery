package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

func main() {
	// Ensure the user provides the correct number of arguments
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <hh:mm> <text message>\n", os.Args[0])
		return
	}

	// Initialize the When package for time parsing
	w := when.New(nil)
	w.Add(en.All...)     // Add English language parsing rules
	w.Add(common.All...) // Add common parsing rules

	text := os.Args[2] // Reminder message
	tm := os.Args[1]   // Time input by user

	// Parse the time input
	r, err := w.Parse(tm, time.Now())
	if err != nil {
		log.Fatalf("Unable to parse time: %v\n", err)
		return
	}

	// Check if the parsed result is nil (no valid time found)
	if r == nil {
		fmt.Println("No matches found for the given time")
		return
	}

	// Calculate the duration until the reminder
	duration := r.Time.Sub(time.Now())

	// Check if the reminder time is in the past
	if duration < 0 {
		fmt.Println("The reminder time is in the past; cannot set reminder.")
		return
	}

	fmt.Printf("Please wait %v for the event\n", duration)

	// Wait for the specified duration before triggering the reminder
	time.Sleep(duration)

	// Display the reminder alert
	err = beeep.Alert("Reminder", text, "")
	if err != nil {
		log.Fatal(err)
	}
}
