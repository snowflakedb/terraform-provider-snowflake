package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ContextClient struct {
	context *TestClientContext
}

func NewContextClient(context *TestClientContext) *ContextClient {
	return &ContextClient{
		context: context,
	}
}

func (c *ContextClient) client() sdk.ContextFunctions {
	return c.context.client.ContextFunctions
}

func (c *ContextClient) CurrentAccount(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	currentAccount, err := c.client().CurrentAccount(ctx)
	require.NoError(t, err)

	return currentAccount
}

func (c *ContextClient) CurrentAccountId(t *testing.T) sdk.AccountIdentifier {
	t.Helper()
	ctx := context.Background()

	currentSessionDetails, err := c.client().CurrentSessionDetails(ctx)
	require.NoError(t, err)

	return sdk.NewAccountIdentifier(currentSessionDetails.OrganizationName, currentSessionDetails.AccountName)
}

func (c *ContextClient) CurrentAccountName(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	currentAccount, err := c.client().CurrentAccountName(ctx)
	require.NoError(t, err)

	return currentAccount
}

func (c *ContextClient) CurrentRole(t *testing.T) sdk.AccountObjectIdentifier {
	t.Helper()
	ctx := context.Background()

	currentRole, err := c.client().CurrentRole(ctx)
	require.NoError(t, err)

	return currentRole
}

func (c *ContextClient) CurrentRegion(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	currentRegion, err := c.client().CurrentRegion(ctx)
	require.NoError(t, err)

	return currentRegion
}

func (c *ContextClient) CurrentUser(t *testing.T) sdk.AccountObjectIdentifier {
	t.Helper()
	ctx := context.Background()

	currentUser, err := c.client().CurrentUser(ctx)
	require.NoError(t, err)

	return currentUser
}

func (c *ContextClient) CurrentAccountIdentifier(t *testing.T) sdk.AccountIdentifier {
	t.Helper()

	details, err := c.client().CurrentSessionDetails(context.Background())
	require.NoError(t, err)

	return sdk.NewAccountIdentifier(details.OrganizationName, details.AccountName)
}

func (c *ContextClient) IsRoleInSession(t *testing.T, id sdk.AccountObjectIdentifier) bool {
	t.Helper()
	ctx := context.Background()

	isInSession, err := c.client().IsRoleInSession(ctx, id)
	require.NoError(t, err)

	return isInSession
}

// ACSURL returns Snowflake Assertion Consumer Service URL
func (c *ContextClient) ACSURL(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("https://%s.snowflakecomputing.com/fed/login", c.CurrentAccount(t))
}

// IssuerURL returns a URL containing the EntityID / Issuer for the Snowflake service provider
func (c *ContextClient) IssuerURL(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("https://%s.snowflakecomputing.com", c.CurrentAccount(t))
}

func (c *ContextClient) LastQueryId(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	id, err := c.client().LastQueryId(ctx)
	require.NoError(t, err)

	return id
}
