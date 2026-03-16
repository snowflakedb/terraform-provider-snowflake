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
	"github.com/hashicorp/go-cty/cty"
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
					Optional:    true,
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
				"storage_aws_access_point_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the access point ARN for the S3 bucket containing your data files. Only applicable for S3 and S3GOV storage providers.",
				},
				"use_privatelink_endpoint": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          BooleanDefault,
					ValidateDiagFunc: validateBooleanString,
					Description:      booleanStringFieldDescription("Specifies whether to use a privatelink endpoint for the storage location. Only applicable for S3, S3GOV, and AZURE storage providers."),
					DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
						return oldValue == "" && newValue == BooleanDefault
					},
				},
				"azure_tenant_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the ID for your Office 365 tenant that the allowed and blocked storage accounts belong to.",
				},
				"storage_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the endpoint for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
				},
				"storage_aws_key_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the AWS key ID for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
				},
				"storage_aws_secret_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "Specifies the AWS secret key for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
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
		Description: "Outputs the result of `DESCRIBE EXTERNAL VOLUME` for the given external volume. Because of Terraform limitations, the changes on storage_location field do not mark this field as computed.",
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
		SchemaVersion: 1,

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
			// storage_location is missing on purpose, because ComputedIfAnyAttributeChanged does not handle nested diffs well enough.
			ComputedIfAnyAttributeChanged(externalVolumeSchema, DescribeOutputAttributeName, "name", "allow_writes", "comment"),
		)),
		Timeouts: defaultTimeouts,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    cty.EmptyObject,
				Upgrade: v2_14_0_ExternalVolumeStateUpgrader,
			},
		},
	}
}

