package resources

import (
	"context"
	"errors"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	failoverGroupCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingCreateWrapper(resources.FailoverGroup, CreateFailoverGroup))
	failoverGroupReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingReadWrapper(resources.FailoverGroup, ReadFailoverGroupFunc(true)))
	failoverGroupUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingUpdateWrapper(resources.FailoverGroup, UpdateFailoverGroup))
	failoverGroupDescriptionExt = "Resource used to manage failover groups in the source account. To manage a failover group in a target account, use `snowflake_failover_group_secondary`. For more information, check [failover group documentation](https://docs.snowflake.com/en/sql-reference/sql/create-failover-group)."
	failoverGroupImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.FailoverGroup, ImportFailoverGroup),
	}
)

// failoverGroupAlterRequestsExt holds the ALTER statements of an update in the order they must run.
// The order matters to Snowflake (RPO assurance, replication schedule, and optimized refresh depend on each other),
// and OBJECT_TYPES with ALLOWED_INTEGRATION_TYPES cannot share a SET with the other properties, so SET and UNSET are not grouped.
type failoverGroupAlterRequestsExt struct {
	steps []func(id sdk.AccountObjectIdentifier) *sdk.AlterFailoverGroupRequest
}

func (r *failoverGroupAlterRequestsExt) add(step func(id sdk.AccountObjectIdentifier) *sdk.AlterFailoverGroupRequest) {
	r.steps = append(r.steps, step)
}

func (r *failoverGroupAlterRequestsExt) set(set *sdk.FailoverGroupSetRequest) {
	r.add(func(id sdk.AccountObjectIdentifier) *sdk.AlterFailoverGroupRequest {
		return sdk.NewAlterFailoverGroupRequest(id).WithSet(*set)
	})
}

func (r *failoverGroupAlterRequestsExt) unset(unset *sdk.FailoverGroupUnsetRequest) {
	r.add(func(id sdk.AccountObjectIdentifier) *sdk.AlterFailoverGroupRequest {
		return sdk.NewAlterFailoverGroupRequest(id).WithUnset(*unset)
	})
}

func failoverGroupNewCreateRequestExt(_ *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateFailoverGroupRequest, error) {
	// OBJECT_TYPES and ALLOWED_ACCOUNTS are constructor arguments. ApplyCreateOptionals fills them.
	return sdk.NewCreateFailoverGroupRequest(id, nil, nil), nil
}

// CREATE has no COMMENT, RPO_ASSURANCE, or READ_ONLY clause. AfterCreate sets them with ALTER.
func failoverGroupApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateFailoverGroupRequest) error {
	return errors.Join(
		setValueCreate(d, "object_types", &request.ObjectTypes, toFailoverGroupObjectType),
		setValueCreate(d, "allowed_databases", &request.AllowedDatabases, toFailoverGroupAccountObjectIdentifier),
		setValueCreate(d, "allowed_external_volumes", &request.AllowedExternalVolumes, toFailoverGroupAccountObjectIdentifier),
		setValueCreate(d, "allowed_shares", &request.AllowedShares, toFailoverGroupAccountObjectIdentifier),
		setValueCreate(d, "allowed_integration_types", &request.AllowedIntegrationTypes, toFailoverGroupIntegrationType),
		setValueCreate(d, "allowed_accounts", &request.AllowedAccounts, toFailoverGroupAccountIdentifier),
		boolAttributeCreateBuilder(d, "ignore_edition_check", request.WithIgnoreEditionCheck),
		stringAttributeCreateBuilder(d, "replication_schedule", request.WithReplicationSchedule),
		boolAttributeCreateBuilder(d, "optimized_refresh", request.WithOptimizedRefresh),
		accountObjectIdentifierAttributeCreateBuilder(d, "error_integration", request.WithErrorIntegration),
	)
}

