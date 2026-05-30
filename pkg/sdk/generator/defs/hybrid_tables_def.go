package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var hybridTableColumn = g.NewQueryStruct("HybridTableColumn").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	PredefinedQueryStructField("Type", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
	PredefinedQueryStructField("InlineConstraint", g.KindOfTPointer[sdkcommons.ColumnInlineConstraint](), g.KeywordOptions()).
	OptionalSQL("NOT NULL").
	PredefinedQueryStructField("DefaultValue", g.KindOfTPointer[sdkcommons.ColumnDefaultValue](), g.KeywordOptions()).
	OptionalTextAssignment("COLLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())

var hybridTableOutOfLineConstraint = g.NewQueryStruct("HybridTableOutOfLineConstraint").
	OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
	PredefinedQueryStructField("Type", g.KindOfT[sdkcommons.ColumnConstraintType](), g.KeywordOptions().Required()).
	PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses()).
	// NOTE: Constraint modifier flags (Enforced, NotEnforced, Deferrable, NotDeferrable,
	// InitiallyDeferred, InitiallyImmediate, Enable, Disable, Validate, Novalidate, Rely, Norely)
	// are not supported on hybrid tables — Snowflake returns "invalid constraint property".
	// Removed from the SDK per PR #4461 review feedback.
	PredefinedQueryStructField("ForeignKey", g.KindOfTPointer[sdkcommons.OutOfLineForeignKey](), g.KeywordOptions())

var hybridTableOutOfLineIndex = g.NewQueryStruct("HybridTableOutOfLineIndex").
	SQL("INDEX").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses().Required()).
	PredefinedQueryStructField("IncludeColumns", "[]string", g.KeywordOptions().Parentheses().SQL("INCLUDE"))

var hybridTableColumnsConstraintsAndIndexes = g.NewQueryStruct("HybridTableColumnsConstraintsAndIndexes").
	ListQueryStructField("Columns", hybridTableColumn, g.KeywordOptions()).
	ListQueryStructField("OutOfLineConstraint", hybridTableOutOfLineConstraint, g.KeywordOptions()).
	ListQueryStructField("OutOfLineIndex", hybridTableOutOfLineIndex, g.KeywordOptions())

var hybridTableAddColumnAction = g.NewQueryStruct("HybridTableAddColumnAction").
	SQL("ADD").
	SQL("COLUMN").
	OptionalSQL("IF NOT EXISTS").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	PredefinedQueryStructField("Type", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
	OptionalTextAssignment("COLLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
	PredefinedQueryStructField("DefaultValue", g.KindOfTPointer[sdkcommons.ColumnDefaultValue](), g.KeywordOptions()).
	PredefinedQueryStructField("InlineConstraint", g.KindOfTPointer[sdkcommons.ColumnInlineConstraint](), g.KeywordOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())

// NOTE: Hybrid tables do not support ALTER TABLE ADD UNIQUE or ADD FOREIGN KEY constraints
// (Snowflake returns: "Unique and foreign-key constraints can only be defined at table creation time").
// The Add action is omitted; only Rename and Drop are supported.
var hybridTableConstraintAction = g.NewQueryStruct("HybridTableConstraintAction").
	OptionalQueryStructField(
		"Rename",
		g.NewQueryStruct("HybridTableConstraintActionRename").
			SQL("RENAME CONSTRAINT").
			Text("OldName", g.KeywordOptions().Required().DoubleQuotes()).
			AssignmentWithFieldName("TO", "string", g.ParameterOptions().NoEquals().DoubleQuotes(), "NewName"),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Drop",
		// NOTE: PRIMARY KEY is not included here — DROP PRIMARY KEY is unsupported on hybrid tables.
		// Snowflake returns an error at runtime; removed per PR #4461 review feedback.
		g.NewQueryStruct("HybridTableConstraintActionDrop").
			SQL("DROP").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "ConstraintName").
			OptionalSQL("UNIQUE").
			OptionalSQL("FOREIGN KEY").
			PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses()).
			OptionalSQL("CASCADE").
			OptionalSQL("RESTRICT").
			WithValidation(g.ExactlyOneValueSet, "ConstraintName", "Unique", "ForeignKey").
			WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
		g.KeywordOptions(),
	).
	WithValidation(g.ExactlyOneValueSet, "Rename", "Drop")

