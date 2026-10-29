package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// catalogIntegrationObjectStorageAttributesModificationExt adds behavior the generator does not emit.
// catalog_source is describe-only. table_format is an SSOT attribute and stays in the generated schema.
func catalogIntegrationObjectStorageAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["catalog_source"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: joinWithSpace("Specifies the type of catalog source.",
			"This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens."),
	}
	return s
}
