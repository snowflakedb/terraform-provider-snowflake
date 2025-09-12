//go:build account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-2312051): Set up proper test level for each test.
func TestInt_SecurityIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	acsURL := testClientHelper().Context.ACSURL(t)
	issuerURL := testClientHelper().Context.IssuerURL(t)
	cert := random.GenerateX509(t)
	rsaKey, rsaKeyHash := random.GenerateRSAPublicKey(t)

	revertParameter := testClientHelper().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterEnableIdentifierFirstLogin, "true")
	t.Cleanup(revertParameter)

	createApiAuthClientCred := func(t *testing.T, with func(*sdk.CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest)) (*sdk.SecurityIntegration, sdk.AccountObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id, false, "foo", "foo")
		if with != nil {
			with(req)
		}
		err := client.SecurityIntegrations.CreateApiAuthenticationWithClientCredentialsFlow(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration, id
	}
	createApiAuthCodeGrant := func(t *testing.T, with func(*sdk.CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest)) (*sdk.SecurityIntegration, sdk.AccountObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id, false, "foo", "foo")
		if with != nil {
			with(req)
		}
		err := client.SecurityIntegrations.CreateApiAuthenticationWithAuthorizationCodeGrantFlow(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration, id
	}
	createApiAuthJwtBearer := func(t *testing.T, with func(*sdk.CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest)) (*sdk.SecurityIntegration, sdk.AccountObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest(id, false, "foo", "foo", "foo")
		if with != nil {
			with(req)
		}
		err := client.SecurityIntegrations.CreateApiAuthenticationWithJwtBearerFlow(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration, id
	}
	createExternalOauth := func(t *testing.T, with func(*sdk.CreateExternalOauthSecurityIntegrationRequest)) (*sdk.SecurityIntegration, sdk.AccountObjectIdentifier, string) {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		issuer := random.String()
		req := sdk.NewCreateExternalOauthSecurityIntegrationRequest(id, false, sdk.ExternalOauthSecurityIntegrationTypeCustom,
			issuer, []sdk.TokenUserMappingClaim{{Claim: "foo"}}, sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName,
		)
		if with != nil {
			with(req)
		}
		err := client.SecurityIntegrations.CreateExternalOauth(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration, id, issuer
	}
	createOauthCustom := func(t *testing.T, with func(*sdk.CreateOauthForCustomClientsSecurityIntegrationRequest)) (*sdk.SecurityIntegration, sdk.AccountObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateOauthForCustomClientsSecurityIntegrationRequest(id, sdk.OauthSecurityIntegrationClientTypePublic, "https://example.com")
		if with != nil {
			with(req)
		}
		err := client.SecurityIntegrations.CreateOauthForCustomClients(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration, id
	}
	createOauthPartner := func(t *testing.T, with func(*sdk.CreateOauthForPartnerApplicationsSecurityIntegrationRequest)) (*sdk.SecurityIntegration, sdk.AccountObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateOauthForPartnerApplicationsSecurityIntegrationRequest(id, sdk.OauthSecurityIntegrationClientLooker).
			WithOauthRedirectUri("http://example.com")

		if with != nil {
			with(req)
		}
		err := client.SecurityIntegrations.CreateOauthForPartnerApplications(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration, id
	}
	createSAML2Integration := func(t *testing.T, with func(*sdk.CreateSaml2SecurityIntegrationRequest)) (*sdk.SecurityIntegration, sdk.AccountObjectIdentifier, string) {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		issuer := testClientHelper().Ids.Alpha()
		saml2Req := sdk.NewCreateSaml2SecurityIntegrationRequest(id, issuer, "https://example.com", sdk.Saml2SecurityIntegrationSaml2ProviderCustom, cert)
		if with != nil {
			with(saml2Req)
		}
		err := client.SecurityIntegrations.CreateSaml2(ctx, saml2Req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration, id, issuer
	}

	createSCIMIntegration := func(t *testing.T, with func(*sdk.CreateScimSecurityIntegrationRequest)) (*sdk.SecurityIntegration, sdk.AccountObjectIdentifier) {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		scimReq := sdk.NewCreateScimSecurityIntegrationRequest(id, sdk.ScimSecurityIntegrationScimClientGeneric, sdk.ScimSecurityIntegrationRunAsRoleGenericScimProvisioner)
		if with != nil {
			with(scimReq)
		}
		err := client.SecurityIntegrations.CreateScim(ctx, scimReq)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration, id
	}

	assertSecurityIntegration := func(t *testing.T, si *sdk.SecurityIntegration, id sdk.AccountObjectIdentifier, siType string, enabled bool, comment string) {
		t.Helper()
		assert.Equal(t, id.Name(), si.Name)
		assert.Equal(t, siType, si.IntegrationType)
		assert.Equal(t, enabled, si.Enabled)
		assert.Equal(t, comment, si.Comment)
		assert.Equal(t, "SECURITY", si.Category)
	}

	// TODO(SNOW-1449579): move to helpers
	assertFieldContainsList := func(details []sdk.SecurityIntegrationProperty, field, value, sep string) {
		found, err := collections.FindFirst(details, func(d sdk.SecurityIntegrationProperty) bool { return d.Name == field })
		assert.NoError(t, err)
		values := strings.Split(found.Value, sep)
		for _, exp := range strings.Split(value, sep) {
			assert.Contains(t, values, exp)
		}
	}

	type apiAuthDetails struct {
		enabled                    string
		oauthAccessTokenValidity   string
		oauthRefreshTokenValidity  string
		oauthClientId              string
		oauthClientAuthMethod      string
		oauthAuthorizationEndpoint string
		oauthTokenEndpoint         string
		oauthAllowedScopes         string
		oauthGrant                 string
		parentIntegration          string
		authType                   string
		oauthAssertionIssuer       string
		comment                    string
	}

	assertApiAuth := func(details []sdk.SecurityIntegrationProperty, d apiAuthDetails) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: d.enabled, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ACCESS_TOKEN_VALIDITY", Type: "Long", Value: d.oauthAccessTokenValidity, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_REFRESH_TOKEN_VALIDITY", Type: "Long", Value: d.oauthRefreshTokenValidity, Default: "7776000"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_ID", Type: "String", Value: d.oauthClientId, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_AUTH_METHOD", Type: "String", Value: d.oauthClientAuthMethod, Default: "CLIENT_SECRET_BASIC"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_AUTHORIZATION_ENDPOINT", Type: "String", Value: d.oauthAuthorizationEndpoint, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_TOKEN_ENDPOINT", Type: "String", Value: d.oauthTokenEndpoint, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ALLOWED_SCOPES", Type: "List", Value: d.oauthAllowedScopes, Default: "[]"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_GRANT", Type: "String", Value: d.oauthGrant, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "PARENT_INTEGRATION", Type: "String", Value: d.parentIntegration, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "AUTH_TYPE", Type: "String", Value: d.authType, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ASSERTION_ISSUER", Type: "String", Value: d.oauthAssertionIssuer, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: d.comment, Default: ""})
	}

	type externalOauthDetails struct {
		enabled                                    string
		externalOauthIssuer                        string
		externalOauthJwsKeysUrl                    string
		externalOauthAnyRoleMode                   string
		externalOauthScopeMappingAttribute         string
		externalOauthRsaPublicKey                  string
		externalOauthRsaPublicKey2                 string
		externalOauthBlockedRolesList              string
		externalOauthAllowedRolesList              string
		externalOauthAudienceList                  string
		externalOauthTokenUserMappingClaim         string
		externalOauthSnowflakeUserMappingAttribute string
		externalOauthScopeDelimiter                string
		comment                                    string
	}

	assertExternalOauth := func(details []sdk.SecurityIntegrationProperty, d externalOauthDetails) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: d.enabled, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_ISSUER", Type: "String", Value: d.externalOauthIssuer, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_JWS_KEYS_URL", Type: "Object", Value: d.externalOauthJwsKeysUrl, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_ANY_ROLE_MODE", Type: "String", Value: d.externalOauthAnyRoleMode, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE", Type: "String", Value: d.externalOauthScopeMappingAttribute, Default: "scp"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_RSA_PUBLIC_KEY", Type: "String", Value: d.externalOauthRsaPublicKey, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2", Type: "String", Value: d.externalOauthRsaPublicKey2, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_ALLOWED_ROLES_LIST", Type: "List", Value: d.externalOauthAllowedRolesList, Default: "[]"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_AUDIENCE_LIST", Type: "List", Value: d.externalOauthAudienceList, Default: "[]"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM", Type: "Object", Value: d.externalOauthTokenUserMappingClaim, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE", Type: "String", Value: d.externalOauthSnowflakeUserMappingAttribute, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_SCOPE_DELIMITER", Type: "String", Value: d.externalOauthScopeDelimiter, Default: ","})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: d.comment, Default: ""})
		assertFieldContainsList(details, "EXTERNAL_OAUTH_BLOCKED_ROLES_LIST", d.externalOauthBlockedRolesList, ",")
	}

	type oauthPartnerDetails struct {
		enabled                 string
		oauthIssueRefreshTokens string
		refreshTokenValidity    string
		useSecondaryRoles       string
		preAuthorizedRolesList  string
		blockedRolesList        string
		networkPolicy           string
		comment                 string
	}

	assertOauthPartner := func(details []sdk.SecurityIntegrationProperty, d oauthPartnerDetails) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: d.enabled, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ISSUE_REFRESH_TOKENS", Type: "Boolean", Value: d.oauthIssueRefreshTokens, Default: "true"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_REFRESH_TOKEN_VALIDITY", Type: "Long", Value: d.refreshTokenValidity, Default: "7776000"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_USE_SECONDARY_ROLES", Type: "String", Value: d.useSecondaryRoles, Default: "NONE"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "PRE_AUTHORIZED_ROLES_LIST", Type: "List", Value: d.preAuthorizedRolesList, Default: "[]"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: d.networkPolicy, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: d.comment, Default: ""})
		assertFieldContainsList(details, "BLOCKED_ROLES_LIST", d.blockedRolesList, ",")
	}

	assertOauthCustom := func(details []sdk.SecurityIntegrationProperty, d oauthPartnerDetails, allowNonTlsRedirectUri, clientType, enforcePkce, key1Hash, key2Hash string) {
		assertOauthPartner(details, d)
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ALLOW_NON_TLS_REDIRECT_URI", Type: "Boolean", Value: allowNonTlsRedirectUri, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_TYPE", Type: "String", Value: clientType, Default: "CONFIDENTIAL"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ENFORCE_PKCE", Type: "Boolean", Value: enforcePkce, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_RSA_PUBLIC_KEY_FP", Type: "String", Value: fmt.Sprintf("SHA256:%s", key1Hash), Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP", Type: "String", Value: fmt.Sprintf("SHA256:%s", key2Hash), Default: ""})
	}

	assertSCIMDescribe := func(details []sdk.SecurityIntegrationProperty, enabled, networkPolicy, runAsRole, syncPassword, comment string) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: enabled, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: networkPolicy, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "RUN_AS_ROLE", Type: "String", Value: runAsRole, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SYNC_PASSWORD", Type: "Boolean", Value: syncPassword, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: comment, Default: ""})
	}

	type saml2Details struct {
		provider                  string
		enableSPInitiated         string
		spInitiatedLoginPageLabel string
		ssoURL                    string
		issuer                    string
		requestedNameIDFormat     string
		forceAuthn                string
		postLogoutRedirectUrl     string
		signrequest               string
		comment                   string
		snowflakeIssuerURL        string
		snowflakeAcsURL           string
		allowedUserDomains        string
		allowedEmailPatterns      string
	}

	assertSAML2Describe := func(details []sdk.SecurityIntegrationProperty, d saml2Details) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_X509_CERT", Type: "String", Value: cert, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_PROVIDER", Type: "String", Value: d.provider, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_ENABLE_SP_INITIATED", Type: "Boolean", Value: d.enableSPInitiated, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SP_INITIATED_LOGIN_PAGE_LABEL", Type: "String", Value: d.spInitiatedLoginPageLabel, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SSO_URL", Type: "String", Value: d.ssoURL, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_ISSUER", Type: "String", Value: d.issuer, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_REQUESTED_NAMEID_FORMAT", Type: "String", Value: d.requestedNameIDFormat, Default: string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress)})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_FORCE_AUTHN", Type: "Boolean", Value: d.forceAuthn, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_POST_LOGOUT_REDIRECT_URL", Type: "String", Value: d.postLogoutRedirectUrl, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SIGN_REQUEST", Type: "Boolean", Value: d.signrequest, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_DIGEST_METHODS_USED", Type: "String", Value: "http://www.w3.org/2001/04/xmlenc#sha256", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SIGNATURE_METHODS_USED", Type: "String", Value: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: d.comment, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SNOWFLAKE_ISSUER_URL", Type: "String", Value: d.snowflakeIssuerURL, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SNOWFLAKE_ACS_URL", Type: "String", Value: d.snowflakeAcsURL, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ALLOWED_USER_DOMAINS", Type: "List", Value: d.allowedUserDomains, Default: "[]"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ALLOWED_EMAIL_PATTERNS", Type: "List", Value: d.allowedEmailPatterns, Default: "[]"})
		// TODO(SNOW-1479617): assert SAML2_SNOWFLAKE_X509_CERT
	}

	t.Run("CreateApiAuthenticationWithClientCredentialsFlow", func(t *testing.T) {
		integration, id := createApiAuthClientCred(t, func(r *sdk.CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest) {
			r.WithComment("a").
				WithOauthAccessTokenValidity(31337).
				WithOauthRefreshTokenValidity(31337).
				WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}}).
				WithOauthClientAuthMethod(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost).
				WithOauthGrantClientCredentials(true).
				WithOauthTokenEndpoint("http://example.com")
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertApiAuth(details, apiAuthDetails{
			enabled:                   "false",
			oauthAccessTokenValidity:  "31337",
			oauthRefreshTokenValidity: "31337",
			oauthClientId:             "foo",
			oauthClientAuthMethod:     string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			oauthTokenEndpoint:        "http://example.com",
			oauthAllowedScopes:        "[foo]",
			oauthGrant:                "CLIENT_CREDENTIALS",
			authType:                  "OAUTH2",
			comment:                   "a",
		})

		assertSecurityIntegration(t, integration, id, "API_AUTHENTICATION", false, "a")
	})

	t.Run("CreateApiAuthenticationWithAuthorizationCodeGrantFlow", func(t *testing.T) {
		integration, id := createApiAuthCodeGrant(t, func(r *sdk.CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest) {
			r.WithComment("a").
				WithOauthAccessTokenValidity(31337).
				WithOauthAuthorizationEndpoint("http://example.com").
				WithOauthClientAuthMethod(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost).
				WithOauthGrantAuthorizationCode(true).
				WithOauthRefreshTokenValidity(31337).
				WithOauthTokenEndpoint("http://example.com")
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertApiAuth(details, apiAuthDetails{
			enabled:                    "false",
			oauthAccessTokenValidity:   "31337",
			oauthRefreshTokenValidity:  "31337",
			oauthClientId:              "foo",
			oauthClientAuthMethod:      string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			oauthAuthorizationEndpoint: "http://example.com",
			oauthTokenEndpoint:         "http://example.com",
			oauthGrant:                 "AUTHORIZATION_CODE",
			authType:                   "OAUTH2",
			comment:                    "a",
		})

		assertSecurityIntegration(t, integration, id, "API_AUTHENTICATION", false, "a")
	})

	t.Run("CreateApiAuthenticationWithJwtBearerFlow", func(t *testing.T) {
		// TODO [SNOW-1452191]: unskip
		t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")
		integration, id := createApiAuthJwtBearer(t, func(r *sdk.CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest) {
			r.WithComment("a").
				WithOauthAccessTokenValidity(31337).
				WithOauthAuthorizationEndpoint("http://example.com").
				WithOauthClientAuthMethod(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost).
				WithOauthGrantJwtBearer(true).
				WithOauthRefreshTokenValidity(31337).
				WithOauthTokenEndpoint("http://example.com")
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertApiAuth(details, apiAuthDetails{
			enabled:                    "false",
			oauthAccessTokenValidity:   "31337",
			oauthRefreshTokenValidity:  "31337",
			oauthClientId:              "foo",
			oauthClientAuthMethod:      string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			oauthAuthorizationEndpoint: "http://example.com",
			oauthTokenEndpoint:         "http://example.com",
			oauthGrant:                 "JWT_BEARER",
			authType:                   "OAUTH2",
			oauthAssertionIssuer:       "foo",
			comment:                    "a",
		})

		assertSecurityIntegration(t, integration, id, "API_AUTHENTICATION", false, "a")
	})

	t.Run("CreateExternalOauth with allowed list and jws keys url", func(t *testing.T) {
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)

		integration, id, _ := createExternalOauth(t, func(r *sdk.CreateExternalOauthSecurityIntegrationRequest) {
			r.WithExternalOauthAllowedRolesList(sdk.AllowedRolesListRequest{AllowedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
				WithExternalOauthJwsKeysUrl([]sdk.JwsKeysUrl{{JwsKeyUrl: "http://example.com"}})
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_JWS_KEYS_URL", Type: "Object", Value: "http://example.com", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_ALLOWED_ROLES_LIST", Type: "List", Value: role1.Name, Default: "[]"})

		assertSecurityIntegration(t, integration, id, "EXTERNAL_OAUTH - CUSTOM", false, "")
	})

	t.Run("CreateExternalOauth with other options", func(t *testing.T) {
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)

		integration, id, issuer := createExternalOauth(t, func(r *sdk.CreateExternalOauthSecurityIntegrationRequest) {
			r.WithExternalOauthBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
				WithExternalOauthRsaPublicKey(rsaKey).
				WithExternalOauthRsaPublicKey2(rsaKey).
				WithExternalOauthAudienceList(sdk.AudienceListRequest{AudienceList: []sdk.AudienceListItem{{Item: "foo"}}}).
				WithExternalOauthAnyRoleMode(sdk.ExternalOauthSecurityIntegrationAnyRoleModeEnable).
				WithExternalOauthScopeDelimiter(" ").
				WithExternalOauthScopeMappingAttribute("scp").
				WithComment("foo")
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertExternalOauth(details, externalOauthDetails{
			enabled:                                    "false",
			externalOauthIssuer:                        issuer,
			externalOauthAnyRoleMode:                   string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeEnable),
			externalOauthScopeMappingAttribute:         "scp",
			externalOauthRsaPublicKey:                  rsaKey,
			externalOauthRsaPublicKey2:                 rsaKey,
			externalOauthBlockedRolesList:              role1.Name,
			externalOauthAudienceList:                  "foo",
			externalOauthTokenUserMappingClaim:         "['foo']",
			externalOauthSnowflakeUserMappingAttribute: string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName),
			externalOauthScopeDelimiter:                " ",
			comment:                                    "foo",
		})

		assertSecurityIntegration(t, integration, id, "EXTERNAL_OAUTH - CUSTOM", false, "foo")
	})

	t.Run("CreateOauthPartner", func(t *testing.T) {
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)

		integration, id := createOauthPartner(t, func(r *sdk.CreateOauthForPartnerApplicationsSecurityIntegrationRequest) {
			r.WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
				WithComment("a").
				WithEnabled(true).
				WithOauthIssueRefreshTokens(true).
				WithOauthRefreshTokenValidity(12345).
				WithOauthUseSecondaryRoles(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertOauthPartner(details, oauthPartnerDetails{
			enabled:                 "true",
			oauthIssueRefreshTokens: "true",
			refreshTokenValidity:    "12345",
			useSecondaryRoles:       string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			blockedRolesList:        role1.Name,
			comment:                 "a",
		})

		assertSecurityIntegration(t, integration, id, "OAUTH - LOOKER", true, "a")
	})

	t.Run("CreateOauthCustom", func(t *testing.T) {
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
		t.Cleanup(networkPolicyCleanup)
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)
		role2, role2Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role2Cleanup)

		integration, id := createOauthCustom(t, func(r *sdk.CreateOauthForCustomClientsSecurityIntegrationRequest) {
			r.WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
				WithComment("a").
				WithEnabled(true).
				WithNetworkPolicy(networkPolicy.ID()).
				WithOauthAllowNonTlsRedirectUri(true).
				WithOauthClientRsaPublicKey(rsaKey).
				WithOauthClientRsaPublicKey2(rsaKey).
				WithOauthEnforcePkce(true).
				WithOauthIssueRefreshTokens(true).
				WithOauthRefreshTokenValidity(12345).
				WithOauthUseSecondaryRoles(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit).
				WithPreAuthorizedRolesList(sdk.PreAuthorizedRolesListRequest{PreAuthorizedRolesList: []sdk.AccountObjectIdentifier{role2.ID()}})
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertOauthCustom(details, oauthPartnerDetails{
			enabled:                 "true",
			oauthIssueRefreshTokens: "true",
			refreshTokenValidity:    "12345",
			useSecondaryRoles:       string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			preAuthorizedRolesList:  role2.Name,
			blockedRolesList:        role1.Name,
			networkPolicy:           networkPolicy.Name,
			comment:                 "a",
		}, "true", string(sdk.OauthSecurityIntegrationClientTypePublic), "true", rsaKeyHash, rsaKeyHash)

		assertSecurityIntegration(t, integration, id, "OAUTH - CUSTOM", true, "a")
	})

	// Prove that creating a security integration with a specified network policy id with lower case characters fails. This is an intended behavior in Snowflake.
	// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3229
	t.Run("CreateOauthCustom_issue3229_lowercase_network_policy_fails_with_double_quotes", func(t *testing.T) {
		networkPolicyId := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix("test")
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, sdk.NewCreateNetworkPolicyRequest(networkPolicyId))
		t.Cleanup(networkPolicyCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE SECURITY INTEGRATION %s TYPE = OAUTH OAUTH_CLIENT = CUSTOM OAUTH_CLIENT_TYPE = PUBLIC OAUTH_REDIRECT_URI = 'https://example.com' NETWORK_POLICY = %s`, id.FullyQualifiedName(), networkPolicy.ID().FullyQualifiedName()))
		require.ErrorContains(t, err, "object does not exist or not authorized")
	})

	// Prove that 3229 is fixed.
	t.Run("CreateOauthCustom_issue3229_fix", func(t *testing.T) {
		networkRule, networkRuleCleanup := testClientHelper().NetworkRule.CreateIP(t)
		t.Cleanup(networkRuleCleanup)

		networkPolicyId := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix("test")
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, sdk.NewCreateNetworkPolicyRequest(networkPolicyId).WithAllowedNetworkRuleList([]sdk.SchemaObjectIdentifier{networkRule.ID()}))
		t.Cleanup(networkPolicyCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateOauthForCustomClientsSecurityIntegrationRequest(id, sdk.OauthSecurityIntegrationClientTypePublic, "https://example.com").WithNetworkPolicy(networkPolicy.ID())
		err := client.SecurityIntegrations.CreateOauthForCustomClients(ctx, req)
		require.NoError(t, err)

		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
	})

	t.Run("CreateSaml2", func(t *testing.T) {
		_, id, issuer := createSAML2Integration(t, func(r *sdk.CreateSaml2SecurityIntegrationRequest) {
			r.WithAllowedEmailPatterns([]sdk.EmailPattern{{Pattern: "^(.+dev)@example.com$"}}).
				WithAllowedUserDomains([]sdk.UserDomain{{Domain: "example.com"}}).
				WithComment("a").
				WithSaml2EnableSpInitiated(true).
				WithSaml2ForceAuthn(true).
				WithSaml2PostLogoutRedirectUrl("http://example.com/logout").
				WithSaml2RequestedNameidFormat(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified).
				WithSaml2SignRequest(true).
				WithSaml2SnowflakeAcsUrl(acsURL).
				WithSaml2SnowflakeIssuerUrl(issuerURL).
				WithSaml2SpInitiatedLoginPageLabel("label").
				WithEnabled(true)
			// TODO(SNOW-1479617): fix after format clarification
			// WithSaml2SnowflakeX509Cert(sdk.Pointer(x509))
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSAML2Describe(details, saml2Details{
			provider:                  string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom),
			enableSPInitiated:         "true",
			spInitiatedLoginPageLabel: "label",
			ssoURL:                    "https://example.com",
			issuer:                    issuer,
			requestedNameIDFormat:     string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified),
			forceAuthn:                "true",
			postLogoutRedirectUrl:     "http://example.com/logout",
			signrequest:               "true",
			comment:                   "a",
			snowflakeIssuerURL:        issuerURL,
			snowflakeAcsURL:           acsURL,
			allowedUserDomains:        "[example.com]",
			allowedEmailPatterns:      "[^(.+dev)@example.com$]",
		})

		si, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSecurityIntegration(t, si, id, "SAML2", true, "a")
	})

	t.Run("CreateScim", func(t *testing.T) {
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
		t.Cleanup(networkPolicyCleanup)

		_, id := createSCIMIntegration(t, func(r *sdk.CreateScimSecurityIntegrationRequest) {
			r.WithComment("a").
				WithNetworkPolicy(networkPolicy.ID()).
				WithSyncPassword(false)
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSCIMDescribe(details, "true", networkPolicy.Name, "GENERIC_SCIM_PROVISIONER", "false", "a")

		si, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSecurityIntegration(t, si, id, "SCIM - GENERIC", true, "a")
	})

	t.Run("AlterApiAuthenticationWithClientCredentialsFlow", func(t *testing.T) {
		_, id := createApiAuthClientCred(t, nil)
		setRequest := sdk.NewAlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id).
			WithSet(
				*sdk.NewApiAuthenticationWithClientCredentialsFlowIntegrationSetRequest().
					WithComment("foo").
					WithEnabled(true).
					WithOauthAccessTokenValidity(31337).
					WithOauthRefreshTokenValidity(31337).
					WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}}).
					WithOauthClientAuthMethod(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost).
					WithOauthClientId("foo").
					WithOauthClientSecret("foo").
					WithOauthGrantClientCredentials(true).
					WithOauthTokenEndpoint("http://example.com"),
			)
		err := client.SecurityIntegrations.AlterApiAuthenticationWithClientCredentialsFlow(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertApiAuth(details, apiAuthDetails{
			enabled:                   "true",
			oauthAccessTokenValidity:  "31337",
			oauthRefreshTokenValidity: "31337",
			oauthClientId:             "foo",
			oauthClientAuthMethod:     string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			oauthTokenEndpoint:        "http://example.com",
			oauthAllowedScopes:        "[foo]",
			oauthGrant:                "CLIENT_CREDENTIALS",
			authType:                  "OAUTH2",
			comment:                   "foo",
		})

		unsetRequest := sdk.NewAlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id).
			WithUnset(
				*sdk.NewApiAuthenticationWithClientCredentialsFlowIntegrationUnsetRequest().
					WithEnabled(true).
					WithComment(true),
			)
		err = client.SecurityIntegrations.AlterApiAuthenticationWithClientCredentialsFlow(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("AlterApiAuthenticationWithAuthorizationCodeGrantFlow", func(t *testing.T) {
		_, id := createApiAuthCodeGrant(t, nil)
		setRequest := sdk.NewAlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id).
			WithSet(
				*sdk.NewApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSetRequest().
					WithComment("foo").
					WithEnabled(true).
					WithOauthAccessTokenValidity(31337).
					WithOauthRefreshTokenValidity(31337).
					WithOauthClientAuthMethod(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost).
					WithOauthClientId("foo").
					WithOauthClientSecret("foo").
					WithOauthGrantAuthorizationCode(true).
					WithOauthAuthorizationEndpoint("http://example.com").
					WithOauthTokenEndpoint("http://example.com"),
			)
		err := client.SecurityIntegrations.AlterApiAuthenticationWithAuthorizationCodeGrantFlow(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertApiAuth(details, apiAuthDetails{
			enabled:                    "true",
			oauthAccessTokenValidity:   "31337",
			oauthRefreshTokenValidity:  "31337",
			oauthClientId:              "foo",
			oauthClientAuthMethod:      string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			oauthAuthorizationEndpoint: "http://example.com",
			oauthTokenEndpoint:         "http://example.com",
			oauthGrant:                 "AUTHORIZATION_CODE",
			authType:                   "OAUTH2",
			comment:                    "foo",
		})

		unsetRequest := sdk.NewAlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id).
			WithUnset(
				*sdk.NewApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnsetRequest().
					WithEnabled(true).
					WithComment(true),
			)
		err = client.SecurityIntegrations.AlterApiAuthenticationWithAuthorizationCodeGrantFlow(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("AlterApiAuthenticationWithJwtBearerFlow", func(t *testing.T) {
		// TODO [SNOW-1452191]: unskip
		t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")

		_, id := createApiAuthJwtBearer(t, nil)
		setRequest := sdk.NewAlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest(id).
			WithSet(
				*sdk.NewApiAuthenticationWithJwtBearerFlowIntegrationSetRequest().
					WithComment("a").
					WithEnabled(true).
					WithOauthAccessTokenValidity(31337).
					WithOauthAuthorizationEndpoint("http://example.com").
					WithOauthClientAuthMethod(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost).
					WithOauthClientId("foo").
					WithOauthClientSecret("foo").
					WithOauthGrantJwtBearer(true).
					WithOauthRefreshTokenValidity(31337).
					WithOauthTokenEndpoint("http://example.com"),
			)
		err := client.SecurityIntegrations.AlterApiAuthenticationWithJwtBearerFlow(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertApiAuth(details, apiAuthDetails{
			enabled:                    "true",
			oauthAccessTokenValidity:   "31337",
			oauthRefreshTokenValidity:  "31337",
			oauthClientId:              "foo",
			oauthClientAuthMethod:      string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost),
			oauthAuthorizationEndpoint: "http://example.com",
			oauthTokenEndpoint:         "http://example.com",
			oauthGrant:                 "JWT_BEARER",
			authType:                   "OAUTH2",
			comment:                    "foo",
		})

		unsetRequest := sdk.NewAlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest(id).
			WithUnset(
				*sdk.NewApiAuthenticationWithJwtBearerFlowIntegrationUnsetRequest().
					WithEnabled(true).
					WithComment(true),
			)
		err = client.SecurityIntegrations.AlterApiAuthenticationWithJwtBearerFlow(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("AlterExternalOauth with other options", func(t *testing.T) {
		_, id, _ := createExternalOauth(t, func(r *sdk.CreateExternalOauthSecurityIntegrationRequest) {
			r.WithExternalOauthRsaPublicKey(rsaKey).
				WithExternalOauthRsaPublicKey2(rsaKey)
		})
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)
		newIssuer := testClientHelper().Ids.Alpha()
		setRequest := sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).
			WithSet(
				*sdk.NewExternalOauthIntegrationSetRequest().
					WithEnabled(true).
					WithExternalOauthIssuer(newIssuer).
					WithExternalOauthTokenUserMappingClaim([]sdk.TokenUserMappingClaim{{Claim: "bar"}}).
					WithExternalOauthSnowflakeUserMappingAttribute(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress).
					WithExternalOauthBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
					WithExternalOauthRsaPublicKey(rsaKey).
					WithExternalOauthRsaPublicKey2(rsaKey).
					WithExternalOauthAudienceList(sdk.AudienceListRequest{AudienceList: []sdk.AudienceListItem{{Item: "foo"}}}).
					WithExternalOauthAnyRoleMode(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable).
					WithExternalOauthScopeDelimiter(" ").
					WithComment(sdk.StringAllowEmpty{Value: "foo"}),
			)
		err := client.SecurityIntegrations.AlterExternalOauth(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertExternalOauth(details, externalOauthDetails{
			enabled:                                    "true",
			externalOauthIssuer:                        newIssuer,
			externalOauthAnyRoleMode:                   string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable),
			externalOauthScopeMappingAttribute:         "scp",
			externalOauthRsaPublicKey:                  rsaKey,
			externalOauthRsaPublicKey2:                 rsaKey,
			externalOauthBlockedRolesList:              role1.Name,
			externalOauthAudienceList:                  "foo",
			externalOauthTokenUserMappingClaim:         "['bar']",
			externalOauthSnowflakeUserMappingAttribute: string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress),
			externalOauthScopeDelimiter:                " ",
			comment:                                    "foo",
		})

		unsetRequest := sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).
			WithUnset(
				*sdk.NewExternalOauthIntegrationUnsetRequest().
					WithEnabled(true).
					WithExternalOauthAudienceList(true),
			)
		err = client.SecurityIntegrations.AlterExternalOauth(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_AUDIENCE_LIST", Type: "List", Value: "", Default: "[]"})
	})

	t.Run("AlterOauthPartner", func(t *testing.T) {
		_, id := createOauthPartner(t, func(r *sdk.CreateOauthForPartnerApplicationsSecurityIntegrationRequest) {
			r.WithOauthRedirectUri("http://example.com")
		})
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)

		setRequest := sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).
			WithSet(
				*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().
					WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
					WithComment("a").
					WithEnabled(true).
					WithOauthIssueRefreshTokens(true).
					WithOauthRedirectUri("http://example2.com").
					WithOauthRefreshTokenValidity(22222).
					WithOauthUseSecondaryRoles(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			)
		err := client.SecurityIntegrations.AlterOauthForPartnerApplications(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertOauthPartner(details, oauthPartnerDetails{
			enabled:                 "true",
			oauthIssueRefreshTokens: "true",
			refreshTokenValidity:    "22222",
			useSecondaryRoles:       string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			preAuthorizedRolesList:  "",
			blockedRolesList:        "ACCOUNTADMIN,SECURITYADMIN",
			networkPolicy:           "",
			comment:                 "a",
		})

		unsetRequest := sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).
			WithUnset(
				*sdk.NewOauthForPartnerApplicationsIntegrationUnsetRequest().
					WithEnabled(true).
					WithOauthUseSecondaryRoles(true),
			)
		err = client.SecurityIntegrations.AlterOauthForPartnerApplications(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_USE_SECONDARY_ROLES", Type: "String", Value: "NONE", Default: "NONE"})
	})

	t.Run("AlterOauthCustom", func(t *testing.T) {
		_, id := createOauthCustom(t, nil)

		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
		t.Cleanup(networkPolicyCleanup)
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)
		role2, role2Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role2Cleanup)

		setRequest := sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).
			WithSet(
				*sdk.NewOauthForCustomClientsIntegrationSetRequest().
					WithEnabled(true).
					WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
					WithComment("a").
					WithNetworkPolicy(networkPolicy.ID()).
					WithOauthAllowNonTlsRedirectUri(true).
					WithOauthClientRsaPublicKey(rsaKey).
					WithOauthClientRsaPublicKey2(rsaKey).
					WithOauthEnforcePkce(true).
					WithOauthIssueRefreshTokens(true).
					WithOauthRedirectUri("http://example2.com").
					WithOauthRefreshTokenValidity(22222).
					WithOauthUseSecondaryRoles(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit).
					WithPreAuthorizedRolesList(sdk.PreAuthorizedRolesListRequest{PreAuthorizedRolesList: []sdk.AccountObjectIdentifier{role2.ID()}}),
			)
		err := client.SecurityIntegrations.AlterOauthForCustomClients(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertOauthCustom(details, oauthPartnerDetails{
			enabled:                 "true",
			oauthIssueRefreshTokens: "true",
			refreshTokenValidity:    "22222",
			useSecondaryRoles:       string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			preAuthorizedRolesList:  role2.Name,
			blockedRolesList:        role1.Name,
			networkPolicy:           networkPolicy.Name,
			comment:                 "a",
		}, "true", string(sdk.OauthSecurityIntegrationClientTypePublic), "true", rsaKeyHash, rsaKeyHash)

		unsetRequest := sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).
			WithUnset(
				*sdk.NewOauthForCustomClientsIntegrationUnsetRequest().
					WithEnabled(true).
					WithOauthUseSecondaryRoles(true).
					WithNetworkPolicy(true).
					WithOauthClientRsaPublicKey(true).
					WithOauthClientRsaPublicKey2(true),
			)
		err = client.SecurityIntegrations.AlterOauthForCustomClients(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_USE_SECONDARY_ROLES", Type: "String", Value: "NONE", Default: "NONE"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_RSA_PUBLIC_KEY_FP", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP", Type: "String", Value: "", Default: ""})
	})

	// Prove that altering a security integration with a specified network policy id with lower case characters fails. This is a bug in Snowflake.
	// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3229
	t.Run("AlterOauthCustom_issue3229_lowercase_network_policy_fails_with_double_quotes", func(t *testing.T) {
		networkPolicyId := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix("test")
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, sdk.NewCreateNetworkPolicyRequest(networkPolicyId))
		t.Cleanup(networkPolicyCleanup)

		_, id := createOauthCustom(t, nil)

		_, err := client.ExecForTests(ctx, fmt.Sprintf(`ALTER SECURITY INTEGRATION %s SET NETWORK_POLICY = %s`, id.FullyQualifiedName(), networkPolicy.ID().FullyQualifiedName()))
		require.ErrorContains(t, err, "object does not exist or not authorized")
	})

	// Prove that 3229 is fixed.
	t.Run("AlterOauthCustom_issue3229_fix", func(t *testing.T) {
		networkRule, networkRuleCleanup := testClientHelper().NetworkRule.CreateIP(t)
		t.Cleanup(networkRuleCleanup)

		networkPolicyId := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix("test")
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, sdk.NewCreateNetworkPolicyRequest(networkPolicyId).WithAllowedNetworkRuleList([]sdk.SchemaObjectIdentifier{networkRule.ID()}))
		t.Cleanup(networkPolicyCleanup)

		_, id := createOauthCustom(t, nil)
		setRequest := sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).
			WithSet(
				*sdk.NewOauthForCustomClientsIntegrationSetRequest().
					WithNetworkPolicy(networkPolicy.ID()),
			)
		err := client.SecurityIntegrations.AlterOauthForCustomClients(ctx, setRequest)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SecurityIntegration.DropSecurityIntegrationFunc(t, id))
	})

	t.Run("AlterSAML2Integration", func(t *testing.T) {
		_, id, issuer := createSAML2Integration(t, nil)

		setRequest := sdk.NewAlterSaml2SecurityIntegrationRequest(id).
			WithSet(
				*sdk.NewSaml2IntegrationSetRequest().
					WithEnabled(true).
					WithSaml2Issuer(issuer).
					WithSaml2SsoUrl("http://example.com").
					WithSaml2Provider("OKTA").
					WithSaml2X509Cert(cert).
					WithComment("a").
					WithSaml2EnableSpInitiated(true).
					WithSaml2ForceAuthn(true).
					WithSaml2PostLogoutRedirectUrl("http://example.com/logout").
					WithSaml2RequestedNameidFormat(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified).
					WithSaml2SignRequest(true).
					WithSaml2SnowflakeAcsUrl(acsURL).
					WithSaml2SnowflakeIssuerUrl(issuerURL).
					WithSaml2SpInitiatedLoginPageLabel("label").
					WithAllowedEmailPatterns([]sdk.EmailPattern{{Pattern: "^(.+dev)@example.com$"}}).
					WithAllowedUserDomains([]sdk.UserDomain{{Domain: "example.com"}}),
				// TODO(SNOW-1479617): fix after format clarification
				// WithSaml2SnowflakeX509Cert(sdk.Pointer(cert)).
			)
		err := client.SecurityIntegrations.AlterSaml2(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSAML2Describe(details, saml2Details{
			provider:                  string(sdk.Saml2SecurityIntegrationSaml2ProviderOkta),
			enableSPInitiated:         "true",
			spInitiatedLoginPageLabel: "label",
			ssoURL:                    "http://example.com",
			issuer:                    issuer,
			requestedNameIDFormat:     string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified),
			forceAuthn:                "true",
			postLogoutRedirectUrl:     "http://example.com/logout",
			signrequest:               "true",
			comment:                   "a",
			snowflakeIssuerURL:        issuerURL,
			snowflakeAcsURL:           acsURL,
			allowedUserDomains:        "[example.com]",
			allowedEmailPatterns:      "[^(.+dev)@example.com$]",
		})

		unsetRequest := sdk.NewAlterSaml2SecurityIntegrationRequest(id).
			WithUnset(
				*sdk.NewSaml2IntegrationUnsetRequest().
					WithSaml2ForceAuthn(true).
					WithSaml2RequestedNameidFormat(true).
					WithSaml2PostLogoutRedirectUrl(true).
					WithComment(true),
			)
		err = client.SecurityIntegrations.AlterSaml2(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_FORCE_AUTHN", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_REQUESTED_NAMEID_FORMAT", Type: "String", Value: string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress), Default: string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress)})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_POST_LOGOUT_REDIRECT_URL", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT_changed", Type: "String", Value: "", Default: ""})
	})

	t.Run("AlterSAML2Integration - REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY", func(t *testing.T) {
		_, id, _ := createSAML2Integration(t, nil)

		setRequest := sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithRefreshSaml2SnowflakePrivateKey(true)
		err := client.SecurityIntegrations.AlterSaml2(ctx, setRequest)
		require.NoError(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		_, id := createSCIMIntegration(t, nil)

		si, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NotNil(t, si)
		require.NoError(t, err)

		err = client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id))
		require.NoError(t, err)

		si, err = client.SecurityIntegrations.ShowByID(ctx, id)
		require.Nil(t, si)
		require.Error(t, err)
	})

	t.Run("Drop non-existing", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")

		err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Describe", func(t *testing.T) {
		_, id := createSCIMIntegration(t, nil)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSCIMDescribe(details, "true", "", "GENERIC_SCIM_PROVISIONER", "false", "")
	})

	t.Run("ShowByID", func(t *testing.T) {
		_, id := createSCIMIntegration(t, nil)

		si, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSecurityIntegration(t, si, id, "SCIM - GENERIC", true, "")
	})

	t.Run("Show ExternalOauth", func(t *testing.T) {
		si1, id1, _ := createExternalOauth(t, func(r *sdk.CreateExternalOauthSecurityIntegrationRequest) {
			r.WithExternalOauthJwsKeysUrl([]sdk.JwsKeysUrl{{JwsKeyUrl: "http://example.com"}})
		})
		si2, _, _ := createExternalOauth(t, func(r *sdk.CreateExternalOauthSecurityIntegrationRequest) {
			r.WithExternalOauthJwsKeysUrl([]sdk.JwsKeysUrl{{JwsKeyUrl: "http://example2.com"}})
		})

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})

	t.Run("Show OauthPartner", func(t *testing.T) {
		si1, id1 := createOauthPartner(t, nil)
		// more than one oauth partner integration is not allowed, create a custom one
		si2, _ := createOauthCustom(t, nil)

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})

	t.Run("Show OauthCustom", func(t *testing.T) {
		si1, id1 := createOauthCustom(t, nil)
		si2, _ := createOauthCustom(t, nil)

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})

	t.Run("Show SAML2", func(t *testing.T) {
		si1, id1, _ := createSAML2Integration(t, nil)
		si2, _, _ := createSAML2Integration(t, nil)

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})

	t.Run("Show SCIM", func(t *testing.T) {
		si1, id1 := createSCIMIntegration(t, nil)
		si2, _ := createSCIMIntegration(t, nil)

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})
}
