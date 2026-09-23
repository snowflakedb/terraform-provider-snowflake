package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (computePoolToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backup_instance_families": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (computePoolToSchemaMapper) additionalToSchema(src *sdk.ComputePool, dst map[string]any) {
	dst["backup_instance_families"] = src.BackupInstanceFamilies
}
