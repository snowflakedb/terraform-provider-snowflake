package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (authenticationPolicyToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"target_scopes": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (authenticationPolicyToSchemaMapper) additionalToSchema(src *sdk.AuthenticationPolicy, dst map[string]any) {
	dst["target_scopes"] = collections.Map(src.TargetScopes, func(v sdk.AuthenticationPolicyTargetScope) string {
		return string(v)
	})
}
