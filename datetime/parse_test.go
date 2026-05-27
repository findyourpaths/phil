package datetime

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/findyourpaths/phil/glr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k0kubun/pp/v3"
	"github.com/sanity-io/litter"
)

// Run with
// rm datetime_glr.go; time GOWORK=off go generate -v ./...
// time go test -v ./...
// DEBUG=true time go test -v -run '^TestExtractDatetimesRanges/102'

// --- Date fixtures (month-only) ---

var DateForFeb = &Date{Month: 2}
var DateForMar = &Date{Month: 3}

// --- Date fixtures (month-day, sorted) ---

var DateForFeb01 = &Date{Month: 2, Day: 1}
var DateForFeb02 = &Date{Month: 2, Day: 2}
var DateForFeb03 = &Date{Month: 2, Day: 3}
var DateForFeb04 = &Date{Month: 2, Day: 4}
var DateForFeb05 = &Date{Month: 2, Day: 5}
var DateForFeb06 = &Date{Month: 2, Day: 6}
var DateForFeb08 = &Date{Month: 2, Day: 8}
var DateForFeb14 = &Date{Month: 2, Day: 14}
var DateForFeb15 = &Date{Month: 2, Day: 15}
var DateForFeb22 = &Date{Month: 2, Day: 22}
var DateForMar01 = &Date{Month: 3, Day: 1}
var DateForMar02 = &Date{Month: 3, Day: 2}
var DateForMar03 = &Date{Month: 3, Day: 3}
var DateForMar04 = &Date{Month: 3, Day: 4}

// --- Date fixtures (year-only and year-month) ---

var DateFor2023 = &Date{Year: 2023}
var DateFor2024 = &Date{Year: 2024}
var DateFor2023Feb = &Date{Year: 2023, Month: 2}
var DateFor2023Mar = &Date{Year: 2023, Month: 3}

// --- Date fixtures (year-month-day, sorted) ---

var DateFor2023Feb01 = NewRawDateFromYMD(2023, 2, 1)
var DateFor2023Feb02 = NewRawDateFromYMD(2023, 2, 2)
var DateFor2023Feb03 = NewRawDateFromYMD(2023, 2, 3)
var DateFor2023Feb04 = NewRawDateFromYMD(2023, 2, 4)
var DateFor2023Feb05 = NewRawDateFromYMD(2023, 2, 5)
var DateFor2023Feb08 = NewRawDateFromYMD(2023, 2, 8)
var DateFor2023Feb10 = NewRawDateFromYMD(2023, 2, 10)
var DateFor2023Feb15 = NewRawDateFromYMD(2023, 2, 15)
var DateFor2023Feb18 = NewRawDateFromYMD(2023, 2, 18)
var DateFor2023Feb19 = NewRawDateFromYMD(2023, 2, 19)
var DateFor2023Feb24 = NewRawDateFromYMD(2023, 2, 24)
var DateFor2023Feb26 = NewRawDateFromYMD(2023, 2, 26)
var DateFor2023Feb22 = NewRawDateFromYMD(2023, 2, 22)
var DateFor2023Mar01 = NewRawDateFromYMD(2023, 3, 1)
var DateFor2023Mar02 = NewRawDateFromYMD(2023, 3, 2)
var DateFor2023Mar03 = NewRawDateFromYMD(2023, 3, 3)
var DateFor2026Jun19 = NewRawDateFromYMD(2026, 6, 19)

// --- Date range shortcuts (heavily reused) ---

var DateRangesForFeb03 = NewRangesWithStartDates(DateForFeb03)
var DateRangesFromFeb02ToFeb05 = NewRangesWithStartEndDates(DateForFeb02, DateForFeb05)
var DateRangesFromFeb03ToFeb04 = NewRangesWithStartEndDates(DateForFeb03, DateForFeb04)
var DateRangesForMar02 = NewRangesWithStartDates(DateForMar02)
var DateRangesFor2023Feb = NewRangesWithStartDates(DateFor2023Feb)
var DateRangesFor2023Feb03 = NewRangesWithStartDates(DateFor2023Feb03)
var DateRangesFor2023Mar02 = NewRangesWithStartDates(DateFor2023Mar02)
var DateRangesFrom2023Feb03To2023Feb04 = NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)

// --- Time fixtures (sorted by hour) ---

var TimeFor07AM = &Time{Hour: 7}
var TimeFor09AM = &Time{Hour: 9}
var TimeFor10AM = &Time{Hour: 10}
var TimeFor11AM = &Time{Hour: 11}
var TimeFor12PM = &Time{Hour: 12}
var TimeFor12_30PM = &Time{Hour: 12, Minute: 30}
var TimeFor02PM = &Time{Hour: 14}
var TimeFor03PM = &Time{Hour: 15}
var TimeFor04PM = &Time{Hour: 16}
var TimeFor04_30PM = &Time{Hour: 16, Minute: 30}
var TimeFor05PM = &Time{Hour: 17}
var TimeFor06PM = &Time{Hour: 18}
var TimeFor06_30PM = &Time{Hour: 18, Minute: 30}
var TimeFor07PM = &Time{Hour: 19}
var TimeFor07_45PM = &Time{Hour: 19, Minute: 45}
var TimeFor08PM = &Time{Hour: 20}
var TimeFor08_30PM = &Time{Hour: 20, Minute: 30}
var TimeFor09PM = &Time{Hour: 21}

// --- TimeZone fixtures (sorted) ---

var TimeZoneForSAST = &TimeZone{Abbreviation: "SAST"}

var TimeZoneForCET = &TimeZone{Abbreviation: "CET"}
var TimeZoneForCT = &TimeZone{Abbreviation: "CT"}
var TimeZoneForCDT = &TimeZone{Abbreviation: "CDT"}
var TimeZoneForET = &TimeZone{Abbreviation: "ET"}
var TimeZoneForEDT = &TimeZone{Abbreviation: "EDT"}
var TimeZoneForEST = &TimeZone{Abbreviation: "EST"}
var TimeZoneForEastern = &TimeZone{Name: "Eastern"}
var TimeZoneForMST = &TimeZone{Abbreviation: "MST"}
var TimeZoneForMT = &TimeZone{Abbreviation: "MT"}
var TimeZoneForPST = &TimeZone{Abbreviation: "PST"}
var TimeZoneForPT = &TimeZone{Abbreviation: "PT"}
var TimeZoneForPacific = &TimeZone{Name: "Pacific"}
var TimeZoneForUTC = &TimeZone{Abbreviation: "UTC"}

// --- Recurrence fixtures ---

var RecurrenceWeeklyFri = &Recurrence{Frequency: FrequencyWeekly, Weekdays: []time.Weekday{time.Friday}}
var RecurrenceWeeklySun = &Recurrence{Frequency: FrequencyWeekly, Weekdays: []time.Weekday{time.Sunday}}
var RecurrenceWeeklyTue = &Recurrence{Frequency: FrequencyWeekly, Weekdays: []time.Weekday{time.Tuesday}}
var RecurrenceWeeklyWed = &Recurrence{Frequency: FrequencyWeekly, Weekdays: []time.Weekday{time.Wednesday}}

// --- Test helpers ---

// dt creates a DateTime from leaf Date, Time, and TimeZone fixtures.
func dt(d *Date, t *Time, tz *TimeZone) *DateTime { return NewDateTime(d, t, tz) }

// withRec sets the Recurrence on a DateTimeRanges and returns it (for inline use in test table).
func withRec(rs *DateTimeRanges, rec *Recurrence) *DateTimeRanges {
	rs.Recurrence = rec
	return rs
}

type parseTest struct {
	dateMode string
	minDT    *DateTime

	in   string
	want *DateTimeRanges
	skip string // non-empty: run parse, if passes t.Fatalf("REMOVE SKIP"), else t.Skip(reason)
}

var parseTestDefaultLocation = mustLoadLocation("America/New_York")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

func parseTestOptions(t *testing.T, tc parseTest) ParseOptions {
	t.Helper()
	loc := parseTestDefaultLocation
	defaultYear := 0
	if tc.minDT != nil {
		if tc.minDT.TimeZone != nil {
			if name := tc.minDT.TimeZone.IANAName(); name != "" {
				loaded, err := time.LoadLocation(name)
				if err != nil {
					t.Fatalf("loading minDT timezone %q: %v", name, err)
				}
				loc = loaded
			}
		}
	}
	if offset := firstExpectedOffset(tc.want); offset == "+00:00" || offset == "-00:00" {
		loc = time.UTC
	}
	return ParseOptions{
		MinDateTime:     tc.minDT,
		DateMode:        tc.dateMode,
		DefaultLocation: loc,
		DefaultYear:     defaultYear,
	}
}

func firstExpectedOffset(rngs *DateTimeRanges) string {
	if rngs == nil {
		return ""
	}
	for _, item := range rngs.Items {
		if item == nil {
			continue
		}
		for _, dt := range []*DateTime{item.Start, item.End} {
			if dt != nil && dt.TimeZone != nil && dt.TimeZone.Offset != "" {
				return dt.TimeZone.Offset
			}
		}
	}
	return ""
}

