package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var tableSetColumnMaskingPolicy = g.NewQueryStruct("TableSetColumnMaskingPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("SET").
	Identifier("MaskingPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	OptionalSQL("FORCE")

var tableUnsetColumnMaskingPolicy = g.NewQueryStruct("TableUnsetColumnMaskingPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("UNSET").
	SQL("MASKING POLICY")

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

var TableSearchMethodEnumDef = g.NewEnum(
	"TableSearchMethod", "TableSearchMethods",
	"SUBSTRING", "EQUALITY", "FULL_TEXT",
)

var tableColumnMaskingPolicy = g.NewQueryStruct("TableColumnMaskingPolicy").
	Identifier("MaskingPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses())

var tableColumnProjectionPolicy = g.NewQueryStruct("TableColumnProjectionPolicy").
	Identifier("ProjectionPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PROJECTION POLICY").Required())

var tableDropColumnAction = g.NewQueryStruct("TableDropColumnAction").
	SQL("DROP COLUMN").
	OptionalSQL("IF EXISTS").
	PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Required())

var tableRenameColumnAction = g.NewQueryStruct("TableRenameColumnAction").
	SQL("RENAME COLUMN").
	Text("OldName", g.KeywordOptions().Required().DoubleQuotes()).
	AssignmentWithFieldName("TO", "string", g.ParameterOptions().NoEquals().DoubleQuotes(), "NewName")

func newTableSearchMethodArgs() *g.QueryStruct {
	return g.NewQueryStruct("TableSearchMethodArgs").
		PredefinedQueryStructField("Targets", "[]string", g.KeywordOptions()).
		OptionalTextAssignment("ANALYZER", g.ParameterOptions().ArrowEquals().SingleQuotes())
}

func newTableSearchMethodWithTarget() *g.QueryStruct {
	return g.NewQueryStruct("TableSearchMethodWithTarget").
		PredefinedQueryStructField("Method", TableSearchMethodEnumDef.Kind(), g.KeywordOptions().Required()).
		QueryStructField("Args", newTableSearchMethodArgs(), g.ListOptions().Parentheses())
}

var tableAddSearchOptimization = g.NewQueryStruct("TableAddSearchOptimization").
	SQL("ADD SEARCH OPTIMIZATION").
	ListQueryStructField("On", newTableSearchMethodWithTarget(), g.KeywordOptions().SQL("ON"))

func tableDropSearchOptimizationOn() *g.QueryStruct {
	return g.NewQueryStruct("TableDropSearchOptimizationOn").
		OptionalQueryStructField("SearchMethodWithTarget", newTableSearchMethodWithTarget(), g.KeywordOptions()).
		OptionalText("ColumnName", g.KeywordOptions()).
		OptionalText("ExpressionId", g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "SearchMethodWithTarget", "ColumnName", "ExpressionId")
}

func tableDropSearchOptimization() *g.QueryStruct {
	return g.NewQueryStruct("TableDropSearchOptimization").
		SQL("DROP SEARCH OPTIMIZATION").
		ListQueryStructField("On", tableDropSearchOptimizationOn(), g.KeywordOptions().SQL("ON")).
		WithSharedToOpts()
}

var tableSearchOptimizationAction = g.NewQueryStruct("TableSearchOptimizationAction").
	OptionalQueryStructField(
		"Add",
		tableAddSearchOptimization,
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Drop",
		tableDropSearchOptimization(),
		g.KeywordOptions(),
	).
	WithValidation(g.ExactlyOneValueSet, "Add", "Drop")

var tableSearchOptimizationDetails = g.StructPair("tableSearchOptimizationDetailsRow", "TableSearchOptimizationDetails").
	Number("expression_id").
	Text("method").
	Text("target").
	DataType("target_data_type").
	BoolFromText("active", g.WithBoolTrueValue("true"))

var tableDescribeSearchOptimization = g.NewQueryStruct("DescribeSearchOptimization").
	Describe().
	SQL("SEARCH OPTIMIZATION").
	SQL("ON").
	Name().
	WithValidation(g.ValidIdentifier, "name")

var TableConstraintTypeDef = g.NewEnum(
	"TableConstraintType", "TableConstraintTypes",
	"PRIMARY KEY", "UNIQUE", "FOREIGN KEY",
)

var tableConstraintDetails = g.StructPair("tableConstraintDetailsRow", "TableConstraintDetails").
	Text("CONSTRAINT_CATALOG").
	Text("CONSTRAINT_SCHEMA").
	Text("CONSTRAINT_NAME").
	Text("TABLE_CATALOG").
	Text("TABLE_SCHEMA").
	Text("TABLE_NAME").
	Enum("CONSTRAINT_TYPE", TableConstraintTypeDef).
	BoolFromText("IS_DEFERRABLE", g.WithBoolTrueValue("YES")).
	BoolFromText("INITIALLY_DEFERRED", g.WithBoolTrueValue("YES")).
	OptionalText("COMMENT").
	Time("CREATED").
	Time("LAST_ALTERED").
	BoolFromText("ENFORCED", g.WithBoolTrueValue("YES")).
	BoolFromText("RELY", g.WithBoolTrueValue("YES"))

var tableShowConstraints = g.NewQueryStruct("ShowTableConstraints").
	SQLWithCustomFieldName("selectAll", "SELECT * FROM").
	Identifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required()).
	SQLWithCustomFieldName("dot", ".").
	SQLWithCustomFieldName("informationSchemaTableConstraints", "INFORMATION_SCHEMA.TABLE_CONSTRAINTS").
	SQLWithCustomFieldName("where", "WHERE").
	TextAssignment("TABLE_SCHEMA", g.ParameterOptions().SingleQuotes().Required()).
	SQLWithCustomFieldName("and", "AND").
	TextAssignment("TABLE_NAME", g.ParameterOptions().SingleQuotes().Required()).
	WithValidation(g.ValidIdentifier, "Database")

var tableCheckConstraintDetails = g.StructPair("tableCheckConstraintDetailsRow", "TableCheckConstraintDetails").
	Text("CONSTRAINT_CATALOG").
	Text("CONSTRAINT_SCHEMA").
	Text("CONSTRAINT_TABLE").
	Text("CONSTRAINT_NAME").
	Text("CHECK_CLAUSE")

var tableSelectCheckConstraints = g.NewQueryStruct("SelectCheckConstraints").
	SQLWithCustomFieldName("selectAll", "SELECT * FROM").
	Identifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required()).
	SQLWithCustomFieldName("dot", ".").
	SQLWithCustomFieldName("informationSchemaCheckConstraints", "INFORMATION_SCHEMA.CHECK_CONSTRAINTS").
	SQLWithCustomFieldName("where", "WHERE").
	TextAssignment("CONSTRAINT_SCHEMA", g.ParameterOptions().SingleQuotes().Required()).
	SQLWithCustomFieldName("and", "AND").
	TextAssignment("CONSTRAINT_TABLE", g.ParameterOptions().SingleQuotes().Required()).
	WithValidation(g.ValidIdentifier, "Database")

