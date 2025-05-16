package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeMaskingPolicySchema represents output of DESCRIBE query for the single masking policy.
var DescribeMaskingPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"signature": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
		Computed: true,
	},
	"return_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"body": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func MaskingPolicyDescriptionToSchema(details sdk.MaskingPolicyDetails) map[string]any {
	return map[string]any{
		"name":        details.Name,
		"signature":   MaskingPolicyArgumentsToSchema(details.Signature),
		"return_type": details.ReturnType.ToSql(),
		"body":        details.Body,
	}
}

func MaskingPolicyArgumentsToSchema(args []sdk.TableColumnSignature) []map[string]any {
	schema := make([]map[string]any, len(args))
	for i, v := range args {
		schema[i] = map[string]any{
			"name": v.Name,
			"type": v.Type.ToSql(),
		}
	}
	return schema
}
