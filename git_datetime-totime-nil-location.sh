#!/bin/bash
# Fix panic in DateTime.ToTime() when Location is nil
#
# Problem: time.Date() panics with "time: missing Location" when
# dt.Location() returns nil (no timezone set).
#
# Fix: Default to UTC when no location is available.

cd /Users/wag/Dropbox/Projects/phil

git add datetime/datetime.go

git commit -m "$(cat <<'EOF'
fix: Handle nil Location in DateTime.ToTime() by defaulting to UTC

Go 1.25+ requires a non-nil location for time.Date(). When DateTime
has no TimeZone set, Location() returns nil, causing a panic.

Fix: Default to time.UTC when Location() returns nil. This is a safe
default for event times without explicit timezone information.

Also includes:
- IANAName() methods for DateTime, DateTimeRange, TimeZone
- NewRange validation to reject year-only to day-only ranges
- Debug flag for Location() warning message

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"

echo "Committed phil datetime changes"
