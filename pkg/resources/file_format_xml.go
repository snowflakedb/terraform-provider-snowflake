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

func fileFormatXmlSchema() map[string]*schema.Schema {
	return collections.MergeMaps(fileFormatCommonSchema, xmlFileFormatFieldsSchema(""), xmlDescOutputSchema())
}

func xmlDescOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE FILE FORMAT` for this file format.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeFileFormatXmlSchema,
			},
		},
	}
}

func FileFormatXml() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.FileFormats.DropSafely
		},
	)

	resourceSchema := fileFormatXmlSchema()

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FileFormatXmlResource), TrackingCreateWrapper(resources.FileFormatXml, CreateFileFormatXml)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FileFormatXmlResource), TrackingReadWrapper(resources.FileFormatXml, GetReadFileFormatXmlFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FileFormatXmlResource), TrackingUpdateWrapper(resources.FileFormatXml, UpdateFileFormatXml)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FileFormatXmlResource), TrackingDeleteWrapper(resources.FileFormatXml, deleteFunc)),
		Description:   "Resource used to manage XML file formats. For more information, check [file format documentation](https://docs.snowflake.com/en/sql-reference/sql/create-file-format).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FileFormatXml, customdiff.All(
			ComputedIfAnyAttributeChanged(resourceSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(
				resourceSchema, DescribeOutputAttributeName,
				"name", "type", "compression", "ignore_utf8_errors", "preserve_space", "strip_outer_element",
				"disable_snowflake_data", "disable_auto_convert", "replace_invalid_characters", "skip_byte_order_mark",
			),
			ComputedIfAnyAttributeChanged(resourceSchema, FullyQualifiedNameAttributeName, "name"),
			RecreateWhenResourceTypeChangedExternally("type", sdk.FileFormatTypeXml, sdk.ToFileFormatType),
		)),

		Schema: resourceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.FileFormatXml, ImportFileFormatXml),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportFileFormatXml(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.FileFormats.DescribeXmlDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.Type != string(sdk.FileFormatTypeXml) {
		return nil, fmt.Errorf("invalid file format type, expected %s, got %s", sdk.FileFormatTypeXml, details.Type)
	}

	var errs []error
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		errs = append(errs, err)
	}

	for key, value := range xmlFileFormatToSchema(details, true) {
		errs = append(errs, d.Set(key, value))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return []*schema.ResourceData{d}, nil
}

func CreateFileFormatXml(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateXmlFileFormatRequest(id)

	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "compression", request.WithCompression, sdk.ToXmlCompression),
		booleanStringAttributeCreateBuilder(d, "ignore_utf8_errors", request.WithIgnoreUtf8Errors),
		booleanStringAttributeCreateBuilder(d, "preserve_space", request.WithPreserveSpace),
		booleanStringAttributeCreateBuilder(d, "strip_outer_element", request.WithStripOuterElement),
		booleanStringAttributeCreateBuilder(d, "disable_snowflake_data", request.WithDisableSnowflakeData),
		booleanStringAttributeCreateBuilder(d, "disable_auto_convert", request.WithDisableAutoConvert),
		booleanStringAttributeCreateBuilder(d, "replace_invalid_characters", request.WithReplaceInvalidCharacters),
		booleanStringAttributeCreateBuilder(d, "skip_byte_order_mark", request.WithSkipByteOrderMark),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.FileFormats.CreateXml(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return GetReadFileFormatXmlFunc(false)(ctx, d, meta)
}

func GetReadFileFormatXmlFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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

		details, err := client.FileFormats.DescribeXmlDetails(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		describeOutputValues := schemas.FileFormatXmlToSchema(details)

		if withExternalChangesMarking {
			valuesToSet := xmlFileFormatToSchema(details, false)
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

func UpdateFileFormatXml(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		if err := client.FileFormats.AlterXml(ctx, sdk.NewAlterXmlFileFormatRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming file format: %w", err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewAlterXmlFileFormatSetRequest()

	errs := errors.Join(
		attributeMappedValueUpdateSetOnlyFallback(d, "compression", &set.Compression, sdk.ToXmlCompression, sdk.XmlCompressionAuto),
		booleanStringAttributeUnsetFallbackUpdate(d, "ignore_utf8_errors", &set.IgnoreUtf8Errors, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "preserve_space", &set.PreserveSpace, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "strip_outer_element", &set.StripOuterElement, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "disable_snowflake_data", &set.DisableSnowflakeData, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "disable_auto_convert", &set.DisableAutoConvert, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, false),
		booleanStringAttributeUnsetFallbackUpdate(d, "skip_byte_order_mark", &set.SkipByteOrderMark, true),
		stringAttributeUpdateSetOnlyNotEmpty(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(set, sdk.NewAlterXmlFileFormatSetRequest()) {
		if err := client.FileFormats.AlterXml(ctx, sdk.NewAlterXmlFileFormatRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadFileFormatXmlFunc(false)(ctx, d, meta)
}

// xmlFileFormatToSchema converts the SDK details for an XML file format to a Terraform schema.
func xmlFileFormatToSchema(xml *sdk.FileFormatXml, setDefaults bool) map[string]any {
	state := map[string]any{
		"compression": xml.Compression,
	}
	if setDefaults {
		state["ignore_utf8_errors"] = BooleanDefault
		state["preserve_space"] = BooleanDefault
		state["strip_outer_element"] = BooleanDefault
		state["disable_snowflake_data"] = BooleanDefault
		state["disable_auto_convert"] = BooleanDefault
		state["replace_invalid_characters"] = BooleanDefault
		state["skip_byte_order_mark"] = BooleanDefault
	} else {
		state["ignore_utf8_errors"] = booleanStringFromBool(xml.IgnoreUtf8Errors)
		state["preserve_space"] = booleanStringFromBool(xml.PreserveSpace)
		state["strip_outer_element"] = booleanStringFromBool(xml.StripOuterElement)
		state["disable_snowflake_data"] = booleanStringFromBool(xml.DisableSnowflakeData)
		state["disable_auto_convert"] = booleanStringFromBool(xml.DisableAutoConvert)
		state["replace_invalid_characters"] = booleanStringFromBool(xml.ReplaceInvalidCharacters)
		state["skip_byte_order_mark"] = booleanStringFromBool(xml.SkipByteOrderMark)
	}
	return state
}
