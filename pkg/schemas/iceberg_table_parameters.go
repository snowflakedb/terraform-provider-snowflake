package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowIcebergTableExternallyManagedParametersSchema contains the Snowflake parameters surfaced for iceberg tables.
var ShowIcebergTableExternallyManagedParametersSchema = map[string]*schema.Schema{
	"external_volume":            ParameterListSchema,
	"catalog":                    ParameterListSchema,
	"replace_invalid_characters": ParameterListSchema,
}

// ShowIcebergTableFromRestParametersSchema extends ShowIcebergTableParametersSchema with the
// additional parameters exposed by the Iceberg table from REST catalog resource.
var ShowIcebergTableFromRestParametersSchema = collections.MergeMaps(
	ShowIcebergTableExternallyManagedParametersSchema,
	map[string]*schema.Schema{
		"target_file_size":               ParameterListSchema,
		"storage_serialization_policy":   ParameterListSchema,
		"enable_iceberg_merge_on_read":   ParameterListSchema,
		"iceberg_merge_on_read_behavior": ParameterListSchema,
	},
)

func IcebergTableExternallyManagedParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return handleCommonIcebergTableExternallyManagedParameter(parameters, providerCtx)
}

// IcebergTableFromRestParametersToSchema maps the parameters surfaced by the Iceberg table from REST
// catalog resource, which includes additional parameters on top of the common ones.
func IcebergTableFromRestParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	result := handleCommonIcebergTableExternallyManagedParameter(parameters, providerCtx)
	for _, param := range parameters {
		parameterSchema := ParameterToSchemaReducedOutput(param, providerCtx)
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.IcebergTableParameterTargetFileSize),
			string(sdk.IcebergTableParameterStorageSerializationPolicy),
			string(sdk.IcebergTableParameterEnableIcebergMergeOnRead),
			string(sdk.IcebergTableParameterIcebergMergeOnReadBehavior):
			result[strings.ToLower(key)] = []map[string]any{parameterSchema}
		}
	}
	return result
}

// handleCommonIcebergTableExternallyManagedParameter maps the parameters common to all Iceberg tables into result.
func handleCommonIcebergTableExternallyManagedParameter(params []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	result := make(map[string]any)
	for _, param := range params {
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.IcebergTableParameterExternalVolume),
			string(sdk.IcebergTableParameterCatalog),
			string(sdk.IcebergTableParameterReplaceInvalidCharacters):
			result[strings.ToLower(key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		default:
		}
	}
	return result
}

var ShowIcebergTableSnowflakeManagedParametersSchema = map[string]*schema.Schema{
	"external_volume":                 ParameterListSchema,
	"catalog":                         ParameterListSchema,
	"target_file_size":                ParameterListSchema,
	"storage_serialization_policy":    ParameterListSchema,
	"catalog_sync":                    ParameterListSchema,
	"data_retention_time_in_days":     ParameterListSchema,
	"max_data_extension_time_in_days": ParameterListSchema,
	"enable_data_compaction":          ParameterListSchema,
	"enable_iceberg_merge_on_read":    ParameterListSchema,
}

// ShowIcebergTableAllTypesParametersSchema is the union of all iceberg table type parameter schemas.
// Used by the data source, which can return mixed table types in a single query.
var ShowIcebergTableAllTypesParametersSchema = collections.MergeMaps(
	ShowIcebergTableSnowflakeManagedParametersSchema,
	ShowIcebergTableFromRestParametersSchema,
	ShowIcebergTableExternallyManagedParametersSchema,
)

// IcebergTableAllTypesParametersToSchema maps parameters from any iceberg table type by
// combining the results of all type-specific mapping functions.
func IcebergTableAllTypesParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return collections.MergeMaps(
		IcebergTableSnowflakeManagedParametersToSchema(parameters, providerCtx),
		IcebergTableFromRestParametersToSchema(parameters, providerCtx),
		IcebergTableExternallyManagedParametersToSchema(parameters, providerCtx),
	)
}

func IcebergTableSnowflakeManagedParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	result := make(map[string]any)
	for _, param := range parameters {
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.IcebergTableParameterExternalVolume),
			string(sdk.IcebergTableParameterCatalog),
			string(sdk.IcebergTableParameterTargetFileSize),
			string(sdk.IcebergTableParameterStorageSerializationPolicy),
			string(sdk.IcebergTableParameterCatalogSync),
			string(sdk.IcebergTableParameterDataRetentionTimeInDays),
			string(sdk.IcebergTableParameterMaxDataExtensionTimeInDays),
			string(sdk.IcebergTableParameterEnableDataCompaction),
			string(sdk.IcebergTableParameterEnableIcebergMergeOnRead):
			result[strings.ToLower(key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return result
}
