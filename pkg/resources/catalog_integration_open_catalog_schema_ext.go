package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// catalogIntegrationOpenCatalogAttributesModificationExt adds fields the generator does not emit.
// Nested blocks are generated. Sensitive fields, empty-string checks, and describe diffs stay here.
// catalog_source is describe-only.
func catalogIntegrationOpenCatalogAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
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
	restConfig["catalog_api_type"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToCatalogIntegrationCatalogApiType), IgnoreChangeToCurrentSnowflakePlainValueInOutput("describe_output.0.rest_config", "catalog_api_type"))
	restConfig["access_delegation_mode"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToCatalogIntegrationAccessDelegationMode), IgnoreChangeToCurrentSnowflakePlainValueInOutput("describe_output.0.rest_config", "access_delegation_mode"))

	restAuth := s["rest_authentication"].Elem.(*schema.Resource).Schema
	restAuth["oauth_token_uri"].ValidateFunc = validation.StringIsNotEmpty
	restAuth["oauth_token_uri"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakePlainValueInOutput("describe_output.0.rest_authentication", "oauth_token_uri")
	restAuth["oauth_client_id"].ValidateFunc = validation.StringIsNotEmpty
	restAuth["oauth_client_id"].Sensitive = true
	restAuth["oauth_client_secret"].ValidateFunc = validation.StringIsNotEmpty
	restAuth["oauth_client_secret"].Sensitive = true
	restAuth["oauth_client_secret"].Description = externalChangesNotDetectedFieldDescription("Specifies the secret of the OAuth2 credential associated with your Open Catalog service connection.")
	return s
}
