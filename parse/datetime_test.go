package parse

import (
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/findyourpaths/phil/glr"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

// Run with
// rm datetime_glr.go; time GOWORK=off go generate -v ./...
// time go test -v ./...
// DEBUG=true time go test -v -run '^TestExtractDatetimesRanges/102'

var DateForFeb = civil.Date{Month: time.February}
var DateForMar = civil.Date{Month: time.March}

var DateForFeb01 = civil.Date{Month: time.February, Day: 1}
var DateForFeb02 = civil.Date{Month: time.February, Day: 2}
var DateForFeb03 = civil.Date{Month: time.February, Day: 3}
var DateForFeb04 = civil.Date{Month: time.February, Day: 4}
var DateForFeb05 = civil.Date{Month: time.February, Day: 5}
var DateForMar02 = civil.Date{Month: time.March, Day: 2}
var DateForMar03 = civil.Date{Month: time.March, Day: 3}
var DateForMar04 = civil.Date{Month: time.March, Day: 4}
var DateForApr03 = civil.Date{Month: time.April, Day: 3}

var DateRangesForFeb03 = NewRangesWithStartDates(DateForFeb03)
var DateRangesFromFeb03ToFeb04 = NewRangesWithStartEndDates(DateForFeb03, DateForFeb04)

var DateFor2023 = civil.Date{Year: 2023}
var DateFor2024 = civil.Date{Year: 2024}

var DateFor2023Feb = civil.Date{Year: 2023, Month: time.February}
var DateFor2023Mar = civil.Date{Year: 2023, Month: time.March}

var DateRangesFor2023Feb = NewRangesWithStartDates(DateFor2023Feb)

var DateRangesFor2023Feb03 = NewRangesWithStartDates(DateFor2023Feb03)

var DateFor2023Feb01 = civil.Date{Year: 2023, Month: time.February, Day: 1}
var DateFor2023Feb02 = civil.Date{Year: 2023, Month: time.February, Day: 2}
var DateFor2023Feb03 = civil.Date{Year: 2023, Month: time.February, Day: 3}
var DateFor2023Feb04 = civil.Date{Year: 2023, Month: time.February, Day: 4}
var DateFor2023Feb05 = civil.Date{Year: 2023, Month: time.February, Day: 5}
var DateFor2023Mar02 = civil.Date{Year: 2023, Month: time.March, Day: 2}
var DateFor2023Mar03 = civil.Date{Year: 2023, Month: time.March, Day: 3}

var DateRangesFrom2023Feb03To2023Feb04 = NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)

var DateTimeForFeb03_09AM = &DateTimeTZ{DateTime: civil.DateTime{Date: DateForFeb03, Time: TimeFor09AM}}
var DateTimeForFeb03_12PM = &DateTimeTZ{DateTime: civil.DateTime{Date: DateForFeb03, Time: TimeFor12PM}}
var DateTimeForFeb03_03PM = &DateTimeTZ{DateTime: civil.DateTime{Date: DateForFeb03, Time: TimeFor03PM}}

var DateTimeFor2023Feb03_09AM = &DateTimeTZ{DateTime: civil.DateTime{Date: DateFor2023Feb03, Time: TimeFor09AM}}
var DateTimeFor2023Feb03_12PM = &DateTimeTZ{DateTime: civil.DateTime{Date: DateFor2023Feb03, Time: TimeFor12PM}}
var DateTimeFor2023Feb03_03PM = &DateTimeTZ{DateTime: civil.DateTime{Date: DateFor2023Feb03, Time: TimeFor03PM}}
var DateTimeFor2023Feb04_03PM = &DateTimeTZ{DateTime: civil.DateTime{Date: DateFor2023Feb04, Time: TimeFor03PM}}

var DateTimeForFeb03_09AM_PST = &DateTimeTZ{DateTime: civil.DateTime{Date: DateForFeb03, Time: TimeFor12PM}, TimeZone: "PST"}
var DateTimeForFeb03_12PM_PST = &DateTimeTZ{DateTime: civil.DateTime{Date: DateForFeb03, Time: TimeFor12PM}, TimeZone: "PST"}

