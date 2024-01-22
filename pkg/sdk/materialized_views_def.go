package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var materializedViewColumn = g.NewQueryStruct("MaterializedViewColumn").
	Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes().NoEquals())

var materializedViewColumnMaskingPolicy = g.NewQueryStruct("MaterializedViewColumnMaskingPolicy").
	Text("Name", g.KeywordOptions().Required()).
	Identifier("MaskingPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	NamedListWithParens("USING", g.KindOfT[string](), nil). // TODO: double quotes here?
	OptionalTags()

var materializedViewRowAccessPolicy = g.NewQueryStruct("MaterializedViewRowAccessPolicy").
	Identifier("RowAccessPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	NamedListWithParens("ON", g.KindOfT[string](), g.KeywordOptions().Required()). // TODO: double quotes here?
	WithValidation(g.ValidIdentifier, "RowAccessPolicy")

var materializedViewSet = g.NewQueryStruct("MaterializedViewSet").
	OptionalSQL("SECURE").
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Secure", "Comment")

var materializedViewUnset = g.NewQueryStruct("MaterializedViewUnset").
	OptionalSQL("SECURE").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Secure", "Comment")

var materializedViewDbRow = g.DbStruct("materializedViewDBRow").
	Text("created_on").
	Text("name").
	OptionalText("reserved").
	Text("database_name").
	Text("schema_name").
	OptionalText("cluster_by").
	Number("rows").
	Number("bytes").
	Text("source_database_name").
	Text("source_schema_name").
	Text("source_table_name").
	Time("refreshed_on").
	Time("compacted_on").
	OptionalText("owner").
	Bool("invalid").
	OptionalText("invalid_reason").
	Text("behind_by").
	OptionalText("comment").
	Text("text").
	Bool("is_secure").
	Bool("automatic_clustering").
	OptionalText("owner_role_type").
	OptionalText("budget")

var materializedView = g.PlainStruct("MaterializedView").
	Text("CreatedOn").
	Text("Name").
	OptionalText("Reserved").
	Text("DatabaseName").
	Text("SchemaName").
	OptionalText("ClusterBy").
	Number("Rows").
	Number("Bytes").
	Text("SourceDatabaseName").
	Text("SourceSchemaName").
	Text("SourceTableName").
	Time("RefreshedOn").
	Time("CompactedOn").
	OptionalText("Owner").
	Bool("Invalid").
	OptionalText("InvalidReason").
	Text("BehindBy").
	OptionalText("Comment").
	Text("Text").
	Bool("IsSecure").
	Bool("AutomaticClustering").
	OptionalText("OwnerRoleType").
	OptionalText("Budget")

var materializedViewDetailsDbRow = g.DbStruct("materializedViewDetailsRow").
	Text("name").
	Field("type", "DataType").
	Text("kind").
	Text("null").
	OptionalText("default").
	Text("primary key").
	Text("unique key").
	OptionalText("check").
	OptionalText("expression").
	OptionalText("comment")

var materializedViewDetails = g.PlainStruct("MaterializedViewDetails").
	Text("Name").
	Field("Type", "DataType").
	Text("Kind").
	Bool("IsNullable").
	OptionalText("Default").
	Bool("IsPrimary").
	Bool("IsUnique").
	OptionalBool("Check").
	OptionalText("Expression").
	OptionalText("Comment")

var MaterializedViewsDef = g.NewInterface(
	"MaterializedViews",
	"MaterializedView",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-materialized-view",
		g.NewQueryStruct("CreateMaterializedView").
			Create().
			OrReplace().
			OptionalSQL("SECURE").
			SQL("MATERIALIZED VIEW").
			IfNotExists().
			Name().
			OptionalCopyGrants().
			ListQueryStructField("Columns", materializedViewColumn, g.ListOptions().Parentheses()).
			ListQueryStructField("ColumnsMaskingPolicies", materializedViewColumnMaskingPolicy, g.ListOptions().NoParentheses().NoEquals()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			// In the current docs ROW ACCESS POLICY and TAG are specified twice.
			// It is a mistake probably so here they are present only once.
			OptionalQueryStructField("RowAccessPolicy", materializedViewRowAccessPolicy, g.KeywordOptions()).
			OptionalTags().
			NamedListWithParens("CLUSTER BY", g.KindOfT[string](), g.KeywordOptions()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-materialized-view",
		g.NewQueryStruct("AlterMaterializedView").
			Alter().
			SQL("MATERIALIZED VIEW").
			Name().
			OptionalIdentifier("RenameTo", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			NamedListWithParens("CLUSTER BY", g.KindOfT[string](), g.KeywordOptions()).
			OptionalSQL("DROP CLUSTERING KEY").
			OptionalSQL("SUSPEND RECLUSTER").
			OptionalSQL("RESUME RECLUSTER").
			OptionalSQL("SUSPEND").
			OptionalSQL("RESUME").
			OptionalQueryStructField("Set", materializedViewSet, g.KeywordOptions().SQL("SET")).
			OptionalQueryStructField("Unset", materializedViewUnset, g.KeywordOptions().SQL("UNSET")).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "ClusterBy", "DropClusteringKey", "SuspendRecluster", "ResumeRecluster", "Suspend", "Resume", "Set", "Unset"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-materialized-view",
		g.NewQueryStruct("DropMaterializedView").
			Drop().
			SQL("MATERIALIZED VIEW").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-materialized-views",
		materializedViewDbRow,
		materializedView,
		g.NewQueryStruct("ShowMaterializedViews").
			Show().
			SQL("MATERIALIZED VIEWS").
			OptionalLike().
			OptionalIn(),
	).
	ShowByIdOperation().
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-materialized-view",
		materializedViewDetailsDbRow,
		materializedViewDetails,
		g.NewQueryStruct("DescribeMaterializedView").
			Describe().
			SQL("MATERIALIZED VIEW").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
