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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO [SNOW-1348345]: split this resource into smaller ones (SNOW-1021713)
// TODO [SNOW-1348345]: remove SQS entirely
// TODO [SNOW-1348345]: support Azure push notifications (AZURE_EVENT_GRID)
var notificationIntegrationSchema = map[string]*schema.Schema{
	// The first part of the schema is shared between all integration vendors
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "QUEUE",
		ValidateFunc: validation.StringInSlice([]string{"QUEUE"}, true),
		Description:  "A type of integration",
		Deprecated:   "Will be removed - it is added automatically on the SDK level.",
	},
	"direction": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"INBOUND", "OUTBOUND"}, true),
		Description:  "Direction of the cloud messaging with respect to Snowflake (required only for error notifications)",
		Deprecated:   "Will be removed - it is added automatically on the SDK level.",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return true
		},
	},
	// This part of the schema is the cloudProviderParams in the Snowflake documentation and differs between vendors
	"notification_provider": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"AZURE_STORAGE_QUEUE", "AWS_SQS", "AWS_SNS", "GCP_PUBSUB"}, true),
		Description:  "The third-party cloud message queuing service (supported values: AZURE_STORAGE_QUEUE, AWS_SNS, GCP_PUBSUB; AWS_SQS is deprecated and will be removed in the future provider versions)",
		ForceNew:     true,
	},
	"azure_storage_queue_primary_uri": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The queue ID for the Azure Queue Storage queue created for Event Grid notifications. Required for AZURE_STORAGE_QUEUE provider",
		// There is no alter SQL for azure_storage_queue_primary_uri for automated data loads, therefore it has to be recreated.
		ForceNew: true,
	},
	"azure_tenant_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the Azure Active Directory tenant used for identity management. Required for AZURE_STORAGE_QUEUE provider",
		// There is no alter SQL for azure_tenant_id for automated data loads, therefore it has to be recreated.
		ForceNew: true,
	},
	"aws_sqs_external_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role",
		Deprecated:  "No longer supported notification method",
	},
	"aws_sqs_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
		Deprecated:  "No longer supported notification method",
	},
	"aws_sqs_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "AWS SQS queue ARN for notification integration to connect to",
		Deprecated:  "No longer supported notification method",
	},
	"aws_sqs_role_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "AWS IAM role ARN for notification integration to assume",
		Deprecated:  "No longer supported notification method",
	},
	"aws_sns_external_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role",
	},
	"aws_sns_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
	},
	"aws_sns_topic_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "AWS SNS Topic ARN for notification integration to connect to. Required for AWS_SNS provider.",
	},
	"aws_sns_role_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "AWS IAM role ARN for notification integration to assume. Required for AWS_SNS provider",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A comment for the integration",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the notification integration was created.",
	},
	"gcp_pubsub_subscription_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The subscription id that Snowflake will listen to when using the GCP_PUBSUB provider.",
		// There is no alter SQL for gcp_pubsub_subscription_name for automated data loads, therefore it has to be recreated.
		ForceNew:      true,
		ConflictsWith: []string{"gcp_pubsub_topic_name"},
	},
	"gcp_pubsub_topic_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The topic id that Snowflake will use to push notifications.",
		// There is no alter SQL for gcp_pubsub_topic_name, therefore it has to be recreated.
		ForceNew:      true,
		ConflictsWith: []string{"gcp_pubsub_subscription_name"},
	},
	"gcp_pubsub_service_account": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The GCP service account identifier that Snowflake will use when assuming the GCP role",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// NotificationIntegration returns a pointer to the resource representing a notification integration.
func NotificationIntegration() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErr[sdk.AccountObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.NotificationIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.NotificationIntegrationResource), TrackingCreateWrapper(resources.NotificationIntegration, CreateNotificationIntegration)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.NotificationIntegrationResource), TrackingReadWrapper(resources.NotificationIntegration, ReadNotificationIntegration)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.NotificationIntegrationResource), TrackingUpdateWrapper(resources.NotificationIntegration, UpdateNotificationIntegration)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.NotificationIntegrationResource), TrackingDeleteWrapper(resources.NotificationIntegration, deleteFunc)),

		Schema: notificationIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateNotificationIntegration implements schema.CreateFunc.
func CreateNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	enabled := d.Get("enabled").(bool)

	createRequest := sdk.NewCreateNotificationIntegrationRequest(id, enabled)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	notificationProvider := strings.ToUpper(d.Get("notification_provider").(string))
	switch notificationProvider {
	case "AWS_SNS":
		topic, ok := d.GetOk("aws_sns_topic_arn")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use AWS_SNS provider you must specify an aws_sns_topic_arn"))
		}
		role, ok := d.GetOk("aws_sns_role_arn")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use AWS_SNS provider you must specify an aws_sns_role_arn"))
		}
		createRequest.WithPushNotificationParams(
			sdk.NewPushNotificationParamsRequest().WithAmazonPushParams(sdk.NewAmazonPushParamsRequest(topic.(string), role.(string))),
		)
	case "GCP_PUBSUB":
		if v, ok := d.GetOk("gcp_pubsub_subscription_name"); ok {
			createRequest.WithAutomatedDataLoadsParams(
				sdk.NewAutomatedDataLoadsParamsRequest().WithGoogleAutoParams(sdk.NewGoogleAutoParamsRequest(v.(string))),
			)
		}
		if v, ok := d.GetOk("gcp_pubsub_topic_name"); ok {
			createRequest.WithPushNotificationParams(
				sdk.NewPushNotificationParamsRequest().WithGooglePushParams(sdk.NewGooglePushParamsRequest(v.(string))),
			)
		}
	case "AZURE_STORAGE_QUEUE":
		uri, ok := d.GetOk("azure_storage_queue_primary_uri")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use AZURE_STORAGE_QUEUE provider you must specify an azure_storage_queue_primary_uri"))
		}
		tenantId, ok := d.GetOk("azure_tenant_id")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use AZURE_STORAGE_QUEUE provider you must specify an azure_tenant_id"))
		}
		createRequest.WithAutomatedDataLoadsParams(
			sdk.NewAutomatedDataLoadsParamsRequest().WithAzureAutoParams(sdk.NewAzureAutoParamsRequest(uri.(string), tenantId.(string))),
		)
	default:
		return diag.FromErr(fmt.Errorf("unexpected provider %v", notificationProvider))
	}

	err := client.NotificationIntegrations.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating notification integration: %w", err))
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadNotificationIntegration(ctx, d, meta)
}

