package datetime

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/findyourpaths/phil/glr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k0kubun/pp/v3"
	"github.com/sanity-io/litter"
)

var acceptBrokenTests = true

// var acceptBrokenTests = false

// Run with
// rm datetime_glr.go; time GOWORK=off go generate -v ./...
// time go test -v ./...
// DEBUG=true time go test -v -run '^TestExtractDatetimesRanges/102'

var DateForFeb = &Date{Month: 2}
var DateForMar = &Date{Month: 3}

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
var DateForApr03 = &Date{Month: time.April, Day: 3}

// var DateForFriday = &Date{Weekday: 6}

var DateRangesForFeb03 = NewRangesWithStartDates(DateForFeb03)
var DateRangesFromFeb02ToFeb05 = NewRangesWithStartEndDates(DateForFeb02, DateForFeb05)
var DateRangesFromFeb03ToFeb04 = NewRangesWithStartEndDates(DateForFeb03, DateForFeb04)
var DateRangesForMar02 = NewRangesWithStartDates(DateForMar02)

var DateFor2023 = &Date{Year: 2023}
var DateFor2024 = &Date{Year: 2024}

var DateFor2023Feb = &Date{Year: 2023, Month: 2}
var DateFor2023Mar = &Date{Year: 2023, Month: 3}

var DateRangesFor2023Feb = NewRangesWithStartDates(DateFor2023Feb)

var DateRangesFor2023Feb03 = NewRangesWithStartDates(DateFor2023Feb03)
var DateRangesFor2023Mar02 = NewRangesWithStartDates(DateFor2023Mar02)

var DateFor2023Jan1 = NewRawDateFromYMD(2023, 1, 1)
var DateFor2023Feb01 = NewRawDateFromYMD(2023, 2, 1)
var DateFor2023Feb02 = NewRawDateFromYMD(2023, 2, 2)
var DateFor2023Feb03 = NewRawDateFromYMD(2023, 2, 3)
var DateFor2023Feb04 = NewRawDateFromYMD(2023, 2, 4)
var DateFor2023Feb05 = NewRawDateFromYMD(2023, 2, 5)
var DateFor2023Feb08 = NewRawDateFromYMD(2023, 2, 8)
var DateFor2023Feb15 = NewRawDateFromYMD(2023, 2, 15)
var DateFor2023Feb22 = NewRawDateFromYMD(2023, 2, 22)
var DateFor2023Mar01 = NewRawDateFromYMD(2023, 3, 1)
var DateFor2023Mar02 = NewRawDateFromYMD(2023, 3, 2)
var DateFor2023Mar03 = NewRawDateFromYMD(2023, 3, 3)

var DateFor2024Jan1 = NewRawDateFromYMD(2024, 1, 1)

var DateRangesFrom2023Feb03To2023Feb04 = NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)

var DateTimeForFeb01_09AM = NewDateTime(DateForFeb01, TimeFor09AM, nil)
var DateTimeForFeb01_12PM = NewDateTime(DateForFeb01, TimeFor12PM, nil)
var DateTimeForFeb01_03PM = NewDateTime(DateForFeb01, TimeFor03PM, nil)

var DateTimeForFeb02_09AM = NewDateTime(DateForFeb02, TimeFor09AM, nil)
var DateTimeForFeb02_12PM = NewDateTime(DateForFeb02, TimeFor12PM, nil)
var DateTimeForFeb02_03PM = NewDateTime(DateForFeb02, TimeFor03PM, nil)

var DateTimeForFeb03_09AM = NewDateTime(DateForFeb03, TimeFor09AM, nil)
var DateTimeForFeb03_12PM = NewDateTime(DateForFeb03, TimeFor12PM, nil)
var DateTimeForFeb03_03PM = NewDateTime(DateForFeb03, TimeFor03PM, nil)

var DateTimeForFeb08_09AM = NewDateTime(DateForFeb08, TimeFor09AM, nil)
var DateTimeForFeb08_12PM = NewDateTime(DateForFeb08, TimeFor12PM, nil)
var DateTimeForFeb08_03PM = NewDateTime(DateForFeb08, TimeFor03PM, nil)

var DateTimeForFeb14_06PM = NewDateTime(DateForFeb14, TimeFor06PM, nil)

