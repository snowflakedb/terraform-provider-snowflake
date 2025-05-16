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

func ProcedureJavascript() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ProcedureJavascriptResource), TrackingCreateWrapper(resources.ProcedureJavascript, CreateContextProcedureJavascript)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ProcedureJavascriptResource), TrackingReadWrapper(resources.ProcedureJavascript, ReadContextProcedureJavascript)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ProcedureJavascriptResource), TrackingUpdateWrapper(resources.ProcedureJavascript, UpdateProcedure("JAVASCRIPT", ReadContextProcedureJavascript))),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ProcedureJavascriptResource), TrackingDeleteWrapper(resources.ProcedureJavascript, DeleteProcedure)),
		Description:   "Resource used to manage javascript procedure objects. For more information, check [procedure documentation](https://docs.snowflake.com/en/sql-reference/sql/create-procedure).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ProcedureJavascript, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(javascriptProcedureSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javascriptProcedureSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(procedureParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllProcedureParameters), strings.ToLower)...),
			procedureParametersCustomDiff,
			// The language check is more for the future.
			// Currently, almost all attributes are marked as forceNew.
			// When language changes, these attributes also change, causing the object to recreate either way.
			// The only option is java staged <-> scala staged (however scala need runtime_version which may interfere).
			RecreateWhenResourceStringFieldChangedExternally("procedure_language", "JAVASCRIPT"),
		)),

		Schema: collections.MergeMaps(javascriptProcedureSchema, procedureParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ProcedureJavascript, ImportProcedure),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextProcedureJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	database := d.Get("database").(string)
	sc := d.Get("schema").(string)
	name := d.Get("name").(string)

	argumentRequests, err := parseProcedureArgumentsCommon(d)
	if err != nil {
		return diag.FromErr(err)
	}
	returnTypeRaw := d.Get("return_type").(string)
	returnDataType, err := datatypes.ParseDataType(returnTypeRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	procedureDefinition := d.Get("procedure_definition").(string)

	argumentDataTypes := collections.Map(argumentRequests, func(r sdk.ProcedureArgumentRequest) datatypes.DataType { return r.ArgDataType })
	id := sdk.NewSchemaObjectIdentifierWithArgumentsNormalized(database, sc, name, argumentDataTypes...)
	request := sdk.NewCreateForJavaScriptProcedureRequestDefinitionWrapped(id.SchemaObjectId(), returnDataType, procedureDefinition).
		WithArguments(argumentRequests)

	errs := errors.Join(
		booleanStringAttributeCreateBuilder(d, "is_secure", request.WithSecure),
		attributeMappedValueCreateBuilder[string](d, "null_input_behavior", request.WithNullInputBehavior, sdk.ToNullInputBehavior),
		attributeMappedValueCreateBuilder[string](d, "execute_as", request.WithExecuteAs, sdk.ToExecuteAs),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.Procedures.CreateForJavaScript(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// parameters do not work in create procedure (query does not fail but parameters stay unchanged)
	setRequest := sdk.NewProcedureSetRequest()
	if parametersCreateDiags := handleProcedureParametersCreate(d, setRequest); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}
	if !reflect.DeepEqual(*setRequest, *sdk.NewProcedureSetRequest()) {
		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithSet(*setRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextProcedureJavascript(ctx, d, meta)
}

func ReadContextProcedureJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	allProcedureDetails, diags := queryAllProcedureDetailsCommon(ctx, d, client, id)
	if diags != nil {
		return diags
	}

	// TODO [SNOW-1348103]: handle external changes marking
	// TODO [SNOW-1348103]: handle setting state to value from config

	errs := errors.Join(
		// not reading is_secure on purpose (handled as external change to show output)
		readFunctionOrProcedureArguments(d, allProcedureDetails.procedureDetails.NormalizedArguments),
		d.Set("return_type", allProcedureDetails.procedureDetails.ReturnDataType.ToSql()),
		// not reading null_input_behavior on purpose (handled as external change to show output)
		// not reading execute_as on purpose (handled as external change to show output)
		d.Set("comment", allProcedureDetails.procedure.Description),
		setOptionalFromStringPtr(d, "procedure_definition", allProcedureDetails.procedureDetails.Body),
		d.Set("procedure_language", allProcedureDetails.procedureDetails.Language),

		handleProcedureParameterRead(d, allProcedureDetails.procedureParameters),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ProcedureToSchema(allProcedureDetails.procedure)}),
		d.Set(ParametersAttributeName, []map[string]any{schemas.ProcedureParametersToSchema(allProcedureDetails.procedureParameters)}),
	)
	if errs != nil {
		return diag.FromErr(err)
	}

	return nil
}
