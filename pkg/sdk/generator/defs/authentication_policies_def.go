package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	AuthenticationMethodsOptionEnumDef = g.NewEnum(
		"AuthenticationMethodsOption", "AuthenticationMethodsOptions",
		"ALL", "SAML", "PASSWORD", "OAUTH", "KEYPAIR", "PROGRAMMATIC_ACCESS_TOKEN", "WORKLOAD_IDENTITY",
	)
	MfaAuthenticationMethodsOptionEnumDef = g.NewEnum(
		"MfaAuthenticationMethodsOption", "MfaAuthenticationMethodsOptions",
		"ALL", "SAML", "PASSWORD",
	)
	MfaAuthenticationMethodsReadOptionEnumDef = g.NewEnum(
		"MfaAuthenticationMethodsReadOption", "MfaAuthenticationMethodsReadOptions",
		"ALL", "SAML", "PASSWORD", "OIDC",
	)
	MfaEnrollmentOptionEnumDef = g.NewEnum(
		"MfaEnrollmentOption", "MfaEnrollmentOptions",
		"REQUIRED", "REQUIRED_PASSWORD_ONLY", "OPTIONAL",
	)
	MfaEnrollmentReadOptionEnumDef = g.NewEnum(
		"MfaEnrollmentReadOption", "MfaEnrollmentReadOptions",
		"REQUIRED", "REQUIRED_PASSWORD_ONLY", "OPTIONAL", "REQUIRED_SNOWFLAKE_UI_PASSWORD_ONLY",
	)
	ClientTypesOptionEnumDef = g.NewEnum(
		"ClientTypesOption", "ClientTypesOptions",
		"ALL", "SNOWFLAKE_UI", "DRIVERS", "SNOWSQL", "SNOWFLAKE_CLI",
	)
	ClientPolicyDriverTypeEnumDef = g.NewEnum(
		"ClientPolicyDriverType", "ClientPolicyDriverTypes",
		"JDBC_DRIVER", "ODBC_DRIVER", "PYTHON_DRIVER", "JAVASCRIPT_DRIVER",
		"C_DRIVER", "GO_DRIVER", "PHP_DRIVER", "DOTNET_DRIVER",
		"SQL_API", "SNOWPIPE_STREAMING_CLIENT_SDK", "PY_CORE",
		"SPROC_PYTHON", "PYTHON_SNOWPARK", "SQL_ALCHEMY", "SNOWPARK", "SNOWFLAKE_CLIENT",
	)
	MfaPolicyAllowedMethodsOptionEnumDef = g.NewEnum(
		"MfaPolicyAllowedMethodsOption", "MfaPolicyAllowedMethodsOptions",
		"ALL", "PASSKEY", "TOTP", "OTP", "DUO",
	)
	NetworkPolicyEvaluationOptionEnumDef = g.NewEnum(
		"NetworkPolicyEvaluationOption", "NetworkPolicyEvaluationOptions",
		"ENFORCED_REQUIRED", "ENFORCED_NOT_REQUIRED", "NOT_ENFORCED",
	)
	AllowedProviderOptionEnumDef = g.NewEnum(
		"AllowedProviderOption", "AllowedProviderOptions",
		"ALL", "AWS", "AZURE", "GCP", "OIDC",
	)
	EnforceMfaOnExternalAuthenticationOptionEnumDef = g.NewEnum(
		"EnforceMfaOnExternalAuthenticationOption", "EnforceMfaOnExternalAuthenticationOptions",
		"ALL", "NONE",
	)

	AuthenticationMethodsOptionDef = g.NewQueryStruct("AuthenticationMethods").Enum("Method", AuthenticationMethodsOptionEnumDef, g.KeywordOptions().SingleQuotes().Required())
	ClientTypesOptionDef           = g.NewQueryStruct("ClientTypes").Enum("ClientType", ClientTypesOptionEnumDef, g.KeywordOptions().SingleQuotes().Required())
	SecurityIntegrationsOptionDef  = g.NewQueryStruct("SecurityIntegrationsOption").
					OptionalSQLWithCustomFieldName("All", "('ALL')").
					UnnamedList("SecurityIntegrations", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.KeywordOptions().Parentheses()).
					WithValidation(g.ExactlyOneValueSet, "All", "SecurityIntegrations")
	AuthenticationPolicyMfaPolicyDef = g.NewQueryStruct("AuthenticationPolicyMfaPolicy").
						OptionalEnum("EnforceMfaOnExternalAuthentication", EnforceMfaOnExternalAuthenticationOptionEnumDef, g.ParameterOptions().SQL("ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION")).
						ListAssignment("ALLOWED_METHODS", "AuthenticationPolicyMfaPolicyListItem", g.ParameterOptions().Parentheses()).
						WithValidation(g.AtLeastOneValueSet, "EnforceMfaOnExternalAuthentication", "AllowedMethods")
	AuthenticationPolicyMfaPolicyListItemDef = g.NewQueryStruct("AuthenticationPolicyMfaPolicyListItem").Enum("Method", MfaPolicyAllowedMethodsOptionEnumDef, g.KeywordOptions().SingleQuotes().Required())
	AuthenticationPolicyPatPolicyDef         = g.NewQueryStruct("AuthenticationPolicyPatPolicy").
							OptionalNumberAssignment("DEFAULT_EXPIRY_IN_DAYS", g.ParameterOptions().NoQuotes()).
							OptionalNumberAssignment("MAX_EXPIRY_IN_DAYS", g.ParameterOptions().NoQuotes()).
							OptionalBooleanAssignment("REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS", g.ParameterOptions()).
							OptionalEnumAssignment("NETWORK_POLICY_EVALUATION", NetworkPolicyEvaluationOptionEnumDef, g.ParameterOptions().NoQuotes()).
							WithValidation(g.AtLeastOneValueSet, "DefaultExpiryInDays", "MaxExpiryInDays", "RequireRoleRestrictionForServiceUsers", "NetworkPolicyEvaluation")
	AuthenticationPolicyAllowedProviderListItemDef = g.NewQueryStruct("AuthenticationPolicyAllowedProviderListItem").Enum("Provider", AllowedProviderOptionEnumDef, g.KeywordOptions().SingleQuotes().Required())
	AuthenticationPolicyClientPolicyEntryDef       = g.NewQueryStruct("AuthenticationPolicyClientPolicyEntry").
							Enum("ClientType", ClientPolicyDriverTypeEnumDef, g.KeywordOptions().NoQuotes().Required()).
							OptionalQueryStructField("Params", AuthenticationPolicyClientPolicyEntryParamsDef, g.ListOptions().SQL("=").Parentheses().Required())
	AuthenticationPolicyClientPolicyEntryParamsDef = g.NewQueryStruct("AuthenticationPolicyClientPolicyEntryParams").
							OptionalTextAssignment("MINIMUM_VERSION", g.ParameterOptions().SingleQuotes())
	AuthenticationPolicyWorkloadIdentityPolicyDef = g.NewQueryStruct("AuthenticationPolicyWorkloadIdentityPolicy").
							ListAssignment("ALLOWED_PROVIDERS", "AuthenticationPolicyAllowedProviderListItem", g.ParameterOptions().Parentheses()).
							ListAssignment("ALLOWED_AWS_ACCOUNTS", "StringListItemWrapper", g.ParameterOptions().Parentheses()).
							ListAssignment("ALLOWED_AZURE_ISSUERS", "StringListItemWrapper", g.ParameterOptions().Parentheses()).
							ListAssignment("ALLOWED_OIDC_ISSUERS", "StringListItemWrapper", g.ParameterOptions().Parentheses()).
							WithValidation(g.AtLeastOneValueSet, "AllowedProviders", "AllowedAwsAccounts", "AllowedAzureIssuers", "AllowedOidcIssuers")
)

