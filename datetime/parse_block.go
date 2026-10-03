package datetime

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ParseBlock extracts every non-overlapping, independently complete date or
// date-range expression from one text block, in document order, instead of
// selecting a single whole-string parse the way Parse does.
//
// Parse picks one GLR parse root for the entire input. That works for a
// single date expression, but a block that interleaves several labeled date
// expressions with prose ("Session 1: Aug 5-8, 2026 / Bridge day: Aug 29,
// 2026 / Session 2: ...") has no single grammar path spanning the whole
// string: the chosen root either truncates the later expressions or, worse,
// reduces a stray month/year token to a silently-incomplete date (a Date
// with Year, Month, or Day left at zero — e.g. "2026-00-00"). ParseBlock
// avoids that failure mode by scanning the block for date-triggering
// anchors (a month name, a plausible year, or "today"/"tomorrow"/
// "yesterday"), and for each anchor, calling the existing single-expression
// Parse on a window of the block starting at that anchor, growing the
// window rightward one whitespace-delimited word at a time until Parse
// returns a result with no incomplete date fields (see parseBlockWindow for
// why growing, rather than starting wide and trimming, is what keeps
// trailing labels and noise from leaking into the result). One bounded
// exception preserves the existing six-word named-month conjunction under
// its trailing year before splitting at the second month; another preserves
// an immediately advertised date-labelled clock range before its date-only
// prefix can be accepted. These reuse the
// existing grammar and its noise tolerance (labels, parentheticals,
// bilingual prose) rather than adding new grammar productions.
//
// opts is applied to every windowed Parse call, so MinDateTime/DateMode/
// DefaultLocation/DefaultYear behave the same as for Parse — except that a
// block's own standalone year label ("2027:") overrides opts.DefaultYear for
// the windows after it, and a month regression between accepted items
// advances that label year; see blockYearContext. A block with no year label
// parses every window with opts unchanged.
func ParseBlock(in string, opts ParseOptions) (*BlockExtraction, error) {
	words := blockWords(in)
	anchors := make([]int, 0, len(words))
	for i, w := range words {
		if isAnchorWord(w.parseText) {
			anchors = append(anchors, i)
		}
	}

	out := &BlockExtraction{}
	yearCtx := &blockYearContext{}
	cursor := 0 // first word index not yet consumed by an accepted item
	for _, ai := range anchors {
		if ai < cursor {
			continue // already inside a previously accepted item's window
		}

		// Bound the window at the next anchor that could start a NEW
		// expression. A bare year is excluded from that search: it is
		// usually the trailing year of the phrase started at ai itself
		// ("August 5-8, 2026"), so treating it as a hard boundary would
		// truncate the window right before the year it needs. A month name
		// or relative-day word, in contrast, can appear at most once per
		// expression, so the next occurrence always marks a new one.
		upper := len(words)
		for _, aj := range anchors {
			if aj > ai && isMonthOrRelativeDayWord(words[aj].parseText) {
				upper = aj
				break
			}
		}
		if upper-ai > maxBlockWindowWords {
			upper = ai + maxBlockWindowWords
		}

		rs, consumedThrough, err := parseBlockWindow(words, ai, upper, yearCtx.options(opts))
		if err != nil {
			return nil, err
		}
		if rs == nil {
			if yearCtx.takeLabel(words[ai]) {
				continue // consumed as the block's year label, not residue
			}
			out.HasResidue = true
			continue
		}
		rs = yearCtx.observe(words[ai:consumedThrough], opts, rs)

		start := words[ai].rawStart
		end := words[consumedThrough-1].rawEnd
		for _, item := range rs.Items {
			out.Items = append(out.Items, BlockItem{Range: item, Start: start, End: end})
		}
		cursor = consumedThrough
	}

	if yearCtx.pendingResidue {
		// A year label that no later expression resolved under is still
		// date-bearing text that produced no item.
		out.HasResidue = true
	}
	return out, nil
}

