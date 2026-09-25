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
	skip := keyedSet(sdkStruct.SkipFields)
	manual := keyedSet(sdkStruct.ManualFields)
	if overlap := intersectingKeys(skip, manual); len(overlap) > 0 {
		panic(fmt.Sprintf("SkipFields and ManualFields for %s overlap: %s", sdkStruct.Name, strings.Join(overlap, ", ")))
	}

	schemaFields := make([]SchemaField, 0, len(sdkStruct.Fields))
	for _, field := range sdkStruct.Fields {
		schemaField := MapToSchemaField(field)
		if _, ok := skip[schemaField.Name]; ok {
			delete(skip, schemaField.Name)
			schemaField.Skipped = true
		}
		if _, ok := manual[schemaField.Name]; ok {
			delete(manual, schemaField.Name)
			schemaField.Manual = true
		}
		schemaFields = append(schemaFields, schemaField)
	}
	if unmatched := sortedKeys(skip); len(unmatched) > 0 {
		panic(fmt.Sprintf("SkipFields for %s contain unknown schema keys: %s", sdkStruct.Name, strings.Join(unmatched, ", ")))
	}
	if unmatched := sortedKeys(manual); len(unmatched) > 0 {
		panic(fmt.Sprintf("ManualFields for %s contain unknown schema keys: %s", sdkStruct.Name, strings.Join(unmatched, ", ")))
	}

	return ShowResultSchemaModel{
		Name:              name,
		SdkType:           sdkStruct.Name,
		IsDescribe:        sdkStruct.IsDescribe,
		UsedAsListEntry:   sdkStruct.UsedAsListEntry,
		AdditionalMapping: len(sdkStruct.ManualFields) > 0,
		SchemaFields:      schemaFields,
		PreambleModel:     preamble,
	}
}

// TODO [next PRs]: considering moving these helper funcs to collections utils with tests
func keyedSet(keys []string) map[string]struct{} {
	set := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		set[key] = struct{}{}
	}
	return set
}

func sortedKeys(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func intersectingKeys(a, b map[string]struct{}) []string {
	var overlap []string
	for key := range a {
		if _, ok := b[key]; ok {
			overlap = append(overlap, key)
		}
	}
	slices.Sort(overlap)
	return overlap
}
