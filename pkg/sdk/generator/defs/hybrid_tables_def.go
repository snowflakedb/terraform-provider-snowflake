package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// hybridTableConstraintAction defines ALTER TABLE ... ADD/DROP/RENAME constraint actions for hybrid tables.
// Per Snowflake docs: https://docs.snowflake.com/en/sql-reference/sql/alter-table#constraint-actions-constraintaction
var hybridTableConstraintAction = g.NewQueryStruct("HybridTableConstraintAction").
	OptionalQueryStructField(
		"Add",
		g.NewQueryStruct("HybridTableConstraintActionAdd").
			SQL("ADD").
			// Uses manually-defined HybridTableOutOfLineConstraint in hybrid_tables_ext.go (rule 13).
			PredefinedQueryStructField("OutOfLineConstraint", "HybridTableOutOfLineConstraint", g.KeywordOptions().Required()),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Drop",
		g.NewQueryStruct("HybridTableConstraintActionDrop").
			SQL("DROP").
			OptionalTextAssignment("ConstraintName", g.ParameterOptions().NoEquals().SQL("CONSTRAINT")).
			PredefinedQueryStructField("ColumnConstraintType", "*ColumnConstraintType", g.KeywordOptions()).
			PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses()),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Rename",
		g.NewQueryStruct("HybridTableConstraintActionRename").
			SQL("RENAME CONSTRAINT").
			Text("OldName", g.KeywordOptions().Required()).
			Text("NewName", g.KeywordOptions().Required().SQL("TO")),
		g.KeywordOptions(),
	)

// hybridTableAlterColumnAction defines ALTER TABLE ... ALTER COLUMN for hybrid tables.
// Per Snowflake docs: "For interactive tables, currently the only clauses that you can use
// with the ALTER TABLE MODIFY COLUMN command are COMMENT and UNSET COMMENT."
var hybridTableAlterColumnAction = g.NewQueryStruct("HybridTableAlterColumnAction").
	SQL("ALTER COLUMN").
	Text("ColumnName", g.KeywordOptions().Required()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	OptionalSQL("UNSET COMMENT")

// hybridTableDropColumnAction defines ALTER TABLE ... DROP COLUMN for hybrid tables.
// Per Snowflake codebase: DROP COLUMN considers index dependencies.
// This is confirmed in the Snowflake codebase but not prominently documented.
// Deviation documented per rule 9.
var hybridTableDropColumnAction = g.NewQueryStruct("HybridTableDropColumnAction").
	SQL("DROP COLUMN").
	Text("ColumnName", g.KeywordOptions().Required())

// hybridTableDropIndexAction defines ALTER TABLE ... DROP INDEX for hybrid tables.
// Per Snowflake codebase: ALTER TABLE <table_name> DROP INDEX <index_name>
// This is an alternative to the standalone DROP INDEX command.
// Confirmed in Snowflake codebase analysis.
var hybridTableDropIndexAction = g.NewQueryStruct("HybridTableDropIndexAction").
	SQL("DROP INDEX").
	Text("IndexName", g.KeywordOptions().Required())

// hybridTableBuildIndexAction defines ALTER TABLE ... BUILD INDEX for hybrid tables.
// Per Snowflake codebase: ALTER TABLE <table_name> BUILD INDEX <index_name> [FENCE | BACKFILL]
// FENCE: Coordinates index builds with ongoing writes.
// BACKFILL: Populates index with existing table data.
// Not documented in official Snowflake docs — discovered via codebase analysis (rule 9).
var hybridTableBuildIndexAction = g.NewQueryStruct("HybridTableBuildIndexAction").
	SQL("BUILD INDEX").
	Text("IndexName", g.KeywordOptions().Required()).
	OptionalSQL("FENCE").
	OptionalSQL("BACKFILL")

// hybridTableSetProperties defines ALTER TABLE ... SET for hybrid tables.
// DATA_RETENTION_TIME_IN_DAYS is confirmed by Snowflake codebase but not in official CREATE docs (rule 9).
var hybridTableSetProperties = g.NewQueryStruct("HybridTableSetProperties").
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())