// blockYearContext carries a block's own year context across ParseBlock's
// windows. Listings often state a year once as a standalone label and list
// year-less expressions under it ("2027: January 26-28, February 11, ...").
// The label word itself can never join a windowed Parse — the window
// starting at it is bounded at the next month anchor, so the label parses
// alone as an incomplete bare year — and without this context it would be
// dropped as residue while every year-less expression after it silently
// resolved under opts.DefaultYear/MinDateTime: the wrong year whenever the
// block's label disagrees. The context:
//
//   - activates when a standalone year-label word (blockYearLabel) fails to
//     resolve as its own window: the label's year then overrides
//     opts.DefaultYear for every subsequent windowed Parse, so year-less
//     expressions resolve under the block's own year. Expressions whose text
//     carries an explicit year are unaffected — DefaultYear only fills dates
//     the grammar left year-less.
//   - advances the year on a month regression between accepted items
//     (December then January), the way a wrapping schedule list reads: the
//     regressed window is re-parsed under the incremented year and adopted
//     only when its dates actually followed the default — an explicit year
//     in the window's own text keeps the window as parsed.
//   - is replaced by a later standalone year label, so a multi-section block
//     ("2026: ... 2027: ...") switches context at each label.
//
// A block with no year label never activates the context and parses every
// window with the caller's options unchanged.
type blockYearContext struct {
	// year is the active label year; 0 means no label has activated the
	// context.
	year int
	// prevMonth is the month the previous accepted item ended in (its end
	// date's month when present, else its start's), used to detect a month
	// regression; 0 until an item is accepted under the active label.
	prevMonth time.Month
	// pendingResidue is true while the active label has not yet been
	// followed by any accepted item, so a trailing or expression-less label
	// still surfaces as residue.
	pendingResidue bool
}

// options returns opts with the active label year (if any) as DefaultYear.
func (c *blockYearContext) options(opts ParseOptions) ParseOptions {
	if c.year != 0 {
		opts.DefaultYear = c.year
	}
	return opts
}

// takeLabel consumes w as a standalone year label, activating (or replacing)
// the block's year context. It reports false when w is not a year-label word,
// in which case the caller treats the failed window as residue as before.
func (c *blockYearContext) takeLabel(w blockWord) bool {
	y, ok := blockYearLabel(w)
	if !ok {
		return false
	}
	c.year = y
	c.prevMonth = 0
	c.pendingResidue = true
	return true
}

// observe folds an accepted window result into the active year context: on a
// month regression against the previous accepted item it re-parses the same
// window text under the incremented year (base opts, so the explicit-year
// absorption already baked into the window text still wins) and adopts the
// re-parse only when it is complete and its dates actually followed the
// incremented default. It returns the (possibly re-parsed) result to
// materialize. With no active label it returns rs unchanged.
func (c *blockYearContext) observe(window []blockWord, opts ParseOptions, rs *DateTimeRanges) *DateTimeRanges {
	if c.year == 0 {
		return rs
	}
	c.pendingResidue = false
	if m := blockStartMonth(rs); c.prevMonth != 0 && m != 0 && m < c.prevMonth {
		next := opts
		next.DefaultYear = c.year + 1
		if rs2, err := Parse(joinParseText(window), next); err == nil &&
			allItemsComplete(rs2) && blockStartYear(rs2) == c.year+1 {
			rs, c.year = rs2, c.year+1
		}
	}
	if m := blockEndMonth(rs); m != 0 {
		c.prevMonth = m
	}
	return rs
}

// blockYearLabel reports the year of a standalone year-label word: a bare
// plausible year immediately followed by a colon in the raw text ("2027:").
// The colon is what distinguishes a label governing the expressions after it
// from a stray year in prose (a copyright line, a "since 1998"), so a bare
// dangling year without one stays residue rather than becoming context.
func blockYearLabel(w blockWord) (int, bool) {
	if !strings.HasSuffix(w.raw, ":") {
		return 0, false
	}
	t := trimWordPunct(w.raw)
	if !blockYearRE.MatchString(t) {
		return 0, false
	}
	y, err := strconv.Atoi(t)
	if err != nil {
		return 0, false
	}
	return y, true
}

// blockStartMonth returns the month of the first item's start date, or 0.
func blockStartMonth(rs *DateTimeRanges) time.Month {
	if rs == nil || len(rs.Items) == 0 {
		return 0
	}
	if dt := rs.Items[0].Start; dt != nil && dt.Date != nil {
		return dt.Date.Month
	}
	return 0
}

// blockStartYear returns the year of the first item's start date, or 0.
func blockStartYear(rs *DateTimeRanges) int {
	if rs == nil || len(rs.Items) == 0 {
		return 0
	}
	if dt := rs.Items[0].Start; dt != nil && dt.Date != nil {
		return dt.Date.Year
	}
	return 0
}

