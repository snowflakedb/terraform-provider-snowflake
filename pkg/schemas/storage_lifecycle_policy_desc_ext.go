package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (storageLifecyclePolicyDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
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

func (storageLifecyclePolicyDetailsToSchemaMapper) additionalToSchema(src *sdk.StorageLifecyclePolicyDetails, dst map[string]any) {
	signatureElem := make([]map[string]any, len(src.Signature))
	for i, v := range src.Signature {
		signatureElem[i] = map[string]any{
			"name": v.Name,
			"type": v.Type.ToSql(),
		}
	}
	dst["signature"] = signatureElem
}
