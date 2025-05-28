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
	OptionalText("git_credentials").
	Text("owner").
	Text("owner_role_type").
	OptionalText("comment").
	Time("last_fetched_at")

var gitRepository = g.PlainStruct("GitRepository").
	Time("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	Text("Origin").
	Field("ApiIntegration", "*AccountObjectIdentifier").
	Field("GitCredentials", "*SchemaObjectIdentifier").
	Text("Owner").
	Text("OwnerRoleType").
	OptionalText("Comment").
	Time("LastFetchedAt")

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
		Identifier("ApiIntegration", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("API_INTEGRATION").Equals().Required()).
		OptionalIdentifier("GitCredentials", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("GIT_CREDENTIALS").Equals()).
		OptionalComment().
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "ApiIntegration").
		WithValidation(g.ValidIdentifierIfSet, "GitCredentials").
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
				OptionalIdentifier("ApiIntegration", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("API_INTEGRATION").Equals()).
				OptionalIdentifier("GitCredentials", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("GIT_CREDENTIALS").Equals()).
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
		OptionalIn().
		OptionalLimit(),
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
		SQL("SHOW GIT BRANCHES").
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
		SQL("SHOW GIT TAGS").
		OptionalLike().
		SQL("IN").
		OptionalSQL("GIT REPOSITORY").
		Name(),
)
