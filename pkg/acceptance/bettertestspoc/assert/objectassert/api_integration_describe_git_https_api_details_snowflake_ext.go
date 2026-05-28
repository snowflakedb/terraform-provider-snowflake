package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ApiIntegrationGitHttpsApiDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ApiIntegrationGitHttpsApiDetails, sdk.AccountObjectIdentifier]
}

func ApiIntegrationGitHttpsApiDetails(t *testing.T, id sdk.AccountObjectIdentifier) *ApiIntegrationGitHttpsApiDetailsAssert {
	t.Helper()
	return &ApiIntegrationGitHttpsApiDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("ApiIntegrationGitHttpsApiDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ApiIntegrationGitHttpsApiDetails, sdk.AccountObjectIdentifier] {
			return testClient.ApiIntegration.DescribeGitHttpsApi
		}),
	}
}

func ApiIntegrationGitHttpsApiDetailsFromObject(t *testing.T, apiIntegrationGitHttpsApiDetails *sdk.ApiIntegrationGitHttpsApiDetails) *ApiIntegrationGitHttpsApiDetailsAssert {
	t.Helper()
	return &ApiIntegrationGitHttpsApiDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("ApiIntegrationGitHttpsApiDetails"), apiIntegrationGitHttpsApiDetails.ID(), apiIntegrationGitHttpsApiDetails),
	}
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasId(expected sdk.AccountObjectIdentifier) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.Id.Name() != expected.Name() {
			return fmt.Errorf("expected id: %v; got: %v", expected.Name(), o.Id.Name())
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasEnabled(expected bool) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.Enabled != expected {
			return fmt.Errorf("expected enabled: %v; got: %v", expected, o.Enabled)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasApiProvider(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.ApiProvider != expected {
			return fmt.Errorf("expected api provider: %v; got: %v", expected, o.ApiProvider)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasAllowedAuthenticationSecrets(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.AllowedAuthenticationSecrets != expected {
			return fmt.Errorf("expected allowed authentication secrets: %v; got: %v", expected, o.AllowedAuthenticationSecrets)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasUserAuthType(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.UserAuthType != expected {
			return fmt.Errorf("expected user auth type: %v; got: %v", expected, o.UserAuthType)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthGrant(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthGrant != expected {
			return fmt.Errorf("expected oauth grant: %v; got: %v", expected, o.OauthGrant)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthClientId(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthClientId != expected {
			return fmt.Errorf("expected oauth client id: %v; got: %v", expected, o.OauthClientId)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthClientAuthMethod(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthClientAuthMethod != expected {
			return fmt.Errorf("expected oauth client auth method: %v; got: %v", expected, o.OauthClientAuthMethod)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthTokenEndpoint(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthTokenEndpoint != expected {
			return fmt.Errorf("expected oauth token endpoint: %v; got: %v", expected, o.OauthTokenEndpoint)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthAuthorizationEndpoint(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthAuthorizationEndpoint != expected {
			return fmt.Errorf("expected oauth authorization endpoint: %v; got: %v", expected, o.OauthAuthorizationEndpoint)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthAccessTokenValidity(expected int) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthAccessTokenValidity != expected {
			return fmt.Errorf("expected oauth access token validity: %v; got: %v", expected, o.OauthAccessTokenValidity)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthRefreshTokenValidity(expected int) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthRefreshTokenValidity != expected {
			return fmt.Errorf("expected oauth refresh token validity: %v; got: %v", expected, o.OauthRefreshTokenValidity)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthAllowedScopes(expected ...string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		mapped := collections.Map(o.OauthAllowedScopes, func(item string) any { return item })
		mappedExpected := collections.Map(expected, func(item string) any { return item })
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected oauth allowed scopes: %v; got: %v", expected, o.OauthAllowedScopes)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthUsername(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthUsername != expected {
			return fmt.Errorf("expected oauth username: %v; got: %v", expected, o.OauthUsername)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthAssertionIssuer(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthAssertionIssuer != expected {
			return fmt.Errorf("expected oauth assertion issuer: %v; got: %v", expected, o.OauthAssertionIssuer)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthResourceUrl(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.OauthResourceUrl != expected {
			return fmt.Errorf("expected oauth resource url: %v; got: %v", expected, o.OauthResourceUrl)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasUsePrivatelinkEndpoint(expected bool) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.UsePrivatelinkEndpoint != expected {
			return fmt.Errorf("expected use privatelink endpoint: %v; got: %v", expected, o.UsePrivatelinkEndpoint)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasTlsTrustedCertificates(expected ...string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		mapped := collections.Map(o.TlsTrustedCertificates, func(item string) any { return item })
		mappedExpected := collections.Map(expected, func(item string) any { return item })
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected tls trusted certificates: %v; got: %v", expected, o.TlsTrustedCertificates)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		mapped := collections.Map(o.AllowedPrefixes, func(item string) any { return item })
		mappedExpected := collections.Map(expected, func(item string) any { return item })
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected allowed prefixes: %v; got: %v", expected, o.AllowedPrefixes)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		mapped := collections.Map(o.BlockedPrefixes, func(item string) any { return item })
		mappedExpected := collections.Map(expected, func(item string) any { return item })
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected blocked prefixes: %v; got: %v", expected, o.BlockedPrefixes)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasComment(expected string) *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasApiProviderNotEmpty() *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.ApiProvider == "" {
			return fmt.Errorf("expected api provider not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasNoBlockedPrefixes() *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if len(o.BlockedPrefixes) != 0 {
			return fmt.Errorf("expected no blocked prefixes; got: %v", o.BlockedPrefixes)
		}
		return nil
	})
	return a
}
