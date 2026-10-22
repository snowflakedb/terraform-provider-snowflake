package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var schemasDef = g.NewInterface(
	"Schemas",
	"Schema",
	g.KindOfT[sdkcommons.DatabaseObjectIdentifier](),
).CreateOperationFromObject(ssotdefs.Schema).
	AlterOperationFromObject(ssotdefs.Schema).
	DropOperationFromObject(ssotdefs.Schema).
	ShowOperationFromObject(ssotdefs.Schema).
	DescribeOperationFromObject(ssotdefs.Schema).
	ShowParametersFromObject(ssotdefs.Schema).
	ShowParametersDetailsFromObject(ssotdefs.Schema).
	CustomOperationsFromObject(ssotdefs.Schema).
	CustomInterfaceMethodsFromObject(ssotdefs.Schema)