// NOTE: Hybrid tables do not support ALTER COLUMN SET/DROP NOT NULL.
var hybridTableAlterColumnAction = g.NewQueryStruct("HybridTableAlterColumnAction").
	SQL("ALTER").
	SQL("COLUMN").
	Text("ColumnName", g.KeywordOptions().Required().DoubleQuotes()).
	OptionalSQL("DROP DEFAULT").
	PredefinedQueryStructField("SetDefault", g.KindOfTPointer[sdkcommons.SequenceName](), g.ParameterOptions().NoEquals().SQL("SET DEFAULT")).
	PredefinedQueryStructField("Type", "*DataType", g.ParameterOptions().NoEquals().SQL("SET DATA TYPE")).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes()).
	OptionalSQL("UNSET COMMENT").
	WithValidation(g.ExactlyOneValueSet, "DropDefault", "SetDefault", "Type", "Comment", "UnsetComment")

var hybridTableDropColumnAction = g.NewQueryStruct("HybridTableDropColumnAction").
	SQL("DROP COLUMN").
	OptionalSQL("IF EXISTS").
	PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Required())

var hybridTableDropIndexAction = g.NewQueryStruct("HybridTableDropIndexAction").
	SQL("DROP INDEX").
	OptionalSQL("IF EXISTS").
	Text("IndexName", g.KeywordOptions().Required().DoubleQuotes())

