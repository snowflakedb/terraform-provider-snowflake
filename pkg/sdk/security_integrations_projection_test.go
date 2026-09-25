package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecurityIntegrationDescribeProjections(t *testing.T) {
	enabled := SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"}
	networkPolicy := SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: "NP", Default: ""}
	runAsRole := SecurityIntegrationProperty{Name: "RUN_AS_ROLE", Type: "String", Value: "GENERIC_SCIM_PROVISIONER", Default: ""}
	syncPassword := SecurityIntegrationProperty{Name: "SYNC_PASSWORD", Type: "Boolean", Value: "false", Default: "true"}
	comment := SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: "scim comment", Default: ""}
	unknown := SecurityIntegrationProperty{Name: "UNEXPECTED", Type: "String", Value: "x", Default: ""}

	scim := AsScim([]SecurityIntegrationProperty{enabled, networkPolicy, runAsRole, syncPassword, comment, unknown})
	require.Equal(t, &enabled, scim.Enabled)
	require.Equal(t, &networkPolicy, scim.NetworkPolicy)
	require.Equal(t, &runAsRole, scim.RunAsRole)
	require.Equal(t, &syncPassword, scim.SyncPassword)
	require.Equal(t, &comment, scim.Comment)
	require.Nil(t, AsScim(nil))

	issuer := SecurityIntegrationProperty{Name: "SAML2_ISSUER", Type: "String", Value: "https://idp.example", Default: ""}
	ssoUrl := SecurityIntegrationProperty{Name: "SAML2_SSO_URL", Type: "String", Value: "https://idp.example/sso", Default: ""}
	domains := SecurityIntegrationProperty{Name: "ALLOWED_USER_DOMAINS", Type: "List", Value: "example.com", Default: ""}
	saml2 := AsSaml2([]SecurityIntegrationProperty{issuer, ssoUrl, domains, comment, unknown})
	require.Equal(t, &issuer, saml2.Saml2Issuer)
	require.Equal(t, &ssoUrl, saml2.Saml2SsoUrl)
	require.Equal(t, &domains, saml2.AllowedUserDomains)
	require.Equal(t, &comment, saml2.Comment)
	require.Nil(t, AsSaml2(nil))

	clientType := SecurityIntegrationProperty{Name: "OAUTH_CLIENT_TYPE", Type: "String", Value: "LOOKER", Default: ""}
	key2Fp := SecurityIntegrationProperty{Name: "OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP", Type: "String", Value: "fp2", Default: ""}
	partner := AsOauthPartner([]SecurityIntegrationProperty{clientType, enabled, key2Fp, comment, unknown})
	require.Equal(t, &clientType, partner.OauthClientType)
	require.Equal(t, &enabled, partner.Enabled)
	require.Equal(t, &key2Fp, partner.OauthClientRsaPublicKey2Fp)
	require.Equal(t, &comment, partner.Comment)
	require.Nil(t, AsOauthPartner(nil))

	custom := AsOauthCustom([]SecurityIntegrationProperty{clientType, enabled, key2Fp, comment, unknown})
	require.Equal(t, &clientType, custom.OauthClientType)
	require.Equal(t, &enabled, custom.Enabled)
	require.Equal(t, &key2Fp, custom.OauthClientRsaPublicKey2Fp)
	require.Equal(t, &comment, custom.Comment)
	require.Nil(t, AsOauthCustom(nil))

	issuerExt := SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_ISSUER", Type: "String", Value: "https://idp.example", Default: ""}
	key2 := SecurityIntegrationProperty{Name: "EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2", Type: "String", Value: "pk2", Default: ""}
	external := AsExternalOauth([]SecurityIntegrationProperty{issuerExt, enabled, key2, comment, unknown})
	require.Equal(t, &issuerExt, external.ExternalOauthIssuer)
	require.Equal(t, &enabled, external.Enabled)
	require.Equal(t, &key2, external.ExternalOauthRsaPublicKey2)
	require.Equal(t, &comment, external.Comment)
	require.Nil(t, AsExternalOauth(nil))

	authType := SecurityIntegrationProperty{Name: "AUTH_TYPE", Type: "String", Value: "OAUTH2", Default: ""}
	grant := SecurityIntegrationProperty{Name: "OAUTH_GRANT", Type: "String", Value: "CLIENT_CREDENTIALS", Default: ""}
	apiAuth := AsApiAuth([]SecurityIntegrationProperty{authType, grant, enabled, comment, unknown})
	require.Equal(t, &authType, apiAuth.AuthType)
	require.Equal(t, &grant, apiAuth.OauthGrant)
	require.Equal(t, &enabled, apiAuth.Enabled)
	require.Equal(t, &comment, apiAuth.Comment)
	require.Nil(t, AsApiAuth(nil))
}
