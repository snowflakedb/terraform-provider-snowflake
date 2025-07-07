//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_NotificationIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// TODO [SNOW-1017580]: replace with real values
	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"
	const gcpPubsubTopicName = "projects/project-1234/topics/top2"
	const azureStorageQueuePrimaryUri = "azure://great-bucket/great-path/"
	const azureTenantId = "00000000-0000-0000-0000-000000000000"
	const azureEventGridTopicEndpoint = "https://apim-hello-world.azure-api.net/dev"
	const awsSnsTopicArn = "arn:aws:sns:us-east-2:123456789012:MyTopic"
	const awsSnsOtherTopicArn = "arn:aws:sns:us-east-2:123456789012:MyOtherTopic"
	const awsSnsRoleArn = "arn:aws:iam::000000000001:/role/test"
	const awsSnsOtherRoleArn = "arn:aws:iam::000000000001:/role/other"

	assertNotificationIntegration := func(t *testing.T, s *sdk.NotificationIntegration, name sdk.AccountObjectIdentifier, notificationType string, comment string) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.True(t, s.Enabled)
		assert.Equal(t, notificationType, s.NotificationType)
		assert.Equal(t, "NOTIFICATION", s.Category)
		assert.Equal(t, comment, s.Comment)
	}

	cleanupNotificationIntegrationProvider := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.NotificationIntegrations.Drop(ctx, sdk.NewDropNotificationIntegrationRequest(id))
			require.NoError(t, err)
		}
	}

	createNotificationIntegrationAutoGoogleRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithAutomatedDataLoadsParams(sdk.NewAutomatedDataLoadsParamsRequest().WithGoogleAutoParams(sdk.NewGoogleAutoParamsRequest(gcpPubsubSubscriptionName)))
	}

	createNotificationIntegrationAutoAzureRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithAutomatedDataLoadsParams(sdk.NewAutomatedDataLoadsParamsRequest().WithAzureAutoParams(sdk.NewAzureAutoParamsRequest(azureStorageQueuePrimaryUri, azureTenantId)))
	}

	createNotificationIntegrationPushAmazonRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithPushNotificationParams(sdk.NewPushNotificationParamsRequest().WithAmazonPushParams(sdk.NewAmazonPushParamsRequest(awsSnsTopicArn, awsSnsRoleArn)))
	}

	createNotificationIntegrationPushGoogleRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithPushNotificationParams(sdk.NewPushNotificationParamsRequest().WithGooglePushParams(sdk.NewGooglePushParamsRequest(gcpPubsubTopicName)))
	}

	createNotificationIntegrationPushAzureRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithPushNotificationParams(sdk.NewPushNotificationParamsRequest().WithAzurePushParams(sdk.NewAzurePushParamsRequest(azureEventGridTopicEndpoint, azureTenantId)))
	}

	createNotificationIntegrationEmailRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		// TODO [SNOW-1007539]: use email of our service user
		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithEmailParams(sdk.NewEmailParamsRequest().WithAllowedRecipients([]sdk.NotificationIntegrationAllowedRecipient{{Email: "artur.sawicki@snowflake.com"}}))
	}

	createNotificationIntegrationWithRequest := func(t *testing.T, request *sdk.CreateNotificationIntegrationRequest) *sdk.NotificationIntegration {
		t.Helper()
		id := request.GetName()

		err := client.NotificationIntegrations.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupNotificationIntegrationProvider(id))

		integration, err := client.NotificationIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration
	}

	createNotificationIntegrationAutoGoogle := func(t *testing.T) *sdk.NotificationIntegration {
		t.Helper()
		return createNotificationIntegrationWithRequest(t, createNotificationIntegrationAutoGoogleRequest(t))
	}

	createNotificationIntegrationAutoAzure := func(t *testing.T) *sdk.NotificationIntegration {
		t.Helper()
		return createNotificationIntegrationWithRequest(t, createNotificationIntegrationAutoAzureRequest(t))
	}

	createNotificationIntegrationPushAmazon := func(t *testing.T) *sdk.NotificationIntegration {
		t.Helper()
		return createNotificationIntegrationWithRequest(t, createNotificationIntegrationPushAmazonRequest(t))
	}

	createNotificationIntegrationEmail := func(t *testing.T) *sdk.NotificationIntegration {
		t.Helper()
		return createNotificationIntegrationWithRequest(t, createNotificationIntegrationEmailRequest(t))
	}

	t.Run("create and describe notification integration - auto google", func(t *testing.T) {
		request := createNotificationIntegrationAutoGoogleRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - GCP_PUBSUB", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "DIRECTION", Type: "String", Value: "INBOUND", Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "GCP_PUBSUB_SUBSCRIPTION_NAME", Type: "String", Value: gcpPubsubSubscriptionName, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})

		prop, err := collections.FindFirst(details, func(property sdk.NotificationIntegrationProperty) bool {
			return property.Name == "GCP_PUBSUB_SERVICE_ACCOUNT"
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, prop.Value)
	})

	t.Run("create and describe notification integration - auto azure", func(t *testing.T) {
		request := createNotificationIntegrationAutoAzureRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - AZURE_STORAGE_QUEUE", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AZURE_STORAGE_QUEUE_PRIMARY_URI", Type: "String", Value: azureStorageQueuePrimaryUri, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})

		prop, err := collections.FindFirst(details, func(property sdk.NotificationIntegrationProperty) bool {
			return property.Name == "AZURE_CONSENT_URL"
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, prop.Value)

		prop, err = collections.FindFirst(details, func(property sdk.NotificationIntegrationProperty) bool {
			return property.Name == "AZURE_MULTI_TENANT_APP_NAME"
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, prop.Value)
	})

	t.Run("create and describe notification integration - push amazon", func(t *testing.T) {
		request := createNotificationIntegrationPushAmazonRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - AWS_SNS", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "NOTIFICATION_PROVIDER", Type: "String", Value: "AWS_SNS", Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "DIRECTION", Type: "String", Value: "OUTBOUND", Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AWS_SNS_TOPIC_ARN", Type: "String", Value: awsSnsTopicArn, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AWS_SNS_ROLE_ARN", Type: "String", Value: awsSnsRoleArn, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})

		prop, err := collections.FindFirst(details, func(property sdk.NotificationIntegrationProperty) bool { return property.Name == "SF_AWS_IAM_USER_ARN" })
		assert.NoError(t, err)
		assert.NotEmpty(t, prop.Value)
		prop, err = collections.FindFirst(details, func(property sdk.NotificationIntegrationProperty) bool { return property.Name == "SF_AWS_EXTERNAL_ID" })
		assert.NoError(t, err)
		assert.NotEmpty(t, prop.Value)
	})

	// TODO [SNOW-1017802]: check the error 001422 (22023): SQL compilation error: invalid value 'OUTBOUND' for property 'Direction'
	t.Run("create and describe notification integration - push google", func(t *testing.T) {
		t.Skip("Skipping because of the error: 001422 (22023): SQL compilation error: invalid value 'OUTBOUND' for property 'Direction'")
		request := createNotificationIntegrationPushGoogleRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - GCP_PUBSUB", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "NOTIFICATION_PROVIDER", Type: "String", Value: "GCP_PUBSUB", Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "DIRECTION", Type: "String", Value: "OUTBOUND", Default: "INBOUND"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "GCP_PUBSUB_TOPIC_NAME", Type: "String", Value: gcpPubsubTopicName, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	// TODO [SNOW-1017802]: check the error 001008 (22023): SQL compilation error: invalid value [QUEUE - AZURE_EVENT_GRID] for parameter 'Integration Type'
	t.Run("create and describe notification integration - push azure", func(t *testing.T) {
		t.Skip("Skipping because of the error: 001008 (22023): SQL compilation error: invalid value [QUEUE - AZURE_EVENT_GRID] for parameter 'Integration Type'")
		request := createNotificationIntegrationPushAzureRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - AZURE_EVENT_GRID", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "NOTIFICATION_PROVIDER", Type: "String", Value: "GCP_PUBSUB", Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "DIRECTION", Type: "String", Value: "OUTBOUND", Default: "INBOUND"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AZURE_EVENT_GRID_TOPIC_ENDPOINT", Type: "String", Value: azureEventGridTopicEndpoint, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AZURE_TENANT_ID", Type: "String", Value: azureTenantId, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("create and describe notification integration - email", func(t *testing.T) {
		request := createNotificationIntegrationEmailRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "EMAIL", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ALLOWED_RECIPIENTS", Type: "List", Value: "artur.sawicki@snowflake.com", Default: "[]"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("create and describe notification integration - email, with empty allowed recipients", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithEmailParams(sdk.NewEmailParamsRequest().WithAllowedRecipients([]sdk.NotificationIntegrationAllowedRecipient{}))

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "EMAIL", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ALLOWED_RECIPIENTS", Type: "List", Value: "", Default: "[]"})
	})

	t.Run("alter notification integration: auto", func(t *testing.T) {
		integration := createNotificationIntegrationAutoGoogle(t)

		setRequest := sdk.NewAlterNotificationIntegrationRequest(integration.ID()).
			WithSet(
				sdk.NewNotificationIntegrationSetRequest().
					WithEnabled(sdk.Bool(false)).
					WithComment(sdk.String("changed comment")),
			)
		err := client.NotificationIntegrations.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "changed comment", Default: ""})

		// only SET is tested because UNSET is unsupported: 000002 (0A000): Unsupported feature 'UNSET'
	})

	t.Run("alter notification integration: push amazon", func(t *testing.T) {
		integration := createNotificationIntegrationPushAmazon(t)

		setRequest := sdk.NewAlterNotificationIntegrationRequest(integration.ID()).
			WithSet(
				sdk.NewNotificationIntegrationSetRequest().
					WithEnabled(sdk.Bool(false)).
					WithSetPushParams(sdk.NewSetPushParamsRequest().WithSetAmazonPush(sdk.NewSetAmazonPushRequest(awsSnsOtherTopicArn, awsSnsOtherRoleArn))).
					WithComment(sdk.String("changed comment")),
			)
		err := client.NotificationIntegrations.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AWS_SNS_TOPIC_ARN", Type: "String", Value: awsSnsOtherTopicArn, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AWS_SNS_ROLE_ARN", Type: "String", Value: awsSnsOtherRoleArn, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "changed comment", Default: ""})

		// only SET is tested because UNSET is unsupported: 000002 (0A000): Unsupported feature 'UNSET'
	})

	// TODO [SNOW-1017802]: implement after "create and describe notification integration - push google" succeeds
	t.Run("alter notification integration: push google", func(t *testing.T) {
		t.Skip(`Skipping because "create and describe notification integration - push google" creation test is failing`)
	})

	// TODO [SNOW-1017802]: implement after "create and describe notification integration - push azure" succeeds
	t.Run("alter notification integration: push azure", func(t *testing.T) {
		t.Skip(`Skipping because "create and describe notification integration - push azure" creation test is failing`)
	})

	t.Run("alter notification integration: email", func(t *testing.T) {
		integration := createNotificationIntegrationEmail(t)

		setRequest := sdk.NewAlterNotificationIntegrationRequest(integration.ID()).
			WithSet(
				sdk.NewNotificationIntegrationSetRequest().
					WithEnabled(sdk.Bool(false)).
					WithSetEmailParams(sdk.NewSetEmailParamsRequest([]sdk.NotificationIntegrationAllowedRecipient{{Email: "jan.cieslak@snowflake.com"}})).
					WithComment(sdk.String("changed comment")),
			)
		err := client.NotificationIntegrations.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ALLOWED_RECIPIENTS", Type: "List", Value: "jan.cieslak@snowflake.com", Default: "[]"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "changed comment", Default: ""})

		unsetRequest := sdk.NewAlterNotificationIntegrationRequest(integration.ID()).
			WithUnsetEmailParams(
				sdk.NewNotificationIntegrationUnsetEmailParamsRequest().
					WithAllowedRecipients(sdk.Bool(true)).
					WithComment(sdk.Bool(true)),
			)
		err = client.NotificationIntegrations.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ALLOWED_RECIPIENTS", Type: "List", Value: "", Default: "[]"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("drop notification integration: existing", func(t *testing.T) {
		request := createNotificationIntegrationEmailRequest(t)
		id := request.GetName()

		err := client.NotificationIntegrations.Create(ctx, request)
		require.NoError(t, err)

		err = client.NotificationIntegrations.Drop(ctx, sdk.NewDropNotificationIntegrationRequest(id))
		require.NoError(t, err)

		_, err = client.NotificationIntegrations.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop notification integration: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		err := client.NotificationIntegrations.Drop(ctx, sdk.NewDropNotificationIntegrationRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	// TODO [SNOW-1017802]: Add missing integrations
	t.Run("show notification integration: default", func(t *testing.T) {
		notificationAutoGoogle := createNotificationIntegrationAutoGoogle(t)
		notificationAutoAzure := createNotificationIntegrationAutoAzure(t)
		notificationPushAmazon := createNotificationIntegrationPushAmazon(t)
		notificationEmail := createNotificationIntegrationEmail(t)

		showRequest := sdk.NewShowNotificationIntegrationRequest()
		returnedIntegrations, err := client.NotificationIntegrations.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *notificationAutoGoogle)
		assert.Contains(t, returnedIntegrations, *notificationAutoAzure)
		assert.Contains(t, returnedIntegrations, *notificationPushAmazon)
		assert.Contains(t, returnedIntegrations, *notificationEmail)
	})

	t.Run("show notification integration: with options", func(t *testing.T) {
		notificationAutoGoogle := createNotificationIntegrationAutoGoogle(t)
		notificationAutoAzure := createNotificationIntegrationAutoAzure(t)

		showRequest := sdk.NewShowNotificationIntegrationRequest().
			WithLike(sdk.Like{Pattern: &notificationAutoGoogle.Name})
		returnedIntegrations, err := client.NotificationIntegrations.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *notificationAutoGoogle)
		assert.NotContains(t, returnedIntegrations, *notificationAutoAzure)
	})

	t.Run("describe notification integration: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		_, err := client.NotificationIntegrations.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
