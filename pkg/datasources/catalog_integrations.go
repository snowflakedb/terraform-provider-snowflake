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

var catalogIntegrationsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC CATALOG INTEGRATION for each catalog integration returned by SHOW CATALOG INTEGRATIONS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"catalog_integrations": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all catalog integrations details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW CATALOG INTEGRATIONS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowCatalogIntegrationSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE CATALOG INTEGRATION.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeCatalogIntegrationSchema,
					},
				},
			},
		},
	},
}

func CatalogIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.CatalogIntegrationsDatasource), TrackingReadWrapper(datasources.CatalogIntegrations, ReadCatalogIntegrations)),
		Schema:      catalogIntegrationsSchema,
		Description: "Data source used to get details of filtered catalog integrations. Filtering is aligned with the current possibilities for [SHOW CATALOG INTEGRATIONS](https://docs.snowflake.com/en/sql-reference/sql/show-integrations) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `catalog_integrations`.",
	}
}

func ReadCatalogIntegrations(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	showRequest := sdk.NewShowCatalogIntegrationRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		showRequest.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	catalogIntegrations, err := client.CatalogIntegrations.Show(ctx, showRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("catalog_integrations_read")

	flattenedCatalogIntegrations := make([]map[string]any, len(catalogIntegrations))

	for i, catalogIntegration := range catalogIntegrations {
		catalogIntegration := catalogIntegration
		var catalogIntegrationDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			descriptions, err := client.CatalogIntegrations.Describe(ctx, catalogIntegration.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			catalogIntegrationDescriptions = schemas.CatalogIntegrationPropertiesToSchema(descriptions)
		}

		flattenedCatalogIntegrations[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.CatalogIntegrationToSchema(&catalogIntegration)},
			resources.DescribeOutputAttributeName: catalogIntegrationDescriptions,
		}
	}

	err = d.Set("catalog_integrations", flattenedCatalogIntegrations)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
