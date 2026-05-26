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
	Items      []*DateTimeRange
	Recurrence *Recurrence // optional, set when parser detects recurrence pattern
}

func (rngs *DateTimeRanges) String() string {
	if rngs == nil {
		return ""
	}
	rs := []string{}
	for _, elt := range rngs.Items {
		rs = append(rs, elt.String())
	}
	s := strings.Join(rs, ", ")
	if rngs.Recurrence != nil {
		s += " (" + rngs.Recurrence.String() + ")"
	}
	return s
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

// NewRangesFromDatesTimeRange creates one time range per date from a list of
// dates, a shared start time, end time, and optional timezone. Used for
// patterns like "February 1, 8, 15 9am-12pm ET" → three ranges.
func NewRangesFromDatesTimeRange(dates []*Date, startTime *Time, endTime *Time, tz *TimeZone) *DateTimeRanges {
	r := &DateTimeRanges{}
	for _, date := range dates {
		r.Items = append(r.Items, NewRange(
			NewDateTime(date, startTime, tz),
			NewDateTime(date, endTime, tz)))
	}
	return r
}

func NewRangesWithStartEndDates(start *Date, end *Date) *DateTimeRanges {
	return &DateTimeRanges{Items: []*DateTimeRange{NewRangeWithStartEndDates(start, end)}}
}

func NewRangesWithStartEndDateTimes(start *DateTime, end *DateTime) *DateTimeRanges {
	return &DateTimeRanges{Items: []*DateTimeRange{NewRange(start, end)}}
}

// Frequency represents how often a recurring event repeats.
type Frequency int

const (
	FrequencyDaily   Frequency = iota + 1
	FrequencyWeekly
	FrequencyMonthly
	FrequencyYearly
)

// String returns a human-readable frequency name.
func (f Frequency) String() string {
	return frequencyStrings[f]
}

var frequencyStrings = map[Frequency]string{
	FrequencyDaily:   "daily",
	FrequencyWeekly:  "weekly",
	FrequencyMonthly: "monthly",
	FrequencyYearly:  "yearly",
}

// frequencySteps maps frequency → (years, months, days) for Occurrences() expansion.
var frequencySteps = map[Frequency][3]int{
	FrequencyDaily:   {0, 0, 1},
	FrequencyWeekly:  {0, 0, 7},
	FrequencyMonthly: {0, 1, 0},
	FrequencyYearly:  {1, 0, 0},
}

// Recurrence captures recurrence metadata for patterns like
// "5 Wednesdays 9:00am-12:00pm February 1st - March 1st".
type Recurrence struct {
	Frequency  Frequency       // how often the event repeats (daily, weekly, etc.)
	Weekdays   []time.Weekday  // e.g. [Tue, Thu] for "Tuesdays and Thursdays"
	NthWeekday []int           // e.g. [2, 4] for "2nd and 4th"
	Count      int             // "5 Wednesdays" → 5 (0 means use Until)
	Until      *Date           // end boundary date (alternative to Count)
}

// String returns a human-readable representation of the recurrence.
func (r *Recurrence) String() string {
	if r == nil {
		return ""
	}
	s := r.Frequency.String()
	if r.Count > 0 {
		s += " x" + strconv.Itoa(r.Count)
	}
	for _, wd := range r.Weekdays {
		s += " " + wd.String()
	}
	if r.Until != nil {
		s += " until " + r.Until.String()
	}
	return s
}

// NewRecurringRanges creates a DateTimeRanges with recurrence metadata.
func NewRecurringRanges(first *DateTimeRange, rec *Recurrence) *DateTimeRanges {
	return &DateTimeRanges{
		Items:      []*DateTimeRange{first},
		Recurrence: rec,
	}
}

// Occurrences returns the expanded list of ranges regardless of representation.
// For recurring: expands by stepping through dates at the given frequency.
// For flat: returns Items as-is.
func (rngs *DateTimeRanges) Occurrences() []*DateTimeRange {
	if rngs.Recurrence == nil {
		return rngs.Items
	}
	rec := rngs.Recurrence
	first := rngs.Items[0]
	if first.Start == nil || first.Start.Date == nil || first.Start.Date.Year == 0 {
		panic("Occurrences: recurring pattern requires year on first occurrence")
	}

	// Compute start-to-end day delta for multi-day ranges (0 for same-day).
	var endDayDelta int
	if first.End != nil && first.End.Date != nil {
		endDayDelta = int(first.End.Date.ToTime().Sub(*first.Start.Date.ToTime()).Hours() / 24)
	}

	step := frequencySteps[rec.Frequency]
	var result []*DateTimeRange
	current := *first.Start.Date.ToTime()
	for {
		if rec.Count > 0 && len(result) >= rec.Count {
			break
		}
		if rec.Until != nil && current.After(*rec.Until.ToTime()) {
			break
		}
		if rec.Count == 0 && rec.Until == nil {
			break
		}

		startDate := &Date{Year: current.Year(), Month: current.Month(), Day: current.Day()}
		var endDT *DateTime
		if first.End != nil {
			endCurrent := current.AddDate(0, 0, endDayDelta)
			endDate := &Date{Year: endCurrent.Year(), Month: endCurrent.Month(), Day: endCurrent.Day()}
			endDT = &DateTime{Date: endDate, Time: first.End.Time, TimeZone: first.End.TimeZone}
		}
		result = append(result, &DateTimeRange{
			Start: &DateTime{Date: startDate, Time: first.Start.Time, TimeZone: first.Start.TimeZone},
			End:   endDT,
		})
		current = current.AddDate(step[0], step[1], step[2])
	}
	return result
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
}

func (rng DateTimeRange) String() string {
	r := rng.Start.String()
	if rng.End != nil {
		r += " - " + rng.End.String()
	}
	return r
}

// IANAName returns the IANA timezone name (e.g., "America/Los_Angeles") for this range.
// Uses the Start datetime's timezone.
func (rng *DateTimeRange) IANAName() string {
	if rng == nil || rng.Start == nil {
		return ""
	}
	return rng.Start.IANAName()
}

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

		if start.Time != nil && end.Time != nil {
			if start.Date == nil && end.Date != nil {
				start.Date = end.Date
			} else if start.Date != nil && end.Date == nil {
				end.Date = start.Date
			}
		}

		if start.Date != nil && end.Date != nil {
			// If both have Days but one Month is missing, copy from the other.
			if start.Date.Day != 0 && end.Date.Day != 0 {
				if start.Date.Month == 0 && end.Date.Month != 0 {
					start.Date.Month = end.Date.Month
				} else if start.Date.Month != 0 && end.Date.Month == 0 {
					end.Date.Month = start.Date.Month
				}
			}

			// If both have Months but one Year is missing, copy from the other.
			// Year straddling: if start month > end month (e.g. "25 Dec - 2 Jan 2016"),
			// the range crosses a year boundary, so start year = end year - 1.
			if start.Date.Month != 0 && end.Date.Month != 0 {
				if start.Date.Year == 0 && end.Date.Year != 0 {
					if start.Date.Month > end.Date.Month {
						start.Date.Year = end.Date.Year - 1
					} else {
						start.Date.Year = end.Date.Year
					}
				} else if start.Date.Year != 0 && end.Date.Year == 0 {
					end.Date.Year = start.Date.Year
				}
			}
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

// IANAName returns the IANA timezone name (e.g., "America/Los_Angeles") for this DateTime.
// Delegates to TimeZone.IANAName().
func (dt *DateTime) IANAName() string {
	if dt == nil || dt.TimeZone == nil {
		return ""
	}
	return dt.TimeZone.IANAName()
}

func (dt *DateTime) Location() *time.Location {
	if dt.TimeZone == nil {
		if DoDebug.Load() {
			fmt.Println("warning: no TimeZone found in DateTime, returning nil Location")
		}
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
	if dt == nil {
		return nil
	}

	t := &Time{}
	if dt.Time != nil {
		t = dt.Time
	}

	loc := dt.Location()
	if loc == nil {
		loc = time.UTC
	}

	r := time.Date(dt.Date.Year, dt.Date.Month, dt.Date.Day, t.Hour, t.Minute, t.Second, t.Nanosecond, loc)
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
	wd := tt.Weekday()
	d.wd = &wd

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
	r := &DateTime{Date: d, Time: t, TimeZone: tz}

	// If we have Date and Time info but no TimeZone, get it from Min if it exists.
	if r.TimeZone == nil &&
		r.Date != nil &&
		r.Time != nil &&
		minimumDateTime != nil &&
		minimumDateTime.TimeZone != nil {
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
	if r.Date != nil {
		r.Date = NewDateFromRaw(r.Date, r.TimeZone)
	}

	// Do semantic check that we don't have gaps in the time scales (e.g. Month and Hour, but not Day).
	if r.Date != nil && r.Date.Day != 0 {
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
				panic(fmt.Sprintf("semantic error: found gap in DateTime at i: %d in YMDH bits: %#v\n", i, bits))
			}
		}
	}

	return r
}

func locationForName(name string) *time.Location {
	r, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("error getting time zone location from name (%q): %v", name, err))
	}
	return r
}

