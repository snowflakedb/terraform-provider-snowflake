package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type ShowResultSchemaModel struct {
	Name         string
	SdkType      string
	SchemaFields []SchemaField

	genhelpers.PreambleModel
}

func ModelFromStructDetails(sdkStruct genhelpers.StructDetails, preamble genhelpers.PreambleModel) ShowResultSchemaModel {
	name, _ := strings.CutPrefix(sdkStruct.Name, "sdk.")
	schemaFields := make([]SchemaField, len(sdkStruct.Fields))
	for idx, field := range sdkStruct.Fields {
		schemaFields[idx] = MapToSchemaField(field)
	}

	return ShowResultSchemaModel{
		Name:          name,
		SdkType:       sdkStruct.Name,
		SchemaFields:  schemaFields,
		PreambleModel: preamble,
	}
}
