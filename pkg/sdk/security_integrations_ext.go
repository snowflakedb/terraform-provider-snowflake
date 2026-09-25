package sdk

import (
	"fmt"
	"strings"
)

func (opts *CreateOauthForPartnerApplicationsSecurityIntegrationOptions) additionalValidations() error {
	if opts.OauthClient == OauthSecurityIntegrationClientOptionLooker && opts.OauthRedirectUri == nil {
		return NewError("OauthRedirectUri is required when OauthClient is LOOKER")
	}
	return nil
}

func (opts *CreateScimSecurityIntegrationOptions) additionalValidations() error {
	if opts.ScimClient == ScimSecurityIntegrationScimClientOptionAzure && opts.SyncPassword != nil {
		return NewError("SyncPassword is not supported for Azure scim client")
	}
	return nil
}

func (r *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateExternalOauthSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateOauthForCustomClientsSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateSaml2SecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateScimSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (s SecurityIntegrationProperty) GetName() string {
	return s.Name
}

func (s SecurityIntegrationProperty) GetDefault() string {
	return s.Default
}

func (s *SecurityIntegration) SubType() (string, error) {
	typeParts := strings.Split(s.IntegrationType, "-")
	if len(typeParts) < 2 {
		return "", fmt.Errorf("expected \"<type> - <subtype>\", got: %s", s.IntegrationType)
	}
	return strings.TrimSpace(typeParts[1]), nil
}

// AsScim projects DESCRIBE output onto SCIM resource columns.
func AsScim(properties []SecurityIntegrationProperty) *ScimSecurityIntegrationDetails {
	if properties == nil {
		return nil
	}
	details := &ScimSecurityIntegrationDetails{}
	for _, property := range properties {
		p := property
		switch p.Name {
		case "ENABLED":
			details.Enabled = &p
		case "NETWORK_POLICY":
			details.NetworkPolicy = &p
		case "RUN_AS_ROLE":
			details.RunAsRole = &p
		case "SYNC_PASSWORD":
			details.SyncPassword = &p
		case "COMMENT":
			details.Comment = &p
		}
	}
	return details
}

// AsSaml2 projects DESCRIBE output onto SAML2 resource columns.
func AsSaml2(properties []SecurityIntegrationProperty) *Saml2SecurityIntegrationDetails {
	if properties == nil {
		return nil
	}
	details := &Saml2SecurityIntegrationDetails{}
	for _, property := range properties {
		p := property
		switch p.Name {
		case "SAML2_ISSUER":
			details.Saml2Issuer = &p
		case "SAML2_SSO_URL":
			details.Saml2SsoUrl = &p
		case "SAML2_PROVIDER":
			details.Saml2Provider = &p
		case "SAML2_SP_INITIATED_LOGIN_PAGE_LABEL":
			details.Saml2SpInitiatedLoginPageLabel = &p
		case "SAML2_ENABLE_SP_INITIATED":
			details.Saml2EnableSpInitiated = &p
		case "SAML2_SIGN_REQUEST":
			details.Saml2SignRequest = &p
		case "SAML2_REQUESTED_NAMEID_FORMAT":
			details.Saml2RequestedNameidFormat = &p
		case "SAML2_POST_LOGOUT_REDIRECT_URL":
			details.Saml2PostLogoutRedirectUrl = &p
		case "SAML2_FORCE_AUTHN":
			details.Saml2ForceAuthn = &p
		case "SAML2_SNOWFLAKE_ISSUER_URL":
			details.Saml2SnowflakeIssuerUrl = &p
		case "SAML2_SNOWFLAKE_ACS_URL":
			details.Saml2SnowflakeAcsUrl = &p
		case "SAML2_SNOWFLAKE_METADATA":
			details.Saml2SnowflakeMetadata = &p
		case "SAML2_DIGEST_METHODS_USED":
			details.Saml2DigestMethodsUsed = &p
		case "SAML2_SIGNATURE_METHODS_USED":
			details.Saml2SignatureMethodsUsed = &p
		case "ALLOWED_USER_DOMAINS":
			details.AllowedUserDomains = &p
		case "ALLOWED_EMAIL_PATTERNS":
			details.AllowedEmailPatterns = &p
		case "COMMENT":
			details.Comment = &p
		}
	}
	return details
}

// AsOauthPartner projects DESCRIBE output onto OAuth partner-application resource columns.
func AsOauthPartner(properties []SecurityIntegrationProperty) *OauthIntegrationForPartnerApplicationsDetails {
	if properties == nil {
		return nil
	}
	details := &OauthIntegrationForPartnerApplicationsDetails{}
	for _, property := range properties {
		p := property
		switch p.Name {
		case "OAUTH_CLIENT_TYPE":
			details.OauthClientType = &p
		case "ENABLED":
			details.Enabled = &p
		case "OAUTH_ALLOW_NON_TLS_REDIRECT_URI":
			details.OauthAllowNonTlsRedirectUri = &p
		case "OAUTH_ENFORCE_PKCE":
			details.OauthEnforcePkce = &p
		case "OAUTH_USE_SECONDARY_ROLES":
			details.OauthUseSecondaryRoles = &p
		case "PRE_AUTHORIZED_ROLES_LIST":
			details.PreAuthorizedRolesList = &p
		case "ALLOWED_ROLES_LIST":
			details.AllowedRolesList = &p
		case "BLOCKED_ROLES_LIST":
			details.BlockedRolesList = &p
		case "OAUTH_ISSUE_REFRESH_TOKENS":
			details.OauthIssueRefreshTokens = &p
		case "OAUTH_REFRESH_TOKEN_VALIDITY":
			details.OauthRefreshTokenValidity = &p
		case "NETWORK_POLICY":
			details.NetworkPolicy = &p
		case "OAUTH_CLIENT_RSA_PUBLIC_KEY_FP":
			details.OauthClientRsaPublicKeyFp = &p
		case "OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP":
			details.OauthClientRsaPublicKey2Fp = &p
		case "COMMENT":
			details.Comment = &p
		case "OAUTH_AUTHORIZATION_ENDPOINT":
			details.OauthAuthorizationEndpoint = &p
		case "OAUTH_TOKEN_ENDPOINT":
			details.OauthTokenEndpoint = &p
		case "OAUTH_ALLOWED_AUTHORIZATION_ENDPOINTS":
			details.OauthAllowedAuthorizationEndpoints = &p
		case "OAUTH_ALLOWED_TOKEN_ENDPOINTS":
			details.OauthAllowedTokenEndpoints = &p
		}
	}
	return details
}

// AsOauthCustom projects DESCRIBE output onto OAuth custom-client resource columns.
func AsOauthCustom(properties []SecurityIntegrationProperty) *OauthIntegrationForCustomClientsDetails {
	if properties == nil {
		return nil
	}
	details := &OauthIntegrationForCustomClientsDetails{}
	for _, property := range properties {
		p := property
		switch p.Name {
		case "OAUTH_CLIENT_TYPE":
			details.OauthClientType = &p
		case "ENABLED":
			details.Enabled = &p
		case "OAUTH_ALLOW_NON_TLS_REDIRECT_URI":
			details.OauthAllowNonTlsRedirectUri = &p
		case "OAUTH_ENFORCE_PKCE":
			details.OauthEnforcePkce = &p
		case "OAUTH_USE_SECONDARY_ROLES":
			details.OauthUseSecondaryRoles = &p
		case "PRE_AUTHORIZED_ROLES_LIST":
			details.PreAuthorizedRolesList = &p
		case "ALLOWED_ROLES_LIST":
			details.AllowedRolesList = &p
		case "BLOCKED_ROLES_LIST":
			details.BlockedRolesList = &p
		case "OAUTH_ISSUE_REFRESH_TOKENS":
			details.OauthIssueRefreshTokens = &p
		case "OAUTH_REFRESH_TOKEN_VALIDITY":
			details.OauthRefreshTokenValidity = &p
		case "NETWORK_POLICY":
			details.NetworkPolicy = &p
		case "OAUTH_CLIENT_RSA_PUBLIC_KEY_FP":
			details.OauthClientRsaPublicKeyFp = &p
		case "OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP":
			details.OauthClientRsaPublicKey2Fp = &p
		case "COMMENT":
			details.Comment = &p
		case "OAUTH_AUTHORIZATION_ENDPOINT":
			details.OauthAuthorizationEndpoint = &p
		case "OAUTH_TOKEN_ENDPOINT":
			details.OauthTokenEndpoint = &p
		case "OAUTH_ALLOWED_AUTHORIZATION_ENDPOINTS":
			details.OauthAllowedAuthorizationEndpoints = &p
		case "OAUTH_ALLOWED_TOKEN_ENDPOINTS":
			details.OauthAllowedTokenEndpoints = &p
		}
	}
	return details
}

// AsExternalOauth projects DESCRIBE output onto external OAuth resource columns.
func AsExternalOauth(properties []SecurityIntegrationProperty) *ExternalOauthSecurityIntegrationDetails {
	if properties == nil {
		return nil
	}
	details := &ExternalOauthSecurityIntegrationDetails{}
	for _, property := range properties {
		p := property
		switch p.Name {
		case "ENABLED":
			details.Enabled = &p
		case "EXTERNAL_OAUTH_ISSUER":
			details.ExternalOauthIssuer = &p
		case "EXTERNAL_OAUTH_JWS_KEYS_URL":
			details.ExternalOauthJwsKeysUrl = &p
		case "EXTERNAL_OAUTH_ANY_ROLE_MODE":
			details.ExternalOauthAnyRoleMode = &p
		case "EXTERNAL_OAUTH_RSA_PUBLIC_KEY":
			details.ExternalOauthRsaPublicKey = &p
		case "EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2":
			details.ExternalOauthRsaPublicKey2 = &p
		case "EXTERNAL_OAUTH_BLOCKED_ROLES_LIST":
			details.ExternalOauthBlockedRolesList = &p
		case "EXTERNAL_OAUTH_ALLOWED_ROLES_LIST":
			details.ExternalOauthAllowedRolesList = &p
		case "EXTERNAL_OAUTH_AUDIENCE_LIST":
			details.ExternalOauthAudienceList = &p
		case "EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM":
			details.ExternalOauthTokenUserMappingClaim = &p
		case "EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE":
			details.ExternalOauthSnowflakeUserMappingAttribute = &p
		case "EXTERNAL_OAUTH_SCOPE_DELIMITER":
			details.ExternalOauthScopeDelimiter = &p
		case "COMMENT":
			details.Comment = &p
		}
	}
	return details
}

// AsApiAuth projects DESCRIBE output onto API authentication resource columns.
func AsApiAuth(properties []SecurityIntegrationProperty) *ApiAuthenticationSecurityIntegrationDetails {
	if properties == nil {
		return nil
	}
	details := &ApiAuthenticationSecurityIntegrationDetails{}
	for _, property := range properties {
		p := property
		switch p.Name {
		case "ENABLED":
			details.Enabled = &p
		case "OAUTH_ACCESS_TOKEN_VALIDITY":
			details.OauthAccessTokenValidity = &p
		case "OAUTH_REFRESH_TOKEN_VALIDITY":
			details.OauthRefreshTokenValidity = &p
		case "OAUTH_CLIENT_AUTH_METHOD":
			details.OauthClientAuthMethod = &p
		case "OAUTH_AUTHORIZATION_ENDPOINT":
			details.OauthAuthorizationEndpoint = &p
		case "OAUTH_TOKEN_ENDPOINT":
			details.OauthTokenEndpoint = &p
		case "OAUTH_ALLOWED_SCOPES":
			details.OauthAllowedScopes = &p
		case "OAUTH_GRANT":
			details.OauthGrant = &p
		case "PARENT_INTEGRATION":
			details.ParentIntegration = &p
		case "AUTH_TYPE":
			details.AuthType = &p
		case "COMMENT":
			details.Comment = &p
		}
	}
	return details
}
