package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type OrganizationAccountClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewOrganizationAccountClient(context *TestClientContext, idsGenerator *IdsGenerator) *OrganizationAccountClient {
	return &OrganizationAccountClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *OrganizationAccountClient) client() sdk.OrganizationAccounts {
	return c.context.client.OrganizationAccounts
}

func (c *OrganizationAccountClient) Alter(t *testing.T, req *sdk.AlterOrganizationAccountRequest) {
	t.Helper()
	err := c.client().Alter(context.Background(), req)
	require.NoError(t, err)
}

func (c *OrganizationAccountClient) ShowCurrent(t *testing.T) sdk.OrganizationAccount {
	t.Helper()
	organizationAccount, err := c.client().Show(context.Background(), sdk.NewShowOrganizationAccountRequest())
	require.NoError(t, err)
	require.Len(t, organizationAccount, 1)
	return organizationAccount[0]
}

func (c *OrganizationAccountClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.OrganizationAccount, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}
