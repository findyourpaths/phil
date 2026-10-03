package datetime

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// Oracles retain the independently authored endpoint years, rather than the
// current parser output. The full ccisg block matters: its isolated January
// range already parsed correctly before this regression was repaired.
func TestMigration072ExplicitRangeYears(t *testing.T) {
	if os.Getenv("DEBUG") == "true" {
		previous := DoDebug.Load()
		SetDebug(true)
		t.Cleanup(func() { SetDebug(previous) })
	}
	cases := []struct {
		name  string
		input string
		dates [][2]Date
	}{
		{
			name:  "ccisg complete catalog block",
			input: "Oct 1-2, 2026 Nov 5-6, 2026 Jan 14-15, 2027",
			dates: [][2]Date{
				{{Year: 2026, Month: time.October, Day: 1}, {Year: 2026, Month: time.October, Day: 2}},
				{{Year: 2026, Month: time.November, Day: 5}, {Year: 2026, Month: time.November, Day: 6}},
				{{Year: 2027, Month: time.January, Day: 14}, {Year: 2027, Month: time.January, Day: 15}},
			},
		},
		{
			name:  "IFS Self-Led authored detail range",
			input: "Wednesdays, February 3-March 10, 2027",
			dates: [][2]Date{{{Year: 2027, Month: time.February, Day: 3}, {Year: 2027, Month: time.March, Day: 10}}},
		},
		{
			name:  "IFS training 1709 complete schedule",
			input: "Schedule September 1-3, 2027 September 13-15, 2027",
			dates: [][2]Date{
				{{Year: 2027, Month: time.September, Day: 1}, {Year: 2027, Month: time.September, Day: 3}},
				{{Year: 2027, Month: time.September, Day: 13}, {Year: 2027, Month: time.September, Day: 15}},
			},
		},
	}
	for _, tc := range cases {
		for _, defaultYear := range []int{0, 2026} {
			t.Run(fmt.Sprintf("%s/default-year-%d", tc.name, defaultYear), func(t *testing.T) {
				result, err := Parse(tc.input, ParseOptions{
					MinDateTime: NewDateTimeForTime(time.Unix(1784685024, 759687000).UTC()),
					DefaultYear: defaultYear,
				})
				if err != nil {
					t.Fatal(err)
				}
				if result == nil || len(result.Items) != len(tc.dates) {
					t.Fatalf("wanted %d advertised ranges, got %v", len(tc.dates), result)
				}
				for i, expected := range tc.dates {
					for endpoint, actual := range []*DateTime{result.Items[i].Start, result.Items[i].End} {
						if actual == nil || actual.Date == nil || actual.Date.String() != expected[endpoint].String() {
							t.Errorf("range %d endpoint %d: got %v, want %v", i, endpoint, actual, expected[endpoint])
						}
					}
				}
			})
		}
	}
}

// This is the complete saved training 1703 schedule, including the two bridge
// days whose one trailing year governs both named-month dates.
func TestMigration072IFS1703BridgeDays(t *testing.T) {
	const input = "Schedule March 11-13, 2027 April 8-10, 2027 May 20-22, 2027 June 10-12, 2027 Bridge days: March 20 and April 23, 2027"
	const want = "2027-03-11 - 2027-03-13, 2027-04-08 - 2027-04-10, 2027-05-20 - 2027-05-22, 2027-06-10 - 2027-06-12, 2027-03-20, 2027-04-23"
	for _, defaultYear := range []int{0, 2026} {
		t.Run(fmt.Sprintf("default-year-%d", defaultYear), func(t *testing.T) {
			result, err := ParseBlock(input, ParseOptions{
				MinDateTime: NewDateTimeForTime(time.Unix(1784685024, 759687000).UTC()),
				DefaultYear: defaultYear,
			})
			if err != nil {
				t.Fatal(err)
			}
			if result == nil || result.String() != want || len(result.Items) != 6 || result.HasResidue {
				t.Fatalf("want every advertised date-only interval and both 2027 bridge days, got %v", result)
			}
		})
	}
}

