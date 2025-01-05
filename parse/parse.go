package parse

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/findyourpaths/phil/glr"
	"github.com/kr/pretty"
)

// daysre is a regexp to match day names, either long or short, regardless of case.
var daysRE = regexp.MustCompile(`(?i:\b` + strings.Join(
	[]string{
		"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
		"sun", "mon", "tue", "wed", "thu", "fri", "sat",
	}, `\b|\b`) + `\b,?)`)

// Match a digit on one side and a letter on another. Used to separate `12pm`.
var boundaryRE1 = regexp.MustCompile(`([[:alpha:]])([[:^alpha:]])`)
var boundaryRE2 = regexp.MustCompile(`([[:^alpha:]])([[:alpha:]])`)

func ExtractDateTimeTZRanges(mode, in string) (*DateTimeTZRanges, error) {
	defer func() {
		if err := recover(); err != nil {
			slog.Warn("in ExtractDatetimeRanges(), got a panic trying to extract", "in", in, "err", err)
		}
	}()

	ambiguousDateMode = mode
	// yyDebug = 3
	if yyDebug == 3 {
		fmt.Printf("in before: %q\n", in)
	}
	in = daysRE.ReplaceAllString(in, ``)
	in = boundaryRE1.ReplaceAllString(in, `$1 $2`)
	in = boundaryRE2.ReplaceAllString(in, `$1 $2`)
	in = strings.TrimPrefix(in, "When ")
	in = strings.TrimPrefix(in, "when ")
	in = CleanTextLine(in)
	if yyDebug == 3 {
		fmt.Printf("in after: %q\n", in)
	}

	g := &glr.Grammar{
		Rules:   parseRules,
		Actions: parseActions,
		States:  parseStates,
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
	return glr.GetParseNodeValue(g, forest[0].Children[1], "").(*DateTimeTZRanges), nil
}
