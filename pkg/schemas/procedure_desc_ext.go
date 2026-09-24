package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (procedureDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"return_data_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (procedureDetailsToSchemaMapper) additionalToSchema(src *sdk.ProcedureDetails, dst map[string]any) {
	// TODO [next PRs]: drop return_data_type from this ext once the generator maps datatypes.DataType via ToSql().
	if src.ReturnDataType != nil {
		dst["return_data_type"] = src.ReturnDataType.ToSql()
	}
}
