package datetime

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

var minimumDateTime *DateTime

var parseDateMode string

var DateModeUnknown = ""

// North American tends to parse dates as month-day-year.
var DateModeNorthAmerican = "na"
var DateModeRest = "rest"

var DateModeNorthAmericanCountryCodes = map[string]bool{"CA": true, "US": true}

// A DateTimes represents a sequence of date and time ranges. It's the
// expected result of parsing a string for datetimes.
//
// This type DOES include location information.
type DateTimeRanges struct {
	Items []*DateTimeRange
}

func (rngs *DateTimeRanges) String() string {
	if rngs == nil {
		return ""
	}
	rs := []string{}
	for _, elt := range rngs.Items {
		rs = append(rs, elt.String())
	}
	return strings.Join(rs, ", ")
}

func AppendDateTimeRanges(rngs *DateTimeRanges, rng *DateTimeRange) *DateTimeRanges {
	rngs.Items = append(rngs.Items, rng)
	return rngs
}

func NewRanges(rngs ...*DateTimeRange) *DateTimeRanges {
	return &DateTimeRanges{Items: rngs}
}

func NewRangesWithStartDateTimes(starts ...*DateTime) *DateTimeRanges {
	r := &DateTimeRanges{}
	for _, start := range starts {
		r.Items = append(r.Items, NewRangeWithStart(start))
	}
	return r
}

func NewRangesWithStartDates(starts ...*Date) *DateTimeRanges {
	r := &DateTimeRanges{}
	for _, start := range starts {
		r.Items = append(r.Items, NewRangeWithStartDate(start))
	}
	return r
}

// NewRangesWithStartEndRanges fills in the days between start and end. For
// example a start of "Feb 1" and end of "Feb 4" is filled in with "Feb 2" and
// "Feb 3".
func NewRangesWithStartEndRanges(start *DateTimeRange, end *DateTimeRange) *DateTimeRanges {
	if start.Start.Date.String() != start.End.Date.String() {
		panic(fmt.Sprintf("semantic error: start range must begin and end on same date: start: %q, end: %q", start.Start, start.End))
	}
	if end.Start.Date.String() != end.End.Date.String() {
		panic(fmt.Sprintf("semantic error: end range must begin and end on same date: start: %q, end: %q", end.Start, end.End))
	}
	if start.Start.Time.String() != end.Start.Time.String() {
		panic(fmt.Sprintf("semantic error: start and end ranges must start with the same time: start: %q, end: %q", start, end))
	}
	if start.End.Time.String() != end.End.Time.String() {
		panic(fmt.Sprintf("semantic error: start and end ranges must end with the same time: start: %q, end: %q", start, end))
	}

	r := &DateTimeRanges{}
	r.Items = append(r.Items, start)
	r.Items = append(r.Items, end)
	return r
}

func NewRangesWithStartEndDates(start *Date, end *Date) *DateTimeRanges {
	return &DateTimeRanges{Items: []*DateTimeRange{NewRangeWithStartEndDates(start, end)}}
}

func NewRangesWithStartEndDateTimes(start *DateTime, end *DateTime) *DateTimeRanges {
	return &DateTimeRanges{Items: []*DateTimeRange{NewRange(start, end)}}
}

func HasStartMonthAndDay(rngs *DateTimeRanges) bool {
	return rngs != nil &&
		len(rngs.Items) > 0 &&
		rngs.Items[0].Start != nil &&
		rngs.Items[0].Start.Date.Month > 0 &&
		rngs.Items[0].Start.Date.Day > 0
}

// A DateTimeRange represents a range of dates and times with time zones.
//
// This type DOES include location information.
type DateTimeRange struct {
	Start *DateTime
	End   *DateTime
	// Frequency Frequency
}

// type Frequency int

// const (
// 	UNSPECIFIED_FREQUENCY Frequency = iota
// 	DAILY
// 	WEEKLY
// )

func (rng DateTimeRange) String() string {
	r := rng.Start.String()
	if rng.End != nil {
		r += " - " + rng.End.String()
	}
	return r
}

func (rng *DateTimeRange) AddDate(years int, months int, days int) *DateTimeRange {
	// return &DateTimeRange {
	// 	Start: rng.Start.Copy().AddDate(years, months, days)
	// End   *DateTime

	// r := rng.Copy()
	// r.Start =

	// r := rng.Start.String()
	// if rng.End != nil {
	// 	r += " - " + rng.End.String()
	// }
	return rng
}

