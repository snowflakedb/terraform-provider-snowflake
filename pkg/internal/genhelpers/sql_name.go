package genhelpers

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	splitSQLPattern   = regexp.MustCompile(`\s+|_`)
	englishLowerCaser = cases.Lower(language.English)
	englishTitleCaser = cases.Title(language.English)
)

// SqlToGoName converts a SQL name to a Go identifier using English title case.
// TIMESTAMP_DAY_IS_ALWAYS_24H becomes TimestampDayIsAlways24H.
// UNSUPPORTED_DDL_ACTION becomes UnsupportedDdlAction.
// When exported is false, the first word stays lowercase.
func SqlToGoName(sql string, exported bool) string {
	words := splitSQLPattern.Split(sql, -1)
	for i, word := range words {
		if word == "=" {
			words[i] = ""
			continue
		}
		if !exported && i == 0 {
			words[i] = englishLowerCaser.String(word)
			continue
		}
		words[i] = englishTitleCaser.String(word)
	}
	return strings.Join(words, "")
}
