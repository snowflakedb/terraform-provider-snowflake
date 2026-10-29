package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func pipeAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["copy_statement"].DiffSuppressFunc = pipeCopyStatementDiffSuppress
	s["auto_ingest"].Default = false
	s["notification_channel"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Amazon Resource Name of the Amazon SQS queue for the stage named in the DEFINITION column.",
	}
	s["owner"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the pipe.",
	}
	return s
}
