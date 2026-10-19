package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var catalogIntegrationIcebergRestIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier (i.e. name) of the catalog integration; must be unique in your account.",
	},
}

var catalogIntegrationIcebergRestAttributesSchemaExt = func() map[string]*schema.Schema {
	authExactlyOneOf := []string{"oauth_rest_authentication", "bearer_rest_authentication", "sigv4_rest_authentication"}

	return map[string]*schema.Schema{
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
			Description:  "Specifies the default namespace for all Iceberg tables that you associate with the catalog integration.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"rest_config": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Specifies information about REST configuration.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"catalog_uri": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Description:  "Specifies the endpoint URL for the catalog REST API.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"prefix": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Specifies an optional prefix appended to all API routes.",
					},
					"catalog_name": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						Description:  "Specifies the catalog or identifier to request from your remote catalog service.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"catalog_api_type": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						Description:      "Specifies the connection type for the catalog API. " + enumValuesDescription(sdk.AllCatalogIntegrationCatalogApiTypes),
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
		"oauth_rest_authentication": {
			Type:         schema.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: authExactlyOneOf,
			Description:  "Specifies OAuth as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"oauth_token_uri": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						Description:      "Specifies URL for the third-party identity provider. If not specified, Snowflake assumes the remote catalog provider is the identity provider.",
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("oauth_rest_authentication.0.oauth_token_uri"),
						ValidateFunc:     validation.StringIsNotEmpty,
					},
					"oauth_client_id": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Sensitive:    true,
						Description:  "Specifies the client ID of the OAuth2 credential.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"oauth_client_secret": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						Description:  externalChangesNotDetectedFieldDescription("Specifies the secret of the OAuth2 credential."),
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
		"bearer_rest_authentication": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: authExactlyOneOf,
			Description:  "Specifies a bearer token as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bearer_token": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						Description:  externalChangesNotDetectedFieldDescription("The bearer token for the identity provider."),
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"sigv4_rest_authentication": {
			Type:         schema.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: authExactlyOneOf,
			Description:  "Specifies Signature Version 4 as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"sigv4_iam_role": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Description:  "Specifies the Amazon Resource Name (ARN) for an IAM role that has permission to access your REST API in API Gateway.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"sigv4_signing_region": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						Description:  "Specifies the AWS Region associated with your API in API Gateway. If you don’t specify this parameter, Snowflake uses the region in which your Snowflake account is deployed.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"sigv4_external_id": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
						Description: externalChangesNotDetectedFieldDescription(joinWithSpace("Specifies an external ID that Snowflake uses to establish a trust relationship with AWS.",
							"If you don’t specify this parameter, Snowflake automatically generates a unique external ID when you create a catalog integration.")),
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}()

var catalogIntegrationIcebergRestShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CATALOG INTEGRATIONS` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowCatalogIntegrationSchema,
		},
	},
}

var catalogIntegrationIcebergRestDescribeOutputSchemaExt = map[string]*schema.Schema{
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE CATALOG INTEGRATION` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeCatalogIntegrationIcebergRestDetailsSchema,
		},
	},
}
