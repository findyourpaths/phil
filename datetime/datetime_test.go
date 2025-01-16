package datetime

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

type toStringTest struct {
	in   *DateTimeTZRanges
	want string
}

func TestToString(t *testing.T) {
	tests := []toStringTest{
		{in: NewRangesWithStartDateTimes(DateTimeFor2023Feb03_12PM), want: "2023-02-03T12:00:00Z"},
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
