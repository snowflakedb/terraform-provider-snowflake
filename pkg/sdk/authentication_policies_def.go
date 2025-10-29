package sdk

import (
	"fmt"
	"slices"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type AuthenticationMethodsOption string

const (
	AuthenticationMethodsAll                     AuthenticationMethodsOption = "ALL"
	AuthenticationMethodsSaml                    AuthenticationMethodsOption = "SAML"
	AuthenticationMethodsPassword                AuthenticationMethodsOption = "PASSWORD"
	AuthenticationMethodsOauth                   AuthenticationMethodsOption = "OAUTH"
	AuthenticationMethodsKeyPair                 AuthenticationMethodsOption = "KEYPAIR"
	AuthenticationMethodsProgrammaticAccessToken AuthenticationMethodsOption = "PROGRAMMATIC_ACCESS_TOKEN" //nolint:gosec
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

type MfaPolicyAllowedMethodsOption string

const (
	MfaPolicyAllowedMethodAll         MfaPolicyAllowedMethodsOption = "ALL"
	MfaPolicyPassAllowedMethodPassKey MfaPolicyAllowedMethodsOption = "PASSKEY"
	MfaPolicyAllowedMethodTotp        MfaPolicyAllowedMethodsOption = "TOTP"
	MfaPolicyAllowedMethodDuo         MfaPolicyAllowedMethodsOption = "DUO"
)

var AllMfaPolicyOptions = []MfaPolicyAllowedMethodsOption{
	MfaPolicyAllowedMethodAll,
	MfaPolicyPassAllowedMethodPassKey,
	MfaPolicyAllowedMethodTotp,
	MfaPolicyAllowedMethodDuo,
}

type NetworkPolicyEvaluationOption string

const (
	NetworkPolicyEvaluationEnforcedRequired    NetworkPolicyEvaluationOption = "ENFORCED_REQUIRED"
	NetworkPolicyEvaluationEnforcedNotRequired NetworkPolicyEvaluationOption = "ENFORCED_NOT_REQUIRED"
	NetworkPolicyEvaluationNotEnforced         NetworkPolicyEvaluationOption = "NOT_ENFORCED"
)

var AllNetworkPolicyEvaluationOptions = []NetworkPolicyEvaluationOption{
	NetworkPolicyEvaluationEnforcedRequired,
	NetworkPolicyEvaluationEnforcedNotRequired,
	NetworkPolicyEvaluationNotEnforced,
}

type AllowedProviderOption string

const (
	AllowedProviderAll   AllowedProviderOption = "ALL"
	AllowedProviderAws   AllowedProviderOption = "AWS"
	AllowedProviderAzure AllowedProviderOption = "AZURE"
	AllowedProviderGcp   AllowedProviderOption = "GCP"
	AllowedProviderOidc  AllowedProviderOption = "OIDC"
)

var AllAllowedProviderOptions = []AllowedProviderOption{
	AllowedProviderAll,
	AllowedProviderAws,
	AllowedProviderAzure,
	AllowedProviderGcp,
	AllowedProviderOidc,
}

type EnforceMfaOnExternalAuthenticationOption string

const (
	EnforceMfaOnExternalAuthenticationAll  EnforceMfaOnExternalAuthenticationOption = "ALL"
	EnforceMfaOnExternalAuthenticationNone EnforceMfaOnExternalAuthenticationOption = "NONE"
)

var AllEnforceMfaOnExternalAuthenticationOptions = []EnforceMfaOnExternalAuthenticationOption{
	EnforceMfaOnExternalAuthenticationAll,
	EnforceMfaOnExternalAuthenticationNone,
}

var (
	AuthenticationMethodsOptionDef    = g.NewQueryStruct("AuthenticationMethods").PredefinedQueryStructField("Method", g.KindOfT[AuthenticationMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	MfaAuthenticationMethodsOptionDef = g.NewQueryStruct("MfaAuthenticationMethods").PredefinedQueryStructField("Method", g.KindOfT[MfaAuthenticationMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	ClientTypesOptionDef              = g.NewQueryStruct("ClientTypes").PredefinedQueryStructField("ClientType", g.KindOfT[ClientTypesOption](), g.KeywordOptions().SingleQuotes().Required())
	SecurityIntegrationsOptionDef     = g.NewQueryStruct("SecurityIntegrationsOption").
						OptionalSQLWithCustomFieldName("All", "('ALL')").
						UnnamedList("SecurityIntegrations", g.KindOfT[AccountObjectIdentifier](), g.KeywordOptions().Parentheses()).
						WithValidation(g.ExactlyOneValueSet, "All", "SecurityIntegrations")
	AuthenticationPolicyMfaPolicyDef = g.NewQueryStruct("AuthenticationPolicyMfaPolicy").
						PredefinedQueryStructField("EnforceMfaOnExternalAuthentication", g.KindOfTPointer[EnforceMfaOnExternalAuthenticationOption](), g.ParameterOptions().SQL("ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION")).
						ListAssignment("ALLOWED_METHODS", "AuthenticationPolicyMfaPolicyListItem", g.ParameterOptions().Parentheses()).
						WithValidation(g.AtLeastOneValueSet, "EnforceMfaOnExternalAuthentication", "AllowedMethods")
	AuthenticationPolicyMfaPolicyListItemDef = g.NewQueryStruct("AuthenticationPolicyMfaPolicyListItem").PredefinedQueryStructField("Method", g.KindOfT[MfaPolicyAllowedMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	AuthenticationPolicyPatPolicyDef         = g.NewQueryStruct("AuthenticationPolicyPatPolicy").
							OptionalNumberAssignment("DEFAULT_EXPIRY_IN_DAYS", g.ParameterOptions().NoQuotes()).
							OptionalNumberAssignment("MAX_EXPIRY_IN_DAYS", g.ParameterOptions().NoQuotes()).
							OptionalAssignment(
			"NETWORK_POLICY_EVALUATION",
			g.KindOfT[NetworkPolicyEvaluationOption](),
			g.ParameterOptions().NoQuotes(),
		).
		WithValidation(g.AtLeastOneValueSet, "DefaultExpiryInDays", "MaxExpiryInDays", "NetworkPolicyEvaluation")
	AuthenticationPolicyAllowedProviderListItemDef = g.NewQueryStruct("AuthenticationPolicyAllowedProviderListItem").PredefinedQueryStructField("Provider", g.KindOfT[AllowedProviderOption](), g.KeywordOptions().SingleQuotes().Required())
	AuthenticationPolicyWorkloadIdentityPolicyDef  = g.NewQueryStruct("AuthenticationPolicyWorkloadIdentityPolicy").
							ListAssignment("ALLOWED_PROVIDERS", "AuthenticationPolicyAllowedProviderListItem", g.ParameterOptions().Parentheses()).
							ListAssignment("ALLOWED_AWS_ACCOUNTS", "StringListItemWrapper", g.ParameterOptions().Parentheses()).
							ListAssignment("ALLOWED_AZURE_ISSUERS", "StringListItemWrapper", g.ParameterOptions().Parentheses()).
							ListAssignment("ALLOWED_OIDC_ISSUERS", "StringListItemWrapper", g.ParameterOptions().Parentheses()).
							WithValidation(g.AtLeastOneValueSet, "AllowedProviders", "AllowedAwsAccounts", "AllowedAzureIssuers", "AllowedOidcIssuers")
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
			OptionalQueryStructField("MfaPolicy", AuthenticationPolicyMfaPolicyDef, g.ListOptions().SQL("MFA_POLICY =").Parentheses().NoComma()).
			ListAssignment("CLIENT_TYPES", "ClientTypes", g.ParameterOptions().Parentheses()).
			OptionalQueryStructField("SecurityIntegrations", SecurityIntegrationsOptionDef, g.ParameterOptions().SQL("SECURITY_INTEGRATIONS")).
			OptionalQueryStructField("PatPolicy", AuthenticationPolicyPatPolicyDef, g.ListOptions().SQL("PAT_POLICY =").Parentheses().NoComma()).
			OptionalQueryStructField("WorkloadIdentityPolicy", AuthenticationPolicyWorkloadIdentityPolicyDef, g.ListOptions().SQL("WORKLOAD_IDENTITY_POLICY =").Parentheses().NoComma()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		AuthenticationMethodsOptionDef,
		MfaAuthenticationMethodsOptionDef,
		ClientTypesOptionDef,
		SecurityIntegrationsOptionDef,
		AuthenticationPolicyMfaPolicyListItemDef,
		AuthenticationPolicyMfaPolicyDef,
		AuthenticationPolicyAllowedProviderListItemDef,
		AuthenticationPolicyWorkloadIdentityPolicyDef,
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
					OptionalQueryStructField("MfaPolicy", AuthenticationPolicyMfaPolicyDef, g.ListOptions().SQL("MFA_POLICY =").Parentheses().NoComma()).
					ListAssignment("CLIENT_TYPES", "ClientTypes", g.ParameterOptions().Parentheses()).
					OptionalQueryStructField("SecurityIntegrations", SecurityIntegrationsOptionDef, g.ParameterOptions().SQL("SECURITY_INTEGRATIONS")).
					OptionalQueryStructField("PatPolicy", AuthenticationPolicyPatPolicyDef, g.ListOptions().SQL("PAT_POLICY =").Parentheses().NoComma()).
					OptionalQueryStructField("WorkloadIdentityPolicy", AuthenticationPolicyWorkloadIdentityPolicyDef, g.ListOptions().SQL("WORKLOAD_IDENTITY_POLICY =").Parentheses().NoComma()).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.AtLeastOneValueSet, "AuthenticationMethods", "MfaAuthenticationMethods", "MfaEnrollment", "ClientTypes", "SecurityIntegrations", "Comment", "MfaPolicy", "PatPolicy", "WorkloadIdentityPolicy"),
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
					OptionalSQL("MFA_POLICY").
					OptionalSQL("PAT_POLICY").
					OptionalSQL("WORKLOAD_IDENTITY_POLICY").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "ClientTypes", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaAuthenticationMethods", "MfaEnrollment", "MfaPolicy", "PatPolicy", "WorkloadIdentityPolicy"),
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

func ToAllowedProviderOption(s string) (AllowedProviderOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllAllowedProviderOptions, AllowedProviderOption(s)) {
		return "", fmt.Errorf("invalid allowed provider: %s", s)
	}
	return AllowedProviderOption(s), nil
}

func ToNetworkPolicyEvaluationOption(s string) (NetworkPolicyEvaluationOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllNetworkPolicyEvaluationOptions, NetworkPolicyEvaluationOption(s)) {
		return "", fmt.Errorf("invalid network policy evaluation option: %s", s)
	}
	return NetworkPolicyEvaluationOption(s), nil
}

func ToEnforceMfaOnExternalAuthenticationOption(s string) (EnforceMfaOnExternalAuthenticationOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllEnforceMfaOnExternalAuthenticationOptions, EnforceMfaOnExternalAuthenticationOption(s)) {
		return "", fmt.Errorf("invalid enforce MFA on external authentication option: %s", s)
	}
	return EnforceMfaOnExternalAuthenticationOption(s), nil
}
