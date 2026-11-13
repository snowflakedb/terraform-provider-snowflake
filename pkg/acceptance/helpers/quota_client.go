package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type QuotaClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewQuotaClient(context *TestClientContext, idsGenerator *IdsGenerator) *QuotaClient {
	return &QuotaClient{context: context, ids: idsGenerator}
}

func (c *QuotaClient) client() sdk.Quotas {
	return c.context.client.Quotas
}

// Create returns the id only (like the budget client), so that tests not about SHOW don't depend on it.
func (c *QuotaClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	return c.CreateWithId(t, c.ids.RandomSchemaObjectIdentifier())
}

func (c *QuotaClient) CreateWithId(t *testing.T, id sdk.SchemaObjectIdentifier) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, sdk.NewCreateQuotaRequest(id))
	require.NoError(t, err)

	return id, c.DropFunc(t, id)
}

func (c *QuotaClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().DropSafely(ctx, id)
		require.NoError(t, err)
	}
}

func (c *QuotaClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Quota, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *QuotaClient) DescribeDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.QuotaDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().DescribeDetails(ctx, id)
}

func (c *QuotaClient) AddSharedResource(t *testing.T, id sdk.SchemaObjectIdentifier, domain sdk.QuotaSharedResourceDomain, instance string) {
	t.Helper()
	ctx := context.Background()

	args, err := sdk.NewQuotaAddSharedResourceArgsRequestFromInstance(domain, instance)
	require.NoError(t, err)
	_, err = c.client().AddSharedResource(ctx, sdk.NewAddSharedResourceQuotaRequest(id, *args))
	require.NoError(t, err)
}

func (c *QuotaClient) SetPerUserLimit(t *testing.T, id sdk.SchemaObjectIdentifier, limit int, cycle sdk.QuotaCycle) {
	t.Helper()
	ctx := context.Background()

	args := sdk.NewQuotaSetPerUserLimitArgsRequest(limit)
	if cycle != sdk.QuotaCycleMonthly {
		args.WithCycle(cycle)
	}
	_, err := c.client().SetPerUserLimit(ctx, sdk.NewSetPerUserLimitQuotaRequest(id, *args))
	require.NoError(t, err)
}

func (c *QuotaClient) IncludeUsers(t *testing.T, id sdk.SchemaObjectIdentifier, users ...string) {
	t.Helper()
	ctx := context.Background()

	_, err := c.client().IncludeUsers(ctx, sdk.NewIncludeUsersQuotaRequest(id, *sdk.NewQuotaIncludeUsersArgsRequest(users)))
	require.NoError(t, err)
}
