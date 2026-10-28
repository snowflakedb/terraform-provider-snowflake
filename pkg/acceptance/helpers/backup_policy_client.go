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
