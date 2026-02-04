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

var externalGcsStageSchema = func() map[string]*schema.Schema {
	gcsStage := map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the URL for the GCS bucket (e.g., 'gcs://bucket/path/').",
		},
		"storage_integration": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Specifies the name of the storage integration used to delegate authentication responsibility to a Snowflake identity. GCS stages require a storage integration.",
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"encryption": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specifies the encryption settings for the GCS external stage.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"gcs_sse_kms": {
						Type:         schema.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.gcs_sse_kms", "encryption.0.none"},
						Description:  "GCS server-side encryption using a KMS key.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kms_key_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the KMS-managed key ID.",
								},
							},
						},
					},
					"none": {
						Type:         schema.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.gcs_sse_kms", "encryption.0.none"},
						Description:  "No encryption.",
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
						Description: "Specifies whether to enable a directory table on the external stage.",
					},
					"refresh_on_create": {
						Type:             schema.TypeString,
						Default:          BooleanDefault,
						ValidateDiagFunc: validateBooleanString,
						Optional:         true,
						Description:      "Specifies whether to automatically refresh the directory table metadata once, immediately after the stage is created." + ignoredAfterCreationDescription(),
						DiffSuppressFunc: IgnoreAfterCreation,
					},
					"auto_refresh": {
						Type:             schema.TypeString,
						Default:          BooleanDefault,
						ValidateDiagFunc: validateBooleanString,
						Optional:         true,
						Description:      "Specifies whether Snowflake should enable triggering automatic refreshes of the directory table metadata.",
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput("describe_output.0.directory_table", "auto_refresh"),
					},
					"notification_integration": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						Description:      "Specifies the name of the notification integration used to automatically refresh the directory table metadata.",
						DiffSuppressFunc: suppressIdentifierQuoting,
					},
				},
			},
		},
		"cloud": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Specifies a cloud provider for the stage. This field is used for checking external changes and recreating the resources if needed.",
		},
	}
	return collections.MergeMaps(stageCommonSchema, gcsStage)
}()

