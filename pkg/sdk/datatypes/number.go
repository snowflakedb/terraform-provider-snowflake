package datatypes

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	DefaultNumberPrecision = 38
	DefaultNumberScale     = 0
)

// NumberDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-numeric#data-types-for-fixed-point-numbers
// It does have synonyms that allow specifying precision and scale; here called synonyms.
// It does have synonyms that does not allow specifying precision and scale; here called subtypes.
type NumberDataType struct {
	precision      int
	scale          int
	underlyingType string

	precisionKnown bool
	scaleKnown     bool
}

func (t *NumberDataType) ToSql() string {
	if slices.Contains(NumberDataTypeSubTypes, t.underlyingType) {
		return t.underlyingType
	}
	return fmt.Sprintf("%s(%d, %d)", t.underlyingType, t.precision, t.scale)
}

func (t *NumberDataType) ToLegacyDataTypeSql() string {
	return NumberLegacyDataType
}

func (t *NumberDataType) Canonical() string {
	return fmt.Sprintf("%s(%d,%d)", NumberLegacyDataType, t.precision, t.scale)
}

func (t *NumberDataType) ToSqlWithoutUnknowns() string {
	switch {
	case slices.Contains(NumberDataTypeSubTypes, t.underlyingType):
		return t.underlyingType
	case t.precisionKnown && t.scaleKnown:
		return fmt.Sprintf("%s(%d, %d)", t.underlyingType, t.precision, t.scale)
	case t.precisionKnown:
		return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
	default:
		return fmt.Sprintf("%s", t.underlyingType)
	}
}

var (
	NumberDataTypeSynonyms = []string{NumberLegacyDataType, "DECIMAL", "DEC", "NUMERIC"}
	NumberDataTypeSubTypes = []string{"INTEGER", "INT", "BIGINT", "SMALLINT", "TINYINT", "BYTEINT"}
	AllNumberDataTypes     = append(NumberDataTypeSynonyms, NumberDataTypeSubTypes...)
)

func parseNumberDataTypeRaw(raw sanitizedDataTypeRaw) (*NumberDataType, error) {
	switch {
	case slices.Contains(NumberDataTypeSynonyms, raw.matchedByType):
		return parseNumberDataTypeWithPrecisionAndScale(raw)
	case slices.Contains(NumberDataTypeSubTypes, raw.matchedByType):
		return parseNumberDataTypeWithoutPrecisionAndScale(raw)
	default:
		return nil, fmt.Errorf("unknown number data type: %s", raw.raw)
	}
}

// parseNumberDataTypeWithPrecisionAndScale extracts precision and scale from the raw number data type input.
// It returns defaults if no arguments were provided. It returns error if any part is not parseable.
func parseNumberDataTypeWithPrecisionAndScale(raw sanitizedDataTypeRaw) (*NumberDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		return &NumberDataType{DefaultNumberPrecision, DefaultNumberScale, raw.matchedByType, false, false}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return nil, fmt.Errorf(`number %s could not be parsed, use "%s(precision, scale)" format`, raw.raw, raw.matchedByType)
	}
	onlyArgs := r[1 : len(r)-1]
	parts := strings.Split(onlyArgs, ",")
	switch l := len(parts); l {
	case 1:
		precision, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf(`could not parse the number's precision: "%s", err: %w`, parts[0], err)
		}
		return &NumberDataType{precision, DefaultNumberScale, raw.matchedByType, true, false}, nil
	case 2:
		precision, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf(`could not parse the number's precision: "%s", err: %w`, parts[0], err)
		}
		scale, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf(`could not parse the number's scale: "%s", err: %w`, parts[1], err)
		}
		return &NumberDataType{precision, scale, raw.matchedByType, true, true}, nil
	default:
		return nil, fmt.Errorf(`number cannot have %d arguments: "%s"; only precision and scale are allowed`, l, onlyArgs)
	}
}

func parseNumberDataTypeWithoutPrecisionAndScale(raw sanitizedDataTypeRaw) (*NumberDataType, error) {
	if raw.raw != raw.matchedByType {
		args := strings.TrimPrefix(raw.raw, raw.matchedByType)
		return nil, fmt.Errorf("number type %s cannot have arguments: %s", raw.matchedByType, args)
	} else {
		return &NumberDataType{DefaultNumberPrecision, DefaultNumberScale, raw.matchedByType, false, false}, nil
	}
}

func areNumberDataTypesTheSame(a, b *NumberDataType) bool {
	return a.precision == b.precision && a.scale == b.scale
}

func areNumberDataTypesDefinitelyDifferent(a, b *NumberDataType) bool {
	var precisionDefinitelyDifferent bool
	if a.precisionKnown && b.precisionKnown {
		precisionDefinitelyDifferent = a.precision != b.precision
	}

	var scaleDefinitelyDifferent bool
	if a.scaleKnown && b.scaleKnown {
		scaleDefinitelyDifferent = a.scale != b.scale
	}

	return precisionDefinitelyDifferent || scaleDefinitelyDifferent
}
