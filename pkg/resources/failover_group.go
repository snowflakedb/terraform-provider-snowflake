package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var failoverGroupSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the failover group. The identifier must start with an alphabetic character and cannot contain spaces or special characters unless the identifier string is enclosed in double quotes (e.g. \"My object\"). Identifiers enclosed in double quotes are also case-sensitive.",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(old, new)
		},
	},
	"object_types": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Type(s) of objects for which you are enabling replication and failover from the source account to the target account. The following object types are supported: \"ACCOUNT PARAMETERS\", \"DATABASES\", \"INTEGRATIONS\", \"NETWORK POLICIES\", \"RESOURCE MONITORS\", \"ROLES\", \"SHARES\", \"USERS\", \"WAREHOUSES\"",
	},
	"allowed_databases": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the database or list of databases for which you are enabling replication and failover from the source account to the target account. The OBJECT_TYPES list must include DATABASES to set this parameter.",
	},
	"allowed_shares": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the share or list of shares for which you are enabling replication and failover from the source account to the target account. The OBJECT_TYPES list must include SHARES to set this parameter.",
	},
	"allowed_integration_types": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Type(s) of integrations for which you are enabling replication and failover from the source account to the target account. This property requires that the OBJECT_TYPES list include INTEGRATIONS to set this parameter. The following integration types are supported: \"SECURITY INTEGRATIONS\", \"API INTEGRATIONS\", \"STORAGE INTEGRATIONS\", \"EXTERNAL ACCESS INTEGRATIONS\", \"NOTIFICATION INTEGRATIONS\"",
	},
	"allowed_accounts": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the target account or list of target accounts to which replication and failover of specified objects from the source account is enabled. Secondary failover groups in the target accounts in this list can be promoted to serve as the primary failover group in case of failover. Expected in the form `<org_name>.<target_account_name>`. This value is case-sensitive.",
	},
	"ignore_edition_check": {
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"from_replica"},
		Description:   "Allows replicating objects to accounts on lower editions.",
	},
	"from_replica": {
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		MaxItems:      1,
		ConflictsWith: []string{"object_types", "allowed_accounts", "allowed_databases", "allowed_shares", "allowed_integration_types", "ignore_edition_check", "replication_schedule"},
		Description:   "Specifies the name of the replica to use as the source for the failover group.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"organization_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of your Snowflake organization.",
				},
				"source_account_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Source account from which you are enabling replication and failover of the specified objects.",
				},
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Identifier for the primary failover group in the source account.",
				},
			},
		},
	},
	"replication_schedule": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		Description:   "Specifies the schedule for refreshing secondary failover groups.",
		ConflictsWith: []string{"from_replica"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cron": {
					Type:          schema.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"replication_schedule.interval"},
					Description:   "Specifies the cron expression for the replication schedule. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"expression": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the cron expression for the replication schedule. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
							},
							"time_zone": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the time zone for secondary group refresh.",
							},
						},
					},
				},
				"interval": {
					Type:          schema.TypeInt,
					Optional:      true,
					ConflictsWith: []string{"replication_schedule.cron"},
					Description:   "Specifies the interval in minutes for the replication schedule. The interval must be greater than 0 and less than 1440 (24 hours).",
				},
			},
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func FailoverGroup() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.AccountObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.FailoverGroups.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingCreateWrapper(resources.FailoverGroup, CreateFailoverGroup)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingReadWrapper(resources.FailoverGroup, ReadFailoverGroup)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingUpdateWrapper(resources.FailoverGroup, UpdateFailoverGroup)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FailoverGroupResource), TrackingDeleteWrapper(resources.FailoverGroup, deleteFunc)),

		Schema: failoverGroupSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateFailoverGroup implements schema.CreateFunc.
func CreateFailoverGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	// getting required attributes
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	// if from_replica is set, then we are creating a failover group from an existing replica
	if v, ok := d.GetOk("from_replica"); ok {
		fromReplica := v.([]interface{})[0].(map[string]interface{})
		organizationName := fromReplica["organization_name"].(string)
		sourceAccountName := fromReplica["source_account_name"].(string)
		sourceFailoverGroupName := fromReplica["name"].(string)

		primaryFailoverGroupID := sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(organizationName, sourceAccountName), sdk.NewAccountObjectIdentifier(sourceFailoverGroupName))
		err := client.FailoverGroups.CreateSecondaryReplicationGroup(ctx, id, primaryFailoverGroupID, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(name)
		return ReadFailoverGroup(ctx, d, meta)
	}

	// these two are required attributes if from_replica is not set
	if _, ok := d.GetOk("object_types"); !ok {
		return diag.FromErr(errors.New("object_types field is required when from_replica is not set"))
	}
	objectTypesList := expandStringList(d.Get("object_types").(*schema.Set).List())
	objectTypes := make([]sdk.PluralObjectType, len(objectTypesList))
	for i, v := range objectTypesList {
		objectTypes[i] = sdk.PluralObjectType(v)
	}

	if _, ok := d.GetOk("allowed_accounts"); !ok {
		return diag.FromErr(errors.New("allowed_accounts field is required when from_replica is not set"))
	}
	aaList := expandStringList(d.Get("allowed_accounts").(*schema.Set).List())
	allowedAccounts := make([]sdk.AccountIdentifier, len(aaList))
	for i, v := range aaList {
		// validation since we cannot do that in the ValidateFunc
		parts := strings.Split(v, ".")
		if len(parts) != 2 {
			return diag.FromErr(fmt.Errorf("allowed_account %s cannot be an account locator and must be of the format <org_name>.<target_account_name>", allowedAccounts[i]))
		}
		organizationName := parts[0]
		accountName := parts[1]
		allowedAccounts[i] = sdk.NewAccountIdentifier(organizationName, accountName)
	}

	var opts sdk.CreateFailoverGroupOptions
	// setting optional attributes
	if v, ok := d.GetOk("allowed_databases"); ok {
		allowedDatabasesList := expandStringList(v.(*schema.Set).List())
		allowedDatabaseIdentifiers := make([]sdk.AccountObjectIdentifier, len(allowedDatabasesList))
		for i, v := range allowedDatabasesList {
			allowedDatabaseIdentifiers[i] = sdk.NewAccountObjectIdentifier(v)
		}
		opts.AllowedDatabases = allowedDatabaseIdentifiers
	}

	if v, ok := d.GetOk("allowed_shares"); ok {
		allowedSharesList := expandStringList(v.(*schema.Set).List())
		allowedShareIdentifiers := make([]sdk.AccountObjectIdentifier, len(allowedSharesList))
		for i, v := range allowedSharesList {
			allowedShareIdentifiers[i] = sdk.NewAccountObjectIdentifier(v)
		}
		opts.AllowedShares = allowedShareIdentifiers
	}

	if v, ok := d.GetOk("allowed_integration_types"); ok {
		allowedIntegrationTypesList := expandStringList(v.(*schema.Set).List())
		allowedIntegrationTypes := make([]sdk.IntegrationType, len(allowedIntegrationTypesList))
		for i, v := range allowedIntegrationTypesList {
			allowedIntegrationTypes[i] = sdk.IntegrationType(v)
		}
		opts.AllowedIntegrationTypes = allowedIntegrationTypes
	}

	if v, ok := d.GetOk("ignore_edition_check"); ok {
		opts.IgnoreEditionCheck = sdk.Bool(v.(bool))
	}

	if v, ok := d.GetOk("replication_schedule"); ok {
		replicationSchedule := v.([]interface{})[0].(map[string]interface{})
		if v, ok := replicationSchedule["cron"]; ok {
			c := v.([]interface{})
			if len(c) > 0 {
				cron := c[0].(map[string]interface{})
				cronExpression := cron["expression"].(string)
				if v, ok := cron["time_zone"]; ok {
					timeZone := v.(string)
					cronExpression += " " + timeZone
				}
				opts.ReplicationSchedule = sdk.String("USING CRON " + cronExpression)
			}
		}
		if v, ok := replicationSchedule["interval"]; ok {
			interval := v.(int)
			if interval > 1 {
				opts.ReplicationSchedule = sdk.String(fmt.Sprintf("%d MINUTE", interval))
			}
		}
	}

	err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	return ReadFailoverGroup(ctx, d, meta)
}

