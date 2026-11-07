package datetime

import (
	"strconv"
	"time"
)

// newRangeWithHourDuration preserves the parsed start's civil/zone semantics.
// Duration is elapsed time, including when a named zone crosses a DST change.
func newRangeWithHourDuration(start *DateTime, text string) *DateTimeRange {
	hours, err := strconv.ParseInt(text, 10, 64)
	if err != nil || hours <= 0 || hours > int64((1<<63-1)/time.Hour) || start == nil || start.Date == nil || start.Time == nil || start.Date.Year == 0 || start.Date.Month == 0 || start.Date.Day == 0 {
		panic("semantic error: hour duration requires a complete date, clock and positive bounded hour count")
	}
	end := NewDateTimeForTime(start.ToTime().Add(time.Duration(hours) * time.Hour))
	end.TimeZone = start.TimeZone
	return NewRange(start, end)
}
