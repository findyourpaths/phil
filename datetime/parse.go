package datetime

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/findyourpaths/phil/glr"
	"github.com/kr/pretty"
	"github.com/tkuchiki/go-timezone"
)

// DoDebug controls debug logging. Use atomic access for goroutine safety.
var DoDebug atomic.Bool

// SetDebug toggles debug logging.
func SetDebug(enabled bool) {
	DoDebug.Store(enabled)
}

// debugf prints debug messages if debug is enabled.
func debugf(format string, args ...any) {
	if DoDebug.Load() {
		fmt.Printf(format, args...)
	}
}

// daysre is a regexp to match day names, either long or short, regardless of case.
// var daysRE = regexp.MustCompile(`(?i:\b` + strings.Join(
// 	[]string{
// 		"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
// 		"sun", "mon", "tue", "wed", "thu", "fri", "sat",
// 	}, `\b|\b`) + `\b,?)`)

var timezoneTZ = timezone.New()

var cache = map[string]*DateTimeRanges{}
var cacheMutex sync.RWMutex

// parseMutex protects the global variables minimumDateTime and parseDateMode
// during Parse(). The GLR action system uses reflection to call action functions
// with only their parsed child values — there is no way to thread a context
// parameter through the action function signatures. The mutex serializes parsing
// to ensure these globals remain stable for the duration of each Parse() call.
// This is acceptable because Parse() is CPU-only (microseconds per call).
var parseMutex sync.Mutex

// looksLikeDate performs a quick pre-validation check to reject strings that clearly
// aren't dates before invoking the GLR parser. This prevents panics from the parser
// trying to interpret random text (like titles or venue names) as dates.
//
// Requirements for a date-like string:
// 1. Must contain at least one digit OR a month name
//
// Examples that should pass:
//   - "Saturday, August 2, 2025 from 11am to 2pm"
//   - "2025-01-15"
//   - "January 15, 2025"
//   - "Feb - Mar"
//
// Examples that should fail:
//   - "Point Four Themes in Four Films" (no digits, no month names)
//   - "" (empty string)
func looksLikeDate(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Check for digits - most dates have numbers
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}

	// Check for month names (handles cases like "Feb - Mar")
	lower := strings.ToLower(s)
	months := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	for _, m := range months {
		if strings.Contains(lower, m) {
			return true
		}
	}

	return false
}

func Parse(minDateTime *DateTime, dateMode string, in string) (*DateTimeRanges, error) {
	// fmt.Printf("datetime.Parse(year: %d, dateMode: %q, timeZone: %#v, in: %q)\n", year, dateMode, timeZone, in)

	// Pre-validation: reject strings that clearly aren't dates before invoking the GLR parser.
	// This prevents panics from the parser trying to interpret random text as timezones.
	if !looksLikeDate(in) {
		return nil, nil
	}

	defer func() {
		if err := recover(); err != nil {
			slog.Warn("in Parse(), got a panic trying to extract", "in", in, "err", err)
		}
	}()

	key := fmt.Sprintf("%q %q %q", minDateTime.String(), dateMode, in)
	cacheMutex.RLock()
	r, found := cache[key]
	cacheMutex.RUnlock()
	if found {
		return r, nil
	}

	// if minDT == nil {
	// 	// Set this so we don't get null pointer references for year and time zone.
	// 	minDT = &DateTime{}
	// }

	// parseYear = year
	// if parseYear == 0 {
	// 	parseYear = time.Now().Year()
	// }

	// Lock to protect global variables (minimumDateTime, parseDateMode) that are
	// read by functions called during GLR parsing via reflection-based actions.
	parseMutex.Lock()
	defer parseMutex.Unlock()

	minimumDateTime = minDateTime
	debugf("minimumDateTime: %q\n", minimumDateTime.String())
	parseDateMode = dateMode
	debugf("parseDateMode: %q\n", parseDateMode)

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
		return nil, nil // fmt.Errorf("no datetime ranges found in %q", in)
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

	var rs *DateTimeRanges
	if rsAny == nil {
		return nil, nil //fmt.Errorf("no parse tree passed semantic checks")
	}

	rs = rsAny.(*DateTimeRanges)
	// fmt.Printf("in datetime.Parse(), rs: %#v\n", rs)
	debugf("rs: %#v\n", rs)
	cacheMutex.Lock()
	cache[key] = rs
	cacheMutex.Unlock()
	return rs, nil
}
