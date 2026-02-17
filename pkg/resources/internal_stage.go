package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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

var internalStageSchema = func() map[string]*schema.Schema {
	internalStage := map[string]*schema.Schema{
		"encryption": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "Specifies the encryption settings for the internal stage.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"snowflake_full": {
						Type:         schema.TypeList,
						Optional:     true,
						ForceNew:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.snowflake_full", "encryption.0.snowflake_sse"},
						Description:  "Client-side and server-side encryption.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{},
						},
					},
					"snowflake_sse": {
						Type:         schema.TypeList,
						Optional:     true,
						ForceNew:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.snowflake_full", "encryption.0.snowflake_sse"},
						Description:  "Server-side encryption only.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{},
						},
					},
				},
			},
		},
		"directory": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Directory tables store a catalog of staged files in cloud storage.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enable": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Specifies whether to enable a directory table on the internal named stage.",
					},
					"auto_refresh": {
						Type:             schema.TypeString,
						Default:          BooleanDefault,
						ValidateDiagFunc: validateBooleanString,
						Optional:         true,
						Description:      "Specifies whether Snowflake should automatically refresh the directory table metadata when new or updated data files are available on the internal named stage.",
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput("describe_output.0.directory_table", "auto_refresh"),
					},
				},
			},
		},
	}
	return collections.MergeMaps(stageCommonSchema(schemas.CommonStageDescribeSchema()), internalStage)
}()

func InternalStage() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.InternalStageResource), TrackingCreateWrapper(resources.InternalStage, CreateInternalStage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.InternalStageResource), TrackingReadWrapper(resources.InternalStage, ReadInternalStageFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.InternalStageResource), TrackingUpdateWrapper(resources.InternalStage, UpdateInternalStage)),
		DeleteContext: DeleteStage(previewfeatures.InternalStageResource, resources.InternalStage),
		Description:   "Resource used to manage internal stages. For more information, check [internal stage documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stage#internal-stage-parameters-internalstageparams).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.InternalStage, customdiff.All(
			ComputedIfAnyAttributeChanged(internalStageSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(internalStageSchema, DescribeOutputAttributeName, "directory.0.enable", "directory.0.auto_refresh", "file_format"),
			ComputedIfAnyAttributeChanged(internalStageSchema, FullyQualifiedNameAttributeName, "name"),
			ForceNewIfChangeToEmptySlice[any]("directory"),
			ForceNewIfNotDefault("directory.0.auto_refresh"),
			RecreateWhenStageTypeChangedExternally(sdk.StageTypeInternal),
		)),

		Schema: internalStageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.InternalStage, ImportInternalStage),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportInternalStage(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		return nil, err
	}
	stageDetails, err := client.Stages.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	details, err := sdk.ParseStageDetails(stageDetails)
	if err != nil {
		return nil, err
	}
	if err := d.Set("directory", []map[string]any{
		{
			"enable":       details.DirectoryTable.Enable,
			"auto_refresh": booleanStringFromBool(details.DirectoryTable.AutoRefresh),
		},
	}); err != nil {
		return nil, err
	}
	if fileFormat := stageFileFormatToSchema(details); fileFormat != nil {
		if err := d.Set("file_format", fileFormat); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func CreateInternalStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	request := sdk.NewCreateInternalStageRequest(id)

	parseEncryption := func(v any) (sdk.InternalStageEncryptionRequest, error) {
		encryptionList := v.([]any)
		if len(encryptionList) == 0 {
			return sdk.InternalStageEncryptionRequest{}, nil
		}
		encryptionConfig := encryptionList[0].(map[string]any)
		encryptionReq := sdk.NewInternalStageEncryptionRequest()

		if snowflakeFull, ok := encryptionConfig["snowflake_full"]; ok {
			if sfList := snowflakeFull.([]any); len(sfList) > 0 {
				encryptionReq.WithSnowflakeFull(*sdk.NewInternalStageEncryptionSnowflakeFullRequest())
			}
		}

		if snowflakeSse, ok := encryptionConfig["snowflake_sse"]; ok {
			if sseList := snowflakeSse.([]any); len(sseList) > 0 {
				encryptionReq.WithSnowflakeSse(*sdk.NewInternalStageEncryptionSnowflakeSseRequest())
			}
		}

		return *encryptionReq, nil
	}
	parseDirectoryTable := func(value any) (sdk.InternalDirectoryTableOptionsRequest, error) {
		directoryList := value.([]any)
		if len(directoryList) == 0 {
			return sdk.InternalDirectoryTableOptionsRequest{}, nil
		}
		directoryConfig := directoryList[0].(map[string]any)
		directoryReq := sdk.NewInternalDirectoryTableOptionsRequest()
		if enable, ok := directoryConfig["enable"]; ok {
			directoryReq.WithEnable(enable.(bool))
		}

		if autoRefresh, ok := directoryConfig["auto_refresh"]; ok && autoRefresh.(string) != BooleanDefault && autoRefresh.(string) != "" {
			autoRefreshBool, err := booleanStringToBool(autoRefresh.(string))
			if err != nil {
				return sdk.InternalDirectoryTableOptionsRequest{}, fmt.Errorf("parsing auto_refresh: %w", err)
			}
			directoryReq.WithAutoRefresh(autoRefreshBool)
		}
		return *directoryReq, nil
	}
	err := errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, "directory", request.WithDirectoryTableOptions, parseDirectoryTable),
		attributeMappedValueCreateBuilder(d, "encryption", request.WithEncryption, parseEncryption),
		attributeMappedValueCreateBuilderNested(d, "file_format", request.WithFileFormat, parseStageFileFormat),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := client.Stages.CreateInternal(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadInternalStageFunc(false)(ctx, d, meta)
}

func ReadInternalStageFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		stage, err := client.Stages.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query internal stage. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Internal stage id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		properties, err := client.Stages.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		details, err := sdk.ParseStageDetails(properties)
		if err != nil {
			return diag.FromErr(err)
		}

		detailsSchema, err := schemas.StageDescribeToSchema(*details)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			directoryTable := []any{
				map[string]any{
					"enable":       details.DirectoryTable.Enable,
					"auto_refresh": details.DirectoryTable.AutoRefresh,
				},
			}
			directoryTableToSet := []any{
				map[string]any{
					"enable":       details.DirectoryTable.Enable,
					"auto_refresh": booleanStringFromBool(details.DirectoryTable.AutoRefresh),
				},
			}
			if err = handleExternalChangesToObjectInFlatDescribeDeepEqual(d,
				outputMapping{"directory_table", "directory", directoryTable, directoryTableToSet, nil},
			); err != nil {
				return diag.FromErr(err)
			}
			if err := handleStageFileFormatRead(d, details); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.StageToSchema(stage)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{detailsSchema}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("comment", stage.Comment),
			d.Set("stage_type", stage.Type),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func UpdateInternalStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	id, err = handleStageRename(ctx, client, d, id)
	if err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	err = handleStageDirectoryTable(ctx, client, d, id)
	if err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	set := sdk.NewAlterInternalStageStageRequest(id)
	err = errors.Join(
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
		attributeMappedValueUpdateSetOnlyFallbackNested(d, "file_format", &set.FileFormat, parseStageFileFormat, sdk.StageFileFormatRequest{FileFormatOptions: &sdk.FileFormatOptions{CsvOptions: &sdk.FileFormatCsvOptions{}}}),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*set, *sdk.NewAlterInternalStageStageRequest(id)) {
		if err := client.Stages.AlterInternalStage(ctx, set); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error updating stage: %w", err))
		}
	}

	return ReadInternalStageFunc(false)(ctx, d, meta)
}