// func NewDailyRange() *DateTimeRange {
// 	return &DateTimeRange{Frequency: DAILY}
// }

// func NewWeeklyRange() *DateTimeRange {
// 	return &DateTimeRange{Frequency: WEEKLY}
// }

func NewRangeWithStartDate(startD *Date) *DateTimeRange {
	return NewRangeWithStartEndDates(startD, nil)
}

func NewRangeWithStartEndDates(startD *Date, endD *Date) *DateTimeRange {
	var startDT *DateTime
	if startD != nil {
		startDT = &DateTime{Date: startD}
	}
	var endDT *DateTime
	if endD != nil {
		endDT = &DateTime{Date: endD}
	}
	return NewRange(startDT, endDT)
}

func NewRangeWithStart(startDT *DateTime) *DateTimeRange {
	return NewRange(startDT, nil)
}

func NewRange(start *DateTime, end *DateTime) *DateTimeRange {
	if start != nil && end != nil {
		// fmt.Printf("checking TimeZones\n")
		// Check that start and end TimeZones don't conflict, and if one is missing, copy from the other.
		if start.TimeZone != nil && end.TimeZone != nil && *(start.TimeZone) != *(end.TimeZone) {
			panic(fmt.Sprintf("semantic error: start TimeZone %#v is different from end TimeZone: %#v\n", start.TimeZone, end.TimeZone))
		}
		if start.TimeZone == nil && end.TimeZone != nil {
			start.TimeZone = end.TimeZone
		}
		if start.TimeZone != nil && end.TimeZone == nil {
			end.TimeZone = start.TimeZone
		}

		// fmt.Printf("checking Times\n")
		if start.Time != nil && end.Time != nil {
			if start.Date == nil && end.Date != nil {
				start.Date = end.Date
			} else if start.Date != nil && end.Date == nil {
				end.Date = start.Date
			}
		}

		// fmt.Printf("checking Dates\n")
		if start.Date != nil && end.Date != nil {
			// fmt.Printf("checking Days\n")
			// If both have Days but one Month is missing, copy from the other.
			if start.Date.Day != 0 && end.Date.Day != 0 {
				if start.Date.Month == 0 && end.Date.Month != 0 {
					start.Date.Month = end.Date.Month
				} else if start.Date.Month != 0 && end.Date.Month == 0 {
					end.Date.Month = start.Date.Month
				}
			}

			// If both have Months but one Year is missing, copy from the other.
			// fmt.Printf("checking Months: %d, %d\n", start.Date.Month, end.Date.Month)
			// fmt.Printf("before checking Months, Years: %d, %d\n", start.Date.Year, end.Date.Year)
			if start.Date.Month != 0 && end.Date.Month != 0 {
				if start.Date.Year == 0 && end.Date.Year != 0 {
					start.Date.Year = end.Date.Year
				} else if start.Date.Year != 0 && end.Date.Year == 0 {
					end.Date.Year = start.Date.Year
				}
			}
			// fmt.Printf("after checking Months, Years: %d, %d\n", start.Date.Year, end.Date.Year)
		}
	}

	return &DateTimeRange{
		Start: start,
		End:   end,
	}
}

// A DateTime represents a date and time with a time zone.
//
// This type DOES include location information.
type DateTime struct {
	Date     *Date
	Time     *Time
	TimeZone *TimeZone
}

func (dt *DateTime) Location() *time.Location {
	if dt.TimeZone == nil {
		fmt.Printf("warning: no TimeZone found in DateTime, returning nil Location")
		return nil
	}

	if dt.TimeZone.Name != "" {
		r := locationForName(dt.TimeZone.Name)
		if r != nil {
			return r
		}
	}

	if dt.TimeZone.Offset != "" {
		r := locationForOffset(dt, dt.TimeZone.Offset)
		if r != nil {
			return r
		}
	}

	if dt.TimeZone.Abbreviation != "" {
		r := locationForAbbreviation(dt, dt.TimeZone.Abbreviation)
		if r != nil {
			return r
		}
	}

	fmt.Printf("warning: no time Location found for TimeZone: %#v, returning nil\n", dt.TimeZone)
	return nil
}