var authenticationPoliciesDef = g.NewInterface(
	"AuthenticationPolicies",
	"AuthenticationPolicy",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
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
			OptionalEnum("MfaEnrollment", MfaEnrollmentOptionEnumDef, g.ParameterOptions().SQL("MFA_ENROLLMENT")).
			OptionalQueryStructField("MfaPolicy", AuthenticationPolicyMfaPolicyDef, g.ListOptions().SQL("MFA_POLICY =").Parentheses().NoComma()).
			ListAssignment("CLIENT_TYPES", "ClientTypes", g.ParameterOptions().Parentheses()).
			ListAssignment("CLIENT_POLICY", "AuthenticationPolicyClientPolicyEntry", g.ParameterOptions().Parentheses()).
			OptionalQueryStructField("SecurityIntegrations", SecurityIntegrationsOptionDef, g.ParameterOptions().SQL("SECURITY_INTEGRATIONS")).
			OptionalQueryStructField("PatPolicy", AuthenticationPolicyPatPolicyDef, g.ListOptions().SQL("PAT_POLICY =").Parentheses().NoComma()).
			OptionalQueryStructField("WorkloadIdentityPolicy", AuthenticationPolicyWorkloadIdentityPolicyDef, g.ListOptions().SQL("WORKLOAD_IDENTITY_POLICY =").Parentheses().NoComma()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		AuthenticationMethodsOptionDef,
		ClientTypesOptionDef,
		SecurityIntegrationsOptionDef,
		AuthenticationPolicyMfaPolicyListItemDef,
		AuthenticationPolicyMfaPolicyDef,
		AuthenticationPolicyClientPolicyEntryParamsDef,
		AuthenticationPolicyClientPolicyEntryDef,
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
					OptionalEnum("MfaEnrollment", MfaEnrollmentOptionEnumDef, g.ParameterOptions().SQL("MFA_ENROLLMENT")).
					OptionalQueryStructField("MfaPolicy", AuthenticationPolicyMfaPolicyDef, g.ListOptions().SQL("MFA_POLICY =").Parentheses().NoComma()).
					ListAssignment("CLIENT_TYPES", "ClientTypes", g.ParameterOptions().Parentheses()).
					ListAssignment("CLIENT_POLICY", "AuthenticationPolicyClientPolicyEntry", g.ParameterOptions().Parentheses()).
					OptionalQueryStructField("SecurityIntegrations", SecurityIntegrationsOptionDef, g.ParameterOptions().SQL("SECURITY_INTEGRATIONS")).
					OptionalQueryStructField("PatPolicy", AuthenticationPolicyPatPolicyDef, g.ListOptions().SQL("PAT_POLICY =").Parentheses().NoComma()).
					OptionalQueryStructField("WorkloadIdentityPolicy", AuthenticationPolicyWorkloadIdentityPolicyDef, g.ListOptions().SQL("WORKLOAD_IDENTITY_POLICY =").Parentheses().NoComma()).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.AtLeastOneValueSet, "AuthenticationMethods", "MfaEnrollment", "ClientTypes", "ClientPolicy", "SecurityIntegrations", "Comment", "MfaPolicy", "PatPolicy", "WorkloadIdentityPolicy"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("AuthenticationPolicyUnset").
					OptionalSQL("CLIENT_TYPES").
					OptionalSQL("CLIENT_POLICY").
					OptionalSQL("AUTHENTICATION_METHODS").
					OptionalSQL("SECURITY_INTEGRATIONS").
					OptionalSQL("MFA_ENROLLMENT").
					OptionalSQL("MFA_POLICY").
					OptionalSQL("PAT_POLICY").
					OptionalSQL("WORKLOAD_IDENTITY_POLICY").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "ClientTypes", "ClientPolicy", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaEnrollment", "MfaPolicy", "PatPolicy", "WorkloadIdentityPolicy"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			Identifier("RenameTo", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
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
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-authentication-policies",
		g.StructPair("showAuthenticationPolicyDBRow", "AuthenticationPolicy").
			OptionalTime("created_on", g.WithRequiredInPlain()).
			Text("name").
			Text("comment").
			OptionalText("database_name", g.WithRequiredInPlain(), g.WithManualConvert()).
			OptionalText("schema_name", g.WithRequiredInPlain(), g.WithManualConvert()).
			Text("kind").
			OptionalText("owner", g.WithRequiredInPlain()).
			OptionalText("owner_role_type", g.WithRequiredInPlain()).
			Text("options"),
		g.NewQueryStruct("ShowAuthenticationPolicies").
			Show().
			SQL("AUTHENTICATION POLICIES").
			OptionalLike().
			OptionalExtendedIn().
			OptionalOn().
			OptionalStartsWith().
			OptionalLimit(),
		g.ShowByIDLikeFiltering,
		g.ShowByIDExtendedInFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-authentication-policy",
		g.StructPair("describeAuthenticationPolicyDBRow", "AuthenticationPolicyDescription").
			Text("property").
			Text("value").
			Text("default").
			Text("description"),
		g.NewQueryStruct("DescribeAuthenticationPolicy").
			Describe().
			SQL("AUTHENTICATION POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithEnums(
		AuthenticationMethodsOptionEnumDef,
		MfaAuthenticationMethodsOptionEnumDef,
		MfaAuthenticationMethodsReadOptionEnumDef,
		MfaEnrollmentOptionEnumDef,
		MfaEnrollmentReadOptionEnumDef,
		ClientTypesOptionEnumDef,
		ClientPolicyDriverTypeEnumDef,
		AllowedProviderOptionEnumDef,
		NetworkPolicyEvaluationOptionEnumDef,
		MfaPolicyAllowedMethodsOptionEnumDef,
		EnforceMfaOnExternalAuthenticationOptionEnumDef,
	)
