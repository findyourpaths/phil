# DateTime Parser Issues and Fixes

## Summary

This document describes the investigation and fixes applied to the datetime parser to reduce test failures from 31 (13.78%) to 28 (12.44%) out of 225 total tests.

## Initial State

**Test Failures:** 31 out of 225 (13.78%)

**Compilation Error:**
- `glr/glr_parse.go:481` - Non-constant format string error blocking all tests

**Major Categories of Failures:**
1. ISO 8601 format issues (year-month, time without seconds)
2. Semantic validation issues (invalid range interpretations)
3. Date list parsing (comma-separated days with month)
4. Complex natural language patterns
5. Timezone and context-dependent parsing

## Fixes Applied

### 1. Compilation Error Fix

**File:** `glr/glr_parse.go:481`

**Issue:**
```go
debugf(fmt.Sprintf("parsing term: %#v\n", term))
```

**Fix:**
```go
debugf("parsing term: %#v\n", term)
```

**Reason:** `debugf` already handles format strings, so `fmt.Sprintf` wrapper was unnecessary and caused compilation error.

---

### 2. ISO 8601 Year-Month Format Support

**File:** `datetime/parse_yacc.y:368-371`

**Issue:** Test 001 ("2023-02") failed because parser didn't recognize year-month format without day.

**Grammar Change:**
```yacc
RFC3339Date:
  Year SUB INT SUB INT {$ = NewRawDateFromDMY($5, $3, $1)}
| Year SUB INT {$ = NewRawDateFromDMY(nil, $3, $1)}  // Added this line
;
```

**Result:** Reduced failures from 31 → 30

---

### 3. ISO 8601 Time Without Seconds

**File:** `datetime/parse_yacc.y:373-376`

**Issue:** Test 004 ("2023-02-03T12:00") failed because parser required seconds in RFC3339 time.

**Grammar Change:**
```yacc
RFC3339Time:
  T INT COLON INT COLON INT {$ = NewTime($2, $4, $6, nil)}
| T INT COLON INT {$ = NewTime($2, $4, nil, nil)}  // Added this line
;
```

---

### 4. Grammar Rule Precedence

**File:** `datetime/parse_yacc.y:303-314`

**Issue:** DateTime rule alternatives were evaluated in wrong order, causing ISO format to not be preferred.

**Grammar Change:**
```yacc
DateTime:
  RFC3339DateTime  // Moved to first position (was further down)

| Date {$ = NewDateTimeWithDate($1)}
| Date DateTimeSepPlus TimeZoneOpt {$ = NewDateTime($1, nil, $3)}
// ... other alternatives
;
```

**Reason:** In Yacc/GLR parsers, rule order affects precedence when resolving ambiguities.

**Result:** Reduced failures from 30 → 29

---

### 5. Semantic Validation for Invalid Ranges

**File:** `datetime/datetime.go:177-199`

**Issue:** Parser accepted semantically invalid ranges like "year 2023 to day 2", which caused test 001 to fail after grammar fix.

**Code Added:**
```go
func NewRange(start *DateTime, end *DateTime) *DateTimeRange {
    if start != nil && end != nil {
        // Reject ranges where one side is only a year and the other is only a day.
        if start.Date != nil && end.Date != nil {
            startOnlyYear := start.Date.Year != 0 && start.Date.Month == 0 && start.Date.Day == 0
            endOnlyDay := end.Date.Year == 0 && end.Date.Month == 0 && end.Date.Day != 0
            if startOnlyYear && endOnlyDay {
                panic(fmt.Sprintf("semantic error: cannot create range from year-only %#v to day-only %#v\n", start.Date, end.Date))
            }
            // Also check the reverse case
            endOnlyYear := end.Date.Year != 0 && end.Date.Month == 0 && end.Date.Day == 0
            startOnlyDay := start.Date.Year == 0 && start.Date.Month == 0 && start.Date.Day != 0
            if endOnlyYear && startOnlyDay {
                panic(fmt.Sprintf("semantic error: cannot create range from day-only %#v to year-only %#v\n", start.Date, end.Date))
            }
        }
        // ... rest of function
    }
    // ...
}
```

**Reason:** Grammar ambiguity allowed "2023-02" to be parsed as both "February 2023" and "range from year 2023 to day 2". Semantic validation rejects the nonsensical interpretation.

