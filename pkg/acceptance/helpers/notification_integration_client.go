package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1017580]: replace with real value
const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"

const VerifiedEmail = "admin@example.com"

type NotificationIntegrationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewNotificationIntegrationClient(context *TestClientContext, idsGenerator *IdsGenerator) *NotificationIntegrationClient {
	return &NotificationIntegrationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *NotificationIntegrationClient) client() sdk.NotificationIntegrations {
	return c.context.client.NotificationIntegrations
}

func (c *NotificationIntegrationClient) CreateWithGcpPubSub(t *testing.T) (*sdk.NotificationIntegration, func()) {
	t.Helper()
	return c.CreateWithRequest(
		t, sdk.NewCreateNotificationIntegrationRequest(c.ids.RandomAccountObjectIdentifier(), true).
			WithAutomatedDataLoadsParams(
				*sdk.NewAutomatedDataLoadsParamsRequest().
					WithGoogleAutoParams(*sdk.NewGoogleAutoParamsRequest(gcpPubsubSubscriptionName)),
			),
	)
}

func (c *NotificationIntegrationClient) CreateWebhook(t *testing.T, webhookUrl string) (*sdk.NotificationIntegration, func()) {
	t.Helper()
	return c.CreateWithRequest(
		t, sdk.NewCreateNotificationIntegrationRequest(c.ids.RandomAccountObjectIdentifier(), true).
			WithWebhookParams(*sdk.NewWebhookParamsRequest(webhookUrl)),
	)
}

func (c *NotificationIntegrationClient) Create(t *testing.T) (*sdk.NotificationIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()

	request := sdk.NewCreateNotificationIntegrationRequest(id, true).
		WithEmailParams(*sdk.NewEmailParamsRequest().WithAllowedRecipients([]sdk.NotificationIntegrationAllowedRecipient{{Email: VerifiedEmail}}))

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	integration, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return integration, c.DropFunc(t, id)
}

func (c *NotificationIntegrationClient) CreateWithRequest(t *testing.T, request *sdk.CreateNotificationIntegrationRequest) (*sdk.NotificationIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	networkRule, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return networkRule, c.DropFunc(t, request.GetName())
}

func (c *NotificationIntegrationClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// assert instead of require, so that a failed drop does not abort the remaining cleanups
		assert.NoError(t, c.client().DropSafely(ctx, id))
	}
}
