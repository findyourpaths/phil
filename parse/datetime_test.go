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

var DateForFeb01 = civil.Date{Month: time.February, Day: 1}
var DateForFeb02 = civil.Date{Month: time.February, Day: 2}
var DateForFeb03 = civil.Date{Month: time.February, Day: 3}
var DateForFeb04 = civil.Date{Month: time.February, Day: 4}
var DateForFeb05 = civil.Date{Month: time.February, Day: 5}
var DateForMar02 = civil.Date{Month: time.March, Day: 2}
var DateForMar03 = civil.Date{Month: time.March, Day: 3}

var DateRangesForFeb03 = NewRangesWithStartDates(DateForFeb03)
var DateRangesFromFeb03ToFeb04 = NewRangesWithStartEndDates(DateForFeb03, DateForFeb04)

var DateFor2023 = civil.Date{Year: 2023}
var DateFor2023Feb = civil.Date{Year: 2023, Month: time.February}

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

		// ISO 8601 format

		{in: "2023", want: NewRangesWithStartDates(DateFor2023)},
		{in: "2023-02", want: NewRangesWithStartDates(DateFor2023Feb)},
		{in: "2023-02-03", want: DateRangesFor2023Feb03},
		{in: "2023-02-03T", want: DateRangesFor2023Feb03},
		{in: "2023-02-03T12", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
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

		//
		// Dates
		//

		// MD
		{in: "Feb 1, 2", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02)},
		{in: "Feb 1, 2, 3", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03)},
		{in: "Feb 1, 2, 3, 4", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04)},
		{in: "Feb 1, 2, 3, 4, 5", want: NewRangesWithStartDates(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb05)},
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
		// DMY
		{in: "1, 2 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02)},
		{in: "1, 2, 3 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03)},
		{in: "1, 2, 3, 4 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04)},
		{in: "1, 2, 3, 4, 5 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04, DateFor2023Feb05)},
		// {input: "1, 2, 3 Feb and 2 Mar 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Mar02)},

		// MD
		{in: "Feb 3 Mar 2", want: NewRangesWithStartDates(DateForFeb03, DateForMar02)},
		// MDY
		{in: "Feb 3 Mar 2 2023", want: NewRangesWithStartDates(DateFor2023Feb03, DateFor2023Mar02)},

		//
		// Date Range
		//

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

		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_PST), timeZone: "PST"},
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_PST), timeZone: "America/Los_Angeles"},
		{in: "Feb 3 12:00 PM", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},

		// {in: "Feb 3 12:00 PM 12:00", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		// {in: "Feb 3 3:00 PM 15:00", want: NewRangesWithStartDateTimes(DateTimeForFeb03_03PM)},

		//
		// Date Time Ranges
		//

		{in: "Feb 3 9am - 12pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},
		{in: "Feb 3 @ 9:00 AM - Feb 3 @ 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},
		// {in: "Feb 3 9:00 AM 09:00 Feb 3 3:00 PM 15:00", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_03PM)},

		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_03PM)},

		{in: "Feb 3 2023 12pm", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "Th , 02.03.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM), dateMode: "na"},
		{in: "Th , 03.02.2023 - 15:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_03PM), dateMode: "rest"},

		{in: "When 3 Feb 2023 9:00 AM - 12:00 PM", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM)},

		// {in: "Thursday, February 3, 2023 9:00 AM 12:00 PM Google Calendar ICS", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM, DateTimeFor2023Feb03_12PM)},

		{in: "3 Feb 9am - 12pm", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM, DateTimeForFeb03_12PM)},

		/*
		   # Ignoring of times
		   ("14th July 1988 06:45", "14/7/1988", None),
		   ("14th July 1988 06.45am", "14/7/1988", None),
		   ("14th July 1988 3:30pm", "14/7/1988", None),
		   ("12:37 1st Jan - 17th Feb 19:00", "1/1/XXXX", "17/2/XXXX"),

		   # Things in different orders
		   ("July 14", "14/7/XXXX", None),
		   ("1990, Dec 29 - 1992, Dec 14", "29/12/1990", "14/12/1992"),
		   ("Thurs Nov 11th 1954", "11/11/1954", None),

		   # Dates with no days
		   ("July", "01/07/XXXX", "31/07/XXXX"),
		   ("Feb 2010", "01/02/2010", "28/02/2010"),
		   ("Feb 1996", "01/02/1996", "29/02/1996"),
		   ("Feb to Nov", "01/02/XXXX", "30/11/XXXX"),
		   ("Feb 2010 - Feb 2012", "01/02/2010", "29/02/2012"),

		   # Bare year
		   ("2013", "01/01/2013", "31/12/2013"),
		   ("1995 - 2010", "01/01/1995", "31/12/2010"),

		   # Straddling end of year
		   ("25 Dec - 2 Jan 2016", "25/12/2015", "02/01/2016"),
		   ("18 Nov 2015 to 14th Feb 2016", "18/11/2015", "14/02/2016"),
		   ("18 Nov 2010 to 14th Feb 2016", "18/11/2010", "14/02/2016"),

		   ("26-29 Oct. 2009", "26/10/2009", "29/10/2009")

		*/

		// {in: "Feb 3 12pm PST", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_PST)},

		// {in: "Feb 3 2023 12pm PST", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_PST)},

		// {in: "Feb 3 2023 9am - 12pm PST", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_PST, DateTimeFor2023Feb03_12PM_PST)},

		// {in: "Feb 3 2023 9am PST to 3pm PST", want: NewRangesWithStartEndDateTimes(DateTimeFor2023Feb03_09AM_PST, DateTimeFor2023Feb03_12PM_PST)},

		// {in: "Feb 3 @ 9:00 AM PST - Feb 3 @ 3:00 PM PST", want: NewRangesWithStartEndDateTimes(DateTimeForFeb03_09AM_PST, DateTimeFor2023Feb03_12PM_PST)},

		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00"
		// {in: "Feb 3 2023 9:00 AM 3:00 PM 09:00 15:00 Google Calendar ICS"
		// {in: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00"
		// {in: "Fri, Apr 14, 2023 9:00 AM 09:00 Sat, Apr 15, 2023 5:00 PM 17:00"
	}

	for i, tc := range tests {
		// for i, tc := range tests[len(tests)-3 : len(tests)-2] {
		// for i, tc := range tests[19:20] {
		// for i, tc := range tests[2:3] {
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
