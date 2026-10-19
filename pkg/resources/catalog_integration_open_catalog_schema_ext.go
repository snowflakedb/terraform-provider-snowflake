package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var catalogIntegrationOpenCatalogIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier (i.e. name) of the catalog integration; must be unique in your account.",
	},
}

var catalogIntegrationOpenCatalogAttributesSchemaExt = map[string]*schema.Schema{
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
		Description: joinWithSpace("Specifies whether the catalog integration is available for use for Iceberg tables.",
			"`true` allows users to create new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration function normally.",
			"`false` prevents users from creating new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration cannot access the catalog in the table definition."),
	},
	"refresh_interval_seconds": {
		Type:     schema.TypeInt,
		Optional: true,
		Description: joinWithSpace("Specifies the number of seconds to wait between attempts to poll the external Iceberg catalog for metadata updates for automated refresh.",
			"For Delta-based tables, specifies the number of seconds to wait between attempts to poll your external cloud storage for new metadata."),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds"),
		ValidateFunc:     validation.IntAtLeast(1),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies a comment for the catalog integration.",
	},
	"catalog_source": {
		Type:     schema.TypeString,
		Computed: true,
		Description: joinWithSpace("Specifies the type of catalog source.",
			"This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens."),
	},
	"catalog_namespace": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Specifies the default Open Catalog namespace for all Iceberg tables that you associate with the catalog integration.",
		ValidateFunc: validation.StringIsNotEmpty,
	},
	"rest_config": {
		Type:        schema.TypeList,
		Required:    true,
		ForceNew:    true,
		MinItems:    1,
		MaxItems:    1,
		Description: "Specifies information about the Open Catalog account and catalog name.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"catalog_uri": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					Description:  "Specifies Open Catalog account URL.",
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"catalog_name": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					Description:  "Specifies the name of the catalog to use in Open Catalog.",
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"catalog_api_type": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "Specifies how Snowflake connects to Open Catalog. " + enumValuesDescription(sdk.AllCatalogIntegrationCatalogApiTypes),
					DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToCatalogIntegrationCatalogApiType), IgnoreChangeToCurrentSnowflakeValueInDescribe("rest_config.0.catalog_api_type")),
					ValidateDiagFunc: sdkValidation(sdk.ToCatalogIntegrationCatalogApiType),
				},
				"access_delegation_mode": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "Specifies the access delegation mode for accessing Iceberg table files in your external cloud storage. " + enumValuesDescription(sdk.AllCatalogIntegrationAccessDelegationModes),
					DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToCatalogIntegrationAccessDelegationMode), IgnoreChangeToCurrentSnowflakeValueInDescribe("rest_config.0.access_delegation_mode")),
					ValidateDiagFunc: sdkValidation(sdk.ToCatalogIntegrationAccessDelegationMode),
				},
			},
		},
	},
	"rest_authentication": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		Description: "Specifies authentication details that Snowflake uses to connect to Open Catalog.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"oauth_token_uri": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "Specifies URL for the third-party identity provider. If not specified, Snowflake assumes the remote catalog provider is the identity provider.",
					DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("rest_authentication.0.oauth_token_uri"),
					ValidateFunc:     validation.StringIsNotEmpty,
				},
				"oauth_client_id": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					Sensitive:    true,
					Description:  "Specifies the client ID of the OAuth2 credential associated with your Open Catalog service connection.",
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"oauth_client_secret": {
					Type:         schema.TypeString,
					Required:     true,
					Sensitive:    true,
					Description:  externalChangesNotDetectedFieldDescription("Specifies the secret of the OAuth2 credential associated with your Open Catalog service connection."),
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"oauth_allowed_scopes": {
					Type:        schema.TypeList,
					Required:    true,
					ForceNew:    true,
					MinItems:    1,
					Description: "Specifies one or more scopes for the OAuth token.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	},
}

var catalogIntegrationOpenCatalogShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CATALOG INTEGRATIONS` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowCatalogIntegrationSchema,
		},
	},
}

var catalogIntegrationOpenCatalogDescribeOutputSchemaExt = map[string]*schema.Schema{
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE CATALOG INTEGRATION` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeCatalogIntegrationOpenCatalogDetailsSchema,
		},
	},
}
