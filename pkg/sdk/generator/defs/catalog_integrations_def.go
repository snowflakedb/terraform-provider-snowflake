package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var catalogIntegrationsDef = g.NewInterface(
	"CatalogIntegrations",
	"CatalogIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateVariantOperationsFromObject(ssotdefs.CatalogIntegration).
	AlterOperationFromObject(ssotdefs.CatalogIntegration).
	DropOperationFromObject(ssotdefs.CatalogIntegration).
	ShowOperationFromObject(ssotdefs.CatalogIntegration).
	DescribeOperationFromObject(ssotdefs.CatalogIntegration).
	WithEnums(ssotdefs.CatalogIntegration.AdditionalEnums...)
