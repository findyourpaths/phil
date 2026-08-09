package datetime

import (
	"strings"
	"testing"
)

// TestParseRejectsBackwardsRange pins the output contract: Parse never
// returns a range whose end date precedes its start date — such a reading is
// a failed parse, and recovery belongs to the caller's general
// parse-failure handling (e.g. ParseBlock over the full text), never to
// repair inside Parse. Source shape: an upstream page whose end-date line
// carries a stale year ("…, 2027 - …, 2026").
func TestParseRejectsBackwardsRange(t *testing.T) {
	minDT := &DateTime{Date: &Date{Year: 2026, Month: 8, Day: 8}}
	opts := ParseOptions{MinDateTime: minDT}

	for _, in := range []string{
		"Saturday, March 6th, 2027 - Sunday, March 8th, 2026",
		"March 6th, 2027 - March 8th, 2026",
		"Saturday, April 17th, 2027 - Sunday, April 19th, 2026",
	} {
		rngs, err := Parse(in, opts)
		if err == nil {
			t.Errorf("Parse(%q) = %s, want backwards-range error", in, rngs)
			continue
		}
		if !strings.Contains(err.Error(), "precedes range start") {
			t.Errorf("Parse(%q) error = %v, want backwards-range semantic error", in, err)
		}
	}

	// Overnight time ranges on valid date order remain parseable: the
	// invariant is date-level only.
	if _, err := Parse("March 6th, 2027 9:00pm - March 7th, 2027 1:00am", opts); err != nil {
		t.Errorf("overnight range parse failed: %v", err)
	}
}
