package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// CatalogIntegration is the base of the catalog integration variants. All of them are the same Snowflake object,
// CATALOG INTEGRATION, and share one SDK interface, which is built from this object.
var CatalogIntegration = ssot.Object{
	Name:            "CatalogIntegration",
	NamePlural:      "CatalogIntegrations",
	IdType:          "AccountObjectIdentifier",
	VariantBaseOnly: true,
}