var DateTimeForFeb15_09AM = NewDateTime(DateForFeb15, TimeFor09AM, nil)
var DateTimeForFeb15_12PM = NewDateTime(DateForFeb15, TimeFor12PM, nil)

var DateTimeForFeb22_09AM = NewDateTime(DateForFeb22, TimeFor09AM, nil)
var DateTimeForFeb22_12PM = NewDateTime(DateForFeb22, TimeFor12PM, nil)

var DateTimeForMar01_09AM = NewDateTime(DateForMar01, TimeFor09AM, nil)
var DateTimeForMar01_12PM = NewDateTime(DateForMar01, TimeFor12PM, nil)

var DateTimeFor2023Feb03_09AM = NewDateTime(DateFor2023Feb03, TimeFor09AM, nil)
var DateTimeFor2023Feb03_12PM = NewDateTime(DateFor2023Feb03, TimeFor12PM, nil)
var DateTimeFor2023Feb03_03PM = NewDateTime(DateFor2023Feb03, TimeFor03PM, nil)
var DateTimeFor2023Feb04_03PM = NewDateTime(DateFor2023Feb04, TimeFor03PM, nil)

var DateTimeFor2023Mar02_03PM = NewDateTime(DateFor2023Mar02, TimeFor03PM, nil)

var DateTimeForFeb03_09AM_ET = NewDateTime(DateForFeb03, TimeFor09AM, TimeZoneForET)
var DateTimeForFeb03_12PM_ET = NewDateTime(DateForFeb03, TimeFor12PM, TimeZoneForET)
var DateTimeForFeb03_03PM_ET = NewDateTime(DateForFeb03, TimeFor03PM, TimeZoneForET)

var DateTimeForFeb03_12PM_Eastern = NewDateTime(DateForFeb03, TimeFor12PM, TimeZoneForEastern)

// var DateTimeForFriday_12PM_ET = NewDateTime(DateForFriday, TimeFor12PM, TimeZoneForET)
// var DateTimeForFriday_03PM_ET = NewDateTime(DateForFriday, TimeFor03PM, TimeZoneForET)

var DateTimeFor2023Jan01_12AM_ET = NewDateTime(DateFor2023Jan1, TimeFor12AM, TimeZoneForET)

var DateTimeFor2023Feb01_09AM_ET = NewDateTime(DateFor2023Feb01, TimeFor09AM, TimeZoneForET)
var DateTimeFor2023Feb01_12PM_ET = NewDateTime(DateFor2023Feb01, TimeFor12PM, TimeZoneForET)
var DateTimeFor2023Feb01_03PM_ET = NewDateTime(DateFor2023Feb01, TimeFor03PM, TimeZoneForET)

var DateTimeFor2023Feb03_09AM_ET = NewDateTime(DateFor2023Feb03, TimeFor09AM, TimeZoneForET)
var DateTimeFor2023Feb03_12PM_ET = NewDateTime(DateFor2023Feb03, TimeFor12PM, TimeZoneForET)
var DateTimeFor2023Feb03_12PM_ADD0 = NewDateTime(DateFor2023Feb03, TimeFor12PM, TimeZoneForADD0)
var DateTimeFor2023Feb03_12PM_SUB0 = NewDateTime(DateFor2023Feb03, TimeFor12PM, TimeZoneForSUB0)
var DateTimeFor2023Feb03_12PM_SUB5 = NewDateTime(DateFor2023Feb03, TimeFor12PM, TimeZoneForSUB5)

var DateTimeFor2023Feb04_09AM_ET = NewDateTime(DateFor2023Feb04, TimeFor09AM, TimeZoneForET)
var DateTimeFor2023Feb04_12PM_ET = NewDateTime(DateFor2023Feb04, TimeFor12PM, TimeZoneForET)

var DateTimeFor2023Feb05_09AM_ET = NewDateTime(DateFor2023Feb05, TimeFor09AM, TimeZoneForET)
var DateTimeFor2023Feb05_12PM_ET = NewDateTime(DateFor2023Feb05, TimeFor12PM, TimeZoneForET)

var DateTimeFor2023Feb08_09AM_ET = NewDateTime(DateFor2023Feb08, TimeFor09AM, TimeZoneForET)
var DateTimeFor2023Feb08_12PM_ET = NewDateTime(DateFor2023Feb08, TimeFor12PM, TimeZoneForET)

