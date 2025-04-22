package datatypes

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// TODO [SNOW-1843440]: generalize definitions for different types; generalize the ParseDataType function
// TODO [SNOW-1843440]: generalize implementation in types (i.e. the internal struct implementing ToLegacyDataTypeSql and containing the underlyingType)
// TODO [SNOW-1843440]: consider known/unknown to use Snowflake defaults and allow better handling in terraform resources
// TODO [SNOW-1843440]: replace old DataTypes

// DataType is the common interface that represents all Snowflake datatypes documented in https://docs.snowflake.com/en/sql-reference/intro-summary-data-types.
type DataType interface {
	// ToSql formats data type explicitly specifying all arguments and using the given type (e.g. CHAR(29) for CHAR(29)).
	ToSql() string
	// ToLegacyDataTypeSql formats data type using its base type without any attributes (e.g. VARCHAR for CHAR(29)).
	ToLegacyDataTypeSql() string
	// Canonical formats the data type between ToSql and ToLegacyDataTypeSql: it uses base type but with arguments (e.g. VARCHAR(29) for CHAR(29)).
	Canonical() string
}

type sanitizedDataTypeRaw struct {
	raw           string
	matchedByType string
}

// ParseDataType is the entry point to get the implementation of the DataType from input raw string.
// TODO [SNOW-1843440]: order currently matters (e.g. HasPrefix(TIME) can match also TIMESTAMP*, make the checks more precise and order-independent)
func ParseDataType(raw string) (DataType, error) {
	dataTypeRaw := strings.TrimSpace(strings.ToUpper(raw))

	if idx := slices.IndexFunc(AllNumberDataTypes, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseNumberDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, AllNumberDataTypes[idx]})
	}
	if slices.Contains(FloatDataTypeSynonyms, dataTypeRaw) {
		return parseFloatDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, dataTypeRaw})
	}
	if idx := slices.IndexFunc(AllTextDataTypes, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTextDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, AllTextDataTypes[idx]})
	}
	if idx := slices.IndexFunc(BinaryDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseBinaryDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, BinaryDataTypeSynonyms[idx]})
	}
	if slices.Contains(BooleanDataTypeSynonyms, dataTypeRaw) {
		return parseBooleanDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, dataTypeRaw})
	}
	if slices.Contains(DateDataTypeSynonyms, dataTypeRaw) {
		return parseDateDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, dataTypeRaw})
	}
	if idx := slices.IndexFunc(TimestampLtzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampLtzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampLtzDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimestampNtzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampNtzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampNtzDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimestampTzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampTzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampTzDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimeDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimeDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimeDataTypeSynonyms[idx]})
	}
	if slices.Contains(VariantDataTypeSynonyms, dataTypeRaw) {
		return parseVariantDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, dataTypeRaw})
	}
	if slices.Contains(ObjectDataTypeSynonyms, dataTypeRaw) {
		return parseObjectDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, dataTypeRaw})
	}
	if slices.Contains(ArrayDataTypeSynonyms, dataTypeRaw) {
		return parseArrayDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, dataTypeRaw})
	}
	if slices.Contains(GeographyDataTypeSynonyms, dataTypeRaw) {
		return parseGeographyDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, dataTypeRaw})
	}
	if slices.Contains(GeometryDataTypeSynonyms, dataTypeRaw) {
		return parseGeometryDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, dataTypeRaw})
	}
	if idx := slices.IndexFunc(VectorDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseVectorDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, VectorDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TableDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTableDataTypeRaw(sanitizedDataTypeRaw{strings.TrimSpace(raw), TableDataTypeSynonyms[idx]})
	}

	return nil, fmt.Errorf("invalid data type: %s", raw)
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
	switch v := a.(type) {
	case *ArrayDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *BinaryDataType:
		return castSuccessfully(v, b, areBinaryDataTypesTheSame)
	case *BooleanDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *DateDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *FloatDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *GeographyDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *GeometryDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *NumberDataType:
		return castSuccessfully(v, b, areNumberDataTypesTheSame)
	case *ObjectDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *TableDataType:
		return castSuccessfully(v, b, areTableDataTypesTheSame)
	case *TextDataType:
		return castSuccessfully(v, b, areTextDataTypesTheSame)
	case *TimeDataType:
		return castSuccessfully(v, b, areTimeDataTypesTheSame)
	case *TimestampLtzDataType:
		return castSuccessfully(v, b, areTimestampLtzDataTypesTheSame)
	case *TimestampNtzDataType:
		return castSuccessfully(v, b, areTimestampNtzDataTypesTheSame)
	case *TimestampTzDataType:
		return castSuccessfully(v, b, areTimestampTzDataTypesTheSame)
	case *VariantDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
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

func noArgsDataTypesAreTheSame[T DataType](_ T, _ T) bool {
	return true
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
// TODO [next PR]: test this function if approved
// TODO [next PR]: implement all missing checks if approved
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
	switch v := a.(type) {
	case *ArrayDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreDefinitelyDifferent)
	case *BinaryDataType:
		return castSuccessfully(v, b, willBeImplemented)
	case *BooleanDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreDefinitelyDifferent)
	case *DateDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreDefinitelyDifferent)
	case *FloatDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreDefinitelyDifferent)
	case *GeographyDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreDefinitelyDifferent)
	case *GeometryDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreDefinitelyDifferent)
	case *NumberDataType:
		return castSuccessfully(v, b, areNumberDataTypesDefinitelyDifferent)
	case *ObjectDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreDefinitelyDifferent)
	case *TableDataType:
		return castSuccessfully(v, b, willBeImplemented)
	case *TextDataType:
		return castSuccessfully(v, b, areTextDataTypesDefinitelyDifferent)
	case *TimeDataType:
		return castSuccessfully(v, b, willBeImplemented)
	case *TimestampLtzDataType:
		return castSuccessfully(v, b, willBeImplemented)
	case *TimestampNtzDataType:
		return castSuccessfully(v, b, willBeImplemented)
	case *TimestampTzDataType:
		return castSuccessfully(v, b, willBeImplemented)
	case *VariantDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreDefinitelyDifferent)
	case *VectorDataType:
		return castSuccessfully(v, b, willBeImplemented)
	}
	return false
}

func noArgsDataTypesAreDefinitelyDifferent[T DataType](_ T, _ T) bool {
	return false
}

func willBeImplemented[T DataType](_ T, _ T) bool {
	return false
}
