package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var gitRepositoryDbRow = g.DbStruct("gitRepositoriesRow").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("origin").
	Text("api_integration").
	Text("git_credentials").
	Text("owner").
	Text("owner_role_type").
	Text("comment")

var gitRepository = g.PlainStruct("GitRepository").
	Time("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	Text("Origin").
	Text("ApiIntegration").
	Text("GitCredentials").
	Text("Owner").
	Text("OwnerRoleType").
	Text("Comment")

var apiIntegrationIdentifierOptions = g.IdentifierOptions().SQL("API_INTEGRATION =")
var gitCredentialsIdentifierOptions = g.IdentifierOptions().SQL("GIT_CREDENTIALS =")

var GitRepositoriesDef = g.NewInterface(
	"GitRepositories",
	"GitRepository",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-git-repository",
	g.NewQueryStruct("CreateGitRepository").
		Create().
		OrReplace().
		SQL("GIT REPOSITORY").
		IfNotExists().
		Name().
		TextAssignment("ORIGIN", g.ParameterOptions().SingleQuotes()).
		Identifier("ApiIntegration", g.KindOfT[AccountObjectIdentifier](), apiIntegrationIdentifierOptions.Required()).
		OptionalIdentifier("GitCredentials", g.KindOfT[AccountObjectIdentifier](), gitCredentialsIdentifierOptions).
		OptionalComment().
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-git-repository",
	g.NewQueryStruct("AlterGitRepository").
		Alter().
		SQL("GIT REPOSITORY").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("GitRepositorySet").
				OptionalIdentifier("ApiIntegration", g.KindOfT[AccountObjectIdentifier](), apiIntegrationIdentifierOptions).
				OptionalIdentifier("GitCredentials", g.KindOfT[AccountObjectIdentifier](), gitCredentialsIdentifierOptions).
				OptionalComment(),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("GitRepositoryUnset").
				OptionalSQL("GIT_CREDENTIALS").
				OptionalSQL("COMMENT"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalSQL("FETCH").
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags", "Fetch"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-git-repository",
	g.NewQueryStruct("DropGitRepository").
		Drop().
		SQL("GIT REPOSITORY").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-git-repository",
	gitRepositoryDbRow,
	gitRepository,
	g.NewQueryStruct("DescribeGitRepository").
		Describe().
		SQL("GIT REPOSITORY").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-git-repositories",
	gitRepositoryDbRow,
	gitRepository,
	g.NewQueryStruct("ShowGitRepositories").
		Show().
		SQL("GIT REPOSITORIES").
		OptionalLike().
		OptionalIn(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).CustomShowOperation(
	"ShowGitBranches",
	"https://docs.snowflake.com/en/sql-reference/sql/show-git-branches",
	g.DbStruct("gitBranchesRow").
		Text("name").
		Text("path").
		Text("checkouts").
		Text("commit_hash"),
	g.PlainStruct("GitBranch").
		Text("Name").
		Text("Path").
		Text("Checkouts").
		Text("CommitHash"),
	g.NewQueryStruct("ShowGitBranches").
		Show().
		SQL("GIT BRANCHES").
		OptionalLike().
		SQL("IN").
		OptionalSQL("GIT REPOSITORY").
		Name(),
).CustomShowOperation(
	"ShowGitTags",
	"https://docs.snowflake.com/en/sql-reference/sql/show-git-tags",
	g.DbStruct("gitTagsRow").
		Text("name").
		Text("path").
		Text("commit_hash").
		Text("author").
		Text("message"),
	g.PlainStruct("GitTag").
		Text("Name").
		Text("Path").
		Text("CommitHash").
		Text("Author").
		Text("Message"),
	g.NewQueryStruct("ShowGitTags").
		Show().
		SQL("GIT TAGS").
		OptionalLike().
		SQL("IN").
		OptionalSQL("GIT REPOSITORY").
		Name(),
)