// blockEndMonth returns the month a result ends in: the last item's end
// date's month when present, else its start date's month, or 0.
func blockEndMonth(rs *DateTimeRanges) time.Month {
	if rs == nil || len(rs.Items) == 0 {
		return 0
	}
	last := rs.Items[len(rs.Items)-1]
	if last.End != nil && last.End.Date != nil && last.End.Date.Month != 0 {
		return last.End.Date.Month
	}
	if last.Start != nil && last.Start.Date != nil {
		return last.Start.Date.Month
	}
	return 0
}

// maxBlockWindowWords bounds how many words a single ParseBlock window may
// span, both between two anchors (a pathologically long run of prose
// between two date-bearing tokens) and past the last anchor. It is sized
// generously for one date/time/timezone/parenthetical expression while
// keeping the number of trial Parse calls (and so worst-case GLR cost)
// bounded per anchor regardless of block size.
const maxBlockWindowWords = 30

// parseBlockWindow grows a window from words[start:start+1] up to
// words[start:upper] one word at a time, calling Parse on the joined
// parseText of each candidate. First it tries the bounded shared-year day
// conjunction as one existing grammar expression, even across the ordinary
// next-month boundary. Otherwise it stops at the first length whose result is
// fully complete (see allItemsComplete). A bounded date-labelled clock range
// is first presented intact so its real hours cannot be lost to a shorter
// date-only prefix. Otherwise, when the word
// immediately after a complete result is itself a bare year (isYearWord),
// it tries absorbing that word too before stopping: a date phrase that
// became "complete" using a defaulted year (from opts.DefaultYear /
// opts.MinDateTime) must still prefer an explicit trailing year over the
// default when the input actually has one, e.g. "August 5-8" alone could
// already be complete under a default year, but the "2026" right after it
// is the real, explicit year and must win. Absorption stops as soon as the
// next word isn't a bare year, the window is exhausted, or absorbing it
// stops the result from being complete.
//
// Growing (rather than shrinking from the widest window and trimming
// trailing words off) matters beyond that: a wide window that happens to
// include the START of the next labeled entry (a session number, a day-2
// numeral) can itself parse as "complete" by misreading a stray digit as an
// hour, silently attaching a wrong time to an otherwise-correct date.
// Growing only as far as completeness requires never reaches that trailing
// text in the first place.
//
// It returns the accepted result and the exclusive word index the window
// consumed through, or a nil result if no window length up to upper
// produced a complete parse. Individual Parse failures/errors for a given
// length are not fatal — they just mean that length doesn't work — so this
// never returns a non-nil error; the return keeps the signature symmetric
// with Parse's.
func parseBlockWindow(words []blockWord, start, upper int, opts ParseOptions) (*DateTimeRanges, int, error) {
	// A date-labelled clock row is one existing grammar expression. Taking
	// its shorter default-complete date would silently discard actual hours.
	if rs := parseBlockDatedClockRange(words[start:upper], opts); rs != nil {
		return rs, upper, nil
	}
	// A conjunction can join two named-month days under one trailing year.
	// Keep that existing grammar expression intact before the ordinary
	// next-month boundary or a default-complete first day can split it.
	if rs, end := parseBlockSharedYearDays(words, start, opts); rs != nil {
		return rs, end, nil
	}
	var lastGood *DateTimeRanges
	lastGoodK := 0
	for k := start + 1; k <= upper; k++ {
		rs, err := Parse(joinParseText(words[start:k]), opts)
		if err != nil || !allItemsComplete(rs) {
			if lastGood != nil {
				break // the previous length was already complete; stop growing.
			}
			continue
		}
		lastGood, lastGoodK = rs, k
		if k >= upper || !isYearWord(words[k].parseText) {
			break // nothing more worth absorbing.
		}
		// Next word is a bare year: try absorbing it before settling.
	}
	if lastGood == nil {
		return nil, 0, nil
	}
	return lastGood, lastGoodK, nil
}