---

### 6. Date List Separator Rules

**File:** `datetime/parse_yacc.y:218-223`

**Issue:** Test 050 ("Feb 1, 2, 3, 4") failed because DayPlus1 rule didn't support comma separators.

**Grammar Change:**
```yacc
DayPlus1:
  Day Day {$$ = []string{$1, $2}}
| DayPlus1 Day {$$ = append($1, $2)}
| Day RangesSepPlus Day {$$ = []string{$1, $3}}              // Added
| DayPlus1 RangesSepPlus Day {$$ = append($1, $3)}           // Added
;
```

Where `RangesSepPlus` includes `AND` and `COMMA` tokens.

**Result:** Test 050 now passes

---

### 7. Ambiguous Date Separator Restriction

**File:** `datetime/parse_yacc.y:467-475`

**Issue:** Patterns like "February 1, 8" were being interpreted as "August 15" instead of "February 1 and February 8" because comma was being used in ambiguous date parsing.

**Grammar Addition:**
```yacc
AmbiguousDateSepPlus:
  AmbiguousDateSep
| AmbiguousDateSepPlus AmbiguousDateSep
;
AmbiguousDateSep:
  DEC      // slash (/)
| PERIOD   // dot (.)
| QUO      // other numeric separators
;
```

**Grammar Change:**
```yacc
Date:
  // "02.03", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| WeekdayOpt Day AmbiguousDateSepPlus Day {$$ = NewRawDateFromAmbiguous($1, $2, $4, nil)}
  // ... (changed from DateSepPlus to AmbiguousDateSepPlus)
| WeekdayOpt Day AmbiguousDateSepPlus Day AmbiguousDateSepPlus Year {$$ = NewRawDateFromAmbiguous($1, $2, $4, $6)}
;
```

**Reason:** Commas should not trigger ambiguous date parsing (which is for patterns like "2/3" that could be Feb 3 or Mar 2). By restricting ambiguous separators to only slashes, dots, etc., we prevent "1, 8" from being interpreted as month-day.

---

### 8. Weekday with Date List Support

**File:** `datetime/parse_yacc.y:150-162`

**Issue:** Test 121 ("Wednesdays, February 1, 8, 15, 22") and similar tests failed because grammar didn't support weekday prefix with date lists.

**Grammar Addition:**
```yacc
DateTimeRanges:
  // ... existing rules ...

  // "Feb 3, 4"
| Month DayPlus1 {$$ = NewRangesWithStartDates(NewRawDateFromMDsYs($1, $2, nil)...)}
  // "Wednesdays, Feb 3, 4"
| WeekdayOpt Month DayPlus1 {$$ = NewRangesWithStartDates(NewRawDateFromMDsYs($2, $3, nil)...)}  // Added

  // "Feb 3, 4 2023"
| Month DayPlus1 Year {$$ = NewRangesWithStartDates(NewRawDateFromMDsYs($1, $2, $3)...)}
  // "Wednesdays, Feb 3, 4 2023"
| WeekdayOpt Month DayPlus1 Year {$$ = NewRangesWithStartDates(NewRawDateFromMDsYs($2, $3, $4)...)}  // Added
;
```

---

### 9. ExplicitTime Nonterminal (Attempted, Partially Successful)

**File:** `datetime/parse_yacc.y:540-548`

**Issue:** Bare integers in date lists (like "8") were being parsed as times instead of days.

**Grammar Addition:**
```yacc
// ExplicitTime - times that are unambiguous (have AM/PM or colons), not bare integers
ExplicitTime:
  INT Am {$$ = NewAMTime($1, nil, nil, nil)}
| INT Pm {$$ = NewPMTime($1, nil, nil, nil)}
| INT TimeSep INT {$$ = NewTime($1, $3, nil, nil)}
| INT TimeSep INT TimeSep INT {$$ = NewTime($1, $3, $5, nil)}
| INT TimeSep INT Am {$$ = NewAMTime($1, $3, nil, nil)}
| INT TimeSep INT Pm {$$ = NewPMTime($1, $3, nil, nil)}
| TIME_NAME {$$ = NewTime($1, nil, nil, nil)}
;
```

**Application (line 190):**
```yacc
// "Wednesdays February 1st & 8th 12:00p-3:00p" - only with explicit times (AM/PM or colons)
| WeekdayOpt Month Day DateSepOpt Day ExplicitTime DateTimeSepOpt ExplicitTime {$$ = ...}
```

