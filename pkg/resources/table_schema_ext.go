package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1348114]: old implementation was quoting every column, SDK is not quoting them, therefore they are quoted here: decide if we quote columns or not
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

	identity := column["identity"].Elem.(*schema.Resource).Schema
	identity["start_num"].Default = 1
	identity["step_num"].Default = 1

	// TODO [SNOW-1348114]: consider adding fully_qualified_name for columns.
	column["schema_evolution_record"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Record of schema evolution.",
	}

	s["primary_key"].Deprecated = "Use snowflake_table_constraint instead"
	s["change_tracking"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("change_tracking", func(x any) any {
		return x.(bool)
	})
	return s
}
