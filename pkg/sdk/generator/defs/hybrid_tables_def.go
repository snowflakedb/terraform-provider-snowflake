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
	OptionalText("Name", g.KeywordOptions().DoubleQuotes().SQL("CONSTRAINT")).
	PredefinedQueryStructField("Type", g.KindOfT[sdkcommons.ColumnConstraintType](), g.KeywordOptions().Required()).
	PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses().DoubleQuotes()).
	PredefinedQueryStructField("ForeignKey", g.KindOfTPointer[sdkcommons.OutOfLineForeignKey](), g.KeywordOptions()).
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

var hybridTableOutOfLineIndex = g.NewQueryStruct("HybridTableOutOfLineIndex").
	SQL("INDEX").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses().Required().DoubleQuotes()).
	PredefinedQueryStructField("IncludeColumns", "[]string", g.KeywordOptions().Parentheses().DoubleQuotes().SQL("INCLUDE"))

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

var hybridTableConstraintAction = g.NewQueryStruct("HybridTableConstraintAction").
	OptionalQueryStructField(
		"Add",
		g.NewQueryStruct("HybridTableConstraintActionAdd").
			SQL("ADD").
			QueryStructField("OutOfLineConstraint", hybridTableOutOfLineConstraint, g.KeywordOptions()),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Rename",
		g.NewQueryStruct("HybridTableConstraintActionRename").
			SQL("RENAME CONSTRAINT").
			Text("OldName", g.KeywordOptions().Required().DoubleQuotes()).
			Text("NewName", g.KeywordOptions().Required().DoubleQuotes().SQL("TO")),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Drop",
		g.NewQueryStruct("HybridTableConstraintActionDrop").
			SQL("DROP").
			OptionalText("ConstraintName", g.KeywordOptions().DoubleQuotes().SQL("CONSTRAINT")).
			OptionalSQL("PRIMARY KEY").
			OptionalSQL("UNIQUE").
			OptionalSQL("FOREIGN KEY").
			PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses().DoubleQuotes()).
			OptionalSQL("CASCADE").
			OptionalSQL("RESTRICT").
			WithValidation(g.ExactlyOneValueSet, "ConstraintName", "PrimaryKey", "Unique", "ForeignKey").
			WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
		g.KeywordOptions(),
	).
	WithValidation(g.ExactlyOneValueSet, "Add", "Rename", "Drop")

var hybridTableAlterColumnAction = g.NewQueryStruct("HybridTableAlterColumnAction").
	SQL("ALTER").
	SQL("COLUMN").
	Text("ColumnName", g.KeywordOptions().Required().DoubleQuotes()).
	OptionalSQL("DROP DEFAULT").
	PredefinedQueryStructField("SetDefault", g.KindOfTPointer[sdkcommons.SequenceName](), g.ParameterOptions().NoEquals().SQL("SET DEFAULT")).
	OptionalQueryStructField(
		"NotNullConstraint",
		g.NewQueryStruct("HybridTableColumnNotNullConstraint").
			OptionalSQL("SET NOT NULL").
			OptionalSQL("DROP NOT NULL").
			WithValidation(g.ExactlyOneValueSet, "SetNotNull", "DropNotNull"),
		g.KeywordOptions(),
	).
	PredefinedQueryStructField("Type", "*DataType", g.ParameterOptions().NoEquals().SQL("SET DATA TYPE")).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes()).
	OptionalSQL("UNSET COMMENT").
	WithValidation(g.ExactlyOneValueSet, "DropDefault", "SetDefault", "NotNullConstraint", "Type", "Comment", "UnsetComment")

var hybridTableDropColumnAction = g.NewQueryStruct("HybridTableDropColumnAction").
	SQL("DROP COLUMN").
	OptionalSQL("IF EXISTS").
	PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Required().DoubleQuotes())

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

var hybridTableSetProperties = g.NewQueryStruct("HybridTableSetProperties").
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
	OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("ENABLE_SCHEMA_EVOLUTION", g.ParameterOptions()).
	PredefinedQueryStructField("Contact", g.KindOfTSlice[sdkcommons.TableContact](), g.KeywordOptions().SQL("CONTACT")).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("ROW_TIMESTAMP", g.ParameterOptions()).
	WithValidation(g.AtLeastOneValueSet, "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "DefaultDdlCollation", "EnableSchemaEvolution", "Contact", "Comment", "RowTimestamp")

var hybridTableUnsetProperties = g.NewQueryStruct("HybridTableUnsetProperties").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
	OptionalSQL("CHANGE_TRACKING").
	OptionalSQL("DEFAULT_DDL_COLLATION").
	OptionalSQL("ENABLE_SCHEMA_EVOLUTION").
	OptionalText("ContactPurpose", g.KeywordOptions().SQL("CONTACT")).
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "DefaultDdlCollation", "EnableSchemaEvolution", "ContactPurpose", "Comment")

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
			g.KeywordOptions().SQL("UNSET"),
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
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-hybrid-tables",
	g.DbStruct("hybridTableRow").
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		OptionalText("owner").
		OptionalNumber("rows").
		OptionalNumber("bytes").
		OptionalText("comment").
		OptionalText("owner_role_type"),
	g.PlainStruct("HybridTable").
		Time("CreatedOn").
		Text("Name").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		OptionalNumber("Rows").
		OptionalNumber("Bytes").
		Text("Comment").
		Text("OwnerRoleType"),
	g.NewQueryStruct("ShowHybridTables").
		Show().
		Terse().
		SQL("HYBRID TABLES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
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
		PredefinedQueryStructField("Columns", "[]string", g.KeywordOptions().Parentheses().Required().DoubleQuotes()).
		PredefinedQueryStructField("IncludeColumns", "[]string", g.KeywordOptions().Parentheses().DoubleQuotes().SQL("INCLUDE")).
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
).CustomShowOperation(
	"ShowIndexes",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-indexes",
	g.DbStruct("hybridTableIndexRow").
		Time("created_on").
		Text("name").
		OptionalText("is_unique").
		OptionalText("columns").
		OptionalText("included_columns").
		Text("table").
		Text("database_name").
		Text("schema_name").
		OptionalText("owner").
		OptionalText("owner_role_type"),
	g.PlainStruct("HybridTableIndex").
		Time("CreatedOn").
		Text("Name").
		OptionalBool("IsUnique").
		OptionalText("Columns").
		Text("IncludedColumns").
		Text("TableName").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		Text("OwnerRoleType"),
	g.NewQueryStruct("ShowHybridTableIndexes").
		Show().
		SQL("INDEXES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
)
