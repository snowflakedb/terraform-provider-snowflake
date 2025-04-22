package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var RowAccessPolicyDescribeSchema = map[string]*schema.Schema{
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

func RowAccessPolicyDescriptionToSchema(description sdk.RowAccessPolicyDescription) map[string]any {
	signatureElem := make([]map[string]any, len(description.Signature))
	for i, v := range description.Signature {
		signatureElem[i] = map[string]any{
			"name": v.Name,
			"type": v.Type.ToSql(),
		}
	}
	return map[string]any{
		"name":        description.Name,
		"signature":   signatureElem,
		"return_type": description.ReturnType,
		"body":        description.Body,
	}
}

// TODO [this PR]: address this func (ToSql vs ToSqlWithoutUnknowns for type)
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
