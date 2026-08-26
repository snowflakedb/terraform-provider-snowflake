package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	providerresources "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataMetricFunction() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(providerresources.DataMetricFunction, CreateDataMetricFunction),
		ReadContext:   TrackingReadWrapper(providerresources.DataMetricFunction, ReadDataMetricFunction),
		DeleteContext: TrackingDeleteWrapper(providerresources.DataMetricFunction, DeleteDataMetricFunction),
		Schema: map[string]*schema.Schema{
			"name":     {Type: schema.TypeString, Required: true, ForceNew: true, DiffSuppressFunc: suppressIdentifierQuoting},
			"database": {Type: schema.TypeString, Required: true, ForceNew: true, DiffSuppressFunc: suppressIdentifierQuoting},
			"schema":   {Type: schema.TypeString, Required: true, ForceNew: true, DiffSuppressFunc: suppressIdentifierQuoting},
			"argument": {Type: schema.TypeList, Required: true, MinItems: 1, ForceNew: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"name": {Type: schema.TypeString, Required: true},
				"column": {Type: schema.TypeList, Required: true, MinItems: 1, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"name": {Type: schema.TypeString, Required: true},
					"type": {Type: schema.TypeString, Required: true, ValidateDiagFunc: IsDataTypeValid, DiffSuppressFunc: DiffSuppressDataTypes},
				}}},
			}}},
			"body":                          {Type: schema.TypeString, Required: true, ForceNew: true, DiffSuppressFunc: DiffSuppressStatement},
			"comment":                       {Type: schema.TypeString, Optional: true, ForceNew: true},
			FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
		},
		Importer: &schema.ResourceImporter{StateContext: ImportDataMetricFunction},
	}
}

func dataMetricFunctionRequest(d *schema.ResourceData) (*sdk.CreateDataMetricFunctionRequest, sdk.SchemaObjectIdentifier, error) {
	id := sdk.NewSchemaObjectIdentifier(d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string))
	arguments := make([]sdk.DataMetricFunctionArgumentRequest, 0, len(d.Get("argument").([]any)))
	for _, rawArgument := range d.Get("argument").([]any) {
		argument := rawArgument.(map[string]any)
		columns := make([]sdk.DataMetricFunctionColumnRequest, 0, len(argument["column"].([]any)))
		for _, rawColumn := range argument["column"].([]any) {
			column := rawColumn.(map[string]any)
			dataType, err := datatypes.ParseDataType(column["type"].(string))
			if err != nil {
				return nil, id, err
			}
			columns = append(columns, sdk.DataMetricFunctionColumnRequest{Name: column["name"].(string), DataType: dataType})
		}
		arguments = append(arguments, sdk.DataMetricFunctionArgumentRequest{Name: argument["name"].(string), Columns: columns})
	}
	request := sdk.NewCreateDataMetricFunctionRequest(id, arguments, d.Get("body").(string))
	if comment, ok := d.GetOk("comment"); ok {
		request.WithComment(comment.(string))
	}
	return request, id, nil
}

func CreateDataMetricFunction(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	request, id, err := dataMetricFunctionRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := meta.(*provider.Context).Client.DataMetricFunctions.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadDataMetricFunction(ctx, d, meta)
}

func ReadDataMetricFunction(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()))
}

func DeleteDataMetricFunction(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	request, _, err := dataMetricFunctionRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := meta.(*provider.Context).Client.DataMetricFunctions.Drop(ctx, sdk.NewDropDataMetricFunctionRequest(request.Name, request.Arguments).WithIfExists(true)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ImportDataMetricFunction(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
