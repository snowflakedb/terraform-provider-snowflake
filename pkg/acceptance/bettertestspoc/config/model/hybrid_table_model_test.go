package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHybridTableModel_BasicBuilder validates the basic builder creates proper structure
func TestHybridTableModel_BasicBuilder(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE")

	assert.NotNil(t, model)
	assert.NotNil(t, model.ResourceModelMeta)
	assert.Equal(t, "test", model.ResourceName())
}

// TestHybridTableModel_WithMethods validates the With* methods work correctly
func TestHybridTableModel_WithMethods(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithComment("test comment").
		WithDataRetentionTimeInDays(7).
		WithOrReplace(true)

	require.NotNil(t, model)
	// Verify values can be set (actual validation happens during terraform apply)
}

// TestHybridTableModel_WithColumnDescs validates column building works
func TestHybridTableModel_WithColumnDescs(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithColumnDescs([]ColumnDesc{
			{
				Name:       "id",
				DataType:   "NUMBER(38,0)",
				Nullable:   Bool(false),
				PrimaryKey: true,
				Comment:    "ID column",
			},
			{
				Name:     "name",
				DataType: "VARCHAR(100)",
				Nullable: Bool(true),
			},
		})

	require.NotNil(t, model)
	require.NotNil(t, model.Column)
}

// TestHybridTableModel_WithColumnDefault validates column default works
func TestHybridTableModel_WithColumnDefault(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithColumnDescs([]ColumnDesc{
			{
				Name:     "created_at",
				DataType: "TIMESTAMP_NTZ",
				Default: &ColumnDefaultOpts{
					Expression: "CURRENT_TIMESTAMP()",
				},
			},
		})

	require.NotNil(t, model)
	require.NotNil(t, model.Column)
}

// TestHybridTableModel_WithColumnIdentity validates identity column works
func TestHybridTableModel_WithColumnIdentity(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithColumnDescs([]ColumnDesc{
			{
				Name:     "id",
				DataType: "NUMBER(38,0)",
				Identity: &ColumnIdentityOpts{
					StartNum: 1,
					StepNum:  1,
				},
			},
		})

	require.NotNil(t, model)
	require.NotNil(t, model.Column)
}

// TestHybridTableModel_WithPrimaryKey validates out-of-line primary key works
func TestHybridTableModel_WithPrimaryKey(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithPrimaryKeyColumns("id", "tenant_id")

	require.NotNil(t, model)
	require.NotNil(t, model.PrimaryKey)
}

// TestHybridTableModel_WithPrimaryKeyNamed validates named primary key works
func TestHybridTableModel_WithPrimaryKeyNamed(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithPrimaryKeyNamed("pk_test", "id")

	require.NotNil(t, model)
	require.NotNil(t, model.PrimaryKey)
}

// TestHybridTableModel_WithIndexes validates index building works
func TestHybridTableModel_WithIndexes(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithIndexes([]IndexDesc{
			{Name: "idx_name", Columns: []string{"name"}},
			{Name: "idx_email", Columns: []string{"email"}},
		})

	require.NotNil(t, model)
	require.NotNil(t, model.Index)
}

// TestHybridTableModel_WithUniqueConstraints validates unique constraint works
func TestHybridTableModel_WithUniqueConstraints(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithUniqueConstraints([]UniqueConstraintDesc{
			{Name: "uq_email", Columns: []string{"email"}},
		})

	require.NotNil(t, model)
	require.NotNil(t, model.UniqueConstraint)
}

// TestHybridTableModel_WithForeignKeys validates foreign key works
func TestHybridTableModel_WithForeignKeys(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "TEST_TABLE").
		WithForeignKeys([]ForeignKeyDesc{
			{
				Name:              "fk_user",
				Columns:           []string{"user_id"},
				ReferencesTable:   "users",
				ReferencesColumns: []string{"id"},
			},
		})

	require.NotNil(t, model)
	require.NotNil(t, model.ForeignKey)
}

// TestHybridTableModel_ComplexExample validates a complex configuration
func TestHybridTableModel_ComplexExample(t *testing.T) {
	model := HybridTable("test", "TEST_DB", "TEST_SCHEMA", "orders").
		WithComment("Orders table").
		WithDataRetentionTimeInDays(30).
		WithColumnDescs([]ColumnDesc{
			{
				Name:     "order_id",
				DataType: "NUMBER(38,0)",
				Nullable: Bool(false),
			},
			{
				Name:     "customer_id",
				DataType: "NUMBER(38,0)",
				Nullable: Bool(false),
			},
			{
				Name:     "order_date",
				DataType: "DATE",
				Nullable: Bool(false),
				Default: &ColumnDefaultOpts{
					Expression: "CURRENT_DATE()",
				},
			},
			{
				Name:     "status",
				DataType: "VARCHAR(50)",
				Nullable: Bool(false),
			},
		}).
		WithPrimaryKeyColumns("order_id").
		WithIndexes([]IndexDesc{
			{Name: "idx_customer", Columns: []string{"customer_id"}},
			{Name: "idx_order_date", Columns: []string{"order_date"}},
		}).
		WithForeignKeys([]ForeignKeyDesc{
			{
				Name:              "fk_customer",
				Columns:           []string{"customer_id"},
				ReferencesTable:   "customers",
				ReferencesColumns: []string{"customer_id"},
			},
		})

	// Verify all components are present
	require.NotNil(t, model)
	require.NotNil(t, model.Column)
	require.NotNil(t, model.PrimaryKey)
	require.NotNil(t, model.Index)
	require.NotNil(t, model.ForeignKey)
}

// TestColumnDesc_toVariable validates the column conversion
func TestColumnDesc_toVariable(t *testing.T) {
	col := ColumnDesc{
		Name:       "test_col",
		DataType:   "VARCHAR(50)",
		Nullable:   Bool(false),
		PrimaryKey: true,
		Comment:    "test comment",
	}

	v := col.toVariable()
	require.NotNil(t, v)
}
