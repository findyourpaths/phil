package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/findyourpaths/phil/datetime"
	"github.com/findyourpaths/phil/ical"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: phil <natural-language-datetime>\n")
		fmt.Fprintf(os.Stderr, "example: phil 'Saturday, August 2, 2025 from 11am to 2pm'\n")
		os.Exit(1)
	}

	input := strings.Join(os.Args[1:], " ")

	now := time.Now()
	minDT := &datetime.DateTime{
		Date: &datetime.Date{Year: now.Year(), Month: now.Month(), Day: now.Day()},
		Time: &datetime.Time{},
	}

	parsed, err := datetime.Parse(minDT, datetime.DateModeNorthAmerican, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	info := &ical.EventInfo{Summary: input}
	cal, err := ical.NewCalendar(parsed, info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ical error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(ical.ICS(cal))
}
