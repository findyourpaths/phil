package ical

import (
	"strings"
	"testing"
	"time"

	"github.com/findyourpaths/phil/datetime"
)

func TestGoogleURL(t *testing.T) {
	tests := []struct {
		name         string
		ranges       *datetime.DateTimeRanges
		info         *EventInfo
		wantContains []string
	}{
		{
			name: "timed_event",
			ranges: datetime.NewRangesWithStartEndDateTimes(
				&datetime.DateTime{
					Date: &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time: &datetime.Time{Hour: 9},
				},
				&datetime.DateTime{
					Date: &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time: &datetime.Time{Hour: 12},
				}),
			info: &EventInfo{Summary: "Workshop", Location: "123 Main St"},
			wantContains: []string{
				"calendar.google.com/calendar/r/eventedit",
				"text=Workshop",
				"20230203T090000",
				"20230203T120000",
			},
		},
		{
			name: "all_day_event",
			ranges: datetime.NewRangesWithStartDates(
				&datetime.Date{Year: 2023, Month: 2, Day: 3}),
			info: &EventInfo{Summary: "All Day"},
			wantContains: []string{
				"calendar.google.com/calendar/r/eventedit",
				"text=All+Day",
				"dates=20230203",
			},
		},
		{
			name: "recurring_event",
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
				}),
			info: &EventInfo{Summary: "Series"},
			wantContains: []string{
				"recur=RRULE",
			},
		},
		{
			name: "with_timezone",
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
			info: &EventInfo{Summary: "ET Event"},
			wantContains: []string{
				"ctz=America",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := GoogleURL(tc.ranges, tc.info)
			if err != nil {
				t.Fatalf("GoogleURL error: %v", err)
			}
			for _, want := range tc.wantContains {
				if !strings.Contains(u, want) {
					t.Errorf("URL missing %q\ngot: %s", want, u)
				}
			}
		})
	}
}

func TestOutlookURL(t *testing.T) {
	tests := []struct {
		name         string
		ranges       *datetime.DateTimeRanges
		info         *EventInfo
		wantContains []string
	}{
		{
			name: "timed_event",
			ranges: datetime.NewRangesWithStartEndDateTimes(
				&datetime.DateTime{
					Date: &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time: &datetime.Time{Hour: 9},
				},
				&datetime.DateTime{
					Date: &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time: &datetime.Time{Hour: 12},
				}),
			info: &EventInfo{Summary: "Workshop"},
			wantContains: []string{
				"outlook.live.com/calendar",
				"subject=Workshop",
			},
		},
		{
			name: "all_day_event",
			ranges: datetime.NewRangesWithStartDates(
				&datetime.Date{Year: 2023, Month: 2, Day: 3}),
			info: &EventInfo{Summary: "All Day"},
			wantContains: []string{
				"outlook.live.com/calendar",
				"allday=true",
				"startdt=2023-02-03",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := OutlookURL(tc.ranges, tc.info)
			if err != nil {
				t.Fatalf("OutlookURL error: %v", err)
			}
			for _, want := range tc.wantContains {
				if !strings.Contains(u, want) {
					t.Errorf("URL missing %q\ngot: %s", want, u)
				}
			}
		})
	}
}

func TestYahooURL(t *testing.T) {
	tests := []struct {
		name         string
		ranges       *datetime.DateTimeRanges
		info         *EventInfo
		wantContains []string
	}{
		{
			name: "timed_event",
			ranges: datetime.NewRangesWithStartEndDateTimes(
				&datetime.DateTime{
					Date: &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time: &datetime.Time{Hour: 9},
				},
				&datetime.DateTime{
					Date: &datetime.Date{Year: 2023, Month: 2, Day: 3},
					Time: &datetime.Time{Hour: 12},
				}),
			info: &EventInfo{Summary: "Workshop"},
			wantContains: []string{
				"calendar.yahoo.com",
				"title=Workshop",
				"v=60",
				"dur=0300",
			},
		},
		{
			name: "all_day_event",
			ranges: datetime.NewRangesWithStartDates(
				&datetime.Date{Year: 2023, Month: 2, Day: 3}),
			info: &EventInfo{Summary: "All Day"},
			wantContains: []string{
				"calendar.yahoo.com",
				"title=All+Day",
				"st=20230203",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := YahooURL(tc.ranges, tc.info)
			if err != nil {
				t.Fatalf("YahooURL error: %v", err)
			}
			for _, want := range tc.wantContains {
				if !strings.Contains(u, want) {
					t.Errorf("URL missing %q\ngot: %s", want, u)
				}
			}
		})
	}
}
