package sdk

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

// DataType is based on https://docs.snowflake.com/en/sql-reference/intro-summary-data-types.
type DataType string

var allowedVectorInnerTypes = []DataType{
	DataTypeInt,
	DataTypeFloat,
}

const (
	DataTypeNumber       DataType = "NUMBER"
	DataTypeInt          DataType = "INT"
	DataTypeFloat        DataType = "FLOAT"
	DataTypeVARCHAR      DataType = "VARCHAR"
	DataTypeString       DataType = "STRING"
	DataTypeBinary       DataType = "BINARY"
	DataTypeBoolean      DataType = "BOOLEAN"
	DataTypeDate         DataType = "DATE"
	DataTypeTime         DataType = "TIME"
	DataTypeTimestampLTZ DataType = "TIMESTAMP_LTZ"
	DataTypeTimestampNTZ DataType = "TIMESTAMP_NTZ"
	DataTypeTimestampTZ  DataType = "TIMESTAMP_TZ"
	DataTypeVariant      DataType = "VARIANT"
	DataTypeObject       DataType = "OBJECT"
	DataTypeArray        DataType = "ARRAY"
	DataTypeGeography    DataType = "GEOGRAPHY"
	DataTypeGeometry     DataType = "GEOMETRY"
)

// IsStringType is a legacy method. datatypes.IsTextDataType should be used instead.
// TODO [SNOW-1348114]: remove with tables rework
func IsStringType(_type string) bool {
	t := strings.ToUpper(_type)
	return strings.HasPrefix(t, "STRING") ||
		strings.HasPrefix(t, "VARCHAR") ||
		strings.HasPrefix(t, "CHAR") ||
		strings.HasPrefix(t, "TEXT") ||
		strings.HasPrefix(t, "NVARCHAR") ||
		strings.HasPrefix(t, "NCHAR")
}

func LegacyDataTypeFrom(newDataType datatypes.DataType) DataType {
	// TODO [SNOW-1850370]: remove this check?
	if newDataType == nil {
		return ""
	}
	return DataType(newDataType.ToLegacyDataTypeSql())
}

func LegacyDataTypeWithAttrs(newDataType datatypes.DataType) DataType {
	// TODO [SNOW-1850370]: remove this check?
	if newDataType == nil {
		return ""
	}
	return DataType(newDataType.ToSql())
}

func LegacyDataTypeWithAttrsCanonical(newDataType datatypes.DataType) DataType {
	// TODO [SNOW-1850370]: remove this check?
	if newDataType == nil {
		return ""
	}
	return DataType(newDataType.Canonical())
}
