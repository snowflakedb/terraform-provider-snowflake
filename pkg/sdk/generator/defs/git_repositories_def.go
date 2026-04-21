package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var gitRepositoryPairs = g.StructPair("gitRepositoriesRow", "GitRepository").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("origin").
	AccountObjectIdentifier("api_integration", g.WithPlainFieldName("ApiIntegration")).
	Field("git_credentials", "sql.NullString", "*SchemaObjectIdentifier", g.WithPlainFieldName("GitCredentials")).
	Text("owner").
	Text("owner_role_type").
	OptionalText("comment").
	OptionalTime("last_fetched_at")

var gitBranchesPairs = g.StructPair("gitBranchesRow", "GitBranch").
	Text("name").
	Text("path").
	Text("checkouts").
	Text("commit_hash")

var gitTagsPairs = g.StructPair("gitTagsRow", "GitTag").
	Text("name").
	Text("path").
	Text("commit_hash").
	Text("author").
	Text("message")

var gitRepositoriesDef = g.NewInterface(
	"GitRepositories",
	"GitRepository",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-git-repository",
	g.NewQueryStruct("CreateGitRepository").
		Create().
		OrReplace().
		SQL("GIT REPOSITORY").
		IfNotExists().
		Name().
		TextAssignment("ORIGIN", g.ParameterOptions().SingleQuotes()).
		Identifier("ApiIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("API_INTEGRATION").Equals().Required()).
		OptionalIdentifier("GitCredentials", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("GIT_CREDENTIALS").Equals()).
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
				OptionalIdentifier("ApiIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("API_INTEGRATION").Equals()).
				OptionalIdentifier("GitCredentials", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("GIT_CREDENTIALS").Equals()).
				OptionalComment().
				WithValidation(g.ValidIdentifierIfSet, "ApiIntegration").
				WithValidation(g.ValidIdentifierIfSet, "GitCredentials"),
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
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-git-repository",
	gitRepositoryPairs,
	g.NewQueryStruct("DescribeGitRepository").
		Describe().
		SQL("GIT REPOSITORY").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-git-repositories",
	gitRepositoryPairs,
	g.NewQueryStruct("ShowGitRepositories").
		Show().
		SQL("GIT REPOSITORIES").
		OptionalLike().
		OptionalIn().
		OptionalLimit(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).CustomShowOperationWithPairedStructs(
	"ShowGitBranches",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-git-branches",
	gitBranchesPairs,
	g.NewQueryStruct("ShowGitBranches").
		SQL("SHOW GIT BRANCHES").
		OptionalLike().
		SQL("IN").
		OptionalSQL("GIT REPOSITORY").
		Name(),
).CustomShowOperationWithPairedStructs(
	"ShowGitTags",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-git-tags",
	gitTagsPairs,
	g.NewQueryStruct("ShowGitTags").
		SQL("SHOW GIT TAGS").
		OptionalLike().
		SQL("IN").
		OptionalSQL("GIT REPOSITORY").
		Name(),
)
