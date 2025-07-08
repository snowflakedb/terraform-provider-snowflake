package helpers

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
	"testing"
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
