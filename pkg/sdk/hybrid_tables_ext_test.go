package sdk

import (
	"testing"
)

func TestHybridTables_Create_AllOptions(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("create with columns and inline primary key", func(t *testing.T) {
		opts := &CreateHybridTableOptions{
			name: id,
			ColumnsAndConstraints: HybridTableColumnsConstraintsAndIndexes{
				Columns: []HybridTableColumn{
					{
						Name: "id",
						Type: DataType("NUMBER(38,0)"),
						InlineConstraint: &ColumnInlineConstraint{
							Type: ColumnConstraintTypePrimaryKey,
						},
					},
					{
						Name:    "name",
						Type:    DataType("VARCHAR(100)"),
						NotNull: Bool(true),
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE HYBRID TABLE %s ("id" NUMBER(38,0) PRIMARY KEY, "name" VARCHAR(100) NOT NULL)`, id.FullyQualifiedName())
	})

	t.Run("create with out-of-line constraint and index", func(t *testing.T) {
		opts := &CreateHybridTableOptions{
			name: id,
			ColumnsAndConstraints: HybridTableColumnsConstraintsAndIndexes{
				Columns: []HybridTableColumn{
					{
						Name: "id",
						Type: DataType("NUMBER(38,0)"),
						InlineConstraint: &ColumnInlineConstraint{
							Type: ColumnConstraintTypePrimaryKey,
						},
					},
					{
						Name: "email",
						Type: DataType("VARCHAR(200)"),
					},
				},
				OutOfLineConstraint: []HybridTableOutOfLineConstraint{
					{
						Type:    ColumnConstraintTypeUnique,
						Columns: []string{"email"},
					},
				},
				OutOfLineIndex: []HybridTableOutOfLineIndex{
					{
						Name:           "idx_email",
						Columns:        []string{"email"},
						IncludeColumns: []string{"id"},
					},
				},
			},
			Comment: String("test table"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE HYBRID TABLE %s ("id" NUMBER(38,0) PRIMARY KEY, "email" VARCHAR(200), UNIQUE ("email"), INDEX "idx_email" ("email") INCLUDE ("id")) COMMENT = 'test table'`, id.FullyQualifiedName())
	})

	t.Run("create with named inline constraint", func(t *testing.T) {
		opts := &CreateHybridTableOptions{
			name: id,
			ColumnsAndConstraints: HybridTableColumnsConstraintsAndIndexes{
				Columns: []HybridTableColumn{
					{
						Name: "id",
						Type: DataType("NUMBER(38,0)"),
						InlineConstraint: &ColumnInlineConstraint{
							Name: String("pk_id"),
							Type: ColumnConstraintTypePrimaryKey,
						},
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE HYBRID TABLE %s ("id" NUMBER(38,0) CONSTRAINT pk_id PRIMARY KEY)`, id.FullyQualifiedName())
	})

	t.Run("create with column comment and collate", func(t *testing.T) {
		opts := &CreateHybridTableOptions{
			name: id,
			ColumnsAndConstraints: HybridTableColumnsConstraintsAndIndexes{
				Columns: []HybridTableColumn{
					{
						Name: "id",
						Type: DataType("NUMBER(38,0)"),
						InlineConstraint: &ColumnInlineConstraint{
							Type: ColumnConstraintTypePrimaryKey,
						},
					},
					{
						Name:    "name",
						Type:    DataType("VARCHAR(100)"),
						NotNull: Bool(true),
						Collate: String("en-ci"),
						Comment: String("name column"),
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE HYBRID TABLE %s ("id" NUMBER(38,0) PRIMARY KEY, "name" VARCHAR(100) NOT NULL COLLATE 'en-ci' COMMENT 'name column')`, id.FullyQualifiedName())
	})
}

func TestHybridTables_CreateIndex_Ext(t *testing.T) {
	indexId := randomSchemaObjectIdentifier()
	tableId := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateIndexHybridTableOptions {
		return &CreateIndexHybridTableOptions{
			name:      indexId,
			TableName: tableId,
			Columns:   []string{"col1"},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateIndexHybridTableOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.TableName]", func(t *testing.T) {
		opts := defaultOpts()
		opts.TableName = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateIndexHybridTableOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("create index", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE INDEX %s ON %s ("col1")`, indexId.FullyQualifiedName(), tableId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Columns = []string{"col1", "col2"}
		opts.IncludeColumns = []string{"col3"}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE INDEX %s ON %s ("col1", "col2") INCLUDE ("col3")`, indexId.FullyQualifiedName(), tableId.FullyQualifiedName())
	})

	t.Run("with if not exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE INDEX IF NOT EXISTS %s ON %s ("col1")`, indexId.FullyQualifiedName(), tableId.FullyQualifiedName())
	})
}

func TestHybridTables_DropIndex_Ext(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DropIndexHybridTableOptions {
		return &DropIndexHybridTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropIndexHybridTableOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("drop index", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP INDEX %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP INDEX IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestHybridTables_ShowIndexes_Ext(t *testing.T) {
	defaultOpts := func() *ShowIndexesHybridTableOptions {
		return &ShowIndexesHybridTableOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowIndexesHybridTableOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("show indexes", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW INDEXES`)
	})

	t.Run("with in schema", func(t *testing.T) {
		schemaId := randomDatabaseObjectIdentifier()
		opts := defaultOpts()
		opts.In = &In{
			Schema: schemaId,
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW INDEXES IN SCHEMA %s`, schemaId.FullyQualifiedName())
	})
}

func TestHybridTables_Alter_DropMultipleColumns(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("alter: drop multiple columns", func(t *testing.T) {
		opts := &AlterHybridTableOptions{
			name:     id,
			IfExists: Bool(true),
			DropColumnAction: &HybridTableDropColumnAction{
				Columns: []string{"col1", "col2", "col3"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF EXISTS %s DROP COLUMN "col1", "col2", "col3"`, id.FullyQualifiedName())
	})

	t.Run("alter: drop column with if exists", func(t *testing.T) {
		opts := &AlterHybridTableOptions{
			name: id,
			DropColumnAction: &HybridTableDropColumnAction{
				IfExists: Bool(true),
				Columns:  []string{"col1"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s DROP COLUMN IF EXISTS "col1"`, id.FullyQualifiedName())
	})
}
