package helpers

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

type IdsGenerator struct {
	context *TestClientContext
}

func NewIdsGenerator(context *TestClientContext) *IdsGenerator {
	return &IdsGenerator{
		context: context,
	}
}

func (c *IdsGenerator) DatabaseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.context.database)
}

func (c *IdsGenerator) SchemaId() sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(c.context.database, c.context.schema)
}

func (c *IdsGenerator) WarehouseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.context.warehouse)
}

func (c *IdsGenerator) SnowflakeWarehouseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier("SNOWFLAKE")
}

func (c *IdsGenerator) AccountIdentifierWithLocator() sdk.AccountIdentifier {
	return sdk.NewAccountIdentifierFromAccountLocator(c.context.client.GetAccountLocator())
}

func (c *IdsGenerator) RandomAccountObjectIdentifier() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.Alpha())
}

func (c *IdsGenerator) RandomSensitiveAccountObjectIdentifier() sdk.AccountObjectIdentifier {
	return c.RandomAccountObjectIdentifierWithPrefix(random.SensitiveAlphanumeric())
}

func (c *IdsGenerator) RandomAccountObjectIdentifierWithPrefix(prefix string) sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.AlphaWithPrefix(prefix))
}

func (c *IdsGenerator) RandomAccountObjectIdentifierContaining(part string) sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.AlphaContaining(part))
}

func (c *IdsGenerator) NewDatabaseObjectIdentifier(name string) sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(c.DatabaseId().Name(), name)
}

func (c *IdsGenerator) NewDatabaseObjectIdentifierInDatabase(name string, databaseId sdk.AccountObjectIdentifier) sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifierInDatabase(databaseId, name)
}

func (c *IdsGenerator) RandomDatabaseObjectIdentifier() sdk.DatabaseObjectIdentifier {
	return c.RandomDatabaseObjectIdentifierInDatabase(c.DatabaseId())
}

func (c *IdsGenerator) RandomDatabaseObjectIdentifierInDatabase(databaseId sdk.AccountObjectIdentifier) sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(databaseId.Name(), c.Alpha())
}

func (c *IdsGenerator) RandomDatabaseObjectIdentifierWithPrefix(prefix string) sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(c.DatabaseId().Name(), c.AlphaWithPrefix(prefix))
}

func (c *IdsGenerator) NewSchemaObjectIdentifier(name string) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierInSchema(c.SchemaId(), name)
}

func (c *IdsGenerator) NewSchemaObjectIdentifierInSchema(name string, schemaId sdk.DatabaseObjectIdentifier) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierInSchema(schemaId, name)
}

func (c *IdsGenerator) RandomSchemaObjectIdentifier() sdk.SchemaObjectIdentifier {
	return c.RandomSchemaObjectIdentifierInSchema(c.SchemaId())
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierContaining(part string) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierInSchema(c.SchemaId(), c.AlphaContaining(part))
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierWithPrefix(prefix string) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierInSchema(c.SchemaId(), c.AlphaWithPrefix(prefix))
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierInSchema(schemaId sdk.DatabaseObjectIdentifier) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierInSchema(schemaId, c.Alpha())
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix string, schemaId sdk.DatabaseObjectIdentifier) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierInSchema(schemaId, c.AlphaWithPrefix(prefix))
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierWithArgumentsOld(arguments ...sdk.DataType) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierWithArgumentsOld(c.SchemaId().DatabaseName(), c.SchemaId().Name(), c.Alpha(), arguments)
}

func (c *IdsGenerator) NewSchemaObjectIdentifierWithArguments(name string, arguments ...sdk.DataType) sdk.SchemaObjectIdentifierWithArguments {
	return sdk.NewSchemaObjectIdentifierWithArguments(c.SchemaId().DatabaseName(), c.SchemaId().Name(), name, arguments...)
}

func (c *IdsGenerator) NewSchemaObjectIdentifierWithArgumentsNewDataTypes(name string, arguments ...datatypes.DataType) sdk.SchemaObjectIdentifierWithArguments {
	legacyDataTypes := collections.Map(arguments, sdk.LegacyDataTypeFrom)
	return sdk.NewSchemaObjectIdentifierWithArguments(c.SchemaId().DatabaseName(), c.SchemaId().Name(), name, legacyDataTypes...)
}

func (c *IdsGenerator) NewSchemaObjectIdentifierWithArgumentsInSchema(name string, schemaId sdk.DatabaseObjectIdentifier, argumentDataTypes ...sdk.DataType) sdk.SchemaObjectIdentifierWithArguments {
	return sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(schemaId, name, argumentDataTypes...)
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierWithArguments(arguments ...sdk.DataType) sdk.SchemaObjectIdentifierWithArguments {
	return sdk.NewSchemaObjectIdentifierWithArguments(c.SchemaId().DatabaseName(), c.SchemaId().Name(), c.Alpha(), arguments...)
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierWithArgumentsInSchema(schemaId sdk.DatabaseObjectIdentifier, arguments ...sdk.DataType) sdk.SchemaObjectIdentifierWithArguments {
	return sdk.NewSchemaObjectIdentifierWithArguments(schemaId.DatabaseName(), schemaId.Name(), c.Alpha(), arguments...)
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(arguments ...datatypes.DataType) sdk.SchemaObjectIdentifierWithArguments {
	legacyDataTypes := collections.Map(arguments, sdk.LegacyDataTypeFrom)
	return sdk.NewSchemaObjectIdentifierWithArguments(c.SchemaId().DatabaseName(), c.SchemaId().Name(), c.Alpha(), legacyDataTypes...)
}

func (c *IdsGenerator) Alpha() string {
	return c.AlphaN(6)
}

func (c *IdsGenerator) AlphaN(n int) string {
	return c.WithTestObjectSuffix(strings.ToUpper(random.AlphaN(n)))
}

func (c *IdsGenerator) AlphaContaining(part string) string {
	return c.WithTestObjectSuffix(c.Alpha() + part)
}

func (c *IdsGenerator) AlphaWithPrefix(prefix string) string {
	return prefix + c.Alpha()
}

func (c *IdsGenerator) WithTestObjectSuffix(text string) string {
	return text + c.context.testObjectSuffix
}
