package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					Description:      fmt.Sprintf("Specifies the cloud storage provider that stores your data files. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllStorageProviders)),
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
						return oldValue == "" && (newValue == BooleanDefault || newValue == BooleanFalse)
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
			Schema: schemas.DescribeExternalVolumeDetailsSchema,
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

		CreateContext: TrackingCreateWrapper(resources.ExternalVolume, CreateContextExternalVolume),
		ReadContext:   TrackingReadWrapper(resources.ExternalVolume, ReadContextExternalVolume(true)),
		UpdateContext: TrackingUpdateWrapper(resources.ExternalVolume, UpdateContextExternalVolume),
		DeleteContext: TrackingDeleteWrapper(resources.ExternalVolume, deleteFunc),

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
			if loc.AzureStorageLocation.UsePrivatelinkEndpoint != nil {
				m["use_privatelink_endpoint"] = booleanStringFromBool(*loc.AzureStorageLocation.UsePrivatelinkEndpoint)
			}
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

	var storageLocations []sdk.ExternalVolumeStorageLocationItem
	if err := listBlockAttributeCreate(d, "storage_location", &storageLocations, parseStorageLocation); err != nil {
		return diag.FromErr(fmt.Errorf("error creating external volume %v err = %w", id.Name(), err))
	}

	storageLocationRequests, err := collections.MapErr(storageLocations, storageLocationItemToRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating external volume %v err = %w", id.Name(), err))
	}

	req := sdk.NewCreateExternalVolumeRequest(id, storageLocationRequests)

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
			if err = handleExternalChangesToObjectInShow(
				d,
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

		detailsSchema := schemas.ExternalVolumeDetailsToSchema(&parsedExternalVolumeDescribed)

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
			// commonPrefixLastIndex is inclusive, so the elements up to and including it are
			// identical between old and new; slicing from the next index excludes them from
			// both removed and added, avoiding an unnecessary remove/re-add of unchanged locations.
			removedLocations = oldLocations[commonPrefixLastIndex+1:]
			addedLocations = newLocations[commonPrefixLastIndex+1:]
		}

		removeLocation := func(location sdk.ExternalVolumeStorageLocationItem) error {
			return removeStorageLocation(location, client, ctx, id)
		}
		addLocation := func(location sdk.ExternalVolumeStorageLocationItem) error {
			return addStorageLocation(location, client, ctx, id)
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

			if err := removeThenAddWithPlaceholder(tempStorageLocation, removedLocations, addedLocations, removeLocation, addLocation); err != nil {
				return diag.FromErr(err)
			}
		} else if err := removeThenAdd(removedLocations, addedLocations, removeLocation, addLocation); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextExternalVolume(false)(ctx, d, meta)
}

func extractStorageLocations(v any) ([]sdk.ExternalVolumeStorageLocationItem, error) {
	if _, ok := v.([]any); !ok {
		return nil, fmt.Errorf("unable to extract storage locations, input is either nil or non expected type (%T): %v", v, v)
	}
	// Add/remove still diffs the raw old and new lists. Each block is parsed with the same prefix helpers as create.
	return listObjects(v, func(block map[string]any) (sdk.ExternalVolumeStorageLocationItem, error) {
		locationData, err := storageLocationResourceData(block)
		if err != nil {
			return sdk.ExternalVolumeStorageLocationItem{}, err
		}
		return parseStorageLocation(locationData, "")
	})
}

func storageLocationResourceData(block map[string]any) (*schema.ResourceData, error) {
	attrs := make(map[string]string, len(block))
	for key, value := range block {
		field, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("unable to extract storage location, %s is not a string", key)
		}
		attrs[key] = field
	}
	locationSchema := externalVolumeSchema["storage_location"].Elem.(*schema.Resource).Schema
	return schema.InternalMap(locationSchema).Data(&terraform.InstanceState{ID: "id", Attributes: attrs}, nil)
}

func parseStorageLocation(d *schema.ResourceData, prefix string) (sdk.ExternalVolumeStorageLocationItem, error) {
	name := d.Get(prefix + "storage_location_name").(string)
	storageProviderRaw := d.Get(prefix + "storage_provider").(string)
	storageBaseUrl := d.Get(prefix + "storage_base_url").(string)
	storageProvider, err := sdk.ToStorageProvider(storageProviderRaw)
	if err != nil {
		return sdk.ExternalVolumeStorageLocationItem{}, err
	}

	switch storageProvider {
	case sdk.StorageProviderS3, sdk.StorageProviderS3gov:
		return parseS3StorageLocation(d, prefix, name, storageBaseUrl)
	case sdk.StorageProviderGcs:
		return parseGCSStorageLocation(d, prefix, name, storageBaseUrl)
	case sdk.StorageProviderAzure:
		return parseAzureStorageLocation(d, prefix, name, storageBaseUrl)
	case sdk.StorageProviderS3compat:
		return parseS3CompatStorageLocation(d, prefix, name, storageBaseUrl)
	default:
		return sdk.ExternalVolumeStorageLocationItem{}, fmt.Errorf("unsupported storage provider: %s", storageProviderRaw)
	}
}

