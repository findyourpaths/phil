package datetime

import (
	"strings"
	"testing"
	"time"
)

// Adjacent complete dates must keep their explicit years. A leading weekday
// belongs to its date, not to a separate occurrence near MinDateTime.
func TestParsePreservesYearBeforeWeekday(t *testing.T) {
	cases := []struct{ text, want string }{
		{"Saturday, March 6th, 2027 Sunday, March 8th, 2026", "2027-03-06, 2026-03-08"},
		{"Saturday, April 17th, 2027 Sunday, April 19th, 2026", "2027-04-17, 2026-04-19"},
		{"March 6, 2027 and March 8, 2026", "2027-03-06, 2026-03-08"},
		{"March 6, 2027, March 8, 2026", "2027-03-06, 2026-03-08"},
		{"Saturday, August 8, 2026 Sunday, August 9, 2026", "2026-08-08, 2026-08-09"},
	}
	for _, refDay := range []int{8, 9} {
		for _, tc := range cases {
			t.Run(tc.text+"/"+time.Date(2026, 8, refDay, 0, 0, 0, 0, time.UTC).Format(time.DateOnly), func(t *testing.T) {
				opts := ParseOptions{MinDateTime: &DateTime{Date: &Date{Year: 2026, Month: time.August, Day: refDay}}, DefaultYear: 2026}
				got, err := Parse(tc.text, opts)
				if err != nil || got == nil || got.String() != tc.want {
					t.Fatalf("Parse(%q) = %v, %v; want %s", tc.text, got, err, tc.want)
				}
			})
		}
	}
}

func TestWeekdayCountNormalizationPreservesExplicitRanges(t *testing.T) {
	opts := ParseOptions{MinDateTime: &DateTime{Date: &Date{Year: 2026, Month: time.August, Day: 8}}, DefaultYear: 2026}
	_, err := Parse("Saturday, March 6th, 2027 - Sunday, March 8th, 2026", opts)
	if err == nil || !strings.Contains(err.Error(), "precedes range start") {
		t.Fatalf("backwards explicit range error = %v", err)
	}
	got, err := Parse("September 29, 2026 - October 1, 2026", opts)
	if err != nil || got == nil || got.String() != "2026-09-29 - 2026-10-01" {
		t.Fatalf("valid explicit range = %v, %v", got, err)
	}
}

func TestWeekdayCountNormalization(t *testing.T) {
	for _, tc := range []struct{ text, want string }{
		{"2 Wednesdays", "Wednesdays"},
		{"5 Tuesdays", "Tuesdays"},
		{"2027 Sunday", "2027 Sunday"},
		{"2026 Saturday", "2026 Saturday"},
	} {
		if got := weekdayCountRE.ReplaceAllString(tc.text, "$1"); got != tc.want {
			t.Errorf("weekday count normalization of %q = %q; want %q", tc.text, got, tc.want)
		}
	}
}
