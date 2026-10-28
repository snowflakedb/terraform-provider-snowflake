package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func alertAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	schedule := s["alert_schedule"].Elem.(*schema.Resource).Schema
	// ConflictsWith is a path from the resource root. The generator emits the sibling name.
	schedule["cron"].ConflictsWith = []string{"alert_schedule.0.interval"}
	schedule["interval"].ConflictsWith = []string{"alert_schedule.0.cron"}
	s["enabled"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies if an alert should be 'started' (enabled) after creation or should remain 'suspended' (default).",
	}
	return s
}
