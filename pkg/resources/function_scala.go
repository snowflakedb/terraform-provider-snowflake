package resources

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FunctionScala() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FunctionScalaResource), TrackingCreateWrapper(resources.FunctionScala, CreateContextFunctionScala)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FunctionScalaResource), TrackingReadWrapper(resources.FunctionScala, ReadContextFunctionScala)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FunctionScalaResource), TrackingUpdateWrapper(resources.FunctionScala, UpdateFunction("SCALA", ReadContextFunctionScala))),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FunctionScalaResource), TrackingDeleteWrapper(resources.FunctionScala, DeleteFunction)),
		Description:   "Resource used to manage scala function objects. For more information, check [function documentation](https://docs.snowflake.com/en/sql-reference/sql/create-function).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionScala, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(scalaFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(scalaFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// The language check is more for the future.
			// Currently, almost all attributes are marked as forceNew.
			// When language changes, these attributes also change, causing the object to recreate either way.
			// The only potential option is java staged <-> scala staged (however scala need runtime_version which may interfere).
			RecreateWhenResourceStringFieldChangedExternally("function_language", "SCALA"),
		)),

		Schema: collections.MergeMaps(scalaFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.FunctionScala, ImportFunction),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextFunctionScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	database := d.Get("database").(string)
	sc := d.Get("schema").(string)
	name := d.Get("name").(string)

	argumentRequests, err := parseFunctionArgumentsCommon(d)
	if err != nil {
		return diag.FromErr(err)
	}
	returnTypeRaw := d.Get("return_type").(string)
	returnDataType, err := datatypes.ParseDataType(returnTypeRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	handler := d.Get("handler").(string)
	runtimeVersion := d.Get("runtime_version").(string)

	argumentDataTypes := collections.Map(argumentRequests, func(r sdk.FunctionArgumentRequest) datatypes.DataType { return r.ArgDataType })
	id := sdk.NewSchemaObjectIdentifierWithArgumentsNormalized(database, sc, name, argumentDataTypes...)
	request := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), returnDataType, runtimeVersion, handler).
		WithArguments(argumentRequests)

	errs := errors.Join(
		booleanStringAttributeCreateBuilder(d, "is_secure", request.WithSecure),
		attributeMappedValueCreateBuilder[string](d, "null_input_behavior", request.WithNullInputBehavior, sdk.ToNullInputBehavior),
		attributeMappedValueCreateBuilder[string](d, "return_results_behavior", request.WithReturnResultsBehavior, sdk.ToReturnResultsBehavior),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		setFunctionImportsInBuilder(d, request.WithImports),
		setFunctionPackagesInBuilder(d, request.WithPackages),
		setExternalAccessIntegrationsInBuilder(d, request.WithExternalAccessIntegrations),
		setSecretsInBuilder(d, request.WithSecrets),
		setFunctionTargetPathInBuilder(d, request.WithTargetPath),
		stringAttributeCreateBuilder(d, "function_definition", request.WithFunctionDefinitionWrapped),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.Functions.CreateForScala(ctx, request); err != nil {
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

	return ReadContextFunctionScala(ctx, d, meta)
}

func ReadContextFunctionScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		setOptionalFromStringPtr(d, "runtime_version", allFunctionDetails.functionDetails.RuntimeVersion),
		d.Set("comment", allFunctionDetails.function.Description),
		readFunctionOrProcedureImports(d, allFunctionDetails.functionDetails.NormalizedImports),
		d.Set("packages", allFunctionDetails.functionDetails.NormalizedPackages),
		setRequiredFromStringPtr(d, "handler", allFunctionDetails.functionDetails.Handler),
		readFunctionOrProcedureExternalAccessIntegrations(d, allFunctionDetails.functionDetails.NormalizedExternalAccessIntegrations),
		readFunctionOrProcedureSecrets(d, allFunctionDetails.functionDetails.NormalizedSecrets),
		readFunctionOrProcedureTargetPath(d, allFunctionDetails.functionDetails.NormalizedTargetPath),
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
