package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var hybridTableParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelHybridTable) {
		hybridTableParametersSchema[p.FieldName()] = parameterSchema(p)
	}
	// Catalog descriptions for these shared parameters are written for databases.
	hybridTableParametersSchema[defs.DataRetentionTimeInDays.FieldName()].Description = enrichWithReferenceToParameterDocs(
		defs.DataRetentionTimeInDays.SqlName,
		"Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the hybrid table.",
	)
	hybridTableParametersSchema[defs.MaxDataExtensionTimeInDays.FieldName()].Description = enrichWithReferenceToParameterDocs(
		defs.MaxDataExtensionTimeInDays.SqlName,
		"Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period for the hybrid table to prevent streams on it from becoming stale.",
	)
}

// handleHybridTableParametersCreate populates retention parameters directly on a
// CreateHybridTableRequest. Both DATA_RETENTION_TIME_IN_DAYS and
// MAX_DATA_EXTENSION_TIME_IN_DAYS are accepted at CREATE HYBRID TABLE time even
// though the public docs omit them from the syntax diagram (verified against
// production via SHOW PARAMETERS).
func handleHybridTableParametersCreate(d *schema.ResourceData, req *sdk.CreateHybridTableRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.DataRetentionTimeInDays.FieldName(), &req.DataRetentionTimeInDays),
		handleParameterCreate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &req.MaxDataExtensionTimeInDays),
	)
}

func handleHybridTableParametersChanges(d *schema.ResourceData, set *sdk.HybridTableSetPropertiesRequest, unset *sdk.HybridTableUnsetPropertiesRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.DataRetentionTimeInDays.FieldName(), &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
	)
}

func handleHybridTableParameterRead(d *schema.ResourceData, parameters *sdk.HybridTableParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.DataRetentionTimeInDays.FieldName(), parameters.DataRetentionTimeInDays.Value),
		setResourceData(d, defs.MaxDataExtensionTimeInDays.FieldName(), parameters.MaxDataExtensionTimeInDays.Value),
	)
}

var hybridTableParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	hybridTableParametersProvider,
	hybridTableParameterDiffFunctions,
)

func hybridTableParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.HybridTableParametersDetails, error) {
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.HybridTables.ShowParametersDetails(ctx, id)
}

func hybridTableParameterDiffFunctions(parameters *sdk.HybridTableParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IntTypedParameterValueComputedIf(defs.DataRetentionTimeInDays.FieldName(), parameters.DataRetentionTimeInDays, sdk.ParameterTypeTable),
		IntTypedParameterValueComputedIf(defs.MaxDataExtensionTimeInDays.FieldName(), parameters.MaxDataExtensionTimeInDays, sdk.ParameterTypeTable),
	}
}
