package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var failoverGroupsSchema = map[string]*schema.Schema{
	"in_account": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the identifier for the account",
	},
	"failover_groups": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of all the failover groups available in the system.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"region_group": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Region group where the account is located. Note: this column is only visible to organizations that span multiple Region Groups.",
				},
				"snowflake_region": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Snowflake Region where the account is located. A Snowflake Region is a distinct location within a cloud platform region that is isolated from other Snowflake Regions. A Snowflake Region can be either multi-tenant or single-tenant (for a Virtual Private Snowflake account).",
				},
				"created_on": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Date and time failover group was created.",
				},
				"account_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the account.",
				},
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Type of group. Valid value is FAILOVER.",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Comment string.",
				},
				"is_primary": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether the failover group is the primary group.",
				},
				"primary": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the primary group.",
				},
				"object_types": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "List of specified object types enabled for replication and failover.",
					Elem:        schema.TypeString,
				},
				"allowed_integration_types": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "A list of integration types that are enabled for replication.",
					Elem:        schema.TypeString,
				},
				"allowed_accounts": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "List of accounts enabled for replication and failover.",
					Elem:        schema.TypeString,
				},
				"organization_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of your Snowflake organization.",
				},
				"account_locator": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Account locator in a region.",
				},
				"replication_schedule": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Scheduled interval for refresh; NULL if no replication schedule is set.",
				},
				"secondary_state": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Current state of scheduled refresh. Valid values are started or suspended. NULL if no replication schedule is set.",
				},
				"next_scheduled_refresh": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Date and time of the next scheduled refresh.",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the role with the OWNERSHIP privilege on the failover group. NULL if the failover group is in a different region.",
				},
			},
		},
	},
}

// FailoverGroups Snowflake FailoverGroups resource.
func FailoverGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.FailoverGroupsDatasource), TrackingReadWrapper(datasources.FailoverGroups, ReadFailoverGroups)),
		Schema:      failoverGroupsSchema,
	}
}

// ReadFailoverGroups lists failover groups.
func ReadFailoverGroups(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	inAccount := d.Get("in_account").(string)
	opts := sdk.ShowFailoverGroupOptions{}
	if inAccount != "" {
		opts.InAccount = sdk.NewAccountIdentifierFromAccountLocator(inAccount)
	}
	failoverGroups, err := client.FailoverGroups.Show(ctx, &opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("failover_groups")
	failoverGroupsFlatten := []map[string]interface{}{}
	for _, failoverGroup := range failoverGroups {
		m := map[string]interface{}{}
		m["region_group"] = failoverGroup.RegionGroup
		m["snowflake_region"] = failoverGroup.SnowflakeRegion
		m["created_on"] = failoverGroup.CreatedOn.String()
		m["account_name"] = failoverGroup.AccountName
		m["type"] = failoverGroup.Type
		m["comment"] = failoverGroup.Comment
		m["is_primary"] = failoverGroup.IsPrimary
		m["primary"] = failoverGroup.Primary.FullyQualifiedName()

		ot := make([]string, len(failoverGroup.ObjectTypes))
		for i, o := range failoverGroup.ObjectTypes {
			ot[i] = string(o)
		}
		m["object_types"] = ot
		ait := make([]string, len(failoverGroup.AllowedIntegrationTypes))
		for i, a := range failoverGroup.AllowedIntegrationTypes {
			ait[i] = string(a)
		}
		m["allowed_integration_types"] = ait
		aa := make([]string, len(failoverGroup.AllowedAccounts))
		for i, a := range failoverGroup.AllowedAccounts {
			aa[i] = a.Name()
		}
		m["allowed_accounts"] = aa
		m["organization_name"] = failoverGroup.OrganizationName
		m["account_locator"] = failoverGroup.AccountLocator
		m["replication_schedule"] = failoverGroup.ReplicationSchedule
		m["secondary_state"] = string(failoverGroup.SecondaryState)
		m["next_scheduled_refresh"] = failoverGroup.NextScheduledRefresh
		m["owner"] = failoverGroup.Owner
		failoverGroupsFlatten = append(failoverGroupsFlatten, m)
	}
	if err := d.Set("failover_groups", failoverGroupsFlatten); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
