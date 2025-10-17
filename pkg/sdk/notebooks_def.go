package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var notebookDbRow = g.DbStruct("notebookRow").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	OptionalText("comment").
	Text("owner").
	OptionalText("query_warehouse").
	Text("url_id").
	Text("owner_role_type").
	Text("code_warehouse")

var notebook = g.PlainStruct("Notebook").
	Time("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	OptionalText("Comment").
	Text("Owner").
	Field("QueryWarehouse", "*AccountObjectIdentifier").
	Text("UrlId").
	Text("OwnerRoleType").
	Field("CodeWarehouse", "AccountObjectIdentifier")

var notebookDetailsDbRow = g.DbStruct("NotebookDetailsRow").
	OptionalText("title").
	Text("main_file").
	OptionalText("query_warehouse").
	Text("url_id").
	Text("default_packages").
	OptionalText("user_packages").
	OptionalText("runtime_name").
	OptionalText("compute_pool").
	Text("owner").
	Text("import_urls").
	Text("external_access_integrations").
	Text("external_access_secrets").
	Text("code_warehouse").
	Number("idle_auto_shutdown_time_seconds").
	Text("runtime_environment_version").
	Text("name").
	OptionalText("comment").
	Text("default_version").
	Text("default_version_name").
	OptionalText("default_version_alias").
	Text("default_version_location_uri").
	OptionalText("default_version_source_location_uri").
	OptionalText("default_version_git_commit_hash").
	Text("last_version_name").
	OptionalText("last_version_alias").
	Text("last_version_location_uri").
	OptionalText("last_version_source_location_uri").
	OptionalText("last_version_git_commit_hash").
	OptionalText("live_version_location_uri")

var notebookDetails = g.PlainStruct("NotebookDetails").
	OptionalText("Title").
	Text("MainFile").
	Field("QueryWarehouse", "*AccountObjectIdentifier").
	Text("UrlId").
	Text("DefaultPackages").
	OptionalText("UserPackages").
	OptionalText("RuntimeName").
	OptionalText("ComputePool").
	Text("Owner").
	Text("ImportUrls").
	Text("ExternalAccessIntegrations").
	Text("ExternalAccessSecrets").
	Text("CodeWarehouse").
	Number("IdleAutoShutdownTimeSeconds").
	Text("RuntimeEnvironmentVersion").
	Text("Name").
	OptionalText("Comment").
	Text("DefaultVersion").
	Text("DefaultVersionName").
	OptionalText("DefaultVersionAlias").
	Text("DefaultVersionLocationUri").
	OptionalText("DefaultVersionSourceLocationUri").
	OptionalText("DefaultVersionGitCommitHash").
	Text("LastVersionName").
	OptionalText("LastVersionAlias").
	Text("LastVersionLocationUri").
	OptionalText("LastVersionSourceLocationUri").
	OptionalText("LastVersionGitCommitHash").
	OptionalText("LiveVersionLocationUri")

var NotebooksDef = g.NewInterface(
	"Notebooks",
	"Notebook",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-notebook",
	g.NewQueryStruct("CreateNotebook").
		Create().
		OrReplace().
		SQL("NOTEBOOK").
		IfNotExists().
		Name().
		PredefinedQueryStructField("From", "*Location", g.ParameterOptions().SQL("FROM").NoQuotes().NoEquals()).
		OptionalTextAssignment("TITLE", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes()).
		OptionalComment().
		OptionalIdentifier("QueryWarehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("QUERY_WAREHOUSE").Equals()).
		OptionalNumberAssignment("IDLE_AUTO_SHUTDOWN_TIME_SECONDS", g.ParameterOptions().NoQuotes()).
		OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("WAREHOUSE").Equals()).
		OptionalTextAssignment("RUNTIME_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalIdentifier("ComputePool", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("COMPUTE_POOL").Equals()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", g.KindOfT[AccountObjectIdentifier](), g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("RUNTIME_ENVIRONMENT_VERSION", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("DEFAULT_VERSION", g.ParameterOptions().NoQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "QueryWarehouse").
		WithValidation(g.ValidIdentifierIfSet, "Warehouse").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
		WithValidation(g.ValidIdentifierIfSet, "ComputePool"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-notebook",
	g.NewQueryStruct("AlterNotebook").
		Alter().
		SQL("NOTEBOOK").
		IfExists().
		Name().
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("NotebookSet").
				OptionalComment().
				OptionalIdentifier("QueryWarehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("QUERY_WAREHOUSE").Equals()).
				OptionalNumberAssignment("IDLE_AUTO_SHUTDOWN_TIME_SECONDS", g.ParameterOptions().NoQuotes()).
				OptionalQueryStructField("Secrets", functionSecretsListWrapper, g.ParameterOptions().SQL("SECRETS").Parentheses()).
				OptionalTextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes()).
				OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("WAREHOUSE").Equals()).
				OptionalTextAssignment("RUNTIME_NAME", g.ParameterOptions().SingleQuotes()).
				OptionalIdentifier("ComputePool", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("COMPUTE_POOL").Equals()).
				ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", g.KindOfT[AccountObjectIdentifier](), g.ParameterOptions().Parentheses()).
				OptionalTextAssignment("RUNTIME_ENVIRONMENT_VERSION", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifierIfSet, "QueryWarehouse").
				WithValidation(g.ValidIdentifierIfSet, "Warehouse").
				WithValidation(g.ValidIdentifierIfSet, "ComputePool").
				WithValidation(g.AtLeastOneValueSet, "Comment", "QueryWarehouse", "IdleAutoShutdownTimeSeconds", "SecretsList", "MainFile", "Warehouse", "RuntimeName", "ComputePool", "ExternalAccessIntegrations", "RuntimeEnvironmentVersion"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("NotebookUnset").
				OptionalSQL("COMMENT").
				OptionalSQL("QUERY_WAREHOUSE").
				OptionalSQL("SECRETS").
				OptionalSQL("WAREHOUSE").
				OptionalSQL("RUNTIME_NAME").
				OptionalSQL("COMPUTE_POOL").
				OptionalSQL("EXTERNAL_ACCESS_INTEGRATIONS").
				OptionalSQL("RUNTIME_ENVIRONMENT_VERSION").
				WithValidation(g.AtLeastOneValueSet, "Comment", "QueryWarehouse", "Secrets", "Warehouse", "RuntimeName", "ComputePool", "ExternalAccessIntegrations", "RuntimeEnvironmentVersion"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-notebook",
	g.NewQueryStruct("DropNotebook").
		Drop().
		SQL("NOTEBOOK").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-notebook",
	notebookDetailsDbRow,
	notebookDetails,
	g.NewQueryStruct("DescribeNotebook").
		Describe().
		SQL("NOTEBOOK").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-notebooks",
	notebookDbRow,
	notebook,
	g.NewQueryStruct("ShowNotebooks").
		Show().
		SQL("NOTEBOOKS").
		OptionalLike().
		OptionalIn().
		OptionalLimitFrom().
		OptionalStartsWith(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
)
