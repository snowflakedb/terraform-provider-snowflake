package defs

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var tableSetColumnMaskingPolicy = g.NewQueryStruct("TableSetColumnMaskingPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("SET").
	Identifier("MaskingPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	OptionalSQL("FORCE").
	WithSharedToOpts()

var tableUnsetColumnMaskingPolicy = g.NewQueryStruct("TableUnsetColumnMaskingPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("UNSET").
	SQL("MASKING POLICY").
	WithSharedToOpts()

var tableSetColumnProjectionPolicy = g.NewQueryStruct("TableSetColumnProjectionPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("SET").
	Identifier("ProjectionPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PROJECTION POLICY").Required()).
	OptionalSQL("FORCE")

var tableUnsetColumnProjectionPolicy = g.NewQueryStruct("TableUnsetColumnProjectionPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("UNSET").
	SQL("PROJECTION POLICY")

var tableSetColumnTags = g.NewQueryStruct("TableSetColumnTags").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SetTags()

var tableUnsetColumnTags = g.NewQueryStruct("TableUnsetColumnTags").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	UnsetTags()

var tableColumnMaskingPolicy = g.NewQueryStruct("TableColumnMaskingPolicy").
	Identifier("MaskingPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	WithValidation(g.ValidIdentifier, "MaskingPolicy").
	WithSharedToOpts()

var tableColumnProjectionPolicy = g.NewQueryStruct("TableColumnProjectionPolicy").
	Identifier("ProjectionPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PROJECTION POLICY").Required())

var tableDropColumnAction = g.NewQueryStruct("TableDropColumnAction").
	SQL("DROP COLUMN").
	OptionalSQL("IF EXISTS").
	PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Required()).
	WithSharedToOpts()

var tableRenameColumnAction = g.NewQueryStruct("TableRenameColumnAction").
	SQL("RENAME COLUMN").
	Text("OldName", g.KeywordOptions().Required().DoubleQuotes()).
	AssignmentWithFieldName("TO", "string", g.ParameterOptions().NoEquals().DoubleQuotes(), "NewName")

func tableSearchOptimizationAction() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableSearchOptimizationAction).WithSharedToOpts()
}

var (
	CloneMomentEnumDef    = g.NewEnum("CloneMoment", "CloneMoments", "AT", "BEFORE")
	ReclusterStateEnumDef = ssotdefs.ReclusterState
	MatchTypeEnumDef      = g.NewEnum("MatchType", "MatchTypes", "FULL", "SIMPLE", "PARTIAL")
)

func tableColumnsAndConstraints() *g.QueryStruct {
	return g.NewQueryStruct("CreateTableColumnsAndConstraints").
		ListQueryStructField("Columns", tableColumn(), g.KeywordOptions()).
		ListQueryStructField("OutOfLineConstraint", legacyTableOutOfLineConstraint(), g.ListOptions().NoParentheses())
}

func tableColumn() *g.QueryStruct {
	return g.NewQueryStruct("TableColumn").
		Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
		PredefinedQueryStructField("ColumnType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
		PredefinedQueryStructField("InlineConstraint", g.KindOfTPointer[sdkcommons.ColumnInlineConstraint](), g.KeywordOptions()).
		OptionalSQL("NOT NULL").
		OptionalTextAssignment("COLLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
		OptionalQueryStructField("DefaultValue", tableColumnDefaultValue(), g.KeywordOptions()).
		OptionalSharedQueryStructField("MaskingPolicy", tableColumnMaskingPolicy, g.KeywordOptions()).
		OptionalTags().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())
}

func tableColumnDefaultValue() *g.QueryStruct {
	return g.NewQueryStruct("ColumnDefaultValue").
		OptionalAssignmentWithFieldName("DEFAULT", "*string", g.ParameterOptions().NoEquals(), "Expression").
		OptionalQueryStructField("Identity", tableColumnIdentity(), g.KeywordOptions().SQL("IDENTITY")).
		WithValidation(g.ExactlyOneValueSet, "Expression", "Identity")
}

func tableColumnIdentity() *g.QueryStruct {
	return g.NewQueryStruct("ColumnIdentity").
		NumberAssignment("START", g.ParameterOptions().NoQuotes().NoEquals().Required()).
		NumberAssignment("INCREMENT", g.ParameterOptions().NoQuotes().NoEquals().Required()).
		OptionalSQL("ORDER").
		OptionalSQL("NOORDER").
		WithValidation(g.ConflictingFields, "Order", "Noorder")
}

// legacyTableOutOfLineConstraint duplicates tableOutOfLineConstraint: legacy does not double-quote the constraint name.
// TODO[SNOW-1007542]: reuse tableOutOfLineConstraint
func legacyTableOutOfLineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("OutOfLineConstraint").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals(), "Name").
		WithField(g.EnumLegacy[sdkcommons.ColumnConstraintType]("ConstraintType", g.KeywordOptions().Required())).
		PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Parentheses()).
		OptionalQueryStructField("ForeignKey", tableOutOfLineForeignKey(), g.KeywordOptions()).
		OptionalSQL("ENFORCED").
		OptionalSQL("NOT ENFORCED").
		OptionalSQL("DEFERRABLE").
		OptionalSQL("NOT DEFERRABLE").
		OptionalSQL("INITIALLY DEFERRED").
		OptionalSQL("INITIALLY IMMEDIATE").
		OptionalSQL("ENABLE").
		OptionalSQL("DISABLE").
		OptionalSQL("VALIDATE").
		OptionalSQL("NOVALIDATE").
		OptionalSQL("RELY").
		OptionalSQL("NORELY").
		WithValidation(g.ValidateValueSet, "Columns").
		WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
		WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
		WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
		WithValidation(g.ConflictingFields, "Enable", "Disable").
		WithValidation(g.ConflictingFields, "Validate", "Novalidate").
		WithValidation(g.ConflictingFields, "Rely", "Norely")
}