func parseS3StorageLocation(d *schema.ResourceData, prefix, name, storageBaseUrl string) (sdk.ExternalVolumeStorageLocationItem, error) {
	s3StorageLocation := &sdk.S3StorageLocationParams{
		StorageBaseUrl:    storageBaseUrl,
		StorageAwsRoleArn: d.Get(prefix + "storage_aws_role_arn").(string),
	}
	err := errors.Join(
		rejectFields(d, prefix, "s3",
			"azure_tenant_id",
			"storage_endpoint",
			"storage_aws_key_id",
			"storage_aws_secret_key",
		),
		requireStorageLocationField(d, prefix+"storage_aws_role_arn", "storage_aws_role_arn is required for s3 storage location"),
		attributeMappedValueCreateBuilder(d, prefix+"storage_provider", func(provider sdk.S3StorageProvider) *sdk.S3StorageLocationParams {
			s3StorageLocation.StorageProvider = provider
			return s3StorageLocation
		}, sdk.ToS3StorageProvider),
		stringAttributeCreate(d, prefix+"storage_aws_external_id", &s3StorageLocation.StorageAwsExternalId),
		stringAttributeCreate(d, prefix+"storage_aws_access_point_arn", &s3StorageLocation.StorageAwsAccessPointArn),
		setOptionalBooleanString(d, prefix+"use_privatelink_endpoint", &s3StorageLocation.UsePrivatelinkEndpoint),
		attributeMappedValueCreate(d, prefix+"encryption_type", &s3StorageLocation.Encryption, func(value any) (*sdk.ExternalVolumeS3Encryption, error) {
			encryptionType, err := sdk.ToS3EncryptionType(value.(string))
			if err != nil {
				return nil, err
			}
			encryption := &sdk.ExternalVolumeS3Encryption{EncryptionType: encryptionType}
			err = stringAttributeCreate(d, prefix+"encryption_kms_key_id", &encryption.KmsKeyId)
			return encryption, err
		}),
	)
	return sdk.ExternalVolumeStorageLocationItem{
		ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name:                    name,
			S3StorageLocationParams: s3StorageLocation,
		},
	}, err
}

func parseGCSStorageLocation(d *schema.ResourceData, prefix, name, storageBaseUrl string) (sdk.ExternalVolumeStorageLocationItem, error) {
	gcsStorageLocation := &sdk.GCSStorageLocationParams{
		StorageBaseUrl: storageBaseUrl,
	}
	err := errors.Join(
		rejectFields(d, prefix, "gcs",
			"azure_tenant_id",
			"storage_aws_role_arn",
			"storage_aws_external_id",
			"storage_aws_access_point_arn",
			"storage_endpoint",
			"storage_aws_key_id",
			"storage_aws_secret_key",
		),
		attributeMappedValueCreate(d, prefix+"encryption_type", &gcsStorageLocation.Encryption, func(value any) (*sdk.ExternalVolumeGCSEncryption, error) {
			encryptionType, err := sdk.ToGCSEncryptionType(value.(string))
			if err != nil {
				return nil, err
			}
			encryption := &sdk.ExternalVolumeGCSEncryption{EncryptionType: encryptionType}
			err = stringAttributeCreate(d, prefix+"encryption_kms_key_id", &encryption.KmsKeyId)
			return encryption, err
		}),
	)
	return sdk.ExternalVolumeStorageLocationItem{
		ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name:                     name,
			GCSStorageLocationParams: gcsStorageLocation,
		},
	}, err
}

func parseAzureStorageLocation(d *schema.ResourceData, prefix, name, storageBaseUrl string) (sdk.ExternalVolumeStorageLocationItem, error) {
	azureStorageLocation := &sdk.AzureStorageLocationParams{
		AzureTenantId:  d.Get(prefix + "azure_tenant_id").(string),
		StorageBaseUrl: storageBaseUrl,
	}
	err := errors.Join(
		rejectFields(d, prefix, "azure",
			"storage_aws_role_arn",
			"storage_aws_external_id",
			"encryption_kms_key_id",
			"storage_aws_access_point_arn",
			"storage_endpoint",
			"storage_aws_key_id",
			"storage_aws_secret_key",
		),
		requireStorageLocationField(d, prefix+"azure_tenant_id", "missing azure_tenant_id provider key in an azure storage location"),
		rejectIgnoredEncryptionType(d, prefix, "azure", func(value string) error {
			_, err := sdk.ToAzureEncryptionType(value)
			return err
		}),
		setOptionalBooleanString(d, prefix+"use_privatelink_endpoint", &azureStorageLocation.UsePrivatelinkEndpoint),
	)
	return sdk.ExternalVolumeStorageLocationItem{
		ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name:                       name,
			AzureStorageLocationParams: azureStorageLocation,
		},
	}, err
}