func (dt *DateTime) String() string {
	if dt == nil {
		return ""
	}
	if dt.Date != nil && dt.Time != nil && dt.TimeZone != nil {
		if t := dt.ToTime(); t != nil {
			return t.Format(time.RFC3339)
		}
	}
	return dt.Date.String() + dt.Time.String() + dt.TimeZone.String()
}

func (dt *DateTime) ToTime() *time.Time {
	// fmt.Printf("dt.ToTime(), dt: %#v\n", dt)
	if dt == nil {
		return nil
	}

	t := &Time{}
	if dt.Time != nil {
		t = dt.Time
	}

	r := time.Date(dt.Date.Year, dt.Date.Month, dt.Date.Day, t.Hour, t.Minute, t.Second, t.Nanosecond, dt.Location())
	// fmt.Printf("in dt.ToTime(), r: %#v\n", r)
	return &r
}

func NewDateTimeWithDate(date *Date) *DateTime {
	return NewDateTime(date, nil, nil)
}

func NewDateTimeForNow() *DateTime {
	return NewDateTimeForTime(time.Now())
}

func NewDateTimeForRFC1123Z(tStr string) *DateTime {
	t, err := time.Parse(time.RFC1123Z, tStr)
	if err != nil {
		panic(fmt.Sprintf("error parsing time string (%q): %v", tStr, err))
	}
	return NewDateTimeForTime(t)
}

func NewDateTimeForRFC3339(tStr string) *DateTime {
	t, err := time.Parse(time.RFC3339, tStr)
	if err != nil {
		panic(fmt.Sprintf("error parsing time string (%q): %v", tStr, err))
	}
	return NewDateTimeForTime(t)
}

func NewDateTimeForTime(t time.Time) *DateTime {
	abbrev, off := t.Zone()
	return NewDateTimeWithTimeAndTimeZone(t, abbrev, &off)
}

func NewDateTimeWithTimeAndTimeZone(tt time.Time, abbreviation string, offset *int) *DateTime {
	d := &Date{}
	d.Year, d.Month, d.Day = tt.Date()
	d.wd = int(tt.Weekday()) + 1

	t := &Time{}
	t.Hour, t.Minute, t.Second = tt.Clock()

	tz := &TimeZone{}
	tz.Abbreviation = abbreviation

	if offset != nil {
		offSign := "+"
		if *offset < 0 {
			offSign = "-"
		}
		offAbs := int(math.Abs(float64(*offset)))
		offHH := offAbs / 3600
		offMM := (offAbs - (offHH * 3600)) / 60
		tz.Offset = fmt.Sprintf("%s%02d:%02d", offSign, offHH, offMM)
	}

	return NewDateTime(d, t, tz)
}

func NewDateTime(d *Date, t *Time, tz *TimeZone) *DateTime {
	// fmt.Printf("NewDateTime(d: %#v, t: %#v, tz %#v)\n", d, t, tz)
	r := &DateTime{Date: d, Time: t, TimeZone: tz}

	// If we have Date and Time info but no TimeZone, get it from Min if it exists.
	if r.TimeZone == nil &&
		r.Date != nil &&
		r.Time != nil &&
		minimumDateTime != nil &&
		minimumDateTime.TimeZone != nil {
		// fmt.Printf("in NewDateTimeFromRaw(), setting time zone to %#v\n", minimumDT.TimeZone)
		r.TimeZone = minimumDateTime.TimeZone
	}

	// If we have an unofficial TimeZone Name, replace it with an Abbreviation (e.g. "Eastern" -> "ET").
	if r.TimeZone != nil &&
		r.TimeZone.Name != "" {
		n := r.TimeZone.Name
		if _, err := time.LoadLocation(n); err != nil {
			if abbrev := timeZoneAbbreviationsByNames[strings.ToLower(n)]; abbrev != "" {
				r.TimeZone.Name = ""
				if r.TimeZone.Abbreviation == "" {
					r.TimeZone.Abbreviation = abbrev
				}
			}
		}
	}

	// Finalize Date based on TimeZone if ambiguous.
	r.Date = NewDateFromRaw(r.Date, r.TimeZone)

	// Do semantic check that we don't have gaps in the time scales (e.g. Month and Hour, but not Day).
	bits := []bool{
		r.Date.Year != 0,
		int(r.Date.Month) != 0,
		r.Date.Day != 0,
		r.Time != nil && r.Time.Hour != 0,
	}
	start := lo.IndexOf(bits, true)
	end := lo.LastIndexOf(bits, true)
	for i := start; i < end; i++ {
		if !bits[i] {
			// fmt.Printf("found gap in DateTime at i: %d in bits: %#v\n", i, bits)
			panic(fmt.Sprintf("semantic error: found gap in DateTime at i: %d in YMDH bits: %#v\n", i, bits))
		}
	}

	return r
}