func ExternalGcsStage() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ExternalGcsStageResource), TrackingCreateWrapper(resources.ExternalGcsStage, CreateExternalGcsStage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ExternalGcsStageResource), TrackingReadWrapper(resources.ExternalGcsStage, ReadExternalGcsStageFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ExternalGcsStageResource), TrackingUpdateWrapper(resources.ExternalGcsStage, UpdateExternalGcsStage)),
		DeleteContext: DeleteStage(previewfeatures.ExternalGcsStageResource, resources.ExternalGcsStage),
		Description:   "Resource used to manage external GCS stages. For more information, check [external stage documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stage#external-stage-parameters-externalstageparams).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ExternalGcsStage, customdiff.All(
			ComputedIfAnyAttributeChanged(externalGcsStageSchema, ShowOutputAttributeName, "name", "comment", "url", "storage_integration", "encryption"),
			ComputedIfAnyAttributeChanged(externalGcsStageSchema, DescribeOutputAttributeName, "directory.0.enable", "directory.0.auto_refresh", "url"),
			ComputedIfAnyAttributeChanged(externalGcsStageSchema, FullyQualifiedNameAttributeName, "name"),
			ForceNewIfChangeToEmptySlice[any]("directory"),
			ForceNewIfChangeToEmptySlice[any]("encryption"),
			ForceNewIfNotDefault("directory.0.auto_refresh"),
			RecreateWhenStageTypeChangedExternally(sdk.StageTypeExternal),
			RecreateWhenStageCloudChangedExternally(sdk.StageCloudGcp),
		)),

		Schema: externalGcsStageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ExternalGcsStage, ImportExternalGcsStage),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportExternalGcsStage(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	if _, err := ImportName[sdk.SchemaObjectIdentifier](context.Background(), d, nil); err != nil {
		return nil, err
	}
	stage, err := client.Stages.ShowByIDSafely(ctx, id)
	if err != nil {
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
	if details.DirectoryTable != nil {
		if err := d.Set("directory", []map[string]any{
			{
				"enable":       details.DirectoryTable.Enable,
				"auto_refresh": booleanStringFromBool(details.DirectoryTable.AutoRefresh),
			},
		}); err != nil {
			return nil, err
		}
	}
	if details.Location != nil {
		if err := d.Set("url", details.Location.Url); err != nil {
			return nil, err
		}
	}
	if stage.StorageIntegration != nil {
		if err := d.Set("storage_integration", stage.StorageIntegration.Name()); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func CreateExternalGcsStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	url := d.Get("url").(string)
	storageIntegrationRaw := d.Get("storage_integration").(string)
	storageIntegrationId := sdk.NewAccountObjectIdentifier(storageIntegrationRaw)

	externalStageParams := sdk.NewExternalGCSStageParamsRequest(url).
		WithStorageIntegration(storageIntegrationId)

	err := errors.Join(
		attributeMappedValueCreateBuilder(d, "encryption", externalStageParams.WithEncryption, parseGcsStageEncryption),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateOnGCSStageRequest(id, *externalStageParams)

	err = errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, "directory", request.WithDirectoryTableOptions, parseGcsStageDirectory),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Stages.CreateOnGCS(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadExternalGcsStageFunc(false)(ctx, d, meta)
}

func ReadExternalGcsStageFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
						Summary:  "Failed to query external GCS stage. Marking the resource as removed.",
						Detail:   fmt.Sprintf("External GCS stage id: %s, Err: %s", id.FullyQualifiedName(), err),
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
		}

		var cloud string
		if stage.Cloud != nil {
			cloud = string(*stage.Cloud)
		}
		var storageIntegrationName string
		if stage.StorageIntegration != nil {
			storageIntegrationName = stage.StorageIntegration.Name()
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.StageToSchema(stage)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{detailsSchema}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("url", stage.Url),
			d.Set("stage_type", stage.Type),
			d.Set("cloud", cloud),
			d.Set("comment", stage.Comment),
			d.Set("storage_integration", storageIntegrationName),
		)

		if errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func UpdateExternalGcsStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	set := sdk.NewAlterExternalGCSStageStageRequest(id)

	needsExternalStageParams := d.HasChanges("url", "storage_integration", "encryption")

	if needsExternalStageParams {
		url := d.Get("url").(string)
		storageIntegrationRaw := d.Get("storage_integration").(string)
		storageIntegrationId := sdk.NewAccountObjectIdentifier(storageIntegrationRaw)

		externalStageParams := sdk.NewExternalGCSStageParamsRequest(url).
			WithStorageIntegration(storageIntegrationId)

		err = errors.Join(
			attributeMappedValueUpdateSetOnly(d, "encryption", &externalStageParams.Encryption, parseGcsStageEncryption),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		set.WithExternalStageParams(*externalStageParams)
	}

	err = errors.Join(
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*set, *sdk.NewAlterExternalGCSStageStageRequest(id)) {
		if err := client.Stages.AlterExternalGCSStage(ctx, set); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error updating external GCS stage: %w", err))
		}
	}

	return ReadExternalGcsStageFunc(false)(ctx, d, meta)
}

func parseGcsStageEncryption(v any) (sdk.ExternalStageGCSEncryptionRequest, error) {
	encryptionList := v.([]any)
	if len(encryptionList) == 0 {
		return sdk.ExternalStageGCSEncryptionRequest{}, nil
	}
	encryptionConfig := encryptionList[0].(map[string]any)
	encryptionReq := sdk.NewExternalStageGCSEncryptionRequest()

	if gcsSseKms, ok := encryptionConfig["gcs_sse_kms"]; ok {
		if kmsList := gcsSseKms.([]any); len(kmsList) > 0 {
			kmsReq := sdk.NewExternalStageGCSEncryptionGcsSseKmsRequest()
			kmsConfig := kmsList[0].(map[string]any)
			if kmsKeyId, ok := kmsConfig["kms_key_id"]; ok && kmsKeyId.(string) != "" {
				kmsReq.WithKmsKeyId(kmsKeyId.(string))
			}
			encryptionReq.WithGcsSseKms(*kmsReq)
		}
	}

	if none, ok := encryptionConfig["none"]; ok {
		if noneList := none.([]any); len(noneList) > 0 {
			encryptionReq.WithNone(*sdk.NewExternalStageGCSEncryptionNoneRequest())
		}
	}

	return *encryptionReq, nil
}

func parseGcsStageDirectory(v any) (sdk.ExternalGCSDirectoryTableOptionsRequest, error) {
	directoryList := v.([]any)
	if len(directoryList) == 0 {
		return sdk.ExternalGCSDirectoryTableOptionsRequest{}, nil
	}
	directoryConfig := directoryList[0].(map[string]any)
	directoryReq := sdk.NewExternalGCSDirectoryTableOptionsRequest().WithEnable(directoryConfig["enable"].(bool))

	if v, ok := directoryConfig["refresh_on_create"]; ok && v.(string) != BooleanDefault {
		refreshOnCreateBool, err := booleanStringToBool(v.(string))
		if err != nil {
			return sdk.ExternalGCSDirectoryTableOptionsRequest{}, fmt.Errorf("parsing refresh_on_create: %w", err)
		}
		directoryReq.WithRefreshOnCreate(refreshOnCreateBool)
	}

	if v, ok := directoryConfig["auto_refresh"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		autoRefreshBool, err := booleanStringToBool(v.(string))
		if err != nil {
			return sdk.ExternalGCSDirectoryTableOptionsRequest{}, fmt.Errorf("parsing auto_refresh: %w", err)
		}
		directoryReq.WithAutoRefresh(autoRefreshBool)
	}

	if notificationIntegration, ok := directoryConfig["notification_integration"]; ok && notificationIntegration.(string) != "" {
		directoryReq.WithNotificationIntegration(notificationIntegration.(string))
	}

	return *directoryReq, nil
}
