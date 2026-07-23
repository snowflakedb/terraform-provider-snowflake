package resources

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func fileFormatOrcSchema() map[string]*schema.Schema {
	return collections.MergeMaps(fileFormatCommonSchema, orcFileFormatSchema(""), orcDescOutputSchema())
}

func orcDescOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE FILE FORMAT` for this file format.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeFileFormatOrcSchema,
			},
		},
	}
}

func FileFormatOrc() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.FileFormats.DropSafely
		},
	)

	resourceSchema := fileFormatOrcSchema()

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FileFormatOrcResource), TrackingCreateWrapper(resources.FileFormatOrc, CreateFileFormatOrc)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FileFormatOrcResource), TrackingReadWrapper(resources.FileFormatOrc, GetReadFileFormatOrcFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FileFormatOrcResource), TrackingUpdateWrapper(resources.FileFormatOrc, UpdateFileFormatOrc)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FileFormatOrcResource), TrackingDeleteWrapper(resources.FileFormatOrc, deleteFunc)),
		Description:   "Resource used to manage ORC file formats. For more information, check [file format documentation](https://docs.snowflake.com/en/sql-reference/sql/create-file-format).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FileFormatOrc, customdiff.All(
			ComputedIfAnyAttributeChanged(resourceSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(
				resourceSchema, DescribeOutputAttributeName,
				"name", "type", "trim_space", "replace_invalid_characters", "null_if",
			),
			ComputedIfAnyAttributeChanged(resourceSchema, FullyQualifiedNameAttributeName, "name"),
			RecreateWhenResourceTypeChangedExternally("type", sdk.FileFormatTypeOrc, sdk.ToFileFormatType),
		)),

		Schema: resourceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.FileFormatOrc, ImportFileFormatOrc),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportFileFormatOrc(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.FileFormats.DescribeOrcDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.Type != string(sdk.FileFormatTypeOrc) {
		return nil, fmt.Errorf("invalid file format type, expected %s, got %s", sdk.FileFormatTypeOrc, details.Type)
	}

	var errs []error
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		errs = append(errs, err)
	}

	for key, value := range orcFileFormatToSchema(details, true) {
		errs = append(errs, d.Set(key, value))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return []*schema.ResourceData{d}, nil
}

func CreateFileFormatOrc(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateOrcFileFormatRequest(id)

	errs := errors.Join(
		booleanStringAttributeCreateBuilder(d, "trim_space", request.WithTrimSpace),
		booleanStringAttributeCreateBuilder(d, "replace_invalid_characters", request.WithReplaceInvalidCharacters),
		attributeMappedValueCreateBuilder(d, "null_if", request.WithNullIf, parseNullIfRequest),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.FileFormats.CreateOrc(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return GetReadFileFormatOrcFunc(false)(ctx, d, meta)
}

func GetReadFileFormatOrcFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		fileFormat, err := client.FileFormats.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query file format. Marking the resource as removed.",
						Detail:   fmt.Sprintf("File format id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		details, err := client.FileFormats.DescribeOrcDetails(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		describeOutputValues := schemas.FileFormatOrcToSchema(details)

		if withExternalChangesMarking {
			valuesToSet := orcFileFormatToSchema(details, false)
			mappings := collections.Map(slices.Collect(maps.Keys(valuesToSet)), func(key string) outputMapping {
				return outputMapping{key, key, describeOutputValues[key], valuesToSet[key], nil}
			})
			if err := handleExternalChangesToObjectInFlatDescribeDeepEqual(d, mappings...); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set("comment", fileFormat.Comment),
			d.Set("type", string(fileFormat.Type)),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.FileFormatToSchema(fileFormat)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{describeOutputValues}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateFileFormatOrc(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		if err := client.FileFormats.AlterOrc(ctx, sdk.NewAlterOrcFileFormatRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming file format: %w", err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewAlterOrcFileFormatSetRequest()

	errs := errors.Join(
		booleanStringAttributeUnsetFallbackUpdate(d, "trim_space", &set.TrimSpace, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, false),
		attributeMappedValueUpdateSetOnlyFallback(d, "null_if", &set.NullIf, parseNullIfRequest, *sdk.NewNullIfListRequest()),
		stringAttributeUpdateSetOnlyNotEmpty(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(set, sdk.NewAlterOrcFileFormatSetRequest()) {
		if err := client.FileFormats.AlterOrc(ctx, sdk.NewAlterOrcFileFormatRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadFileFormatOrcFunc(false)(ctx, d, meta)
}

// orcFileFormatToSchema converts the SDK details for an ORC file format to a Terraform schema.
func orcFileFormatToSchema(orc *sdk.FileFormatOrc, setDefaults bool) map[string]any {
	state := map[string]any{
		"null_if": collections.Map(orc.NullIf, func(v string) any { return v }),
	}
	if setDefaults {
		state["trim_space"] = BooleanDefault
		state["replace_invalid_characters"] = BooleanDefault
	} else {
		state["trim_space"] = booleanStringFromBool(orc.TrimSpace)
		state["replace_invalid_characters"] = booleanStringFromBool(orc.ReplaceInvalidCharacters)
	}
	return state
}
