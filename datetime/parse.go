package datetime

import (
	"fmt"
	"log/slog"
	"regexp"
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
		// log.Printf(format, args...)
	}
}

// daysre is a regexp to match day names, either long or short, regardless of case.
// var daysRE = regexp.MustCompile(`(?i:\b` + strings.Join(
// 	[]string{
// 		"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
// 		"sun", "mon", "tue", "wed", "thu", "fri", "sat",
// 	}, `\b|\b`) + `\b,?)`)

var timezoneTZ = timezone.New()

var cache = map[string]*DateTimeTZRanges{}
var cacheMutex sync.RWMutex

func Parse(year int, dateMode string, timeZone *TimeZone, in string) (*DateTimeTZRanges, error) {
	// fmt.Printf("datetime.Parse(year: %d, dateMode: %q, timeZone: %#v, in: %q)\n", year, dateMode, timeZone, in)
	defer func() {
		if err := recover(); err != nil {
			slog.Warn("in Parse(), got a panic trying to extract", "in", in, "err", err)
		}
	}()

	key := fmt.Sprintf("%d %q %#v %q", year, dateMode, timeZone, in)
	cacheMutex.RLock()
	r, found := cache[key]
	cacheMutex.RUnlock()
	if found {
		return r, nil
	}

	parseYear = year
	debugf("parseYear: %d\n", parseYear)

	if dateMode == "" {
		if timeZone != nil && strings.HasPrefix(timeZone.Name, "America/") {
			dateMode = "na"
		} else {
			dateMode = "rest"
		}
	}
	parseDateMode = dateMode
	debugf("parseDateMode: %q\n", parseDateMode)

	parseTimeZone = timeZone
	debugf("parseTimeZone: %q\n", parseTimeZone)

	g := &glr.Grammar{
		Rules:   datetimeRules,
		Actions: datetimeActions,
		States:  datetimeStates,
	}
	roots, err := glr.Parse(g, NewDatetimeLexer(in))
	if yyDebug == 3 {
		fmt.Printf("tree:\n%# v\n", pretty.Formatter(roots))
		fmt.Printf("err: %v\n", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse datetime ranges from %q", in)
	}
	if len(roots) == 0 || len(roots[0].Children) == 0 {
		return nil, fmt.Errorf("no datetime ranges found in %q", in)
	}

	var rsAny any
	for _, root := range roots {
		var err error
		rsAny, err = glr.GetParseNodeValue(g, root, "")
		if err == nil {
			break
		}
		// return nil, err
	}

	// for _, rng := range rngs.Items {
	// 	if rng.Start.TimeZone == "" {
	// 		rng.Start.TimeZone = timeZone
	// 		// fmt.Printf("set time zone to: %q\n", tz)
	// 	}
	// 	if rng.End != nil && rng.End.TimeZone == "" {
	// 		rng.End.TimeZone = timeZone
	// 	}
	// }

	var rs *DateTimeTZRanges
	if rsAny == nil {
		return nil, fmt.Errorf("no parse tree passed semantic checks")
	}

	rs = rsAny.(*DateTimeTZRanges)
	// fmt.Printf("in datetime.Parse(), rs: %#q\n", rs.String())
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
