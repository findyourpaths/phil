package phil_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/findyourpaths/phil/datetime"
	"github.com/findyourpaths/phil/ical"

	"github.com/google/go-cmp/cmp"
)

// ---------------------------------------------------------------------------
// Integration test: A→C pipeline (input string → ICS golden file)
//
// Each golden file in testdata/ is a self-contained test case.
//
// Normal mode: walks testdata/*.ics, parses each file's comment header to
// extract input string, dateMode, minDT, then runs the pipeline and compares
// the output to the file body. The golden files ARE the test table.
//
// UPDATE mode: uses the input list below to regenerate golden files:
//   UPDATE=1 go test . -run TestIntegration -v
//
// To add a test: add an entry to testInputs, run UPDATE=1, verify the golden
// file, commit it. From then on the golden file drives the test.
// ---------------------------------------------------------------------------

// defaultMinDT is used when no minDT comment is present in a golden file.
// All inputs without an explicit year get 2023 from this.
var defaultMinDT = &datetime.DateTime{
	Date:     &datetime.Date{Year: 2023, Month: 1, Day: 1},
	Time:     &datetime.Time{},
	TimeZone: &datetime.TimeZone{Abbreviation: "ET"},
}

// ---------------------------------------------------------------------------
// Golden file comment protocol
// ---------------------------------------------------------------------------
//
// Line 1:        # <input string>               — REQUIRED
// Optional:      # dateMode: rest | na
// Optional:      # minDT: 2023-01-01T00:00 ET   — date, time, timezone
// Optional:      # nil                           — parser returns nil
// Optional:      # broken                        — known broken, skip
// Optional:      # ERROR: <msg>                  — ical conversion error
//
// First non-comment line begins the expected ICS body.

// goldenTestCase holds the data parsed from a golden file's comment header.
type goldenTestCase struct {
	input     string
	dateMode  string
	minDT     *datetime.DateTime // nil = use defaultMinDT
	wantNil   bool
	broken    bool
	wantError string // non-empty if "# ERROR: ..."
	body      string // everything after the comment header
}

// isBodyMarker returns true if a line marks the start of the golden file body.
func isBodyMarker(line string) bool {
	return strings.HasPrefix(line, "BEGIN:") ||
		line == "# nil" ||
		line == "# broken" ||
		strings.HasPrefix(line, "# ERROR: ")
}

// parseGoldenFile reads a golden .ics file and extracts the test case from
// its comment header.
//
// The file has two sections:
//   - Header: input string (line 1), optional directives (dateMode, minDT)
//   - Body: everything from the first body marker onward (BEGIN:VCALENDAR,
//     "# nil", "# broken", or "# ERROR: ...")
//
// The input string may span multiple lines (e.g., embedded \n). Lines between
// the first comment and the first body marker that aren't directives are
// treated as input continuations.
func parseGoldenFile(path string) (*goldenTestCase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")

	// Find the body start: first line that is a body marker.
	bodyStart := len(lines)
	for i, line := range lines {
		if i == 0 {
			continue // line 0 is always the input comment
		}
		if isBodyMarker(line) {
			bodyStart = i
			break
		}
	}

	tc := &goldenTestCase{}

	// Parse the header (lines 0..bodyStart-1).
	for i := 0; i < bodyStart && i < len(lines); i++ {
		line := lines[i]

		if i == 0 {
			tc.input = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			continue
		}

		if !strings.HasPrefix(line, "# ") {
			// Non-comment line before body marker: input continuation.
			tc.input += "\n" + line
			continue
		}

		content := strings.TrimSpace(strings.TrimPrefix(line, "# "))
		switch {
		case strings.HasPrefix(content, "dateMode: "):
			mode := strings.TrimPrefix(content, "dateMode: ")
			switch mode {
			case "rest":
				tc.dateMode = datetime.DateModeRest
			case "na":
				tc.dateMode = datetime.DateModeNorthAmerican
			}
		case strings.HasPrefix(content, "minDT: "):
			val := strings.TrimPrefix(content, "minDT: ")
			dt, parseErr := parseMinDT(val)
			if parseErr != nil {
				return nil, fmt.Errorf("parsing minDT in %s: %w", path, parseErr)
			}
			tc.minDT = dt
		default:
			// Unknown comment: treat as input continuation.
			tc.input += "\n" + content
		}
	}

	// Parse the body.
	if bodyStart < len(lines) {
		bodyLine := lines[bodyStart]
		switch {
		case bodyLine == "# nil":
			tc.wantNil = true
		case bodyLine == "# broken":
			tc.broken = true
		default:
			// ICS content or ERROR line — store as body for comparison.
			tc.body = strings.Join(lines[bodyStart:], "\n")
		}
	}

	return tc, nil
}

