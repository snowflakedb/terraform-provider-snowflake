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
	AlterOperationFromObject(ssotdefs.Schema).
	DropOperationFromObject(ssotdefs.Schema).
	ShowOperationFromObject(ssotdefs.Schema).
	DescribeOperationFromObject(ssotdefs.Schema).
	ShowParameters("DatabaseObjectIdentifier").
	ShowParametersDetails(schemaParameters...).
	CustomOperationsFromObject(ssotdefs.Schema).
	CustomInterfaceMethodsFromObject(ssotdefs.Schema)
