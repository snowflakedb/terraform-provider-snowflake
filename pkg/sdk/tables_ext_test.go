package sdk

import (
	"fmt"
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
	storageLifecyclePolicyId := randomSchemaObjectIdentifier()
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
		withModify(case_Tables_validation_Create_opts_StageFileFormat_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateTableOptions) {
			opts.StageFileFormat = &LegacyFileFormat{
				FormatName:     new("fmt"),
				FileFormatType: new(FileFormatTypeCsv),
			}
		}).
		withAdditionalValidationCase(
			"validation_Create_noColumns",
			func(opts *CreateTableOptions) { opts.ColumnsAndConstraints.Columns = nil },
			errNotSet("CreateTableOptions", "Columns"),
		).
		withAdditionalValidationCase(
			"validation_Create_DefaultValue_ExactlyOneOf_Expression_Identity",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:       "FIRST_COLUMN",
					ColumnType: DataTypeVARCHAR,
					DefaultValue: &ColumnDefaultValue{
						Expression: new("1"),
						Identity:   &ColumnIdentity{Start: 1, Increment: 1},
					},
				}}
			},
			errExactlyOneOf("DefaultValue", "Expression", "Identity"),
		).
		withAdditionalValidationCase(
			"validation_Create_Identity_Order_Noorder",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:       "FIRST_COLUMN",
					ColumnType: DataTypeVARCHAR,
					DefaultValue: &ColumnDefaultValue{
						Identity: &ColumnIdentity{Start: 1, Increment: 1, Order: new(true), Noorder: new(true)},
					},
				}}
			},
			errMoreThanOneOf("Identity", "Order", "Noorder"),
		).
		withAdditionalValidationCase(
			"validation_Create_MaskingPolicy_invalidIdentifier",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:          "FIRST_COLUMN",
					ColumnType:    DataTypeVARCHAR,
					MaskingPolicy: &ColumnMaskingPolicy{Name: emptySchemaObjectIdentifier},
				}}
			},
			errInvalidIdentifier("ColumnMaskingPolicy", "Name"),
		).
		withAdditionalValidationCase(
			"validation_Create_columnTag_invalidIdentifier",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:       "FIRST_COLUMN",
					ColumnType: DataTypeVARCHAR,
					Tag:        []TagAssociation{{Name: emptySchemaObjectIdentifier, Value: "v"}},
				}}
			},
			errInvalidIdentifier("TagAssociation", "Name"),
		).
		withAdditionalValidationCase(
			"validation_Create_RowAccessPolicy_invalidIdentifier",
			func(opts *CreateTableOptions) {
				opts.RowAccessPolicy = &TableRowAccessPolicyLegacy{Name: emptySchemaObjectIdentifier, On: []string{"COLUMN_1"}}
			},
			errInvalidIdentifier("TableRowAccessPolicy", "Name"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_emptyType",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Name: new("INLINE_CONSTRAINT")},
				}}
			},
			errInvalidValue("ColumnInlineConstraint", "Type", ""),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_invalidType",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: "not existing type"},
				}}
			},
			errInvalidValue("ColumnInlineConstraint", "Type", "not existing type"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_foreignKeyNotSet",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeForeignKey},
				}}
			},
			errNotSet("ColumnInlineConstraint", "ForeignKey"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_foreignKeyTableNameNotSet",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeForeignKey, ForeignKey: &InlineForeignKey{}},
				}}
			},
			errNotSet("InlineForeignKey", "TableName"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_foreignKeySet",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeUnique, ForeignKey: &InlineForeignKey{TableName: "t"}},
				}}
			},
			errSet("ColumnInlineConstraint", "ForeignKey"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_Enforced_NotEnforced",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeUnique, Enforced: new(true), NotEnforced: new(true)},
				}}
			},
			errMoreThanOneOf("ColumnInlineConstraint", "Enforced", "NotEnforced"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_Deferrable_NotDeferrable",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeUnique, Deferrable: new(true), NotDeferrable: new(true)},
				}}
			},
			errMoreThanOneOf("ColumnInlineConstraint", "Deferrable", "NotDeferrable"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_InitiallyDeferred_InitiallyImmediate",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeUnique, InitiallyDeferred: new(true), InitiallyImmediate: new(true)},
				}}
			},
			errMoreThanOneOf("ColumnInlineConstraint", "InitiallyDeferred", "InitiallyImmediate"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_Enable_Disable",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeUnique, Enable: new(true), Disable: new(true)},
				}}
			},
			errMoreThanOneOf("ColumnInlineConstraint", "Enable", "Disable"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_Validate_Novalidate",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeUnique, Validate: new(true), NoValidate: new(true)},
				}}
			},
			errMoreThanOneOf("ColumnInlineConstraint", "Validate", "Novalidate"),
		).
		withAdditionalValidationCase(
			"validation_Create_inlineConstraint_Rely_Norely",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.Columns = []TableColumn{{
					Name:             "FIRST_COLUMN",
					ColumnType:       DataTypeVARCHAR,
					InlineConstraint: &ColumnInlineConstraint{Type: ColumnConstraintTypeUnique, Rely: new(true), NoRely: new(true)},
				}}
			},
			errMoreThanOneOf("ColumnInlineConstraint", "Rely", "Norely"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_noColumns",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: ColumnConstraintTypeUnique}}
			},
			errNotSet("OutOfLineConstraint", "Columns"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_emptyType",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{Name: new("OUT_OF_LINE_CONSTRAINT"), Columns: []string{"COLUMN_1"}}}
			},
			errInvalidValue("OutOfLineConstraint", "ConstraintType", ""),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_invalidType",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: "not existing type", Columns: []string{"COLUMN_1"}}}
			},
			errInvalidValue("OutOfLineConstraint", "ConstraintType", "not existing type"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_foreignKeyNotSet",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: ColumnConstraintTypeForeignKey, Columns: []string{"COLUMN_1"}}}
			},
			errNotSet("OutOfLineConstraint", "ForeignKey"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_foreignKeyTableNameNotSet",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{
					ConstraintType: ColumnConstraintTypeForeignKey,
					Columns:        []string{"COLUMN_1"},
					ForeignKey:     &OutOfLineForeignKey{},
				}}
			},
			errNotSet("OutOfLineForeignKey", "TableName"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_foreignKeySet",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{
					ConstraintType: ColumnConstraintTypeUnique,
					Columns:        []string{"COLUMN_1"},
					ForeignKey:     &OutOfLineForeignKey{TableName: fkRefId},
				}}
			},
			errSet("OutOfLineConstraint", "ForeignKey"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_Enforced_NotEnforced",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: ColumnConstraintTypeUnique, Columns: []string{"COLUMN_1"}, Enforced: new(true), NotEnforced: new(true)}}
			},
			errMoreThanOneOf("OutOfLineConstraint", "Enforced", "NotEnforced"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_Deferrable_NotDeferrable",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: ColumnConstraintTypeUnique, Columns: []string{"COLUMN_1"}, Deferrable: new(true), NotDeferrable: new(true)}}
			},
			errMoreThanOneOf("OutOfLineConstraint", "Deferrable", "NotDeferrable"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_InitiallyDeferred_InitiallyImmediate",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: ColumnConstraintTypeUnique, Columns: []string{"COLUMN_1"}, InitiallyDeferred: new(true), InitiallyImmediate: new(true)}}
			},
			errMoreThanOneOf("OutOfLineConstraint", "InitiallyDeferred", "InitiallyImmediate"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_Enable_Disable",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: ColumnConstraintTypeUnique, Columns: []string{"COLUMN_1"}, Enable: new(true), Disable: new(true)}}
			},
			errMoreThanOneOf("OutOfLineConstraint", "Enable", "Disable"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_Validate_Novalidate",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: ColumnConstraintTypeUnique, Columns: []string{"COLUMN_1"}, Validate: new(true), Novalidate: new(true)}}
			},
			errMoreThanOneOf("OutOfLineConstraint", "Validate", "Novalidate"),
		).
		withAdditionalValidationCase(
			"validation_Create_outOfLineConstraint_Rely_Norely",
			func(opts *CreateTableOptions) {
				opts.ColumnsAndConstraints.OutOfLineConstraint = []OutOfLineConstraint{{ConstraintType: ColumnConstraintTypeUnique, Columns: []string{"COLUMN_1"}, Rely: new(true), Norely: new(true)}}
			},
			errMoreThanOneOf("OutOfLineConstraint", "Rely", "Norely"),
		).
		withExpectedSqlf(
			case_Tables_sql_Create_basic,
			`CREATE TABLE %s (FIRST_COLUMN VARCHAR)`, id.FullyQualifiedName(),
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
						MaskingPolicy: &ColumnMaskingPolicy{
							Name:  maskingPolicyId,
							Using: []string{"FOO", "BAR"},
						},
						Tag: []TagAssociation{
							{Name: columnTagId1, Value: "v1"},
							{Name: columnTagId2, Value: "v2"},
						},
						InlineConstraint: &ColumnInlineConstraint{Name: new("INLINE_CONSTRAINT"), Type: ColumnConstraintTypePrimaryKey},
					}},
					OutOfLineConstraint: []OutOfLineConstraint{
						{
							Name:           new("OUT_OF_LINE_CONSTRAINT"),
							ConstraintType: ColumnConstraintTypeForeignKey,
							Columns:        []string{"COLUMN_1", "COLUMN_2"},
							ForeignKey: &OutOfLineForeignKey{
								TableName:   fkRefId,
								ColumnNames: []string{"COLUMN_3", "COLUMN_4"},
								Match:       new(MatchTypeFull),
								On: &ForeignKeyOnAction{
									OnUpdate: new(ForeignKeySetNullAction),
									OnDelete: new(ForeignKeyRestrictAction),
								},
							},
						},
						{
							ConstraintType:    ColumnConstraintTypeUnique,
							Columns:           []string{"COLUMN_1"},
							Enforced:          new(true),
							Deferrable:        new(true),
							InitiallyDeferred: new(true),
							Enable:            new(true),
							Rely:              new(true),
						},
					},
				}
				opts.ClusterBy = []string{"COLUMN_1", "COLUMN_2"}
				opts.EnableSchemaEvolution = new(true)
				opts.StageFileFormat = &LegacyFileFormat{
					FileFormatType: new(FileFormatTypeCsv),
					Options:        &FileFormatTypeOptionsLegacy{CSVCompression: new(CsvCompressionAuto)},
				}
				opts.StageCopyOptions = &LegacyTableCopyOptions{OnError: &LegacyTableCopyOnErrorOptions{SkipFile: new("SKIP_FILE")}}
				opts.DataRetentionTimeInDays = new(10)
				opts.MaxDataExtensionTimeInDays = new(100)
				opts.ChangeTracking = new(true)
				opts.DefaultDdlCollation = new("en")
				opts.RowAccessPolicy = &TableRowAccessPolicyLegacy{Name: rowAccessPolicyId, On: []string{"COLUMN_1", "COLUMN_2"}}
				opts.Tag = []TagAssociation{
					{Name: tableTagId1, Value: "v1"},
					{Name: tableTagId2, Value: "v2"},
				}
				opts.Comment = &tableComment
			},
			`CREATE TABLE IF NOT EXISTS %s (FIRST_COLUMN VARCHAR CONSTRAINT INLINE_CONSTRAINT PRIMARY KEY NOT NULL COLLATE 'de' IDENTITY START 10 INCREMENT 1 ORDER MASKING POLICY %s USING (FOO, BAR) TAG (%s = 'v1', %s = 'v2') COMMENT '%s', CONSTRAINT OUT_OF_LINE_CONSTRAINT FOREIGN KEY (COLUMN_1, COLUMN_2) REFERENCES %s (COLUMN_3, COLUMN_4) MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT, UNIQUE (COLUMN_1) ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE RELY) CLUSTER BY (COLUMN_1, COLUMN_2) ENABLE_SCHEMA_EVOLUTION = true STAGE_FILE_FORMAT = (TYPE = CSV COMPRESSION = AUTO) STAGE_COPY_OPTIONS = (ON_ERROR = SKIP_FILE) DATA_RETENTION_TIME_IN_DAYS = 10 MAX_DATA_EXTENSION_TIME_IN_DAYS = 100 CHANGE_TRACKING = true DEFAULT_DDL_COLLATION = 'en' ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) TAG (%s = 'v1', %s = 'v2') COMMENT = '%s'`,
			id.FullyQualifiedName(),
			maskingPolicyId.FullyQualifiedName(),
			columnTagId1.FullyQualifiedName(),
			columnTagId2.FullyQualifiedName(),
			columnComment,
			fkRefId.FullyQualifiedName(),
			rowAccessPolicyId.FullyQualifiedName(),
			tableTagId1.FullyQualifiedName(),
			tableTagId2.FullyQualifiedName(),
			tableComment,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateTableOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE TABLE %s (FIRST_COLUMN VARCHAR) COPY GRANTS`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_skipFileX",
			func(opts *CreateTableOptions) {
				opts.StageCopyOptions = &LegacyTableCopyOptions{OnError: &LegacyTableCopyOnErrorOptions{SkipFile: new(fmt.Sprintf("SKIP_FILE_%d", 5))}}
			},
			`CREATE TABLE %s (FIRST_COLUMN VARCHAR) STAGE_COPY_OPTIONS = (ON_ERROR = SKIP_FILE_5)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_skipFileXPercent",
			func(opts *CreateTableOptions) {
				opts.StageCopyOptions = &LegacyTableCopyOptions{OnError: &LegacyTableCopyOnErrorOptions{SkipFile: new(fmt.Sprintf("'SKIP_FILE_%d%%'", 10))}}
			},
			`CREATE TABLE %s (FIRST_COLUMN VARCHAR) STAGE_COPY_OPTIONS = (ON_ERROR = 'SKIP_FILE_10%%')`,
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
			`CREATE TABLE %s (a) AS SELECT 1`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_CreateAsSelect_all,
			func(opts *CreateAsSelectTableOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
				opts.Columns = []TableAsSelectColumn{{
					Name:          "FIRST_COLUMN",
					ColumnType:    new(DataTypeVARCHAR),
					MaskingPolicy: &maskingPolicyId,
				}}
				opts.ClusterBy = []string{"COLUMN_1", "COLUMN_2"}
				opts.RowAccessPolicy = &TableRowAccessPolicyLegacy{Name: rowAccessPolicyId, On: []string{"COLUMN_1", "COLUMN_2"}}
				opts.Query = "SELECT * FROM ANOTHER_TABLE"
			},
			`CREATE OR REPLACE TABLE %s (FIRST_COLUMN VARCHAR MASKING POLICY %s) CLUSTER BY (COLUMN_1, COLUMN_2) COPY GRANTS ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) AS SELECT * FROM ANOTHER_TABLE`,
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
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE TABLE %s COPY GRANTS USING TEMPLATE (sample_data)`,
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
				opts.CopyGrants = new(true)
				opts.ClusterBy = []string{"date", "id"}
			},
			`CREATE OR REPLACE TABLE %s LIKE %s CLUSTER BY (date, id) COPY GRANTS`,
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
				opts.CopyGrants = new(true)
				opts.ClonePoint = &ClonePoint{
					Moment: CloneMomentAt,
					At:     TimeTravel{Offset: new(0)},
				}
			},
			`CREATE OR REPLACE TABLE %s CLONE %s AT (OFFSET => 0) COPY GRANTS`,
			id.FullyQualifiedName(), sourceTable.FullyQualifiedName(),
		)

	tablesTests.Alter.
		withModify(case_Tables_validation_Alter_opts_ClusteringAction_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterTableOptions) {
			opts.ClusteringAction = &TableClusteringAction{ClusterBy: []string{"a"}, DropClusteringKey: new(true)}
		}).
		withModify(case_Tables_validation_Alter_opts_SearchOptimizationAction_Drop_On_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterTableOptions) {
			opts.SearchOptimizationAction = &TableSearchOptimizationActionLegacy{Drop: &TableDropSearchOptimization{
				On: []TableDropSearchOptimizationOn{{ColumnName: new("c"), ExpressionId: new("e")}},
			}}
		}).
		withModify(case_Tables_validation_Alter_opts_SearchOptimizationAction_Drop_On_ExactlyOneValueSet_OneValidOneInvalid, func(opts *AlterTableOptions) {
			opts.SearchOptimizationAction = &TableSearchOptimizationActionLegacy{Drop: &TableDropSearchOptimization{
				On: []TableDropSearchOptimizationOn{{ColumnName: new("c")}, {}},
			}}
		}).
		withModify(case_Tables_validation_Alter_opts_Set_StageFileFormat_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterTableOptions) {
			opts.Set = &TableSet{StageFileFormat: &LegacyFileFormat{FormatName: new("fmt"), FileFormatType: new(FileFormatTypeCsv)}}
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
		withAdditionalValidationCase(
			"validation_Alter_ConstraintAction_Add_ColumnsNotSet",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Add: &OutOfLineConstraint{ConstraintType: ColumnConstraintTypeUnique}}
			},
			errNotSet("OutOfLineConstraint", "Columns"),
		).
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
			`ALTER TABLE %s ADD COLUMN IF NOT EXISTS NEXT_COLUMN VARCHAR COLLATE 'utf8' IDENTITY START 10 INCREMENT 1`, id.FullyQualifiedName(),
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
			`ALTER TABLE %s ADD COLUMN NEW_BOOLEAN_COLUMN BOOLEAN DEFAULT FALSE`, id.FullyQualifiedName(),
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
			`ALTER TABLE %s ADD COLUMN NEXT_COLUMN VARCHAR IDENTITY START 10 INCREMENT 1`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_renameColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{Rename: &TableColumnRenameAction{OldName: "OLD_NAME", NewName: "NEW_NAME"}}
			},
			`ALTER TABLE %s RENAME COLUMN OLD_NAME TO NEW_NAME`, id.FullyQualifiedName(),
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
					{Name: "COLUMN_2", DataType: new(DataTypeVARCHAR), Collate: new("utf8")},
					{Name: "COLUMN_2", NotNullConstraint: &TableColumnNotNullConstraint{Drop: new(true)}},
				}}
			},
			`ALTER TABLE %s ALTER COLUMN COLUMN_1 DROP DEFAULT, COLUMN COLUMN_1 SET DEFAULT SEQUENCE_1.NEXTVAL, COLUMN COLUMN_1 UNSET COMMENT, COLUMN COLUMN_2 DROP DEFAULT, COLUMN COLUMN_2 SET DEFAULT SEQUENCE_2.NEXTVAL, COLUMN COLUMN_2 COMMENT 'comment', COLUMN COLUMN_2 SET DATA TYPE VARCHAR COLLATE 'utf8', COLUMN COLUMN_2 DROP NOT NULL`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setMaskingPolicyOnColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{SetMaskingPolicy: &TableColumnAlterSetMaskingPolicyAction{
					ColumnName:        "COLUMN_1",
					MaskingPolicyName: maskingPolicyId,
					Using:             []string{"FOO", "BAR"},
					Force:             new(true),
				}}
			},
			`ALTER TABLE %s ALTER COLUMN COLUMN_1 SET MASKING POLICY %s USING (FOO, BAR) FORCE`,
			id.FullyQualifiedName(), maskingPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_unsetMaskingPolicyOnColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{UnsetMaskingPolicy: &TableColumnAlterUnsetMaskingPolicyAction{ColumnName: "COLUMN_1"}}
			},
			`ALTER TABLE %s ALTER COLUMN COLUMN_1 UNSET MASKING POLICY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setTagsOnColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{SetTags: &TableColumnAlterSetTagsAction{
					ColumnName: "COLUMN_1",
					SetTags: []TagAssociation{
						{Name: columnTagId1, Value: "v1"},
						{Name: columnTagId2, Value: "v2"},
					},
				}}
			},
			`ALTER TABLE %s ALTER COLUMN COLUMN_1 SET TAG %s = 'v1', %s = 'v2'`,
			id.FullyQualifiedName(), columnTagId1.FullyQualifiedName(), columnTagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_unsetTagsOnColumn",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{UnsetTags: &TableColumnAlterUnsetTagsAction{
					ColumnName: "COLUMN_1",
					UnsetTags:  []ObjectIdentifier{columnTagId1, columnTagId2},
				}}
			},
			`ALTER TABLE %s ALTER COLUMN COLUMN_1 UNSET TAG %s, %s`,
			id.FullyQualifiedName(), columnTagId1.FullyQualifiedName(), columnTagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_dropColumns",
			func(opts *AlterTableOptions) {
				opts.ColumnAction = &TableColumnAction{DropColumns: &TableColumnAlterDropColumns{IfExists: new(true), Columns: []string{"COLUMN_1", "COLUMN_2"}}}
			},
			`ALTER TABLE %s DROP COLUMN IF EXISTS COLUMN_1, COLUMN_2`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_ConstraintAction,
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Add: &OutOfLineConstraint{
					Name:           new("OUT_OF_LINE_CONSTRAINT"),
					ConstraintType: ColumnConstraintTypeForeignKey,
					Columns:        []string{"COLUMN_1", "COLUMN_2"},
					ForeignKey: &OutOfLineForeignKey{
						TableName:   fkRefId,
						ColumnNames: []string{"COLUMN_3", "COLUMN_4"},
						Match:       new(MatchTypeFull),
						On: &ForeignKeyOnAction{
							OnUpdate: new(ForeignKeySetNullAction),
							OnDelete: new(ForeignKeyRestrictAction),
						},
					},
				}}
			},
			`ALTER TABLE %s ADD CONSTRAINT OUT_OF_LINE_CONSTRAINT FOREIGN KEY (COLUMN_1, COLUMN_2) REFERENCES %s (COLUMN_3, COLUMN_4) MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT`,
			id.FullyQualifiedName(), fkRefId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_renameConstraint",
			func(opts *AlterTableOptions) {
				opts.ConstraintAction = &TableConstraintAction{Rename: &TableConstraintRenameAction{OldName: "OLD_NAME_CONSTRAINT", NewName: "NEW_NAME_CONSTRAINT"}}
			},
			`ALTER TABLE %s RENAME CONSTRAINT OLD_NAME_CONSTRAINT TO NEW_NAME_CONSTRAINT`, id.FullyQualifiedName(),
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
			`ALTER TABLE %s ALTER CONSTRAINT OUT_OF_LINE_CONSTRAINT (COLUMN_3, COLUMN_4) NOT ENFORCED VALIDATE RELY`, id.FullyQualifiedName(),
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
			`ALTER TABLE %s DROP CONSTRAINT OUT_OF_LINE_CONSTRAINT (COLUMN_3, COLUMN_4) CASCADE`, id.FullyQualifiedName(),
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
		withAdditionalSqlCasef(
			"sql_Alter_renameExternalTableColumn",
			func(opts *AlterTableOptions) {
				opts.ExternalTableAction = &TableExternalTableAction{Rename: &TableExternalTableColumnRenameAction{OldName: "OLD_NAME_COLUMN", NewName: "NEW_NAME_COLUMN"}}
			},
			`ALTER TABLE %s RENAME COLUMN OLD_NAME_COLUMN TO NEW_NAME_COLUMN`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_dropExternalTableColumn",
			func(opts *AlterTableOptions) {
				opts.ExternalTableAction = &TableExternalTableAction{Drop: &TableExternalTableColumnDropAction{IfExists: new(true), Names: []string{"COLUMN_3", "COLUMN_4"}}}
			},
			`ALTER TABLE %s DROP COLUMN IF EXISTS COLUMN_3, COLUMN_4`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_SearchOptimizationAction,
			func(opts *AlterTableOptions) {
				opts.SearchOptimizationAction = &TableSearchOptimizationActionLegacy{Add: &AddSearchOptimization{On: []string{"SUBSTRING(*)", "GEO(*)"}}}
			},
			`ALTER TABLE %s ADD SEARCH OPTIMIZATION ON SUBSTRING(*), GEO(*)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_dropSearchOptimization",
			func(opts *AlterTableOptions) {
				opts.SearchOptimizationAction = &TableSearchOptimizationActionLegacy{Drop: &TableDropSearchOptimization{
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
					StageFileFormat:            &LegacyFileFormat{FileFormatType: new(FileFormatTypeCsv)},
					StageCopyOptions:           &LegacyTableCopyOptions{OnError: &LegacyTableCopyOnErrorOptions{SkipFile: new("SKIP_FILE")}},
					DataRetentionTimeInDays:    new(30),
					MaxDataExtensionTimeInDays: new(90),
					ChangeTracking:             new(false),
					DefaultDdlCollation:        new("us"),
					Comment:                    &tableComment,
				}
			},
			`ALTER TABLE %s SET ENABLE_SCHEMA_EVOLUTION = true STAGE_FILE_FORMAT = (TYPE = CSV) STAGE_COPY_OPTIONS = (ON_ERROR = SKIP_FILE) DATA_RETENTION_TIME_IN_DAYS = 30 MAX_DATA_EXTENSION_TIME_IN_DAYS = 90 CHANGE_TRACKING = false DEFAULT_DDL_COLLATION = 'us' COMMENT = '%s'`,
			id.FullyQualifiedName(), tableComment,
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
					EnableSchemaEvolution:      new(true),
					Comment:                    new(true),
				}
			},
			`ALTER TABLE %s UNSET DATA_RETENTION_TIME_IN_DAYS MAX_DATA_EXTENSION_TIME_IN_DAYS CHANGE_TRACKING DEFAULT_DDL_COLLATION ENABLE_SCHEMA_EVOLUTION COMMENT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tables_sql_Alter_AddRowAccessPolicy,
			func(opts *AlterTableOptions) {
				opts.AddRowAccessPolicy = &TableAddRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId, On: []string{"FIRST_COLUMN"}}
			},
			`ALTER TABLE %s ADD ROW ACCESS POLICY %s ON (FIRST_COLUMN)`,
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
					Add:  TableAddRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId2, On: []string{"FIRST_COLUMN"}},
				}
			},
			`ALTER TABLE %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON (FIRST_COLUMN)`,
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
				opts.AddStorageLifecyclePolicy = &TableAddStorageLifecyclePolicy{
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
				opts.AddStorageLifecyclePolicy = &TableAddStorageLifecyclePolicy{
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
				opts.In = &ExtendedIn{In: In{Database: showDatabaseId}}
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
			func(opts *ShowTableOptions) { opts.In = &ExtendedIn{In: In{Database: showDatabaseId}} },
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
		row := tableColumnDetailsRow{Type: DataType("VARCHAR(10) COLLATE 'utf8'")}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "utf8", *actualCollation)
	})
	t.Run("with locale", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: DataType("VARCHAR(10) COLLATE 'en_US'")}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "en_US", *actualCollation)
	})
	t.Run("with multiple specifiers", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: DataType("VARCHAR(10) COLLATE 'fr_CA-ai-pi-trim'")}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "fr_CA-ai-pi-trim", *actualCollation)
	})
	t.Run("with empty collation", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: DataType("VARCHAR(10) COLLATE ''")}
		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "", *actualCollation)
	})
	t.Run("without collation", func(t *testing.T) {
		row := tableColumnDetailsRow{Type: DataType("NUMBER(38, 0)")}
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
