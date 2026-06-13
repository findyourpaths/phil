package datetime

import (
	"strings"
	"testing"
)

// TestRangeStartInheritsEndMeridiem: a same-day time range whose start carries
// no AM/PM marker resolves against the end's explicit PM — "6:00 - 7:30 pm"
// means 6 PM (the visible-schedule convention) — but only when that keeps
// start <= end, so the cross-noon "11:00 - 1:30 pm" stays 11 AM, and an
// explicit AM start is never overridden.
func TestRangeStartInheritsEndMeridiem(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		// The page says CST but June is daylight time; phil resolves the
		// abbreviation to America/Chicago and renders the in-effect offset.
		{"Monday, June 15, 2026 6:00 - 7:30 pm CST", "2026-06-15T18:00:00-05:00 - 2026-06-15T19:30:00-05:00"},
		{"June 15, 2026 11:00 - 1:30 pm", "2026-06-15T11:00:00 - 2026-06-15T13:30:00"},
		{"June 15, 2026 6:00 am - 7:30 pm", "2026-06-15T06:00:00 - 2026-06-15T19:30:00"},
		{"June 15, 2026 17:00 - 19:30", "2026-06-15T17:00:00 - 2026-06-15T19:30:00"},
	}
	for _, tc := range cases {
		rngs, err := Parse(tc.in, ParseOptions{})
		if err != nil {
			t.Errorf("Parse(%q) returned error: %v", tc.in, err)
			continue
		}
		if rngs == nil || !strings.Contains(rngs.String(), tc.want) {
			t.Errorf("Parse(%q) = %v, want substring %q", tc.in, rngs, tc.want)
		}
	}
}
