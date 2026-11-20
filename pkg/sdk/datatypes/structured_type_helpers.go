package datatypes

import (
	"fmt"
	"regexp"
	"strings"
)

// Common parts of the structured types based on https://docs.snowflake.com/en/sql-reference/data-types-structured.

var notNullSuffix = regexp.MustCompile(`(?i)\s+NOT\s+NULL$`)

// structuredTypeArguments extracts the content of the parentheses following the matched type.
// It returns hasArguments=false when there are no parentheses (the semi-structured variant of the type)
// and ok=false when the input after the matched type is not wrapped in parentheses.
func structuredTypeArguments(raw sanitizedDataTypeRaw) (args string, hasArguments bool, ok bool) {
	r := strings.TrimSpace(raw.raw[len(raw.matchedByType):])
	if r == "" {
		return "", false, true
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return "", true, false
	}
	return strings.TrimSpace(r[1 : len(r)-1]), true, true
}

// trimNotNullSuffix removes the trailing NOT NULL constraint and reports if it was present.
func trimNotNullSuffix(s string) (string, bool) {
	if loc := notNullSuffix.FindStringIndex(s); loc != nil {
		return strings.TrimSpace(s[:loc[0]]), true
	}
	return s, false
}

func formatWithNotNull(s string, notNull bool) string {
	if notNull {
		return s + " NOT NULL"
	}
	return s
}

// parseNestedDataType parses the data type nested in the structured type (e.g. array's element or object's value).
func parseNestedDataType(raw string) (DataType, error) {
	dataType, err := ParseDataType(raw)
	if err != nil {
		return nil, err
	}
	if _, ok := dataType.(*TableDataType); ok {
		return nil, fmt.Errorf("table cannot be nested in structured types: %s", raw)
	}
	return dataType, nil
}
