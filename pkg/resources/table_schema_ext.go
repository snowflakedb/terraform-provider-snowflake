package resources

import (
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1031688]: move data manipulation logic to the SDK - SQL generation or builders part (e.g. different default types/identity)

// tableAttributesModificationExt restores schema chrome the generator does not emit.
func tableAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	column := s["column"].Elem.(*schema.Resource).Schema
	column["type"].ValidateDiagFunc = IsDataTypeValid
	column["type"].DiffSuppressFunc = DiffSuppressDataTypes
	column["nullable"].Default = true
	column["comment"].Default = ""
	column["masking_policy"].Default = ""
	column["collate"].Default = ""
	// A column without an explicit collation inherits the table's default_ddl_collation at the moment it is created,
	// and DESCRIBE keeps reporting it after the table default is changed or unset. An omitted collate is therefore
	// treated as "not managed", so it must not plan a change to the collation Snowflake assigned. An explicit
	// `collate = ""` is managed: it is a change that clears the collation, so it is not suppressed.
	column["collate"].DiffSuppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
		return new == "" && old != "" && isColumnCollateOmitted(k, d)
	}

	identity := column["identity"].Elem.(*schema.Resource).Schema
	identity["start_num"].Default = 1
	identity["step_num"].Default = 1

	s["primary_key"].Deprecated = "Use snowflake_table_constraint instead"

	// Identifier validation is not generated for nested-block fields (identifierSchemaType skips InBlock).
	for _, key := range []string{"row_access_policy", "aggregation_policy", "join_policy"} {
		s[key].Elem.(*schema.Resource).Schema["policy_name"].ValidateDiagFunc = IsValidIdentifier[sdk.SchemaObjectIdentifier]()
	}

	return s
}

// isColumnCollateOmitted reports whether the collate of the column under the given schema key (column.<index>.collate)
// is absent from the configuration. When the configuration cannot be read, collate is treated as omitted,
// which keeps the previous behavior of not planning a change to the collation Snowflake assigned.
func isColumnCollateOmitted(key string, d *schema.ResourceData) bool {
	parts := strings.Split(key, ".")
	if len(parts) != 3 || parts[0] != "column" || parts[2] != "collate" {
		return true
	}
	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return true
	}
	value, diags := d.GetRawConfigAt(cty.GetAttrPath("column").IndexInt(index).GetAttr("collate"))
	if diags.HasError() {
		return true
	}
	return value.IsNull()
}
