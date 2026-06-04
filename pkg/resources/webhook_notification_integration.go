package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var webhookNotificationIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
	},
	"webhook_url": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The URL of the webhook endpoint.",
	},
	"webhook_secret": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: `Fully qualified name of the GENERIC_STRING secret ("database"."schema"."name") containing the webhook credentials.`,
	},
	"webhook_body_template": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The body template for the webhook call. Must contain the placeholder SNOWFLAKE_WEBHOOK_MESSAGE.",
	},
	"webhook_headers": {
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "HTTP headers to include in the webhook call, expressed as a map of header name to header value.",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func WebhookNotificationIntegration() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.AccountObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.NotificationIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.WebhookNotificationIntegrationResource), TrackingCreateWrapper(resources.WebhookNotificationIntegration, CreateWebhookNotificationIntegration)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.WebhookNotificationIntegrationResource), TrackingReadWrapper(resources.WebhookNotificationIntegration, ReadWebhookNotificationIntegration)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.WebhookNotificationIntegrationResource), TrackingUpdateWrapper(resources.WebhookNotificationIntegration, UpdateWebhookNotificationIntegration)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.WebhookNotificationIntegrationResource), TrackingDeleteWrapper(resources.WebhookNotificationIntegration, deleteFunc)),

		Schema: webhookNotificationIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateWebhookNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	enabled := d.Get("enabled").(bool)

	createRequest := sdk.NewCreateNotificationIntegrationRequest(id, enabled)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(v.(string))
	}

	webhookUrl := d.Get("webhook_url").(string)
	webhookParamsRequest := sdk.NewWebhookParamsRequest(webhookUrl)

	if v, ok := d.GetOk("webhook_secret"); ok {
		secretId, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid webhook_secret identifier %q: %w", v, err))
		}
		webhookParamsRequest.WithWebhookSecret(secretId)
	}

	if v, ok := d.GetOk("webhook_body_template"); ok {
		webhookParamsRequest.WithWebhookBodyTemplate(v.(string))
	}

	if v, ok := d.GetOk("webhook_headers"); ok {
		webhookParamsRequest.WithWebhookHeaders(toWebhookHeaders(v.(map[string]interface{})))
	}

	createRequest.WithWebhookParams(*webhookParamsRequest)

	err := client.NotificationIntegrations.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating webhook notification integration: %w", err))
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadWebhookNotificationIntegration(ctx, d, meta)
}

func ReadWebhookNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.NotificationIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query webhook notification integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Webhook notification integration id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", integration.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", integration.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", integration.Comment); err != nil {
		return diag.FromErr(err)
	}

	integrationProperties, err := client.NotificationIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe webhook notification integration: %w", err))
	}

	for _, property := range integrationProperties {
		name := property.Name
		value := property.Value

		switch name {
		case "ENABLED":
			// set via SHOW
		case "WEBHOOK_URL":
			if err := d.Set("webhook_url", value); err != nil {
				return diag.FromErr(err)
			}
		case "WEBHOOK_SECRET":
			if err := d.Set("webhook_secret", value); err != nil {
				return diag.FromErr(err)
			}
		case "WEBHOOK_BODY_TEMPLATE":
			if err := d.Set("webhook_body_template", value); err != nil {
				return diag.FromErr(err)
			}
		case "WEBHOOK_HEADERS":
			headers, err := parseWebhookHeadersMap(value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("webhook_headers", headers); err != nil {
				return diag.FromErr(err)
			}
		default:
			log.Printf("[WARN] unexpected webhook notification integration property %v returned from Snowflake", name)
		}
	}

	return diag.FromErr(err)
}

func UpdateWebhookNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	var runSetStatement bool
	var runUnsetStatement bool

	setRequest := sdk.NewNotificationIntegrationSetRequest()
	setWebhookRequest := sdk.NewSetWebhookParamsRequest()
	var setWebhookChanged bool
	unsetRequest := sdk.NewNotificationIntegrationUnsetWebhookParamsRequest()

	if d.HasChange("enabled") {
		runSetStatement = true
		setRequest.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("comment") {
		v := d.Get("comment").(string)
		if v == "" {
			runUnsetStatement = true
			unsetRequest.WithComment(true)
		} else {
			runSetStatement = true
			setRequest.WithComment(v)
		}
	}

	if d.HasChange("webhook_url") {
		runSetStatement = true
		setWebhookChanged = true
		setWebhookRequest.WithWebhookUrl(d.Get("webhook_url").(string))
	}

	if d.HasChange("webhook_secret") {
		v := d.Get("webhook_secret").(string)
		if v == "" {
			runUnsetStatement = true
			unsetRequest.WithWebhookSecret(true)
		} else {
			secretId, err := sdk.ParseSchemaObjectIdentifier(v)
			if err != nil {
				return diag.FromErr(fmt.Errorf("invalid webhook_secret identifier %q: %w", v, err))
			}
			runSetStatement = true
			setWebhookChanged = true
			setWebhookRequest.WithWebhookSecret(secretId)
		}
	}

	if d.HasChange("webhook_body_template") {
		v := d.Get("webhook_body_template").(string)
		if v == "" {
			runUnsetStatement = true
			unsetRequest.WithWebhookBodyTemplate(true)
		} else {
			runSetStatement = true
			setWebhookChanged = true
			setWebhookRequest.WithWebhookBodyTemplate(v)
		}
	}

	if d.HasChange("webhook_headers") {
		v := d.Get("webhook_headers").(map[string]interface{})
		if len(v) == 0 {
			runUnsetStatement = true
			unsetRequest.WithWebhookHeaders(true)
		} else {
			runSetStatement = true
			setWebhookChanged = true
			setWebhookRequest.WithWebhookHeaders(toWebhookHeaders(v))
		}
	}

	if setWebhookChanged {
		setRequest.WithSetWebhookParams(*setWebhookRequest)
	}

	if runSetStatement {
		err := client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithSet(*setRequest))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating webhook notification integration: %w", err))
		}
	}

	if runUnsetStatement {
		err := client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithUnsetWebhookParams(*unsetRequest))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating webhook notification integration: %w", err))
		}
	}

	return ReadWebhookNotificationIntegration(ctx, d, meta)
}

func toWebhookHeaders(raw map[string]interface{}) []sdk.WebhookHeaderRequest {
	headers := make([]sdk.WebhookHeaderRequest, 0, len(raw))
	for header, val := range raw {
		headers = append(headers, sdk.WebhookHeaderRequest{Header: header, Value: val.(string)})
	}
	return headers
}

// parseWebhookHeadersMap parses the Snowflake DESCRIBE output for WEBHOOK_HEADERS.
// Snowflake returns headers as a map string, e.g. "{Content-Type=application/json, X-Custom=value}".
func parseWebhookHeadersMap(s string) (map[string]string, error) {
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
