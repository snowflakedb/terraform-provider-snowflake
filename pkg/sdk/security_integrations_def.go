package sdk

import (
	"fmt"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

const (
	SecurityIntegrationCategory                                     = "SECURITY"
	ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode = "AUTHORIZATION_CODE"
	ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials = "CLIENT_CREDENTIALS" //nolint:gosec
	ApiAuthenticationSecurityIntegrationOauthGrantJwtBearer         = "JWT_BEARER"
)

type ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption string

const (
	ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption = "CLIENT_SECRET_POST"
)

var AllApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption = []ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption{
	ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost,
}

func ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption(s string) (ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost):
		return ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost, nil
	default:
		return "", fmt.Errorf("invalid ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption: %s", s)
	}
}

type ExternalOauthSecurityIntegrationTypeOption string

const (
	ExternalOauthSecurityIntegrationTypeOkta         ExternalOauthSecurityIntegrationTypeOption = "OKTA"
	ExternalOauthSecurityIntegrationTypeAzure        ExternalOauthSecurityIntegrationTypeOption = "AZURE"
	ExternalOauthSecurityIntegrationTypePingFederate ExternalOauthSecurityIntegrationTypeOption = "PING_FEDERATE"
	ExternalOauthSecurityIntegrationTypeCustom       ExternalOauthSecurityIntegrationTypeOption = "CUSTOM"
)

var AllExternalOauthSecurityIntegrationTypes = []ExternalOauthSecurityIntegrationTypeOption{
	ExternalOauthSecurityIntegrationTypeOkta,
	ExternalOauthSecurityIntegrationTypeAzure,
	ExternalOauthSecurityIntegrationTypePingFederate,
	ExternalOauthSecurityIntegrationTypeCustom,
}

func ToExternalOauthSecurityIntegrationTypeOption(s string) (ExternalOauthSecurityIntegrationTypeOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ExternalOauthSecurityIntegrationTypeOkta):
		return ExternalOauthSecurityIntegrationTypeOkta, nil
	case string(ExternalOauthSecurityIntegrationTypeAzure):
		return ExternalOauthSecurityIntegrationTypeAzure, nil
	case string(ExternalOauthSecurityIntegrationTypePingFederate):
		return ExternalOauthSecurityIntegrationTypePingFederate, nil
	case string(ExternalOauthSecurityIntegrationTypeCustom):
		return ExternalOauthSecurityIntegrationTypeCustom, nil
	default:
		return "", fmt.Errorf("invalid ExternalOauthSecurityIntegrationTypeOption: %s", s)
	}
}

type ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption string

const (
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName    ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption = "LOGIN_NAME"
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption = "EMAIL_ADDRESS"
)

var AllExternalOauthSecurityIntegrationSnowflakeUserMappingAttributes = []ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption{
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName,
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress,
}

func ToExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption(s string) (ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName):
		return ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName, nil
	case string(ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress):
		return ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress, nil
	default:
		return "", fmt.Errorf("invalid ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption: %s", s)
	}
}

type ExternalOauthSecurityIntegrationAnyRoleModeOption string

const (
	ExternalOauthSecurityIntegrationAnyRoleModeDisable            ExternalOauthSecurityIntegrationAnyRoleModeOption = "DISABLE"
	ExternalOauthSecurityIntegrationAnyRoleModeEnable             ExternalOauthSecurityIntegrationAnyRoleModeOption = "ENABLE"
	ExternalOauthSecurityIntegrationAnyRoleModeEnableForPrivilege ExternalOauthSecurityIntegrationAnyRoleModeOption = "ENABLE_FOR_PRIVILEGE"
)

var AllExternalOauthSecurityIntegrationAnyRoleModes = []ExternalOauthSecurityIntegrationAnyRoleModeOption{
	ExternalOauthSecurityIntegrationAnyRoleModeDisable,
	ExternalOauthSecurityIntegrationAnyRoleModeEnable,
	ExternalOauthSecurityIntegrationAnyRoleModeEnableForPrivilege,
}

func ToExternalOauthSecurityIntegrationAnyRoleModeOption(s string) (ExternalOauthSecurityIntegrationAnyRoleModeOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ExternalOauthSecurityIntegrationAnyRoleModeDisable):
		return ExternalOauthSecurityIntegrationAnyRoleModeDisable, nil
	case string(ExternalOauthSecurityIntegrationAnyRoleModeEnable):
		return ExternalOauthSecurityIntegrationAnyRoleModeEnable, nil
	case string(ExternalOauthSecurityIntegrationAnyRoleModeEnableForPrivilege):
		return ExternalOauthSecurityIntegrationAnyRoleModeEnableForPrivilege, nil
	default:
		return "", fmt.Errorf("invalid ExternalOauthSecurityIntegrationAnyRoleModeOption: %s", s)
	}
}

type OauthSecurityIntegrationUseSecondaryRolesOption string

const (
	OauthSecurityIntegrationUseSecondaryRolesImplicit OauthSecurityIntegrationUseSecondaryRolesOption = "IMPLICIT"
	OauthSecurityIntegrationUseSecondaryRolesNone     OauthSecurityIntegrationUseSecondaryRolesOption = "NONE"
)

