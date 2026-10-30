package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DESC REPLICATION GROUP is not a published command, so there is no with_describe.
// SHOW supports IN ACCOUNT only. LIKE, STARTS WITH, LIMIT, and database or schema IN are not on the page.
var replicationGroupsSchema = map[string]*schema.Schema{
	"in_account": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Limits results to replication groups in this account, as organization and account name. Matches [SHOW REPLICATION GROUPS](https://docs.snowflake.com/en/sql-reference/sql/show-replication-groups) `IN ACCOUNT`.",
		ValidateDiagFunc: resources.IsValidAccountIdentifier(),
	},
	"replication_groups": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the output of SHOW REPLICATION GROUPS.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW REPLICATION GROUPS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowReplicationGroupSchema,
					},
				},
			},
		},
	},
}

func ReplicationGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.ReplicationGroups, ReadReplicationGroups),
		Schema:      replicationGroupsSchema,
		Description: "Data source used to get details of replication groups. Filtering is aligned with the current possibilities for [SHOW REPLICATION GROUPS](https://docs.snowflake.com/en/sql-reference/sql/show-replication-groups) query. The results of SHOW are encapsulated in one output collection `replication_groups`.",
	}
}

func ReadReplicationGroups(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowReplicationGroupRequest()
	if v, ok := d.GetOk("in_account"); ok {
		account, err := helpers.DecodeSnowflakeAccountIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithInAccount(account)
	}
	groups, err := client.ReplicationGroups.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("replication_groups_read")

	flattened := make([]map[string]any, len(groups))
	for i, group := range groups {
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.ReplicationGroupToSchema(&group)},
		}
	}
	if err := d.Set("replication_groups", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
