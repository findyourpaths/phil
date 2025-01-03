package parse

import (
	"testing"
	"time"

	"cloud.google.com/go/civil"
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

var DateFor2023Feb01 = civil.Date{Year: 2023, Month: time.February, Day: 1}
var DateFor2023Feb02 = civil.Date{Year: 2023, Month: time.February, Day: 2}
var DateFor2023Feb03 = civil.Date{Year: 2023, Month: time.February, Day: 3}
var DateFor2023Feb04 = civil.Date{Year: 2023, Month: time.February, Day: 4}
var DateFor2023Feb05 = civil.Date{Year: 2023, Month: time.February, Day: 5}
var DateFor2023Mar02 = civil.Date{Year: 2023, Month: time.March, Day: 2}
var DateFor2023Mar03 = civil.Date{Year: 2023, Month: time.March, Day: 3}

func TestExtractDatetimesRanges(t *testing.T) {
	type test struct {
		mode  string
		input string
		want  *DateTimeTZRanges
	}

	tests := []test{
		// MD
		{mode: "", input: "Feb 3", want: NewRangesWithStarts(DateForFeb03)},
		{mode: "", input: "Thu Feb 3", want: NewRangesWithStarts(DateForFeb03)},
		{mode: "", input: "Thu 3 Feb", want: NewRangesWithStarts(DateForFeb03)},
		// DM
		{mode: "", input: "3 Feb", want: NewRangesWithStarts(DateForFeb03)},

		// MDY
		{mode: "", input: "Feb 3 2023", want: NewRangesWithStarts(DateFor2023Feb03)},
		{mode: "", input: "February 3 2023", want: NewRangesWithStarts(DateFor2023Feb03)},
		{mode: "", input: "February 3, 2023", want: NewRangesWithStarts(DateFor2023Feb03)},
		{mode: "", input: "February 3rd, 2023", want: NewRangesWithStarts(DateFor2023Feb03)},
		{mode: "", input: "Thu Feb 3, 2023", want: NewRangesWithStarts(DateFor2023Feb03)},
		// DMY
		{mode: "", input: "3 Feb 2023", want: NewRangesWithStarts(DateFor2023Feb03)},
		{mode: "", input: "3rd Feb 2023", want: NewRangesWithStarts(DateFor2023Feb03)},
		{mode: "", input: "3 February, 2023", want: NewRangesWithStarts(DateFor2023Feb03)},

		// Both
		{mode: "na", input: "2/3/2023", want: NewRangesWithStarts(DateFor2023Feb03)},
		{mode: "", input: "2/3/2023", want: NewRangesWithStarts(DateFor2023Mar02)},

		// MD
		{mode: "", input: "Feb 3-4", want: NewRangesWithStartEnd(DateForFeb03, DateForFeb04)},
		{mode: "", input: "Feb 3 - Mar 2", want: NewRangesWithStartEnd(DateForFeb03, DateForMar02)},
		// DM
		{mode: "", input: "3-4 Feb", want: NewRangesWithStartEnd(DateForFeb03, DateForFeb04)},
		{mode: "", input: "Thu Feb 3 - Fri Feb 4", want: NewRangesWithStartEnd(DateForFeb03, DateForFeb04)},
		{mode: "", input: "3 February - 2 March", want: NewRangesWithStartEnd(DateForFeb03, DateForMar02)},

		// MDY
		{mode: "", input: "Feb 3-4 2023", want: NewRangesWithStartEnd(DateFor2023Feb03, DateFor2023Feb04)},
		{mode: "", input: "Feb 3 - 4 2023", want: NewRangesWithStartEnd(DateFor2023Feb03, DateFor2023Feb04)},
		{mode: "", input: "Feb 3 2023 - Feb 4 2023", want: NewRangesWithStartEnd(DateFor2023Feb03, DateFor2023Feb04)},
		{mode: "", input: "Thu Feb 3, 2023 - Fri Feb 4, 2023", want: NewRangesWithStartEnd(DateFor2023Feb03, DateFor2023Feb04)},
		{mode: "", input: "February 3 - March 2, 2023", want: NewRangesWithStartEnd(DateFor2023Feb03, DateFor2023Mar02)},
		// DMY
		{mode: "", input: "3-4 Feb 2023", want: NewRangesWithStartEnd(DateFor2023Feb03, DateFor2023Feb04)},
		{mode: "", input: "3-4 February, 2023", want: NewRangesWithStartEnd(DateFor2023Feb03, DateFor2023Feb04)},

		// MD
		{mode: "", input: "Feb 1, 2", want: NewRangesWithStarts(DateForFeb01, DateForFeb02)},
		{mode: "", input: "Feb 1, 2, 3", want: NewRangesWithStarts(DateForFeb01, DateForFeb02, DateForFeb03)},
		{mode: "", input: "Feb 1, 2, 3, 4", want: NewRangesWithStarts(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04)},
		{mode: "", input: "Feb 1, 2, 3, 4, 5", want: NewRangesWithStarts(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb05)},
		// DM
		{mode: "", input: "1, 2 Feb", want: NewRangesWithStarts(DateForFeb01, DateForFeb02)},
		{mode: "", input: "1, 2, 3 Feb", want: NewRangesWithStarts(DateForFeb01, DateForFeb02, DateForFeb03)},
		{mode: "", input: "1, 2, 3, 4 Feb", want: NewRangesWithStarts(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04)},
		{mode: "", input: "1, 2, 3, 4, 5 Feb", want: NewRangesWithStarts(DateForFeb01, DateForFeb02, DateForFeb03, DateForFeb04, DateForFeb05)},

		// MDY
		{mode: "", input: "Feb 1, 2 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02)},
		{mode: "", input: "Feb 1, 2, 3 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03)},
		{mode: "", input: "Feb 1, 2, 3, 4 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04)},
		{mode: "", input: "Feb 1, 2, 3, 4, 5 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04, DateFor2023Feb05)},
		// DMY
		{mode: "", input: "1, 2 Feb 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02)},
		{mode: "", input: "1, 2, 3 Feb 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03)},
		{mode: "", input: "1, 2, 3, 4 Feb 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04)},
		{mode: "", input: "1, 2, 3, 4, 5 Feb 2023", want: NewRangesWithStarts(DateFor2023Feb01, DateFor2023Feb02, DateFor2023Feb03, DateFor2023Feb04, DateFor2023Feb05)},

		// MD
		{mode: "", input: "Feb 3 Mar 2", want: NewRangesWithStarts(DateForFeb03, DateForMar02)},
		// MDY
		{mode: "", input: "Feb 3 Mar 2 2023", want: NewRangesWithStarts(DateFor2023Feb03, DateFor2023Mar02)},

		// MD
		{mode: "", input: "Feb 1-2, 3-4", want: NewRanges(NewRangeWithStartEnd(DateForFeb01, DateForFeb02), NewRangeWithStartEnd(DateForFeb03, DateForFeb04))},
		{mode: "", input: "Feb 1-2, 3-4 2023", want: NewRanges(NewRangeWithStartEnd(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEnd(DateFor2023Feb03, DateFor2023Feb04))},

		// MD
		{mode: "", input: "Feb 1-2; Mar 2-3", want: NewRanges(NewRangeWithStartEnd(DateForFeb01, DateForFeb02), NewRangeWithStartEnd(DateForMar02, DateForMar03))},
		// DM
		{mode: "", input: "1-2 Feb; 2-3 Mar", want: NewRanges(NewRangeWithStartEnd(DateForFeb01, DateForFeb02), NewRangeWithStartEnd(DateForMar02, DateForMar03))},

		// MDY
		{mode: "", input: "Feb 1-2; Mar 2-3 2023", want: NewRanges(NewRangeWithStartEnd(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEnd(DateFor2023Mar02, DateFor2023Mar03))},
		// DMY
		{mode: "", input: "1-2 Feb; 2-3 Mar 2023", want: NewRanges(NewRangeWithStartEnd(DateFor2023Feb01, DateFor2023Feb02), NewRangeWithStartEnd(DateFor2023Mar02, DateFor2023Mar03))},

		// {mode: "", input: "Feb 3 2023 12pm PST", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03Noon)}}},
		// {mode: "", input: "Feb 3 2023 9am - 12pm", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03Noon)}}},
		// {mode: "", input: "Feb 3 2023 9am PST to 3pm PST", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03ThreePM)}}},
		// {mode: "", input: "Feb 3 @ 9:00 AM - Feb 3 @ 3:00 PM", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03ThreePM)}}},
		// {mode: "", input: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03ThreePM)}}},
		// {mode: "", input: "Feb 3 2023 9:00 AM 3:00 PM 09:00 15:00 Google Calendar ICS", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03ThreePM)}}},
		// {mode: "", input: "When 3 Feb 2023 12:00 PM - 3:00 PM", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03Noon, TimeFor2023Feb03ThreePM)}}},
		// {mode: "", input: "Thursday, February 3, 2023 12:00 PM 3:00 PM Google Calendar ICS", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03Noon, TimeFor2023Feb03ThreePM)}}},

		// {mode: "", input: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03NineAMPST, TimeFor2023Feb03ThreePMPST)}}},

		// "Tue May 9 2023, 07:00pm MDT to 09:00pm MDT"
		// "29 April 2023 10 am - 12 pm"
		// "Fri, Apr 14, 2023 9:00 AM 09:00 Sat, Apr 15, 2023 5:00 PM 17:00"
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ExtractDateTimeTZRanges(tc.mode, tc.input)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected difference:\n%v", diff)
			}
		})
	}
}
