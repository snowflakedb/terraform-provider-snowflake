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

var externalS3CompatStageSchema = func() map[string]*schema.Schema {
	s3CompatStage := map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specifies the URL for the S3-compatible storage location (e.g., 's3compat://bucket/path/').",
		},
		"endpoint": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specifies the endpoint for the S3-compatible storage provider.",
		},
		"credentials": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "Specifies the AWS credentials for the S3-compatible external stage.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_key_id": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "Specifies the AWS access key ID.",
					},
					"aws_secret_key": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						Description: "Specifies the AWS secret access key.",
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
				},
			},
		},
		"cloud": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Specifies a cloud provider for the stage. This field is used for checking external changes and recreating the resources if needed.",
		},
	}
	return collections.MergeMaps(stageCommonSchema, s3CompatStage)
}()

func ExternalS3CompatibleStage() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ExternalS3CompatStageResource), TrackingCreateWrapper(resources.ExternalS3CompatibleStage, CreateExternalS3CompatStage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ExternalS3CompatStageResource), TrackingReadWrapper(resources.ExternalS3CompatibleStage, ReadExternalS3CompatStageFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ExternalS3CompatStageResource), TrackingUpdateWrapper(resources.ExternalS3CompatibleStage, UpdateExternalS3CompatStage)),
		DeleteContext: DeleteStage(previewfeatures.ExternalS3CompatStageResource, resources.ExternalS3CompatibleStage),
		Description:   "Resource used to manage external S3-compatible stages. For more information, check [external stage documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stage#external-stage-parameters-externalstageparams).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ExternalS3CompatibleStage, customdiff.All(
			ComputedIfAnyAttributeChanged(externalS3CompatStageSchema, ShowOutputAttributeName, "name", "comment", "url", "endpoint"),
			ComputedIfAnyAttributeChanged(externalS3CompatStageSchema, DescribeOutputAttributeName, "directory.0.enable", "directory.0.auto_refresh", "url"),
			ComputedIfAnyAttributeChanged(externalS3CompatStageSchema, FullyQualifiedNameAttributeName, "name"),
			ForceNewIfChangeToEmptySlice[any]("directory"),
			ForceNewIfChangeToEmptySlice[any]("credentials"),
			ForceNewIfNotDefault("directory.0.auto_refresh"),
			RecreateWhenStageTypeChangedExternally(sdk.StageTypeExternal),
			RecreateWhenStageCloudChangedExternally(sdk.StageCloudAws),
		)),

		Schema: externalS3CompatStageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ExternalS3CompatibleStage, ImportExternalS3CompatStage),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportExternalS3CompatStage(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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
	if stage.Endpoint != nil {
		if err := d.Set("endpoint", *stage.Endpoint); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func CreateExternalS3CompatStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	url := d.Get("url").(string)
	endpoint := d.Get("endpoint").(string)
	externalStageParams := sdk.NewExternalS3CompatibleStageParamsRequest(url, endpoint)

	err := attributeMappedValueCreateBuilder(d, "credentials", externalStageParams.WithCredentials, parseS3CompatStageCredentials)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateOnS3CompatibleStageRequest(id, *externalStageParams)

	err = errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, "directory", request.WithDirectoryTableOptions, parseS3CompatStageDirectory),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Stages.CreateOnS3Compatible(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadExternalS3CompatStageFunc(false)(ctx, d, meta)
}

func ReadExternalS3CompatStageFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
						Summary:  "Failed to query external S3-compatible stage. Marking the resource as removed.",
						Detail:   fmt.Sprintf("External S3-compatible stage id: %s, Err: %s", id.FullyQualifiedName(), err),
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

		var endpoint string
		if stage.Endpoint != nil {
			endpoint = *stage.Endpoint
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.StageToSchema(stage)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{detailsSchema}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("url", stage.Url),
			d.Set("endpoint", endpoint),
			d.Set("stage_type", stage.Type),
			d.Set("cloud", cloud),
			d.Set("comment", stage.Comment),
		)

		if errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func UpdateExternalS3CompatStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	// Note: url, endpoint, and credentials changes trigger ForceNew due to schema definitions,
	// so they won't reach here. Only name and directory.enable can be updated in-place.
	// For now, s3 compatible stage does not have dedicated syntax in the docs, but we can reuse the comment syntax.
	set := sdk.NewAlterExternalS3StageStageRequest(id)
	err = errors.Join(
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*set, *sdk.NewAlterExternalS3StageStageRequest(id)) {
		if err := client.Stages.AlterExternalS3Stage(ctx, set); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error updating external S3 stage: %w", err))
		}
	}

	return ReadExternalS3CompatStageFunc(false)(ctx, d, meta)
}

func parseS3CompatStageCredentials(v any) (sdk.ExternalStageS3CompatibleCredentialsRequest, error) {
	credentialsList := v.([]any)
	if len(credentialsList) == 0 {
		return sdk.ExternalStageS3CompatibleCredentialsRequest{}, nil
	}
	credentialsConfig := credentialsList[0].(map[string]any)
	awsKeyId := credentialsConfig["aws_key_id"].(string)
	awsSecretKey := credentialsConfig["aws_secret_key"].(string)
	return *sdk.NewExternalStageS3CompatibleCredentialsRequest(awsKeyId, awsSecretKey), nil
}

// TODO: rm and use s3
func parseS3CompatStageDirectory(v any) (sdk.StageS3CommonDirectoryTableOptionsRequest, error) {
	directoryList := v.([]any)
	if len(directoryList) == 0 {
		return sdk.StageS3CommonDirectoryTableOptionsRequest{}, nil
	}
	directoryConfig := directoryList[0].(map[string]any)
	directoryReq := sdk.NewStageS3CommonDirectoryTableOptionsRequest().WithEnable(directoryConfig["enable"].(bool))

	if v, ok := directoryConfig["refresh_on_create"]; ok && v.(string) != BooleanDefault {
		refreshOnCreateBool, err := booleanStringToBool(v.(string))
		if err != nil {
			return sdk.StageS3CommonDirectoryTableOptionsRequest{}, fmt.Errorf("parsing refresh_on_create: %w", err)
		}
		directoryReq.WithRefreshOnCreate(refreshOnCreateBool)
	}

	if v, ok := directoryConfig["auto_refresh"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		autoRefreshBool, err := booleanStringToBool(v.(string))
		if err != nil {
			return sdk.StageS3CommonDirectoryTableOptionsRequest{}, fmt.Errorf("parsing auto_refresh: %w", err)
		}
		directoryReq.WithAutoRefresh(autoRefreshBool)
	}

	return *directoryReq, nil
}
