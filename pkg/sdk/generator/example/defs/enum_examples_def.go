//go:build sdk_generation_examples

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// ExampleStatusDef demonstrates a simple enum without aliases.
var ExampleStatusDef = g.NewEnum(
	"ExampleStatus",
	"ExampleStatuses",
	"ACTIVE", "INACTIVE", "EXPIRED",
)

// ExampleSizeDef demonstrates an enum with aliases, similar to WarehouseSize.
var ExampleSizeDef = g.NewEnum(
	"ExampleSize",
	"ExampleSizes",
	"XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "XXLARGE",
).WithAliases("XSMALL", "X-SMALL").
	WithAliases("XLARGE", "X-LARGE").
	WithAliases("XXLARGE", "X2LARGE", "2X-LARGE")

var EnumExamplesDef = g.NewInterface(
	"EnumExamples",
	"EnumExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).WithAllowedGenerationParts("enums").
	WithEnums(ExampleStatusDef, ExampleSizeDef)
