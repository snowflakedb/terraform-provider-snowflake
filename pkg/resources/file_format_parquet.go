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

func fileFormatParquetSchema() map[string]*schema.Schema {
	return collections.MergeMaps(fileFormatCommonSchema, parquetFileFormatSchema(""), parquetDescOutputSchema())
}

func parquetDescOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE FILE FORMAT` for this file format.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeFileFormatParquetSchema,
			},
		},
	}
}

func FileFormatParquet() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.FileFormats.DropSafely
		},
	)

	resourceSchema := fileFormatParquetSchema()

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FileFormatParquetResource), TrackingCreateWrapper(resources.FileFormatParquet, CreateFileFormatParquet)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FileFormatParquetResource), TrackingReadWrapper(resources.FileFormatParquet, GetReadFileFormatParquetFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FileFormatParquetResource), TrackingUpdateWrapper(resources.FileFormatParquet, UpdateFileFormatParquet)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FileFormatParquetResource), TrackingDeleteWrapper(resources.FileFormatParquet, deleteFunc)),
		Description:   "Resource used to manage Parquet file formats. For more information, check [file format documentation](https://docs.snowflake.com/en/sql-reference/sql/create-file-format).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FileFormatParquet, customdiff.All(
			ComputedIfAnyAttributeChanged(resourceSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(
				resourceSchema, DescribeOutputAttributeName,
				"name", "type", "compression", "snappy_compression", "binary_as_text", "use_logical_type",
				"trim_space", "use_vectorized_scanner", "replace_invalid_characters", "null_if",
			),
			ComputedIfAnyAttributeChanged(resourceSchema, FullyQualifiedNameAttributeName, "name"),
			RecreateWhenResourceTypeChangedExternally("type", sdk.FileFormatTypeParquet, sdk.ToFileFormatType),
		)),

		Schema: resourceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.FileFormatParquet, ImportFileFormatParquet),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportFileFormatParquet(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.FileFormats.DescribeParquetDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.Type != string(sdk.FileFormatTypeParquet) {
		return nil, fmt.Errorf("invalid file format type, expected %s, got %s", sdk.FileFormatTypeParquet, details.Type)
	}

	var errs []error
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		errs = append(errs, err)
	}

	for key, value := range parquetFileFormatToSchema(details, true) {
		errs = append(errs, d.Set(key, value))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return []*schema.ResourceData{d}, nil
}

func CreateFileFormatParquet(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateParquetFileFormatRequest(id)

	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "compression", request.WithCompression, sdk.ToParquetCompression),
		booleanStringAttributeCreateBuilder(d, "snappy_compression", request.WithSnappyCompression),
		booleanStringAttributeCreateBuilder(d, "binary_as_text", request.WithBinaryAsText),
		booleanStringAttributeCreateBuilder(d, "use_logical_type", request.WithUseLogicalType),
		booleanStringAttributeCreateBuilder(d, "trim_space", request.WithTrimSpace),
		booleanStringAttributeCreateBuilder(d, "use_vectorized_scanner", request.WithUseVectorizedScanner),
		booleanStringAttributeCreateBuilder(d, "replace_invalid_characters", request.WithReplaceInvalidCharacters),
		attributeMappedValueCreateBuilder(d, "null_if", request.WithNullIf, parseNullIfRequest),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.FileFormats.CreateParquet(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return GetReadFileFormatParquetFunc(false)(ctx, d, meta)
}

func GetReadFileFormatParquetFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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

		details, err := client.FileFormats.DescribeParquetDetails(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		describeOutputValues := schemas.FileFormatParquetToSchema(details)

		if withExternalChangesMarking {
			valuesToSet := parquetFileFormatToSchema(details, false)
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

func UpdateFileFormatParquet(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		if err := client.FileFormats.AlterParquet(ctx, sdk.NewAlterParquetFileFormatRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming file format: %w", err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewAlterParquetFileFormatSetRequest()

	errs := errors.Join(
		attributeMappedValueUpdateSetOnlyFallback(d, "compression", &set.Compression, sdk.ToParquetCompression, sdk.ParquetCompressionAuto),
		booleanStringAttributeUnsetFallbackUpdate(d, "snappy_compression", &set.SnappyCompression, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "binary_as_text", &set.BinaryAsText, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "use_logical_type", &set.UseLogicalType, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "trim_space", &set.TrimSpace, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "use_vectorized_scanner", &set.UseVectorizedScanner, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, false),
		attributeMappedValueUpdateSetOnlyFallback(d, "null_if", &set.NullIf, parseNullIfRequest, *sdk.NewNullIfListRequest()),
		stringAttributeUpdateSetOnlyNotEmpty(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(set, sdk.NewAlterParquetFileFormatSetRequest()) {
		if err := client.FileFormats.AlterParquet(ctx, sdk.NewAlterParquetFileFormatRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadFileFormatParquetFunc(false)(ctx, d, meta)
}

// parquetFileFormatToSchema converts the SDK details for a Parquet file format to a Terraform schema.
func parquetFileFormatToSchema(parquet *sdk.FileFormatParquet, setDefaults bool) map[string]any {
	state := map[string]any{
		"compression": parquet.Compression,
		"null_if":     collections.Map(parquet.NullIf, func(v string) any { return v }),
	}
	// SNAPPY_COMPRESSION is not reported back by DESCRIBE FILE FORMAT (it is folded into COMPRESSION = SNAPPY),
	// so it cannot be reliably read back and is always reset to its default.
	state["snappy_compression"] = BooleanDefault
	if setDefaults {
		state["binary_as_text"] = BooleanDefault
		state["use_logical_type"] = BooleanDefault
		state["trim_space"] = BooleanDefault
		state["use_vectorized_scanner"] = BooleanDefault
		state["replace_invalid_characters"] = BooleanDefault
	} else {
		state["binary_as_text"] = booleanStringFromBool(parquet.BinaryAsText)
		state["use_logical_type"] = booleanStringFromBool(parquet.UseLogicalType)
		state["trim_space"] = booleanStringFromBool(parquet.TrimSpace)
		state["use_vectorized_scanner"] = booleanStringFromBool(parquet.UseVectorizedScanner)
		state["replace_invalid_characters"] = booleanStringFromBool(parquet.ReplaceInvalidCharacters)
	}
	return state
}
