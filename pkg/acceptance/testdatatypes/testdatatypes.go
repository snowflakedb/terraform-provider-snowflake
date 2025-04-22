package testdatatypes

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"

// TODO [SNOW-1843440]: create using constructors (when we add them)?
var (
	DataTypeNumber_36_2, _ = datatypes.ParseDataType("NUMBER(36, 2)")
	DataTypeVarchar_100, _ = datatypes.ParseDataType("VARCHAR(100)")
	DataTypeVarchar, _     = datatypes.ParseDataType("VARCHAR")
	DataTypeNumber, _      = datatypes.ParseDataType("NUMBER")
	DataTypeFloat, _       = datatypes.ParseDataType("FLOAT")
	DataTypeVariant, _     = datatypes.ParseDataType("VARIANT")
)
