package datetime

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

var parseDateMode string
var parseYear int
var parseTimeZone *TimeZone
var parseTimeZoneAbbrev string

// A DateTimeTZs represents a sequence of date and time ranges. It's the
// expected result of parsing a string for datetimes.
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

func NewRangesWithStartDates(starts ...*Date) *DateTimeTZRanges {
	r := &DateTimeTZRanges{}
	for _, start := range starts {
		r.Items = append(r.Items, &DateTimeTZRange{Start: NewDateTimeTZWithDate(start, nil)})
	}
	return r
}

func NewRangesWithStartEndDates(start *Date, end *Date) *DateTimeTZRanges {
	return &DateTimeTZRanges{Items: []*DateTimeTZRange{NewRangeWithStartEndDates(start, end)}}
}

func NewRangesWithStartEndDateTimes(start *DateTimeTZ, end *DateTimeTZ) *DateTimeTZRanges {
	return &DateTimeTZRanges{Items: []*DateTimeTZRange{NewRangeWithStartEndDateTimes(start, end)}}
}

func HasStartMonthAndDay(rngs *DateTimeTZRanges) bool {
	return rngs != nil &&
		len(rngs.Items) > 0 &&
		rngs.Items[0].Start != nil &&
		rngs.Items[0].Start.Date.Month > 0 &&
		rngs.Items[0].Start.Date.Day > 0
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

func NewRangeWithStart(start *Date) *DateTimeTZRange {
	return &DateTimeTZRange{Start: &DateTimeTZ{Date: start}}
}

func NewRangeWithStartEndDates(start *Date, end *Date) *DateTimeTZRange {
	return &DateTimeTZRange{
		Start: &DateTimeTZ{Date: start},
		End:   &DateTimeTZ{Date: end},
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
	Date     *Date
	Time     *Time
	TimeZone *TimeZone
}

func (dttz *DateTimeTZ) String() string {
	if dttz == nil {
		return ""
	}
	r := dttz.Date.String() + "T" + dttz.Time.String()
	if tzStr := dttz.TimeZone.String(); tzStr != "" {
		r += tzStr
	}
	return r
}

func NewDateTimeTZ(date *Date, time *Time, timeZone *TimeZone) *DateTimeTZ {
	// fmt.Printf("NewDateTime(date: %#v, time: %#v, timeZone: %#v)\n", date, time, timeZone)
	if timeZone == nil {
		timeZone = parseTimeZone
	}
	return &DateTimeTZ{Date: date, Time: time, TimeZone: timeZone}
}

func NewDateTimeTZWithDate(date *Date, timeZone *TimeZone) *DateTimeTZ {
	if timeZone == nil {
		timeZone = parseTimeZone
	}
	return &DateTimeTZ{Date: date, TimeZone: timeZone}
}

// A Date represents a date (year, month, day, weekday).
//
// This type does not include location information, and therefore does not
// describe a unique 24-hour timespan.
type Date struct {
	Year    int        // Year (e.g., 2014).
	Month   time.Month // Month of the year (January = 1, ...).
	Day     int        // Day of the month, starting at 1.
	Weekday int        // Weekday, starting at 1
}

// String returns the date in RFC3339 full-date format.
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// A Time represents a time with nanosecond precision.
//
// This type does not include location information, and therefore does not
// describe a unique moment in time.
//
// This type exists to represent the TIME type in storage-based APIs like BigQuery.
// Most operations on Times are unlikely to be meaningful. Prefer the DateTime type.
type Time struct {
	Hour       int // The hour of the day in 24-hour format; range [0-23]
	Minute     int // The minute of the hour; range [0-59]
	Second     int // The second of the minute; range [0-59]
	Nanosecond int // The nanosecond of the second; range [0-999999999]
}

// String returns the date in the format described in ParseTime. If Nanoseconds
// is zero, no fractional part will be generated. Otherwise, the result will
// end with a fractional part consisting of a decimal point and nine digits.
func (t Time) String() string {
	s := fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
	if t.Nanosecond == 0 {
		return s
	}
	return s + fmt.Sprintf(".%09d", t.Nanosecond)
}

type TimeZone struct {
	Name   string
	Abbrev string
	Offset string
}

func (tz *TimeZone) String() string {
	if tz == nil {
		return "Z"
	}
	if tz.Offset != "" {
		return tz.Offset
	}
	if tz.Name != "" {
		if tz, _ := timezoneTZ.GetTzInfo(tz.Name); tz != nil {
			return tz.StandardOffsetHHMM()
		}
	}
	if tz.Abbrev != "" {
		tzs, _ := timezoneTZ.GetTzAbbreviationInfo(tz.Abbrev)
		if len(tzs) == 0 {
			return ""
		}
		if len(tzs) > 1 {
			slog.Debug("got multiple time zones", "tz", fmt.Sprintf("%#v", tz), "tzs", tzs)
		}
		return tzs[0].OffsetHHMM()
	}
	return ""
}

func NewTimeZone(nameAny any, abbrevAny any, offsetAny any) *TimeZone {
	name, nok := nameAny.(string)
	abbrev, aok := abbrevAny.(string)
	offset, ook := offsetAny.(string)
	if !nok && !aok && !ook {
		slog.Debug("didn't make new time zone because of failure to parse name or abbrev or offset", "nameAny", nameAny, "abbrevAny", abbrevAny, "offsetAny", offsetAny)
		return nil
	}
	if name == "" && abbrev == "" && offset == "" {
		return nil
	}
	return &TimeZone{Name: name, Abbrev: abbrev, Offset: offset}
}

var weekdaysByNames = map[string]int{
	"su":        1,
	"sun":       1,
	"sunday":    1,
	"mo":        2,
	"mon":       2,
	"monday":    2,
	"tu":        3,
	"tue":       3,
	"tues":      3,
	"tuesday":   3,
	"we":        4,
	"wed":       4,
	"weds":      4,
	"wednesday": 4,
	//	"th":        4, recognize this separately because it also appears in "4th"
	"thu":      5,
	"thus":     5,
	"thursday": 5,
	"fr":       6,
	"fri":      6,
	"friday":   6,
	"sa":       7,
	"sat":      7,
	"saturday": 7,
}

var monthsByNames = map[string]int{
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

type timeUnit struct {
	name          string
	min           int
	max           int
	emptyVal      any
	fixFn         func(any, int) (int, bool)
	stringToIntFn func(string) int
}

var yearUnit = timeUnit{name: "year", fixFn: fixYear}
var monthUnit = timeUnit{name: "month", min: 1, max: 12, stringToIntFn: monthNameToMonth}
var dayUnit = timeUnit{name: "day", min: 1, max: 31}
var weekdayUnit = timeUnit{name: "weekday", min: 1, max: 7}
var hourUnit = timeUnit{name: "hour", min: 0, max: 24, stringToIntFn: hourNameToHour}
var minuteUnit = timeUnit{name: "minute", min: 0, max: 59}
var secondUnit = timeUnit{name: "second", min: 0, max: 59}
var nsUnit = timeUnit{name: "ns", min: 0, max: 999}

var noTime = &Time{}
var noTimeZone = &TimeZone{}
var noYear = 0
var noMonth = 0
var noDay = 0
var noHour = 0
var noMinute = 0
var noSecond = 0
var noNS = 0

// func findFloat(valAny any) (float32, bool) {
// 	if valAny == nil {
// 		return 0.0, true
// 	}
// 	var r float32
// 	switch valAny.(type) {
// 	case float32:
// 		rFloat, ok := valAny.(float32)
// 		if ok {
// 			r = rFloat
// 		}
// 	case string:
// 		rFloat, err := strconv.ParseFloat(valAny.(string), 32)
// 		if err == nil {
// 			r = float32(rFloat)
// 		}
// 	default:
// 		// panic(fmt.Sprintln("failed to find float", "valAny", valAny))
// 		return r, false
// 	}

// 	return r, true
// }

func findInt(tunit timeUnit, valAny any) int {
	r := -1
	if valAny != nil {
		switch valAny.(type) {
		case int:
			rInt, ok := valAny.(int)
			if ok {
				r = rInt
			}
		case *int:
			rPtr, ok := valAny.(*int)
			if ok {
				r = *rPtr
			}
		case string:
			rInt, err := strconv.Atoi(valAny.(string))
			if err == nil {
				r = rInt
			} else {
				if tunit.stringToIntFn != nil {
					r = tunit.stringToIntFn(valAny.(string))
				}
			}
		}
	}

	var ok bool
	if tunit.fixFn != nil {
		r, ok = tunit.fixFn(valAny, r)
		// fmt.Println("r", r, "ok", ok)
	} else if valAny == nil {
		return 0
	}

	// debugf("in findInt(), r: %d, ok: %t\n", r, ok)
	if !ok && (r < tunit.min || r > tunit.max) {
		panic(fmt.Sprintln("found int but failed bounds check", "tunit", tunit, "valAny", valAny))
	}
	// debugf("in findInt(), returning: %d\n", r)
	return r
}

func fixYear(yearAny any, year int) (int, bool) {
	// fmt.Printf("fixYear(yearAny: %#v, year: %d), parseYear: %d\n", yearAny, year, parseYear)
	if parseYear != 0 && year == -1 {
		return parseYear, true
	}
	if parseYear != 0 && year >= 0 && year <= 99 {
		return 100*(parseYear/100) + year, true
	}
	if yearAny == nil {
		return 0, true
	}
	return year, (year >= 1700 && year <= 2100)
}

func monthNameToMonth(monthName string) int {
	month, found := monthsByNames[strings.ToLower(monthName)]
	if !found {
		return -1
	}
	return month
}

func hourNameToHour(hourName string) int {
	l := strings.ToLower(hourName)
	if l == "noon" {
		return 12
	}
	if l == "midnight" {
		return 0
	}
	return -1
}

// func mustAtoi(str string) *int {
// 	if str == "" {
// 		return nil
// 	}
// 	r, err := strconv.Atoi(str)
// 	if err != nil {
// 		return nil
// 	}
// 	return &r
// }

func NewAmbiguousDate(first string, second string, yearAny any) *Date {
	// North American tends to parse dates as month-day-year.
	if parseDateMode == "na" {
		return NewMDYDate(first, second, yearAny)
	}
	return NewDMYDate(first, second, yearAny)
}

func NewDsMYDates(daysAny []string, monthAny any, yearAny any) []*Date {
	rs := []*Date{}
	for _, dayAny := range daysAny {
		rs = append(rs, NewDMYDate(dayAny, monthAny, yearAny))
	}
	return rs
}

func NewDMYDate(dayAny any, monthAny any, yearAny any) *Date {
	// fmt.Printf("NewDMYDate(dayAny: %#v, monthAny %#v, yearAny %#v)\n", dayAny, monthAny, yearAny)
	day := findInt(dayUnit, dayAny)
	month := findInt(monthUnit, monthAny)
	year := findInt(yearUnit, yearAny)
	return &Date{Day: day, Month: time.Month(month), Year: year}
}

func NewDMYWDate(dayAny any, monthAny any, yearAny any, weekdayAny any) *Date {
	// fmt.Printf("NewDMYDate(dayAny: %#v, monthAny %#v, yearAny %#v)\n", dayAny, monthAny, yearAny)
	day := findInt(dayUnit, dayAny)
	month := findInt(monthUnit, monthAny)
	year := findInt(yearUnit, yearAny)
	weekday := findInt(weekdayUnit, weekdayAny)
	return &Date{Day: day, Month: time.Month(month), Year: year, Weekday: weekday}
}

func NewMDsYDates(monthAny any, daysAny []string, yearAny any) []*Date {
	rs := []*Date{}
	for _, dayAny := range daysAny {
		rs = append(rs, NewMDYDate(monthAny, dayAny, yearAny))
	}
	return rs
}

func NewMDYDate(monthAny any, dayAny any, yearAny any) *Date {
	return NewDMYDate(dayAny, monthAny, yearAny)
}

func NewMDYWDate(monthAny any, dayAny any, yearAny any, weekdayAny any) *Date {
	return NewDMYWDate(dayAny, monthAny, yearAny, weekdayAny)
}

func NewAMTime(hourAny any, minuteAny any, secondAny any, nsAny any) *Time {
	r := NewTime(hourAny, minuteAny, secondAny, nsAny)
	if r.Hour > 12 {
		panic(fmt.Sprintln("found hour but failed AM bounds check", "r.Hour", r.Hour))
	}
	r.Hour = r.Hour % 12
	return r
}

func NewPMTime(hourAny any, minuteAny any, secondAny any, nsAny any) *Time {
	r := NewTime(hourAny, minuteAny, secondAny, nsAny)
	if r.Hour > 12 {
		panic(fmt.Sprintln("found hour but failed PM bounds check", "r.Hour", r.Hour))
	}
	r.Hour = (r.Hour % 12) + 12
	return r
}

func NewTime(hourAny any, minuteAny any, secondAny any, nsAny any) *Time {
	hour := findInt(hourUnit, hourAny)
	minute := findInt(minuteUnit, minuteAny)
	second := findInt(secondUnit, secondAny)
	ns := findInt(nsUnit, nsAny)
	return &Time{Hour: hour, Minute: minute, Second: second, Nanosecond: ns}
}