**Status:** Successfully prevents this specific pattern from causing ambiguity, but doesn't solve the general problem.

---

## Current State

**Test Failures:** 28 out of 225 (12.44%)

**Improvement:** 3 tests fixed (from 31 to 28)

### Remaining Failing Tests

```
044  Workshop_Update_(2/3/23)
053  Our_next_cohort_kicks_off_on_March_2nd_and_we're_accepting_applications_through_February_1st.
059  1-3_Feb_and_2_Mar
060  1-3_&_5_February
061  1-4_&_6_February
120  5_Wednesdays_2/1_&_2/8_&_2/15_&_2/22,_3/2
121  Wednesdays,_February_1,_8,_15,_22
122  Wednesdays,_February_1,_8,_15,_and_22
123  Wednesdays,_February_1,_8,_15_and_22
124  Wednesdays,_February_1,_8,_15_&_22
130  Join_today_for_Day_2_at_12pm
132  Today_Friday,_12pm_ET
141  Friday_2/14:_**Love_is_Listening_and_Art:_Social_+_Listening_Art_Sessions**_at_6pm_facilitated_by_Lauren_V
148  Feb_3_12pm_US/Eastern
152  Starts_Friday_2/3_at_9:00_am_ET
162  Feb._3,_2023_12:00pm,_3:00pm
170  Date:_Friday,_February_3rd_-_February_4th,_2023_ET
177  February_3_+_4,_9_am_-_12_pm_each_day
178  February_1_+_2:_In-person_at_SeekHealing_Asheville,12_pm_-_3_pm_each_day
180  2_Wednesdays_Feb_1st_&_8th_12:00p-3:00p
181  5_Wednesdays_9:00am-12:00pm_February_1st_-_March_1st
193  February_3rd,_12:00-3:00pm_Eastern_(New_York)_time
194  February_3rd_&_4th,_9:00_am_-_noon_Eastern_time
195  February_3rd_-_5th,_9:00_am_-_noon_ET_each_day
197  Wednesdays,_February_1st,_8th,_and_15th_9:00am_-_12:00pm_(ET)
198  Wednesdays_-_February_1,_8_9:00_AM_-_12:00_PM_ET
199  Wednesdays_-_February_1,_8,_15,_22,_and_March_1_9:00_AM_-_12:00_PM_ET
206  Fr._3._Feb._2023,_9:00-ca.12:00
```

---

## Root Causes of Remaining Issues

### 1. Fundamental Grammar Ambiguity: Bare Integers

**The Core Problem:**

The grammar allows bare integers to be parsed as either:
- **Day numbers** (in date contexts like "February 1, 8")
- **Hour numbers** (in time contexts like "12pm" or standalone "8")

This is specified in `parse_yacc.y:559`:
```yacc
Time:
  TimePrefixPlus Time {$$ = $2}
  // "11am"
| INT Am {$$ = NewAMTime($1, nil, nil, nil)}
| INT Pm {$$ = NewPMTime($1, nil, nil, nil)}
  // "12"
| INT {$$ = NewTime($1, nil, nil, nil)}  // ← This rule is the problem
```

And the DateTime rule (line 315):
```yacc
DateTime:
  RFC3339DateTime
| Date {$$ = NewDateTimeWithDate($1)}
| Date DateTimeSepPlus TimeZoneOpt {$$ = NewDateTime($1, nil, $3)}
| Date Time TimeZoneOpt {$$ = NewDateTime($1, $2, $3)}  // ← Combines Date with bare INT as Time
```

**Example Failure (Test 121):**

Input: `"Wednesdays, February 1, 8, 15, 22"`

Expected: Four separate dates (Feb 1, Feb 8, Feb 15, Feb 22)

Actual parsing:
1. "Wednesdays, February" → Weekday + Month
2. "1" → Day (matches `WeekdayOpt Month Day` → Date)
3. "," → consumed as separator
4. "8" → INT (ambiguous: could be Day or Time)
5. Parser chooses `Date Time` rule → DateTime with February 1 at 08:00
6. Remaining "15, 22" become separate dates

**Why This Happens:**

The GLR parser explores all possible parse trees for ambiguous inputs. When it encounters "February 1, 8", it finds multiple valid interpretations:

