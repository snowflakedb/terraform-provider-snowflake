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
			// Uses manually-defined HybridTableOutOfLineConstraint in hybrid_tables_ext.go.
			PredefinedQueryStructField("OutOfLineConstraint", "HybridTableOutOfLineConstraint", g.KeywordOptions().Required()),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Drop",
		g.NewQueryStruct("HybridTableConstraintActionDrop").
			SQL("DROP").
			// One of: ConstraintName, PrimaryKey, Unique, ForeignKey (with optional Columns)
			OptionalText("ConstraintName", g.KeywordOptions().SQL("CONSTRAINT")).
			OptionalSQL("PRIMARY KEY").
			OptionalSQL("UNIQUE").
			OptionalSQL("FOREIGN KEY").
			PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses()).
			// CASCADE or RESTRICT
			OptionalSQL("CASCADE").
			OptionalSQL("RESTRICT").
			WithValidation(g.ExactlyOneValueSet, "ConstraintName", "PrimaryKey", "Unique", "ForeignKey").
			WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
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
// Per Snowflake docs: "For hybrid tables, currently the only clauses that you can use
// with the ALTER TABLE MODIFY COLUMN command are COMMENT and UNSET COMMENT."
var hybridTableAlterColumnAction = g.NewQueryStruct("HybridTableAlterColumnAction").
	SQL("ALTER COLUMN").
	Text("ColumnName", g.KeywordOptions().Required()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	OptionalSQL("UNSET COMMENT")

// hybridTableDropColumnAction defines ALTER TABLE ... DROP COLUMN for hybrid tables.
var hybridTableDropColumnAction = g.NewQueryStruct("HybridTableDropColumnAction").
	SQL("DROP COLUMN").
	Text("ColumnName", g.KeywordOptions().Required())

// hybridTableDropIndexAction defines ALTER TABLE ... DROP INDEX for hybrid tables.
// Syntax: ALTER TABLE <table_name> DROP INDEX <index_name>
// This is an alternative to the standalone DROP INDEX command.
var hybridTableDropIndexAction = g.NewQueryStruct("HybridTableDropIndexAction").
	SQL("DROP INDEX").
	Text("IndexName", g.KeywordOptions().Required())


// hybridTableSetProperties defines ALTER TABLE ... SET for hybrid tables.
var hybridTableSetProperties = g.NewQueryStruct("HybridTableSetProperties").
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	WithValidation(g.AtLeastOneValueSet, "DataRetentionTimeInDays", "Comment")

// hybridTableUnsetProperties defines ALTER TABLE ... UNSET for hybrid tables.
var hybridTableUnsetProperties = g.NewQueryStruct("HybridTableUnsetProperties").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "DataRetentionTimeInDays", "Comment")

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
		// because the column/constraint/index structure is too complex for the generator DSL.
		PredefinedQueryStructField("ColumnsAndConstraints", "HybridTableColumnsConstraintsAndIndexes", g.ListOptions().Parentheses().Required()).
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
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
		WithValidation(g.ExactlyOneValueSet, "ConstraintAction", "AlterColumnAction", "DropColumnAction", "DropIndexAction", "Set", "Unset"),
).DropOperation(
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
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		OptionalText("owner").
		Number("rows").
		Number("bytes").
		OptionalText("comment").
		OptionalText("owner_role_type"),
	g.PlainStruct("HybridTable").
		Time("CreatedOn").
		Text("Name").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		Number("Rows").
		Number("Bytes").
		Text("Comment").
		Text("OwnerRoleType"),
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
		// This MUST be manually adjusted to db:"null?" in the generated file.
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
		Text("Name").
		Text("Type").
		Text("Kind").
		Text("IsNullable").
		Text("Default").
		Text("PrimaryKey").
		Text("UniqueKey").
		Text("Check").
		Text("Expression").
		Text("Comment").
		Text("PolicyName").
		Text("PrivacyDomain").
		Text("SchemaEvolutionRecord"),
	g.NewQueryStruct("DescribeHybridTable").
		Describe().
		SQL("TABLE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