// failoverGroupAfterCreateExt runs after the id is set, so a failure here still leaves the created group in the state.
func failoverGroupAfterCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	// RPO assurance has its own ALTER. The comment and the read-only mode share the last one, so the group enters the read-only mode last.
	var sets []*sdk.FailoverGroupSetRequest
	if v, ok := d.GetOk("rpo_assurance"); ok && v.(bool) {
		sets = append(sets, sdk.NewFailoverGroupSetRequest().WithRpoAssurance(true))
	}
	last := sdk.NewFailoverGroupSetRequest()
	hasLast := false
	if v, ok := d.GetOk("comment"); ok {
		last.WithComment(v.(string))
		hasLast = true
	}
	if v, ok := d.GetOk("read_only"); ok && v.(bool) {
		last.WithReadOnly(true)
		hasLast = true
	}
	if hasLast {
		sets = append(sets, last)
	}
	for _, set := range sets {
		if err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func failoverGroupApplySetUnsetExt(d *schema.ResourceData) (*failoverGroupAlterRequestsExt, error) {
	steps := &failoverGroupAlterRequestsExt{}

	if err := failoverGroupObjectTypesStepsExt(d, steps); err != nil {
		return nil, err
	}

	if d.HasChange("rpo_assurance") && !d.Get("rpo_assurance").(bool) {
		steps.unset(sdk.NewFailoverGroupUnsetRequest().WithRpoAssurance(true))
		// RPO auto-enables Optimized Refresh and leaves it on after UNSET RPO.
		// Read stores optimized_refresh as false while RPO is on, so HasChange
		// may not fire even though SHOW still reports it enabled.
		if !d.Get("optimized_refresh").(bool) {
			steps.unset(sdk.NewFailoverGroupUnsetRequest().WithOptimizedRefresh(true))
		}
	}

	if d.HasChange("optimized_refresh") && !d.Get("optimized_refresh").(bool) {
		steps.unset(sdk.NewFailoverGroupUnsetRequest().WithOptimizedRefresh(true))
	}

	if d.HasChange("replication_schedule") {
		if schedule := d.Get("replication_schedule").(string); schedule != "" {
			steps.set(sdk.NewFailoverGroupSetRequest().WithReplicationSchedule(schedule))
		} else {
			steps.unset(sdk.NewFailoverGroupUnsetRequest().WithReplicationSchedule(true))
		}
	}

	if d.HasChange("optimized_refresh") && d.Get("optimized_refresh").(bool) {
		steps.set(sdk.NewFailoverGroupSetRequest().WithOptimizedRefresh(true))
	}

	if d.HasChange("rpo_assurance") && d.Get("rpo_assurance").(bool) {
		steps.set(sdk.NewFailoverGroupSetRequest().WithRpoAssurance(true))
	}

	if d.HasChange("comment") {
		if comment := d.Get("comment").(string); comment != "" {
			steps.set(sdk.NewFailoverGroupSetRequest().WithComment(comment))
		} else {
			steps.unset(sdk.NewFailoverGroupUnsetRequest().WithComment(true))
		}
	}

	if d.HasChange("error_integration") {
		if raw := d.Get("error_integration").(string); raw != "" {
			errorIntegration, err := sdk.ParseAccountObjectIdentifier(raw)
			if err != nil {
				return nil, err
			}
			steps.set(sdk.NewFailoverGroupSetRequest().WithErrorIntegration(errorIntegration))
		} else {
			steps.unset(sdk.NewFailoverGroupUnsetRequest().WithErrorIntegration(true))
		}
	}

	if d.HasChange("read_only") {
		if d.Get("read_only").(bool) {
			steps.set(sdk.NewFailoverGroupSetRequest().WithReadOnly(true))
		} else {
			steps.unset(sdk.NewFailoverGroupUnsetRequest().WithReadOnly(true))
		}
	}

	return steps, errors.Join(
		failoverGroupAddRemoveStepsExt(d, steps, "allowed_databases", toFailoverGroupAccountObjectIdentifier,
			func(ids []sdk.AccountObjectIdentifier) *sdk.FailoverGroupRemoveRequest {
				return sdk.NewFailoverGroupRemoveRequest().WithAllowedDatabases(ids)
			},
			func(ids []sdk.AccountObjectIdentifier) *sdk.FailoverGroupAddRequest {
				return sdk.NewFailoverGroupAddRequest().WithAllowedDatabases(ids)
			}),
		failoverGroupAddRemoveStepsExt(d, steps, "allowed_shares", toFailoverGroupAccountObjectIdentifier,
			func(ids []sdk.AccountObjectIdentifier) *sdk.FailoverGroupRemoveRequest {
				return sdk.NewFailoverGroupRemoveRequest().WithAllowedShares(ids)
			},
			func(ids []sdk.AccountObjectIdentifier) *sdk.FailoverGroupAddRequest {
				return sdk.NewFailoverGroupAddRequest().WithAllowedShares(ids)
			}),
		failoverGroupAddRemoveStepsExt(d, steps, "allowed_external_volumes", toFailoverGroupAccountObjectIdentifier,
			func(ids []sdk.AccountObjectIdentifier) *sdk.FailoverGroupRemoveRequest {
				return sdk.NewFailoverGroupRemoveRequest().WithAllowedExternalVolumes(ids)
			},
			func(ids []sdk.AccountObjectIdentifier) *sdk.FailoverGroupAddRequest {
				return sdk.NewFailoverGroupAddRequest().WithAllowedExternalVolumes(ids)
			}),
		failoverGroupAddRemoveStepsExt(d, steps, "allowed_accounts", toFailoverGroupAccountIdentifier,
			func(ids []sdk.AccountIdentifier) *sdk.FailoverGroupRemoveRequest {
				return sdk.NewFailoverGroupRemoveRequest().WithAllowedAccounts(ids)
			},
			func(ids []sdk.AccountIdentifier) *sdk.FailoverGroupAddRequest {
				return sdk.NewFailoverGroupAddRequest().WithAllowedAccounts(ids)
			}),
	)
}

// failoverGroupObjectTypesStepsExt is one SET with OBJECT_TYPES and ALLOWED_INTEGRATION_TYPES. The schedule and the other properties cannot share it (GH #2544).
func failoverGroupObjectTypesStepsExt(d *schema.ResourceData, steps *failoverGroupAlterRequestsExt) error {
	set := sdk.NewFailoverGroupSetRequest()
	if d.HasChange("object_types") {
		objectTypes, err := failoverGroupSetValues(d, "object_types", toFailoverGroupObjectType)
		if err != nil {
			return err
		}
		set.WithObjectTypes(objectTypes)
		if slices.Contains(objectTypes, sdk.FailoverGroupObjectTypeIntegrations) {
			allowedIntegrationTypes, err := failoverGroupSetValues(d, "allowed_integration_types", toFailoverGroupIntegrationType)
			if err != nil {
				return err
			}
			set.WithAllowedIntegrationTypes(allowedIntegrationTypes)
		}
	}
	if d.HasChange("allowed_integration_types") {
		allowedIntegrationTypes, err := failoverGroupSetValues(d, "allowed_integration_types", toFailoverGroupIntegrationType)
		if err != nil {
			return err
		}
		set.WithAllowedIntegrationTypes(allowedIntegrationTypes)
	}
	if set.ObjectTypes != nil || set.AllowedIntegrationTypes != nil {
		steps.set(set)
	}
	return nil
}

// failoverGroupAddRemoveStepsExt removes the values dropped from the config before it adds the new ones. ADD and REMOVE take one list per statement.
func failoverGroupAddRemoveStepsExt[T comparable](
	d *schema.ResourceData,
	steps *failoverGroupAlterRequestsExt,
	key string,
	mapper func(any) (T, error),
	remove func([]T) *sdk.FailoverGroupRemoveRequest,
	add func([]T) *sdk.FailoverGroupAddRequest,
) error {
	if !d.HasChange(key) {
		return nil
	}
	oldRaw, newRaw := d.GetChange(key)
	oldValues, err := collections.MapErr(oldRaw.(*schema.Set).List(), mapper)
	if err != nil {
		return err
	}
	newValues, err := collections.MapErr(newRaw.(*schema.Set).List(), mapper)
	if err != nil {
		return err
	}
	removed := slices.DeleteFunc(slices.Clone(oldValues), func(v T) bool { return slices.Contains(newValues, v) })
	added := slices.DeleteFunc(slices.Clone(newValues), func(v T) bool { return slices.Contains(oldValues, v) })
	if len(removed) > 0 {
		steps.add(func(id sdk.AccountObjectIdentifier) *sdk.AlterFailoverGroupRequest {
			return sdk.NewAlterFailoverGroupRequest(id).WithRemove(*remove(removed))
		})
	}
	if len(added) > 0 {
		steps.add(func(id sdk.AccountObjectIdentifier) *sdk.AlterFailoverGroupRequest {
			return sdk.NewAlterFailoverGroupRequest(id).WithAdd(*add(added))
		})
	}
	return nil
}

func failoverGroupSetValues[T any](d *schema.ResourceData, key string, mapper func(any) (T, error)) ([]T, error) {
	return collections.MapErr(d.Get(key).(*schema.Set).List(), mapper)
}

func failoverGroupAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, steps *failoverGroupAlterRequestsExt) error {
	client := meta.(*provider.Context).Client
	for _, step := range steps.steps {
		if err := client.FailoverGroups.Alter(ctx, step(id)); err != nil {
			return err
		}
	}
	return nil
}

