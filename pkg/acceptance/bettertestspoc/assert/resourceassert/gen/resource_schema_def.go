package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceSchemaDef struct {
	name   string
	schema map[string]*schema.Schema
}

func GetResourceSchemaDetails() []genhelpers.ResourceSchemaDetails {
	allResourceSchemas := allResourceSchemaDefs
	allResourceSchemasDetails := make([]genhelpers.ResourceSchemaDetails, len(allResourceSchemas))
	for idx, s := range allResourceSchemas {
		allResourceSchemasDetails[idx] = genhelpers.ExtractResourceSchemaDetails(s.name, s.schema)
	}
	return allResourceSchemasDetails
}

var allResourceSchemaDefs = []ResourceSchemaDef{
	{
		name:   "Warehouse",
		schema: resources.Warehouse().Schema,
	},
	{
		name:   "User",
		schema: resources.User().Schema,
	},
	{
		name:   "ServiceUser",
		schema: resources.ServiceUser().Schema,
	},
	{
		name:   "LegacyServiceUser",
		schema: resources.LegacyServiceUser().Schema,
	},
	{
		name:   "View",
		schema: resources.View().Schema,
	},
	{
		name:   "Database",
		schema: resources.Database().Schema,
	},
	{
		name:   "DatabaseRole",
		schema: resources.DatabaseRole().Schema,
	},
	{
		name:   "ResourceMonitor",
		schema: resources.ResourceMonitor().Schema,
	},
	{
		name:   "RowAccessPolicy",
		schema: resources.RowAccessPolicy().Schema,
	},
	{
		name:   "Schema",
		schema: resources.Schema().Schema,
	},
	{
		name:   "MaskingPolicy",
		schema: resources.MaskingPolicy().Schema,
	},
	{
		name:   "StreamOnTable",
		schema: resources.StreamOnTable().Schema,
	},
	{
		name:   "StreamOnExternalTable",
		schema: resources.StreamOnExternalTable().Schema,
	},
	{
		name:   "SecretWithAuthorizationCodeGrant",
		schema: resources.SecretWithAuthorizationCodeGrant().Schema,
	},
	{
		name:   "SecretWithBasicAuthentication",
		schema: resources.SecretWithBasicAuthentication().Schema,
	},
	{
		name:   "SecretWithClientCredentials",
		schema: resources.SecretWithClientCredentials().Schema,
	},
	{
		name:   "SecretWithGenericString",
		schema: resources.SecretWithGenericString().Schema,
	},
	{
		name:   "StreamOnDirectoryTable",
		schema: resources.StreamOnDirectoryTable().Schema,
	},
	{
		name:   "StreamOnView",
		schema: resources.StreamOnView().Schema,
	},
	{
		name:   "PrimaryConnection",
		schema: resources.PrimaryConnection().Schema,
	},
	{
		name:   "SecondaryConnection",
		schema: resources.SecondaryConnection().Schema,
	},
	{
		name:   "Tag",
		schema: resources.Tag().Schema,
	},
	{
		name:   "TagAssociation",
		schema: resources.TagAssociation().Schema,
	},
	{
		name:   "Task",
		schema: resources.Task().Schema,
	},
	{
		name:   "Account",
		schema: resources.Account().Schema,
	},
	{
		name:   "AccountParameter",
		schema: resources.AccountParameter().Schema,
	},
	{
		name:   "OauthIntegrationForCustomClients",
		schema: resources.OauthIntegrationForCustomClients().Schema,
	},
	{
		name:   "OauthIntegrationForPartnerApplications",
		schema: resources.OauthIntegrationForPartnerApplications().Schema,
	},
	{
		name:   "FunctionJava",
		schema: resources.FunctionJava().Schema,
	},
	{
		name:   "FunctionJavascript",
		schema: resources.FunctionJavascript().Schema,
	},
	{
		name:   "FunctionPython",
		schema: resources.FunctionPython().Schema,
	},
	{
		name:   "FunctionScala",
		schema: resources.FunctionScala().Schema,
	},
	{
		name:   "FunctionSql",
		schema: resources.FunctionSql().Schema,
	},
	{
		name:   "ProcedureJava",
		schema: resources.ProcedureJava().Schema,
	},
	{
		name:   "ProcedureJavascript",
		schema: resources.ProcedureJavascript().Schema,
	},
	{
		name:   "ProcedurePython",
		schema: resources.ProcedurePython().Schema,
	},
	{
		name:   "ProcedureScala",
		schema: resources.ProcedureScala().Schema,
	},
	{
		name:   "ProcedureSql",
		schema: resources.ProcedureSql().Schema,
	},
	{
		name:   "ManagedAccount",
		schema: resources.ManagedAccount().Schema,
	},
	{
		name:   "NetworkPolicy",
		schema: resources.NetworkPolicy().Schema,
	},
	{
		name:   "ScimSecurityIntegration",
		schema: resources.SCIMIntegration().Schema,
	},
	{
		name:   "Saml2SecurityIntegration",
		schema: resources.SAML2Integration().Schema,
	},
	{
		name:   "ApiAuthenticationIntegrationWithAuthorizationCodeGrant",
		schema: resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant().Schema,
	},
	{
		name:   "ApiAuthenticationIntegrationWithClientCredentials",
		schema: resources.ApiAuthenticationIntegrationWithClientCredentials().Schema,
	},
	{
		name:   "ExternalOauthSecurityIntegration",
		schema: resources.ExternalOauthIntegration().Schema,
	},
	{
		name:   "Streamlit",
		schema: resources.Streamlit().Schema,
	},
	{
		name:   "AccountRole",
		schema: resources.AccountRole().Schema,
	},
	{
		name:   "SecondaryDatabase",
		schema: resources.SecondaryDatabase().Schema,
	},
	{
		name:   "SharedDatabase",
		schema: resources.SharedDatabase().Schema,
	},
	{
		name:   "Table",
		schema: resources.Table().Schema,
	},
	{
		name:   "ImageRepository",
		schema: resources.ImageRepository().Schema,
	},
}