func tableOutOfLineForeignKey() *g.QueryStruct {
	return g.NewQueryStruct("OutOfLineForeignKey").
		SQL("REFERENCES").
		Identifier("TableName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("ColumnNames", "[]Column", g.ParameterOptions().NoEquals().Parentheses()).
		OptionalEnum("Match", MatchTypeEnumDef, g.ParameterOptions().NoEquals().SQL("MATCH")).
		PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions())
}

func tableCopyOptions() *g.QueryStruct {
	return g.NewQueryStruct("LegacyTableCopyOptions").
		OptionalQueryStructField("OnError", tableCopyOnErrorOptions(), g.ParameterOptions().SQL("ON_ERROR")).
		OptionalNumberAssignment("SIZE_LIMIT", g.ParameterOptions()).
		OptionalBooleanAssignment("PURGE", g.ParameterOptions()).
		OptionalBooleanAssignment("RETURN_FAILED_ONLY", g.ParameterOptions()).
		OptionalAssignment("MATCH_BY_COLUMN_NAME", "*StageCopyColumnMapOption", g.ParameterOptions()).
		OptionalBooleanAssignment("ENFORCE_LENGTH", g.ParameterOptions()).
		OptionalBooleanAssignment("TRUNCATECOLUMNS", g.ParameterOptions()).
		OptionalBooleanAssignment("FORCE", g.ParameterOptions())
}

func tableCopyOnErrorOptions() *g.QueryStruct {
	return g.NewQueryStruct("LegacyTableCopyOnErrorOptions").
		OptionalSQLWithCustomFieldName("Continue_", "CONTINUE").
		OptionalText("SkipFile", g.KeywordOptions().SQL("SKIP_FILE")).
		OptionalSQL("ABORT_STATEMENT")
}

func tableAsSelectColumn() *g.QueryStruct {
	return g.NewQueryStruct("TableAsSelectColumn").
		Text("Name", g.KeywordOptions().Required()).
		PredefinedQueryStructField("ColumnType", g.KindOfTPointer[sdkcommons.DataType](), g.KeywordOptions()).
		OptionalIdentifier("MaskingPolicy", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY"))
}

func tableClonePoint() *g.QueryStruct {
	return g.NewQueryStruct("ClonePoint").
		Enum("Moment", CloneMomentEnumDef, g.ParameterOptions().NoEquals().Required()).
		PredefinedQueryStructField("At", "TimeTravel", g.ListOptions().Parentheses().NoComma().Required())
}

func tableClusteringAction() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableClusteringAction)
}

func tableColumnAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAction").
		OptionalQueryStructField("Add", tableColumnAddAction(), g.KeywordOptions().SQL("ADD")).
		OptionalQueryStructField("Rename", legacyTableColumnRenameAction(), g.KeywordOptions()).
		ListQueryStructField("Alter", tableColumnAlterAction(), g.KeywordOptions().SQL("ALTER")).
		OptionalSharedQueryStructField("SetMaskingPolicy", tableSetColumnMaskingPolicy, g.KeywordOptions()).
		OptionalSharedQueryStructField("UnsetMaskingPolicy", tableUnsetColumnMaskingPolicy, g.KeywordOptions()).
		OptionalQueryStructField("SetTags", legacyTableColumnAlterSetTagsAction(), g.KeywordOptions()).
		OptionalQueryStructField("UnsetTags", legacyTableColumnAlterUnsetTagsAction(), g.KeywordOptions()).
		OptionalSharedQueryStructField("DropColumns", tableDropColumnAction, g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "Add", "Rename", "Alter", "SetMaskingPolicy", "UnsetMaskingPolicy", "SetTags", "UnsetTags", "DropColumns")
}

func tableColumnAddAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAddAction").
		SQL("COLUMN").
		IfNotExists().
		Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
		PredefinedQueryStructField("ColumnType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
		OptionalTextAssignment("COLLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
		OptionalQueryStructField("DefaultValue", tableColumnDefaultValue(), g.KeywordOptions()).
		OptionalQueryStructField("InlineConstraint", legacyTableColumnAddInlineConstraint(), g.KeywordOptions()).
		OptionalSharedQueryStructField("MaskingPolicy", tableColumnMaskingPolicy, g.KeywordOptions()).
		OptionalTags().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())
}

// legacyTableColumnAddInlineConstraint duplicates tableColumnInlineConstraint: legacy does not double-quote
// the constraint name and takes the foreign key target as string.
// TODO[SNOW-1007542]: reuse tableColumnInlineConstraint
func legacyTableColumnAddInlineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAddInlineConstraint").
		OptionalSQL("NOT NULL").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals(), "Name").
		WithField(g.EnumLegacy[sdkcommons.ColumnConstraintType]("ConstraintType", g.KeywordOptions().Required())).
		OptionalQueryStructField("ForeignKey", tableColumnAddForeignKey(), g.KeywordOptions())
}

func tableColumnAddForeignKey() *g.QueryStruct {
	return g.NewQueryStruct("ColumnAddForeignKey").
		Text("TableName", g.KeywordOptions().SQL("REFERENCES").Required()).
		Text("ColumnName", g.KeywordOptions().Parentheses().Required())
}

// legacyTableColumnRenameAction duplicates tableRenameColumnAction: legacy does not double-quote
// the old and new column names.
// TODO[SNOW-1007542]: reuse tableRenameColumnAction
func legacyTableColumnRenameAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnRenameAction").
		AssignmentWithFieldName("RENAME COLUMN", "string", g.ParameterOptions().NoEquals().Required(), "OldName").
		AssignmentWithFieldName("TO", "string", g.ParameterOptions().NoEquals().Required(), "NewName")
}

func tableColumnAlterAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAlterAction").
		SQL("COLUMN").
		Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
		OptionalSQL("DROP DEFAULT").
		WithField(g.OptionalEnumLegacy[sdkcommons.SequenceName]("SetDefault", g.ParameterOptions().NoEquals().SQL("SET DEFAULT"))).
		OptionalInlineQueryStructField("NotNullConstraint", tableColumnNotNullConstraint()).
		PredefinedQueryStructField("DataType", "*DataType", g.ParameterOptions().NoEquals().SQL("SET DATA TYPE")).
		OptionalTextAssignment("COLLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes()).
		OptionalSQL("UNSET COMMENT").
		WithValidation(g.ExactlyOneValueSet, "DropDefault", "SetDefault", "NotNullConstraint", "DataType", "Comment", "UnsetComment")
}

