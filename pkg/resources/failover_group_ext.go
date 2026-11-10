package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	failoverGroupCreateExt = PreviewFeatureCreateContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingCreateWrapper(resources.FailoverGroup, CreateFailoverGroup))
	failoverGroupReadExt   = PreviewFeatureReadContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingReadWrapper(resources.FailoverGroup, ReadFailoverGroup))
	failoverGroupUpdateExt = PreviewFeatureUpdateContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingUpdateWrapper(resources.FailoverGroup, UpdateFailoverGroup))
	// Delete decodes the id with the error-returning variant, while Read and Update use failoverGroupParseIdExt (the original behavior of each).
	failoverGroupDeleteExt = PreviewFeatureDeleteContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingDeleteWrapper(resources.FailoverGroup, ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.AccountObjectIdentifier],
		failoverGroupDropSafelyInSdkExt,
	)))
	failoverGroupDescriptionExt = docs.WithPreviewFeatureMark()
	failoverGroupImporterExt    = &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
)

// =============================================================================
// Request / mapping / Read-Update helpers
// =============================================================================

func failoverGroupParseIdFromConfigExt(d *schema.ResourceData) (sdk.AccountObjectIdentifier, error) {
	return sdk.NewAccountObjectIdentifier(d.Get("name").(string)), nil
}

// failoverGroupParseIdExt keeps the original Read and Update behavior (the decoding panics on a non-account-object identifier).
func failoverGroupParseIdExt(id string) (sdk.AccountObjectIdentifier, error) {
	return helpers.DecodeSnowflakeIDLegacy(id).(sdk.AccountObjectIdentifier), nil
}

// failoverGroupBeforeCreateExt creates the secondary failover group (from_replica); the primary one is created by the regular Create flow.
func failoverGroupBeforeCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (bool, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	name := id.Name()

	if v, ok := d.GetOk("from_replica"); ok {
		fromReplica := v.([]any)[0].(map[string]any)
		organizationName := fromReplica["organization_name"].(string)
		sourceAccountName := fromReplica["source_account_name"].(string)
		sourceFailoverGroupName := fromReplica["name"].(string)

		primaryFailoverGroupID := sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(organizationName, sourceAccountName), sdk.NewAccountObjectIdentifier(sourceFailoverGroupName))
		err := client.FailoverGroups.CreateSecondaryReplicationGroup(ctx, sdk.NewCreateSecondaryReplicationGroupFailoverGroupRequest(id, primaryFailoverGroupID))
		if err != nil {
			return true, diag.FromErr(err)
		}
		d.SetId(name)
		return true, ReadFailoverGroup(ctx, d, meta)
	}

	return false, nil
}

func failoverGroupNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateFailoverGroupRequest, error) {
	if _, ok := d.GetOk("object_types"); !ok {
		return nil, errors.New("object_types field is required when from_replica is not set")
	}
	objectTypesList := expandStringList(d.Get("object_types").(*schema.Set).List())
	objectTypes := make([]sdk.PluralObjectType, len(objectTypesList))
	for i, v := range objectTypesList {
		objectType, err := sdk.ToPluralObjectType(v)
		if err != nil {
			return nil, err
		}
		objectTypes[i] = objectType
	}

	if _, ok := d.GetOk("allowed_accounts"); !ok {
		return nil, errors.New("allowed_accounts field is required when from_replica is not set")
	}
	aaList := expandStringList(d.Get("allowed_accounts").(*schema.Set).List())
	allowedAccounts := make([]sdk.AccountIdentifier, len(aaList))
	for i, v := range aaList {
		parts := strings.Split(v, ".")
		if len(parts) != 2 {
			return nil, fmt.Errorf("allowed_account %s cannot be an account locator and must be of the format <org_name>.<target_account_name>", allowedAccounts[i])
		}
		allowedAccounts[i] = sdk.NewAccountIdentifier(parts[0], parts[1])
	}

	req := sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts)

	return req, nil
}

