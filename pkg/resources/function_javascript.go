package resources

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FunctionJavascript() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FunctionJavascriptResource), TrackingCreateWrapper(resources.FunctionJavascript, CreateContextFunctionJavascript)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FunctionJavascriptResource), TrackingReadWrapper(resources.FunctionJavascript, ReadContextFunctionJavascript)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FunctionJavascriptResource), TrackingUpdateWrapper(resources.FunctionJavascript, UpdateFunction("JAVASCRIPT", ReadContextFunctionJavascript))),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FunctionJavascriptResource), TrackingDeleteWrapper(resources.FunctionJavascript, DeleteFunction)),
		Description:   "Resource used to manage javascript function objects. For more information, check [function documentation](https://docs.snowflake.com/en/sql-reference/sql/create-function).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionJavascript, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(javascriptFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javascriptFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// The language check is more for the future.
			// Currently, almost all attributes are marked as forceNew.
			// When language changes, these attributes also change, causing the object to recreate either way.
			// The only potential option is java staged <-> scala staged (however scala need runtime_version which may interfere).
			RecreateWhenResourceStringFieldChangedExternally("function_language", "JAVASCRIPT"),
		)),

		Schema: collections.MergeMaps(javascriptFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.FunctionJavascript, ImportFunction),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextFunctionJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	database := d.Get("database").(string)
	sc := d.Get("schema").(string)
	name := d.Get("name").(string)

	argumentRequests, err := parseFunctionArgumentsCommon(d)
	if err != nil {
		return diag.FromErr(err)
	}
	returns, err := parseFunctionReturnsCommon(d)
	if err != nil {
		return diag.FromErr(err)
	}
	functionDefinition := d.Get("function_definition").(string)

	argumentDataTypes := collections.Map(argumentRequests, func(r sdk.FunctionArgumentRequest) datatypes.DataType { return r.ArgDataType })
	id := sdk.NewSchemaObjectIdentifierWithArgumentsNormalized(database, sc, name, argumentDataTypes...)
	request := sdk.NewCreateForJavascriptFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, functionDefinition).
		WithArguments(argumentRequests)

	errs := errors.Join(
		booleanStringAttributeCreateBuilder(d, "is_secure", request.WithSecure),
		attributeMappedValueCreateBuilder[string](d, "null_input_behavior", request.WithNullInputBehavior, sdk.ToNullInputBehavior),
		attributeMappedValueCreateBuilder[string](d, "return_results_behavior", request.WithReturnResultsBehavior, sdk.ToReturnResultsBehavior),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.Functions.CreateForJavascript(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// parameters do not work in create function (query does not fail but parameters stay unchanged)
	setRequest := sdk.NewFunctionSetRequest()
	if parametersCreateDiags := handleFunctionParametersCreate(d, setRequest); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}
	if !reflect.DeepEqual(*setRequest, *sdk.NewFunctionSetRequest()) {
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSet(*setRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextFunctionJavascript(ctx, d, meta)
}

func ReadContextFunctionJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	allFunctionDetails, diags := queryAllFunctionDetailsCommon(ctx, d, client, id)
	if diags != nil {
		return diags
	}

	// TODO [SNOW-1348103]: handle external changes marking
	// TODO [SNOW-1348103]: handle setting state to value from config

	errs := errors.Join(
		// not reading is_secure on purpose (handled as external change to show output)
		readFunctionOrProcedureArguments(d, allFunctionDetails.functionDetails.NormalizedArguments),
		d.Set("return_type", allFunctionDetails.functionDetails.ReturnDataType.ToSql()),
		// not reading null_input_behavior on purpose (handled as external change to show output)
		// not reading return_results_behavior on purpose (handled as external change to show output)
		d.Set("comment", allFunctionDetails.function.Description),
		setRequiredFromStringPtr(d, "handler", allFunctionDetails.functionDetails.Handler),
		setOptionalFromStringPtr(d, "function_definition", allFunctionDetails.functionDetails.Body),
		d.Set("function_language", allFunctionDetails.functionDetails.Language),

		handleFunctionParameterRead(d, allFunctionDetails.functionParameters),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.FunctionToSchema(allFunctionDetails.function)}),
		d.Set(ParametersAttributeName, []map[string]any{schemas.FunctionParametersToSchema(allFunctionDetails.functionParameters)}),
	)
	if errs != nil {
		return diag.FromErr(err)
	}

	return nil
}