func tableColumnNotNullConstraint() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnNotNullConstraint").
		OptionalSQLWithCustomFieldName("Set", "SET NOT NULL").
		OptionalSQLWithCustomFieldName("Drop", "DROP NOT NULL")
}

// legacyTableColumnAlterSetTagsAction duplicates tableSetColumnTags: legacy does not double-quote
// the column name.
// TODO[SNOW-1007542]: reuse tableSetColumnTags
func legacyTableColumnAlterSetTagsAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAlterSetTagsAction").
		SQL("ALTER COLUMN").
		Text("ColumnName", g.KeywordOptions().Required()).
		SetTags()
}

// legacyTableColumnAlterUnsetTagsAction duplicates tableUnsetColumnTags: legacy does not double-quote
// the column name.
// TODO[SNOW-1007542]: reuse tableUnsetColumnTags
func legacyTableColumnAlterUnsetTagsAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAlterUnsetTagsAction").
		SQL("ALTER COLUMN").
		Text("ColumnName", g.KeywordOptions().Required()).
		UnsetTags()
}

func tableConstraintAction() *g.QueryStruct {
	return g.NewQueryStruct("TableConstraintAction").
		OptionalQueryStructField("Add", legacyTableOutOfLineConstraint(), g.KeywordOptions().SQL("ADD")).
		OptionalQueryStructField("Rename", tableConstraintRenameAction(), g.KeywordOptions().SQL("RENAME CONSTRAINT")).
		OptionalQueryStructField("Alter", tableConstraintAlterAction(), g.KeywordOptions().SQL("ALTER")).
		OptionalQueryStructField("Drop", tableConstraintDropAction(), g.KeywordOptions().SQL("DROP")).
		WithValidation(g.ExactlyOneValueSet, "Add", "Rename", "Alter", "Drop")
}

func tableConstraintRenameAction() *g.QueryStruct {
	return g.NewQueryStruct("TableConstraintRenameAction").
		Text("OldName", g.KeywordOptions().Required()).
		AssignmentWithFieldName("TO", "string", g.ParameterOptions().NoEquals().Required(), "NewName")
}

func tableConstraintAlterAction() *g.QueryStruct {
	return g.NewQueryStruct("TableConstraintAlterAction").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals(), "ConstraintName").
		OptionalSQL("PRIMARY KEY").
		OptionalSQL("UNIQUE").
		OptionalSQL("FOREIGN KEY").
		PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses()).
		OptionalSQL("ENFORCED").
		OptionalSQL("NOT ENFORCED").
		OptionalSQL("VALIDATE").
		OptionalSQL("NOVALIDATE").
		OptionalSQL("RELY").
		OptionalSQL("NORELY").
		WithValidation(g.ExactlyOneValueSet, "ConstraintName", "PrimaryKey", "Unique", "ForeignKey")
}

func tableConstraintDropAction() *g.QueryStruct {
	return g.NewQueryStruct("TableConstraintDropAction").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals(), "ConstraintName").
		OptionalSQL("PRIMARY KEY").
		OptionalSQL("UNIQUE").
		OptionalSQL("FOREIGN KEY").
		PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses()).
		OptionalSQL("CASCADE").
		OptionalSQL("RESTRICT").
		WithValidation(g.ExactlyOneValueSet, "ConstraintName", "PrimaryKey", "Unique", "ForeignKey")
}

func tableExternalTableAction() *g.QueryStruct {
	return g.NewQueryStruct("TableExternalTableAction").
		OptionalQueryStructField("Add", tableExternalTableColumnAddAction(), g.KeywordOptions()).
		OptionalQueryStructField("Rename", tableExternalTableColumnRenameAction(), g.KeywordOptions()).
		OptionalQueryStructField("Drop", tableExternalTableColumnDropAction(), g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "Add", "Rename", "Drop")
}

