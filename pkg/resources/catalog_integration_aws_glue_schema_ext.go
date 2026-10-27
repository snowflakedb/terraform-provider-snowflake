package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// catalogIntegrationAwsGlueAttributesModificationExt adds the computed catalog_source field
// and the describe-output diff suppressors the attribute generator does not emit.
// AWS Glue has no nested blocks.
func catalogIntegrationAwsGlueAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["refresh_interval_seconds"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds")
	s["glue_aws_role_arn"].ValidateFunc = validation.StringIsNotEmpty
	s["glue_catalog_id"].ValidateFunc = validation.StringIsNotEmpty
	s["glue_region"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInDescribe("glue_region")
	s["glue_region"].ValidateFunc = validation.StringIsNotEmpty
	s["catalog_namespace"].ValidateFunc = validation.StringIsNotEmpty
	s["catalog_source"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: joinWithSpace("Specifies the type of catalog source.",
			"This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens."),
	}
	return s
}
