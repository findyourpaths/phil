package parse

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

var TimeFor2023Feb03 = time.Date(2023, time.February, 3, 0, 0, 0, 0, time.UTC)
var TimeFor2023Feb03NineAM = time.Date(2023, time.February, 3, 9, 0, 0, 0, time.UTC)   // time.PST)
var TimeFor2023Feb03Noon = time.Date(2023, time.February, 3, 12, 0, 0, 0, time.UTC)    // time.PST)
var TimeFor2023Feb03ThreePM = time.Date(2023, time.February, 3, 15, 0, 0, 0, time.UTC) // time.PST)
var TimeFor2023Feb06 = time.Date(2023, time.February, 6, 0, 0, 0, 0, time.UTC)
var TimeFor2023Mar04 = time.Date(2023, time.March, 4, 0, 0, 0, 0, time.UTC)
var TimeFor2023Mar07 = time.Date(2023, time.March, 7, 0, 0, 0, 0, time.UTC)

func TestExtractDatetimesRanges(t *testing.T) {
	type test struct {
		mode  string
		input string
		want  *DatetimeRanges
	}

	tests := []test{
		{mode: "", input: "Feb 3 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03)}}},
		{mode: "", input: "February 3 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03)}}},
		{mode: "", input: "February 3, 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03)}}},
		{mode: "", input: "February 3rd, 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03)}}},

		{mode: "", input: "3 Feb 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03)}}},
		{mode: "", input: "3 February, 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03)}}},

		{mode: "us", input: "2/3/2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03)}}},
		{mode: "eu", input: "4/3/2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Mar04)}}},

		{mode: "", input: "Thu Feb 3, 2023 - Sun Feb 6, 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Feb06)}}},
		{mode: "", input: "February 3-6, 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Feb06)}}},
		{mode: "", input: "February 3 - March 4, 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Mar04)}}},

		{mode: "", input: "3-6 Feb 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Feb06)}}},
		{mode: "", input: "3-6 February, 2023", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Feb06)}}},

		{mode: "", input: "Feb 3-6", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Feb06)}}},
		{mode: "", input: "Feb 3 Mar 4", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Mar04)}}},
		{mode: "", input: "Feb 3 - Mar 4", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Mar04)}}},

		{mode: "", input: "Feb 3 2023 12pm PST", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03Noon)}}},
		{mode: "", input: "Feb 3 2023 9am - 12pm", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03Noon)}}},
		{mode: "", input: "Feb 3 2023 9am PST to 3pm PST", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03ThreePM)}}},
		{mode: "", input: "Feb 3 @ 9:00 AM - Feb 3 @ 3:00 PM", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03ThreePM)}}},
		{mode: "", input: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03ThreePM)}}},
		{mode: "", input: "Feb 3 2023 9:00 AM 3:00 PM 09:00 15:00 Google Calendar ICS", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03NineAM, TimeFor2023Feb03ThreePM)}}},
		{mode: "", input: "When 3 Feb 2023 12:00 PM - 3:00 PM", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03Noon, TimeFor2023Feb03ThreePM)}}},
		{mode: "", input: "Thursday, February 3, 2023 12:00 PM 3:00 PM Google Calendar ICS", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03Noon, TimeFor2023Feb03ThreePM)}}},

		{mode: "", input: "Feb 3-6; Mar 4-7", want: &DatetimeRanges{Items: []*DatetimeRange{
			NewDatetimeRange(TimeFor2023Feb03, TimeFor2023Feb06),
			NewDatetimeRange(TimeFor2023Mar04, TimeFor2023Mar07),
		}}},

		// {mode: "", input: "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00", want: &DatetimeRanges{Items: []*DatetimeRange{
		// 	NewDatetimeRange(TimeFor2023Feb03NineAMPST, TimeFor2023Feb03ThreePMPST)}}},

		// "Tue May 9 2023, 07:00pm MDT to 09:00pm MDT"
		// "29 April 2023 10 am - 12 pm"
		// "Fri, Apr 14, 2023 9:00 AM 09:00 Sat, Apr 15, 2023 5:00 PM 17:00"
	}

	for _, tc := range tests[0:1] {
		t.Run(tc.mode, func(t *testing.T) {
			// fmt.Printf(tc.input + "\n")
			got, err := ExtractDatetimeRanges(tc.mode, tc.input)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			for _, dtr := range append(got.Items, tc.want.Items...) {
				dtr.Id = ""
			}
			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected difference:\n%v", diff)
			}
		})
	}
}