var hybridTableClusteringAction = g.NewQueryStruct("HybridTableClusteringAction").
	PredefinedQueryStructField("ClusterBy", "[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
	OptionalQueryStructField(
		"Recluster",
		g.NewQueryStruct("HybridTableReclusterAction").
			SQL("RECLUSTER").
			OptionalNumberAssignment("MAX_SIZE", g.ParameterOptions()).
			OptionalTextAssignment("WHERE", g.ParameterOptions().NoEquals()),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"ChangeReclusterState",
		g.NewQueryStruct("HybridTableReclusterChangeState").
			PredefinedQueryStructField("State", g.KindOfTPointer[sdkcommons.ReclusterState](), g.KeywordOptions()).
			SQL("RECLUSTER"),
		g.KeywordOptions(),
	).
	OptionalSQL("DROP CLUSTERING KEY").
	WithValidation(g.ExactlyOneValueSet, "ClusterBy", "Recluster", "ChangeReclusterState", "DropClusteringKey")

// NOTE: Hybrid tables do not support CHANGE_TRACKING, DEFAULT_DDL_COLLATION, ENABLE_SCHEMA_EVOLUTION,
// CONTACT, or ROW_TIMESTAMP in ALTER TABLE SET (per Snowflake documentation and runtime behavior).
var hybridTableSetProperties = g.NewQueryStruct("HybridTableSetProperties").
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	WithValidation(g.AtLeastOneValueSet, "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "Comment")

// NOTE: Multi-property `ALTER TABLE ... UNSET` on hybrid tables requires comma-separated
// property names (`UNSET A, B, C`) — the bare-keyword form (`UNSET A B C`) emitted by the
// generator's default `keyword` rendering is rejected by the parser. Applying
// `g.ListOptions().NoParentheses().SQL("UNSET")` to the parent field on
// AlterHybridTableOptions causes the SQL builder to comma-join the children.
// Mirrors NetworkPolicyUnset in pkg/sdk/network_policies_gen.go:74. Verified on preprod6.
var hybridTableUnsetProperties = g.NewQueryStruct("HybridTableUnsetProperties").
	OptionalSQL("COMMENT").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
	WithValidation(g.AtLeastOneValueSet, "Comment", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays")

// NOTE: After running make generate-sdk, the hybridTableDetailsRow.Null field in
// hybrid_tables_gen.go will have tag db:"null" (the generator strips '?'). It must
// be manually corrected back to db:"null?" because the DESCRIBE TABLE output column
// is literally named "null?" in Snowflake. The convert() method in hybrid_tables_impl_gen.go
// uses r.Null == "Y" to derive the IsNullable bool.
//
// NOTE: HybridTableDetails carries a Collation *string field that is not derivable from
// the DESCRIBE TABLE output via the generator (DESCRIBE returns the collation glued onto
// the "type" column, e.g. "VARCHAR(200) COLLATE 'en-ci'"). After running make generate-sdk,
// the Collation field must be re-added manually to HybridTableDetails in hybrid_tables_gen.go,
// and the convert() method in hybrid_tables_impl_gen.go must call splitTypeAndCollation()
// (in hybrid_tables_ext.go) to populate Type (without the suffix) and Collation. This
// mirrors pkg/sdk/tables.go:736 for classic tables.
var hybridTablesDef = g.NewInterface(
	"HybridTables",
	"HybridTable",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table",
	g.NewQueryStruct("CreateHybridTable").
		Create().
		OrReplace().
		SQL("HYBRID TABLE").
		IfNotExists().
		Name().
		QueryStructField("ColumnsAndConstraints", hybridTableColumnsConstraintsAndIndexes, g.ListOptions().Parentheses()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-table",
	g.NewQueryStruct("AlterHybridTable").
		Alter().
		SQL("TABLE").
		IfExists().
		Name().
		OptionalIdentifier("NewName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField(
			"AddColumnAction",
			hybridTableAddColumnAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"ConstraintAction",
			hybridTableConstraintAction,
			g.KeywordOptions(),
		).
		ListQueryStructField(
			"AlterColumnAction",
			hybridTableAlterColumnAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"DropColumnAction",
			hybridTableDropColumnAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"DropIndexAction",
			hybridTableDropIndexAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"ClusteringAction",
			hybridTableClusteringAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"Set",
			hybridTableSetProperties,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			hybridTableUnsetProperties,
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "NewName", "AddColumnAction", "ConstraintAction", "AlterColumnAction", "DropColumnAction", "DropIndexAction", "ClusteringAction", "Set", "Unset"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-table",
	g.NewQueryStruct("DropHybridTable").
		Drop().
		SQL("TABLE").
		IfExists().
		Name().
		OptionalSQL("CASCADE").
		OptionalSQL("RESTRICT").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-hybrid-tables",
	g.StructPair("hybridTableRow", "HybridTable").
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		OptionalText("owner", g.WithRequiredInPlain()).
		OptionalNumber("rows").
		OptionalNumber("bytes").
		OptionalText("comment", g.WithRequiredInPlain()).
		OptionalText("owner_role_type", g.WithRequiredInPlain()),
	g.NewQueryStruct("ShowHybridTables").
		Show().
		Terse().
		SQL("HYBRID TABLES").
		OptionalLike().
		OptionalTableIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-table",
	g.StructPair("hybridTableDetailsRow", "HybridTableDetails").
		Text("name").
		Text("type").
		Text("kind").
		Text("null", g.WithPlainFieldName("IsNullable")).
		OptionalText("default", g.WithRequiredInPlain()).
		Text("primary key", g.WithPlainFieldName("PrimaryKey")).
		Text("unique key", g.WithPlainFieldName("UniqueKey")).
		OptionalText("check", g.WithRequiredInPlain()).
		OptionalText("expression", g.WithRequiredInPlain()).
		OptionalText("comment", g.WithRequiredInPlain()).
		OptionalText("policy name", g.WithPlainFieldName("PolicyName"), g.WithRequiredInPlain()).
		OptionalText("privacy domain", g.WithPlainFieldName("PrivacyDomain"), g.WithRequiredInPlain()).
		OptionalText("schema_evolution_record", g.WithRequiredInPlain()),
	g.NewQueryStruct("DescribeHybridTable").
		Describe().
		SQL("TABLE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateIndex",
	"https://docs.snowflake.com/en/sql-reference/sql/create-index",
	g.NewQueryStruct("CreateHybridTableIndex").
		Create().
		OrReplace().
		SQL("INDEX").
		IfNotExists().
		Name().
		SQL("ON").
		Identifier("TableName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses().Required()).
		PredefinedQueryStructField("IncludeColumns", "[]string", g.KeywordOptions().Parentheses().SQL("INCLUDE")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "TableName").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).CustomOperation(
	"DropIndex",
	"https://docs.snowflake.com/en/sql-reference/sql/drop-index",
	g.NewQueryStruct("DropHybridTableIndex").
		Drop().
		SQL("INDEX").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).CustomShowOperationWithPairedStructs(
	"ShowIndexes",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-indexes",
	g.StructPair("hybridTableIndexRow", "HybridTableIndex").
		Time("created_on").
		Text("name").
		Field("is_unique", "sql.NullString", "*bool").
		OptionalText("columns").
		OptionalText("included_columns", g.WithRequiredInPlain()).
		Text("table", g.WithPlainFieldName("TableName")).
		Text("database_name").
		Text("schema_name").
		OptionalText("owner", g.WithRequiredInPlain()).
		OptionalText("owner_role_type", g.WithRequiredInPlain()),
	g.NewQueryStruct("ShowHybridTableIndexes").
		Show().
		SQL("INDEXES").
		OptionalLike().
		OptionalTableIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
).WithCustomInterfaceMethod(
	"ShowParameters",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
	"[]*Parameter", "error",
)
