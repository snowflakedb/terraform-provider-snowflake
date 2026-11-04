package genhelpers

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var splitSQLPattern = regexp.MustCompile(`\s+|_`)

// SqlToGoName converts a SQL name to a Go identifier using English title case.
// TIMESTAMP_DAY_IS_ALWAYS_24H becomes TimestampDayIsAlways24H.
// UNSUPPORTED_DDL_ACTION becomes UnsupportedDdlAction.
// When exported is false, the first word stays lowercase.
//
// A cases.Caser keeps word state and cannot be shared between goroutines.
// Generators convert names for several objects at once, so each call builds its own.
func SqlToGoName(sql string, exported bool) string {
	lower := cases.Lower(language.English)
	title := cases.Title(language.English)
	words := splitSQLPattern.Split(sql, -1)
	for i, word := range words {
		if word == "=" {
			words[i] = ""
			continue
		}
		if !exported && i == 0 {
			words[i] = lower.String(word)
			continue
		}
		words[i] = title.String(word)
	}
	return strings.Join(words, "")
}
