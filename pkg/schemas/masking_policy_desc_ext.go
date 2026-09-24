package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (maskingPolicyDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
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
		"return_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (maskingPolicyDetailsToSchemaMapper) additionalToSchema(src *sdk.MaskingPolicyDetails, dst map[string]any) {
	dst["signature"] = MaskingPolicyArgumentsToSchema(src.Signature)
	// TODO [next PRs]: drop return_type from this ext once the generator maps datatypes.DataType via ToSql().
	dst["return_type"] = src.ReturnType.ToSql()
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
