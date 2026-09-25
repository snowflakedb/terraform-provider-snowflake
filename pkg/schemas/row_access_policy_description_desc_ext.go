package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (rowAccessPolicyDescriptionToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}

func (rowAccessPolicyDescriptionToSchemaMapper) additionalToSchema(src *sdk.RowAccessPolicyDescription, dst map[string]any) {
	dst["signature"] = RowAccessPolicyArgumentsToSchema(src.Signature)
}

func RowAccessPolicyArgumentsToSchema(args []sdk.TableColumnSignature) []map[string]any {
	schema := make([]map[string]any, len(args))
	for i, v := range args {
		schema[i] = map[string]any{
			"name": v.Name,
			"type": v.Type.ToSql(),
		}
	}
	return schema
}