var DateTimeFor2023Feb03_09AM_PST = &DateTimeTZ{DateTime: civil.DateTime{Date: DateFor2023Feb03, Time: TimeFor09AM}, TimeZone: "PST"}
var DateTimeFor2023Feb03_12PM_PST = &DateTimeTZ{DateTime: civil.DateTime{Date: DateFor2023Feb03, Time: TimeFor12PM}, TimeZone: "PST"}

var TimeFor09AM = civil.Time{Hour: 9}
var TimeFor12PM = civil.Time{Hour: 12}
var TimeFor03PM = civil.Time{Hour: 15}

func TestExtractDatetimesRanges(t *testing.T) {
	if os.Getenv("DEBUG") == "true" {
		glr.DoDebug = true
		DoDebug = true
	}

	type test struct {
		dateMode string
		year     int
		timeZone string

		in   string
		want *DateTimeTZRanges
	}

	tests := []test{

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
		// {in: "2023-02-03T12", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "2023-02-03T12:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "2023-02-03T12:00:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},

		//
		// Date
		//

		// MD
		{in: "Feb 3", want: DateRangesForFeb03},
		{in: "Thu Feb 3", want: DateRangesForFeb03},
		{in: "Thu 3 Feb", want: DateRangesForFeb03},
		// DM
		{in: "3 Feb", want: DateRangesForFeb03},
		{in: "3, Feb", want: DateRangesForFeb03},
		{in: "3rd of Feb", want: DateRangesForFeb03},
		{in: "thu 3 Feb", want: DateRangesForFeb03},

		// MDY
		{in: "Feb 3 2023", want: DateRangesFor2023Feb03},
		{in: "February 3 2023", want: DateRangesFor2023Feb03},
		{in: "February 3, 2023", want: DateRangesFor2023Feb03},
		{in: "February 3rd, 2023", want: DateRangesFor2023Feb03},
		{in: "Thu Feb 3, 2023", want: DateRangesFor2023Feb03},
		{in: "Thursday Feb 3rd 2023", want: DateRangesFor2023Feb03},
		{in: "Time:Feb 3 2023", want: DateRangesFor2023Feb03},
		// DMY
		{in: "3 Feb 2023", want: DateRangesFor2023Feb03},
		{in: "3rd Feb 2023", want: DateRangesFor2023Feb03},
		{in: "3 February, 2023", want: DateRangesFor2023Feb03},
		{in: "Thursday 3rd Feb 2023", want: DateRangesFor2023Feb03},

		// Both
		{in: "02.03", want: DateRangesForFeb03, dateMode: "na"},
		{in: "02.03", want: NewRangesWithStartDates(DateForMar02), dateMode: "rest"},
		{in: "02.03.", want: DateRangesForFeb03, dateMode: "na"},
		{in: "02.03.", want: NewRangesWithStartDates(DateForMar02), dateMode: "rest"},
		{in: "2/3/2023", want: DateRangesFor2023Feb03, dateMode: "na"},
		{in: "2/3/2023", want: NewRangesWithStartDates(DateFor2023Mar02), dateMode: "rest"},

		{in: "Feb 2023", want: DateRangesFor2023Feb},

		//
		// Dates
		//

		// MD
		{in: "Feb 1, 2", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02)},
		{in: "Feb 1, 2, 3", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03)},
		{in: "Feb 1, 2, 3, 4", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04)},
		{in: "Feb 1, 2, 3, 4, 5", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb05)},
		//		{in: "February 1, 2, March 2, 3, and 4, April 3.", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForMar02, DateForMar03, DateForMar04, DateForApr03)},
		// DM
		{in: "1, 2 Feb", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02)},
		{in: "1, 2, 3 Feb", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03)},
		{in: "1, 2, 3, 4 Feb", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04)},
		{in: "1, 2, 3, 4, 5 Feb", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb05)},
		{in: "1, 2, 3 Feb and 2 Mar", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForMar02)},

		// MDY
		{in: "Feb 1, 2 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02)},
		{in: "Feb 1, 2, 3 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03)},
		{in: "Feb 1, 2, 3, 4 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04)},
		{in: "Feb 1, 2, 3, 4, 5 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04, DateFor2023Feb05)},
		{in: "Feb 1, 2, 3 and Mar 2 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Mar02)},
		// DMY
		{in: "1, 2 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02)},
		{in: "1, 2, 3 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03)},
		{in: "1, 2, 3, 4 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04)},
		{in: "1, 2, 3, 4, 5 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04, DateFor2023Feb05)},
		{in: "1, 2, 3 Feb and 2 Mar 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Mar02)},

		// MD
		{in: "Feb 3 Mar 2", want: NewRangesWithStartDates(DateForFeb03, DateForMar02)},
		// MDY
		{in: "Feb 3 Mar 2 2023", want: NewRangesWithStartDates(DateFor2023Feb03, DateFor2023Mar02)},

		//
		// Date Range
		//

		{in: "2023 - 2024", want: NewRangesWithStartEndDates(DateFor2023, DateFor2024)},
		{in: "Feb - Mar", want: NewRangesWithStartEndDates(DateForFeb, DateForMar)},
		{in: "Feb 2023 - Mar 2023", want: NewRangesWithStartEndDates(DateFor2023Feb, DateFor2023Mar)},

		// MD
		{in: "Feb 3rd-4th", want: DateRangesFromFeb03ToFeb04},
		{in: "Feb 3 - Mar 2", want: NewRangesWithStartEndDates(DateForFeb03, DateForMar02)},
		{in: "Feb 3 to Mar 2", want: NewRangesWithStartEndDates(DateForFeb03, DateForMar02)},
		// DM
		{in: "3-4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "3 Feb - 4 Feb", want: DateRangesFromFeb03ToFeb04},
		{in: "Thu Feb 3 - Fri Feb 4", want: DateRangesFromFeb03ToFeb04},
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
		{in: "Thu Feb 3, 2023 - Fri Feb 4, 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri 3 Feb - Sat 4 February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri 3rd Feb - Sat 4th February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri Feb 3rd - Sat 4th February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri 3rd Feb - 4th Sat February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "Fri Feb 3rd - 4th Sat February 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "February 3 - March 2, 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Mar02)},
		// DMY
		{in: "3-4 Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "3-4 Feb. 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "3-4 February, 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "3rd-4th Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "3 Feb 2023 - 4 Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "From 3rd to 4th, Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},
		{in: "beginning 3rd to 4th Feb 2023", want: DateRangesFrom2023Feb03To2023Feb04},

		// YMD
		{in: "2023, Feb 3 - 2023, Feb 4", want: DateRangesFrom2023Feb03To2023Feb04},

		//
		// Date Ranges
		//

		// MD
		{in: "Feb 1-2, 3-4", want: NewRanges(NewRangeWithStartEndDates(DateForFeb01, DateForFeb02), NewRangeWithStartEndDates(DateForFeb03, DateForFeb04))},
		{in: "Feb 1-2, 3-4 2023", want: NewRanges(NewRangeWithStartEndDates(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04))},

		// MD
		{in: "Feb 1-2; Mar 2-3", want: NewRanges(NewRangeWithStartEndDates(DateForFeb01, DateForFeb02), NewRangeWithStartEndDates(DateForMar02, DateForMar03))},
		// DM
		{in: "1-2 Feb; 2-3 Mar", want: NewRanges(NewRangeWithStartEndDates(DateForFeb01, DateForFeb02), NewRangeWithStartEndDates(DateForMar02, DateForMar03))},

		// MDY
		{in: "Feb 1-2; Mar 2-3 2023", want: NewRanges(NewRangeWithStartEndDates(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEndDates(DateFor2023Mar02, DateFor2023Mar03))},
		// DMY
		{in: "1-2 Feb; 2-3 Mar 2023", want: NewRanges(NewRangeWithStartEndDates(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEndDates(DateFor2023Mar02, DateFor2023Mar03))},

		//
		// Date Times
		//

		// MD
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_PST), timeZone: "PST"},
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_PST), timeZone: "America/Los_Angeles"},
		{in: "Feb 3 12:00 PM", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "Date:Thu 03 Feb, Time:12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},

		// DM
		{in: "Date:Thu 03 Feb, Time:3.00pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_03PM)},

		// MDY
		{in: "Feb. 3, 2023 12:00pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "Thursday, February 3rd 2023 from 12:00 PM", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},

		// DMY
		{in: "3rd Feb 2023 9:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_09AM)},
		{in: "3rd Feb 2023 9:00am", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_09AM)},
		{in: "3rd Feb 2023 3:00pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM)},

		//
		// Date Time Ranges
		//

		{in: "Feb 3 9am - 12pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},
		{in: "Feb 3 @ 9:00 AM - Feb 3 @ 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},
		{in: "February, 3 9:00 - 15:00", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_03PM)},

		{in: "Feb 3 2023 12pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "Th , 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM), dateMode: "na"},
		{in: "Th , 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM), dateMode: "rest"},

		{in: "When 3 Feb 2023 9:00 AM - 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM)},

		{in: "Thursday, February 3, 2023 9:00 AM 12:00 PM Google Calendar ICS", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM)},
		{in: "Thursday, February 3rd 2023 from 9:00 AM to 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM)},

		{in: "3 Feb 9am - 12pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},

		{in: "9:00am 3rd Feb - 4th Feb 3:00pm 2023", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb04_03PM)},
		{in: "9:00am on 3rd Feb - 4th Feb at 3:00pm 2023", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb04_03PM)},

		// # Straddling end of year
		// ("25 Dec - 2 Jan 2016", "25/12/2015", "02/01/2016"),
		// ("18 Nov 2015 to 14th Feb 2016", "18/11/2015", "14/02/2016"),
		// ("18 Nov 2010 to 14th Feb 2016", "18/11/2010", "14/02/2016"),

		// {in: "Feb 3 12:00 PM 12:00", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		// {in: "Feb 3 3:00 PM 15:00", want: NewRangesWithStartDateTimes(DateTimeForFeb03_03PM)},
		// {in: "Feb 3 9:00 AM 09:00 Feb 3 3:00 PM 15:00", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_03PM)},
		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_03PM)},
		// {in: "Feb 3 12pm PST", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_PST)},
		// {in: "Feb 3 2023 12pm PST", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_PST)},
		// {in: "Feb 3 2023 9am - 12pm PST", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_PST, DateTimeFor2023Feb03_12PM_PST)},
		// {in: "Feb 3 2023 9am PST to 3pm PST", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_PST, DateTimeFor2023Feb03_12PM_PST)},
		// {in: "Feb 3 @ 9:00 AM PST - Feb 3 @ 3:00 PM PST", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM_PST, DateTimeFor2023Feb03_12PM_PST)},
		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00"
		// {in: "Feb 3 2023 9:00 AM 3:00 PM 09:00 15:00 Google Calendar ICS"
		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00"
		// {in: "Fri, Apr 14, 2023 9:00 AM 09:00 Sat, Apr 15, 2023 5:00 PM 17:00"

		// "Sept. 10, 2017 11:00am, 3:00pm" in expected(
		// "Saturday 24 Jun 6:45pm (doors) | 11pm (curfew)" in expected(
		// "(1 Jan 2017 - 3 Jan 2017) 11:00 13:00" in expected(
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

		// FAILS
		// "269-781-9094 x13"
		// "1,095\u00a0Individuals (617 Park City, 478 Heber)"
		// "Burlingame, 2300 Adeline Dr, Burlingame, CA 94010, USA"
		// "317-270-4214"
		// "CCPC Update (5/18/20)"
		// "The Wisdom of the Enneagram Workshop: February 1-2, 2025 VIRTUAL | Enneagram Institute of Ohio"
		// "January 31st from 9 am-12.30 pm PST."
		// "January 28, 2025 @ 6:00 pm"
		// "Updated: 04/18/2024"
		// "We may request cookies to be set on your device."
		// "Apr 25, 2024, 7:00 PM PDT – Apr 28, 2024, 3:00 PM PDT"
		// "Winter Retreat for 6-12th graders!"
		// "For 6th-12th grade students @ SpringHill Camp"
		// "October 8th - 10.00am- 3pm\u00a0MST"
		// "The Three Instincts Workshop: April 5-6, 2025  VIRTUAL"
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%03d__%s", i, tc.in), func(t *testing.T) {
			got, err := ExtractDateTimeTZRanges(tc.year, tc.dateMode, tc.timeZone, tc.in)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected difference:\n%v", diff)
			}
		})
	}
}
