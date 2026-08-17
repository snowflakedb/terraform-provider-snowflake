package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func (c *ContextClient) CurrentOrganizationName(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	organizationName, err := c.client().CurrentOrganizationName(ctx)
	require.NoError(t, err)

	return organizationName
}

func (c *ContextClient) IsRoleInSession(t *testing.T, id sdk.AccountObjectIdentifier) bool {
	t.Helper()
	ctx := context.Background()

	isInSession, err := c.client().IsRoleInSession(ctx, id)
	require.NoError(t, err)

	return isInSession
}

// ACSURL returns Snowflake Assertion Consumer Service URL.
// Uses the org-based URL format which is valid for all deployments (prod and non-prod).
func (c *ContextClient) ACSURL(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s/fed/login", c.AccountURL(t))
}

// IssuerURL returns a URL containing the EntityID / Issuer for the Snowflake service provider.
// Uses the org-based URL format which is valid for all deployments (prod and non-prod).
func (c *ContextClient) IssuerURL(t *testing.T) string {
	t.Helper()
	return c.AccountURL(t)
}

// AccountURL returns the org-based account URL in the standard DNS-safe format:
// https://<org>-<account_name>.snowflakecomputing.com (lowercase, underscores replaced with hyphens).
func (c *ContextClient) AccountURL(t *testing.T) string {
	t.Helper()
	org := strings.ToLower(c.CurrentOrganizationName(t))
	accountName := strings.ToLower(strings.ReplaceAll(c.CurrentAccountName(t), "_", "-"))
	return fmt.Sprintf("https://%s-%s.snowflakecomputing.com", org, accountName)
}

func (c *ContextClient) LastQueryId(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	id, err := c.client().LastQueryId(ctx)
	require.NoError(t, err)

	return id
}

func (c *ContextClient) DefaultConsumptionBillingEntity(t *testing.T) sdk.AccountObjectIdentifier {
	t.Helper()
	orgName, err := c.context.client.ContextFunctions.CurrentOrganizationName(context.Background())
	require.NoError(t, err)
	return sdk.NewAccountObjectIdentifier(fmt.Sprintf("%s_DefaultBE", orgName))
}