func applyParseTestDefaults(t *testing.T, rngs *DateTimeRanges, opts ParseOptions, input string) {
	t.Helper()
	if rngs == nil {
		return
	}
	hasInputTZ := parseTestInputHasTimezone(input)
	for _, item := range rngs.Items {
		if item == nil {
			continue
		}
		normalizeExpectedInheritedTimezone(item.Start, opts, hasInputTZ)
		normalizeExpectedInheritedTimezone(item.End, opts, hasInputTZ)
		applyParseTestDateDefaults(item.Start, opts)
		applyParseTestDateDefaults(item.End, opts)
	}
	if err := stampDateTimeRangesDefaultTZ(rngs, opts.DefaultLocation); err != nil {
		t.Fatalf("applying expected parse defaults: %v", err)
	}
}

func normalizeExpectedInheritedTimezone(dt *DateTime, opts ParseOptions, hasInputTZ bool) {
	if hasInputTZ || dt == nil || dt.TimeZone == nil ||
		opts.MinDateTime == nil || opts.MinDateTime.TimeZone == nil {
		return
	}
	if dt.TimeZone.IANAName() == opts.MinDateTime.TimeZone.IANAName() {
		dt.TimeZone = nil
	}
}

func parseTestInputHasTimezone(input string) bool {
	upper := strings.ToUpper(input)
	for _, marker := range []string{
		" ET", " EDT", " EST", " EASTERN",
		" PT", " PDT", " PST", " PACIFIC",
		" CT", " CDT", " CST", " CENTRAL",
		" MT", " MDT", " MST", " MOUNTAIN",
		" CET", " CEST", " UTC", " GMT", " SAST",
		" US/", " AUSTRALIAN ",
	} {
		if strings.Contains(upper, marker) {
			return true
		}
	}
	return false
}

func applyParseTestDateDefaults(dt *DateTime, opts ParseOptions) {
	if dt == nil || dt.Date == nil ||
		dt.Date.Year != 0 || dt.Date.Month == 0 || dt.Date.Day == 0 {
		return
	}
	if opts.DefaultYear != 0 {
		dt.Date.Year = opts.DefaultYear
		return
	}
	if opts.MinDateTime != nil && opts.MinDateTime.Date != nil {
		dt.Date.Year = opts.MinDateTime.Date.Year
	}
}

func cloneDateTimeRanges(rngs *DateTimeRanges) *DateTimeRanges {
	if rngs == nil {
		return nil
	}
	out := &DateTimeRanges{
		Items:      make([]*DateTimeRange, 0, len(rngs.Items)),
		Recurrence: cloneRecurrence(rngs.Recurrence),
	}
	for _, item := range rngs.Items {
		out.Items = append(out.Items, cloneDateTimeRange(item))
	}
	return out
}

func cloneDateTimeRange(rng *DateTimeRange) *DateTimeRange {
	if rng == nil {
		return nil
	}
	return &DateTimeRange{
		Start: cloneDateTime(rng.Start),
		End:   cloneDateTime(rng.End),
	}
}

func cloneDateTime(dt *DateTime) *DateTime {
	if dt == nil {
		return nil
	}
	return &DateTime{
		Date:     cloneDate(dt.Date),
		Time:     cloneTime(dt.Time),
		TimeZone: cloneTimeZone(dt.TimeZone),
	}
}

func cloneDate(d *Date) *Date {
	if d == nil {
		return nil
	}
	out := *d
	if d.unknown != nil {
		out.unknown = append([]any(nil), d.unknown...)
	}
	return &out
}

func cloneTime(tm *Time) *Time {
	if tm == nil {
		return nil
	}
	out := *tm
	return &out
}

func cloneTimeZone(tz *TimeZone) *TimeZone {
	if tz == nil {
		return nil
	}
	out := *tz
	return &out
}

func cloneRecurrence(rec *Recurrence) *Recurrence {
	if rec == nil {
		return nil
	}
	out := *rec
	out.Weekdays = append([]time.Weekday(nil), rec.Weekdays...)
	out.NthWeekday = append([]int(nil), rec.NthWeekday...)
	out.Until = cloneDate(rec.Until)
	return &out
}

