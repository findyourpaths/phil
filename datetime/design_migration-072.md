# Authored ranges and bounded weekly schedules
Status: ACTIVE
Kind: REFERENCE

## Contract

The grammar remains the syntax authority. Propagate a date range's authored year before constructing resolved date/time endpoints. A complete sequence of explicit-year month/day ranges consumes each trailing year locally, rather than allowing a preceding range's year to seed the next range.

Separate `Dates` and `Time` labels can describe one complete expression. Compose raw date bounds first, then combine their dates with the supplied clock endpoints and shared zone. This requires neither a second parser nor a caller-specific string rewrite.

An explicitly plural weekday following bounded dates and clock hours denotes a weekly series. Store the first stated date with the session's clock duration in its first range; retain the inclusive final date in `Recurrence.Until`. A singular weekday is not authority for this recurrence production. `Times` is lexical normalization of the existing `Time` label, not a substitute for missing grammar.

## Real-input verification

The complete ccisg block, Self-Led detail range and IFS training #1709 schedule retain their authored endpoint years (/Users/wag/Dropbox/Projects/phil/datetime/migration_072_test.go:27-49). Combined Self-Led dates/hours and bounded Roxannemanning Thursday sessions preserve their clock duration and timezone (/Users/wag/Dropbox/Projects/phil/datetime/migration_072_test.go:83; /Users/wag/Dropbox/Projects/phil/datetime/recurrence_072_test.go:11).

Verified by: `TestMigration072ExplicitRangeYears`, `TestMigration072IFSSelfLedDatesAndTime`, `TestMigration072RoxannemanningBoundedThursdays`, and the existing `TestParse` suite. Paths additionally verifies its unchanged caller and existing civil Schedule conversion with `TestMigration072SavedScheduleComposition`.

## Block-scoped shared-year days

The saved IFS training #1703 block ends in `Bridge days: March 20 and April 23, 2027`. Both bridge days belong to the same advertised year; a block window must not stop at the second month and accept the first day under a 2026 default. `parseBlockSharedYearDays` admits only the bounded six-word, two-named-month conjunction with a trailing year, then delegates the complete expression to the already existing grammar. A full valid result is required before adopting the group. Ordinary next-anchor boundaries, label-year context and completeness-triggered parsing remain unchanged for all other input shapes; no year is carried across unrelated entries or arbitrary prose.

Verified by: `TestMigration072IFS1703BridgeDays`, including the full four-module schedule and both date-only bridge days with default year zero and 2026; the existing `TestParse` and block-parser suites remain compatibility oracles.

## Date-labelled clock windows

The saved Equitable masthead has a year-bearing summary and four explicitly dated Times rows, while the body advertises their Pacific zone. The copied Source owns representation cleanup: retain the observed year as a block label, keep the dated clock rows, and select the body's actual zone. Phil's existing block window then preserves `Month Day: clock-clock [zone]` as one bounded grammar expression before accepting its shorter default-complete date prefix. Reuse the existing clock-range recognizer and delegate date/time syntax to `Parse`; require complete start/end clocks and reject trailing prose or an unresolved zone token. Empty Times placeholders remain ordinary date-only evidence. No event identifier, fixed session count, custom date parser or grammar regeneration is introduced.

Verified by: `TestMigration072EquitableDatedClockWindows` under matching and different reference-year defaults, `TestMigration072DatedClockWindowDoesNotInventHours`, existing parser/block suites, and Paths `TestMigration072EquitableDatedClockComposition` plus the saved four-detail replay.

## Comma-separated multi-month day lists

The saved Aephoria Professional Pairs input is `27 October, 4, 9, 12 November 2026`. The existing `DayPlus Month AND DayPlus Month Year` grammar already composes both groups with the explicit year; its conjunction control passes before repair. Extend only the existing comma normalization to admit a bounded second comma-separated day list (at most 31 day tokens), preserving the existing month vocabulary and grammar result selection. This is demonstrated punctuation normalization, not a new syntax parser or generated grammar change. The four dates are October 27 and November 4, 9 and 12, 2026; none advertises a clock, end or zone. A fully dated but wrong prior parse is not acceptance evidence. Paths retains its completeness-triggered block fallback unchanged and documents this lexical correction separately.

Verified by: `TestMigration072AephoriaMultiMonthDays`, including the exact saved input, existing conjunction control, matching and conflicting default years, four distinct date-only items and no invented clock/end/zone. The existing parser and block suites remain compatibility oracles.

## Generated actions and compatibility

Semantic grammar actions must be single result expressions that the GLR generator can translate. Multi-statement composition belongs in private semantic constructors, not hand-edited generated code. Regenerate `parse_yacc.go`, `parse_glr.go` and `parse_yacc.states.txt` with `make generate`; generated artifacts are read-only.

The complete explicit-year-list production requires at least two complete ranges; it is not an arbitrary partial-date shortcut. Date/time composition and bounded weekly parsing are independent of the earlier range-year increment and each must preserve existing parser results. No public API or Schedule schema is changed. Missing timezone behavior and ordinary parse-error handling remain unchanged.
