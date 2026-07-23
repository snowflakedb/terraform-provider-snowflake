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

func fileFormatAvroSchema() map[string]*schema.Schema {
	return collections.MergeMaps(fileFormatCommonSchema, avroFileFormatSchema(""), avroDescOutputSchema())
}

func avroDescOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE FILE FORMAT` for this file format.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeFileFormatAvroSchema,
			},
		},
	}
}

func FileFormatAvro() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.FileFormats.DropSafely
		},
	)

	resourceSchema := fileFormatAvroSchema()

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FileFormatAvroResource), TrackingCreateWrapper(resources.FileFormatAvro, CreateFileFormatAvro)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FileFormatAvroResource), TrackingReadWrapper(resources.FileFormatAvro, GetReadFileFormatAvroFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FileFormatAvroResource), TrackingUpdateWrapper(resources.FileFormatAvro, UpdateFileFormatAvro)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FileFormatAvroResource), TrackingDeleteWrapper(resources.FileFormatAvro, deleteFunc)),
		Description:   "Resource used to manage AVRO file formats. For more information, check [file format documentation](https://docs.snowflake.com/en/sql-reference/sql/create-file-format).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FileFormatAvro, customdiff.All(
			ComputedIfAnyAttributeChanged(resourceSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(
				resourceSchema, DescribeOutputAttributeName,
				"name", "type", "compression", "trim_space", "replace_invalid_characters", "null_if",
			),
			ComputedIfAnyAttributeChanged(resourceSchema, FullyQualifiedNameAttributeName, "name"),
			RecreateWhenResourceTypeChangedExternally("type", sdk.FileFormatTypeAvro, sdk.ToFileFormatType),
		)),

		Schema: resourceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.FileFormatAvro, ImportFileFormatAvro),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportFileFormatAvro(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.FileFormats.DescribeAvroDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.Type != string(sdk.FileFormatTypeAvro) {
		return nil, fmt.Errorf("invalid file format type, expected %s, got %s", sdk.FileFormatTypeAvro, details.Type)
	}

	var errs []error
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		errs = append(errs, err)
	}

	for key, value := range avroFileFormatToSchema(details, true) {
		errs = append(errs, d.Set(key, value))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return []*schema.ResourceData{d}, nil
}

func CreateFileFormatAvro(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateAvroFileFormatRequest(id)

	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "compression", request.WithCompression, sdk.ToAvroCompression),
		booleanStringAttributeCreateBuilder(d, "trim_space", request.WithTrimSpace),
		booleanStringAttributeCreateBuilder(d, "replace_invalid_characters", request.WithReplaceInvalidCharacters),
		attributeMappedValueCreateBuilder(d, "null_if", request.WithNullIf, func(v any) (sdk.NullIfListRequest, error) {
			nullIf, err := parseNullIf(v)
			if err != nil {
				return sdk.NullIfListRequest{}, err
			}
			return *sdk.NewNullIfListRequest().WithNullIf(nullIf), nil
		}),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.FileFormats.CreateAvro(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return GetReadFileFormatAvroFunc(false)(ctx, d, meta)
}

func GetReadFileFormatAvroFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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

		details, err := client.FileFormats.DescribeAvroDetails(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		describeOutputValues := schemas.FileFormatAvroToSchema(details)

		if withExternalChangesMarking {
			valuesToSet := avroFileFormatToSchema(details, false)
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

func UpdateFileFormatAvro(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		if err := client.FileFormats.AlterAvro(ctx, sdk.NewAlterAvroFileFormatRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming file format: %w", err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewAlterAvroFileFormatSetRequest()

	errs := errors.Join(
		attributeMappedValueUpdateSetOnlyFallback(d, "compression", &set.Compression, sdk.ToAvroCompression, sdk.AvroCompressionAuto),
		booleanStringAttributeUnsetFallbackUpdate(d, "trim_space", &set.TrimSpace, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, false),
		attributeMappedValueUpdateSetOnlyFallback(d, "null_if", &set.NullIf, parseNullIfRequest, *sdk.NewNullIfListRequest()),
		stringAttributeUpdateSetOnlyNotEmpty(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(set, sdk.NewAlterAvroFileFormatSetRequest()) {
		if err := client.FileFormats.AlterAvro(ctx, sdk.NewAlterAvroFileFormatRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadFileFormatAvroFunc(false)(ctx, d, meta)
}

// avroFileFormatToSchema converts the SDK details for an AVRO file format to a Terraform schema.
func avroFileFormatToSchema(avro *sdk.FileFormatAvro, setDefaults bool) map[string]any {
	state := map[string]any{
		"compression": avro.Compression,
		"null_if":     collections.Map(avro.NullIf, func(v string) any { return v }),
	}
	if setDefaults {
		state["trim_space"] = BooleanDefault
		state["replace_invalid_characters"] = BooleanDefault
	} else {
		state["trim_space"] = booleanStringFromBool(avro.TrimSpace)
		state["replace_invalid_characters"] = booleanStringFromBool(avro.ReplaceInvalidCharacters)
	}
	return state
}
