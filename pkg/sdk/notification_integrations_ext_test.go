package sdk

func init() {
	id := notificationIntegrationsTestIdAccountObjectIdentifier
	secretId := NewSchemaObjectIdentifier("metrics_catalog", "metric_config", "slack_integration_webhook")

	apiAwsRoleArn := "arn:aws:iam::000000000001:/role/test"
	azureTenantId := "00000000-0000-0000-0000-000000000000"
	gcpPubsubSubscriptionName := "projects/project-1234/subscriptions/sub2"
	gcpPubsubTopicName := "projects/project-1234/topics/top2"
	azureStorageQueuePrimaryUri := "azure://great-bucket/great-path/"
	azureEventGridTopicEndpoint := "https://apim-hello-world.azure-api.net/dev"
	awsSnsTopicArn := "arn:aws:sns:us-east-2:123456789012:MyTopic"
	webhookUrl := "https://hooks.slack.com/services/SNOWFLAKE_WEBHOOK_SECRET"
	email := "some.email@some.com"
	otherEmail := "some.other.email@some.com"

	notificationIntegrationsTests.Create.
		withDefaultOpts(func() *CreateNotificationIntegrationOptions {
			return &CreateNotificationIntegrationOptions{
				name:        id,
				Enabled:     true,
				EmailParams: &EmailParams{},
			}
		}).
		withExpectedSqlf(
			case_NotificationIntegrations_sql_Create_basic,
			"CREATE NOTIFICATION INTEGRATION %s ENABLED = true TYPE = EMAIL",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Create_all,
			func(opts *CreateNotificationIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
				opts.EmailParams.AllowedRecipients = []NotificationIntegrationAllowedRecipient{
					{Email: email},
					{Email: otherEmail},
				}
			},
			"CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true TYPE = EMAIL ALLOWED_RECIPIENTS = ('%s', '%s') COMMENT = 'some comment'",
			id.FullyQualifiedName(), email, otherEmail,
		)

	notificationIntegrationsTests.Create.
		withAdditionalSqlCasef(
			"sql_Create_autoGoogle",
			func(opts *CreateNotificationIntegrationOptions) {
				opts.EmailParams = nil
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
				opts.AutomatedDataLoadsParams = &AutomatedDataLoadsParams{
					GoogleAutoParams: &GoogleAutoParams{GcpPubsubSubscriptionName: gcpPubsubSubscriptionName},
				}
			},
			"CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true TYPE = QUEUE NOTIFICATION_PROVIDER = GCP_PUBSUB GCP_PUBSUB_SUBSCRIPTION_NAME = '%s' COMMENT = 'some comment'",
			id.FullyQualifiedName(), gcpPubsubSubscriptionName,
		).
		withAdditionalSqlCasef(
			"sql_Create_autoAzure",
			func(opts *CreateNotificationIntegrationOptions) {
				opts.EmailParams = nil
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
				opts.AutomatedDataLoadsParams = &AutomatedDataLoadsParams{
					AzureAutoParams: &AzureAutoParams{
						AzureStorageQueuePrimaryUri: azureStorageQueuePrimaryUri,
						AzureTenantId:               azureTenantId,
					},
				}
			},
			"CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true TYPE = QUEUE NOTIFICATION_PROVIDER = AZURE_STORAGE_QUEUE AZURE_STORAGE_QUEUE_PRIMARY_URI = '%s' AZURE_TENANT_ID = '%s' COMMENT = 'some comment'",
			id.FullyQualifiedName(), azureStorageQueuePrimaryUri, azureTenantId,
		).
		withAdditionalSqlCasef(
			"sql_Create_pushAmazon",
			func(opts *CreateNotificationIntegrationOptions) {
				opts.EmailParams = nil
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
				opts.PushNotificationParams = &PushNotificationParams{
					AmazonPushParams: &AmazonPushParams{
						AwsSnsTopicArn: awsSnsTopicArn,
						AwsSnsRoleArn:  apiAwsRoleArn,
					},
				}
			},
			"CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true DIRECTION = OUTBOUND TYPE = QUEUE NOTIFICATION_PROVIDER = AWS_SNS AWS_SNS_TOPIC_ARN = '%s' AWS_SNS_ROLE_ARN = '%s' COMMENT = 'some comment'",
			id.FullyQualifiedName(), awsSnsTopicArn, apiAwsRoleArn,
		).
		withAdditionalSqlCasef(
			"sql_Create_pushGoogle",
			func(opts *CreateNotificationIntegrationOptions) {
				opts.EmailParams = nil
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
				opts.PushNotificationParams = &PushNotificationParams{
					GooglePushParams: &GooglePushParams{GcpPubsubTopicName: gcpPubsubTopicName},
				}
			},
			"CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true DIRECTION = OUTBOUND TYPE = QUEUE NOTIFICATION_PROVIDER = GCP_PUBSUB GCP_PUBSUB_TOPIC_NAME = '%s' COMMENT = 'some comment'",
			id.FullyQualifiedName(), gcpPubsubTopicName,
		).
		withAdditionalSqlCasef(
			"sql_Create_pushAzure",
			func(opts *CreateNotificationIntegrationOptions) {
				opts.EmailParams = nil
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
				opts.PushNotificationParams = &PushNotificationParams{
					AzurePushParams: &AzurePushParams{
						AzureEventGridTopicEndpoint: azureEventGridTopicEndpoint,
						AzureTenantId:               azureTenantId,
					},
				}
			},
			"CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true DIRECTION = OUTBOUND TYPE = QUEUE NOTIFICATION_PROVIDER = AZURE_EVENT_GRID AZURE_EVENT_GRID_TOPIC_ENDPOINT = '%s' AZURE_TENANT_ID = '%s' COMMENT = 'some comment'",
			id.FullyQualifiedName(), azureEventGridTopicEndpoint, azureTenantId,
		).
		withAdditionalSqlCasef(
			"sql_Create_webhookMinimal",
			func(opts *CreateNotificationIntegrationOptions) {
				opts.EmailParams = nil
				opts.WebhookParams = &WebhookParams{WebhookUrl: webhookUrl}
			},
			"CREATE NOTIFICATION INTEGRATION %s ENABLED = true TYPE = WEBHOOK WEBHOOK_URL = '%s'",
			id.FullyQualifiedName(), webhookUrl,
		).
		withAdditionalSqlCasef(
			"sql_Create_webhookWithSecretAndHeaders",
			func(opts *CreateNotificationIntegrationOptions) {
				opts.EmailParams = nil
				opts.IfNotExists = new(true)
				opts.Comment = new("slack webhook")
				opts.WebhookParams = &WebhookParams{
					WebhookUrl:          webhookUrl,
					WebhookSecret:       &secretId,
					WebhookBodyTemplate: new("SNOWFLAKE_WEBHOOK_MESSAGE"),
					WebhookHeaders: []WebhookHeader{
						{Header: "Content-Type", Value: "application/json"},
					},
				}
			},
			`CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true TYPE = WEBHOOK WEBHOOK_URL = '%s' WEBHOOK_SECRET = %s WEBHOOK_BODY_TEMPLATE = 'SNOWFLAKE_WEBHOOK_MESSAGE' WEBHOOK_HEADERS = ('Content-Type' = 'application/json') COMMENT = 'slack webhook'`,
			id.FullyQualifiedName(), webhookUrl, secretId.FullyQualifiedName(),
		)

	notificationIntegrationsTests.Alter.
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Alter_Set,
			func(opts *AlterNotificationIntegrationOptions) {
				opts.Set = &NotificationIntegrationSet{
					Enabled: new(true),
					Comment: new("some comment"),
				}
			},
			"ALTER NOTIFICATION INTEGRATION %s SET ENABLED = true COMMENT = 'some comment'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Alter_UnsetEmailParams,
			func(opts *AlterNotificationIntegrationOptions) {
				opts.UnsetEmailParams = &NotificationIntegrationUnsetEmailParams{
					Comment: new(true),
				}
			},
			"ALTER NOTIFICATION INTEGRATION %s UNSET COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Alter_UnsetWebhookParams,
			func(opts *AlterNotificationIntegrationOptions) {
				opts.UnsetWebhookParams = &NotificationIntegrationUnsetWebhookParams{
					WebhookSecret:       new(true),
					WebhookBodyTemplate: new(true),
					WebhookHeaders:      new(true),
					Comment:             new(true),
				}
			},
			"ALTER NOTIFICATION INTEGRATION %s UNSET WEBHOOK_SECRET, WEBHOOK_BODY_TEMPLATE, WEBHOOK_HEADERS, COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Alter_SetTags,
			func(opts *AlterNotificationIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER NOTIFICATION INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Alter_UnsetTags,
			func(opts *AlterNotificationIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER NOTIFICATION INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	notificationIntegrationsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Set_pushAmazon",
			func(opts *AlterNotificationIntegrationOptions) {
				opts.Set = &NotificationIntegrationSet{
					Enabled: new(true),
					SetPushParams: &SetPushParams{
						SetAmazonPush: &SetAmazonPush{
							AwsSnsTopicArn: awsSnsTopicArn,
							AwsSnsRoleArn:  apiAwsRoleArn,
						},
					},
					Comment: new("some comment"),
				}
			},
			"ALTER NOTIFICATION INTEGRATION %s SET ENABLED = true AWS_SNS_TOPIC_ARN = '%s' AWS_SNS_ROLE_ARN = '%s' COMMENT = 'some comment'",
			id.FullyQualifiedName(), awsSnsTopicArn, apiAwsRoleArn,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_pushGoogle",
			func(opts *AlterNotificationIntegrationOptions) {
				opts.Set = &NotificationIntegrationSet{
					Enabled: new(true),
					SetPushParams: &SetPushParams{
						SetGooglePush: &SetGooglePush{
							GcpPubsubSubscriptionName: gcpPubsubSubscriptionName,
						},
					},
					Comment: new("some comment"),
				}
			},
			"ALTER NOTIFICATION INTEGRATION %s SET ENABLED = true GCP_PUBSUB_SUBSCRIPTION_NAME = '%s' COMMENT = 'some comment'",
			id.FullyQualifiedName(), gcpPubsubSubscriptionName,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_pushAzure",
			func(opts *AlterNotificationIntegrationOptions) {
				opts.Set = &NotificationIntegrationSet{
					Enabled: new(true),
					SetPushParams: &SetPushParams{
						SetAzurePush: &SetAzurePush{
							AzureStorageQueuePrimaryUri: azureStorageQueuePrimaryUri,
							AzureTenantId:               azureTenantId,
						},
					},
					Comment: new("some comment"),
				}
			},
			"ALTER NOTIFICATION INTEGRATION %s SET ENABLED = true AZURE_STORAGE_QUEUE_PRIMARY_URI = '%s' AZURE_TENANT_ID = '%s' COMMENT = 'some comment'",
			id.FullyQualifiedName(), azureStorageQueuePrimaryUri, azureTenantId,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_email",
			func(opts *AlterNotificationIntegrationOptions) {
				opts.Set = &NotificationIntegrationSet{
					Enabled: new(true),
					SetEmailParams: &SetEmailParams{
						AllowedRecipients: []NotificationIntegrationAllowedRecipient{
							{Email: email},
							{Email: otherEmail},
						},
					},
					Comment: new("some comment"),
				}
			},
			"ALTER NOTIFICATION INTEGRATION %s SET ENABLED = true ALLOWED_RECIPIENTS = ('%s', '%s') COMMENT = 'some comment'",
			id.FullyQualifiedName(), email, otherEmail,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_webhook",
			func(opts *AlterNotificationIntegrationOptions) {
				opts.Set = &NotificationIntegrationSet{
					Enabled: new(true),
					SetWebhookParams: &SetWebhookParams{
						WebhookUrl:          new(webhookUrl),
						WebhookSecret:       &secretId,
						WebhookBodyTemplate: new("SNOWFLAKE_WEBHOOK_MESSAGE"),
						WebhookHeaders: []WebhookHeader{
							{Header: "Content-Type", Value: "application/json"},
						},
					},
					Comment: new("some comment"),
				}
			},
			`ALTER NOTIFICATION INTEGRATION %s SET ENABLED = true WEBHOOK_URL = '%s' WEBHOOK_SECRET = %s WEBHOOK_BODY_TEMPLATE = 'SNOWFLAKE_WEBHOOK_MESSAGE' WEBHOOK_HEADERS = ('Content-Type' = 'application/json') COMMENT = 'some comment'`,
			id.FullyQualifiedName(), webhookUrl, secretId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_UnsetEmailParams_multiple",
			func(opts *AlterNotificationIntegrationOptions) {
				opts.UnsetEmailParams = &NotificationIntegrationUnsetEmailParams{
					AllowedRecipients: new(true),
					Comment:           new(true),
				}
			},
			"ALTER NOTIFICATION INTEGRATION %s UNSET ALLOWED_RECIPIENTS, COMMENT",
			id.FullyQualifiedName(),
		)

	notificationIntegrationsTests.Drop.
		withExpectedSqlf(
			case_NotificationIntegrations_sql_Drop_basic,
			"DROP NOTIFICATION INTEGRATION %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Drop_all,
			func(opts *DropNotificationIntegrationOptions) { opts.IfExists = new(true) },
			"DROP NOTIFICATION INTEGRATION IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	notificationIntegrationsTests.Show.
		withExpectedSql(case_NotificationIntegrations_sql_Show_basic, "SHOW NOTIFICATION INTEGRATIONS").
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Show_all,
			func(opts *ShowNotificationIntegrationOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
			},
			"SHOW NOTIFICATION INTEGRATIONS LIKE '%s'",
			id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_NotificationIntegrations_sql_Show_Like,
			func(opts *ShowNotificationIntegrationOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
			},
			"SHOW NOTIFICATION INTEGRATIONS LIKE '%s'",
			id.Name(),
		)

	notificationIntegrationsTests.Describe.
		withExpectedSqlf(
			case_NotificationIntegrations_sql_Describe_basic,
			"DESCRIBE NOTIFICATION INTEGRATION %s",
			id.FullyQualifiedName(),
		)
}
