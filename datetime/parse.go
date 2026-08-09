package datetime

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/findyourpaths/phil/glr"
	"github.com/kr/pretty"
	"github.com/tkuchiki/go-timezone"
)

// yearMonthRE matches ISO 8601 YYYY-MM format (e.g. "2023-02").
var yearMonthRE = regexp.MustCompile(`^\d{4}-\d{2}$`)

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

// ClearParseCache drops cached Parse results. Long-running authoring processes
// use this to isolate exploratory parser calls from later deterministic passes.
func ClearParseCache() {
	cacheMutex.Lock()
	cache = map[string]*DateTimeRanges{}
	cacheMutex.Unlock()
}

// ParseOptions configures parse-time context that the grammar cannot infer
// from the input text alone.
type ParseOptions struct {
	MinDateTime *DateTime
	DateMode    string
	// DefaultLocation anchors otherwise civil results to an IANA location. When
	// nil, Parse preserves civil date/time output instead of inventing UTC.
	DefaultLocation *time.Location
	DefaultYear     int
}

type parserContext struct {
	minimumDateTime *DateTime
	dateMode        string
	defaultLocation *time.Location
	defaultYear     int
}

var parseCtx parserContext

// parseMutex protects parseCtx
// during Parse(). The GLR action system uses reflection to call action functions
// with only their parsed child values — there is no way to thread a context
// parameter through the action function signatures. The mutex serializes parsing
// to ensure the context remains stable for the duration of each Parse() call.
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

func Parse(in string, opts ParseOptions) (rngs *DateTimeRanges, err error) {
	if err := validateDateMode(opts.DateMode); err != nil {
		return nil, err
	}

	// Pre-validation: reject strings that clearly aren't dates before invoking the GLR parser.
	// This prevents panics from the parser trying to interpret random text as timezones.
	if !looksLikeDate(in) {
		return nil, nil
	}

	// Handle ISO 8601 YYYY-MM format directly (e.g. "2023-02") since the GLR
	// parser misinterprets the "-" as a range separator.
	if yearMonthRE.MatchString(strings.TrimSpace(in)) {
		year, _ := strconv.Atoi(in[:4])
		month, _ := strconv.Atoi(in[5:7])
		rs := NewRangesWithStartDates(&Date{Year: year, Month: time.Month(month)})
		if err := stampDateTimeRangesDefaultTZ(rs, opts.DefaultLocation); err != nil {
			return nil, err
		}
		return rs, nil
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			slog.Warn("in Parse(), got a panic trying to extract", "in", in, "err", recovered)
			rngs = nil
			err = fmt.Errorf("phil.Parse: panic parsing %q: %v", in, recovered)
		}
	}()

	key := parseCacheKey(in, opts)
	cacheMutex.RLock()
	r, found := cache[key]
	cacheMutex.RUnlock()
	if found {
		return r.Clone(), nil
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
	parseCtx = parserContext{
		minimumDateTime: opts.MinDateTime,
		dateMode:        normalizedDateMode(opts.DateMode),
		defaultLocation: opts.DefaultLocation,
		defaultYear:     opts.DefaultYear,
	}
	defer func() {
		parseCtx = parserContext{}
		parseMutex.Unlock()
	}()

	debugf("minimumDateTime: %q\n", parseCtx.minimumDateTime.String())
	debugf("parseDateMode: %q\n", parseCtx.dateMode)

	g := &glr.Grammar{
		Rules:   datetimeRules,
		Actions: datetimeActions,
		States:  datetimeStates,
	}
	lexer := NewDatetimeLexer(in)
	if lexer.root != nil {
		return finalizeParseResult(lexer.root, lexer.recurrence, opts, key)
	}
	roots, err := glr.Parse(g, lexer)
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
	for i, root := range roots {
		var err error
		rsAny, err = glr.GetParseNodeValue(g, root, "")
		debugf("root[%d]: score=%s err=%v rsAny=%v\n", i, root, err, rsAny)
		if err == nil {
			break
		}
	}

	var rs *DateTimeRanges
	if rsAny == nil {
		return nil, nil
	}

	rs = rsAny.(*DateTimeRanges)
	return finalizeParseResult(rs, lexer.recurrence, opts, key)
}

