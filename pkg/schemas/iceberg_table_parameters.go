package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var icebergTableExternallyManagedParameters = []parameterdefs.ParameterDef{
	defs.ExternalVolume,
	defs.Catalog,
	defs.ReplaceInvalidCharacters,
}

var icebergTableFromRestParameters = append(slices.Clone(icebergTableExternallyManagedParameters),
	defs.TargetFileSize,
	defs.StorageSerializationPolicy,
	defs.EnableIcebergMergeOnRead,
	defs.IcebergMergeOnReadBehavior,
)

var icebergTableSnowflakeManagedParameters = []parameterdefs.ParameterDef{
	defs.ExternalVolume,
	defs.Catalog,
	defs.TargetFileSize,
	defs.StorageSerializationPolicy,
	defs.CatalogSync,
	defs.DataRetentionTimeInDays,
	defs.MaxDataExtensionTimeInDays,
	defs.EnableDataCompaction,
	defs.EnableIcebergMergeOnRead,
}

var (
	ShowIcebergTableExternallyManagedParametersSchema = make(map[string]*schema.Schema)
	ShowIcebergTableFromRestParametersSchema          = make(map[string]*schema.Schema)
	ShowIcebergTableSnowflakeManagedParametersSchema  = make(map[string]*schema.Schema)

	ShowIcebergTableAllTypesParametersSchema = make(map[string]*schema.Schema)
)

func init() {
	for _, def := range icebergTableExternallyManagedParameters {
		ShowIcebergTableExternallyManagedParametersSchema[def.FieldName()] = ParameterListSchema
	}
	for _, def := range icebergTableFromRestParameters {
		ShowIcebergTableFromRestParametersSchema[def.FieldName()] = ParameterListSchema
	}
	for _, def := range icebergTableSnowflakeManagedParameters {
		ShowIcebergTableSnowflakeManagedParametersSchema[def.FieldName()] = ParameterListSchema
	}
	ShowIcebergTableAllTypesParametersSchema = collections.MergeMaps(
		ShowIcebergTableSnowflakeManagedParametersSchema,
		ShowIcebergTableFromRestParametersSchema,
		ShowIcebergTableExternallyManagedParametersSchema,
	)
}

func icebergTableParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context, defsForVariant []parameterdefs.ParameterDef) map[string]any {
	result := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(defsForVariant, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			result[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return result
}

func IcebergTableExternallyManagedParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return icebergTableParametersToSchema(parameters, providerCtx, icebergTableExternallyManagedParameters)
}

func IcebergTableFromRestParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return icebergTableParametersToSchema(parameters, providerCtx, icebergTableFromRestParameters)
}

func IcebergTableSnowflakeManagedParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return icebergTableParametersToSchema(parameters, providerCtx, icebergTableSnowflakeManagedParameters)
}

func IcebergTableAllTypesParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return collections.MergeMaps(
		IcebergTableSnowflakeManagedParametersToSchema(parameters, providerCtx),
		IcebergTableFromRestParametersToSchema(parameters, providerCtx),
		IcebergTableExternallyManagedParametersToSchema(parameters, providerCtx),
	)
}
