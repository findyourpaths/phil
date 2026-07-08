package ical

import (
	"testing"
	"time"

	"github.com/findyourpaths/phil/datetime"
)

func TestFormatRRule(t *testing.T) {
	tests := []struct {
		name string
		rec  *datetime.Recurrence
		want string
	}{
		{
			name: "weekly_interval",
			rec: &datetime.Recurrence{
				Frequency: datetime.FrequencyWeekly,
				Interval:  2,
				Count:     3,
				Weekdays:  []time.Weekday{time.Wednesday},
			},
			want: "FREQ=WEEKLY;INTERVAL=2;COUNT=3;BYDAY=WE",
		},
		{
			name: "monthly_ordinal_weekday",
			rec: &datetime.Recurrence{
				Frequency:  datetime.FrequencyMonthly,
				Weekdays:   []time.Weekday{time.Wednesday},
				NthWeekday: []int{1},
			},
			want: "FREQ=MONTHLY;BYDAY=1WE",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := FormatRRule(tc.rec); got != tc.want {
				t.Fatalf("FormatRRule() = %q, want %q", got, tc.want)
			}
		})
	}
}