func finalizeParseResult(rs *DateTimeRanges, recurrence *Recurrence, opts ParseOptions, key string) (*DateTimeRanges, error) {
	// Attach recurrence captured from preprocessing (e.g., stripped plural weekday).
	// Only attach to single-range results — multi-range results are already expanded
	// and the recurrence is just informational (would cause ICS expansion errors).
	if recurrence != nil && rs.Recurrence == nil && len(rs.Items) <= 1 {
		rs.Recurrence = recurrence
	}

	// Resolve missing years using minimumDateTime. Grammar actions like
	// NewRangesWithStartDates create Dates that bypass NewDateTime (and thus
	// setNewDateYear). NewRange propagates year across start/end within a
	// range, but day-list items (DayPlus1 Month without Year) have no range
	// partner to inherit from. Fix them here after all parsing is complete.
	if parseCtx.defaultYear != 0 || parseCtx.minimumDateTime != nil {
		for _, item := range rs.Items {
			if item.Start != nil && item.Start.Date != nil && item.Start.Date.Year == 0 &&
				item.Start.Date.Month != 0 && item.Start.Date.Day != 0 {
				setNewDateYear(item.Start.Date)
			}
			if item.End != nil && item.End.Date != nil && item.End.Date.Year == 0 &&
				item.End.Date.Month != 0 && item.End.Date.Day != 0 {
				setNewDateYear(item.End.Date)
			}
		}
	}

	if err := stampDateTimeRangesDefaultTZ(rs, opts.DefaultLocation); err != nil {
		return nil, err
	}

	// A result containing a backwards range is a failed parse, never a
	// returned value: a schedule is not printed with its end date before its
	// start date, so such a reading is invalid. (Item ORDER across separate
	// items is deliberately not checked: prose legitimately mentions dates
	// out of order — "kicks off March 2nd, applications through February
	// 1st".) Recovery belongs to the caller's general parse-failure handling
	// (e.g. re-presenting the full text through ParseBlock), never to repair
	// here.
	for _, item := range rs.Items {
		if rangeEndsBeforeStart(item) {
			return nil, fmt.Errorf("semantic error: range end %s precedes range start %s", item.End, item.Start)
		}
	}

	debugf("rs: %#v\n", rs)
	cacheMutex.Lock()
	cache[key] = rs.Clone()
	cacheMutex.Unlock()
	return rs, nil
}

// rangeEndsBeforeStart reports whether rng carries complete start and end
// dates with the end date strictly before the start date. Comparison is
// date-level only: an end time earlier than the start time on the same or a
// later date remains valid (overnight ranges).
func rangeEndsBeforeStart(rng *DateTimeRange) bool {
	if rng == nil || rng.Start == nil || rng.End == nil {
		return false
	}
	return dateStrictlyBefore(rng.End.Date, rng.Start.Date)
}

// dateStrictlyBefore reports whether a is strictly before b, comparing only
// when both dates are complete (year, month, and day all set).
func dateStrictlyBefore(a, b *Date) bool {
	if a == nil || b == nil {
		return false
	}
	if a.Year == 0 || a.Month == 0 || a.Day == 0 ||
		b.Year == 0 || b.Month == 0 || b.Day == 0 {
		return false
	}
	if a.Year != b.Year {
		return a.Year < b.Year
	}
	if a.Month != b.Month {
		return a.Month < b.Month
	}
	return a.Day < b.Day
}

func validateDateMode(mode string) error {
	switch mode {
	case DateModeUnknown, DateModeNorthAmerican, DateModeRest:
		return nil
	default:
		return fmt.Errorf("phil.Parse: invalid DateMode %q; want '', 'na', or 'rest'", mode)
	}
}

func normalizedDateMode(mode string) string {
	if mode == DateModeUnknown {
		return DateModeNorthAmerican
	}
	return mode
}

func parseCacheKey(in string, opts ParseOptions) string {
	tzKey := ""
	if opts.DefaultLocation != nil {
		tzKey = opts.DefaultLocation.String()
	}
	yearKey := ""
	if opts.DefaultYear != 0 {
		yearKey = strconv.Itoa(opts.DefaultYear)
	}
	return fmt.Sprintf("%q %q %q %q %q", opts.MinDateTime.String(), opts.DateMode, tzKey, yearKey, in)
}

func stampDateTimeRangesDefaultTZ(rs *DateTimeRanges, loc *time.Location) error {
	if rs == nil {
		return nil
	}
	for _, item := range rs.Items {
		if item == nil {
			continue
		}
		// Meridiem inheritance has already run at range construction; drop the
		// parse-internal ambiguity marker so Parse output is canonical.
		clearMeridiemImplied(item.Start)
		clearMeridiemImplied(item.End)
		if err := stampDefaultTZ(item.Start, loc); err != nil {
			return err
		}
		if err := stampDefaultTZ(item.End, loc); err != nil {
			return err
		}
	}
	return nil
}

func clearMeridiemImplied(dt *DateTime) {
	if dt != nil && dt.Time != nil {
		dt.Time.MeridiemImplied = false
	}
}

func stampDefaultTZ(dt *DateTime, loc *time.Location) error {
	if dt == nil || dt.Date == nil ||
		dt.Date.Year == 0 || dt.Date.Month == 0 || dt.Date.Day == 0 {
		return nil
	}
	if dt.TimeZone != nil && dt.TimeZone.IANAName() != "" {
		return nil
	}
	if loc == nil {
		return nil
	}

	if dt.TimeZone != nil && dt.TimeZone.Offset != "" {
		want := offsetForLocation(loc, dt)
		if !offsetsEqual(dt.TimeZone.Offset, want) {
			return fmt.Errorf("phil.Parse: parsed datetime %q has numeric offset %s, which does not match DefaultLocation %s offset %s",
				dt, dt.TimeZone.Offset, loc.String(), want)
		}
	}

	dt.TimeZone = timeZoneForLocation(loc, dt.Date, dt.Time)
	return nil
}

func offsetForLocation(loc *time.Location, dt *DateTime) string {
	t := timeForLocation(loc, dt.Date, dt.Time)
	_, offsetSec := t.Zone()
	return formatOffset(offsetSec)
}

func offsetsEqual(got, want string) bool {
	if got == "-00:00" {
		got = "+00:00"
	}
	if want == "-00:00" {
		want = "+00:00"
	}
	return got == want
}