**Interpretation A (Correct):**
- `WeekdayOpt Month DayPlus1` where DayPlus1 = [1, 8, 15, 22]
- Creates 4 DateTimeRanges

**Interpretation B (Chosen by parser):**
- `WeekdayOpt Month Day` → Date (February 1)
- `INT` → Time (8:00)
- `Date Time` → DateTime (February 1, 08:00)
- Then continues with remaining tokens

**Conflict Resolution:**

The GLR parser's conflict resolution strategy (based on shift/reduce precedence and rule ordering) chooses Interpretation B. This is likely because:
1. The `Date Time` rule appears earlier in the DateTime alternatives
2. The DateTime rule is more specific (matches fewer tokens per reduction)
3. The parser's scoring heuristics favor smaller, more granular rules

---

### 2. Range Notation with Ampersands and Plus Signs

**Examples:**
- Test 059: "1-3 Feb and 2 Mar"
- Test 060: "1-3 & 5 February"
- Test 177: "February 3 + 4, 9 am - 12 pm each day"

**Issue:** The grammar has limited support for:
- Day ranges within a month ("1-3 Feb")
- Mixed operators ("&", "+", "and")
- Combining date ranges with multiple months

**Current Support:**

Limited rules exist (lines 165-182) for specific patterns:
```yacc
// "Feb 1-2, 3-4"
| Month Day RangeSep Day Day RangeSep Day {$ = NewRanges(...)}
// "1-2, 3-4 Feb"
| Day RangeSep Day Day RangeSep Day Month {$ = NewRanges(...)}
```

But these don't cover:
- Single range followed by single day ("1-3 & 5")
- Ranges across different months ("1-3 Feb and 2 Mar")
- Plus signs as separators

---

### 3. Contextual Date Resolution

**Examples:**
- Test 130: "Join today for Day 2 at 12pm"
- Test 132: "Today Friday, 12pm ET"
- Test 053: "...through February 1st"

**Issue:** Parser needs context (`minimumDateTime`) to resolve:
- "today" relative to current date
- Day numbers without month when month is implicit
- "Day 2" in sequences

**Current Implementation:**

`minimumDateTime` is a global variable used in `NewDateFromRaw()` (line 515-517):
```go
dm := DateMode(tz)
if dm == DateModeUnknown {
    dm = parseDateMode
}
if dm == DateModeUnknown && minimumDateTime != nil {
    dm = DateMode(minimumDateTime.TimeZone)
}
```

**Problems:**
- Insufficient context propagation through parsing
- "Day N" patterns not recognized as dates
- Relative dates in complex strings with other dates

---

### 4. Timezone Parsing

**Examples:**
- Test 148: "Feb 3 12pm US/Eastern"
- Test 193: "February 3rd, 12:00-3:00pm Eastern (New York) time"

**Issue:** Limited support for:
- Region-based timezone names ("US/Eastern" vs "Eastern")
- Timezone names with parenthetical clarifications
- Timezone names separated from time

---

### 5. Natural Language Prefixes/Suffixes

**Examples:**
- Test 152: "Starts Friday 2/3 at 9:00 am ET"
- Test 178: "February 1 + 2: In-person at SeekHealing..."

**Issue:** Text before/after dates interferes with parsing:
- Prefix words ("Starts", "Date:")
- Suffix descriptions (": In-person...")
- Embedded formatting ("**text**")

**Current Handling:**

Limited prefix support via `RangePrefixPlus` (line 194-200):
```yacc
RangesPrefixPlus:
  RangesPrefix
| RangesPrefixPlus RangesPrefix
;
RangesPrefix:
  WHEN
;
```

Only "WHEN" is recognized. Other prefix words are not handled.

---

### 6. Abbreviated Day/Month Names

**Example:**
- Test 206: "Fr. 3. Feb. 2023, 9:00-ca.12:00"

**Issue:**
- Abbreviated day names with periods ("Fr.")
- Day numbers with trailing periods ("3.")
- German-style abbreviations ("ca." for circa)

---

### 7. Complex Date-Time Combinations

**Examples:**
- Test 162: "Feb. 3, 2023 12:00pm, 3:00pm" (two times, one date)
- Test 180: "2 Wednesdays Feb 1st & 8th 12:00p-3:00p" (count prefix)
- Test 181: "5 Wednesdays 9:00am-12:00pm February 1st - March 1st"

