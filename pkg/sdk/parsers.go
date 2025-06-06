package sdk

import (
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// fix timestamp merge
func ParseTimestampWithOffset(s string, dateTimeFormat string) (string, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err.Error(), err
	}
	_, offset := t.Zone()
	adjustedTime := t.Add(-time.Duration(offset) * time.Second)
	adjustedTimeFormat := adjustedTime.Format(dateTimeFormat)
	return adjustedTimeFormat, nil
}

// ParseCommaSeparatedStringArray can be used to parse Snowflake output containing a list in the format of "[item1, item2, ...]",
// the assumptions are that:
// 1. The list is enclosed by [] brackets, and they shouldn't be a part of any item's value
// 2. Items are separated by commas, and they shouldn't be a part of any item's value
// 3. Items can have as many spaces in between, but after separation they will be trimmed and shouldn't be a part of any item's value
func ParseCommaSeparatedStringArray(value string, trimQuotes bool) []string {
	value = strings.Trim(value, "[]")
	if value == "" {
		return make([]string, 0)
	}
	listItems := strings.Split(value, ",")
	trimmedListItems := make([]string, len(listItems))
	for i, item := range listItems {
		trimmedListItems[i] = strings.TrimSpace(item)
		if trimQuotes {
			trimmedListItems[i] = strings.Trim(trimmedListItems[i], "'\"")
		}
	}
	return trimmedListItems
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
