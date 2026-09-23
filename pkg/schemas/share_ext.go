package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (shareToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (shareToSchemaMapper) additionalToSchema(src *sdk.Share, dst map[string]any) {
	dst["name"] = src.ID().FullyQualifiedName()
}
