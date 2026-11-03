package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func pipeAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["copy_statement"].DiffSuppressFunc = pipeCopyStatementDiffSuppress
	s["auto_ingest"].Default = false
	s["notification_channel"].Deprecated = "Use the `show_output` field instead (`show_output.0.notification_channel`). It will be removed in the next major version."
	s["owner"].Deprecated = "Use the `show_output` field instead (`show_output.0.owner`). It will be removed in the next major version."
	return s
}
