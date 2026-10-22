package resources

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func computePoolAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["auto_resume"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume")
	s["initially_suspended"].DiffSuppressFunc = IgnoreAfterCreation
	s["auto_suspend_secs"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend_secs")
	return s
}
