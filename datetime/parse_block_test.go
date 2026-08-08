package datetime

import (
	"testing"
	"time"
)

// ifs1494Block is the real IFS Institute training #1494 detail-page
// schedule (English/Spanish bilingual, two labeling styles: "Session N" and
// "Bridge day"). Five complete date/date-range expressions.
const ifs1494Block = "Session 1 / sesión 1: August/Agosto 5-8, 2026 (online/en línea)\n" +
	"    Bridge day / Día Puente: Agosto 29 de 2026 (online/en línea)\n" +
	"    Session 2 / Sesión 2: September/septiembre 9-12, 2026 (online/en línea)\n" +
	"    Bridge day / Día Puente: October/Octubre 2, 2026 (online/en línea)\n" +
	"    Session 3 / Sesión 3: October/Octubre 14-17, 2026 (ONSITE/PRESENCIAL)"

const wantIfs1494String = "2026-08-05 - 2026-08-08, 2026-08-29, 2026-09-09 - 2026-09-12, 2026-10-02, 2026-10-14 - 2026-10-17"

// ifs1463Block is the real IFS Institute training #1463 detail-page
// schedule: four ranges separated by "/" (not newlines), with parenthetical
// delivery-mode notes and one label-prefixed ("Bridge day") entry with no
// session number. Four complete date/date-range expressions.
const ifs1463Block = "October 22-27, 2026 (Onsite) / November 19-22, 2026 (Online) / " +
	"December 17-20, 2026 (online) / Bridge day November 7, 2026 (online)"

const wantIfs1463String = "2026-10-22 - 2026-10-27, 2026-11-19 - 2026-11-22, 2026-12-17 - 2026-12-20, 2026-11-07"

// TestParseLabeledMultirangeBlock exercises ParseBlock against real-world
// labeled multi-range schedule blocks: a text block that interleaves
// several complete date/date-range expressions with session labels,
// bilingual month spellings, and parenthetical delivery-mode notes. Parse
// selects a single whole-string GLR root, which on inputs like these either
// truncates the later expressions or reduces a stray month/year token to a
// silently-incomplete date (e.g. "2026-00-00"); see the ParseBlock
// docstring in parse_block.go for the full failure mode. ParseBlock must
// recover every complete expression instead.
func TestParseLabeledMultirangeBlock(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		wantItems   int
		wantString  string
		wantResidue bool
	}{
		{
			name:       "ifs_1494_five_labeled_sessions",
			in:         ifs1494Block,
			wantItems:  5,
			wantString: wantIfs1494String,
		},
		{
			name:       "ifs_1463_four_slash_separated_ranges",
			in:         ifs1463Block,
			wantItems:  4,
			wantString: wantIfs1463String,
		},
		{
			// The same #1494 schedule, but pulled verbatim (modulo the \x1e
			// -> "\n" join a caller performs) from the committed evidence
			// under testdata/test_web_ifs-institute-com/.../*1494*
			// PageRecordCards's "schedules" raw_text fields — markdown bold
			// wrapping ("**...**", "***...***") and all. Same five items:
			// the label text differs in shape from the plain-text fixture
			// above, but ParseBlock's anchors start after the label either
			// way, so the markdown noise shouldn't matter.
			name: "ifs_1494_real_evidence_markdown_labels",
			in: "**Session 1 / sesión 1:** August/Agosto 5-8, 2026 (online/en línea)\n" +
				"***Bridge day / Día Puente***: Agosto 29 de 2026  (online/en línea)\n" +
				"**Session 2 / Sesión 2:** September/septiembre 9-12, 2026  (online/en línea)\n" +
				"***Bridge day/ Día Puente***: October/Octubre 2, 2026  (online/en línea)\n" +
				"**Session 3 / Sesión 3:** October/Octubre 14-17, 2026 (ONSITE/PRESENCIAL)",
			wantItems:  5,
			wantString: wantIfs1494String,
		},
		{
			// Negative: the only date-like content is a bare year. It must
			// not surface as an incomplete "2026-00-00" item; it must
			// surface as zero items plus a residue diagnostic.
			name:        "bare_year_only_is_residue_not_an_item",
			in:          "2026",
			wantItems:   0,
			wantResidue: true,
		},
		{
			// Negative: the only date-like content is a bare month name (no
			// day, no year). Same contract: zero items, residue flagged.
			name:        "bare_month_only_is_residue_not_an_item",
			in:          "Point Four Themes in Four Films August",
			wantItems:   0,
			wantResidue: true,
		},
		{
			// Negative: no date-bearing content at all. Zero items, and no
			// false-positive residue either.
			name:        "no_date_bearing_content_no_residue",
			in:          "Point Four Themes in Four Films",
			wantItems:   0,
			wantResidue: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			be, err := ParseBlock(tc.in, ParseOptions{})
			if err != nil {
				t.Fatalf("ParseBlock(%q) returned error: %v", tc.in, err)
			}
			if len(be.Items) != tc.wantItems {
				t.Errorf("ParseBlock(%q) got %d items, want %d\ngot string: %s", tc.in, len(be.Items), tc.wantItems, be.String())
			}
			if tc.wantString != "" && be.String() != tc.wantString {
				t.Errorf("ParseBlock(%q).String()\n got  %s\n want %s", tc.in, be.String(), tc.wantString)
			}
			if be.HasResidue != tc.wantResidue {
				t.Errorf("ParseBlock(%q).HasResidue = %v, want %v", tc.in, be.HasResidue, tc.wantResidue)
			}

			// Item spans must be valid, non-overlapping, and in document
			// order: each item's Start:End must slice the actual input,
			// and later items must not start before an earlier item ends.
			prevEnd := 0
			for i, item := range be.Items {
				if item.Start < 0 || item.End > len(tc.in) || item.Start >= item.End {
					t.Errorf("item %d has invalid span [%d:%d] for input of length %d", i, item.Start, item.End, len(tc.in))
					continue
				}
				if item.Start < prevEnd {
					t.Errorf("item %d span [%d:%d] overlaps the previous item, which ended at %d", i, item.Start, item.End, prevEnd)
				}
				prevEnd = item.End
			}
		})
	}
}