func TestParse(t *testing.T) {
	if os.Getenv("DEBUG") == "true" {
		glr.DoDebug = true
		DoDebug.Store(true)
	}

	tests := []parseTest{

		// Adapts tests from:
		// https://github.com/robintw/daterangeparser/blob/master/daterangeparser/test.py
		// https://github.com/vitalcode/date-time-range-parser/blob/master/src/test/scala/uk/vitalcode/dateparser/Examples.scala
		// https://github.com/waltzofpearls/dateparser
		//

		// ISO 8601 format

		{in: "2023", want: NewRangesWithStartDates(DateFor2023)},
		{in: "2023-02", want: NewRangesWithStartDates(DateFor2023Feb)},
		{in: "2023-02-03", want: DateRangesFor2023Feb03},
		{in: "2023-02-03T", want: DateRangesFor2023Feb03},
		// {in: "2023-02-03T12", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil))},
		{in: "2023-02-03T12:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil))},
		{in: "2023-02-03T12:00:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil))},
		{in: "2023-02-03T12:00:00Z", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil))},
		{in: "2023-02-03T12:00:00+00:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, &TimeZone{Offset: "+00:00"}))},
		{in: "2023-02-03T12:00:00-00:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, &TimeZone{Offset: "-00:00"}))},
		{in: "2023-02-03T12:00:00-05:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, &TimeZone{Offset: "-05:00"}))},

		//
		// Date
		//

		// MD
		{in: "Feb 3", want: DateRangesForFeb03},
		{in: "February 3", want: DateRangesForFeb03},
		{in: "Fri Feb 3", want: DateRangesForFeb03},
		{in: "Fri 3 Feb", want: DateRangesForFeb03},

		// DM
		{in: "3 Feb", want: DateRangesForFeb03},
		{in: "3, Feb", want: DateRangesForFeb03},
		{in: "3rd of Feb", want: DateRangesForFeb03},
		{in: "fri 3 Feb", want: DateRangesForFeb03},

		// MDY
		{in: "Feb 3 2023", want: DateRangesFor2023Feb03},
		{in: "February 3 2023", want: DateRangesFor2023Feb03},
		{in: "February 3, 2023", want: DateRangesFor2023Feb03},
		{in: "February 3rd, 2023", want: DateRangesFor2023Feb03},
		{in: "Fri Feb 3, 2023", want: DateRangesFor2023Feb03},
		{in: "Friday Feb 3rd 2023", want: DateRangesFor2023Feb03},
		{in: "Time:Feb 3 2023", want: DateRangesFor2023Feb03},
		// DMY
		{in: "3 Feb 2023", want: DateRangesFor2023Feb03},
		{in: "3rd Feb 2023", want: DateRangesFor2023Feb03},
		{in: "3 February, 2023", want: DateRangesFor2023Feb03},
		{in: "Friday 3rd Feb 2023", want: DateRangesFor2023Feb03},

		// MY
		{in: "Feb 2023", want: DateRangesFor2023Feb},

		// Both
		{in: "02.03", want: DateRangesForFeb03, dateMode: DateModeNorthAmerican},
		{in: "02.03", want: DateRangesForMar02, dateMode: DateModeRest},
		{in: "02.03", want: DateRangesFor2023Feb03, dateMode: DateModeNorthAmerican, minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},
		{in: "02.03", want: DateRangesFor2023Mar02, dateMode: DateModeRest, minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},
		{in: "02.03.", want: DateRangesForFeb03, dateMode: DateModeNorthAmerican},
		{in: "02.03.", want: DateRangesForMar02, dateMode: DateModeRest},
		{in: "02.03.", want: DateRangesFor2023Feb03, dateMode: DateModeNorthAmerican, minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},
		{in: "02.03.", want: DateRangesFor2023Mar02, dateMode: DateModeRest, minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},
		{in: "2/3/2023", want: DateRangesFor2023Feb03, dateMode: DateModeNorthAmerican},
		{in: "2/3/2023", want: DateRangesFor2023Mar02, dateMode: DateModeRest},

		// Extra tokens
		{in: "Feb 3 Google Calendar ICS", want: DateRangesForFeb03},
		{in: "Updated: Feb 3", want: DateRangesForFeb03},
		{in: "Workshop Update (2/3/23)", want: DateRangesFor2023Feb03, dateMode: DateModeNorthAmerican},
		{in: "Workshop: Feb 3 2023  VIRTUAL", want: DateRangesFor2023Feb03},
		{in: "Release date: February 3, 2023", want: DateRangesFor2023Feb03},
		// Need to replace scanner for these.
		// {in: "http://musicvenue.de/event/id/2023/02/03", want: DateRangesFor2023Feb03},
		// {in: "http://beatricechestnut.com/calendar/skills-mar-2021-5pnfs", want: DateRangesFor2023Feb03},

		//
		// Dates
		//

		// MD
		{in: "Feb 1, 2", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02)},
		{in: "Feb 1, 2, 3", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03)},
		{in: "Feb 1, 2, 3, 4", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04)},
		{in: "Feb 1, 2, 3, 4, 5", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb05)},
		//		{in: "February 1, 2, March 2, 3, and 4, April 3.", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForMar02, DateForMar03, DateForMar04, DateForApr03)},
		{in: "Feb 3 Mar 2", want: NewRangesWithStartDates(DateForFeb03, DateForMar02)},
		{in: "Our next cohort kicks off on March 2nd and we're accepting applications through February 1st.", want: NewRangesWithStartDates(DateForMar02, DateForFeb01)},
		// DM
		{in: "1, 2 Feb", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02)},
		{in: "1, 2, 3 Feb", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03)},
		{in: "1, 2, 3, 4 Feb", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04)},
		{in: "1, 2, 3, 4, 5 Feb", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb05)},
		{in: "1, 2, 3 Feb and 2 Mar", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForMar02)},
		{in: "1-3 Feb and 2 Mar", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForMar02)},
		{in: "1-3 & 5 February", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb05)},
		{in: "1-4 & 6 February", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb06)},

		// MDY
		{in: "Feb 1, 2 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02)},
		{in: "Feb 1, 2, 3 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03)},
		{in: "Feb 1, 2, 3, 4 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04)},
		{in: "Feb 1, 2, 3, 4, 5 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04, DateFor2023Feb05)},
		{in: "Feb 1, 2, 3 and Mar 2 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Mar02)},
		{in: "Feb 3 Mar 2 2023", want: NewRangesWithStartDates(DateFor2023Feb03, DateFor2023Mar02)},
		// DMY
		{in: "1, 2 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02)},
		{in: "1, 2, 3 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03)},
		{in: "1, 2, 3, 4 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04)},
		{in: "1, 2, 3, 4, 5 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04, DateFor2023Feb05)},
		{in: "1, 2, 3 Feb and 2 Mar 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Mar02)},

		// Day list with time range and timezone (single month)
		{in: "Course Schedule: 1, 3, 8, 10 February 2023 07:00 AM \u2013 11:00 AM (SAST)", want: NewRangesFromDatesTimeRange(
			[]*Date{DateFor2023Feb01, DateFor2023Feb03, DateFor2023Feb08, DateFor2023Feb10},
			TimeFor07AM, TimeFor11AM, TimeZoneForSAST)},
		// Broken time colon from HTML strip ("07: 00" instead of "07:00")
		{in: "Course Schedule: 1, 3, 8, 10 February 2023 07: 00 AM \u2013 11: 00 AM (SAST)", want: NewRangesFromDatesTimeRange(
			[]*Date{DateFor2023Feb01, DateFor2023Feb03, DateFor2023Feb08, DateFor2023Feb10},
			TimeFor07AM, TimeFor11AM, TimeZoneForSAST)},
		// Multi-timezone lines — only first timezone (SAST) should be captured
		{in: "Course Schedule: 18, 19, 24, 26 February 2023 17:00 PM \u2013 21:00 PM (SAST) 07: 00 AM \u2013 11:00 PM (PST) 11:00 AM \u2013 15:00 PM (EDT)", want: NewRangesFromDatesTimeRange(
			[]*Date{DateFor2023Feb18, DateFor2023Feb19, DateFor2023Feb24, DateFor2023Feb26},
			TimeFor05PM, TimeFor09PM, TimeZoneForSAST)},
		// Multi-month day list with time range and timezone (DM order)
		{in: "23, 25, 30 March, and 1 April 2023 17:00 PM \u2013 21:00 PM (SAST)", want: NewRangesFromDatesTimeRange(
			[]*Date{NewRawDateFromYMD(2023, 3, 23), NewRawDateFromYMD(2023, 3, 25), NewRawDateFromYMD(2023, 3, 30), NewRawDateFromYMD(2023, 4, 1)},
			TimeFor05PM, TimeFor09PM, TimeZoneForSAST)},
		// Multi-month with full noisy description text (aephoriagroup real-world input)
		{in: "Accredit in Aephoria Colleagues advanced accreditations. Course schedule: 23, 25, 30 March, and 1 April 2023 17:00 PM \u2013 21:00 PM (SAST) 08:00 AM \u2013 12:00 PM (PDT) 11:00 AM \u2013 15:00 PM (EDT)", want: NewRangesFromDatesTimeRange(
			[]*Date{NewRawDateFromYMD(2023, 3, 23), NewRawDateFromYMD(2023, 3, 25), NewRawDateFromYMD(2023, 3, 30), NewRawDateFromYMD(2023, 4, 1)},
			TimeFor05PM, TimeFor09PM, TimeZoneForSAST)},
		// Multi-month without "and" — comma between months: "30 June, 2 July 2023"
		{in: "Course schedule: 23, 25, 30 June, 2 July 2023 17:00 PM \u2013 21:00 PM (SAST)", want: NewRangesFromDatesTimeRange(
			[]*Date{NewRawDateFromYMD(2023, 6, 23), NewRawDateFromYMD(2023, 6, 25), NewRawDateFromYMD(2023, 6, 30), NewRawDateFromYMD(2023, 7, 2)},
			TimeFor05PM, TimeFor09PM, TimeZoneForSAST)},

		//
		// Date Range
		//

		// MD
		{in: "Feb 3rd-4th", want: DateRangesFromFeb03ToFeb04},
		{in: "Feb 3 - Mar 2", want: NewRangesWithStartEndDates(DateForFeb03, DateForMar02)},
		{in: "Feb 3 to Mar 2", want: NewRangesWithStartEndDates(DateForFeb03, DateForMar02)},
		{in: "February 2 - 5 (TH-SU)", want: DateRangesFromFeb02ToFeb05},
		// DM
		{in: "3-4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 Feb - 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "Fri Feb 3 - Sat Feb 4", want: DateRangesFromFeb03ToFeb04},
		{in: "3 February - 2 March", want: NewRangesWithStartEndDates(DateForFeb03, DateForMar02)},
		// Various separators
		{in: "Feb 3-4", want: DateRangesFromFeb03ToFeb04},
		{in: "3--4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 - 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 -- 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 to 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 until 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 through 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 \u2013 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 \u2014 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3-> 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "From 3 - 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "from 3rd till 4th of Feb", want: DateRangesFromFeb03ToFeb04},

		// MDY
		{in: "Feb 3-4 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Feb 3 - 4 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Feb 3 2023 - Feb 4 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri Feb 3, 2023 - Sat Feb 4, 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri 3 Feb - Sat 4 February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri 3rd Feb - Sat 4th February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri Feb 3rd - Sat 4th February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri 3rd Feb - 4th Sat February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri Feb 3rd - 4th Sat February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "February 3 - March 2, 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Mar02)},
		{in: "SAVE THE DATES: Feb 3-4, 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		// DMY
		{in: "3-4 Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "3-4 Feb. 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "3-4 February, 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "3rd-4th Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "3 Feb 2023 - 4 Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "From 3rd to 4th, Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "beginning 3rd to 4th Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},

		// M
		{in: "Feb - Mar", want: NewRangesWithStartEndDates(DateForFeb, DateForMar)},
		// Y
		{in: "2023 - 2024", want: NewRangesWithStartEndDates(DateFor2023, DateFor2024)},
		// MY
		{in: "Feb 2023 - Mar 2023", want: NewRangesWithStartEndDates(DateFor2023Feb, DateFor2023Mar)},
		// YMD
		{in: "2023, Feb 3 - 2023, Feb 4", want: DateRangesFrom2023Feb03To2023Feb04},

		//
		// Date Ranges
		//

		// MD
		{in: "Feb 1-2, 3-4", want: NewRanges(NewRangeWithStartEndDates(DateForFeb01, DateForFeb02), NewRangeWithStartEndDates(DateForFeb03, DateForFeb04))},
		{in: "Feb 1-2, 3-4 2023", want: NewRanges(NewRangeWithStartEndDates(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04))},
		{in: "Feb 1-2; Mar 2-3", want: NewRanges(NewRangeWithStartEndDates(DateForFeb01, DateForFeb02), NewRangeWithStartEndDates(DateForMar02, DateForMar03))},
		{in: "2/1, 2/2, 3/2, 3/3", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForMar02, DateForMar03), dateMode: DateModeNorthAmerican},
		{in: "1/2, 2/2, 2/3, 3/3", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForMar02, DateForMar03), dateMode: DateModeRest},
		// 5 Mondays 3/17 & 3/31, 4/14 & 4/28, 5/12
		{in: "5 Wednesdays 2/1 & 2/8 & 2/15 & 2/22, 3/2", want: NewRangesWithStartDates(DateForFeb01, DateForFeb08, DateForFeb15, DateForFeb22, DateForMar02)},

		// DM
		{in: "1-2 Feb; 2-3 Mar", want: NewRanges(NewRangeWithStartEndDates(DateForFeb01, DateForFeb02), NewRangeWithStartEndDates(DateForMar02, DateForMar03))},

		// MDY
		{in: "Feb 1-2; Mar 2-3 2023", want: NewRanges(NewRangeWithStartEndDates(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEndDates(DateFor2023Mar02, DateFor2023Mar03))},
		// DMY
		{in: "1-2 Feb; 2-3 Mar 2023", want: NewRanges(NewRangeWithStartEndDates(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEndDates(DateFor2023Mar02, DateFor2023Mar03))},
		{in: "Part 1: 1st-2nd February 2023", want: NewRanges(NewRangeWithStartEndDates(DateFor2023Feb01, DateFor2023Feb02))},

		//
		// Date Time
		//

		// Relative
		{in: "Join today at 12pm", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET)), minDT: dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET)},
		// Join today for Day 2 at 10am PST
		{in: "Join today for Day 2 at 12pm", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET)), minDT: dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET)},
		{in: "Tomorrow at 12pm", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb04, TimeFor12PM, TimeZoneForET)), minDT: dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET)},

		// Relative Z
		{in: "Today Friday, 12pm ET", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET)), minDT: dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET)},

		// MDT
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, nil))},
		{in: "Feb 3 12:00 PM", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, nil))},
		{in: "February 3 @ 12:00 PM", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, nil))},
		{in: "February 3  12 p.m.", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, nil))},
		{in: "Date:Fri 03 Feb, Time:12pm", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, nil))},
		{in: "Starting February 3rd at 12pm", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, nil))},

		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET)), minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},

		// WMDT
		//    Friday, 2/14: **Love is Listening and Art: Social + Listening Art Sessions** at 6pm facilitated by Lauren V
		{in: "Friday 2/14: **Love is Listening and Art: Social + Listening Art Sessions**", want: NewRangesWithStartDates(DateForFeb14), dateMode: DateModeNorthAmerican},
		{in: "Friday 2/14: **Love is Listening and Art: Social + Listening Art Sessions** at 6pm facilitated by Lauren V", want: NewRangesWithStartDateTimes(dt(DateForFeb14, TimeFor06PM, nil)), dateMode: DateModeNorthAmerican},

		// WMDTZ
		// SAT | 5/17 @ 11:00am CST
		{in: "FRI 2/3 @ 12:00pm ET", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},

		// MDTZ
		{in: "Feb 3 12pm ET", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "Feb 3 12pm (ET)", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "Feb 3 12pm - ET", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "Feb 3 12pm in ET", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "Feb 3 12pm Eastern", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForEastern))},
		{in: "Feb 3 12pm US/Eastern", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "FEBRUARY 3RD 12 PM ET, ON FRIDAY", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		// {in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForEast)), minDT: DateTimeForEast, wantDiff: true},
		{in: "Starting February 3rd at 12pm (ET) - Virtually.", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "Starts Friday 2/3 at 9:00 am ET", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor09AM, TimeZoneForET)), dateMode: DateModeNorthAmerican},

		// DMT
		{in: "Date:Fri 03 Feb, Time:3.00pm", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor03PM, nil))},
		{in: "Friday 3 Feb 3:00pm (doors) | 11pm (curfew)", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor03PM, nil))},
		// Need to update sorting algorithm for this.
		{in: "Fri, 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor03PM, nil))},
		{in: "Thu, 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Mar02, TimeFor03PM, nil)), dateMode: DateModeRest},
		{in: "Fri, 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor03PM, nil)), dateMode: DateModeRest},
		{in: "Thu, 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Mar02, TimeFor03PM, nil))},

		// MDYT
		{in: "Feb. 3, 2023 12:00pm", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil))},
		{in: "Feb 3, 2023 @ 12:00 PM", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil))},
		{in: "Friday, February 3rd 2023 from 12:00 PM", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil))},

		// MDYTT
		// Not sure if this is a range or multiple.
		{in: "Feb. 3, 2023 12:00pm, 3:00pm", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil), dt(DateFor2023Feb03, TimeFor03PM, nil))},

		// MDYTZ
		{in: "Feb 3 2023 12pm ET", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET))},
		//    Sep 18, 2024 11:30 AM Eastern Time (US and Canada)
		{in: "Feb 3, 2023 12:00 PM Eastern Time (US and Canada)", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET))},

		// DMY
		{in: "3rd Feb 2023 9:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil))},
		{in: "3rd Feb 2023 9:00am", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil))},
		{in: "3rd Feb 2023 3:00pm", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor03PM, nil))},

		// TZMD
		{in: "12:00 pm ET February 3rd", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},

		//
		// Date Time Ranges

		// MD-MDY
		// Date: Thursday, September 7th - September 11th, 2025 EST
		{in: "Date: Friday, February 3rd - February 4th, 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)},
		{in: "Date: Friday, February 3rd - February 4th, 2023 ET", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, nil, TimeZoneForET), dt(DateFor2023Feb04, nil, TimeZoneForET))},

		// MDTT
		// Need to fix parser for this.
		{in: "February 3: 9am - 12pm", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, nil), dt(DateForFeb03, TimeFor12PM, nil))},
		{in: "Feb 3 9am - 12pm", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, nil), dt(DateForFeb03, TimeFor12PM, nil))},
		{in: "Feb 3 @ 9:00 AM - Feb 3 @ 12:00 PM", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, nil), dt(DateForFeb03, TimeFor12PM, nil))},
		{in: "February, 3 9:00 - 15:00", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, nil), dt(DateForFeb03, TimeFor03PM, nil))},
		{in: "Friday, February 3rd from 12 - 3pm", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor12PM, nil), dt(DateForFeb03, TimeFor03PM, nil))},
		{in: "Feb, 3rd from 9 am-3.00 pm", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, nil), dt(DateForFeb03, TimeFor03PM, nil))},
		{in: "February 3 + 4, 9 am - 12 pm each day", want: NewRanges(
			NewRange(dt(DateForFeb03, TimeFor09AM, nil), dt(DateForFeb03, TimeFor12PM, nil)),
			NewRange(dt(DateForFeb04, TimeFor09AM, nil), dt(DateForFeb04, TimeFor12PM, nil)))},
		//    November 9 + 10: In-person at SeekHealing Asheville,10 am - 7 pm each day
		{in: "February 1 + 2: In-person at SeekHealing Asheville,12 pm - 3 pm each day", want: NewRanges(
			NewRange(dt(DateForFeb01, TimeFor12PM, nil), dt(DateForFeb01, TimeFor03PM, nil)),
			NewRange(dt(DateForFeb02, TimeFor12PM, nil), dt(DateForFeb02, TimeFor03PM, nil)))},
		{in: "THIS Friday: February 3 \n 12-3:00pm", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor12PM, nil), dt(DateForFeb03, TimeFor03PM, nil))},
		//    2 Tuesdays Jan 7th & 21st 6:30p-8:30p
		{in: "2 Wednesdays Feb 1st & 8th 12:00p-3:00p", want: NewRanges(
			NewRange(dt(DateForFeb01, TimeFor12PM, nil), dt(DateForFeb01, TimeFor03PM, nil)),
			NewRange(dt(DateForFeb08, TimeFor12PM, nil), dt(DateForFeb08, TimeFor03PM, nil)))},
		//    10 Mondays 6:30pm-8:30pm March 3rd - May 5th
		{in: "5 Wednesdays 9:00am-12:00pm February 1st - March 1st", want: NewRanges(
			NewRange(dt(DateForFeb01, TimeFor09AM, nil), dt(DateForFeb01, TimeFor12PM, nil)),
			NewRange(dt(DateForFeb08, TimeFor09AM, nil), dt(DateForFeb08, TimeFor12PM, nil)),
			NewRange(dt(DateForFeb15, TimeFor09AM, nil), dt(DateForFeb15, TimeFor12PM, nil)),
			NewRange(dt(DateForFeb22, TimeFor09AM, nil), dt(DateForFeb22, TimeFor12PM, nil)),
			NewRange(dt(DateForMar01, TimeFor09AM, nil), dt(DateForMar01, TimeFor12PM, nil))), skip: "parser: recurrence"},

		// MDTTZ
		{in: "Feb 3rd - 9.00 AM- 12pm ET", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, TimeZoneForET), dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "February 3rd, 9-12pm ET", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, TimeZoneForET), dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "Feb 3 2023 9am - 12pm ET", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET))},
		{in: "Feb 3 2023 9am ET to 12pm ET", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET))},
		{in: "Feb 3 @ 9:00 AM ET - Feb 3 @ 12:00 PM ET", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, TimeZoneForET), dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "Feb 3, 2023, 9:00 AM ET - Feb 3, 2023, 12:00 PM ET", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET))},
		{in: "Friday, 2/3 12-3pm ET", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET), dt(DateForFeb03, TimeFor03PM, TimeZoneForET)), dateMode: DateModeNorthAmerican},
		// Tuesday 8/13 from 6:30 - 8:30 pm EST
		{in: "Friday 2/3 from 12:00 - 3 pm ET", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET), dt(DateForFeb03, TimeFor03PM, TimeZoneForET))},
		//    Sunday, October 6, 9am - 4pm Pacific Time
		{in: "Friday, February 3, 9am - 12pm Eastern Time", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, TimeZoneForET), dt(DateForFeb03, TimeFor12PM, TimeZoneForET))},
		{in: "Sat, February 4 · 4:00 PM PST", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb04, TimeFor04PM, TimeZoneForPST)), minDT: dt(DateFor2023Feb10, TimeFor09AM, TimeZoneForCT)},
		{in: "Saturday, February 4 · 6:00 PM to 8:00 PM CDT", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb04, TimeFor06PM, TimeZoneForCDT), dt(DateFor2023Feb04, TimeFor08PM, TimeZoneForCDT)), minDT: dt(DateFor2023Feb10, TimeFor09AM, TimeZoneForCT)},
		{in: "Saturday, February 4, 2023 · 6:00 PM to 8:00 PM CDT", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb04, TimeFor06PM, TimeZoneForCDT), dt(DateFor2023Feb04, TimeFor08PM, TimeZoneForCDT))},
		{in: "Venue Online, via Zoom Starts Fri Jun 19 2026, 3:00pm EDT Ends Fri Jun 19 2026, 4:30pm EDT", want: NewRangesWithStartEndDateTimes(dt(DateFor2026Jun19, TimeFor03PM, TimeZoneForEDT), dt(DateFor2026Jun19, TimeFor04_30PM, TimeZoneForEDT))},
		{in: "February 3, 2023 from 9:00 am to noon ET", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET))},
		// {in: "February 3, 2023 / 9:00 AM", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor09AM, nil)), wantDiff: true},
		{in: "February 3, 2023 / 9:00 AM - 12:00 PM ET", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET))},
		{in: "February 3rd, 12:00-3:00pm Eastern (New York) time", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor12PM, TimeZoneForET), dt(DateForFeb03, TimeFor03PM, TimeZoneForET))},
		{in: "February 3rd & 4th, 9:00 am - noon Eastern time", want: NewRanges(
			NewRange(dt(DateForFeb03, TimeFor09AM, TimeZoneForET), dt(DateForFeb03, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateForFeb04, TimeFor09AM, TimeZoneForET), dt(DateForFeb04, TimeFor12PM, TimeZoneForET)))},
		{in: "February 3rd - 5th, 9:00 am - noon ET each day", want: NewRanges(
			NewRange(dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateFor2023Feb04, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb04, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateFor2023Feb05, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb05, TimeFor12PM, TimeZoneForET))), skip: "grammar: day range + time"},
		{in: "Wednesdays February 1st & 8th 12:00p-3:00p", want: NewRanges(
			NewRange(dt(DateForFeb01, TimeFor12PM, nil), dt(DateForFeb01, TimeFor03PM, nil)),
			NewRange(dt(DateForFeb08, TimeFor12PM, nil), dt(DateForFeb08, TimeFor03PM, nil)))},
		{in: "Wednesdays, February 1st, 8th, and 15th 9:00am - 12:00pm (ET)", want: NewRanges(
			NewRange(dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb01, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateFor2023Feb08, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb08, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateFor2023Feb15, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb15, TimeFor12PM, TimeZoneForET))), skip: "grammar: complex pattern"},
		{in: "Wednesdays - February 1, 8 9:00 AM - 12:00 PM ET", want: NewRanges(
			NewRange(dt(DateForFeb01, TimeFor09AM, TimeZoneForET), dt(DateForFeb01, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateForFeb08, TimeFor09AM, TimeZoneForET), dt(DateForFeb08, TimeFor12PM, TimeZoneForET)))},
		//    Tuesdays - March 18, 25, and April 1, 8, 15, 22 10:00 AM - 12:30 PM PST
		{in: "Wednesdays - February 1, 8, 15, 22, and March 1 9:00 AM - 12:00 PM ET", want: NewRanges(
			NewRange(dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb01, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateFor2023Feb08, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb08, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateFor2023Feb15, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb15, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateFor2023Feb22, TimeFor09AM, TimeZoneForET), dt(DateFor2023Feb22, TimeFor12PM, TimeZoneForET)),
			NewRange(dt(DateFor2023Mar01, TimeFor09AM, TimeZoneForET), dt(DateFor2023Mar01, TimeFor12PM, TimeZoneForET))), skip: "grammar: complex pattern"},

		// DTT
		{in: "Friday 12 to 3 PM Eastern", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor12PM, TimeZoneForET), dt(DateFor2023Feb03, TimeFor03PM, TimeZoneForET)), minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},

		// DMTT
		{in: "3 Feb 9am - 12pm", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, nil), dt(DateForFeb03, TimeFor12PM, nil))},

		// MDYT
		{in: "Feb 3 2023 12pm", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor12PM, nil))},
		// MDYTT
		{in: "Friday, February 3, 2023 9:00 AM 12:00 PM", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil), dt(DateFor2023Feb03, TimeFor12PM, nil))},
		{in: "Friday, February 3rd 2023 from 9:00 AM to 12:00 PM", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil), dt(DateFor2023Feb03, TimeFor12PM, nil))},
		// DMYTT
		{in: "When 3 Feb 2023 9:00 AM - 12:00 PM", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil), dt(DateFor2023Feb03, TimeFor12PM, nil))},
		{in: "Fr. 3. Feb. 2023, 9:00-ca.12:00", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil), dt(DateFor2023Feb03, TimeFor12PM, nil))},

		// TDMY
		{in: "9:00am 3rd Feb - 4th Feb 3:00pm 2023", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil), dt(DateFor2023Feb04, TimeFor03PM, nil))},
		{in: "9:00am on 3rd Feb - 4th Feb at 3:00pm 2023", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil), dt(DateFor2023Feb04, TimeFor03PM, nil))},
		// Not sure how to parse this one.
		// {in: "(2 Feb 2023 - 3 Feb 2023) 09:00 15:00", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, nil), dt(DateFor2023Feb04, TimeFor03PM, nil))},

		// Both
		{in: "02.03.2023", want: DateRangesFor2023Feb03, dateMode: DateModeNorthAmerican},
		{in: "02.03.2023", want: DateRangesFor2023Mar02, dateMode: DateModeRest},
		{in: "02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor03PM, nil)), dateMode: DateModeNorthAmerican},
		{in: "02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Mar02, TimeFor03PM, nil)), dateMode: DateModeRest},
		{in: "Th , 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Mar02, TimeFor03PM, nil)), dateMode: DateModeNorthAmerican},
		{in: "Fr , 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor03PM, nil)), dateMode: DateModeNorthAmerican},
		{in: "Th , 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Mar02, TimeFor03PM, nil)), dateMode: DateModeRest},
		{in: "Fr , 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor03PM, nil)), dateMode: DateModeRest},
		{in: "Th , 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Mar02, TimeFor03PM, nil))},
		{in: "Fr , 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(dt(DateFor2023Feb03, TimeFor03PM, nil))},

		//
		// Failures

		{in: "814-555-1212", want: nil},
		{in: "814-555-1212 x123", want: nil},
		{in: "102 W. Mahoning Street. Punxsutawney, PA 15767", want: nil},
		// Need to fix these
		{in: "We may request cookies to be set on your device.", want: nil},
		{in: "Winter Retreat for 6-12th graders!", want: nil},
		{in: "For 6th-12th grade students @ SpringHill Camp", want: nil},

		//
		// Category A: Simple Weekday + Time (needs minDT)
		//

		// A1: AM applies to both sides of range
		{in: "Sunday 9 to 10:00 AM Eastern", want: NewRangesWithStartEndDateTimes(
			dt(DateFor2023Feb05, TimeFor09AM, TimeZoneForEastern),
			dt(DateFor2023Feb05, TimeFor10AM, TimeZoneForEastern)),
			minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},
		// A2: TimeZone inherited from minDT when Date+Time present but no explicit TZ
		{in: "Sunday 3:00 PM", want: NewRangesWithStartDateTimes(
			dt(DateFor2023Feb05, TimeFor03PM, TimeZoneForET)),
			minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},
		// A3
		{in: "Thursday 9AM", want: NewRangesWithStartDateTimes(
			dt(DateFor2023Feb02, TimeFor09AM, TimeZoneForET)),
			minDT: dt(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)},

		//
		// Category C: Date + Time Range + Timezone
		//

		// C1: "1st -" looks like date range but is date + time range separator
		{in: "February 1st - 10.00am- 3pm\u00a0MST", want: NewRangesWithStartEndDateTimes(
			dt(DateForFeb01, TimeFor10AM, TimeZoneForMST),
			dt(DateForFeb01, TimeFor03PM, TimeZoneForMST))},
		// C2: Recurrence with start date (simplified to start date + time range)
		{in: "Beginning February 3, 2023, Fridays 3:00 - 5:00 pm EASTERN", want: NewRangesWithStartEndDateTimes(
			dt(DateFor2023Feb03, TimeFor03PM, TimeZoneForEastern),
			dt(DateFor2023Feb03, TimeFor05PM, TimeZoneForEastern)),
			skip: "grammar: complex pattern"},
		// C3: Date range with weekday filter (simplified to date range)
		{in: "February 1 - 15 (M-W-F; M-W-F)", want: NewRangesWithStartEndDates(DateForFeb01, DateForFeb15)},
		// C4
		{in: "February 3 - 14 (M-W-F; M-W-F)", want: NewRangesWithStartEndDates(DateForFeb03, DateForFeb14)},

		//
		// Category D: Multi-Timezone
		//

		// D1: Two TZ variants of same event
		{in: "Thursday night 2 February, 7-8:30pm Australian Eastern Standard Time, 11am-12:30pm CET", want: NewRanges(
			NewRange(dt(DateForFeb02, TimeFor07PM, &TimeZone{Abbreviation: "AEST"}), dt(DateForFeb02, TimeFor08_30PM, &TimeZone{Abbreviation: "AEST"})),
			NewRange(dt(DateForFeb02, TimeFor11AM, TimeZoneForCET), dt(DateForFeb02, TimeFor12_30PM, TimeZoneForCET))),
			skip: "grammar: complex pattern"},
		// D2: Slash-separated TZ alternatives
		{in: "Sunday, February 5th\n09:00 - 11:00 CET / 10:00 - 12:00 UTC", want: NewRanges(
			NewRange(dt(DateForFeb05, TimeFor10AM, TimeZoneForCET), dt(DateForFeb05, TimeFor12PM, TimeZoneForCET)),
			NewRange(dt(DateForFeb05, TimeFor09AM, TimeZoneForUTC), dt(DateForFeb05, TimeFor11AM, TimeZoneForUTC))),
			skip: "grammar: complex pattern"},
		// D3: 7 TZ variants (simplified to 4 US TZs)
		{in: "Friday, February 3\nPT (AZ): 9:00 am, MT: 10:00 am, CT: 11:00 pm, ET: 12:00 pm, London: 4:00 pm, Sweden and France: 5:00 pm, Israel: 6:00 pm", want: NewRangesWithStartDateTimes(
			dt(DateForFeb03, TimeFor09AM, TimeZoneForPT),
			dt(DateForFeb03, TimeFor10AM, TimeZoneForMT),
			dt(DateForFeb03, TimeFor11AM, TimeZoneForCT),
			dt(DateForFeb03, TimeFor12PM, TimeZoneForET)),
			skip: "grammar: complex pattern"},
		// D4: Two TZ variants + recurrence in date range (simplified)
		{in: "February 5 - March 5\nSundays 11:00\u201312:30 CT / 18.00-19.00 CET", want: NewRanges(
			NewRange(dt(DateForFeb05, TimeFor11AM, TimeZoneForCT), dt(DateForFeb05, TimeFor12_30PM, TimeZoneForCT)),
			NewRange(dt(DateForFeb05, TimeFor06PM, TimeZoneForCET), dt(DateForFeb05, TimeFor07PM, TimeZoneForCET))),
			skip: "grammar: complex pattern"},

		//
		// Category E: Recurrence Patterns
		//

		// E1: Multiple weekdays + month range
		{in: "Tues/Thurs 6:30p-9:00p February/March", want: NewRangesWithStartEndDateTimes(
			dt(DateForFeb, TimeFor06_30PM, nil),
			dt(DateForMar, TimeFor09PM, nil)),
			skip: "parser: recurrence"},
		// E2: Count + month range
		{in: "5 Wednesdays February thru March", want: NewRangesWithStartEndDates(DateForFeb, DateForMar),
			skip: "parser: recurrence"},
		// E3: Similar to E1
		{in: "Tuesdays and Thursdays February or March 6:30pm - 9:00pm", want: NewRangesWithStartEndDateTimes(
			dt(DateForFeb, TimeFor06_30PM, nil),
			dt(DateForMar, TimeFor09PM, nil)),
			skip: "parser: recurrence"},
		// E4: Date range + every + weekday
		{in: "February 2 - March 2: Via Zoom with SeekHealing Online, every Thursday from 12 - 2 pm EST", want: NewRangesWithStartEndDateTimes(
			dt(DateForFeb02, TimeFor12PM, TimeZoneForEST),
			dt(DateForMar02, TimeFor02PM, TimeZoneForEST)),
			skip: "parser: recurrence"},
		// E4b: Variant with different dates
		{in: "February 1 - March 1: Via Zoom with SeekHealing Online, every Wednesday from 9 - 12 pm ET", want: NewRangesWithStartEndDateTimes(
			dt(DateForFeb01, TimeFor09AM, TimeZoneForET),
			dt(DateForMar01, TimeFor12PM, TimeZoneForET)),
			skip: "parser: recurrence"},
		// E7: "every" + "through" bound
		{in: "every Wednesday from 12 - 3pm ET through March 1st", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor12PM, TimeZoneForET),
			dt(DateForMar01, TimeFor03PM, TimeZoneForET)),
			skip: "parser: recurrence"},
		// E8: "beginning" as start date
		{in: "every Wednesday from 9 - 12 pm beginning February 1st", want: NewRangesWithStartEndDateTimes(
			dt(DateForFeb01, TimeFor09AM, nil),
			dt(DateForFeb01, TimeFor12PM, nil)),
			skip: "parser: recurrence"},
		// E9: Range + single date, same month
		{in: "1-3 February (in person) & Weds 8th Integration eve (online)", want: NewRanges(
			NewRangeWithStartEndDates(DateForFeb01, DateForFeb03),
			NewRangeWithStart(dt(DateForFeb08, nil, nil))),
			skip: "grammar: complex pattern"},
		// E11: "Every" + "beginning"
		{in: "Every Wednesday beginning February 1st from 12:30 - 2 pm", want: withRec(NewRangesWithStartEndDateTimes(
			dt(DateForFeb01, TimeFor12_30PM, nil),
			dt(DateForFeb01, TimeFor02PM, nil)), RecurrenceWeeklyWed)},
		// E12: Same as E4 with comma separator
		{in: "February 2 - March 2, every Thursday from 12 - 2 pm EST", want: NewRangesWithStartEndDateTimes(
			dt(DateForFeb02, TimeFor12PM, TimeZoneForEST),
			dt(DateForMar02, TimeFor02PM, TimeZoneForEST)),
			skip: "parser: recurrence"},
		// E13: "Courses Begin" is noise; two separate date items
		{in: "February 1 - 8\nCourses Begin February 15", want: NewRanges(
			NewRangeWithStartEndDates(DateForFeb01, DateForFeb08),
			NewRangeWithStart(dt(DateForFeb15, nil, nil))),
			skip: "grammar: complex pattern"},
		// E14: "NEW" prefix, "Starts" keyword — noise stripping + plural weekday prefix
		{in: "NEW Sundays at 9 am ET - Starts February 5", want: withRec(NewRangesWithStartDateTimes(
			dt(DateForFeb05, TimeFor09AM, TimeZoneForET)), RecurrenceWeeklySun)},
		// E15: Too complex (multi-TZ + 36 class recurrence)
		{in: "Program Begins Thursday, February 2, 2023 | 36 Classes PT: 5:00 pm, MT: 6:00 pm, CT: 7:00, ET: 8:00 pm 90 Minute Sessions.", want: NewRangesWithStartDateTimes(
			dt(DateFor2023Feb02, TimeFor05PM, TimeZoneForPT)),
			skip: "parser: recurrence"},
		// E16: "Part N:" prefix
		{in: "Part 1: 1st\u20138th February 2023, Part 2: 15th-22nd February 2023", want: NewRanges(
			NewRangeWithStartEndDates(DateFor2023Feb01, DateFor2023Feb08),
			NewRangeWithStartEndDates(DateFor2023Feb15, DateFor2023Feb22)),
			skip: "grammar: complex pattern"},
		// E17: Structured fields with noise
		{in: "CRN 66932, SLFO NC025, Dates: 8 Fri, Feb 3 - Mar 3, Time, 6:00 pm \u2013 8:30pm, Pacific", want: NewRangesWithStartEndDateTimes(
			dt(DateForFeb03, TimeFor06PM, TimeZoneForPacific),
			dt(DateForMar03, TimeFor08_30PM, TimeZoneForPacific)),
			skip: "parser: recurrence"},
		// E18: Structured Starts/Ends/Meets fields
		{in: "Starts: Wednesday, February 1, 2023 Ends: Wednesday, March 1, 2023 Meets: Online, for 5 consecutive Wednesday evenings, from 7:00 PM to 8:30 PM PST", want: NewRangesWithStartEndDateTimes(
			dt(DateFor2023Feb01, TimeFor07PM, TimeZoneForPST),
			dt(DateFor2023Mar01, TimeFor08_30PM, TimeZoneForPST)),
			skip: "parser: recurrence"},
		// E19: Bare weekday + time range
		{in: "Tuesdays 7 to 9 pm ET", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor07PM, TimeZoneForET),
			dt(nil, TimeFor09PM, TimeZoneForET)),
			skip: "parser: recurrence"},
		// E20: Bare recurrence, no time (needs Recurrence type)
		{in: "Weekly on Mondays", want: nil},
		// E21: Bare weekday + time
		{in: "Wednesdays at 6:30pm ET", want: withRec(NewRangesWithStartDateTimes(
			dt(nil, TimeFor06_30PM, TimeZoneForET)), RecurrenceWeeklyWed)},
		// E22: Ordinal recurrence (needs NthWeekday extension)
		{in: "2nd and 4th Tuesdays 7 to 9 pm ET", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor07PM, TimeZoneForET),
			dt(nil, TimeFor09PM, TimeZoneForET)),
			skip: "parser: recurrence"},
		// E23: Suffix noise
		{in: "Tuesdays 10:00am-12:00pm ET, Mindful Pause", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor10AM, TimeZoneForET),
			dt(nil, TimeFor12PM, TimeZoneForET)),
			skip: "parser: recurrence"},
		// E25: "Day N" stripped by dayNumberRE
		{in: "Join today for Day 2 at 10am PST", want: NewRangesWithStartDateTimes(
			dt(DateFor2023Feb03, TimeFor10AM, TimeZoneForPST)),
			minDT: dt(DateFor2023Feb03, TimeFor09AM, TimeZoneForET)},
		// E26: Bare weekday + time range
		{in: "Fridays 3:00 - 5:00 pm EASTERN", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor03PM, TimeZoneForEastern),
			dt(nil, TimeFor05PM, TimeZoneForEastern)),
			skip: "parser: recurrence"},
		// E27: Date range + session count
		{in: "Wednesdays, February 1 - March 1 (4 sessions)", want: NewRangesWithStartEndDates(
			DateForFeb01, DateForMar01)},
		// E30: "Starting" prefix + ordinal recurrence
		{in: "Starting 2nd and 4th Tuesdays 7 to 9 pm ET", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor07PM, TimeZoneForET),
			dt(nil, TimeFor09PM, TimeZoneForET)),
			skip: "parser: recurrence"},

		//
		// Category F: Redundant 12h + 24h Times (Google Calendar ICS)
		//

		// F1
		{in: "Feb 3 12:00 PM 12:00", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor12PM, nil))},
		// F2
		{in: "Feb 3 3:00 PM 15:00", want: NewRangesWithStartDateTimes(dt(DateForFeb03, TimeFor03PM, nil))},
		// F3
		{in: "Feb 3 9:00 AM 09:00 Feb 3 3:00 PM 15:00", want: NewRangesWithStartEndDateTimes(dt(DateForFeb03, TimeFor09AM, nil), dt(DateForFeb03, TimeFor03PM, nil))},
		// F4
		{in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil), dt(DateFor2023Feb03, TimeFor03PM, nil))},
		// F5
		{in: "Feb 3 2023 9:00 AM 3:00 PM 09:00 15:00 Google Calendar ICS", want: NewRangesWithStartEndDateTimes(dt(DateFor2023Feb03, TimeFor09AM, nil), dt(DateFor2023Feb03, TimeFor03PM, nil))},
		// F6: Multi-day with redundant times — also needs grammar: multi-day without separator
		{in: "Fri, Feb 3, 2023 9:00 AM 09:00 Sat, Feb 4, 2023 5:00 PM 17:00", want: NewRangesWithStartEndDateTimes(
			dt(DateFor2023Feb03, TimeFor09AM, nil),
			dt(DateFor2023Feb04, TimeFor05PM, nil)),
			skip: "grammar: multi-day without separator"},

		//
		// Category G: Complex Multi-Day/Multi-Time
		//

		// G1
		{in: "(1 Feb 2023 - 3 Feb 2023) Wednesday 9:00 15:00", want: NewRangesWithStartEndDateTimes(
			dt(DateFor2023Feb01, TimeFor09AM, nil),
			dt(DateFor2023Feb03, TimeFor03PM, nil)),
			skip: "grammar: complex pattern"},
		// G2: Very complex multi-weekday schedule (deferred)
		{in: "(1 Feb 2023 - 4 Feb 2023) Wednesday 11:00 13:00 Thursday 14:00 15:00 Friday 16:05 17:20 Saturday 19:30 20:45", want: nil,
			skip: "grammar: complex pattern"},
		// G3
		{in: "(3 Feb 2023) Friday 19:45 21:00", want: NewRangesWithStartEndDateTimes(
			dt(DateFor2023Feb03, TimeFor07_45PM, nil),
			dt(DateFor2023Feb03, TimeFor09PM, nil))},
		// G4: Multiple time slots same day
		{in: "Select date Wed 1 February 12:00pm Wed 1 February 3:00pm Wed 1 Feb 6:00pm (last few)", want: NewRangesWithStartDateTimes(
			dt(DateForFeb01, TimeFor12PM, nil),
			dt(DateForFeb01, TimeFor03PM, nil),
			dt(DateForFeb01, TimeFor06PM, nil)),
			skip: "grammar: complex pattern"},
		// G5: Multi-day multi-time
		{in: "Select date Thu 2 February 7:45pm 8:30pm Fri 3 February 7:45pm - 20:30 Sat 4 February 7:45pm to 21:00", want: NewRanges(
			NewRange(dt(DateForFeb02, TimeFor07_45PM, nil), dt(DateForFeb02, TimeFor08_30PM, nil)),
			NewRange(dt(DateForFeb03, TimeFor07_45PM, nil), dt(DateForFeb03, TimeFor08_30PM, nil)),
			NewRange(dt(DateForFeb04, TimeFor07_45PM, nil), dt(DateForFeb04, TimeFor09PM, nil))),
			skip: "grammar: complex pattern"},

		//
		// Category H: Venue/Door Format + Midnight Crossing
		//

		// H1: Crosses midnight
		{in: "Sat 5pm to 9am", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor05PM, nil),
			dt(nil, TimeFor09AM, nil))},
		// H2: Multiple weekdays crossing midnight
		{in: "Sun/Thu 5pm\u20139am", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor05PM, nil),
			dt(nil, TimeFor09AM, nil))},
		// H3: Crosses midnight
		{in: "Fri 5pm to 9am", want: NewRangesWithStartEndDateTimes(
			dt(nil, TimeFor05PM, nil),
			dt(nil, TimeFor09AM, nil))},
		// H4: Two labeled ranges, "21+" is noise
		{in: "Doors: 8PM / Show: 9PM / 21+", want: NewRangesWithStartDateTimes(
			dt(nil, TimeFor08PM, nil),
			dt(nil, TimeFor09PM, nil))},
		// H5: "21+" and "Free" are noise
		{in: "12PM / 21+ / Free", want: NewRangesWithStartDateTimes(
			dt(nil, TimeFor12PM, nil))},
		// H6: More suffix noise
		{in: "Doors: 8PM / Show: 9PM / 21+RSVP DOES NOT GUARANTEE ENTRY", want: NewRangesWithStartDateTimes(
			dt(nil, TimeFor08PM, nil),
			dt(nil, TimeFor09PM, nil))},

		//
		// Category I: Year Straddling
		//

		// I1: Month decrease = year crossed — year straddling inference in NewRange()
		{in: "25 Dec - 2 Jan 2016", want: NewRangesWithStartEndDates(
			NewRawDateFromYMD(2015, 12, 25), NewRawDateFromYMD(2016, 1, 2))},
		// I2: Both years explicit
		{in: "18 Nov 2015 to 14th Feb 2016", want: NewRangesWithStartEndDates(
			NewRawDateFromYMD(2015, 11, 18), NewRawDateFromYMD(2016, 2, 14))},
		// I3: Both years explicit, large gap
		{in: "18 Nov 2010 to 14th Feb 2016", want: NewRangesWithStartEndDates(
			NewRawDateFromYMD(2010, 11, 18), NewRawDateFromYMD(2016, 2, 14))},
	}

	failed := 0
	for i, tc := range tests {
		if !t.Run(fmt.Sprintf("%03d__%s", i, tc.in), testParseFn(t, tc)) {
			failed++
		}
	}

	if len(tests) == 0 {
		fmt.Println("No tests were run")
		return
	}

	percent := float64(failed) / float64(len(tests)) * 100
	fmt.Printf("TestParse: %.2f%% of tests failed (%d/%d)\n", percent, failed, len(tests))
}

