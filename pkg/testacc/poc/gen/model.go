package gen

import (
	"os"
	"slices"
)

type PluginFrameworkProviderModel struct {
	ModelFields   []ProviderModelField
	SchemaEntries []ProviderSchemaEntry
	PackageName   string
}

func (m PluginFrameworkProviderModel) SomeFunc() {
}

func ModelFromSdkV2Schema(sdkV2ProviderSchema SdkV2ProviderSchema) PluginFrameworkProviderModel {
	orderedAttributeNames := make([]string, 0)
	for key := range sdkV2ProviderSchema {
		orderedAttributeNames = append(orderedAttributeNames, key)
	}
	slices.Sort(orderedAttributeNames)

	modelFields := make([]ProviderModelField, 0)
	schemaEntries := make([]ProviderSchemaEntry, 0)
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
	}
}
