package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var semanticViewDbRow = g.DbStruct("semanticViewsRow").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("owner").
	Text("owner_role_type").
	OptionalText("comment")

var semanticView = g.PlainStruct("semanticView").
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
		OptionalComment().
		OptionalCopyGrants().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-semantic-view",
	g.NewQueryStruct("DropSemanticView").
		Drop().
		SQL("SEMANTIC VIEW").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-semantic-view",
	semanticViewDbRow,
	semanticView,
	g.NewQueryStruct("DescribeSemanticView").
		Describe().
		SQL("SEMANTIC VIEW").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views",
	gitRepositoryDbRow,
	gitRepository,
	g.NewQueryStruct("ShowSemanticViews").
		Show().
		Terse().
		SQL("SEMANTIC VIEWS").
		OptionalLike().
		OptionalIn().
		OptionalLimit(),
)
