package ical

import (
	"strings"
	"testing"
	"time"

	"github.com/findyourpaths/phil/datetime"
)


// ---------------------------------------------------------------------------
// Conversion tests
// ---------------------------------------------------------------------------

func TestNewCalendar_Conversion(t *testing.T) {
	tests := []struct {
		name       string
		ranges     *datetime.DateTimeRanges
		info       *EventInfo
		wantEvents int
		wantAllDay bool
		wantRRule  bool
	}{
		{
			name:       "all_day_single",
			ranges:     datetime.NewRangesWithStartDates(&datetime.Date{Year: 2023, Month: 2, Day: 3}),
			info:       &EventInfo{Summary: "Workshop"},
			wantEvents: 1,
			wantAllDay: true,
		},
		{
			name: "timed_single",
			ranges: datetime.NewRangesWithStartDateTimes(
				&datetime.DateTime{
					Date: &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time: &datetime.Time{Hour: 9},
				}),
			info:       &EventInfo{Summary: "Morning Session"},
			wantEvents: 1,
		},
		{
			name: "timed_single_with_timezone",
			ranges: datetime.NewRangesWithStartDateTimes(
				&datetime.DateTime{
					Date:     &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time:     &datetime.Time{Hour: 9},
					TimeZone: &datetime.TimeZone{Abbreviation: "ET"},
				}),
			info:       &EventInfo{Summary: "Morning Session ET"},
			wantEvents: 1,
		},
		{
			name: "date_range",
			ranges: datetime.NewRangesWithStartEndDates(
				&datetime.Date{Year: 2023, Month: 2, Day: 3},
				&datetime.Date{Year: 2023, Month: 2, Day: 4}),
			info:       &EventInfo{Summary: "Retreat"},
			wantEvents: 1,
			wantAllDay: true,
		},
		{
			name: "time_range",
			ranges: datetime.NewRangesWithStartEndDateTimes(
				&datetime.DateTime{
					Date:     &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time:     &datetime.Time{Hour: 9},
					TimeZone: &datetime.TimeZone{Abbreviation: "ET"},
				},
				&datetime.DateTime{
					Date:     &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time:     &datetime.Time{Hour: 12},
					TimeZone: &datetime.TimeZone{Abbreviation: "ET"},
				}),
			info:       &EventInfo{Summary: "Morning Workshop"},
			wantEvents: 1,
		},
		{
			name: "multiple_dates",
			ranges: datetime.NewRangesWithStartDates(
				&datetime.Date{Year: 2023, Month: 2, Day: 1},
				&datetime.Date{Year: 2023, Month: 2, Day: 8},
				&datetime.Date{Year: 2023, Month: 2, Day: 15}),
			info:       &EventInfo{Summary: "Class"},
			wantEvents: 3,
			wantAllDay: true,
		},
		{
			name: "weekly_recurrence",
			ranges: datetime.NewRecurringRanges(
				datetime.NewRange(
					&datetime.DateTime{
						Date:     &datetime.Date{Year: 2023, Month: 2, Day: 1},
						Time:     &datetime.Time{Hour: 9},
						TimeZone: &datetime.TimeZone{Abbreviation: "ET"},
					},
					&datetime.DateTime{
						Date:     &datetime.Date{Year: 2023, Month: 2, Day: 1},
						Time:     &datetime.Time{Hour: 12},
						TimeZone: &datetime.TimeZone{Abbreviation: "ET"},
					}),
				&datetime.Recurrence{
					Frequency: datetime.FrequencyWeekly,
					Count:     5,
					Weekdays:  []time.Weekday{time.Wednesday},
				}),
			info:       &EventInfo{Summary: "5 Week Series"},
			wantEvents: 1,
			wantRRule:  true,
		},
		{
			name: "full_event_info",
			ranges: datetime.NewRangesWithStartDateTimes(
				&datetime.DateTime{
					Date: &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time: &datetime.Time{Hour: 9},
				}),
			info: &EventInfo{
				Summary:     "Workshop",
				Description: "A great workshop",
				Location:    "123 Main St",
				URL:         "https://example.com",
			},
			wantEvents: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cal, err := NewCalendar(tc.ranges, tc.info)
			if err != nil {
				t.Fatalf("NewCalendar error: %v", err)
			}
			events := cal.Events()
			if len(events) != tc.wantEvents {
				t.Fatalf("got %d events, want %d", len(events), tc.wantEvents)
			}

			icsStr := ICS(cal)
			if icsStr == "" {
				t.Fatal("ICS returned empty string")
			}
			if !strings.Contains(icsStr, "BEGIN:VCALENDAR") {
				t.Error("ICS missing BEGIN:VCALENDAR")
			}
			if !strings.Contains(icsStr, "BEGIN:VEVENT") {
				t.Error("ICS missing BEGIN:VEVENT")
			}
			if tc.info.Summary != "" && !strings.Contains(icsStr, tc.info.Summary) {
				t.Errorf("ICS missing SUMMARY %q", tc.info.Summary)
			}
			if tc.wantAllDay && !strings.Contains(icsStr, "VALUE=DATE") {
				t.Error("expected all-day event with VALUE=DATE")
			}
			if tc.wantRRule && !strings.Contains(icsStr, "RRULE:") {
				t.Error("expected RRULE in output")
			}
			if tc.info.Location != "" && !strings.Contains(icsStr, tc.info.Location) {
				t.Errorf("ICS missing LOCATION %q", tc.info.Location)
			}
			if tc.info.Description != "" && !strings.Contains(icsStr, tc.info.Description) {
				t.Errorf("ICS missing DESCRIPTION %q", tc.info.Description)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Validation tests (via NewCalendar errors)
// ---------------------------------------------------------------------------

func TestNewCalendar_Validation(t *testing.T) {
	tests := []struct {
		name    string
		ranges  *datetime.DateTimeRanges
		info    *EventInfo
		wantErr string
	}{
		{
			name: "weekday_mismatch",
			ranges: datetime.NewRecurringRanges(
				datetime.NewRange(
					&datetime.DateTime{Date: &datetime.Date{Year: 2025, Month: 2, Day: 1}},
					nil),
				&datetime.Recurrence{
					Frequency: datetime.FrequencyWeekly,
					Count:     5,
					Weekdays:  []time.Weekday{time.Wednesday},
				}),
			info:    &EventInfo{Summary: "Test"},
			wantErr: "Saturday, not in [Wednesday]",
		},
		{
			name: "count_mismatch",
			ranges: datetime.NewRecurringRanges(
				datetime.NewRange(
					&datetime.DateTime{Date: &datetime.Date{Year: 2023, Month: 2, Day: 1}},
					nil),
				&datetime.Recurrence{
					Frequency: datetime.FrequencyWeekly,
					Count:     5,
					Weekdays:  []time.Weekday{time.Wednesday},
					Until:     &datetime.Date{Year: 2023, Month: 2, Day: 28},
				}),
			info:    &EventInfo{Summary: "Test"},
			wantErr: "count 5 but only 4",
		},
		{
			name: "valid_recurrence",
			ranges: datetime.NewRecurringRanges(
				datetime.NewRange(
					&datetime.DateTime{
						Date: &datetime.Date{Year: 2023, Month: 2, Day: 1},
						Time: &datetime.Time{Hour: 9},
					},
					&datetime.DateTime{
						Date: &datetime.Date{Year: 2023, Month: 2, Day: 1},
						Time: &datetime.Time{Hour: 12},
					}),
				&datetime.Recurrence{
					Frequency: datetime.FrequencyWeekly,
					Count:     5,
					Weekdays:  []time.Weekday{time.Wednesday},
					Until:     &datetime.Date{Year: 2023, Month: 3, Day: 1},
				}),
			info:    &EventInfo{Summary: "Test"},
			wantErr: "",
		},
		{
			name: "invalid_date_feb30",
			ranges: datetime.NewRangesWithStartDates(
				&datetime.Date{Year: 2023, Month: 2, Day: 30}),
			info:    &EventInfo{Summary: "Test"},
			wantErr: "invalid date",
		},
		{
			name: "start_after_end",
			ranges: datetime.NewRangesWithStartEndDates(
				&datetime.Date{Year: 2023, Month: 3, Day: 4},
				&datetime.Date{Year: 2023, Month: 2, Day: 3}),
			info:    &EventInfo{Summary: "Test"},
			wantErr: "start after end",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCalendar(tc.ranges, tc.info)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got: %v", tc.wantErr, err)
			}
		})
	}
}