func locationForOffset(dt *DateTime, offset string) *time.Location {
	if dt.Date == nil || dt.Date.Year == 0 || dt.Date.Month == 0 || dt.Date.Day == 0 {
		return nil
	}
	fakeStr := dt.Date.String() + "T00:00:00" + offset
	t, err := time.Parse(time.RFC3339, fakeStr)
	if err != nil {
		// Malformed date (e.g., "Nov 31") — caller falls back through other
		// timezone lookups when this returns nil.
		return nil
	}
	r := t.Location()
	return r
}

func locationForAbbreviation(dt *DateTime, abbrev string) *time.Location {
	name := PreferredLocationNamesByAbbrev[abbrev]
	if name != "" {
		if loc := locationForName(name); loc != nil {
			return loc
		}
	}

	tzs, _ := timezoneTZ.GetTzAbbreviationInfo(abbrev)
	if len(tzs) == 1 {
		r := locationForOffset(dt, tzs[0].OffsetHHMM())
		return r
	}

	panic(fmt.Sprintf("no preferred time Location found for time zone abbreviation %q — add it to PreferredLocationNamesByAbbrev", abbrev))
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
	unknown []any          // Unprocessed Day and Month, with order depending upon the DateMode.
	wd      *time.Weekday // nil = unset, non-nil = specific weekday
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
	// If there's no ambiguity with the Month and Day, just process the raw Date.
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

	// If we now know the DateMode for the Month and Day, process the raw Date.
	u0 := d.unknown[0]
	u1 := d.unknown[1]

	if dm == DateModeRest {
		setNewDateMonthAndDay(d, findInt(monthUnit, u1), findInt(dayUnit, u0))
		r, err := maybeNewDateFromRaw(d, tz)
		if err != nil {
			panic(err.Error())
		}
		return r
	}
	if dm == DateModeNorthAmerican {
		setNewDateMonthAndDay(d, findInt(monthUnit, u0), findInt(dayUnit, u1))
		r, err := maybeNewDateFromRaw(d, tz)
		if err != nil {
			panic(err.Error())
		}
		return r
	}

	// At this point, we have the ambiguous Month and Day, and we don't know the DateMode. Try NA first.
	setNewDateMonthAndDay(d, findInt(monthUnit, u0), findInt(dayUnit, u1))
	r, err := maybeNewDateFromRaw(d, tz)
	if err == nil {
		return r
	}
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
	// if Month and Day are both set, check consistency and set Year and Weekday.
	if d.Month == 0 || d.Day == 0 {
		return d, nil
	}

	// Fix year by setting date to be no earlier than minimumDateTime.
	if d.Year == 0 {
		setNewDateYear(d)
	}

	// Check the extracted Weekday with the computed Weekday.
	if d.Year != 0 && d.wd != nil {
		computed := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).Weekday()
		if *d.wd != computed {
			return nil, fmt.Errorf("semantic error: weekday %s doesn't match computed %s for %s",
				d.wd, computed, d.String())
		}
	}

	return d, nil
}