func TestParseDefaultTimezone(t *testing.T) {
	pacific := mustLoadLocation("America/Los_Angeles")
	budapest := mustLoadLocation("Europe/Budapest")

	if got, err := Parse("June 1, 2026", ParseOptions{DefaultLocation: pacific}); err != nil {
		t.Fatalf("bare date with default timezone: %v", err)
	} else {
		assertStartDate(t, got, 2026, time.June, 1)
		assertStartTZ(t, got, "America/Los_Angeles", "PDT", "-07:00")
	}

	if got, err := Parse("December 1, 2026", ParseOptions{DefaultLocation: pacific}); err != nil {
		t.Fatalf("winter bare date with default timezone: %v", err)
	} else {
		assertStartDate(t, got, 2026, time.December, 1)
		assertStartTZ(t, got, "America/Los_Angeles", "PST", "-08:00")
	}

	if _, err := Parse("June 1, 2026", ParseOptions{}); err == nil {
		t.Fatalf("bare date without default timezone returned nil error")
	}

	if got, err := Parse("01 June 2026 10:30 am PT", ParseOptions{DefaultLocation: budapest}); err != nil {
		t.Fatalf("explicit Pacific timezone with default timezone: %v", err)
	} else if name := got.Items[0].Start.TimeZone.IANAName(); name != "America/Los_Angeles" {
		t.Fatalf("explicit Pacific IANAName = %q, want America/Los_Angeles", name)
	}

	if got, err := Parse("01 June 2026 10:30 am ET", ParseOptions{DefaultLocation: pacific}); err != nil {
		t.Fatalf("explicit Eastern timezone with default timezone: %v", err)
	} else if name := got.Items[0].Start.TimeZone.IANAName(); name != "America/New_York" {
		t.Fatalf("explicit Eastern IANAName = %q, want America/New_York", name)
	}

	if got, err := Parse("2026-06-01T10:30:00-07:00", ParseOptions{DefaultLocation: pacific}); err != nil {
		t.Fatalf("numeric offset with matching default timezone: %v", err)
	} else {
		assertStartTZ(t, got, "America/Los_Angeles", "PDT", "-07:00")
	}

	if _, err := Parse("2026-06-01T10:30:00-07:00", ParseOptions{}); err == nil {
		t.Fatalf("numeric offset without default timezone returned nil error")
	}

	if _, err := Parse("2026-12-01T10:30:00-07:00", ParseOptions{DefaultLocation: pacific}); err == nil {
		t.Fatalf("conflicting numeric offset returned nil error")
	}

	if _, err := Parse("June 1, 2026 10:30 am PT", ParseOptions{}); err != nil {
		t.Fatalf("explicit timezone without default timezone: %v", err)
	}

	if _, err := Parse("June 1, 2026 10:30 am", ParseOptions{}); err == nil {
		t.Fatalf("timed text without timezone or default returned nil error")
	}

	if got, err := Parse("01/02/2026", ParseOptions{DefaultLocation: budapest}); err != nil {
		t.Fatalf("NA date mode parse: %v", err)
	} else {
		assertStartDate(t, got, 2026, time.January, 2)
	}

	if got, err := Parse("01/02/2026", ParseOptions{DateMode: DateModeRest, DefaultLocation: budapest}); err != nil {
		t.Fatalf("REST date mode parse: %v", err)
	} else {
		assertStartDate(t, got, 2026, time.February, 1)
	}

	if _, err := Parse("01/02/2026", ParseOptions{DateMode: "middle-endian", DefaultLocation: budapest}); err == nil {
		t.Fatalf("invalid date mode returned nil error")
	}

	if got, err := Parse("March 1, 2026", ParseOptions{DefaultLocation: pacific}); err != nil {
		t.Fatalf("Pacific cache variation parse: %v", err)
	} else {
		assertStartTZ(t, got, "America/Los_Angeles", "PST", "-08:00")
	}
	if got, err := Parse("March 1, 2026", ParseOptions{DefaultLocation: budapest}); err != nil {
		t.Fatalf("Budapest cache variation parse: %v", err)
	} else {
		assertStartTZ(t, got, "Europe/Budapest", "CET", "+01:00")
	}
}

