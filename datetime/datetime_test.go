package datetime

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

type toStringTest struct {
	in   *DateTimeRanges
	want string
}

func TestToString(t *testing.T) {
	tests := []toStringTest{
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM), want: "2023-02-03T12:00:00Z"},
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ET), want: "2023-02-03T12:00:00Z"},
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_ADD0), want: "2023-02-03T12:00:00+00:00"},
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM_SUB0), want: "2023-02-03T12:00:00-00:00"},
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

// type newRangesTest struct {
// 	got  *DateTimeRanges
// 	want *DateTimeRanges
// }

// func TestNewRanges(t *testing.T) {
// 	tests := []newRangesTest{{
// 		got: NewRangesWithStartEndRanges(
// 			NewRange(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET),
// 			NewRange(DateTimeFor2023Feb05_09AM_ET, DateTimeFor2023Feb05_12PM_ET)),
// 		want: NewRanges(
// 			NewRange(DateTimeFor2023Feb03_09AM_ET, DateTimeFor2023Feb03_12PM_ET),
// 			NewRange(DateTimeFor2023Feb04_09AM_ET, DateTimeFor2023Feb04_12PM_ET),
// 			NewRange(DateTimeFor2023Feb05_09AM_ET, DateTimeFor2023Feb05_12PM_ET)),
// 	},
// 	}

// 	failed := 0
// 	for i, tc := range tests {
// 		if !t.Run(fmt.Sprintf("%03d__%s", i, tc.want), testNewRangesFn(t, tc)) {
// 			failed++
// 		}
// 	}

// 	if len(tests) == 0 {
// 		fmt.Println("No tests were run")
// 		return
// 	}

// 	percent := float64(failed) / float64(len(tests)) * 100
// 	fmt.Printf("TestToString: %.2f%% of tests failed (%d/%d)\n", percent, failed, len(tests))
// }

// func testNewRangesFn(t *testing.T, tc newRangesTest) func(*testing.T) {
// 	return func(t *testing.T) {
// 		if diff := cmp.Diff(tc.got, tc.want, protocmp.Transform()); diff != "" {
// 			t.Errorf("unexpected difference:\n%v", diff)
// 		}
// 	}
// }
