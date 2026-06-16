package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	hybridTableParametersSchema     = make(map[string]*schema.Schema)
	hybridTableParametersCustomDiff = ParametersCustomDiff(
		hybridTableParametersProvider,
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterDataRetentionTimeInDays, valueTypeInt, sdk.ParameterTypeTable},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterMaxDataExtensionTimeInDays, valueTypeInt, sdk.ParameterTypeTable},
	)
)

func init() {
	hybridTableParameterFields := []parameterDef[sdk.ObjectParameter]{
		{
			Name:         sdk.ObjectParameterDataRetentionTimeInDays,
			Type:         schema.TypeInt,
			ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)),
			Description:  "Specifies the retention period for the hybrid table so that Time Travel actions can be performed on historical data.",
		},
		{
			Name:         sdk.ObjectParameterMaxDataExtensionTimeInDays,
			Type:         schema.TypeInt,
			ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)),
			Description:  "Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period for the hybrid table to prevent streams on it from becoming stale.",
		},
	}

	for _, field := range hybridTableParameterFields {
		fieldName := strings.ToLower(string(field.Name))
		hybridTableParametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      enrichWithReferenceToParameterDocs(field.Name, field.Description),
			Computed:         true,
			Optional:         true,
			ValidateDiagFunc: field.ValidateDiag,
			DiffSuppressFunc: field.DiffSuppress,
		}
	}
}

func hybridTableParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), hybridTableParametersProviderFunc, sdk.ParseSchemaObjectIdentifier)
}

func hybridTableParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.SchemaObjectIdentifier] {
	return c.HybridTables.ShowParameters
}

// handleHybridTableParametersCreate populates retention parameters directly on a
// CreateHybridTableRequest. Both DATA_RETENTION_TIME_IN_DAYS and
// MAX_DATA_EXTENSION_TIME_IN_DAYS are accepted at CREATE HYBRID TABLE time even
// though the public docs omit them from the syntax diagram (verified against
// production via SHOW PARAMETERS).
func handleHybridTableParametersCreate(d *schema.ResourceData, req *sdk.CreateHybridTableRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, sdk.ObjectParameterDataRetentionTimeInDays, &req.DataRetentionTimeInDays),
		handleParameterCreate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &req.MaxDataExtensionTimeInDays),
	)
}

func handleHybridTableParametersChanges(d *schema.ResourceData, set *sdk.HybridTableSetPropertiesRequest, unset *sdk.HybridTableUnsetPropertiesRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.ObjectParameterDataRetentionTimeInDays, &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
	)
}

func handleHybridTableParameterRead(d *schema.ResourceData, parameters []*sdk.Parameter) diag.Diagnostics {
	for _, parameter := range parameters {
		switch parameter.Key {
		case string(sdk.ObjectParameterDataRetentionTimeInDays),
			string(sdk.ObjectParameterMaxDataExtensionTimeInDays):
			value, err := strconv.Atoi(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}
