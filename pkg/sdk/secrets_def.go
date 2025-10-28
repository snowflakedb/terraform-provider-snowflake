package sdk

import (
	"fmt"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type (
	SecretType string
)

func ToSecretType(s string) (SecretType, error) {
	switch strings.ToUpper(s) {
	case string(SecretTypePassword):
		return SecretTypePassword, nil
	case string(SecretTypeOAuth2):
		return SecretTypeOAuth2, nil
	case string(SecretTypeGenericString):
		return SecretTypeGenericString, nil
	case string(SecretTypeOAuth2ClientCredentials):
		return SecretTypeOAuth2ClientCredentials, nil
	case string(SecretTypeOAuth2AuthorizationCodeGrant):
		return SecretTypeOAuth2AuthorizationCodeGrant, nil
	default:
		return "", fmt.Errorf("invalid secret type: %s", s)
	}
}

const (
	SecretTypePassword                     SecretType = "PASSWORD"
	SecretTypeOAuth2                       SecretType = "OAUTH2"
	SecretTypeGenericString                SecretType = "GENERIC_STRING"
	SecretTypeOAuth2ClientCredentials      SecretType = "OAUTH2_CLIENT_CREDENTIALS"       // #nosec G101
	SecretTypeOAuth2AuthorizationCodeGrant SecretType = "OAUTH2_AUTHORIZATION_CODE_GRANT" // #nosec G101
)

var AcceptableSecretTypes = map[SecretType][]SecretType{
	SecretTypePassword:                     {SecretTypePassword},
	SecretTypeOAuth2:                       {SecretTypeOAuth2},
	SecretTypeGenericString:                {SecretTypeGenericString},
	SecretTypeOAuth2ClientCredentials:      {SecretTypeOAuth2ClientCredentials, SecretTypeOAuth2},
	SecretTypeOAuth2AuthorizationCodeGrant: {SecretTypeOAuth2ClientCredentials, SecretTypeOAuth2},
}

var secretDbRow = g.DbStruct("secretDBRow").
	Field("created_on", "time.Time").
	Field("name", "string").
	Field("schema_name", "string").
	Field("database_name", "string").
	Field("owner", "string").
	Field("comment", "sql.NullString").
	Field("secret_type", "string").
	Field("oauth_scopes", "sql.NullString").
	Field("owner_role_type", "string")

var secret = g.PlainStruct("Secret").
	Field("CreatedOn", "time.Time").
	Field("Name", "string").
	Field("SchemaName", "string").
	Field("DatabaseName", "string").
	Field("Owner", "string").
	Field("Comment", "*string").
	Field("SecretType", "string").
	Field("OauthScopes", "[]string").
	Field("OwnerRoleType", "string")

var secretDetailsDbRow = g.DbStruct("secretDetailsDBRow").
	Field("created_on", "time.Time").
	Field("name", "string").
	Field("schema_name", "string").
	Field("database_name", "string").
	Field("owner", "string").
	Field("comment", "sql.NullString").
	Field("secret_type", "string").
	Field("username", "sql.NullString").
	Field("oauth_access_token_expiry_time", "*time.Time").
	Field("oauth_refresh_token_expiry_time", "*time.Time").
	Field("oauth_scopes", "sql.NullString").
	Field("integration_name", "sql.NullString")

var secretDetails = g.PlainStruct("SecretDetails").
	Field("CreatedOn", "time.Time").
	Field("Name", "string").
	Field("SchemaName", "string").
	Field("DatabaseName", "string").
	Field("Owner", "string").
	Field("Comment", "*string").
	Field("SecretType", "string").
	Field("Username", "*string").
	Field("OauthAccessTokenExpiryTime", "*time.Time").
	Field("OauthRefreshTokenExpiryTime", "*time.Time").
	Field("OauthScopes", "[]string").
	Field("IntegrationName", "*string")

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

var SecretsDef = g.NewInterface(
	"Secrets",
	"Secret",
	g.KindOfT[SchemaObjectIdentifier](),
).CustomOperation(
	"CreateWithOAuthClientCredentialsFlow",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithOAuthClientCredentialsFlow").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		PredefinedQueryStructField("secretType", "string", g.StaticOptions().SQL(fmt.Sprintf("TYPE = %s", SecretTypeOAuth2))).
		Identifier("ApiIntegration", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required().Equals().SQL("API_AUTHENTICATION")).
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
		PredefinedQueryStructField("secretType", "string", g.StaticOptions().SQL(fmt.Sprintf("TYPE = %s", SecretTypeOAuth2))).
		TextAssignment("OAUTH_REFRESH_TOKEN", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		TextAssignment("OAUTH_REFRESH_TOKEN_EXPIRY_TIME", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		Identifier("ApiIntegration", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required().Equals().SQL("API_AUTHENTICATION")).
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
		PredefinedQueryStructField("secretType", "string", g.StaticOptions().SQL(fmt.Sprintf("TYPE = %s", SecretTypePassword))).
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
		PredefinedQueryStructField("secretType", "string", g.StaticOptions().SQL(fmt.Sprintf("TYPE = %s", SecretTypeGenericString))).
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
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-secrets",
	secretDbRow,
	secret,
	g.NewQueryStruct("ShowSecret").
		Show().
		SQL("SECRETS").
		OptionalLike().
		OptionalExtendedIn(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDExtendedInFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-secret",
	secretDetailsDbRow,
	secretDetails,
	g.NewQueryStruct("DescribeSecret").
		Describe().
		SQL("SECRET").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