func tableExternalTableColumnAddAction() *g.QueryStruct {
	return g.NewQueryStruct("TableExternalTableColumnAddAction").
		SQL("ADD COLUMN").
		IfNotExists().
		Text("Name", g.KeywordOptions().Required()).
		PredefinedQueryStructField("ColumnType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
		ListAssignmentWithFieldName("AS", "string", g.ParameterOptions().NoEquals().Parentheses(), "Expression").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())
}

func tableExternalTableColumnRenameAction() *g.QueryStruct {
	return g.NewQueryStruct("TableExternalTableColumnRenameAction").
		AssignmentWithFieldName("RENAME COLUMN", "string", g.ParameterOptions().NoEquals().Required(), "OldName").
		AssignmentWithFieldName("TO", "string", g.ParameterOptions().NoEquals().Required(), "NewName")
}

func tableExternalTableColumnDropAction() *g.QueryStruct {
	return g.NewQueryStruct("TableExternalTableColumnDropAction").
		SQL("DROP COLUMN").
		IfExists().
		PredefinedQueryStructField("Names", "[]string", g.KeywordOptions().Required()).
		WithValidation(g.ValidateValueSet, "Names")
}

func tableSet() *g.QueryStruct {
	return g.NewQueryStruct("TableSet").
		OptionalBooleanAssignment("ENABLE_SCHEMA_EVOLUTION", g.ParameterOptions()).
		OptionalQueryStructField("StageCopyOptions", tableCopyOptions(), g.ListOptions().Parentheses().SQL("STAGE_COPY_OPTIONS =")).
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
		OptionalAssignment("DEFAULT_DDL_COLLATION", g.KindOfT[sdkcommons.StringAllowEmpty](), g.ParameterOptions()).
		OptionalComment()
}

func tableUnset() *g.QueryStruct {
	return g.NewQueryStruct("TableUnset").
		OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
		OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
		OptionalSQL("CHANGE_TRACKING").
		OptionalSQL("DEFAULT_DDL_COLLATION").
		OptionalSQL("ENABLE_SCHEMA_EVOLUTION").
		OptionalSQL("COMMENT")
}

func tableAddRowAccessPolicy() *g.QueryStruct {
	return g.NewQueryStruct("TableAddRowAccessPolicy").
		SQL("ADD").
		Identifier("RowAccessPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
		PredefinedQueryStructField("On", "[]string", g.KeywordOptions().Parentheses().SQL("ON"))
}

func tableDropRowAccessPolicy() *g.QueryStruct {
	return g.NewQueryStruct("TableDropRowAccessPolicy").
		SQL("DROP").
		Identifier("RowAccessPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required())
}

func tableDropAndAddRowAccessPolicy() *g.QueryStruct {
	return g.NewQueryStruct("TableDropAndAddRowAccessPolicy").
		QueryStructField("Drop", tableDropRowAccessPolicy(), g.KeywordOptions().Required()).
		QueryStructField("Add", tableAddRowAccessPolicy(), g.KeywordOptions().Required())
}

func tableAddStorageLifecyclePolicy() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableStorageLifecyclePolicy)
}

func tableSetAggregationPolicy() *g.QueryStruct {
	return g.QueryStructFromAttributeUse("Table", ssotdefs.TableAggregationPolicy, ssot.ClauseSet).WithSharedToOpts()
}

func tableUnsetAggregationPolicy() *g.QueryStruct {
	return g.QueryStructFromAttributeUse("Table", ssotdefs.TableAggregationPolicy, ssot.ClauseUnset).WithSharedToOpts()
}

func tableSetJoinPolicy() *g.QueryStruct {
	return g.QueryStructFromAttributeUse("Table", ssotdefs.TableJoinPolicy, ssot.ClauseSet).WithSharedToOpts()
}

func tableUnsetJoinPolicy() *g.QueryStruct {
	return g.QueryStructFromAttributeUse("Table", ssotdefs.TableJoinPolicy, ssot.ClauseUnset).WithSharedToOpts()
}

// DESCRIBE TABLE returns the collation glued onto the "type" column (e.g. "VARCHAR(200) COLLATE 'en-ci'"),
// so Type is converted manually and Collation is derived from it in additionalConvert().