**Issues:**
- Multiple times with single date
- Numeric count before weekday ("2 Wednesdays", "5 Wednesdays")
- Date ranges with weekday and time ranges

---

## Technical Analysis

### GLR Parser Behavior

The parser uses a **Generalized LR (GLR)** algorithm implemented in `glr/glr_parse.go`. Key characteristics:

1. **Parallel Parsing:** Maintains multiple parse stacks simultaneously to explore all possible interpretations
2. **Conflict Resolution:** When multiple interpretations exist, uses heuristics to choose one:
   - Rule ordering (earlier rules preferred)
   - Shift/reduce precedence
   - Custom scoring (not currently implemented)
3. **Ambiguity Reporting:** Current conflict count: 175 shift/reduce, 223 reduce/reduce

### Grammar Conflicts

From `go generate` output:
```
conflicts: 175 shift/reduce, 223 reduce/reduce
```

These conflicts represent points where the parser must make choices. Each conflict is a potential source of incorrect parsing.

**Shift/Reduce Conflicts:** Parser must choose between:
- **Shift:** Read another token before reducing
- **Reduce:** Apply a grammar rule now

**Reduce/Reduce Conflicts:** Parser must choose between multiple applicable grammar rules.

High conflict counts indicate significant grammar ambiguity.

---

## Attempted Solutions That Didn't Work

### 1. Using ExplicitTime in DateTime Rules

**Attempted Change:**
```yacc
DateTime:
  RFC3339DateTime
| Date {$$ = NewDateTimeWithDate($1)}
| Date ExplicitTime TimeZoneOpt {$$ = NewDateTime($1, $2, $3)}  // Changed from Time
```

**Result:**
- Broke 7 tests that rely on bare integers as times ("12 - 3pm")
- Failures increased from 28 to 35
- **Reverted**

**Why It Failed:** Too restrictive. Valid patterns like "February 3, 12 - 3pm" (where "12" is clearly a time in the context of a range with "3pm") were rejected.

---

### 2. Removing Bare INT from Time Rule

**Attempted Change:**
```yacc
Time:
  INT Am {$$ = NewAMTime($1, nil, nil, nil)}
| INT Pm {$$ = NewPMTime($1, nil, nil, nil)}
// | INT {$$ = NewTime($1, nil, nil, nil)}  // Removed this line
```

**Result:**
- Fixed date list tests
- Broke many time range tests
- Net negative impact

**Why It Failed:** Bare integer times are actually common and valid in many contexts (time ranges, military time, etc.)

---

## Potential Solutions

### Solution 1: Context-Aware Lexing (Recommended)

**Approach:** Make the lexer track parsing context and emit different tokens for integers based on context.

**Implementation:**
1. Add context states to lexer: `CONTEXT_DATE_LIST`, `CONTEXT_TIME_RANGE`, etc.
2. When lexer sees INT after Month+Day+Comma, emit `DAY_INT` instead of `INT`
3. Update grammar to use `DAY_INT` vs `TIME_INT` tokens

**Grammar Changes:**
```yacc
Time:
  TIME_INT {$$ = NewTime($1, nil, nil, nil)}
| TIME_INT Am {$$ = NewAMTime($1, nil, nil, nil)}
// ...
;

DayPlus1:
  Day DAY_INT {$$ = []string{$1, $2}}
| DayPlus1 DAY_INT {$$ = append($1, $2)}
// ...
;
```

**Pros:**
- Eliminates ambiguity at lexer level
- No grammar conflicts for INT tokens
- Can handle all date list patterns correctly

**Cons:**
- Complex lexer state management
- Requires careful definition of context boundaries
- May need lookahead in lexer

---

### Solution 2: Semantic Filtering with GLR Scoring

**Approach:** Allow GLR to parse all interpretations, but add semantic scoring to prefer correct ones.

**Implementation:**
1. Implement `Score()` method on parse tree nodes
2. Penalize unlikely interpretations:
   - DateTime with single-digit hour in middle of date list: -100
   - Day lists that create monotonic sequences: +50
3. GLR parser uses scores to choose best interpretation

