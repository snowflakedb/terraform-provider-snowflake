package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// external_volume and catalog are ForceNew: Snowflake does not support altering them after creation
// for any Iceberg table type.
func icebergTableCommonParametersSchema() map[string]*schema.Schema {
	parametersSchema := make(map[string]*schema.Schema)
	for _, p := range []parameterdefs.ParameterDef{defs.ExternalVolume, defs.Catalog} {
		s := parameterSchema(p)
		s.ForceNew = true
		parametersSchema[p.FieldName()] = s
	}
	return parametersSchema
}

func icebergTableExternalManagedParametersSchema() map[string]*schema.Schema {
	parametersSchema := icebergTableCommonParametersSchema()
	parametersSchema[defs.ReplaceInvalidCharacters.FieldName()] = parameterSchema(defs.ReplaceInvalidCharacters)
	return parametersSchema
}

func icebergTableFromRestParametersSchema() map[string]*schema.Schema {
	parametersSchema := icebergTableExternalManagedParametersSchema()
	parametersSchema[defs.TargetFileSize.FieldName()] = parameterSchema(defs.TargetFileSize)
	parametersSchema[defs.EnableIcebergMergeOnRead.FieldName()] = parameterSchema(defs.EnableIcebergMergeOnRead)

	storageSerializationPolicy := parameterSchema(defs.StorageSerializationPolicy)
	storageSerializationPolicy.ForceNew = true
	parametersSchema[defs.StorageSerializationPolicy.FieldName()] = storageSerializationPolicy

	// TODO (next PRs): this is now available in ALTER ... SET - add to sdk and make it non-force-new here.
	icebergMergeOnReadBehavior := parameterSchema(defs.IcebergMergeOnReadBehavior)
	icebergMergeOnReadBehavior.ForceNew = true
	parametersSchema[defs.IcebergMergeOnReadBehavior.FieldName()] = icebergMergeOnReadBehavior

	return parametersSchema
}

func icebergTableSnowflakeManagedParametersSchema() map[string]*schema.Schema {
	parametersSchema := icebergTableCommonParametersSchema()

	for _, p := range []parameterdefs.ParameterDef{
		defs.TargetFileSize,
		defs.CatalogSync,
		defs.DataRetentionTimeInDays,
		defs.MaxDataExtensionTimeInDays,
		defs.EnableDataCompaction,
		defs.EnableIcebergMergeOnRead,
	} {
		parametersSchema[p.FieldName()] = parameterSchema(p)
	}

	// storage_serialization_policy is only accepted at CREATE time (there is no ALTER ... SET for it), so it is ForceNew.
	storageSerializationPolicy := parameterSchema(defs.StorageSerializationPolicy)
	storageSerializationPolicy.ForceNew = true
	parametersSchema[defs.StorageSerializationPolicy.FieldName()] = storageSerializationPolicy

	return parametersSchema
}

func icebergTableParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.IcebergTableParametersDetails, error) {
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.IcebergTables.ShowParametersDetails(ctx, id)
}

func icebergTableExternalManagedParameterDiffFunctions(parameters *sdk.IcebergTableParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IdentifierTypedParameterValueComputedIf(defs.ExternalVolume.FieldName(), parameters.ExternalVolume, sdk.ParameterTypeTable),
		IdentifierTypedParameterValueComputedIf(defs.Catalog.FieldName(), parameters.Catalog, sdk.ParameterTypeTable),
		BoolTypedParameterValueComputedIf(defs.ReplaceInvalidCharacters.FieldName(), parameters.ReplaceInvalidCharacters, sdk.ParameterTypeTable),
	}
}

var icebergTableExternalManagedParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	icebergTableParametersProvider,
	icebergTableExternalManagedParameterDiffFunctions,
)

var icebergTableFromRestParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	icebergTableParametersProvider,
	icebergTableFromRestParameterDiffFunctions,
)

func icebergTableFromRestParameterDiffFunctions(parameters *sdk.IcebergTableParametersDetails) []schema.CustomizeDiffFunc {
	return append(
		icebergTableExternalManagedParameterDiffFunctions(parameters),
		StringTypedParameterValueComputedIf(defs.TargetFileSize.FieldName(), parameters.TargetFileSize, sdk.ParameterTypeTable),
		BoolTypedParameterValueComputedIf(defs.EnableIcebergMergeOnRead.FieldName(), parameters.EnableIcebergMergeOnRead, sdk.ParameterTypeTable),
	)
}

func handleIcebergTableExternalManagedParametersCreate(d *schema.ResourceData, externalVolume, catalog **sdk.AccountObjectIdentifier, replaceInvalidCharacters **bool) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreateWithMapping(d, defs.ExternalVolume.FieldName(), externalVolume, sdk.ParseAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, defs.Catalog.FieldName(), catalog, sdk.ParseAccountObjectIdentifier),
		handleParameterCreate(d, defs.ReplaceInvalidCharacters.FieldName(), replaceInvalidCharacters),
	)
}

