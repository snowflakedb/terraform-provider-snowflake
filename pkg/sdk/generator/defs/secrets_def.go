package defs

import (
	"fmt"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var secretPairs = g.StructPair("secretDBRow", "Secret").
	Time("created_on").
	Text("name").
	Text("schema_name").
	Text("database_name").
	Text("owner").
	OptionalText("comment").
	Text("secret_type").
	Field("oauth_scopes", "sql.NullString", "[]string").
	Text("owner_role_type")

var secretDetailsPairs = g.StructPair("secretDetailsDBRow", "SecretDetails").
	Time("created_on").
	Text("name").
	Text("schema_name").
	Text("database_name").
	Text("owner").
	OptionalText("comment").
	Text("secret_type").
	OptionalText("username").
	Field("oauth_access_token_expiry_time", "*time.Time", "*time.Time").
	Field("oauth_refresh_token_expiry_time", "*time.Time", "*time.Time").
	Field("oauth_scopes", "sql.NullString", "[]string").
	OptionalText("integration_name")

var secretsApiIntegrationScopeDef = g.NewQueryStruct("ApiIntegrationScope").
	Text("Scope", g.KeywordOptions().SingleQuotes().Required())

var oauthScopesListDef = g.NewQueryStruct("OauthScopesList").List("OauthScopesList", "ApiIntegrationScope", g.ListOptions().Required().MustParentheses())

var secretSet = g.NewQueryStruct("SecretSet").
	OptionalComment().
	OptionalQueryStructField(
		"SetForFlow",
		g.NewQueryStruct("SetForFlow").
			OptionalQueryStructField(
				"SetForOAuthClientCredentials",
				g.NewQueryStruct("SetForOAuthClientCredentials").
					OptionalQueryStructField("OauthScopes", oauthScopesListDef, g.ParameterOptions().SQL("OAUTH_SCOPES").Parentheses()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"SetForOAuthAuthorization",
				g.NewQueryStruct("SetForOAuthAuthorization").
					OptionalTextAssignment("OAUTH_REFRESH_TOKEN", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("OAUTH_REFRESH_TOKEN_EXPIRY_TIME", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"SetForBasicAuthentication",
				g.NewQueryStruct("SetForBasicAuthentication").
					OptionalTextAssignment("USERNAME", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("PASSWORD", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"SetForGenericString",
				g.NewQueryStruct("SetForGenericString").
					OptionalTextAssignment("SECRET_STRING", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			WithValidation(g.ExactlyOneValueSet, "SetForOAuthClientCredentials", "SetForOAuthAuthorization", "SetForBasicAuthentication", "SetForGenericString"),
		g.KeywordOptions(),
	).
	WithValidation(g.AtLeastOneValueSet, "SetForFlow", "Comment")

// TODO [SNOW-1678749]: Change to use UNSET when it will be possible
var secretUnset = g.NewQueryStruct("SecretUnset").
	PredefinedQueryStructField("Comment", "*bool", g.KeywordOptions().SQL("SET COMMENT = NULL"))

var secretsDef = g.NewInterface(
	"Secrets",
	"Secret",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CustomOperation(
	"CreateWithOAuthClientCredentialsFlow",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithOAuthClientCredentialsFlow").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		SQLWithCustomFieldName("secretType", fmt.Sprintf("TYPE = %s", sdkcommons.SecretTypeOAuth2)).
		Identifier("ApiIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required().Equals().SQL("API_AUTHENTICATION")).
		OptionalQueryStructField("OauthScopes", oauthScopesListDef, g.ParameterOptions().SQL("OAUTH_SCOPES").Parentheses()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	secretsApiIntegrationScopeDef,
).CustomOperation(
	"CreateWithOAuthAuthorizationCodeFlow",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithOAuthAuthorizationCodeFlow").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		SQLWithCustomFieldName("secretType", fmt.Sprintf("TYPE = %s", sdkcommons.SecretTypeOAuth2)).
		TextAssignment("OAUTH_REFRESH_TOKEN", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		TextAssignment("OAUTH_REFRESH_TOKEN_EXPIRY_TIME", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		Identifier("ApiIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required().Equals().SQL("API_AUTHENTICATION")).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).CustomOperation(
	"CreateWithBasicAuthentication",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithBasicAuthentication").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		SQLWithCustomFieldName("secretType", fmt.Sprintf("TYPE = %s", sdkcommons.SecretTypePassword)).
		TextAssignment("USERNAME", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		TextAssignment("PASSWORD", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).CustomOperation(
	"CreateWithGenericString",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithGenericString").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		SQLWithCustomFieldName("secretType", fmt.Sprintf("TYPE = %s", sdkcommons.SecretTypeGenericString)).
		TextAssignment("SECRET_STRING", g.ParameterOptions().SingleQuotes().Required()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-secret",
	g.NewQueryStruct("AlterSecret").
		Alter().
		SQL("SECRET").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			secretSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			secretUnset,
			g.KeywordOptions(),
		).
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-secret",
	g.NewQueryStruct("DropSecret").
		Drop().
		SQL("SECRET").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-secrets",
	secretPairs,
	g.NewQueryStruct("ShowSecret").
		Show().
		SQL("SECRETS").
		OptionalLike().
		OptionalExtendedIn(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDExtendedInFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-secret",
	secretDetailsPairs,
	g.NewQueryStruct("DescribeSecret").
		Describe().
		SQL("SECRET").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
