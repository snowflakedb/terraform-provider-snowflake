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

var externalAzureStageSchema = func() map[string]*schema.Schema {
	azureStage := map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the URL for the Azure storage container (e.g., 'azure://account.blob.core.windows.net/container').",
		},
		"storage_integration": {
			Type:             schema.TypeString,
			Optional:         true,
			ConflictsWith:    []string{"use_privatelink_endpoint", "credentials"},
			Description:      "Specifies the name of the storage integration used to delegate authentication responsibility to a Snowflake identity.",
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"credentials": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"storage_integration"},
			Description:   "Specifies the Azure SAS token credentials for the external stage.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"azure_sas_token": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "Specifies the shared access signature (SAS) token for Azure.",
					},
				},
			},
		},
		"encryption": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specifies the encryption settings for the Azure external stage.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"azure_cse": {
						Type:         schema.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.azure_cse", "encryption.0.none"},
						Description:  "Azure client-side encryption using a master key.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"master_key": {
									Type:        schema.TypeString,
									Required:    true,
									Sensitive:   true,
									Description: "Specifies the 128-bit or 256-bit client-side master key.",
								},
							},
						},
					},
					"none": {
						Type:         schema.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.azure_cse", "encryption.0.none"},
						Description:  "No encryption.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{},
						},
					},
				},
			},
		},
		"use_privatelink_endpoint": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			ConflictsWith:    []string{"storage_integration"},
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("use_privatelink_endpoint"),
			Description:      "Specifies whether to use a private link endpoint for Azure storage.",
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
						Description:      "Specifies whether to enable a directory table on the external stage.",
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("directory_table.0.enable"),
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
						ForceNew:         true,
						Description:      "Specifies whether Snowflake should enable triggering automatic refreshes of the directory table metadata.",
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("directory_table.0.auto_refresh"),
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
	return collections.MergeMaps(stageCommonSchema, azureStage)
}()

