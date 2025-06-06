package resources

import (
	"context"
	"errors"
	"fmt"

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

var externalVolumeSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		ForceNew:         true,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the external volume; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	// A list is used as the order of storage locations matter. Storage location position in the list is used to select
	// the active storage location - https://docs.snowflake.com/en/user-guide/tables-iceberg-storage#active-storage-location
	// This is also why it has been left as one list with optional cloud dependent parameters, rather than splitting into
	// one list per cloud provider.
	"storage_location": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Description: "List of named cloud storage locations in different regions and, optionally, cloud platforms. Minimum 1 required. The order of the list is important as it impacts the active storage location, and updates will be triggered if it changes. Note that not all parameter combinations are valid as they depend on the given storage_provider. Consult [the docs](https://docs.snowflake.com/en/sql-reference/sql/create-external-volume#cloud-provider-parameters-cloudproviderparams) for more details on this.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"storage_location_name": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      blocklistedCharactersFieldDescription("Name of the storage location. Must be unique for the external volume. Do not use the name `terraform_provider_sentinel_storage_location` - this is reserved for the provider for performing update operations."),
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
				"storage_provider": {
					Type:             schema.TypeString,
					Required:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToStorageProvider),
					DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToStorageProvider)),
					Description:      fmt.Sprintf("Specifies the cloud storage provider that stores your data files. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllStorageProviderValues)),
				},
				"storage_base_url": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the base URL for your cloud storage location.",
				},
				"storage_aws_role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the case-sensitive Amazon Resource Name (ARN) of the AWS identity and access management (IAM) role that grants privileges on the S3 bucket containing your data files.",
				},
				"storage_aws_external_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "External ID that Snowflake uses to establish a trust relationship with AWS.",
				},
				"encryption_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the encryption type used.",
					DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
						return oldValue == "NONE" && newValue == ""
					},
				},
				"encryption_kms_key_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the ID for the KMS-managed key used to encrypt files.",
				},
				"azure_tenant_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the ID for your Office 365 tenant that the allowed and blocked storage accounts belong to.",
				},
			},
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the external volume.",
	},
	"allow_writes": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     BooleanDefault,
		Description: booleanStringFieldDescription("Specifies whether write operations are allowed for the external volume; must be set to TRUE for Iceberg tables that use Snowflake as the catalog."),
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW EXTERNAL VOLUMES` for the given external volume.",
		Elem: &schema.Resource{
			Schema: schemas.ShowExternalVolumeSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE EXTERNAL VOLUME` for the given external volume.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeExternalVolumeSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func ExternalVolume() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ExternalVolumes.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ExternalVolumeResource), TrackingCreateWrapper(resources.ExternalVolume, CreateContextExternalVolume)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ExternalVolumeResource), TrackingReadWrapper(resources.ExternalVolume, ReadContextExternalVolume(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ExternalVolumeResource), TrackingUpdateWrapper(resources.ExternalVolume, UpdateContextExternalVolume)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ExternalVolumeResource), TrackingDeleteWrapper(resources.ExternalVolume, deleteFunc)),

		Description: "Resource used to manage external volume objects. For more information, check [external volume documentation](https://docs.snowflake.com/en/sql-reference/commands-data-loading#external-volume).",

		Schema: externalVolumeSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ExternalVolume, ImportExternalVolume),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ExternalVolume, customdiff.All(
			ComputedIfAnyAttributeChanged(externalVolumeSchema, ShowOutputAttributeName, "name", "allow_writes", "comment"),
			ComputedIfAnyAttributeChanged(externalVolumeSchema, DescribeOutputAttributeName, "name", "allow_writes", "comment", "storage_location"),
		)),
		Timeouts: defaultTimeouts,
	}
}

