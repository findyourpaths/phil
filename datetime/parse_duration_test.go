package datetime

import (
	"strings"
	"testing"
)

func TestEnneagramBusinessOneHourDuration(t *testing.T) {
	const input = "Thursday, November 12th, 12 pm Pacific Time Zone for one hour."
	ranges, err := Parse(input, ParseOptions{DefaultYear: 2026})
	if err != nil {
		t.Fatalf("Parse(%q): %v", input, err)
	}
	const want = "2026-11-12T12:00:00-08:00 - 2026-11-12T13:00:00-08:00"
	if ranges == nil || !strings.Contains(ranges.String(), want) {
		t.Fatalf("Parse(%q) = %v, want %s", input, ranges, want)
	}
}

func TestHourDurationPreservesCivilAndZonedDates(t *testing.T) {
	for input, want := range map[string]string{
		"November 12, 2026 12 pm for 2 hours":                    "2026-11-12T12:00:00 - 2026-11-12T14:00:00",
		"December 31, 2026 11 pm for 2 hours":                    "2026-12-31T23:00:00 - 2027-01-01T01:00:00",
		"March 8, 2026 1:30 am America/Los_Angeles for one hour": "2026-03-08T01:30:00-08:00 - 2026-03-08T03:30:00-07:00",
	} {
		ranges, err := Parse(input, ParseOptions{})
		if err != nil || ranges == nil || !strings.Contains(ranges.String(), want) {
			t.Errorf("Parse(%q) = %v, %v; want %s", input, ranges, err, want)
		}
	}
}
