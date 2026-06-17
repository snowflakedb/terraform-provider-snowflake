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

var postgresForkSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the forked Postgres instance; must be unique for your account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"fork_from": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the identifier of the source Postgres instance to fork from.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"at": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Description: "Specifies the point-in-time for the fork using AT. Exactly one of `timestamp` or `offset` must be set.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"timestamp": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Specifies an exact date and time to use for Time Travel. The value must be explicitly cast to a TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ, or TIMESTAMP_TZ data type.",
					ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset"},
				},
				"offset": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Specifies the difference in seconds from the current time to use for Time Travel, in the form -N where N can be an integer or arithmetic expression (e.g. -120 is 120 seconds, -30*60 is 1800 seconds or 30 minutes).",
					ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset"},
				},
			},
		},
		ConflictsWith: []string{"before"},
	},
	"before": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Description: "Specifies the point-in-time for the fork using BEFORE. Exactly one of `timestamp` or `offset` must be set.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"timestamp": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Specifies an exact date and time to use for Time Travel. The value must be explicitly cast to a TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ, or TIMESTAMP_TZ data type.",
					ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset"},
				},
				"offset": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Specifies the difference in seconds from the current time to use for Time Travel, in the form -N where N can be an integer or arithmetic expression (e.g. -120 is 120 seconds, -30*60 is 1800 seconds or 30 minutes).",
					ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset"},
				},
			},
		},
		ConflictsWith: []string{"at"},
	},
	"compute_family": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToPostgresInstanceComputeFamily),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToPostgresInstanceComputeFamily),
		Description:      fmt.Sprintf("Specifies the compute family for the forked Postgres instance. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllPostgresInstanceComputeFamilies)),
	},
	"storage_size_gb": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("storage_size"),
		Description:      "Specifies the storage size in GB for the forked Postgres instance.",
	},
	"high_availability": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_ha"),
		Description:      booleanStringFieldDescription("Specifies whether the Postgres instance should be configured for high availability."),
		Default:          BooleanDefault,
	},
	"postgres_settings": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompare(sdk.NormalizePostgresSettings),
		Description:      "Specifies custom Postgres settings as a JSON string.",
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

func PostgresFork() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.PostgresInstances.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.PostgresForkResource), TrackingCreateWrapper(resources.PostgresFork, CreatePostgresFork)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.PostgresForkResource), TrackingReadWrapper(resources.PostgresFork, ReadPostgresForkFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.PostgresForkResource), TrackingUpdateWrapper(resources.PostgresFork, UpdatePostgresFork)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.PostgresForkResource), TrackingDeleteWrapper(resources.PostgresFork, deleteFunc)),
		Description:   "Resource used to manage forked Postgres instance objects. For more information, check [Postgres instance documentation](https://docs.snowflake.com/en/sql-reference/sql/create-postgres-instance).",

		Schema: postgresForkSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.PostgresFork, ImportPostgresFork),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.PostgresFork, customdiff.All(
			ComputedIfAnyAttributeChanged(postgresForkSchema, ShowOutputAttributeName, "name", "compute_family", "storage_size_gb", "comment", "high_availability", "postgres_settings"),
			ComputedIfAnyAttributeChanged(postgresForkSchema, DescribeOutputAttributeName, "name", "compute_family", "storage_size_gb", "comment", "high_availability", "postgres_settings"),
			ComputedIfAnyAttributeChanged(postgresForkSchema, FullyQualifiedNameAttributeName, "name"),
		)),
		Timeouts: defaultTimeouts,
	}
}

func ImportPostgresFork(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	pi, err := client.PostgresInstances.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if pi.Origin == nil || *pi.Origin == "" {
		return nil, fmt.Errorf("postgres instance %s is not a fork (origin is empty); use the snowflake_postgres_instance resource to import non-fork instances", id.FullyQualifiedName())
	}
	forkFrom := *pi.Origin

	errs := errors.Join(
		d.Set("name", pi.Name),
		d.Set("fork_from", forkFrom),
		d.Set("compute_family", pi.ComputeFamily),
		d.Set("storage_size_gb", pi.StorageSize),
		d.Set("high_availability", booleanStringFromBool(pi.IsHighlyAvailable)),
		setOptionalFromPtr(d, "comment", pi.Comment),
		setOptionalFromPtr(d, "postgres_settings", normalizePostgresSettings(pi.PostgresSettings)),
	)
	if errs != nil {
		return nil, errs
	}

	return []*schema.ResourceData{d}, nil
}

