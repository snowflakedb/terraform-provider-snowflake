package testdatatypes

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

// TODO [SNOW-1843440]: create using constructors (when we add them)?
var (
	DataTypeNumber_36_2, _  = datatypes.ParseDataType("NUMBER(36, 2)")
	DataTypeNumber_2_0, _   = datatypes.ParseDataType("NUMBER(2, 0)")
	DataTypeVarchar_100, _  = datatypes.ParseDataType("VARCHAR(100)")
	DataTypeVarchar_200, _  = datatypes.ParseDataType("VARCHAR(200)")
	DataTypeVarchar, _      = datatypes.ParseDataType("VARCHAR")
	DataTypeText, _         = datatypes.ParseDataType("TEXT")
	DataTypeChar, _         = datatypes.ParseDataType("CHAR")
	DataTypeString, _       = datatypes.ParseDataType("STRING")
	DataTypeBoolean, _      = datatypes.ParseDataType("BOOLEAN")
	DataTypeNumber, _       = datatypes.ParseDataType("NUMBER")
	DataTypeInteger, _      = datatypes.ParseDataType("INTEGER")
	DataTypeDecimal, _      = datatypes.ParseDataType("DECIMAL")
	DataTypeFloat, _        = datatypes.ParseDataType("FLOAT")
	DataTypeDouble, _       = datatypes.ParseDataType("DOUBLE")
	DataTypeBinary, _       = datatypes.ParseDataType("BINARY")
	DataTypeVarbinary, _    = datatypes.ParseDataType("VARBINARY")
	DataTypeVariant, _      = datatypes.ParseDataType("VARIANT")
	DataTypeObject, _       = datatypes.ParseDataType("OBJECT")
	DataTypeArray, _        = datatypes.ParseDataType("ARRAY")
	DataTypeGeography, _    = datatypes.ParseDataType("GEOGRAPHY")
	DataTypeGeometry, _     = datatypes.ParseDataType("GEOMETRY")
	DataTypeTime, _         = datatypes.ParseDataType("TIME")
	DataTypeDate, _         = datatypes.ParseDataType("DATE")
	DataTypeDatetime, _     = datatypes.ParseDataType("DATETIME")
	DataTypeTimestampNTZ, _ = datatypes.ParseDataType("TIMESTAMP_NTZ")
	DataTypeTimestampLTZ, _ = datatypes.ParseDataType("TIMESTAMP_LTZ")
	DataTypeTimestampTZ, _  = datatypes.ParseDataType("TIMESTAMP_TZ")
)

var DefaultVarcharAsString = fmt.Sprintf("VARCHAR(%d)", datatypes.DefaultVarcharLength)