// blockBenchInput is a ~620-byte realistic block (the #1494 five-session
// fixture plus the #1463 four-range fixture, repeated once more to push the
// size toward the ~700-byte figure used for the cost-bound check) used to
// measure ParseBlock's wall time on input of realistic size and density.
var blockBenchInput = ifs1494Block + "\n" + ifs1463Block + "\n" + ifs1463Block

// BenchmarkParseBlock measures wall time on a realistic block (see the
// design constraint in the block-scoped extraction task: block-scoped
// extraction must not exhibit pathological GLR blowup on realistic input).
// Run with: go test -run '^$' -bench BenchmarkParseBlock -benchtime 20x -v ./datetime/
func BenchmarkParseBlock(b *testing.B) {
	b.Logf("input size: %d bytes", len(blockBenchInput))
	ClearParseCache() // measure GLR cost, not the (minDT,dateMode,tz,year,input) result cache
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ParseBlock(blockBenchInput, ParseOptions{}); err != nil {
			b.Fatalf("ParseBlock: %v", err)
		}
	}
}

// TestParseBlockPrefersExplicitYearOverDefault guards a specific ordering
// choice in parseBlockWindow: it grows each candidate window word by word
// and stops at the FIRST complete result, which is what lets it skip past
// trailing noise (a next entry's label) instead of accidentally consuming
// it. But when a caller supplies DefaultYear/MinDateTime, a bare "Month
// Day-Day" prefix can already be "complete" under the default BEFORE the
// text's own explicit trailing year is reached. ParseBlock must still
// prefer that explicit year over the default rather than stopping early.
func TestParseBlockPrefersExplicitYearOverDefault(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("LoadLocation: %v", err)
	}
	opts := ParseOptions{
		DefaultYear:     2020,
		DefaultLocation: loc,
		MinDateTime:     &DateTime{Date: &Date{Year: 2020, Month: 1, Day: 1}},
	}

	t.Run("explicit_year_in_text_wins", func(t *testing.T) {
		in := "Session 1: August 5-8, 2026 (online)\nSession 2: September 9-12, 2026 (online)"
		be, err := ParseBlock(in, opts)
		if err != nil {
			t.Fatalf("ParseBlock: %v", err)
		}
		want := "2026-08-05 - 2026-08-08, 2026-09-09 - 2026-09-12"
		if got := be.String(); got != want {
			t.Errorf("got  %s\nwant %s", got, want)
		}
	})

	t.Run("default_year_used_when_text_has_none", func(t *testing.T) {
		in := "Session 1: August 5-8 (online)\nSession 2: September 9-12 (online)"
		be, err := ParseBlock(in, opts)
		if err != nil {
			t.Fatalf("ParseBlock: %v", err)
		}
		want := "2020-08-05 - 2020-08-08, 2020-09-09 - 2020-09-12"
		if got := be.String(); got != want {
			t.Errorf("got  %s\nwant %s", got, want)
		}
	})
}

