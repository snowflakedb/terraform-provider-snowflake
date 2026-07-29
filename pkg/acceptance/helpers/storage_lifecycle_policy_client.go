package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StorageLifecyclePolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewStorageLifecyclePolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *StorageLifecyclePolicyClient {
	return &StorageLifecyclePolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *StorageLifecyclePolicyClient) client() sdk.StorageLifecyclePolicies {
	return c.context.client.StorageLifecyclePolicies
}

// DefaultArgs is a minimal valid signature for a storage lifecycle policy. The body below references it.
func (c *StorageLifecyclePolicyClient) DefaultArgs() []sdk.CreateStorageLifecyclePolicyArgsRequest {
	return []sdk.CreateStorageLifecyclePolicyArgsRequest{
		*sdk.NewCreateStorageLifecyclePolicyArgsRequest("VAL", testdatatypes.DataTypeVarchar_200),
	}
}

// DefaultBody is a minimal valid boolean expression that references the default signature argument.
func (c *StorageLifecyclePolicyClient) DefaultBody() string {
	return "LENGTH(VAL) > 0"
}

func (c *StorageLifecyclePolicyClient) Create(t *testing.T) (*sdk.StorageLifecyclePolicy, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	cleanup := c.CreateWithId(t, id)

	policy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return policy, cleanup
}

func (c *StorageLifecyclePolicyClient) CreateWithId(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	return c.CreateWithRequest(t, id, sdk.NewCreateStorageLifecyclePolicyRequest(id, c.DefaultArgs(), c.DefaultBody()))
}

func (c *StorageLifecyclePolicyClient) CreateWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, req *sdk.CreateStorageLifecyclePolicyRequest) func() {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	return c.DropFunc(t, id)
}

func (c *StorageLifecyclePolicyClient) Alter(t *testing.T, req *sdk.AlterStorageLifecyclePolicyRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *StorageLifecyclePolicyClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().DropSafely(ctx, id)
		assert.NoError(t, err)
	}
}

func (c *StorageLifecyclePolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.StorageLifecyclePolicy, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *StorageLifecyclePolicyClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.StorageLifecyclePolicyDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, id)
}