func DateMode(tz *TimeZone) string {
	if tz == nil {
		return DateModeUnknown
	}

	abbrev := tz.Abbreviation
	name := tz.Name

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
	if minimumDateTime == nil || minimumDateTime.Date == nil {
		d.Year = 0
		return d
	}

	d.Year = minimumDateTime.Date.Year
	return d
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

// IANAName returns the IANA timezone name (e.g., "America/Los_Angeles") for this TimeZone.
// It checks Name first (if it's already an IANA name), then falls back to looking up
// the Abbreviation in PreferredLocationNamesByAbbrev.
// Returns empty string if no IANA name can be determined.
func (tz *TimeZone) IANAName() string {
	if tz == nil {
		return ""
	}

	// If Name is already set and is a valid IANA name, use it
	if tz.Name != "" {
		if _, err := time.LoadLocation(tz.Name); err == nil {
			return tz.Name
		}
	}

	// Look up from abbreviation
	if tz.Abbreviation != "" {
		if name := PreferredLocationNamesByAbbrev[tz.Abbreviation]; name != "" {
			return name
		}
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
	return &TimeZone{Name: name, Abbreviation: abbrev, Offset: offset}
}

func findInt(tUnit timeUnit, val any) int {
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

	if !ok && (r < tUnit.min || r > tUnit.max) {
		panic(fmt.Sprintln("found int but failed bounds check", "tunit", tUnit, "val", val))
	}

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
var hourUnit = timeUnit{name: "hour", min: 0, max: 24, stringToIntFn: hourNameToHour}
var minuteUnit = timeUnit{name: "minute", min: 0, max: 59}
var secondUnit = timeUnit{name: "second", min: 0, max: 59}
var nsUnit = timeUnit{name: "ns", min: 0, max: 999}

func fixYear(yearAny any, year int) (int, bool) {
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

var weekdaysByNames = map[string]time.Weekday{
	"su":          time.Sunday,
	"sun":         time.Sunday,
	"sunday":      time.Sunday,
	"sundays":     time.Sunday,
	"mo":          time.Monday,
	"mon":         time.Monday,
	"monday":      time.Monday,
	"mondays":     time.Monday,
	"tu":          time.Tuesday,
	"tue":         time.Tuesday,
	"tues":        time.Tuesday,
	"tuesday":     time.Tuesday,
	"tuesdays":    time.Tuesday,
	"we":          time.Wednesday,
	"wed":         time.Wednesday,
	"weds":        time.Wednesday,
	"wednesday":   time.Wednesday,
	"wednesdays":  time.Wednesday,
	"th":          time.Thursday,
	"thu":         time.Thursday,
	"thus":        time.Thursday,
	"thursday":    time.Thursday,
	"thursdays":   time.Thursday,
	"fr":          time.Friday,
	"fri":         time.Friday,
	"friday":      time.Friday,
	"fridays":     time.Friday,
	"sa":          time.Saturday,
	"sat":         time.Saturday,
	"saturday":    time.Saturday,
	"saturdays":   time.Saturday,
}

// weekdayFromName converts a weekday name string to *time.Weekday.
// Returns nil for empty/unrecognized names.
func weekdayFromName(name string) *time.Weekday {
	if name == "" {
		return nil
	}
	wd, ok := weekdaysByNames[strings.ToLower(name)]
	if !ok {
		return nil
	}
	return &wd
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
	"AEST": "Australia/Sydney",
	"AEDT": "Australia/Sydney",
	"AWST": "Australia/Perth",
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
	"NZST": "Pacific/Auckland",
	"NZDT": "Pacific/Auckland",
	"PT":   "America/Los_Angeles",
	"PDT":  "America/Los_Angeles",
	"PST":  "America/Los_Angeles",
	"SAST": "Africa/Johannesburg",
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

// Single-char timezone abbreviations are now ignored via length check in parse_lex.go

func NewRawDateFromRelative(relativeName string) *Date {
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

	wd, ok := weekdaysByNames[relName]
	if !ok {
		panic(fmt.Sprintf("semantic error: found unknown relativeName: %q\n", relativeName))
	}

	daysUntilNext := (int(wd) - int(minT.Weekday()) + 7) % 7
	y, m, d := minT.AddDate(0, 0, daysUntilNext).Date()
	return &Date{Day: d, Month: m, Year: y}
}

func NewRawDateFromAmbiguous(weekdayName string, first string, last string, yearAny any) *Date {
	year := findInt(yearUnit, yearAny)
	return &Date{unknown: []any{first, last}, Year: year, wd: weekdayFromName(weekdayName)}
}

func NewRawDateFromDsMYs(daysAny []string, monthAny any, yearAny any) []*Date {
	rs := []*Date{}
	for _, dayAny := range daysAny {
		rs = append(rs, NewRawDateFromDMY(dayAny, monthAny, yearAny))
	}
	return rs
}

func NewRawDateFromDMY(dayAny any, monthAny any, yearAny any) *Date {
	day := findInt(dayUnit, dayAny)
	month := findInt(monthUnit, monthAny)
	year := findInt(yearUnit, yearAny)
	return &Date{Day: day, Month: time.Month(month), Year: year}
}

func NewRawDateFromWDMY(weekdayName string, dayAny any, monthAny any, yearAny any) *Date {
	day := findInt(dayUnit, dayAny)
	month := findInt(monthUnit, monthAny)
	year := findInt(yearUnit, yearAny)
	return &Date{Day: day, Month: time.Month(month), Year: year, wd: weekdayFromName(weekdayName)}
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

func NewRawDateFromWMDY(weekdayName string, monthAny any, dayAny any, yearAny any) *Date {
	return NewRawDateFromWDMY(weekdayName, dayAny, monthAny, yearAny)
}

func NewRawDateFromYMD(yearAny any, monthAny any, dayAny any) *Date {
	return NewRawDateFromDMY(dayAny, monthAny, yearAny)
}

func NewAMTime(hourAny any, minuteAny any, secondAny any, nsAny any) *Time {
	r := NewTime(hourAny, minuteAny, secondAny, nsAny)
	if r.Hour > 12 {
		// 24-hour time with redundant AM suffix (e.g., "07:00 AM" parsed as
		// hour=7 is fine, but "00:00 AM" where hour=0 is also valid).
		// Hours >12 with AM are invalid — reject this parse alternative.
		panic(fmt.Sprintf("semantic error: found hour %#v but failed AM bounds check\n", r.Hour))
	}
	r.Hour = r.Hour % 12
	return r
}

func NewPMTime(hourAny any, minuteAny any, secondAny any, nsAny any) *Time {
	r := NewTime(hourAny, minuteAny, secondAny, nsAny)
	if r.Hour > 12 {
		// 24-hour time with redundant PM suffix (e.g., "17:00 PM").
		// Treat as 24-hour time — ignore the PM indicator.
		return r
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
