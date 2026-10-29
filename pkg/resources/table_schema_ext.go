package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO [SNOW-1348114]: old implementation was quoting every column, SDK is not quoting them, therefore they are quoted here: decide if we quote columns or not
// TODO [SNOW-1031688]: move data manipulation logic to the SDK - SQL generation or builders part (e.g. different default types/identity)

// tableAttributesModificationExt restores schema chrome the generator does not emit.
// data_retention_time_in_days Min and Max are -1 and 90. Those emit ValidateDiagFunc;
// this resource keeps ValidateFunc.
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

	s["owner"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the table.",
	}
	s["primary_key"].Deprecated = "Use snowflake_table_constraint instead"

	retention := s["data_retention_time_in_days"]
	retention.ValidateDiagFunc = nil
	retention.ValidateFunc = validation.IntBetween(-1, 90)
	s["change_tracking"].Default = false
	s["tag"] = tagReferenceSchema
	return s
}
