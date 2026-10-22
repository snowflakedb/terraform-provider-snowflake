package gen

import (
	"maps"
	"os"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type PluginFrameworkProviderModel struct {
	ModelFields   []ProviderModelField
	SchemaEntries []ProviderSchemaEntry
	PackageName   string

	*genhelpers.PreambleModel
}

func ModelFromSdkV2Schema(sdkV2ProviderSchema SdkV2ProviderSchema, preamble *genhelpers.PreambleModel) PluginFrameworkProviderModel {
	orderedAttributeNames := slices.Collect(maps.Keys(sdkV2ProviderSchema))
	slices.Sort(orderedAttributeNames)

	modelFields := make([]ProviderModelField, 0, len(orderedAttributeNames))
	schemaEntries := make([]ProviderSchemaEntry, 0, len(orderedAttributeNames))
	for _, key := range orderedAttributeNames {
		sdkV2SchemaAttribute := sdkV2ProviderSchema[key]
		modelField, err := MapToPluginFrameworkProviderModelField(key, sdkV2SchemaAttribute)
		if err != nil {
			panic(err)
		}
		schemaEntry, err := MapToPluginFrameworkProviderSchema(key, sdkV2SchemaAttribute)
		if err != nil {
			panic(err)
		}
		modelFields = append(modelFields, *modelField)
		schemaEntries = append(schemaEntries, *schemaEntry)
	}
	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return PluginFrameworkProviderModel{
		ModelFields:   modelFields,
		SchemaEntries: schemaEntries,
		PackageName:   packageWithGenerateDirective,
		PreambleModel: preamble,
	}
}