func parseS3CompatStorageLocation(d *schema.ResourceData, prefix, name, storageBaseUrl string) (sdk.ExternalVolumeStorageLocationItem, error) {
	s3CompatStorageLocation := &sdk.S3CompatStorageLocationParams{
		StorageBaseUrl:  storageBaseUrl,
		StorageEndpoint: d.Get(prefix + "storage_endpoint").(string),
		Credentials: sdk.ExternalVolumeS3CompatCredentials{
			AwsKeyId:     d.Get(prefix + "storage_aws_key_id").(string),
			AwsSecretKey: d.Get(prefix + "storage_aws_secret_key").(string),
		},
	}
	err := errors.Join(
		rejectFields(d, prefix, "s3compat",
			"storage_aws_role_arn",
			"storage_aws_external_id",
			"azure_tenant_id",
			"storage_aws_access_point_arn",
		),
		requireStorageLocationField(d, prefix+"storage_endpoint", "storage_endpoint is required for s3compat storage location"),
		requireStorageLocationField(d, prefix+"storage_aws_key_id", "storage_aws_key_id is required for s3compat storage location"),
		rejectIgnoredEncryptionType(d, prefix, "s3compat", func(value string) error {
			_, err := sdk.ToS3CompatEncryptionType(value)
			return err
		}),
	)
	return sdk.ExternalVolumeStorageLocationItem{
		ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name:                          name,
			S3CompatStorageLocationParams: s3CompatStorageLocation,
		},
	}, err
}

func rejectFields(d *schema.ResourceData, prefix, provider string, fields ...string) error {
	for _, field := range fields {
		if _, ok := d.GetOk(prefix + field); ok {
			return fmt.Errorf("unable to extract storage location, %s is not supported for %s storage location", field, provider)
		}
	}
	return nil
}

func requireStorageLocationField(d *schema.ResourceData, key, message string) error {
	if _, ok := d.GetOk(key); !ok {
		return fmt.Errorf("unable to extract storage location, %s", message)
	}
	return nil
}

// rejectIgnoredEncryptionType allows an empty type or the provider's NONE type, which Snowflake reports back but does not accept on write.
func rejectIgnoredEncryptionType(d *schema.ResourceData, prefix, provider string, validate func(string) error) error {
	raw, ok := d.GetOk(prefix + "encryption_type")
	if !ok {
		return nil
	}
	if err := validate(raw.(string)); err != nil {
		return fmt.Errorf("unable to extract storage location, encryption_type is not supported for %s storage location: %w", provider, err)
	}
	return nil
}

// setOptionalBooleanString leaves the field unset for an empty string or BooleanDefault.
// Read can persist "" for this tri-state field, whose schema default is BooleanDefault.
func setOptionalBooleanString(d *schema.ResourceData, key string, createField **bool) error {
	value := d.Get(key).(string)
	if value == "" || value == BooleanDefault {
		return nil
	}
	parsed, err := booleanStringToBool(value)
	if err != nil {
		return err
	}
	*createField = sdk.Bool(parsed)
	return nil
}

func storageLocationItemToRequest(addedLocationItem sdk.ExternalVolumeStorageLocationItem) (sdk.ExternalVolumeStorageLocationItemRequest, error) {
	storageProvider, err := sdk.GetStorageLocationStorageProvider(addedLocationItem)
	if err != nil {
		return sdk.ExternalVolumeStorageLocationItemRequest{}, err
	}

	var newStorageLocationreq *sdk.ExternalVolumeStorageLocationRequest
	switch storageProvider {
	case sdk.StorageProviderS3, sdk.StorageProviderS3gov:
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
	case sdk.StorageProviderGcs:
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
		if addedLocation.UsePrivatelinkEndpoint != nil {
			azureParamsRequest = azureParamsRequest.WithUsePrivatelinkEndpoint(*addedLocation.UsePrivatelinkEndpoint)
		}
		newStorageLocationreq = sdk.NewExternalVolumeStorageLocationRequest(addedLocationItem.ExternalVolumeStorageLocation.Name).WithAzureStorageLocationParams(*azureParamsRequest)
	case sdk.StorageProviderS3compat:
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
	default:
		return sdk.ExternalVolumeStorageLocationItemRequest{}, fmt.Errorf("unsupported storage provider: %s", storageProvider)
	}

	return sdk.ExternalVolumeStorageLocationItemRequest{ExternalVolumeStorageLocation: *newStorageLocationreq}, nil
}

func addStorageLocation(
	addedLocationItem sdk.ExternalVolumeStorageLocationItem,
	client *sdk.Client,
	ctx context.Context,
	id sdk.AccountObjectIdentifier,
) error {
	req, err := storageLocationItemToRequest(addedLocationItem)
	if err != nil {
		return err
	}

	return client.ExternalVolumes.Alter(ctx, sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(req))
}

func removeStorageLocation(
	removedLocation sdk.ExternalVolumeStorageLocationItem,
	client *sdk.Client,
	ctx context.Context,
	id sdk.AccountObjectIdentifier,
) error {
	return client.ExternalVolumes.Alter(ctx, sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(removedLocation.ExternalVolumeStorageLocation.Name))
}
