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

var functionParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelFunction) {
		functionParametersSchema[p.FieldName()] = parameterSchema(p)
	}
}

// They do not work in create, that's why are set in alter
func handleFunctionParametersCreate(d *schema.ResourceData, set *sdk.FunctionSetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.EnableConsoleOutput.FieldName(), &set.EnableConsoleOutput),
		handleParameterCreateWithMapping(d, defs.LogLevel.FieldName(), &set.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreateWithMapping(d, defs.LogEventLevel.FieldName(), &set.LogEventLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreateWithMapping(d, defs.MetricLevel.FieldName(), &set.MetricLevel, stringToStringEnumProvider(sdk.ToMetricLevel)),
		handleParameterCreateWithMapping(d, defs.TraceLevel.FieldName(), &set.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
	)
}

func handleFunctionParametersUpdate(d *schema.ResourceData, set *sdk.FunctionSetRequest, unset *sdk.FunctionUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.EnableConsoleOutput.FieldName(), &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
		handleParameterUpdateWithMapping(d, defs.LogLevel.FieldName(), &set.LogLevel, &unset.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdateWithMapping(d, defs.LogEventLevel.FieldName(), &set.LogEventLevel, &unset.LogEventLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdateWithMapping(d, defs.MetricLevel.FieldName(), &set.MetricLevel, &unset.MetricLevel, stringToStringEnumProvider(sdk.ToMetricLevel)),
		handleParameterUpdateWithMapping(d, defs.TraceLevel.FieldName(), &set.TraceLevel, &unset.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
	)
}

func handleFunctionParameterRead(d *schema.ResourceData, parameters *sdk.FunctionParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.EnableConsoleOutput.FieldName(), parameters.EnableConsoleOutput.Value),
		setResourceData(d, defs.LogLevel.FieldName(), parameters.LogLevel.Value),
		setResourceData(d, defs.LogEventLevel.FieldName(), parameters.LogEventLevel.Value),
		setResourceData(d, defs.MetricLevel.FieldName(), parameters.MetricLevel.Value),
		setResourceData(d, defs.TraceLevel.FieldName(), parameters.TraceLevel.Value),
	)
}

var functionParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	functionParametersProvider,
	functionParameterDiffFunctions,
)

func functionParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.FunctionParametersDetails, error) {
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Functions.ShowParametersDetails(ctx, id)
}

func functionParameterDiffFunctions(parameters *sdk.FunctionParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		BoolTypedParameterValueComputedIf(defs.EnableConsoleOutput.FieldName(), parameters.EnableConsoleOutput, sdk.ParameterTypeFunction),
		StringTypedParameterValueComputedIf(defs.LogLevel.FieldName(), parameters.LogLevel, sdk.ParameterTypeFunction),
		StringTypedParameterValueComputedIf(defs.LogEventLevel.FieldName(), parameters.LogEventLevel, sdk.ParameterTypeFunction),
		StringTypedParameterValueComputedIf(defs.MetricLevel.FieldName(), parameters.MetricLevel, sdk.ParameterTypeFunction),
		StringTypedParameterValueComputedIf(defs.TraceLevel.FieldName(), parameters.TraceLevel, sdk.ParameterTypeFunction),
	}
}
