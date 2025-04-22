package datatypes

import (
	"fmt"
	"strconv"
	"strings"
)

// TimestampTzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does have optional precision attribute.
type TimestampTzDataType struct {
	precision      int
	underlyingType string

	precisionKnown bool
}

func (t *TimestampTzDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
}

func (t *TimestampTzDataType) ToLegacyDataTypeSql() string {
	return TimestampTzLegacyDataType
}

func (t *TimestampTzDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", TimestampTzLegacyDataType, t.precision)
}

var TimestampTzDataTypeSynonyms = []string{TimestampTzLegacyDataType, "TIMESTAMPTZ", "TIMESTAMP WITH TIME ZONE"}

func parseTimestampTzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampTzDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		return &TimestampTzDataType{DefaultTimestampPrecision, raw.matchedByType, false}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return nil, fmt.Errorf(`timestamp tz %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
	}
	precisionRaw := r[1 : len(r)-1]
	precision, err := strconv.Atoi(strings.TrimSpace(precisionRaw))
	if err != nil {
		return nil, fmt.Errorf(`could not parse the timestamp's precision: "%s", err: %w`, precisionRaw, err)
	}
	return &TimestampTzDataType{precision, raw.matchedByType, true}, nil
}

func areTimestampTzDataTypesTheSame(a, b *TimestampTzDataType) bool {
	return a.precision == b.precision
}

func areTimestampTzDataTypesDefinitelyDifferent(a, b *TimestampTzDataType) bool {
	var precisionDefinitelyDifferent bool
	if a.precisionKnown && b.precisionKnown {
		precisionDefinitelyDifferent = a.precision != b.precision
	}
	return precisionDefinitelyDifferent
}
