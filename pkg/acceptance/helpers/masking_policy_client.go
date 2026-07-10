package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MaskingPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewMaskingPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *MaskingPolicyClient {
	return &MaskingPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *MaskingPolicyClient) client() sdk.MaskingPolicies {
	return c.context.client.MaskingPolicies
}

func (c *MaskingPolicyClient) CreateMaskingPolicy(t *testing.T) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	signature := []sdk.CreateMaskingPolicySignatureRequest{
		*sdk.NewCreateMaskingPolicySignatureRequest(c.ids.Alpha(), testdatatypes.DataTypeVarchar),
		*sdk.NewCreateMaskingPolicySignatureRequest(c.ids.Alpha(), testdatatypes.DataTypeVarchar),
	}
	expression := "REPLACE('X', 1, 2)"
	return c.CreateMaskingPolicyWithRequest(t, signature, testdatatypes.DataTypeVarchar, expression)
}

func (c *MaskingPolicyClient) CreateMaskingPolicyIdentity(t *testing.T, columnType datatypes.DataType) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	signature := []sdk.CreateMaskingPolicySignatureRequest{
		*sdk.NewCreateMaskingPolicySignatureRequest("a", columnType),
	}
	expression := "a"
	return c.CreateMaskingPolicyWithRequest(t, signature, columnType, expression)
}

func (c *MaskingPolicyClient) CreateMaskingPolicyWithRequest(t *testing.T, signature []sdk.CreateMaskingPolicySignatureRequest, returns datatypes.DataType, expression string) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	ctx := context.Background()
	id := c.ids.RandomSchemaObjectIdentifier()

	err := c.client().Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, returns, expression))
	require.NoError(t, err)

	maskingPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return maskingPolicy, c.DropMaskingPolicyFunc(t, id)
}

func (c *MaskingPolicyClient) CreateOrReplaceMaskingPolicyWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, signature []sdk.CreateMaskingPolicySignatureRequest, returns datatypes.DataType, expression string) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, sdk.NewCreateMaskingPolicyRequest(id, signature, returns, expression).WithOrReplace(true))
	require.NoError(t, err)

	maskingPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return maskingPolicy, c.DropMaskingPolicyFunc(t, id)
}

func (c *MaskingPolicyClient) Alter(t *testing.T, req *sdk.AlterMaskingPolicyRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *MaskingPolicyClient) DropMaskingPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropMaskingPolicyRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	}
}

func (c *MaskingPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.MaskingPolicy, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
