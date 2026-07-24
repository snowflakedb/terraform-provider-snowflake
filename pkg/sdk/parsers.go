package sdk

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// ParseCommaSeparatedStringArray can be used to parse Snowflake output containing a list in the format of "[item1, item2, ...]",
// the assumptions are that:
// 1. The list may be enclosed by outer [] brackets which are stripped. Brackets shouldn't be a part of any item's value
// 2. Items are separated by commas. Commas enclosed in double quotes are treated as part of the item's value rather than separators (e.g. `"a,b"` is a single item)
// 3. Items can have as many spaces in between, but after separation they will be trimmed and shouldn't be a part of any item's value
func ParseCommaSeparatedStringArray(value string, trimQuotes bool) []string {
	value = stripOuterBrackets(value)
	if value == "" {
		return make([]string, 0)
	}
	listItems := splitCommaSeparated(value)
	return applyFormatting(listItems, trimQuotes)
}

// ParseOuterCommaSeparatedStringArray can be used to parse Snowflake output containing a list in the format of "[item1, item2, ...]".
// It also supports nested lists like "[[item1, item2], [item3, item4]]".
// The assumptions are that:
// 1. The list may be enclosed by outer [] brackets which are stripped
// 2. Brackets have special meaning (they denote list boundaries) and shouldn't be a part of any item's value
// 3. Items are separated by top-level commas (commas inside nested brackets or double quotes are preserved)
// 4. Items can have as many spaces in between, but after separation they will be trimmed and shouldn't be a part of any item's value
// 5. Quote trimming only applies to top-level items. Nested items are returned as a raw string
func ParseOuterCommaSeparatedStringArray(value string, trimQuotes bool) []string {
	value = stripOuterBrackets(value)
	if value == "" {
		return make([]string, 0)
	}
	listItems := splitOuter(value)
	return applyFormatting(listItems, trimQuotes)
}

// ParseCommaSeparatedSchemaObjectIdentifierArray can be used to parse Snowflake output containing a list of schema-level object identifiers
// in the format of ["db".SCHEMA."name", "db"."schema2"."name2", ...],
func ParseCommaSeparatedSchemaObjectIdentifierArray(value string) ([]SchemaObjectIdentifier, error) {
	return collections.MapErr(ParseCommaSeparatedStringArray(value, false), ParseSchemaObjectIdentifier)
}

// ParseCommaSeparatedAccountIdentifierArray can be used to parse Snowflake output containing a list of account identifiers
// in the format of ["organization1.account1", "organization2.account2", ...],
func ParseCommaSeparatedAccountIdentifierArray(value string) ([]AccountIdentifier, error) {
	return collections.MapErr(ParseCommaSeparatedStringArray(value, false), ParseAccountIdentifier)
}

// ParseCommaSeparatedAccountObjectIdentifierArray can be used to parse Snowflake output containing a list of account object identifiers
// in the format of ["object1", "object2", ...],
func ParseCommaSeparatedAccountObjectIdentifierArray(value string) ([]AccountObjectIdentifier, error) {
	return collections.MapErr(ParseCommaSeparatedStringArray(value, false), ParseAccountObjectIdentifier)
}

func stripOuterBrackets(value string) string {
	value = strings.TrimPrefix(value, "[")
	return strings.TrimSuffix(value, "]")
}

func applyFormatting(listItems []string, trimQuotes bool) []string {
	trimmedListItems := make([]string, len(listItems))
	for i, item := range listItems {
		trimmedListItems[i] = strings.TrimSpace(item)
		if trimQuotes {
			trimmedListItems[i] = strings.Trim(trimmedListItems[i], "'\"")
		}
	}
	return trimmedListItems
}

func splitOuter(value string) []string {
	depth := 0
	idx := 0
	inQuotes := false
	var parts []string
	for i, ch := range value {
		switch ch {
		case '"':
			inQuotes = !inQuotes
		case '[':
			if !inQuotes {
				depth++
			}
		case ']':
			if !inQuotes {
				depth--
			}
		case ',':
			if !inQuotes && depth <= 0 {
				parts = append(parts, value[idx:i])
				idx = i + 1
			}
		}
	}
	return append(parts, value[idx:])
}

func emptyIfNull(s string) string {
	if s == "null" {
		return ""
	}
	return s
}

// splitCommaSeparated splits a comma-separated string, treating commas enclosed in double
// quotes as part of a value rather than separators. Snowflake identifiers may contain commas
// when quoted (e.g. `DB.SCHEMA."a,b"`), so a naive strings.Split(",") would break such values
// apart. The returned parts are raw (not trimmed or parsed).
func splitCommaSeparated(s string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	for _, r := range s {
		switch {
		case r == '"':
			inQuotes = !inQuotes
			current.WriteRune(r)
		case r == ',' && !inQuotes:
			parts = append(parts, current.String())
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	return append(parts, current.String())
}
