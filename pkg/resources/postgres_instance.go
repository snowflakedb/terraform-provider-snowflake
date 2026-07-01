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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var postgresInstanceSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
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
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.NormalizePostgresSettings), schema.SchemaDiffSuppressFunc(func(k, old, new string, d *schema.ResourceData) bool {
			return old == "{}" && new == ""
		})),
		Description: "Specifies custom Postgres settings as a JSON string.",
	},
	"maintenance_window_start": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 23)),
		Description:      "Specifies the hour (0-23 UTC) at which the maintenance window starts.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("maintenance_window_start"),
	},
	"comment": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("comment"),
		Description:      "Specifies a comment for the Postgres instance.",
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
			StateContext: TrackingImportWrapper(resources.PostgresInstance, ImportPostgresInstance),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.PostgresInstance, customdiff.All(
			ComputedIfAnyAttributeChanged(postgresInstanceSchema, ShowOutputAttributeName, "name", "compute_family", "storage_size_gb", "authentication_authority", "comment", "high_availability", "postgres_settings"),
			ComputedIfAnyAttributeChanged(postgresInstanceSchema, DescribeOutputAttributeName, "name", "compute_family", "storage_size_gb", "authentication_authority", "comment", "high_availability", "network_policy", "storage_integration", "postgres_version", "maintenance_window_start", "postgres_settings"),
			ComputedIfAnyAttributeChanged(postgresInstanceSchema, FullyQualifiedNameAttributeName, "name"),
		)),
		Timeouts: defaultTimeouts,
	}
}

