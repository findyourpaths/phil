package datetime

import (
	"strings"
	"testing"
)

// Regression for the GLR recover handler (glr/glr_parse.go). The Enneagram
// Institute detail-page schedule format ("<Weekday>, <Month> <Day> | <h:mm>am –
// <h:mm>pm") drives the parser down a candidate path whose semantic action
// panics with a non-string value. The recover handler used to assert the panic
// value as a string, which itself panicked and aborted the whole parse, so
// these valid dates failed. They must now parse cleanly.
func TestPipeSeparatedScheduleParses(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Saturday, September 12 | 10:00am – 5:00pm", "09-12T10:00:00Z - 0000-09-12T17:00:00Z"},
		{"Friday, November 6 | 6:00pm – 9:00pm", "11-06T18:00:00Z - 0000-11-06T21:00:00Z"},
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
