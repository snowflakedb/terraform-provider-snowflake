package datatypes

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// DataType is the common interface that represents all Snowflake datatypes documented in https://docs.snowflake.com/en/sql-reference/intro-summary-data-types.
type DataType interface {
	// ToSql formats data type explicitly specifying all arguments and using the given type (e.g. CHAR(29) for CHAR(29)).
	ToSql() string
	// ToLegacyDataTypeSql formats data type using its base type without any attributes (e.g. VARCHAR for CHAR(29)).
	ToLegacyDataTypeSql() string
	// Canonical formats the data type between ToSql and ToLegacyDataTypeSql: it uses base type but with arguments (e.g. VARCHAR(29) for CHAR(29)).
	Canonical() string
	// ToSqlWithoutUnknowns formats data type explicitly specifying all KNOWN arguments and using the given type.
	ToSqlWithoutUnknowns() string
}

type sanitizedDataTypeRaw struct {
	raw           string
	matchedByType string
}

// ParseDataType is the entry point to get the implementation of the DataType from input raw string.
func ParseDataType(raw string) (DataType, error) {
	dataTypeRaw := strings.TrimSpace(strings.ToUpper(raw))
	namePart, _, hasArgs := strings.Cut(dataTypeRaw, "(")
	name := strings.TrimSpace(namePart)
	r := sanitizedDataTypeRaw{dataTypeRaw, name}

	switch {
	case slices.Contains(DecfloatDataTypeSynonyms, name):
		return parseDecfloatDataTypeRaw(r)
	case slices.Contains(AllNumberDataTypes, name):
		return parseNumberDataTypeRaw(r)
	case slices.Contains(FloatDataTypeSynonyms, name) && !hasArgs:
		return parseFloatDataTypeRaw(r)
	case slices.Contains(AllTextDataTypes, name):
		return parseTextDataTypeRaw(r)
	case slices.Contains(BinaryDataTypeSynonyms, name):
		return parseBinaryDataTypeRaw(r)
	case slices.Contains(BooleanDataTypeSynonyms, name) && !hasArgs:
		return parseBooleanDataTypeRaw(r)
	case slices.Contains(DateDataTypeSynonyms, name) && !hasArgs:
		return parseDateDataTypeRaw(r)
	case slices.Contains(TimestampLtzDataTypeSynonyms, name):
		return parseTimestampLtzDataTypeRaw(r)
	case slices.Contains(TimestampNtzDataTypeSynonyms, name):
		return parseTimestampNtzDataTypeRaw(r)
	case slices.Contains(TimestampTzDataTypeSynonyms, name):
		return parseTimestampTzDataTypeRaw(r)
	case slices.Contains(TimeDataTypeSynonyms, name):
		return parseTimeDataTypeRaw(r)
	case slices.Contains(VariantDataTypeSynonyms, name) && !hasArgs:
		return parseVariantDataTypeRaw(r)
	case slices.Contains(ObjectDataTypeSynonyms, name):
		return parseObjectDataTypeRaw(sanitizedDataTypeRaw{strings.TrimSpace(raw), name})
	case slices.Contains(ArrayDataTypeSynonyms, name):
		return parseArrayDataTypeRaw(sanitizedDataTypeRaw{strings.TrimSpace(raw), name})
	case slices.Contains(MapDataTypeSynonyms, name):
		return parseMapDataTypeRaw(sanitizedDataTypeRaw{strings.TrimSpace(raw), name})
	case slices.Contains(GeographyDataTypeSynonyms, name) && !hasArgs:
		return parseGeographyDataTypeRaw(r)
	case slices.Contains(GeometryDataTypeSynonyms, name) && !hasArgs:
		return parseGeometryDataTypeRaw(r)
	case slices.Contains(VectorDataTypeSynonyms, name):
		return parseVectorDataTypeRaw(r)
	case slices.Contains(TableDataTypeSynonyms, name):
		return parseTableDataTypeRaw(sanitizedDataTypeRaw{strings.TrimSpace(raw), name})
	default:
		return nil, fmt.Errorf("invalid data type: %s", raw)
	}
}

// AreTheSame compares any two data types.
// If both data types are nil it returns true.
// If only one data type is nil it returns false.
// It returns false for different underlying types.
// For the same type it performs type-specific comparison.
func AreTheSame(a DataType, b DataType) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	if _, ok := a.(plain); ok {
		return true
	}
	if v, ok := a.(sized); ok {
		return areSizedDataTypesTheSame(v.asSized(), b.(sized).asSized())
	}
	switch v := a.(type) {
	case *ArrayDataType:
		return castSuccessfully(v, b, areArrayDataTypesTheSame)
	case *MapDataType:
		return castSuccessfully(v, b, areMapDataTypesTheSame)
	case *NumberDataType:
		return castSuccessfully(v, b, areNumberDataTypesTheSame)
	case *ObjectDataType:
		return castSuccessfully(v, b, areObjectDataTypesTheSame)
	case *TableDataType:
		return castSuccessfully(v, b, areTableDataTypesTheSame)
	case *VectorDataType:
		return castSuccessfully(v, b, areVectorDataTypesTheSame)
	}
	return false
}

func IsTextDataType(a DataType) bool {
	_, ok := a.(*TextDataType)
	return ok
}

func castSuccessfully[T any](a T, b DataType, invoke func(a T, b T) bool) bool {
	if dCasted, ok := b.(T); ok {
		return invoke(a, dCasted)
	}
	return false
}

// AreDefinitelyDifferent compares any two data types.
// The logic for the equality check has 3 values: YES, NO, and MAYBE.
// MAYBE is a result of Snowflake not always returning the attributes of the complex data type,
// e.g. NUMBER(20, 4) can be returned as NUMBER and NUMBER(38, 0) (the defaults) can also be returned as NUMBER.
// Because of that, they are MAYBE equal.
// This method is meant to be an entry point in cases where we need to be sure, that the data types are different (the NO part).
// Example could be the handling of the external Snowflake changes where the data type is for sure different.
//
// The logic goes like this:
// - If both data types are nil it returns false.
// - If only one data type is nil it returns true.
// - It returns true for different underlying types.
// - For the same type it performs type-specific check.
func AreDefinitelyDifferent(a DataType, b DataType) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil || b == nil {
		return true
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return true
	}
	if _, ok := a.(plain); ok {
		return false
	}
	if v, ok := a.(sized); ok {
		return areSizedDataTypesDefinitelyDifferent(v.asSized(), b.(sized).asSized())
	}
	switch v := a.(type) {
	case *ArrayDataType:
		return castSuccessfully(v, b, areArrayDataTypesDefinitelyDifferent)
	case *MapDataType:
		return castSuccessfully(v, b, areMapDataTypesDefinitelyDifferent)
	case *NumberDataType:
		return castSuccessfully(v, b, areNumberDataTypesDefinitelyDifferent)
	case *ObjectDataType:
		return castSuccessfully(v, b, areObjectDataTypesDefinitelyDifferent)
	case *TableDataType:
		return castSuccessfully(v, b, areTableDataTypesDefinitelyDifferent)
	case *VectorDataType:
		return castSuccessfully(v, b, areVectorDataTypesDefinitelyDifferent)
	}
	return false
}