func failoverGroupApplyCreateOptionalsExt(d *schema.ResourceData, req *sdk.CreateFailoverGroupRequest) error {
	if v, ok := d.GetOk("allowed_databases"); ok {
		allowedDatabasesList := expandStringList(v.(*schema.Set).List())
		allowedDatabaseIdentifiers := make([]sdk.AccountObjectIdentifier, len(allowedDatabasesList))
		for i, v := range allowedDatabasesList {
			allowedDatabaseIdentifiers[i] = sdk.NewAccountObjectIdentifier(v)
		}
		req.WithAllowedDatabases(allowedDatabaseIdentifiers)
	}

	if v, ok := d.GetOk("allowed_shares"); ok {
		allowedSharesList := expandStringList(v.(*schema.Set).List())
		allowedShareIdentifiers := make([]sdk.AccountObjectIdentifier, len(allowedSharesList))
		for i, v := range allowedSharesList {
			allowedShareIdentifiers[i] = sdk.NewAccountObjectIdentifier(v)
		}
		req.WithAllowedShares(allowedShareIdentifiers)
	}

	if v, ok := d.GetOk("allowed_integration_types"); ok {
		allowedIntegrationTypesList := expandStringList(v.(*schema.Set).List())
		allowedIntegrationTypes := make([]sdk.IntegrationType, len(allowedIntegrationTypesList))
		for i, v := range allowedIntegrationTypesList {
			allowedIntegrationTypes[i] = sdk.IntegrationType(v)
		}
		req.WithAllowedIntegrationTypes(allowedIntegrationTypes)
	}

	if v, ok := d.GetOk("ignore_edition_check"); ok {
		req.WithIgnoreEditionCheck(v.(bool))
	}

	if v, ok := d.GetOk("replication_schedule"); ok {
		replicationSchedule := v.([]any)[0].(map[string]any)
		if v, ok := replicationSchedule["cron"]; ok {
			c := v.([]any)
			if len(c) > 0 {
				cron := c[0].(map[string]any)
				cronExpression := cron["expression"].(string)
				if v, ok := cron["time_zone"]; ok {
					timeZone := v.(string)
					cronExpression += " " + timeZone
				}
				req.WithReplicationSchedule("USING CRON " + cronExpression)
			}
		}
		if v, ok := replicationSchedule["interval"]; ok {
			interval := v.(int)
			if interval > 1 {
				req.WithReplicationSchedule(fmt.Sprintf("%d MINUTE", interval))
			}
		}
	}

	if v, ok := d.GetOk("optimized_refresh"); ok {
		req.WithOptimizedRefresh(v.(bool))
	}

	return nil
}

