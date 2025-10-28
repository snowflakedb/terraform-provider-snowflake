package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	sampleColumnName := "FIRST_COLUMN"
	sampleColumnType := DataTypeVARCHAR

	defaultOpts := func() *createTableOptions {
		return &createTableOptions{
			name: id,
		}
	}

	defaultOptsWithColumnInlineConstraint := func(inlineConstraint *ColumnInlineConstraint) *createTableOptions {
		columns := []TableColumn{{
			Name:             sampleColumnName,
			Type:             sampleColumnType,
			InlineConstraint: inlineConstraint,
		}}
		return &createTableOptions{
			name:                  id,
			ColumnsAndConstraints: CreateTableColumnsAndConstraints{Columns: columns},
		}
	}

	defaultOptsWithColumnOutOfLineConstraint := func(outOfLineConstraint *OutOfLineConstraint) *createTableOptions {
		columns := []TableColumn{{
			Name: sampleColumnName,
			Type: sampleColumnType,
		}}
		return &createTableOptions{
			name:                  id,
			ColumnsAndConstraints: CreateTableColumnsAndConstraints{Columns: columns, OutOfLineConstraint: []OutOfLineConstraint{*outOfLineConstraint}},
		}
	}

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableOptions", "name"))
	})

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableOptions", "name"))
	})

	t.Run("validation: no columns", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("createTableOptions", "Columns"))
	})

	t.Run("validation: both expression and identity of a column are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = randomSchemaObjectIdentifier()
		opts.ColumnsAndConstraints = CreateTableColumnsAndConstraints{
			Columns: []TableColumn{{
				Name: "a",
				DefaultValue: &ColumnDefaultValue{
					Expression: String("expr"),
					Identity: &ColumnIdentity{
						Start:     10,
						Increment: 1,
					},
				},
			}},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DefaultValue", "Expression", "Identity"))
	})

	t.Run("validation: both order and noorder are present for identity", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = randomSchemaObjectIdentifier()
		opts.ColumnsAndConstraints = CreateTableColumnsAndConstraints{
			Columns: []TableColumn{{
				Name: "a",
				DefaultValue: &ColumnDefaultValue{
					Identity: &ColumnIdentity{
						Start:     10,
						Increment: 1,
						Order:     Bool(true),
						Noorder:   Bool(true),
					},
				},
			}},
		}
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("Identity", "Order", "Noorder"))
	})

	t.Run("validation: column masking policy incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = randomSchemaObjectIdentifier()
		opts.ColumnsAndConstraints = CreateTableColumnsAndConstraints{
			Columns: []TableColumn{{
				Name: "a",
				MaskingPolicy: &ColumnMaskingPolicy{
					Name: emptySchemaObjectIdentifier,
				},
			}},
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("ColumnMaskingPolicy", "Name"))
	})

	t.Run("validation: column tag association's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = randomSchemaObjectIdentifier()
		opts.ColumnsAndConstraints = CreateTableColumnsAndConstraints{
			Columns: []TableColumn{{
				Name: "a",
				Tags: []TagAssociation{
					{
						Name:  emptySchemaObjectIdentifier,
						Value: "v1",
					},
				},
			}},
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("TagAssociation", "Name"))
	})

	t.Run("validation: stageFileFormat's both format name and format type are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.StageFileFormat = &StageFileFormat{
			FormatName: String("some_format"),
			Type:       Pointer(FileFormatTypeCSV),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("StageFileFormat", "FormatName", "FormatType"))
	})

	t.Run("validation: rowAccessPolicy's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.RowAccessPolicy = &TableRowAccessPolicy{
			Name: emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("TableRowAccessPolicy", "Name"))
	})

	t.Run("validation: inline constraint - constraint name empty", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type: "",
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidValue("ColumnInlineConstraint", "Type", ""))
	})

	t.Run("validation: inline constraint - constraint ", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type: "not existing type",
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidValue("ColumnInlineConstraint", "Type", "not existing type"))
	})

	t.Run("validation: inline constraint - foreign key present for foreign key constraint", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type:       ColumnConstraintTypeForeignKey,
			ForeignKey: nil,
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("ColumnInlineConstraint", "ForeignKey"))
	})

	t.Run("validation: inline constraint - foreign key validation", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type: ColumnConstraintTypeForeignKey,
			ForeignKey: &InlineForeignKey{
				TableName: "",
			},
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("InlineForeignKey", "TableName"))
	})

	t.Run("validation: inline constraint - foreign key absent for constraint other than foreign key", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type: ColumnConstraintTypeUnique,
			ForeignKey: &InlineForeignKey{
				TableName: "table",
			},
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errSet("ColumnInlineConstraint", "ForeignKey"))
	})

	t.Run("validation: inline constraint - enforced and not enforced both present", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type:        ColumnConstraintTypeUnique,
			Enforced:    Bool(true),
			NotEnforced: Bool(true),
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("ColumnInlineConstraint", "Enforced", "NotEnforced"))
	})

	t.Run("validation: inline constraint - deferrable and not deferrable both present", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type:          ColumnConstraintTypeUnique,
			Deferrable:    Bool(true),
			NotDeferrable: Bool(true),
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("ColumnInlineConstraint", "Deferrable", "NotDeferrable"))
	})

	t.Run("validation: inline constraint - initially deferred and initially immediate both present", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type:               ColumnConstraintTypeUnique,
			InitiallyDeferred:  Bool(true),
			InitiallyImmediate: Bool(true),
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("ColumnInlineConstraint", "InitiallyDeferred", "InitiallyImmediate"))
	})

	t.Run("validation: inline constraint - enable and disable both present", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type:    ColumnConstraintTypeUnique,
			Enable:  Bool(true),
			Disable: Bool(true),
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("ColumnInlineConstraint", "Enable", "Disable"))
	})

	t.Run("validation: inline constraint - validate and novalidate both present", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type:       ColumnConstraintTypeUnique,
			Validate:   Bool(true),
			NoValidate: Bool(true),
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("ColumnInlineConstraint", "Validate", "Novalidate"))
	})

	t.Run("validation: inline constraint - rely and norely both present", func(t *testing.T) {
		inlineConstraint := ColumnInlineConstraint{
			Type:   ColumnConstraintTypeUnique,
			Rely:   Bool(true),
			NoRely: Bool(true),
		}
		opts := defaultOptsWithColumnInlineConstraint(&inlineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("ColumnInlineConstraint", "Rely", "Norely"))
	})

	t.Run("validation: out of line constraint - no columns", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type: ColumnConstraintTypeUnique,
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("OutOfLineConstraint", "Columns"))
	})

	t.Run("validation: out of line constraint - constraint name empty", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type: "",
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidValue("OutOfLineConstraint", "Type", ""))
	})

	t.Run("validation: out of line constraint - constraint ", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type: "not existing type",
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidValue("OutOfLineConstraint", "Type", "not existing type"))
	})

	t.Run("validation: out of line constraint - foreign key present for foreign key constraint", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type:       ColumnConstraintTypeForeignKey,
			ForeignKey: nil,
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("OutOfLineConstraint", "ForeignKey"))
	})

	t.Run("validation: out of line constraint - foreign key validation", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type: ColumnConstraintTypeForeignKey,
			ForeignKey: &OutOfLineForeignKey{
				TableName: emptySchemaObjectIdentifier,
			},
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("OutOfLineForeignKey", "TableName"))
	})

	t.Run("validation: out of line constraint - foreign key absent for constraint other than foreign key", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type: ColumnConstraintTypeUnique,
			ForeignKey: &OutOfLineForeignKey{
				TableName: randomSchemaObjectIdentifier(),
			},
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errSet("OutOfLineConstraint", "ForeignKey"))
	})

	t.Run("validation: out of line constraint - enforced and not enforced both present", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type:        ColumnConstraintTypeUnique,
			Enforced:    Bool(true),
			NotEnforced: Bool(true),
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("OutOfLineConstraint", "Enforced", "NotEnforced"))
	})

	t.Run("validation: out of line constraint - deferrable and not deferrable both present", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type:          ColumnConstraintTypeUnique,
			Deferrable:    Bool(true),
			NotDeferrable: Bool(true),
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("OutOfLineConstraint", "Deferrable", "NotDeferrable"))
	})

	t.Run("validation: out of line constraint - initially deferred and initially immediate both present", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type:               ColumnConstraintTypeUnique,
			InitiallyDeferred:  Bool(true),
			InitiallyImmediate: Bool(true),
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("OutOfLineConstraint", "InitiallyDeferred", "InitiallyImmediate"))
	})

	t.Run("validation: out of line constraint - enable and disable both present", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type:    ColumnConstraintTypeUnique,
			Enable:  Bool(true),
			Disable: Bool(true),
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("OutOfLineConstraint", "Enable", "Disable"))
	})

	t.Run("validation: out of line constraint - validate and novalidate both present", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type:       ColumnConstraintTypeUnique,
			Validate:   Bool(true),
			NoValidate: Bool(true),
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("OutOfLineConstraint", "Validate", "Novalidate"))
	})

	t.Run("validation: out of line constraint - rely and norely both present", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type:   ColumnConstraintTypeUnique,
			Rely:   Bool(true),
			NoRely: Bool(true),
		}
		opts := defaultOptsWithColumnOutOfLineConstraint(&outOfLineConstraint)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("OutOfLineConstraint", "Rely", "Norely"))
	})

	t.Run("with complete options", func(t *testing.T) {
		columnTagId1 := randomSchemaObjectIdentifier()
		columnTagId2 := randomSchemaObjectIdentifierInSchema(columnTagId1.SchemaId())
		tableTagId1 := randomSchemaObjectIdentifierInSchema(columnTagId1.SchemaId())
		tableTagId2 := randomSchemaObjectIdentifierInSchema(columnTagId1.SchemaId())
		columnComment := random.Comment()
		tableComment := random.Comment()
		collation := "de"
		columnName := "FIRST_COLUMN"
		columnTypeRaw, err := datatypes.ParseDataType("VARCHAR")
		require.NoError(t, err)
		columnType := LegacyDataTypeFrom(columnTypeRaw)
		maskingPolicy := ColumnMaskingPolicy{
			Name:  randomSchemaObjectIdentifier(),
			Using: []string{"FOO", "BAR"},
		}
		columnTags := []TagAssociation{
			{
				Name:  columnTagId1,
				Value: "v1",
			},
			{
				Name:  columnTagId2,
				Value: "v2",
			},
		}

		tableTags := []TagAssociation{
			{
				Name:  tableTagId1,
				Value: "v1",
			},
			{
				Name:  tableTagId2,
				Value: "v2",
			},
		}
		inlineConstraint := ColumnInlineConstraint{
			Name: String("INLINE_CONSTRAINT"),
			Type: ColumnConstraintTypePrimaryKey,
		}
		require.NoError(t, err)
		outOfLineConstraint1 := OutOfLineConstraint{
			Name:    String("OUT_OF_LINE_CONSTRAINT"),
			Type:    ColumnConstraintTypeForeignKey,
			Columns: []string{"COLUMN_1", "COLUMN_2"},
			ForeignKey: &OutOfLineForeignKey{
				TableName:   randomSchemaObjectIdentifier(),
				ColumnNames: []string{"COLUMN_3", "COLUMN_4"},
				Match:       Pointer(FullMatchType),
				On: &ForeignKeyOnAction{
					OnUpdate: Pointer(ForeignKeySetNullAction),
					OnDelete: Pointer(ForeignKeyRestrictAction),
				},
			},
		}
		outOfLineConstraint2 := OutOfLineConstraint{
			Type:              ColumnConstraintTypeUnique,
			Columns:           []string{"COLUMN_1"},
			Enforced:          Bool(true),
			Deferrable:        Bool(true),
			InitiallyDeferred: Bool(true),
			Enable:            Bool(true),
			Rely:              Bool(true),
		}
		stageFileFormat := StageFileFormat{
			Type: Pointer(FileFormatTypeCSV),
			Options: &FileFormatTypeOptions{
				CSVCompression: Pointer(CSVCompressionAuto),
			},
		}
		stageCopyOptions := StageCopyOptions{
			OnError: &StageCopyOnErrorOptions{SkipFile: String("SKIP_FILE")},
		}
		rowAccessPolicy := TableRowAccessPolicy{
			Name: randomSchemaObjectIdentifier(),
			On:   []string{"COLUMN_1", "COLUMN_2"},
		}
		columns := []TableColumn{{
			Name:    columnName,
			Type:    columnType,
			Collate: &collation,
			Comment: &columnComment,
			DefaultValue: &ColumnDefaultValue{
				Identity: &ColumnIdentity{
					Start:     10,
					Increment: 1,
					Order:     Bool(true),
				},
			},
			NotNull:          Bool(true),
			MaskingPolicy:    &maskingPolicy,
			Tags:             columnTags,
			InlineConstraint: &inlineConstraint,
		}}
		opts := &createTableOptions{
			name:                       id,
			ColumnsAndConstraints:      CreateTableColumnsAndConstraints{columns, []OutOfLineConstraint{outOfLineConstraint1, outOfLineConstraint2}},
			ClusterBy:                  []string{"COLUMN_1", "COLUMN_2"},
			EnableSchemaEvolution:      Bool(true),
			StageFileFormat:            &stageFileFormat,
			StageCopyOptions:           &stageCopyOptions,
			DataRetentionTimeInDays:    Int(10),
			MaxDataExtensionTimeInDays: Int(100),
			ChangeTracking:             Bool(true),
			DefaultDDLCollation:        String("en"),
			CopyGrants:                 Bool(true),
			RowAccessPolicy:            &rowAccessPolicy,
			Tags:                       tableTags,
			Comment:                    &tableComment,
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE TABLE %s (%s %s CONSTRAINT INLINE_CONSTRAINT PRIMARY KEY NOT NULL COLLATE 'de' IDENTITY START 10 INCREMENT 1 ORDER MASKING POLICY %s USING (FOO, BAR) TAG (%s = 'v1', %s = 'v2') COMMENT '%s', CONSTRAINT OUT_OF_LINE_CONSTRAINT FOREIGN KEY (COLUMN_1, COLUMN_2) REFERENCES %s (COLUMN_3, COLUMN_4) MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT, UNIQUE (COLUMN_1) ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE RELY) CLUSTER BY (COLUMN_1, COLUMN_2) ENABLE_SCHEMA_EVOLUTION = true STAGE_FILE_FORMAT = (TYPE = CSV COMPRESSION = AUTO) STAGE_COPY_OPTIONS = (ON_ERROR = SKIP_FILE) DATA_RETENTION_TIME_IN_DAYS = 10 MAX_DATA_EXTENSION_TIME_IN_DAYS = 100 CHANGE_TRACKING = true DEFAULT_DDL_COLLATION = 'en' COPY GRANTS ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) TAG (%s = 'v1', %s = 'v2') COMMENT = '%s'`,
			id.FullyQualifiedName(),
			columnName,
			columnType,
			maskingPolicy.Name.FullyQualifiedName(),
			columnTagId1.FullyQualifiedName(),
			columnTagId2.FullyQualifiedName(),
			columnComment,
			outOfLineConstraint1.ForeignKey.TableName.FullyQualifiedName(),
			rowAccessPolicy.Name.FullyQualifiedName(),
			tableTagId1.FullyQualifiedName(),
			tableTagId2.FullyQualifiedName(),
			tableComment,
		)
	})

	t.Run("with skip file x", func(t *testing.T) {
		columns := []TableColumnRequest{
			{name: "FIRST_COLUMN", type_: DataTypeVARCHAR},
		}
		request := NewCreateTableRequest(id, columns).
			WithStageCopyOptions(*NewStageCopyOptionsRequest().WithOnError(*NewStageCopyOnErrorOptionsRequest().WithSkipFileX(5)))
		assertOptsValidAndSQLEquals(t, request.toOpts(), `CREATE TABLE %s (FIRST_COLUMN VARCHAR) STAGE_COPY_OPTIONS = (ON_ERROR = SKIP_FILE_5)`, id.FullyQualifiedName())
	})

	t.Run("with skip file x %", func(t *testing.T) {
		columns := []TableColumnRequest{
			{name: "FIRST_COLUMN", type_: DataTypeVARCHAR},
		}
		request := NewCreateTableRequest(id, columns).
			WithStageCopyOptions(*NewStageCopyOptionsRequest().WithOnError(*NewStageCopyOnErrorOptionsRequest().WithSkipFileXPercent(10)))
		assertOptsValidAndSQLEquals(t, request.toOpts(), `CREATE TABLE %s (FIRST_COLUMN VARCHAR) STAGE_COPY_OPTIONS = (ON_ERROR = 'SKIP_FILE_10%%')`, id.FullyQualifiedName())
	})
}

func TestTableCreateAsSelect(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *createTableAsSelectOptions {
		return &createTableAsSelectOptions{
			name:    id,
			Columns: []TableAsSelectColumn{{Name: "a"}},
		}
	}

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableAsSelectOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableAsSelectOptions", "name"))
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableAsSelectOptions", "name"))
	})

	t.Run("validation: no columns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Columns = []TableAsSelectColumn{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("createTableAsSelectOptions", "Columns"))
	})

	t.Run("validation: no query", func(t *testing.T) {
		opts := defaultOpts()
		opts.Query = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("createTableAsSelectOptions", "Query"))
	})

	t.Run("with complete options", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		columnName := "FIRST_COLUMN"
		columnTypeRaw, err := datatypes.ParseDataType("VARCHAR")
		require.NoError(t, err)
		columnType := LegacyDataTypeFrom(columnTypeRaw)
		maskingPolicy := TableAsSelectColumnMaskingPolicy{
			Name: randomSchemaObjectIdentifier(),
		}
		rowAccessPolicy := TableRowAccessPolicy{
			Name: randomSchemaObjectIdentifier(),
			On:   []string{"COLUMN_1", "COLUMN_2"},
		}
		opts := &createTableAsSelectOptions{
			OrReplace: Bool(true),
			name:      id,
			Columns: []TableAsSelectColumn{
				{
					Name:          columnName,
					Type:          Pointer(columnType),
					MaskingPolicy: &maskingPolicy,
				},
			},
			ClusterBy:  []string{"COLUMN_1", "COLUMN_2"},
			CopyGrants: Bool(true),

			RowAccessPolicy: &rowAccessPolicy,
			Query:           "SELECT * FROM ANOTHER_TABLE",
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE TABLE %s (FIRST_COLUMN VARCHAR MASKING POLICY %s) CLUSTER BY (COLUMN_1, COLUMN_2) COPY GRANTS ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) AS SELECT * FROM ANOTHER_TABLE",
			id.FullyQualifiedName(),
			maskingPolicy.Name.FullyQualifiedName(),
			rowAccessPolicy.Name.FullyQualifiedName(),
		)
	})
}

func TestTableCreateUsingTemplate(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *createTableUsingTemplateOptions {
		return &createTableUsingTemplateOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createTableUsingTemplateOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableUsingTemplateOptions", "name"))
	})

	t.Run("with complete options", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &createTableUsingTemplateOptions{
			OrReplace:  Bool(true),
			name:       id,
			CopyGrants: Bool(true),
			Query:      []string{"sample_data"},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE TABLE %s COPY GRANTS USING TEMPLATE (sample_data)", id.FullyQualifiedName())
	})
}

func TestTableCreateLike(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *createTableLikeOptions {
		return &createTableLikeOptions{
			name:        id,
			SourceTable: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createTableLikeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableLikeOptions", "name"))
	})

	t.Run("validation: source table's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.SourceTable = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableLikeOptions", "SourceTable"))
	})

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableLikeOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableLikeOptions", "name"))
	})

	t.Run("with complete options", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		sourceTable := randomSchemaObjectIdentifier()
		opts := &createTableLikeOptions{
			OrReplace:   Bool(true),
			name:        id,
			SourceTable: sourceTable,
			ClusterBy:   []string{"date", "id"},
			CopyGrants:  Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE TABLE %s LIKE %s CLUSTER BY (date, id) COPY GRANTS", id.FullyQualifiedName(), sourceTable.FullyQualifiedName())
	})
}

func TestTableCreateClone(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *createTableCloneOptions {
		return &createTableCloneOptions{
			name:        id,
			SourceTable: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createTableCloneOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableCloneOptions", "name"))
	})

	t.Run("validation: source table's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.SourceTable = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableCloneOptions", "SourceTable"))
	})

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableCloneOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("createTableCloneOptions", "name"))
	})

	t.Run("with complete options", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		sourceTable := randomSchemaObjectIdentifier()
		clonePoint := ClonePoint{
			Moment: CloneMomentAt,
			At: TimeTravel{
				Offset: Int(0),
			},
		}
		opts := &createTableCloneOptions{
			OrReplace:   Bool(true),
			name:        id,
			SourceTable: sourceTable,
			ClonePoint:  &clonePoint,
			CopyGrants:  Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE TABLE %s CLONE %s AT (OFFSET => 0) COPY GRANTS", id.FullyQualifiedName(), sourceTable.FullyQualifiedName())
	})
}

func TestTableAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *alterTableOptions {
		return &alterTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *alterTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: no action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterTableOptions", "NewName", "SwapWith", "ClusteringAction", "ColumnAction", "ConstraintAction", "ExternalTableAction", "SearchOptimizationAction", "Set", "SetTags", "UnsetTags", "Unset", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllAccessRowPolicies"))
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("alterTableOptions", "name"))
	})

	t.Run("validation: both NewName and SwapWith are present ", func(t *testing.T) {
		opts := defaultOpts()
		opts.NewName = Pointer(randomSchemaObjectIdentifier())
		opts.SwapWith = Pointer(randomSchemaObjectIdentifier())

		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterTableOptions", "NewName", "SwapWith", "ClusteringAction", "ColumnAction", "ConstraintAction", "ExternalTableAction", "SearchOptimizationAction", "Set", "SetTags", "UnsetTags", "Unset", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllAccessRowPolicies"))
	})

	t.Run("validation: NewName's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.NewName = Pointer(emptySchemaObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("alterTableOptions", "NewName"))
	})

	t.Run("validation: SwapWith incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.SwapWith = Pointer(emptySchemaObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("alterTableOptions", "SwapWith"))
	})

	t.Run("validation: clustering action - no option present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ClusteringAction = &TableClusteringAction{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ClusteringAction", "ClusterBy", "Recluster", "ChangeReclusterState", "DropClusteringKey"))
	})

	t.Run("validation: clustering action - two options present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ClusteringAction = &TableClusteringAction{
			ClusterBy: []string{"date"},
			Recluster: &TableReclusterAction{
				MaxSize:   Int(10),
				Condition: String("true"),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ClusteringAction", "ClusterBy", "Recluster", "ChangeReclusterState", "DropClusteringKey"))
	})

	t.Run("validation: column action - no option present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ColumnAction = &TableColumnAction{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ColumnAction", "Add", "Rename", "Alter", "SetMaskingPolicy", "UnsetMaskingPolicy", "SetTags", "UnsetTags", "DropColumns"))
	})

	t.Run("validation: column action - two options present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ColumnAction = &TableColumnAction{
			Add: &TableColumnAddAction{},
			Rename: &TableColumnRenameAction{
				NewName: "new",
				OldName: "old",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ColumnAction", "Add", "Rename", "Alter", "SetMaskingPolicy", "UnsetMaskingPolicy", "SetTags", "UnsetTags", "DropColumns"))
	})

	t.Run("validation: column action alter - no option present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ColumnAction = &TableColumnAction{
			Alter: []TableColumnAlterAction{{}},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableColumnAlterAction", "DropDefault", "SetDefault", "NotNullConstraint", "Type", "Comment", "UnsetComment"))
	})

	t.Run("validation: column action alter - two options present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ColumnAction = &TableColumnAction{
			Alter: []TableColumnAlterAction{{
				DropDefault: Bool(true),
				SetDefault:  Pointer(SequenceName("sequence")),
			}},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableColumnAlterAction", "DropDefault", "SetDefault", "NotNullConstraint", "Type", "Comment", "UnsetComment"))
	})

	t.Run("validation: constraint alter action - no option present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ConstraintAction = &TableConstraintAction{
			Alter: &TableConstraintAlterAction{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableConstraintAlterAction", "ConstraintName", "PrimaryKey", "Unique", "ForeignKey", "Columns"))
	})

	t.Run("validation: constraint alter action - two options present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ConstraintAction = &TableConstraintAction{
			Alter: &TableConstraintAlterAction{
				ConstraintName: String("constraint"),
				Unique:         Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableConstraintAlterAction", "ConstraintName", "PrimaryKey", "Unique", "ForeignKey", "Columns"))
	})

	t.Run("validation: constraint drop action - no option present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ConstraintAction = &TableConstraintAction{
			Drop: &TableConstraintDropAction{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableConstraintDropAction", "ConstraintName", "PrimaryKey", "Unique", "ForeignKey", "Columns"))
	})

	t.Run("validation: constraint drop action - two options present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ConstraintAction = &TableConstraintAction{
			Drop: &TableConstraintDropAction{
				ConstraintName: String("constraint"),
				Unique:         Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableConstraintDropAction", "ConstraintName", "PrimaryKey", "Unique", "ForeignKey", "Columns"))
	})

	t.Run("validation: external action - no option present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalTableAction = &TableExternalTableAction{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableExternalTableAction", "Add", "Rename", "Drop"))
	})

	t.Run("validation: external action - two options present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalTableAction = &TableExternalTableAction{
			Add: &TableExternalTableColumnAddAction{},
			Rename: &TableExternalTableColumnRenameAction{
				OldName: "old_name",
				NewName: "new_name",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableExternalTableAction", "Add", "Rename", "Drop"))
	})

	t.Run("validation: external action - drop with no columns", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalTableAction = &TableExternalTableAction{
			Drop: &TableExternalTableColumnDropAction{
				Names: []string{},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("TableExternalTableColumnDropAction", "Names"))
	})

	t.Run("validation: search optimization - no option present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SearchOptimizationAction = &TableSearchOptimizationAction{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableSearchOptimizationAction", "Add", "Drop"))
	})

	t.Run("validation: search optimization - two options present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SearchOptimizationAction = &TableSearchOptimizationAction{
			Add:  &AddSearchOptimization{},
			Drop: &DropSearchOptimization{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TableSearchOptimizationAction", "Add", "Drop"))
	})

	t.Run("empty options", func(t *testing.T) {
		opts := &alterTableOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("alterTableOptions", "name"))
	})

	t.Run("rename", func(t *testing.T) {
		newID := randomSchemaObjectIdentifierInSchema(id.SchemaId())
		opts := &alterTableOptions{
			name:    id,
			NewName: &newID,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
	})

	t.Run("swap with", func(t *testing.T) {
		targetTableId := randomSchemaObjectIdentifierInSchema(id.SchemaId())
		opts := &alterTableOptions{
			name:     id,
			SwapWith: &targetTableId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s SWAP WITH %s", id.FullyQualifiedName(), targetTableId.FullyQualifiedName())
	})

	t.Run("cluster by", func(t *testing.T) {
		clusterByColumns := []string{"date", "id"}
		opts := &alterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				ClusterBy: clusterByColumns,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s CLUSTER BY (date, id)", id.FullyQualifiedName())
	})

	t.Run("recluster", func(t *testing.T) {
		condition := "name = 'John'"
		opts := &alterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				Recluster: &TableReclusterAction{
					MaxSize:   Int(1024),
					Condition: &condition,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s RECLUSTER MAX_SIZE = 1024 WHERE name = 'John'", id.FullyQualifiedName())
	})

	t.Run("suspend recluster", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				ChangeReclusterState: &TableReclusterChangeState{
					State: Pointer(ReclusterStateSuspend),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s SUSPEND RECLUSTER", id.FullyQualifiedName())
	})

	t.Run("drop clustering key", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				DropClusteringKey: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s DROP CLUSTERING KEY", id.FullyQualifiedName())
	})

	t.Run("add new column", func(t *testing.T) {
		columnName := "NEXT_COLUMN"
		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				Add: &TableColumnAddAction{
					IfNotExists: Bool(true),
					Name:        columnName,
					Type:        DataTypeVARCHAR,
					Collate:     String("utf8"),
					DefaultValue: &ColumnDefaultValue{
						Identity: &ColumnIdentity{
							Start:     10,
							Increment: 1,
						},
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s ADD COLUMN IF NOT EXISTS NEXT_COLUMN VARCHAR COLLATE 'utf8' IDENTITY START 10 INCREMENT 1", id.FullyQualifiedName())
	})

	t.Run("rename column", func(t *testing.T) {
		oldColumn := "OLD_NAME"
		newColumnName := "NEW_NAME"
		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				Rename: &TableColumnRenameAction{
					OldName: oldColumn,
					NewName: newColumnName,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s RENAME COLUMN OLD_NAME TO NEW_NAME", id.FullyQualifiedName())
	})

	t.Run("alter column", func(t *testing.T) {
		// column_1
		columnOneName := "COLUMN_1"
		alterActionsForColumnOne := []TableColumnAlterAction{
			{
				Name:        columnOneName,
				DropDefault: Bool(true),
			},
			{
				Name:       columnOneName,
				SetDefault: Pointer(SequenceName("SEQUENCE_1")),
			},
			{
				Name:         columnOneName,
				UnsetComment: Bool(true),
			},
		}

		columnTwoName := "COLUMN_2"
		alterActionsForColumnTwo := []TableColumnAlterAction{
			{
				Name:        columnTwoName,
				DropDefault: Bool(true),
			},
			{
				Name:       columnTwoName,
				SetDefault: Pointer(SequenceName("SEQUENCE_2")),
			},
			{
				Name:    columnTwoName,
				Comment: String("comment"),
			},
			{
				Name:    columnTwoName,
				Type:    Pointer(DataTypeVARCHAR),
				Collate: String("utf8"),
			},
			{
				Name:              columnTwoName,
				NotNullConstraint: &TableColumnNotNullConstraint{Drop: Bool(true)},
			},
		}
		var actions []TableColumnAlterAction
		actions = append(actions, alterActionsForColumnOne...)
		actions = append(actions, alterActionsForColumnTwo...)

		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				Alter: actions,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s ALTER COLUMN COLUMN_1 DROP DEFAULT, COLUMN COLUMN_1 SET DEFAULT SEQUENCE_1.NEXTVAL, COLUMN COLUMN_1 UNSET COMMENT, COLUMN COLUMN_2 DROP DEFAULT, COLUMN COLUMN_2 SET DEFAULT SEQUENCE_2.NEXTVAL, COLUMN COLUMN_2 COMMENT 'comment', COLUMN COLUMN_2 SET DATA TYPE VARCHAR COLLATE 'utf8', COLUMN COLUMN_2 DROP NOT NULL", id.FullyQualifiedName())
	})

	t.Run("alter: set masking policy", func(t *testing.T) {
		maskingPolicyName := randomSchemaObjectIdentifier()
		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				SetMaskingPolicy: &TableColumnAlterSetMaskingPolicyAction{
					ColumnName:        "COLUMN_1",
					MaskingPolicyName: maskingPolicyName,
					Using:             []string{"FOO", "BAR"},
					Force:             Bool(true),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s ALTER COLUMN COLUMN_1 SET MASKING POLICY %s USING (FOO, BAR) FORCE", id.FullyQualifiedName(), maskingPolicyName.FullyQualifiedName())
	})

	t.Run("alter: unset masking policy", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				UnsetMaskingPolicy: &TableColumnAlterUnsetMaskingPolicyAction{
					ColumnName: "COLUMN_1",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s ALTER COLUMN COLUMN_1 UNSET MASKING POLICY", id.FullyQualifiedName())
	})

	t.Run("alter: set tags", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		columnTags := []TagAssociation{
			{
				Name:  tagId1,
				Value: "v1",
			},
			{
				Name:  tagId2,
				Value: "v2",
			},
		}
		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				SetTags: &TableColumnAlterSetTagsAction{
					ColumnName: "COLUMN_1",
					Tags:       columnTags,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s ALTER COLUMN COLUMN_1 SET TAG %s = 'v1', %s = 'v2'`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("alter: unset tags", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		columnTags := []ObjectIdentifier{
			tagId1,
			tagId2,
		}
		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				UnsetTags: &TableColumnAlterUnsetTagsAction{
					ColumnName: "COLUMN_1",
					Tags:       columnTags,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s ALTER COLUMN COLUMN_1 UNSET TAG %s, %s`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("alter: drop columns", func(t *testing.T) {
		columns := []string{"COLUMN_1", "COLUMN_2"}
		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				DropColumns: &TableColumnAlterDropColumns{
					IfExists: Bool(true),
					Columns:  columns,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s DROP COLUMN IF EXISTS COLUMN_1, COLUMN_2", id.FullyQualifiedName())
	})

	t.Run("validation: alter constraint: no option", func(t *testing.T) {
		opts := &alterTableOptions{
			name:             id,
			ConstraintAction: &TableConstraintAction{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ConstraintAction", "Add", "Rename", "Alter", "Drop"))
	})

	t.Run("validation: alter constraint: more than one option", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type:    ColumnConstraintTypeUnique,
			Columns: []string{"COLUMN_1"},
		}
		opts := &alterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Add: &outOfLineConstraint,
				Rename: &TableConstraintRenameAction{
					OldName: "OLD_NAME_CONSTRAINT",
					NewName: "NEW_NAME_CONSTRAINT",
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ConstraintAction", "Add", "Rename", "Alter", "Drop"))
	})

	t.Run("validation: alter constraint: validation", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Type: ColumnConstraintTypeUnique,
		}
		opts := &alterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Add: &outOfLineConstraint,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("OutOfLineConstraint", "Columns"))
	})

	t.Run("alter constraint: add", func(t *testing.T) {
		outOfLineConstraint := OutOfLineConstraint{
			Name:    String("OUT_OF_LINE_CONSTRAINT"),
			Type:    ColumnConstraintTypeForeignKey,
			Columns: []string{"COLUMN_1", "COLUMN_2"},
			ForeignKey: &OutOfLineForeignKey{
				TableName:   randomSchemaObjectIdentifier(),
				ColumnNames: []string{"COLUMN_3", "COLUMN_4"},
				Match:       Pointer(FullMatchType),
				On: &ForeignKeyOnAction{
					OnUpdate: Pointer(ForeignKeySetNullAction),
					OnDelete: Pointer(ForeignKeyRestrictAction),
				},
			},
		}
		opts := &alterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Add: &outOfLineConstraint,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s ADD CONSTRAINT OUT_OF_LINE_CONSTRAINT FOREIGN KEY (COLUMN_1, COLUMN_2) REFERENCES %s (COLUMN_3, COLUMN_4) MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT", id.FullyQualifiedName(), outOfLineConstraint.ForeignKey.TableName.FullyQualifiedName())
	})

	t.Run("alter constraint: rename", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Rename: &TableConstraintRenameAction{
					OldName: "OLD_NAME_CONSTRAINT",
					NewName: "NEW_NAME_CONSTRAINT",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s RENAME CONSTRAINT OLD_NAME_CONSTRAINT TO NEW_NAME_CONSTRAINT", id.FullyQualifiedName())
	})

	t.Run("alter constraint: alter", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Alter: &TableConstraintAlterAction{
					ConstraintName: String("OUT_OF_LINE_CONSTRAINT"),
					Columns:        []string{"COLUMN_3", "COLUMN_4"},
					NotEnforced:    Bool(true),
					Validate:       Bool(true),
					Rely:           Bool(true),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s ALTER CONSTRAINT OUT_OF_LINE_CONSTRAINT (COLUMN_3, COLUMN_4) NOT ENFORCED VALIDATE RELY", id.FullyQualifiedName())
	})

	t.Run("alter constraint: drop", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Drop: &TableConstraintDropAction{
					ConstraintName: String("OUT_OF_LINE_CONSTRAINT"),
					Columns:        []string{"COLUMN_3", "COLUMN_4"},
					Cascade:        Bool(true),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s DROP CONSTRAINT OUT_OF_LINE_CONSTRAINT (COLUMN_3, COLUMN_4) CASCADE", id.FullyQualifiedName())
	})

	t.Run("external table: add", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ExternalTableAction: &TableExternalTableAction{
				Add: &TableExternalTableColumnAddAction{
					IfNotExists: Bool(true),
					Name:        "COLUMN_1",
					Type:        DataTypeBoolean,
					Expression:  []string{"SELECT 1"},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s ADD COLUMN IF NOT EXISTS COLUMN_1 BOOLEAN AS (SELECT 1)", id.FullyQualifiedName())
	})

	t.Run("external table: rename", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ExternalTableAction: &TableExternalTableAction{
				Rename: &TableExternalTableColumnRenameAction{
					OldName: "OLD_NAME_COLUMN",
					NewName: "NEW_NAME_COLUMN",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s RENAME COLUMN OLD_NAME_COLUMN TO NEW_NAME_COLUMN", id.FullyQualifiedName())
	})

	t.Run("external table: drop", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			ExternalTableAction: &TableExternalTableAction{
				Drop: &TableExternalTableColumnDropAction{
					IfExists: Bool(true),
					Names:    []string{"COLUMN_3", "COLUMN_4"},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s DROP COLUMN IF EXISTS COLUMN_3, COLUMN_4", id.FullyQualifiedName())
	})

	t.Run("add search optimization", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			SearchOptimizationAction: &TableSearchOptimizationAction{
				Add: &AddSearchOptimization{
					On: []string{"SUBSTRING(*)", "GEO(*)"},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s ADD SEARCH OPTIMIZATION ON SUBSTRING(*), GEO(*)", id.FullyQualifiedName())
	})

	t.Run("drop search optimization", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			SearchOptimizationAction: &TableSearchOptimizationAction{
				Drop: &DropSearchOptimization{
					On: []string{"SUBSTRING(*)", "FOO"},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s DROP SEARCH OPTIMIZATION ON SUBSTRING(*), FOO", id.FullyQualifiedName())
	})

	t.Run("drop search optimization", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			SearchOptimizationAction: &TableSearchOptimizationAction{
				Drop: &DropSearchOptimization{
					On: []string{"SUBSTRING(*)", "FOO"},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE %s DROP SEARCH OPTIMIZATION ON SUBSTRING(*), FOO", id.FullyQualifiedName())
	})

	t.Run("set: with complete options", func(t *testing.T) {
		comment := random.Comment()
		opts := &alterTableOptions{
			name: id,
			Set: &TableSet{
				EnableSchemaEvolution: Bool(true),
				StageFileFormat: &StageFileFormat{
					Type: Pointer(FileFormatTypeCSV),
				},
				StageCopyOptions: &StageCopyOptions{
					OnError: &StageCopyOnErrorOptions{SkipFile: String("SKIP_FILE")},
				},
				DataRetentionTimeInDays:    Int(30),
				MaxDataExtensionTimeInDays: Int(90),
				ChangeTracking:             Bool(false),
				DefaultDDLCollation:        String("us"),
				Comment:                    &comment,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s SET ENABLE_SCHEMA_EVOLUTION = true STAGE_FILE_FORMAT = (TYPE = CSV) STAGE_COPY_OPTIONS = (ON_ERROR = SKIP_FILE) DATA_RETENTION_TIME_IN_DAYS = 30 MAX_DATA_EXTENSION_TIME_IN_DAYS = 90 CHANGE_TRACKING = false DEFAULT_DDL_COLLATION = 'us' COMMENT = '%s'`, id.FullyQualifiedName(), comment)
	})

	t.Run("set tags", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := &alterTableOptions{
			name: id,
			SetTags: []TagAssociation{
				{
					Name:  tagId1,
					Value: "v1",
				},
				{
					Name:  tagId2,
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s SET TAG %s = 'v1', %s = 'v2'`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := &alterTableOptions{
			name: id,
			UnsetTags: []ObjectIdentifier{
				tagId1,
				tagId2,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s UNSET TAG %s, %s`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("unset: complete options", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			Unset: &TableUnset{
				DataRetentionTimeInDays:    Bool(true),
				MaxDataExtensionTimeInDays: Bool(true),
				ChangeTracking:             Bool(true),
				DefaultDDLCollation:        Bool(true),
				EnableSchemaEvolution:      Bool(true),
				Comment:                    Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s UNSET DATA_RETENTION_TIME_IN_DAYS MAX_DATA_EXTENSION_TIME_IN_DAYS CHANGE_TRACKING DEFAULT_DDL_COLLATION ENABLE_SCHEMA_EVOLUTION COMMENT`, id.FullyQualifiedName())
	})

	t.Run("add row access policy", func(t *testing.T) {
		rowAccessPolicyId := randomSchemaObjectIdentifier()

		opts := &alterTableOptions{
			name: id,
			AddRowAccessPolicy: &TableAddRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicyId,
				On:              []string{"FIRST_COLUMN"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s ADD ROW ACCESS POLICY %s ON (FIRST_COLUMN)`, id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("drop row access policy", func(t *testing.T) {
		rowAccessPolicyId := randomSchemaObjectIdentifier()

		opts := &alterTableOptions{
			name: id,
			DropRowAccessPolicy: &TableDropRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicyId,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s DROP ROW ACCESS POLICY %s`, id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("drop and add row access policy", func(t *testing.T) {
		rowAccessPolicyId1 := randomSchemaObjectIdentifier()
		rowAccessPolicyId2 := randomSchemaObjectIdentifier()

		opts := &alterTableOptions{
			name: id,
			DropAndAddRowAccessPolicy: &TableDropAndAddRowAccessPolicy{
				Drop: TableDropRowAccessPolicy{
					RowAccessPolicy: rowAccessPolicyId1,
				},
				Add: TableAddRowAccessPolicy{
					RowAccessPolicy: rowAccessPolicyId2,
					On:              []string{"FIRST_COLUMN"},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON (FIRST_COLUMN)`, id.FullyQualifiedName(), rowAccessPolicyId1.FullyQualifiedName(), rowAccessPolicyId2.FullyQualifiedName())
	})

	t.Run("drop all row access policies", func(t *testing.T) {
		opts := &alterTableOptions{
			name:                     id,
			DropAllAccessRowPolicies: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s DROP ALL ROW ACCESS POLICIES`, id.FullyQualifiedName())
	})
}

func TestTableDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *dropTableOptions {
		return &dropTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("dropTableOptions", "name"))
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP TABLE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP TABLE IF EXISTS %s`, id.FullyQualifiedName())
	})

	t.Run("validation: both cascade and restrict present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Cascade = Bool(true)
		opts.Restrict = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("dropTableOptions", "Cascade", "Restrict"))
	})
}

func TestTableShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *showTableOptions {
		return &showTableOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowPipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("show", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW TABLES`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TABLES LIKE '%s'`, id.Name())
	})
}

func TestTableDescribeColumns(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *describeTableColumnsOptions {
		return &describeTableColumnsOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *describeTableColumnsOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("describeTableColumnsOptions", "name"))
	})

	t.Run("describe", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE TABLE %s TYPE = COLUMNS`, id.FullyQualifiedName())
	})
}

func TestTableColumnDetailsRow_SplitTypeAndCollation(t *testing.T) {
	t.Run("with utf8", func(t *testing.T) {
		row := tableColumnDetailsRow{
			Type: DataType("VARCHAR(10) COLLATE 'utf8'"),
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "utf8", *actualCollation)
	})

	t.Run("with locale", func(t *testing.T) {
		row := tableColumnDetailsRow{
			Type: DataType("VARCHAR(10) COLLATE 'en_US'"),
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "en_US", *actualCollation)
	})

	t.Run("with multiple specifiers", func(t *testing.T) {
		row := tableColumnDetailsRow{
			Type: DataType("VARCHAR(10) COLLATE 'fr_CA-ai-pi-trim'"),
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "fr_CA-ai-pi-trim", *actualCollation)
	})

	t.Run("with empty collation", func(t *testing.T) {
		row := tableColumnDetailsRow{
			Type: DataType("VARCHAR(10) COLLATE ''"),
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("VARCHAR(10)"), actualType)
		assert.Equal(t, "", *actualCollation)
	})

	t.Run("without collation", func(t *testing.T) {
		row := tableColumnDetailsRow{
			Type: DataType("NUMBER(38, 0)"),
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, DataType("NUMBER(38, 0)"), actualType)
		assert.Nil(t, actualCollation)
	})
}

func TestTableDescribeStage(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *describeTableStageOptions {
		return &describeTableStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *describeTableStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("describeTableStageOptions", "name"))
	})

	t.Run("describe", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE TABLE %s TYPE = STAGE`, id.FullyQualifiedName())
	})
}

