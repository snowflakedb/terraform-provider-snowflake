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

func fileFormatJsonSchema() map[string]*schema.Schema {
	return collections.MergeMaps(fileFormatCommonSchema, jsonFileFormatSchema(""), jsonDescOutputSchema())
}

func FileFormatJson() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.FileFormats.DropSafely
		},
	)

	resourceSchema := fileFormatJsonSchema()

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FileFormatJsonResource), TrackingCreateWrapper(resources.FileFormatJson, CreateFileFormatJson)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FileFormatJsonResource), TrackingReadWrapper(resources.FileFormatJson, GetReadFileFormatJsonFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FileFormatJsonResource), TrackingUpdateWrapper(resources.FileFormatJson, UpdateFileFormatJson)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FileFormatJsonResource), TrackingDeleteWrapper(resources.FileFormatJson, deleteFunc)),
		Description:   "Resource used to manage JSON file formats. For more information, check [file format documentation](https://docs.snowflake.com/en/sql-reference/sql/create-file-format).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FileFormatJson, customdiff.All(
			ComputedIfAnyAttributeChanged(resourceSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(
				resourceSchema, DescribeOutputAttributeName,
				"name", "type", "compression", "date_format", "time_format", "timestamp_format", "binary_format",
				"trim_space", "multi_line", "null_if", "file_extension", "enable_octal", "allow_duplicate",
				"strip_outer_array", "strip_null_values", "replace_invalid_characters", "ignore_utf8_errors", "skip_byte_order_mark",
			),
			ComputedIfAnyAttributeChanged(resourceSchema, FullyQualifiedNameAttributeName, "name"),
			RecreateWhenResourceTypeChangedExternally("type", sdk.FileFormatTypeJson, sdk.ToFileFormatType),
		)),

		Schema: resourceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.FileFormatJson, ImportFileFormatJson),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportFileFormatJson(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.FileFormats.DescribeJsonDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.Type != string(sdk.FileFormatTypeJson) {
		return nil, fmt.Errorf("invalid file format type, expected %s, got %s", sdk.FileFormatTypeJson, details.Type)
	}

	var errs []error
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		errs = append(errs, err)
	}

	for key, value := range jsonFileFormatToSchema(details, true) {
		errs = append(errs, d.Set(key, value))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return []*schema.ResourceData{d}, nil
}

var fileFormatStringOrAutoFallback = *sdk.NewStageFileFormatStringOrAutoRequest().WithAuto(true)

func CreateFileFormatJson(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateJsonFileFormatRequest(id)

	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "compression", request.WithCompression, sdk.ToJsonCompression),
		attributeMappedValueCreateBuilder(d, "date_format", request.WithDateFormat, fileFormatStringOrAutoMapper),
		attributeMappedValueCreateBuilder(d, "time_format", request.WithTimeFormat, fileFormatStringOrAutoMapper),
		attributeMappedValueCreateBuilder(d, "timestamp_format", request.WithTimestampFormat, fileFormatStringOrAutoMapper),
		attributeMappedValueCreateBuilder(d, "binary_format", request.WithBinaryFormat, sdk.ToBinaryFormat),
		booleanStringAttributeCreateBuilder(d, "trim_space", request.WithTrimSpace),
		booleanStringAttributeCreateBuilder(d, "multi_line", request.WithMultiLine),
		stringAttributeCreateBuilder(d, "file_extension", request.WithFileExtension),
		booleanStringAttributeCreateBuilder(d, "enable_octal", request.WithEnableOctal),
		booleanStringAttributeCreateBuilder(d, "allow_duplicate", request.WithAllowDuplicate),
		booleanStringAttributeCreateBuilder(d, "strip_outer_array", request.WithStripOuterArray),
		booleanStringAttributeCreateBuilder(d, "strip_null_values", request.WithStripNullValues),
		booleanStringAttributeCreateBuilder(d, "replace_invalid_characters", request.WithReplaceInvalidCharacters),
		booleanStringAttributeCreateBuilder(d, "ignore_utf8_errors", request.WithIgnoreUtf8Errors),
		booleanStringAttributeCreateBuilder(d, "skip_byte_order_mark", request.WithSkipByteOrderMark),
		attributeMappedValueCreateBuilder(d, "null_if", func(nullIf []sdk.NullString) *sdk.CreateJsonFileFormatRequest {
			request.WithNullIf(nullIf)
			return request
		}, parseNullIf),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.FileFormats.CreateJson(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return GetReadFileFormatJsonFunc(false)(ctx, d, meta)
}

func GetReadFileFormatJsonFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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

		details, err := client.FileFormats.DescribeJsonDetails(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		describeOutputValues := schemas.FileFormatJsonToSchema(details)

		if withExternalChangesMarking {
			valuesToSet := jsonFileFormatToSchema(details, false)
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

func UpdateFileFormatJson(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		if err := client.FileFormats.AlterJson(ctx, sdk.NewAlterJsonFileFormatRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming file format: %w", err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewAlterJsonFileFormatSetRequest()

	errs := errors.Join(
		attributeMappedValueUpdateSetOnlyFallback(d, "compression", &set.Compression, sdk.ToJsonCompression, sdk.JsonCompressionAuto),
		attributeMappedValueUpdateSetOnlyFallback(d, "date_format", &set.DateFormat, fileFormatStringOrAutoMapper, fileFormatStringOrAutoFallback),
		attributeMappedValueUpdateSetOnlyFallback(d, "time_format", &set.TimeFormat, fileFormatStringOrAutoMapper, fileFormatStringOrAutoFallback),
		attributeMappedValueUpdateSetOnlyFallback(d, "timestamp_format", &set.TimestampFormat, fileFormatStringOrAutoMapper, fileFormatStringOrAutoFallback),
		attributeMappedValueUpdateSetOnlyFallback(d, "binary_format", &set.BinaryFormat, sdk.ToBinaryFormat, sdk.BinaryFormatHex),
		booleanStringAttributeUnsetFallbackUpdate(d, "trim_space", &set.TrimSpace, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "multi_line", &set.MultiLine, true),
		stringAttributeUpdateSetOnlyNotEmpty(d, "file_extension", &set.FileExtension),
		booleanStringAttributeUnsetFallbackUpdate(d, "enable_octal", &set.EnableOctal, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "allow_duplicate", &set.AllowDuplicate, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "strip_outer_array", &set.StripOuterArray, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "strip_null_values", &set.StripNullValues, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "ignore_utf8_errors", &set.IgnoreUtf8Errors, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "skip_byte_order_mark", &set.SkipByteOrderMark, true),
		attributeMappedValueUpdateSetOnlySliceFallback(d, "null_if", &set.NullIf, parseNullIf, []sdk.NullString{{S: `\N`}}),
		stringAttributeUpdateSetOnlyNotEmpty(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(set, sdk.NewAlterJsonFileFormatSetRequest()) {
		if err := client.FileFormats.AlterJson(ctx, sdk.NewAlterJsonFileFormatRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadFileFormatJsonFunc(false)(ctx, d, meta)
}

// jsonFileFormatToSchema converts the SDK details for a JSON file format to a Terraform schema.
func jsonFileFormatToSchema(json *sdk.FileFormatJson, setDefaults bool) map[string]any {
	state := map[string]any{
		"compression":      json.Compression,
		"date_format":      json.DateFormat,
		"time_format":      json.TimeFormat,
		"timestamp_format": json.TimestampFormat,
		"binary_format":    json.BinaryFormat,
		"null_if":          collections.Map(json.NullIf, func(v string) any { return v }),
		"file_extension":   json.FileExtension,
	}
	if setDefaults {
		state["ignore_utf8_errors"] = BooleanDefault
		state["skip_byte_order_mark"] = BooleanDefault
		state["trim_space"] = BooleanDefault
		state["multi_line"] = BooleanDefault
		state["allow_duplicate"] = BooleanDefault
		state["strip_outer_array"] = BooleanDefault
		state["strip_null_values"] = BooleanDefault
		state["replace_invalid_characters"] = BooleanDefault
		state["enable_octal"] = BooleanDefault
	} else {
		state["ignore_utf8_errors"] = booleanStringFromBool(json.IgnoreUtf8Errors)
		state["skip_byte_order_mark"] = booleanStringFromBool(json.SkipByteOrderMark)
		state["trim_space"] = booleanStringFromBool(json.TrimSpace)
		state["multi_line"] = booleanStringFromBool(json.MultiLine)
		state["allow_duplicate"] = booleanStringFromBool(json.AllowDuplicate)
		state["strip_outer_array"] = booleanStringFromBool(json.StripOuterArray)
		state["strip_null_values"] = booleanStringFromBool(json.StripNullValues)
		state["replace_invalid_characters"] = booleanStringFromBool(json.ReplaceInvalidCharacters)
		state["enable_octal"] = booleanStringFromBool(json.EnableOctal)
	}
	return state
}

func jsonDescOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE FILE FORMAT` for this file format.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeFileFormatJsonSchema,
			},
		},
	}
}
