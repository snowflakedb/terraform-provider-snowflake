package helpers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TableClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewTableClient(context *TestClientContext, idsGenerator *IdsGenerator) *TableClient {
	return &TableClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *TableClient) client() sdk.Tables {
	return c.context.client.Tables
}

func (c *TableClient) Create(t *testing.T) (*sdk.Table, func()) {
	t.Helper()
	return c.CreateInSchema(t, c.ids.SchemaId())
}

func (c *TableClient) CreateWithName(t *testing.T, name string) (*sdk.Table, func()) {
	t.Helper()

	columns := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
	}
	return c.CreateWithRequest(t, sdk.NewCreateTableRequest(c.ids.NewSchemaObjectIdentifier(name), columns))
}

func (c *TableClient) CreateInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Table, func()) {
	t.Helper()

	columns := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
	}
	return c.CreateWithRequest(t, sdk.NewCreateTableRequest(c.ids.RandomSchemaObjectIdentifierInSchema(schemaId), columns))
}

func (c *TableClient) CreateWithColumns(t *testing.T, columns []sdk.TableColumnRequest) (*sdk.Table, func()) {
	t.Helper()

	return c.CreateWithRequest(t, sdk.NewCreateTableRequest(c.ids.RandomSchemaObjectIdentifier(), columns))
}

func (c *TableClient) CreateWithPredefinedColumns(t *testing.T) (*sdk.Table, func()) {
	t.Helper()

	columns := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", "NUMBER"),
		*sdk.NewTableColumnRequest("some_text_column", "VARCHAR"),
		*sdk.NewTableColumnRequest("some_other_text_column", "VARCHAR"),
	}

	return c.CreateWithRequest(t, sdk.NewCreateTableRequest(c.ids.RandomSchemaObjectIdentifier(), columns))
}

func (c *TableClient) CreateWithChangeTrackingInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Table, func()) {
	t.Helper()

	columns := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", "NUMBER"),
	}

	return c.CreateWithRequest(t, sdk.NewCreateTableRequest(c.ids.RandomSchemaObjectIdentifierInSchema(schemaId), columns).WithChangeTracking(sdk.Pointer(true)))
}

func (c *TableClient) CreateWithChangeTracking(t *testing.T) (*sdk.Table, func()) {
	t.Helper()

	columns := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", "NUMBER"),
	}

	return c.CreateWithRequest(t, sdk.NewCreateTableRequest(c.ids.RandomSchemaObjectIdentifier(), columns).WithChangeTracking(sdk.Pointer(true)))
}

func (c *TableClient) CreateWithRequest(t *testing.T, req *sdk.CreateTableRequest) (*sdk.Table, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	table, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return table, c.DropFunc(t, req.GetName())
}

func (c *TableClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// to prevent error when schema was removed before the table
		_, err := c.context.client.Schemas.ShowByID(ctx, id.SchemaId())
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			return
		}

		dropErr := c.client().Drop(ctx, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, dropErr)
	}
}

func (c *TableClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Table, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *TableClient) SetDataRetentionTime(t *testing.T, id sdk.SchemaObjectIdentifier, days int) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterTableRequest(id).WithSet(sdk.NewTableSetRequest().WithDataRetentionTimeInDays(sdk.Int(days))))
	require.NoError(t, err)
}

// GetTableColumnsFor is based on https://docs.snowflake.com/en/sql-reference/info-schema/columns.
// TODO: extract getting table columns as resource (like getting tag in system functions)
func (c *TableClient) GetTableColumnsFor(t *testing.T, tableId sdk.SchemaObjectIdentifier) []InformationSchemaColumns {
	t.Helper()
	ctx := context.Background()

	var columns []InformationSchemaColumns
	query := fmt.Sprintf("SELECT * FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s' ORDER BY ordinal_position", tableId.SchemaName(), tableId.Name())
	err := c.context.client.QueryForTests(ctx, &columns, query)
	require.NoError(t, err)

	return columns
}

func (c *TableClient) InsertInt(t *testing.T, tableId sdk.SchemaObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf("INSERT INTO %s VALUES(1);", tableId.FullyQualifiedName()))
	require.NoError(t, err)
}

type InformationSchemaColumns struct {
	TableCatalog           string         `db:"TABLE_CATALOG"`
	TableSchema            string         `db:"TABLE_SCHEMA"`
	TableName              string         `db:"TABLE_NAME"`
	ColumnName             string         `db:"COLUMN_NAME"`
	OrdinalPosition        string         `db:"ORDINAL_POSITION"`
	ColumnDefault          sql.NullString `db:"COLUMN_DEFAULT"`
	IsNullable             string         `db:"IS_NULLABLE"`
	DataType               string         `db:"DATA_TYPE"`
	CharacterMaximumLength sql.NullString `db:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   sql.NullString `db:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       sql.NullString `db:"NUMERIC_PRECISION"`
	NumericPrecisionRadix  sql.NullString `db:"NUMERIC_PRECISION_RADIX"`
	NumericScale           sql.NullString `db:"NUMERIC_SCALE"`
	DatetimePrecision      sql.NullString `db:"DATETIME_PRECISION"`
	IntervalType           sql.NullString `db:"INTERVAL_TYPE"`
	IntervalPrecision      sql.NullString `db:"INTERVAL_PRECISION"`
	CharacterSetCatalog    sql.NullString `db:"CHARACTER_SET_CATALOG"`
	CharacterSetSchema     sql.NullString `db:"CHARACTER_SET_SCHEMA"`
	CharacterSetName       sql.NullString `db:"CHARACTER_SET_NAME"`
	CollationCatalog       sql.NullString `db:"COLLATION_CATALOG"`
	CollationSchema        sql.NullString `db:"COLLATION_SCHEMA"`
	CollationName          sql.NullString `db:"COLLATION_NAME"`
	DomainCatalog          sql.NullString `db:"DOMAIN_CATALOG"`
	DomainSchema           sql.NullString `db:"DOMAIN_SCHEMA"`
	DomainName             sql.NullString `db:"DOMAIN_NAME"`
	UdtCatalog             sql.NullString `db:"UDT_CATALOG"`
	UdtSchema              sql.NullString `db:"UDT_SCHEMA"`
	UdtName                sql.NullString `db:"UDT_NAME"`
	ScopeCatalog           sql.NullString `db:"SCOPE_CATALOG"`
	ScopeSchema            sql.NullString `db:"SCOPE_SCHEMA"`
	ScopeName              sql.NullString `db:"SCOPE_NAME"`
	MaximumCardinality     sql.NullString `db:"MAXIMUM_CARDINALITY"`
	DtdIdentifier          sql.NullString `db:"DTD_IDENTIFIER"`
	IsSelfReferencing      string         `db:"IS_SELF_REFERENCING"`
	IsIdentity             string         `db:"IS_IDENTITY"`
	IdentityGeneration     sql.NullString `db:"IDENTITY_GENERATION"`
	IdentityStart          sql.NullString `db:"IDENTITY_START"`
	IdentityIncrement      sql.NullString `db:"IDENTITY_INCREMENT"`
	IdentityMaximum        sql.NullString `db:"IDENTITY_MAXIMUM"`
	IdentityMinimum        sql.NullString `db:"IDENTITY_MINIMUM"`
	IdentityCycle          sql.NullString `db:"IDENTITY_CYCLE"`
	IdentityOrdered        sql.NullString `db:"IDENTITY_ORDERED"`
	Comment                sql.NullString `db:"COMMENT"`
	SchemaEvolutionRecord  sql.NullString `db:"SCHEMA_EVOLUTION_RECORD"`
}
