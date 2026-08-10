package sdk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	id := securityIntegrationsTestIdAccountObjectIdentifier
	// Shared across operations: each case only ever assigns one role/policy to a given field,
	// so reusing the same random identifier for "the allowed role", "the blocked role", etc.
	// across Create/Alter variants doesn't create cross-case ambiguity.
	allowedRoleId := randomAccountObjectIdentifier()
	blockedRoleId := randomAccountObjectIdentifier()
	preAuthorizedRoleId := randomAccountObjectIdentifier()
	networkPolicyId := randomAccountObjectIdentifier()

	securityIntegrationsTests.CreateApiAuthenticationWithClientCredentialsFlow.
		withDefaultOpts(func() *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions {
			return &CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions{
				name:              id,
				Enabled:           true,
				OauthClientId:     "foo",
				OauthClientSecret: "bar",
			}
		}).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateApiAuthenticationWithClientCredentialsFlow_basic,
			func(opts *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions) {
				opts.OrReplace = new(true)
			},
			"CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_CLIENT_ID = 'foo' OAUTH_CLIENT_SECRET = 'bar'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateApiAuthenticationWithClientCredentialsFlow_all,
			func(opts *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.OauthTokenEndpoint = new("foo")
				opts.OauthClientAuthMethod = new(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOptionClientSecretPost)
				opts.OauthGrantClientCredentials = new(true)
				opts.OauthAccessTokenValidity = new(42)
				opts.OauthRefreshTokenValidity = new(42)
				opts.OauthAllowedScopes = []AllowedScope{{Scope: "bar"}}
				opts.Comment = new("foo")
			},
			"CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_TOKEN_ENDPOINT = 'foo'"+
				" OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST OAUTH_CLIENT_ID = 'foo' OAUTH_CLIENT_SECRET = 'bar' OAUTH_GRANT = CLIENT_CREDENTIALS"+
				" OAUTH_ACCESS_TOKEN_VALIDITY = 42 OAUTH_REFRESH_TOKEN_VALIDITY = 42 OAUTH_ALLOWED_SCOPES = ('bar') COMMENT = 'foo'",
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.CreateApiAuthenticationWithAuthorizationCodeGrantFlow.
		withDefaultOpts(func() *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions {
			return &CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions{
				name:              id,
				Enabled:           true,
				OauthClientId:     "foo",
				OauthClientSecret: "bar",
			}
		}).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateApiAuthenticationWithAuthorizationCodeGrantFlow_basic,
			func(opts *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions) {
				opts.OrReplace = new(true)
			},
			"CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_CLIENT_ID = 'foo' OAUTH_CLIENT_SECRET = 'bar'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateApiAuthenticationWithAuthorizationCodeGrantFlow_all,
			func(opts *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.OauthAuthorizationEndpoint = new("foo")
				opts.OauthTokenEndpoint = new("foo")
				opts.OauthClientAuthMethod = new(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOptionClientSecretPost)
				opts.OauthGrantAuthorizationCode = new(true)
				opts.OauthAccessTokenValidity = new(42)
				opts.OauthRefreshTokenValidity = new(42)
				opts.OauthAllowedScopes = []AllowedScope{{Scope: "bar"}}
				opts.Comment = new("foo")
			},
			"CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_AUTHORIZATION_ENDPOINT = 'foo'"+
				" OAUTH_TOKEN_ENDPOINT = 'foo' OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST OAUTH_CLIENT_ID = 'foo' OAUTH_CLIENT_SECRET = 'bar' OAUTH_GRANT = AUTHORIZATION_CODE"+
				" OAUTH_ACCESS_TOKEN_VALIDITY = 42 OAUTH_REFRESH_TOKEN_VALIDITY = 42 OAUTH_ALLOWED_SCOPES = ('bar') COMMENT = 'foo'",
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.CreateApiAuthenticationWithJwtBearerFlow.
		withDefaultOpts(func() *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions {
			return &CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions{
				name:                 id,
				Enabled:              true,
				OauthClientId:        "foo",
				OauthClientSecret:    "bar",
				OauthAssertionIssuer: "foo",
			}
		}).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateApiAuthenticationWithJwtBearerFlow_basic,
			func(opts *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions) {
				opts.OrReplace = new(true)
			},
			"CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_ASSERTION_ISSUER = 'foo' OAUTH_CLIENT_ID = 'foo'"+
				" OAUTH_CLIENT_SECRET = 'bar'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateApiAuthenticationWithJwtBearerFlow_all,
			func(opts *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.OauthAuthorizationEndpoint = new("foo")
				opts.OauthTokenEndpoint = new("foo")
				opts.OauthClientAuthMethod = new(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOptionClientSecretPost)
				opts.OauthGrantJwtBearer = new(true)
				opts.OauthAccessTokenValidity = new(42)
				opts.OauthRefreshTokenValidity = new(42)
				opts.Comment = new("foo")
			},
			"CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_ASSERTION_ISSUER = 'foo'"+
				" OAUTH_AUTHORIZATION_ENDPOINT = 'foo' OAUTH_TOKEN_ENDPOINT = 'foo' OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST OAUTH_CLIENT_ID = 'foo' OAUTH_CLIENT_SECRET = 'bar' OAUTH_GRANT = JWT_BEARER"+
				" OAUTH_ACCESS_TOKEN_VALIDITY = 42 OAUTH_REFRESH_TOKEN_VALIDITY = 42 COMMENT = 'foo'",
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.CreateExternalOauth.
		withDefaultOpts(func() *CreateExternalOauthSecurityIntegrationOptions {
			return &CreateExternalOauthSecurityIntegrationOptions{
				name:                               id,
				Enabled:                            false,
				ExternalOauthType:                  ExternalOauthSecurityIntegrationTypeOptionCustom,
				ExternalOauthIssuer:                "foo",
				ExternalOauthTokenUserMappingClaim: []TokenUserMappingClaim{{Claim: "foo"}},
				ExternalOauthSnowflakeUserMappingAttribute: ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress,
			}
		}).
		withModify(
			case_SecurityIntegrations_validation_CreateExternalOauth_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateExternalOauthSecurityIntegrationOptions) {
				opts.ExternalOauthJwsKeysUrl = []JwsKeysUrl{{JwsKeyUrl: "foo"}}
				opts.ExternalOauthRsaPublicKey = new("key")
			},
		).
		withModify(
			case_SecurityIntegrations_validation_CreateExternalOauth_opts_ConflictingFields_ExternalOauthJwsKeysUrl_ExternalOauthRsaPublicKey2,
			func(opts *CreateExternalOauthSecurityIntegrationOptions) {
				opts.ExternalOauthJwsKeysUrl = []JwsKeysUrl{{JwsKeyUrl: "foo"}}
				opts.ExternalOauthRsaPublicKey2 = new("key")
			},
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateExternalOauth_basic,
			func(opts *CreateExternalOauthSecurityIntegrationOptions) {
				opts.OrReplace = new(true)
				opts.ExternalOauthJwsKeysUrl = []JwsKeysUrl{{JwsKeyUrl: "foo"}}
				opts.ExternalOauthBlockedRolesList = &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{blockedRoleId}}
			},
			"CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = EXTERNAL_OAUTH ENABLED = false EXTERNAL_OAUTH_TYPE = CUSTOM EXTERNAL_OAUTH_ISSUER = 'foo'"+
				" EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM = ('foo') EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE = 'EMAIL_ADDRESS' EXTERNAL_OAUTH_JWS_KEYS_URL = ('foo')"+
				" EXTERNAL_OAUTH_BLOCKED_ROLES_LIST = (%s)",
			id.FullyQualifiedName(), blockedRoleId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateExternalOauth_all,
			func(opts *CreateExternalOauthSecurityIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.ExternalOauthAllowedRolesList = &AllowedRolesList{AllowedRolesList: []AccountObjectIdentifier{allowedRoleId}}
				opts.ExternalOauthRsaPublicKey = new("foo")
				opts.ExternalOauthRsaPublicKey2 = new("foo")
				opts.ExternalOauthAudienceList = &AudienceList{AudienceList: []AudienceListItem{{Item: "foo"}}}
				opts.ExternalOauthAnyRoleMode = new(ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable)
				opts.ExternalOauthScopeDelimiter = new(" ")
				opts.ExternalOauthScopeMappingAttribute = new("foo")
				opts.Comment = new("foo")
			},
			"CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = EXTERNAL_OAUTH ENABLED = false EXTERNAL_OAUTH_TYPE = CUSTOM EXTERNAL_OAUTH_ISSUER = 'foo'"+
				" EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM = ('foo') EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE = 'EMAIL_ADDRESS' EXTERNAL_OAUTH_ALLOWED_ROLES_LIST = (%s)"+
				" EXTERNAL_OAUTH_RSA_PUBLIC_KEY = 'foo' EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2 = 'foo' EXTERNAL_OAUTH_AUDIENCE_LIST = ('foo') EXTERNAL_OAUTH_ANY_ROLE_MODE = DISABLE"+
				" EXTERNAL_OAUTH_SCOPE_DELIMITER = ' ' EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE = 'foo' COMMENT = 'foo'",
			id.FullyQualifiedName(), allowedRoleId.FullyQualifiedName(),
		)

	securityIntegrationsTests.CreateOauthForPartnerApplications.
		withDefaultOpts(func() *CreateOauthForPartnerApplicationsSecurityIntegrationOptions {
			return &CreateOauthForPartnerApplicationsSecurityIntegrationOptions{
				name:        id,
				OauthClient: OauthSecurityIntegrationClientOptionTableauDesktop,
			}
		}).
		withAdditionalValidationCase(
			"validation_CreateOauthForPartnerApplications_OauthRedirectUri_requiredForLooker",
			func(opts *CreateOauthForPartnerApplicationsSecurityIntegrationOptions) {
				opts.OauthClient = OauthSecurityIntegrationClientOptionLooker
				opts.OauthRedirectUri = nil
			},
			NewError("OauthRedirectUri is required when OauthClient is LOOKER"),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateOauthForPartnerApplications_basic,
			func(opts *CreateOauthForPartnerApplicationsSecurityIntegrationOptions) {
				opts.OrReplace = new(true)
			},
			"CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = OAUTH OAUTH_CLIENT = TABLEAU_DESKTOP",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateOauthForPartnerApplications_all,
			func(opts *CreateOauthForPartnerApplicationsSecurityIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.OauthClient = OauthSecurityIntegrationClientOptionLooker
				opts.OauthRedirectUri = new("uri")
				opts.Enabled = new(true)
				opts.OauthIssueRefreshTokens = new(true)
				opts.OauthRefreshTokenValidity = new(42)
				opts.OauthUseSecondaryRoles = new(OauthSecurityIntegrationUseSecondaryRolesOptionNone)
				opts.AllowedRolesList = &AllowedRolesList{AllowedRolesList: []AccountObjectIdentifier{allowedRoleId}}
				opts.BlockedRolesList = &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{blockedRoleId}}
				opts.Comment = new("a")
			},
			"CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = OAUTH OAUTH_CLIENT = LOOKER OAUTH_REDIRECT_URI = 'uri' ENABLED = true OAUTH_ISSUE_REFRESH_TOKENS = true"+
				" OAUTH_REFRESH_TOKEN_VALIDITY = 42 OAUTH_USE_SECONDARY_ROLES = NONE ALLOWED_ROLES_LIST = (%s) BLOCKED_ROLES_LIST = (%s) COMMENT = 'a'",
			id.FullyQualifiedName(), allowedRoleId.FullyQualifiedName(), blockedRoleId.FullyQualifiedName(),
		)

	securityIntegrationsTests.CreateOauthForCustomClients.
		withDefaultOpts(func() *CreateOauthForCustomClientsSecurityIntegrationOptions {
			return &CreateOauthForCustomClientsSecurityIntegrationOptions{
				name:             id,
				OauthClientType:  OauthSecurityIntegrationClientTypeOptionPublic,
				OauthRedirectUri: "uri",
			}
		}).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateOauthForCustomClients_basic,
			func(opts *CreateOauthForCustomClientsSecurityIntegrationOptions) {
				opts.OrReplace = new(true)
			},
			"CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = OAUTH OAUTH_CLIENT = CUSTOM OAUTH_CLIENT_TYPE = 'PUBLIC' OAUTH_REDIRECT_URI = 'uri'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateOauthForCustomClients_all,
			func(opts *CreateOauthForCustomClientsSecurityIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.Enabled = new(true)
				opts.OauthAllowNonTlsRedirectUri = new(true)
				opts.OauthEnforcePkce = new(true)
				opts.OauthUseSecondaryRoles = new(OauthSecurityIntegrationUseSecondaryRolesOptionNone)
				opts.PreAuthorizedRolesList = &PreAuthorizedRolesList{PreAuthorizedRolesList: []AccountObjectIdentifier{preAuthorizedRoleId}}
				opts.AllowedRolesList = &AllowedRolesList{AllowedRolesList: []AccountObjectIdentifier{allowedRoleId}}
				opts.BlockedRolesList = &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{blockedRoleId}}
				opts.OauthIssueRefreshTokens = new(true)
				opts.OauthRefreshTokenValidity = new(42)
				opts.NetworkPolicy = &networkPolicyId
				opts.OauthClientRsaPublicKey = new("key")
				opts.OauthClientRsaPublicKey2 = new("key2")
				opts.Comment = new("a")
			},
			"CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = OAUTH OAUTH_CLIENT = CUSTOM OAUTH_CLIENT_TYPE = 'PUBLIC' OAUTH_REDIRECT_URI = 'uri' ENABLED = true"+
				" OAUTH_ALLOW_NON_TLS_REDIRECT_URI = true OAUTH_ENFORCE_PKCE = true OAUTH_USE_SECONDARY_ROLES = NONE PRE_AUTHORIZED_ROLES_LIST = (%s) ALLOWED_ROLES_LIST = (%s) BLOCKED_ROLES_LIST = (%s)"+
				" OAUTH_ISSUE_REFRESH_TOKENS = true OAUTH_REFRESH_TOKEN_VALIDITY = 42 NETWORK_POLICY = '\\\"%s\\\"' OAUTH_CLIENT_RSA_PUBLIC_KEY = 'key' OAUTH_CLIENT_RSA_PUBLIC_KEY_2 = 'key2' COMMENT = 'a'",
			id.FullyQualifiedName(), preAuthorizedRoleId.FullyQualifiedName(), allowedRoleId.FullyQualifiedName(),
			blockedRoleId.FullyQualifiedName(), networkPolicyId.Name(),
		)

	securityIntegrationsTests.CreateSaml2.
		withDefaultOpts(func() *CreateSaml2SecurityIntegrationOptions {
			return &CreateSaml2SecurityIntegrationOptions{
				name:          id,
				Saml2Issuer:   "issuer",
				Saml2SsoUrl:   "url",
				Saml2Provider: "provider",
				Saml2X509Cert: "cert",
			}
		}).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateSaml2_basic,
			func(opts *CreateSaml2SecurityIntegrationOptions) {
				opts.OrReplace = new(true)
			},
			"CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = SAML2 SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateSaml2_all,
			func(opts *CreateSaml2SecurityIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.Enabled = new(true)
				opts.AllowedEmailPatterns = []EmailPattern{{Pattern: "pattern"}}
				opts.AllowedUserDomains = []UserDomain{{Domain: "domain"}}
				opts.Comment = new("a")
				opts.Saml2EnableSpInitiated = new(true)
				opts.Saml2ForceAuthn = new(true)
				opts.Saml2PostLogoutRedirectUrl = new("redirect")
				opts.Saml2RequestedNameidFormat = new(Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos)
				opts.Saml2SignRequest = new(true)
				opts.Saml2SnowflakeAcsUrl = new("acs")
				opts.Saml2SnowflakeIssuerUrl = new("issuer")
				opts.Saml2SpInitiatedLoginPageLabel = new("label")
				opts.Saml2SnowflakeX509Cert = new("cert")
			},
			"CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = SAML2 ENABLED = true SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'"+
				" ALLOWED_USER_DOMAINS = ('domain') ALLOWED_EMAIL_PATTERNS = ('pattern') SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'label' SAML2_ENABLE_SP_INITIATED = true SAML2_SNOWFLAKE_X509_CERT = 'cert' SAML2_SIGN_REQUEST = true"+
				" SAML2_REQUESTED_NAMEID_FORMAT = '%s' SAML2_POST_LOGOUT_REDIRECT_URL = 'redirect' SAML2_FORCE_AUTHN = true SAML2_SNOWFLAKE_ISSUER_URL = 'issuer' SAML2_SNOWFLAKE_ACS_URL = 'acs'"+
				" COMMENT = 'a'",
			id.FullyQualifiedName(), Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos,
		)

	securityIntegrationsTests.CreateScim.
		withDefaultOpts(func() *CreateScimSecurityIntegrationOptions {
			return &CreateScimSecurityIntegrationOptions{
				name:       id,
				ScimClient: "GENERIC",
				RunAsRole:  AccountObjectIdentifier{"GENERIC_SCIM_PROVISIONER"},
			}
		}).
		withAdditionalValidationCase(
			"validation_CreateScim_conflictingSyncPasswordForAzureScimClient",
			func(opts *CreateScimSecurityIntegrationOptions) {
				opts.ScimClient = ScimSecurityIntegrationScimClientOptionAzure
				opts.SyncPassword = new(true)
			},
			NewError("SyncPassword is not supported for Azure scim client"),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateScim_basic,
			func(opts *CreateScimSecurityIntegrationOptions) {
				opts.OrReplace = new(true)
			},
			`CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = SCIM SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = '\"GENERIC_SCIM_PROVISIONER\"'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_CreateScim_all,
			func(opts *CreateScimSecurityIntegrationOptions) {
				opts.Enabled = new(true)
				opts.IfNotExists = new(true)
				opts.NetworkPolicy = &networkPolicyId
				opts.SyncPassword = new(true)
				opts.Comment = new("a")
			},
			`CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = SCIM ENABLED = true SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = '\"GENERIC_SCIM_PROVISIONER\"'`+
				` NETWORK_POLICY = '\"%s\"' SYNC_PASSWORD = true COMMENT = 'a'`,
			id.FullyQualifiedName(), networkPolicyId.Name(),
		)

	securityIntegrationsTests.AlterApiAuthenticationWithClientCredentialsFlow.
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithClientCredentialsFlow_Set,
			func(opts *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions) {
				opts.Set = &ApiAuthenticationWithClientCredentialsFlowIntegrationSet{
					Enabled:                     new(true),
					OauthTokenEndpoint:          new("foo"),
					OauthClientAuthMethod:       new(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOptionClientSecretPost),
					OauthClientId:               new("foo"),
					OauthClientSecret:           new("foo"),
					OauthGrantClientCredentials: new(true),
					OauthAccessTokenValidity:    new(42),
					OauthRefreshTokenValidity:   new(42),
					OauthAllowedScopes:          []AllowedScope{{Scope: "foo"}},
					Comment:                     new("foo"),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_TOKEN_ENDPOINT = 'foo', OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST,"+
				" OAUTH_CLIENT_ID = 'foo', OAUTH_CLIENT_SECRET = 'foo', OAUTH_GRANT = CLIENT_CREDENTIALS, OAUTH_ACCESS_TOKEN_VALIDITY = 42,"+
				" OAUTH_REFRESH_TOKEN_VALIDITY = 42, OAUTH_ALLOWED_SCOPES = ('foo'), COMMENT = 'foo'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithClientCredentialsFlow_Unset,
			func(opts *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions) {
				opts.Unset = &ApiAuthenticationWithClientCredentialsFlowIntegrationUnset{
					Enabled: new(true),
					Comment: new(true),
				}
			},
			"ALTER SECURITY INTEGRATION %s UNSET ENABLED, COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithClientCredentialsFlow_SetTags,
			func(opts *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithClientCredentialsFlow_UnsetTags,
			func(opts *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.AlterApiAuthenticationWithAuthorizationCodeGrantFlow.
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithAuthorizationCodeGrantFlow_Set,
			func(opts *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions) {
				opts.Set = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSet{
					Enabled:                     new(true),
					OauthTokenEndpoint:          new("foo"),
					OauthClientAuthMethod:       new(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOptionClientSecretPost),
					OauthClientId:               new("foo"),
					OauthClientSecret:           new("foo"),
					OauthGrantAuthorizationCode: new(true),
					OauthAccessTokenValidity:    new(42),
					OauthRefreshTokenValidity:   new(42),
					OauthAllowedScopes:          []AllowedScope{{Scope: "bar"}},
					Comment:                     new("foo"),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_TOKEN_ENDPOINT = 'foo', OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST,"+
				" OAUTH_CLIENT_ID = 'foo', OAUTH_CLIENT_SECRET = 'foo', OAUTH_GRANT = AUTHORIZATION_CODE, OAUTH_ACCESS_TOKEN_VALIDITY = 42,"+
				" OAUTH_REFRESH_TOKEN_VALIDITY = 42, OAUTH_ALLOWED_SCOPES = ('bar'), COMMENT = 'foo'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithAuthorizationCodeGrantFlow_Unset,
			func(opts *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions) {
				opts.Unset = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnset{
					Enabled: new(true),
					Comment: new(true),
				}
			},
			"ALTER SECURITY INTEGRATION %s UNSET ENABLED, COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithAuthorizationCodeGrantFlow_SetTags,
			func(opts *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithAuthorizationCodeGrantFlow_UnsetTags,
			func(opts *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.AlterApiAuthenticationWithJwtBearerFlow.
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithJwtBearerFlow_Set,
			func(opts *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions) {
				opts.Set = &ApiAuthenticationWithJwtBearerFlowIntegrationSet{
					Enabled:                   new(true),
					OauthTokenEndpoint:        new("foo"),
					OauthClientAuthMethod:     new(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOptionClientSecretPost),
					OauthClientId:             new("foo"),
					OauthClientSecret:         new("foo"),
					OauthGrantJwtBearer:       new(true),
					OauthAccessTokenValidity:  new(42),
					OauthRefreshTokenValidity: new(42),
					Comment:                   new("foo"),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_TOKEN_ENDPOINT = 'foo', OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST,"+
				" OAUTH_CLIENT_ID = 'foo', OAUTH_CLIENT_SECRET = 'foo', OAUTH_GRANT = JWT_BEARER, OAUTH_ACCESS_TOKEN_VALIDITY = 42,"+
				" OAUTH_REFRESH_TOKEN_VALIDITY = 42, COMMENT = 'foo'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithJwtBearerFlow_Unset,
			func(opts *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions) {
				opts.Unset = &ApiAuthenticationWithJwtBearerFlowIntegrationUnset{
					Enabled: new(true),
					Comment: new(true),
				}
			},
			"ALTER SECURITY INTEGRATION %s UNSET ENABLED, COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithJwtBearerFlow_SetTags,
			func(opts *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterApiAuthenticationWithJwtBearerFlow_UnsetTags,
			func(opts *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.AlterExternalOauth.
		withModify(
			case_SecurityIntegrations_validation_AlterExternalOauth_opts_Set_ConflictingFields_ExternalOauthJwsKeysUrl_ExternalOauthRsaPublicKey,
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.Set = &ExternalOauthIntegrationSet{
					ExternalOauthJwsKeysUrl:   []JwsKeysUrl{{JwsKeyUrl: "foo"}},
					ExternalOauthRsaPublicKey: new("key"),
				}
			},
		).
		withModify(
			case_SecurityIntegrations_validation_AlterExternalOauth_opts_Set_ConflictingFields_ExternalOauthJwsKeysUrl_ExternalOauthRsaPublicKey2,
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.Set = &ExternalOauthIntegrationSet{
					ExternalOauthJwsKeysUrl:    []JwsKeysUrl{{JwsKeyUrl: "foo"}},
					ExternalOauthRsaPublicKey2: new("key"),
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterExternalOauth_Set,
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.Set = &ExternalOauthIntegrationSet{
					ExternalOauthBlockedRolesList: &BlockedRolesList{},
					ExternalOauthAudienceList:     &AudienceList{},
				}
			},
			"ALTER SECURITY INTEGRATION %s SET EXTERNAL_OAUTH_BLOCKED_ROLES_LIST = (), EXTERNAL_OAUTH_AUDIENCE_LIST = ()",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalOauth_Set_emptyAllowedRolesList",
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.Set = &ExternalOauthIntegrationSet{
					ExternalOauthAllowedRolesList: &AllowedRolesList{},
				}
			},
			"ALTER SECURITY INTEGRATION %s SET EXTERNAL_OAUTH_ALLOWED_ROLES_LIST = ()",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalOauth_Set_all",
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.Set = &ExternalOauthIntegrationSet{
					Enabled:                            new(true),
					ExternalOauthType:                  new(ExternalOauthSecurityIntegrationTypeOptionCustom),
					ExternalOauthIssuer:                new("foo"),
					ExternalOauthTokenUserMappingClaim: []TokenUserMappingClaim{{Claim: "foo"}},
					ExternalOauthSnowflakeUserMappingAttribute: new(ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress),
					ExternalOauthAllowedRolesList:              &AllowedRolesList{AllowedRolesList: []AccountObjectIdentifier{allowedRoleId}},
					ExternalOauthRsaPublicKey:                  new("foo"),
					ExternalOauthRsaPublicKey2:                 new("foo"),
					ExternalOauthAudienceList:                  &AudienceList{AudienceList: []AudienceListItem{{Item: "foo"}}},
					ExternalOauthAnyRoleMode:                   new(ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable),
					ExternalOauthScopeDelimiter:                new(" "),
					ExternalOauthScopeMappingAttribute:         new("foo"),
					Comment:                                    new(StringAllowEmpty{Value: "foo"}),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ENABLED = true, EXTERNAL_OAUTH_TYPE = CUSTOM, EXTERNAL_OAUTH_ISSUER = 'foo',"+
				" EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM = ('foo'), EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE = 'EMAIL_ADDRESS', EXTERNAL_OAUTH_ALLOWED_ROLES_LIST = (%s),"+
				" EXTERNAL_OAUTH_RSA_PUBLIC_KEY = 'foo', EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2 = 'foo', EXTERNAL_OAUTH_AUDIENCE_LIST = ('foo'), EXTERNAL_OAUTH_ANY_ROLE_MODE = DISABLE,"+
				" EXTERNAL_OAUTH_SCOPE_DELIMITER = ' ', EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE = 'foo', COMMENT = 'foo'",
			id.FullyQualifiedName(), allowedRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalOauth_Set_blockedRolesListAndJwsKeysUrl",
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.Set = &ExternalOauthIntegrationSet{
					ExternalOauthBlockedRolesList: &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{blockedRoleId}},
					ExternalOauthJwsKeysUrl:       []JwsKeysUrl{{JwsKeyUrl: "foo"}},
				}
			},
			"ALTER SECURITY INTEGRATION %s SET EXTERNAL_OAUTH_JWS_KEYS_URL = ('foo'), EXTERNAL_OAUTH_BLOCKED_ROLES_LIST = (%s)",
			id.FullyQualifiedName(), blockedRoleId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterExternalOauth_Unset,
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.Unset = &ExternalOauthIntegrationUnset{
					Enabled:                   new(true),
					ExternalOauthAudienceList: new(true),
				}
			},
			"ALTER SECURITY INTEGRATION %s UNSET ENABLED, EXTERNAL_OAUTH_AUDIENCE_LIST",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalOauth_Set_emptyComment",
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.Set = &ExternalOauthIntegrationSet{
					Comment: new(StringAllowEmpty{Value: ""}),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET COMMENT = ''",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterExternalOauth_SetTags,
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterExternalOauth_UnsetTags,
			func(opts *AlterExternalOauthSecurityIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.AlterOauthForPartnerApplications.
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterOauthForPartnerApplications_Set,
			func(opts *AlterOauthForPartnerApplicationsSecurityIntegrationOptions) {
				opts.Set = &OauthForPartnerApplicationsIntegrationSet{
					AllowedRolesList: &AllowedRolesList{},
					BlockedRolesList: &BlockedRolesList{},
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ALLOWED_ROLES_LIST = (), BLOCKED_ROLES_LIST = ()",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterOauthForPartnerApplications_Set_all",
			func(opts *AlterOauthForPartnerApplicationsSecurityIntegrationOptions) {
				opts.Set = &OauthForPartnerApplicationsIntegrationSet{
					Enabled:                   new(true),
					OauthRedirectUri:          new("uri"),
					OauthIssueRefreshTokens:   new(true),
					OauthRefreshTokenValidity: new(42),
					OauthUseSecondaryRoles:    new(OauthSecurityIntegrationUseSecondaryRolesOptionNone),
					AllowedRolesList:          &AllowedRolesList{AllowedRolesList: []AccountObjectIdentifier{allowedRoleId}},
					BlockedRolesList:          &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{blockedRoleId}},
					Comment:                   new(StringAllowEmpty{Value: ""}),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_ISSUE_REFRESH_TOKENS = true, OAUTH_REDIRECT_URI = 'uri', OAUTH_REFRESH_TOKEN_VALIDITY = 42,"+
				" OAUTH_USE_SECONDARY_ROLES = NONE, ALLOWED_ROLES_LIST = (%s), BLOCKED_ROLES_LIST = (%s), COMMENT = ''",
			id.FullyQualifiedName(), allowedRoleId.FullyQualifiedName(), blockedRoleId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterOauthForPartnerApplications_Unset,
			func(opts *AlterOauthForPartnerApplicationsSecurityIntegrationOptions) {
				opts.Unset = &OauthForPartnerApplicationsIntegrationUnset{
					Enabled:                new(true),
					OauthUseSecondaryRoles: new(true),
				}
			},
			"ALTER SECURITY INTEGRATION %s UNSET ENABLED, OAUTH_USE_SECONDARY_ROLES",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterOauthForPartnerApplications_SetTags,
			func(opts *AlterOauthForPartnerApplicationsSecurityIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterOauthForPartnerApplications_UnsetTags,
			func(opts *AlterOauthForPartnerApplicationsSecurityIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.AlterOauthForCustomClients.
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterOauthForCustomClients_Set,
			func(opts *AlterOauthForCustomClientsSecurityIntegrationOptions) {
				opts.Set = &OauthForCustomClientsIntegrationSet{
					PreAuthorizedRolesList: &PreAuthorizedRolesList{},
					AllowedRolesList:       &AllowedRolesList{},
					BlockedRolesList:       &BlockedRolesList{},
				}
			},
			"ALTER SECURITY INTEGRATION %s SET PRE_AUTHORIZED_ROLES_LIST = (), ALLOWED_ROLES_LIST = (), BLOCKED_ROLES_LIST = ()",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterOauthForCustomClients_Set_all",
			func(opts *AlterOauthForCustomClientsSecurityIntegrationOptions) {
				opts.Set = &OauthForCustomClientsIntegrationSet{
					Enabled:                     new(true),
					OauthRedirectUri:            new("uri"),
					OauthAllowNonTlsRedirectUri: new(true),
					OauthEnforcePkce:            new(true),
					OauthUseSecondaryRoles:      new(OauthSecurityIntegrationUseSecondaryRolesOptionNone),
					PreAuthorizedRolesList:      &PreAuthorizedRolesList{PreAuthorizedRolesList: []AccountObjectIdentifier{preAuthorizedRoleId}},
					AllowedRolesList:            &AllowedRolesList{AllowedRolesList: []AccountObjectIdentifier{allowedRoleId}},
					BlockedRolesList:            &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{blockedRoleId}},
					OauthIssueRefreshTokens:     new(true),
					OauthRefreshTokenValidity:   new(42),
					NetworkPolicy:               &networkPolicyId,
					OauthClientRsaPublicKey:     new("key"),
					OauthClientRsaPublicKey2:    new("key2"),
					Comment:                     new("a"),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_REDIRECT_URI = 'uri', OAUTH_ALLOW_NON_TLS_REDIRECT_URI = true, OAUTH_ENFORCE_PKCE = true,"+
				" PRE_AUTHORIZED_ROLES_LIST = (%s), ALLOWED_ROLES_LIST = (%s), BLOCKED_ROLES_LIST = (%s), OAUTH_ISSUE_REFRESH_TOKENS = true, OAUTH_REFRESH_TOKEN_VALIDITY = 42, OAUTH_USE_SECONDARY_ROLES = NONE,"+
				" NETWORK_POLICY = '\\\"%s\\\"', OAUTH_CLIENT_RSA_PUBLIC_KEY = 'key', OAUTH_CLIENT_RSA_PUBLIC_KEY_2 = 'key2', COMMENT = 'a'",
			id.FullyQualifiedName(), preAuthorizedRoleId.FullyQualifiedName(), allowedRoleId.FullyQualifiedName(),
			blockedRoleId.FullyQualifiedName(), networkPolicyId.Name(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterOauthForCustomClients_Unset,
			func(opts *AlterOauthForCustomClientsSecurityIntegrationOptions) {
				opts.Unset = &OauthForCustomClientsIntegrationUnset{
					Enabled:                  new(true),
					OauthUseSecondaryRoles:   new(true),
					NetworkPolicy:            new(true),
					OauthClientRsaPublicKey:  new(true),
					OauthClientRsaPublicKey2: new(true),
				}
			},
			"ALTER SECURITY INTEGRATION %s UNSET ENABLED, NETWORK_POLICY, OAUTH_CLIENT_RSA_PUBLIC_KEY, OAUTH_CLIENT_RSA_PUBLIC_KEY_2, OAUTH_USE_SECONDARY_ROLES",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterOauthForCustomClients_SetTags,
			func(opts *AlterOauthForCustomClientsSecurityIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterOauthForCustomClients_UnsetTags,
			func(opts *AlterOauthForCustomClientsSecurityIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.AlterSaml2.
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterSaml2_Set,
			func(opts *AlterSaml2SecurityIntegrationOptions) {
				opts.Set = &Saml2IntegrationSet{
					Enabled:                        new(true),
					Saml2Issuer:                    new("issuer"),
					Saml2SsoUrl:                    new("url"),
					Saml2Provider:                  new(Saml2SecurityIntegrationSaml2ProviderOptionCustom),
					Saml2X509Cert:                  new("cert"),
					AllowedUserDomains:             []UserDomain{{Domain: "domain"}},
					AllowedEmailPatterns:           []EmailPattern{{Pattern: "pattern"}},
					Saml2SpInitiatedLoginPageLabel: new("label"),
					Saml2EnableSpInitiated:         new(true),
					Saml2SnowflakeX509Cert:         new("cert"),
					Saml2SignRequest:               new(true),
					Saml2RequestedNameidFormat:     new(Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos),
					Saml2PostLogoutRedirectUrl:     new("redirect"),
					Saml2ForceAuthn:                new(true),
					Saml2SnowflakeIssuerUrl:        new("issuer"),
					Saml2SnowflakeAcsUrl:           new("acs"),
					Comment:                        new("a"),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ENABLED = true, SAML2_ISSUER = 'issuer', SAML2_SSO_URL = 'url', SAML2_PROVIDER = 'CUSTOM', SAML2_X509_CERT = 'cert',"+
				" ALLOWED_USER_DOMAINS = ('domain'), ALLOWED_EMAIL_PATTERNS = ('pattern'), SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'label', SAML2_ENABLE_SP_INITIATED = true, SAML2_SNOWFLAKE_X509_CERT = 'cert', SAML2_SIGN_REQUEST = true,"+
				" SAML2_REQUESTED_NAMEID_FORMAT = '%s', SAML2_POST_LOGOUT_REDIRECT_URL = 'redirect', SAML2_FORCE_AUTHN = true, SAML2_SNOWFLAKE_ISSUER_URL = 'issuer', SAML2_SNOWFLAKE_ACS_URL = 'acs',"+
				" COMMENT = 'a'",
			id.FullyQualifiedName(), Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos,
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterSaml2_Unset,
			func(opts *AlterSaml2SecurityIntegrationOptions) {
				opts.Unset = &Saml2IntegrationUnset{
					Saml2ForceAuthn:            new(true),
					Saml2RequestedNameidFormat: new(true),
					Saml2PostLogoutRedirectUrl: new(true),
					Comment:                    new(true),
				}
			},
			"ALTER SECURITY INTEGRATION %s UNSET SAML2_FORCE_AUTHN, SAML2_REQUESTED_NAMEID_FORMAT, SAML2_POST_LOGOUT_REDIRECT_URL, COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterSaml2_RefreshSaml2SnowflakePrivateKey,
			func(opts *AlterSaml2SecurityIntegrationOptions) {
				opts.RefreshSaml2SnowflakePrivateKey = new(true)
			},
			"ALTER SECURITY INTEGRATION %s REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterSaml2_SetTags,
			func(opts *AlterSaml2SecurityIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterSaml2_UnsetTags,
			func(opts *AlterSaml2SecurityIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.AlterScim.
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterScim_Set,
			func(opts *AlterScimSecurityIntegrationOptions) {
				opts.Set = &ScimIntegrationSet{
					Enabled:       new(true),
					NetworkPolicy: &networkPolicyId,
					SyncPassword:  new(true),
					Comment:       new(StringAllowEmpty{Value: "test"}),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET ENABLED = true, NETWORK_POLICY = '\\\"%s\\\"', SYNC_PASSWORD = true, COMMENT = 'test'",
			id.FullyQualifiedName(), networkPolicyId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_AlterScim_Set_emptyComment",
			func(opts *AlterScimSecurityIntegrationOptions) {
				opts.Set = &ScimIntegrationSet{
					Comment: new(StringAllowEmpty{Value: ""}),
				}
			},
			"ALTER SECURITY INTEGRATION %s SET COMMENT = ''",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterScim_Unset,
			func(opts *AlterScimSecurityIntegrationOptions) {
				opts.Unset = &ScimIntegrationUnset{
					Enabled:       new(true),
					NetworkPolicy: new(true),
					SyncPassword:  new(true),
				}
			},
			"ALTER SECURITY INTEGRATION %s UNSET ENABLED, NETWORK_POLICY, SYNC_PASSWORD",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterScim_SetTags,
			func(opts *AlterScimSecurityIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_AlterScim_UnsetTags,
			func(opts *AlterScimSecurityIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.Drop.
		withExpectedSqlf(
			case_SecurityIntegrations_sql_Drop_basic,
			"DROP SECURITY INTEGRATION %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_Drop_all,
			func(opts *DropSecurityIntegrationOptions) { opts.IfExists = new(true) },
			"DROP SECURITY INTEGRATION IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.Describe.
		withExpectedSqlf(
			case_SecurityIntegrations_sql_Describe_basic,
			"DESCRIBE SECURITY INTEGRATION %s",
			id.FullyQualifiedName(),
		)

	securityIntegrationsTests.Show.
		withExpectedSql(case_SecurityIntegrations_sql_Show_basic, "SHOW SECURITY INTEGRATIONS").
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_Show_all,
			func(opts *ShowSecurityIntegrationOptions) {
				opts.Like = &Like{Pattern: new("some pattern")}
			},
			"SHOW SECURITY INTEGRATIONS LIKE 'some pattern'",
		).
		withModifyAndExpectedSqlf(
			case_SecurityIntegrations_sql_Show_Like,
			func(opts *ShowSecurityIntegrationOptions) {
				opts.Like = &Like{Pattern: new("some pattern")}
			},
			"SHOW SECURITY INTEGRATIONS LIKE 'some pattern'",
		)
}

func TestSecurityIntegration_SubType(t *testing.T) {
	testCases := map[string]struct {
		integration SecurityIntegration
		subType     string
		err         error
	}{
		"subtype for scim integration": {
			integration: SecurityIntegration{IntegrationType: "SCIM - AZURE"},
			subType:     "AZURE",
		},
		"invalid integration type": {
			integration: SecurityIntegration{IntegrationType: "invalid"},
			err:         errors.New("expected \"<type> - <subtype>\", got: invalid"),
		},
		"empty integration type": {
			integration: SecurityIntegration{IntegrationType: ""},
			err:         errors.New("expected \"<type> - <subtype>\", got: "),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			subType, err := tc.integration.SubType()
			if err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, tc.err)
				require.Equal(t, tc.subType, subType)
			}
		})
	}
}