// failoverGroupAfterCreateExt runs after the id is set, so a failure here still leaves the created group in the state.
func failoverGroupAfterCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if v, ok := d.GetOk("rpo_assurance"); ok && v.(bool) {
		if err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithRpoAssurance(true))); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func failoverGroupSetConfigFieldsExt(ctx context.Context, d *schema.ResourceData, meta any, failoverGroup *sdk.FailoverGroup) error {
	client := meta.(*provider.Context).Client
	// The id is decoded the same way as in Read (not taken from the SHOW output).
	id, err := failoverGroupParseIdExt(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("name", failoverGroup.Name); err != nil {
		return err
	}
	if _, ok := d.GetOk("from_replica"); ok {
		return nil
	}

	if err := d.Set("rpo_assurance", failoverGroup.RpoAssurance); err != nil {
		return err
	}
	// RPO Assurance enables Optimized Refresh automatically. Do not write SHOW's
	// implied-true value back when RPO is on, or config default false will drift.
	if failoverGroup.RpoAssurance {
		if err := d.Set("optimized_refresh", false); err != nil {
			return err
		}
	} else if err := d.Set("optimized_refresh", failoverGroup.IsOptimizedRefreshEnabled); err != nil {
		return err
	}

	replicationSchedule := failoverGroup.ReplicationSchedule
	if replicationSchedule != "" {
		if strings.Contains(replicationSchedule, "MINUTE") {
			interval, err := strconv.Atoi(strings.TrimSuffix(replicationSchedule, " MINUTE"))
			if err != nil {
				return err
			}
			err = d.Set("replication_schedule", []any{
				map[string]any{
					"interval": interval,
				},
			})
			if err != nil {
				return err
			}
		} else {
			repScheduleParts := strings.Split(replicationSchedule, " ")
			timeZone := repScheduleParts[len(repScheduleParts)-1]
			expression := strings.TrimSuffix(strings.TrimPrefix(replicationSchedule, "USING CRON "), " "+timeZone)
			err = d.Set("replication_schedule", []any{
				map[string]any{
					"cron": []any{
						map[string]any{
							"expression": expression,
							"time_zone":  timeZone,
						},
					},
				},
			})
			if err != nil {
				return err
			}
		}
	}

	objectTypes := make([]any, len(failoverGroup.ObjectTypes))
	for i, v := range failoverGroup.ObjectTypes {
		objectTypes[i] = string(v)
	}
	objectTypesSet := schema.NewSet(schema.HashString, objectTypes)
	if err := d.Set("object_types", objectTypesSet); err != nil {
		return err
	}

	allowedIntegrationTypes := make([]any, len(failoverGroup.AllowedIntegrationTypes))
	for i, v := range failoverGroup.AllowedIntegrationTypes {
		allowedIntegrationTypes[i] = string(v)
	}
	allowedIntegrationsTypesSet := schema.NewSet(schema.HashString, allowedIntegrationTypes)
	if err := d.Set("allowed_integration_types", allowedIntegrationsTypesSet); err != nil {
		return err
	}

	allowedAccounts := make([]any, len(failoverGroup.AllowedAccounts))
	for i, v := range failoverGroup.AllowedAccounts {
		allowedAccounts[i] = v.Name()
	}
	allowedAccountsSet := schema.NewSet(schema.HashString, allowedAccounts)
	if err := d.Set("allowed_accounts", allowedAccountsSet); err != nil {
		return err
	}

	databases, err := client.FailoverGroups.ShowDatabases(ctx, id)
	if err != nil {
		return err
	}
	allowedDatabases := make([]any, len(databases))
	for i, database := range databases {
		allowedDatabases[i] = database.Name()
	}
	allowedDatabasesSet := schema.NewSet(schema.HashString, allowedDatabases)
	if len(allowedDatabases) > 0 {
		if err := d.Set("allowed_databases", allowedDatabasesSet); err != nil {
			return err
		}
	} else {
		if err := d.Set("allowed_databases", nil); err != nil {
			return err
		}
	}

	shares, err := client.FailoverGroups.ShowShares(ctx, id)
	if err != nil {
		return err
	}
	allowedShares := make([]any, len(shares))
	for i, share := range shares {
		allowedShares[i] = share.Name()
	}
	allowedSharesSet := schema.NewSet(schema.HashString, allowedShares)
	if len(allowedShares) > 0 {
		if err := d.Set("allowed_shares", allowedSharesSet); err != nil {
			return err
		}
	} else {
		if err := d.Set("allowed_shares", nil); err != nil {
			return err
		}
	}

	return nil
}

// failoverGroupAlterStepsExt applies all the changes in the original order; the order matters to Snowflake (e.g. RPO assurance vs replication schedule),
// so SET and UNSET cannot be grouped into one SET and one UNSET.
func failoverGroupAlterStepsExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	set := sdk.NewFailoverGroupSetRequest()

	if d.HasChange("object_types") {
		_, n := d.GetChange("object_types")
		newObjectTypes := expandStringList(n.(*schema.Set).List())
		objectTypes := make([]sdk.PluralObjectType, len(newObjectTypes))
		for i, v := range newObjectTypes {
			objectType, err := sdk.ToPluralObjectType(v)
			if err != nil {
				return diag.FromErr(err)
			}
			objectTypes[i] = objectType
		}
		set.WithObjectTypes(objectTypes)
		if slices.Contains(objectTypes, sdk.PluralObjectTypeIntegrations) {
			ait := expandStringList(d.Get("allowed_integration_types").(*schema.Set).List())
			allowedIntegrationTypes := make([]sdk.IntegrationType, len(ait))
			for i, v := range ait {
				allowedIntegrationTypes[i] = sdk.IntegrationType(v)
			}
			set.WithAllowedIntegrationTypes(allowedIntegrationTypes)
		}
	}

	if d.HasChange("allowed_integration_types") {
		_, n := d.GetChange("allowed_integration_types")
		newAllowedIntegrationTypes := expandStringList(n.(*schema.Set).List())
		allowedIntegrationTypes := make([]sdk.IntegrationType, len(newAllowedIntegrationTypes))
		for i, v := range newAllowedIntegrationTypes {
			allowedIntegrationTypes[i] = sdk.IntegrationType(v)
		}
		set.WithAllowedIntegrationTypes(allowedIntegrationTypes)
	}
	if !reflect.DeepEqual(set, sdk.NewFailoverGroupSetRequest()) {
		req := sdk.NewAlterSourceFailoverGroupRequest(id).WithSet(*set)
		if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("rpo_assurance") && !d.Get("rpo_assurance").(bool) {
		req := sdk.NewAlterSourceFailoverGroupRequest(id).WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithRpoAssurance(true))
		if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
			return diag.FromErr(err)
		}
		// RPO auto-enables Optimized Refresh and leaves it on after UNSET RPO.
		// Read stores optimized_refresh as false while RPO is on, so HasChange
		// may not fire even though SHOW still reports it enabled.
		if !d.Get("optimized_refresh").(bool) {
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithOptimizedRefresh(true))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("optimized_refresh") && !d.Get("optimized_refresh").(bool) {
		req := sdk.NewAlterSourceFailoverGroupRequest(id).WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithOptimizedRefresh(true))
		if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("replication_schedule") {
		replicationSchedules := d.Get("replication_schedule").([]any)
		if len(replicationSchedules) > 0 {
			replicationSchedule := replicationSchedules[0].(map[string]any)
			crons := replicationSchedule["cron"].([]any)
			var updatedReplicationSchedule string
			if len(crons) > 0 {
				cron := crons[0].(map[string]any)
				cronExpression := cron["expression"].(string)
				cronExpression = "USING CRON " + cronExpression
				if v, ok := cron["time_zone"]; ok {
					timeZone := v.(string)
					if timeZone != "" {
						cronExpression = cronExpression + " " + timeZone
					}
				}
				updatedReplicationSchedule = cronExpression
			} else {
				updatedReplicationSchedule = fmt.Sprintf("%d MINUTE", replicationSchedule["interval"].(int))
			}
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithSet(*sdk.NewFailoverGroupSetRequest().WithReplicationSchedule(updatedReplicationSchedule))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		} else {
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithReplicationSchedule(true))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("optimized_refresh") && d.Get("optimized_refresh").(bool) {
		req := sdk.NewAlterSourceFailoverGroupRequest(id).WithSet(*sdk.NewFailoverGroupSetRequest().WithOptimizedRefresh(true))
		if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("rpo_assurance") && d.Get("rpo_assurance").(bool) {
		req := sdk.NewAlterSourceFailoverGroupRequest(id).WithSet(*sdk.NewFailoverGroupSetRequest().WithRpoAssurance(true))
		if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("allowed_databases") {
		o, n := d.GetChange("allowed_databases")
		oad := expandStringList(o.(*schema.Set).List())
		oldAllowedDatabases := make([]sdk.AccountObjectIdentifier, len(oad))
		for i, v := range oad {
			oldAllowedDatabases[i] = sdk.NewAccountObjectIdentifier(v)
		}
		nad := expandStringList(n.(*schema.Set).List())
		newAllowedDatabases := make([]sdk.AccountObjectIdentifier, len(nad))
		for i, v := range nad {
			newAllowedDatabases[i] = sdk.NewAccountObjectIdentifier(v)
		}

		var removedDatabases []sdk.AccountObjectIdentifier
		for _, v := range oldAllowedDatabases {
			if !slices.Contains(newAllowedDatabases, v) {
				removedDatabases = append(removedDatabases, v)
			}
		}
		if len(removedDatabases) > 0 {
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedDatabases(removedDatabases))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(fmt.Errorf("error removing allowed databases for failover group %v err = %w", id.Name(), err))
			}
		}

		var addedDatabases []sdk.AccountObjectIdentifier
		for _, v := range newAllowedDatabases {
			if !slices.Contains(oldAllowedDatabases, v) {
				addedDatabases = append(addedDatabases, v)
			}
		}
		if len(addedDatabases) > 0 {
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedDatabases(addedDatabases))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(fmt.Errorf("error removing allowed databases for failover group %v err = %w", id.Name(), err))
			}
		}
	}

	if d.HasChange("allowed_shares") {
		o, n := d.GetChange("allowed_shares")
		oad := expandStringList(o.(*schema.Set).List())
		oldAllowedShares := make([]sdk.AccountObjectIdentifier, len(oad))
		for i, v := range oad {
			oldAllowedShares[i] = sdk.NewAccountObjectIdentifier(v)
		}
		nad := expandStringList(n.(*schema.Set).List())
		newAllowedShares := make([]sdk.AccountObjectIdentifier, len(nad))
		for i, v := range nad {
			newAllowedShares[i] = sdk.NewAccountObjectIdentifier(v)
		}

		var removedShares []sdk.AccountObjectIdentifier
		for _, v := range oldAllowedShares {
			if !slices.Contains(newAllowedShares, v) {
				removedShares = append(removedShares, v)
			}
		}
		if len(removedShares) > 0 {
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedShares(removedShares))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(fmt.Errorf("error removing allowed shares for failover group %v err = %w", id.Name(), err))
			}
		}

		var addedShares []sdk.AccountObjectIdentifier
		for _, v := range newAllowedShares {
			if !slices.Contains(oldAllowedShares, v) {
				addedShares = append(addedShares, v)
			}
		}
		if len(addedShares) > 0 {
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedShares(addedShares))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(fmt.Errorf("error removing allowed shares for failover group %v err = %w", id.Name(), err))
			}
		}
	}

	if d.HasChange("allowed_accounts") {
		o, n := d.GetChange("allowed_accounts")
		oad := expandStringList(o.(*schema.Set).List())
		oldAllowedAccounts := make([]sdk.AccountIdentifier, len(oad))
		for i, v := range oad {
			parts := strings.Split(v, ".")
			oldAllowedAccounts[i] = sdk.NewAccountIdentifier(parts[0], parts[1])
		}
		nad := expandStringList(n.(*schema.Set).List())
		newAllowedAccounts := make([]sdk.AccountIdentifier, len(nad))
		for i, v := range nad {
			parts := strings.Split(v, ".")
			newAllowedAccounts[i] = sdk.NewAccountIdentifier(parts[0], parts[1])
		}

		var removedAccounts []sdk.AccountIdentifier
		for _, v := range oldAllowedAccounts {
			if !slices.Contains(newAllowedAccounts, v) {
				removedAccounts = append(removedAccounts, v)
			}
		}
		if len(removedAccounts) > 0 {
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedAccounts(removedAccounts))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(fmt.Errorf("error removing allowed accounts for failover group %v err = %w", id.Name(), err))
			}
		}

		var addedAccounts []sdk.AccountIdentifier
		for _, v := range newAllowedAccounts {
			if !slices.Contains(oldAllowedAccounts, v) {
				addedAccounts = append(addedAccounts, v)
			}
		}
		if len(addedAccounts) > 0 {
			req := sdk.NewAlterSourceFailoverGroupRequest(id).WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedAccounts(addedAccounts))
			if err := client.FailoverGroups.AlterSource(ctx, req); err != nil {
				return diag.FromErr(fmt.Errorf("error adding allowed accounts for failover group %v err = %w", id.Name(), err))
			}
		}
	}

	return nil
}

// =============================================================================
// SDK incision points
// =============================================================================

func failoverGroupCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateFailoverGroupRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.FailoverGroups.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func failoverGroupShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.FailoverGroup, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	failoverGroup, err := client.FailoverGroups.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query failover group. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Failover group id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return failoverGroup, nil
}

func failoverGroupDropSafelyInSdkExt(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
	return client.FailoverGroups.DropSafely
}