var AllOauthSecurityIntegrationUseSecondaryRoles = []OauthSecurityIntegrationUseSecondaryRolesOption{
	OauthSecurityIntegrationUseSecondaryRolesImplicit,
	OauthSecurityIntegrationUseSecondaryRolesNone,
}

func ToOauthSecurityIntegrationUseSecondaryRolesOption(s string) (OauthSecurityIntegrationUseSecondaryRolesOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(OauthSecurityIntegrationUseSecondaryRolesImplicit):
		return OauthSecurityIntegrationUseSecondaryRolesImplicit, nil
	case string(OauthSecurityIntegrationUseSecondaryRolesNone):
		return OauthSecurityIntegrationUseSecondaryRolesNone, nil
	default:
		return "", fmt.Errorf("invalid OauthSecurityIntegrationUseSecondaryRolesOption: %s", s)
	}
}

type OauthSecurityIntegrationClientTypeOption string

const (
	OauthSecurityIntegrationClientTypePublic       OauthSecurityIntegrationClientTypeOption = "PUBLIC"
	OauthSecurityIntegrationClientTypeConfidential OauthSecurityIntegrationClientTypeOption = "CONFIDENTIAL"
)

var AllOauthSecurityIntegrationClientTypes = []OauthSecurityIntegrationClientTypeOption{
	OauthSecurityIntegrationClientTypePublic,
	OauthSecurityIntegrationClientTypeConfidential,
}

func ToOauthSecurityIntegrationClientTypeOption(s string) (OauthSecurityIntegrationClientTypeOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(OauthSecurityIntegrationClientTypePublic):
		return OauthSecurityIntegrationClientTypePublic, nil
	case string(OauthSecurityIntegrationClientTypeConfidential):
		return OauthSecurityIntegrationClientTypeConfidential, nil
	default:
		return "", fmt.Errorf("invalid OauthSecurityIntegrationClientTypeOption: %s", s)
	}
}

type OauthSecurityIntegrationClientOption string

const (
	OauthSecurityIntegrationClientLooker         OauthSecurityIntegrationClientOption = "LOOKER"
	OauthSecurityIntegrationClientTableauDesktop OauthSecurityIntegrationClientOption = "TABLEAU_DESKTOP"
	OauthSecurityIntegrationClientTableauServer  OauthSecurityIntegrationClientOption = "TABLEAU_SERVER"
)

var AllOauthSecurityIntegrationClients = []OauthSecurityIntegrationClientOption{
	OauthSecurityIntegrationClientLooker,
	OauthSecurityIntegrationClientTableauDesktop,
	OauthSecurityIntegrationClientTableauServer,
}

func ToOauthSecurityIntegrationClientOption(s string) (OauthSecurityIntegrationClientOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(OauthSecurityIntegrationClientLooker):
		return OauthSecurityIntegrationClientLooker, nil
	case string(OauthSecurityIntegrationClientTableauDesktop):
		return OauthSecurityIntegrationClientTableauDesktop, nil
	case string(OauthSecurityIntegrationClientTableauServer):
		return OauthSecurityIntegrationClientTableauServer, nil
	default:
		return "", fmt.Errorf("invalid OauthSecurityIntegrationClientOption: %s", s)
	}
}

type Saml2SecurityIntegrationSaml2ProviderOption string

const (
	Saml2SecurityIntegrationSaml2ProviderOkta   Saml2SecurityIntegrationSaml2ProviderOption = "OKTA"
	Saml2SecurityIntegrationSaml2ProviderAdfs   Saml2SecurityIntegrationSaml2ProviderOption = "ADFS"
	Saml2SecurityIntegrationSaml2ProviderCustom Saml2SecurityIntegrationSaml2ProviderOption = "CUSTOM"
)

var AllSaml2SecurityIntegrationSaml2Providers = []Saml2SecurityIntegrationSaml2ProviderOption{
	Saml2SecurityIntegrationSaml2ProviderOkta,
	Saml2SecurityIntegrationSaml2ProviderAdfs,
	Saml2SecurityIntegrationSaml2ProviderCustom,
}

func ToSaml2SecurityIntegrationSaml2ProviderOption(s string) (Saml2SecurityIntegrationSaml2ProviderOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(Saml2SecurityIntegrationSaml2ProviderOkta):
		return Saml2SecurityIntegrationSaml2ProviderOkta, nil
	case string(Saml2SecurityIntegrationSaml2ProviderAdfs):
		return Saml2SecurityIntegrationSaml2ProviderAdfs, nil
	case string(Saml2SecurityIntegrationSaml2ProviderCustom):
		return Saml2SecurityIntegrationSaml2ProviderCustom, nil
	default:
		return "", fmt.Errorf("invalid Saml2SecurityIntegrationSaml2ProviderOption: %s", s)
	}
}

type Saml2SecurityIntegrationSaml2RequestedNameidFormatOption string

const (
	Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified                Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress               Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatX509SubjectName            Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:1.1:nameid-format:X509SubjectName"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatWindowsDomainQualifiedName Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:1.1:nameid-format:WindowsDomainQualifiedName"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos                   Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:2.0:nameid-format:kerberos"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatPersistent                 Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatTransient                  Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:2.0:nameid-format:transient"
)

var AllSaml2SecurityIntegrationSaml2RequestedNameidFormats = []Saml2SecurityIntegrationSaml2RequestedNameidFormatOption{
	Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatX509SubjectName,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatWindowsDomainQualifiedName,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatPersistent,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatTransient,
}