// This is the combined Dates/Time contribution from the authoritative Self-Led
// Lawyer detail, not the conflicting list-page date.
func TestMigration072IFSSelfLedDatesAndTime(t *testing.T) {
	result, err := Parse("Dates:  Wednesdays, February 3-March 10, 2027 Time : 12:00–3:00 p.m. ET", ParseOptions{
		MinDateTime: NewDateTimeForTime(time.Unix(1784685024, 759687000).UTC()),
		DefaultYear: 2026,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("want one combined advertised detail range, got %v", result)
	}
	rangeValue := result.Items[0]
	if rangeValue.Start == nil || rangeValue.End == nil || rangeValue.Start.Date == nil || rangeValue.End.Date == nil {
		t.Fatalf("detail range lacks dated endpoints: %v", rangeValue)
	}
	if rangeValue.Start.Date.String() != "2027-02-03" || rangeValue.End.Date.String() != "2027-03-10" {
		t.Errorf("want authoritative detail dates 2027-02-03–2027-03-10, got %v", rangeValue)
	}
	if rangeValue.Start.Time == nil || rangeValue.End.Time == nil || rangeValue.Start.Time.Hour != 12 || rangeValue.End.Time.Hour != 15 {
		t.Errorf("want combined 12:00–15:00 hours, got %v", rangeValue)
	}
	if rangeValue.IANAName() != "America/New_York" {
		t.Errorf("detail hours must retain Eastern zone, got %q", rangeValue.IANAName())
	}
}

// The copied config retains the observed heading year and dated Times rows,
// with the body's advertised PT appended. Neither the summary date range nor
// the capture's reference year may replace these four clock-bearing sessions.
func TestMigration072EquitableDatedClockWindows(t *testing.T) {
	const input = "2026: July 9: 10:30 am-12:30 pm July 16: 10:30 am-12:30 pm July 23: 10:30 am-12:30 pm July 30: 10:30 am-12:30 pm PT"
	for _, defaultYear := range []int{2026, 2027} {
		t.Run(fmt.Sprintf("default-year-%d", defaultYear), func(t *testing.T) {
			result, err := ParseBlock(input, ParseOptions{
				MinDateTime: NewDateTimeForTime(time.Unix(1783518509, 0).UTC()),
				DefaultYear: defaultYear,
			})
			if err != nil {
				t.Fatal(err)
			}
			if result == nil || len(result.Items) != 4 || result.HasResidue {
				t.Fatalf("want four advertised sessions without residue, got %v", result)
			}
			for i, day := range []int{9, 16, 23, 30} {
				r := result.Items[i].Range
				want := fmt.Sprintf("2026-07-%02d", day)
				if r.Start == nil || r.End == nil || r.Start.Date == nil || r.End.Date == nil || r.Start.Date.String() != want || r.End.Date.String() != want {
					t.Fatalf("session %d must retain its advertised date %s, got %v", i, want, r)
				}
				if r.Start.Time == nil || r.End.Time == nil || r.Start.Time.Hour != 10 || r.Start.Time.Minute != 30 || r.End.Time.Hour != 12 || r.End.Time.Minute != 30 {
					t.Errorf("session %d must retain 10:30–12:30, got %v", i, r)
				}
			}
			if got := result.Items[3].Range.IANAName(); got != "America/Los_Angeles" {
				t.Errorf("advertised shared Pacific zone = %q", got)
			}
		})
	}
}

// Professional Pairs advertises four dates, not the 4 and 9 as clocks.
// The conjunction control exercises the already existing grammar unchanged.
func TestMigration072AephoriaMultiMonthDays(t *testing.T) {
	const want = "2026-10-27, 2026-11-04, 2026-11-09, 2026-11-12"
	for _, input := range []string{
		"27 October, 4, 9, 12 November 2026",
		"27 October and 4, 9, 12 November 2026",
	} {
		for _, defaultYear := range []int{2026, 2027} {
			t.Run(fmt.Sprintf("%s/default-year-%d", input, defaultYear), func(t *testing.T) {
				result, err := Parse(input, ParseOptions{DefaultYear: defaultYear})
				if err != nil {
					t.Fatal(err)
				}
				if result == nil || result.String() != want || len(result.Items) != 4 {
					t.Fatalf("want four advertised date-only 2026 items %s, got %v", want, result)
				}
				for _, item := range result.Items {
					if item.Start.Time != nil || item.Start.TimeZone != nil || item.End != nil {
						t.Fatalf("unadvertised clock, zone or end: %v", item)
					}
				}
			})
		}
	}
}

func TestMigration072DatedClockWindowDoesNotInventHours(t *testing.T) {
	for _, input := range []string{
		"July 9:",
		"July 9: 10:30 am",
		"July 9: 10:30 am-12:30 pm Workshop",
		"July 9: 10:30 am-12:30 pm Session 2:",
	} {
		t.Run(input, func(t *testing.T) {
			if got := parseBlockDatedClockRange(blockWords(input), ParseOptions{DefaultYear: 2026}); got != nil {
				t.Fatalf("incomplete clocks or trailing prose must not extend the date window: %v", got)
			}
		})
	}
}
