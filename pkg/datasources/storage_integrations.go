package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var storageIntegrationsSchema = map[string]*schema.Schema{
	"storage_integrations": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The storage integrations in the database",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func StorageIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.StorageIntegrationsDatasource), TrackingReadWrapper(datasources.StorageIntegrations, ReadStorageIntegrations)),
		Schema:      storageIntegrationsSchema,
	}
}

func ReadStorageIntegrations(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	account, err := client.ContextFunctions.CurrentAccount(ctx)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("[DEBUG] unable to retrieve current account"))
	}

	region, err := client.ContextFunctions.CurrentRegion(ctx)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("[DEBUG] unable to retrieve current region"))
	}

	d.SetId(fmt.Sprintf("%s.%s", account, region))

	storageIntegrations, err := client.StorageIntegrations.Show(ctx, sdk.NewShowStorageIntegrationRequest())
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("unable to retrieve storage integrations in account (%s), err = %w", d.Id(), err))
	}

	storageIntegrationMaps := make([]map[string]any, len(storageIntegrations))

	for i, storageIntegration := range storageIntegrations {
		storageIntegrationMaps[i] = map[string]any{
			"name":    storageIntegration.Name,
			"type":    storageIntegration.StorageType,
			"enabled": storageIntegration.Enabled,
			"comment": storageIntegration.Comment,
		}
	}

	return diag.FromErr(d.Set("storage_integrations", storageIntegrationMaps))
}