// TestParseBlockYearLabelContext exercises the block-scoped year context
// (blockYearContext): a listing that states its year once as a standalone
// leading label ("2027:") and lists year-less expressions under it. The
// label word always fails its own window (it is bounded at the next month
// anchor, so it parses alone as an incomplete bare year), and before the
// context existed it was dropped as residue while every year-less
// expression resolved under MinDateTime/DefaultYear — the wrong year
// whenever the label disagrees. The label's year must govern subsequent
// year-less expressions, a month regression (December then January) within
// the labeled block must advance the year, explicit years in the text must
// still win, and a block with no label must keep the caller's year context
// unchanged. The first two inputs are the real IFS Institute training-list
// blocks that surfaced the defect.
func TestParseBlockYearLabelContext(t *testing.T) {
	// The pipeline's reference-date year context: mid-2026, no operator
	// DefaultYear — the exact options under which the defect surfaced.
	opts := ParseOptions{
		MinDateTime: &DateTime{Date: &Date{Year: 2026, Month: 7, Day: 7}},
	}

	tests := []struct {
		name        string
		in          string
		opts        ParseOptions
		wantString  string
		wantResidue bool
	}{
		{
			// The leading "2027:" label must set the year for all six
			// year-less expressions; without the context they resolve to
			// MinDateTime's 2026.
			name:       "leading_year_label_governs_yearless_expressions",
			in:         "2027: January 26-28, February 11, March 2-4, March 25, April 13-15, May 18-20",
			opts:       opts,
			wantString: "2027-01-26 - 2027-01-28, 2027-02-11, 2027-03-02 - 2027-03-04, 2027-03-25, 2027-04-13 - 2027-04-15, 2027-05-18 - 2027-05-20",
		},
		{
			// The trailing January entry regresses from December and must
			// advance to 2027 while every earlier item stays under the 2026
			// label. ("September 29 - October 1" splits at the October month
			// anchor into two single-day items — the pre-existing windowing
			// shape, unchanged by the year context.)
			name:       "december_to_january_regression_advances_year",
			in:         "2026: August 4-6, September 2, September 29 - October 1, November 2, December 1-3, January 5-7",
			opts:       opts,
			wantString: "2026-08-04 - 2026-08-06, 2026-09-02, 2026-09-29, 2026-10-01, 2026-11-02, 2026-12-01 - 2026-12-03, 2027-01-05 - 2027-01-07",
		},
		{
			// No label: year-less expressions keep resolving under the
			// caller's MinDateTime year, exactly as before the context
			// existed — including a month regression, which must NOT advance
			// the year without a label.
			name:       "no_label_keeps_caller_year_context",
			in:         "December 1-3, January 5-7",
			opts:       opts,
			wantString: "2026-12-01 - 2026-12-03, 2026-01-05 - 2026-01-07",
		},
		{
			// Expressions carrying explicit years are unaffected by the
			// label: DefaultYear only fills dates the grammar left year-less.
			name:       "explicit_years_win_over_label",
			in:         "2027: August 5-8, 2026 (online) September 9-12, 2026",
			opts:       opts,
			wantString: "2026-08-05 - 2026-08-08, 2026-09-09 - 2026-09-12",
		},
		{
			// A label no expression ever resolves under is still
			// date-bearing text that produced no item: residue.
			name:        "label_with_no_following_expressions_is_residue",
			in:          "2027:",
			opts:        opts,
			wantString:  "",
			wantResidue: true,
		},
		{
			// The block's own label is more specific than the operator-level
			// DefaultYear and must override it.
			name: "label_overrides_caller_default_year",
			in:   "2027: January 26-28",
			opts: ParseOptions{
				DefaultYear: 2020,
				MinDateTime: &DateTime{Date: &Date{Year: 2026, Month: 7, Day: 7}},
			},
			wantString: "2027-01-26 - 2027-01-28",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			be, err := ParseBlock(tc.in, tc.opts)
			if err != nil {
				t.Fatalf("ParseBlock(%q) returned error: %v", tc.in, err)
			}
			if got := be.String(); got != tc.wantString {
				t.Errorf("ParseBlock(%q).String()\n got  %s\n want %s", tc.in, got, tc.wantString)
			}
			if be.HasResidue != tc.wantResidue {
				t.Errorf("ParseBlock(%q).HasResidue = %v, want %v", tc.in, be.HasResidue, tc.wantResidue)
			}
		})
	}
}
