package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var backupPoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC BACKUP POLICY for each backup policy returned by SHOW BACKUP POLICIES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          inSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"backup_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all backup policy details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW BACKUP POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowBackupPolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE BACKUP POLICY.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeBackupPolicyDetailsSchema,
					},
				},
			},
		},
	},
}

func BackupPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.BackupPolicies, ReadBackupPolicies),
		Schema:      backupPoliciesSchema,
		Description: "Data source used to get details of filtered backup policies. Filtering is aligned with the current possibilities for [SHOW BACKUP POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-backup-policies) query. " +
			"The results of SHOW and DESCRIBE are encapsulated in one output collection `backup_policies`.",
	}
}

func ReadBackupPolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowBackupPolicyRequest{}

	handleLike(d, &req.Like)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	backupPolicies, err := client.BackupPolicies.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("backup_policies_read")

	flattened := make([]map[string]any, len(backupPolicies))
	for i := range backupPolicies {
		policy := backupPolicies[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.BackupPolicies.Describe(ctx, policy.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.BackupPolicyDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.BackupPolicyToSchema(&policy)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("backup_policies", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
