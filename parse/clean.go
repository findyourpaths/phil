package parse

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var whitespacesRE = regexp.MustCompile(`\s+`)

func CleanTextLine(s string) string {
	r := bluemonday.StrictPolicy().AddSpaceWhenStrippingTag(true).Sanitize(s)
	r = strings.ReplaceAll(r, "\u00a0", " ")
	r = whitespacesRE.ReplaceAllString(r, " ")
	r = strings.TrimSpace(r)
	return r
}
