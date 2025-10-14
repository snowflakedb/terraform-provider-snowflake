package sdk

import (
	"fmt"
	"slices"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

type AuthenticationMethodsOption string

const (
	AuthenticationMethodsAll                     AuthenticationMethodsOption = "ALL"
	AuthenticationMethodsSaml                    AuthenticationMethodsOption = "SAML"
	AuthenticationMethodsPassword                AuthenticationMethodsOption = "PASSWORD"
	AuthenticationMethodsOauth                   AuthenticationMethodsOption = "OAUTH"
	AuthenticationMethodsKeyPair                 AuthenticationMethodsOption = "KEYPAIR"
	AuthenticationMethodsProgrammaticAccessToken AuthenticationMethodsOption = "PROGRAMMATIC_ACCESS_TOKEN"
	AuthenticationMethodsWorkloadIdentity        AuthenticationMethodsOption = "WORKLOAD_IDENTITY"
)

var AllAuthenticationMethods = []AuthenticationMethodsOption{
	AuthenticationMethodsAll,
	AuthenticationMethodsSaml,
	AuthenticationMethodsPassword,
	AuthenticationMethodsOauth,
	AuthenticationMethodsKeyPair,
	AuthenticationMethodsProgrammaticAccessToken,
	AuthenticationMethodsWorkloadIdentity,
}

type MfaAuthenticationMethodsOption string

const (
	MfaAuthenticationMethodsAll      MfaAuthenticationMethodsOption = "ALL"
	MfaAuthenticationMethodsSaml     MfaAuthenticationMethodsOption = "SAML"
	MfaAuthenticationMethodsPassword MfaAuthenticationMethodsOption = "PASSWORD"
)

var AllMfaAuthenticationMethods = []MfaAuthenticationMethodsOption{
	MfaAuthenticationMethodsAll,
	MfaAuthenticationMethodsSaml,
	MfaAuthenticationMethodsPassword,
}

type MfaEnrollmentOption string

const (
	MfaEnrollmentRequired             MfaEnrollmentOption = "REQUIRED"
	MfaEnrollmentRequiredPasswordOnly MfaEnrollmentOption = "REQUIRED_PASSWORD_ONLY"
	MfaEnrollmentOptional             MfaEnrollmentOption = "OPTIONAL"
)

var AllMfaEnrollmentOptions = []MfaEnrollmentOption{
	MfaEnrollmentRequired,
	MfaEnrollmentRequiredPasswordOnly,
	MfaEnrollmentOptional,
}

type ClientTypesOption string

const (
	ClientTypesAll          ClientTypesOption = "ALL"
	ClientTypesSnowflakeUi  ClientTypesOption = "SNOWFLAKE_UI"
	ClientTypesDrivers      ClientTypesOption = "DRIVERS"
	ClientTypesSnowSql      ClientTypesOption = "SNOWSQL"
	ClientTypesSnowflakeCli ClientTypesOption = "SNOWFLAKE_CLI"
)

var AllClientTypes = []ClientTypesOption{
	ClientTypesAll,
	ClientTypesSnowflakeUi,
	ClientTypesDrivers,
	ClientTypesSnowSql,
	ClientTypesSnowflakeCli,
}

var (
	AuthenticationMethodsOptionDef    = g.NewQueryStruct("AuthenticationMethods").PredefinedQueryStructField("Method", g.KindOfT[AuthenticationMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	MfaAuthenticationMethodsOptionDef = g.NewQueryStruct("MfaAuthenticationMethods").PredefinedQueryStructField("Method", g.KindOfT[MfaAuthenticationMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	ClientTypesOptionDef              = g.NewQueryStruct("ClientTypes").PredefinedQueryStructField("ClientType", g.KindOfT[ClientTypesOption](), g.KeywordOptions().SingleQuotes().Required())
	SecurityIntegrationsOptionDef     = g.NewQueryStruct("SecurityIntegrationsOption").Identifier("Name", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required())
)

var AuthenticationPoliciesDef = g.NewInterface(
	"AuthenticationPolicies",
	"AuthenticationPolicy",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy",
		g.NewQueryStruct("CreateAuthenticationPolicy").
			Create().
			OrReplace().
			SQL("AUTHENTICATION POLICY").
			IfNotExists().
			Name().
			ListAssignment("AUTHENTICATION_METHODS", "AuthenticationMethods", g.ParameterOptions().Parentheses()).
			ListAssignment("MFA_AUTHENTICATION_METHODS", "MfaAuthenticationMethods", g.ParameterOptions().Parentheses()).
			PredefinedQueryStructField("MfaEnrollment", g.KindOfTPointer[MfaEnrollmentOption](), g.ParameterOptions().SQL("MFA_ENROLLMENT")).
			ListAssignment("CLIENT_TYPES", "ClientTypes", g.ParameterOptions().Parentheses()).
			ListAssignment("SECURITY_INTEGRATIONS", "SecurityIntegrationsOption", g.ParameterOptions().Parentheses()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		AuthenticationMethodsOptionDef,
		MfaAuthenticationMethodsOptionDef,
		ClientTypesOptionDef,
		SecurityIntegrationsOptionDef,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-authentication-policy",
		g.NewQueryStruct("AlterAuthenticationPolicy").
			Alter().
			SQL("AUTHENTICATION POLICY").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("AuthenticationPolicySet").
					ListAssignment("AUTHENTICATION_METHODS", "AuthenticationMethods", g.ParameterOptions().Parentheses()).
					ListAssignment("MFA_AUTHENTICATION_METHODS", "MfaAuthenticationMethods", g.ParameterOptions().Parentheses()).
					PredefinedQueryStructField("MfaEnrollment", g.KindOfTPointer[MfaEnrollmentOption](), g.ParameterOptions().SQL("MFA_ENROLLMENT")).
					ListAssignment("CLIENT_TYPES", "ClientTypes", g.ParameterOptions().Parentheses()).
					ListAssignment("SECURITY_INTEGRATIONS", "SecurityIntegrationsOption", g.ParameterOptions().Parentheses()).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.AtLeastOneValueSet, "AuthenticationMethods", "MfaAuthenticationMethods", "MfaEnrollment", "ClientTypes", "SecurityIntegrations", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("AuthenticationPolicyUnset").
					OptionalSQL("CLIENT_TYPES").
					OptionalSQL("AUTHENTICATION_METHODS").
					OptionalSQL("SECURITY_INTEGRATIONS").
					OptionalSQL("MFA_AUTHENTICATION_METHODS").
					OptionalSQL("MFA_ENROLLMENT").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "ClientTypes", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaAuthenticationMethods", "MfaEnrollment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "RenameTo").
			WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-authentication-policy",
		g.NewQueryStruct("DropAuthenticationPolicy").
			Drop().
			SQL("AUTHENTICATION POLICY").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-authentication-policies",
		g.DbStruct("showAuthenticationPolicyDBRow").
			Time("created_on").
			Text("name").
			Text("comment").
			Text("database_name").
			Text("schema_name").
			Text("kind").
			Text("owner").
			Text("owner_role_type").
			Text("options"),
		g.PlainStruct("AuthenticationPolicy").
			Time("CreatedOn").
			Text("Name").
			Text("Comment").
			Text("DatabaseName").
			Text("SchemaName").
			Text("Kind").
			Text("Owner").
			Text("OwnerRoleType").
			Text("Options"),
		g.NewQueryStruct("ShowAuthenticationPolicies").
			Show().
			SQL("AUTHENTICATION POLICIES").
			OptionalLike().
			OptionalExtendedIn().
			OptionalOn().
			OptionalStartsWith().
			OptionalLimit(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
		g.ShowByIDExtendedInFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-authentication-policy",
		g.DbStruct("describeAuthenticationPolicyDBRow").
			Text("property").
			Text("value").
			Text("default").
			Text("description"),
		g.PlainStruct("AuthenticationPolicyDescription").
			Text("Property").
			Text("Value").
			Text("Default").
			Text("Description"),
		g.NewQueryStruct("DescribeAuthenticationPolicy").
			Describe().
			SQL("AUTHENTICATION POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)

func ToAuthenticationMethodsOption(s string) (AuthenticationMethodsOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllAuthenticationMethods, AuthenticationMethodsOption(s)) {
		return "", fmt.Errorf("invalid authentication method: %s", s)
	}
	return AuthenticationMethodsOption(s), nil
}

func ToMfaAuthenticationMethodsOption(s string) (MfaAuthenticationMethodsOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllMfaAuthenticationMethods, MfaAuthenticationMethodsOption(s)) {
		return "", fmt.Errorf("invalid MFA authentication method: %s", s)
	}
	return MfaAuthenticationMethodsOption(s), nil
}

func ToMfaEnrollmentOption(s string) (MfaEnrollmentOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllMfaEnrollmentOptions, MfaEnrollmentOption(s)) {
		return "", fmt.Errorf("invalid MFA enrollment option: %s", s)
	}
	return MfaEnrollmentOption(s), nil
}

func ToClientTypesOption(s string) (ClientTypesOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllClientTypes, ClientTypesOption(s)) {
		return "", fmt.Errorf("invalid client type: %s", s)
	}
	return ClientTypesOption(s), nil
}
