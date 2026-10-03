package datetime

import (
	"fmt"
	"strings"
	"time"
)

// newRangeWithDatesAndTimes propagates the range's explicit year before
// combining separately labelled dates and times into resolved endpoints.
func newRangeWithDatesAndTimes(startDate, endDate *Date, startTime, endTime *Time, zone *TimeZone) *DateTimeRange {
	bounds := NewRangeWithStartEndDates(startDate, endDate)
	return NewRange(
		NewDateTime(bounds.Start.Date, startTime, zone),
		NewDateTime(bounds.End.Date, endTime, zone),
	)
}

// newBoundedWeeklyRanges composes a grammar-recognized plural weekday series.
// The date range bounds the series; its first occurrence has only the stated
// clock duration. A singular weekday does not authorize a recurrence reading.
func newBoundedWeeklyRanges(startDate, untilDate *Date, startTime, endTime *Time, zone *TimeZone, weekdayName string) *DateTimeRanges {
	weekday, ok := recurrenceWeekday(weekdayName)
	if !ok || !strings.HasSuffix(strings.ToLower(weekdayName), "days") {
		panic(fmt.Sprintf("semantic error: bounded weekly series requires a plural weekday, got %q", weekdayName))
	}
	bounds := NewRangeWithStartEndDates(startDate, untilDate)
	first := NewRange(
		NewDateTime(bounds.Start.Date, startTime, zone),
		NewDateTime(bounds.Start.Date, endTime, zone),
	)
	return NewRecurringRanges(first, &Recurrence{
		Frequency: FrequencyWeekly,
		Weekdays:  []time.Weekday{weekday},
		Until:     bounds.End.Date,
	})
}
