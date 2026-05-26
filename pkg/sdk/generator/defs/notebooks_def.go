package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var notebookPairs = g.StructPair("notebookRow", "Notebook").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	OptionalText("comment").
	Text("owner").
	Field("query_warehouse", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("QueryWarehouse")).
	Text("url_id").
	Text("owner_role_type").
	Field("code_warehouse", "string", "AccountObjectIdentifier", g.WithPlainFieldName("CodeWarehouse")).
	WithConvertGeneration()

var notebookDetailsPairs = g.StructPair("NotebookDetailsRow", "NotebookDetails").
	OptionalText("title").
	Text("main_file").
	Field("query_warehouse", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("QueryWarehouse")).
	Text("url_id").
	Text("default_packages").
	OptionalText("user_packages").
	OptionalText("runtime_name").
	Field("compute_pool", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("ComputePool")).
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
	OptionalText("live_version_location_uri").
	WithConvertGeneration()

var notebooksDef = g.NewInterface(
	"Notebooks",
	"Notebook",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-notebook",
	g.NewQueryStruct("CreateNotebook").
		Create().
		OrReplace().
		SQL("NOTEBOOK").
		IfNotExists().
		Name().
		PredefinedQueryStructField("From", "*Location", g.ParameterOptions().SQL("FROM").SingleQuotes().NoEquals()).
		OptionalTextAssignment("TITLE", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes()).
		OptionalComment().
		OptionalIdentifier("QueryWarehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("QUERY_WAREHOUSE").Equals()).
		OptionalNumberAssignment("IDLE_AUTO_SHUTDOWN_TIME_SECONDS", g.ParameterOptions().NoQuotes()).
		OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("WAREHOUSE").Equals()).
		OptionalTextAssignment("RUNTIME_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalIdentifier("ComputePool", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("COMPUTE_POOL").Equals()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("RUNTIME_ENVIRONMENT_VERSION", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("DEFAULT_VERSION", g.ParameterOptions().NoQuotes()).
		OptionalQueryStructField("Secrets", functionSecretsListWrapper, g.ParameterOptions().SQL("SECRETS").Parentheses()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "QueryWarehouse").
		WithValidation(g.ValidIdentifierIfSet, "Warehouse").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
		WithValidation(g.ValidIdentifierIfSet, "ComputePool").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-notebook",
	g.NewQueryStruct("AlterNotebook").
		Alter().
		SQL("NOTEBOOK").
		IfExists().
		Name().
		Identifier("RenameTo", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("NotebookSet").
				OptionalComment().
				OptionalIdentifier("QueryWarehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("QUERY_WAREHOUSE").Equals()).
				OptionalNumberAssignment("IDLE_AUTO_SHUTDOWN_TIME_SECONDS", g.ParameterOptions().NoQuotes()).
				OptionalQueryStructField("Secrets", functionSecretsListWrapper, g.ParameterOptions().SQL("SECRETS").Parentheses()).
				OptionalTextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes()).
				OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("WAREHOUSE").Equals()).
				OptionalTextAssignment("RUNTIME_NAME", g.ParameterOptions().SingleQuotes()).
				OptionalIdentifier("ComputePool", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("COMPUTE_POOL").Equals()).
				ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.ParameterOptions().Parentheses()).
				OptionalTextAssignment("RUNTIME_ENVIRONMENT_VERSION", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifierIfSet, "QueryWarehouse").
				WithValidation(g.ValidIdentifierIfSet, "Warehouse").
				WithValidation(g.ValidIdentifierIfSet, "ComputePool").
				WithValidation(g.AtLeastOneValueSet, "Comment", "QueryWarehouse", "IdleAutoShutdownTimeSeconds", "Secrets", "MainFile", "Warehouse", "RuntimeName", "ComputePool", "ExternalAccessIntegrations", "RuntimeEnvironmentVersion").
				WithAdditionalValidations(),
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
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags", "RenameTo"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-notebook",
	g.NewQueryStruct("DropNotebook").
		Drop().
		SQL("NOTEBOOK").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-notebook",
	notebookDetailsPairs,
	g.NewQueryStruct("DescribeNotebook").
		Describe().
		SQL("NOTEBOOK").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-notebooks",
	notebookPairs,
	g.NewQueryStruct("ShowNotebooks").
		Show().
		SQL("NOTEBOOKS").
		OptionalLike().
		OptionalIn().
		OptionalLimitFrom().
		OptionalStartsWith(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
)