func TestParseDefaultYear(t *testing.T) {
	pacific := mustLoadLocation("America/Los_Angeles")
	budapest := mustLoadLocation("Europe/Budapest")
	min2026 := dt(NewRawDateFromYMD(2026, 1, 1), TimeFor09AM, nil)

	if got, err := Parse("May 28, Wednesday • 16:15 – 16:45", ParseOptions{
		DefaultLocation: pacific,
		DefaultYear:     2025,
	}); err != nil {
		t.Fatalf("default year timed parse: %v", err)
	} else {
		assertStartDate(t, got, 2025, time.May, 28)
		assertStartTime(t, got, 16, 15)
		assertStartTZ(t, got, "America/Los_Angeles", "PDT", "-07:00")
	}

	if got, err := Parse("May 28, Sunday", ParseOptions{
		DefaultLocation: budapest,
		DefaultYear:     2023,
	}); err != nil {
		t.Fatalf("default year all-day parse: %v", err)
	} else {
		assertStartDate(t, got, 2023, time.May, 28)
		assertStartTZ(t, got, "Europe/Budapest", "CEST", "+02:00")
	}

	if got, err := Parse("2026-06-01", ParseOptions{
		DefaultLocation: pacific,
		DefaultYear:     2025,
	}); err != nil {
		t.Fatalf("explicit year parse: %v", err)
	} else {
		assertStartDate(t, got, 2026, time.June, 1)
	}

	if got, err := Parse("May 28", ParseOptions{
		MinDateTime:     min2026,
		DefaultLocation: pacific,
	}); err != nil {
		t.Fatalf("minDT year fallback parse: %v", err)
	} else {
		assertStartDate(t, got, 2026, time.May, 28)
	}

	if got, err := Parse("May 28", ParseOptions{
		MinDateTime:     min2026,
		DefaultLocation: pacific,
		DefaultYear:     2025,
	}); err != nil {
		t.Fatalf("default year overrides minDT parse: %v", err)
	} else {
		assertStartDate(t, got, 2025, time.May, 28)
	}

	if got, err := Parse("May 29", ParseOptions{
		DefaultLocation: pacific,
		DefaultYear:     2025,
	}); err != nil {
		t.Fatalf("default year cache variation 2025 parse: %v", err)
	} else {
		assertStartDate(t, got, 2025, time.May, 29)
	}
	if got, err := Parse("May 29", ParseOptions{
		DefaultLocation: pacific,
		DefaultYear:     2026,
	}); err != nil {
		t.Fatalf("default year cache variation 2026 parse: %v", err)
	} else {
		assertStartDate(t, got, 2026, time.May, 29)
	}
}