func ImportExternalVolume(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("name", id.Name()); err != nil {
		return nil, err
	}

	externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("allow_writes", booleanStringFromBool(externalVolume.AllowWrites)); err != nil {
		return nil, err
	}

	externalVolumeDescribe, err := client.ExternalVolumes.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	parsedExternalVolumeDescribed, err := helpers.ParseExternalVolumeDescribed(externalVolumeDescribe)
	if err != nil {
		return nil, err
	}

	storageLocations := make([]map[string]any, len(parsedExternalVolumeDescribed.StorageLocations))
	for i, storageLocation := range parsedExternalVolumeDescribed.StorageLocations {
		storageLocations[i] = map[string]any{
			"storage_location_name":   storageLocation.Name,
			"storage_provider":        storageLocation.StorageProvider,
			"storage_base_url":        storageLocation.StorageBaseUrl,
			"storage_aws_role_arn":    storageLocation.StorageAwsRoleArn,
			"storage_aws_external_id": storageLocation.StorageAwsExternalId,
			"encryption_type":         storageLocation.EncryptionType,
			"encryption_kms_key_id":   storageLocation.EncryptionKmsKeyId,
			"azure_tenant_id":         storageLocation.AzureTenantId,
		}
	}

	if err = d.Set("storage_location", storageLocations); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextExternalVolume(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	storageLocations, err := extractStorageLocations(d.Get("storage_location"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating external volume %v err = %w", id.Name(), err))
	}

	req := sdk.NewCreateExternalVolumeRequest(id, storageLocations)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v := d.Get("allow_writes").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithAllowWrites(parsed)
	}

	err = client.ExternalVolumes.Create(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating external volume %v err = %w", id.Name(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadContextExternalVolume(false)(ctx, d, meta)
}

func ReadContextExternalVolume(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		externalVolume, err := client.ExternalVolumes.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query external volume. Marking the resource as removed.",
						Detail:   fmt.Sprintf("External Volume id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}

			return diag.FromErr(err)
		}

		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"allow_writes", "allow_writes", externalVolume.AllowWrites, booleanStringFromBool(externalVolume.AllowWrites), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, externalVolumeSchema, []string{
			"allow_writes",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("comment", externalVolume.Comment); err != nil {
			return diag.FromErr(err)
		}

		externalVolumeDescribe, err := client.ExternalVolumes.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		parsedExternalVolumeDescribed, err := helpers.ParseExternalVolumeDescribed(externalVolumeDescribe)
		if err != nil {
			return diag.FromErr(err)
		}

		storageLocations := make([]map[string]any, len(parsedExternalVolumeDescribed.StorageLocations))
		for i, storageLocation := range parsedExternalVolumeDescribed.StorageLocations {
			storageLocations[i] = map[string]any{
				"storage_location_name":   storageLocation.Name,
				"storage_provider":        storageLocation.StorageProvider,
				"storage_base_url":        storageLocation.StorageBaseUrl,
				"storage_aws_role_arn":    storageLocation.StorageAwsRoleArn,
				"storage_aws_external_id": storageLocation.StorageAwsExternalId,
				"encryption_type":         storageLocation.EncryptionType,
				"encryption_kms_key_id":   storageLocation.EncryptionKmsKeyId,
				"azure_tenant_id":         storageLocation.AzureTenantId,
			}
		}

		if err = d.Set("storage_location", storageLocations); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(DescribeOutputAttributeName, schemas.ExternalVolumeDescriptionToSchema(externalVolumeDescribe)); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.ExternalVolumeToSchema(externalVolume)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateContextExternalVolume(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewAlterExternalVolumeSetRequest()

	if d.HasChange("comment") {
		// not using d.GetOk as that doesn't let comments be reset to the empty string
		set.WithComment(d.Get("comment").(string))
	}

	if d.HasChange("allow_writes") {
		if v := d.Get("allow_writes").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithAllowWrites(parsed)
		} else {
			// no way to unset allow writes - set to false as a default
			set.WithAllowWrites(false)
		}
	}

	if (*set != sdk.AlterExternalVolumeSetRequest{}) {
		if err := client.ExternalVolumes.Alter(ctx, sdk.NewAlterExternalVolumeRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("storage_location") {
		old, new := d.GetChange("storage_location")
		oldLocations, err := extractStorageLocations(old)
		if err != nil {
			return diag.FromErr(err)
		}

		newLocations, err := extractStorageLocations(new)
		if err != nil {
			return diag.FromErr(err)
		}

		// Storage locations can only be added to the tail of the list, but can be
		// removed at any position. Given this limitation, to keep the configuration order
		// matching that on Snowflake the list needs to be partially recreated. For example, if a location
		// is added in the configuration at index 5 in the list, all existing storage locations from index 5
		// need to be removed, then the new location can be added, and then the removed locations
		// can be added back. The storage locations lower than index 5 don't need to be modified.
		// The removal process could be done without the above recreation, but it handles this case
		// too so it's used for both actions.
		commonPrefixLastIndex, err := sdk.CommonPrefixLastIndex(newLocations, oldLocations)
		if err != nil {
			return diag.FromErr(err)
		}

		var removedLocations []sdk.ExternalVolumeStorageLocation
		var addedLocations []sdk.ExternalVolumeStorageLocation
		if commonPrefixLastIndex == -1 {
			removedLocations = oldLocations
			addedLocations = newLocations
		} else {
			// Could +1 on the prefix here as the lists until and including this index
			// are identical, would need to add some more checks for list length to avoid
			// an array index out of bounds error
			removedLocations = oldLocations[commonPrefixLastIndex:]
			addedLocations = newLocations[commonPrefixLastIndex:]
		}

		if len(removedLocations) == len(oldLocations) {
			// Create a temporary storage location, which is a copy of a storage location currently existing
			// except with a different name. This is done to avoid recreating the external volume, which
			// would otherwise be necessary as a minimum of 1 storage location per external volume is required.
			// The alternative solution of adding volumes before removing them isn't possible as
			// name must be unique for storage locations
			tempStorageLocation, err := sdk.CopySentinelStorageLocation(removedLocations[0])
			if err != nil {
				return diag.FromErr(err)
			}

			addTempErr := addStorageLocation(tempStorageLocation, client, ctx, id)
			if addTempErr != nil {
				return diag.FromErr(addTempErr)
			}

			updateErr := updateStorageLocations(removedLocations, addedLocations, client, ctx, id)
			// TODO use defer for the removal of the temp storage location
			if updateErr != nil {
				// Try to remove the temp location and then return with error
				removeErr := removeStorageLocation(tempStorageLocation, client, ctx, id)
				if removeErr != nil {
					return diag.FromErr(errors.Join(updateErr, removeErr))
				}

				return diag.FromErr(updateErr)
			}

			removeErr := removeStorageLocation(tempStorageLocation, client, ctx, id)
			if removeErr != nil {
				return diag.FromErr(removeErr)
			}
		} else {
			updateErr := updateStorageLocations(removedLocations, addedLocations, client, ctx, id)
			if updateErr != nil {
				return diag.FromErr(updateErr)
			}
		}
	}

	return ReadContextExternalVolume(false)(ctx, d, meta)
}

func extractStorageLocations(v any) ([]sdk.ExternalVolumeStorageLocation, error) {
	_, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("unable to extract storage locations, input is either nil or non expected type (%T): %v", v, v)
	}

	storageLocations := make([]sdk.ExternalVolumeStorageLocation, len(v.([]any)))
	for i, storageLocationConfigRaw := range v.([]any) {
		storageLocationConfig, ok := storageLocationConfigRaw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unable to extract storage location, non expected type of %T: %v", storageLocationConfigRaw, storageLocationConfigRaw)
		}

		name, ok := storageLocationConfig["storage_location_name"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to extract storage location, missing storage_location_name key in storage location")
		}

		storageProvider, ok := storageLocationConfig["storage_provider"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to extract storage location, missing storage_provider key in storage location")
		}

		storageBaseUrl, ok := storageLocationConfig["storage_base_url"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to extract storage location, missing storage_base_url key in storage location")
		}

		storageProviderParsed, err := sdk.ToStorageProvider(storageProvider)
		if err != nil {
			return nil, err
		}

		var storageLocation sdk.ExternalVolumeStorageLocation
		switch storageProviderParsed {
		case sdk.StorageProviderS3, sdk.StorageProviderS3GOV:
			// Test that azure_tenant_id is not given
			// If given non empty plans will be produced
			azureTenantId, ok := storageLocationConfig["azure_tenant_id"].(string)
			if ok && len(azureTenantId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, azure_tenant_id provided for s3 storage location")
			}

			storageAwsRoleArn, ok := storageLocationConfig["storage_aws_role_arn"].(string)
			if !ok || len(storageAwsRoleArn) == 0 {
				return nil, fmt.Errorf("unable to extract storage location, missing storage_aws_role_arn key in an s3 storage location")
			}

			s3StorageProvider, err := sdk.ToS3StorageProvider(storageProvider)
			if err != nil {
				return nil, err
			}

			s3StorageLocation := &sdk.S3StorageLocationParams{
				Name:              name,
				StorageProvider:   s3StorageProvider,
				StorageBaseUrl:    storageBaseUrl,
				StorageAwsRoleArn: storageAwsRoleArn,
			}

			encryptionType, ok := storageLocationConfig["encryption_type"].(string)
			if ok && len(encryptionType) > 0 {
				encryptionTypeParsed, err := sdk.ToS3EncryptionType(encryptionType)
				if err != nil {
					return nil, err
				}

				encryptionKmsKeyId, ok := storageLocationConfig["encryption_kms_key_id"].(string)
				if ok && len(encryptionKmsKeyId) > 0 {
					s3StorageLocation.Encryption = &sdk.ExternalVolumeS3Encryption{
						Type:     encryptionTypeParsed,
						KmsKeyId: &encryptionKmsKeyId,
					}
				} else {
					s3StorageLocation.Encryption = &sdk.ExternalVolumeS3Encryption{
						Type: encryptionTypeParsed,
					}
				}
			}

			storageLocation = sdk.ExternalVolumeStorageLocation{
				S3StorageLocationParams: s3StorageLocation,
			}
		case sdk.StorageProviderGCS:
			// Test that azure_tenant_id and storage_aws_role_arn are not given
			// If given non empty plans will be produced
			azureTenantId, ok := storageLocationConfig["azure_tenant_id"].(string)
			if ok && len(azureTenantId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, azure_tenant_id provided for gcs storage location")
			}

			storageAwsRoleArn, ok := storageLocationConfig["storage_aws_role_arn"].(string)
			if ok && len(storageAwsRoleArn) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_role_arn provided for gcs storage location")
			}

			gcsStorageLocation := &sdk.GCSStorageLocationParams{
				Name:           name,
				StorageBaseUrl: storageBaseUrl,
			}
			encryptionType, ok := storageLocationConfig["encryption_type"].(string)
			if ok && len(encryptionType) > 0 {
				encryptionTypeParsed, err := sdk.ToGCSEncryptionType(encryptionType)
				if err != nil {
					return nil, err
				}
				encryptionKmsKeyId, ok := storageLocationConfig["encryption_kms_key_id"].(string)
				if ok && len(encryptionKmsKeyId) > 0 {
					gcsStorageLocation.Encryption = &sdk.ExternalVolumeGCSEncryption{
						Type:     encryptionTypeParsed,
						KmsKeyId: &encryptionKmsKeyId,
					}
				} else {
					gcsStorageLocation.Encryption = &sdk.ExternalVolumeGCSEncryption{
						Type: encryptionTypeParsed,
					}
				}
			}

			storageLocation = sdk.ExternalVolumeStorageLocation{
				GCSStorageLocationParams: gcsStorageLocation,
			}
		case sdk.StorageProviderAzure:
			// Test that storage_aws_role_arn and encryption_kms_key_id is not given
			// If given non empty plans will be produced
			storageAwsRolArn, ok := storageLocationConfig["storage_aws_role_arn"].(string)
			if ok && len(storageAwsRolArn) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_role_arn provided for azure storage location")
			}

			encryptionKmsKeyId, ok := storageLocationConfig["encryption_kms_key_id"].(string)
			if ok && len(encryptionKmsKeyId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, encryption_kms_key_id provided for azure storage location")
			}

			// TODO add check that encryption_type is not set in the config for azure storage locations
			// This may be more difficult as NONE is returned as the encyption type from Snowflake for azure
			// storage locations, although it's not documented as a parameter

			azureTenantId, ok := storageLocationConfig["azure_tenant_id"].(string)
			if !ok || len(azureTenantId) == 0 {
				return nil, fmt.Errorf("unable to extract storage location, missing azure_tenant_id provider key in an azure storage location")
			}

			storageLocation = sdk.ExternalVolumeStorageLocation{
				AzureStorageLocationParams: &sdk.AzureStorageLocationParams{
					Name:           name,
					AzureTenantId:  azureTenantId,
					StorageBaseUrl: storageBaseUrl,
				},
			}
		}
		storageLocations[i] = storageLocation
	}
	return storageLocations, nil
}

func addStorageLocation(
	addedLocation sdk.ExternalVolumeStorageLocation,
	client *sdk.Client,
	ctx context.Context,
	id sdk.AccountObjectIdentifier,
) error {
	storageProvider, err := sdk.GetStorageLocationStorageProvider(addedLocation)
	if err != nil {
		return err
	}

	var newStorageLocationreq *sdk.ExternalVolumeStorageLocationRequest
	switch storageProvider {
	case sdk.StorageProviderS3, sdk.StorageProviderS3GOV:
		addedLocation := addedLocation.S3StorageLocationParams
		s3ParamsRequest := sdk.NewS3StorageLocationParamsRequest(
			addedLocation.Name,
			addedLocation.StorageProvider,
			addedLocation.StorageAwsRoleArn,
			addedLocation.StorageBaseUrl,
		)
		if addedLocation.Encryption != nil {
			encryptionRequest := sdk.NewExternalVolumeS3EncryptionRequest(addedLocation.Encryption.Type)
			if addedLocation.Encryption.KmsKeyId != nil {
				encryptionRequest = encryptionRequest.WithKmsKeyId(*addedLocation.Encryption.KmsKeyId)
			}

			s3ParamsRequest = s3ParamsRequest.WithEncryption(*encryptionRequest)
		}

		newStorageLocationreq = sdk.NewExternalVolumeStorageLocationRequest().WithS3StorageLocationParams(*s3ParamsRequest)
	case sdk.StorageProviderGCS:
		addedLocation := addedLocation.GCSStorageLocationParams
		gcsParamsRequest := sdk.NewGCSStorageLocationParamsRequest(
			addedLocation.Name,
			addedLocation.StorageBaseUrl,
		)

		if addedLocation.Encryption != nil {
			encryptionRequest := sdk.NewExternalVolumeGCSEncryptionRequest(addedLocation.Encryption.Type)
			if addedLocation.Encryption.KmsKeyId != nil {
				encryptionRequest = encryptionRequest.WithKmsKeyId(*addedLocation.Encryption.KmsKeyId)
			}

			gcsParamsRequest = gcsParamsRequest.WithEncryption(*encryptionRequest)
		}

		newStorageLocationreq = sdk.NewExternalVolumeStorageLocationRequest().WithGCSStorageLocationParams(*gcsParamsRequest)
	case sdk.StorageProviderAzure:
		addedLocation := addedLocation.AzureStorageLocationParams
		azureParamsRequest := sdk.NewAzureStorageLocationParamsRequest(
			addedLocation.Name,
			addedLocation.AzureTenantId,
			addedLocation.StorageBaseUrl,
		)
		newStorageLocationreq = sdk.NewExternalVolumeStorageLocationRequest().WithAzureStorageLocationParams(*azureParamsRequest)
	}

	return client.ExternalVolumes.Alter(ctx, sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(*newStorageLocationreq))
}

func removeStorageLocation(
	removedLocation sdk.ExternalVolumeStorageLocation,
	client *sdk.Client,
	ctx context.Context,
	id sdk.AccountObjectIdentifier,
) error {
	removedName, err := sdk.GetStorageLocationName(removedLocation)
	if err != nil {
		return err
	}

	return client.ExternalVolumes.Alter(ctx, sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(removedName))
}

// Process the removal / addition storage location requests.
// to avoid creating storage locations with duplicate names.
// len(removedLocations) should be less than the total number
// of storage locations the external volume has, else this function will fail.
func updateStorageLocations(
	removedLocations []sdk.ExternalVolumeStorageLocation,
	addedLocations []sdk.ExternalVolumeStorageLocation,
	client *sdk.Client,
	ctx context.Context,
	id sdk.AccountObjectIdentifier,
) error {
	for _, removedLocation := range removedLocations {
		err := removeStorageLocation(removedLocation, client, ctx, id)
		if err != nil {
			return err
		}
	}
	for _, addedLocation := range addedLocations {
		err := addStorageLocation(addedLocation, client, ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}
