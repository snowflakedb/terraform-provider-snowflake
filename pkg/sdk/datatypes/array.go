package datatypes

import (
	"fmt"
)

// ArrayDataType is based on:
// - https://docs.snowflake.com/en/sql-reference/data-types-semistructured#array (semi-structured ARRAY),
// - https://docs.snowflake.com/en/sql-reference/data-types-structured#label-structured-types-specifying-array (structured ARRAY).
// It does not have synonyms.
// Semi-structured ARRAY does not have any attributes.
// Structured ARRAY has the required element type and the optional NOT NULL constraint for the elements.
type ArrayDataType struct {
	elementType    DataType
	elementNotNull bool
	underlyingType string
}

func (t *ArrayDataType) ToSql() string {
	return t.format(t.underlyingType, func(dt DataType) string { return dt.ToSql() })
}

func (t *ArrayDataType) ToLegacyDataTypeSql() string {
	return ArrayLegacyDataType
}

func (t *ArrayDataType) Canonical() string {
	return t.format(ArrayLegacyDataType, func(dt DataType) string { return dt.Canonical() })
}

func (t *ArrayDataType) ToSqlWithoutUnknowns() string {
	return t.format(t.underlyingType, func(dt DataType) string { return dt.ToSqlWithoutUnknowns() })
}

func (t *ArrayDataType) IsStructured() bool {
	return t.elementType != nil
}

func (t *ArrayDataType) ElementType() DataType {
	return t.elementType
}

func (t *ArrayDataType) ElementNotNull() bool {
	return t.elementNotNull
}

func (t *ArrayDataType) format(arrayType string, formatElementType func(DataType) string) string {
	if !t.IsStructured() {
		return arrayType
	}
	return fmt.Sprintf("%s(%s)", arrayType, formatWithNotNull(formatElementType(t.elementType), t.elementNotNull))
}

var ArrayDataTypeSynonyms = []string{ArrayLegacyDataType}

// parseArrayDataTypeRaw returns the semi-structured ARRAY when no arguments are given
// and the structured ARRAY when the element type (with optional NOT NULL) is given in parentheses.
func parseArrayDataTypeRaw(raw sanitizedDataTypeRaw) (*ArrayDataType, error) {
	args, hasArguments, ok := structuredTypeArguments(raw)
	if !ok {
		return nil, fmt.Errorf(`array %s could not be parsed, use "%s" or "%s(element_type [ NOT NULL ])" format`, raw.raw, raw.matchedByType, raw.matchedByType)
	}
	if !hasArguments {
		return &ArrayDataType{underlyingType: raw.matchedByType}, nil
	}
	elementTypeRaw, elementNotNull := trimNotNullSuffix(args)
	if elementTypeRaw == "" {
		return nil, fmt.Errorf(`array %s could not be parsed, element type is missing; use "%s(element_type [ NOT NULL ])" format`, raw.raw, raw.matchedByType)
	}
	elementType, err := parseNestedDataType(elementTypeRaw)
	if err != nil {
		return nil, fmt.Errorf(`could not parse the array's element type: "%s", err: %w`, elementTypeRaw, err)
	}
	return &ArrayDataType{
		elementType:    elementType,
		elementNotNull: elementNotNull,
		underlyingType: raw.matchedByType,
	}, nil
}

func areArrayDataTypesTheSame(a, b *ArrayDataType) bool {
	if a.IsStructured() != b.IsStructured() {
		return false
	}
	if !a.IsStructured() {
		return true
	}
	return a.elementNotNull == b.elementNotNull && AreTheSame(a.elementType, b.elementType)
}

// arrays are different if:
// - one is structured and the other is semi-structured
// - NOT NULL constraint of the elements differs
// - element types are definitely different
func areArrayDataTypesDefinitelyDifferent(a, b *ArrayDataType) bool {
	if a.IsStructured() != b.IsStructured() {
		return true
	}
	if !a.IsStructured() {
		return false
	}
	return a.elementNotNull != b.elementNotNull || AreDefinitelyDifferent(a.elementType, b.elementType)
}
