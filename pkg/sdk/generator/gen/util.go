package gen

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// Split by empty space or underscore
	splitSQLPattern   = regexp.MustCompile(`\s+|_`)
	englishLowerCaser = cases.Lower(language.English)
	englishTitleCaser = cases.Title(language.English)
)

// IsNil is used for special cases where x != nil might not work (e.g. passing nil instead of interface implementation)
func IsNil(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}

	return false
}

func startingWithLowerCase(s string) string {
	firstLetter, _ := utf8.DecodeRuneInString(s)
	return strings.ToLower(string(firstLetter)) + s[1:]
}

func startingWithUpperCase(s string) string {
	firstLetter, _ := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(firstLetter)) + s[1:]
}

func wrapWith(s string, with string) string {
	return fmt.Sprintf("%s%s%s", with, s, with)
}

func sqlToFieldName(sql string, shouldExport bool) string {
	sqlWords := splitSQLPattern.Split(sql, -1)
	for i, s := range sqlWords {
		if s == "=" {
			sqlWords[i] = ""
			continue
		}
		if !shouldExport && i == 0 {
			sqlWords[i] = englishLowerCaser.String(s)
			continue
		}
		sqlWords[i] = englishTitleCaser.String(s)
	}
	return strings.Join(sqlWords, "")
}
