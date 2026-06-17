package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var postgresInstanceSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Postgres instance; must be unique for your account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"compute_family": {
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToPostgresInstanceComputeFamily),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToPostgresInstanceComputeFamily),
		Description:      fmt.Sprintf("Specifies the compute family for the Postgres instance. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllPostgresInstanceComputeFamilies)),
	},
	"storage_size_gb": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the storage size in GB for the Postgres instance.",
	},
	"authentication_authority": {
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToPostgresInstanceAuthenticationAuthority),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToPostgresInstanceAuthenticationAuthority),
		Description:      fmt.Sprintf("Specifies the authentication authority for the Postgres instance. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllPostgresInstanceAuthenticationAuthorities)),
	},
	"postgres_version": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("postgres_version"),
		Description:      "Specifies the Postgres version for the instance.",
	},
	"network_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies the network policy to associate with the Postgres instance.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"high_availability": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_ha"),
		Description:      booleanStringFieldDescription("Specifies whether the Postgres instance should be configured for high availability."),
		Default:          BooleanDefault,
	},
	"storage_integration": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies the storage integration for the Postgres instance.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"postgres_settings": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompare(sdk.NormalizePostgresSettings),
		Description:      "Specifies custom Postgres settings as a JSON string.",
	},
	"maintenance_window_start": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 23)),
		Description:      "Specifies the hour (0-23 UTC) at which the maintenance window starts.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the Postgres instance.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW POSTGRES INSTANCES` for the given Postgres instance.",
		Elem: &schema.Resource{
			Schema: schemas.ShowPostgresInstanceSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE POSTGRES INSTANCE` for the given Postgres instance.",
		Elem: &schema.Resource{
			Schema: schemas.DescribePostgresInstanceSchema,
		},
	},
}

func PostgresInstance() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.PostgresInstances.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.PostgresInstanceResource), TrackingCreateWrapper(resources.PostgresInstance, CreatePostgresInstance)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.PostgresInstanceResource), TrackingReadWrapper(resources.PostgresInstance, ReadPostgresInstanceFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.PostgresInstanceResource), TrackingUpdateWrapper(resources.PostgresInstance, UpdatePostgresInstance)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.PostgresInstanceResource), TrackingDeleteWrapper(resources.PostgresInstance, deleteFunc)),
		Description:   "Resource used to manage Postgres instance objects. For more information, check [Postgres instance documentation](https://docs.snowflake.com/en/sql-reference/sql/create-postgres-instance).",

		Schema: postgresInstanceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.PostgresInstance, ImportName[sdk.AccountObjectIdentifier]),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.PostgresInstance, customdiff.All(
			ComputedIfAnyAttributeChanged(postgresInstanceSchema, ShowOutputAttributeName, "name", "compute_family", "storage_size_gb", "authentication_authority", "comment", "high_availability", "postgres_settings"),
			ComputedIfAnyAttributeChanged(postgresInstanceSchema, DescribeOutputAttributeName, "name", "compute_family", "storage_size_gb", "authentication_authority", "comment", "high_availability", "network_policy", "storage_integration", "postgres_version", "maintenance_window_start", "postgres_settings"),
			ComputedIfAnyAttributeChanged(postgresInstanceSchema, FullyQualifiedNameAttributeName, "name"),
		)),
		Timeouts: defaultTimeouts,
	}
}

func CreatePostgresInstance(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}

	computeFamily := d.Get("compute_family").(string)
	storageSizeGb := d.Get("storage_size_gb").(int)
	authAuthorityRaw := d.Get("authentication_authority").(string)
	authAuthority, err := sdk.ToPostgresInstanceAuthenticationAuthority(authAuthorityRaw)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreatePostgresInstanceRequest(id, computeFamily, storageSizeGb, authAuthority)
	errs := errors.Join(
		intAttributeCreateBuilder(d, "postgres_version", request.WithPostgresVersion),
		attributeMappedValueCreateBuilder(d, "network_policy", request.WithNetworkPolicy, sdk.ParseAccountObjectIdentifier),
		attributeMappedValueCreateBuilder(d, "storage_integration", request.WithStorageIntegration, sdk.ParseAccountObjectIdentifier),
		stringAttributeCreateBuilder(d, "postgres_settings", request.WithPostgresSettings),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		booleanStringAttributeCreateBuilder(d, "high_availability", request.WithHighAvailability),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.PostgresInstances.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Postgres instance %s: %w", id.FullyQualifiedName(), err))
	}

	if err := util.Retry(5, 3*time.Second, func() (error, bool) {
		_, err = client.PostgresInstances.ShowByID(ctx, id)
		if err != nil {
			log.Printf("[DEBUG] retryable operation resulted in error: %v", err)
			if errors.Is(err, sdk.ErrObjectNotFound) {
				return nil, false
			} else {
				return err, true
			}
		}
		return nil, true
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to query Postgres instance (%s) after creation, err: %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadPostgresInstanceFunc(false)(ctx, d, meta)
}

func ReadPostgresInstanceFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		pi, err := client.PostgresInstances.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Failed to query Postgres instance. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Postgres instance id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		details, err := client.PostgresInstances.DescribeDetails(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"compute_family", "compute_family", pi.ComputeFamily, pi.ComputeFamily, nil},
				outputMapping{"storage_size", "storage_size_gb", pi.StorageSize, pi.StorageSize, nil},
				outputMapping{"authentication_authority", "authentication_authority", pi.AuthenticationAuthority, pi.AuthenticationAuthority, nil},
				outputMapping{"is_ha", "high_availability", pi.IsHighlyAvailable, booleanStringFromBool(pi.IsHighlyAvailable), nil},
			); err != nil {
				return diag.FromErr(err)
			}
			var networkPolicy, storageIntegration string
			if details.NetworkPolicy != nil {
				networkPolicy = details.NetworkPolicy.Name()
			}
			if details.StorageIntegration != nil {
				storageIntegration = details.StorageIntegration.Name()
			}
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"network_policy", "network_policy", networkPolicy, networkPolicy, nil},
				outputMapping{"storage_integration", "storage_integration", storageIntegration, storageIntegration, nil},
				outputMapping{"postgres_version", "postgres_version", details.PostgresVersion, details.PostgresVersion, nil},
				outputMapping{"maintenance_window_start", "maintenance_window_start", details.MaintenanceWindowStart, details.MaintenanceWindowStart, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, postgresInstanceSchema, []string{
			"authentication_authority",
			"compute_family",
			"high_availability",
		}); err != nil {
			return diag.FromErr(err)
		}

		postgresSettings := normalizePostgresSettings(pi.PostgresSettings)
		errs := errors.Join(
			d.Set("name", pi.Name),
			d.Set("compute_family", pi.ComputeFamily),
			d.Set("storage_size_gb", pi.StorageSize),
			d.Set("authentication_authority", pi.AuthenticationAuthority),
			d.Set("high_availability", booleanStringFromBool(pi.IsHighlyAvailable)),
			setOptionalFromPtr(d, "comment", pi.Comment),
			setOptionalFromPtr(d, "postgres_settings", postgresSettings),
			d.Set("postgres_version", details.PostgresVersion),
			setOptionalFromAccountObjectIdentifierPtr(d, "network_policy", details.NetworkPolicy),
			setOptionalFromAccountObjectIdentifierPtr(d, "storage_integration", details.StorageIntegration),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.PostgresInstanceToSchema(pi)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.PostgresInstanceDetailsToSchema(details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		if details.MaintenanceWindowStart != 0 || d.Get("maintenance_window_start").(int) != IntDefault {
			if err := d.Set("maintenance_window_start", details.MaintenanceWindowStart); err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	}
}

func UpdatePostgresInstance(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithRenameTo(newId)
		if err := client.PostgresInstances.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming Postgres instance: %w", err))
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// POSTGRES_SETTINGS cannot be combined with COMPUTE_FAMILY/STORAGE_SIZE_GB/POSTGRES_VERSION/HIGH_AVAILABILITY.
	// HIGH_AVAILABILITY cannot be combined with COMPUTE_FAMILY/STORAGE_SIZE_GB/POSTGRES_VERSION/POSTGRES_SETTINGS.
	// Split SET operations into separate ALTER calls by parameter group:
	//   1. POSTGRES_SETTINGS (alone)
	//   2. Upgrade ops (COMPUTE_FAMILY, STORAGE_SIZE_GB, POSTGRES_VERSION) + non-conflicting params
	//   3. HIGH_AVAILABILITY (alone)

	// Group 1: POSTGRES_SETTINGS
	pgSettingsSet := sdk.NewPostgresInstanceSetRequest()
	pgSettingsUnset := sdk.NewPostgresInstanceUnsetRequest()
	errs := errors.Join(
		stringAttributeUpdate(d, "postgres_settings", &pgSettingsSet.PostgresSettings, &pgSettingsUnset.PostgresSettings),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(pgSettingsSet, &sdk.PostgresInstanceSetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithSet(*pgSettingsSet)
		if err := client.PostgresInstances.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error setting Postgres instance postgres_settings: %w", err))
		}
	}

	// Group 2: Upgrade ops + non-conflicting params
	set := sdk.NewPostgresInstanceSetRequest()
	unset := sdk.NewPostgresInstanceUnsetRequest()
	errs = errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		accountObjectIdentifierAttributeUpdate(d, "network_policy", &set.NetworkPolicy, &unset.NetworkPolicy),
		accountObjectIdentifierAttributeUpdate(d, "storage_integration", &set.StorageIntegration, &unset.StorageIntegration),
		intAttributeUpdateSetOnly(d, "storage_size_gb", &set.StorageSizeGb),
		intAttributeUpdateSetOnly(d, "postgres_version", &set.PostgresVersion),
		intAttributeWithSpecialDefaultUpdate(d, "maintenance_window_start", &set.MaintenanceWindowStart, &unset.MaintenanceWindowStart),
		attributeMappedValueUpdateSetOnly(d, "authentication_authority", &set.AuthenticationAuthority, sdk.ToPostgresInstanceAuthenticationAuthority),
		stringAttributeUpdateSetOnlyNotEmpty(d, "compute_family", &set.ComputeFamily),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(set, &sdk.PostgresInstanceSetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithSet(*set)
		if err := client.PostgresInstances.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error setting Postgres instance properties: %w", err))
		}
	}

	// Group 3: HIGH_AVAILABILITY
	highAvailabilitySet := sdk.NewPostgresInstanceSetRequest()
	errs = errors.Join(
		booleanStringAttributeUpdateSetOnly(d, "high_availability", &highAvailabilitySet.HighAvailability),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(highAvailabilitySet, &sdk.PostgresInstanceSetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithSet(*highAvailabilitySet)
		if err := client.PostgresInstances.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error setting Postgres instance high_availability: %w", err))
		}
	}

	// Unset operations (non-conflicting, can be sent together)
	if !reflect.DeepEqual(pgSettingsUnset, &sdk.PostgresInstanceUnsetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithUnset(*pgSettingsUnset)
		if err := client.PostgresInstances.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting Postgres instance postgres_settings: %w", err))
		}
	}

	if !reflect.DeepEqual(unset, &sdk.PostgresInstanceUnsetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithUnset(*unset)
		if err := client.PostgresInstances.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting Postgres instance properties: %w", err))
		}
	}

	return ReadPostgresInstanceFunc(false)(ctx, d, meta)
}

// normalizePostgresSettings returns nil if the postgres_settings value is
// an empty JSON object ("{}") or empty string, treating it as unset.
func normalizePostgresSettings(s *string) *string {
	if s == nil {
		return nil
	}
	normalized, err := sdk.NormalizePostgresSettings(*s)
	if err != nil || normalized == "" {
		return nil
	}
	return &normalized
}

// setOptionalFromNonEmptyStringPtr sets a key in resource data only if the
// pointer is non-nil and the pointed-to string is non-empty.
func setOptionalFromNonEmptyStringPtr(d *schema.ResourceData, key string, ptr *string) error {
	if ptr != nil && *ptr != "" {
		return d.Set(key, *ptr)
	}
	return nil
}

// setOptionalFromAccountObjectIdentifierPtr sets a key in resource data only if the
// pointer is non-nil.
func setOptionalFromAccountObjectIdentifierPtr(d *schema.ResourceData, key string, ptr *sdk.AccountObjectIdentifier) error {
	if ptr != nil {
		return d.Set(key, ptr.Name())
	}
	return nil
}
