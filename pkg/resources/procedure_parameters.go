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
	procedureParametersSchema     = make(map[string]*schema.Schema)
	procedureParametersCustomDiff = ParametersCustomDiff(
		procedureParametersProvider,
		parameter[sdk.ProcedureParameter]{sdk.ProcedureParameterEnableConsoleOutput, valueTypeBool, sdk.ParameterTypeProcedure},
		parameter[sdk.ProcedureParameter]{sdk.ProcedureParameterLogLevel, valueTypeString, sdk.ParameterTypeProcedure},
		parameter[sdk.ProcedureParameter]{sdk.ProcedureParameterMetricLevel, valueTypeString, sdk.ParameterTypeProcedure},
		parameter[sdk.ProcedureParameter]{sdk.ProcedureParameterTraceLevel, valueTypeString, sdk.ParameterTypeProcedure},
	)
)

func init() {
	procedureParameterFields := []parameterDef[sdk.ProcedureParameter]{
		// session params
		{Name: sdk.ProcedureParameterEnableConsoleOutput, Type: schema.TypeBool, Description: "Enable stdout/stderr fast path logging for anonyous stored procs. This is a public parameter (similar to LOG_LEVEL)."},
		{Name: sdk.ProcedureParameterLogLevel, Type: schema.TypeString, Description: "LOG_LEVEL to use when filtering events"},
		{Name: sdk.ProcedureParameterMetricLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToMetricLevel), DiffSuppress: NormalizeAndCompare(sdk.ToMetricLevel), Description: "METRIC_LEVEL value to control whether to emit metrics to Event Table"},
		{Name: sdk.ProcedureParameterTraceLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToTraceLevel), DiffSuppress: NormalizeAndCompare(sdk.ToTraceLevel), Description: "Trace level value to use when generating/filtering trace events"},
	}

	for _, field := range procedureParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		procedureParametersSchema[fieldName] = &schema.Schema{
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

func procedureParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), procedureParametersProviderFunc, sdk.ParseSchemaObjectIdentifierWithArguments)
}

func procedureParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.SchemaObjectIdentifierWithArguments] {
	return c.Procedures.ShowParameters
}

func handleProcedureParameterRead(d *schema.ResourceData, procedureParameters []*sdk.Parameter) error {
	for _, p := range procedureParameters {
		switch p.Key {
		case
			string(sdk.ProcedureParameterLogLevel),
			string(sdk.ProcedureParameterMetricLevel),
			string(sdk.ProcedureParameterTraceLevel):
			if err := d.Set(strings.ToLower(p.Key), p.Value); err != nil {
				return err
			}
		case
			string(sdk.ProcedureParameterEnableConsoleOutput):
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
func handleProcedureParametersCreate(d *schema.ResourceData, set *sdk.ProcedureSetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, sdk.ProcedureParameterEnableConsoleOutput, &set.EnableConsoleOutput),
		handleParameterCreateWithMapping(d, sdk.ProcedureParameterLogLevel, &set.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreateWithMapping(d, sdk.ProcedureParameterMetricLevel, &set.MetricLevel, stringToStringEnumProvider(sdk.ToMetricLevel)),
		handleParameterCreateWithMapping(d, sdk.ProcedureParameterTraceLevel, &set.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
	)
}

func handleProcedureParametersUpdate(d *schema.ResourceData, set *sdk.ProcedureSetRequest, unset *sdk.ProcedureUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.ProcedureParameterEnableConsoleOutput, &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
		handleParameterUpdateWithMapping(d, sdk.ProcedureParameterLogLevel, &set.LogLevel, &unset.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdateWithMapping(d, sdk.ProcedureParameterMetricLevel, &set.MetricLevel, &unset.MetricLevel, stringToStringEnumProvider(sdk.ToMetricLevel)),
		handleParameterUpdateWithMapping(d, sdk.ProcedureParameterTraceLevel, &set.TraceLevel, &unset.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
	)
}