// ReadFailoverGroup implements schema.ReadFunc.
func ReadFailoverGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)
	failoverGroup, err := client.FailoverGroups.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query failover group. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Failover group id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", failoverGroup.Name); err != nil {
		return diag.FromErr(err)
	}
	// if the failover group is created from a replica, then we do not want to get the other values
	if _, ok := d.GetOk("from_replica"); ok {
		return nil
	}

	replicationSchedule := failoverGroup.ReplicationSchedule
	if replicationSchedule != "" {
		if strings.Contains(replicationSchedule, "MINUTE") {
			interval, err := strconv.Atoi(strings.TrimSuffix(replicationSchedule, " MINUTE"))
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set("replication_schedule", []interface{}{
				map[string]interface{}{
					"interval": interval,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			repScheduleParts := strings.Split(replicationSchedule, " ")
			timeZone := repScheduleParts[len(repScheduleParts)-1]
			expression := strings.TrimSuffix(strings.TrimPrefix(replicationSchedule, "USING CRON "), " "+timeZone)
			err = d.Set("replication_schedule", []interface{}{
				map[string]interface{}{
					"cron": []interface{}{
						map[string]interface{}{
							"expression": expression,
							"time_zone":  timeZone,
						},
					},
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	// object types
	objectTypes := make([]interface{}, len(failoverGroup.ObjectTypes))
	for i, v := range failoverGroup.ObjectTypes {
		objectTypes[i] = string(v)
	}
	objectTypesSet := schema.NewSet(schema.HashString, objectTypes)
	if err := d.Set("object_types", objectTypesSet); err != nil {
		return diag.FromErr(err)
	}

	// integration types
	allowedIntegrationTypes := make([]interface{}, len(failoverGroup.AllowedIntegrationTypes))
	for i, v := range failoverGroup.AllowedIntegrationTypes {
		allowedIntegrationTypes[i] = string(v)
	}

	allowedIntegrationsTypesSet := schema.NewSet(schema.HashString, allowedIntegrationTypes)
	if err := d.Set("allowed_integration_types", allowedIntegrationsTypesSet); err != nil {
		return diag.FromErr(err)
	}

	// allowed accounts
	allowedAccounts := make([]interface{}, len(failoverGroup.AllowedAccounts))
	for i, v := range failoverGroup.AllowedAccounts {
		allowedAccounts[i] = v.Name()
	}
	allowedAccountsSet := schema.NewSet(schema.HashString, allowedAccounts)
	if err := d.Set("allowed_accounts", allowedAccountsSet); err != nil {
		return diag.FromErr(err)
	}

	// allowed databases
	databases, err := client.FailoverGroups.ShowDatabases(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	allowedDatabases := make([]interface{}, len(databases))
	for i, database := range databases {
		allowedDatabases[i] = database.Name()
	}
	allowedDatabasesSet := schema.NewSet(schema.HashString, allowedDatabases)
	if len(allowedDatabases) > 0 {
		if err := d.Set("allowed_databases", allowedDatabasesSet); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("allowed_databases", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	// allowed shares
	shares, err := client.FailoverGroups.ShowShares(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	allowedShares := make([]interface{}, len(shares))
	for i, share := range shares {
		allowedShares[i] = share.Name()
	}
	allowedSharesSet := schema.NewSet(schema.HashString, allowedShares)
	if len(allowedShares) > 0 {
		if err := d.Set("allowed_shares", allowedSharesSet); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("allowed_shares", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// UpdateFailoverGroup implements schema.UpdateFunc.
func UpdateFailoverGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	// alter failover group <name> set ...
	opts := &sdk.AlterSourceFailoverGroupOptions{
		Set: &sdk.FailoverGroupSet{},
	}
	runSet := false

	if d.HasChange("object_types") {
		_, n := d.GetChange("object_types")
		newObjectTypes := expandStringList(n.(*schema.Set).List())
		objectTypes := make([]sdk.PluralObjectType, len(newObjectTypes))
		for i, v := range newObjectTypes {
			objectTypes[i] = sdk.PluralObjectType(v)
		}
		opts.Set.ObjectTypes = objectTypes
		if slices.Contains(objectTypes, sdk.PluralObjectTypeIntegrations) {
			ait := expandStringList(d.Get("allowed_integration_types").(*schema.Set).List())
			allowedIntegrationTypes := make([]sdk.IntegrationType, len(ait))
			for i, v := range ait {
				allowedIntegrationTypes[i] = sdk.IntegrationType(v)
			}
			opts.Set.AllowedIntegrationTypes = allowedIntegrationTypes
		}
		runSet = true
	}

	if d.HasChange("allowed_integration_types") {
		_, n := d.GetChange("allowed_integration_types")
		newAllowedIntegrationTypes := expandStringList(n.(*schema.Set).List())
		allowedIntegrationTypes := make([]sdk.IntegrationType, len(newAllowedIntegrationTypes))
		for i, v := range newAllowedIntegrationTypes {
			allowedIntegrationTypes[i] = sdk.IntegrationType(v)
		}
		opts.Set.AllowedIntegrationTypes = allowedIntegrationTypes
		runSet = true
	}
	if runSet {
		if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
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
			err := client.FailoverGroups.AlterSource(ctx, id, &sdk.AlterSourceFailoverGroupOptions{
				Set: &sdk.FailoverGroupSet{
					ReplicationSchedule: sdk.String(updatedReplicationSchedule),
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := client.FailoverGroups.AlterSource(ctx, id, &sdk.AlterSourceFailoverGroupOptions{
				Unset: &sdk.FailoverGroupUnset{
					ReplicationSchedule: sdk.Bool(true),
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
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
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Remove: &sdk.FailoverGroupRemove{
					AllowedDatabases: removedDatabases,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
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
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Add: &sdk.FailoverGroupAdd{
					AllowedDatabases: addedDatabases,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
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
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Remove: &sdk.FailoverGroupRemove{
					AllowedShares: removedShares,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
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
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Add: &sdk.FailoverGroupAdd{
					AllowedShares: addedShares,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
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
			organizationName := parts[0]
			accountName := parts[1]
			accountIdentifier := sdk.NewAccountIdentifier(accountName, organizationName)
			oldAllowedAccounts[i] = accountIdentifier
		}
		nad := expandStringList(n.(*schema.Set).List())
		newAllowedAccounts := make([]sdk.AccountIdentifier, len(nad))
		for i, v := range nad {
			parts := strings.Split(v, ".")
			organizationName := parts[0]
			accountName := parts[1]
			accountIdentifier := sdk.NewAccountIdentifier(accountName, organizationName)
			newAllowedAccounts[i] = accountIdentifier
		}

		var removedAccounts []sdk.AccountIdentifier
		for _, v := range oldAllowedAccounts {
			if !slices.Contains(newAllowedAccounts, v) {
				removedAccounts = append(removedAccounts, v)
			}
		}
		if len(removedAccounts) > 0 {
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Remove: &sdk.FailoverGroupRemove{
					AllowedAccounts: removedAccounts,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
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
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Add: &sdk.FailoverGroupAdd{
					AllowedAccounts: addedAccounts,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
				return diag.FromErr(fmt.Errorf("error adding allowed accounts for failover group %v err = %w", id.Name(), err))
			}
		}
	}

	return ReadFailoverGroup(ctx, d, meta)
}