func locationForName(name string) *time.Location {
	r, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("error getting time zone location from name (%q): %v", name, err))
	}
	// fmt.Printf("in locationForName(dt, name: %q), returning %#v\n", name, r)
	return r
}

func locationForOffset(dt *DateTime, offset string) *time.Location {
	if dt.Date == nil || dt.Date.Year == 0 || dt.Date.Month == 0 || dt.Date.Day == 0 {
		return nil
	}
	fakeStr := dt.Date.String() + "T00:00:00" + offset
	t, err := time.Parse(time.RFC3339, fakeStr)
	if err != nil {
		panic(fmt.Sprintf("error parsing fake time string (%q) for offset %q: %v", fakeStr, offset, err))
	}
	r := t.Location()
	// fmt.Printf("in locationForOffset(dt, offset: %q), returning %#v\n", offset, r)
	return r
}

func locationForAbbreviation(dt *DateTime, abbrev string) *time.Location {
	// slog.Debug("locationForAbbreviation(dt, abbrev)", "abbrev", abbrev)

	name := PreferredLocationNamesByAbbrev[abbrev]
	// slog.Debug("in locationForAbbreviation()", "name", name)
	if name != "" {
		if loc := locationForName(name); loc != nil {
			return loc
		}
	}

	tzs, _ := timezoneTZ.GetTzAbbreviationInfo(abbrev)
	// for i, tz := range tzs {
	// 	fmt.Printf("in locationForAbbreviation(dt, abbrev: %q), got tz[%d]: %#v\n", abbrev, i, tz)
	// }

	if len(tzs) == 1 {
		r := locationForOffset(dt, tzs[0].OffsetHHMM())
		// fmt.Printf("in locationForAbbreviation(dt, abbrev: %q), returning %#v\n", abbrev, r)
		return r
	}

	panic(fmt.Sprintf("error: no preferred time Location found for time zone abbreviation %q\n", abbrev))
}

// A Date represents a date (year, month, day, weekday).
//
// This type does not include location information, and therefore does not
// describe a unique 24-hour timespan.
//
// When a field is unspecified, it holds 0.
type Date struct {
	Year  int        // Year (e.g., 2014), starting at 1.
	Month time.Month // Month of the year, starting at 1 for January.
	Day   int        // Day of the month, starting at 1.
	// Weekday int        // Day of the week, starting at 1 for Sunday.
	unknown []any // Unprocessed Day and Month, with order depending upon the DateMode.
	wd      any   // Unprocessed Weekday, to be confirmed with computed Weekday.
}