func ImportPostgresInstance(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	pi, err := client.PostgresInstances.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if pi.Origin != nil && *pi.Origin != "" {
		return nil, fmt.Errorf("postgres instance %s is a fork (origin: %s); use the snowflake_postgres_fork resource to import fork instances", id.FullyQualifiedName(), *pi.Origin)
	}

	details, err := client.PostgresInstances.DescribeDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	errs := errors.Join(
		d.Set("name", pi.Name),
		d.Set("compute_family", pi.ComputeFamily),
		d.Set("storage_size_gb", pi.StorageSize),
		d.Set("authentication_authority", pi.AuthenticationAuthority),
		d.Set("high_availability", booleanStringFromBool(pi.IsHighlyAvailable)),
		d.Set("postgres_version", details.PostgresVersion),
		d.Set("maintenance_window_start", details.MaintenanceWindowStart),
	)
	if errs != nil {
		return nil, errs
	}

	if pi.Comment != nil && *pi.Comment != "" {
		if err := d.Set("comment", *pi.Comment); err != nil {
			return nil, err
		}
	}
	if details.NetworkPolicy != nil {
		if err := d.Set("network_policy", details.NetworkPolicy.Name()); err != nil {
			return nil, err
		}
	}
	if details.StorageIntegration != nil {
		if err := d.Set("storage_integration", details.StorageIntegration.Name()); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
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
		// postgres_settings omitted from CREATE — Snowflake returns an internal error when
		// POSTGRES_SETTINGS is present in CREATE POSTGRES INSTANCE. Set post-create via ALTER.
		// maintenance_window_start omitted from CREATE — not supported in CREATE DDL.
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		booleanStringAttributeCreateBuilder(d, "high_availability", request.WithHighAvailability),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if _, err := client.PostgresInstances.CreateSafely(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Postgres instance %s: %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	// postgres_settings and maintenance_window_start cannot be set at CREATE time — apply via ALTER.
	postCreateSet := sdk.NewPostgresInstanceSetRequest()
	needsPostCreateSet := false
	if pgSettings := d.Get("postgres_settings").(string); pgSettings != "" {
		postCreateSet.PostgresSettings = &pgSettings
		needsPostCreateSet = true
	}
	if mws := d.Get("maintenance_window_start").(int); mws != IntDefault {
		postCreateSet.MaintenanceWindowStart = &mws
		needsPostCreateSet = true
	}
	if needsPostCreateSet {
		if err := client.PostgresInstances.AlterSafely(ctx, sdk.NewAlterPostgresInstanceRequest(id).WithSet(*postCreateSet)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting post-create properties for Postgres instance %s: %w", id.FullyQualifiedName(), err))
		}
	}

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
			var comment string
			if pi.Comment != nil {
				comment = *pi.Comment
			}
			if err = handleExternalChangesToObjectInShow(
				d,
				outputMapping{"compute_family", "compute_family", pi.ComputeFamily, pi.ComputeFamily, nil},
				outputMapping{"storage_size", "storage_size_gb", pi.StorageSize, pi.StorageSize, nil},
				outputMapping{"authentication_authority", "authentication_authority", pi.AuthenticationAuthority, pi.AuthenticationAuthority, nil},
				outputMapping{"is_ha", "high_availability", pi.IsHighlyAvailable, booleanStringFromBool(pi.IsHighlyAvailable), nil},
				outputMapping{"comment", "comment", comment, comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}

			var networkPolicy, storageIntegration, postgresSettings string
			if details.NetworkPolicy != nil {
				networkPolicy = details.NetworkPolicy.Name()
			}
			if details.StorageIntegration != nil {
				storageIntegration = details.StorageIntegration.Name()
			}
			if details.PostgresSettings != nil {
				postgresSettings = *details.PostgresSettings
			}
			if err = handleExternalChangesToObjectInFlatDescribe(
				d,
				outputMapping{"network_policy", "network_policy", networkPolicy, networkPolicy, nil},
				outputMapping{"storage_integration", "storage_integration", storageIntegration, storageIntegration, nil},
				outputMapping{"postgres_version", "postgres_version", details.PostgresVersion, details.PostgresVersion, nil},
				outputMapping{"postgres_settings", "postgres_settings", postgresSettings, postgresSettings, nil},
				outputMapping{"maintenance_window_start", "maintenance_window_start", details.MaintenanceWindowStart, details.MaintenanceWindowStart, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, postgresInstanceSchema, []string{
			"authentication_authority",
			"compute_family",
			"high_availability",
			"storage_size_gb",
			"postgres_version",
			"postgres_settings",
			"network_policy",
			"storage_integration",
			"maintenance_window_start",
			"comment",
		}); err != nil {
			return diag.FromErr(err)
		}

		if errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.PostgresInstanceToSchema(pi)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.PostgresInstanceDetailsToSchema(details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		); errs != nil {
			return diag.FromErr(errs)
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

	// POSTGRES_SETTINGS cannot be combined with COMPUTE_FAMILY/STORAGE_SIZE_GB/POSTGRES_VERSION/HIGH_AVAILABILITY.
	// HIGH_AVAILABILITY cannot be combined with COMPUTE_FAMILY/STORAGE_SIZE_GB/POSTGRES_VERSION/POSTGRES_SETTINGS.
	// Split SET operations into separate ALTER calls by parameter group:
	//   1. POSTGRES_SETTINGS (alone)
	//   2. Upgrade ops (COMPUTE_FAMILY, STORAGE_SIZE_GB, POSTGRES_VERSION) + non-conflicting params
	//   3. HIGH_AVAILABILITY (alone)
	// UNSET fields have no mutual conflicts and are sent in a single combined call.
	unset := sdk.NewPostgresInstanceUnsetRequest()

	// Group 1: POSTGRES_SETTINGS
	pgSettingsSet := sdk.NewPostgresInstanceSetRequest()
	errs := errors.Join(
		stringAttributeUpdate(d, "postgres_settings", &pgSettingsSet.PostgresSettings, &unset.PostgresSettings),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(pgSettingsSet, &sdk.PostgresInstanceSetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithSet(*pgSettingsSet)
		if err := client.PostgresInstances.AlterSafely(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error setting Postgres instance postgres_settings: %w", err))
		}
	}

	// Group 2: Upgrade ops + non-conflicting params
	set := sdk.NewPostgresInstanceSetRequest()
	errs = errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		accountObjectIdentifierAttributeUpdate(d, "network_policy", &set.NetworkPolicy, &unset.NetworkPolicy),
		accountObjectIdentifierAttributeUpdate(d, "storage_integration", &set.StorageIntegration, &unset.StorageIntegration),
		intAttributeUpdateSetOnly(d, "storage_size_gb", &set.StorageSizeGb),
		intAttributeWithSpecialDefaultUpdate(d, "maintenance_window_start", &set.MaintenanceWindowStart, &unset.MaintenanceWindowStart),
		attributeMappedValueUpdateSetOnly(d, "authentication_authority", &set.AuthenticationAuthority, sdk.ToPostgresInstanceAuthenticationAuthority),
		stringAttributeUpdateSetOnlyNotEmpty(d, "compute_family", &set.ComputeFamily),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	// No UNSET for postgres_version; SET the Snowflake default (18) when config removes it.
	// Skip when the effective new value equals the old state value: a no-op SET POSTGRES_VERSION
	// triggers a multi-minute async Snowflake operation that blocks the subsequent HIGH_AVAILABILITY ALTER.
	if d.HasChange("postgres_version") {
		old, _ := d.GetChange("postgres_version")
		effectiveNew := d.Get("postgres_version").(int)
		if effectiveNew == 0 {
			effectiveNew = 18
		}
		if old.(int) != effectiveNew {
			set.PostgresVersion = sdk.Int(effectiveNew)
		}
	}

	if !reflect.DeepEqual(set, &sdk.PostgresInstanceSetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithSet(*set)
		if err := client.PostgresInstances.AlterSafely(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error setting Postgres instance properties: %w", err))
		}
	}

	// Group 3: HIGH_AVAILABILITY (no UNSET; setting Snowflake default false when config removes it).
	// Skip when the effective new value equals the old state value: a no-op SET HIGH_AVAILABILITY
	// triggered after a POSTGRES_VERSION/COMPUTE_FAMILY/STORAGE_SIZE_GB ALTER will fail with
	// "must be complete before issuing ALTER SET HIGH_AVAILABILITY".
	highAvailabilitySet := sdk.NewPostgresInstanceSetRequest()
	if d.HasChange("high_availability") {
		var effectiveNew bool
		if v := d.Get("high_availability").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			effectiveNew = parsed
		}
		oldRaw, _ := d.GetChange("high_availability")
		var oldBool bool
		if s := oldRaw.(string); s != BooleanDefault {
			parsed, err := booleanStringToBool(s)
			if err == nil {
				oldBool = parsed
			}
		}
		if oldBool != effectiveNew {
			highAvailabilitySet.HighAvailability = sdk.Bool(effectiveNew)
		}
	}

	if !reflect.DeepEqual(highAvailabilitySet, &sdk.PostgresInstanceSetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithSet(*highAvailabilitySet)
		if err := client.PostgresInstances.AlterSafely(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error setting Postgres instance high_availability: %w", err))
		}
	}

	if !reflect.DeepEqual(unset, &sdk.PostgresInstanceUnsetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithUnset(*unset)
		// AlterSafely waits for READY and retries on "must be complete" errors.
		// We intentionally do not wait for DESCRIBE propagation after UNSET: both SHOW and
		// DESCRIBE can take >10 minutes to reflect UNSET, exceeding the update timeout.
		// ReadPostgresInstanceFunc reads comment/postgres_settings from the desired config on
		// the post-update pass (withExternalChangesMarking=false), so the resource state is
		// correct immediately without a propagation wait.
		if err := client.PostgresInstances.AlterSafely(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting Postgres instance properties: %w", err))
		}
	}

	return ReadPostgresInstanceFunc(false)(ctx, d, meta)
}
