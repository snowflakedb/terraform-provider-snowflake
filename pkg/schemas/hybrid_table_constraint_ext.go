package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// additionalSchema contributes the list and FK fields of ShowHybridTableConstraintSchema.
// That schema is for a single row of the hybrid table's merged constraints view
// (PRIMARY KEY / UNIQUE / FOREIGN KEY) built by GetConstraints from SHOW PRIMARY KEYS,
// SHOW UNIQUE KEYS, and SHOW IMPORTED KEYS.
func (hybridTableConstraintToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"columns": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"referenced_table": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"referenced_columns": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
	dst["columns"] = src.Columns
	// referenced_table/referenced_columns/delete_rule/update_rule are FK-only,
	// so they are left unset (empty) for PRIMARY KEY / UNIQUE constraints.
	if src.Kind == sdk.ColumnConstraintTypeForeignKey {
		dst["referenced_table"] = src.ReferencedTable.FullyQualifiedName()
		dst["referenced_columns"] = src.ReferencedColumns
		dst["delete_rule"] = src.DeleteRule
		dst["update_rule"] = src.UpdateRule
	}
}