// hybridTableUnsetProperties defines ALTER TABLE ... UNSET for hybrid tables.
var hybridTableUnsetProperties = g.NewQueryStruct("HybridTableUnsetProperties").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("COMMENT")

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
		// Columns, out-of-line constraints, and indexes are in a parenthesized body.
		// Uses manually-defined HybridTableColumnsConstraintsAndIndexes in hybrid_tables_ext.go
		// because the column/constraint/index structure is too complex for the generator DSL (rule 13).
		PredefinedQueryStructField("ColumnsAndConstraints", "HybridTableColumnsConstraintsAndIndexes", g.ListOptions().Parentheses().Required()).
		// DATA_RETENTION_TIME_IN_DAYS: not in official CREATE HYBRID TABLE docs but confirmed
		// by Snowflake codebase (rule 9). Needs integration test verification.
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-table",
	g.NewQueryStruct("AlterHybridTable").
		Alter().
		SQL("TABLE").
		IfNotExists().
		Name().
		OptionalQueryStructField(
			"ConstraintAction",
			hybridTableConstraintAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
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
			"BuildIndexAction",
			hybridTableBuildIndexAction,
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
			g.KeywordOptions().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "ConstraintAction", "AlterColumnAction", "DropColumnAction", "DropIndexAction", "BuildIndexAction", "Set", "Unset"),
).DropOperation(
	// Note: Snowflake codebase shows DROP HYBRID TABLE syntax, but official docs use DROP TABLE.
	// We use DROP TABLE here (matching docs). Integration tests should verify if DROP HYBRID TABLE
	// also works and document the finding (rule 9).
	"https://docs.snowflake.com/en/sql-reference/sql/drop-table",
	g.NewQueryStruct("DropHybridTable").
		Drop().
		SQL("TABLE").
		IfExists().
		Name().
		OptionalSQL("RESTRICT").
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-hybrid-tables",
	g.DbStruct("hybridTableRow").
		Field("created_on", "time.Time").
		Text("name").
		Text("database_name").
		Text("schema_name").
		OptionalText("owner").
		Number("rows").
		Number("bytes").
		OptionalText("comment").
		OptionalText("owner_role_type"),
	g.PlainStruct("HybridTable").
		Field("CreatedOn", "time.Time").
		Field("Name", "string").
		Field("DatabaseName", "string").
		Field("SchemaName", "string").
		Field("Owner", "string").
		Field("Rows", "int").
		Field("Bytes", "int").
		Field("Comment", "string").
		Field("OwnerRoleType", "string"),
	g.NewQueryStruct("ShowHybridTables").
		Show().
		Terse().
		SQL("HYBRID TABLES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimit(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-table",
	g.DbStruct("hybridTableDetailsRow").
		Text("name").
		Text("type").
		Text("kind").
		// The actual Snowflake column name is "null?" but the generator produces db:"null" from Text("null").
		// This MUST be manually adjusted to db:"null?" in the generated file (rule 13).
		Text("null").
		OptionalText("default").
		Text("primary key").
		Text("unique key").
		OptionalText("check").
		OptionalText("expression").
		OptionalText("comment").
		OptionalText("policy name").
		OptionalText("privacy domain").
		OptionalText("schema_evolution_record"),
	g.PlainStruct("HybridTableDetails").
		Field("Name", "string").
		Field("Type", "string").
		Field("Kind", "string").
		Field("IsNullable", "string").
		Field("Default", "string").
		Field("PrimaryKey", "string").
		Field("UniqueKey", "string").
		Field("Check", "string").
		Field("Expression", "string").
		Field("Comment", "string").
		Field("PolicyName", "string").
		Field("PrivacyDomain", "string").
		Field("SchemaEvolutionRecord", "string"),
	g.NewQueryStruct("DescribeHybridTable").
		Describe().
		SQL("TABLE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
