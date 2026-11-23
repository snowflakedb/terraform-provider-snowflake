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

	// TODO [SNOW-1843440]: create using constructors (when we add them)?
	dataTypeNumber_36_2, _                = datatypes.ParseDataType("NUMBER(36, 2)")
	dataTypeNumber, _                     = datatypes.ParseDataType("NUMBER")
	dataTypeNumber_38_0, _                = datatypes.ParseDataType("NUMBER(38, 0)")
	dataTypeVarchar_100, _                = datatypes.ParseDataType("VARCHAR(100)")
	dataTypeVarchar_200, _                = datatypes.ParseDataType("VARCHAR(200)")
	dataTypeVarchar, _                    = datatypes.ParseDataType("VARCHAR")
	dataTypeBoolean, _                    = datatypes.ParseDataType("BOOLEAN")
	dataTypeFloat, _                      = datatypes.ParseDataType("FLOAT")
	dataTypeVariant, _                    = datatypes.ParseDataType("VARIANT")
	dataTypeChar, _                       = datatypes.ParseDataType("CHAR")
	dataTypeChar_100, _                   = datatypes.ParseDataType("CHAR(100)")
	dataTypeDoublePrecision, _            = datatypes.ParseDataType("DOUBLE PRECISION")
	dataTypeTimestampWithoutTimeZone_5, _ = datatypes.ParseDataType("TIMESTAMP WITHOUT TIME ZONE(5)")
	dataTypeTime, _                       = datatypes.ParseDataType("TIME")
	dataTypeTimestampTz, _                = datatypes.ParseDataType("TIMESTAMP_TZ")
	dataTypeVectorInt_10, _               = datatypes.ParseDataType("VECTOR(INT, 10)")
	dataTypeVectorInt_20, _               = datatypes.ParseDataType("VECTOR(INT, 20)")
	dataTypeVectorFloat_10, _             = datatypes.ParseDataType("VECTOR(FLOAT, 10)")
	dataTypeTableDate, _                  = datatypes.ParseDataType("TABLE(DATE)")
	dataTypeObjectStructured, _           = datatypes.ParseDataType("OBJECT(a VARCHAR, b NUMBER NOT NULL)")
	dataTypeArrayNumber, _                = datatypes.ParseDataType("ARRAY(NUMBER)")
	dataTypeMapVarcharNumber, _           = datatypes.ParseDataType("MAP(VARCHAR, NUMBER)")
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
