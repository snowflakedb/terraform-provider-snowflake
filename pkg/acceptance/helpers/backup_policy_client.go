package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type BackupPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewBackupPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *BackupPolicyClient {
	return &BackupPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *BackupPolicyClient) client() sdk.BackupPolicies {
	return c.context.client.BackupPolicies
}

func (c *BackupPolicyClient) CreateWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, req *sdk.CreateBackupPolicyRequest) (*sdk.BackupPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	require.NoError(t, c.client().Create(ctx, req))
	backupPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return backupPolicy, c.DropFunc(t, id)
}

func (c *BackupPolicyClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		require.NoError(t, c.client().DropSafely(ctx, id))
	}
}

func (c *BackupPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.BackupPolicy, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *BackupPolicyClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.BackupPolicyDetails, error) {
	t.Helper()
	return c.client().Describe(context.Background(), id)
}

func (c *BackupPolicyClient) Alter(t *testing.T, request *sdk.AlterBackupPolicyRequest) {
	t.Helper()
	require.NoError(t, c.client().Alter(context.Background(), request))
}
