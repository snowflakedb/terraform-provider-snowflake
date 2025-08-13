package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var semanticViewDbRow = g.DbStruct("semanticViewDBRow").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("owner").
	Text("owner_role_type").
	OptionalText("comment")

var semanticView = g.PlainStruct("SemanticView").
	Time("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	Text("Owner").
	Text("OwnerRoleType").
	OptionalText("Comment")

var semanticViewDetailsDbRow = g.DbStruct("semanticViewDetailsRow").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("owner").
	Text("owner_role_type").
	OptionalText("comment")

var semanticViewDetails = g.PlainStruct("SemanticViewDetails").
	Time("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	Text("Owner").
	Text("OwnerRoleType").
	OptionalText("Comment")

var SemanticViewsDef = g.NewInterface(
	"SemanticViews",
	"SemanticView",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-semantic-view",
	g.NewQueryStruct("CreateSemanticView").
		Create().
		OrReplace().
		SQL("SEMANTIC VIEW").
		IfNotExists().
		Name().
		SQL("TABLES").
		ListQueryStructField("logicalTables", logicalTable, g.ListOptions().Required().Parentheses()).
		OptionalComment().
		OptionalCopyGrants().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"), // both can't be used at the same time
	logicalTable,
	semanticViewColumn,
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-semantic-view",
	g.NewQueryStruct("DropSemanticView").
		Drop().
		SQL("SEMANTIC VIEW").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-semantic-view",
	semanticViewDetailsDbRow,
	semanticViewDetails,
	g.NewQueryStruct("DescribeSemanticView").
		Describe().
		SQL("SEMANTIC VIEW").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views",
	semanticViewDbRow,
	semanticView,
	g.NewQueryStruct("ShowSemanticViews").
		Show().
		Terse().
		SQL("SEMANTIC VIEWS").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimit(),
)

var primaryKey = g.NewQueryStruct("PrimaryKeys").
	ListAssignment("PRIMARY KEY", "SemanticViewColumn", g.ParameterOptions().Parentheses().NoEquals())

var uniqueKey = g.NewQueryStruct("UniqueKeys").
	ListAssignment("UNIQUE", "SemanticViewColumn", g.ParameterOptions().Parentheses().NoEquals())

var synonym = g.NewQueryStruct("Synonyms").
	ListAssignment("WITH SYNONYMS", "string", g.ParameterOptions().NoEquals().Parentheses())

var logicalTableAlias = g.NewQueryStruct("LogicalTableAlias").
	Text("LogicalTableAlias", g.KeywordOptions()).
	SQL("AS")

var logicalTable = g.NewQueryStruct("LogicalTable").
	OptionalQueryStructField("logicalTableAlias", logicalTableAlias, g.KeywordOptions()).
	Identifier("TableName", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
	OptionalQueryStructField("primaryKeys", primaryKey, g.ParameterOptions().NoEquals()).
	ListQueryStructField("uniqueKeys", uniqueKey, g.ListOptions().NoEquals()).
	OptionalQueryStructField("synonyms", synonym, g.ParameterOptions().NoEquals()).
	OptionalComment()

var semanticViewColumn = g.NewQueryStruct("SemanticViewColumn").
	Text("Name", g.KeywordOptions().Required())
