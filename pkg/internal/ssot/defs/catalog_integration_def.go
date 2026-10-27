package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var catalogIntegrationTableFormatIceberg = ssot.Attribute{
	Name:      "table_format_iceberg",
	SqlOnly:   true,
	CreateSql: ssot.StaticSQL("tableFormat", "TABLE_FORMAT = ICEBERG"),
}

// catalogIntegrationRestAuthTypeOauth is the fixed TYPE of an OAuth REST_AUTHENTICATION struct.
var catalogIntegrationRestAuthTypeOauth = ssot.Attribute{
	Name:      "type_oauth",
	SqlOnly:   true,
	CreateSql: ssot.StaticSQL("restAuthType", "TYPE = OAUTH"),
}

// catalogIntegrationSetAttrs are the fields that SET changes on every catalog integration variant. They only carry the SQL of
// the shared ALTER, the variants describe the same fields for Terraform. A REST authentication is set by its secret only.
var catalogIntegrationSetAttrs = struct {
	Enabled                  ssot.Attribute
	RefreshIntervalSeconds   ssot.Attribute
	Comment                  ssot.Attribute
	OauthRestAuthentication  ssot.Attribute
	BearerRestAuthentication ssot.Attribute
}{
	Enabled:                ssot.Attribute{Name: "enabled", Type: ssot.AttrBool},
	RefreshIntervalSeconds: ssot.Attribute{Name: "refresh_interval_seconds", Type: ssot.AttrInt, SetSql: ssot.IntParameter(ssot.NoQuotes())},
	// TODO(SNOW-3121221): use COMMENT in unset and here use a plain string.
	Comment: ssot.Attribute{Name: "comment", Type: ssot.AttrString, SetSql: ssot.StringParameter(ssot.AllowEmpty())},
	OauthRestAuthentication: ssot.Attribute{
		Name:    "set_oauth_rest_authentication",
		Type:    ssot.AttrObject,
		SqlOnly: true,
		SetSql:  ssot.ObjectParameter(ssot.Sql("REST_AUTHENTICATION ="), ssot.Parentheses(), ssot.FieldName("SetOAuthRestAuthentication")),
		Children: []ssot.Attribute{
			{Name: "oauth_client_secret", Type: ssot.AttrString, Required: true},
		},
	},
	BearerRestAuthentication: ssot.Attribute{
		Name:    "set_bearer_rest_authentication",
		Type:    ssot.AttrObject,
		SqlOnly: true,
		SetSql:  ssot.ObjectParameter(ssot.Sql("REST_AUTHENTICATION ="), ssot.Parentheses(), ssot.FieldName("SetBearerRestAuthentication")),
		Children: []ssot.Attribute{
			{Name: "bearer_token", Type: ssot.AttrString, Required: true},
		},
	},
}

// CatalogIntegration is the base of the catalog integration variants. All of them are the same Snowflake object,
// CATALOG INTEGRATION, and share one SDK interface, which is built from this object.
var CatalogIntegration = ssot.Object{
	Name:            "CatalogIntegration",
	NamePlural:      "CatalogIntegrations",
	IdType:          "AccountObjectIdentifier",
	VariantBaseOnly: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithSdk(ssot.SdkCreateForVariants(), ssot.SdkAlter()),
	},
	Set: []*ssot.Attribute{
		&catalogIntegrationSetAttrs.OauthRestAuthentication,
		&catalogIntegrationSetAttrs.BearerRestAuthentication,
		&catalogIntegrationSetAttrs.Enabled,
		&catalogIntegrationSetAttrs.RefreshIntervalSeconds,
		&catalogIntegrationSetAttrs.Comment,
	},
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{ssot.AlterSet(), ssot.AlterSetTags(), ssot.AlterUnsetTags()},
		SetConflictingFields: [][]*ssot.Attribute{{
			&catalogIntegrationSetAttrs.OauthRestAuthentication,
			&catalogIntegrationSetAttrs.BearerRestAuthentication,
		}},
	},
}
