package datetime

import (
	"testing"
	"time"
)

// The saved Roxannemanning page independently advertises the date bounds,
// two-hour Pacific clock interval and Thursday cadence. Until is the inclusive
// series stop, not the end of one nine-month occurrence.
func TestMigration072RoxannemanningBoundedThursdays(t *testing.T) {
	result, err := Parse("Date February 5 – October 15, 2026 Times 8 - 10am PT Thursdays", ParseOptions{
		MinDateTime: NewDateTimeForTime(time.Date(2026, time.July, 8, 0, 0, 0, 0, time.UTC)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("want one first-session interval with recurrence, got %v", result)
	}
	first := result.Items[0]
	if first.Start == nil || first.End == nil || first.Start.Date == nil || first.End.Date == nil {
		t.Fatalf("first session lacks its date endpoints: %v", first)
	}
	if first.Start.Date.String() != "2026-02-05" || first.End.Date.String() != "2026-02-05" {
		t.Errorf("first session must be February 5 only: %v", first)
	}
	if first.Start.Time == nil || first.End.Time == nil || first.Start.Time.Hour != 8 || first.End.Time.Hour != 10 {
		t.Errorf("first session must run 08:00–10:00: %v", first)
	}
	if first.IANAName() != "America/Los_Angeles" {
		t.Errorf("Pacific hours must retain their zone, got %q", first.IANAName())
	}
	recurrence := result.Recurrence
	if recurrence == nil {
		t.Fatal("Thursday series has no recurrence")
	}
	if recurrence.Frequency != FrequencyWeekly || len(recurrence.Weekdays) != 1 || recurrence.Weekdays[0] != time.Thursday {
		t.Errorf("want weekly Thursdays, got %v", recurrence)
	}
	if recurrence.Until == nil || recurrence.Until.String() != "2026-10-15" {
		t.Errorf("want inclusive series stop 2026-10-15, got %v", recurrence)
	}
}