func ToSaml2SecurityIntegrationSaml2RequestedNameidFormatOption(s string) (Saml2SecurityIntegrationSaml2RequestedNameidFormatOption, error) {
	switch s {
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatX509SubjectName):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatX509SubjectName, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatWindowsDomainQualifiedName):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatWindowsDomainQualifiedName, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatPersistent):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatPersistent, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatTransient):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatTransient, nil
	default:
		return "", fmt.Errorf("invalid Saml2SecurityIntegrationSaml2RequestedNameidFormatOption: %s", s)
	}
}

type ScimSecurityIntegrationScimClientOption string

const (
	ScimSecurityIntegrationScimClientOkta    ScimSecurityIntegrationScimClientOption = "OKTA"
	ScimSecurityIntegrationScimClientAzure   ScimSecurityIntegrationScimClientOption = "AZURE"
	ScimSecurityIntegrationScimClientGeneric ScimSecurityIntegrationScimClientOption = "GENERIC"
)

var AllScimSecurityIntegrationScimClients = []ScimSecurityIntegrationScimClientOption{
	ScimSecurityIntegrationScimClientOkta,
	ScimSecurityIntegrationScimClientAzure,
	ScimSecurityIntegrationScimClientGeneric,
}

func ToScimSecurityIntegrationScimClientOption(s string) (ScimSecurityIntegrationScimClientOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ScimSecurityIntegrationScimClientOkta):
		return ScimSecurityIntegrationScimClientOkta, nil
	case string(ScimSecurityIntegrationScimClientAzure):
		return ScimSecurityIntegrationScimClientAzure, nil
	case string(ScimSecurityIntegrationScimClientGeneric):
		return ScimSecurityIntegrationScimClientGeneric, nil
	default:
		return "", fmt.Errorf("invalid ScimSecurityIntegrationScimClientOption: %s", s)
	}
}

type ScimSecurityIntegrationRunAsRoleOption string

const (
	ScimSecurityIntegrationRunAsRoleOktaProvisioner        ScimSecurityIntegrationRunAsRoleOption = "OKTA_PROVISIONER"
	ScimSecurityIntegrationRunAsRoleAadProvisioner         ScimSecurityIntegrationRunAsRoleOption = "AAD_PROVISIONER"
	ScimSecurityIntegrationRunAsRoleGenericScimProvisioner ScimSecurityIntegrationRunAsRoleOption = "GENERIC_SCIM_PROVISIONER"
)

var AllScimSecurityIntegrationRunAsRoles = []ScimSecurityIntegrationRunAsRoleOption{
	ScimSecurityIntegrationRunAsRoleOktaProvisioner,
	ScimSecurityIntegrationRunAsRoleAadProvisioner,
	ScimSecurityIntegrationRunAsRoleGenericScimProvisioner,
}

func ToScimSecurityIntegrationRunAsRoleOption(s string) (ScimSecurityIntegrationRunAsRoleOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ScimSecurityIntegrationRunAsRoleOktaProvisioner):
		return ScimSecurityIntegrationRunAsRoleOktaProvisioner, nil
	case string(ScimSecurityIntegrationRunAsRoleAadProvisioner):
		return ScimSecurityIntegrationRunAsRoleAadProvisioner, nil
	case string(ScimSecurityIntegrationRunAsRoleGenericScimProvisioner):
		return ScimSecurityIntegrationRunAsRoleGenericScimProvisioner, nil
	default:
		return "", fmt.Errorf("invalid ScimSecurityIntegrationRunAsRoleOption: %s", s)
	}
}

var (
	allowedScopeDef           = g.NewQueryStruct("AllowedScope").Text("Scope", g.KeywordOptions().SingleQuotes().Required())
	userDomainDef             = g.NewQueryStruct("UserDomain").Text("Domain", g.KeywordOptions().SingleQuotes().Required())
	emailPatternDef           = g.NewQueryStruct("EmailPattern").Text("Pattern", g.KeywordOptions().SingleQuotes().Required())
	preAuthorizedRolesListDef = g.NewQueryStruct("PreAuthorizedRolesList").
					List("PreAuthorizedRolesList", "AccountObjectIdentifier", g.ListOptions().MustParentheses())
	blockedRolesListDef = g.NewQueryStruct("BlockedRolesList").
				List("BlockedRolesList", "AccountObjectIdentifier", g.ListOptions().Required().MustParentheses())
	allowedRolesListDef = g.NewQueryStruct("AllowedRolesList").
				List("AllowedRolesList", "AccountObjectIdentifier", g.ListOptions().Required().MustParentheses())
	jwsKeysUrlDef       = g.NewQueryStruct("JwsKeysUrl").Text("JwsKeyUrl", g.KeywordOptions().SingleQuotes().Required())
	audienceListItemDef = g.NewQueryStruct("AudienceListItem").Text("Item", g.KeywordOptions().SingleQuotes().Required())
	audienceListDef     = g.NewQueryStruct("AudienceList").
				List("AudienceList", "AudienceListItem", g.ListOptions().Required().MustParentheses())
	tokenUserMappingClaimDef = g.NewQueryStruct("TokenUserMappingClaim").Text("Claim", g.KeywordOptions().SingleQuotes().Required())
)

