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

var DateFor2023 = civil.Date{Year: 2023}
var DateFor2023Feb = civil.Date{Year: 2023, Month: time.February}

var DateFor2023Feb01 = civil.Date{Year: 2023, Month: time.February, Day: 1}
var DateFor2023Feb02 = civil.Date{Year: 2023, Month: time.February, Day: 2}
var DateFor2023Feb03 = civil.Date{Year: 2023, Month: time.February, Day: 3}
var DateFor2023Feb04 = civil.Date{Year: 2023, Month: time.February, Day: 4}
var DateFor2023Feb05 = civil.Date{Year: 2023, Month: time.February, Day: 5}
var DateFor2023Mar02 = civil.Date{Year: 2023, Month: time.March, Day: 2}
var DateFor2023Mar03 = civil.Date{Year: 2023, Month: time.March, Day: 3}

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
		{in: "2023-02-03", want: NewRangesWithStartDates(DateFor2023Feb03)},
		{in: "2023-02-03T", want: NewRangesWithStartDates(DateFor2023Feb03)},
		{in: "2023-02-03T12", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "2023-02-03T12:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},
		{in: "2023-02-03T12:00:00", want: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM)},

		//
		// Dates and Date Ranges
		//

		// MD
		{in: "Feb 3", want: NewRangesWithStartDates(DateForFeb03)},
		{in: "Thu Feb 3", want: NewRangesWithStartDates(DateForFeb03)},
		{in: "Thu 3 Feb", want: NewRangesWithStartDates(DateForFeb03)},
		// DM
		{in: "3 Feb", want: NewRangesWithStartDates(DateForFeb03)},

		// MDY
		{in: "Feb 3 2023", want: NewRangesWithStartDates(DateFor2023Feb03)},
		{in: "February 3 2023", want: NewRangesWithStartDates(DateFor2023Feb03)},
		{in: "February 3, 2023", want: NewRangesWithStartDates(DateFor2023Feb03)},
		{in: "February 3rd, 2023", want: NewRangesWithStartDates(DateFor2023Feb03)},
		{in: "Thu Feb 3, 2023", want: NewRangesWithStartDates(DateFor2023Feb03)},
		// DMY
		{in: "3 Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb03)},
		{in: "3rd Feb 2023", want: NewRangesWithStartDates(DateFor2023Feb03)},
		{in: "3 February, 2023", want: NewRangesWithStartDates(DateFor2023Feb03)},

		// Both
		{in: "02.03", want: NewRangesWithStartDates(DateForFeb03), dateMode: "na"},
		{in: "02.03", want: NewRangesWithStartDates(DateForMar02), dateMode: "rest"},
		{in: "02.03.", want: NewRangesWithStartDates(DateForFeb03), dateMode: "na"},
		{in: "02.03.", want: NewRangesWithStartDates(DateForMar02), dateMode: "rest"},
		{in: "2/3/2023", want: NewRangesWithStartDates(DateFor2023Feb03), dateMode: "na"},
		{in: "2/3/2023", want: NewRangesWithStartDates(DateFor2023Mar02), dateMode: "rest"},

		// MD
		{in: "Feb 3-4", want: NewRangesWithStartEndDates(DateForFeb03, DateForFeb04)},
		{in: "Feb 3 - Mar 2", want: NewRangesWithStartEndDates(DateForFeb03, DateForMar02)},
		{in: "Feb 3 to Mar 2", want: NewRangesWithStartEndDates(DateForFeb03, DateForMar02)},
		// DM
		{in: "3-4 Feb", want: NewRangesWithStartEndDates(DateForFeb03, DateForFeb04)},
		{in: "Thu Feb 3 - Fri Feb 4", want: NewRangesWithStartEndDates(DateForFeb03, DateForFeb04)},
		{in: "3 February - 2 March", want: NewRangesWithStartEndDates(DateForFeb03, DateForMar02)},

		// MDY
		{in: "Feb 3-4 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)},
		{in: "Feb 3 - 4 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)},
		{in: "Feb 3 2023 - Feb 4 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)},
		{in: "Thu Feb 3, 2023 - Fri Feb 4, 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)},
		{in: "February 3 - March 2, 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Mar02)},
		// DMY
		{in: "3-4 Feb 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)},
		{in: "3-4 February, 2023", want: NewRangesWithStartEndDates(DateFor2023Feb03, DateFor2023Feb04)},

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
		// DateTimes and DateTime Ranges
		//
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_PST), timeZone: "PST"},
		{in: "Feb 3 12pm", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM_PST), timeZone: "America/Los_Angeles"},
		{in: "Feb 3 12:00 PM", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},

		// {in: "Feb 3 12:00 PM 12:00", want: NewRangesWithStartDateTimes(DateTimeForFeb03_12PM)},
		// {in: "Feb 3 3:00 PM 15:00", want: NewRangesWithStartDateTimes(DateTimeForFeb03_03PM)},

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
		// for i, tc := range tests[22:23] {
		// for i, tc := range tests[70:71] {
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