func TestTable_GetClusterByKeys(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		table := Table{ClusterBy: ""}

		assert.Nil(t, table.GetClusterByKeys())
	})

	t.Run("one param", func(t *testing.T) {
		table := Table{ClusterBy: "LINEAR(abc)"}

		assert.Equal(t, []string{"abc"}, table.GetClusterByKeys())
	})

	t.Run("more params", func(t *testing.T) {
		table := Table{ClusterBy: "LINEAR(abc,def)"}

		assert.Equal(t, []string{"abc", "def"}, table.GetClusterByKeys())
	})

	t.Run("white space", func(t *testing.T) {
		table := Table{ClusterBy: "   LINEAR(  abc  , def )"}

		assert.Equal(t, []string{"abc", "def"}, table.GetClusterByKeys())
	})

	t.Run("with function with one param", func(t *testing.T) {
		table := Table{ClusterBy: "LINEAR(some_func(some_param),other_param)"}

		assert.Equal(t, []string{"some_func(some_param)", "other_param"}, table.GetClusterByKeys())
	})

	t.Run("with function with more than one param", func(t *testing.T) {
		table := Table{ClusterBy: "LINEAR(date_trunc('HOUR',TIMESTAMP_HR),CLIENT_ID)"}

		assert.Equal(t, []string{"date_trunc('HOUR',TIMESTAMP_HR)", "CLIENT_ID"}, table.GetClusterByKeys())
	})

	t.Run("with nested functions", func(t *testing.T) {
		table := Table{ClusterBy: "LINEAR(some_func(some_param, some_other_param, other_func(some_param, some_other_param)),other_param)"}

		assert.Equal(t, []string{"some_func(some_param, some_other_param, other_func(some_param, some_other_param))", "other_param"}, table.GetClusterByKeys())
	})
}
