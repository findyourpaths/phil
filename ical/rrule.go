package ical

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/findyourpaths/phil/datetime"
)

// iCal FREQ values by Frequency.
var frequencyICalStrings = map[datetime.Frequency]string{
	datetime.FrequencyDaily:   "DAILY",
	datetime.FrequencyWeekly:  "WEEKLY",
	datetime.FrequencyMonthly: "MONTHLY",
	datetime.FrequencyYearly:  "YEARLY",
}

// iCal BYDAY abbreviations by time.Weekday.
var weekdayICalStrings = map[time.Weekday]string{
	time.Sunday: "SU", time.Monday: "MO", time.Tuesday: "TU",
	time.Wednesday: "WE", time.Thursday: "TH", time.Friday: "FR",
	time.Saturday: "SA",
}

// FormatRRule formats a Recurrence as an iCal RRULE string.
// Returns empty string for nil Recurrence.
func FormatRRule(rec *datetime.Recurrence) string {
	if rec == nil {
		return ""
	}
	s := "FREQ=" + frequencyICalStrings[rec.Frequency]
	if rec.Count > 0 {
		s += ";COUNT=" + strconv.Itoa(rec.Count)
	}
	if len(rec.Weekdays) > 0 {
		var days []string
		for _, wd := range rec.Weekdays {
			days = append(days, weekdayICalStrings[wd])
		}
		s += ";BYDAY=" + strings.Join(days, ",")
	}
	if rec.Until != nil {
		s += ";UNTIL=" + FormatICalDate(rec.Until)
	}
	return s
}

// FormatICalDate formats a Date in YYYYMMDD format for iCal values.
func FormatICalDate(d *datetime.Date) string {
	return fmt.Sprintf("%04d%02d%02d", d.Year, d.Month, d.Day)
}
