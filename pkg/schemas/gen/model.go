package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
)

type ShowResultSchemaModel struct {
	Name         string
	SdkType      string
	SchemaFields []SchemaField
}

func (m ShowResultSchemaModel) SomeFunc() {
}

func ModelFromStructDetails(sdkStruct genhelpers.StructDetails) ShowResultSchemaModel {
	name, _ := strings.CutPrefix(sdkStruct.Name, "sdk.")
	schemaFields := make([]SchemaField, len(sdkStruct.Fields))
	for idx, field := range sdkStruct.Fields {
		schemaFields[idx] = MapToSchemaField(field)
	}

	return ShowResultSchemaModel{
		Name:         name,
		SdkType:      sdkStruct.Name,
		SchemaFields: schemaFields,
	}
}
