package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
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

	// TODO [SNOW-1843440]: create using constructors (when we add them)?
	dataTypeNumber_36_2, _                = datatypes.ParseDataType("NUMBER(36, 2)")
	dataTypeNumber, _                     = datatypes.ParseDataType("NUMBER")
	dataTypeVarchar_100, _                = datatypes.ParseDataType("VARCHAR(100)")
	dataTypeVarchar, _                    = datatypes.ParseDataType("VARCHAR")
	dataTypeBoolean, _                    = datatypes.ParseDataType("BOOLEAN")
	dataTypeFloat, _                      = datatypes.ParseDataType("FLOAT")
	dataTypeVariant, _                    = datatypes.ParseDataType("VARIANT")
	dataTypeChar, _                       = datatypes.ParseDataType("CHAR")
	dataTypeChar_100, _                   = datatypes.ParseDataType("CHAR(100)")
	dataTypeDoublePrecision, _            = datatypes.ParseDataType("DOUBLE PRECISION")
	dataTypeTimestampWithoutTimeZone_5, _ = datatypes.ParseDataType("TIMESTAMP WITHOUT TIME ZONE(5)")
)

func randomSchemaObjectIdentifierWithArguments(argumentDataTypes ...DataType) SchemaObjectIdentifierWithArguments {
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
