package datetime

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

var parseDateMode string
var parseDTTZ *DateTimeTZ

// A DateTimeTZs represents a sequence of date and time ranges. It's the
// expected result of parsing a string for datetimes.
//
// This type DOES include location information.
type DateTimeTZRanges struct {
	Items []*DateTimeTZRange
}

func (rngs *DateTimeTZRanges) String() string {
	if rngs == nil {
		return ""
	}
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
	// if dttz.Date != nil {
	// 	r += dttz.Date.String()
	// }
	// r := "T"
	// if dttz.Time != nil {
	// 	r += dttz.Time.String()
	// }
	// if dttz.TimeZone != nil {
	// 	r += dttz.TimeZone.String()
	// }
	// // if tzStr := dttz.TimeZone.String(); tzStr != "" {
	// // 	r += tzStr
	// // }
	return dttz.Date.String() + dttz.Time.String() + dttz.TimeZone.String()
}

func NewDateTimeTZ(date *Date, time *Time, timeZone *TimeZone) *DateTimeTZ {
	// fmt.Printf("NewDateTime(date: %#v, time: %#v, timeZone: %#v)\n", date, time, timeZone)
	if timeZone == nil && parseDTTZ != nil {
		timeZone = parseDTTZ.TimeZone
	}
	return &DateTimeTZ{Date: date, Time: time, TimeZone: timeZone}
}

func NewDateTimeTZWithDate(date *Date, timeZone *TimeZone) *DateTimeTZ {
	if timeZone == nil && parseDTTZ != nil {
		timeZone = parseDTTZ.TimeZone
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
	Weekday int        // Day of the week, starting at 1 for Sunday
}

// String returns the date in RFC3339 full-date format.
func (d *Date) String() string {
	if d == nil {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func NewDate(date *Date) *Date {
	return date
	if date.Year == 0 || date.Month == 0 || date.Day == 0 {
		return date
	}
	t := time.Date(date.Year, date.Month, date.Day, 0, 0, 0, 0, time.UTC)
	tw := weekdaysByNames[strings.ToLower(t.Weekday().String())]
	if date.Weekday == 0 {
		date.Weekday = tw
		return date
	}
	if date.Weekday != tw {
		panic(fmt.Sprintf("extracted weekday of %q doesn't actual match weekday of %q for %4d-%2d-%2d", weekdayNames[tw], weekdayNames[date.Weekday], date.Year, date.Month, date.Day))
	}
	return date
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
func (t *Time) String() string {
	if t == nil {
		return "T"
	}
	s := fmt.Sprintf("T%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
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

func findInt(tunit timeUnit, valAny any) int {
	// fmt.Printf("findInt(tunit: %#v, valAny: %#v)\n", tunit, valAny)
	r := -1
	var ok bool
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
			if valAny.(string) == "" {
				ok = false
				r = 0
				break
			}
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
	// fmt.Printf("in findInt(): r: %d\n", r)

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
var weekdayUnit = timeUnit{name: "weekday", min: 1, max: 7, stringToIntFn: weekdayNameToWeekday}
var hourUnit = timeUnit{name: "hour", min: 0, max: 24, stringToIntFn: hourNameToHour}
var minuteUnit = timeUnit{name: "minute", min: 0, max: 59}
var secondUnit = timeUnit{name: "second", min: 0, max: 59}
var nsUnit = timeUnit{name: "ns", min: 0, max: 999}

// var noTime = &Time{}
// var noTimeZone = &TimeZone{}
// var noYear = 0
// var noMonth = 0
// var noDay = 0
// var noHour = 0
// var noMinute = 0
// var noSecond = 0
// var noNS = 0

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

var weekdayNames = []string{
	"",
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
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

func fixYear(yearAny any, year int) (int, bool) {
	// fmt.Printf("fixYear(yearAny: %#v, year: %d), parseYear: %d\n", yearAny, year, parseYear)
	if parseDTTZ.Date != nil && parseDTTZ.Date.Year != 0 && year == -1 {
		return parseDTTZ.Date.Year, true
	}
	if parseDTTZ.Date != nil && parseDTTZ.Date.Year != 0 && year >= 0 && year <= 99 {
		return 100*(parseDTTZ.Date.Year/100) + year, true
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

func weekdayNameToWeekday(weekdayName string) int {
	weekday, found := weekdaysByNames[strings.ToLower(weekdayName)]
	if !found {
		return 0
	}
	return weekday
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

func NewAmbiguousDate(weekdayAny any, first string, second string, yearAny any) *Date {
	// North American tends to parse dates as month-day-year.
	if parseDateMode == "na" {
		return NewWMDYDate(weekdayAny, first, second, yearAny)
	}
	return NewWDMYDate(weekdayAny, first, second, yearAny)
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
	return NewDate(&Date{Day: day, Month: time.Month(month), Year: year})
}

func NewWDMYDate(weekdayAny any, dayAny any, monthAny any, yearAny any) *Date {
	// fmt.Printf("NewDMYDate(dayAny: %#v, monthAny %#v, yearAny %#v)\n", dayAny, monthAny, yearAny)
	weekday := 0 // findInt(weekdayUnit, weekdayAny)
	day := findInt(dayUnit, dayAny)
	month := findInt(monthUnit, monthAny)
	year := findInt(yearUnit, yearAny)
	return NewDate(&Date{Day: day, Month: time.Month(month), Year: year, Weekday: weekday})
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

func NewWMDYDate(weekdayAny any, monthAny any, dayAny any, yearAny any) *Date {
	return NewWDMYDate(weekdayAny, dayAny, monthAny, yearAny)
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
