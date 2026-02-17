package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var storageIntegrationsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC STORAGE INTEGRATION for each storage integration returned by SHOW STORAGE INTEGRATIONS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"storage_integrations": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all storage integrations details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW STORAGE INTEGRATIONS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowStorageIntegrationSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the aggregated output of DESCRIBE STORAGE INTEGRATIONS.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeStorageIntegrationAllDetailsSchema,
					},
				},
			},
		},
	},
}

func StorageIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.StorageIntegrationsDatasource), TrackingReadWrapper(datasources.StorageIntegrations, ReadStorageIntegrations)),
		Schema:      storageIntegrationsSchema,
		Description: "Data source used to get details of filtered storage integrations. Filtering is aligned with the current possibilities for [SHOW STORAGE INTEGRATIONS](https://docs.snowflake.com/en/sql-reference/sql/show-integrations) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `storage_integrations`.",
	}
}

func ReadStorageIntegrations(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	showRequest := sdk.NewShowStorageIntegrationRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		showRequest.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	storageIntegrations, err := client.StorageIntegrations.Show(ctx, showRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("storage_integrations_read")

	flattenedStorageIntegrations := make([]map[string]any, len(storageIntegrations))

	for i, storageIntegration := range storageIntegrations {
		storageIntegration := storageIntegration
		var storageIntegrationDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			descriptions, err := client.StorageIntegrations.DescribeDetails(ctx, storageIntegration.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			storageIntegrationDescriptions = make([]map[string]any, 1)
			storageIntegrationDescriptions[0] = schemas.StorageIntegrationAllDetailsToSchema(descriptions)
		}

		flattenedStorageIntegrations[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.StorageIntegrationToSchema(&storageIntegration)},
			resources.DescribeOutputAttributeName: storageIntegrationDescriptions,
		}
	}

	err = d.Set("storage_integrations", flattenedStorageIntegrations)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
