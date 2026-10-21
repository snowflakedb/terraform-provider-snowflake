package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

var schemaParameters = ParameterDefsForLevel(parameterdefs.ParameterLevelSchema)

var schemasDef = g.NewInterface(
	"Schemas",
	"Schema",
	g.KindOfT[sdkcommons.DatabaseObjectIdentifier](),
).CreateOperationFromObject(ssotdefs.Schema).
	CustomOperation(
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
	).AlterOperationFromObject(ssotdefs.Schema).
	DropOperationFromObject(ssotdefs.Schema).
	CustomOperation(
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
