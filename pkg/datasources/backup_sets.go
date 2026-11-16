package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var backupSetsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC BACKUP SET for each backup set returned by SHOW BACKUP SETS. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          inSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"backup_sets": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all backup set details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW BACKUP SETS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowBackupSetSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE BACKUP SET.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeBackupSetDetailsSchema,
					},
				},
			},
		},
	},
}

func BackupSets() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.BackupSetsDatasource), TrackingReadWrapper(datasources.BackupSets, ReadBackupSets)),
		Schema:      backupSetsSchema,
		Description: docs.WithPreviewFeatureMark() + "Data source used to get details of filtered backup sets. Filtering is aligned with the current possibilities for [SHOW BACKUP SETS](https://docs.snowflake.com/en/sql-reference/sql/show-backup-sets) query. " +
			"The results of SHOW and DESCRIBE are encapsulated in one output collection `backup_sets`.",
	}
}

func ReadBackupSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowBackupSetRequest{}

	handleLike(d, &req.Like)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	backupSets, err := client.BackupSets.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("backup_sets_read")

	flattened := make([]map[string]any, len(backupSets))
	for i := range backupSets {
		backupSet := backupSets[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.BackupSets.Describe(ctx, backupSet.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.BackupSetDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.BackupSetToSchema(&backupSet)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("backup_sets", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
