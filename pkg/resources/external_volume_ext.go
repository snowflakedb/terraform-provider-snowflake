package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	externalVolumeCreateExt        = TrackingCreateWrapper(resources.ExternalVolume, CreateExternalVolume)
	externalVolumeReadExt          = TrackingReadWrapper(resources.ExternalVolume, ReadExternalVolumeFunc(true))
	externalVolumeUpdateExt        = TrackingUpdateWrapper(resources.ExternalVolume, UpdateExternalVolume)
	externalVolumeDescriptionExt   = "Resource used to manage external volume objects. For more information, check [external volume documentation](https://docs.snowflake.com/en/sql-reference/commands-data-loading#external-volume)."
	externalVolumeCustomizeDiffExt = TrackingCustomDiffWrapper(resources.ExternalVolume, customdiff.All(
		ComputedIfAnyAttributeChanged(externalVolumeSchema, ShowOutputAttributeName, "name", "allow_writes", "comment"),
		// storage_location is missing on purpose, because ComputedIfAnyAttributeChanged does not handle nested diffs well enough.
		ComputedIfAnyAttributeChanged(externalVolumeSchema, DescribeOutputAttributeName, "name", "allow_writes", "comment"),
	))
	externalVolumeImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.ExternalVolume, ImportExternalVolume),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

type externalVolumeAlterRequestsExt struct {
	set *sdk.ExternalVolumeSetRequest
}

func externalVolumeNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateExternalVolumeRequest, error) {
	var storageLocations []sdk.ExternalVolumeStorageLocationItem
	if err := listBlockAttributeCreate(d, "storage_location", &storageLocations, parseStorageLocation); err != nil {
		return nil, fmt.Errorf("error creating external volume %v err = %w", id.Name(), err)
	}

	storageLocationRequests, err := collections.MapErr(storageLocations, storageLocationItemToRequest)
	if err != nil {
		return nil, fmt.Errorf("error creating external volume %v err = %w", id.Name(), err)
	}

	return sdk.NewCreateExternalVolumeRequest(id, storageLocationRequests), nil
}

func externalVolumeHandleExternalChangesExt(d *schema.ResourceData, externalVolume *sdk.ExternalVolume, withExternalChangesMarking bool) error {
	if !withExternalChangesMarking {
		return nil
	}
	return handleExternalChangesToObjectInShow(d,
		outputMapping{"allow_writes", "allow_writes", externalVolume.AllowWrites, booleanStringFromBool(externalVolume.AllowWrites), nil},
	)
}

func externalVolumeApplySetUnsetExt(d *schema.ResourceData) (*externalVolumeAlterRequestsExt, error) {
	set := sdk.NewExternalVolumeSetRequest()
	errs := errors.Join(
		stringAttributeUpdateSetOnlyNotEmpty(d, "comment", &set.Comment),
		booleanStringAttributeUnsetFallbackUpdate(d, "allow_writes", &set.AllowWrites, true),
	)
	if errs != nil {
		return nil, errs
	}
	return &externalVolumeAlterRequestsExt{set: set}, nil
}

func externalVolumeSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, externalVolume *sdk.ExternalVolume, importPermissiveMode bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		importBooleanStringAttribute(d, "allow_writes", importPermissiveMode, externalVolume.AllowWrites),
	)
}

func externalVolumeSetDescribeOutputExt(d *schema.ResourceData, _ sdk.AccountObjectIdentifier, details *sdk.ExternalVolumeDetails, describeErr error, withExternalChangesMarking bool) error {
	if describeErr != nil {
		return describeErr
	}
	return d.Set("storage_location", readStorageLocations(d, *details, withExternalChangesMarking))
}

// =============================================================================
// SDK incision points
// =============================================================================

func externalVolumeCreateInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, req *sdk.CreateExternalVolumeRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.ExternalVolumes.Create(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating external volume %v err = %w", id.Name(), err))
	}
	return nil
}

func externalVolumeShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.ExternalVolume, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	externalVolume, err := client.ExternalVolumes.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query external volume. Marking the resource as removed.",
					Detail:   fmt.Sprintf("External Volume id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return externalVolume, nil
}

func externalVolumeShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.ExternalVolume, error) {
	client := meta.(*provider.Context).Client
	return client.ExternalVolumes.ShowByID(ctx, id)
}

func externalVolumeDescribeInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.ExternalVolumeDetails, error) {
	client := meta.(*provider.Context).Client
	described, err := client.ExternalVolumes.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	parsed, err := sdk.ParseExternalVolumeDescribed(described)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func externalVolumeAlterSetInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *externalVolumeAlterRequestsExt) error {
	client := meta.(*provider.Context).Client
	if *reqs.set != (sdk.ExternalVolumeSetRequest{}) {
		if err := client.ExternalVolumes.Alter(ctx, sdk.NewAlterExternalVolumeRequest(id).WithSet(*reqs.set)); err != nil {
			return err
		}
	}

	if !d.HasChange("storage_location") {
		return nil
	}

	old, newLocationsRaw := d.GetChange("storage_location")
	oldLocations, err := extractStorageLocations(old)
	if err != nil {
		return err
	}
	newLocations, err := extractStorageLocations(newLocationsRaw)
	if err != nil {
		return err
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
			return err
		}
		return removeThenAddWithPlaceholder(tempStorageLocation, removedLocations, addedLocations, removeLocation, addLocation)
	}
	return removeThenAdd(removedLocations, addedLocations, removeLocation, addLocation)
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

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
		attributeMappedValueCreate(d, prefix+"encryption_type", &gcsStorageLocation.Encryption, func(value any) (*sdk.ExternalVolumeGcsEncryption, error) {
			encryptionType, err := sdk.ToGCSEncryptionType(value.(string))
			if err != nil {
				return nil, err
			}
			encryption := &sdk.ExternalVolumeGcsEncryption{EncryptionType: encryptionType}
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
			encryptionRequest := sdk.NewExternalVolumeGcsEncryptionRequest(addedLocation.Encryption.EncryptionType)
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
