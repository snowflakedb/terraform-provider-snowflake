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

var procedureParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelProcedure) {
		procedureParametersSchema[p.FieldName()] = parameterSchema(p)
	}
}

// They do not work in create, that's why are set in alter
func handleProcedureParametersCreate(d *schema.ResourceData, set *sdk.ProcedureSetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.EnableConsoleOutput.FieldName(), &set.EnableConsoleOutput),
		handleParameterCreateWithMapping(d, defs.LogLevel.FieldName(), &set.LogLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, defs.LogEventLevel.FieldName(), &set.LogEventLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, defs.MetricLevel.FieldName(), &set.MetricLevel, sdk.ToMetricLevel),
		handleParameterCreateWithMapping(d, defs.TraceLevel.FieldName(), &set.TraceLevel, sdk.ToTraceLevel),
	)
}

func handleProcedureParametersUpdate(d *schema.ResourceData, set *sdk.ProcedureSetRequest, unset *sdk.ProcedureUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.EnableConsoleOutput.FieldName(), &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
		handleParameterUpdateWithMapping(d, defs.LogLevel.FieldName(), &set.LogLevel, &unset.LogLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, defs.LogEventLevel.FieldName(), &set.LogEventLevel, &unset.LogEventLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, defs.MetricLevel.FieldName(), &set.MetricLevel, &unset.MetricLevel, sdk.ToMetricLevel),
		handleParameterUpdateWithMapping(d, defs.TraceLevel.FieldName(), &set.TraceLevel, &unset.TraceLevel, sdk.ToTraceLevel),
	)
}

func handleProcedureParameterRead(d *schema.ResourceData, parameters *sdk.ProcedureParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.EnableConsoleOutput.FieldName(), parameters.EnableConsoleOutput.Value),
		setResourceData(d, defs.LogLevel.FieldName(), parameters.LogLevel.Value),
		setResourceData(d, defs.LogEventLevel.FieldName(), parameters.LogEventLevel.Value),
		setResourceData(d, defs.MetricLevel.FieldName(), parameters.MetricLevel.Value),
		setResourceData(d, defs.TraceLevel.FieldName(), parameters.TraceLevel.Value),
	)
}

var procedureParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	procedureParametersProvider,
	procedureParameterDiffFunctions,
)

func procedureParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.ProcedureParametersDetails, error) {
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Procedures.ShowParametersDetails(ctx, id)
}

func procedureParameterDiffFunctions(parameters *sdk.ProcedureParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		BoolTypedParameterValueComputedIf(defs.EnableConsoleOutput.FieldName(), parameters.EnableConsoleOutput, sdk.ParameterTypeProcedure),
		StringTypedParameterValueComputedIf(defs.LogLevel.FieldName(), parameters.LogLevel, sdk.ParameterTypeProcedure),
		StringTypedParameterValueComputedIf(defs.LogEventLevel.FieldName(), parameters.LogEventLevel, sdk.ParameterTypeProcedure),
		StringTypedParameterValueComputedIf(defs.MetricLevel.FieldName(), parameters.MetricLevel, sdk.ParameterTypeProcedure),
		StringTypedParameterValueComputedIf(defs.TraceLevel.FieldName(), parameters.TraceLevel, sdk.ParameterTypeProcedure),
	}
}