func storageLocationDetailsToStateMaps(locations []sdk.ExternalVolumeStorageLocationDetails) []map[string]any {
	result := make([]map[string]any, len(locations))
	for i, loc := range locations {
		m := map[string]any{
			"storage_location_name": loc.Name,
			"storage_provider":      loc.StorageProvider,
			"storage_base_url":      loc.StorageBaseUrl,
			"encryption_type":       loc.EncryptionType,
		}
		switch {
		case loc.S3StorageLocation != nil:
			m["storage_aws_role_arn"] = loc.S3StorageLocation.StorageAwsRoleArn

			m["storage_aws_access_point_arn"] = loc.S3StorageLocation.StorageAwsAccessPointArn
			m["encryption_kms_key_id"] = loc.S3StorageLocation.EncryptionKmsKeyId
			if loc.S3StorageLocation.UsePrivatelinkEndpoint != nil {
				m["use_privatelink_endpoint"] = booleanStringFromBool(*loc.S3StorageLocation.UsePrivatelinkEndpoint)
			}
		case loc.AzureStorageLocation != nil:
			m["azure_tenant_id"] = loc.AzureStorageLocation.AzureTenantId
		case loc.S3CompatStorageLocation != nil:
			m["storage_endpoint"] = loc.S3CompatStorageLocation.Endpoint
			m["storage_aws_key_id"] = loc.S3CompatStorageLocation.AwsAccessKeyId
			m["encryption_kms_key_id"] = loc.S3CompatStorageLocation.EncryptionKmsKeyId
		case loc.GCSStorageLocation != nil:
			m["encryption_kms_key_id"] = loc.GCSStorageLocation.EncryptionKmsKeyId
		}
		result[i] = m
	}
	return result
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

	parsedExternalVolumeDescribed, err := sdk.ParseExternalVolumeDescribed(externalVolumeDescribe)
	if err != nil {
		return nil, err
	}

	storageLocations := storageLocationDetailsToStateMaps(parsedExternalVolumeDescribed.StorageLocations)

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

	errs := errors.Join(
		stringAttributeCreateBuilder(d, "comment", req.WithComment),
		booleanStringAttributeCreateBuilder(d, "allow_writes", req.WithAllowWrites),
	)
	if errs != nil {
		return diag.FromErr(errs)
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

		externalVolumeDescribe, err := client.ExternalVolumes.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		parsedExternalVolumeDescribed, err := sdk.ParseExternalVolumeDescribed(externalVolumeDescribe)
		if err != nil {
			return diag.FromErr(err)
		}

		storageLocations := readStorageLocations(d, parsedExternalVolumeDescribed, withExternalChangesMarking)

		detailsSchema := schemas.ExternalVolumeDetailsToSchema(parsedExternalVolumeDescribed)

		errs := errors.Join(
			d.Set("comment", externalVolume.Comment),
			d.Set("storage_location", storageLocations),
			d.Set(DescribeOutputAttributeName, []map[string]any{detailsSchema}),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.ExternalVolumeToSchema(externalVolume)}),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func readStorageLocations(d *schema.ResourceData, parsedExternalVolumeDescribed sdk.ExternalVolumeDetails, withExternalChangesMarking bool) []map[string]any {
	storageLocations := storageLocationDetailsToStateMaps(parsedExternalVolumeDescribed.StorageLocations)

	// Preserve fields not returned by the API (secret key) and user-configured fields
	// not tracked from DESCRIBE output (external id) from the previous state.
	oldSecretKeys := make(map[string]string)
	oldExternalIds := make(map[string]string)
	for i := range d.Get("storage_location.#").(int) {
		name := d.Get(fmt.Sprintf("storage_location.%d.storage_location_name", i)).(string)
		oldSecretKeys[name] = d.Get(fmt.Sprintf("storage_location.%d.storage_aws_secret_key", i)).(string)
		oldExternalIds[name] = d.Get(fmt.Sprintf("storage_location.%d.storage_aws_external_id", i)).(string)
	}
	for i := range len(storageLocations) {
		locName := storageLocations[i]["storage_location_name"].(string)
		if v, ok := oldSecretKeys[locName]; ok {
			storageLocations[i]["storage_aws_secret_key"] = v
		}
		if v, ok := oldExternalIds[locName]; ok {
			storageLocations[i]["storage_aws_external_id"] = v
		}
	}

	// Handle external changes for storage_aws_external_id in S3 storage locations.
	// Build a map from the previous describe_output keyed by location name, then compare
	// against the current SF value to detect external changes. Using names instead of
	// positional indexes avoids mismatches when locations are added/removed between reads.
	if withExternalChangesMarking {
		prevDescExternalIds := make(map[string]string)
		for i := range d.Get("describe_output.0.storage_locations.#").(int) {
			locName := d.Get(fmt.Sprintf("describe_output.0.storage_locations.%d.name", i)).(string)
			s3Locs := d.Get(fmt.Sprintf("describe_output.0.storage_locations.%d.s3_storage_location", i)).([]any)
			if len(s3Locs) > 0 {
				if s3Map, ok := s3Locs[0].(map[string]any); ok {
					prevDescExternalIds[locName] = s3Map["storage_aws_external_id"].(string)
				}
			}
		}
		for i, loc := range parsedExternalVolumeDescribed.StorageLocations {
			if loc.S3StorageLocation == nil || i >= len(storageLocations) {
				continue
			}
			if prev, ok := prevDescExternalIds[loc.Name]; ok && prev != loc.S3StorageLocation.StorageAwsExternalId {
				storageLocations[i]["storage_aws_external_id"] = loc.S3StorageLocation.StorageAwsExternalId
			}
		}
	}

	return storageLocations
}

func UpdateContextExternalVolume(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewAlterExternalVolumeSetRequest()

	errs := errors.Join(
		stringAttributeUpdateSetOnlyNotEmpty(d, "comment", &set.Comment),
		booleanStringAttributeUnsetFallbackUpdate(d, "allow_writes", &set.AllowWrites, false),
	)
	if errs != nil {
		return diag.FromErr(errs)
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
		commonPrefixLastIndex := collections.CommonPrefixLastIndex(newLocations, oldLocations, func(a, b sdk.ExternalVolumeStorageLocationItem) bool {
			return reflect.DeepEqual(a, b)
		})

		var removedLocations []sdk.ExternalVolumeStorageLocationItem
		var addedLocations []sdk.ExternalVolumeStorageLocationItem
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
			tempStorageLocation, err := sdk.CopySentinelStorageLocationItem(removedLocations[0])
			if err != nil {
				return diag.FromErr(err)
			}

			if err := updateStorageLocationsWithTemp(tempStorageLocation, removedLocations, addedLocations, client, ctx, id); err != nil {
				return diag.FromErr(err)
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

func extractStorageLocations(v any) ([]sdk.ExternalVolumeStorageLocationItem, error) {
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
			// Validate that provider-incompatible fields are not given
			azureTenantId, ok := storageLocationConfig["azure_tenant_id"].(string)
			if ok && len(azureTenantId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, azure_tenant_id is not supported for s3 storage location")
			}
			storageEndpoint, ok := storageLocationConfig["storage_endpoint"].(string)
			if ok && len(storageEndpoint) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_endpoint is not supported for s3 storage location")
			}
			storageAwsKeyId, ok := storageLocationConfig["storage_aws_key_id"].(string)
			if ok && len(storageAwsKeyId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_key_id is not supported for s3 storage location")
			}
			storageAwsSecretKey, ok := storageLocationConfig["storage_aws_secret_key"].(string)
			if ok && len(storageAwsSecretKey) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_secret_key is not supported for s3 storage location")
			}

			storageAwsRoleArn, ok := storageLocationConfig["storage_aws_role_arn"].(string)
			if !ok || len(storageAwsRoleArn) == 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_role_arn is required for s3 storage location")
			}

			s3StorageProvider, err := sdk.ToS3StorageProvider(storageProvider)
			if err != nil {
				return nil, err
			}

			s3StorageLocation := &sdk.S3StorageLocationParams{
				StorageProvider:   s3StorageProvider,
				StorageBaseUrl:    storageBaseUrl,
				StorageAwsRoleArn: storageAwsRoleArn,
			}

			storageAwsExternalId, ok := storageLocationConfig["storage_aws_external_id"].(string)
			if ok && len(storageAwsExternalId) > 0 {
				s3StorageLocation.StorageAwsExternalId = &storageAwsExternalId
			}

			storageAwsAccessPointArn, ok := storageLocationConfig["storage_aws_access_point_arn"].(string)
			if ok && len(storageAwsAccessPointArn) > 0 {
				s3StorageLocation.StorageAwsAccessPointArn = &storageAwsAccessPointArn
			}

			usePrivatelinkEndpoint, ok := storageLocationConfig["use_privatelink_endpoint"].(string)
			if ok && usePrivatelinkEndpoint != BooleanDefault && len(usePrivatelinkEndpoint) > 0 {
				b, err := booleanStringToBool(usePrivatelinkEndpoint)
				if err != nil {
					return nil, err
				}
				s3StorageLocation.UsePrivatelinkEndpoint = &b
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
						EncryptionType: encryptionTypeParsed,
						KmsKeyId:       &encryptionKmsKeyId,
					}
				} else {
					s3StorageLocation.Encryption = &sdk.ExternalVolumeS3Encryption{
						EncryptionType: encryptionTypeParsed,
					}
				}
			}

			storageLocation = sdk.ExternalVolumeStorageLocation{
				Name:                    name,
				S3StorageLocationParams: s3StorageLocation,
			}
		case sdk.StorageProviderGCS:
			// Validate that provider-incompatible fields are not given
			azureTenantId, ok := storageLocationConfig["azure_tenant_id"].(string)
			if ok && len(azureTenantId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, azure_tenant_id is not supported for gcs storage location")
			}
			storageAwsRoleArn, ok := storageLocationConfig["storage_aws_role_arn"].(string)
			if ok && len(storageAwsRoleArn) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_role_arn is not supported for gcs storage location")
			}
			storageAwsExternalId, ok := storageLocationConfig["storage_aws_external_id"].(string)
			if ok && len(storageAwsExternalId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_external_id is not supported for gcs storage location")
			}
			storageAwsAccessPointArn, ok := storageLocationConfig["storage_aws_access_point_arn"].(string)
			if ok && len(storageAwsAccessPointArn) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_access_point_arn is not supported for gcs storage location")
			}
			storageEndpoint, ok := storageLocationConfig["storage_endpoint"].(string)
			if ok && len(storageEndpoint) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_endpoint is not supported for gcs storage location")
			}
			storageAwsKeyId, ok := storageLocationConfig["storage_aws_key_id"].(string)
			if ok && len(storageAwsKeyId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_key_id is not supported for gcs storage location")
			}
			storageAwsSecretKey, ok := storageLocationConfig["storage_aws_secret_key"].(string)
			if ok && len(storageAwsSecretKey) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_secret_key is not supported for gcs storage location")
			}

			gcsStorageLocation := &sdk.GCSStorageLocationParams{
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
						EncryptionType: encryptionTypeParsed,
						KmsKeyId:       &encryptionKmsKeyId,
					}
				} else {
					gcsStorageLocation.Encryption = &sdk.ExternalVolumeGCSEncryption{
						EncryptionType: encryptionTypeParsed,
					}
				}
			}

			storageLocation = sdk.ExternalVolumeStorageLocation{
				Name:                     name,
				GCSStorageLocationParams: gcsStorageLocation,
			}
		case sdk.StorageProviderAzure:
			// Validate that provider-incompatible fields are not given
			storageAwsRolArn, ok := storageLocationConfig["storage_aws_role_arn"].(string)
			if ok && len(storageAwsRolArn) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_role_arn is not supported for azure storage location")
			}
			storageAwsExternalId, ok := storageLocationConfig["storage_aws_external_id"].(string)
			if ok && len(storageAwsExternalId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_external_id is not supported for azure storage location")
			}
			encryptionKmsKeyId, ok := storageLocationConfig["encryption_kms_key_id"].(string)
			if ok && len(encryptionKmsKeyId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, encryption_kms_key_id is not supported for azure storage location")
			}
			storageAwsAccessPointArn, ok := storageLocationConfig["storage_aws_access_point_arn"].(string)
			if ok && len(storageAwsAccessPointArn) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_access_point_arn is not supported for azure storage location")
			}
			storageEndpoint, ok := storageLocationConfig["storage_endpoint"].(string)
			if ok && len(storageEndpoint) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_endpoint is not supported for azure storage location")
			}
			storageAwsKeyId, ok := storageLocationConfig["storage_aws_key_id"].(string)
			if ok && len(storageAwsKeyId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_key_id is not supported for azure storage location")
			}
			storageAwsSecretKey, ok := storageLocationConfig["storage_aws_secret_key"].(string)
			if ok && len(storageAwsSecretKey) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_secret_key is not supported for azure storage location")
			}
			azureTenantId, ok := storageLocationConfig["azure_tenant_id"].(string)
			if !ok || len(azureTenantId) == 0 {
				return nil, fmt.Errorf("unable to extract storage location, missing azure_tenant_id provider key in an azure storage location")
			}
			// Azure location doesn't support setting encryption, so we want to disallow it. However, Snowflake returns NONE for encryption_type.
			// So, here we allow NONE, but it is not used in the request anyway.
			encryptionTypeRaw, _ := storageLocationConfig["encryption_type"].(string)
			if encryptionTypeRaw != "" {
				_, err := sdk.ToAzureEncryptionType(encryptionTypeRaw)
				if err != nil {
					return nil, fmt.Errorf("unable to extract storage location, encryption_type is not supported for azure storage location: %w", err)
				}
			}

			// TODO(SNOW-2356128): handle use_privatelink_endpoint for Azure once testing on Azure deployment is possible
			storageLocation = sdk.ExternalVolumeStorageLocation{
				Name: name,
				AzureStorageLocationParams: &sdk.AzureStorageLocationParams{
					AzureTenantId:  azureTenantId,
					StorageBaseUrl: storageBaseUrl,
				},
			}
		case sdk.StorageProviderS3Compatible:
			// Validate that provider-incompatible fields are not given
			storageAwsRoleArn, ok := storageLocationConfig["storage_aws_role_arn"].(string)
			if ok && len(storageAwsRoleArn) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_role_arn is not supported for s3compat storage location")
			}
			storageAwsExternalId, ok := storageLocationConfig["storage_aws_external_id"].(string)
			if ok && len(storageAwsExternalId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_external_id is not supported for s3compat storage location")
			}
			azureTenantId, ok := storageLocationConfig["azure_tenant_id"].(string)
			if ok && len(azureTenantId) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, azure_tenant_id is not supported for s3compat storage location")
			}
			storageAwsAccessPointArn, ok := storageLocationConfig["storage_aws_access_point_arn"].(string)
			if ok && len(storageAwsAccessPointArn) > 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_access_point_arn is not supported for s3compat storage location")
			}

			storageEndpoint, ok := storageLocationConfig["storage_endpoint"].(string)
			if !ok || len(storageEndpoint) == 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_endpoint is required for s3compat storage location")
			}
			storageAwsKeyId, ok := storageLocationConfig["storage_aws_key_id"].(string)
			if !ok || len(storageAwsKeyId) == 0 {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_key_id is required for s3compat storage location")
			}
			// storage_aws_secret_key is required, but it is not returned from Snowflake. That's why we don't validate it here.
			storageAwsSecretKey, ok := storageLocationConfig["storage_aws_secret_key"].(string)
			if !ok {
				return nil, fmt.Errorf("unable to extract storage location, storage_aws_secret_key is not a string")
			}
			// s3compat location doesn't support setting encryption, so we want to disallow it. However, Snowflake returns NONE for encryption_type.
			// So, here we allow NONE, but it is not used in the request anyway.
			encryptionTypeRaw, ok := storageLocationConfig["encryption_type"].(string)
			if !ok {
				return nil, fmt.Errorf("unable to extract storage location, encryption_type is not a string")
			}
			if encryptionTypeRaw != "" {
				_, err := sdk.ToS3CompatEncryptionType(encryptionTypeRaw)
				if err != nil {
					return nil, fmt.Errorf("unable to extract storage location, encryption_type is not supported for s3compat storage location: %w", err)
				}
			}

			storageLocation = sdk.ExternalVolumeStorageLocation{
				Name: name,
				S3CompatStorageLocationParams: &sdk.S3CompatStorageLocationParams{
					StorageBaseUrl:  storageBaseUrl,
					StorageEndpoint: storageEndpoint,
					Credentials: sdk.ExternalVolumeS3CompatCredentials{
						AwsKeyId:     storageAwsKeyId,
						AwsSecretKey: storageAwsSecretKey,
					},
				},
			}
		}
		storageLocations[i] = storageLocation
	}
	return collections.Map(storageLocations, func(storageLocation sdk.ExternalVolumeStorageLocation) sdk.ExternalVolumeStorageLocationItem {
		return sdk.ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: storageLocation}
	}), nil
}

