package helpers

import (
	"context"
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
	return c.CreateSchemaWithOpts(t, c.ids.SchemaId(), &sdk.CreateSchemaOptions{IfNotExists: sdk.Bool(true)})
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
	return c.CreateSchemaWithOpts(t, id, nil)
}

func (c *SchemaClient) CreateSchemaWithOpts(t *testing.T, id sdk.DatabaseObjectIdentifier, opts *sdk.CreateSchemaOptions) (*sdk.Schema, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)
	schema, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return schema, c.DropSchemaFunc(t, id)
}

func (c *SchemaClient) DropSchemaFunc(t *testing.T, id sdk.DatabaseObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropSchemaOptions{IfExists: sdk.Bool(true)})
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

	c.Alter(t, id, &sdk.AlterSchemaOptions{
		Set: &sdk.SchemaSet{
			DataRetentionTimeInDays: sdk.Int(days),
		},
	})
}

func (c *SchemaClient) UnsetDataRetentionTime(t *testing.T, id sdk.DatabaseObjectIdentifier) {
	t.Helper()

	c.Alter(t, id, &sdk.AlterSchemaOptions{
		Unset: &sdk.SchemaUnset{
			DataRetentionTimeInDays: sdk.Bool(true),
		},
	})
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

func (c *SchemaClient) ShowWithOptions(t *testing.T, opts *sdk.ShowSchemaOptions) []sdk.Schema {
	t.Helper()
	ctx := context.Background()

	schemas, err := c.client().Show(ctx, opts)
	require.NoError(t, err)
	return schemas
}

func (c *SchemaClient) Alter(t *testing.T, id sdk.DatabaseObjectIdentifier, opts *sdk.AlterSchemaOptions) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, opts)
	require.NoError(t, err)
}