**Code Changes in `glr/glr_parse.go`:**
```go
type ParseNodeScorer interface {
    Score() int
}

// In parse tree comparison
func (p *Parser) compareParseTrees(tree1, tree2 *ParseNode) *ParseNode {
    score1 := p.scoreTree(tree1)
    score2 := p.scoreTree(tree2)
    if score1 > score2 {
        return tree1
    }
    return tree2
}
```

**Grammar Actions:**
```yacc
// In DateTimeRanges for date list
| WeekdayOpt Month DayPlus1 {
    $$ = NewRangesWithStartDates(NewRawDateFromMDsYs($2, $3, nil)...)
    $$.score = 50  // Prefer date lists
}

// In DateTime for date+time
| Date Time TimeZoneOpt {
    $$ = NewDateTime($1, $2, $3)
    if isAmbiguousContext() {
        $$.score = -50  // Penalize in ambiguous contexts
    }
}
```

**Pros:**
- Leverages existing GLR infrastructure
- Doesn't require lexer changes
- Can handle complex ambiguity patterns

**Cons:**
- Requires defining comprehensive scoring rules
- May still fail on edge cases
- Debugging is harder (why did parser choose X?)

---

### Solution 3: Two-Pass Parsing

**Approach:** Parse with lenient grammar, then apply semantic analysis and re-parse if needed.

**Implementation:**

**Pass 1:** Parse with current grammar, collecting all ambiguous points

**Pass 2:** Apply heuristics:
```go
func resolveAmbiguity(input string, ambiguousParse *DateTimeRanges) *DateTimeRanges {
    // If we see "Month Day, INT, INT, INT" pattern
    if matchesDateListPattern(ambiguousParse) {
        // Re-parse with date-list-only grammar
        return reparseAsDateList(input)
    }

    // If we see "Date Time" where Time is single digit
    if hasAmbiguousTime(ambiguousParse) {
        // Check if interpretation as day makes more sense
        if shouldBeDayInstead(ambiguousParse) {
            return reparseWithoutTime(input)
        }
    }

    return ambiguousParse
}
```

**Pros:**
- Clean separation of concerns
- Can use different grammars for different patterns
- Easier to debug and maintain

**Cons:**
- Performance overhead (parsing twice)
- Requires maintaining multiple grammar variants
- May still miss some cases

---

### Solution 4: Explicit Disambiguation Tokens

**Approach:** Add optional explicit markers to grammar that user can use to disambiguate.

**Implementation:**

