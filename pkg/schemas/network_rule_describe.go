package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeNetworkRuleSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"mode": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"value_list": {
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
}

func NetworkRuleDetailsToSchema(details *sdk.NetworkRuleDetails) map[string]any {
	return map[string]any{
		"created_on":    details.CreatedOn.String(),
		"name":          details.Name,
		"database_name": details.DatabaseName,
		"schema_name":   details.SchemaName,
		"owner":         details.Owner,
		"comment":       details.Comment,
		"type":          string(details.Type),
		"mode":          string(details.Mode),
		"value_list":    details.ValueList,
	}
}
