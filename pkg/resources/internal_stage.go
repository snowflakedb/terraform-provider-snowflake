package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var internalStageSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the stage; must be unique for the database and schema in which the stage is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the stage."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the stage."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
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
					Type:             schema.TypeBool,
					Required:         true,
					Description:      "Specifies whether to enable a directory table on the internal named stage.",
					DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("directory_table.0.enable"),
				},
				"auto_refresh": {
					Type:             schema.TypeString,
					Default:          BooleanDefault,
					ValidateDiagFunc: validateBooleanString,
					Optional:         true,
					ForceNew:         true,
					Description:      "Specifies whether Snowflake should automatically refresh the directory table metadata when new or updated data files are available on the internal named stage.",
					DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("directory_table.0.auto_refresh"),
				},
			},
		},
	},
	"stage_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies a type for the stage. This field is used for checking external changes and recreating the resources if needed.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the stage.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW STAGES` for the given stage.",
		Elem: &schema.Resource{
			Schema: schemas.ShowStageSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STAGE` for the given stage.",
		Elem: &schema.Resource{
			Schema: schemas.StageDescribeSchema,
		},
	},
}

func InternalStage() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Stages.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.InternalStageResource), TrackingCreateWrapper(resources.InternalStage, CreateInternalStage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.InternalStageResource), TrackingReadWrapper(resources.InternalStage, ReadInternalStageFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.InternalStageResource), TrackingUpdateWrapper(resources.InternalStage, UpdateInternalStage)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.InternalStageResource), TrackingDeleteWrapper(resources.InternalStage, deleteFunc)),
		Description:   "Resource used to manage internal stages. For more information, check [internal stage documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stage#internal-stage-parameters-internalstageparams).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.InternalStage, customdiff.All(
			ComputedIfAnyAttributeChanged(internalStageSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(internalStageSchema, DescribeOutputAttributeName, "directory"),
			ComputedIfAnyAttributeChanged(internalStageSchema, FullyQualifiedNameAttributeName, "name"),
			ForceNewIfChangeToEmptySlice[any]("directory"),
			RecreateWhenResourceTypeChangedExternally("stage_type", sdk.StageTypeInternal, sdk.ToStageType),
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
	if _, err := ImportName[sdk.SchemaObjectIdentifier](context.Background(), d, nil); err != nil {
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

		if autoRefresh, ok := directoryConfig["auto_refresh"]; ok && autoRefresh.(string) != BooleanDefault {
			autoRefreshBool, err := booleanStringToBool(autoRefresh.(string))
			if err != nil {
				return sdk.InternalDirectoryTableOptionsRequest{}, err
			}
			directoryReq.WithAutoRefresh(autoRefreshBool)
		}
		return *directoryReq, nil
	}
	err := errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, "directory", request.WithDirectoryTableOptions, parseDirectoryTable),
		attributeMappedValueCreateBuilder(d, "encryption", request.WithEncryption, parseEncryption),
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
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"comment", "comment", stage.Comment, stage.Comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
			directoryTable := []any{
				map[string]any{
					"enable":       details.DirectoryTable.Enable,
					"auto_refresh": details.DirectoryTable.AutoRefresh,
				},
			}
			if err = handleExternalChangesToObjectInFlatDescribeDeepEqual(d,
				outputMapping{"directory_table", "directory", directoryTable, directoryTable, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.StageToSchema(stage)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{detailsSchema}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("comment", stage.Comment),
			d.Set("stage_type", stage.Type.Canonical()),
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

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), newName)

		err := client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming stage %v to %v: %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewAlterInternalStageStageRequest(id)
	err = errors.Join(
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*set, sdk.AlterInternalStageStageRequest{}) {
		if err := client.Stages.AlterInternalStage(ctx, set); err != nil {
			return diag.FromErr(fmt.Errorf("error updating stage: %w", err))
		}
	}
	setDirectoryTable := sdk.NewAlterDirectoryTableStageRequest(id)
	parseDirectoryTable := func(value any) (sdk.DirectoryTableSetRequest, error) {
		directoryList := value.([]any)
		if len(directoryList) == 0 {
			return sdk.DirectoryTableSetRequest{}, nil
		}
		directoryConfig := directoryList[0].(map[string]any)
		directoryReq := sdk.NewDirectoryTableSetRequest(directoryConfig["enable"].(bool))
		return *directoryReq, nil
	}
	err = errors.Join(
		attributeMappedValueUpdateSetOnly(d, "directory", &setDirectoryTable.SetDirectory, parseDirectoryTable),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	if !reflect.DeepEqual(setDirectoryTable, sdk.NewAlterDirectoryTableStageRequest(id)) {
		if err := client.Stages.AlterDirectoryTable(ctx, setDirectoryTable); err != nil {
			return diag.FromErr(fmt.Errorf("error updating stage: %w", err))
		}
	}

	return ReadInternalStageFunc(false)(ctx, d, meta)
}
