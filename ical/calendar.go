package ical

import (
	"fmt"
	"time"

	"github.com/findyourpaths/phil/datetime"

	ics "github.com/arran4/golang-ical"
	"github.com/segmentio/ksuid"
)

// EventInfo provides metadata for iCal VEVENT creation.
type EventInfo struct {
	Summary     string // SUMMARY — event title (required)
	Description string // DESCRIPTION — event description
	Location    string // LOCATION — venue/address
	URL         string // URL — event page link
}

// resolvedEvent is the fully-resolved intermediate representation.
// All output formats (ics, Google URL, Outlook URL, Yahoo URL) derive from this.
// Produced by Resolve(), which does all validation.
type resolvedEvent struct {
	Start       time.Time
	End         time.Time
	AllDay      bool
	Summary     string
	Description string
	Location    string
	URL         string
	RRule       string // empty if non-recurring, e.g. "FREQ=WEEKLY;COUNT=5;BYDAY=WE"
	TimeZone    string // IANA name, e.g. "America/New_York"
}

// ResolvedEvents holds the validated, fully-concrete intermediate result.
// All output functions (NewCalendarFrom, GoogleURLFrom, etc.) derive from this.
// Use Resolve() to create.
type ResolvedEvents struct {
	events []*resolvedEvent
}

// Resolve converts DateTimeRanges + EventInfo into a validated intermediate form.
// This is the single point where:
//   - DateTimeRanges are interpreted (recurring vs flat)
//   - Semantic validation runs (weekday, count, ordering, validity)
//   - Timezone is resolved to IANA name
//   - All-day vs timed is determined
//
// Callers generating multiple outputs (ICS + URLs) should call Resolve() once
// and pass the result to NewCalendarFrom(), GoogleURLFrom(), etc.
func Resolve(ranges *datetime.DateTimeRanges, info *EventInfo) (*ResolvedEvents, error) {
	events, err := resolve(ranges, info)
	if err != nil {
		return nil, err
	}
	return &ResolvedEvents{events: events}, nil
}

func resolve(ranges *datetime.DateTimeRanges, info *EventInfo) ([]*resolvedEvent, error) {
	if ranges == nil || len(ranges.Items) == 0 {
		return nil, fmt.Errorf("no date/time ranges provided")
	}

	if ranges.Recurrence != nil {
		return resolveRecurring(ranges, info)
	}
	return resolveFlat(ranges, info)
}

func resolveFlat(ranges *datetime.DateTimeRanges, info *EventInfo) ([]*resolvedEvent, error) {
	var events []*resolvedEvent
	for _, item := range ranges.Items {
		ev, err := resolveItem(item, "", info)
		if err != nil {
			return nil, err
		}
		events = append(events, ev)
	}
	return events, nil
}

func resolveRecurring(ranges *datetime.DateTimeRanges, info *EventInfo) ([]*resolvedEvent, error) {
	rec := ranges.Recurrence
	first := ranges.Items[0]

	if rec.Frequency == 0 {
		return nil, fmt.Errorf("recurrence missing frequency")
	}
	if rec.Count == 0 && rec.Until == nil {
		return nil, fmt.Errorf("recurrence requires at least one of Count or Until")
	}

	// Validate weekday-date consistency.
	if first.Start != nil && first.Start.Date != nil && first.Start.Date.Year != 0 && rec.Weekday != nil {
		computed := time.Date(first.Start.Date.Year, first.Start.Date.Month, first.Start.Date.Day,
			0, 0, 0, 0, time.UTC).Weekday()
		if *rec.Weekday != computed {
			return nil, fmt.Errorf("weekday mismatch: start date is %s, not %s", computed, *rec.Weekday)
		}
	}

	// Validate count vs until consistency.
	if rec.Count > 0 && rec.Until != nil {
		occs := ranges.Occurrences()
		if len(occs) != rec.Count {
			return nil, fmt.Errorf("count %d but only %d occurrences fit within the date range", rec.Count, len(occs))
		}
	}

	rrule := FormatRRule(rec)
	ev, err := resolveItem(first, rrule, info)
	if err != nil {
		return nil, err
	}
	return []*resolvedEvent{ev}, nil
}