// ReadNotificationIntegration implements schema.ReadFunc.
func ReadNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.NotificationIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query notification integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Notification integration id: %s, Err: %s", d.Id(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	// Note: category must be NOTIFICATION or something is broken
	if c := integration.Category; c != "NOTIFICATION" {
		return diag.FromErr(fmt.Errorf("expected %v to be a NOTIFICATION integration, got %v", id, c))
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", integration.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", integration.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_on", integration.CreatedOn.String()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", integration.Enabled); err != nil {
		return diag.FromErr(err)
	}

	// Snowflake returns "QUEUE - AZURE_STORAGE_QUEUE" instead of simple "QUEUE" as a type
	// It needs to be parsed in order to not show a diff in Terraform.
	typeParts := strings.Split(integration.NotificationType, "-")
	parsedType := strings.TrimSpace(typeParts[0])
	if err := d.Set("type", parsedType); err != nil {
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
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "DIRECTION":
			if err := d.Set("direction", value); err != nil {
				return diag.FromErr(err)
			}
		case "NOTIFICATION_PROVIDER":
			if err := d.Set("notification_provider", value); err != nil {
				return diag.FromErr(err)
			}
		case "AZURE_STORAGE_QUEUE_PRIMARY_URI":
			if err := d.Set("azure_storage_queue_primary_uri", value); err != nil {
				return diag.FromErr(err)
			}
			// NOTIFICATION_PROVIDER is not returned for azure automated data load, so we set it manually in such a case
			if err := d.Set("notification_provider", "AZURE_STORAGE_QUEUE"); err != nil {
				return diag.FromErr(err)
			}
		case "AZURE_TENANT_ID":
			if err := d.Set("azure_tenant_id", value); err != nil {
				return diag.FromErr(err)
			}
		case "AWS_SNS_TOPIC_ARN":
			if err := d.Set("aws_sns_topic_arn", value); err != nil {
				return diag.FromErr(err)
			}
		case "AWS_SNS_ROLE_ARN":
			if err := d.Set("aws_sns_role_arn", value); err != nil {
				return diag.FromErr(err)
			}
		case "SF_AWS_EXTERNAL_ID":
			if err := d.Set("aws_sns_external_id", value); err != nil {
				return diag.FromErr(err)
			}
		case "SF_AWS_IAM_USER_ARN":
			if err := d.Set("aws_sns_iam_user_arn", value); err != nil {
				return diag.FromErr(err)
			}
		case "GCP_PUBSUB_SUBSCRIPTION_NAME":
			if err := d.Set("gcp_pubsub_subscription_name", value); err != nil {
				return diag.FromErr(err)
			}
			// NOTIFICATION_PROVIDER is not returned for gcp, so we set it manually in such a case
			if err := d.Set("notification_provider", "GCP_PUBSUB"); err != nil {
				return diag.FromErr(err)
			}
		case "GCP_PUBSUB_TOPIC_NAME":
			if err := d.Set("gcp_pubsub_topic_name", value); err != nil {
				return diag.FromErr(err)
			}
			// NOTIFICATION_PROVIDER is not returned for gcp, so we set it manually in such a case
			if err := d.Set("notification_provider", "GCP_PUBSUB"); err != nil {
				return diag.FromErr(err)
			}
		case "GCP_PUBSUB_SERVICE_ACCOUNT":
			if err := d.Set("gcp_pubsub_service_account", value); err != nil {
				return diag.FromErr(err)
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", name)
		}
	}

	return diag.FromErr(err)
}

// UpdateNotificationIntegration implements schema.UpdateFunc.
func UpdateNotificationIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	var runSetStatement bool
	setRequest := sdk.NewNotificationIntegrationSetRequest()
	if d.HasChange("comment") {
		runSetStatement = true
		setRequest.WithComment(sdk.String(d.Get("comment").(string)))
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		setRequest.WithEnabled(sdk.Bool(d.Get("enabled").(bool)))
	}

	notificationProvider := strings.ToUpper(d.Get("notification_provider").(string))
	switch notificationProvider {
	case "AWS_SNS":
		if d.HasChange("aws_sns_topic_arn") || d.HasChange("aws_sns_role_arn") {
			runSetStatement = true
			setAmazonPush := sdk.NewSetAmazonPushRequest(d.Get("aws_sns_topic_arn").(string), d.Get("aws_sns_role_arn").(string))
			setRequest.WithSetPushParams(sdk.NewSetPushParamsRequest().WithSetAmazonPush(setAmazonPush))
		}
	case "GCP_PUBSUB":
		log.Printf("[WARN] all GCP_PUBSUB properties should recreate the resource")
	case "AZURE_STORAGE_QUEUE":
		log.Printf("[WARN] all AZURE_STORAGE_QUEUE properties should recreate the resource")
	default:
		return diag.FromErr(fmt.Errorf("unexpected provider %v", notificationProvider))
	}

	if runSetStatement {
		err := client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithSet(setRequest))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating notification integration: %w", err))
		}
	}

	return ReadNotificationIntegration(ctx, d, meta)
}
