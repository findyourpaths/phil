package datetime

import (
	"testing"
	"time"
)

// Regression: a newline between date/time expressions is a list separator, so a
// multi-day schedule ("Sat … | 10:00am – 5:00pm \n Sun … | 10:00am – 5:00pm")
// parses as separate day-long events. Before the fix the newline was collapsed to
// a space, letting a trailing time bind to the next line's weekday — and with a
// MinDateTime that bare weekday resolved to a bogus near-min date. See
// preprocess() / newlineToAndRE.
func TestNewlineSeparatesDayList(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	opts := ParseOptions{
		DateMode:        DateModeNorthAmerican,
		DefaultLocation: loc,
		DefaultYear:     2026,
		MinDateTime:     &DateTime{Date: &Date{Year: 2026, Month: time.June, Day: 8}},
	}
	rngs, err := Parse("Saturday, June 27 | 10:00am – 5:00pm\nSunday, June 28 | 10:00am – 5:00pm", opts)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	want := "2026-06-27T10:00:00-04:00 - 2026-06-27T17:00:00-04:00, 2026-06-28T10:00:00-04:00 - 2026-06-28T17:00:00-04:00"
	if rngs == nil || rngs.String() != want {
		t.Errorf("got  %v\nwant %s", rngs, want)
	}
}