// parseBlockDatedClockRange preserves "Month Day: clock-clock [zone]" as one
// bounded window. The existing clock-range recognizer rejects empty labels
// and trailing prose; the grammar remains the date/time syntax authority.
func parseBlockDatedClockRange(words []blockWord, opts ParseOptions) *DateTimeRanges {
	if len(words) < 3 || !strings.HasSuffix(words[1].raw, ":") {
		return nil
	}
	if _, ok := monthsByNames[strings.ToLower(trimWordPunct(words[0].parseText))]; !ok {
		return nil
	}
	day, err := strconv.Atoi(strings.TrimSuffix(words[1].raw, ":"))
	if err != nil || day < 1 || day > 31 {
		return nil
	}
	clock := simpleTimeRangeRE.FindStringSubmatch(joinParseText(words[2:]))
	if clock == nil {
		return nil
	}
	if clock[7] != "" {
		zone := recurrenceTimeZone(clock[7])
		if zone == nil || zone.IANAName() == "" {
			return nil // a trailing prose word is not zone evidence
		}
	}
	rs, err := Parse(joinParseText(words), opts)
	if err != nil || !allItemsComplete(rs) || len(rs.Items) != 1 {
		return nil
	}
	r := rs.Items[0]
	if r.Start.Time == nil || r.End == nil || r.End.Time == nil {
		return nil
	}
	return rs
}

// parseBlockSharedYearDays admits only the existing six-word grammar shape
// "March 20 and April 23, 2027". It neither carries a year across unrelated
// entries nor extends a window into labels, clocks or arbitrary prose.
func parseBlockSharedYearDays(words []blockWord, start int, opts ParseOptions) (*DateTimeRanges, int) {
	end := start + 6
	if end > len(words) || !strings.EqualFold(words[start+2].parseText, "and") || !isYearWord(words[end-1].parseText) {
		return nil, 0
	}
	for _, offset := range []int{0, 3} {
		if _, ok := monthsByNames[strings.ToLower(trimWordPunct(words[start+offset].parseText))]; !ok {
			return nil, 0
		}
		day, err := strconv.Atoi(trimWordPunct(words[start+offset+1].parseText))
		if err != nil || day < 1 || day > 31 {
			return nil, 0
		}
	}
	rs, err := Parse(joinParseText(words[start:end]), opts)
	if err != nil || !allItemsComplete(rs) || len(rs.Items) != 2 {
		return nil, 0
	}
	return rs, end
}

// joinParseText joins the words' parse-view text with single spaces — the
// exact string a windowed Parse call sees for that word span.
func joinParseText(words []blockWord) string {
	parts := make([]string, len(words))
	for i, w := range words {
		parts[i] = w.parseText
	}
	return strings.Join(parts, " ")
}

// allItemsComplete reports whether every item in rs has a fully-specified
// date (Year, Month, and Day all set): the block-extraction contract never
// emits a bare year or month fragment (e.g. "2026-00-00") as a result item.
func allItemsComplete(rs *DateTimeRanges) bool {
	if rs == nil || len(rs.Items) == 0 {
		return false
	}
	for _, item := range rs.Items {
		if item == nil || !dateTimeComplete(item.Start) {
			return false
		}
		if item.End != nil && item.End.Date != nil && !dateComplete(item.End.Date) {
			return false
		}
	}
	return true
}

func dateTimeComplete(dt *DateTime) bool {
	return dt != nil && dateComplete(dt.Date)
}

func dateComplete(d *Date) bool {
	return d != nil && d.Year != 0 && d.Month != 0 && d.Day != 0
}

// BlockExtraction is the result of ParseBlock: every complete date/date-range
// expression found in one text block, in document order, plus a coverage
// diagnostic.
type BlockExtraction struct {
	// Items holds one entry per complete date/date-range expression found,
	// in document order.
	Items []BlockItem
	// HasResidue is true when some date-triggering text (a year, a month
	// name, "today"/"tomorrow"/"yesterday") never resolved into a complete
	// item — for example a bare year or a month name with no day. Callers
	// can use this to detect silent partial success: the block had
	// date-bearing content that ParseBlock could not turn into a complete
	// date.
	HasResidue bool
}

// BlockItem is one extracted date/date-range expression, plus the byte
// offsets of the span in ParseBlock's input that produced it.
type BlockItem struct {
	Range *DateTimeRange
	// Start and End are byte offsets into the original ParseBlock input
	// ([Start, End)) covering the anchor word through the last word Parse
	// needed to produce Range. Trailing noise within that span (a
	// parenthetical, a dangling separator) that Parse tolerated as part of
	// the same expression is included; text before Start or at/after End is
	// not part of this item.
	Start, End int
}

func (be *BlockExtraction) String() string {
	if be == nil {
		return ""
	}
	rs := make([]string, len(be.Items))
	for i, item := range be.Items {
		rs[i] = item.Range.String()
	}
	return strings.Join(rs, ", ")
}

