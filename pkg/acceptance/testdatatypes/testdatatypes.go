package testdatatypes

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

var (
	DataTypeNumber_36_2      = datatypes.NewNumberDataTypeWithPrecisionAndScale(36, 2)
	DataTypeNumber_2_0       = datatypes.NewNumberDataTypeWithPrecisionAndScale(2, 0)
	DataTypeVarchar_50       = datatypes.NewVarcharDataTypeWithLength(50)
	DataTypeVarchar_100      = datatypes.NewVarcharDataTypeWithLength(100)
	DataTypeVarchar_200      = datatypes.NewVarcharDataTypeWithLength(200)
	DataTypeVarchar_255      = datatypes.NewVarcharDataTypeWithLength(255)
	DataTypeVarchar_256      = datatypes.NewVarcharDataTypeWithLength(256)
	DataTypeVarchar_500      = datatypes.NewVarcharDataTypeWithLength(500)
	DataTypeVarchar          = datatypes.NewVarcharDataType()
	DataTypeDecfloat         = datatypes.NewDecfloatDataType()
	DataTypeText, _          = datatypes.ParseDataType("TEXT")
	DataTypeChar             = datatypes.NewCharDataType()
	DataTypeString, _        = datatypes.ParseDataType("STRING")
	DataTypeBoolean          = datatypes.NewBooleanDataType()
	DataTypeNumber           = datatypes.NewNumberDataType()
	DataTypeDecimal_38_0, _  = datatypes.ParseDataType("DECIMAL(38, 0)")
	DataTypeNumber_38_0      = datatypes.NewNumberDataTypeWithPrecisionAndScale(38, 0)
	DataTypeNumber_36_0      = datatypes.NewNumberDataTypeWithPrecisionAndScale(36, 0)
	DataTypeNumber_19_0      = datatypes.NewNumberDataTypeWithPrecisionAndScale(19, 0)
	DataTypeInteger, _       = datatypes.ParseDataType("INTEGER")
	DataTypeDecimal, _       = datatypes.ParseDataType("DECIMAL")
	DataTypeFloat            = datatypes.NewFloatDataType()
	DataTypeDouble, _        = datatypes.ParseDataType("DOUBLE")
	DataTypeBinary           = datatypes.NewBinaryDataType()
	DataTypeVarbinary, _     = datatypes.ParseDataType("VARBINARY")
	DataTypeVariant          = datatypes.NewVariantDataType()
	DataTypeObject           = datatypes.NewObjectDataType()
	DataTypeStructuredObject = datatypes.NewStructuredObjectDataType(
		datatypes.ObjectDataTypeField{Name: "foo", DataType: datatypes.NewVarcharDataTypeWithLength(datatypes.MaxVarcharLength)},
	)
	DataTypeArray           = datatypes.NewArrayDataType()
	DataTypeGeography       = datatypes.NewGeographyDataType()
	DataTypeGeometry        = datatypes.NewGeometryDataType()
	DataTypeTime            = datatypes.NewTimeDataType()
	DataTypeDate            = datatypes.NewDateDataType()
	DataTypeDatetime, _     = datatypes.ParseDataType("DATETIME")
	DataTypeTimestampNTZ    = datatypes.NewTimestampNtzDataType()
	DataTypeTimestampLTZ    = datatypes.NewTimestampLtzDataType()
	DataTypeTimestampTZ     = datatypes.NewTimestampTzDataType()
	DataTypeVarcharIceberg  = datatypes.NewVarcharDataTypeWithLength(datatypes.MaxVarcharLength)
	DataTypeTimestampNTZ_6  = datatypes.NewTimestampNtzDataTypeWithPrecision(6)
	DataTypeTimestampNTZ_9  = datatypes.NewTimestampNtzDataTypeWithPrecision(9)
	DataTypeVectorFloat_768 = datatypes.NewVectorFloatDataType(768)
)

var DefaultVarcharAsString = fmt.Sprintf("VARCHAR(%d)", datatypes.DefaultVarcharLength)
