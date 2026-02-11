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

var externalS3StageSchema = func() map[string]*schema.Schema {
	s3Stage := map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the URL for the S3 bucket (e.g., 's3://bucket-name/path/').",
		},
		"storage_integration": {
			Type:             schema.TypeString,
			Optional:         true,
			ConflictsWith:    []string{"use_privatelink_endpoint", "credentials"},
			Description:      blocklistedCharactersFieldDescription("Specifies the name of the storage integration used to delegate authentication responsibility to a Snowflake identity."),
			DiffSuppressFunc: suppressIdentifierQuoting,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		"aws_access_point_arn": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies the ARN for an AWS S3 Access Point to use for data transfer.",
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("location.0.aws_access_point_arn"),
		},
		"credentials": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"storage_integration"},
			Description:   "Specifies the AWS credentials for the external stage.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_key_id": {
						Type:          schema.TypeString,
						Optional:      true,
						Sensitive:     true,
						ConflictsWith: []string{"credentials.0.aws_role"},
						Description:   "Specifies the AWS access key ID.",
					},
					"aws_secret_key": {
						Type:          schema.TypeString,
						Optional:      true,
						Sensitive:     true,
						ConflictsWith: []string{"credentials.0.aws_role"},
						Description:   "Specifies the AWS secret access key.",
					},
					"aws_token": {
						Type:          schema.TypeString,
						Optional:      true,
						Sensitive:     true,
						ConflictsWith: []string{"credentials.0.aws_role"},
						Description:   "Specifies the AWS session token for temporary credentials.",
					},
					"aws_role": {
						Type:          schema.TypeString,
						Optional:      true,
						ConflictsWith: []string{"credentials.0.aws_key_id", "credentials.0.aws_secret_key", "credentials.0.aws_token"},
						Description:   "Specifies the AWS IAM role ARN to use for accessing the bucket.",
					},
				},
			},
		},
		"encryption": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specifies the encryption settings for the S3 external stage.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_cse": {
						Type:         schema.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.aws_cse", "encryption.0.aws_sse_s3", "encryption.0.aws_sse_kms", "encryption.0.none"},
						Description:  "AWS client-side encryption using a master key.",
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
					"aws_sse_s3": {
						Type:         schema.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.aws_cse", "encryption.0.aws_sse_s3", "encryption.0.aws_sse_kms", "encryption.0.none"},
						Description:  "AWS server-side encryption using S3-managed keys.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{},
						},
					},
					"aws_sse_kms": {
						Type:         schema.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"encryption.0.aws_cse", "encryption.0.aws_sse_s3", "encryption.0.aws_sse_kms", "encryption.0.none"},
						Description:  "AWS server-side encryption using KMS-managed keys.",
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
						ExactlyOneOf: []string{"encryption.0.aws_cse", "encryption.0.aws_sse_s3", "encryption.0.aws_sse_kms", "encryption.0.none"},
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
			Description:      "Specifies whether to use a private link endpoint for S3 storage.",
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("privatelink.0.use_privatelink_endpoint"),
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
	return collections.MergeMaps(stageCommonSchema(schemas.AwsStageDescribeSchema()), s3Stage)
}()

