package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderSchemaDef struct {
	name   string
	schema map[string]*schema.Schema
}

// TODO [SNOW-1501905]: rename ResourceSchemaDetails (because it is used for the datasources and provider too)
func GetProviderSchemaDetails() []genhelpers.ResourceSchemaDetails {
	allProvidersSchemas := allProviderSchemaDefs
	allProvidersSchemasDetails := make([]genhelpers.ResourceSchemaDetails, len(allProvidersSchemas))
	for idx, s := range allProvidersSchemas {
		allProvidersSchemasDetails[idx] = genhelpers.ExtractResourceSchemaDetails(s.name, s.schema)
	}
	return allProvidersSchemasDetails
}

var allProviderSchemaDefs = []ProviderSchemaDef{
	{
		name:   "Snowflake",
		schema: provider.Provider().Schema,
	},
}
