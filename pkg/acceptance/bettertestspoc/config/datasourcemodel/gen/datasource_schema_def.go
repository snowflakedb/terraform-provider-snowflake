package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DatasourceSchemaDef struct {
	name   string
	schema map[string]*schema.Schema
}

// TODO [SNOW-1501905]: rename ResourceSchemaDetails (because it is used for the datasources and provider too)
func GetDatasourceSchemaDetails() []genhelpers.ResourceSchemaDetails {
	allDatasourcesSchemas := allDatasourcesSchemaDefs
	allDatasourcesSchemasDetails := make([]genhelpers.ResourceSchemaDetails, len(allDatasourcesSchemas))
	for idx, s := range allDatasourcesSchemas {
		allDatasourcesSchemasDetails[idx] = genhelpers.ExtractResourceSchemaDetails(s.name, s.schema)
	}
	return allDatasourcesSchemasDetails
}

var allDatasourcesSchemaDefs = []DatasourceSchemaDef{
	{
		name:   "Accounts",
		schema: datasources.Accounts().Schema,
	},
	{
		name:   "ComputePools",
		schema: datasources.ComputePools().Schema,
	},
	{
		name:   "Database",
		schema: datasources.Database().Schema,
	},
	{
		name:   "DatabaseRole",
		schema: datasources.DatabaseRole().Schema,
	},
	{
		name:   "DatabaseRoles",
		schema: datasources.DatabaseRoles().Schema,
	},
	{
		name:   "Databases",
		schema: datasources.Databases().Schema,
	},
	{
		name:   "Grants",
		schema: datasources.Grants().Schema,
	},
	{
		name:   "ImageRepositories",
		schema: datasources.ImageRepositories().Schema,
	},
	{
		name:   "MaskingPolicies",
		schema: datasources.MaskingPolicies().Schema,
	},
	{
		name:   "NetworkPolicies",
		schema: datasources.NetworkPolicies().Schema,
	},
	{
		name:   "ResourceMonitors",
		schema: datasources.ResourceMonitors().Schema,
	},
	{
		name:   "Schemas",
		schema: datasources.Schemas().Schema,
	},
	{
		name:   "Secrets",
		schema: datasources.Secrets().Schema,
	},
	{
		name:   "SecurityIntegrations",
		schema: datasources.SecurityIntegrations().Schema,
	},
	{
		name:   "Streamlits",
		schema: datasources.Streamlits().Schema,
	},
	{
		name:   "Streams",
		schema: datasources.Streams().Schema,
	},
	{
		name:   "Tags",
		schema: datasources.Tags().Schema,
	},
	{
		name:   "Tasks",
		schema: datasources.Tasks().Schema,
	},
	{
		name:   "Users",
		schema: datasources.Users().Schema,
	},
	{
		name:   "Views",
		schema: datasources.Views().Schema,
	},
	{
		name:   "Warehouses",
		schema: datasources.Warehouses().Schema,
	},
}
