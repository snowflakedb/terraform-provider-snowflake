package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var CatalogIntegrationCatalogSourceType = ssot.NewEnum(
	"CatalogIntegrationCatalogSourceType", "CatalogIntegrationCatalogSourceTypes",
	"GLUE", "OBJECT_STORE", "POLARIS", "ICEBERG_REST", "SAP_BDC",
)

var CatalogIntegrationTableFormat = ssot.NewEnum(
	"CatalogIntegrationTableFormat", "CatalogIntegrationTableFormats",
	"ICEBERG", "DELTA",
)

var CatalogIntegrationCatalogApiType = ssot.NewEnum(
	"CatalogIntegrationCatalogApiType", "CatalogIntegrationCatalogApiTypes",
	"PUBLIC", "PRIVATE", "AWS_API_GATEWAY", "AWS_PRIVATE_API_GATEWAY", "AWS_GLUE", "AWS_PRIVATE_GLUE",
)

var CatalogIntegrationAccessDelegationMode = ssot.NewEnum(
	"CatalogIntegrationAccessDelegationMode", "CatalogIntegrationAccessDelegationModes",
	"VENDED_CREDENTIALS", "EXTERNAL_VOLUME_CREDENTIALS",
)

// CatalogIntegrationOpenCatalogAttrs holds reusable Open Catalog integration fields so SHOW
// and DESCRIBE Columns can point at the same definitions instead of restating names.
//
// rest_config and rest_authentication are nested blocks. Sensitive fields, empty-string checks,
// and describe diff suppressors stay in the schema ext.
var CatalogIntegrationOpenCatalogAttrs = struct {
	Enabled                ssot.Attribute
	RefreshIntervalSeconds ssot.Attribute
	Comment                ssot.Attribute
	CatalogNamespace       ssot.Attribute
	RestConfig             ssot.Attribute
	RestAuthentication     ssot.Attribute
	Id                     ssot.Attribute
	CatalogSource          ssot.Attribute
	TableFormat            ssot.Attribute
	Name                   ssot.Attribute
	Type                   ssot.Attribute
	Category               ssot.Attribute
	CreatedOn              ssot.Attribute
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
		CreateMapping: ssot.IntAttributeCreateBuilder,
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
		Description: "Specifies the default Open Catalog namespace for all Iceberg tables that you associate with the catalog integration.",
	},
	RestConfig: ssot.Attribute{
		Name:        "rest_config",
		Type:        ssot.AttrObject,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies information about the Open Catalog account and catalog name.",
		Children: []ssot.Attribute{
			{
				Name:        "catalog_uri",
				Type:        ssot.AttrString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies Open Catalog account URL.",
			},
			{
				Name:        "catalog_name",
				Type:        ssot.AttrString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the name of the catalog to use in Open Catalog.",
			},
			{
				Name:        "catalog_api_type",
				Type:        ssot.AttrEnum,
				ForceNew:    true,
				Description: "Specifies how Snowflake connects to Open Catalog.",
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
	RestAuthentication: ssot.Attribute{
		Name:        "rest_authentication",
		Type:        ssot.AttrObject,
		Required:    true,
		Description: "Specifies authentication details that Snowflake uses to connect to Open Catalog.",
		Children: []ssot.Attribute{
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
				Description: "Specifies the client ID of the OAuth2 credential associated with your Open Catalog service connection.",
			},
			{
				Name:        "oauth_client_secret",
				Type:        ssot.AttrString,
				Required:    true,
				Description: "Specifies the secret of the OAuth2 credential associated with your Open Catalog service connection.",
			},
			{
				Name:        "oauth_allowed_scopes",
				Type:        ssot.AttrList,
				Required:    true,
				ForceNew:    true,
				MinItems:    new(1),
				Description: "Specifies one or more scopes for the OAuth token.",
				Elem:        &ssot.Attribute{Type: ssot.AttrString},
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

var CatalogIntegrationOpenCatalog = ssot.Object{
	Name:                   "CatalogIntegrationOpenCatalog",
	NamePlural:             "CatalogIntegrations",
	VariantOf:              &CatalogIntegration,
	DescribeFailRead:       true,
	IdType:                 "AccountObjectIdentifier",
	IdentityNameForceNew:   true,
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		// ShowOutputShow omitted. SHOW rows are sdk.CatalogIntegration, so tests use CatalogIntegrationShowOutput.
		// WithSchema omitted. SHOW is the shared CatalogIntegration schema alias in the ext.
		// DESCRIBE stays on the SDK schema list until ManualFields exist in SSOT.
		// WithModel omitted. The model generator cannot emit setters for nested blocks.
		// WithResourceAssert omitted. catalog_source is added in the schema ext, so the schema-based assert still emits it.
		// ObjectAssertShow omitted. There is no sdk.CatalogIntegrationOpenCatalog type.
		// ObjectAssertDescribe omitted. The SDK list keeps DescribeOpenCatalog and the nested rest asserts.
		ssot.WithShowOutputAssert(ssot.ShowOutputDescribe()),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		// Snowflake's command is DESCRIBE CATALOG INTEGRATION. The generated text would use the resource name.
		ssot.WithDescribeOutputSchema(ssot.ExtCall()),
		ssot.WithHandleExternalChangesInDescribe(ssot.ExtCall()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		// SHOW has name, enabled, and comment only. catalog_source and catalog_namespace
		// come from DESCRIBE, so SetConfigFields stays on the describe hook.
		ssot.WithSetConfigFields(ssot.Omitted()),
	},
	Attributes: []ssot.Attribute{
		CatalogIntegrationOpenCatalogAttrs.Enabled,
		CatalogIntegrationOpenCatalogAttrs.RefreshIntervalSeconds,
		CatalogIntegrationOpenCatalogAttrs.Comment,
		CatalogIntegrationOpenCatalogAttrs.CatalogNamespace,
		CatalogIntegrationOpenCatalogAttrs.RestConfig,
		CatalogIntegrationOpenCatalogAttrs.RestAuthentication,
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationOpenCatalogAttrs.Name,
			&CatalogIntegrationOpenCatalogAttrs.Enabled,
			&CatalogIntegrationOpenCatalogAttrs.Type,
			&CatalogIntegrationOpenCatalogAttrs.Category,
			&CatalogIntegrationOpenCatalogAttrs.Comment,
			&CatalogIntegrationOpenCatalogAttrs.CreatedOn,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationOpenCatalogAttrs.Id,
			&CatalogIntegrationOpenCatalogAttrs.CatalogSource,
			&CatalogIntegrationOpenCatalogAttrs.TableFormat,
			&CatalogIntegrationOpenCatalogAttrs.Enabled,
			&CatalogIntegrationOpenCatalogAttrs.RefreshIntervalSeconds,
			&CatalogIntegrationOpenCatalogAttrs.Comment,
			&CatalogIntegrationOpenCatalogAttrs.CatalogNamespace,
		},
	},
}
