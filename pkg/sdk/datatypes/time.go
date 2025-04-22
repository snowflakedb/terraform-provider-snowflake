package datatypes

import (
	"fmt"
	"strconv"
	"strings"
)

const DefaultTimePrecision = 9

// TimeDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#time
// It does not have synonyms. It does have optional precision attribute.
type TimeDataType struct {
	precision      int
	underlyingType string

	precisionKnown bool
}

func (t *TimeDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
}

func (t *TimeDataType) ToLegacyDataTypeSql() string {
	return TimeLegacyDataType
}

func (t *TimeDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", TimeLegacyDataType, t.precision)
}

func (t *TimeDataType) ToSqlWithoutUnknowns() string {
	switch {
	case t.precisionKnown:
		return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
	default:
		return fmt.Sprintf("%s", t.underlyingType)
	}
}

var TimeDataTypeSynonyms = []string{TimeLegacyDataType}

func parseTimeDataTypeRaw(raw sanitizedDataTypeRaw) (*TimeDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		return &TimeDataType{DefaultTimePrecision, raw.matchedByType, false}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return nil, fmt.Errorf(`time %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
	}
	precisionRaw := r[1 : len(r)-1]
	precision, err := strconv.Atoi(strings.TrimSpace(precisionRaw))
	if err != nil {
		return nil, fmt.Errorf(`could not parse the time's precision: "%s", err: %w`, precisionRaw, err)
	}
	return &TimeDataType{precision, raw.matchedByType, true}, nil
}

func areTimeDataTypesTheSame(a, b *TimeDataType) bool {
	return a.precision == b.precision
}

func areTimeDataTypesDefinitelyDifferent(a, b *TimeDataType) bool {
	var precisionDefinitelyDifferent bool
	if a.precisionKnown && b.precisionKnown {
		precisionDefinitelyDifferent = a.precision != b.precision
	}
	return precisionDefinitelyDifferent
}
