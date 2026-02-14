package ical

import (
	"fmt"
	"net/url"
	"time"

	"github.com/findyourpaths/phil/datetime"
)

// Time format layouts for calendar URL providers.
const (
	layoutCompact = "20060102T150405"   // Google, Yahoo
	layoutISO     = "2006-01-02T15:04:05" // Outlook
)

// formatCalTime formats a time for calendar URLs. If a timezone is present,
// returns local time in the given layout. Otherwise returns UTC with Z suffix.
func formatCalTime(t time.Time, tz, layout string) string {
	if tz != "" {
		return t.Format(layout)
	}
	return t.UTC().Format(layout + "Z")
}

// GoogleURL returns a Google Calendar event creation URL.
// Convenience wrapper: calls Resolve() then GoogleURLFrom().
func GoogleURL(ranges *datetime.DateTimeRanges, info *EventInfo) (string, error) {
	resolved, err := Resolve(ranges, info)
	if err != nil {
		return "", err
	}
	return GoogleURLFrom(resolved), nil
}

// GoogleURLFrom returns a Google Calendar event creation URL from pre-resolved events.
// Uses the first resolved event only. For recurring events, the RRULE covers
// all occurrences. For flat multi-event lists, callers should generate URLs
// per-range (calendar URLs can only represent a single event).
func GoogleURLFrom(resolved *ResolvedEvents) string {
	ev := resolved.events[0]

	params := url.Values{}
	params.Set("text", ev.Summary)

	if ev.AllDay {
		params.Set("dates", ev.Start.Format("20060102")+"/"+ev.End.Format("20060102"))
	} else {
		params.Set("dates", formatCalTime(ev.Start, ev.TimeZone, layoutCompact)+"/"+formatCalTime(ev.End, ev.TimeZone, layoutCompact))
	}

	if ev.Description != "" {
		params.Set("details", ev.Description)
	}
	if ev.Location != "" {
		params.Set("location", ev.Location)
	}
	if ev.TimeZone != "" {
		params.Set("ctz", ev.TimeZone)
	}
	if ev.RRule != "" {
		params.Set("recur", "RRULE:"+ev.RRule)
	}

	return "https://calendar.google.com/calendar/r/eventedit?" + params.Encode()
}

// OutlookURL returns an Outlook Web calendar event creation URL.
// Convenience wrapper: calls Resolve() then OutlookURLFrom().
func OutlookURL(ranges *datetime.DateTimeRanges, info *EventInfo) (string, error) {
	resolved, err := Resolve(ranges, info)
	if err != nil {
		return "", err
	}
	return OutlookURLFrom(resolved), nil
}

// OutlookURLFrom returns an Outlook Web calendar event creation URL from pre-resolved events.
// Uses the first resolved event only.
func OutlookURLFrom(resolved *ResolvedEvents) string {
	ev := resolved.events[0]

	params := url.Values{}
	params.Set("rru", "addevent")
	params.Set("subject", ev.Summary)

	if ev.AllDay {
		params.Set("startdt", ev.Start.Format("2006-01-02"))
		params.Set("enddt", ev.End.Format("2006-01-02"))
		params.Set("allday", "true")
	} else {
		params.Set("startdt", formatCalTime(ev.Start, ev.TimeZone, layoutISO))
		params.Set("enddt", formatCalTime(ev.End, ev.TimeZone, layoutISO))
	}

	if ev.Description != "" {
		params.Set("body", ev.Description)
	}
	if ev.Location != "" {
		params.Set("location", ev.Location)
	}

	return "https://outlook.live.com/calendar/0/action/compose?" + params.Encode()
}

// YahooURL returns a Yahoo Calendar event creation URL.
// Convenience wrapper: calls Resolve() then YahooURLFrom().
func YahooURL(ranges *datetime.DateTimeRanges, info *EventInfo) (string, error) {
	resolved, err := Resolve(ranges, info)
	if err != nil {
		return "", err
	}
	return YahooURLFrom(resolved), nil
}

// YahooURLFrom returns a Yahoo Calendar event creation URL from pre-resolved events.
// Uses the first resolved event only.
func YahooURLFrom(resolved *ResolvedEvents) string {
	ev := resolved.events[0]

	params := url.Values{}
	params.Set("v", "60")
	params.Set("title", ev.Summary)

	if ev.AllDay {
		params.Set("st", ev.Start.Format("20060102"))
		// Yahoo uses duration for all-day; calculate days.
		days := int(ev.End.Sub(ev.Start).Hours() / 24)
		params.Set("dur", fmt.Sprintf("allday%d", days))
	} else {
		params.Set("st", formatCalTime(ev.Start, ev.TimeZone, layoutCompact))
		// Yahoo uses duration in HHMM format.
		dur := ev.End.Sub(ev.Start)
		hours := int(dur.Hours())
		minutes := int(dur.Minutes()) % 60
		params.Set("dur", fmt.Sprintf("%02d%02d", hours, minutes))
	}

	if ev.Description != "" {
		params.Set("desc", ev.Description)
	}
	if ev.Location != "" {
		params.Set("in_loc", ev.Location)
	}

	return "https://calendar.yahoo.com/?" + params.Encode()
}
