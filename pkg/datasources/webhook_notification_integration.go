package datasources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var webhookNotificationIntegrationDataSourceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the webhook notification integration to read.",
	},
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"webhook_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"webhook_secret": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"webhook_body_template": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"webhook_headers": {
		Type:     schema.TypeMap,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func WebhookNotificationIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.WebhookNotificationIntegrationDatasource), TrackingReadWrapper(datasources.WebhookNotificationIntegration, ReadWebhookNotificationIntegration)),
		Schema:      webhookNotificationIntegrationDataSourceSchema,
	}
}

func ReadWebhookNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	integration, err := client.NotificationIntegrations.ShowByID(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading webhook notification integration %q: %w", name, err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	if err := d.Set("name", integration.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", integration.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", integration.Comment); err != nil {
		return diag.FromErr(err)
	}

	properties, err := client.NotificationIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error describing webhook notification integration %q: %w", name, err))
	}
	for _, property := range properties {
		switch property.Name {
		case "ENABLED":
			// set via SHOW
		case "WEBHOOK_URL":
			if err := d.Set("webhook_url", property.Value); err != nil {
				return diag.FromErr(err)
			}
		case "WEBHOOK_SECRET":
			if err := d.Set("webhook_secret", property.Value); err != nil {
				return diag.FromErr(err)
			}
		case "WEBHOOK_BODY_TEMPLATE":
			if err := d.Set("webhook_body_template", property.Value); err != nil {
				return diag.FromErr(err)
			}
		case "WEBHOOK_HEADERS":
			headers, err := webhookDatasourceParseHeadersMap(property.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("webhook_headers", headers); err != nil {
				return diag.FromErr(err)
			}
		default:
			log.Printf("[WARN] unexpected webhook notification integration property %v returned from Snowflake", property.Name)
		}
	}

	return nil
}

// webhookDatasourceParseHeadersMap parses Snowflake DESCRIBE output for WEBHOOK_HEADERS,
// e.g. "{Content-Type=application/json, X-Custom=value}" → map.
func webhookDatasourceParseHeadersMap(s string) (map[string]string, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	s = strings.TrimSpace(s)
	if s == "" {
		return map[string]string{}, nil
	}
	result := make(map[string]string)
	for _, pair := range strings.Split(s, ", ") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("unexpected WEBHOOK_HEADERS format: %q", s)
		}
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return result, nil
}
