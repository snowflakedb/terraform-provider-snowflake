package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeStorageLifecyclePolicyDetailsSchema represents output of DESCRIBE query for the single StorageLifecyclePolicy.
var DescribeStorageLifecyclePolicyDetailsSchema = map[string]*schema.Schema{
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
	"archive_for_days": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"archive_tier": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func StorageLifecyclePolicyDetailsToSchema(details sdk.StorageLifecyclePolicyDetails) map[string]any {
	signatureElem := make([]map[string]any, len(details.Signature))
	for i, v := range details.Signature {
		signatureElem[i] = map[string]any{
			"name": v.Name,
			"type": v.Type.ToSql(),
		}
	}
	result := map[string]any{
		"name":         details.Name,
		"signature":    signatureElem,
		"return_type":  details.ReturnType.ToSql(),
		"body":         details.Body,
		"archive_tier": details.ArchiveTier,
	}
	if details.ArchiveForDays != nil {
		result["archive_for_days"] = *details.ArchiveForDays
	}
	return result
}
