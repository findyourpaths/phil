package parse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
)

var parseDateMode string
var parseYear int
var parseTimeZone string

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
		r.Items = append(r.Items, &DateTimeTZRange{Start: NewDateTimeWithDate(start)})
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
	r := dttz.DateTime.String()
	if dttz.TimeZone != "" {
		tzs, err := singletonTZ.GetTzAbbreviationInfo(dttz.TimeZone)
		if err != nil {
			panic(fmt.Sprintf("got error looking up time zone for %#v: %v", dttz, err))
		}
		if len(tzs) > 1 {
			panic(fmt.Sprintf("got multiple time zones for abbreviation: %q for %#v", dttz.TimeZone, dttz))
		}
		r += tzs[0].OffsetHHMM()
	}
	return r
}

func NewDateTime(date civil.Date, time civil.Time, tz string) *DateTimeTZ {
	if tz == "" {
		tz = parseTimeZone
	}
	return &DateTimeTZ{DateTime: civil.DateTime{Date: date, Time: time}, TimeZone: tz}
}

func NewDateTimeWithDate(date civil.Date) *DateTimeTZ {
	return &DateTimeTZ{DateTime: civil.DateTime{Date: date}, TimeZone: parseTimeZone}
}

var weekdaysByNames = map[string]int{
	"su":        0,
	"sun":       0,
	"sunday":    0,
	"mo":        1,
	"mon":       1,
	"monday":    1,
	"tu":        2,
	"tue":       2,
	"tues":      2,
	"tuesday":   2,
	"we":        3,
	"wed":       3,
	"weds":      3,
	"wednesday": 3,
	//	"th":        4, recognize this separately because it also appears in "4th"
	"thu":      4,
	"thus":     4,
	"thursday": 4,
	"fr":       5,
	"fri":      5,
	"friday":   5,
	"sa":       6,
	"sat":      6,
	"saturday": 6,
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
	// "th": true, recognize this separately because it also shortens Thursday
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

func NewAmbiguousDate(first string, second string, yearAny any) civil.Date {
	year := findInt("year", yearAny)
	if year == 0 {
		year = parseYear
	}

	// North American tends to parse dates as month-day-year.
	if parseDateMode == "na" {
		return civil.Date{Month: time.Month(mustAtoi(first)), Day: mustAtoi(second), Year: year}
	}
	return civil.Date{Day: mustAtoi(first), Month: time.Month(mustAtoi(second)), Year: year}
}

func NewDsMYDates(daysAny []string, monthAny any, yearAny any) []civil.Date {
	rs := []civil.Date{}
	for _, dayAny := range daysAny {
		rs = append(rs, NewDMYDate(dayAny, monthAny, yearAny))
	}
	return rs
}

func NewDMYDate(dayAny any, monthAny any, yearAny any) civil.Date {
	day := findInt("day", dayAny)

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

	year := findInt("year", yearAny)
	if year == 0 {
		year = parseYear
	}

	return civil.Date{Day: day, Month: month, Year: year}
}

func NewMDsYDates(monthAny any, daysAny []string, yearAny any) []civil.Date {
	rs := []civil.Date{}
	for _, dayAny := range daysAny {
		rs = append(rs, NewMDYDate(monthAny, dayAny, yearAny))
	}
	return rs
}

func NewMDYDate(monthAny any, dayAny any, yearAny any) civil.Date {
	return NewDMYDate(dayAny, monthAny, yearAny)
}

func NewTime(hourAny any, minuteAny any, secondAny any, nsAny any) civil.Time {
	hour := findInt("hour", hourAny)
	minute := findInt("minute", minuteAny)
	second := findInt("second", secondAny)
	ns := findInt("ns", nsAny)
	return civil.Time{Hour: hour, Minute: minute, Second: second, Nanosecond: ns}
}

func findInt(name string, valAny any) int {
	switch valAny.(type) {
	case int:
		return valAny.(int)
	case string:
		return mustAtoi(valAny.(string))
	default:
		panic(fmt.Sprintf("can't handle %s in unknown format: %#v", name, valAny))
	}
}
