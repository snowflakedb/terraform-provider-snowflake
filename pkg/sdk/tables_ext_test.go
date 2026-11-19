package sdk

import (
	"database/sql"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/stretchr/testify/assert"
)

func init() {
	id := tablesTestIdSchemaObjectIdentifier
	databaseId := NewAccountObjectIdentifier(id.DatabaseName())
	sourceTable := randomSchemaObjectIdentifier()
	renameTarget := randomSchemaObjectIdentifierInSchema(id.SchemaId())
	swapWith := randomSchemaObjectIdentifierInSchema(id.SchemaId())
	fkRefId := randomSchemaObjectIdentifier()
	maskingPolicyId := randomSchemaObjectIdentifier()
	rowAccessPolicyId := randomSchemaObjectIdentifier()
	rowAccessPolicyId2 := randomSchemaObjectIdentifier()
	aggregationPolicyId := randomSchemaObjectIdentifier()
	joinPolicyId := randomSchemaObjectIdentifier()
	storageLifecyclePolicyId := randomSchemaObjectIdentifier()
	contactId := randomSchemaObjectIdentifier()
	dataMetricFunctionId := randomSchemaObjectIdentifier()
	executeAsRole := randomAccountObjectIdentifier()
	columnTagId1 := randomSchemaObjectIdentifier()
	columnTagId2 := randomSchemaObjectIdentifierInSchema(columnTagId1.SchemaId())
	tableTagId1 := randomSchemaObjectIdentifierInSchema(columnTagId1.SchemaId())
	tableTagId2 := randomSchemaObjectIdentifierInSchema(columnTagId1.SchemaId())
	columnComment := random.Comment()
	tableComment := random.Comment()
	likePattern := id.Name()
	showDatabaseId := NewAccountObjectIdentifier("database")

	tablesTests.Create.
		withDefaultOpts(func() *CreateTableOptions {
			return &CreateTableOptions{
				name: id,
				ColumnsAndConstraints: CreateTableColumnsAndConstraints{
					Columns: []TableColumn{tableTestColumn()},
				},
			}
		}).
		withModify(case_Tables_validation_Create_opts_ColumnsAndConstraints_OutOfLineConstraint_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateTableOptions) {
			opts.ColumnsAndConstraints.OutOfLineConstraint = []TableOutOfLineConstraint{{
				UniquePK: &TableOutOfLineUniquePK{Unique: new(true), Columns: []Column{{Value: "COLUMN_1"}}},
				FK:       &TableOutOfLineFK{References: fkRefId, Columns: []Column{{Value: "COLUMN_1"}}},
			}}
		}).
		withModify(case_Tables_validation_Create_opts_ColumnsAndConstraints_OutOfLineConstraint_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateTableOptions) {
			opts.ColumnsAndConstraints.OutOfLineConstraint = []TableOutOfLineConstraint{
				{UniquePK: &TableOutOfLineUniquePK{Unique: new(true), Columns: []Column{{Value: "COLUMN_1"}}}},
				{},
			}
		}).
		withAdditionalValidationCase(
			"validation_Create_noColumns",
			func(opts *CreateTableOptions) { opts.ColumnsAndConstraints.Columns = nil },
			errNotSet("CreateTableOptions", "Columns"),
		).
		withAdditionalValidationCase(
			"validation_Create_RowAccessPolicy_invalidIdentifier",
			func(opts *CreateTableOptions) {
				opts.RowAccessPolicy = &TableRowAccessPolicy{Name: emptySchemaObjectIdentifier, On: []string{"COLUMN_1"}}
			},
			errInvalidIdentifier("CreateTableOptions.RowAccessPolicy", "Name"),
		).
		withExpectedSqlf(
			case_Tables_sql_Create_basic,
			`CREATE TABLE %s ("FIRST_COLUMN" VARCHAR)`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Create_all,
			func(opts *CreateTableOptions) {
				opts.IfNotExists = new(true)
				opts.ColumnsAndConstraints = CreateTableColumnsAndConstraints{
					Columns: []TableColumn{{
						Name:       "FIRST_COLUMN",
						ColumnType: DataTypeVARCHAR,
						Collate:    new("de"),
						Comment:    &columnComment,
						DefaultValue: &ColumnDefaultValue{
							Identity: &ColumnIdentity{Start: 10, Increment: 1, Order: new(true)},
						},
						NotNull: new(true),
						MaskingPolicy: &TableColumnMaskingPolicy{
							MaskingPolicy: maskingPolicyId,
							Using:         []Column{{Value: "FOO"}, {Value: "BAR"}},
						},
						Tag: []TagAssociation{
							{Name: columnTagId1, Value: "v1"},
							{Name: columnTagId2, Value: "v2"},
						},
						InlineConstraint: &TableColumnInlineConstraint{UniquePK: &TableColumnInlineUniquePK{Name: new("Inline_Constraint"), PrimaryKey: new(true)}},
					}, {
						Name:       "DOUBLE_ID",
						ColumnType: DataTypeNumber,
						Virtual:    &TableColumnVirtual{GeneratedAlways: new(true), As: "(1 + 1)", Virtual: new(true)},
					}},
					OutOfLineConstraint: []TableOutOfLineConstraint{
						{FK: &TableOutOfLineFK{
							Name:       new("Out_Of_Line_Constraint"),
							Columns:    []Column{{Value: "COLUMN_1"}, {Value: "COLUMN_2"}},
							References: fkRefId,
							RefColumns: []Column{{Value: "COLUMN_3"}, {Value: "COLUMN_4"}},
							Match:      new(MatchTypeFull),
							On: &ForeignKeyOnAction{
								OnUpdate: new(ForeignKeySetNullAction),
								OnDelete: new(ForeignKeyRestrictAction),
							},
						}},
						{UniquePK: &TableOutOfLineUniquePK{
							Unique:            new(true),
							Columns:           []Column{{Value: "COLUMN_1"}},
							Enforced:          new(true),
							Deferrable:        new(true),
							InitiallyDeferred: new(true),
							Enable:            new(true),
							Rely:              new(true),
							Comment:           new("unique comment"),
						}},
						{CH: &TableOutOfLineCH{
							Name:             new("check_Positive"),
							Expression:       "COLUMN_1 > 0",
							EnableNovalidate: new(true),
						}},
					},
				}
				opts.ClusterBy = []string{"COLUMN_1", "COLUMN_2"}
				opts.EnableSchemaEvolution = new(true)
				opts.DataRetentionTimeInDays = new(10)
				opts.MaxDataExtensionTimeInDays = new(100)
				opts.ChangeTracking = new(true)
				opts.DefaultDdlCollation = new(StringAllowEmpty{Value: "en"})
				opts.ErrorLogging = new(true)
				opts.Copy = []TableCopyOption{TableCopyOptionComments, TableCopyOptionTags, TableCopyOptionGrants}
				opts.RowAccessPolicy = &TableRowAccessPolicy{Name: rowAccessPolicyId, On: []string{"COLUMN_1", "COLUMN_2"}}
				opts.AggregationPolicy = &TableAggregationPolicy{AggregationPolicy: aggregationPolicyId, EntityKey: []Column{{Value: "COLUMN_1"}}}
				opts.JoinPolicy = &TableJoinPolicy{JoinPolicy: joinPolicyId, AllowedJoinKeys: []Column{{Value: "COLUMN_1"}}}
				opts.StorageLifecyclePolicy = &CreateTableStorageLifecyclePolicy{
					StorageLifecyclePolicy: storageLifecyclePolicyId,
					On:                     []Column{{Value: "FIRST_COLUMN"}},
				}
				opts.Tag = []TagAssociation{
					{Name: tableTagId1, Value: "v1"},
					{Name: tableTagId2, Value: "v2"},
				}
				opts.Contact = []TableContact{{Purpose: "SUPPORT", Contact: contactId}}
				opts.Comment = &tableComment
				opts.RowTimestamp = new(true)
			},
			`CREATE TABLE IF NOT EXISTS %s ("FIRST_COLUMN" VARCHAR CONSTRAINT "Inline_Constraint" PRIMARY KEY NOT NULL COLLATE 'de' IDENTITY START 10 INCREMENT 1 ORDER MASKING POLICY %s USING ("FOO", "BAR") TAG (%s = 'v1', %s = 'v2') COMMENT '%s', "DOUBLE_ID" NUMBER GENERATED ALWAYS AS (1 + 1) VIRTUAL, CONSTRAINT "Out_Of_Line_Constraint" FOREIGN KEY ("COLUMN_1", "COLUMN_2") REFERENCES %s ("COLUMN_3", "COLUMN_4") MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT, UNIQUE ("COLUMN_1") ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE RELY COMMENT 'unique comment', CONSTRAINT "check_Positive" CHECK ( COLUMN_1 > 0 ) ENABLE NOVALIDATE) CLUSTER BY (COLUMN_1, COLUMN_2) ENABLE_SCHEMA_EVOLUTION = true DATA_RETENTION_TIME_IN_DAYS = 10 MAX_DATA_EXTENSION_TIME_IN_DAYS = 100 CHANGE_TRACKING = true DEFAULT_DDL_COLLATION = 'en' ERROR_LOGGING = true COPY = (COMMENTS, TAGS, GRANTS) ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) AGGREGATION POLICY %s ENTITY KEY ("COLUMN_1") JOIN POLICY %s ALLOWED JOIN KEYS ("COLUMN_1") STORAGE LIFECYCLE POLICY %s ON ("FIRST_COLUMN") TAG (%s = 'v1', %s = 'v2') WITH CONTACT (SUPPORT = %s) COMMENT = '%s' ROW_TIMESTAMP = true`,
			id.FullyQualifiedName(),
			maskingPolicyId.FullyQualifiedName(),
			columnTagId1.FullyQualifiedName(),
			columnTagId2.FullyQualifiedName(),
			columnComment,
			fkRefId.FullyQualifiedName(),
			rowAccessPolicyId.FullyQualifiedName(),
			aggregationPolicyId.FullyQualifiedName(),
			joinPolicyId.FullyQualifiedName(),
			storageLifecyclePolicyId.FullyQualifiedName(),
			tableTagId1.FullyQualifiedName(),
			tableTagId2.FullyQualifiedName(),
			contactId.FullyQualifiedName(),
			tableComment,
		).
		withAdditionalSqlCasef(
			"sql_Create_emptyDefaultDdlCollation",
			func(opts *CreateTableOptions) {
				opts.DefaultDdlCollation = new(StringAllowEmpty{Value: ""})
			},
			`CREATE TABLE %s ("FIRST_COLUMN" VARCHAR) DEFAULT_DDL_COLLATION = ''`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateTableOptions) {
				opts.OrReplace = new(true)
				opts.Copy = []TableCopyOption{TableCopyOptionGrants}
			},
			`CREATE OR REPLACE TABLE %s ("FIRST_COLUMN" VARCHAR) COPY = (GRANTS)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_lowercaseColumnNamesOnConstraints",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints = CreateTableColumnsAndConstraints{
					Columns: []TableColumn{{
						Name:             "id",
						ColumnType:       DataTypeNumber,
						InlineConstraint: &TableColumnInlineConstraint{UniquePK: &TableColumnInlineUniquePK{Name: new("pk_id"), PrimaryKey: new(true)}},
					}},
					OutOfLineConstraint: []TableOutOfLineConstraint{
						{UniquePK: &TableOutOfLineUniquePK{
							Name:    new("uq_code"),
							Unique:  new(true),
							Columns: []Column{{Value: "code"}},
						}},
					},
				}
			},
			`CREATE TABLE %s ("id" NUMBER CONSTRAINT "pk_id" PRIMARY KEY, CONSTRAINT "uq_code" UNIQUE ("code"))`,
			id.FullyQualifiedName(),
		)

	tablesTests.CreateAsSelect.
		withDefaultOpts(func() *CreateAsSelectTableOptions {
			return &CreateAsSelectTableOptions{
				name:    id,
				Columns: []TableAsSelectColumn{{Name: "a"}},
				Query:   "SELECT 1",
			}
		}).
		withExpectedSqlf(
			case_Tables_sql_CreateAsSelect_basic,
			`CREATE TABLE %s ("a") AS SELECT 1`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_CreateAsSelect_all,
			func(opts *CreateAsSelectTableOptions) {
				opts.OrReplace = new(true)
				opts.Columns = []TableAsSelectColumn{{
					Name:          "FIRST_COLUMN",
					ColumnType:    new(DataTypeVARCHAR),
					MaskingPolicy: &maskingPolicyId,
				}}
				opts.ClusterBy = []string{"COLUMN_1", "COLUMN_2"}
				opts.Copy = []TableCopyOption{TableCopyOptionComments, TableCopyOptionTags, TableCopyOptionGrants}
				opts.RowAccessPolicy = &TableRowAccessPolicy{Name: rowAccessPolicyId, On: []string{"COLUMN_1", "COLUMN_2"}}
				opts.Query = "SELECT * FROM ANOTHER_TABLE"
			},
			`CREATE OR REPLACE TABLE %s ("FIRST_COLUMN" VARCHAR MASKING POLICY %s) CLUSTER BY (COLUMN_1, COLUMN_2) COPY = (COMMENTS, TAGS, GRANTS) ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) AS SELECT * FROM ANOTHER_TABLE`,
			id.FullyQualifiedName(),
			maskingPolicyId.FullyQualifiedName(),
			rowAccessPolicyId.FullyQualifiedName(),
		)

	tablesTests.CreateUsingTemplate.
		withDefaultOpts(func() *CreateUsingTemplateTableOptions {
			return &CreateUsingTemplateTableOptions{
				name:  id,
				Query: []string{"sample_data"},
			}
		}).
		withExpectedSqlf(
			case_Tables_sql_CreateUsingTemplate_basic,
			`CREATE TABLE %s USING TEMPLATE (sample_data)`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_CreateUsingTemplate_all,
			func(opts *CreateUsingTemplateTableOptions) {
				opts.OrReplace = new(true)
				opts.Copy = []TableCopyOption{TableCopyOptionComments, TableCopyOptionTags, TableCopyOptionGrants}
			},
			`CREATE OR REPLACE TABLE %s COPY = (COMMENTS, TAGS, GRANTS) USING TEMPLATE (sample_data)`,
			id.FullyQualifiedName(),
		)

	tablesTests.CreateLike.
		withDefaultOpts(func() *CreateLikeTableOptions {
			return &CreateLikeTableOptions{
				name:        id,
				SourceTable: sourceTable,
			}
		}).
		withExpectedSqlf(
			case_Tables_sql_CreateLike_basic,
			`CREATE TABLE %s LIKE %s`, id.FullyQualifiedName(), sourceTable.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_CreateLike_all,
			func(opts *CreateLikeTableOptions) {
				opts.OrReplace = new(true)
				opts.Copy = []TableCopyOption{TableCopyOptionComments, TableCopyOptionTags, TableCopyOptionGrants}
				opts.ClusterBy = []string{"date", "id"}
			},
			`CREATE OR REPLACE TABLE %s LIKE %s CLUSTER BY (date, id) COPY = (COMMENTS, TAGS, GRANTS)`,
			id.FullyQualifiedName(), sourceTable.FullyQualifiedName(),
		)

	tablesTests.CreateClone.
		withDefaultOpts(func() *CreateCloneTableOptions {
			return &CreateCloneTableOptions{
				name:        id,
				SourceTable: sourceTable,
			}
		}).
		withExpectedSqlf(
			case_Tables_sql_CreateClone_basic,
			`CREATE TABLE %s CLONE %s`, id.FullyQualifiedName(), sourceTable.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_CreateClone_all,
			func(opts *CreateCloneTableOptions) {
				opts.OrReplace = new(true)
				opts.Copy = []TableCopyOption{TableCopyOptionComments, TableCopyOptionTags, TableCopyOptionGrants}
				opts.ClonePoint = &TableClonePoint{
					Moment: CloneMomentAt,
					At:     TimeTravel{Offset: new(0)},
				}
			},
			`CREATE OR REPLACE TABLE %s CLONE %s AT (OFFSET => 0) COPY = (COMMENTS, TAGS, GRANTS)`,
			id.FullyQualifiedName(), sourceTable.FullyQualifiedName(),
		)

	tablesTests.Alter.
		withModify(case_Tables_validation_Alter_opts_ClusteringAction_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterTableOptions) {
			opts.ClusteringAction = &TableClusteringAction{ClusterBy: []string{"a"}, DropClusteringKey: new(true)}
		}).
		withModify(case_Tables_validation_Alter_opts_SearchOptimizationAction_Drop_On_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterTableOptions) {
			opts.SearchOptimizationAction = &TableSearchOptimizationAction{Drop: &TableDropSearchOptimization{
				On: []TableDropSearchOptimizationOn{{ColumnName: new("c"), ExpressionId: new("e")}},
			}}
		}).
		withModify(case_Tables_validation_Alter_opts_SearchOptimizationAction_Drop_On_ExactlyOneValueSet_OneValidOneInvalid, func(opts *AlterTableOptions) {
			opts.SearchOptimizationAction = &TableSearchOptimizationAction{Drop: &TableDropSearchOptimization{
				On: []TableDropSearchOptimizationOn{{ColumnName: new("c")}, {}},
			}}
		}).
		withModify(case_Tables_validation_Alter_opts_ColumnAction_Alter_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterTableOptions) {
			opts.ColumnAction = &TableColumnAction{Alter: []TableColumnAlterAction{{
				DropDefault: new(true),
				SetDefault:  new(SequenceName("sequence")),
			}}}
		}).
		withModify(case_Tables_validation_Alter_opts_ColumnAction_Alter_ExactlyOneValueSet_OneValidOneInvalid, func(opts *AlterTableOptions) {
			opts.ColumnAction = &TableColumnAction{Alter: []TableColumnAlterAction{
				{DropDefault: new(true)},
				{},
			}}
		}).
		withModify(case_Tables_validation_Alter_opts_ModifyDataMetricFunction_ModifiedFunctions_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterTableOptions) {
			opts.ModifyDataMetricFunction = &TableModifyDataMetricFunction{
				ModifiedFunctions: []TableModifiedFunctions{{
					DataMetricFunction: dataMetricFunctionId,
					On:                 []Column{{Value: "ID"}},
					Status:             new(TableDataMetricScheduleStatusOperationSuspend),
					Filter:             new("ID > 0"),
				}},
			}
		}).
		withModify(case_Tables_validation_Alter_opts_ModifyDataMetricFunction_ModifiedFunctions_ExactlyOneValueSet_OneValidOneInvalid, func(opts *AlterTableOptions) {
			opts.ModifyDataMetricFunction = &TableModifyDataMetricFunction{
				ModifiedFunctions: []TableModifiedFunctions{
					{
						DataMetricFunction: dataMetricFunctionId,
						On:                 []Column{{Value: "ID"}},
						Status:             new(TableDataMetricScheduleStatusOperationSuspend),
					},
					{
						DataMetricFunction: dataMetricFunctionId,
						On:                 []Column{{Value: "ID"}},
					},
				},
			}
		}).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_RenameTo,
			func(opts *AlterTableOptions) { opts.RenameTo = &renameTarget },
			`ALTER TABLE %s RENAME TO %s`, id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_SwapWith,
			func(opts *AlterTableOptions) { opts.SwapWith = &swapWith },
			`ALTER TABLE %s SWAP WITH %s`, id.FullyQualifiedName(), swapWith.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_ClusteringAction,
			func(opts *AlterTableOptions) {
				opts.ClusteringAction = &TableClusteringAction{ClusterBy: []string{"date", "id"}}
			},
			`ALTER TABLE %s CLUSTER BY (date, id)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_recluster",
			func(opts *AlterTableOptions) {
				opts.ClusteringAction = &TableClusteringAction{Recluster: &TableReclusterAction{MaxSize: new(1024), Condition: new("name = 'John'")}}
			},
			`ALTER TABLE %s RECLUSTER MAX_SIZE = 1024 WHERE name = 'John'`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_suspendRecluster",
			func(opts *AlterTableOptions) {
				opts.ClusteringAction = &TableClusteringAction{ChangeReclusterState: &TableReclusterChangeState{State: new(ReclusterStateSuspend)}}
			},
			`ALTER TABLE %s SUSPEND RECLUSTER`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_dropClusteringKey",
			func(opts *AlterTableOptions) {
				opts.ClusteringAction = &TableClusteringAction{DropClusteringKey: new(true)}
			},
			`ALTER TABLE %s DROP CLUSTERING KEY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_ColumnAction,
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{Add: &TableColumnAddAction{
					IfNotExists: new(true),
					Name:        "NEXT_COLUMN",
					ColumnType:  DataTypeVARCHAR,
					Collate:     new("utf8"),
					DefaultValue: &ColumnDefaultValue{
						Identity: &ColumnIdentity{Start: 10, Increment: 1},
					},
				}}
			},
			`ALTER TABLE %s ADD COLUMN IF NOT EXISTS "NEXT_COLUMN" VARCHAR COLLATE 'utf8' IDENTITY START 10 INCREMENT 1`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addVirtualColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{Add: &TableColumnAddAction{
					Name:       "DOUBLE_ID",
					ColumnType: DataTypeNumber,
					Virtual:    &TableColumnVirtual{GeneratedAlways: new(true), As: "(ID * 2)", Virtual: new(true)},
				}}
			},
			`ALTER TABLE %s ADD COLUMN "DOUBLE_ID" NUMBER GENERATED ALWAYS AS (ID * 2) VIRTUAL`,
			id.FullyQualifiedName(),
		).
		// https://github.com/snowflakedb/terraform-provider-snowflake/issues/4730
		// Adding a column with a non-identity default used to panic in TableColumnActionRequest.toOpts.
		withAdditionalSqlCasef(
			"sql_Alter_addColumn_constantDefault",
			func(opts *AlterTableOptions) {
				*opts = *NewAlterTableRequest(id).WithColumnAction(
					*NewTableColumnActionRequest().WithAdd(
						*NewTableColumnAddActionRequest("NEW_BOOLEAN_COLUMN", DataTypeBoolean).
							WithDefaultValue(*NewColumnDefaultValueRequest().WithExpression("FALSE")),
					),
				).toOpts()
			},
			`ALTER TABLE %s ADD COLUMN "NEW_BOOLEAN_COLUMN" BOOLEAN DEFAULT FALSE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addColumn_identityDefault",
			func(opts *AlterTableOptions) {
				*opts = *NewAlterTableRequest(id).WithColumnAction(
					*NewTableColumnActionRequest().WithAdd(
						*NewTableColumnAddActionRequest("NEXT_COLUMN", DataTypeVARCHAR).
							WithDefaultValue(*NewColumnDefaultValueRequest().WithIdentity(*NewColumnIdentityRequest(10, 1))),
					),
				).toOpts()
			},
			`ALTER TABLE %s ADD COLUMN "NEXT_COLUMN" VARCHAR IDENTITY START 10 INCREMENT 1`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addColumn_notNull",
			func(opts *AlterTableOptions) {
				*opts = *NewAlterTableRequest(id).WithColumnAction(
					*NewTableColumnActionRequest().WithAdd(
						*NewTableColumnAddActionRequest("NEXT_COLUMN", DataTypeVARCHAR).WithNotNull(true),
					),
				).toOpts()
			},
			`ALTER TABLE %s ADD COLUMN "NEXT_COLUMN" VARCHAR NOT NULL`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_renameColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{Rename: &TableRenameColumnAction{OldName: "OLD_NAME", NewName: "NEW_NAME"}}
			},
			`ALTER TABLE %s RENAME COLUMN "OLD_NAME" TO "NEW_NAME"`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_alterColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{Alter: []TableColumnAlterAction{
					{Name: "COLUMN_1", DropDefault: new(true)},
					{Name: "COLUMN_1", SetDefault: new(SequenceName("SEQUENCE_1"))},
					{Name: "COLUMN_1", UnsetComment: new(true)},
					{Name: "COLUMN_2", DropDefault: new(true)},
					{Name: "COLUMN_2", SetDefault: new(SequenceName("SEQUENCE_2"))},
					{Name: "COLUMN_2", Comment: new("comment")},
					{Name: "COLUMN_2", DataType: &TableColumnAlterDataType{DataType: DataTypeVARCHAR, Collate: new("utf8")}},
					{Name: "COLUMN_2", NotNullConstraint: &TableColumnNotNullConstraint{Drop: new(true)}},
				}}
			},
			`ALTER TABLE %s ALTER COLUMN "COLUMN_1" DROP DEFAULT, COLUMN "COLUMN_1" SET DEFAULT SEQUENCE_1.NEXTVAL, COLUMN "COLUMN_1" UNSET COMMENT, COLUMN "COLUMN_2" DROP DEFAULT, COLUMN "COLUMN_2" SET DEFAULT SEQUENCE_2.NEXTVAL, COLUMN "COLUMN_2" COMMENT 'comment', COLUMN "COLUMN_2" SET DATA TYPE VARCHAR COLLATE 'utf8', COLUMN "COLUMN_2" DROP NOT NULL`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setMaskingPolicyOnColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{SetMaskingPolicy: &TableSetColumnMaskingPolicy{
					Name:          "COLUMN_1",
					MaskingPolicy: maskingPolicyId,
					Using:         []Column{{Value: "FOO"}, {Value: "BAR"}},
					Force:         new(true),
				}}
			},
			`ALTER TABLE %s ALTER COLUMN "COLUMN_1" SET MASKING POLICY %s USING ("FOO", "BAR") FORCE`,
			id.FullyQualifiedName(), maskingPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_unsetMaskingPolicyOnColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{UnsetMaskingPolicy: &TableUnsetColumnMaskingPolicy{Name: "COLUMN_1"}}
			},
			`ALTER TABLE %s ALTER COLUMN "COLUMN_1" UNSET MASKING POLICY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setTagsOnColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{SetTags: &TableSetColumnTags{
					Name: "COLUMN_1",
					SetTags: []TagAssociation{
						{Name: columnTagId1, Value: "v1"},
						{Name: columnTagId2, Value: "v2"},
					},
				}}
			},
			`ALTER TABLE %s ALTER COLUMN "COLUMN_1" SET TAG %s = 'v1', %s = 'v2'`,
			id.FullyQualifiedName(), columnTagId1.FullyQualifiedName(), columnTagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_unsetTagsOnColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{UnsetTags: &TableUnsetColumnTags{
					Name:      "COLUMN_1",
					UnsetTags: []ObjectIdentifier{columnTagId1, columnTagId2},
				}}
			},
			`ALTER TABLE %s ALTER COLUMN "COLUMN_1" UNSET TAG %s, %s`,
			id.FullyQualifiedName(), columnTagId1.FullyQualifiedName(), columnTagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_dropColumns",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{DropColumns: &TableDropColumnAction{IfExists: new(true), Columns: []Column{{Value: "COLUMN_1"}, {Value: "COLUMN_2"}}}}
			},
			`ALTER TABLE %s DROP COLUMN IF EXISTS "COLUMN_1", "COLUMN_2"`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_ConstraintAction,
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Add: &TableOutOfLineConstraint{FK: &TableOutOfLineFK{
					Name:       new("OUT_OF_LINE_CONSTRAINT"),
					Columns:    []Column{{Value: "COLUMN_1"}, {Value: "COLUMN_2"}},
					References: fkRefId,
					RefColumns: []Column{{Value: "COLUMN_3"}, {Value: "COLUMN_4"}},
					Match:      new(MatchTypeFull),
					On: &ForeignKeyOnAction{
						OnUpdate: new(ForeignKeySetNullAction),
						OnDelete: new(ForeignKeyRestrictAction),
					},
				}}}
			},
			`ALTER TABLE %s ADD CONSTRAINT "OUT_OF_LINE_CONSTRAINT" FOREIGN KEY ("COLUMN_1", "COLUMN_2") REFERENCES %s ("COLUMN_3", "COLUMN_4") MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT`,
			id.FullyQualifiedName(), fkRefId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addPrimaryKeyConstraint",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Add: &TableOutOfLineConstraint{UniquePK: &TableOutOfLineUniquePK{
					Name:              new("pk_Constraint"),
					PrimaryKey:        new(true),
					Columns:           []Column{{Value: "COLUMN_1"}},
					Enforced:          new(true),
					Deferrable:        new(true),
					InitiallyDeferred: new(true),
					Enable:            new(true),
					Validate:          new(true),
					Rely:              new(true),
					Comment:           new("pk comment"),
				}}}
			},
			`ALTER TABLE %s ADD CONSTRAINT "pk_Constraint" PRIMARY KEY ("COLUMN_1") ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE VALIDATE RELY COMMENT 'pk comment'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addPrimaryKeyConstraintLowercaseColumns",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Add: &TableOutOfLineConstraint{UniquePK: &TableOutOfLineUniquePK{
					Name:       new("pk_id"),
					PrimaryKey: new(true),
					Columns:    []Column{{Value: "id"}},
				}}}
			},
			`ALTER TABLE %s ADD CONSTRAINT "pk_id" PRIMARY KEY ("id")`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addUniqueConstraint",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Add: &TableOutOfLineConstraint{UniquePK: &TableOutOfLineUniquePK{
					Unique:  new(true),
					Columns: []Column{{Value: "COLUMN_1"}, {Value: "COLUMN_2"}},
				}}}
			},
			`ALTER TABLE %s ADD UNIQUE ("COLUMN_1", "COLUMN_2")`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addCheckConstraint",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Add: &TableOutOfLineConstraint{CH: &TableOutOfLineCH{
					Name:           new("check_positive"),
					Expression:     "COLUMN_1 > 0",
					EnableValidate: new(true),
				}}}
			},
			`ALTER TABLE %s ADD CONSTRAINT "check_positive" CHECK ( COLUMN_1 > 0 ) ENABLE VALIDATE`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_renameConstraint",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Rename: &TableConstraintRenameAction{OldName: "OLD_NAME_CONSTRAINT", NewName: "NEW_NAME_CONSTRAINT"}}
			},
			`ALTER TABLE %s RENAME CONSTRAINT "OLD_NAME_CONSTRAINT" TO "NEW_NAME_CONSTRAINT"`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_alterConstraint",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Alter: &TableConstraintAlterAction{
					ConstraintName: new("OUT_OF_LINE_CONSTRAINT"),
					Columns:        []string{"COLUMN_3", "COLUMN_4"},
					NotEnforced:    new(true),
					Validate:       new(true),
					Rely:           new(true),
				}}
			},
			`ALTER TABLE %s ALTER CONSTRAINT "OUT_OF_LINE_CONSTRAINT" (COLUMN_3, COLUMN_4) NOT ENFORCED VALIDATE RELY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_dropConstraint",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Drop: &TableConstraintDropAction{
					ConstraintName: new("OUT_OF_LINE_CONSTRAINT"),
					Columns:        []string{"COLUMN_3", "COLUMN_4"},
					Cascade:        new(true),
				}}
			},
			`ALTER TABLE %s DROP CONSTRAINT "OUT_OF_LINE_CONSTRAINT" (COLUMN_3, COLUMN_4) CASCADE`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_ExternalTableAction,
			func(opts *AlterTableOptions) {
				opts.ExternalTableAction = &TableExternalTableAction{Add: &TableExternalTableColumnAddAction{
					IfNotExists: new(true),
					Name:        "COLUMN_1",
					ColumnType:  DataTypeBoolean,
					Expression:  []string{"SELECT 1"},
				}}
			},
			`ALTER TABLE %s ADD COLUMN IF NOT EXISTS COLUMN_1 BOOLEAN AS (SELECT 1)`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_SearchOptimizationAction,
			func(opts *AlterTableOptions) {
				opts.SearchOptimizationAction = &TableSearchOptimizationAction{
					Add: &TableAddSearchOptimization{
						On: []TableSearchMethodWithTarget{
							{Method: TableSearchMethodSubstring, Args: TableSearchMethodArgs{Targets: []string{"*"}}},
							{Method: TableSearchMethodEquality, Args: TableSearchMethodArgs{Targets: []string{"COLUMN_1"}}},
						},
					},
				}
			},
			`ALTER TABLE %s ADD SEARCH OPTIMIZATION ON SUBSTRING (*), EQUALITY (COLUMN_1)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_dropSearchOptimization",
			func(opts *AlterTableOptions) {
				opts.SearchOptimizationAction = &TableSearchOptimizationAction{Drop: &TableDropSearchOptimization{
					On: []TableDropSearchOptimizationOn{
						{ColumnName: new("SUBSTRING(*)")},
						{ColumnName: new("FOO")},
					},
				}}
			},
			`ALTER TABLE %s DROP SEARCH OPTIMIZATION ON SUBSTRING(*), FOO`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_Set,
			func(opts *AlterTableOptions) {
				opts.Set = &TableSet{
					EnableSchemaEvolution:      new(true),
					DataRetentionTimeInDays:    new(30),
					MaxDataExtensionTimeInDays: new(90),
					ChangeTracking:             new(false),
					DefaultDdlCollation:        new(StringAllowEmpty{Value: "us"}),
					IcebergDefaultDdlCollation: new(StringAllowEmpty{Value: "en-ci"}),
					ErrorLogging:               new(true),
					Contact:                    []TableContact{{Purpose: "SUPPORT", Contact: contactId}},
					Comment:                    &tableComment,
					RowTimestamp:               new(true),
				}
			},
			`ALTER TABLE %s SET ENABLE_SCHEMA_EVOLUTION = true DATA_RETENTION_TIME_IN_DAYS = 30 MAX_DATA_EXTENSION_TIME_IN_DAYS = 90 CHANGE_TRACKING = false DEFAULT_DDL_COLLATION = 'us' ICEBERG_DEFAULT_DDL_COLLATION = 'en-ci' ERROR_LOGGING = true CONTACT SUPPORT = %s COMMENT = '%s' ROW_TIMESTAMP = true`,
			id.FullyQualifiedName(), contactId.FullyQualifiedName(), tableComment,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_emptyDefaultDdlCollation",
			func(opts *AlterTableOptions) {
				opts.Set = &TableSet{DefaultDdlCollation: new(StringAllowEmpty{Value: ""})}
			},
			`ALTER TABLE %s SET DEFAULT_DDL_COLLATION = ''`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_SetTags,
			func(opts *AlterTableOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tableTagId1, Value: "v1"},
					{Name: tableTagId2, Value: "v2"},
				}
			},
			`ALTER TABLE %s SET TAG %s = 'v1', %s = 'v2'`,
			id.FullyQualifiedName(), tableTagId1.FullyQualifiedName(), tableTagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_UnsetTags,
			func(opts *AlterTableOptions) {
				opts.UnsetTags = []ObjectIdentifier{tableTagId1, tableTagId2}
			},
			`ALTER TABLE %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tableTagId1.FullyQualifiedName(), tableTagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_Unset,
			func(opts *AlterTableOptions) {
				opts.Unset = &TableUnset{
					DataRetentionTimeInDays:    new(true),
					MaxDataExtensionTimeInDays: new(true),
					ChangeTracking:             new(true),
					DefaultDdlCollation:        new(true),
					IcebergDefaultDdlCollation: new(true),
					EnableSchemaEvolution:      new(true),
					ErrorLogging:               new(true),
					Comment:                    new(true),
					RowTimestamp:               new(true),
				}
			},
			`ALTER TABLE %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, CHANGE_TRACKING, DEFAULT_DDL_COLLATION, ICEBERG_DEFAULT_DDL_COLLATION, ENABLE_SCHEMA_EVOLUTION, ERROR_LOGGING, COMMENT, ROW_TIMESTAMP`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_UnsetContact",
			func(opts *AlterTableOptions) {
				opts.Unset = &TableUnset{Contact: new("SUPPORT")}
			},
			`ALTER TABLE %s UNSET CONTACT SUPPORT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_AddRowAccessPolicy,
			func(opts *AlterTableOptions) {
				opts.AddRowAccessPolicy = &TableAddRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId, On: []Column{{"FIRST_COLUMN"}}}
			},
			`ALTER TABLE %s ADD ROW ACCESS POLICY %s ON ("FIRST_COLUMN")`,
			id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_DropRowAccessPolicy,
			func(opts *AlterTableOptions) {
				opts.DropRowAccessPolicy = &TableDropRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId}
			},
			`ALTER TABLE %s DROP ROW ACCESS POLICY %s`,
			id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_DropAndAddRowAccessPolicy,
			func(opts *AlterTableOptions) {
				opts.DropAndAddRowAccessPolicy = &TableDropAndAddRowAccessPolicy{
					Drop: TableDropRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId},
					Add:  TableAddRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId2, On: []Column{{"FIRST_COLUMN"}}},
				}
			},
			`ALTER TABLE %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON ("FIRST_COLUMN")`,
			id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), rowAccessPolicyId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_DropAllRowAccessPolicies,
			func(opts *AlterTableOptions) { opts.DropAllRowAccessPolicies = new(true) },
			`ALTER TABLE %s DROP ALL ROW ACCESS POLICIES`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_AddStorageLifecyclePolicy,
			func(opts *AlterTableOptions) {
				opts.AddStorageLifecyclePolicy = &TableStorageLifecyclePolicy{
					StorageLifecyclePolicy: storageLifecyclePolicyId,
					On:                     []Column{{Value: "FIRST_COLUMN"}},
				}
			},
			`ALTER TABLE %s ADD STORAGE LIFECYCLE POLICY %s ON ("FIRST_COLUMN")`,
			id.FullyQualifiedName(), storageLifecyclePolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addStorageLifecyclePolicy_multipleColumns",
			func(opts *AlterTableOptions) {
				opts.AddStorageLifecyclePolicy = &TableStorageLifecyclePolicy{
					StorageLifecyclePolicy: storageLifecyclePolicyId,
					On:                     []Column{{Value: "FIRST_COLUMN"}, {Value: "SECOND_COLUMN"}},
				}
			},
			`ALTER TABLE %s ADD STORAGE LIFECYCLE POLICY %s ON ("FIRST_COLUMN", "SECOND_COLUMN")`,
			id.FullyQualifiedName(), storageLifecyclePolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_DropStorageLifecyclePolicy,
			func(opts *AlterTableOptions) { opts.DropStorageLifecyclePolicy = new(true) },
			`ALTER TABLE %s DROP STORAGE LIFECYCLE POLICY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_SetAggregationPolicy,
			func(opts *AlterTableOptions) {
				opts.SetAggregationPolicy = &TableSetAggregationPolicy{
					AggregationPolicy: aggregationPolicyId,
					EntityKey:         []Column{{Value: "FIRST_COLUMN"}},
					Force:             new(true),
				}
			},
			`ALTER TABLE %s SET AGGREGATION POLICY %s ENTITY KEY ("FIRST_COLUMN") FORCE`,
			id.FullyQualifiedName(), aggregationPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_UnsetAggregationPolicy,
			func(opts *AlterTableOptions) { opts.UnsetAggregationPolicy = &TableUnsetAggregationPolicy{} },
			`ALTER TABLE %s UNSET AGGREGATION POLICY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_SetJoinPolicy,
			func(opts *AlterTableOptions) {
				opts.SetJoinPolicy = &TableSetJoinPolicy{JoinPolicy: joinPolicyId, Force: new(true)}
			},
			`ALTER TABLE %s SET JOIN POLICY %s FORCE`,
			id.FullyQualifiedName(), joinPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_UnsetJoinPolicy,
			func(opts *AlterTableOptions) { opts.UnsetJoinPolicy = &TableUnsetJoinPolicy{} },
			`ALTER TABLE %s UNSET JOIN POLICY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_SetDataMetricSchedule,
			func(opts *AlterTableOptions) {
				opts.SetDataMetricSchedule = &TableSetDataMetricSchedule{DataMetricSchedule: "TRIGGER_ON_CHANGES"}
			},
			`ALTER TABLE %s SET DATA_METRIC_SCHEDULE = 'TRIGGER_ON_CHANGES'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_UnsetDataMetricSchedule,
			func(opts *AlterTableOptions) { opts.UnsetDataMetricSchedule = new(true) },
			`ALTER TABLE %s UNSET DATA_METRIC_SCHEDULE`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_AddDataMetricFunction,
			func(opts *AlterTableOptions) {
				opts.AddDataMetricFunction = &TableAddDataMetricFunction{
					AddedFunctions: []TableAddedFunctions{{
						DataMetricFunction: dataMetricFunctionId,
						On:                 []Column{{Value: "ID"}},
						WithinGroup:        []Column{{Value: "REGION"}},
						GroupLimit:         new(10),
						Filter:             new("(ID > 0)"),
						Expectations:       []TableExpectations{{Name: "nonzero", Expression: "(COUNT(*) > 0)"}},
						ExecuteAsRole:      &executeAsRole,
						AnomalyDetection:   new(true),
						Sensitivity:        new(TableDataMetricSensitivityHigh),
					}},
				}
			},
			`ALTER TABLE %s ADD DATA METRIC FUNCTION %s ON ("ID") WITHIN GROUP ("REGION") GROUP LIMIT 10 FILTER (ID > 0) EXPECTATION "nonzero" (COUNT(*) > 0) EXECUTE AS ROLE %s ANOMALY_DETECTION = true SENSITIVITY = 'HIGH'`,
			id.FullyQualifiedName(), dataMetricFunctionId.FullyQualifiedName(), executeAsRole.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_DropDataMetricFunction,
			func(opts *AlterTableOptions) {
				opts.DropDataMetricFunction = &TableDropDataMetricFunction{
					DroppedFunctions: []TableDroppedFunctions{{
						DataMetricFunction: dataMetricFunctionId,
						On:                 []Column{{Value: "ID"}},
					}},
				}
			},
			`ALTER TABLE %s DROP DATA METRIC FUNCTION %s ON ("ID")`,
			id.FullyQualifiedName(), dataMetricFunctionId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_ModifyDataMetricFunction,
			func(opts *AlterTableOptions) {
				opts.ModifyDataMetricFunction = &TableModifyDataMetricFunction{
					ModifiedFunctions: []TableModifiedFunctions{{
						DataMetricFunction: dataMetricFunctionId,
						On:                 []Column{{Value: "ID"}},
						Status:             new(TableDataMetricScheduleStatusOperationSuspend),
					}},
				}
			},
			`ALTER TABLE %s MODIFY DATA METRIC FUNCTION %s ON ("ID") SUSPEND`,
			id.FullyQualifiedName(), dataMetricFunctionId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_modifyDataMetricFunction_addExpectation",
			func(opts *AlterTableOptions) {
				opts.ModifyDataMetricFunction = &TableModifyDataMetricFunction{
					ModifiedFunctions: []TableModifiedFunctions{{
						DataMetricFunction: dataMetricFunctionId,
						On:                 []Column{{Value: "ID"}},
						AddExpectation:     []TableAddExpectation{{Name: "nonzero", Expression: "(COUNT(*) > 0)"}},
					}},
				}
			},
			`ALTER TABLE %s MODIFY DATA METRIC FUNCTION %s ON ("ID") ADD EXPECTATION "nonzero" (COUNT(*) > 0)`,
			id.FullyQualifiedName(), dataMetricFunctionId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_modifyDataMetricFunction_set",
			func(opts *AlterTableOptions) {
				opts.ModifyDataMetricFunction = &TableModifyDataMetricFunction{
					ModifiedFunctions: []TableModifiedFunctions{{
						DataMetricFunction: dataMetricFunctionId,
						On:                 []Column{{Value: "ID"}},
						Set: &TableModifyDataMetricFunctionSet{
							AnomalyDetection:        new(true),
							Sensitivity:             new(TableDataMetricSensitivityMedium),
							DataQualityNotification: new(true),
						},
					}},
				}
			},
			`ALTER TABLE %s MODIFY DATA METRIC FUNCTION %s ON ("ID") SET ANOMALY_DETECTION = true SENSITIVITY = 'MEDIUM' DATA_QUALITY_NOTIFICATION = true`,
			id.FullyQualifiedName(), dataMetricFunctionId.FullyQualifiedName(),
		)

	tablesTests.Drop.
		withExpectedSqlf(
			case_Tables_sql_Drop_basic,
			`DROP TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Drop_all,
			func(opts *DropTableOptions) {
				opts.IfExists = new(true)
				opts.Cascade = new(true)
			},
			`DROP TABLE IF EXISTS %s CASCADE`, id.FullyQualifiedName(),
		)

	tablesTests.Show.
		withExpectedSql(case_Tables_sql_Show_basic, `SHOW TABLES`).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Show_all,
			func(opts *ShowTableOptions) {
				opts.Terse = new(true)
				opts.History = new(true)
				opts.Like = &Like{Pattern: &likePattern}
				opts.In = &ExtendedIn{Database: showDatabaseId}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			`SHOW TERSE TABLES HISTORY LIKE '%s' IN DATABASE %s STARTS WITH 'prefix' LIMIT 10`,
			likePattern, showDatabaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Show_Like,
			func(opts *ShowTableOptions) { opts.Like = &Like{Pattern: &likePattern} },
			`SHOW TABLES LIKE '%s'`, likePattern,
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Show_In,
			func(opts *ShowTableOptions) { opts.In = &ExtendedIn{Database: showDatabaseId} },
			`SHOW TABLES IN DATABASE %s`, showDatabaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Show_StartsWith,
			func(opts *ShowTableOptions) { opts.StartsWith = new("prefix") },
			`SHOW TABLES STARTS WITH 'prefix'`,
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Show_Limit,
			func(opts *ShowTableOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			`SHOW TABLES LIMIT 10`,
		).
		withAdditionalValidationCase(
			"validation_Show_emptyLike",
			func(opts *ShowTableOptions) { opts.Like = &Like{} },
			ErrPatternRequiredForLikeKeyword,
		)

	// TODO [SNOW-4532966]: add CreateFromBackupSet unit tests once the CREATE TABLE ... FROM BACKUP SET variant is added back to the SSOT definition.

	tablesTests.DescribeColumns.
		withExpectedSqlf(
			case_Tables_sql_DescribeColumns_basic,
			`DESCRIBE TABLE %s TYPE = COLUMNS`, id.FullyQualifiedName(),
		)

	tablesTests.DescribeStage.
		withExpectedSqlf(
			case_Tables_sql_DescribeStage_basic,
			`DESCRIBE TABLE %s TYPE = STAGE`, id.FullyQualifiedName(),
		)

	tablesTests.DescribeSearchOptimization.
		withExpectedSqlf(
			case_Tables_sql_DescribeSearchOptimization_basic,
			"DESCRIBE SEARCH OPTIMIZATION ON %s", id.FullyQualifiedName(),
		)

	tablesTests.SelectTableConstraints.
		withDefaultOpts(func() *SelectTableConstraintsTableOptions {
			return &SelectTableConstraintsTableOptions{
				Database:    databaseId,
				TableSchema: id.SchemaName(),
				TableName:   id.Name(),
			}
		}).
		withExpectedSqlf(
			case_Tables_sql_SelectTableConstraints_basic,
			"SELECT * FROM %s . INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'",
			databaseId.FullyQualifiedName(), id.SchemaName(), id.Name(),
		)

	tablesTests.SelectCheckConstraints.
		withDefaultOpts(func() *SelectCheckConstraintsTableOptions {
			return &SelectCheckConstraintsTableOptions{
				Database:         databaseId,
				ConstraintSchema: id.SchemaName(),
				ConstraintTable:  id.Name(),
			}
		}).
		withExpectedSqlf(
			case_Tables_sql_SelectCheckConstraints_basic,
			"SELECT * FROM %s . INFORMATION_SCHEMA.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = '%s' AND CONSTRAINT_TABLE = '%s'",
			databaseId.FullyQualifiedName(), id.SchemaName(), id.Name(),
		)
}

func tableTestColumn() TableColumn {
	return TableColumn{Name: "FIRST_COLUMN", ColumnType: DataTypeVARCHAR}
}

func TestTableColumnDetailsRow_SplitTypeAndCollation(t *testing.T) {
	t.Run("with utf8", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: "VARCHAR(10) COLLATE 'utf8'"}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "utf8", *actualCollation)
	})
	t.Run("with locale", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: "VARCHAR(10) COLLATE 'en_US'"}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "en_US", *actualCollation)
	})
	t.Run("with multiple specifiers", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: "VARCHAR(10) COLLATE 'fr_CA-ai-pi-trim'"}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "fr_CA-ai-pi-trim", *actualCollation)
	})
	t.Run("with empty collation", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: "VARCHAR(10) COLLATE ''"}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "", *actualCollation)
	})
	t.Run("without collation", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: "NUMBER(38, 0)"}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("NUMBER(38, 0)"), actualType)
		assert.Nil(t, actualCollation)
	})
}

func TestTable_GetClusterByKeys(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, (&Table{ClusterBy: ""}).GetClusterByKeys())
	})
	t.Run("one param", func(t *testing.T) {
		assert.Equal(t, []string{"abc"}, (&Table{ClusterBy: "LINEAR(abc)"}).GetClusterByKeys())
	})
	t.Run("more params", func(t *testing.T) {
		assert.Equal(t, []string{"abc", "def"}, (&Table{ClusterBy: "LINEAR(abc,def)"}).GetClusterByKeys())
	})
	t.Run("white space", func(t *testing.T) {
		assert.Equal(t, []string{"abc", "def"}, (&Table{ClusterBy: "   LINEAR(  abc  , def )"}).GetClusterByKeys())
	})
	t.Run("with function with one param", func(t *testing.T) {
		assert.Equal(t, []string{"some_func(some_param)", "other_param"}, (&Table{ClusterBy: "LINEAR(some_func(some_param),other_param)"}).GetClusterByKeys())
	})
	t.Run("with function with more than one param", func(t *testing.T) {
		assert.Equal(t, []string{"date_trunc('HOUR',TIMESTAMP_HR)", "CLIENT_ID"}, (&Table{ClusterBy: "LINEAR(date_trunc('HOUR',TIMESTAMP_HR),CLIENT_ID)"}).GetClusterByKeys())
	})
	t.Run("with nested functions", func(t *testing.T) {
		assert.Equal(t, []string{"some_func(some_param, some_other_param, other_func(some_param, some_other_param))", "other_param"}, (&Table{ClusterBy: "LINEAR(some_func(some_param, some_other_param, other_func(some_param, some_other_param)),other_param)"}).GetClusterByKeys())
	})
}

func Test_Table_convertIsDynamic(t *testing.T) {
	convert := func(t *testing.T, value sql.NullString) *Table {
		t.Helper()
		table, err := (tablesRow{IsDynamic: value}).convert()
		assert.NoError(t, err)
		return table
	}

	t.Run("Y", func(t *testing.T) {
		assert.True(t, convert(t, sql.NullString{String: "Y", Valid: true}).IsDynamic)
	})
	t.Run("N", func(t *testing.T) {
		assert.False(t, convert(t, sql.NullString{String: "N", Valid: true}).IsDynamic)
	})
	t.Run("other text", func(t *testing.T) {
		assert.False(t, convert(t, sql.NullString{String: "YES", Valid: true}).IsDynamic)
	})
	t.Run("unset", func(t *testing.T) {
		assert.False(t, convert(t, sql.NullString{}).IsDynamic)
	})
}
