package datatypes

import (
	"fmt"
	"strconv"
	"strings"
)

// TimestampLtzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does have optional precision attribute.
type TimestampLtzDataType struct {
	precision      int
	underlyingType string

	precisionKnown bool
}

func (t *TimestampLtzDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
}

func (t *TimestampLtzDataType) ToLegacyDataTypeSql() string {
	return TimestampLtzLegacyDataType
}

func (t *TimestampLtzDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", TimestampLtzLegacyDataType, t.precision)
}

func (t *TimestampLtzDataType) ToSqlWithoutUnknowns() string {
	switch {
	case t.precisionKnown:
		return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
	default:
		return fmt.Sprintf("%s", t.underlyingType)
	}
}

var TimestampLtzDataTypeSynonyms = []string{TimestampLtzLegacyDataType, "TIMESTAMPLTZ", "TIMESTAMP WITH LOCAL TIME ZONE"}

func parseTimestampLtzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampLtzDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		return &TimestampLtzDataType{DefaultTimestampPrecision, raw.matchedByType, false}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return nil, fmt.Errorf(`timestamp ltz %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
	}
	precisionRaw := r[1 : len(r)-1]
	precision, err := strconv.Atoi(strings.TrimSpace(precisionRaw))
	if err != nil {
		return nil, fmt.Errorf(`could not parse the timestamp's precision: "%s", err: %w`, precisionRaw, err)
	}
	return &TimestampLtzDataType{precision, raw.matchedByType, true}, nil
}

func areTimestampLtzDataTypesTheSame(a, b *TimestampLtzDataType) bool {
	return a.precision == b.precision
}

func areTimestampLtzDataTypesDefinitelyDifferent(a, b *TimestampLtzDataType) bool {
	var precisionDefinitelyDifferent bool
	if a.precisionKnown && b.precisionKnown {
		precisionDefinitelyDifferent = a.precision != b.precision
	}
	return precisionDefinitelyDifferent
}
