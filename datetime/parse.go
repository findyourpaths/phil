package datetime

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/findyourpaths/phil/glr"
	"github.com/kr/pretty"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tkuchiki/go-timezone"
)

// Debug flags
// var DoDebug = true

var DoDebug = false

// SetDebug toggles debug logging
func SetDebug(enabled bool) {
	DoDebug = enabled
}

// debugf prints debug messages if debug is enabled
func debugf(format string, args ...any) {
	if DoDebug {
		fmt.Printf(format, args...)
	}
}

// daysre is a regexp to match day names, either long or short, regardless of case.
// var daysRE = regexp.MustCompile(`(?i:\b` + strings.Join(
// 	[]string{
// 		"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
// 		"sun", "mon", "tue", "wed", "thu", "fri", "sat",
// 	}, `\b|\b`) + `\b,?)`)

// Match a digit on one side and a letter on another. Used to separate `12pm`.
var boundaryRE1 = regexp.MustCompile(`([[:alpha:]])([[:^alpha:]])`)
var boundaryRE2 = regexp.MustCompile(`([[:^alpha:]])([[:alpha:]])`)

// var spacifyRE = regexp.MustCompile(`\s*\b(.|-)\b\s*`)
var spacifyRE = regexp.MustCompile(`(?:^|\s*\b)(.|-)(?:\b\s*|$)`)

var singletonTZ = timezone.New()

var cache = map[[4]string]*DateTimeTZRanges{}
var cacheMutex sync.RWMutex

func Parse(year int, dateMode, timeZone, in string) (*DateTimeTZRanges, error) {
	defer func() {
		if err := recover(); err != nil {
			slog.Warn("in ExtractDatetimeRanges(), got a panic trying to extract", "in", in, "err", err)
		}
	}()
	key := [4]string{strconv.Itoa(year), dateMode, timeZone, in}
	cacheMutex.RLock()
	r, found := cache[key]
	cacheMutex.RUnlock()
	if found {
		return r, nil
	}

	parseYear = year
	debugf("parseYear: %d\n", parseYear)

	if dateMode == "" {
		if strings.HasPrefix(timeZone, "America/") {
			dateMode = "na"
		} else {
			dateMode = "rest"
		}
	}
	parseDateMode = dateMode
	debugf("parseDateMode: %q\n", parseDateMode)

	if strings.Index(timeZone, "/") > 0 {
		timeZone, _ = singletonTZ.GetTimezoneAbbreviation(timeZone)
	}
	parseTimeZone = timeZone
	debugf("parseTimeZone: %q\n", parseTimeZone)

	// yyDebug = 3
	debugf("in before processing: %q\n", in)
	// in = daysRE.ReplaceAllString(in, ``)
	in = boundaryRE1.ReplaceAllString(in, `$1 $2`)
	in = boundaryRE2.ReplaceAllString(in, `$1 $2`)
	in = spacifyRE.ReplaceAllString(in, ` $1 `)
	// in = strings.TrimPrefix(in, "When ")
	// in = strings.TrimPrefix(in, "when ")
	in = CleanTextLine(in)
	debugf("in after processing: %q\n", in)

	g := &glr.Grammar{
		Rules:   datetimeRules,
		Actions: datetimeActions,
		States:  datetimeStates,
	}
	forest, err := glr.Parse(g, NewDatetimeLexer(in))
	if yyDebug == 3 {
		fmt.Printf("tree:\n%# v\n", pretty.Formatter(forest))
		fmt.Printf("err: %v\n", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse datetime ranges from %q", in)
	}
	if len(forest) == 0 || len(forest[0].Children) == 0 {
		return nil, fmt.Errorf("no datetime ranges found in %q", in)
	}

	rs := glr.GetParseNodeValue(g, forest[0], "").(*DateTimeTZRanges)
	// for _, rng := range rngs.Items {
	// 	if rng.Start.TimeZone == "" {
	// 		rng.Start.TimeZone = timeZone
	// 		// fmt.Printf("set time zone to: %q\n", tz)
	// 	}
	// 	if rng.End != nil && rng.End.TimeZone == "" {
	// 		rng.End.TimeZone = timeZone
	// 	}
	// }
	debugf("rs: %#v\n", rs)
	cacheMutex.Lock()
	cache[key] = rs
	cacheMutex.Unlock()
	return rs, nil
}

var whitespacesRE = regexp.MustCompile(`\s+`)

func CleanTextLine(s string) string {
	r := bluemonday.StrictPolicy().AddSpaceWhenStrippingTag(true).Sanitize(s)
	r = strings.ReplaceAll(r, "\u00a0", " ")
	r = whitespacesRE.ReplaceAllString(r, " ")
	r = strings.TrimSpace(r)
	return r
}
