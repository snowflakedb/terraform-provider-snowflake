package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func alertAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	// enabled keeps its description here. Optional+computed fields stay out of the generator.
	s["enabled"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies if an alert should be 'started' (enabled) after creation or should remain 'suspended' (default).",
	}
	return s
}