func addStorageLocation(
	addedLocationItem sdk.ExternalVolumeStorageLocationItem,
	client *sdk.Client,
	ctx context.Context,
	id sdk.AccountObjectIdentifier,
) error {
	storageProvider, err := sdk.GetStorageLocationStorageProvider(addedLocationItem)
	if err != nil {
		return err
	}

	var newStorageLocationreq *sdk.ExternalVolumeStorageLocationRequest
	switch storageProvider {
	case sdk.StorageProviderS3, sdk.StorageProviderS3GOV:
		addedLocation := addedLocationItem.ExternalVolumeStorageLocation.S3StorageLocationParams
		s3ParamsRequest := sdk.NewS3StorageLocationParamsRequest(
			addedLocation.StorageProvider,
			addedLocation.StorageAwsRoleArn,
			addedLocation.StorageBaseUrl,
		)
		if addedLocation.StorageAwsAccessPointArn != nil {
			s3ParamsRequest = s3ParamsRequest.WithStorageAwsAccessPointArn(*addedLocation.StorageAwsAccessPointArn)
		}
		if addedLocation.UsePrivatelinkEndpoint != nil {
			s3ParamsRequest = s3ParamsRequest.WithUsePrivatelinkEndpoint(*addedLocation.UsePrivatelinkEndpoint)
		}
		if addedLocation.StorageAwsExternalId != nil {
			s3ParamsRequest = s3ParamsRequest.WithStorageAwsExternalId(*addedLocation.StorageAwsExternalId)
		}
		if addedLocation.Encryption != nil {
			encryptionRequest := sdk.NewExternalVolumeS3EncryptionRequest(addedLocation.Encryption.EncryptionType)
			if addedLocation.Encryption.KmsKeyId != nil {
				encryptionRequest = encryptionRequest.WithKmsKeyId(*addedLocation.Encryption.KmsKeyId)
			}

			s3ParamsRequest = s3ParamsRequest.WithEncryption(*encryptionRequest)
		}

		newStorageLocationreq = sdk.NewExternalVolumeStorageLocationRequest(addedLocationItem.ExternalVolumeStorageLocation.Name).WithS3StorageLocationParams(*s3ParamsRequest)
	case sdk.StorageProviderGCS:
		addedLocation := addedLocationItem.ExternalVolumeStorageLocation.GCSStorageLocationParams
		gcsParamsRequest := sdk.NewGCSStorageLocationParamsRequest(
			addedLocation.StorageBaseUrl,
		)

		if addedLocation.Encryption != nil {
			encryptionRequest := sdk.NewExternalVolumeGCSEncryptionRequest(addedLocation.Encryption.EncryptionType)
			if addedLocation.Encryption.KmsKeyId != nil {
				encryptionRequest = encryptionRequest.WithKmsKeyId(*addedLocation.Encryption.KmsKeyId)
			}

			gcsParamsRequest = gcsParamsRequest.WithEncryption(*encryptionRequest)
		}

		newStorageLocationreq = sdk.NewExternalVolumeStorageLocationRequest(addedLocationItem.ExternalVolumeStorageLocation.Name).WithGCSStorageLocationParams(*gcsParamsRequest)
	case sdk.StorageProviderAzure:
		addedLocation := addedLocationItem.ExternalVolumeStorageLocation.AzureStorageLocationParams
		azureParamsRequest := sdk.NewAzureStorageLocationParamsRequest(
			addedLocation.AzureTenantId,
			addedLocation.StorageBaseUrl,
		)
		// TODO(SNOW-2356128): handle use_privatelink_endpoint for Azure once testing on Azure deployment is possible
		newStorageLocationreq = sdk.NewExternalVolumeStorageLocationRequest(addedLocationItem.ExternalVolumeStorageLocation.Name).WithAzureStorageLocationParams(*azureParamsRequest)
	case sdk.StorageProviderS3Compatible:
		addedLocation := addedLocationItem.ExternalVolumeStorageLocation.S3CompatStorageLocationParams
		s3CompatParamsRequest := sdk.NewS3CompatStorageLocationParamsRequest(
			addedLocation.StorageBaseUrl,
			addedLocation.StorageEndpoint,
			sdk.ExternalVolumeS3CompatCredentialsRequest{
				AwsKeyId:     addedLocation.Credentials.AwsKeyId,
				AwsSecretKey: addedLocation.Credentials.AwsSecretKey,
			},
		)
		newStorageLocationreq = sdk.NewExternalVolumeStorageLocationRequest(addedLocationItem.ExternalVolumeStorageLocation.Name).WithS3CompatStorageLocationParams(*s3CompatParamsRequest)
	}

	return client.ExternalVolumes.Alter(ctx, sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(sdk.ExternalVolumeStorageLocationItemRequest{ExternalVolumeStorageLocation: *newStorageLocationreq}))
}