var DateTimeFor2023Feb15_09AM_ET = NewDateTime(DateFor2023Feb15, TimeFor09AM, TimeZoneForET)
var DateTimeFor2023Feb15_12PM_ET = NewDateTime(DateFor2023Feb15, TimeFor12PM, TimeZoneForET)

var DateTimeFor2023Feb22_09AM_ET = NewDateTime(DateFor2023Feb22, TimeFor09AM, TimeZoneForET)
var DateTimeFor2023Feb22_12PM_ET = NewDateTime(DateFor2023Feb22, TimeFor12PM, TimeZoneForET)

var DateTimeFor2023Mar01_09AM_ET = NewDateTime(DateFor2023Mar01, TimeFor09AM, TimeZoneForET)
var DateTimeFor2023Mar01_12PM_ET = NewDateTime(DateFor2023Mar01, TimeFor12PM, TimeZoneForET)

var TimeFor12AM = &Time{}
var TimeFor09AM = &Time{Hour: 9}
var TimeFor12PM = &Time{Hour: 12}
var TimeFor03PM = &Time{Hour: 15}
var TimeFor06PM = &Time{Hour: 18}

var TimeZoneForEastern = &TimeZone{Name: "Eastern"}
var TimeZoneForET = &TimeZone{Abbreviation: "ET"}
var TimeZoneForADD0 = &TimeZone{Offset: "+00:00"}
var TimeZoneForSUB0 = &TimeZone{Offset: "-00:00"}
var TimeZoneForSUB5 = &TimeZone{Offset: "-05:00"}

type parseTest struct {
	dateMode string
	minDT    *DateTime

	in       string
	want     *DateTimeRanges
	isBroken bool
}

