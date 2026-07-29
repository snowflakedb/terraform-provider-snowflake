//go:build sdk_generation

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var wifTypeEnum = g.NewEnum("WIFType", "WIFTypes", "AWS", "AZURE", "GCP", "OIDC")

var secondaryRolesOptionEnum = g.NewEnum("SecondaryRolesOption", "SecondaryRolesOptions", "DEFAULT", "NONE", "ALL")

var userTypeEnum = g.NewEnum("UserType", "UserTypes", "PERSON", "SERVICE", "LEGACY_SERVICE")

var wifMethodsPairs = g.StructPair("userWorkloadIdentityAuthenticationMethodsDBRow", "UserWorkloadIdentityAuthenticationMethod").
	Text("name").
	Field("type", "string", "WIFType", g.WithManualConvert()).
	OptionalText("comment", g.WithRequiredInPlain()).
	OptionalTime("last_used", g.WithRequiredInPlain()).
	Time("created_on").
	// additional_info maps to 4 typed sub-structs via additionalConvert(); added to db row for sqlx scanning.
	Field("additional_info", "sql.NullString", "sql.NullString", g.WithManualConvert()).
	PlainOnlyField("AwsAdditionalInfo", "*UserWorkloadIdentityAuthenticationMethodsAwsAdditionalInfo").
	PlainOnlyField("AzureAdditionalInfo", "*UserWorkloadIdentityAuthenticationMethodsAzureAdditionalInfo").
	PlainOnlyField("GcpAdditionalInfo", "*UserWorkloadIdentityAuthenticationMethodsGcpAdditionalInfo").
	PlainOnlyField("OidcAdditionalInfo", "*UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo")

func showWifMethodsQueryStruct() *g.QueryStruct {
	return g.NewQueryStruct("ShowUserWorkloadIdentityAuthenticationMethodOptions").
		Show().
		SQL("USER WORKLOAD IDENTITY AUTHENTICATION METHODS").
		Identifier("ForUser", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FOR USER").NoEquals().Required()).
		WithValidation(g.ValidIdentifier, "ForUser")
}

var describeUserPropertyPairs = g.StructPair("describeUserPropertyRow", "UserProperty").
	Text("property").
	Text("value").
	Text("default", g.WithPlainFieldName("DefaultValue")).
	Text("description")

func describeUserQueryStruct() *g.QueryStruct {
	return g.NewQueryStruct("DescribeUser").
		Describe().
		SQL("USER").
		Name().
		WithValidation(g.ValidIdentifier, "name")
}

var userPairs = g.StructPair("userDBRow", "User").
	Text("name").
	Time("created_on").
	OptionalText("login_name", g.WithRequiredInPlain()).
	OptionalText("display_name", g.WithRequiredInPlain()).
	OptionalText("first_name", g.WithRequiredInPlain()).
	OptionalText("last_name", g.WithRequiredInPlain()).
	OptionalText("email", g.WithRequiredInPlain()).
	OptionalText("mins_to_unlock", g.WithRequiredInPlain()).
	OptionalText("days_to_expiry", g.WithRequiredInPlain()).
	OptionalText("comment", g.WithRequiredInPlain()).
	Field("disabled", "sql.NullString", "bool", g.WithManualConvert()).
	Field("must_change_password", "sql.NullString", "bool", g.WithManualConvert()).
	Field("snowflake_lock", "sql.NullString", "bool", g.WithManualConvert()).
	OptionalText("default_warehouse", g.WithRequiredInPlain()).
	OptionalText("default_namespace", g.WithRequiredInPlain()).
	OptionalText("default_role", g.WithRequiredInPlain()).
	OptionalText("default_secondary_roles", g.WithRequiredInPlain()).
	Field("ext_authn_duo", "sql.NullString", "bool", g.WithManualConvert()).
	OptionalText("ext_authn_uid", g.WithRequiredInPlain()).
	OptionalText("mins_to_bypass_mfa", g.WithRequiredInPlain()).
	Text("owner").
	OptionalTime("last_success_login", g.WithRequiredInPlain()).
	OptionalTime("expires_at_time", g.WithRequiredInPlain()).
	OptionalTime("locked_until_time", g.WithRequiredInPlain()).
	OptionalBool("has_password", g.WithRequiredInPlain()).
	OptionalBool("has_rsa_public_key", g.WithRequiredInPlain()).
	OptionalText("type", g.WithRequiredInPlain()).
	OptionalBool("has_mfa", g.WithRequiredInPlain()).
	OptionalBool("has_workload_identity", g.WithRequiredInPlain())

func secondaryRolesStruct() *g.QueryStruct {
	return g.NewQueryStruct("SecondaryRoles").
		OptionalSQLWithCustomFieldName("None", "()").
		OptionalSQLWithCustomFieldName("All", "('ALL')").
		WithValidation(g.ExactlyOneValueSet, "None", "All")
}

func userObjectWorkloadIdentityAwsStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserObjectWorkloadIdentityAws").
		PredefinedQueryStructField("wifType", "string", g.StaticOptions().SQL("TYPE = AWS")).
		OptionalTextAssignment("ARN", g.ParameterOptions().SingleQuotes())
}

func userObjectWorkloadIdentityAzureStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserObjectWorkloadIdentityAzure").
		PredefinedQueryStructField("wifType", "string", g.StaticOptions().SQL("TYPE = AZURE")).
		OptionalTextAssignment("ISSUER", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("SUBJECT", g.ParameterOptions().SingleQuotes())
}

func userObjectWorkloadIdentityGcpStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserObjectWorkloadIdentityGcp").
		PredefinedQueryStructField("wifType", "string", g.StaticOptions().SQL("TYPE = GCP")).
		OptionalTextAssignment("SUBJECT", g.ParameterOptions().SingleQuotes())
}

func userObjectWorkloadIdentityOidcStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserObjectWorkloadIdentityOidc").
		PredefinedQueryStructField("wifType", "string", g.StaticOptions().SQL("TYPE = OIDC")).
		OptionalTextAssignment("ISSUER", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("SUBJECT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("OidcAudienceList", "[]StringListItemWrapper", g.ParameterOptions().Parentheses().SQL("OIDC_AUDIENCE_LIST"))
}

func userObjectWorkloadIdentityPropertiesStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserObjectWorkloadIdentityProperties").
		OptionalQueryStructField("AwsType", userObjectWorkloadIdentityAwsStruct(), g.KeywordOptions()).
		OptionalQueryStructField("AzureType", userObjectWorkloadIdentityAzureStruct(), g.KeywordOptions()).
		OptionalQueryStructField("GcpType", userObjectWorkloadIdentityGcpStruct(), g.KeywordOptions()).
		OptionalQueryStructField("OidcType", userObjectWorkloadIdentityOidcStruct(), g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "AwsType", "AzureType", "GcpType", "OidcType")
}

