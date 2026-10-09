package resources

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["with_managed_access"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("options", func(x any) any {
		return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "MANAGED ACCESS")
	})
	s["is_transient"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("options", func(x any) any {
		return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "TRANSIENT")
	})
	return s
}