func assertStartDate(t *testing.T, rngs *DateTimeRanges, year int, month time.Month, day int) {
	t.Helper()
	if rngs == nil || len(rngs.Items) == 0 || rngs.Items[0].Start == nil || rngs.Items[0].Start.Date == nil {
		t.Fatalf("missing start date in %#v", rngs)
	}
	got := rngs.Items[0].Start.Date
	if got.Year != year || got.Month != month || got.Day != day {
		t.Fatalf("start date = %04d-%02d-%02d, want %04d-%02d-%02d",
			got.Year, got.Month, got.Day, year, month, day)
	}
}

func assertStartTime(t *testing.T, rngs *DateTimeRanges, hour, minute int) {
	t.Helper()
	if rngs == nil || len(rngs.Items) == 0 || rngs.Items[0].Start == nil || rngs.Items[0].Start.Time == nil {
		t.Fatalf("missing start time in %#v", rngs)
	}
	got := rngs.Items[0].Start.Time
	if got.Hour != hour || got.Minute != minute {
		t.Fatalf("start time = %02d:%02d, want %02d:%02d", got.Hour, got.Minute, hour, minute)
	}
}

func assertStartTZ(t *testing.T, rngs *DateTimeRanges, name, abbr, offset string) {
	t.Helper()
	if rngs == nil || len(rngs.Items) == 0 || rngs.Items[0].Start == nil || rngs.Items[0].Start.TimeZone == nil {
		t.Fatalf("missing start timezone in %#v", rngs)
	}
	got := rngs.Items[0].Start.TimeZone
	if got.Name != name || got.Abbreviation != abbr || got.Offset != offset {
		t.Fatalf("start timezone = %#v, want Name=%q Abbreviation=%q Offset=%q", got, name, abbr, offset)
	}
}

func testParseFn(t *testing.T, tc parseTest) func(*testing.T) {
	return func(t *testing.T) {
		opts := parseTestOptions(t, tc)
		got, err := Parse(tc.in, opts)
		if err != nil {
			if tc.skip != "" {
				t.Skip(tc.skip)
			}
			t.Fatalf("error: %v", err)
		}
		if got == nil && tc.want == nil {
			if tc.skip != "" {
				t.Fatalf("REMOVE SKIP: %q — test passes now", tc.skip)
			}
			return
		}

		want := cloneDateTimeRanges(tc.want)
		applyParseTestDefaults(t, want, opts, tc.in)
		diff := cmp.Diff(got, want, cmpopts.IgnoreUnexported(Date{}))

		if tc.skip != "" {
			if diff == "" {
				t.Fatalf("REMOVE SKIP: %q — test passes now", tc.skip)
			}
			t.Skip(tc.skip)
			return
		}
		if diff != "" {
			pp.Default.SetColoringEnabled(false)
			fmt.Printf("got vs. want:\n%s\n", litter.Sdump([]*DateTimeRanges{got, want}))
			t.Fatalf("unexpected difference:\n%v", diff)
		}
	}
}
