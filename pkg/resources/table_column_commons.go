package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// basicColumn is a minimal column definition (name + type) shared by table-like resources.
// It intentionally omits constraints, policies, defaults, etc. so it can be reused by resources
// that only need the essential column shape.
type basicColumn struct {
	Name     string
	DataType datatypes.DataType
}

func basicColumnSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		ForceNew:    true,
		MinItems:    1,
		Description: "Definitions of the columns to create in the table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Column name.",
				},
				"type": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         true,
					Description:      "Column type, e.g. VARIANT. For a full list of column types, see [Summary of Data Types](https://docs.snowflake.com/en/sql-reference/intro-summary-data-types).",
					ValidateDiagFunc: IsDataTypeValid,
					DiffSuppressFunc: DiffSuppressDataTypes,
				},
			},
		},
	}
}

func parseBasicColumns(raw []any) ([]basicColumn, error) {
	return collections.MapErr(raw, func(r any) (basicColumn, error) {
		c := r.(map[string]any)
		name := c["name"].(string)
		dataType, err := datatypes.ParseDataType(c["type"].(string))
		if err != nil {
			return basicColumn{}, fmt.Errorf("parsing data type of column %q: %w", name, err)
		}
		return basicColumn{Name: name, DataType: dataType}, nil
	})
}
