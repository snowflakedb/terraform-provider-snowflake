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

func fileFormatCsvSchema() map[string]*schema.Schema {
	return collections.MergeMaps(fileFormatCommonSchema, csvFileFormatSchema(""), csvDescOutputSchema())
}

func csvDescOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE FILE FORMAT` for this file format.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeFileFormatCsvSchema,
			},
		},
	}
}

func FileFormatCsv() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.FileFormats.DropSafely
		},
	)

	resourceSchema := fileFormatCsvSchema()

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FileFormatCsvResource), TrackingCreateWrapper(resources.FileFormatCsv, CreateFileFormatCsv)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FileFormatCsvResource), TrackingReadWrapper(resources.FileFormatCsv, GetReadFileFormatCsvFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FileFormatCsvResource), TrackingUpdateWrapper(resources.FileFormatCsv, UpdateFileFormatCsv)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FileFormatCsvResource), TrackingDeleteWrapper(resources.FileFormatCsv, deleteFunc)),
		Description:   "Resource used to manage CSV file formats. For more information, check [file format documentation](https://docs.snowflake.com/en/sql-reference/sql/create-file-format).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FileFormatCsv, customdiff.All(
			ComputedIfAnyAttributeChanged(resourceSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(
				resourceSchema, DescribeOutputAttributeName,
				"name", "type", "compression", "record_delimiter", "field_delimiter", "multi_line", "file_extension",
				"parse_header", "skip_header", "skip_blank_lines", "date_format", "time_format", "timestamp_format",
				"binary_format", "escape", "escape_unenclosed_field", "trim_space", "field_optionally_enclosed_by",
				"null_if", "error_on_column_count_mismatch", "replace_invalid_characters", "empty_field_as_null",
				"skip_byte_order_mark", "encoding",
			),
			ComputedIfAnyAttributeChanged(resourceSchema, FullyQualifiedNameAttributeName, "name"),
			RecreateWhenResourceTypeChangedExternally("type", sdk.FileFormatTypeCsv, sdk.ToFileFormatType),
		)),

		Schema: resourceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.FileFormatCsv, ImportFileFormatCsv),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportFileFormatCsv(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.FileFormats.DescribeCsvDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.Type != sdk.FileFormatTypeCsv {
		return nil, fmt.Errorf("invalid file format type, expected %s, got %s", sdk.FileFormatTypeCsv, details.Type)
	}

	var errs []error
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		errs = append(errs, err)
	}

	valuesToSet := csvFileFormatToSchema(details, true)
	// skip_header uses IntDefault as its "not managed" marker (like the boolean fields use BooleanDefault),
	// so it is not imported to avoid a permanent diff for configurations that do not set it.
	valuesToSet["skip_header"] = IntDefault

	for key, value := range valuesToSet {
		errs = append(errs, d.Set(key, value))
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateFileFormatCsv(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateCsvFileFormatRequest(id)

	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "compression", request.WithCompression, sdk.ToCsvCompression),
		attributeMappedValueCreateBuilder(d, "record_delimiter", request.WithRecordDelimiter, fileFormatStringOrNoneMapper),
		attributeMappedValueCreateBuilder(d, "field_delimiter", request.WithFieldDelimiter, fileFormatStringOrNoneMapper),
		booleanStringAttributeCreateBuilder(d, "multi_line", request.WithMultiLine),
		stringAttributeCreateBuilder(d, "file_extension", request.WithFileExtension),
		booleanStringAttributeCreateBuilder(d, "parse_header", request.WithParseHeader),
		intAttributeWithSpecialDefaultCreateBuilder(d, "skip_header", request.WithSkipHeader),
		booleanStringAttributeCreateBuilder(d, "skip_blank_lines", request.WithSkipBlankLines),
		attributeMappedValueCreateBuilder(d, "date_format", request.WithDateFormat, fileFormatStringOrAutoMapper),
		attributeMappedValueCreateBuilder(d, "time_format", request.WithTimeFormat, fileFormatStringOrAutoMapper),
		attributeMappedValueCreateBuilder(d, "timestamp_format", request.WithTimestampFormat, fileFormatStringOrAutoMapper),
		attributeMappedValueCreateBuilder(d, "binary_format", request.WithBinaryFormat, sdk.ToBinaryFormat),
		attributeMappedValueCreateBuilder(d, "escape", request.WithEscape, fileFormatStringOrNoneMapper),
		attributeMappedValueCreateBuilder(d, "escape_unenclosed_field", request.WithEscapeUnenclosedField, fileFormatStringOrNoneMapper),
		booleanStringAttributeCreateBuilder(d, "trim_space", request.WithTrimSpace),
		attributeMappedValueCreateBuilder(d, "field_optionally_enclosed_by", request.WithFieldOptionallyEnclosedBy, fileFormatStringOrNoneMapper),
		// Here we must use RawConfig because the default value for null_if is not an empty list.
		attributeMappedValueCreateBuilderRawConfig(d, "null_if", request.WithNullIf, parseNullIfRequest),
		booleanStringAttributeCreateBuilder(d, "error_on_column_count_mismatch", request.WithErrorOnColumnCountMismatch),
		booleanStringAttributeCreateBuilder(d, "replace_invalid_characters", request.WithReplaceInvalidCharacters),
		booleanStringAttributeCreateBuilder(d, "empty_field_as_null", request.WithEmptyFieldAsNull),
		booleanStringAttributeCreateBuilder(d, "skip_byte_order_mark", request.WithSkipByteOrderMark),
		attributeMappedValueCreateBuilder(d, "encoding", request.WithEncoding, sdk.ToCsvEncoding),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.FileFormats.CreateCsv(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return GetReadFileFormatCsvFunc(false)(ctx, d, meta)
}

func GetReadFileFormatCsvFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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

		details, err := client.FileFormats.DescribeCsvDetails(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		describeOutputValues := schemas.FileFormatCsvToSchema(details)

		if withExternalChangesMarking {
			valuesToSet := csvFileFormatToSchema(details, false)
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

func UpdateFileFormatCsv(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		if err := client.FileFormats.AlterCsv(ctx, sdk.NewAlterCsvFileFormatRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming file format: %w", err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewAlterCsvFileFormatSetRequest()
	fileFormatStringOrAutoFallback := *sdk.NewStageFileFormatStringOrAutoRequest().WithAuto(true)
	fileFormatStringOrNoneFallback := *sdk.NewStageFileFormatStringOrNoneRequest().WithNone(true)

	errs := errors.Join(
		attributeMappedValueUpdateSetOnlyFallback(d, "compression", &set.Compression, sdk.ToCsvCompression, sdk.CsvCompressionAuto),
		attributeMappedValueUpdateSetOnlyFallback(d, "record_delimiter", &set.RecordDelimiter, fileFormatStringOrNoneMapper, *sdk.NewStageFileFormatStringOrNoneRequest().WithValue(`\n`)),
		attributeMappedValueUpdateSetOnlyFallback(d, "field_delimiter", &set.FieldDelimiter, fileFormatStringOrNoneMapper, *sdk.NewStageFileFormatStringOrNoneRequest().WithValue(",")),
		booleanStringAttributeUnsetFallbackUpdate(d, "multi_line", &set.MultiLine, true),
		stringAttributeUpdateSetOnlyNotEmpty(d, "file_extension", &set.FileExtension),
		intAttributeWithSpecialDefaultUnsetFallbackUpdate(d, "skip_header", &set.SkipHeader, 0),
		booleanStringAttributeUnsetFallbackUpdate(d, "skip_blank_lines", &set.SkipBlankLines, false),
		attributeMappedValueUpdateSetOnlyFallback(d, "date_format", &set.DateFormat, fileFormatStringOrAutoMapper, fileFormatStringOrAutoFallback),
		attributeMappedValueUpdateSetOnlyFallback(d, "time_format", &set.TimeFormat, fileFormatStringOrAutoMapper, fileFormatStringOrAutoFallback),
		attributeMappedValueUpdateSetOnlyFallback(d, "timestamp_format", &set.TimestampFormat, fileFormatStringOrAutoMapper, fileFormatStringOrAutoFallback),
		attributeMappedValueUpdateSetOnlyFallback(d, "binary_format", &set.BinaryFormat, sdk.ToBinaryFormat, sdk.BinaryFormatHex),
		attributeMappedValueUpdateSetOnlyFallback(d, "escape", &set.Escape, fileFormatStringOrNoneMapper, fileFormatStringOrNoneFallback),
		attributeMappedValueUpdateSetOnlyFallback(d, "escape_unenclosed_field", &set.EscapeUnenclosedField, fileFormatStringOrNoneMapper, *sdk.NewStageFileFormatStringOrNoneRequest().WithValue(`\\`)),
		booleanStringAttributeUnsetFallbackUpdate(d, "trim_space", &set.TrimSpace, false),
		attributeMappedValueUpdateSetOnlyFallback(d, "field_optionally_enclosed_by", &set.FieldOptionallyEnclosedBy, fileFormatStringOrNoneMapper, fileFormatStringOrNoneFallback),
		attributeMappedValueUpdateSetOnlyRawConfigFallback(d, "null_if", &set.NullIf, parseNullIfRequest, *sdk.NewNullIfListRequest().WithNullIf([]sdk.NullString{{S: `\N`}})),
		booleanStringAttributeUnsetFallbackUpdate(d, "error_on_column_count_mismatch", &set.ErrorOnColumnCountMismatch, true),
		booleanStringAttributeUnsetFallbackUpdate(d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "empty_field_as_null", &set.EmptyFieldAsNull, true),
		booleanStringAttributeUnsetFallbackUpdate(d, "skip_byte_order_mark", &set.SkipByteOrderMark, true),
		attributeMappedValueUpdateSetOnlyFallback(d, "encoding", &set.Encoding, sdk.ToCsvEncoding, sdk.CsvEncodingUtf8),
		stringAttributeUpdateSetOnlyNotEmpty(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(set, sdk.NewAlterCsvFileFormatSetRequest()) {
		if err := client.FileFormats.AlterCsv(ctx, sdk.NewAlterCsvFileFormatRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	set = sdk.NewAlterCsvFileFormatSetRequest()
	if err := booleanStringAttributeUnsetFallbackUpdate(d, "parse_header", &set.ParseHeader, false); err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(set, sdk.NewAlterCsvFileFormatSetRequest()) {
		if err := client.FileFormats.AlterCsv(ctx, sdk.NewAlterCsvFileFormatRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadFileFormatCsvFunc(false)(ctx, d, meta)
}

// csvFileFormatToSchema converts the SDK details for a CSV file format to a Terraform schema.
func csvFileFormatToSchema(csv *sdk.FileFormatCsv, setDefaults bool) map[string]any {
	state := map[string]any{
		"record_delimiter":             csv.RecordDelimiter,
		"field_delimiter":              csv.FieldDelimiter,
		"file_extension":               csv.FileExtension,
		"skip_header":                  csv.SkipHeader,
		"date_format":                  csv.DateFormat,
		"time_format":                  csv.TimeFormat,
		"timestamp_format":             csv.TimestampFormat,
		"binary_format":                string(csv.BinaryFormat),
		"escape":                       csv.Escape,
		"escape_unenclosed_field":      csv.EscapeUnenclosedField,
		"field_optionally_enclosed_by": csv.FieldOptionallyEnclosedBy,
		"null_if":                      collections.Map(csv.NullIf, func(v string) any { return v }),
		"compression":                  string(csv.Compression),
		"encoding":                     string(csv.Encoding),
	}
	if setDefaults {
		state["parse_header"] = BooleanDefault
		state["trim_space"] = BooleanDefault
		state["error_on_column_count_mismatch"] = BooleanDefault
		state["skip_blank_lines"] = BooleanDefault
		state["replace_invalid_characters"] = BooleanDefault
		state["empty_field_as_null"] = BooleanDefault
		state["skip_byte_order_mark"] = BooleanDefault
		state["multi_line"] = BooleanDefault
	} else {
		state["parse_header"] = booleanStringFromBool(csv.ParseHeader)
		state["trim_space"] = booleanStringFromBool(csv.TrimSpace)
		state["error_on_column_count_mismatch"] = booleanStringFromBool(csv.ErrorOnColumnCountMismatch)
		state["skip_blank_lines"] = booleanStringFromBool(csv.SkipBlankLines)
		state["replace_invalid_characters"] = booleanStringFromBool(csv.ReplaceInvalidCharacters)
		state["empty_field_as_null"] = booleanStringFromBool(csv.EmptyFieldAsNull)
		state["skip_byte_order_mark"] = booleanStringFromBool(csv.SkipByteOrderMark)
		state["multi_line"] = booleanStringFromBool(csv.MultiLine)
	}
	return state
}