func CreatePostgresFork(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}

	forkFromRaw := d.Get("fork_from").(string)
	forkFromId, err := sdk.ParseAccountObjectIdentifier(forkFromRaw)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewForkPostgresInstanceRequest(id, forkFromId)

	// Handle AT time travel
	if v := d.Get("at").([]any); len(v) > 0 {
		atConfig := v[0].(map[string]any)
		at := sdk.NewPostgresInstanceForkAtRequest()
		if ts, ok := atConfig["timestamp"].(string); ok && len(ts) > 0 {
			at.WithTimestamp(ts)
		}
		if offset, ok := atConfig["offset"].(string); ok && len(offset) > 0 {
			at.WithOffset(offset)
		}
		request.WithAt(*at)
	}

	// Handle BEFORE time travel
	if v := d.Get("before").([]any); len(v) > 0 {
		beforeConfig := v[0].(map[string]any)
		before := sdk.NewPostgresInstanceForkBeforeRequest()
		if ts, ok := beforeConfig["timestamp"].(string); ok && len(ts) > 0 {
			before.WithTimestamp(ts)
		}
		if offset, ok := beforeConfig["offset"].(string); ok && len(offset) > 0 {
			before.WithOffset(offset)
		}
		request.WithBefore(*before)
	}

	// Handle optional fork-time parameters
	errs := errors.Join(
		stringAttributeCreateBuilder(d, "compute_family", request.WithComputeFamily),
		booleanStringAttributeCreateBuilder(d, "high_availability", request.WithHighAvailability),
		stringAttributeCreateBuilder(d, "postgres_settings", request.WithPostgresSettings),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if v, ok := d.GetOk("storage_size_gb"); ok {
		request.WithStorageSizeGb(v.(int))
	}
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.PostgresInstances.Fork(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating forked Postgres instance %s: %w", id.FullyQualifiedName(), err))
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
	return ReadPostgresForkFunc(false)(ctx, d, meta)
}

func ReadPostgresForkFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
				outputMapping{"is_ha", "high_availability", pi.IsHighlyAvailable, booleanStringFromBool(pi.IsHighlyAvailable), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, postgresForkSchema, []string{
			"compute_family",
			"high_availability",
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set("name", pi.Name),
			d.Set("compute_family", pi.ComputeFamily),
			d.Set("storage_size_gb", pi.StorageSize),
			d.Set("high_availability", booleanStringFromBool(pi.IsHighlyAvailable)),
			setOptionalFromPtr(d, "comment", pi.Comment),
			setOptionalFromPtr(d, "postgres_settings", normalizePostgresSettings(pi.PostgresSettings)),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.PostgresInstanceToSchema(pi)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.PostgresInstanceDetailsToSchema(details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func UpdatePostgresFork(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	// POSTGRES_SETTINGS cannot be combined with COMPUTE_FAMILY/STORAGE_SIZE_GB/HIGH_AVAILABILITY.
	// HIGH_AVAILABILITY cannot be combined with COMPUTE_FAMILY/STORAGE_SIZE_GB/POSTGRES_SETTINGS.
	// Split SET operations into separate ALTER calls by parameter group:
	//   1. POSTGRES_SETTINGS (alone)
	//   2. Upgrade ops (COMPUTE_FAMILY, STORAGE_SIZE_GB) + non-conflicting params
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
		intAttributeUpdateSetOnly(d, "storage_size_gb", &set.StorageSizeGb),
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

	// Unset operations
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

	return ReadPostgresForkFunc(false)(ctx, d, meta)
}