// blockWord is one whitespace-delimited token from a ParseBlock input, with
// its raw byte span in the original input and the text to feed to Parse
// (which differs from the raw text only for a glued bilingual month pair,
// e.g. "August/Agosto" -> "August"; see bilingualMonthWord).
type blockWord struct {
	rawStart, rawEnd int
	raw              string
	parseText        string
}

// blockWords splits s into whitespace-delimited words, preserving each
// word's byte offsets so ParseBlock can report source spans in terms of the
// original input.
func blockWords(s string) []blockWord {
	var words []blockWord
	start := -1
	for i, r := range s {
		if unicode.IsSpace(r) {
			if start >= 0 {
				words = append(words, blockWord{rawStart: start, rawEnd: i, raw: s[start:i]})
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		words = append(words, blockWord{rawStart: start, rawEnd: len(s), raw: s[start:]})
	}
	for i := range words {
		words[i].parseText = bilingualMonthWord(words[i].raw)
	}
	return words
}

// wordPunctRE strips punctuation commonly attached to a date-bearing word
// (leading "(", "[" and trailing ")", "]", ",", ".", ":", ";") before
// classifying it as a month name, a year, or a relative-day word.
var wordPunctRE = regexp.MustCompile(`^[(\[]+|[)\],.:;]+$`)

func trimWordPunct(w string) string {
	return wordPunctRE.ReplaceAllString(w, "")
}

// bilingualMonthWord collapses a glued bilingual month pair like
// "August/Agosto" or "October/Octubre" to just the first spelling. Event
// listings that give a month name in two languages join the pair with a
// slash and no space (a lexical/punctuation-classification quirk, not a new
// date shape), and left as-is the second spelling derails the grammar (it
// has no production for "MONTH_NAME QUO MONTH_NAME" or "MONTH_NAME QUO
// IDENT"). Only collapses when the first segment is a recognized month name
// and the second segment either isn't a recognized month name (so it can't
// be a legitimate second month in a real cross-month expression) or
// resolves to the SAME month number (a genuine translation pair, e.g.
// "septiembre" for September even though "septiembre" itself isn't in
// monthsByNames). A real cross-month expression like "August/September" is
// left untouched.
func bilingualMonthWord(w string) string {
	slash := strings.IndexByte(w, '/')
	if slash < 0 {
		return w
	}
	first := w[:slash]
	second := w[slash+1:]
	if idx := strings.IndexByte(second, '/'); idx >= 0 {
		second = second[:idx]
	}
	firstMonth, firstOK := monthsByNames[strings.ToLower(trimWordPunct(first))]
	if !firstOK {
		return w
	}
	if secondMonth, secondOK := monthsByNames[strings.ToLower(trimWordPunct(second))]; secondOK && secondMonth != firstMonth {
		return w // genuine cross-month expression; keep both spellings
	}
	return first
}

var blockYearRE = regexp.MustCompile(`^(?:17|18|19|20|21)\d{2}$`)

var blockRelativeDayWords = map[string]bool{
	"today":     true,
	"tomorrow":  true,
	"yesterday": true,
}

// isAnchorWord reports whether a word (already run through
// bilingualMonthWord) is a date-triggering token: a recognized month name,
// a plausible 4-digit year, or a relative-day word. These are the same
// token classes the datetimeLexer itself classifies as MONTH_NAME, YEAR,
// and RELATIVE_DAY — ParseBlock's anchors are exactly the positions where
// the grammar's own vocabulary says a date could start, so the windowing
// stays general instead of keying off any site-specific label shape.
func isAnchorWord(w string) bool {
	return isMonthOrRelativeDayWord(w) || isYearWord(w)
}

// isMonthOrRelativeDayWord reports whether w is a recognized month name or a
// relative-day word ("today"/"tomorrow"/"yesterday"). Unlike isYearWord,
// this excludes bare years: the next month ordinarily starts a new block
// expression, whereas a year usually completes the expression in progress.
// parseBlockSharedYearDays preserves the bounded conjoined two-month shape
// as the explicit exception to that ordinary block-window boundary.
func isMonthOrRelativeDayWord(w string) bool {
	trimmed := strings.ToLower(trimWordPunct(w))
	if trimmed == "" {
		return false
	}
	if _, ok := monthsByNames[trimmed]; ok {
		return true
	}
	return blockRelativeDayWords[trimmed]
}

func isYearWord(w string) bool {
	return blockYearRE.MatchString(trimWordPunct(w))
}
