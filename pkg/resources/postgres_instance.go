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
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Postgres instance; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"compute_family": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the compute family for the Postgres instance (e.g. STANDARD_1).",
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
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the Postgres version for the instance.",
	},
	"network_policy": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the network policy to associate with the Postgres instance.",
	},
	"high_availability": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the Postgres instance should be configured for high availability.",
	},
	"storage_integration": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the storage integration for the Postgres instance.",
	},
	"postgres_settings": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies custom Postgres settings as a JSON string.",
	},
	"maintenance_window_start": {
		Type:             schema.TypeInt,
		Optional:         true,
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

	details, err := client.PostgresInstances.DescribeDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	errs := errors.Join(
		d.Set("name", pi.Name),
		d.Set("compute_family", pi.ComputeFamily),
		d.Set("storage_size_gb", pi.StorageSize),
		d.Set("authentication_authority", pi.AuthenticationAuthority),
		d.Set("high_availability", pi.IsHa),
		setOptionalFromPtr(d, "comment", pi.Comment),
		setOptionalFromPtr(d, "postgres_settings", pi.PostgresSettings),
		d.Set("postgres_version", details.PostgresVersion),
		setOptionalFromPtr(d, "network_policy", details.NetworkPolicy),
		setOptionalFromPtr(d, "storage_integration", details.StorageIntegration),
		d.Set("maintenance_window_start", details.MaintenanceWindowStart),
	)
	if errs != nil {
		return nil, errs
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
		stringAttributeCreateBuilder(d, "network_policy", request.WithNetworkPolicy),
		stringAttributeCreateBuilder(d, "storage_integration", request.WithStorageIntegration),
		stringAttributeCreateBuilder(d, "postgres_settings", request.WithPostgresSettings),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if v, ok := d.GetOk("high_availability"); ok {
		request.WithHighAvailability(v.(bool))
	}
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.PostgresInstances.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Postgres instance %s: %w", id.FullyQualifiedName(), err))
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
				outputMapping{"is_ha", "high_availability", pi.IsHa, pi.IsHa, nil},
			); err != nil {
				return diag.FromErr(err)
			}
			var networkPolicy, storageIntegration string
			if details.NetworkPolicy != nil {
				networkPolicy = *details.NetworkPolicy
			}
			if details.StorageIntegration != nil {
				storageIntegration = *details.StorageIntegration
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
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set("name", pi.Name),
			d.Set("compute_family", pi.ComputeFamily),
			d.Set("storage_size_gb", pi.StorageSize),
			d.Set("authentication_authority", pi.AuthenticationAuthority),
			d.Set("high_availability", pi.IsHa),
			setOptionalFromPtr(d, "comment", pi.Comment),
			setOptionalFromPtr(d, "postgres_settings", pi.PostgresSettings),
			d.Set("postgres_version", details.PostgresVersion),
			setOptionalFromPtr(d, "network_policy", details.NetworkPolicy),
			setOptionalFromPtr(d, "storage_integration", details.StorageIntegration),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.PostgresInstanceToSchema(pi)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.PostgresInstanceDetailsToSchema(details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("maintenance_window_start", details.MaintenanceWindowStart),
		)
		if errs != nil {
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

	set := sdk.NewPostgresInstanceSetRequest()
	unset := sdk.NewPostgresInstanceUnsetRequest()

	errs := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		stringAttributeUpdate(d, "network_policy", &set.NetworkPolicy, &unset.NetworkPolicy),
		stringAttributeUpdate(d, "storage_integration", &set.StorageIntegration, &unset.StorageIntegration),
		stringAttributeUpdate(d, "postgres_settings", &set.PostgresSettings, &unset.PostgresSettings),
		intAttributeUpdateSetOnly(d, "storage_size_gb", &set.StorageSizeGb),
		intAttributeUpdateSetOnly(d, "postgres_version", &set.PostgresVersion),
		intAttributeUpdate(d, "maintenance_window_start", &set.MaintenanceWindowStart, &unset.MaintenanceWindowStart),
		booleanAttributeUpdateSetOnly(d, "high_availability", &set.HighAvailability),
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

	if !reflect.DeepEqual(unset, &sdk.PostgresInstanceUnsetRequest{}) {
		alterReq := sdk.NewAlterPostgresInstanceRequest(id).WithUnset(*unset)
		if err := client.PostgresInstances.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting Postgres instance properties: %w", err))
		}
	}

	return ReadPostgresInstanceFunc(false)(ctx, d, meta)
}