// replace_invalid_characters is the only alterable parameter here: external_volume and catalog are
// both ForceNew.
func handleIcebergTableExternalManagedParametersUpdate(set *sdk.IcebergTableSetPropertiesRequest, unset *sdk.IcebergTableUnsetPropertiesRequest, d *schema.ResourceData) diag.Diagnostics {
	return handleParameterUpdate(d, defs.ReplaceInvalidCharacters.FieldName(), &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters)
}

// rawParameterValue returns the as-reported SHOW PARAMETERS value for key, or "" if absent.
func rawParameterValue(parameters []*sdk.Parameter, key string) string {
	for _, p := range parameters {
		if strings.EqualFold(p.Key, key) {
			return p.Value
		}
	}
	return ""
}

// external_volume/catalog are read from the raw SHOW PARAMETERS value rather than the typed
// AccountObjectIdentifier: Snowflake reports these inconsistently quoted depending on whether the
// value is inherited or explicitly set, and the typed identifier can only render one fixed form.
func handleIcebergTableExternalManagedParameterRead(d *schema.ResourceData, parameters []*sdk.Parameter, parameterDetails *sdk.IcebergTableParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.ExternalVolume.FieldName(), rawParameterValue(parameters, defs.ExternalVolume.SqlName)),
		setResourceData(d, defs.Catalog.FieldName(), rawParameterValue(parameters, defs.Catalog.SqlName)),
		setResourceData(d, defs.ReplaceInvalidCharacters.FieldName(), parameterDetails.ReplaceInvalidCharacters.Value),
	)
}

func handleIcebergTableFromRestParametersCreate(d *schema.ResourceData, req *sdk.CreateFromIcebergRestIcebergTableRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreateWithMapping(d, defs.TargetFileSize.FieldName(), &req.TargetFileSize, stringToStringEnumProvider(sdk.ToIcebergTableTargetFileSize)),
		handleParameterCreateWithMapping(d, defs.StorageSerializationPolicy.FieldName(), &req.StorageSerializationPolicy, stringToStringEnumProvider(sdk.ToStorageSerializationPolicy)),
		handleParameterCreateWithMapping(d, defs.IcebergMergeOnReadBehavior.FieldName(), &req.IcebergMergeOnReadBehavior, stringToStringEnumProvider(sdk.ToIcebergTableIcebergMergeOnReadBehavior)),
		handleParameterCreate(d, defs.EnableIcebergMergeOnRead.FieldName(), &req.EnableIcebergMergeOnRead),
	)
}

// storage_serialization_policy and iceberg_merge_on_read_behavior are omitted: both are create-only
// (ForceNew) and cannot be altered.
//
// NOTE(SNOW-3735539): altering replace_invalid_characters and target_file_size together with comment
// does not cause any changes, so changes to these parameters must be applied via a separate ALTER call
// from the comment update - see UpdateIcebergTableFromRest.
func handleIcebergTableFromRestParametersUpdate(d *schema.ResourceData, set *sdk.IcebergTableSetPropertiesRequest, unset *sdk.IcebergTableUnsetPropertiesRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdateWithMapping(d, defs.TargetFileSize.FieldName(), &set.TargetFileSize, &unset.TargetFileSize, stringToStringEnumProvider(sdk.ToIcebergTableTargetFileSize)),
		handleParameterUpdate(d, defs.ReplaceInvalidCharacters.FieldName(), &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleParameterUpdate(d, defs.EnableIcebergMergeOnRead.FieldName(), &set.EnableIcebergMergeOnRead, &unset.EnableIcebergMergeOnRead),
	)
}

func handleIcebergTableFromRestParameterRead(d *schema.ResourceData, parameters []*sdk.Parameter, parameterDetails *sdk.IcebergTableParametersDetails) diag.Diagnostics {
	return JoinDiags(
		handleIcebergTableExternalManagedParameterRead(d, parameters, parameterDetails),
		setResourceData(d, defs.TargetFileSize.FieldName(), parameterDetails.TargetFileSize.Value),
		setResourceData(d, defs.StorageSerializationPolicy.FieldName(), parameterDetails.StorageSerializationPolicy.Value),
		setResourceData(d, defs.EnableIcebergMergeOnRead.FieldName(), parameterDetails.EnableIcebergMergeOnRead.Value),
		setResourceData(d, defs.IcebergMergeOnReadBehavior.FieldName(), parameterDetails.IcebergMergeOnReadBehavior.Value),
	)
}