func userObjectPropertiesFields(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalTextAssignment("PASSWORD", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("LOGIN_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("FIRST_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("MIDDLE_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("LAST_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("EMAIL", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("MUST_CHANGE_PASSWORD", nil).
		OptionalBooleanAssignment("DISABLED", nil).
		OptionalNumberAssignment("DAYS_TO_EXPIRY", g.ParameterOptions()).
		OptionalNumberAssignment("MINS_TO_UNLOCK", g.ParameterOptions()).
		OptionalIdentifier("DefaultWarehouse", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("DEFAULT_WAREHOUSE").Equals()).
		OptionalIdentifier("DefaultNamespace", "ObjectIdentifier", g.IdentifierOptions().SQL("DEFAULT_NAMESPACE").Equals()).
		OptionalIdentifier("DefaultRole", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("DEFAULT_ROLE").Equals()).
		OptionalQueryStructField("DefaultSecondaryRoles", secondaryRolesStruct(), g.ParameterOptions().SQL("DEFAULT_SECONDARY_ROLES")).
		OptionalNumberAssignment("MINS_TO_BYPASS_MFA", g.ParameterOptions()).
		OptionalTextAssignment("RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("RSA_PUBLIC_KEY_FP", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("RSA_PUBLIC_KEY_2", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("RSA_PUBLIC_KEY_2_FP", g.ParameterOptions().SingleQuotes()).
		OptionalAssignmentWithFieldName("TYPE", userTypeEnum.KindPtr(), g.ParameterOptions().NoQuotes(), "UserType").
		OptionalQueryStructField("WorkloadIdentity", userObjectWorkloadIdentityPropertiesStruct(), g.ListOptions().SQL("WORKLOAD_IDENTITY =").Parentheses().NoComma()).
		OptionalComment()
}

func userObjectPropertiesStruct() *g.QueryStruct {
	return userObjectPropertiesFields(g.NewQueryStruct("UserObjectProperties"))
}

func userAlterObjectPropertiesStruct() *g.QueryStruct {
	return userObjectPropertiesFields(g.NewQueryStruct("UserAlterObjectProperties")).
		OptionalBooleanAssignment("DISABLE_MFA", nil)
}

func userObjectParametersStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserObjectParameters").
		OptionalBooleanAssignment("ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR", nil).
		OptionalIdentifier("NetworkPolicy", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("NETWORK_POLICY").Equals()).
		OptionalBooleanAssignment("PREVENT_UNLOAD_TO_INTERNAL_STAGES", nil)
}

func userObjectPropertiesUnsetStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserObjectPropertiesUnset").
		OptionalSQL("PASSWORD").
		OptionalSQL("LOGIN_NAME").
		OptionalSQL("DISPLAY_NAME").
		OptionalSQL("FIRST_NAME").
		OptionalSQL("MIDDLE_NAME").
		OptionalSQL("LAST_NAME").
		OptionalSQL("EMAIL").
		OptionalSQL("MUST_CHANGE_PASSWORD").
		OptionalSQL("DISABLED").
		OptionalSQL("DAYS_TO_EXPIRY").
		OptionalSQL("MINS_TO_UNLOCK").
		OptionalSQL("DEFAULT_WAREHOUSE").
		OptionalSQL("DEFAULT_NAMESPACE").
		OptionalSQL("DEFAULT_ROLE").
		OptionalSQL("DEFAULT_SECONDARY_ROLES").
		OptionalSQL("MINS_TO_BYPASS_MFA").
		OptionalSQL("DISABLE_MFA").
		OptionalSQL("RSA_PUBLIC_KEY").
		OptionalSQL("RSA_PUBLIC_KEY_2").
		OptionalSQLWithCustomFieldName("UserType", "TYPE").
		OptionalSQL("WORKLOAD_IDENTITY").
		OptionalSQL("COMMENT")
}

func userObjectParametersUnsetStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserObjectParametersUnset").
		OptionalSQL("ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR").
		OptionalSQL("NETWORK_POLICY").
		OptionalSQL("PREVENT_UNLOAD_TO_INTERNAL_STAGES")
}

func addDelegatedAuthorizationStruct() *g.QueryStruct {
	return g.NewQueryStruct("AddDelegatedAuthorization").
		PredefinedQueryStructField("Role", "string", g.ParameterOptions().NoEquals().SQL("ADD DELEGATED AUTHORIZATION OF ROLE")).
		PredefinedQueryStructField("Integration", "string", g.ParameterOptions().NoEquals().SQL("TO SECURITY INTEGRATION"))
}

func removeDelegatedAuthorizationStruct() *g.QueryStruct {
	return g.NewQueryStruct("RemoveDelegatedAuthorization").
		OptionalAssignment("REMOVE DELEGATED AUTHORIZATION OF ROLE", "string", g.ParameterOptions().NoEquals()).
		PredefinedQueryStructField("Authorizations", "*bool", g.ParameterOptions().NoEquals().SQL("REMOVE DELEGATED AUTHORIZATIONS")).
		PredefinedQueryStructField("Integration", "string", g.ParameterOptions().NoEquals().SQL("FROM SECURITY INTEGRATION")).
		WithValidation(g.ExactlyOneValueSet, "RemoveDelegatedAuthorizationOfRole", "Authorizations").
		WithValidation(g.ValidateValueSet, "Integration")
}

func userSetStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserSet").
		OptionalIdentifier("PasswordPolicy", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PASSWORD POLICY")).
		OptionalIdentifier("SessionPolicy", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("SESSION POLICY")).
		OptionalIdentifier("AuthenticationPolicy", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("AUTHENTICATION POLICY")).
		OptionalQueryStructField("ObjectProperties", userAlterObjectPropertiesStruct(), g.KeywordOptions()).
		OptionalQueryStructField("ObjectParameters", userObjectParametersStruct(), g.KeywordOptions()).
		PredefinedQueryStructField("SessionParameters", "*SessionParameters", g.KeywordOptions()).
		WithValidation(g.AtLeastOneValueSet, "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters").
		WithValidation(g.MoreThanOneValueSet, "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy").
		WithAdditionalValidations()
}

func userUnsetStruct() *g.QueryStruct {
	return g.NewQueryStruct("UserUnset").
		OptionalSQL("PASSWORD POLICY").
		OptionalSQL("SESSION POLICY").
		OptionalSQL("AUTHENTICATION POLICY").
		OptionalQueryStructField("ObjectProperties", userObjectPropertiesUnsetStruct(), g.ListOptions()).
		OptionalQueryStructField("ObjectParameters", userObjectParametersUnsetStruct(), g.ListOptions()).
		PredefinedQueryStructField("SessionParameters", "*SessionParametersUnset", g.ListOptions()).
		WithValidation(g.AtLeastOneValueSet, "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters").
		WithValidation(g.MoreThanOneValueSet, "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy").
		WithAdditionalValidations()
}

var usersDef = g.NewInterface(
	"Users",
	"User",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-user",
	g.NewQueryStruct("CreateUser").
		Create().
		OrReplace().
		SQL("USER").
		IfNotExists().
		Name().
		OptionalQueryStructField("ObjectProperties", userObjectPropertiesStruct(), g.KeywordOptions()).
		OptionalQueryStructField("ObjectParameters", userObjectParametersStruct(), g.KeywordOptions()).
		PredefinedQueryStructField("SessionParameters", "*SessionParameters", g.KeywordOptions()).
		OptionalSQL("WITH").
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-user",
	g.NewQueryStruct("AlterUser").
		Alter().
		SQL("USER").
		IfExists().
		Name().
		OptionalIdentifier("NewName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalSQL("RESET PASSWORD").
		OptionalSQL("ABORT ALL QUERIES").
		OptionalQueryStructField("AddDelegatedAuthorization", addDelegatedAuthorizationStruct(), g.KeywordOptions()).
		OptionalQueryStructField("RemoveDelegatedAuthorization", removeDelegatedAuthorizationStruct(), g.KeywordOptions()).
		OptionalQueryStructField("Set", userSetStruct(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset", userUnsetStruct(), g.ListOptions().NoParentheses().SQL("UNSET")).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "NewName", "ResetPassword", "AbortAllQueries", "AddDelegatedAuthorization", "RemoveDelegatedAuthorization", "Set", "Unset", "SetTags", "UnsetTags"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-user",
	g.NewQueryStruct("DropUser").
		Drop().
		SQL("USER").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-users",
	userPairs,
	g.NewQueryStruct("ShowUsers").
		Show().
		Terse().
		SQL("USERS").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimit().
		WithAdditionalValidations(),
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-user",
	describeUserPropertyPairs,
	describeUserQueryStruct(),
	g.PlainStruct("UserDetails").
		Field("Name", "*StringProperty").
		Field("Comment", "*StringProperty").
		Field("DisplayName", "*StringProperty").
		Field("Type", "*StringProperty").
		Field("LoginName", "*StringProperty").
		Field("FirstName", "*StringProperty").
		Field("MiddleName", "*StringProperty").
		Field("LastName", "*StringProperty").
		Field("Email", "*StringProperty").
		Field("Password", "*StringProperty").
		Field("MustChangePassword", "*BoolProperty").
		Field("Disabled", "*BoolProperty").
		Field("SnowflakeLock", "*BoolProperty").
		Field("SnowflakeSupport", "*BoolProperty").
		Field("DaysToExpiry", "*FloatProperty").
		Field("MinsToUnlock", "*IntProperty").
		Field("DefaultWarehouse", "*StringProperty").
		Field("DefaultNamespace", "*StringProperty").
		Field("DefaultRole", "*StringProperty").
		Field("DefaultSecondaryRoles", "*StringProperty").
		Field("ExtAuthnDuo", "*BoolProperty").
		Field("ExtAuthnUid", "*StringProperty").
		Field("MinsToBypassMfa", "*IntProperty").
		Field("MinsToBypassNetworkPolicy", "*IntProperty").
		Field("RsaPublicKey", "*StringProperty").
		Field("RsaPublicKeyFp", "*StringProperty").
		Field("RsaPublicKeyLastSetTime", "*StringProperty").
		Field("RsaPublicKey2", "*StringProperty").
		Field("RsaPublicKey2Fp", "*StringProperty").
		Field("RsaPublicKey2LastSetTime", "*StringProperty").
		Field("PasswordLastSetTime", "*StringProperty").
		Field("CustomLandingPageUrl", "*StringProperty").
		Field("CustomLandingPageUrlFlushNextUiLoad", "*BoolProperty").
		Field("HasMfa", "*BoolProperty").
		Field("HasWorkloadIdentity", "*BoolProperty"),
).ShowParameters("AccountObjectIdentifier").
	WithCustomInterfaceMethod(
		"DescribeDetails",
		"DescribeDetails aggregates the []UserProperty result of Describe into *UserDetails. Callers should migrate from Describe to DescribeDetails.",
		[]*g.MethodParameter{g.NewMethodParameter("id", "AccountObjectIdentifier")},
		"*UserDetails", "error",
	).
	WithCustomInterfaceMethod(
		"AddProgrammaticAccessToken", "",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*AddUserProgrammaticAccessTokenRequest")},
		"*AddProgrammaticAccessTokenResult", "error",
	).
	WithCustomInterfaceMethod(
		"ModifyProgrammaticAccessToken", "",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*ModifyUserProgrammaticAccessTokenRequest")},
		"error",
	).
	WithCustomInterfaceMethod(
		"RotateProgrammaticAccessToken", "",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*RotateUserProgrammaticAccessTokenRequest")},
		"*RotateProgrammaticAccessTokenResult", "error",
	).
	WithCustomInterfaceMethod(
		"RemoveProgrammaticAccessToken", "",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*RemoveUserProgrammaticAccessTokenRequest")},
		"error",
	).
	WithCustomInterfaceMethod(
		"RemoveProgrammaticAccessTokenSafely", "",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*RemoveUserProgrammaticAccessTokenRequest")},
		"error",
	).
	WithCustomInterfaceMethod(
		"ShowProgrammaticAccessTokens", "",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*ShowUserProgrammaticAccessTokenRequest")},
		"[]ProgrammaticAccessToken", "error",
	).
	WithCustomInterfaceMethod(
		"ShowProgrammaticAccessTokenByName", "",
		[]*g.MethodParameter{
			g.NewMethodParameter("userId", "AccountObjectIdentifier"),
			g.NewMethodParameter("tokenName", "AccountObjectIdentifier"),
		},
		"*ProgrammaticAccessToken", "error",
	).
	WithCustomInterfaceMethod(
		"ShowProgrammaticAccessTokenByNameSafely", "",
		[]*g.MethodParameter{
			g.NewMethodParameter("userId", "AccountObjectIdentifier"),
			g.NewMethodParameter("tokenName", "AccountObjectIdentifier"),
		},
		"*ProgrammaticAccessToken", "error",
	).
	CustomShowOperationWithPairedStructs(
		"ShowUserWorkloadIdentityAuthenticationMethodOptions",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/show-user-workload-identity-authentication-methods",
		wifMethodsPairs,
		showWifMethodsQueryStruct(),
	).
	WithEnums(
		wifTypeEnum,
		secondaryRolesOptionEnum,
		userTypeEnum,
	)
