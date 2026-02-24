package datatypes

import (
	"fmt"
	"strconv"
	"strings"
)

const DefaultDecfloatPrecision = 38

// DecfloatDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-numeric#decfloat
// It doesn't have synonyms. It does have optional precision attribute.
// Currently, the only allowed value is 38. To be consistent with other data types we are not validating this value (e.g., it's possible to set number's precision to 40 and it will fail on Snowflake).
// Precision can be known or unknown.
type DecfloatDataType struct {
	precision      int
	underlyingType string

	precisionKnown bool
}

func (t *DecfloatDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
}

func (t *DecfloatDataType) ToLegacyDataTypeSql() string {
	return DecfloatLegacyDataType
}

func (t *DecfloatDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", DecfloatLegacyDataType, t.precision)
}

func (t *DecfloatDataType) ToSqlWithoutUnknowns() string {
	switch {
	case t.precisionKnown:
		return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
	default:
		return fmt.Sprintf("%s", t.underlyingType)
	}
}

var DecfloatDataTypeSynonyms = []string{DecfloatLegacyDataType}

func parseDecfloatDataTypeRaw(raw sanitizedDataTypeRaw) (*DecfloatDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		return &DecfloatDataType{DefaultDecfloatPrecision, raw.matchedByType, false}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return nil, fmt.Errorf(`decfloat %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
	}
	precisionRaw := r[1 : len(r)-1]
	precision, err := strconv.Atoi(strings.TrimSpace(precisionRaw))
	if err != nil {
		return nil, fmt.Errorf(`could not parse the decfloat's precision: "%s", err: %w`, precisionRaw, err)
	}
	return &DecfloatDataType{precision, raw.matchedByType, true}, nil
}

func areDecfloatDataTypesTheSame(a, b *DecfloatDataType) bool {
	return a.precision == b.precision
}

func areDecfloatDataTypesDefinitelyDifferent(a, b *DecfloatDataType) bool {
	var precisionDefinitelyDifferent bool
	if a.precisionKnown && b.precisionKnown {
		precisionDefinitelyDifferent = a.precision != b.precision
	}
	return precisionDefinitelyDifferent
}