// TODO [SNOW-1007542]: add missing features:
// - show columns (https://docs.snowflake.com/en/sql-reference/sql/show-columns)
// - show primary keys (https://docs.snowflake.com/en/sql-reference/sql/show-primary-keys)
// - truncate table (https://docs.snowflake.com/en/sql-reference/sql/truncate-table)
// - undrop table (https://docs.snowflake.com/en/sql-reference/sql/undrop-table)
var tablesDef = g.NewInterface(
	"Tables",
	"Table",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	// CREATE comes from ssotdefs.Table. The structs it names stay here because other operations reuse them.
	CreateOperationFromObject(ssotdefs.Table, tableColumnsAndConstraints(), tableCopyOptions()).
	// TODO: check if [...] in the docs (like in https://docs.snowflake.com/en/sql-reference/sql/create-table#create-table-using-template) mean that we can reuse all parameters from "normal" CreateTableOptions
	CustomOperation(
		"CreateAsSelect",
		"https://docs.snowflake.com/en/sql-reference/sql/create-table",
		g.NewQueryStruct("CreateTableAsSelect").
			Create().
			OrReplace().
			SQL("TABLE").
			Name().
			ListQueryStructField("Columns", tableAsSelectColumn(), g.ListOptions().Parentheses().Required()).
			PredefinedQueryStructField("ClusterBy", "[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
			OptionalCopyGrants().
			PredefinedQueryStructField("RowAccessPolicy", "*TableRowAccessPolicy", g.KeywordOptions()).
			AssignmentWithFieldName("AS", "string", g.ParameterOptions().NoEquals().Required(), "Query").
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidateValueSet, "Columns").
			WithValidation(g.ValidateValueSet, "Query"),
	).
	CustomOperation(
		"CreateUsingTemplate",
		"https://docs.snowflake.com/en/sql-reference/sql/create-table#create-table-using-template",
		g.NewQueryStruct("CreateTableUsingTemplate").
			Create().
			OrReplace().
			SQL("TABLE").
			Name().
			OptionalCopyGrants().
			ListAssignmentWithFieldName("USING TEMPLATE", "string", g.ParameterOptions().NoEquals().Parentheses(), "Query").
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomOperation(
		"CreateLike",
		"https://docs.snowflake.com/en/sql-reference/sql/create-table#create-table-like",
		g.NewQueryStruct("CreateTableLike").
			Create().
			OrReplace().
			SQL("TABLE").
			Name().
			Identifier("SourceTable", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("LIKE").Required()).
			PredefinedQueryStructField("ClusterBy", "[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
			OptionalCopyGrants().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifier, "SourceTable"),
	).
	CustomOperation(
		"CreateClone",
		"https://docs.snowflake.com/en/sql-reference/sql/create-clone",
		g.NewQueryStruct("CreateTableClone").
			Create().
			OrReplace().
			SQL("TABLE").
			Name().
			Identifier("SourceTable", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("CLONE").Required()).
			OptionalQueryStructField("ClonePoint", tableClonePoint(), g.KeywordOptions()).
			OptionalCopyGrants().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifier, "SourceTable"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-table",
		g.NewQueryStruct("AlterTable").
			Alter().
			SQL("TABLE").
			IfExists().
			Name().
			RenameTo().
			OptionalIdentifier("SwapWith", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("SWAP WITH")).
			OptionalQueryStructField("ClusteringAction", tableClusteringAction(), g.KeywordOptions()).
			OptionalQueryStructField("ColumnAction", tableColumnAction(), g.KeywordOptions()).
			OptionalQueryStructField("ConstraintAction", tableConstraintAction(), g.KeywordOptions()).
			OptionalQueryStructField("ExternalTableAction", tableExternalTableAction(), g.KeywordOptions()).
			OptionalQueryStructField("SearchOptimizationAction", tableSearchOptimizationAction(), g.KeywordOptions()).
			OptionalQueryStructField("Set", tableSet(), g.KeywordOptions().SQL("SET")).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalQueryStructField("Unset", tableUnset(), g.KeywordOptions().SQL("UNSET")).
			OptionalQueryStructField("AddRowAccessPolicy", tableAddRowAccessPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("DropRowAccessPolicy", tableDropRowAccessPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("DropAndAddRowAccessPolicy", tableDropAndAddRowAccessPolicy(), g.ListOptions().NoParentheses()).
			OptionalSQL("DROP ALL ROW ACCESS POLICIES").
			OptionalQueryStructField("AddStorageLifecyclePolicy", tableAddStorageLifecyclePolicy(), g.KeywordOptions()).
			OptionalSQL("DROP STORAGE LIFECYCLE POLICY").
			OptionalQueryStructField("SetAggregationPolicy", tableSetAggregationPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("UnsetAggregationPolicy", tableUnsetAggregationPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("SetJoinPolicy", tableSetJoinPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("UnsetJoinPolicy", tableUnsetJoinPolicy(), g.KeywordOptions()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "RenameTo").
			WithValidation(g.ValidIdentifierIfSet, "SwapWith").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "SwapWith", "ClusteringAction", "ColumnAction", "ConstraintAction", "ExternalTableAction", "SearchOptimizationAction", "Set", "SetTags", "UnsetTags", "Unset", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "AddStorageLifecyclePolicy", "DropStorageLifecyclePolicy", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetJoinPolicy", "UnsetJoinPolicy").
			WithAdditionalValidations(),
	).
	DropOperationFromObject(ssotdefs.Table).
	ShowOperationFromObject(ssotdefs.Table).
	// Describe and select come from ssotdefs.Table. CREATE AS SELECT, USING TEMPLATE, LIKE, and CLONE stay handwritten,
	// so SdkDefineFromSsot stays off. Follow-up: convert those four and turn the setting on.
	CustomOperationsFromObject(ssotdefs.Table).
	WithEnums(ssotdefs.TableConstraintType, ssotdefs.TableScope, ssotdefs.TableKind, CloneMomentEnumDef, ReclusterStateEnumDef, MatchTypeEnumDef, ssotdefs.TableSearchMethod)

func tableOutOfLineUniquePK() *g.QueryStruct {
	return withOutOfLineConstraintTail(
		g.NewQueryStruct("TableOutOfLineUniquePK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			OptionalSQL("UNIQUE").
			OptionalSQL("PRIMARY KEY").
			PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Parentheses()).
			WithValidation(g.ExactlyOneValueSet, "Unique", "PrimaryKey"),
	)
}

func tableOutOfLineFK() *g.QueryStruct {
	return withOutOfLineConstraintTail(
		g.NewQueryStruct("TableOutOfLineFK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			SQL("FOREIGN KEY").
			PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Parentheses()).
			Identifier("References", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("REFERENCES").Required()).
			PredefinedQueryStructField("RefColumns", "[]Column", g.KeywordOptions().Parentheses()).
			OptionalEnum("Match", MatchTypeEnumDef, g.ParameterOptions().NoEquals().SQL("MATCH")).
			PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions()).
			WithValidation(g.ValidIdentifier, "References"),
	)
}

func tableOutOfLineCH() *g.QueryStruct {
	return g.NewQueryStruct("TableOutOfLineCH").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
		SQL("CHECK").
		SQLWithCustomFieldName("openParen", "(").
		Text("Expression", g.KeywordOptions().NoQuotes().Required()).
		SQLWithCustomFieldName("closeParen", ")").
		OptionalSQL("ENABLE VALIDATE").
		OptionalSQL("ENABLE NOVALIDATE").
		WithValidation(g.ConflictingFields, "EnableValidate", "EnableNovalidate")
}

func tableOutOfLineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("TableOutOfLineConstraint").
		OptionalQueryStructField("UniquePK", tableOutOfLineUniquePK(), g.KeywordOptions()).
		OptionalQueryStructField("FK", tableOutOfLineFK(), g.KeywordOptions()).
		OptionalQueryStructField("CH", tableOutOfLineCH(), g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "UniquePK", "FK", "CH")
}

// withOutOfLineConstraintTail appends the tail clauses shared by out-of-line UNIQUE/PK and FK
// constraints (ENFORCED / DEFERRABLE / INITIALLY / ENABLE / VALIDATE / RELY pairs plus COMMENT)
// and their ConflictingFields validations. Applied as a wrapper so tail fields are emitted
// after the caller's struct-specific fields. Out-of-line CHECK uses its own ENABLE VALIDATE
// pair and is not wrapped.
func withOutOfLineConstraintTail(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalSQL("ENFORCED").
		OptionalSQL("NOT ENFORCED").
		OptionalSQL("DEFERRABLE").
		OptionalSQL("NOT DEFERRABLE").
		OptionalSQL("INITIALLY DEFERRED").
		OptionalSQL("INITIALLY IMMEDIATE").
		OptionalSQL("ENABLE").
		OptionalSQL("DISABLE").
		OptionalSQL("VALIDATE").
		OptionalSQL("NOVALIDATE").
		OptionalSQL("RELY").
		OptionalSQL("NORELY").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes()).
		WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
		WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
		WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
		WithValidation(g.ConflictingFields, "Enable", "Disable").
		WithValidation(g.ConflictingFields, "Validate", "Novalidate").
		WithValidation(g.ConflictingFields, "Rely", "Norely")
}

func tableColumnInlineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnInlineConstraint").
		OptionalQueryStructField("UniquePK", tableColumnInlineUniquePK(), g.KeywordOptions()).
		OptionalQueryStructField("FK", tableColumnInlineFK(), g.KeywordOptions()).
		OptionalQueryStructField("CH", tableColumnInlineCH(), g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "UniquePK", "FK", "CH")
}

func tableColumnInlineUniquePK() *g.QueryStruct {
	return withInlineConstraintTail(
		g.NewQueryStruct("TableColumnInlineUniquePK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			OptionalSQL("UNIQUE").
			OptionalSQL("PRIMARY KEY").
			WithValidation(g.ExactlyOneValueSet, "Unique", "PrimaryKey"),
	)
}

func tableColumnInlineFK() *g.QueryStruct {
	return withInlineConstraintTail(
		g.NewQueryStruct("TableColumnInlineFK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			OptionalSQL("FOREIGN KEY").
			Identifier("References", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("REFERENCES").Required()).
			PredefinedQueryStructField("RefColumn", "[]Column", g.KeywordOptions().Parentheses()).
			OptionalEnum("Match", MatchTypeEnumDef, g.ParameterOptions().NoEquals().SQL("MATCH")).
			PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions()).
			WithValidation(g.ValidIdentifier, "References"),
	)
}

func tableColumnInlineCH() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnInlineCH").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
		SQL("CHECK").
		SQLWithCustomFieldName("openParen", "(").
		Text("Expression", g.KeywordOptions().NoQuotes().Required()).
		SQLWithCustomFieldName("closeParen", ")").
		OptionalSQL("ENABLE VALIDATE").
		OptionalSQL("ENABLE NOVALIDATE").
		WithValidation(g.ConflictingFields, "EnableValidate", "EnableNovalidate")
}

// withInlineConstraintTail appends the tail clauses shared by inline UNIQUE/PK and FK
// constraints (ENFORCED / DEFERRABLE / INITIALLY / ENABLE / VALIDATE / RELY pairs) and
// their ConflictingFields validations. Applied as a wrapper so tail fields are emitted
// after the caller's struct-specific fields.
func withInlineConstraintTail(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalSQL("ENFORCED").
		OptionalSQL("NOT ENFORCED").
		OptionalSQL("DEFERRABLE").
		OptionalSQL("NOT DEFERRABLE").
		OptionalSQL("INITIALLY DEFERRED").
		OptionalSQL("INITIALLY IMMEDIATE").
		OptionalSQL("ENABLE").
		OptionalSQL("DISABLE").
		OptionalSQL("VALIDATE").
		OptionalSQL("NOVALIDATE").
		OptionalSQL("RELY").
		OptionalSQL("NORELY").
		WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
		WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
		WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
		WithValidation(g.ConflictingFields, "Enable", "Disable").
		WithValidation(g.ConflictingFields, "Validate", "Novalidate").
		WithValidation(g.ConflictingFields, "Rely", "Norely")
}