func TestParse(t *testing.T) {
	if os.Getenv("DEBUG") == "true" {
		glr.DoDebug = true
		DoDebug = true
	}

	tests := []parseTest{

		// Adapts tests from:
		// https://github.com/robintw/daterangeparser/blob/master/daterangeparser/test.py
		// https://github.com/vitalcode/date-time-range-parser/blob/master/src/test/scala/uk/vitalcode/dateparser/Examples.scala
		// https://github.com/waltzofpearls/dateparser
		//

		// ISO 8601 format

		{in: "2023", want: NewRangesWithStartDates(DateFor2023)},
		{in: "2023-02", want: NewRangesWithStartDates(DateFor2023Feb), isBroken: true},
		{in: "2023-02-03", want: DateRangesFor2023Feb03},
		{in: "2023-02-03T", want: DateRangesFor2023Feb03},
		// {in: "2023-02-03T12", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "2023-02-03T12:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM), isBroken: true},
		{in: "2023-02-03T12:00:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "2023-02-03T12:00:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "2023-02-03T12:00:00Z", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "2023-02-03T12:00:00+00:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ADD0)},
		{in: "2023-02-03T12:00:00-00:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_SUB0)},
		{in: "2023-02-03T12:00:00-05:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_SUB5)},

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
		{in: "Feb 2023", want: DateRangesFor2023Feb},
		{in: "02.03", want: DateRangesForFeb03, dateMode: "na"},
		{in: "02.03", want: DateRangesForMar02, dateMode: "rest"},
		{in: "02.03", want: DateRangesFor2023Feb03, dateMode: "na", minDT: DateTimeFor2023Jan01_12AM_ET},
		{in: "02.03", want: DateRangesFor2023Mar02, dateMode: "rest", minDT: DateTimeFor2023Jan01_12AM_ET},
		{in: "02.03.", want: DateRangesForFeb03, dateMode: "na"},
		{in: "02.03.", want: DateRangesForMar02, dateMode: "rest"},
		{in: "02.03.", want: DateRangesFor2023Feb03, dateMode: "na", minDT: DateTimeFor2023Jan01_12AM_ET},
		{in: "02.03.", want: DateRangesFor2023Mar02, dateMode: "rest", minDT: DateTimeFor2023Jan01_12AM_ET},
		{in: "2/3/2023", want: DateRangesFor2023Feb03, dateMode: "na"},
		{in: "2/3/2023", want: DateRangesFor2023Mar02, dateMode: "rest"},

		// Extra tokens
		{in: "Feb 3 Google Calendar ICS", want: DateRangesForFeb03},
		{in: "Updated: Feb 3", want: DateRangesForFeb03},
		{in: "Workshop Update (2/3/23)", want: DateRangesFor2023Feb03, dateMode: "na", isBroken: true},
		{in: "Workshop: Feb 3 2023  VIRTUAL", want: DateRangesFor2023Feb03},
		{in: "Release date: February 3, 2023", want: DateRangesFor2023Feb03},
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
		{in: "Feb 1, 2, 3, 4", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04), isBroken: true},
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
		{in: "1-3 Feb and 2 Mar", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForMar02), isBroken: true},
		{in: "1-3 & 5 February", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb05), isBroken: true},
		{in: "1-4 & 6 February", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb06), isBroken: true},

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
		{in: "2/1, 2/2, 3/2, 3/3", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForMar02, DateForMar03), dateMode: "na"},
		{in: "1/2, 2/2, 2/3, 3/3", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForMar02, DateForMar03), dateMode: "rest"},
		// 5 Mondays 3/17 & 3/31, 4/14 & 4/28, 5/12
		{in: "5 Wednesdays 2/1 & 2/8 & 2/15 & 2/22, 3/2", want: NewRangesWithStartDates(DateForFeb01, DateForFeb08, DateForFeb15, DateForFeb22, DateForMar01), isBroken: true},

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
		{in: "Join today at 12pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ET), minDT: DateTimeFor2023Feb03_09AM_ET},
		// Join today for Day 2 at 10am PST
		{in: "Join today for Day 2 at 12pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ET), minDT: DateTimeFor2023Feb03_09AM_ET, isBroken: true},
		{in: "Tomorrow at 12pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb04_12PM_ET), minDT: DateTimeFor2023Feb03_09AM_ET},

		// Relative Z
		{in: "Today Friday, 12pm ET", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET), minDT: DateTimeFor2023Feb03_09AM_ET, isBroken: true},

		// MDT
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "Feb 3 12:00 PM", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "February 3 @ 12:00 PM", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "February 3  12 p.m.", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "Date:Thu 03 Feb, Time:12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "Starting February 3rd at 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		//    Friday, 2/14: **Love is Listening and Art: Social + Listening Art Sessions** at 6pm facilitated by Lauren V
		{in: "Friday 2/14: **Love is Listening and Art: Social + Listening Art Sessions** at 6pm facilitated by Lauren V", want: NewRangesWithStartDateTimes(DateTimeForFeb14_06PM), dateMode: "na"},

		// MDT with min
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ET), minDT: DateTimeFor2023Jan01_12AM_ET},

		// MDTZ
		{in: "Feb 3 12pm ET", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET)},
		{in: "Feb 3 12pm (ET)", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET)},
		{in: "Feb 3 12pm - ET", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET)},
		{in: "Feb 3 12pm in ET", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET)},
		{in: "Feb 3 12pm Eastern", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_Eastern)},
		{in: "Feb 3 12pm US/Eastern", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET)},
		{in: "FEBRUARY 3RD 12 PM ET, ON FRIDAY", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET)},
		// {in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_East), minDT: DateTimeForEast, wantDiff: true},
		{in: "Starting February 3rd at 12pm (ET) - Virtually.", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET)},
		{in: "Starts Friday 2/3 at 9:00 am ET", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET), dateMode: "na", isBroken: true},

		// DMT
		{in: "Date:Fri 03 Feb, Time:3.00pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_03PM)},
		{in: "Friday 3 Feb 3:00pm (doors) | 11pm (curfew)", want: NewRangesWithStartDateTimes(DateTimeForFeb03_03PM)},
		// Need to update sorting algorithm for this.
		{in: "Fri, 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM)},
		{in: "Thu, 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Mar02_03PM)},
		{in: "Fri, 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM)},
		{in: "Thu, 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Mar02_03PM)},

		// MDYT
		{in: "Feb. 3, 2023 12:00pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "Feb 3, 2023 @ 12:00 PM", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "Friday, February 3rd 2023 from 12:00 PM", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},

		// MDYTT
		// Not sure if this is a range or multiple.
		{in: "Feb. 3, 2023 12:00pm, 3:00pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM, DateTimeFor2023Feb03_03PM), isBroken: true},

		// MDYTZ
		{in: "Feb 3 2023 12pm ET", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ET)},
		//    Sep 18, 2024 11:30 AM Eastern Time (US and Canada)
		{in: "Feb 3, 2023 12:00 PM Eastern Time (US and Canada)", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ET)},

		// DMY
		{in: "3rd Feb 2023 9:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_09AM)},
		{in: "3rd Feb 2023 9:00am", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_09AM)},
		{in: "3rd Feb 2023 3:00pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM)},

		// TZMD
		{in: "12:00 pm ET February 3rd", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_ET)},

		//
		// Date Time Ranges

		// MDTT
		// Need to fix parser for this.
		{in: "February 3: 9am - 12pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},
		{in: "Feb 3 9am - 12pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},
		{in: "Feb 3 @ 9:00 AM - Feb 3 @ 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},
		{in: "February, 3 9:00 - 15:00", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_03PM)},
		{in: "Friday, February 3rd from 12 - 3pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_12PM, DateTimeForFeb03_03PM)},
		{in: "Feb, 3rd from 9 am-3.00 pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_03PM)},
		{in: "February 3 + 4, 9 am - 12 pm each day", want: NewRanges(
			NewRange(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET),
			NewRange(DateTimeFor2023Feb04_09AM_ET, DateTimeFor2023Feb04_12PM_ET)), isBroken: true},
		//    November 9 + 10: In-person at SeekHealing Asheville,10 am - 7 pm each day
		{in: "February 1 + 2: In-person at SeekHealing Asheville,12 pm - 3 pm each day", want: NewRanges(
			NewRange(DateTimeForFeb01_12PM, DateTimeForFeb01_03PM),
			NewRange(DateTimeForFeb02_12PM, DateTimeForFeb02_03PM)), isBroken: true},
		{in: "THIS Friday: February 3 \n 12-3:00pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_12PM, DateTimeForFeb03_03PM)},
		//    2 Tuesdays Jan 7th & 21st 6:30p-8:30p
		{in: "2 Wednesdays Feb 1st & 8th 12:00p-3:00p", want: NewRanges(
			NewRange(DateTimeForFeb01_12PM, DateTimeForFeb01_03PM),
			NewRange(DateTimeForFeb08_12PM, DateTimeForFeb08_03PM)), isBroken: true},
		//    10 Mondays 6:30pm-8:30pm March 3rd - May 5th
		{in: "5 Wednesdays 9:00am-12:00pm February 1st - March 1st", want: NewRanges(
			NewRange(DateTimeForFeb01_09AM, DateTimeForFeb01_12PM),
			NewRange(DateTimeForFeb08_09AM, DateTimeForFeb08_12PM),
			NewRange(DateTimeForFeb15_09AM, DateTimeForFeb15_12PM),
			NewRange(DateTimeForFeb22_09AM, DateTimeForFeb22_12PM),
			NewRange(DateTimeForMar01_09AM, DateTimeForMar01_12PM)), isBroken: true},

		// MDTTZ
		{in: "Feb 3rd - 9.00 AM- 12pm ET", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM_ET, DateTimeForFeb03_12PM_ET)},
		{in: "February 3rd, 9-12pm ET", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM_ET, DateTimeForFeb03_12PM_ET)},
		{in: "Feb 3 2023 9am - 12pm ET", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET)},
		{in: "Feb 3 2023 9am ET to 12pm ET", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET)},
		{in: "Feb 3 @ 9:00 AM ET - Feb 3 @ 12:00 PM ET", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM_ET, DateTimeForFeb03_12PM_ET)},
		{in: "Feb 3, 2023, 9:00 AM ET - Feb 3, 2023, 12:00 PM ET", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET)},
		{in: "Friday, 2/3 12-3pm ET", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_12PM_ET, DateTimeForFeb03_03PM_ET), dateMode: "na"},
		// Tuesday 8/13 from 6:30 - 8:30 pm EST
		{in: "Friday 2/3 from 12:00 - 3 pm ET", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_12PM_ET, DateTimeForFeb03_03PM_ET)}, //, isBroken: true},
		//    Sunday, October 6, 9am - 4pm Pacific Time
		{in: "Friday, February 3, 9am - 12pm Eastern Time", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM_ET, DateTimeForFeb03_12PM_ET)},
		{in: "February 3, 2023 from 9:00 am to noon ET", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET)},
		// {in: "February 3, 2023 / 9:00 AM", want: NewRangesWithStartDateTimes(DateTimeForFeb03_09AM), wantDiff: true},
		{in: "February 3, 2023 / 9:00 AM - 12:00 PM ET", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET)},
		{in: "February 3rd, 12:00-3:00pm Eastern (New York) time", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_12PM, DateTimeForFeb03_03PM), isBroken: true},
		{in: "February 3rd & 4th, 9:00 am - noon Eastern time", want: NewRanges(
			NewRange(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET),
			NewRange(DateTimeFor2023Feb04_09AM_ET, DateTimeFor2023Feb04_12PM_ET)), isBroken: true},
		{in: "February 3rd - 5th, 9:00 am - noon ET each day", want: NewRanges(
			NewRange(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET),
			NewRange(DateTimeFor2023Feb04_09AM_ET, DateTimeFor2023Feb04_12PM_ET),
			NewRange(DateTimeFor2023Feb05_09AM_ET, DateTimeFor2023Feb05_12PM_ET)), isBroken: true},
		{in: "Wednesdays February 1st & 8th 12:00p-3:00p", want: NewRanges(
			NewRange(DateTimeForFeb01_12PM, DateTimeForFeb01_03PM),
			NewRange(DateTimeForFeb08_12PM, DateTimeForFeb08_03PM))},
		{in: "Wednesdays, February 1st, 8th, and 15th 9:00am - 12:00pm (ET)", want: NewRanges(
			NewRange(DateTimeFor2023Feb01_09AM_ET, DateTimeFor2023Feb01_12PM_ET),
			NewRange(DateTimeFor2023Feb08_09AM_ET, DateTimeFor2023Feb08_12PM_ET),
			NewRange(DateTimeFor2023Feb15_09AM_ET, DateTimeFor2023Feb15_12PM_ET)), isBroken: true},
		{in: "Wednesdays - February 1, 8 9:00 AM - 12:00 PM ET", want: NewRanges(
			NewRange(DateTimeFor2023Feb01_09AM_ET, DateTimeFor2023Feb01_12PM_ET),
			NewRange(DateTimeFor2023Feb08_09AM_ET, DateTimeFor2023Feb08_12PM_ET)), isBroken: true},
		//    Tuesdays – March 18, 25, and April 1, 8, 15, 22 10:00 AM – 12:30 PM PST
		{in: "Wednesdays - February 1, 8, 15, 22, and March 1 9:00 AM - 12:00 PM ET", want: NewRanges(
			NewRange(DateTimeFor2023Feb01_09AM_ET, DateTimeFor2023Feb01_12PM_ET),
			NewRange(DateTimeFor2023Feb08_09AM_ET, DateTimeFor2023Feb08_12PM_ET),
			NewRange(DateTimeFor2023Feb15_09AM_ET, DateTimeFor2023Feb15_12PM_ET),
			NewRange(DateTimeFor2023Feb22_09AM_ET, DateTimeFor2023Feb22_12PM_ET),
			NewRange(DateTimeFor2023Mar01_09AM_ET, DateTimeFor2023Mar01_12PM_ET)), isBroken: true},

		// DTT
		{in: "Friday 12 to 3 PM Eastern", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_12PM_ET, DateTimeFor2023Feb03_12PM_ET), minDT: DateTimeFor2023Feb01_09AM_ET}, //, isBroken: true},

		// DMTT
		{in: "3 Feb 9am - 12pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},

		// MDYT
		{in: "Feb 3 2023 12pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		// MDYTT
		{in: "Thursday, February 3, 2023 9:00 AM 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM)},
		{in: "Thursday, February 3rd 2023 from 9:00 AM to 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM)},
		// DMYTT
		{in: "When 3 Feb 2023 9:00 AM - 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM)},
		{in: "Fr. 3. Feb. 2023, 9:00-ca.12:00", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM), isBroken: true},

		// TDMY
		{in: "9:00am 3rd Feb - 4th Feb 3:00pm 2023", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb04_03PM)},
		{in: "9:00am on 3rd Feb - 4th Feb at 3:00pm 2023", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb04_03PM), isBroken: true},
		// Not sure how to parse this one.
		// {in: "(2 Feb 2023 - 3 Feb 2023) 09:00 15:00", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeFor2023Feb04_03PM)},

		// Both
		// Need to fix parser for these
		{in: "02.03.2023", want: DateRangesFor2023Feb03, dateMode: "na"},
		{in: "02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM), dateMode: "na"},
		{in: "Th , 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM), dateMode: "na"},
		{in: "Th , 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM), dateMode: "rest"},

		//
		// Failures

		{in: "814-555-1212", want: nil},
		{in: "814-555-1212 x123", want: nil},
		{in: "102 W. Mahoning Street. Punxsutawney, PA 15767", want: nil},
		// Need to fix these
		{in: "We may request cookies to be set on your device.", want: nil, isBroken: true},
		{in: "Winter Retreat for 6-12th graders!", want: nil},
		{in: "For 6th-12th grade students @ SpringHill Camp", want: nil},

		//
		// TODO
		//

		// 5 Mondays 3/17 & 3/31, 4/14 & 4/28, 5/12
		// Tues/Thurs 6:30p-9:00p March/April
		// 5 Mondays March thru May

		// Tuesdays and Thursdays March or April 6:30pm - 9:00pm
		// September 19 - October 24: Via Zoom with SeekHealing Online, every Thursday from 12 - 2 pm EST
		// October 12 + 13: In-person at SeekHealing Waynesville, 10 am - 7 pm each day
		// November 9 + 10: In-person at SeekHealing Asheville,10 am - 7 pm each day
		// every Thursday from 12 - 1:30 pm EST through August 29th

		// Wait until we are outputing ICAL before deciding how to parse repeating dates like the ones below.

		// September 22 - November 3, 2024\nSundays 11:00–12:30 CT / 18.00-19.30 CET
		// * August 8 - September 12th: Via Zoom with SeekHealing Online, every Thursday from 6:30 - 8:30 pm EST
		// * August 10 + 11: In-person at SeekHealing Waynesville, 10 am - 7 pm each day
		// * November 9 + 10: In-person at SeekHealing Asheville,10 am - 7 pm each day
		// every Thursday from 5 - 6:15 pm at SeekHealing Asheville beginning August 8th
		// 24-26 October (in person) & Weds 31st Integration eve (online)
		// Tuesdays January 7th & 21st 6:30p-8:30p
		// Every Tuesday beginning September 10th from 12:30 - 2 pm
		// September 19 - October 24, every Thursday from 12 - 2 pm EST
		// September 4 - 12 \n Courses Begin September 17
		// NEW Sundays at 9 am ET - Starts October 20
		// Monday, September 2 \n PT (AZ): 10:30 am, MT: 11:30 am, CT: 12:30 pm, ET: 1:30 pm, London: 5:30 pm, Sweden and France: 6:30 pm, Israel: 7:30 pm
		// Sunday 10 to 11:10 AM Eastern
		// Sunday 1:30 PM
		// Thursday 8AM
		// Program Begins Thursday, January 16, 2025 | 36 Classes PT: 5:00 pm, MT: 6:00 pm, CT: 7:00, ET: 8:00 pm 90 Minute Sessions.
		// Part 1: 15th–20th March 2025, Part 2: 19th-24th October 2025
		// CRN 66932, SLFO NC025, Dates: 8 Tues, Jan 28 - Mar 18, Time, 6:15 pm – 8:30pm, Pacific
		// Starts: Tuesday, April 1, 2025 Ends: Tuesday, May 27, 2025 Meets: Online, for 9 consecutive Tuesday evenings, from 7:00 PM to 8:30 PM PST
		// Tuesdays 7 to 9 pm ET
		// April 14 - 25 (M-W-F; M-W-F)
		// "October 8th - 10.00am- 3pm\u00a0MST"
		// Weekly on Mondays
		// Wednesdays at 6:30pm ET
		// 2nd and 4th Tuesdays 7 to 9 pm ET
		// Tuesdays 11:45am-12:00pm ET, Mindful Pause
		// Tuesdays 7 to 9 pm ET
		// Join today for Day 2 at 10am PST
		// Fridays 3:00 - 5:00 pm EASTERN
		// Fridays, February 7 - December 5, 2025 (45 sessions)
		// 5 Mondays 3/17 & 3/31, 4/14 & 4/28, 5/12
		// Online 5 Mondays 3/17 & 3/31, 4/14 & 4/28, 5/12
		// 1/16/2025
		// 3/14/2025
		// Sunday, February 16th, 2025\n17:00 - 19:00 CET / 16:00 - 18:00 UTC (Find your local start time here)
		// Beginning February 7, 2025, Fridays 3:00 - 5:00 pm EASTERN"
		// October 20 - 31 (M-W-F; M-W-F)"
		// Starting 2nd and 4th Tuesdays 7 to 9 pm ET"

		// # Straddling end of year
		// ("25 Dec - 2 Jan 2016", "25/12/2015", "02/01/2016"),
		// ("18 Nov 2015 to 14th Feb 2016", "18/11/2015", "14/02/2016"),
		// ("18 Nov 2010 to 14th Feb 2016", "18/11/2010", "14/02/2016"),

		// {in: "Feb 3 12:00 PM 12:00", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		// {in: "Feb 3 3:00 PM 15:00", want: NewRangesWithStartDateTimes(DateTimeForFeb03_03PM)},
		// {in: "Feb 3 9:00 AM 09:00 Feb 3 3:00 PM 15:00", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_03PM)},
		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_03PM)},
		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00"
		// {in: "Feb 3 2023 9:00 AM 3:00 PM 09:00 15:00 Google Calendar ICS"
		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00"
		// {in: "Fri, Apr 14, 2023 9:00 AM 09:00 Sat, Apr 15, 2023 5:00 PM 17:00"

		// "(1 Jan 2017 - 3 Jan 2017) Tuesday 11:00 13:00" in expected(
		// "(1 Jan 2016 - 4 Jan 2016) Monday 11:00 13:00 Tuesday 14:00 15:00 Friday 16:05 17:20 Sunday 19:30 20:45" in expected(
		// "(3 Feb 2017) Friday 19:30 21:30" in expected(
		// "Select date Tue 19 September 12:00pm Tue 19 September 2:00pm Tue 19 Sep 4:00pm (last few)" in expected(
		// "Select date Thu 15 September 7:45pm 8:50pm Fri 16 September 7:45pm - 20:45 Sat 17 October 7:45pm to 21:10" in expected(

		// Sat 5pm to 2am
		// Sun/Thu 5pm–1am"
		// Fri 5pm to 2am
		// "Doors: 8PM / Show: 9PM / 21+"
		// "12PM / 21+ / Free"
		// "Doors: 8PM / Show: 9PM / 21+RSVP DOES NOT GUARANTEE ENTRY"
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

	fmt.Printf("acceptBrokenTests: %t\n", acceptBrokenTests)
	percent := float64(failed) / float64(len(tests)) * 100
	fmt.Printf("TestParse: %.2f%% of tests failed (%d/%d)\n", percent, failed, len(tests))
}

func testParseFn(t *testing.T, tc parseTest) func(*testing.T) {
	return func(t *testing.T) {
		got, err := Parse(tc.minDT, tc.dateMode, tc.in)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if !tc.isBroken && (got == nil && tc.want == nil) {
			return
		}

		// Clear Weekdays which would just be annoying to test. Trust that we would find with semantic checks.
		// if got != nil {
		// 	for _, dtr := range got.Items {
		// 		if dtr.Start != nil && dtr.Start.Date != nil {
		// 			dtr.Start.Date.Weekday = 0
		// 		}
		// 		if dtr.End != nil && dtr.End.Date != nil {
		// 			dtr.End.Date.Weekday = 0
		// 		}
		// 	}
		// }

		diff := cmp.Diff(got, tc.want, cmpopts.IgnoreUnexported(Date{})) //, protocmp.Transform())
		// fmt.Printf("diff == \"\": %t\n", diff == "")
		// fmt.Printf("tc.isBroken: %t\n", tc.isBroken)

		if acceptBrokenTests && tc.isBroken {
			if diff == "" {
				t.Fatal("unexpected sameness")
			}
			return
		}
		if diff != "" &&
			(!acceptBrokenTests ||
				(acceptBrokenTests && !tc.isBroken)) {
			pp.Default.SetColoringEnabled(false)
			// fmt.Printf("got:\n%s\n", litter.Sdump(got))
			// fmt.Printf("want:\n%s\n", litter.Sdump(tc.want))
			fmt.Printf("got vs. want:\n%s\n", litter.Sdump([]*DateTimeRanges{got, tc.want}))
			t.Fatalf("unexpected difference:\n%v", diff)
		}
	}
}
