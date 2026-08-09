# datetime

Natural language date/time parser using a GLR parser with a yacc-generated grammar.

## Key Types/Functions

- `Parse(minDateTime, dateMode, input) (*DateTimeRanges, error)` — main entry point; parses a string into structured date/time ranges
- `DateTimeRanges` — top-level result: a list of `DateTimeRange` items with optional `Recurrence`
- `DateTimeRange` — a start/end `DateTime` pair
- `DateTime` — a `Date` + `Time` + `TimeZone`
- `datetimeLexer` — tokenizer that preprocesses input (boundary splitting, weekday stripping, IANA tz extraction) then scans tokens for the GLR parser

## Architecture

```
input string
    │
    ▼
datetimeLexer (parse_lex.go)
  - HTML unescape, strip markdown bold
  - Boundary splitting: "12pm" → "12 pm", "2023-02-03" → "2023 - 02 - 03"
  - Weekday plural prefix stripping → Recurrence
  - IANA timezone extraction
  - Token classification: YEAR, INT, MONTH_NAME, AM/PM, AT, SUB, etc.
    │
    ▼
GLR parser (../glr/glr_parse.go)
  - Reads grammar rules + parse states from generated code
  - Explores ALL shift/reduce paths in parallel (handles ambiguity)
  - Returns parse trees ranked by numTerms (more tokens consumed = better)
    │
    ▼
Semantic actions (parse_yacc.y → parse_yacc.go)
  - Each grammar rule has a Go action that builds DateTime structs
  - Called via reflection by glr.GetParseNodeValue
    │
    ▼
Parse() result selection (parse.go)
  - Iterates roots (sorted by numTerms desc), takes first without error
  - Resolves missing years via minimumDateTime
  - Caches results by (minDateTime, dateMode, input)
```

## Files

| File | Role | Editable? |
|------|------|-----------|
| `parse_yacc.y` | Grammar (yacc) — the source of truth for what date formats are recognized | YES |
| `parse_lex.go` | Lexer — tokenization and input preprocessing | YES |
| `parse.go` | `Parse()` entry point, result selection, caching | YES |
| `datetime.go` | Domain types (`DateTimeRanges`, `Date`, `Time`, etc.) and constructors | YES |
| `generate.go` | `//go:generate` directives | YES |
| `parse_yacc.go` | Generated from `parse_yacc.y` by goyacc | NO — `go generate` |
| `parse_glr.go` | Generated GLR tables from yacc grammar + states | NO — `go generate` |

## Grammar Structure (parse_yacc.y)

The grammar is a hierarchy of nonterminals. Understanding this hierarchy is essential for debugging:

```
root
 └─ DateTimeRanges
     └─ DateTimeRange          (single or start-end pair)
         ├─ DateTime            (Date + Time + TimeZone)
         │   ├─ Date            (many formats: "Feb 3 2023", "2023-02-03", "3rd Feb", etc.)
         │   │   └─ RFC3339Date (Year SUB INT SUB INT — "2023-02-03")
         │   ├─ Time            (INT:INT AM/PM, TimePrefixPlus Time, etc.)
         │   │   └─ TimePrefix: AT | "TIME :"
         │   ├─ TimeZone        (abbreviation or name)
         │   └─ RFC3339DateTime (RFC3339Date + RFC3339Time + RFC3339TimeZone)
         │       └─ RFC3339Time (T INT:INT:INT — strict ISO format)
         ├─ RangeSep            (SUB, TO, THROUGH, TILL, UNTIL, DEC)
         └─ DateTimeSep         (COLON, COMMA, DEC, QUO, SUB, T)
```

Key design decisions:
- **`RFC3339Date` is a `Date` alternative** (line 455): ISO dates like `2023-02-03` reduce to `Date`, so ALL existing `Date Time` rules work automatically (`Date Time TimeZoneOpt`, `Date INT Am`, etc.)
- **`RFC3339DateTime` is a `DateTime` alternative** (line 363): for strict RFC3339 (`2023-02-03T12:00:00Z`)
- **Time-only rules live at `DateTimeRange` level** (lines 273-285), not `DateTime`, to avoid GLR ambiguity
- **Explicit `Date INT Am/Pm` inlines** (lines 353-356): yacc reduces INT to Day after Date, blocking `Time: INT Am`. Inline rules bypass this.
- **`DateSep` excludes `SUB`** but `DateTimeSep` includes it: prevents `2023-02` from being parsed as `Year DateSep Day` (ambiguous with range separator)

## Conventions

- **Regenerate after grammar changes**: `cd phil/datetime && go generate -v ./...` — runs goyacc then glr-generate
- **Conflicts are expected**: the GLR parser handles shift/reduce and reduce/reduce conflicts by exploring all paths; conflict counts in `go generate` output are informational
- **"never reduced" warnings** for `*PrefixPlus` rules are expected (the Plus variant handles repetition)
- **Test naming**: tests are numbered `TestParse/NNN__input_text`; run specific tests with `go test -v -run 'TestParse/192'`
- **Debug a parse**: `DEBUG=true go test -v -run 'TestParse/192' ./phil/datetime/` — prints full GLR trace (token shifts, reductions, parser forks, result ranking)
- **`skip` field in test cases**: if non-empty, test runs parse but skips on failure; fatals with "REMOVE SKIP" if it unexpectedly passes
- **Global state**: `minimumDateTime` and `parseDateMode` are protected by `parseMutex` — Parse() is goroutine-safe but serialized
- **Backwards ranges fail parsing**: when one range has complete start/end dates and its end date precedes its start date, `Parse` returns a semantic error. Separate items may be mentioned out of order, and an overnight time ending earlier on a later date remains valid; callers own ordinary parse-failure recovery rather than repairing dates here.

## Debugging a Parse Failure

When a date string fails to parse or returns wrong results:

1. **Add a test case** in `parse_test.go` with the input and expected output
2. **Run with DEBUG**: `DEBUG=true go test -v -run 'TestParse/NNN'` — check the token stream ("input after processing") and the GLR trace
3. **Check the lexer** (`parse_lex.go`): is the input tokenized correctly? Look at "scanned literal" lines. Common issues: `@` not lexed as `AT`, boundary splitting wrong
4. **Check the grammar** (`parse_yacc.y`): does a rule exist for this date format? Trace the nonterminal hierarchy — which `Date` rule matches? Which `DateTime` rule combines Date + Time?
5. **Check result ranking**: the GLR parser returns multiple results sorted by `numTerms` (tokens consumed). If a partial parse ranks above the full parse, the wrong result is selected. Look for "result[N]" lines at the end of the debug trace — `score=[numTerms, size, depths]`

Common root causes:
- **Missing grammar rule**: a date format has no path through the nonterminal hierarchy (e.g., RFC3339Date + regular Time was missing before adding `Date: RFC3339Date`)
- **Wrong token type**: lexer assigns wrong token (e.g., `@` as `AMP` instead of `AT`)
- **Partial parse wins**: GLR accepts a short parse; full parse exists but as a lower-ranked partial result — indicates a missing grammar path for the full parse, NOT a GLR bug

## Design Docs

- [ISSUES.md](ISSUES.md) — known parser issues
- [../glr/](../glr/) — GLR parser implementation (`glr_parse.go`)
- [../AGENTS.md](../AGENTS.md) — Phil project context
