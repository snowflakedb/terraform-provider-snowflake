package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

var (
	_    = randomTableColumnIdentifier()
	_, _ = emptyAccountIdentifier, emptyTableColumnIdentifier
)

var (
	invalidAccountObjectIdentifier = NewAccountObjectIdentifier(random.StringN(256))
	longSchemaObjectIdentifier     = NewSchemaObjectIdentifier(random.StringN(255), random.StringN(255), random.StringN(255))

	// TODO: Add to the generator
	emptyAccountObjectIdentifier             = NewAccountObjectIdentifier("")
	emptyExternalObjectIdentifier            = NewExternalObjectIdentifier(NewAccountIdentifier("", ""), NewObjectIdentifierFromFullyQualifiedName(""))
	emptyDatabaseObjectIdentifier            = NewDatabaseObjectIdentifier("", "")
	emptySchemaObjectIdentifier              = NewSchemaObjectIdentifier("", "", "")
	emptySchemaObjectIdentifierWithArguments = NewSchemaObjectIdentifierWithArguments("", "", "")
	emptyAccountIdentifier                   = NewAccountIdentifier("", "")
	emptyTableColumnIdentifier               = NewTableColumnIdentifier("", "", "", "")

	dataTypeNumber_36_2                   = datatypes.NewNumberDataTypeWithPrecisionAndScale(36, 2)
	dataTypeNumber                        = datatypes.NewNumberDataType()
	dataTypeNumber_38_0                   = datatypes.NewNumberDataTypeWithPrecisionAndScale(38, 0)
	dataTypeVarchar_100                   = datatypes.NewVarcharDataTypeWithLength(100)
	dataTypeVarchar_200                   = datatypes.NewVarcharDataTypeWithLength(200)
	dataTypeVarchar                       = datatypes.NewVarcharDataType()
	dataTypeBoolean                       = datatypes.NewBooleanDataType()
	dataTypeFloat                         = datatypes.NewFloatDataType()
	dataTypeVariant                       = datatypes.NewVariantDataType()
	dataTypeChar                          = datatypes.NewCharDataType()
	dataTypeChar_100                      = datatypes.NewCharDataTypeWithLength(100)
	dataTypeDoublePrecision, _            = datatypes.ParseDataType("DOUBLE PRECISION")
	dataTypeTimestampWithoutTimeZone_5, _ = datatypes.ParseDataType("TIMESTAMP WITHOUT TIME ZONE(5)")
	dataTypeTime                          = datatypes.NewTimeDataType()
	dataTypeTimestampTz                   = datatypes.NewTimestampTzDataType()
	dataTypeVectorInt_10                  = datatypes.NewVectorIntDataType(10)
	dataTypeVectorInt_20                  = datatypes.NewVectorIntDataType(20)
	dataTypeVectorFloat_10                = datatypes.NewVectorFloatDataType(10)
	dataTypeTableDate                     = datatypes.NewTableDataType(datatypes.NewDateDataType())
	dataTypeObjectStructured              = datatypes.NewStructuredObjectDataType(
		datatypes.ObjectDataTypeField{Name: "a", DataType: datatypes.NewVarcharDataType()},
		datatypes.ObjectDataTypeField{Name: "b", DataType: datatypes.NewNumberDataType(), NotNull: true},
	)
	dataTypeArrayNumber      = datatypes.NewStructuredArrayDataType(datatypes.NewNumberDataType())
	dataTypeMapVarcharNumber = datatypes.NewMapDataType(datatypes.NewVarcharDataType(), datatypes.NewNumberDataType())
)

func randomSchemaObjectIdentifierWithArguments(argumentDataTypes ...datatypes.DataType) SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(random.StringN(12), random.StringN(12), random.StringN(12), argumentDataTypes...)
}

func randomSchemaObjectIdentifier() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(random.StringN(12), random.StringN(12), random.StringN(12))
}

func randomSchemaObjectIdentifierInSchema(schemaId DatabaseObjectIdentifier) SchemaObjectIdentifier {
	return NewSchemaObjectIdentifierInSchema(schemaId, random.StringN(12))
}

func randomExternalObjectIdentifier() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator(random.StringN(12)), randomAccountObjectIdentifier())
}

func randomDatabaseObjectIdentifier() DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(random.StringN(12), random.StringN(12))
}

func randomDatabaseObjectIdentifierInDatabase(databaseId AccountObjectIdentifier) DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(databaseId.Name(), random.StringN(12))
}

func randomAccountIdentifier() AccountIdentifier {
	return NewAccountIdentifier(random.StringN(12), random.StringN(12))
}

func randomAccountObjectIdentifier() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(random.StringN(12))
}

func randomTableColumnIdentifier() TableColumnIdentifier {
	return NewTableColumnIdentifier(random.StringN(12), random.StringN(12), random.StringN(12), random.StringN(12))
}

func randomTableColumnIdentifierInSchemaObject(objectId SchemaObjectIdentifier) TableColumnIdentifier {
	return NewTableColumnIdentifier(objectId.DatabaseName(), objectId.SchemaName(), objectId.Name(), random.StringN(12))
}
