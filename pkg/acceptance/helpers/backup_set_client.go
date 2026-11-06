package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type BackupSetClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewBackupSetClient(context *TestClientContext, idsGenerator *IdsGenerator) *BackupSetClient {
	return &BackupSetClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *BackupSetClient) client() sdk.BackupSets {
	return c.context.client.BackupSets
}

func (c *BackupSetClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.BackupSet, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *BackupSetClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.BackupSetDetails, error) {
	t.Helper()
	return c.client().Describe(context.Background(), id)
}

func (c *BackupSetClient) CreateWithOptions(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateBackupSetRequest) (*sdk.BackupSet, func()) {
	t.Helper()
	ctx := context.Background()
	require.NoError(t, c.client().Create(ctx, request))
	backupSet, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return backupSet, c.DropFunc(t, id)
}

func (c *BackupSetClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	return func() {
		require.NoError(t, c.client().Drop(context.Background(), sdk.NewDropBackupSetRequest(id).WithIfExists(true)))
	}
}

func (c *BackupSetClient) Alter(t *testing.T, req *sdk.AlterBackupSetRequest) {
	t.Helper()
	require.NoError(t, c.client().Alter(context.Background(), req))
}

// BackupID reads the first backup_id from SHOW BACKUPS IN BACKUP SET.
// MODIFY and DELETE BACKUP take that id, and there is no backups SDK.
func (c *BackupSetClient) BackupID(t *testing.T, id sdk.SchemaObjectIdentifier) string {
	t.Helper()
	rows, err := c.context.client.QueryUnsafe(context.Background(), fmt.Sprintf("SHOW BACKUPS IN BACKUP SET %s", id.FullyQualifiedName()))
	require.NoError(t, err)
	require.NotEmpty(t, rows)
	value, ok := (*rows[0]["backup_id"]).(string)
	require.True(t, ok)
	require.NotEmpty(t, value)
	return value
}
