package testdatatypes

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

// TODO [SNOW-1843440]: create using constructors (when we add them)?
var (
	DataTypeNumber_36_2, _    = datatypes.ParseDataType("NUMBER(36, 2)")
	DataTypeNumber_2_0, _     = datatypes.ParseDataType("NUMBER(2, 0)")
	DataTypeVarchar_100, _    = datatypes.ParseDataType("VARCHAR(100)")
	DataTypeVarchar_200, _    = datatypes.ParseDataType("VARCHAR(200)")
	DataTypeVarchar, _        = datatypes.ParseDataType("VARCHAR")
	DataTypeDecfloat, _       = datatypes.ParseDataType("DECFLOAT")
	DataTypeText, _           = datatypes.ParseDataType("TEXT")
	DataTypeChar, _           = datatypes.ParseDataType("CHAR")
	DataTypeString, _         = datatypes.ParseDataType("STRING")
	DataTypeBoolean, _        = datatypes.ParseDataType("BOOLEAN")
	DataTypeNumber, _         = datatypes.ParseDataType("NUMBER")
	DataTypeDecimal_38_0, _   = datatypes.ParseDataType("DECIMAL(38, 0)")
	DataTypeNumber_38_0, _    = datatypes.ParseDataType("NUMBER(38, 0)")
	DataTypeInteger, _        = datatypes.ParseDataType("INTEGER")
	DataTypeDecimal, _        = datatypes.ParseDataType("DECIMAL")
	DataTypeFloat, _          = datatypes.ParseDataType("FLOAT")
	DataTypeDouble, _         = datatypes.ParseDataType("DOUBLE")
	DataTypeBinary, _         = datatypes.ParseDataType("BINARY")
	DataTypeVarbinary, _      = datatypes.ParseDataType("VARBINARY")
	DataTypeVariant, _        = datatypes.ParseDataType("VARIANT")
	DataTypeObject, _         = datatypes.ParseDataType("OBJECT")
	DataTypeArray, _          = datatypes.ParseDataType("ARRAY")
	DataTypeGeography, _      = datatypes.ParseDataType("GEOGRAPHY")
	DataTypeGeometry, _       = datatypes.ParseDataType("GEOMETRY")
	DataTypeTime, _           = datatypes.ParseDataType("TIME")
	DataTypeDate, _           = datatypes.ParseDataType("DATE")
	DataTypeDatetime, _       = datatypes.ParseDataType("DATETIME")
	DataTypeTimestampNTZ, _   = datatypes.ParseDataType("TIMESTAMP_NTZ")
	DataTypeTimestampLTZ, _   = datatypes.ParseDataType("TIMESTAMP_LTZ")
	DataTypeTimestampTZ, _    = datatypes.ParseDataType("TIMESTAMP_TZ")
	DataTypeVarcharIceberg, _ = datatypes.ParseDataType("VARCHAR(134217728)")
	DataTypeTimestampNTZ_6, _ = datatypes.ParseDataType("TIMESTAMP_NTZ(6)")
)

var DefaultVarcharAsString = fmt.Sprintf("VARCHAR(%d)", datatypes.DefaultVarcharLength)
