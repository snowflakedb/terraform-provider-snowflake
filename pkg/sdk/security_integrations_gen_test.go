package sdk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecurityIntegrations_CreateApiAuthenticationWithClientCredentialsFlow(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions
	defaultOpts := func() *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions {
		return &CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions{
			name:              id,
			Enabled:           true,
			OauthClientId:     "foo",
			OauthClientSecret: "bar",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_CLIENT_ID = 'foo'"+
			" OAUTH_CLIENT_SECRET = 'bar'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OauthTokenEndpoint = Pointer("foo")
		opts.OauthClientAuthMethod = Pointer(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)
		opts.OauthGrantClientCredentials = Pointer(true)
		opts.OauthAccessTokenValidity = Pointer(42)
		opts.OauthRefreshTokenValidity = Pointer(42)
		opts.OauthAllowedScopes = []AllowedScope{{Scope: "bar"}}
		opts.Comment = Pointer("foo")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_TOKEN_ENDPOINT = 'foo'"+
			" OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST OAUTH_CLIENT_ID = 'foo' OAUTH_CLIENT_SECRET = 'bar' OAUTH_GRANT = CLIENT_CREDENTIALS"+
			" OAUTH_ACCESS_TOKEN_VALIDITY = 42 OAUTH_REFRESH_TOKEN_VALIDITY = 42 OAUTH_ALLOWED_SCOPES = ('bar') COMMENT = 'foo'", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateApiAuthenticationWithAuthorizationCodeGrantFlow(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions
	defaultOpts := func() *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions {
		return &CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions{
			name:              id,
			Enabled:           true,
			OauthClientId:     "foo",
			OauthClientSecret: "bar",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})
	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_CLIENT_ID = 'foo'"+
			" OAUTH_CLIENT_SECRET = 'bar'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OauthAuthorizationEndpoint = Pointer("foo")
		opts.OauthTokenEndpoint = Pointer("foo")
		opts.OauthClientAuthMethod = Pointer(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)
		opts.OauthGrantAuthorizationCode = Pointer(true)
		opts.OauthAccessTokenValidity = Pointer(42)
		opts.OauthRefreshTokenValidity = Pointer(42)
		opts.OauthAllowedScopes = []AllowedScope{{Scope: "bar"}}
		opts.Comment = Pointer("foo")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_AUTHORIZATION_ENDPOINT = 'foo'"+
			" OAUTH_TOKEN_ENDPOINT = 'foo' OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST OAUTH_CLIENT_ID = 'foo' OAUTH_CLIENT_SECRET = 'bar' OAUTH_GRANT = AUTHORIZATION_CODE"+
			" OAUTH_ACCESS_TOKEN_VALIDITY = 42 OAUTH_REFRESH_TOKEN_VALIDITY = 42 OAUTH_ALLOWED_SCOPES = ('bar') COMMENT = 'foo'", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateApiAuthenticationWithJwtBearerFlow(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions
	defaultOpts := func() *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions {
		return &CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions{
			name:                 id,
			Enabled:              true,
			OauthClientId:        "foo",
			OauthClientSecret:    "bar",
			OauthAssertionIssuer: "foo",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})
	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_ASSERTION_ISSUER = 'foo' OAUTH_CLIENT_ID = 'foo'"+
			" OAUTH_CLIENT_SECRET = 'bar'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OauthAuthorizationEndpoint = Pointer("foo")
		opts.OauthTokenEndpoint = Pointer("foo")
		opts.OauthClientAuthMethod = Pointer(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)
		opts.OauthGrantJwtBearer = Pointer(true)
		opts.OauthAccessTokenValidity = Pointer(42)
		opts.OauthRefreshTokenValidity = Pointer(42)
		opts.Comment = Pointer("foo")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = API_AUTHENTICATION AUTH_TYPE = OAUTH2 ENABLED = true OAUTH_ASSERTION_ISSUER = 'foo'"+
			" OAUTH_AUTHORIZATION_ENDPOINT = 'foo' OAUTH_TOKEN_ENDPOINT = 'foo' OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST OAUTH_CLIENT_ID = 'foo' OAUTH_CLIENT_SECRET = 'bar' OAUTH_GRANT = JWT_BEARER"+
			" OAUTH_ACCESS_TOKEN_VALIDITY = 42 OAUTH_REFRESH_TOKEN_VALIDITY = 42 COMMENT = 'foo'", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateExternalOauth(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateExternalOauthSecurityIntegrationOptions
	defaultOpts := func() *CreateExternalOauthSecurityIntegrationOptions {
		return &CreateExternalOauthSecurityIntegrationOptions{
			name:                               id,
			Enabled:                            false,
			ExternalOauthType:                  ExternalOauthSecurityIntegrationTypeCustom,
			ExternalOauthIssuer:                "foo",
			ExternalOauthTokenUserMappingClaim: []TokenUserMappingClaim{{Claim: "foo"}},
			ExternalOauthSnowflakeUserMappingAttribute: ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateExternalOauthSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateExternalOauthSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: exactly one fields in [opts.ExternalOauthJwsKeysUrl opts.ExternalOauthRsaPublicKey]", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateExternalOauthSecurityIntegrationOptions", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey"))
		opts.ExternalOauthJwsKeysUrl = []JwsKeysUrl{{JwsKeyUrl: "foo"}}
		opts.ExternalOauthRsaPublicKey = Pointer("key")
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateExternalOauthSecurityIntegrationOptions", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey"))
	})
	t.Run("validation: conflicting fields for [opts.ExternalOauthJwsKeysUrl opts.ExternalOauthRsaPublicKey2]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalOauthJwsKeysUrl = []JwsKeysUrl{{JwsKeyUrl: "foo"}}
		opts.ExternalOauthRsaPublicKey2 = Pointer("key")
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateExternalOauthSecurityIntegrationOptions", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey2"))
	})
	t.Run("validation: conflicting fields for [opts.ExternalOauthAllowedRolesList opts.ExternalOauthBlockedRolesList]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalOauthAllowedRolesList = &AllowedRolesList{}
		opts.ExternalOauthBlockedRolesList = &BlockedRolesList{}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateExternalOauthSecurityIntegrationOptions", "ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList"))
	})
	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomAccountObjectIdentifier()
		opts.OrReplace = Bool(true)
		opts.ExternalOauthJwsKeysUrl = []JwsKeysUrl{{JwsKeyUrl: "foo"}}
		opts.ExternalOauthBlockedRolesList = &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{roleID}}
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = EXTERNAL_OAUTH ENABLED = false EXTERNAL_OAUTH_TYPE = CUSTOM EXTERNAL_OAUTH_ISSUER = 'foo'"+
			" EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM = ('foo') EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE = 'EMAIL_ADDRESS' EXTERNAL_OAUTH_JWS_KEYS_URL = ('foo')"+
			" EXTERNAL_OAUTH_BLOCKED_ROLES_LIST = (%s)", id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomAccountObjectIdentifier()
		opts.IfNotExists = Bool(true)
		opts.ExternalOauthAllowedRolesList = &AllowedRolesList{AllowedRolesList: []AccountObjectIdentifier{roleID}}
		opts.ExternalOauthRsaPublicKey = Pointer("foo")
		opts.ExternalOauthRsaPublicKey2 = Pointer("foo")
		opts.ExternalOauthAudienceList = &AudienceList{AudienceList: []AudienceListItem{{Item: "foo"}}}
		opts.ExternalOauthAnyRoleMode = Pointer(ExternalOauthSecurityIntegrationAnyRoleModeDisable)
		opts.ExternalOauthScopeDelimiter = Pointer(" ")
		opts.ExternalOauthScopeMappingAttribute = Pointer("foo")
		opts.Comment = Pointer("foo")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = EXTERNAL_OAUTH ENABLED = false EXTERNAL_OAUTH_TYPE = CUSTOM EXTERNAL_OAUTH_ISSUER = 'foo'"+
			" EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM = ('foo') EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE = 'EMAIL_ADDRESS' EXTERNAL_OAUTH_ALLOWED_ROLES_LIST = (%s)"+
			" EXTERNAL_OAUTH_RSA_PUBLIC_KEY = 'foo' EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2 = 'foo' EXTERNAL_OAUTH_AUDIENCE_LIST = ('foo') EXTERNAL_OAUTH_ANY_ROLE_MODE = DISABLE"+
			" EXTERNAL_OAUTH_SCOPE_DELIMITER = ' ' EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE = 'foo' COMMENT = 'foo'", id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateOauthForCustomClients(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateOauthForCustomClientsSecurityIntegrationOptions
	defaultOpts := func() *CreateOauthForCustomClientsSecurityIntegrationOptions {
		return &CreateOauthForCustomClientsSecurityIntegrationOptions{
			name:             id,
			OauthClientType:  OauthSecurityIntegrationClientTypePublic,
			OauthRedirectUri: "uri",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOauthForCustomClientsSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOauthForCustomClientsSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = OAUTH OAUTH_CLIENT = CUSTOM OAUTH_CLIENT_TYPE = 'PUBLIC' OAUTH_REDIRECT_URI = 'uri'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		roleID, role2ID, npID := randomAccountObjectIdentifier(), randomAccountObjectIdentifier(), randomAccountObjectIdentifier()
		opts.IfNotExists = Bool(true)
		opts.OauthClientType = OauthSecurityIntegrationClientTypePublic
		opts.Enabled = Pointer(true)
		opts.OauthAllowNonTlsRedirectUri = Pointer(true)
		opts.OauthEnforcePkce = Pointer(true)
		opts.OauthUseSecondaryRoles = Pointer(OauthSecurityIntegrationUseSecondaryRolesNone)
		opts.PreAuthorizedRolesList = &PreAuthorizedRolesList{PreAuthorizedRolesList: []AccountObjectIdentifier{roleID}}
		opts.BlockedRolesList = &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{role2ID}}
		opts.OauthIssueRefreshTokens = Pointer(true)
		opts.OauthRefreshTokenValidity = Pointer(42)
		opts.NetworkPolicy = securityIntegrationNetworkPolicyQuoted(&npID)
		opts.OauthClientRsaPublicKey = Pointer("key")
		opts.OauthClientRsaPublicKey2 = Pointer("key2")
		opts.Comment = Pointer("a")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = OAUTH OAUTH_CLIENT = CUSTOM OAUTH_CLIENT_TYPE = 'PUBLIC' OAUTH_REDIRECT_URI = 'uri' ENABLED = true"+
			" OAUTH_ALLOW_NON_TLS_REDIRECT_URI = true OAUTH_ENFORCE_PKCE = true OAUTH_USE_SECONDARY_ROLES = NONE PRE_AUTHORIZED_ROLES_LIST = (%s) BLOCKED_ROLES_LIST = (%s)"+
			" OAUTH_ISSUE_REFRESH_TOKENS = true OAUTH_REFRESH_TOKEN_VALIDITY = 42 NETWORK_POLICY = '%s' OAUTH_CLIENT_RSA_PUBLIC_KEY = 'key' OAUTH_CLIENT_RSA_PUBLIC_KEY_2 = 'key2' COMMENT = 'a'",
			id.FullyQualifiedName(), roleID.FullyQualifiedName(), role2ID.FullyQualifiedName(), npID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateOauthForPartnerApplications(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateOauthForPartnerApplicationsSecurityIntegrationOptions
	defaultOpts := func() *CreateOauthForPartnerApplicationsSecurityIntegrationOptions {
		return &CreateOauthForPartnerApplicationsSecurityIntegrationOptions{
			name:        id,
			OauthClient: OauthSecurityIntegrationClientTableauDesktop,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOauthForPartnerApplicationsSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOauthForPartnerApplicationsSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: OAUTH_REDIRECT_URI is required when OAUTH_CLIENT=LOOKER", func(t *testing.T) {
		opts := &CreateOauthForPartnerApplicationsSecurityIntegrationOptions{
			name:        id,
			OauthClient: OauthSecurityIntegrationClientLooker,
		}
		assertOptsInvalidJoinedErrors(t, opts, NewError("OauthRedirectUri is required when OauthClient is LOOKER"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = OAUTH OAUTH_CLIENT = TABLEAU_DESKTOP", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		blockedRoleID := randomAccountObjectIdentifier()
		opts.IfNotExists = Bool(true)
		opts.OauthClient = OauthSecurityIntegrationClientLooker
		opts.OauthRedirectUri = Pointer("uri")
		opts.Enabled = Pointer(true)
		opts.OauthIssueRefreshTokens = Pointer(true)
		opts.OauthRefreshTokenValidity = Pointer(42)
		opts.OauthUseSecondaryRoles = Pointer(OauthSecurityIntegrationUseSecondaryRolesNone)
		opts.BlockedRolesList = &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{blockedRoleID}}
		opts.Comment = Pointer("a")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = OAUTH OAUTH_CLIENT = LOOKER OAUTH_REDIRECT_URI = 'uri' ENABLED = true OAUTH_ISSUE_REFRESH_TOKENS = true"+
			" OAUTH_REFRESH_TOKEN_VALIDITY = 42 OAUTH_USE_SECONDARY_ROLES = NONE BLOCKED_ROLES_LIST = (%s) COMMENT = 'a'", id.FullyQualifiedName(), blockedRoleID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateSaml2(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateSaml2SecurityIntegrationOptions
	defaultOpts := func() *CreateSaml2SecurityIntegrationOptions {
		return &CreateSaml2SecurityIntegrationOptions{
			name:          id,
			Saml2Issuer:   "issuer",
			Saml2SsoUrl:   "url",
			Saml2Provider: "provider",
			Saml2X509Cert: "cert",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSaml2SecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSaml2SecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = SAML2 SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Enabled = Pointer(true)
		opts.AllowedEmailPatterns = []EmailPattern{{Pattern: "pattern"}}
		opts.AllowedUserDomains = []UserDomain{{Domain: "domain"}}
		opts.Comment = Pointer("a")
		opts.Saml2EnableSpInitiated = Pointer(true)
		opts.Saml2ForceAuthn = Pointer(true)
		opts.Saml2PostLogoutRedirectUrl = Pointer("redirect")
		opts.Saml2RequestedNameidFormat = Pointer(Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos)
		opts.Saml2SignRequest = Pointer(true)
		opts.Saml2SnowflakeAcsUrl = Pointer("acs")
		opts.Saml2SnowflakeIssuerUrl = Pointer("issuer")
		opts.Saml2SpInitiatedLoginPageLabel = Pointer("label")
		opts.Saml2SnowflakeX509Cert = Pointer("cert")

		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = SAML2 ENABLED = true SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'"+
			" ALLOWED_USER_DOMAINS = ('domain') ALLOWED_EMAIL_PATTERNS = ('pattern') SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'label' SAML2_ENABLE_SP_INITIATED = true SAML2_SNOWFLAKE_X509_CERT = 'cert' SAML2_SIGN_REQUEST = true"+
			" SAML2_REQUESTED_NAMEID_FORMAT = '%s' SAML2_POST_LOGOUT_REDIRECT_URL = 'redirect' SAML2_FORCE_AUTHN = true SAML2_SNOWFLAKE_ISSUER_URL = 'issuer' SAML2_SNOWFLAKE_ACS_URL = 'acs'"+
			" COMMENT = 'a'", id.FullyQualifiedName(), Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos)
	})
}

func TestSecurityIntegrations_CreateScim(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateScimSecurityIntegrationOptions
	defaultOpts := func() *CreateScimSecurityIntegrationOptions {
		return &CreateScimSecurityIntegrationOptions{
			name:       id,
			ScimClient: "GENERIC",
			RunAsRole:  "GENERIC_SCIM_PROVISIONER",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateScimSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateScimSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: conflicting SyncPassword for Azure ScimClient", func(t *testing.T) {
		opts := defaultOpts()
		opts.ScimClient = ScimSecurityIntegrationScimClientAzure
		opts.SyncPassword = Pointer(true)
		assertOptsInvalidJoinedErrors(t, opts, NewError("SyncPassword is not supported for Azure scim client"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = SCIM SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = 'GENERIC_SCIM_PROVISIONER'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		networkPolicyID := randomAccountObjectIdentifier()
		opts.Enabled = Pointer(true)
		opts.IfNotExists = Pointer(true)
		opts.NetworkPolicy = securityIntegrationNetworkPolicyQuoted(&networkPolicyID)
		opts.SyncPassword = Pointer(true)
		opts.Comment = Pointer("a")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = SCIM ENABLED = true SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = 'GENERIC_SCIM_PROVISIONER'"+
			" NETWORK_POLICY = '%s' SYNC_PASSWORD = true COMMENT = 'a'", id.FullyQualifiedName(), networkPolicyID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterApiAuthenticationWithClientCredentialsFlow(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions
	defaultOpts := func() *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions {
		return &AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithClientCredentialsFlowIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithClientCredentialsFlowIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions.Set", "Enabled", "OauthTokenEndpoint",
			"OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantClientCredentials", "OauthAccessTokenValidity", "OauthRefreshTokenValidity", "OauthAllowedScopes", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiAuthenticationWithClientCredentialsFlowIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions.Unset",
			"Enabled", "Comment"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithClientCredentialsFlowIntegrationSet{}
		opts.Unset = &ApiAuthenticationWithClientCredentialsFlowIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithClientCredentialsFlowIntegrationSet{
			Enabled:                     Pointer(true),
			OauthTokenEndpoint:          Pointer("foo"),
			OauthClientAuthMethod:       Pointer(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			OauthClientId:               Pointer("foo"),
			OauthClientSecret:           Pointer("foo"),
			OauthGrantClientCredentials: Pointer(true),
			OauthAccessTokenValidity:    Pointer(42),
			OauthRefreshTokenValidity:   Pointer(42),
			OauthAllowedScopes:          []AllowedScope{{Scope: "foo"}},
			Comment:                     Pointer("foo"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_TOKEN_ENDPOINT = 'foo', OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST,"+
			" OAUTH_CLIENT_ID = 'foo', OAUTH_CLIENT_SECRET = 'foo', OAUTH_GRANT = CLIENT_CREDENTIALS, OAUTH_ACCESS_TOKEN_VALIDITY = 42,"+
			" OAUTH_REFRESH_TOKEN_VALIDITY = 42, OAUTH_ALLOWED_SCOPES = ('foo'), COMMENT = 'foo'", id.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiAuthenticationWithClientCredentialsFlowIntegrationUnset{
			Enabled: Pointer(true),
			Comment: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, COMMENT", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterApiAuthenticationWithAuthorizationCodeFlow(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions
	defaultOpts := func() *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions {
		return &AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions.Set", "Enabled", "OauthAuthorizationEndpoint", "OauthTokenEndpoint",
			"OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantAuthorizationCode", "OauthAccessTokenValidity", "OauthRefreshTokenValidity", "OauthAllowedScopes", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions.Unset",
			"Enabled", "Comment"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSet{}
		opts.Unset = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSet{
			Enabled:                     Pointer(true),
			OauthTokenEndpoint:          Pointer("foo"),
			OauthClientAuthMethod:       Pointer(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			OauthClientId:               Pointer("foo"),
			OauthClientSecret:           Pointer("foo"),
			OauthGrantAuthorizationCode: Pointer(true),
			OauthAccessTokenValidity:    Pointer(42),
			OauthRefreshTokenValidity:   Pointer(42),
			OauthAllowedScopes:          []AllowedScope{{Scope: "bar"}},
			Comment:                     Pointer("foo"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_TOKEN_ENDPOINT = 'foo', OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST,"+
			" OAUTH_CLIENT_ID = 'foo', OAUTH_CLIENT_SECRET = 'foo', OAUTH_GRANT = AUTHORIZATION_CODE, OAUTH_ACCESS_TOKEN_VALIDITY = 42,"+
			" OAUTH_REFRESH_TOKEN_VALIDITY = 42, OAUTH_ALLOWED_SCOPES = ('bar'), COMMENT = 'foo'", id.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnset{
			Enabled: Pointer(true),
			Comment: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, COMMENT", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterApiAuthenticationWithJwtBearerFlow(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions
	defaultOpts := func() *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions {
		return &AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithJwtBearerFlowIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithJwtBearerFlowIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions.Set", "Enabled", "OauthAuthorizationEndpoint", "OauthTokenEndpoint",
			"OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantJwtBearer", "OauthAccessTokenValidity", "OauthRefreshTokenValidity", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiAuthenticationWithJwtBearerFlowIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions.Unset",
			"Enabled", "Comment"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithJwtBearerFlowIntegrationSet{}
		opts.Unset = &ApiAuthenticationWithJwtBearerFlowIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiAuthenticationWithJwtBearerFlowIntegrationSet{
			Enabled:                   Pointer(true),
			OauthTokenEndpoint:        Pointer("foo"),
			OauthClientAuthMethod:     Pointer(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			OauthClientId:             Pointer("foo"),
			OauthClientSecret:         Pointer("foo"),
			OauthGrantJwtBearer:       Pointer(true),
			OauthAccessTokenValidity:  Pointer(42),
			OauthRefreshTokenValidity: Pointer(42),
			Comment:                   Pointer("foo"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_TOKEN_ENDPOINT = 'foo', OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST,"+
			" OAUTH_CLIENT_ID = 'foo', OAUTH_CLIENT_SECRET = 'foo', OAUTH_GRANT = JWT_BEARER, OAUTH_ACCESS_TOKEN_VALIDITY = 42,"+
			" OAUTH_REFRESH_TOKEN_VALIDITY = 42, COMMENT = 'foo'", id.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiAuthenticationWithJwtBearerFlowIntegrationUnset{
			Enabled: Pointer(true),
			Comment: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, COMMENT", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterExternalOauth(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterExternalOauthSecurityIntegrationOptions
	defaultOpts := func() *AlterExternalOauthSecurityIntegrationOptions {
		return &AlterExternalOauthSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterExternalOauthSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalOauthIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalOauthSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalOauthIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterExternalOauthSecurityIntegrationOptions.Set", "Enabled", "ExternalOauthType",
			"ExternalOauthIssuer", "ExternalOauthTokenUserMappingClaim", "ExternalOauthSnowflakeUserMappingAttribute", "ExternalOauthJwsKeysUrl",
			"ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList", "ExternalOauthRsaPublicKey", "ExternalOauthRsaPublicKey2", "ExternalOauthAudienceList",
			"ExternalOauthAnyRoleMode", "ExternalOauthScopeDelimiter", "ExternalOauthScopeMappingAttribute", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ExternalOauthIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterExternalOauthSecurityIntegrationOptions.Unset",
			"Enabled", "ExternalOauthAudienceList"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalOauthIntegrationSet{}
		opts.Unset = &ExternalOauthIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalOauthSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalOauthJwsKeysUrl opts.ExternalOauthRsaPublicKey]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalOauthIntegrationSet{
			ExternalOauthJwsKeysUrl:   []JwsKeysUrl{{JwsKeyUrl: "foo"}},
			ExternalOauthRsaPublicKey: Pointer("key"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterExternalOauthSecurityIntegrationOptions.Set", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey"))
	})
	t.Run("validation: conflicting fields for [opts.ExternalOauthJwsKeysUrl opts.ExternalOauthRsaPublicKey2]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalOauthIntegrationSet{
			ExternalOauthJwsKeysUrl:    []JwsKeysUrl{{JwsKeyUrl: "foo"}},
			ExternalOauthRsaPublicKey2: Pointer("key"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterExternalOauthSecurityIntegrationOptions.Set", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey2"))
	})
	t.Run("validation: conflicting fields for [opts.ExternalOauthAllowedRolesList opts.ExternalOauthBlockedRolesList]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalOauthIntegrationSet{
			ExternalOauthAllowedRolesList: &AllowedRolesList{},
			ExternalOauthBlockedRolesList: &BlockedRolesList{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterExternalOauthSecurityIntegrationOptions.Set", "ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList"))
	})
	t.Run("empty lists", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalOauthIntegrationSet{
			ExternalOauthBlockedRolesList: &BlockedRolesList{},
			ExternalOauthAudienceList:     &AudienceList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET EXTERNAL_OAUTH_BLOCKED_ROLES_LIST = (), EXTERNAL_OAUTH_AUDIENCE_LIST = ()", id.FullyQualifiedName())
		opts.Set = &ExternalOauthIntegrationSet{
			ExternalOauthAllowedRolesList: &AllowedRolesList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET EXTERNAL_OAUTH_ALLOWED_ROLES_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomAccountObjectIdentifier()
		opts.Set = &ExternalOauthIntegrationSet{
			Enabled:                            Pointer(true),
			ExternalOauthType:                  Pointer(ExternalOauthSecurityIntegrationTypeCustom),
			ExternalOauthIssuer:                Pointer("foo"),
			ExternalOauthTokenUserMappingClaim: []TokenUserMappingClaim{{Claim: "foo"}},
			ExternalOauthSnowflakeUserMappingAttribute: Pointer(ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress),
			ExternalOauthAllowedRolesList:              &AllowedRolesList{AllowedRolesList: []AccountObjectIdentifier{roleID}},
			ExternalOauthRsaPublicKey:                  Pointer("foo"),
			ExternalOauthRsaPublicKey2:                 Pointer("foo"),
			ExternalOauthAudienceList:                  &AudienceList{AudienceList: []AudienceListItem{{Item: "foo"}}},
			ExternalOauthAnyRoleMode:                   Pointer(ExternalOauthSecurityIntegrationAnyRoleModeDisable),
			ExternalOauthScopeDelimiter:                Pointer(" "),
			ExternalOauthScopeMappingAttribute:         Pointer("foo"),
			Comment:                                    Pointer(StringAllowEmpty{Value: "foo"}),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, EXTERNAL_OAUTH_TYPE = CUSTOM, EXTERNAL_OAUTH_ISSUER = 'foo',"+
			" EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM = ('foo'), EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE = 'EMAIL_ADDRESS', EXTERNAL_OAUTH_ALLOWED_ROLES_LIST = (%s),"+
			" EXTERNAL_OAUTH_RSA_PUBLIC_KEY = 'foo', EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2 = 'foo', EXTERNAL_OAUTH_AUDIENCE_LIST = ('foo'), EXTERNAL_OAUTH_ANY_ROLE_MODE = DISABLE,"+
			" EXTERNAL_OAUTH_SCOPE_DELIMITER = ' ', EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE = 'foo', COMMENT = 'foo'", id.FullyQualifiedName(), roleID.FullyQualifiedName())
		opts.Set = &ExternalOauthIntegrationSet{
			ExternalOauthBlockedRolesList: &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{roleID}},
			ExternalOauthJwsKeysUrl:       []JwsKeysUrl{{JwsKeyUrl: "foo"}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET EXTERNAL_OAUTH_JWS_KEYS_URL = ('foo'), EXTERNAL_OAUTH_BLOCKED_ROLES_LIST = (%s)", id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ExternalOauthIntegrationUnset{
			Enabled:                   Pointer(true),
			ExternalOauthAudienceList: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, EXTERNAL_OAUTH_AUDIENCE_LIST", id.FullyQualifiedName())
	})

	t.Run("set empty comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalOauthIntegrationSet{
			Comment: Pointer(StringAllowEmpty{Value: ""}),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET COMMENT = ''", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterOauthForPartnerApplications(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterOauthForPartnerApplicationsSecurityIntegrationOptions
	defaultOpts := func() *AlterOauthForPartnerApplicationsSecurityIntegrationOptions {
		return &AlterOauthForPartnerApplicationsSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterOauthForPartnerApplicationsSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions.Set", "Enabled", "OauthIssueRefreshTokens",
			"OauthRedirectUri", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles", "BlockedRolesList", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OauthForPartnerApplicationsIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions.Unset",
			"Enabled", "OauthUseSecondaryRoles"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{}
		opts.Unset = &OauthForPartnerApplicationsIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("empty roles lists", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{
			BlockedRolesList: &BlockedRolesList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET BLOCKED_ROLES_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomAccountObjectIdentifier()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{
			Enabled:                   Pointer(true),
			OauthRedirectUri:          Pointer("uri"),
			OauthIssueRefreshTokens:   Pointer(true),
			OauthRefreshTokenValidity: Pointer(42),
			OauthUseSecondaryRoles:    Pointer(OauthSecurityIntegrationUseSecondaryRolesNone),
			BlockedRolesList:          &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{roleID}},
			Comment:                   Pointer(StringAllowEmpty{""}),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_ISSUE_REFRESH_TOKENS = true, OAUTH_REDIRECT_URI = 'uri', OAUTH_REFRESH_TOKEN_VALIDITY = 42,"+
			" OAUTH_USE_SECONDARY_ROLES = NONE, BLOCKED_ROLES_LIST = (%s), COMMENT = ''", id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OauthForPartnerApplicationsIntegrationUnset{
			Enabled:                Pointer(true),
			OauthUseSecondaryRoles: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, OAUTH_USE_SECONDARY_ROLES", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterOauthForCustomClients(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterOauthForCustomClientsSecurityIntegrationOptions
	defaultOpts := func() *AlterOauthForCustomClientsSecurityIntegrationOptions {
		return &AlterOauthForCustomClientsSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterOauthForCustomClientsSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForCustomClientsIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForCustomClientsIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions.Set", "Enabled", "OauthRedirectUri", "OauthAllowNonTlsRedirectUri",
			"OauthEnforcePkce", "PreAuthorizedRolesList", "BlockedRolesList", "OauthIssueRefreshTokens", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles",
			"NetworkPolicy", "OauthClientRsaPublicKey", "OauthClientRsaPublicKey2", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OauthForCustomClientsIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions.Unset",
			"Enabled", "NetworkPolicy", "OauthUseSecondaryRoles", "OauthClientRsaPublicKey", "OauthClientRsaPublicKey2"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForCustomClientsIntegrationSet{}
		opts.Unset = &OauthForCustomClientsIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("empty roles lists", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForCustomClientsIntegrationSet{
			PreAuthorizedRolesList: &PreAuthorizedRolesList{},
			BlockedRolesList:       &BlockedRolesList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET PRE_AUTHORIZED_ROLES_LIST = (), BLOCKED_ROLES_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		roleID, role2ID, npID := randomAccountObjectIdentifier(), randomAccountObjectIdentifier(), randomAccountObjectIdentifier()
		opts.Set = &OauthForCustomClientsIntegrationSet{
			Enabled:                     Pointer(true),
			OauthRedirectUri:            Pointer("uri"),
			OauthAllowNonTlsRedirectUri: Pointer(true),
			OauthEnforcePkce:            Pointer(true),
			OauthUseSecondaryRoles:      Pointer(OauthSecurityIntegrationUseSecondaryRolesNone),
			PreAuthorizedRolesList:      &PreAuthorizedRolesList{PreAuthorizedRolesList: []AccountObjectIdentifier{roleID}},
			BlockedRolesList:            &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{role2ID}},
			OauthIssueRefreshTokens:     Pointer(true),
			OauthRefreshTokenValidity:   Pointer(42),
			NetworkPolicy:               securityIntegrationNetworkPolicyQuoted(&npID),
			OauthClientRsaPublicKey:     Pointer("key"),
			OauthClientRsaPublicKey2:    Pointer("key2"),
			Comment:                     Pointer("a"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_REDIRECT_URI = 'uri', OAUTH_ALLOW_NON_TLS_REDIRECT_URI = true, OAUTH_ENFORCE_PKCE = true,"+
			" PRE_AUTHORIZED_ROLES_LIST = (%s), BLOCKED_ROLES_LIST = (%s), OAUTH_ISSUE_REFRESH_TOKENS = true, OAUTH_REFRESH_TOKEN_VALIDITY = 42, OAUTH_USE_SECONDARY_ROLES = NONE,"+
			" NETWORK_POLICY = '%s', OAUTH_CLIENT_RSA_PUBLIC_KEY = 'key', OAUTH_CLIENT_RSA_PUBLIC_KEY_2 = 'key2', COMMENT = 'a'", id.FullyQualifiedName(), roleID.FullyQualifiedName(), role2ID.FullyQualifiedName(), npID.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OauthForCustomClientsIntegrationUnset{
			Enabled:                  Pointer(true),
			OauthUseSecondaryRoles:   Pointer(true),
			NetworkPolicy:            Pointer(true),
			OauthClientRsaPublicKey:  Pointer(true),
			OauthClientRsaPublicKey2: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, NETWORK_POLICY, OAUTH_CLIENT_RSA_PUBLIC_KEY, OAUTH_CLIENT_RSA_PUBLIC_KEY_2, OAUTH_USE_SECONDARY_ROLES", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterSaml2(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterSaml2IntegrationSecurityIntegrationOptions
	defaultOpts := func() *AlterSaml2SecurityIntegrationOptions {
		return &AlterSaml2SecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSaml2SecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Saml2IntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSaml2SecurityIntegrationOptions", "Set", "Unset", "RefreshSaml2SnowflakePrivateKey", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Saml2IntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSaml2SecurityIntegrationOptions.Set", "Enabled", "Saml2Issuer", "Saml2SsoUrl", "Saml2Provider",
			"Saml2X509Cert", "AllowedUserDomains", "AllowedEmailPatterns", "Saml2SpInitiatedLoginPageLabel", "Saml2EnableSpInitiated", "Saml2SnowflakeX509Cert", "Saml2SignRequest",
			"Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Saml2ForceAuthn", "Saml2SnowflakeIssuerUrl", "Saml2SnowflakeAcsUrl", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &Saml2IntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSaml2SecurityIntegrationOptions.Unset",
			"Saml2ForceAuthn", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Comment"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Saml2IntegrationSet{}
		opts.Unset = &Saml2IntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSaml2SecurityIntegrationOptions", "Set", "Unset", "RefreshSaml2SnowflakePrivateKey", "SetTags", "UnsetTags"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Saml2IntegrationSet{
			Enabled:                        Pointer(true),
			Saml2Issuer:                    Pointer("issuer"),
			Saml2SsoUrl:                    Pointer("url"),
			Saml2Provider:                  Pointer(Saml2SecurityIntegrationSaml2ProviderCustom),
			Saml2X509Cert:                  Pointer("cert"),
			AllowedUserDomains:             []UserDomain{{Domain: "domain"}},
			AllowedEmailPatterns:           []EmailPattern{{Pattern: "pattern"}},
			Saml2SpInitiatedLoginPageLabel: Pointer("label"),
			Saml2EnableSpInitiated:         Pointer(true),
			Saml2SnowflakeX509Cert:         Pointer("cert"),
			Saml2SignRequest:               Pointer(true),
			Saml2RequestedNameidFormat:     Pointer(Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos),
			Saml2PostLogoutRedirectUrl:     Pointer("redirect"),
			Saml2ForceAuthn:                Pointer(true),
			Saml2SnowflakeIssuerUrl:        Pointer("issuer"),
			Saml2SnowflakeAcsUrl:           Pointer("acs"),
			Comment:                        Pointer("a"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, SAML2_ISSUER = 'issuer', SAML2_SSO_URL = 'url', SAML2_PROVIDER = 'CUSTOM', SAML2_X509_CERT = 'cert',"+
			" ALLOWED_USER_DOMAINS = ('domain'), ALLOWED_EMAIL_PATTERNS = ('pattern'), SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'label', SAML2_ENABLE_SP_INITIATED = true, SAML2_SNOWFLAKE_X509_CERT = 'cert', SAML2_SIGN_REQUEST = true,"+
			" SAML2_REQUESTED_NAMEID_FORMAT = '%s', SAML2_POST_LOGOUT_REDIRECT_URL = 'redirect', SAML2_FORCE_AUTHN = true, SAML2_SNOWFLAKE_ISSUER_URL = 'issuer', SAML2_SNOWFLAKE_ACS_URL = 'acs',"+
			" COMMENT = 'a'", id.FullyQualifiedName(), Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos)
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &Saml2IntegrationUnset{
			Saml2ForceAuthn:            Pointer(true),
			Saml2RequestedNameidFormat: Pointer(true),
			Saml2PostLogoutRedirectUrl: Pointer(true),
			Comment:                    Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET SAML2_FORCE_AUTHN, SAML2_REQUESTED_NAMEID_FORMAT, SAML2_POST_LOGOUT_REDIRECT_URL, COMMENT", id.FullyQualifiedName())
	})

	t.Run("refresh SAML2_SNOWFLAKE_PRIVATE_KEY", func(t *testing.T) {
		opts := defaultOpts()
		opts.RefreshSaml2SnowflakePrivateKey = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterScim(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterScimSecurityIntegrationOptions
	defaultOpts := func() *AlterScimSecurityIntegrationOptions {
		return &AlterScimSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterScimSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ScimIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterScimSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ScimIntegrationSet{}
		opts.Unset = &ScimIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterScimSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ScimIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterScimSecurityIntegrationOptions.Set", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ScimIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterScimSecurityIntegrationOptions.Unset", "Enabled", "NetworkPolicy", "SyncPassword"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		networkPolicyID := randomAccountObjectIdentifier()
		opts.Set = &ScimIntegrationSet{
			Enabled:       Pointer(true),
			NetworkPolicy: securityIntegrationNetworkPolicyQuoted(&networkPolicyID),
			SyncPassword:  Pointer(true),
			Comment:       Pointer(StringAllowEmpty{Value: "test"}),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, NETWORK_POLICY = '%s', SYNC_PASSWORD = true, COMMENT = 'test'",
			id.FullyQualifiedName(), networkPolicyID.FullyQualifiedName())
	})

	t.Run("set empty comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ScimIntegrationSet{
			Comment: Pointer(StringAllowEmpty{Value: ""}),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET COMMENT = ''", id.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ScimIntegrationUnset{
			Enabled:       Pointer(true),
			NetworkPolicy: Pointer(true),
			SyncPassword:  Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, NETWORK_POLICY, SYNC_PASSWORD", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid DropSecurityIntegrationOptions
	defaultOpts := func() *DropSecurityIntegrationOptions {
		return &DropSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP SECURITY INTEGRATION IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid DescribeSecurityIntegrationOptions
	defaultOpts := func() *DescribeSecurityIntegrationOptions {
		return &DescribeSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE SECURITY INTEGRATION %s", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Show(t *testing.T) {
	// Minimal valid ShowSecurityIntegrationOptions
	defaultOpts := func() *ShowSecurityIntegrationOptions {
		return &ShowSecurityIntegrationOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECURITY INTEGRATIONS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("some pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECURITY INTEGRATIONS LIKE 'some pattern'")
	})
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
