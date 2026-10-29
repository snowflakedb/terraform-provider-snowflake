package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// CatalogIntegrationIcebergRestAttrs holds reusable Iceberg REST catalog integration fields
// so SHOW and DESCRIBE Columns can point at the same definitions instead of restating names.
//
// rest_config and the authentication blocks are nested. Sensitive is on the attributes.
// Empty-string checks and describe diff suppressors stay in the schema ext.
var CatalogIntegrationIcebergRestAttrs = struct {
	Enabled                  ssot.Attribute
	RefreshIntervalSeconds   ssot.Attribute
	Comment                  ssot.Attribute
	CatalogNamespace         ssot.Attribute
	RestConfig               ssot.Attribute
	OauthRestAuthentication  ssot.Attribute
	BearerRestAuthentication ssot.Attribute
	Sigv4RestAuthentication  ssot.Attribute
	Id                       ssot.Attribute
	CatalogSource            ssot.Attribute
	TableFormat              ssot.Attribute
	Name                     ssot.Attribute
	Type                     ssot.Attribute
	Category                 ssot.Attribute
	CreatedOn                ssot.Attribute
}{
	Enabled: ssot.Attribute{
		Name:     "enabled",
		Type:     ssot.AttrBool,
		Required: true,
		Description: "Specifies whether the catalog integration is available for use for Iceberg tables. " +
			"`true` allows users to create new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration function normally. " +
			"`false` prevents users from creating new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration cannot access the catalog in the table definition.",
	},
	RefreshIntervalSeconds: ssot.Attribute{
		Name: "refresh_interval_seconds",
		Type: ssot.AttrInt,
		Description: "Specifies the number of seconds to wait between attempts to poll the external Iceberg catalog for metadata updates for automated refresh. " +
			"For Delta-based tables, specifies the number of seconds to wait between attempts to poll your external cloud storage for new metadata.",
		Min:           new(1),
		ServerDefault: true,
		CreateMapping: ssot.IntAttributeCreateBuilder,
		CreateSql:     ssot.IntParameter(ssot.NoQuotes()),
		DiffSuppress:  ssot.SuppressDescribe,
	},
	Comment: ssot.Attribute{
		Name:          "comment",
		Type:          ssot.AttrString,
		Description:   "Specifies a comment for the catalog integration.",
		OptionalInDb:  true,
		CreateMapping: ssot.StringAttributeCreateBuilder,
	},
	CatalogNamespace: ssot.Attribute{
		Name:        "catalog_namespace",
		Type:        ssot.AttrString,
		ForceNew:    true,
		Description: "Specifies the default namespace for all Iceberg tables that you associate with the catalog integration.",
	},
	RestConfig: ssot.Attribute{
		Name:        "rest_config",
		Type:        ssot.AttrObject,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies information about REST configuration.",
		CreateSql:   ssot.ObjectParameter(ssot.Sql("REST_CONFIG ="), ssot.Parentheses()),
		Children: []ssot.Attribute{
			{
				Name:        "catalog_uri",
				Type:        ssot.AttrString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the endpoint URL for the catalog REST API.",
			},
			{
				Name:        "prefix",
				Type:        ssot.AttrString,
				ForceNew:    true,
				Description: "Specifies an optional prefix appended to all API routes.",
			},
			{
				Name:        "catalog_name",
				Type:        ssot.AttrString,
				ForceNew:    true,
				Description: "Specifies the catalog or identifier to request from your remote catalog service.",
			},
			{
				Name:        "catalog_api_type",
				Type:        ssot.AttrEnum,
				ForceNew:    true,
				Description: "Specifies the connection type for the catalog API.",
				Enum:        CatalogIntegrationCatalogApiType,
			},
			{
				Name:        "access_delegation_mode",
				Type:        ssot.AttrEnum,
				ForceNew:    true,
				Description: "Specifies the access delegation mode for accessing Iceberg table files in your external cloud storage.",
				Enum:        CatalogIntegrationAccessDelegationMode,
			},
		},
	},
	OauthRestAuthentication: ssot.Attribute{
		Name:        "oauth_rest_authentication",
		Type:        ssot.AttrObject,
		ForceNew:    true,
		Description: "Specifies OAuth as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
		CreateSql:   ssot.ObjectParameter(ssot.Sql("REST_AUTHENTICATION ="), ssot.Parentheses(), ssot.FieldName("OAuthRestAuthentication")),
		Children: []ssot.Attribute{
			catalogIntegrationRestAuthTypeOauth,
			{
				Name:        "oauth_token_uri",
				Type:        ssot.AttrString,
				ForceNew:    true,
				Description: "Specifies URL for the third-party identity provider. If not specified, Snowflake assumes the remote catalog provider is the identity provider.",
			},
			{
				Name:        "oauth_client_id",
				Type:        ssot.AttrString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Specifies the client ID of the OAuth2 credential.",
			},
			{
				Name:        "oauth_client_secret",
				Type:        ssot.AttrString,
				Required:    true,
				Sensitive:   true,
				Description: "Specifies the secret of the OAuth2 credential.",
			},
			{
				Name:        "oauth_allowed_scopes",
				Type:        ssot.AttrList,
				Required:    true,
				ForceNew:    true,
				MinItems:    new(1),
				Description: "Specifies one or more scopes for the OAuth token.",
				Elem:        &ssot.Attribute{Type: ssot.AttrString},
				CreateSql:   ssot.ListParameter("StringListItemWrapper"),
			},
		},
	},
	BearerRestAuthentication: ssot.Attribute{
		Name:        "bearer_rest_authentication",
		Type:        ssot.AttrObject,
		Description: "Specifies a bearer token as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
		CreateSql:   ssot.ObjectParameter(ssot.Sql("REST_AUTHENTICATION ="), ssot.Parentheses()),
		Children: []ssot.Attribute{
			catalogIntegrationIcebergRestAuthTypeBearer,
			{
				Name:        "bearer_token",
				Type:        ssot.AttrString,
				Required:    true,
				Sensitive:   true,
				Description: "The bearer token for the identity provider.",
			},
		},
	},
	Sigv4RestAuthentication: ssot.Attribute{
		Name:        "sigv4_rest_authentication",
		Type:        ssot.AttrObject,
		ForceNew:    true,
		Description: "Specifies Signature Version 4 as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
		CreateSql:   ssot.ObjectParameter(ssot.Sql("REST_AUTHENTICATION ="), ssot.Parentheses(), ssot.FieldName("SigV4RestAuthentication")),
		Children: []ssot.Attribute{
			catalogIntegrationIcebergRestAuthTypeSigv4,
			{
				Name:        "sigv4_iam_role",
				Type:        ssot.AttrString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the Amazon Resource Name (ARN) for an IAM role that has permission to access your REST API in API Gateway.",
			},
			{
				Name:        "sigv4_signing_region",
				Type:        ssot.AttrString,
				ForceNew:    true,
				Description: "Specifies the AWS Region associated with your API in API Gateway. If you don’t specify this parameter, Snowflake uses the region in which your Snowflake account is deployed.",
			},
			{
				Name:        "sigv4_external_id",
				Type:        ssot.AttrString,
				ForceNew:    true,
				Description: "Specifies an external ID that Snowflake uses to establish a trust relationship with AWS. If you don’t specify this parameter, Snowflake automatically generates a unique external ID when you create a catalog integration.",
			},
		},
	},
	Id: ssot.Attribute{Name: "id", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", Readonly: true},
	CatalogSource: ssot.Attribute{
		Name:        "catalog_source",
		Type:        ssot.AttrEnum,
		Description: "Specifies the type of catalog source. This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens.",
		Enum:        CatalogIntegrationCatalogSourceType,
		Readonly:    true,
	},
	TableFormat: ssot.Attribute{Name: "table_format", Type: ssot.AttrEnum, Enum: CatalogIntegrationTableFormat, Readonly: true},
	Name:        ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	Type:        ssot.Attribute{Name: "type", Type: ssot.AttrString, Readonly: true},
	Category:    ssot.Attribute{Name: "category", Type: ssot.AttrString, Readonly: true},
	CreatedOn:   ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
}

var (
	catalogIntegrationIcebergRestCatalogSourceStatic = ssot.Attribute{
		Name:      "catalog_source_iceberg_rest",
		SqlOnly:   true,
		CreateSql: ssot.StaticSQL("catalogSource", "CATALOG_SOURCE = ICEBERG_REST"),
	}
	catalogIntegrationIcebergRestAuthTypeBearer = ssot.Attribute{
		Name:      "type_bearer",
		SqlOnly:   true,
		CreateSql: ssot.StaticSQL("restAuthType", "TYPE = BEARER"),
	}
	catalogIntegrationIcebergRestAuthTypeSigv4 = ssot.Attribute{
		Name:      "type_sigv4",
		SqlOnly:   true,
		CreateSql: ssot.StaticSQL("restAuthType", "TYPE = SIGV4"),
	}
)

var CatalogIntegrationIcebergRest = ssot.Object{
	Name:                   "CatalogIntegrationIcebergRest",
	NamePlural:             "CatalogIntegrations",
	VariantOf:              &CatalogIntegration,
	DescribeFailRead:       true,
	IdType:                 "AccountObjectIdentifier",
	IdentityNameForceNew:   true,
	AttributesModification: true,
	// catalog_source is describe-only. The schema ext adds it; assertions read it from here.
	ComputedAttributes: []ssot.Attribute{
		CatalogIntegrationIcebergRestAttrs.CatalogSource,
	},
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithSchema(
			ssot.SchemaShow(ssot.SchemaShowFromBase()),
			ssot.SchemaDescribe(ssot.SchemaDescribeManualFields(
				&CatalogIntegrationIcebergRestAttrs.RestConfig,
				&CatalogIntegrationIcebergRestAttrs.OauthRestAuthentication,
				&CatalogIntegrationIcebergRestAttrs.BearerRestAuthentication,
				&CatalogIntegrationIcebergRestAttrs.Sigv4RestAuthentication,
			)),
		),
		ssot.WithModel(ssot.ModelRequestType(
			&CatalogIntegrationIcebergRestAttrs.RestConfig,
			"[]sdk.CatalogIntegrationIcebergRestRestConfigRequest",
		)),
		ssot.WithResourceAssert(),
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(ssot.ObjectAssertShowRowFromBase()),
			ssot.ObjectAssertDescribe(
				ssot.ObjectAssertDescribeFromBase("DescribeIcebergRest"),
				ssot.ObjectAssertDescribeNestedField(&CatalogIntegrationIcebergRestAttrs.RestConfig),
				ssot.ObjectAssertDescribeNestedField(&CatalogIntegrationIcebergRestAttrs.OauthRestAuthentication),
				ssot.ObjectAssertDescribeNestedField(&CatalogIntegrationIcebergRestAttrs.Sigv4RestAuthentication),
				ssot.ObjectAssertDescribeSkipFields(&CatalogIntegrationIcebergRestAttrs.BearerRestAuthentication),
			),
		),
		// DESCRIBE CATALOG INTEGRATION is shared by every catalog integration variant. The details are parsed in the SDK ext.
		ssot.WithSdk(ssot.SdkDescribe(ssot.SdkDescribeSql("CATALOG INTEGRATION"), ssot.SdkDescribeDetails())),
		ssot.WithShowOutputAssert(
			ssot.ShowOutputShow(),
			ssot.ShowOutputDescribe(ssot.ShowOutputDescribeSkipFields(
				&CatalogIntegrationIcebergRestAttrs.RestConfig,
				&CatalogIntegrationIcebergRestAttrs.OauthRestAuthentication,
				&CatalogIntegrationIcebergRestAttrs.BearerRestAuthentication,
				&CatalogIntegrationIcebergRestAttrs.Sigv4RestAuthentication,
			))),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithAlterUnsetInSdk(ssot.Omitted()),
		ssot.WithHandleExternalChangesInDescribe(ssot.Generated()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		CatalogIntegrationIcebergRestAttrs.Enabled,
		CatalogIntegrationIcebergRestAttrs.RefreshIntervalSeconds,
		CatalogIntegrationIcebergRestAttrs.Comment,
		CatalogIntegrationIcebergRestAttrs.CatalogNamespace,
		CatalogIntegrationIcebergRestAttrs.RestConfig,
		CatalogIntegrationIcebergRestAttrs.OauthRestAuthentication,
		CatalogIntegrationIcebergRestAttrs.BearerRestAuthentication,
		CatalogIntegrationIcebergRestAttrs.Sigv4RestAuthentication,
	},
	ExactlyOneOf: [][]*ssot.Attribute{{
		&CatalogIntegrationIcebergRestAttrs.OauthRestAuthentication,
		&CatalogIntegrationIcebergRestAttrs.BearerRestAuthentication,
		&CatalogIntegrationIcebergRestAttrs.Sigv4RestAuthentication,
	}},
	Create: ssot.CreateSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-rest",
		Clauses: []*ssot.Attribute{
			&catalogIntegrationIcebergRestCatalogSourceStatic,
			&catalogIntegrationTableFormatIceberg,
			&CatalogIntegrationIcebergRestAttrs.CatalogNamespace,
			&CatalogIntegrationIcebergRestAttrs.RestConfig,
			&CatalogIntegrationIcebergRestAttrs.OauthRestAuthentication,
			&CatalogIntegrationIcebergRestAttrs.BearerRestAuthentication,
			&CatalogIntegrationIcebergRestAttrs.Sigv4RestAuthentication,
			&CatalogIntegrationIcebergRestAttrs.Enabled,
			&CatalogIntegrationIcebergRestAttrs.RefreshIntervalSeconds,
			&CatalogIntegrationIcebergRestAttrs.Comment,
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationIcebergRestAttrs.Name,
			&CatalogIntegrationIcebergRestAttrs.Enabled,
			&CatalogIntegrationIcebergRestAttrs.Type,
			&CatalogIntegrationIcebergRestAttrs.Category,
			&CatalogIntegrationIcebergRestAttrs.Comment,
			&CatalogIntegrationIcebergRestAttrs.CreatedOn,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationIcebergRestAttrs.Id,
			&CatalogIntegrationIcebergRestAttrs.CatalogSource,
			&CatalogIntegrationIcebergRestAttrs.TableFormat,
			&CatalogIntegrationIcebergRestAttrs.Enabled,
			&CatalogIntegrationIcebergRestAttrs.RefreshIntervalSeconds,
			&CatalogIntegrationIcebergRestAttrs.Comment,
			&CatalogIntegrationIcebergRestAttrs.CatalogNamespace,
			&CatalogIntegrationIcebergRestAttrs.RestConfig,
			&CatalogIntegrationIcebergRestAttrs.OauthRestAuthentication,
			&CatalogIntegrationIcebergRestAttrs.BearerRestAuthentication,
			&CatalogIntegrationIcebergRestAttrs.Sigv4RestAuthentication,
		},
	},
}