// storage_serialization_policy is intentionally omitted: it is ForceNew and create-only, and Snowflake
// always reports it at the TABLE parameter level once the table exists (there is no "inherited" state
// for it to drift from). Running it through ParameterValueComputedIf would make the
// `parameter.Level == objectParameterLevel` branch always true, marking it computed - and being
// ForceNew - forcing a spurious replace on every plan.
func icebergTableSnowflakeManagedParameterDiffFunctions(parameters *sdk.IcebergTableParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IdentifierTypedParameterValueComputedIf(defs.ExternalVolume.FieldName(), parameters.ExternalVolume, sdk.ParameterTypeTable),
		IdentifierTypedParameterValueComputedIf(defs.Catalog.FieldName(), parameters.Catalog, sdk.ParameterTypeTable),
		StringTypedParameterValueComputedIf(defs.TargetFileSize.FieldName(), parameters.TargetFileSize, sdk.ParameterTypeTable),
		StringTypedParameterValueComputedIf(defs.CatalogSync.FieldName(), parameters.CatalogSync, sdk.ParameterTypeTable),
		IntTypedParameterValueComputedIf(defs.DataRetentionTimeInDays.FieldName(), parameters.DataRetentionTimeInDays, sdk.ParameterTypeTable),
		IntTypedParameterValueComputedIf(defs.MaxDataExtensionTimeInDays.FieldName(), parameters.MaxDataExtensionTimeInDays, sdk.ParameterTypeTable),
		BoolTypedParameterValueComputedIf(defs.EnableDataCompaction.FieldName(), parameters.EnableDataCompaction, sdk.ParameterTypeTable),
		BoolTypedParameterValueComputedIf(defs.EnableIcebergMergeOnRead.FieldName(), parameters.EnableIcebergMergeOnRead, sdk.ParameterTypeTable),
	}
}

var icebergTableSnowflakeManagedParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	icebergTableParametersProvider,
	icebergTableSnowflakeManagedParameterDiffFunctions,
)

func handleIcebergTableSnowflakeManagedParametersCreate(d *schema.ResourceData, req *sdk.CreateIcebergTableRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreateWithMapping(d, defs.ExternalVolume.FieldName(), &req.ExternalVolume, sdk.ParseAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, defs.Catalog.FieldName(), &req.Catalog, stringToStringEnumProvider(sdk.ToIcebergTableCatalog)),
		handleParameterCreateWithMapping(d, defs.TargetFileSize.FieldName(), &req.TargetFileSize, stringToStringEnumProvider(sdk.ToIcebergTableTargetFileSize)),
		handleParameterCreateWithMapping(d, defs.StorageSerializationPolicy.FieldName(), &req.StorageSerializationPolicy, stringToStringEnumProvider(sdk.ToStorageSerializationPolicy)),
		handleParameterCreate(d, defs.CatalogSync.FieldName(), &req.CatalogSync),
		handleParameterCreate(d, defs.DataRetentionTimeInDays.FieldName(), &req.DataRetentionTimeInDays),
		handleParameterCreate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &req.MaxDataExtensionTimeInDays),
		handleParameterCreate(d, defs.EnableDataCompaction.FieldName(), &req.EnableDataCompaction),
		handleParameterCreate(d, defs.EnableIcebergMergeOnRead.FieldName(), &req.EnableIcebergMergeOnRead),
	)
}

// storage_serialization_policy is omitted: it is create-only (ForceNew) and cannot be altered.
func handleIcebergTableSnowflakeManagedParametersUpdate(d *schema.ResourceData, set *sdk.IcebergTableSetPropertiesRequest, unset *sdk.IcebergTableUnsetPropertiesRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.CatalogSync.FieldName(), &set.CatalogSync, &unset.CatalogSync),
		handleParameterUpdate(d, defs.DataRetentionTimeInDays.FieldName(), &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleParameterUpdate(d, defs.EnableDataCompaction.FieldName(), &set.EnableDataCompaction, &unset.EnableDataCompaction),
		handleParameterUpdate(d, defs.EnableIcebergMergeOnRead.FieldName(), &set.EnableIcebergMergeOnRead, &unset.EnableIcebergMergeOnRead),
		handleParameterUpdateWithMapping(d, defs.TargetFileSize.FieldName(), &set.TargetFileSize, &unset.TargetFileSize, stringToStringEnumProvider(sdk.ToIcebergTableTargetFileSize)),
	)
}

// external_volume/catalog come from the raw SHOW PARAMETERS value (see rawParameterValue), not the
// typed ShowParametersDetails output.
func handleIcebergTableSnowflakeManagedParameterRead(d *schema.ResourceData, parameters []*sdk.Parameter, parameterDetails *sdk.IcebergTableParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.ExternalVolume.FieldName(), rawParameterValue(parameters, defs.ExternalVolume.SqlName)),
		setResourceData(d, defs.Catalog.FieldName(), rawParameterValue(parameters, defs.Catalog.SqlName)),
		setResourceData(d, defs.TargetFileSize.FieldName(), parameterDetails.TargetFileSize.Value),
		setResourceData(d, defs.StorageSerializationPolicy.FieldName(), parameterDetails.StorageSerializationPolicy.Value),
		setResourceData(d, defs.CatalogSync.FieldName(), parameterDetails.CatalogSync.Value),
		setResourceData(d, defs.DataRetentionTimeInDays.FieldName(), parameterDetails.DataRetentionTimeInDays.Value),
		setResourceData(d, defs.MaxDataExtensionTimeInDays.FieldName(), parameterDetails.MaxDataExtensionTimeInDays.Value),
		setResourceData(d, defs.EnableDataCompaction.FieldName(), parameterDetails.EnableDataCompaction.Value),
		setResourceData(d, defs.EnableIcebergMergeOnRead.FieldName(), parameterDetails.EnableIcebergMergeOnRead.Value),
	)
}
