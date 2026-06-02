package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var NotificationIntegrationAllowedRecipientDef = g.NewQueryStruct("NotificationIntegrationAllowedRecipient").
	Text("Email", g.KeywordOptions().SingleQuotes().Required())

var NotificationIntegrationWebhookHeaderDef = g.NewQueryStruct("WebhookHeader").
	Text("Header", g.KeywordOptions().SingleQuotes().Required()).
	SQLWithCustomFieldName("equals", "=").
	Text("Value", g.KeywordOptions().SingleQuotes().Required())

// TODO [SNOW-1016561]: all integrations reuse almost the same show, drop, and describe. For now we are copying it. Consider reusing in linked issue.
var notificationIntegrationsDef = g.NewInterface(
	"NotificationIntegrations",
	"NotificationIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration",
		g.NewQueryStruct("CreateNotificationIntegration").
			Create().
			OrReplace().
			SQL("NOTIFICATION INTEGRATION").
			IfNotExists().
			Name().
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalQueryStructField(
				"AutomatedDataLoadsParams",
				g.NewQueryStruct("AutomatedDataLoadsParams").
					SQLWithCustomFieldName("notificationType", "TYPE = QUEUE").
					OptionalQueryStructField(
						"GoogleAutoParams",
						g.NewQueryStruct("GoogleAutoParams").
							SQLWithCustomFieldName("notificationProvider", "NOTIFICATION_PROVIDER = GCP_PUBSUB").
							TextAssignment("GCP_PUBSUB_SUBSCRIPTION_NAME", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"AzureAutoParams",
						g.NewQueryStruct("AzureAutoParams").
							SQLWithCustomFieldName("notificationProvider", "NOTIFICATION_PROVIDER = AZURE_STORAGE_QUEUE").
							TextAssignment("AZURE_STORAGE_QUEUE_PRIMARY_URI", g.ParameterOptions().SingleQuotes().Required()).
							TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					WithValidation(g.ExactlyOneValueSet, "GoogleAutoParams", "AzureAutoParams"),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"PushNotificationParams",
				g.NewQueryStruct("PushNotificationParams").
					SQLWithCustomFieldName("direction", "DIRECTION = OUTBOUND").
					SQLWithCustomFieldName("notificationType", "TYPE = QUEUE").
					OptionalQueryStructField(
						"AmazonPushParams",
						g.NewQueryStruct("AmazonPushParams").
							SQLWithCustomFieldName("notificationProvider", "NOTIFICATION_PROVIDER = AWS_SNS").
							TextAssignment("AWS_SNS_TOPIC_ARN", g.ParameterOptions().SingleQuotes().Required()).
							TextAssignment("AWS_SNS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"GooglePushParams",
						g.NewQueryStruct("GooglePushParams").
							SQLWithCustomFieldName("notificationProvider", "NOTIFICATION_PROVIDER = GCP_PUBSUB").
							TextAssignment("GCP_PUBSUB_TOPIC_NAME", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"AzurePushParams",
						g.NewQueryStruct("AzurePushParams").
							SQLWithCustomFieldName("notificationProvider", "NOTIFICATION_PROVIDER = AZURE_EVENT_GRID").
							TextAssignment("AZURE_EVENT_GRID_TOPIC_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
							TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					WithValidation(g.ExactlyOneValueSet, "AmazonPushParams", "GooglePushParams", "AzurePushParams"),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"EmailParams",
				g.NewQueryStruct("EmailParams").
					SQLWithCustomFieldName("notificationType", "TYPE = EMAIL").
					ListAssignment("ALLOWED_RECIPIENTS", "NotificationIntegrationAllowedRecipient", g.ParameterOptions().Parentheses()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"WebhookParams",
				g.NewQueryStruct("WebhookParams").
					SQLWithCustomFieldName("webhookType", "TYPE = WEBHOOK").
					TextAssignment("WEBHOOK_URL", g.ParameterOptions().SingleQuotes().Required()).
					OptionalIdentifier("WebhookSecret", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WEBHOOK_SECRET")).
					OptionalTextAssignment("WEBHOOK_BODY_TEMPLATE", g.ParameterOptions().SingleQuotes()).
					ListQueryStructField("WebhookHeaders", NotificationIntegrationWebhookHeaderDef, g.ParameterOptions().SQL("WEBHOOK_HEADERS").Parentheses()),
				g.KeywordOptions(),
			).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "AutomatedDataLoadsParams", "PushNotificationParams", "EmailParams", "WebhookParams"),
		NotificationIntegrationAllowedRecipientDef,
		NotificationIntegrationWebhookHeaderDef,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration",
		g.NewQueryStruct("AlterNotificationIntegration").
			Alter().
			SQL("NOTIFICATION INTEGRATION").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("NotificationIntegrationSet").
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					OptionalQueryStructField(
						"SetPushParams",
						g.NewQueryStruct("SetPushParams").
							OptionalQueryStructField(
								"SetAmazonPush",
								g.NewQueryStruct("SetAmazonPush").
									TextAssignment("AWS_SNS_TOPIC_ARN", g.ParameterOptions().SingleQuotes().Required()).
									TextAssignment("AWS_SNS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()),
								g.KeywordOptions(),
							).
							OptionalQueryStructField(
								"SetGooglePush",
								g.NewQueryStruct("SetGooglePush").
									TextAssignment("GCP_PUBSUB_SUBSCRIPTION_NAME", g.ParameterOptions().SingleQuotes().Required()),
								g.KeywordOptions(),
							).
							OptionalQueryStructField(
								"SetAzurePush",
								g.NewQueryStruct("SetAzurePush").
									TextAssignment("AZURE_STORAGE_QUEUE_PRIMARY_URI", g.ParameterOptions().SingleQuotes().Required()).
									TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()),
								g.KeywordOptions(),
							).
							WithValidation(g.ExactlyOneValueSet, "SetAmazonPush", "SetGooglePush", "SetAzurePush"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"SetEmailParams",
						g.NewQueryStruct("SetEmailParams").
							ListAssignment("ALLOWED_RECIPIENTS", "NotificationIntegrationAllowedRecipient", g.ParameterOptions().Parentheses().Required()).
							WithValidation(g.ValidateValueSet, "AllowedRecipients"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"SetWebhookParams",
						g.NewQueryStruct("SetWebhookParams").
							OptionalTextAssignment("WEBHOOK_URL", g.ParameterOptions().SingleQuotes()).
							OptionalIdentifier("WebhookSecret", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WEBHOOK_SECRET")).
							OptionalTextAssignment("WEBHOOK_BODY_TEMPLATE", g.ParameterOptions().SingleQuotes()).
							ListQueryStructField("WebhookHeaders", NotificationIntegrationWebhookHeaderDef, g.ParameterOptions().SQL("WEBHOOK_HEADERS").Parentheses()),
						g.KeywordOptions(),
					).
					OptionalComment().
					WithValidation(g.MoreThanOneValueSet, "SetPushParams", "SetEmailParams", "SetWebhookParams").
					WithValidation(g.AtLeastOneValueSet, "Enabled", "SetPushParams", "SetEmailParams", "SetWebhookParams", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			// UNSET is supported only for the email notifications
			OptionalQueryStructField(
				"UnsetEmailParams",
				g.NewQueryStruct("NotificationIntegrationUnsetEmailParams").
					OptionalSQL("ALLOWED_RECIPIENTS").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "AllowedRecipients", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalQueryStructField(
				"UnsetWebhookParams",
				g.NewQueryStruct("NotificationIntegrationUnsetWebhookParams").
					OptionalSQL("WEBHOOK_SECRET").
					OptionalSQL("WEBHOOK_BODY_TEMPLATE").
					OptionalSQL("WEBHOOK_HEADERS").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "WebhookSecret", "WebhookBodyTemplate", "WebhookHeaders", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "UnsetEmailParams", "UnsetWebhookParams", "SetTags", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropNotificationIntegration").
			Drop().
			SQL("NOTIFICATION INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.StructPair("showNotificationIntegrationsDbRow", "NotificationIntegration").
			Text("name").
			Text("type", g.WithPlainFieldName("NotificationType")).
			Text("category").
			Bool("enabled").
			OptionalText("comment", g.WithRequiredInPlain()).
			Time("created_on"),
		g.NewQueryStruct("ShowNotificationIntegrations").
			Show().
			SQL("NOTIFICATION INTEGRATIONS").
			OptionalLike(),
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.StructPair("descNotificationIntegrationsDbRow", "NotificationIntegrationProperty").
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_type", g.WithPlainFieldName("Type")).
			Text("property_value", g.WithPlainFieldName("Value")).
			Text("property_default", g.WithPlainFieldName("Default")),
		g.NewQueryStruct("DescribeNotificationIntegration").
			Describe().
			SQL("NOTIFICATION INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithShowObjectType("Integration")
