package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Snowflake describes every catalog integration with DESCRIBE CATALOG INTEGRATION.
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

// Nested rest_config and authentication blocks have no SSOT struct type.
// Scalar fields are generated; this only adds those blocks and the describe-only catalog_source attribute.
func catalogIntegrationIcebergRestAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	authExactlyOneOf := []string{"oauth_rest_authentication", "bearer_rest_authentication", "sigv4_rest_authentication"}

	s["refresh_interval_seconds"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds")
	s["catalog_namespace"].ValidateFunc = validation.StringIsNotEmpty

	s["catalog_source"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: joinWithSpace("Specifies the type of catalog source.",
			"This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens."),
	}
	s["rest_config"] = &schema.Schema{
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
	}
	s["oauth_rest_authentication"] = &schema.Schema{
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
	}
	s["bearer_rest_authentication"] = &schema.Schema{
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
	}
	s["sigv4_rest_authentication"] = &schema.Schema{
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
	}
	return s
}
