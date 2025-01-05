package parse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
)

// A DateTimeTZs represents a sequence of date and time ranges.
//
// This type DOES include location information.
type DateTimeTZRanges struct {
	Items []*DateTimeTZRange
}

func (rngs DateTimeTZRanges) String() string {
	rs := []string{}
	for _, elt := range rngs.Items {
		rs = append(rs, elt.String())
	}
	return strings.Join(rs, ", ")
}

func AppendDateTimeTZRanges(rngs *DateTimeTZRanges, rng *DateTimeTZRange) *DateTimeTZRanges {
	rngs.Items = append(rngs.Items, rng)
	return rngs
}

func NewRanges(rngs ...*DateTimeTZRange) *DateTimeTZRanges {
	return &DateTimeTZRanges{Items: rngs}
}

func NewRangesWithStartDateTimes(starts ...*DateTimeTZ) *DateTimeTZRanges {
	r := &DateTimeTZRanges{}
	for _, start := range starts {
		r.Items = append(r.Items, &DateTimeTZRange{Start: start})
	}
	return r
}

func NewRangesWithStartDates(starts ...civil.Date) *DateTimeTZRanges {
	r := &DateTimeTZRanges{}
	for _, start := range starts {
		r.Items = append(r.Items, &DateTimeTZRange{Start: &DateTimeTZ{DateTime: civil.DateTime{Date: start}}})
	}
	return r
}

func NewRangesWithStartEndDates(start civil.Date, end civil.Date) *DateTimeTZRanges {
	return &DateTimeTZRanges{Items: []*DateTimeTZRange{NewRangeWithStartEndDates(start, end)}}
}

func NewRangesWithStartEndDateTimes(start *DateTimeTZ, end *DateTimeTZ) *DateTimeTZRanges {
	return &DateTimeTZRanges{Items: []*DateTimeTZRange{NewRangeWithStartEndDateTimes(start, end)}}
}

// A DateTimeTZRange represents a range of dates and times with time zones.
//
// This type DOES include location information.
type DateTimeTZRange struct {
	Start *DateTimeTZ
	End   *DateTimeTZ
}

func (rng DateTimeTZRange) String() string {
	r := rng.Start.String()
	if rng.End != nil {
		r += " - " + rng.End.String()
	}
	return r
}

func NewRangeWithStart(start civil.Date) *DateTimeTZRange {
	return &DateTimeTZRange{Start: &DateTimeTZ{DateTime: civil.DateTime{Date: start}}}
}

func NewRangeWithStartEndDates(start civil.Date, end civil.Date) *DateTimeTZRange {
	return &DateTimeTZRange{
		Start: &DateTimeTZ{DateTime: civil.DateTime{Date: start}},
		End:   &DateTimeTZ{DateTime: civil.DateTime{Date: end}},
	}
}

func NewRangeWithStartEndDateTimes(start *DateTimeTZ, end *DateTimeTZ) *DateTimeTZRange {
	return &DateTimeTZRange{
		Start: start,
		End:   end,
	}
}

// A DateTimeTZ represents a date and time with a time zone.
//
// This type DOES include location information.
type DateTimeTZ struct {
	civil.DateTime
	TimeZone string
}

func (dttz DateTimeTZ) String() string {
	return dttz.DateTime.String()
}

func NewDateTime(date civil.Date, time civil.Time, tz string) *DateTimeTZ {
	// fmt.Printf("NewDateTime(date: %#v, time: %#v, tz: %q)\n", date, time, tz)
	return &DateTimeTZ{DateTime: civil.DateTime{Date: date, Time: time}, TimeZone: tz}
}

func NewDateTimeWithDate(date civil.Date) *DateTimeTZ {
	return &DateTimeTZ{DateTime: civil.DateTime{Date: date}}
}

var weekdaysByNames = map[string]int{
	"sun":       0,
	"sunday":    0,
	"mon":       1,
	"monday":    1,
	"tue":       2,
	"tues":      2,
	"tuesday":   2,
	"wed":       3,
	"weds":      3,
	"wednesday": 3,
	"thu":       4,
	"thus":      4,
	"thursday":  4,
	"fri":       5,
	"friday":    5,
	"saturday":  6,
	"sat":       6,
}

