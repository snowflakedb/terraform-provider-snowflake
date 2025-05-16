package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var emailNotificationIntegrationSchema = map[string]*schema.Schema{
	// The first part of the schema is shared between all integration vendors
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
	},
	"allowed_recipients": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "List of email addresses that should receive notifications.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A comment for the email integration.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func EmailNotificationIntegration() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErr[sdk.AccountObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.NotificationIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.EmailNotificationIntegrationResource), TrackingCreateWrapper(resources.EmailNotificationIntegration, CreateEmailNotificationIntegration)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.EmailNotificationIntegrationResource), TrackingReadWrapper(resources.EmailNotificationIntegration, ReadEmailNotificationIntegration)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.EmailNotificationIntegrationResource), TrackingUpdateWrapper(resources.EmailNotificationIntegration, UpdateEmailNotificationIntegration)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.EmailNotificationIntegrationResource), TrackingDeleteWrapper(resources.EmailNotificationIntegration, deleteFunc)),

		Schema: emailNotificationIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func toAllowedRecipients(emails []string) []sdk.NotificationIntegrationAllowedRecipient {
	allowedRecipients := make([]sdk.NotificationIntegrationAllowedRecipient, len(emails))
	for i, prefix := range emails {
		allowedRecipients[i] = sdk.NotificationIntegrationAllowedRecipient{Email: prefix}
	}
	return allowedRecipients
}

// CreateEmailNotificationIntegration implements schema.CreateFunc.
func CreateEmailNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	enabled := d.Get("enabled").(bool)

	createRequest := sdk.NewCreateNotificationIntegrationRequest(id, enabled)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	emailParamsRequest := sdk.NewEmailParamsRequest()
	if v, ok := d.GetOk("allowed_recipients"); ok {
		emailParamsRequest.WithAllowedRecipients(toAllowedRecipients(expandStringList(v.(*schema.Set).List())))
	}
	createRequest.WithEmailParams(emailParamsRequest)

	err := client.NotificationIntegrations.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating notification integration: %w", err))
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadEmailNotificationIntegration(ctx, d, meta)
}

// ReadEmailNotificationIntegration implements schema.ReadFunc.
func ReadEmailNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.NotificationIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to email notification integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Email notification integration id: %s, Err: %s", id.FullyQualifiedName(), err),
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

	// Some properties come from the DESCRIBE INTEGRATION call
	integrationProperties, err := client.NotificationIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe notification integration: %w", err))
	}
	for _, property := range integrationProperties {
		name := property.Name
		value := property.Value

		switch name {
		case "ALLOWED_RECIPIENTS":
			if value == "" {
				if err := d.Set("allowed_recipients", make([]string, 0)); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err := d.Set("allowed_recipients", strings.Split(value, ",")); err != nil {
					return diag.FromErr(err)
				}
			}
		default:
			log.Printf("[WARN] unexpected notification integration property %v returned from Snowflake", name)
		}
	}

	return diag.FromErr(err)
}

// UpdateEmailNotificationIntegration implements schema.UpdateFunc.
func UpdateEmailNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	var runSetStatement bool
	var runUnsetStatement bool
	setRequest := sdk.NewNotificationIntegrationSetRequest()
	unsetRequest := sdk.NewNotificationIntegrationUnsetEmailParamsRequest()
	if d.HasChange("comment") {
		v := d.Get("comment").(string)
		if v == "" {
			runUnsetStatement = true
			unsetRequest.WithComment(sdk.Bool(true))
		} else {
			runSetStatement = true
			setRequest.WithComment(sdk.String(d.Get("comment").(string)))
		}
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		setRequest.WithEnabled(sdk.Bool(d.Get("enabled").(bool)))
	}

	if d.HasChange("allowed_recipients") {
		v := d.Get("allowed_recipients").(*schema.Set).List()
		if len(v) == 0 {
			runUnsetStatement = true
			unsetRequest.WithAllowedRecipients(sdk.Bool(true))
		} else {
			runSetStatement = true
			setRequest.WithSetEmailParams(sdk.NewSetEmailParamsRequest(toAllowedRecipients(expandStringList(v))))
		}
	}

	if runSetStatement {
		err := client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithSet(setRequest))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating notification integration: %w", err))
		}
	}

	if runUnsetStatement {
		err := client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithUnsetEmailParams(unsetRequest))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating notification integration: %w", err))
		}
	}

	return ReadEmailNotificationIntegration(ctx, d, meta)
}