var (
	TableScopeEnumDef     = g.NewEnum("TableScope", "TableScopes", "GLOBAL", "LOCAL")
	TableKindEnumDef      = g.NewEnum("TableKind", "TableKinds", "TEMPORARY", "VOLATILE", "TRANSIENT")
	CloneMomentEnumDef    = g.NewEnum("CloneMoment", "CloneMoments", "AT", "BEFORE")
	ReclusterStateEnumDef = g.NewEnum("ReclusterState", "ReclusterStates", "RESUME", "SUSPEND")
	MatchTypeEnumDef      = g.NewEnum("MatchType", "MatchTypes", "FULL", "SIMPLE", "PARTIAL")
)

func tableColumnsAndConstraints() *g.QueryStruct {
	return g.NewQueryStruct("CreateTableColumnsAndConstraints").
		ListQueryStructField("Columns", tableColumn(), g.KeywordOptions()).
		ListQueryStructField("OutOfLineConstraint", legacyTableOutOfLineConstraint(), g.ListOptions().NoParentheses())
}

func tableColumn() *g.QueryStruct {
	return g.NewQueryStruct("TableColumn").
		Text("Name", g.KeywordOptions().Required()).
		PredefinedQueryStructField("ColumnType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
		PredefinedQueryStructField("InlineConstraint", g.KindOfTPointer[sdkcommons.ColumnInlineConstraint](), g.KeywordOptions()).
		OptionalSQL("NOT NULL").
		OptionalTextAssignment("COLLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
		OptionalQueryStructField("DefaultValue", tableColumnDefaultValue(), g.KeywordOptions()).
		OptionalQueryStructField("MaskingPolicy", legacyTableColumnMaskingPolicy(), g.KeywordOptions()).
		OptionalTags().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())
}

func tableColumnDefaultValue() *g.QueryStruct {
	return g.NewQueryStruct("ColumnDefaultValue").
		OptionalAssignmentWithFieldName("DEFAULT", "*string", g.ParameterOptions().NoEquals(), "Expression").
		OptionalQueryStructField("Identity", tableColumnIdentity(), g.KeywordOptions().SQL("IDENTITY"))
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations below for now
	// WithValidation(g.ExactlyOneValueSet, "Expression", "Identity")
}

func tableColumnIdentity() *g.QueryStruct {
	return g.NewQueryStruct("ColumnIdentity").
		NumberAssignment("START", g.ParameterOptions().NoQuotes().NoEquals().Required()).
		NumberAssignment("INCREMENT", g.ParameterOptions().NoQuotes().NoEquals().Required()).
		OptionalSQL("ORDER").
		OptionalSQL("NOORDER")
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations below for now
	// WithValidation(g.ConflictingFields, "Order", "Noorder")
}

// legacyTableColumnMaskingPolicy duplicates tableColumnMaskingPolicy: legacy uses []string rather than []Column
// (so column names are not double-quoted).
// TODO[SNOW-1007542]: reuse tableColumnMaskingPolicy
func legacyTableColumnMaskingPolicy() *g.QueryStruct {
	return g.NewQueryStruct("ColumnMaskingPolicy").
		OptionalSQL("WITH").
		SQL("MASKING POLICY").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("Using", "[]string", g.KeywordOptions().Parentheses().SQL("USING"))
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations below for now
	// WithValidation(g.ValidIdentifier, "Name")
}

// legacyTableOutOfLineConstraint duplicates tableOutOfLineConstraint: legacy does not double-quote the constraint name
// and takes columns as []string rather than []Column.
// TODO[SNOW-1007542]: reuse tableOutOfLineConstraint
func legacyTableOutOfLineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("OutOfLineConstraint").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals(), "Name").
		WithField(g.EnumLegacy[sdkcommons.ColumnConstraintType]("ConstraintType", g.KeywordOptions().Required())).
		PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses()).
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
		OptionalSQL("NORELY")
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations below for now
	// WithValidation(g.ValidateValueSet, "Columns").
	// WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
	// WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
	// WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
	// WithValidation(g.ConflictingFields, "Enable", "Disable").
	// WithValidation(g.ConflictingFields, "Validate", "Novalidate").
	// WithValidation(g.ConflictingFields, "Rely", "Norely")
}

