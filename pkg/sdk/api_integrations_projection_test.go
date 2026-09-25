package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiIntegrationDescribeProjections(t *testing.T) {
	id := NewAccountObjectIdentifier("API_INT")
	git := &ApiIntegrationGitHttpsApiDetails{ //nolint:gosec // test credentials
		Id:                           id,
		Enabled:                      true,
		ApiProvider:                  "GIT_HTTPS_API",
		AllowedAuthenticationSecrets: "ALL",
		UserAuthType:                 "OAUTH2",
		OauthGrant:                   "AUTHORIZATION_CODE",
		OauthClientId:                "client-id",
		OauthClientAuthMethod:        "CLIENT_SECRET_POST",
		OauthTokenEndpoint:           "https://auth.example.com/token",
		OauthAuthorizationEndpoint:   "https://auth.example.com/authorize",
		OauthAccessTokenValidity:     3600,
		OauthRefreshTokenValidity:    86400,
		OauthAllowedScopes:           []string{"repo"},
		OauthUsername:                "user",
		OauthAssertionIssuer:         "issuer",
		OauthResourceUrl:             "https://resource.example.com",
		UsePrivatelinkEndpoint:       true,
		TlsTrustedCertificates:       []string{"DB.SCH.CERT"},
		AllowedPrefixes:              []string{"https://github.com/org/"},
		BlockedPrefixes:              []string{"https://github.com/blocked/"},
		Comment:                      "git comment",
	}

	token := git.AsToken()
	require.Equal(t, git.Enabled, token.Enabled)
	require.Equal(t, git.ApiProvider, token.ApiProvider)
	require.Equal(t, git.AllowedAuthenticationSecrets, token.AllowedAuthenticationSecrets)
	require.Equal(t, git.AllowedPrefixes, token.AllowedPrefixes)
	require.Equal(t, git.BlockedPrefixes, token.BlockedPrefixes)
	require.Equal(t, git.Comment, token.Comment)
	require.Nil(t, (*ApiIntegrationGitHttpsApiDetails)(nil).AsToken())

	githubApp := git.AsGithubApp()
	require.Equal(t, git.UserAuthType, githubApp.UserAuthType)
	require.Equal(t, git.ApiProvider, githubApp.ApiProvider)
	require.Nil(t, (*ApiIntegrationGitHttpsApiDetails)(nil).AsGithubApp())

	oauth2 := git.AsOauth2()
	require.Equal(t, git.OauthClientId, oauth2.OauthClientId)
	require.Equal(t, git.OauthAllowedScopes, oauth2.OauthAllowedScopes)
	require.Equal(t, git.OauthUsername, oauth2.OauthUsername)
	require.Nil(t, (*ApiIntegrationGitHttpsApiDetails)(nil).AsOauth2())

	privateLink := git.AsPrivateLink()
	require.Equal(t, git.UsePrivatelinkEndpoint, privateLink.UsePrivatelinkEndpoint)
	require.Equal(t, git.TlsTrustedCertificates, privateLink.TlsTrustedCertificates)
	require.Equal(t, git.AllowedAuthenticationSecrets, privateLink.AllowedAuthenticationSecrets)
	require.Nil(t, (*ApiIntegrationGitHttpsApiDetails)(nil).AsPrivateLink())

	mcp := &ApiIntegrationExternalMcpDetails{ //nolint:gosec // test credentials
		Id:                         id,
		Enabled:                    true,
		ApiProvider:                "EXTERNAL_MCP",
		UserAuthType:               "OAUTH2",
		OauthGrant:                 "CLIENT_CREDENTIALS",
		OauthClientId:              "mcp-client",
		OauthClientAuthMethod:      "CLIENT_SECRET_BASIC",
		OauthTokenEndpoint:         "https://mcp.example.com/token",
		OauthAuthorizationEndpoint: "https://mcp.example.com/authorize",
		OauthAccessTokenValidity:   1800,
		OauthRefreshTokenValidity:  7200,
		OauthAllowedScopes:         []string{"mcp"},
		OauthUsername:              "mcp-user",
		OauthAssertionIssuer:       "mcp-issuer",
		OauthResourceUrl:           "https://mcp.example.com/resource",
		AllowedPrefixes:            []string{"https://mcp.example.com/"},
		BlockedPrefixes:            []string{"https://blocked.example.com/"},
		Comment:                    "mcp comment",
	}

	mcpOauth2 := mcp.AsOauth2()
	require.Equal(t, mcp.OauthGrant, mcpOauth2.OauthGrant)
	require.Equal(t, mcp.OauthClientAuthMethod, mcpOauth2.OauthClientAuthMethod)
	require.Equal(t, mcp.OauthAssertionIssuer, mcpOauth2.OauthAssertionIssuer)
	require.Equal(t, mcp.ApiProvider, mcpOauth2.ApiProvider)
	require.Nil(t, (*ApiIntegrationExternalMcpDetails)(nil).AsOauth2())

	dynamicClient := mcp.AsDynamicClient()
	require.Equal(t, mcp.OauthResourceUrl, dynamicClient.OauthResourceUrl)
	require.Equal(t, mcp.UserAuthType, dynamicClient.UserAuthType)
	require.Equal(t, mcp.AllowedPrefixes, dynamicClient.AllowedPrefixes)
	require.Nil(t, (*ApiIntegrationExternalMcpDetails)(nil).AsDynamicClient())
}
