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
	signature := []sdk.TableColumnSignature{
		{
			Name: c.ids.Alpha(),
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: c.ids.Alpha(),
			Type: testdatatypes.DataTypeVarchar,
		},
	}
	expression := "REPLACE('X', 1, 2)"
	return c.CreateMaskingPolicyWithOptions(t, signature, testdatatypes.DataTypeVarchar, expression, &sdk.CreateMaskingPolicyOptions{})
}

func (c *MaskingPolicyClient) CreateMaskingPolicyIdentity(t *testing.T, columnType datatypes.DataType) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	name := "a"
	signature := []sdk.TableColumnSignature{
		{
			Name: name,
			Type: columnType,
		},
	}
	expression := "a"
	return c.CreateMaskingPolicyWithOptions(t, signature, columnType, expression, &sdk.CreateMaskingPolicyOptions{})
}

func (c *MaskingPolicyClient) CreateMaskingPolicyWithOptions(t *testing.T, signature []sdk.TableColumnSignature, returns datatypes.DataType, expression string, options *sdk.CreateMaskingPolicyOptions) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	ctx := context.Background()
	id := c.ids.RandomSchemaObjectIdentifier()

	err := c.client().Create(ctx, id, signature, returns, expression, options)
	require.NoError(t, err)

	maskingPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return maskingPolicy, c.DropMaskingPolicyFunc(t, id)
}

func (c *MaskingPolicyClient) CreateOrReplaceMaskingPolicyWithOptions(t *testing.T, id sdk.SchemaObjectIdentifier, signature []sdk.TableColumnSignature, returns datatypes.DataType, expression string, options *sdk.CreateMaskingPolicyOptions) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	options.OrReplace = sdk.Pointer(true)

	err := c.client().Create(ctx, id, signature, returns, expression, options)
	require.NoError(t, err)

	maskingPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return maskingPolicy, c.DropMaskingPolicyFunc(t, id)
}

func (c *MaskingPolicyClient) Alter(t *testing.T, id sdk.SchemaObjectIdentifier, req *sdk.AlterMaskingPolicyOptions) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, req)
	require.NoError(t, err)
}

func (c *MaskingPolicyClient) DropMaskingPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Bool(true)})
		assert.NoError(t, err)
	}
}

func (c *MaskingPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.MaskingPolicy, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
