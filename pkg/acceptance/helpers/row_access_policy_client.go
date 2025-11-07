package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/require"
)

type RowAccessPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewRowAccessPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *RowAccessPolicyClient {
	return &RowAccessPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *RowAccessPolicyClient) client() sdk.RowAccessPolicies {
	return c.context.client.RowAccessPolicies
}

func (c *RowAccessPolicyClient) CreateRowAccessPolicy(t *testing.T) (*sdk.RowAccessPolicy, func()) {
	t.Helper()
	return c.CreateRowAccessPolicyWithDataType(t, testdatatypes.DataTypeNumber)
}

func (c *RowAccessPolicyClient) CreateRowAccessPolicyWithDataType(t *testing.T, datatype datatypes.DataType) (*sdk.RowAccessPolicy, func()) {
	t.Helper()

	arg := sdk.NewCreateRowAccessPolicyArgsRequest("A", datatype)
	return c.CreateRowAccessPolicyWithArguments(t, []sdk.CreateRowAccessPolicyArgsRequest{*arg})
}

func (c *RowAccessPolicyClient) CreateRowAccessPolicyWithArguments(t *testing.T, args []sdk.CreateRowAccessPolicyArgsRequest) (*sdk.RowAccessPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	body := "true"
	createRequest := sdk.NewCreateRowAccessPolicyRequest(id, args, body)

	err := c.client().Create(ctx, createRequest)
	require.NoError(t, err)

	rowAccessPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return rowAccessPolicy, c.DropRowAccessPolicyFunc(t, id)
}

func (c *RowAccessPolicyClient) CreateRowAccessPolicyWithRequest(t *testing.T, req sdk.CreateRowAccessPolicyRequest) (*sdk.RowAccessPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, &req)
	require.NoError(t, err)

	rowAccessPolicy, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return rowAccessPolicy, c.DropRowAccessPolicyFunc(t, req.GetName())
}

func (c *RowAccessPolicyClient) Alter(t *testing.T, req sdk.AlterRowAccessPolicyRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, &req)
	require.NoError(t, err)
}

func (c *RowAccessPolicyClient) DropRowAccessPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *RowAccessPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.RowAccessPolicy, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