func ExternalAzureStage() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ExternalAzureStageResource), TrackingCreateWrapper(resources.ExternalAzureStage, CreateExternalAzureStage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ExternalAzureStageResource), TrackingReadWrapper(resources.ExternalAzureStage, ReadExternalAzureStageFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ExternalAzureStageResource), TrackingUpdateWrapper(resources.ExternalAzureStage, UpdateExternalAzureStage)),
		DeleteContext: DeleteStage(previewfeatures.ExternalAzureStageResource, resources.ExternalAzureStage),
		Description:   "Resource used to manage external Azure stages. For more information, check [external stage documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stage#external-stage-parameters-externalstageparams).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ExternalAzureStage, customdiff.All(
			ComputedIfAnyAttributeChanged(externalAzureStageSchema, ShowOutputAttributeName, "name", "comment", "url", "storage_integration", "encryption"),
			ComputedIfAnyAttributeChanged(externalAzureStageSchema, DescribeOutputAttributeName, "directory", "encryption"),
			ComputedIfAnyAttributeChanged(externalAzureStageSchema, FullyQualifiedNameAttributeName, "name"),
			ForceNewIfChangeToEmptySlice[any]("directory"),
			ForceNewIfChangeToEmptySlice[any]("credentials"),
			ForceNewIfChangeToEmptySlice[any]("encryption"),
			RecreateWhenStageTypeChangedExternally(sdk.StageTypeExternal),
			RecreateWhenStageCloudChangedExternally(sdk.StageCloudAzure),
		)),

		Schema: externalAzureStageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ExternalAzureStage, ImportExternalAzureStage),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportExternalAzureStage(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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
	if details.PrivateLink != nil {
		if err := d.Set("use_privatelink_endpoint", booleanStringFromBool(details.PrivateLink.UsePrivatelinkEndpoint)); err != nil {
			return nil, err
		}
	} else {
		if err := d.Set("use_privatelink_endpoint", BooleanFalse); err != nil {
			return nil, err
		}
	}
	if details.Location != nil {
		if err := d.Set("url", details.Location.Url); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func CreateExternalAzureStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	url := d.Get("url").(string)
	externalStageParams := sdk.NewExternalAzureStageParamsRequest(url)

	err := errors.Join(
		attributeMappedValueCreateBuilder(d, "credentials", externalStageParams.WithCredentials, parseAzureStageCredentials),
		attributeMappedValueCreateBuilder(d, "encryption", externalStageParams.WithEncryption, parseAzureStageEncryption),
		booleanStringAttributeCreateBuilder(d, "use_privatelink_endpoint", externalStageParams.WithUsePrivatelinkEndpoint),
		accountObjectIdentifierAttributeCreate(d, "storage_integration", &externalStageParams.StorageIntegration),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateOnAzureStageRequest(id, *externalStageParams)

	err = errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, "directory", request.WithDirectoryTableOptions, parseAzureStageDirectory),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Stages.CreateOnAzure(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadExternalAzureStageFunc(false)(ctx, d, meta)
}

func ReadExternalAzureStageFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
						Summary:  "Failed to query external Azure stage. Marking the resource as removed.",
						Detail:   fmt.Sprintf("External Azure stage id: %s, Err: %s", id.FullyQualifiedName(), err),
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
			var storageIntegrationName string
			if stage.StorageIntegration != nil {
				storageIntegrationName = stage.StorageIntegration.Name()
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"comment", "comment", stage.Comment, stage.Comment, nil},
				outputMapping{"storage_integration", "storage_integration", storageIntegrationName, storageIntegrationName, nil},
			); err != nil {
				return diag.FromErr(err)
			}
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

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.StageToSchema(stage)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{detailsSchema}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("url", stage.Url),
			d.Set("comment", stage.Comment),
			d.Set("stage_type", stage.Type),
		)

		if errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func UpdateExternalAzureStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	id, err = handleStageRename(ctx, client, d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = handleStageDirectoryTable(ctx, client, d, id)
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewAlterExternalAzureStageStageRequest(id)

	needsExternalStageParams := d.HasChanges("url", "storage_integration", "credentials", "encryption", "use_privatelink_endpoint")

	if needsExternalStageParams {
		url := d.Get("url").(string)
		externalStageParams := sdk.NewExternalAzureStageParamsRequest(url)

		err = errors.Join(
			booleanStringAttributeUpdateSetOnly(d, "use_privatelink_endpoint", &externalStageParams.UsePrivatelinkEndpoint),
			accountObjectIdentifierAttributeSetOnly(d, "storage_integration", &externalStageParams.StorageIntegration),
			attributeMappedValueUpdateSetOnly(d, "credentials", &externalStageParams.Credentials, parseAzureStageCredentials),
			attributeMappedValueUpdateSetOnly(d, "encryption", &externalStageParams.Encryption, parseAzureStageEncryption),
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

	if !reflect.DeepEqual(*set, sdk.AlterExternalAzureStageStageRequest{}) {
		if err := client.Stages.AlterExternalAzureStage(ctx, set); err != nil {
			return diag.FromErr(fmt.Errorf("error updating external Azure stage: %w", err))
		}
	}

	return ReadExternalAzureStageFunc(false)(ctx, d, meta)
}

func parseAzureStageCredentials(v any) (sdk.ExternalStageAzureCredentialsRequest, error) {
	credentialsList := v.([]any)
	if len(credentialsList) == 0 {
		return sdk.ExternalStageAzureCredentialsRequest{}, nil
	}
	credentialsConfig := credentialsList[0].(map[string]any)
	sasToken := credentialsConfig["azure_sas_token"].(string)
	return *sdk.NewExternalStageAzureCredentialsRequest(sasToken), nil
}

func parseAzureStageEncryption(v any) (sdk.ExternalStageAzureEncryptionRequest, error) {
	encryptionList := v.([]any)
	if len(encryptionList) == 0 {
		return sdk.ExternalStageAzureEncryptionRequest{}, nil
	}
	encryptionConfig := encryptionList[0].(map[string]any)
	encryptionReq := sdk.NewExternalStageAzureEncryptionRequest()

	if azureCse, ok := encryptionConfig["azure_cse"]; ok {
		if cseList := azureCse.([]any); len(cseList) > 0 {
			cseConfig := cseList[0].(map[string]any)
			masterKey := cseConfig["master_key"].(string)
			encryptionReq.WithAzureCse(*sdk.NewExternalStageAzureEncryptionAzureCseRequest(masterKey))
		}
	}

	if none, ok := encryptionConfig["none"]; ok {
		if noneList := none.([]any); len(noneList) > 0 {
			encryptionReq.WithNone(*sdk.NewExternalStageAzureEncryptionNoneRequest())
		}
	}

	return *encryptionReq, nil
}

func parseAzureStageDirectory(v any) (sdk.ExternalAzureDirectoryTableOptionsRequest, error) {
	directoryList := v.([]any)
	if len(directoryList) == 0 {
		return sdk.ExternalAzureDirectoryTableOptionsRequest{}, nil
	}
	directoryConfig := directoryList[0].(map[string]any)
	directoryReq := sdk.NewExternalAzureDirectoryTableOptionsRequest().WithEnable(directoryConfig["enable"].(bool))

	if v, ok := directoryConfig["refresh_on_create"]; ok && v.(string) != BooleanDefault {
		refreshOnCreateBool, err := booleanStringToBool(v.(string))
		if err != nil {
			return sdk.ExternalAzureDirectoryTableOptionsRequest{}, err
		}
		directoryReq.WithRefreshOnCreate(refreshOnCreateBool)
	}

	if v, ok := directoryConfig["auto_refresh"]; ok && v.(string) != BooleanDefault {
		autoRefreshBool, err := booleanStringToBool(v.(string))
		if err != nil {
			return sdk.ExternalAzureDirectoryTableOptionsRequest{}, err
		}
		directoryReq.WithAutoRefresh(autoRefreshBool)
	}

	if notificationIntegration, ok := directoryConfig["notification_integration"]; ok && notificationIntegration.(string) != "" {
		directoryReq.WithNotificationIntegration(notificationIntegration.(string))
	}

	return *directoryReq, nil
}
