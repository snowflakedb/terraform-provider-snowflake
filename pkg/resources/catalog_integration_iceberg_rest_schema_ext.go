package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Nested blocks and ExactlyOneOf are generated. This adds sensitive fields, empty-string checks,
// describe diff suppressors, and the describe-only catalog_source attribute.
func catalogIntegrationIcebergRestAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["refresh_interval_seconds"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds")
	s["catalog_namespace"].ValidateFunc = validation.StringIsNotEmpty

	s["catalog_source"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: joinWithSpace("Specifies the type of catalog source.",
			"This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens."),
	}

	restConfig := s["rest_config"].Elem.(*schema.Resource).Schema
	restConfig["catalog_uri"].ValidateFunc = validation.StringIsNotEmpty
	restConfig["catalog_name"].ValidateFunc = validation.StringIsNotEmpty
	restConfig["catalog_api_type"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToCatalogIntegrationCatalogApiType), IgnoreChangeToCurrentSnowflakeValueInDescribe("rest_config.0.catalog_api_type"))
	restConfig["access_delegation_mode"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToCatalogIntegrationAccessDelegationMode), IgnoreChangeToCurrentSnowflakeValueInDescribe("rest_config.0.access_delegation_mode"))

	oauth := s["oauth_rest_authentication"].Elem.(*schema.Resource).Schema
	oauth["oauth_token_uri"].ValidateFunc = validation.StringIsNotEmpty
	oauth["oauth_token_uri"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInDescribe("oauth_rest_authentication.0.oauth_token_uri")
	oauth["oauth_client_id"].ValidateFunc = validation.StringIsNotEmpty
	oauth["oauth_client_id"].Sensitive = true
	oauth["oauth_client_secret"].ValidateFunc = validation.StringIsNotEmpty
	oauth["oauth_client_secret"].Sensitive = true
	oauth["oauth_client_secret"].Description = externalChangesNotDetectedFieldDescription("Specifies the secret of the OAuth2 credential.")

	bearer := s["bearer_rest_authentication"].Elem.(*schema.Resource).Schema
	bearer["bearer_token"].ValidateFunc = validation.StringIsNotEmpty
	bearer["bearer_token"].Sensitive = true
	bearer["bearer_token"].Description = externalChangesNotDetectedFieldDescription("The bearer token for the identity provider.")

	sigv4 := s["sigv4_rest_authentication"].Elem.(*schema.Resource).Schema
	sigv4["sigv4_iam_role"].ValidateFunc = validation.StringIsNotEmpty
	sigv4["sigv4_signing_region"].ValidateFunc = validation.StringIsNotEmpty
	sigv4["sigv4_external_id"].ValidateFunc = validation.StringIsNotEmpty
	sigv4["sigv4_external_id"].Description = externalChangesNotDetectedFieldDescription(joinWithSpace(
		"Specifies an external ID that Snowflake uses to establish a trust relationship with AWS.",
		"If you don’t specify this parameter, Snowflake automatically generates a unique external ID when you create a catalog integration.",
	))
	return s
}