// SHOW lists the allowed databases and shares only through separate commands. SHOW returns NULL in read_only whatever is set,
// and no command lists the allowed external volumes or the error integration, so those keep the configured values.
func failoverGroupSetConfigFieldsExt(ctx context.Context, d *schema.ResourceData, meta any, group *sdk.FailoverGroup) error {
	client := meta.(*provider.Context).Client
	allowedAccounts, err := allowedAccountsWithoutCurrent(ctx, meta, group.AllowedAccounts)
	if err != nil {
		return err
	}
	databases, err := client.FailoverGroups.ShowDatabases(ctx, group.ID())
	if err != nil {
		return err
	}
	shares, err := client.FailoverGroups.ShowShares(ctx, group.ID())
	if err != nil {
		return err
	}

	// RPO Assurance enables Optimized Refresh automatically. Do not write SHOW's
	// implied-true value back when RPO is on, or the config without optimized_refresh will drift.
	optimizedRefresh := group.IsOptimizedRefreshEnabled && !group.RpoAssurance

	return errors.Join(
		d.Set("object_types", enumStrings(group.ObjectTypes)),
		d.Set("allowed_integration_types", enumStrings(group.AllowedIntegrationTypes)),
		d.Set("allowed_accounts", collections.Map(allowedAccounts, sdk.AccountIdentifier.Name)),
		d.Set("allowed_databases", collections.Map(databases, sdk.AccountObjectIdentifier.Name)),
		d.Set("allowed_shares", collections.Map(shares, sdk.AccountObjectIdentifier.Name)),
		d.Set("comment", group.Comment),
		d.Set("replication_schedule", group.ReplicationSchedule),
		d.Set("rpo_assurance", group.RpoAssurance),
		d.Set("optimized_refresh", optimizedRefresh),
		setOptionalValueWithMapping(d, "error_integration", group.ErrorIntegration, (*sdk.AccountObjectIdentifier).Name),
		d.Set("read_only", util.DerefOrZero(group.IsPrimaryReadOnly)),
	)
}

func failoverGroupSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, _ *sdk.FailoverGroup, _ bool) error {
	return d.Set("name", id.Name())
}

func toFailoverGroupObjectType(value any) (sdk.FailoverGroupObjectType, error) {
	return sdk.ToFailoverGroupObjectType(value.(string))
}

func toFailoverGroupIntegrationType(value any) (sdk.FailoverGroupIntegrationType, error) {
	return sdk.ToFailoverGroupIntegrationType(value.(string))
}

func toFailoverGroupAccountObjectIdentifier(value any) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(value.(string))
}

func toFailoverGroupAccountIdentifier(value any) (sdk.AccountIdentifier, error) {
	return sdk.ParseAccountIdentifier(value.(string))
}
