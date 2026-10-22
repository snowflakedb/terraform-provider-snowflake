package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var alertIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the alert; must be unique for the database and schema in which the alert is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the alert.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the alert.",
		ForceNew:    true,
	},
}

func alertAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	schedule := s["alert_schedule"].Elem.(*schema.Resource).Schema
	// ConflictsWith is a path from the resource root. The generator emits the sibling name.
	schedule["cron"].ConflictsWith = []string{"alert_schedule.interval"}
	schedule["interval"].ConflictsWith = []string{"alert_schedule.cron"}
	s["enabled"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies if an alert should be 'started' (enabled) after creation or should remain 'suspended' (default).",
	}
	return s
}
