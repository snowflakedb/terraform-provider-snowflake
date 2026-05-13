package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type ShowResultSchemaModel struct {
	Name         string
	SdkType      string
	SchemaFields []SchemaField
	// Prefix overrides the default "Show" prefix used in variable/function names.
	// When empty, the default "Show" prefix is used (e.g. ShowPostgresInstanceSchema).
	// Set to "Describe" for DESCRIBE output schemas (e.g. DescribePostgresInstanceSchema).
	Prefix string
	// MapperName overrides the name used in the ToSchema mapper function.
	// When empty, .Name is used (e.g. PostgresInstanceToSchema).
	// Set differently when the schema variable name should differ from the mapper function name.
	MapperName string

	*genhelpers.PreambleModel
}

// SchemaPrefix returns the prefix for schema variable names. Defaults to "Show" if not set.
func (m ShowResultSchemaModel) SchemaPrefix() string {
	if m.Prefix != "" {
		return m.Prefix
	}
	return "Show"
}

// EffectiveMapperName returns the name used for the ToSchema function. Defaults to .Name.
func (m ShowResultSchemaModel) EffectiveMapperName() string {
	if m.MapperName != "" {
		return m.MapperName
	}
	return m.Name
}

func ModelFromStructDetails(sdkStruct genhelpers.StructDetails, preamble *genhelpers.PreambleModel) ShowResultSchemaModel {
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

// ModelFromStructDetailsWithPrefix creates a model with a custom schema prefix (e.g. "Describe").
// It also strips common suffixes like "Details" from the struct name to produce cleaner schema
// variable names (e.g. PostgresInstanceDetails → DescribePostgresInstanceSchema), while keeping
// the full name for the ToSchema mapper function (e.g. PostgresInstanceDetailsToSchema).
func ModelFromStructDetailsWithPrefix(prefix string) func(genhelpers.StructDetails, *genhelpers.PreambleModel) ShowResultSchemaModel {
	return func(sdkStruct genhelpers.StructDetails, preamble *genhelpers.PreambleModel) ShowResultSchemaModel {
		model := ModelFromStructDetails(sdkStruct, preamble)
		model.Prefix = prefix
		// Keep original name for the mapper function (e.g. PostgresInstanceDetailsToSchema).
		model.MapperName = model.Name
		// Strip "Details" suffix from name for cleaner schema variable naming.
		model.Name = strings.TrimSuffix(model.Name, "Details")
		return model
	}
}