Allow (but don't require) users to use markers like:
- `"February {1, 8, 15, 22}"` - braces indicate date list
- `"February 1 @ 8:00"` - @ indicates time
- `"February 1 | 8 | 15"` - pipes indicate list separator

**Grammar:**
```yacc
DateTimeRanges:
| Month LBRACE DayList RBRACE {$$ = NewRangesWithStartDates(...)}
| Date AT Time {$$ = NewDateTime($1, $3, nil)}
;
```

**Pros:**
- User has explicit control
- No ambiguity when markers used
- Backward compatible (markers are optional)

**Cons:**
- Doesn't solve problem for unmarked inputs
- Users need to learn marker syntax
- Not suitable for natural language extraction use case

---

### Solution 5: Grammar Restriction + Post-Processing

**Approach:** Simplify grammar to remove ambiguity, then use post-processing to handle complex cases.

**Implementation:**

**Grammar Changes:**
```yacc
// Remove bare INT from Time in certain contexts
Time:
  INT Am {$$ = NewAMTime($1, nil, nil, nil)}
| INT Pm {$$ = NewPMTime($1, nil, nil, nil)}
| INT TimeSep INT {$$ = NewTime($1, $3, nil, nil)}
// NO: | INT {$$ = NewTime($1, nil, nil, nil)}
;

// Separate BareTime for explicit contexts
BareTime:
  INT {$$ = NewTime($1, nil, nil, nil)}
;

// Use BareTime only in unambiguous contexts
DateTimeRange:
| Date RangeSepPlus BareTime {$$ = NewRange(...)}  // "Feb 3 - 12" (in range context)
;
```

**Post-Processing:**
```go
func postProcess(parsed *DateTimeRanges) *DateTimeRanges {
    // Look for patterns that need fixup
    for i, item := range parsed.Items {
        // If we have single-date items that should be part of a list
        if isPartOfDateList(parsed, i) {
            // Merge into date list
            parsed.Items = mergeToDateList(parsed.Items, i)
        }
    }
    return parsed
}
```

**Pros:**
- Grammar becomes cleaner and less ambiguous
- Post-processing can use more sophisticated logic
- Easier to maintain and extend

**Cons:**
- Some patterns may not parse at all initially
- Post-processing logic can become complex
- Risk of introducing new bugs in post-processing

---

## Recommended Implementation Plan

### Phase 1: Low-Hanging Fruit (Already Done)
- ✅ Fix compilation errors
- ✅ Add missing ISO format support
- ✅ Add semantic validation
- ✅ Fix obvious grammar issues

### Phase 2: Context-Aware Improvements
1. **Enhance DayPlus rules** for better range support
   - Add support for "&" and "+" separators in day lists
   - Support "Day-Day & Day" patterns

2. **Improve prefix handling**
   - Add common prefix tokens ("Starts", "Date:", etc.)
   - Make them optional in DateTimeRanges rules

3. **Better timezone support**
   - Add region-based timezone tokens
   - Handle parenthetical clarifications

### Phase 3: Core Ambiguity Resolution (Choose One Approach)

**Recommended: Solution 2 (GLR Scoring)**

Rationale:
- Leverages existing GLR infrastructure
- Doesn't require lexer rewrite
- Can be incrementally improved
- Handles complex real-world cases

Implementation steps:
1. Add scoring interface to parse tree nodes
2. Implement basic scores for common patterns
3. Update GLR parser to use scores in tree selection
4. Iterate on scoring rules based on test results

### Phase 4: Edge Cases
- Handle abbreviated names
- Support count prefixes ("2 Wednesdays")
- Multiple times with single date
- Complex natural language patterns

---

## Testing Recommendations

### Current Test Status
- `acceptBrokenTests` flag set to `false` in `parse_test.go:16-18`
- All 28 failing tests marked with `isBroken: true`

### Testing Strategy
1. **Keep broken tests visible**: Continue running with `acceptBrokenTests = false` to track progress
2. **Add unit tests for each fix**: When fixing a pattern, add simpler unit tests
3. **Regression prevention**: Run full suite after every change
4. **Debug mode**: Use `DEBUG=true` environment variable for detailed parse traces

### Example Debug Command
```bash
env DEBUG=true go test -v ./datetime -run 'TestParse/121__Wednesdays,_February_1,_8,_15,_22$'
```

---

## Files Modified

### Core Changes
- `glr/glr_parse.go:481` - Fixed debugf format string
- `datetime/parse_yacc.y` - Multiple grammar improvements
- `datetime/datetime.go:177-199` - Added semantic validation

### Generated Files (Auto-regenerated)
- `datetime/parse_yacc.go`
- `datetime/parse.go`
- `datetime/parse_glr.go`

### Test Configuration
- `datetime/parse_test.go:16-18` - Changed acceptBrokenTests flag

---

## Conclusion

The datetime parser has fundamental ambiguity around bare integer interpretation (day vs hour). The fixes applied so far have:

1. Resolved build-blocking issues
2. Added missing format support
3. Improved date list parsing
4. Reduced ambiguity in specific cases

However, the core issue remains: **the grammar allows valid interpretations that are semantically incorrect in context**. Solving this requires one of:

- Context-aware lexing (complex but comprehensive)
- GLR scoring (moderate complexity, incremental improvement)
- Two-pass parsing (clean but slower)
- Grammar restriction + post-processing (simpler grammar, complex post-processing)

**Recommended next step:** Implement GLR scoring (Solution 2) to handle the date list vs datetime ambiguity, as it provides the best balance of effectiveness and maintainability.

---

## References

### Related Code Sections
- GLR parser implementation: `glr/glr_parse.go`
- Grammar specification: `datetime/parse_yacc.y`
- Date/Time constructors: `datetime/datetime.go:960-1042`
- Lexer: `datetime/parse_lex.go`
- Test cases: `datetime/parse_test.go`

### External Documentation
- GLR parsing: [Wikipedia - GLR parser](https://en.wikipedia.org/wiki/GLR_parser)
- Yacc/Bison: [GNU Bison Manual](https://www.gnu.org/software/bison/manual/)
- Go generate: [The Go Blog - Generating code](https://go.dev/blog/generate)

---

*Document generated: 2025-10-08*
*Parser version: commit e1b1455*
*Test suite: 225 tests, 28 failing (12.44%)*
