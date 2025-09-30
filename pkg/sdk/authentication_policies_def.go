package sdk

import (
	"fmt"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

type AuthenticationMethodsOption string

const (
	AuthenticationMethodsAll      AuthenticationMethodsOption = "ALL"
	AuthenticationMethodsSaml     AuthenticationMethodsOption = "SAML"
	AuthenticationMethodsPassword AuthenticationMethodsOption = "PASSWORD"
	AuthenticationMethodsOauth    AuthenticationMethodsOption = "OAUTH"
	AuthenticationMethodsKeyPair  AuthenticationMethodsOption = "KEYPAIR"
)

var AllAuthenticationMethods = []AuthenticationMethodsOption{
	AuthenticationMethodsAll,
	AuthenticationMethodsSaml,
	AuthenticationMethodsPassword,
	AuthenticationMethodsOauth,
	AuthenticationMethodsKeyPair,
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

type ClientTypesOption string

const (
	ClientTypesAll         ClientTypesOption = "ALL"
	ClientTypesSnowflakeUi ClientTypesOption = "SNOWFLAKE_UI"
	ClientTypesDrivers     ClientTypesOption = "DRIVERS"
	ClientTypesSnowSql     ClientTypesOption = "SNOWSQL"
)

var AllClientTypes = []ClientTypesOption{
	ClientTypesAll,
	ClientTypesSnowflakeUi,
	ClientTypesDrivers,
	ClientTypesSnowSql,
}

var (
	AuthenticationMethodsOptionDef    = g.NewQueryStruct("AuthenticationMethods").PredefinedQueryStructField("Method", g.KindOfT[AuthenticationMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	MfaAuthenticationMethodsOptionDef = g.NewQueryStruct("MfaAuthenticationMethods").PredefinedQueryStructField("Method", g.KindOfT[MfaAuthenticationMethods](), g.KeywordOptions().SingleQuotes().Required())
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
			Field("created_on", "string").
			Field("name", "string").
			Field("comment", "string").
			Field("database_name", "string").
			Field("schema_name", "string").
			Field("owner", "string").
			Field("owner_role_type", "string").
			Field("options", "string"),
		g.PlainStruct("AuthenticationPolicy").
			Field("CreatedOn", "string").
			Field("Name", "string").
			Field("Comment", "string").
			Field("DatabaseName", "string").
			Field("SchemaName", "string").
			Field("Owner", "string").
			Field("OwnerRoleType", "string").
			Field("Options", "string"),
		g.NewQueryStruct("ShowAuthenticationPolicies").
			Show().
			SQL("AUTHENTICATION POLICIES").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalLimit(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
		g.ShowByIDInFiltering,
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

func ToAuthenticationMethodsOption(s string) (*AuthenticationMethodsOption, error) {
	switch authenticationMethodsOption := AuthenticationMethodsOption(strings.ToUpper(s)); authenticationMethodsOption {
	case AuthenticationMethodsAll,
		AuthenticationMethodsSaml,
		AuthenticationMethodsPassword,
		AuthenticationMethodsOauth,
		AuthenticationMethodsKeyPair:
		return &authenticationMethodsOption, nil
	default:
		return nil, fmt.Errorf("invalid authentication method type: %s", s)
	}
}

func ToMfaAuthenticationMethodsOption(s string) (*MfaAuthenticationMethodsOption, error) {
	switch mfaAuthenticationMethodsOption := MfaAuthenticationMethodsOption(strings.ToUpper(s)); mfaAuthenticationMethodsOption {
	case MfaAuthenticationMethodsAll,
		MfaAuthenticationMethodsSaml,
		MfaAuthenticationMethodsPassword:
		return &mfaAuthenticationMethodsOption, nil
	default:
		return nil, fmt.Errorf("invalid MFA authentication method type: %s", s)
	}
}

func ToMfaEnrollmentOption(s string) (*MfaEnrollmentOption, error) {
	switch mfaEnrollmentOption := MfaEnrollmentOption(strings.ToUpper(s)); mfaEnrollmentOption {
	case MfaEnrollmentRequired,
		MfaEnrollmentRequiredPasswordOnly,
		MfaEnrollmentOptional:
		return &mfaEnrollmentOption, nil
	default:
		return nil, fmt.Errorf("invalid enrollment option type: %s", s)
	}
}

func ToClientTypesOption(s string) (*ClientTypesOption, error) {
	switch clientTypesOption := ClientTypesOption(strings.ToUpper(s)); clientTypesOption {
	case ClientTypesAll,
		ClientTypesSnowflakeUi,
		ClientTypesDrivers,
		ClientTypesSnowSql:
		return &clientTypesOption, nil
	default:
		return nil, fmt.Errorf("invalid client type: %s", s)
	}
}
