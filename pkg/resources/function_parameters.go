package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	functionParametersSchema     = make(map[string]*schema.Schema)
	functionParametersCustomDiff = ParametersCustomDiff(
		functionParametersProvider,
		parameter[sdk.FunctionParameter]{sdk.FunctionParameterEnableConsoleOutput, valueTypeBool, sdk.ParameterTypeFunction},
		parameter[sdk.FunctionParameter]{sdk.FunctionParameterLogLevel, valueTypeString, sdk.ParameterTypeFunction},
		parameter[sdk.FunctionParameter]{sdk.FunctionParameterMetricLevel, valueTypeString, sdk.ParameterTypeFunction},
		parameter[sdk.FunctionParameter]{sdk.FunctionParameterTraceLevel, valueTypeString, sdk.ParameterTypeFunction},
	)
)

func init() {
	functionParameterFields := []parameterDef[sdk.FunctionParameter]{
		// session params
		{Name: sdk.FunctionParameterEnableConsoleOutput, Type: schema.TypeBool, Description: "Enable stdout/stderr fast path logging for anonymous stored procs. This is a public parameter (similar to LOG_LEVEL)."},
		{Name: sdk.FunctionParameterLogLevel, Type: schema.TypeString, Description: "LOG_LEVEL to use when filtering events"},
		{Name: sdk.FunctionParameterMetricLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToMetricLevel), DiffSuppress: NormalizeAndCompare(sdk.ToMetricLevel), Description: "METRIC_LEVEL value to control whether to emit metrics to Event Table"},
		{Name: sdk.FunctionParameterTraceLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToTraceLevel), DiffSuppress: NormalizeAndCompare(sdk.ToTraceLevel), Description: "Trace level value to use when generating/filtering trace events"},
	}

	for _, field := range functionParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		functionParametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      enrichWithReferenceToParameterDocs(field.Name, field.Description),
			Computed:         true,
			Optional:         true,
			ValidateDiagFunc: field.ValidateDiag,
			DiffSuppressFunc: field.DiffSuppress,
			ConflictsWith:    field.ConflictsWith,
		}
	}
}

func functionParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), functionParametersProviderFunc, sdk.ParseSchemaObjectIdentifierWithArguments)
}

func functionParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.SchemaObjectIdentifierWithArguments] {
	return c.Functions.ShowParameters
}

func handleFunctionParameterRead(d *schema.ResourceData, functionParameters []*sdk.Parameter) error {
	for _, p := range functionParameters {
		switch p.Key {
		case
			string(sdk.FunctionParameterLogLevel),
			string(sdk.FunctionParameterMetricLevel),
			string(sdk.FunctionParameterTraceLevel):
			if err := d.Set(strings.ToLower(p.Key), p.Value); err != nil {
				return err
			}
		case
			string(sdk.FunctionParameterEnableConsoleOutput):
			value, err := strconv.ParseBool(p.Value)
			if err != nil {
				return err
			}
			if err := d.Set(strings.ToLower(p.Key), value); err != nil {
				return err
			}
		}
	}

	return nil
}

// They do not work in create, that's why are set in alter
func handleFunctionParametersCreate(d *schema.ResourceData, set *sdk.FunctionSetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, sdk.FunctionParameterEnableConsoleOutput, &set.EnableConsoleOutput),
		handleParameterCreateWithMapping(d, sdk.FunctionParameterLogLevel, &set.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreateWithMapping(d, sdk.FunctionParameterMetricLevel, &set.MetricLevel, stringToStringEnumProvider(sdk.ToMetricLevel)),
		handleParameterCreateWithMapping(d, sdk.FunctionParameterTraceLevel, &set.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
	)
}

func handleFunctionParametersUpdate(d *schema.ResourceData, set *sdk.FunctionSetRequest, unset *sdk.FunctionUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.FunctionParameterEnableConsoleOutput, &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
		handleParameterUpdateWithMapping(d, sdk.FunctionParameterLogLevel, &set.LogLevel, &unset.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdateWithMapping(d, sdk.FunctionParameterMetricLevel, &set.MetricLevel, &unset.MetricLevel, stringToStringEnumProvider(sdk.ToMetricLevel)),
		handleParameterUpdateWithMapping(d, sdk.FunctionParameterTraceLevel, &set.TraceLevel, &unset.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
	)
}
