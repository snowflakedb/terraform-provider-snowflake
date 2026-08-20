package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	id := hybridTablesTestIdSchemaObjectIdentifier

	hybridTablesTests.Create.
		withExpectedSqlf(
			case_HybridTables_sql_Create_basic,
			`CREATE HYBRID TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Create_all,
			func(opts *CreateHybridTableOptions) {
				opts.OrReplace = new(true)
				opts.ColumnsAndConstraints = HybridTableColumnsConstraintsAndIndexes{
					Columns: []HybridTableColumn{
						{
							Name:     "ID",
							DataType: DataTypeNumber,
							InlineConstraint: &ColumnInlineConstraint{
								Name: new("pk_id"),
								Type: ColumnConstraintTypePrimaryKey,
							},
						},
						{
							Name:     "NAME",
							DataType: DataTypeVARCHAR,
							NotNull:  new(true),
							Collate:  new("en-ci"),
							Comment:  new("the name"),
						},
						{
							Name:     "REF_ID",
							DataType: DataTypeNumber,
							InlineConstraint: &ColumnInlineConstraint{
								Type: ColumnConstraintTypeForeignKey,
								ForeignKey: &InlineForeignKey{
									TableName:  "other_table",
									ColumnName: []string{"ID"},
								},
							},
						},
					},
					OutOfLineConstraint: []HybridTableOutOfLineConstraint{
						{
							ColumnConstraintType: ColumnConstraintTypeUnique,
							Columns:              []Column{{Value: "NAME"}},
						},
					},
					OutOfLineIndex: []HybridTableOutOfLineIndex{
						{
							Name:           "idx_name",
							Columns:        []Column{{Value: "NAME"}},
							IncludeColumns: []Column{{Value: "ID"}},
						},
					},
				}
				opts.Comment = new("test comment")
			},
			`CREATE OR REPLACE HYBRID TABLE %s ("ID" NUMBER CONSTRAINT pk_id PRIMARY KEY, "NAME" VARCHAR NOT NULL COLLATE 'en-ci' COMMENT 'the name', "REF_ID" NUMBER FOREIGN KEY REFERENCES other_table (ID), UNIQUE ("NAME"), INDEX "idx_name" ("NAME") INCLUDE ("ID")) COMMENT = 'test comment'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withColumnsAndConstraints",
			func(opts *CreateHybridTableOptions) {
				opts.ColumnsAndConstraints = HybridTableColumnsConstraintsAndIndexes{
					Columns: []HybridTableColumn{
						{
							Name:     "ID",
							DataType: DataTypeNumber,
							InlineConstraint: &ColumnInlineConstraint{
								Type: ColumnConstraintTypePrimaryKey,
							},
						},
						{
							Name:     "NAME",
							DataType: DataTypeVARCHAR,
							NotNull:  new(true),
							Comment:  new("the name"),
						},
					},
					OutOfLineConstraint: []HybridTableOutOfLineConstraint{
						{
							ColumnConstraintType: ColumnConstraintTypeUnique,
							Columns:              []Column{{Value: "NAME"}},
						},
					},
					OutOfLineIndex: []HybridTableOutOfLineIndex{
						{
							Name:    "idx_name",
							Columns: []Column{{Value: "NAME"}},
						},
					},
				}
			},
			`CREATE HYBRID TABLE %s ("ID" NUMBER PRIMARY KEY, "NAME" VARCHAR NOT NULL COMMENT 'the name', UNIQUE ("NAME"), INDEX "idx_name" ("NAME"))`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withColumnsAndInlinePrimaryKey",
			func(opts *CreateHybridTableOptions) {
				opts.ColumnsAndConstraints = HybridTableColumnsConstraintsAndIndexes{
					Columns: []HybridTableColumn{
						{
							Name:     "id",
							DataType: DataType("NUMBER(38,0)"),
							InlineConstraint: &ColumnInlineConstraint{
								Type: ColumnConstraintTypePrimaryKey,
							},
						},
						{
							Name:     "name",
							DataType: DataType("VARCHAR(100)"),
							NotNull:  new(true),
						},
					},
				}
			},
			`CREATE HYBRID TABLE %s ("id" NUMBER(38,0) PRIMARY KEY, "name" VARCHAR(100) NOT NULL)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withOutOfLineConstraintAndIndex",
			func(opts *CreateHybridTableOptions) {
				opts.ColumnsAndConstraints = HybridTableColumnsConstraintsAndIndexes{
					Columns: []HybridTableColumn{
						{
							Name:     "id",
							DataType: DataType("NUMBER(38,0)"),
							InlineConstraint: &ColumnInlineConstraint{
								Type: ColumnConstraintTypePrimaryKey,
							},
						},
						{
							Name:     "email",
							DataType: DataType("VARCHAR(200)"),
						},
					},
					OutOfLineConstraint: []HybridTableOutOfLineConstraint{
						{
							ColumnConstraintType: ColumnConstraintTypeUnique,
							Columns:              []Column{{Value: "email"}},
						},
					},
					OutOfLineIndex: []HybridTableOutOfLineIndex{
						{
							Name:           "idx_email",
							Columns:        []Column{{Value: "email"}},
							IncludeColumns: []Column{{Value: "id"}},
						},
					},
				}
				opts.Comment = new("test table")
			},
			`CREATE HYBRID TABLE %s ("id" NUMBER(38,0) PRIMARY KEY, "email" VARCHAR(200), UNIQUE ("email"), INDEX "idx_email" ("email") INCLUDE ("id")) COMMENT = 'test table'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withOutOfLineForeignKey",
			func(opts *CreateHybridTableOptions) {
				opts.ColumnsAndConstraints = HybridTableColumnsConstraintsAndIndexes{
					Columns: []HybridTableColumn{
						{
							Name:     "id",
							DataType: DataType("NUMBER(38,0)"),
							InlineConstraint: &ColumnInlineConstraint{
								Type: ColumnConstraintTypePrimaryKey,
							},
						},
						{
							Name:     "parent_id",
							DataType: DataType("NUMBER(38,0)"),
						},
					},
					OutOfLineConstraint: []HybridTableOutOfLineConstraint{
						{
							Name:                 new("fk_parent"),
							ColumnConstraintType: ColumnConstraintTypeForeignKey,
							Columns:              []Column{{Value: "parent_id"}},
							ForeignKey: &HybridTableOutOfLineForeignKey{
								TableName:   hybridTablesTestIdSchemaObjectIdentifier,
								ColumnNames: []Column{{Value: "id"}},
							},
						},
					},
				}
			},
			`CREATE HYBRID TABLE %s ("id" NUMBER(38,0) PRIMARY KEY, "parent_id" NUMBER(38,0), CONSTRAINT "fk_parent" FOREIGN KEY ("parent_id") REFERENCES %s ("id"))`,
			id.FullyQualifiedName(), hybridTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withNamedInlineConstraint",
			func(opts *CreateHybridTableOptions) {
				opts.ColumnsAndConstraints = HybridTableColumnsConstraintsAndIndexes{
					Columns: []HybridTableColumn{
						{
							Name:     "id",
							DataType: DataType("NUMBER(38,0)"),
							InlineConstraint: &ColumnInlineConstraint{
								Name: new("pk_id"),
								Type: ColumnConstraintTypePrimaryKey,
							},
						},
					},
				}
			},
			`CREATE HYBRID TABLE %s ("id" NUMBER(38,0) CONSTRAINT pk_id PRIMARY KEY)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withColumnCommentAndCollate",
			func(opts *CreateHybridTableOptions) {
				opts.ColumnsAndConstraints = HybridTableColumnsConstraintsAndIndexes{
					Columns: []HybridTableColumn{
						{
							Name:     "id",
							DataType: DataType("NUMBER(38,0)"),
							InlineConstraint: &ColumnInlineConstraint{
								Type: ColumnConstraintTypePrimaryKey,
							},
						},
						{
							Name:     "name",
							DataType: DataType("VARCHAR(100)"),
							NotNull:  new(true),
							Collate:  new("en-ci"),
							Comment:  new("name column"),
						},
					},
				}
			},
			`CREATE HYBRID TABLE %s ("id" NUMBER(38,0) PRIMARY KEY, "name" VARCHAR(100) NOT NULL COLLATE 'en-ci' COMMENT 'name column')`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withRetentionParameters",
			func(opts *CreateHybridTableOptions) {
				opts.ColumnsAndConstraints = HybridTableColumnsConstraintsAndIndexes{
					Columns: []HybridTableColumn{
						{
							Name:     "id",
							DataType: DataType("NUMBER(38,0)"),
							InlineConstraint: &ColumnInlineConstraint{
								Type: ColumnConstraintTypePrimaryKey,
							},
						},
					},
				}
				opts.DataRetentionTimeInDays = new(7)
				opts.MaxDataExtensionTimeInDays = new(14)
				opts.Comment = new("with retention")
			},
			`CREATE HYBRID TABLE %s ("id" NUMBER(38,0) PRIMARY KEY) DATA_RETENTION_TIME_IN_DAYS = 7 MAX_DATA_EXTENSION_TIME_IN_DAYS = 14 COMMENT = 'with retention'`,
			id.FullyQualifiedName(),
		)

	newID := randomSchemaObjectIdentifierInSchema(id.SchemaId())
	alterDataType := DataType("VARCHAR(200)")

	hybridTablesTests.Alter.
		withModify(
			case_HybridTables_validation_Alter_opts_AlterColumnAction_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterHybridTableOptions) {
				opts.AlterColumnAction = []HybridTableAlterColumnAction{
					{
						ColumnName:  "col1",
						DropDefault: new(true),
						Comment:     new("c"),
					},
				}
			},
		).
		withModify(
			case_HybridTables_validation_Alter_opts_AlterColumnAction_ExactlyOneValueSet_OneValidOneInvalid,
			func(opts *AlterHybridTableOptions) {
				opts.AlterColumnAction = []HybridTableAlterColumnAction{
					{
						ColumnName:  "col1",
						DropDefault: new(true),
					},
					{
						ColumnName: "col2",
					},
				}
			},
		).
		withModify(
			case_HybridTables_validation_Alter_opts_ClusteringAction_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterHybridTableOptions) {
				opts.ClusteringAction = &HybridTableClusteringAction{
					ClusterBy: []string{"col1"},
					Recluster: &HybridTableReclusterAction{},
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_RenameTo,
			func(opts *AlterHybridTableOptions) { opts.RenameTo = &newID },
			`ALTER TABLE %s RENAME TO %s`, id.FullyQualifiedName(), newID.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_AddColumnAction,
			func(opts *AlterHybridTableOptions) {
				opts.AddColumnAction = &HybridTableAddColumnAction{
					Name:     "NEW_COLUMN",
					DataType: DataTypeVARCHAR,
				}
			},
			`ALTER TABLE %s ADD COLUMN "NEW_COLUMN" VARCHAR`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_AddColumnAction_allOptions",
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.AddColumnAction = &HybridTableAddColumnAction{
					IfNotExists: new(true),
					Name:        "NEW_COLUMN",
					DataType:    DataTypeVARCHAR,
					Collate:     new("utf8"),
					InlineConstraint: &ColumnInlineConstraint{
						Name: new("uq_new"),
						Type: ColumnConstraintTypeUnique,
					},
					Comment: new("new column comment"),
				}
			},
			`ALTER TABLE IF EXISTS %s ADD COLUMN IF NOT EXISTS "NEW_COLUMN" VARCHAR COLLATE 'utf8' CONSTRAINT uq_new UNIQUE COMMENT 'new column comment'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_ConstraintAction,
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.ConstraintAction = &HybridTableConstraintAction{
					Drop: &HybridTableConstraintActionDrop{
						ConstraintName: new("my_constraint"),
						Restrict:       new(true),
					},
				}
			},
			`ALTER TABLE IF EXISTS %s DROP CONSTRAINT "my_constraint" RESTRICT`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ConstraintAction_dropUniqueWithColumns",
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.ConstraintAction = &HybridTableConstraintAction{
					Drop: &HybridTableConstraintActionDrop{
						Unique:  new(true),
						Columns: []Column{{Value: "col1"}},
					},
				}
			},
			`ALTER TABLE IF EXISTS %s DROP UNIQUE ("col1")`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ConstraintAction_dropForeignKeyWithCascade",
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.ConstraintAction = &HybridTableConstraintAction{
					Drop: &HybridTableConstraintActionDrop{
						ForeignKey: new(true),
						Columns:    []Column{{Value: "col1"}, {Value: "col2"}},
						Cascade:    new(true),
					},
				}
			},
			`ALTER TABLE IF EXISTS %s DROP FOREIGN KEY ("col1", "col2") CASCADE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ConstraintAction_renameConstraint",
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.ConstraintAction = &HybridTableConstraintAction{
					Rename: &HybridTableConstraintActionRename{
						OldName: "old_constraint_name",
						NewName: "new_constraint_name",
					},
				}
			},
			`ALTER TABLE IF EXISTS %s RENAME CONSTRAINT "old_constraint_name" TO "new_constraint_name"`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_AlterColumnAction,
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.AlterColumnAction = []HybridTableAlterColumnAction{
					{
						ColumnName: "column1",
						Comment:    new("column comment"),
					},
				}
			},
			`ALTER TABLE IF EXISTS %s ALTER COLUMN "column1" COMMENT 'column comment'`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_AlterColumnAction_unsetComment",
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.AlterColumnAction = []HybridTableAlterColumnAction{
					{
						ColumnName:   "column1",
						UnsetComment: new(true),
					},
				}
			},
			`ALTER TABLE IF EXISTS %s ALTER COLUMN "column1" UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_AlterColumnAction_dropDefault",
			func(opts *AlterHybridTableOptions) {
				opts.AlterColumnAction = []HybridTableAlterColumnAction{
					{
						ColumnName:  "column1",
						DropDefault: new(true),
					},
				}
			},
			`ALTER TABLE %s ALTER COLUMN "column1" DROP DEFAULT`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_AlterColumnAction_setDataType",
			func(opts *AlterHybridTableOptions) {
				opts.AlterColumnAction = []HybridTableAlterColumnAction{
					{
						ColumnName: "column1",
						DataType:   &alterDataType,
					},
				}
			},
			`ALTER TABLE %s ALTER COLUMN "column1" SET DATA TYPE VARCHAR(200)`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_DropColumnAction,
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.DropColumnAction = &HybridTableDropColumnAction{
					Columns: []Column{{Value: "column_to_drop"}},
				}
			},
			`ALTER TABLE IF EXISTS %s DROP COLUMN "column_to_drop"`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_DropColumnAction_multipleColumns",
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.DropColumnAction = &HybridTableDropColumnAction{
					Columns: []Column{{Value: "col1"}, {Value: "col2"}, {Value: "col3"}},
				}
			},
			`ALTER TABLE IF EXISTS %s DROP COLUMN "col1", "col2", "col3"`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_DropColumnAction_ifExistsOnColumn",
			func(opts *AlterHybridTableOptions) {
				opts.DropColumnAction = &HybridTableDropColumnAction{
					IfExists: new(true),
					Columns:  []Column{{Value: "col1"}},
				}
			},
			`ALTER TABLE %s DROP COLUMN IF EXISTS "col1"`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_DropIndexAction,
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.DropIndexAction = &HybridTableDropIndexAction{
					IndexName: "idx_name",
				}
			},
			`ALTER TABLE IF EXISTS %s DROP INDEX "idx_name"`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_ClusteringAction,
			func(opts *AlterHybridTableOptions) {
				opts.ClusteringAction = &HybridTableClusteringAction{
					ClusterBy: []string{"col1", "col2"},
				}
			},
			`ALTER TABLE %s CLUSTER BY (col1, col2)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ClusteringAction_recluster",
			func(opts *AlterHybridTableOptions) {
				opts.ClusteringAction = &HybridTableClusteringAction{
					Recluster: &HybridTableReclusterAction{},
				}
			},
			`ALTER TABLE %s RECLUSTER`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ClusteringAction_reclusterWithOptions",
			func(opts *AlterHybridTableOptions) {
				opts.ClusteringAction = &HybridTableClusteringAction{
					Recluster: &HybridTableReclusterAction{
						MaxSize: new(1000),
						Where:   new("col1 > 100"),
					},
				}
			},
			`ALTER TABLE %s RECLUSTER MAX_SIZE = 1000 WHERE col1 > 100`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ClusteringAction_suspendRecluster",
			func(opts *AlterHybridTableOptions) {
				opts.ClusteringAction = &HybridTableClusteringAction{
					ChangeReclusterState: &HybridTableReclusterChangeState{
						State: Pointer(ReclusterStateSuspend),
					},
				}
			},
			`ALTER TABLE %s SUSPEND RECLUSTER`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ClusteringAction_resumeRecluster",
			func(opts *AlterHybridTableOptions) {
				opts.ClusteringAction = &HybridTableClusteringAction{
					ChangeReclusterState: &HybridTableReclusterChangeState{
						State: Pointer(ReclusterStateResume),
					},
				}
			},
			`ALTER TABLE %s RESUME RECLUSTER`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ClusteringAction_dropClusteringKey",
			func(opts *AlterHybridTableOptions) {
				opts.ClusteringAction = &HybridTableClusteringAction{
					DropClusteringKey: new(true),
				}
			},
			`ALTER TABLE %s DROP CLUSTERING KEY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_Set,
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.Set = &HybridTableSetProperties{
					DataRetentionTimeInDays:    new(14),
					MaxDataExtensionTimeInDays: new(28),
					Comment:                    new("updated comment"),
				}
			},
			`ALTER TABLE IF EXISTS %s SET DATA_RETENTION_TIME_IN_DAYS = 14 MAX_DATA_EXTENSION_TIME_IN_DAYS = 28 COMMENT = 'updated comment'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Alter_Unset,
			func(opts *AlterHybridTableOptions) {
				opts.Unset = &HybridTableUnsetProperties{
					Comment: new(true),
				}
			},
			`ALTER TABLE %s UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_allPropertiesInSingleAlter",
			func(opts *AlterHybridTableOptions) {
				opts.IfExists = new(true)
				opts.Unset = &HybridTableUnsetProperties{
					Comment:                    new(true),
					DataRetentionTimeInDays:    new(true),
					MaxDataExtensionTimeInDays: new(true),
				}
			},
			`ALTER TABLE IF EXISTS %s UNSET COMMENT, DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS`,
			id.FullyQualifiedName(),
		)

	hybridTablesTests.Drop.
		withExpectedSqlf(
			case_HybridTables_sql_Drop_basic,
			`DROP TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Drop_all,
			func(opts *DropHybridTableOptions) {
				opts.IfExists = new(true)
				opts.Cascade = new(true)
			},
			`DROP TABLE IF EXISTS %s CASCADE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Drop_allOptionsRestrict",
			func(opts *DropHybridTableOptions) {
				opts.IfExists = new(true)
				opts.Restrict = new(true)
			},
			`DROP TABLE IF EXISTS %s RESTRICT`, id.FullyQualifiedName(),
		)

	hybridTablesTests.Show.
		withExpectedSqlf(
			case_HybridTables_sql_Show_basic,
			`SHOW HYBRID TABLES`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Show_all,
			func(opts *ShowHybridTableOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: new("some_pattern")}
				opts.In = &TableIn{In: In{Schema: NewDatabaseObjectIdentifier("db", "schema")}}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			`SHOW TERSE HYBRID TABLES LIKE 'some_pattern' IN SCHEMA "db"."schema" STARTS WITH 'prefix' LIMIT 10`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Show_Like,
			func(opts *ShowHybridTableOptions) { opts.Like = &Like{Pattern: new("pattern_test")} },
			`SHOW HYBRID TABLES LIKE 'pattern_test'`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Show_In,
			func(opts *ShowHybridTableOptions) {
				opts.In = &TableIn{In: In{Database: NewAccountObjectIdentifier("test_db")}}
			},
			`SHOW HYBRID TABLES IN DATABASE "test_db"`,
		).
		withAdditionalSqlCasef(
			"sql_Show_inSchema",
			func(opts *ShowHybridTableOptions) {
				opts.In = &TableIn{In: In{Schema: NewDatabaseObjectIdentifier("test_db", "test_schema")}}
			},
			`SHOW HYBRID TABLES IN SCHEMA "test_db"."test_schema"`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Show_StartsWith,
			func(opts *ShowHybridTableOptions) { opts.StartsWith = new("prefix") },
			`SHOW HYBRID TABLES STARTS WITH 'prefix'`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_Show_Limit,
			func(opts *ShowHybridTableOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("table_name")}
			},
			`SHOW HYBRID TABLES LIMIT 10 FROM 'table_name'`,
		)

	hybridTablesTests.Describe.
		withExpectedSqlf(
			case_HybridTables_sql_Describe_basic,
			`DESCRIBE TABLE %s`, id.FullyQualifiedName(),
		)

	indexId := randomSchemaObjectIdentifier()
	tableId := randomSchemaObjectIdentifier()

	hybridTablesTests.CreateIndex.
		withDefaultOpts(func() *CreateIndexHybridTableOptions {
			return &CreateIndexHybridTableOptions{
				name:      indexId,
				TableName: tableId,
				Columns:   []Column{{Value: "col1"}},
			}
		}).
		withExpectedSqlf(
			case_HybridTables_sql_CreateIndex_basic,
			`CREATE INDEX %s ON %s ("col1")`, indexId.FullyQualifiedName(), tableId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_CreateIndex_all,
			func(opts *CreateIndexHybridTableOptions) {
				opts.OrReplace = new(true)
				opts.Columns = []Column{{Value: "col1"}, {Value: "col2"}}
				opts.IncludeColumns = []Column{{Value: "col3"}}
			},
			`CREATE OR REPLACE INDEX %s ON %s ("col1", "col2") INCLUDE ("col3")`, indexId.FullyQualifiedName(), tableId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateIndex_ifNotExists",
			func(opts *CreateIndexHybridTableOptions) { opts.IfNotExists = new(true) },
			`CREATE INDEX IF NOT EXISTS %s ON %s ("col1")`, indexId.FullyQualifiedName(), tableId.FullyQualifiedName(),
		)

	hybridTablesTests.DropIndex.
		withExpectedSqlf(
			case_HybridTables_sql_DropIndex_basic,
			`DROP INDEX %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_DropIndex_all,
			func(opts *DropIndexHybridTableOptions) { opts.IfExists = new(true) },
			`DROP INDEX IF EXISTS %s`, id.FullyQualifiedName(),
		)

	hybridTablesTests.ShowIndexes.
		withExpectedSqlf(
			case_HybridTables_sql_ShowIndexes_basic,
			`SHOW INDEXES`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_ShowIndexes_all,
			func(opts *ShowIndexesHybridTableOptions) {
				opts.Like = &Like{Pattern: new("idx_pattern")}
				opts.In = &TableIn{In: In{Schema: NewDatabaseObjectIdentifier("db", "schema")}}
				opts.StartsWith = new("idx_")
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			`SHOW INDEXES LIKE 'idx_pattern' IN SCHEMA "db"."schema" STARTS WITH 'idx_' LIMIT 10`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_ShowIndexes_Like,
			func(opts *ShowIndexesHybridTableOptions) { opts.Like = &Like{Pattern: new("idx_pattern")} },
			`SHOW INDEXES LIKE 'idx_pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_ShowIndexes_In,
			func(opts *ShowIndexesHybridTableOptions) {
				opts.In = &TableIn{In: In{Database: NewAccountObjectIdentifier("test_db")}}
			},
			`SHOW INDEXES IN DATABASE "test_db"`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_ShowIndexes_StartsWith,
			func(opts *ShowIndexesHybridTableOptions) { opts.StartsWith = new("idx_") },
			`SHOW INDEXES STARTS WITH 'idx_'`,
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_ShowIndexes_Limit,
			func(opts *ShowIndexesHybridTableOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			`SHOW INDEXES LIMIT 10`,
		)

	hybridTablesTests.ShowPrimaryKeys.
		withExpectedSqlf(
			case_HybridTables_sql_ShowPrimaryKeys_basic,
			`SHOW PRIMARY KEYS IN TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_ShowPrimaryKeys_all,
			func(opts *ShowPrimaryKeysHybridTableOptions) {},
			`SHOW PRIMARY KEYS IN TABLE %s`, id.FullyQualifiedName(),
		)

	hybridTablesTests.ShowUniqueKeys.
		withExpectedSqlf(
			case_HybridTables_sql_ShowUniqueKeys_basic,
			`SHOW UNIQUE KEYS IN TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_ShowUniqueKeys_all,
			func(opts *ShowUniqueKeysHybridTableOptions) {},
			`SHOW UNIQUE KEYS IN TABLE %s`, id.FullyQualifiedName(),
		)

	hybridTablesTests.ShowImportedKeys.
		withExpectedSqlf(
			case_HybridTables_sql_ShowImportedKeys_basic,
			`SHOW IMPORTED KEYS IN TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_HybridTables_sql_ShowImportedKeys_all,
			func(opts *ShowImportedKeysHybridTableOptions) {},
			`SHOW IMPORTED KEYS IN TABLE %s`, id.FullyQualifiedName(),
		)
}

func TestHybridTableDetailsRow_SplitTypeAndCollation(t *testing.T) {
	testCases := []struct {
		Name              string
		Value             string
		ExpectedType      string
		ExpectedCollation *string
	}{
		{
			Name:              "with utf8",
			Value:             "VARCHAR(10) COLLATE 'utf8'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: new("utf8"),
		},
		{
			Name:              "with locale",
			Value:             "VARCHAR(10) COLLATE 'en_US'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: new("en_US"),
		},
		{
			Name:              "with multiple specifiers",
			Value:             "VARCHAR(10) COLLATE 'fr_CA-ai-pi-trim'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: new("fr_CA-ai-pi-trim"),
		},
		{
			Name:              "with empty collation",
			Value:             "VARCHAR(10) COLLATE ''",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: new(""),
		},
		{
			Name:              "without collation",
			Value:             "NUMBER(38, 0)",
			ExpectedType:      "NUMBER(38, 0)",
			ExpectedCollation: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			row := hybridTableDetailsRow{Type: tc.Value}
			actualType, actualCollation := row.splitTypeAndCollation()
			assert.Equal(t, tc.ExpectedType, actualType)
			if tc.ExpectedCollation == nil {
				assert.Nil(t, actualCollation)
			} else {
				assert.Equal(t, *tc.ExpectedCollation, *actualCollation)
			}
		})
	}
}

func TestHybridTableConstraints_mergeKeyRows(t *testing.T) {
	testCases := []struct {
		name     string
		rows     []keyRow
		kind     ColumnConstraintType
		expected []HybridTableConstraint
	}{
		{
			name:     "empty input returns nil",
			rows:     nil,
			kind:     ColumnConstraintTypePrimaryKey,
			expected: nil,
		},
		{
			name: "single-column primary key",
			rows: []keyRow{
				{ConstraintName: "PK_T", ColumnName: "ID", KeySequence: 1},
			},
			kind: ColumnConstraintTypePrimaryKey,
			expected: []HybridTableConstraint{
				{Name: "PK_T", Kind: ColumnConstraintTypePrimaryKey, Columns: []string{"ID"}},
			},
		},
		{
			name: "two distinct unique constraints, columns ordered by key_sequence",
			rows: []keyRow{
				// Intentionally out of key_sequence order to prove sorting.
				{ConstraintName: "UQ_A", ColumnName: "COL_A2", KeySequence: 2},
				{ConstraintName: "UQ_A", ColumnName: "COL_A1", KeySequence: 1},
				{ConstraintName: "UQ_B", ColumnName: "COL_B", KeySequence: 1},
			},
			kind: ColumnConstraintTypeUnique,
			expected: []HybridTableConstraint{
				{Name: "UQ_A", Kind: ColumnConstraintTypeUnique, Columns: []string{"COL_A1", "COL_A2"}},
				{Name: "UQ_B", Kind: ColumnConstraintTypeUnique, Columns: []string{"COL_B"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, mergeKeyRows(tc.rows, tc.kind))
		})
	}
}

func TestHybridTableConstraints_mergeForeignKeyRows(t *testing.T) {
	t.Run("empty input returns nil", func(t *testing.T) {
		assert.Nil(t, mergeForeignKeyRows(nil))
	})

	t.Run("multi-column foreign key, columns and referenced columns ordered by key_sequence", func(t *testing.T) {
		// Two rows of the same FK, intentionally out of key_sequence order.
		rows := []TableImportedKey{
			{
				FkName:         "FK_T",
				FkColumnName:   "PARENT_B",
				KeySequence:    2,
				PkDatabaseName: "DB",
				PkSchemaName:   "SCH",
				PkTableName:    "PARENT",
				PkColumnName:   "B",
				DeleteRule:     "NO ACTION",
				UpdateRule:     "NO ACTION",
			},
			{
				FkName:         "FK_T",
				FkColumnName:   "PARENT_A",
				KeySequence:    1,
				PkDatabaseName: "DB",
				PkSchemaName:   "SCH",
				PkTableName:    "PARENT",
				PkColumnName:   "A",
				DeleteRule:     "NO ACTION",
				UpdateRule:     "NO ACTION",
			},
		}

		expected := []HybridTableConstraint{
			{
				Name:              "FK_T",
				Kind:              ColumnConstraintTypeForeignKey,
				Columns:           []string{"PARENT_A", "PARENT_B"},
				ReferencedTable:   NewSchemaObjectIdentifier("DB", "SCH", "PARENT"),
				ReferencedColumns: []string{"A", "B"},
				DeleteRule:        "NO ACTION",
				UpdateRule:        "NO ACTION",
			},
		}

		assert.Equal(t, expected, mergeForeignKeyRows(rows))
	})
}