func removeStorageLocation(
	removedLocation sdk.ExternalVolumeStorageLocationItem,
	client *sdk.Client,
	ctx context.Context,
	id sdk.AccountObjectIdentifier,
) error {
	return client.ExternalVolumes.Alter(ctx, sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(removedLocation.ExternalVolumeStorageLocation.Name))
}

// updateStorageLocationsWithTemp adds a temporary storage location, performs the update, and ensures
// the temporary location is always cleaned up. This is used when all existing storage locations
// need to be replaced — because an external volume requires at least one storage location at all times.
func updateStorageLocationsWithTemp(
	tempStorageLocation sdk.ExternalVolumeStorageLocationItem,
	removedLocations []sdk.ExternalVolumeStorageLocationItem,
	addedLocations []sdk.ExternalVolumeStorageLocationItem,
	client *sdk.Client,
	ctx context.Context,
	id sdk.AccountObjectIdentifier,
) (err error) {
	if err := addStorageLocation(tempStorageLocation, client, ctx, id); err != nil {
		return err
	}
	defer func() {
		removeErr := removeStorageLocation(tempStorageLocation, client, ctx, id)
		err = errors.Join(err, removeErr)
	}()

	return updateStorageLocations(removedLocations, addedLocations, client, ctx, id)
}

// Process the removal / addition storage location requests.
// to avoid creating storage locations with duplicate names.
// len(removedLocations) should be less than the total number
// of storage locations the external volume has, else this function will fail.
func updateStorageLocations(
	removedLocations []sdk.ExternalVolumeStorageLocationItem,
	addedLocations []sdk.ExternalVolumeStorageLocationItem,
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
