package defs

import (
	"slices"

	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

var schemaParameters = ParameterDefsForLevel(parameterdefs.ParameterLevelSchema)

var schemaParameterFieldNames = collections.Map(schemaParameters, g.ParameterSqlToFieldName)

var schemaSetStruct = g.NewQueryStruct("SchemaSet").
	WithParameters(schemaParameters...).
	OptionalComment().
	WithValidation(g.ValidIdentifierIfSet, "ExternalVolume").
	WithValidation(g.ValidIdentifierIfSet, "Catalog").
	WithValidation(g.AtLeastOneValueSet, append(slices.Clone(schemaParameterFieldNames), "Comment")...)

var schemaUnsetStruct = g.NewQueryStruct("SchemaUnset").
	WithParametersUnset(schemaParameters...).
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, append(slices.Clone(schemaParameterFieldNames), "Comment")...)

var schemasDef = g.NewInterface(
	"Schemas",
	"Schema",
	g.KindOfT[sdkcommons.DatabaseObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-schema",
	g.NewQueryStruct("CreateSchema").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("SCHEMA").
		IfNotExists().
		Name().
		OptionalSQL("WITH MANAGED ACCESS").
		WithParameters(schemaParameters...).
		OptionalComment().
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ValidIdentifierIfSet, "ExternalVolume").
		WithValidation(g.ValidIdentifierIfSet, "Catalog"),
).CustomOperation(
	"Clone",
	"https://docs.snowflake.com/en/sql-reference/sql/create-schema",
	g.NewQueryStruct("CloneSchema").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("SCHEMA").
		IfNotExists().
		Name().
		PredefinedQueryStructField("Clone", "Clone", g.KeywordOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-schema",
	g.NewQueryStruct("AlterSchema").
		Alter().
		SQL("SCHEMA").
		IfExists().
		Name().
		RenameTo().
		Identifier("SwapWith", g.KindOfTPointer[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("SWAP WITH")).
		OptionalQueryStructField("Set", schemaSetStruct, g.ListOptions().NoParentheses().SQL("SET")).
		OptionalQueryStructField("Unset", schemaUnsetStruct, g.ListOptions().NoParentheses().SQL("UNSET")).
		OptionalSetTags().
		OptionalUnsetTags().
		OptionalSQL("ENABLE MANAGED ACCESS").
		OptionalSQL("DISABLE MANAGED ACCESS").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "SwapWith", "Set", "Unset", "SetTags", "UnsetTags", "EnableManagedAccess", "DisableManagedAccess").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ValidIdentifierIfSet, "SwapWith"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-schema",
	g.NewQueryStruct("DropSchema").
		Drop().
		SQL("SCHEMA").
		IfExists().
		Name().
		OptionalSQL("CASCADE").
		OptionalSQL("RESTRICT").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
).CustomOperation(
	"Undrop",
	"https://docs.snowflake.com/en/sql-reference/sql/undrop-schema",
	g.NewQueryStruct("UndropSchema").
		SQL("UNDROP").
		SQL("SCHEMA").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationFromObject(ssotdefs.Schema).
	DescribeOperationFromObject(ssotdefs.Schema).
	ShowParameters("DatabaseObjectIdentifier").
	ShowParametersDetails(schemaParameters...).
	WithCustomInterfaceMethod(
		"Use", "Use is based on https://docs.snowflake.com/en/sql-reference/sql/use-schema",
		[]*g.MethodParameter{g.NewMethodParameter("id", "DatabaseObjectIdentifier")},
		"error",
	)