func tableOutOfLineForeignKey() *g.QueryStruct {
	return g.NewQueryStruct("OutOfLineForeignKey").
		SQL("REFERENCES").
		Identifier("TableName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("ColumnNames", "[]string", g.ParameterOptions().NoEquals().Parentheses()).
		OptionalEnum("Match", MatchTypeEnumDef, g.ParameterOptions().NoEquals().SQL("MATCH")).
		PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions())
}

func tableFileFormat() *g.QueryStruct {
	return g.NewQueryStruct("LegacyFileFormat").
		OptionalTextAssignment("FORMAT_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalAssignmentWithFieldName("TYPE", "*FileFormatType", g.ParameterOptions(), "FileFormatType").
		PredefinedQueryStructField("Options", "*FileFormatTypeOptionsLegacy", g.ListOptions().NoComma()).
		WithValidation(g.ExactlyOneValueSet, "FormatName", "FileFormatType")
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
	return g.NewQueryStruct("TableClusteringAction").
		PredefinedQueryStructField("ClusterBy", "[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
		OptionalQueryStructField("Recluster", tableReclusterAction(), g.KeywordOptions().SQL("RECLUSTER")).
		OptionalQueryStructField("ChangeReclusterState", tableReclusterChangeState(), g.KeywordOptions()).
		OptionalSQL("DROP CLUSTERING KEY").
		WithValidation(g.ExactlyOneValueSet, "ClusterBy", "Recluster", "ChangeReclusterState", "DropClusteringKey")
}

func tableReclusterAction() *g.QueryStruct {
	return g.NewQueryStruct("TableReclusterAction").
		OptionalNumberAssignment("MAX_SIZE", g.ParameterOptions()).
		OptionalAssignmentWithFieldName("WHERE", "*string", g.ParameterOptions().NoEquals(), "Condition")
}

func tableReclusterChangeState() *g.QueryStruct {
	return g.NewQueryStruct("TableReclusterChangeState").
		OptionalEnum("State", ReclusterStateEnumDef, g.KeywordOptions()).
		SQL("RECLUSTER")
}

func tableColumnAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAction").
		OptionalQueryStructField("Add", tableColumnAddAction(), g.KeywordOptions().SQL("ADD")).
		OptionalQueryStructField("Rename", legacyTableColumnRenameAction(), g.KeywordOptions()).
		ListQueryStructField("Alter", tableColumnAlterAction(), g.KeywordOptions().SQL("ALTER")).
		OptionalQueryStructField("SetMaskingPolicy", legacyTableColumnAlterSetMaskingPolicyAction(), g.KeywordOptions()).
		OptionalQueryStructField("UnsetMaskingPolicy", legacyTableColumnAlterUnsetMaskingPolicyAction(), g.KeywordOptions()).
		OptionalQueryStructField("SetTags", legacyTableColumnAlterSetTagsAction(), g.KeywordOptions()).
		OptionalQueryStructField("UnsetTags", legacyTableColumnAlterUnsetTagsAction(), g.KeywordOptions()).
		OptionalQueryStructField("DropColumns", legacyTableColumnAlterDropColumns(), g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "Add", "Rename", "Alter", "SetMaskingPolicy", "UnsetMaskingPolicy", "SetTags", "UnsetTags", "DropColumns")
}

func tableColumnAddAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAddAction").
		SQL("COLUMN").
		IfNotExists().
		Text("Name", g.KeywordOptions().Required()).
		PredefinedQueryStructField("ColumnType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
		OptionalTextAssignment("COLLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
		OptionalQueryStructField("DefaultValue", tableColumnDefaultValue(), g.KeywordOptions()).
		OptionalQueryStructField("InlineConstraint", legacyTableColumnAddInlineConstraint(), g.KeywordOptions()).
		OptionalQueryStructField("MaskingPolicy", legacyTableColumnMaskingPolicy(), g.KeywordOptions()).
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
		Text("Name", g.KeywordOptions().Required()).
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

// legacyTableColumnAlterSetMaskingPolicyAction duplicates tableSetColumnMaskingPolicy: legacy does
// not double-quote the column name, and takes Using as []string rather than []Column.
// TODO[SNOW-1007542]: reuse tableSetColumnMaskingPolicy
func legacyTableColumnAlterSetMaskingPolicyAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAlterSetMaskingPolicyAction").
		SQL("ALTER COLUMN").
		Text("ColumnName", g.KeywordOptions().Required()).
		SQL("SET MASKING POLICY").
		Identifier("MaskingPolicyName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("Using", "[]string", g.KeywordOptions().Parentheses().SQL("USING")).
		OptionalSQL("FORCE")
}

// legacyTableColumnAlterUnsetMaskingPolicyAction duplicates tableUnsetColumnMaskingPolicy: legacy
// does not double-quote the column name.
// TODO[SNOW-1007542]: reuse tableUnsetColumnMaskingPolicy
func legacyTableColumnAlterUnsetMaskingPolicyAction() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAlterUnsetMaskingPolicyAction").
		SQL("ALTER COLUMN").
		Text("ColumnName", g.KeywordOptions().Required()).
		SQL("UNSET MASKING POLICY")
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

// legacyTableColumnAlterDropColumns duplicates tableDropColumnAction: legacy takes the columns as
// []string rather than []Column, so their names are not double-quoted.
// TODO[SNOW-1007542]: reuse tableDropColumnAction
func legacyTableColumnAlterDropColumns() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnAlterDropColumns").
		SQL("DROP COLUMN").
		IfExists().
		PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Required())
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

// legacyTableSearchOptimizationAction duplicates tableSearchOptimizationAction, differing only in
// that its Add side still needs the legacy helper below.
// TODO[SNOW-1007542]: reuse tableSearchOptimizationAction
func legacyTableSearchOptimizationAction() *g.QueryStruct {
	return g.NewQueryStruct("TableSearchOptimizationActionLegacy").
		OptionalQueryStructField("Add", legacyTableAddSearchOptimization(), g.KeywordOptions()).
		OptionalSharedQueryStructField("Drop", tableDropSearchOptimization(), g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "Add", "Drop")
}

// legacyTableAddSearchOptimization duplicates tableAddSearchOptimization. The shared helper
// requires a search method per target, so it emits ON EQUALITY(c1) and cannot produce legacy's bare ON c1, c2.
// TODO[SNOW-1007542]: reuse tableAddSearchOptimization
func legacyTableAddSearchOptimization() *g.QueryStruct {
	return g.NewQueryStruct("AddSearchOptimization").
		SQL("ADD SEARCH OPTIMIZATION").
		PredefinedQueryStructField("On", "[]string", g.KeywordOptions().SQL("ON"))
}

func tableSet() *g.QueryStruct {
	return g.NewQueryStruct("TableSet").
		OptionalBooleanAssignment("ENABLE_SCHEMA_EVOLUTION", g.ParameterOptions()).
		OptionalQueryStructField("StageFileFormat", tableFileFormat(), g.ListOptions().Parentheses().SQL("STAGE_FILE_FORMAT =")).
		OptionalQueryStructField("StageCopyOptions", tableCopyOptions(), g.ListOptions().Parentheses().SQL("STAGE_COPY_OPTIONS =")).
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
		OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
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
	return g.NewQueryStruct("TableAddStorageLifecyclePolicy").
		SQL("ADD").
		Identifier("StorageLifecyclePolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("STORAGE LIFECYCLE POLICY").Required()).
		ListAssignment("ON", "Column", g.ParameterOptions().NoEquals().Parentheses().Required()).
		WithValidation(g.ValidIdentifier, "StorageLifecyclePolicy").
		WithValidation(g.ValidateValueSet, "On")
}

var tableShowOutput = g.StructPair("tableDBRow", "Table").
	Text("created_on").
	Text("name").
	Text("schema_name").
	Text("database_name").
	Text("kind").
	OptionalText("comment", g.WithRequiredInPlain()).
	OptionalText("cluster_by", g.WithRequiredInPlain()).
	OptionalNumber("rows", g.WithRequiredInPlain()).
	OptionalNumber("bytes").
	Text("owner").
	OptionalNumber("retention_time", g.WithRequiredInPlain()).
	OptionalText("dropped_on").
	OptionalBoolFromText("automatic_clustering", g.WithBoolTrueValue("ON"), g.WithRequiredInPlain()).
	OptionalBoolFromText("change_tracking", g.WithBoolTrueValue("ON"), g.WithRequiredInPlain()).
	OptionalBoolFromText("search_optimization", g.WithBoolTrueValue("ON"), g.WithRequiredInPlain()).
	OptionalText("search_optimization_progress", g.WithRequiredInPlain()).
	OptionalNumber("search_optimization_bytes").
	OptionalBoolFromText("is_external", g.WithRequiredInPlain()).
	OptionalBoolFromText("enable_schema_evolution", g.WithRequiredInPlain()).
	OptionalText("owner_role_type", g.WithRequiredInPlain()).
	OptionalBoolFromText("is_event", g.WithRequiredInPlain()).
	OptionalText("budget")

// DESCRIBE TABLE returns the collation glued onto the "type" column (e.g. "VARCHAR(200) COLLATE 'en-ci'"),
// so Type is converted manually and Collation is derived from it in additionalConvert().
var tableColumnDetails = g.StructPair("tableColumnDetailsRow", "TableColumnDetails").
	Text("name").
	Field("type", "DataType", "DataType", g.WithManualConvert()).
	Text("kind").
	Field("null?", "string", "bool", g.WithDbFieldName("IsNullable"), g.WithPlainFieldName("IsNullable")).
	OptionalText("default").
	Field("primary key", "string", "bool", g.WithDbFieldName("IsPrimary"), g.WithPlainFieldName("IsPrimary")).
	Field("unique key", "string", "bool", g.WithDbFieldName("IsUnique"), g.WithPlainFieldName("IsUnique")).
	OptionalBoolFromText("check").
	OptionalText("expression").
	OptionalText("comment").
	OptionalText("policy name").
	PlainOnlyField("Collation", "*string").
	OptionalText("schema evolution record")

var tableStageDetails = g.StructPair("tableStageDetailsRow", "TableStageDetails").
	Text("parent_property").
	Text("property").
	Text("property_type").
	Text("property_value").
	Text("property_default")

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
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-table",
		g.NewQueryStruct("CreateTable").
			Create().
			OrReplace().
			OptionalEnum("Scope", TableScopeEnumDef, g.KeywordOptions()).
			OptionalEnum("Kind", TableKindEnumDef, g.KeywordOptions()).
			SQL("TABLE").
			IfNotExists().
			Name().
			QueryStructField("ColumnsAndConstraints", tableColumnsAndConstraints(), g.ListOptions().Parentheses().Required()).
			PredefinedQueryStructField("ClusterBy", "[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
			OptionalBooleanAssignment("ENABLE_SCHEMA_EVOLUTION", g.ParameterOptions()).
			OptionalQueryStructField("StageFileFormat", tableFileFormat(), g.ListOptions().Parentheses().NoComma().SQL("STAGE_FILE_FORMAT =")).
			OptionalQueryStructField("StageCopyOptions", tableCopyOptions(), g.ListOptions().Parentheses().NoComma().SQL("STAGE_COPY_OPTIONS =")).
			OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
			OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
			OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
			OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
			OptionalCopyGrants().
			PredefinedQueryStructField("RowAccessPolicy", "*TableRowAccessPolicyLegacy", g.KeywordOptions()).
			OptionalTags().
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			// Per-column and per-out-of-line-constraint checks cannot be generated; see the TODOs on
			// tableColumnIdentity, tableColumnDefaultValue, legacyTableColumnMaskingPolicy, and
			// legacyTableOutOfLineConstraint.
			WithAdditionalValidations(),
	).
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
			PredefinedQueryStructField("RowAccessPolicy", "*TableRowAccessPolicyLegacy", g.KeywordOptions()).
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
			OptionalQueryStructField("SearchOptimizationAction", legacyTableSearchOptimizationAction(), g.KeywordOptions()).
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
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "RenameTo").
			WithValidation(g.ValidIdentifierIfSet, "SwapWith").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "SwapWith", "ClusteringAction", "ColumnAction", "ConstraintAction", "ExternalTableAction", "SearchOptimizationAction", "Set", "SetTags", "UnsetTags", "Unset", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "AddStorageLifecyclePolicy", "DropStorageLifecyclePolicy").
			WithAdditionalValidations(),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-table",
		g.NewQueryStruct("DropTable").
			Drop().
			SQL("TABLE").
			IfExists().
			Name().
			OptionalSQL("CASCADE").
			OptionalSQL("RESTRICT").
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-tables",
		tableShowOutput,
		g.NewQueryStruct("ShowTable").
			Show().
			Terse().
			SQL("TABLES").
			OptionalSQL("HISTORY").
			OptionalLike().
			OptionalExtendedIn().
			OptionalStartsWith().
			OptionalLimitFrom().
			// ErrPatternRequiredForLikeKeyword cannot be generated.
			WithAdditionalValidations(),
		g.ShowByIDExtendedInFiltering,
		g.ShowByIDLikeFiltering,
	).
	CustomShowOperationWithPairedStructs(
		"DescribeColumns",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-table",
		tableColumnDetails,
		g.NewQueryStruct("DescribeTableColumns").
			Describe().
			SQL("TABLE").
			Name().
			SQLWithCustomFieldName("columnsType", "TYPE = COLUMNS").
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomShowOperationWithPairedStructs(
		"DescribeStage",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-table",
		tableStageDetails,
		g.NewQueryStruct("DescribeTableStage").
			Describe().
			SQL("TABLE").
			Name().
			SQLWithCustomFieldName("stageType", "TYPE = STAGE").
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomShowOperationWithPairedStructs(
		"DescribeSearchOptimization",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-search-optimization",
		tableSearchOptimizationDetails,
		tableDescribeSearchOptimization,
	).
	CustomShowOperationWithPairedStructs(
		"SelectTableConstraints",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/info-schema/table_constraints",
		tableConstraintDetails,
		tableShowConstraints,
	).
	CustomShowOperationWithPairedStructs(
		"SelectCheckConstraints",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/info-schema/check_constraints",
		tableCheckConstraintDetails,
		tableSelectCheckConstraints,
	).
	WithEnums(TableConstraintTypeDef, TableScopeEnumDef, TableKindEnumDef, CloneMomentEnumDef, ReclusterStateEnumDef, MatchTypeEnumDef).
	// TODO(next-pr): drop this allow-list once SDK unit tests move to generated + *_ext_test.go
	WithAllowedGenerationParts(
		g.PartDefault,
		g.PartDto,
		g.PartDtoBuilders,
		g.PartImpl,
		g.PartValidations,
		g.PartEnums,
	)

var tableSetAggregationPolicy = g.NewQueryStruct("TableSetAggregationPolicy").
	SQL("SET").
	Identifier("AggregationPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("AGGREGATION POLICY").Required()).
	ListAssignment("ENTITY KEY", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	OptionalSQL("FORCE").
	WithValidation(g.ValidIdentifier, "AggregationPolicy")

var tableUnsetAggregationPolicy = g.NewQueryStruct("TableUnsetAggregationPolicy").
	SQL("UNSET AGGREGATION POLICY")

var tableSetJoinPolicy = g.NewQueryStruct("TableSetJoinPolicy").
	SQL("SET").
	Identifier("JoinPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("JOIN POLICY").Required()).
	OptionalSQL("FORCE").
	WithValidation(g.ValidIdentifier, "JoinPolicy")

var tableUnsetJoinPolicy = g.NewQueryStruct("TableUnsetJoinPolicy").
	SQL("UNSET JOIN POLICY")

func tableOutOfLineUniquePK() *g.QueryStruct {
	return withOutOfLineConstraintTail(
		g.NewQueryStruct("TableOutOfLineUniquePK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			OptionalSQL("UNIQUE").
			OptionalSQL("PRIMARY KEY").
			PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Parentheses()),
		// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
		// WithValidation(g.ExactlyOneValueSet, "Unique", "PrimaryKey")
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
			PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions()),
		// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
		// WithValidation(g.ValidIdentifier, "References")
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
		OptionalSQL("ENABLE NOVALIDATE")
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ConflictingFields, "EnableValidate", "EnableNovalidate")
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
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
	// WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
	// WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
	// WithValidation(g.ConflictingFields, "Enable", "Disable").
	// WithValidation(g.ConflictingFields, "Validate", "Novalidate").
	// WithValidation(g.ConflictingFields, "Rely", "Norely")
}

func tableColumnInlineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnInlineConstraint").
		OptionalQueryStructField("UniquePK", tableColumnInlineUniquePK(), g.KeywordOptions()).
		OptionalQueryStructField("FK", tableColumnInlineFK(), g.KeywordOptions()).
		OptionalQueryStructField("CH", tableColumnInlineCH(), g.KeywordOptions())
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ExactlyOneValueSet, "UniquePK", "FK", "CH")
}

func tableColumnInlineUniquePK() *g.QueryStruct {
	return withInlineConstraintTail(
		g.NewQueryStruct("TableColumnInlineUniquePK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			OptionalSQL("UNIQUE").
			OptionalSQL("PRIMARY KEY"),
		// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
		// WithValidation(g.ExactlyOneValueSet, "Unique", "PrimaryKey")
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
			PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions()),
		// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
		// WithValidation(g.ValidIdentifier, "References")
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
		OptionalSQL("ENABLE NOVALIDATE")
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ConflictingFields, "EnableValidate", "EnableNovalidate")
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
		OptionalSQL("NORELY")
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
	// WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
	// WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
	// WithValidation(g.ConflictingFields, "Enable", "Disable").
	// WithValidation(g.ConflictingFields, "Validate", "Novalidate").
	// WithValidation(g.ConflictingFields, "Rely", "Norely")
}
