package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (streamToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// TODO [SNOW-3113144]: should it be list?
		"base_tables": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		},
	}
}

func (streamToSchemaMapper) additionalToSchema(src *sdk.Stream, dst map[string]any) {
	// Identifier slice: map each table to FullyQualifiedName; skip when unset.
	if src.BaseTables != nil {
		dst["base_tables"] = collections.Map(src.BaseTables, sdk.SchemaObjectIdentifier.FullyQualifiedName)
	}
}