func resolveItem(item *datetime.DateTimeRange, rrule string, info *EventInfo) (*resolvedEvent, error) {
	if item.Start == nil || item.Start.Date == nil {
		return nil, fmt.Errorf("date/time range missing start date")
	}

	sd := item.Start.Date
	allDay := item.Start.Time == nil

	// Validate start date.
	if sd.Month != 0 && sd.Day != 0 && sd.Year != 0 {
		check := time.Date(sd.Year, sd.Month, sd.Day, 0, 0, 0, 0, time.UTC)
		if check.Day() != sd.Day {
			return nil, fmt.Errorf("invalid date: %04d-%02d-%02d", sd.Year, sd.Month, sd.Day)
		}
	}

	// Resolve timezone.
	tzName := item.IANAName()
	loc := time.UTC
	if tzName != "" {
		var err error
		loc, err = time.LoadLocation(tzName)
		if err != nil {
			return nil, fmt.Errorf("unknown timezone: %s", tzName)
		}
	}

	// Build start time.
	var start time.Time
	if allDay {
		start = time.Date(sd.Year, sd.Month, sd.Day, 0, 0, 0, 0, loc)
	} else {
		t := item.Start.Time
		start = time.Date(sd.Year, sd.Month, sd.Day, t.Hour, t.Minute, t.Second, 0, loc)
	}

	// Build end time.
	var end time.Time
	if item.End != nil && item.End.Date != nil {
		ed := item.End.Date
		// Validate end date.
		if ed.Month != 0 && ed.Day != 0 && ed.Year != 0 {
			check := time.Date(ed.Year, ed.Month, ed.Day, 0, 0, 0, 0, time.UTC)
			if check.Day() != ed.Day {
				return nil, fmt.Errorf("invalid date: %04d-%02d-%02d", ed.Year, ed.Month, ed.Day)
			}
		}
		if allDay {
			// iCal: DTEND is exclusive for all-day events.
			end = time.Date(ed.Year, ed.Month, ed.Day, 0, 0, 0, 0, loc).AddDate(0, 0, 1)
		} else if item.End.Time != nil {
			et := item.End.Time
			end = time.Date(ed.Year, ed.Month, ed.Day, et.Hour, et.Minute, et.Second, 0, loc)
		} else {
			end = time.Date(ed.Year, ed.Month, ed.Day, 0, 0, 0, 0, loc).AddDate(0, 0, 1)
		}
	} else if allDay {
		// Single all-day event: DTEND = DTSTART + 1 day.
		end = start.AddDate(0, 0, 1)
	}

	// Validate ordering.
	if !end.IsZero() && start.After(end) {
		return nil, fmt.Errorf("start after end: %s > %s", start.Format(time.RFC3339), end.Format(time.RFC3339))
	}

	return &resolvedEvent{
		Start:       start,
		End:         end,
		AllDay:      allDay,
		Summary:     info.Summary,
		Description: info.Description,
		Location:    info.Location,
		URL:         info.URL,
		RRule:       rrule,
		TimeZone:    tzName,
	}, nil
}

// newVEvent creates a single VEVENT from a resolvedEvent.
func newVEvent(ev *resolvedEvent) *ics.VEvent {
	event := ics.NewEvent(ksuid.New().String())
	event.SetProperty(ics.ComponentPropertyDtstamp,
		time.Now().UTC().Format("20060102T150405Z"))

	if ev.AllDay {
		event.SetAllDayStartAt(ev.Start)
		if !ev.End.IsZero() {
			event.SetAllDayEndAt(ev.End)
		}
	} else if ev.TimeZone != "" {
		// DTSTART/DTEND with TZID parameter (local time, no Z suffix).
		tzParam := &ics.KeyValues{Key: "TZID", Value: []string{ev.TimeZone}}
		event.SetProperty(ics.ComponentPropertyDtStart,
			ev.Start.Format("20060102T150405"), tzParam)
		if !ev.End.IsZero() {
			event.SetProperty(ics.ComponentPropertyDtEnd,
				ev.End.Format("20060102T150405"), tzParam)
		}
	} else {
		event.SetStartAt(ev.Start)
		if !ev.End.IsZero() {
			event.SetEndAt(ev.End)
		}
	}

	if ev.Summary != "" {
		event.SetSummary(ev.Summary)
	}
	if ev.Description != "" {
		event.SetDescription(ev.Description)
	}
	if ev.Location != "" {
		event.SetLocation(ev.Location)
	}
	if ev.URL != "" {
		event.SetURL(ev.URL)
	}
	if ev.RRule != "" {
		event.AddRrule(ev.RRule)
	}

	return event
}

// NewCalendar converts DateTimeRanges and event metadata into an iCal Calendar.
// Convenience wrapper: calls Resolve() then NewCalendarFrom().
func NewCalendar(ranges *datetime.DateTimeRanges, info *EventInfo) (*ics.Calendar, error) {
	resolved, err := Resolve(ranges, info)
	if err != nil {
		return nil, err
	}
	return NewCalendarFrom(resolved), nil
}

// NewCalendarFrom creates an iCal Calendar from pre-resolved events.
// If the original ranges had Recurrence, the single VEVENT contains an RRULE.
// If flat, produces one VEVENT per resolved event.
func NewCalendarFrom(resolved *ResolvedEvents) *ics.Calendar {
	cal := ics.NewCalendar()
	cal.SetProductId("-//Phil//EN")
	for _, ev := range resolved.events {
		cal.AddVEvent(newVEvent(ev))
	}
	return cal
}

// ICS returns the serialized .ics content for the calendar.
func ICS(cal *ics.Calendar) string {
	return cal.Serialize()
}