// parseMinDT parses a minDT string like "2023-01-01T00:00 ET" or "2023-02-03T09:00 ET".
func parseMinDT(s string) (*datetime.DateTime, error) {
	parts := strings.Fields(s)
	if len(parts) < 1 {
		return nil, fmt.Errorf("empty minDT")
	}

	datePart := parts[0]
	tz := "ET" // default timezone
	if len(parts) >= 2 {
		tz = parts[1]
	}

	var year, month, day, hour, minute int
	if idx := strings.Index(datePart, "T"); idx >= 0 {
		// Has time component: 2023-02-03T09:00
		dateStr := datePart[:idx]
		timeStr := datePart[idx+1:]

		dateParts := strings.Split(dateStr, "-")
		if len(dateParts) != 3 {
			return nil, fmt.Errorf("bad date format: %q", dateStr)
		}
		year, _ = strconv.Atoi(dateParts[0])
		month, _ = strconv.Atoi(dateParts[1])
		day, _ = strconv.Atoi(dateParts[2])

		timeParts := strings.Split(timeStr, ":")
		if len(timeParts) >= 1 {
			hour, _ = strconv.Atoi(timeParts[0])
		}
		if len(timeParts) >= 2 {
			minute, _ = strconv.Atoi(timeParts[1])
		}
	} else {
		// Date only: 2023-01-01
		dateParts := strings.Split(datePart, "-")
		if len(dateParts) != 3 {
			return nil, fmt.Errorf("bad date format: %q", datePart)
		}
		year, _ = strconv.Atoi(dateParts[0])
		month, _ = strconv.Atoi(dateParts[1])
		day, _ = strconv.Atoi(dateParts[2])
	}

	return &datetime.DateTime{
		Date:     &datetime.Date{Year: year, Month: time.Month(month), Day: day},
		Time:     &datetime.Time{Hour: hour, Minute: minute},
		TimeZone: &datetime.TimeZone{Abbreviation: tz},
	}, nil
}

func parseOptions(input string, minDT *datetime.DateTime, dateMode string) datetime.ParseOptions {
	var loc *time.Location
	if minDT != nil && minDT.TimeZone != nil {
		if name := minDT.TimeZone.IANAName(); name != "" {
			loaded, err := time.LoadLocation(name)
			if err == nil {
				loc = loaded
			}
		}
	}
	if strings.Contains(input, "+00:00") || strings.Contains(input, "-00:00") {
		loc = time.UTC
	}
	defaultYear := 0
	if minDT != nil && minDT.Date != nil {
		defaultYear = minDT.Date.Year
	}
	return datetime.ParseOptions{
		MinDateTime:     minDT,
		DateMode:        dateMode,
		DefaultLocation: loc,
		DefaultYear:     defaultYear,
	}
}

// ---------------------------------------------------------------------------
// ICS normalization
// ---------------------------------------------------------------------------

var uidLineRE = regexp.MustCompile(`(?m)^UID:.*\n`)
var dtstampLineRE = regexp.MustCompile(`(?m)^DTSTAMP:.*\n`)

func normalizeICS(s string) string {
	s = uidLineRE.ReplaceAllString(s, "")
	s = dtstampLineRE.ReplaceAllString(s, "")
	return s
}

// ---------------------------------------------------------------------------
// Golden file naming (used only in UPDATE mode)
// ---------------------------------------------------------------------------

var slugRE = regexp.MustCompile(`[^a-z0-9_-]+`)

type testInput struct {
	in       string
	dateMode string
	minDT    *datetime.DateTime
	wantNil  bool
	isBroken bool
}

