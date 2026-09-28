package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// additionalSchema contributes the FK-only scalars of ShowHybridTableConstraintSchema.
// That schema is for a single row of the hybrid table's merged constraints view
// (PRIMARY KEY / UNIQUE / FOREIGN KEY) built by GetConstraints from SHOW PRIMARY KEYS,
// SHOW UNIQUE KEYS, and SHOW IMPORTED KEYS.
func (hybridTableConstraintToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"referenced_table": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"delete_rule": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"update_rule": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (hybridTableConstraintToSchemaMapper) additionalToSchema(src *sdk.HybridTableConstraint, dst map[string]any) {
	// referenced_table/delete_rule/update_rule are FK-only,
	// so they are left unset (empty) for PRIMARY KEY / UNIQUE constraints.
	if src.Kind == sdk.ColumnConstraintTypeForeignKey {
		dst["referenced_table"] = src.ReferencedTable.FullyQualifiedName()
		dst["delete_rule"] = src.DeleteRule
		dst["update_rule"] = src.UpdateRule
	}
}
