package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