func goldenName(i int, tc testInput) string {
	s := strings.ToLower(tc.in)
	s = strings.ReplaceAll(s, " ", "_")
	s = slugRE.ReplaceAllString(s, "")
	if len(s) > 100 {
		s = s[:100]
	}
	name := fmt.Sprintf("%03d_%s", i, s)
	if tc.dateMode == datetime.DateModeNorthAmerican {
		name += "_na"
	} else if tc.dateMode == datetime.DateModeRest {
		name += "_rest"
	}
	if tc.minDT != nil {
		name += "_mindt"
	}
	return name
}

// buildCommentHeader builds the self-contained comment header for a golden file.
func buildCommentHeader(tc testInput) string {
	var lines []string
	lines = append(lines, "# "+tc.in)
	if tc.dateMode == datetime.DateModeRest {
		lines = append(lines, "# dateMode: rest")
	} else if tc.dateMode == datetime.DateModeNorthAmerican {
		lines = append(lines, "# dateMode: na")
	}
	if tc.minDT != nil {
		d := tc.minDT.Date
		t := tc.minDT.Time
		tz := "ET"
		if tc.minDT.TimeZone != nil && tc.minDT.TimeZone.Abbreviation != "" {
			tz = tc.minDT.TimeZone.Abbreviation
		}
		lines = append(lines, fmt.Sprintf("# minDT: %04d-%02d-%02dT%02d:%02d %s",
			d.Year, d.Month, d.Day, t.Hour, t.Minute, tz))
	}
	return strings.Join(lines, "\n") + "\n"
}

// ---------------------------------------------------------------------------
// TestIntegration_ParseToICS
// ---------------------------------------------------------------------------