var monthsByNames = map[string]time.Month{
	"jan":       1,
	"january":   1,
	"feb":       2,
	"february":  2,
	"mar":       3,
	"march":     3,
	"apr":       4,
	"april":     4,
	"may":       5,
	"jun":       6,
	"june":      6,
	"jul":       7,
	"july":      7,
	"aug":       8,
	"august":    8,
	"sep":       9,
	"sept":      9,
	"september": 9,
	"oct":       10,
	"october":   10,
	"nov":       11,
	"november":  11,
	"dec":       12,
	"december":  12,
}

var ordinals = map[string]bool{
	"st": true,
	"nd": true,
	"rd": true,
	"th": true,
}

func mustAtoi(str string) int {
	if str == "" {
		return 0
	}
	r, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return r
}

var ambiguousDateMode string

func NewAmbiguousDate(mode string, first string, second string, year string) civil.Date {
	// North American tends to parse dates as month-day-year.
	if mode == "na" {
		return civil.Date{Month: time.Month(mustAtoi(first)), Day: mustAtoi(second), Year: mustAtoi(year)}
	}
	return civil.Date{Day: mustAtoi(first), Month: time.Month(mustAtoi(second)), Year: mustAtoi(year)}
}

func NewDMYDate(dayAny any, monthAny any, yearAny any) civil.Date {
	var day int
	switch dayAny.(type) {
	case int:
		day = dayAny.(int)
	case string:
		day = mustAtoi(dayAny.(string))
	default:
		panic(fmt.Sprintf("can't handle day in unknown format: %#v", dayAny))
	}

	var month time.Month
	switch monthAny.(type) {
	case int:
		month = monthAny.(time.Month)
	case time.Month:
		month = monthAny.(time.Month)
	case string:
		var found bool
		month, found = monthsByNames[strings.ToLower(monthAny.(string))]
		if !found {
			month = time.Month(mustAtoi(monthAny.(string)))
		}
	default:
		panic(fmt.Sprintf("can't handle month %#v of unknown type: %s", monthAny, reflect.TypeOf(monthAny)))
	}

	var year int
	switch yearAny.(type) {
	case int:
		year = yearAny.(int)
	case string:
		year = mustAtoi(yearAny.(string))
	default:
		panic(fmt.Sprintf("can't handle year in unknown format: %#v", yearAny))
	}

	return civil.Date{Day: day, Month: month, Year: year}
}

func NewMDYDate(monthAny any, dayAny any, yearAny any) civil.Date {
	return NewDMYDate(dayAny, monthAny, yearAny)
}

func NewTime(hourAny any, minuteAny any, secondAny any, nsAny any) civil.Time {
	var hour int
	switch hourAny.(type) {
	case int:
		hour = hourAny.(int)
	case string:
		hour = mustAtoi(hourAny.(string))
	default:
		panic(fmt.Sprintf("can't handle hour in unknown format: %#v", hourAny))
	}

	var minute int
	switch minuteAny.(type) {
	case int:
		minute = minuteAny.(int)
	case string:
		minute = mustAtoi(minuteAny.(string))
	default:
		panic(fmt.Sprintf("can't handle minute in unknown format: %#v", minuteAny))
	}

	var second int
	switch secondAny.(type) {
	case int:
		second = secondAny.(int)
	case string:
		second = mustAtoi(secondAny.(string))
	default:
		panic(fmt.Sprintf("can't handle second in unknown format: %#v", secondAny))
	}

	var ns int
	switch nsAny.(type) {
	case int:
		ns = nsAny.(int)
	case string:
		ns = mustAtoi(nsAny.(string))
	default:
		panic(fmt.Sprintf("can't handle ns in unknown format: %#v", nsAny))
	}

	return civil.Time{Hour: hour, Minute: minute, Second: second, Nanosecond: ns}
}
