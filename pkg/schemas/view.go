package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ViewDescribeSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_nullable": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"default": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_primary": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_unique": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"check": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"expression": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"policy_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"privacy_domain": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ViewDescriptionToSchema(description []sdk.ViewDetails) []map[string]any {
	result := make([]map[string]any, len(description))
	for i, row := range description {
		result[i] = map[string]any{
			"name":           row.Name,
			"type":           row.Type,
			"kind":           row.Kind,
			"is_nullable":    row.IsNullable,
			"default":        row.Default,
			"is_primary":     row.IsPrimary,
			"is_unique":      row.IsUnique,
			"check":          row.Check,
			"expression":     row.Expression,
			"comment":        row.Comment,
			"policy_name":    row.PolicyName,
			"privacy_domain": row.PrivacyDomain,
		}
	}
	return result
}
