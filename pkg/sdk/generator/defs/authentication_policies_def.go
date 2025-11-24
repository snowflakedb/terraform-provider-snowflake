package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	AuthenticationMethodsOptionDef    = g.NewQueryStruct("AuthenticationMethods").PredefinedQueryStructField("Method", g.KindOfT[sdkcommons.AuthenticationMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	MfaAuthenticationMethodsOptionDef = g.NewQueryStruct("MfaAuthenticationMethods").PredefinedQueryStructField("Method", g.KindOfT[sdkcommons.MfaAuthenticationMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	ClientTypesOptionDef              = g.NewQueryStruct("ClientTypes").PredefinedQueryStructField("ClientType", g.KindOfT[sdkcommons.ClientTypesOption](), g.KeywordOptions().SingleQuotes().Required())
	SecurityIntegrationsOptionDef     = g.NewQueryStruct("SecurityIntegrationsOption").
						OptionalSQLWithCustomFieldName("All", "('ALL')").
						UnnamedList("SecurityIntegrations", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.KeywordOptions().Parentheses()).
						WithValidation(g.ExactlyOneValueSet, "All", "SecurityIntegrations")
	AuthenticationPolicyMfaPolicyDef = g.NewQueryStruct("AuthenticationPolicyMfaPolicy").
						PredefinedQueryStructField("EnforceMfaOnExternalAuthentication", g.KindOfTPointer[sdkcommons.EnforceMfaOnExternalAuthenticationOption](), g.ParameterOptions().SQL("ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION")).
						ListAssignment("ALLOWED_METHODS", "AuthenticationPolicyMfaPolicyListItem", g.ParameterOptions().Parentheses()).
						WithValidation(g.AtLeastOneValueSet, "EnforceMfaOnExternalAuthentication", "AllowedMethods")
	AuthenticationPolicyMfaPolicyListItemDef = g.NewQueryStruct("AuthenticationPolicyMfaPolicyListItem").PredefinedQueryStructField("Method", g.KindOfT[sdkcommons.MfaPolicyAllowedMethodsOption](), g.KeywordOptions().SingleQuotes().Required())
	AuthenticationPolicyPatPolicyDef         = g.NewQueryStruct("AuthenticationPolicyPatPolicy").
							OptionalNumberAssignment("DEFAULT_EXPIRY_IN_DAYS", g.ParameterOptions().NoQuotes()).
							OptionalNumberAssignment("MAX_EXPIRY_IN_DAYS", g.ParameterOptions().NoQuotes()).
							OptionalAssignment(
			"NETWORK_POLICY_EVALUATION",
			g.KindOfT[sdkcommons.NetworkPolicyEvaluationOption](),
			g.ParameterOptions().NoQuotes(),
		).
		WithValidation(g.AtLeastOneValueSet, "DefaultExpiryInDays", "MaxExpiryInDays", "NetworkPolicyEvaluation")
	AuthenticationPolicyAllowedProviderListItemDef = g.NewQueryStruct("AuthenticationPolicyAllowedProviderListItem").PredefinedQueryStructField("Provider", g.KindOfT[sdkcommons.AllowedProviderOption](), g.KeywordOptions().SingleQuotes().Required())
	AuthenticationPolicyWorkloadIdentityPolicyDef  = g.NewQueryStruct("AuthenticationPolicyWorkloadIdentityPolicy").
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
			ListAssignment("MFA_AUTHENTICATION_METHODS", "MfaAuthenticationMethods", g.ParameterOptions().Parentheses()).
			PredefinedQueryStructField("MfaEnrollment", g.KindOfTPointer[sdkcommons.MfaEnrollmentOption](), g.ParameterOptions().SQL("MFA_ENROLLMENT")).
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
					PredefinedQueryStructField("MfaEnrollment", g.KindOfTPointer[sdkcommons.MfaEnrollmentOption](), g.ParameterOptions().SQL("MFA_ENROLLMENT")).
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