// String returns the date in RFC3339 full-date format.
func (d *Date) String() string {
	if d == nil {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d *Date) ToTime() *time.Time {
	r := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
	return &r
}

func NewDateFromRaw(d *Date, tz *TimeZone) *Date {
	// fmt.Printf("NewDateFromRaw(d: %#v, tz: %#v)\n", d, tz)

	// If there's no ambiguity with the Month and Day, just process the raw Date.
	// fmt.Printf("checking whether Month and Day are ambiguous")
	if len(d.unknown) != 2 {
		r, err := maybeNewDateFromRaw(d, tz)
		if err != nil {
			panic(err.Error())
		}
		return r
	}

	dm := DateMode(tz)
	if dm == DateModeUnknown {
		dm = parseDateMode
	}
	if dm == DateModeUnknown && minimumDateTime != nil {
		dm = DateMode(minimumDateTime.TimeZone)
	}
	// fmt.Printf("found DateMode: %q\n", dm)

	// If we now know the DateMode for the Month and Day, process the raw Date.
	u0 := d.unknown[0]
	u1 := d.unknown[1]

	// fmt.Printf("going with DateModeRest\n")
	if dm == DateModeRest {
		setNewDateMonthAndDay(d, findInt(monthUnit, u1), findInt(dayUnit, u0))
		r, err := maybeNewDateFromRaw(d, tz)
		if err != nil {
			panic(err.Error())
		}
		return r
	}
	// fmt.Printf("going with DateModeNorthAmerican\n")
	if dm == DateModeNorthAmerican {
		setNewDateMonthAndDay(d, findInt(monthUnit, u0), findInt(dayUnit, u1))
		r, err := maybeNewDateFromRaw(d, tz)
		if err != nil {
			panic(err.Error())
		}
		return r
	}

	// At this point, we have the ambiguous Month and Day, and we don't know the DateMode. Try NA first.
	// fmt.Printf("trying with DateModeNorthAmerican\n")
	setNewDateMonthAndDay(d, findInt(monthUnit, u0), findInt(dayUnit, u1))
	r, err := maybeNewDateFromRaw(d, tz)
	if err == nil {
		return r
	}
	// fmt.Printf("trying with DateModeRest\n")
	setNewDateMonthAndDay(d, findInt(monthUnit, u1), findInt(dayUnit, u0))
	r, err = maybeNewDateFromRaw(d, tz)
	if err == nil {
		return r
	}
	panic(err.Error())
}

func setNewDateMonthAndDay(d *Date, month int, day int) *Date {
	d.Month = time.Month(month)
	d.Day = day
	return d
}

func maybeNewDateFromRaw(d *Date, tz *TimeZone) (*Date, error) {
	// fmt.Printf("maybeNewDateFromRaw(d: %#v, tz: %#v)\n", d, tz)

	// if Month and Day are both set, check consistency and set Year and Weekday.
	if d.Month == 0 || d.Day == 0 {
		return d, nil
	}

	// Fix year by setting date to be no earlier than minimumDateTime.
	if d.Year == 0 {
		setNewDateYear(d)
	}

	// Check the extracted Weekday with the computed Weekday.
	extWD := 0
	switch dwd := d.wd.(type) {
	case string:
		if dwd != "" {
			extWD = findInt(weekdayUnit, d.wd)
		}
	default:
		// TODO: should we do a type-specific nil check here?
		if dwd == nil {
			extWD = findInt(weekdayUnit, d.wd)
		}
	}

	wd := extWD
	if d.Year != 0 {
		t := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
		wd = weekdaysByNames[strings.ToLower(t.Weekday().String())]
		if extWD != 0 && extWD != wd {
			return nil, fmt.Errorf("semantic error: extracted weekday of %q doesn't match computed weekday of %q for %s\n",
				weekdayNames[extWD], weekdayNames[wd], d.String())
		}
	}
	// d.Weekday = wd

	return d, nil
}

func DateMode(tz *TimeZone) string {
	// fmt.Printf("DateMode(tz: %#v)\n", tz)
	if tz == nil {
		return DateModeUnknown
	}

	abbrev := tz.Abbreviation
	name := tz.Name

	// fmt.Printf("in DateMode(), name: %q, abbrev: %q\n", name, abbrev)
	if abbrev == "" {
		abbrev = timeZoneAbbreviationsByNames[name]
	}
	if northAmericanTimeZoneAbbreviations[abbrev] {
		return DateModeNorthAmerican
	}

	tzInfo, err := timezoneTZ.GetTzInfo(name)
	if err != nil {
		if DateModeNorthAmericanCountryCodes[tzInfo.CountryCode()] {
			return DateModeNorthAmerican
		}
	}

	return DateModeUnknown
}

func setNewDateYear(d *Date) *Date {
	// fmt.Printf("checking minimumDT: %#v\n", minimumDateTime)
	if minimumDateTime == nil {
		d.Year = 0
		return d
	}

	minTime := minimumDateTime.ToTime()
	// fmt.Printf("checking minTime: %#v\n", minTime)
	if minTime == nil {
		d.Year = 0
		return d
	}

	d.Year = minimumDateTime.Date.Year
	dateTime := d.ToTime()
	// fmt.Printf("checking same year dateTime: %#v\n", dateTime)
	if dateTime.After(*minTime) {
		return d
	}

	d.Year = minimumDateTime.Date.Year + 1
	dateTime = d.ToTime()
	// fmt.Printf("checking next year dateTime: %#v\n", dateTime)
	if dateTime.After(*minTime) {
		return d
	}

	panic(fmt.Sprintf("semantic error: unclear how to set year with minimumDateTime: %s\n", minimumDateTime.String()))
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
	Name         string
	Abbreviation string
	Offset       string
}

func (tz *TimeZone) String() string {
	if tz == nil {
		return "Z"
	}
	return tz.Offset
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
	return &TimeZone{Name: name, Abbreviation: abbrev, Offset: offset}
}

func findInt(tUnit timeUnit, val any) int {
	// fmt.Printf("findInt(tUnit: %#v, valAny (%T): %#v)\n", tUnit, val, val)
	r := -1
	var ok bool
	if val != nil {
		switch val := val.(type) {
		case int:
			r = val
		case time.Month:
			r = int(val)
		case *int:
			r = *val
		case string:
			if val == "" {
				ok = false
				r = 0
				break
			}
			rInt, err := strconv.Atoi(val)
			if err == nil {
				r = rInt
			} else {
				if tUnit.stringToIntFn != nil {
					r = tUnit.stringToIntFn(val)
				}
			}
		}
	}

	if tUnit.fixFn != nil {
		r, ok = tUnit.fixFn(val, r)
	} else if val == nil {
		return 0
	}

	// debugf("in findInt(), r: %d, ok: %t\n", r, ok)
	if !ok && (r < tUnit.min || r > tUnit.max) {
		// fmt.Printf("ok: %#v\n", ok)
		// fmt.Printf("r: %#v\n", r)
		// fmt.Printf("r < tunit.min: %#v\n", r < tUnit.min)
		// fmt.Printf("r > tunit.max: %#v\n", r > tUnit.max)
		panic(fmt.Sprintln("found int but failed bounds check", "tunit", tUnit, "val", val))
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

var yearUnit = timeUnit{name: "year", min: 1, max: math.MaxInt, fixFn: fixYear}
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

func fixYear(yearAny any, year int) (int, bool) {
	// fmt.Printf("fixYear(yearAny: %#v, year: %d)\n", yearAny, year)
	if yearAny == nil {
		return 0, true
	}

	return year, (year >= 1700 && year <= 2100)
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

func monthNameToMonth(monthName string) int {
	month, found := monthsByNames[strings.ToLower(monthName)]
	if !found {
		return 0
	}
	return month
}

var ordinals = map[string]bool{
	"st": true,
	"nd": true,
	"rd": true,
	// "th": true, recognize this separately because it also shortens Thursday
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
	"th":        5,
	"thu":       5,
	"thus":      5,
	"thursday":  5,
	"fr":        6,
	"fri":       6,
	"friday":    6,
	"sa":        7,
	"sat":       7,
	"saturday":  7,
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

var PreferredLocationNamesByAbbrev = map[string]string{
	"BST":  "Europe/London",
	"ET":   "America/New_York",
	"EDT":  "America/New_York",
	"EST":  "America/New_York",
	"CET":  "Europe/Paris",
	"CEST": "Europe/Paris",
	"CT":   "America/Chicago",
	"CDT":  "America/Chicago",
	"CST":  "America/Chicago",
	"GMT":  "Europe/London",
	"MT":   "America/Denver",
	"MDT":  "America/Denver",
	"MST":  "America/Denver",
	"PT":   "America/Los_Angeles",
	"PDT":  "America/Los_Angeles",
	"PST":  "America/Los_Angeles",
}

var timeZoneAbbreviationsByNames = map[string]string{
	"eastern":                      "ET",
	"eastern (new york)":           "ET",
	"eastern time (us and canada)": "ET",
	"us/eastern":                   "ET",
	"pacific":                      "PT",
	"pacific (los angeles)":        "PT",
	"pacific time (us and canada)": "PT",
	"us/pacific":                   "PT",
}

var northAmericanTimeZoneAbbreviations = map[string]bool{
	"CT": true,
	"ET": true,
	"MT": true,
	"PT": true,
}

var ignorableTimeZoneAbbreviations = map[string]bool{
	"M": true,
	"V": true,
}

func NewRawDateFromRelative(relativeName string) *Date {
	// fmt.Printf("NewRawDateFromRelative(relativeName: %q), minimumDateTime: %s\n", relativeName, minimumDateTime.String())
	if minimumDateTime == nil {
		panic(fmt.Sprintf("semantic error: found relativeName %q but minimumDateTimeName is nil\n", relativeName))
	}
	minT := time.Date(minimumDateTime.Date.Year, minimumDateTime.Date.Month, minimumDateTime.Date.Day, 0, 0, 0, 0, time.UTC)

	relName := strings.ToLower(relativeName)
	if relName == "yesterday" {
		y, m, d := minT.AddDate(0, 0, -1).Date()
		return &Date{Day: d, Month: m, Year: y}
	}
	if relName == "today" {
		return minimumDateTime.Date
	}
	if relName == "tomorrow" {
		y, m, d := minT.AddDate(0, 0, 1).Date()
		return &Date{Day: d, Month: m, Year: y}
	}

	wd := weekdaysByNames[relName]
	if wd == 0 {
		panic(fmt.Sprintf("semantic error: found unknown relativeName: %q\n", relativeName))
	}

	// fmt.Printf("in NewRawDateFromRelative(), wd name: %q\n", weekdayNames[wd])
	daysUntilNext := ((int(wd) - int(minT.Weekday()) + 7) % 7) - 1
	y, m, d := minT.AddDate(0, 0, daysUntilNext).Date()
	// fmt.Printf("in NewRawDateFromRelative(), y: %#v, m: %#v, d: %#v\n", y, m, d)
	return &Date{Day: d, Month: m, Year: y}
}

func NewRawDateFromAmbiguous(weekdayAny any, first string, last string, yearAny any) *Date {
	year := findInt(yearUnit, yearAny)
	return &Date{unknown: []any{first, last}, Year: year, wd: weekdayAny}
}

func NewRawDateFromDsMYs(daysAny []string, monthAny any, yearAny any) []*Date {
	rs := []*Date{}
	for _, dayAny := range daysAny {
		rs = append(rs, NewRawDateFromDMY(dayAny, monthAny, yearAny))
	}
	return rs
}

func NewRawDateFromDMY(dayAny any, monthAny any, yearAny any) *Date {
	// fmt.Printf("NewRawDateFromDMY(dayAny: %#v, monthAny %#v, yearAny %#v)\n", dayAny, monthAny, yearAny)
	day := findInt(dayUnit, dayAny)
	month := findInt(monthUnit, monthAny)
	year := findInt(yearUnit, yearAny)
	return &Date{Day: day, Month: time.Month(month), Year: year}
}

func NewRawDateFromWDMY(weekdayAny any, dayAny any, monthAny any, yearAny any) *Date {
	// fmt.Printf("NewRawDateFromWDMY(weekdayAny: %#v, dayAny: %#v, monthAny %#v, yearAny %#v)\n", weekdayAny, dayAny, monthAny, yearAny)
	day := findInt(dayUnit, dayAny)
	month := findInt(monthUnit, monthAny)
	year := findInt(yearUnit, yearAny)
	return &Date{Day: day, Month: time.Month(month), Year: year, wd: weekdayAny}
}

func NewRawDateFromMDsYs(monthAny any, daysAny []string, yearAny any) []*Date {
	rs := []*Date{}
	for _, dayAny := range daysAny {
		rs = append(rs, NewRawDateFromMDY(monthAny, dayAny, yearAny))
	}
	return rs
}

func NewRawDateFromMDY(monthAny any, dayAny any, yearAny any) *Date {
	return NewRawDateFromDMY(dayAny, monthAny, yearAny)
}

func NewRawDateFromWMDY(weekdayAny any, monthAny any, dayAny any, yearAny any) *Date {
	return NewRawDateFromWDMY(weekdayAny, dayAny, monthAny, yearAny)
}

func NewRawDateFromYMD(yearAny any, monthAny any, dayAny any) *Date {
	return NewRawDateFromDMY(dayAny, monthAny, yearAny)
}

func NewAMTime(hourAny any, minuteAny any, secondAny any, nsAny any) *Time {
	r := NewTime(hourAny, minuteAny, secondAny, nsAny)
	if r.Hour > 12 {
		panic(fmt.Sprintf("semantic error: found hour %#v but failed AM bounds check\n", r.Hour))
	}
	r.Hour = r.Hour % 12
	return r
}

func NewPMTime(hourAny any, minuteAny any, secondAny any, nsAny any) *Time {
	r := NewTime(hourAny, minuteAny, secondAny, nsAny)
	if r.Hour > 12 {
		panic(fmt.Sprintf("semantic error: found hour %#v but failed PM bounds check\n", r.Hour))
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