func ExternalS3Stage() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ExternalS3StageResource), TrackingCreateWrapper(resources.ExternalS3Stage, CreateExternalS3Stage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ExternalS3StageResource), TrackingReadWrapper(resources.ExternalS3Stage, ReadExternalS3StageFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ExternalS3StageResource), TrackingUpdateWrapper(resources.ExternalS3Stage, UpdateExternalS3Stage)),
		DeleteContext: DeleteStage(previewfeatures.ExternalS3StageResource, resources.ExternalS3Stage),
		Description:   "Resource used to manage external S3 stages. For more information, check [external stage documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stage#external-stage-parameters-externalstageparams).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ExternalS3Stage, customdiff.All(
			ComputedIfAnyAttributeChanged(externalS3StageSchema, ShowOutputAttributeName, "name", "comment", "url", "storage_integration", "encryption"),
			ComputedIfAnyAttributeChanged(externalS3StageSchema, DescribeOutputAttributeName, "directory.0.enable", "directory.0.auto_refresh", "url", "use_privatelink_endpoint", "aws_access_point_arn", "file_format"),
			ComputedIfAnyAttributeChanged(externalS3StageSchema, FullyQualifiedNameAttributeName, "name"),
			ForceNewIfChangeToEmptySlice[any]("directory"),
			ForceNewIfChangeToEmptySlice[any]("credentials"),
			ForceNewIfChangeToEmptySlice[any]("encryption"),
			ForceNewIfChangeToEmptyString("storage_integration"),
			ForceNewIfNotDefault("directory.0.auto_refresh"),
			ForceNewIfChangeToEmptyString("aws_access_point_arn"),
			ForceNewIfUrlIsS3Compatible(),
			RecreateWhenStageTypeChangedExternally(sdk.StageTypeExternal),
			RecreateWhenStageCloudChangedExternally(sdk.StageCloudAws),
			ForceNewIfChangeToEmptyString("aws_access_point_arn"),
			// This is a similar configuration as for external S3-compatible stage, but the additional differences are:
			// - endpoint is required for S3-compatible stages, but it's null for S3 stages
			// - url starts with s3compat:// instead of s3://
			// changes on both of these fields trigger ForceNew.
			RecreateWhenCredentialsAndStorageIntegrationChangedOnExternalStage(),
		)),

		Schema: externalS3StageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ExternalS3Stage, ImportExternalS3Stage),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportExternalS3Stage(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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
	if fileFormat := stageFileFormatToSchema(details); fileFormat != nil {
		if err := d.Set("file_format", fileFormat); err != nil {
			return nil, err
		}
	}
	if details.PrivateLink != nil {
		if err := d.Set("use_privatelink_endpoint", booleanStringFromBool(details.PrivateLink.UsePrivatelinkEndpoint)); err != nil {
			return nil, err
		}
	}
	if details.Location != nil {
		if err := d.Set("url", details.Location.Url); err != nil {
			return nil, err
		}
		if details.Location.AwsAccessPointArn != "" {
			if err := d.Set("aws_access_point_arn", details.Location.AwsAccessPointArn); err != nil {
				return nil, err
			}
		}
	}
	if stage.StorageIntegration != nil {
		if err := d.Set("storage_integration", stage.StorageIntegration.Name()); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func CreateExternalS3Stage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	url := d.Get("url").(string)
	externalStageParams := sdk.NewExternalS3StageParamsRequest(url)

	err := errors.Join(
		attributeMappedValueCreateBuilder(d, "credentials", externalStageParams.WithCredentials, parseS3StageCredentials),
		attributeMappedValueCreateBuilder(d, "encryption", externalStageParams.WithEncryption, parseS3StageEncryption),
		booleanStringAttributeCreateBuilder(d, "use_privatelink_endpoint", externalStageParams.WithUsePrivatelinkEndpoint),
		accountObjectIdentifierAttributeCreate(d, "storage_integration", &externalStageParams.StorageIntegration),
		stringAttributeCreate(d, "aws_access_point_arn", &externalStageParams.AwsAccessPointArn),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateOnS3StageRequest(id, *externalStageParams)

	err = errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, "directory", request.WithDirectoryTableOptions, parseS3StageDirectory),
		attributeMappedValueCreateBuilderNested(d, "file_format", request.WithFileFormat, parseStageFileFormat),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Stages.CreateOnS3(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadExternalS3StageFunc(false)(ctx, d, meta)
}

func ReadExternalS3StageFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
						Summary:  "Failed to query external S3 stage. Marking the resource as removed.",
						Detail:   fmt.Sprintf("External S3 stage id: %s, Err: %s", id.FullyQualifiedName(), err),
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

		detailsSchema, err := schemas.AwsStageDescribeToSchema(*details)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			var storageIntegrationName string
			if stage.StorageIntegration != nil {
				storageIntegrationName = stage.StorageIntegration.Name()
			}
			if err = handleExternalChangesToObjectInShow(d,
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
				outputMapping{"location.0.aws_access_point_arn", "aws_access_point_arn", details.Location.AwsAccessPointArn, details.Location.AwsAccessPointArn, nil},
			); err != nil {
				return diag.FromErr(err)
			}
			if err := handleStageFileFormatRead(d, details); err != nil {
				return diag.FromErr(err)
			}
			var usePrivatelinkEndpoint bool
			if details.PrivateLink != nil {
				usePrivatelinkEndpoint = details.PrivateLink.UsePrivatelinkEndpoint
			}
			if err = handleExternalChangesToObject(d,
				"describe_output.0.privatelink",
				outputMapping{"use_privatelink_endpoint", "use_privatelink_endpoint", usePrivatelinkEndpoint, booleanStringFromBool(usePrivatelinkEndpoint), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		var cloud string
		if stage.Cloud != nil {
			cloud = string(*stage.Cloud)
		}
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.StageToSchema(stage)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{detailsSchema}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("url", stage.Url),
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

func UpdateExternalS3Stage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	set := sdk.NewAlterExternalS3StageStageRequest(id)

	needsExternalStageParams := d.HasChanges("url", "storage_integration", "credentials", "encryption", "use_privatelink_endpoint", "aws_access_point_arn")

	if needsExternalStageParams {
		url := d.Get("url").(string)
		externalStageParams := sdk.NewExternalS3StageParamsRequest(url)

		err = errors.Join(
			booleanStringAttributeUpdateSetOnly(d, "use_privatelink_endpoint", &externalStageParams.UsePrivatelinkEndpoint),
			accountObjectIdentifierAttributeSetOnly(d, "storage_integration", &externalStageParams.StorageIntegration),
			attributeMappedValueUpdateSetOnly(d, "credentials", &externalStageParams.Credentials, parseS3StageCredentials),
			attributeMappedValueUpdateSetOnly(d, "encryption", &externalStageParams.Encryption, parseS3StageEncryption),
			attributeMappedValueUpdateSetOnly(d, "aws_access_point_arn", &externalStageParams.AwsAccessPointArn, identityMapping),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		set.WithExternalStageParams(*externalStageParams)
	}

	err = errors.Join(
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
		attributeMappedValueUpdateSetOnlyFallbackNested(d, "file_format", &set.FileFormat, parseStageFileFormat, sdk.StageFileFormatRequest{FileFormatOptions: &sdk.FileFormatOptions{CsvOptions: &sdk.FileFormatCsvOptions{}}}),
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

	return ReadExternalS3StageFunc(false)(ctx, d, meta)
}

func parseS3StageCredentials(v any) (sdk.ExternalStageS3CredentialsRequest, error) {
	credentialsList := v.([]any)
	if len(credentialsList) == 0 {
		return sdk.ExternalStageS3CredentialsRequest{}, nil
	}
	credentialsConfig := credentialsList[0].(map[string]any)
	credentialsReq := sdk.NewExternalStageS3CredentialsRequest()

	if awsKeyId, ok := credentialsConfig["aws_key_id"]; ok && awsKeyId.(string) != "" {
		credentialsReq.WithAwsKeyId(awsKeyId.(string))
	}
	if awsSecretKey, ok := credentialsConfig["aws_secret_key"]; ok && awsSecretKey.(string) != "" {
		credentialsReq.WithAwsSecretKey(awsSecretKey.(string))
	}
	if awsToken, ok := credentialsConfig["aws_token"]; ok && awsToken.(string) != "" {
		credentialsReq.WithAwsToken(awsToken.(string))
	}
	if awsRole, ok := credentialsConfig["aws_role"]; ok && awsRole.(string) != "" {
		credentialsReq.WithAwsRole(awsRole.(string))
	}

	return *credentialsReq, nil
}

func parseS3StageEncryption(v any) (sdk.ExternalStageS3EncryptionRequest, error) {
	encryptionList := v.([]any)
	if len(encryptionList) == 0 {
		return sdk.ExternalStageS3EncryptionRequest{}, nil
	}
	encryptionConfig := encryptionList[0].(map[string]any)
	encryptionReq := sdk.NewExternalStageS3EncryptionRequest()

	if awsCse, ok := encryptionConfig["aws_cse"]; ok {
		if cseList := awsCse.([]any); len(cseList) > 0 {
			cseConfig := cseList[0].(map[string]any)
			masterKey := cseConfig["master_key"].(string)
			encryptionReq.WithAwsCse(*sdk.NewExternalStageS3EncryptionAwsCseRequest(masterKey))
		}
	}

	if awsSseS3, ok := encryptionConfig["aws_sse_s3"]; ok {
		if sseS3List := awsSseS3.([]any); len(sseS3List) > 0 {
			encryptionReq.WithAwsSseS3(*sdk.NewExternalStageS3EncryptionAwsSseS3Request())
		}
	}

	if awsSseKms, ok := encryptionConfig["aws_sse_kms"]; ok {
		if sseKmsList := awsSseKms.([]any); len(sseKmsList) > 0 {
			kmsReq := sdk.NewExternalStageS3EncryptionAwsSseKmsRequest()
			kmsConfig := sseKmsList[0].(map[string]any)
			if kmsKeyId, ok := kmsConfig["kms_key_id"]; ok && kmsKeyId.(string) != "" {
				kmsReq.WithKmsKeyId(kmsKeyId.(string))
			}
			encryptionReq.WithAwsSseKms(*kmsReq)
		}
	}

	if none, ok := encryptionConfig["none"]; ok {
		if noneList := none.([]any); len(noneList) > 0 {
			encryptionReq.WithNone(*sdk.NewExternalStageS3EncryptionNoneRequest())
		}
	}

	return *encryptionReq, nil
}

func parseS3StageDirectory(v any) (sdk.StageS3CommonDirectoryTableOptionsRequest, error) {
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