func TestIntegration_ParseToICS(t *testing.T) {
	update := os.Getenv("UPDATE") == "1"

	if update {
		generateGoldenFiles(t)
		return
	}

	// Normal mode: golden files drive the test.
	info := &ical.EventInfo{Summary: "Test Event"}

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata/: %v", err)
	}

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".ics") {
			continue
		}
		e := e // capture loop var
		t.Run(e.Name(), func(t *testing.T) {
			path := filepath.Join("testdata", e.Name())
			tc, err := parseGoldenFile(path)
			if err != nil {
				t.Fatalf("failed to parse golden file: %v", err)
			}

			if tc.broken {
				t.Skip("broken test")
				return
			}

			// Resolve minDT: use parsed value or default.
			minDT := tc.minDT
			if minDT == nil {
				minDT = defaultMinDT
			}

			parsed, parseErr := datetime.Parse(tc.input, parseOptions(tc.input, minDT, tc.dateMode))
			if parseErr != nil {
				t.Fatalf("Parse error: %v", parseErr)
			}

			if tc.wantNil {
				if parsed != nil {
					t.Fatalf("expected nil parse result for %q, got %v", tc.input, parsed)
				}
				return
			}
			if parsed == nil {
				t.Fatalf("unexpected nil parse result for %q", tc.input)
			}

			cal, calErr := ical.NewCalendar(parsed, info)

			var gotBody string
			if calErr != nil {
				gotBody = "# ERROR: " + calErr.Error() + "\n"
			} else {
				gotBody = normalizeICS(ical.ICS(cal))
			}

			if diff := cmp.Diff(tc.body, gotBody); diff != "" {
				t.Errorf("ICS mismatch for %q (-want +got):\n%s", tc.input, diff)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// UPDATE mode: generate golden files from test inputs
// ---------------------------------------------------------------------------

func generateGoldenFiles(t *testing.T) {
	t.Helper()

	etTZ := &datetime.TimeZone{Abbreviation: "ET"}
	minDT_2023Feb01_09AM := &datetime.DateTime{
		Date:     &datetime.Date{Year: 2023, Month: 2, Day: 1},
		Time:     &datetime.Time{Hour: 9},
		TimeZone: etTZ,
	}
	minDT_2023Feb03_09AM := &datetime.DateTime{
		Date:     &datetime.Date{Year: 2023, Month: 2, Day: 3},
		Time:     &datetime.Time{Hour: 9},
		TimeZone: etTZ,
	}

	inputs := []testInput{
		// ISO 8601 format
		{in: "2023"},
		{in: "2023-02"},
		{in: "2023-02-03"},
		{in: "2023-02-03T"},
		{in: "2023-02-03T12:00"},
		{in: "2023-02-03T12:00:00"},
		{in: "2023-02-03T12:00:00Z"},
		{in: "2023-02-03T12:00:00+00:00"},
		{in: "2023-02-03T12:00:00-00:00"},
		{in: "2023-02-03T12:00:00-05:00"},

		// Date - MD
		{in: "Feb 3"},
		{in: "February 3"},
		{in: "Fri Feb 3"},
		{in: "Fri 3 Feb"},

		// Date - DM
		{in: "3 Feb"},
		{in: "3, Feb"},
		{in: "3rd of Feb"},
		{in: "fri 3 Feb"},

		// Date - MDY
		{in: "Feb 3 2023"},
		{in: "February 3 2023"},
		{in: "February 3, 2023"},
		{in: "February 3rd, 2023"},
		{in: "Fri Feb 3, 2023"},
		{in: "Friday Feb 3rd 2023"},
		{in: "Time:Feb 3 2023"},

		// Date - DMY
		{in: "3 Feb 2023"},
		{in: "3rd Feb 2023"},
		{in: "3 February, 2023"},
		{in: "Friday 3rd Feb 2023"},

		// Date - MY
		{in: "Feb 2023"},

		// Date - Both
		{in: "02.03", dateMode: datetime.DateModeNorthAmerican},
		{in: "02.03", dateMode: datetime.DateModeRest},
		{in: "02.03.", dateMode: datetime.DateModeNorthAmerican},
		{in: "02.03.", dateMode: datetime.DateModeRest},
		{in: "2/3/2023", dateMode: datetime.DateModeNorthAmerican},
		{in: "2/3/2023", dateMode: datetime.DateModeRest},

		// Extra tokens
		{in: "Feb 3 Google Calendar ICS"},
		{in: "Updated: Feb 3"},
		{in: "Workshop Update (2/3/23)", dateMode: datetime.DateModeNorthAmerican},
		{in: "Workshop: Feb 3 2023  VIRTUAL"},
		{in: "Release date: February 3, 2023"},

		// Dates - MD
		{in: "Feb 1, 2"},
		{in: "Feb 1, 2, 3"},
		{in: "Feb 1, 2, 3, 4"},
		{in: "Feb 1, 2, 3, 4, 5"},
		{in: "Feb 3 Mar 2"},
		{in: "Our next cohort kicks off on March 2nd and we're accepting applications through February 1st."},

		// Dates - DM
		{in: "1, 2 Feb"},
		{in: "1, 2, 3 Feb"},
		{in: "1, 2, 3, 4 Feb"},
		{in: "1, 2, 3, 4, 5 Feb"},
		{in: "1, 2, 3 Feb and 2 Mar"},
		{in: "1-3 Feb and 2 Mar", isBroken: true},
		{in: "1-3 & 5 February", isBroken: true},
		{in: "1-4 & 6 February", isBroken: true},

		// Dates - MDY
		{in: "Feb 1, 2 2023"},
		{in: "Feb 1, 2, 3 2023"},
		{in: "Feb 1, 2, 3, 4 2023"},
		{in: "Feb 1, 2, 3, 4, 5 2023"},
		{in: "Feb 1, 2, 3 and Mar 2 2023"},
		{in: "Feb 3 Mar 2 2023"},

		// Dates - DMY
		{in: "1, 2 Feb 2023"},
		{in: "1, 2, 3 Feb 2023"},
		{in: "1, 2, 3, 4 Feb 2023"},
		{in: "1, 2, 3, 4, 5 Feb 2023"},
		{in: "1, 2, 3 Feb and 2 Mar 2023"},

		// Date Range - MD
		{in: "Feb 3rd-4th"},
		{in: "Feb 3 - Mar 2"},
		{in: "Feb 3 to Mar 2"},
		{in: "February 2 - 5 (TH-SU)"},

		// Date Range - DM
		{in: "3-4 Feb"},
		{in: "3 Feb - 4 Feb"},
		{in: "Fri Feb 3 - Sat Feb 4"},
		{in: "3 February - 2 March"},

		// Date Range - Various separators
		{in: "Feb 3-4"},
		{in: "3--4 Feb"},
		{in: "3 - 4 Feb"},
		{in: "3 -- 4 Feb"},
		{in: "3 to 4 Feb"},
		{in: "3 until 4 Feb"},
		{in: "3 through 4 Feb"},
		{in: "3 \u2013 4 Feb"},
		{in: "3 \u2014 4 Feb"},
		{in: "3-> 4 Feb"},
		{in: "From 3 - 4 Feb"},
		{in: "from 3rd till 4th of Feb"},

		// Date Range - MDY
		{in: "Feb 3-4 2023"},
		{in: "Feb 3 - 4 2023"},
		{in: "Feb 3 2023 - Feb 4 2023"},
		{in: "Fri Feb 3, 2023 - Sat Feb 4, 2023"},
		{in: "Fri 3 Feb - Sat 4 February 2023"},
		{in: "Fri 3rd Feb - Sat 4th February 2023"},
		{in: "Fri Feb 3rd - Sat 4th February 2023"},
		{in: "Fri 3rd Feb - 4th Sat February 2023"},
		{in: "Fri Feb 3rd - 4th Sat February 2023"},
		{in: "February 3 - March 2, 2023"},
		{in: "SAVE THE DATES: Feb 3-4, 2023"},

		// Date Range - DMY
		{in: "3-4 Feb 2023"},
		{in: "3-4 Feb. 2023"},
		{in: "3-4 February, 2023"},
		{in: "3rd-4th Feb 2023"},
		{in: "3 Feb 2023 - 4 Feb 2023"},
		{in: "From 3rd to 4th, Feb 2023"},
		{in: "beginning 3rd to 4th Feb 2023"},

		// Date Range - M
		{in: "Feb - Mar"},
		// Date Range - Y
		{in: "2023 - 2024"},
		// Date Range - MY
		{in: "Feb 2023 - Mar 2023"},
		// Date Range - YMD
		{in: "2023, Feb 3 - 2023, Feb 4"},

		// Date Ranges
		{in: "Feb 1-2, 3-4"},
		{in: "Feb 1-2, 3-4 2023"},
		{in: "Feb 1-2; Mar 2-3"},
		{in: "2/1, 2/2, 3/2, 3/3", dateMode: datetime.DateModeNorthAmerican},
		{in: "1/2, 2/2, 2/3, 3/3", dateMode: datetime.DateModeRest},
		{in: "5 Wednesdays 2/1 & 2/8 & 2/15 & 2/22, 3/2"},

		// Date Ranges - DM
		{in: "1-2 Feb; 2-3 Mar"},

		// Date Ranges - MDY
		{in: "Feb 1-2; Mar 2-3 2023"},
		// Date Ranges - DMY
		{in: "1-2 Feb; 2-3 Mar 2023"},
		{in: "Part 1: 1st-2nd February 2023"},

		// Date Time - Relative
		{in: "Join today at 12pm", minDT: minDT_2023Feb03_09AM},
		{in: "Join today for Day 2 at 12pm", minDT: minDT_2023Feb03_09AM},
		{in: "Tomorrow at 12pm", minDT: minDT_2023Feb03_09AM},

		// Date Time - Relative Z
		{in: "Today Friday, 12pm ET", minDT: minDT_2023Feb03_09AM},

		// Date Time - MDT
		{in: "Feb 3 12pm"},
		{in: "Feb 3 12:00 PM"},
		{in: "February 3 @ 12:00 PM"},
		{in: "February 3  12 p.m."},
		{in: "Date:Fri 03 Feb, Time:12pm"},
		{in: "Starting February 3rd at 12pm"},

		// Date Time - WMDT
		{in: "Friday 2/14: **Love is Listening and Art: Social + Listening Art Sessions**", dateMode: datetime.DateModeNorthAmerican},
		{in: "Friday 2/14: **Love is Listening and Art: Social + Listening Art Sessions** at 6pm facilitated by Lauren V", dateMode: datetime.DateModeNorthAmerican},

		// Date Time - WMDTZ
		{in: "FRI 2/3 @ 12:00pm ET"},

		// Date Time - MDTZ
		{in: "Feb 3 12pm ET"},
		{in: "Feb 3 12pm (ET)"},
		{in: "Feb 3 12pm - ET"},
		{in: "Feb 3 12pm in ET"},
		{in: "Feb 3 12pm Eastern"},
		{in: "Feb 3 12pm US/Eastern"},
		{in: "FEBRUARY 3RD 12 PM ET, ON FRIDAY"},
		{in: "Starting February 3rd at 12pm (ET) - Virtually."},
		{in: "Starts Friday 2/3 at 9:00 am ET", dateMode: datetime.DateModeNorthAmerican},

		// Date Time - DMT
		{in: "Date:Fri 03 Feb, Time:3.00pm"},
		{in: "Friday 3 Feb 3:00pm (doors) | 11pm (curfew)"},
		{in: "Fri, 02.03.2023 - 15:00"},
		{in: "Thu, 02.03.2023 - 15:00"},
		{in: "Fri, 03.02.2023 - 15:00"},
		{in: "Thu, 03.02.2023 - 15:00"},

		// Date Time - MDYT
		{in: "Feb. 3, 2023 12:00pm"},
		{in: "Feb 3, 2023 @ 12:00 PM"},
		{in: "Friday, February 3rd 2023 from 12:00 PM"},

		// Date Time - MDYTT
		{in: "Feb. 3, 2023 12:00pm, 3:00pm"},

		// Date Time - MDYTZ
		{in: "Feb 3 2023 12pm ET"},
		{in: "Feb 3, 2023 12:00 PM Eastern Time (US and Canada)"},

		// Date Time - DMY
		{in: "3rd Feb 2023 9:00"},
		{in: "3rd Feb 2023 9:00am"},
		{in: "3rd Feb 2023 3:00pm"},

		// Date Time - TZMD
		{in: "12:00 pm ET February 3rd"},

		// Date Time Ranges - MD-MDY
		{in: "Date: Friday, February 3rd - February 4th, 2023"},
		{in: "Date: Friday, February 3rd - February 4th, 2023 ET"},

		// Date Time Ranges - MDTT
		{in: "February 3: 9am - 12pm"},
		{in: "Feb 3 9am - 12pm"},
		{in: "Feb 3 @ 9:00 AM - Feb 3 @ 12:00 PM"},
		{in: "February, 3 9:00 - 15:00"},
		{in: "Friday, February 3rd from 12 - 3pm"},
		{in: "Feb, 3rd from 9 am-3.00 pm"},
		{in: "February 3 + 4, 9 am - 12 pm each day"},
		{in: "February 1 + 2: In-person at SeekHealing Asheville,12 pm - 3 pm each day"},
		{in: "THIS Friday: February 3 \n 12-3:00pm"},
		{in: "2 Wednesdays Feb 1st & 8th 12:00p-3:00p"},
		{in: "5 Wednesdays 9:00am-12:00pm February 1st - March 1st", isBroken: true},

		// Date Time Ranges - MDTTZ
		{in: "Feb 3rd - 9.00 AM- 12pm ET"},
		{in: "February 3rd, 9-12pm ET"},
		{in: "Feb 3 2023 9am - 12pm ET"},
		{in: "Feb 3 2023 9am ET to 12pm ET"},
		{in: "Feb 3 @ 9:00 AM ET - Feb 3 @ 12:00 PM ET"},
		{in: "Feb 3, 2023, 9:00 AM ET - Feb 3, 2023, 12:00 PM ET"},
		{in: "Friday, 2/3 12-3pm ET", dateMode: datetime.DateModeNorthAmerican},
		{in: "Friday 2/3 from 12:00 - 3 pm ET"},
		{in: "Friday, February 3, 9am - 12pm Eastern Time"},
		{in: "February 3, 2023 from 9:00 am to noon ET"},
		{in: "February 3, 2023 / 9:00 AM - 12:00 PM ET"},
		{in: "February 3rd, 12:00-3:00pm Eastern (New York) time"},
		{in: "February 3rd & 4th, 9:00 am - noon Eastern time"},
		{in: "February 3rd - 5th, 9:00 am - noon ET each day", isBroken: true},
		{in: "Wednesdays February 1st & 8th 12:00p-3:00p"},
		{in: "Wednesdays, February 1st, 8th, and 15th 9:00am - 12:00pm (ET)", isBroken: true},
		{in: "Wednesdays - February 1, 8 9:00 AM - 12:00 PM ET"},
		{in: "Wednesdays - February 1, 8, 15, 22, and March 1 9:00 AM - 12:00 PM ET", isBroken: true},

		// Date Time Ranges - DTT
		{in: "Friday 12 to 3 PM Eastern", minDT: minDT_2023Feb01_09AM},

		// Date Time Ranges - DMTT
		{in: "3 Feb 9am - 12pm"},

		// Date Time Ranges - MDYT
		{in: "Feb 3 2023 12pm"},
		// Date Time Ranges - MDYTT
		{in: "Friday, February 3, 2023 9:00 AM 12:00 PM"},
		{in: "Friday, February 3rd 2023 from 9:00 AM to 12:00 PM"},
		// Date Time Ranges - DMYTT
		{in: "When 3 Feb 2023 9:00 AM - 12:00 PM"},
		{in: "Fr. 3. Feb. 2023, 9:00-ca.12:00"},

		// Date Time Ranges - TDMY
		{in: "9:00am 3rd Feb - 4th Feb 3:00pm 2023"},
		{in: "9:00am on 3rd Feb - 4th Feb at 3:00pm 2023"},

		// Date Time Ranges - Both
		{in: "02.03.2023", dateMode: datetime.DateModeNorthAmerican},
		{in: "02.03.2023", dateMode: datetime.DateModeRest},
		{in: "02.03.2023 - 15:00", dateMode: datetime.DateModeNorthAmerican},
		{in: "02.03.2023 - 15:00", dateMode: datetime.DateModeRest},
		{in: "Th , 03.02.2023 - 15:00", dateMode: datetime.DateModeNorthAmerican},
		{in: "Fr , 02.03.2023 - 15:00", dateMode: datetime.DateModeNorthAmerican},
		{in: "Th , 02.03.2023 - 15:00", dateMode: datetime.DateModeRest},
		{in: "Fr , 03.02.2023 - 15:00", dateMode: datetime.DateModeRest},
		{in: "Th , 03.02.2023 - 15:00"},
		{in: "Fr , 02.03.2023 - 15:00"},

		// Failures (parse returns nil)
		{in: "814-555-1212", wantNil: true},
		{in: "814-555-1212 x123", wantNil: true},
		{in: "102 W. Mahoning Street. Punxsutawney, PA 15767", wantNil: true},
		{in: "We may request cookies to be set on your device.", wantNil: true},
		{in: "Winter Retreat for 6-12th graders!", wantNil: true},
		{in: "For 6th-12th grade students @ SpringHill Camp", wantNil: true},
	}

	info := &ical.EventInfo{Summary: "Test Event"}

	// Track generated files for orphan detection.
	generatedFiles := make(map[string]bool)

	for i, tc := range inputs {
		name := goldenName(i, tc)
		golden := filepath.Join("testdata", name+".ics")
		generatedFiles[name+".ics"] = true

		// Apply default minDT.
		minDT := tc.minDT
		if minDT == nil {
			minDT = defaultMinDT
		}

		header := buildCommentHeader(tc)

		if tc.isBroken {
			content := header + "# broken\n"
			if err := os.WriteFile(golden, []byte(content), 0o644); err != nil {
				t.Fatalf("failed to write golden file: %v", err)
			}
			continue
		}

		parsed, err := datetime.Parse(tc.in, parseOptions(tc.in, minDT, tc.dateMode))
		if err != nil {
			t.Fatalf("[%d] Parse error for %q: %v", i, tc.in, err)
		}

		if tc.wantNil {
			if parsed != nil {
				t.Fatalf("[%d] expected nil parse result for %q", i, tc.in)
			}
			content := header + "# nil\n"
			if err := os.WriteFile(golden, []byte(content), 0o644); err != nil {
				t.Fatalf("failed to write golden file: %v", err)
			}
			continue
		}
		if parsed == nil {
			t.Fatalf("[%d] unexpected nil parse result for %q", i, tc.in)
		}

		cal, calErr := ical.NewCalendar(parsed, info)
		var body string
		if calErr != nil {
			body = "# ERROR: " + calErr.Error() + "\n"
		} else {
			body = normalizeICS(ical.ICS(cal))
		}

		content := header + body
		if err := os.WriteFile(golden, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to write golden file: %v", err)
		}
	}

	// Warn about orphaned golden files.
	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata/: %v", err)
	}
	var orphaned []string
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".ics") {
			continue
		}
		if !generatedFiles[e.Name()] {
			orphaned = append(orphaned, e.Name())
		}
	}
	if len(orphaned) > 0 {
		sort.Strings(orphaned)
		t.Logf("WARNING: orphaned golden files (not in input list):\n  %s", strings.Join(orphaned, "\n  "))
	}

	t.Logf("Generated %d golden files", len(inputs))
}
