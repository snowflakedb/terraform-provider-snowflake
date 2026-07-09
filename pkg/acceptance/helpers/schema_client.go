package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type SchemaClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSchemaClient(context *TestClientContext, idsGenerator *IdsGenerator) *SchemaClient {
	return &SchemaClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SchemaClient) client() sdk.Schemas {
	return c.context.client.Schemas
}

func (c *SchemaClient) CreateSchema(t *testing.T) (*sdk.Schema, func()) {
	t.Helper()
	return c.CreateSchemaInDatabase(t, c.ids.DatabaseId())
}

// CreateTestSchemaIfNotExists should be used to create the main schema used throughout the acceptance tests.
// It's created only if it does not exist already.
func (c *SchemaClient) CreateTestSchemaIfNotExists(t *testing.T) (*sdk.Schema, func()) {
	t.Helper()
	return c.CreateSchemaWithRequest(t, c.ids.SchemaId(), sdk.NewCreateSchemaRequest(c.ids.SchemaId()).WithIfNotExists(true))
}

func (c *SchemaClient) CreateSchemaInDatabase(t *testing.T, databaseId sdk.AccountObjectIdentifier) (*sdk.Schema, func()) {
	t.Helper()
	return c.CreateSchemaWithIdentifier(t, c.ids.RandomDatabaseObjectIdentifierInDatabase(databaseId))
}

func (c *SchemaClient) CreateSchemaWithName(t *testing.T, name string) (*sdk.Schema, func()) {
	t.Helper()
	return c.CreateSchemaWithIdentifier(t, c.ids.NewDatabaseObjectIdentifier(name))
}

func (c *SchemaClient) CreateSchemaWithIdentifier(t *testing.T, id sdk.DatabaseObjectIdentifier) (*sdk.Schema, func()) {
	t.Helper()
	return c.CreateSchemaWithRequest(t, id, sdk.NewCreateSchemaRequest(id))
}

func (c *SchemaClient) CreateSchemaWithRequest(t *testing.T, id sdk.DatabaseObjectIdentifier, req *sdk.CreateSchemaRequest) (*sdk.Schema, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	schema, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return schema, c.DropSchemaFunc(t, id)
}

func (c *SchemaClient) DropSchemaFunc(t *testing.T, id sdk.DatabaseObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropSchemaRequest(id).WithIfExists(true))
		require.NoError(t, err)
		err = c.context.client.Sessions.UseSchema(ctx, c.ids.SchemaId())
		require.NoError(t, err)
	}
}

func (c *SchemaClient) UseDefaultSchema(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	err := c.context.client.Sessions.UseSchema(ctx, c.ids.SchemaId())
	require.NoError(t, err)
}

func (c *SchemaClient) UpdateDataRetentionTime(t *testing.T, id sdk.DatabaseObjectIdentifier, days int) {
	t.Helper()

	c.Alter(t, sdk.NewAlterSchemaRequest(id).WithSet(sdk.SchemaSetRequest{DataRetentionTimeInDays: sdk.Int(days)}))
}

func (c *SchemaClient) UpdateLogLevel(t *testing.T, id sdk.DatabaseObjectIdentifier, level sdk.LogLevel) {
	t.Helper()

	c.Alter(t, sdk.NewAlterSchemaRequest(id).WithSet(sdk.SchemaSetRequest{LogLevel: &level}))
}

func (c *SchemaClient) UnsetDataRetentionTime(t *testing.T, id sdk.DatabaseObjectIdentifier) {
	t.Helper()

	c.Alter(t, sdk.NewAlterSchemaRequest(id).WithUnset(sdk.SchemaUnsetRequest{DataRetentionTimeInDays: sdk.Bool(true)}))
}

func (c *SchemaClient) UnsetLogLevel(t *testing.T, id sdk.DatabaseObjectIdentifier) {
	t.Helper()

	c.Alter(t, sdk.NewAlterSchemaRequest(id).WithUnset(sdk.SchemaUnsetRequest{LogLevel: sdk.Bool(true)}))
}

func (c *SchemaClient) Show(t *testing.T, id sdk.DatabaseObjectIdentifier) (*sdk.Schema, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *SchemaClient) ShowParameters(t *testing.T, id sdk.DatabaseObjectIdentifier) ([]*sdk.Parameter, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowParameters(ctx, id)
}

func (c *SchemaClient) ShowWithOptions(t *testing.T, req *sdk.ShowSchemaRequest) []sdk.Schema {
	t.Helper()
	ctx := context.Background()

	schemas, err := c.client().Show(ctx, req)
	require.NoError(t, err)
	return schemas
}

func (c *SchemaClient) Alter(t *testing.T, req *sdk.AlterSchemaRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *SchemaClient) AlterDefaultStreamlitNotebookWarehouse(t *testing.T, id sdk.DatabaseObjectIdentifier, warehouse sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	query := fmt.Sprintf(`ALTER SCHEMA %s SET DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE = '%s'`, id.FullyQualifiedName(), warehouse.Name())

	_, err := c.context.client.ExecForTests(ctx, query)
	require.NoError(t, err)
}
