package gen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type ShowResultSchemaModel struct {
	Name              string
	SdkType           string
	IsDescribe        bool
	UsedAsListEntry   bool
	AdditionalMapping bool
	SchemaFields      []SchemaField

	*genhelpers.PreambleModel
}

func (m ShowResultSchemaModel) Filename() string {
	snake := genhelpers.ToSnakeCase(m.Name)
	if m.IsDescribe {
		return strings.TrimSuffix(snake, "_details") + "_desc_gen.go"
	}
	return snake + "_gen.go"
}

func ModelFromStructDetails(sdkStruct ShowResultSchemaDetails, preamble *genhelpers.PreambleModel) ShowResultSchemaModel {
	if sdkStruct.IsDescribe && sdkStruct.UsedAsListEntry {
		panic(fmt.Sprintf("%s: IsDescribe and UsedAsListEntry are mutually exclusive", sdkStruct.Name))
	}
	name, _ := strings.CutPrefix(sdkStruct.Name, "sdk.")
	skip := make(map[string]struct{}, len(sdkStruct.SkipFields))
	for _, fieldName := range sdkStruct.SkipFields {
		skip[fieldName] = struct{}{}
	}

	schemaFields := make([]SchemaField, 0, len(sdkStruct.Fields))
	for _, field := range sdkStruct.Fields {
		schemaField := MapToSchemaField(field)
		if _, ok := skip[schemaField.Name]; ok {
			delete(skip, schemaField.Name)
			continue
		}
		schemaFields = append(schemaFields, schemaField)
	}
	if len(skip) > 0 {
		unmatched := make([]string, 0, len(skip))
		for fieldName := range skip {
			unmatched = append(unmatched, fieldName)
		}
		slices.Sort(unmatched)
		panic(fmt.Sprintf("SkipFields for %s contain unknown schema keys: %s", sdkStruct.Name, strings.Join(unmatched, ", ")))
	}

	return ShowResultSchemaModel{
		Name:              name,
		SdkType:           sdkStruct.Name,
		IsDescribe:        sdkStruct.IsDescribe,
		UsedAsListEntry:   sdkStruct.UsedAsListEntry,
		AdditionalMapping: sdkStruct.AdditionalMapping,
		SchemaFields:      schemaFields,
		PreambleModel:     preamble,
	}
}
