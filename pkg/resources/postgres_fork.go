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

var postgresForkSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the forked Postgres instance; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"fork_from": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier of the source Postgres instance to fork from.",
	},
	"at_timestamp": {
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		Description:   "Specifies the timestamp for the fork point-in-time (AT TIMESTAMP).",
		ConflictsWith: []string{"at_offset", "before_timestamp", "before_offset"},
	},
	"at_offset": {
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		Description:   "Specifies the offset in seconds for the fork point-in-time (AT OFFSET).",
		ConflictsWith: []string{"at_timestamp", "before_timestamp", "before_offset"},
	},
	"before_timestamp": {
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		Description:   "Specifies the timestamp for the fork point-in-time (BEFORE TIMESTAMP).",
		ConflictsWith: []string{"at_timestamp", "at_offset", "before_offset"},
	},
	"before_offset": {
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		Description:   "Specifies the offset in seconds for the fork point-in-time (BEFORE OFFSET).",
		ConflictsWith: []string{"at_timestamp", "at_offset", "before_timestamp"},
	},
	"compute_family": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("compute_family"),
		Description:      "Specifies the compute family for the forked Postgres instance (e.g. STANDARD_M).",
	},
	"storage_size_gb": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("storage_size"),
		Description:      "Specifies the storage size in GB for the forked Postgres instance.",
	},
	"high_availability": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the Postgres instance should be configured for high availability.",
	},
	"postgres_settings": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies custom Postgres settings as a JSON string.",
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

	// Set fork_from from the origin field if available
	forkFrom := ""
	if pi.Origin != nil {
		forkFrom = *pi.Origin
	}

	errs := errors.Join(
		d.Set("name", pi.Name),
		d.Set("fork_from", forkFrom),
		d.Set("compute_family", pi.ComputeFamily),
		d.Set("storage_size_gb", pi.StorageSize),
		d.Set("high_availability", pi.IsHighlyAvailable()),
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
	if v, ok := d.GetOk("at_timestamp"); ok {
		at := sdk.NewPostgresInstanceForkAtRequest().WithTimestamp(v.(string))
		request.WithAt(*at)
	} else if v, ok := d.GetOk("at_offset"); ok {
		at := sdk.NewPostgresInstanceForkAtRequest().WithOffset(v.(string))
		request.WithAt(*at)
	}

	// Handle BEFORE time travel
	if v, ok := d.GetOk("before_timestamp"); ok {
		before := sdk.NewPostgresInstanceForkBeforeRequest().WithTimestamp(v.(string))
		request.WithBefore(*before)
	} else if v, ok := d.GetOk("before_offset"); ok {
		before := sdk.NewPostgresInstanceForkBeforeRequest().WithOffset(v.(string))
		request.WithBefore(*before)
	}

	// Handle optional fork-time parameters
	if v, ok := d.GetOk("compute_family"); ok {
		request.WithComputeFamily(v.(string))
	}
	if v, ok := d.GetOk("storage_size_gb"); ok {
		request.WithStorageSizeGb(v.(int))
	}
	if v, ok := d.GetOk("high_availability"); ok {
		request.WithHighAvailability(v.(bool))
	}
	if v, ok := d.GetOk("postgres_settings"); ok {
		request.WithPostgresSettings(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.PostgresInstances.Fork(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating forked Postgres instance %s: %w", id.FullyQualifiedName(), err))
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
				outputMapping{"is_ha", "high_availability", pi.IsHighlyAvailable(), pi.IsHighlyAvailable(), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set("name", pi.Name),
			d.Set("compute_family", pi.ComputeFamily),
			d.Set("storage_size_gb", pi.StorageSize),
			d.Set("high_availability", pi.IsHighlyAvailable()),
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
		booleanAttributeUpdateSetOnly(d, "high_availability", &highAvailabilitySet.HighAvailability),
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