func createSecurityIntegrationOperation(structName string, opts func(qs *g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(structName).
		Create().
		OrReplace().
		SQL("SECURITY INTEGRATION").
		IfNotExists().
		Name()
	qs = opts(qs)
	return qs.
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists")
}

func alterSecurityIntegrationOperation(structName string, opts func(qs *g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(structName).
		Alter().
		SQL("SECURITY INTEGRATION").
		IfExists().
		Name().
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name")
	qs = opts(qs)
	return qs
}

var apiAuthClientCredentialsFlowIntegrationSetDef = g.NewQueryStruct("ApiAuthenticationWithClientCredentialsFlowIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalTextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes()).
	OptionalAssignment(
		"OAUTH_CLIENT_AUTH_METHOD",
		g.KindOfT[ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption](),
		g.ParameterOptions(),
	).
	OptionalTextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes()).
	OptionalSQL("OAUTH_GRANT = CLIENT_CREDENTIALS").
	OptionalNumberAssignment("OAUTH_ACCESS_TOKEN_VALIDITY", g.ParameterOptions()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
	ListAssignment("OAUTH_ALLOWED_SCOPES", "AllowedScope", g.ParameterOptions().Parentheses()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthTokenEndpoint", "OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantClientCredentials",
		"OauthAccessTokenValidity", "OauthRefreshTokenValidity", "OauthAllowedScopes", "Comment")

var apiAuthClientCredentialsFlowIntegrationUnsetDef = g.NewQueryStruct("ApiAuthenticationWithClientCredentialsFlowIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "Comment")

var apiAuthCodeGrantFlowIntegrationSetDef = g.NewQueryStruct("ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalTextAssignment("OAUTH_AUTHORIZATION_ENDPOINT", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes()).
	OptionalAssignment(
		"OAUTH_CLIENT_AUTH_METHOD",
		g.KindOfT[ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption](),
		g.ParameterOptions(),
	).
	OptionalTextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes()).
	OptionalSQL("OAUTH_GRANT = AUTHORIZATION_CODE").
	OptionalNumberAssignment("OAUTH_ACCESS_TOKEN_VALIDITY", g.ParameterOptions()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
	ListAssignment("OAUTH_ALLOWED_SCOPES", "AllowedScope", g.ParameterOptions().Parentheses()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthAuthorizationEndpoint", "OauthTokenEndpoint", "OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantAuthorizationCode",
		"OauthAccessTokenValidity", "OauthRefreshTokenValidity", "OauthAllowedScopes", "Comment")

var apiAuthCodeGrantFlowIntegrationUnsetDef = g.NewQueryStruct("ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "Comment")

var apiAuthJwtBearerFlowIntegrationSetDef = g.NewQueryStruct("ApiAuthenticationWithJwtBearerFlowIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalTextAssignment("OAUTH_AUTHORIZATION_ENDPOINT", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes()).
	OptionalAssignment(
		"OAUTH_CLIENT_AUTH_METHOD",
		g.KindOfT[ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption](),
		g.ParameterOptions(),
	).
	OptionalTextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes()).
	OptionalSQL("OAUTH_GRANT = JWT_BEARER").
	OptionalNumberAssignment("OAUTH_ACCESS_TOKEN_VALIDITY", g.ParameterOptions()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthAuthorizationEndpoint", "OauthTokenEndpoint", "OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantJwtBearer",
		"OauthAccessTokenValidity", "OauthRefreshTokenValidity", "Comment")

var apiAuthJwtBearerFlowIntegrationUnsetDef = g.NewQueryStruct("ApiAuthenticationWithJwtBearerFlowIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "Comment")

var externalOauthIntegrationSetDef = g.NewQueryStruct("ExternalOauthIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalAssignment(
		"EXTERNAL_OAUTH_TYPE",
		g.KindOfT[ExternalOauthSecurityIntegrationTypeOption](),
		g.ParameterOptions(),
	).
	OptionalTextAssignment("EXTERNAL_OAUTH_ISSUER", g.ParameterOptions().SingleQuotes()).
	ListAssignment("EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM", "TokenUserMappingClaim", g.ParameterOptions().Parentheses()).
	OptionalAssignment(
		"EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE",
		g.KindOfT[ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption](),
		g.ParameterOptions().SingleQuotes(),
	).
	ListAssignment("EXTERNAL_OAUTH_JWS_KEYS_URL", "JwsKeysUrl", g.ParameterOptions().Parentheses()).
	OptionalQueryStructField("ExternalOauthBlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("EXTERNAL_OAUTH_BLOCKED_ROLES_LIST").Parentheses()).
	OptionalQueryStructField("ExternalOauthAllowedRolesList", allowedRolesListDef, g.ParameterOptions().SQL("EXTERNAL_OAUTH_ALLOWED_ROLES_LIST").Parentheses()).
	OptionalTextAssignment("EXTERNAL_OAUTH_RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2", g.ParameterOptions().SingleQuotes()).
	OptionalQueryStructField("ExternalOauthAudienceList", audienceListDef, g.ParameterOptions().SQL("EXTERNAL_OAUTH_AUDIENCE_LIST").Parentheses()).
	OptionalAssignment(
		"EXTERNAL_OAUTH_ANY_ROLE_MODE",
		g.KindOfT[ExternalOauthSecurityIntegrationAnyRoleModeOption](),
		g.ParameterOptions(),
	).
	OptionalTextAssignment("EXTERNAL_OAUTH_SCOPE_DELIMITER", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE", g.ParameterOptions().SingleQuotes()).
	// TODO(SNOW-1461780): use COMMENT in unset and here use OptionalComment
	OptionalAssignment("COMMENT", "StringAllowEmpty", g.ParameterOptions()).
	WithValidation(g.ConflictingFields, "ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList").
	WithValidation(g.ConflictingFields, "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey").
	WithValidation(g.ConflictingFields, "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey2").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "ExternalOauthType", "ExternalOauthIssuer", "ExternalOauthTokenUserMappingClaim", "ExternalOauthSnowflakeUserMappingAttribute",
		"ExternalOauthJwsKeysUrl", "ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList", "ExternalOauthRsaPublicKey", "ExternalOauthRsaPublicKey2",
		"ExternalOauthAudienceList", "ExternalOauthAnyRoleMode", "ExternalOauthScopeDelimiter", "ExternalOauthScopeMappingAttribute", "Comment")

var externalOauthIntegrationUnsetDef = g.NewQueryStruct("ExternalOauthIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("EXTERNAL_OAUTH_AUDIENCE_LIST").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "ExternalOauthAudienceList")

var oauthForPartnerApplicationsIntegrationSetDef = g.NewQueryStruct("OauthForPartnerApplicationsIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalBooleanAssignment("OAUTH_ISSUE_REFRESH_TOKENS", g.ParameterOptions()).
	OptionalTextAssignment("OAUTH_REDIRECT_URI", g.ParameterOptions().SingleQuotes()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
	OptionalAssignment(
		"OAUTH_USE_SECONDARY_ROLES",
		g.KindOfT[OauthSecurityIntegrationUseSecondaryRolesOption](),
		g.ParameterOptions(),
	).
	OptionalQueryStructField("BlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("BLOCKED_ROLES_LIST").Parentheses()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthIssueRefreshTokens", "OauthRedirectUri", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles",
		"BlockedRolesList", "Comment")

var oauthForPartnerApplicationsIntegrationUnsetDef = g.NewQueryStruct("OauthForPartnerApplicationsIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("OAUTH_USE_SECONDARY_ROLES").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthUseSecondaryRoles")

var oauthForCustomClientsIntegrationSetDef = g.NewQueryStruct("OauthForCustomClientsIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalTextAssignment("OAUTH_REDIRECT_URI", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("OAUTH_ALLOW_NON_TLS_REDIRECT_URI", g.ParameterOptions()).
	OptionalBooleanAssignment("OAUTH_ENFORCE_PKCE", g.ParameterOptions()).
	OptionalQueryStructField("PreAuthorizedRolesList", preAuthorizedRolesListDef, g.ParameterOptions().SQL("PRE_AUTHORIZED_ROLES_LIST").Parentheses()).
	OptionalQueryStructField("BlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("BLOCKED_ROLES_LIST").Parentheses()).
	OptionalBooleanAssignment("OAUTH_ISSUE_REFRESH_TOKENS", g.ParameterOptions()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
	OptionalAssignment(
		"OAUTH_USE_SECONDARY_ROLES",
		g.KindOfT[OauthSecurityIntegrationUseSecondaryRolesOption](),
		g.ParameterOptions(),
	).
	OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
	OptionalTextAssignment("OAUTH_CLIENT_RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("OAUTH_CLIENT_RSA_PUBLIC_KEY_2", g.ParameterOptions().SingleQuotes()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthRedirectUri", "OauthAllowNonTlsRedirectUri", "OauthEnforcePkce", "PreAuthorizedRolesList",
		"BlockedRolesList", "OauthIssueRefreshTokens", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles", "NetworkPolicy", "OauthClientRsaPublicKey",
		"OauthClientRsaPublicKey2", "Comment")

var oauthForCustomClientsIntegrationUnsetDef = g.NewQueryStruct("OauthForCustomClientsIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("NETWORK_POLICY").
	OptionalSQL("OAUTH_CLIENT_RSA_PUBLIC_KEY").
	OptionalSQL("OAUTH_CLIENT_RSA_PUBLIC_KEY_2").
	OptionalSQL("OAUTH_USE_SECONDARY_ROLES").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "NetworkPolicy", "OauthUseSecondaryRoles", "OauthClientRsaPublicKey", "OauthClientRsaPublicKey2")

var saml2IntegrationSetDef = g.NewQueryStruct("Saml2IntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalTextAssignment("SAML2_ISSUER", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SAML2_SSO_URL", g.ParameterOptions().SingleQuotes()).
	OptionalAssignment(
		"SAML2_PROVIDER",
		g.KindOfT[Saml2SecurityIntegrationSaml2ProviderOption](),
		g.ParameterOptions().SingleQuotes(),
	).
	OptionalTextAssignment("SAML2_X509_CERT", g.ParameterOptions().SingleQuotes()).
	ListAssignment("ALLOWED_USER_DOMAINS", "UserDomain", g.ParameterOptions().Parentheses()).
	ListAssignment("ALLOWED_EMAIL_PATTERNS", "EmailPattern", g.ParameterOptions().Parentheses()).
	OptionalTextAssignment("SAML2_SP_INITIATED_LOGIN_PAGE_LABEL", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("SAML2_ENABLE_SP_INITIATED", g.ParameterOptions()).
	OptionalTextAssignment("SAML2_SNOWFLAKE_X509_CERT", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("SAML2_SIGN_REQUEST", g.ParameterOptions()).
	OptionalAssignment(
		"SAML2_REQUESTED_NAMEID_FORMAT",
		g.KindOfT[Saml2SecurityIntegrationSaml2RequestedNameidFormatOption](),
		g.ParameterOptions().SingleQuotes(),
	).
	OptionalTextAssignment("SAML2_POST_LOGOUT_REDIRECT_URL", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("SAML2_FORCE_AUTHN", g.ParameterOptions()).
	OptionalTextAssignment("SAML2_SNOWFLAKE_ISSUER_URL", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SAML2_SNOWFLAKE_ACS_URL", g.ParameterOptions().SingleQuotes()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "Saml2Issuer", "Saml2SsoUrl", "Saml2Provider", "Saml2X509Cert", "AllowedUserDomains", "AllowedEmailPatterns",
		"Saml2SpInitiatedLoginPageLabel", "Saml2EnableSpInitiated", "Saml2SnowflakeX509Cert", "Saml2SignRequest", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl",
		"Saml2ForceAuthn", "Saml2SnowflakeIssuerUrl", "Saml2SnowflakeAcsUrl", "Comment")

var saml2IntegrationUnsetDef = g.NewQueryStruct("Saml2IntegrationUnset").
	OptionalSQL("SAML2_FORCE_AUTHN").
	OptionalSQL("SAML2_REQUESTED_NAMEID_FORMAT").
	OptionalSQL("SAML2_POST_LOGOUT_REDIRECT_URL").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Saml2ForceAuthn", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Comment")

var scimIntegrationSetDef = g.NewQueryStruct("ScimIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
	OptionalBooleanAssignment("SYNC_PASSWORD", g.ParameterOptions()).
	// TODO(SNOW-1461780): use COMMENT in unset and here use OptionalComment
	OptionalAssignment("COMMENT", "StringAllowEmpty", g.ParameterOptions()).
	WithValidation(g.AtLeastOneValueSet, "Enabled", "NetworkPolicy", "SyncPassword", "Comment")

var scimIntegrationUnsetDef = g.NewQueryStruct("ScimIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("NETWORK_POLICY").
	OptionalSQL("SYNC_PASSWORD").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "NetworkPolicy", "SyncPassword")

var SecurityIntegrationsDef = g.NewInterface(
	"SecurityIntegrations",
	"SecurityIntegration",
	g.KindOfT[AccountObjectIdentifier](),
).
	CustomOperation(
		"CreateApiAuthenticationWithClientCredentialsFlow",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth",
		createSecurityIntegrationOperation("CreateApiAuthenticationWithClientCredentialsFlow", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = API_AUTHENTICATION")).
				PredefinedQueryStructField("authType", "string", g.StaticOptions().SQL("AUTH_TYPE = OAUTH2")).
				BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
				OptionalTextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes()).
				OptionalAssignment(
					"OAUTH_CLIENT_AUTH_METHOD",
					g.KindOfT[ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption](),
					g.ParameterOptions(),
				).
				TextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().Required().SingleQuotes()).
				OptionalSQL("OAUTH_GRANT = CLIENT_CREDENTIALS").
				OptionalNumberAssignment("OAUTH_ACCESS_TOKEN_VALIDITY", g.ParameterOptions()).
				OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
				ListAssignment("OAUTH_ALLOWED_SCOPES", "AllowedScope", g.ParameterOptions().Parentheses())
		}),
		allowedScopeDef,
	).
	CustomOperation(
		"CreateApiAuthenticationWithAuthorizationCodeGrantFlow",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth",
		createSecurityIntegrationOperation("CreateApiAuthenticationWithAuthorizationCodeGrantFlow", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = API_AUTHENTICATION")).
				PredefinedQueryStructField("authType", "string", g.StaticOptions().SQL("AUTH_TYPE = OAUTH2")).
				BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
				OptionalTextAssignment("OAUTH_AUTHORIZATION_ENDPOINT", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes()).
				OptionalAssignment(
					"OAUTH_CLIENT_AUTH_METHOD",
					g.KindOfT[ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption](),
					g.ParameterOptions(),
				).
				TextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().Required().SingleQuotes()).
				OptionalSQL("OAUTH_GRANT = AUTHORIZATION_CODE").
				OptionalNumberAssignment("OAUTH_ACCESS_TOKEN_VALIDITY", g.ParameterOptions()).
				OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
				ListAssignment("OAUTH_ALLOWED_SCOPES", "AllowedScope", g.ParameterOptions().Parentheses())
		}),
	).
	CustomOperation(
		"CreateApiAuthenticationWithJwtBearerFlow",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth",
		createSecurityIntegrationOperation("CreateApiAuthenticationWithJwtBearerFlow", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = API_AUTHENTICATION")).
				PredefinedQueryStructField("authType", "string", g.StaticOptions().SQL("AUTH_TYPE = OAUTH2")).
				BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
				TextAssignment("OAUTH_ASSERTION_ISSUER", g.ParameterOptions().Required().SingleQuotes()).
				OptionalTextAssignment("OAUTH_AUTHORIZATION_ENDPOINT", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes()).
				OptionalAssignment(
					"OAUTH_CLIENT_AUTH_METHOD",
					g.KindOfT[ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption](),
					g.ParameterOptions(),
				).
				TextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().Required().SingleQuotes()).
				OptionalSQL("OAUTH_GRANT = JWT_BEARER").
				OptionalNumberAssignment("OAUTH_ACCESS_TOKEN_VALIDITY", g.ParameterOptions()).
				OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions())
		}),
	).
	CustomOperation(
		"CreateExternalOauth",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-external",
		createSecurityIntegrationOperation("CreateExternalOauth", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = EXTERNAL_OAUTH")).
				BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
				Assignment(
					"EXTERNAL_OAUTH_TYPE",
					g.KindOfT[ExternalOauthSecurityIntegrationTypeOption](),
					g.ParameterOptions().Required(),
				).
				TextAssignment("EXTERNAL_OAUTH_ISSUER", g.ParameterOptions().Required().SingleQuotes()).
				ListAssignment("EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM", "TokenUserMappingClaim", g.ParameterOptions().Required().Parentheses()).
				Assignment(
					"EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE",
					g.KindOfT[ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption](),
					g.ParameterOptions().SingleQuotes().Required(),
				).
				ListAssignment("EXTERNAL_OAUTH_JWS_KEYS_URL", "JwsKeysUrl", g.ParameterOptions().Parentheses()).
				OptionalQueryStructField("ExternalOauthBlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("EXTERNAL_OAUTH_BLOCKED_ROLES_LIST").Parentheses()).
				OptionalQueryStructField("ExternalOauthAllowedRolesList", allowedRolesListDef, g.ParameterOptions().SQL("EXTERNAL_OAUTH_ALLOWED_ROLES_LIST").Parentheses()).
				OptionalTextAssignment("EXTERNAL_OAUTH_RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2", g.ParameterOptions().SingleQuotes()).
				OptionalQueryStructField("ExternalOauthAudienceList", audienceListDef, g.ParameterOptions().SQL("EXTERNAL_OAUTH_AUDIENCE_LIST").Parentheses()).
				OptionalAssignment(
					"EXTERNAL_OAUTH_ANY_ROLE_MODE",
					g.KindOfT[ExternalOauthSecurityIntegrationAnyRoleModeOption](),
					g.ParameterOptions(),
				).
				OptionalTextAssignment("EXTERNAL_OAUTH_SCOPE_DELIMITER", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ConflictingFields, "ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList").
				WithValidation(g.ExactlyOneValueSet, "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey").
				WithValidation(g.ConflictingFields, "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey2")
		}),
		allowedRolesListDef,
		blockedRolesListDef,
		jwsKeysUrlDef,
		audienceListDef,
		audienceListItemDef,
		tokenUserMappingClaimDef,
	).
	CustomOperation(
		"CreateOauthForPartnerApplications",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake",
		createSecurityIntegrationOperation("CreateOauthForPartnerApplications", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = OAUTH")).
				Assignment(
					"OAUTH_CLIENT",
					g.KindOfT[OauthSecurityIntegrationClientOption](),
					g.ParameterOptions().Required(),
				).
				OptionalTextAssignment("OAUTH_REDIRECT_URI", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
				OptionalBooleanAssignment("OAUTH_ISSUE_REFRESH_TOKENS", g.ParameterOptions()).
				OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
				OptionalAssignment(
					"OAUTH_USE_SECONDARY_ROLES",
					g.KindOfT[OauthSecurityIntegrationUseSecondaryRolesOption](),
					g.ParameterOptions(),
				).
				OptionalQueryStructField("BlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("BLOCKED_ROLES_LIST").Parentheses())
		}),
		preAuthorizedRolesListDef,
		blockedRolesListDef,
	).
	CustomOperation(
		"CreateOauthForCustomClients",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake",
		createSecurityIntegrationOperation("CreateOauthForCustomClients", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = OAUTH")).
				PredefinedQueryStructField("oauthClient", "string", g.StaticOptions().SQL("OAUTH_CLIENT = CUSTOM")).
				Assignment(
					"OAUTH_CLIENT_TYPE",
					g.KindOfT[OauthSecurityIntegrationClientTypeOption](),
					g.ParameterOptions().Required().SingleQuotes(),
				).
				OptionalTextAssignment("OAUTH_REDIRECT_URI", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
				OptionalBooleanAssignment("OAUTH_ALLOW_NON_TLS_REDIRECT_URI", g.ParameterOptions()).
				OptionalBooleanAssignment("OAUTH_ENFORCE_PKCE", g.ParameterOptions()).
				OptionalAssignment(
					"OAUTH_USE_SECONDARY_ROLES",
					g.KindOfT[OauthSecurityIntegrationUseSecondaryRolesOption](),
					g.ParameterOptions(),
				).
				OptionalQueryStructField("PreAuthorizedRolesList", preAuthorizedRolesListDef, g.ParameterOptions().SQL("PRE_AUTHORIZED_ROLES_LIST").Parentheses()).
				OptionalQueryStructField("BlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("BLOCKED_ROLES_LIST").Parentheses()).
				OptionalBooleanAssignment("OAUTH_ISSUE_REFRESH_TOKENS", g.ParameterOptions()).
				OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
				OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
				OptionalTextAssignment("OAUTH_CLIENT_RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("OAUTH_CLIENT_RSA_PUBLIC_KEY_2", g.ParameterOptions().SingleQuotes())
		}),
	).
	CustomOperation(
		"CreateSaml2",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-saml2",
		createSecurityIntegrationOperation("CreateSaml2", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = SAML2")).
				OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
				TextAssignment("SAML2_ISSUER", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("SAML2_SSO_URL", g.ParameterOptions().Required().SingleQuotes()).
				Assignment(
					"SAML2_PROVIDER",
					g.KindOfT[Saml2SecurityIntegrationSaml2ProviderOption](),
					g.ParameterOptions().Required().SingleQuotes(),
				).
				TextAssignment("SAML2_X509_CERT", g.ParameterOptions().Required().SingleQuotes()).
				ListAssignment("ALLOWED_USER_DOMAINS", "UserDomain", g.ParameterOptions().Parentheses()).
				ListAssignment("ALLOWED_EMAIL_PATTERNS", "EmailPattern", g.ParameterOptions().Parentheses()).
				OptionalTextAssignment("SAML2_SP_INITIATED_LOGIN_PAGE_LABEL", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("SAML2_ENABLE_SP_INITIATED", g.ParameterOptions()).
				OptionalTextAssignment("SAML2_SNOWFLAKE_X509_CERT", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("SAML2_SIGN_REQUEST", g.ParameterOptions()).
				OptionalAssignment(
					"SAML2_REQUESTED_NAMEID_FORMAT",
					g.KindOfT[Saml2SecurityIntegrationSaml2RequestedNameidFormatOption](),
					g.ParameterOptions().SingleQuotes(),
				).
				OptionalTextAssignment("SAML2_POST_LOGOUT_REDIRECT_URL", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("SAML2_FORCE_AUTHN", g.ParameterOptions()).
				OptionalTextAssignment("SAML2_SNOWFLAKE_ISSUER_URL", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("SAML2_SNOWFLAKE_ACS_URL", g.ParameterOptions().SingleQuotes())
		}),
		userDomainDef,
		emailPatternDef,
	).
	CustomOperation(
		"CreateScim",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-scim",
		createSecurityIntegrationOperation("CreateScim", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = SCIM")).
				OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
				Assignment(
					"SCIM_CLIENT",
					g.KindOfT[ScimSecurityIntegrationScimClientOption](),
					g.ParameterOptions().SingleQuotes().Required(),
				).
				Assignment(
					"RUN_AS_ROLE",
					g.KindOfT[ScimSecurityIntegrationRunAsRoleOption](),
					g.ParameterOptions().SingleQuotes().Required(),
				).
				OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
				OptionalBooleanAssignment("SYNC_PASSWORD", g.ParameterOptions())
		}),
	).
	CustomOperation(
		"AlterApiAuthenticationWithClientCredentialsFlow",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-api-auth",
		alterSecurityIntegrationOperation("AlterApiAuthenticationWithClientCredentialsFlow", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				apiAuthClientCredentialsFlowIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				apiAuthClientCredentialsFlowIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterApiAuthenticationWithAuthorizationCodeGrantFlow",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-api-auth",
		alterSecurityIntegrationOperation("AlterApiAuthenticationWithAuthorizationCodeGrantFlow", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				apiAuthCodeGrantFlowIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				apiAuthCodeGrantFlowIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterApiAuthenticationWithJwtBearerFlow",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-api-auth",
		alterSecurityIntegrationOperation("AlterApiAuthenticationWithJwtBearerFlow", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				apiAuthJwtBearerFlowIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				apiAuthJwtBearerFlowIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterExternalOauth",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-external",
		alterSecurityIntegrationOperation("AlterExternalOauth", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				externalOauthIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				externalOauthIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterOauthForPartnerApplications",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake",
		alterSecurityIntegrationOperation("AlterOauthForPartnerApplications", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				oauthForPartnerApplicationsIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				oauthForPartnerApplicationsIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterOauthForCustomClients",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake",
		alterSecurityIntegrationOperation("AlterOauthForCustomClients", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				oauthForCustomClientsIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				oauthForCustomClientsIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterSaml2",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-saml2",
		alterSecurityIntegrationOperation("AlterSaml2", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				saml2IntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				saml2IntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).OptionalSQL("REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY").
				WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "RefreshSaml2SnowflakePrivateKey", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterScim",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-scim",
		alterSecurityIntegrationOperation("AlterScim", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				scimIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				scimIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropSecurityIntegration").
			Drop().
			SQL("SECURITY INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.DbStruct("securityIntegrationDescRow").
			Field("property", "string").
			Field("property_type", "string").
			Field("property_value", "string").
			Field("property_default", "string"),
		g.PlainStruct("SecurityIntegrationProperty").
			Field("Name", "string").
			Field("Type", "string").
			Field("Value", "string").
			Field("Default", "string"),
		g.NewQueryStruct("DescSecurityIntegration").
			Describe().
			SQL("SECURITY INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.DbStruct("securityIntegrationShowRow").
			Text("name").
			Text("type").
			Text("category").
			Bool("enabled").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("SecurityIntegration").
			Text("Name").
			Text("IntegrationType").
			Text("Category").
			Bool("Enabled").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowSecurityIntegrations").
			Show().
			SQL("SECURITY INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
	)
