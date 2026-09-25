package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (fileFormatParquetToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"null_if": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (fileFormatParquetToSchemaMapper) additionalToSchema(src *sdk.FileFormatParquet, dst map[string]any) {
	// TODO: drop this []any coerce when handleExternalChangesToObjectDeepEqual compares slices in a Terraform-aware way.
	dst["null_if"] = collections.Map(src.NullIf, func(v string) any { return v })
}
