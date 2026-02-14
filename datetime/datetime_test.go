package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

type toStringTest struct {
	in   *DateTimeRanges
	want string
}

func TestToString(t *testing.T) {
	tests := []toStringTest{
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM), want: "2023-02-03T12:00:00Z"},
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ET), want: "2023-02-03T12:00:00-05:00"},
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ADD0), want: "2023-02-03T12:00:00Z"},
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_SUB0), want: "2023-02-03T12:00:00Z"},
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_SUB5), want: "2023-02-03T12:00:00-05:00"},
	}

	failed := 0
	for i, tc := range tests {
		if !t.Run(fmt.Sprintf("%03d__%s", i, tc.in), testToStringFn(t, tc)) {
			failed++
		}
	}

	if len(tests) == 0 {
		fmt.Println("No tests were run")
		return
	}

	percent := float64(failed) / float64(len(tests)) * 100
	fmt.Printf("TestToString: %.2f%% of tests failed (%d/%d)\n", percent, failed, len(tests))
}

func testToStringFn(t *testing.T, tc toStringTest) func(*testing.T) {
	return func(t *testing.T) {
		got := tc.in.String()
		if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
			t.Errorf("unexpected difference:\n%v", diff)
		}
	}
}

func weekdayPtr(wd time.Weekday) *time.Weekday { return &wd }

func TestOccurrences(t *testing.T) {
	tests := []struct {
		name      string
		ranges    *DateTimeRanges
		wantCount int
		wantDates []Date // start dates of each occurrence
	}{
		{
			name:      "flat_occurrences",
			ranges:    NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb08),
			wantCount: 2,
		},
		{
			name: "weekly_count_occurrences",
			ranges: NewRecurringRanges(
				NewRange(
					&DateTime{Date: DateFor2023Feb01, Time: TimeFor09AM},
					&DateTime{Date: DateFor2023Feb01, Time: TimeFor12PM}),
				&Recurrence{Frequency: FrequencyWeekly, Count: 5, Weekday: weekdayPtr(time.Wednesday)}),
			wantCount: 5,
			wantDates: []Date{
				{Year: 2023, Month: 2, Day: 1},
				{Year: 2023, Month: 2, Day: 8},
				{Year: 2023, Month: 2, Day: 15},
				{Year: 2023, Month: 2, Day: 22},
				{Year: 2023, Month: 3, Day: 1},
			},
		},
		{
			name: "weekly_until_occurrences",
			ranges: NewRecurringRanges(
				NewRange(
					&DateTime{Date: DateFor2023Feb01, Time: TimeFor09AM},
					&DateTime{Date: DateFor2023Feb01, Time: TimeFor12PM}),
				&Recurrence{
					Frequency: FrequencyWeekly,
					Weekday:   weekdayPtr(time.Wednesday),
					Until:     &Date{Year: 2023, Month: 3, Day: 1},
				}),
			wantCount: 5,
			wantDates: []Date{
				{Year: 2023, Month: 2, Day: 1},
				{Year: 2023, Month: 2, Day: 8},
				{Year: 2023, Month: 2, Day: 15},
				{Year: 2023, Month: 2, Day: 22},
				{Year: 2023, Month: 3, Day: 1},
			},
		},
		{
			name: "daily_count_occurrences",
			ranges: NewRecurringRanges(
				NewRange(
					&DateTime{Date: DateFor2023Feb01},
					nil),
				&Recurrence{Frequency: FrequencyDaily, Count: 3}),
			wantCount: 3,
			wantDates: []Date{
				{Year: 2023, Month: 2, Day: 1},
				{Year: 2023, Month: 2, Day: 2},
				{Year: 2023, Month: 2, Day: 3},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ranges.Occurrences()
			if len(got) != tc.wantCount {
				t.Fatalf("got %d occurrences, want %d", len(got), tc.wantCount)
			}
			if tc.wantDates != nil {
				for i, wantDate := range tc.wantDates {
					gotDate := *got[i].Start.Date
					if diff := cmp.Diff(gotDate, wantDate, cmpopts.IgnoreUnexported(Date{})); diff != "" {
						t.Errorf("occurrence[%d] date mismatch:\n%s", i, diff)
					}
				}
			}
		})
	}
}
