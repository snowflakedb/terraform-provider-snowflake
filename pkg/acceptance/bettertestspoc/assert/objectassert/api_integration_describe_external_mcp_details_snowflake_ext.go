package objectassert

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ApiIntegrationExternalMcpDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ApiIntegrationExternalMcpDetails, sdk.AccountObjectIdentifier]
}

func ApiIntegrationExternalMcpDetails(t *testing.T, id sdk.AccountObjectIdentifier) *ApiIntegrationExternalMcpDetailsAssert {
	t.Helper()
	return &ApiIntegrationExternalMcpDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("ApiIntegrationExternalMcpDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ApiIntegrationExternalMcpDetails, sdk.AccountObjectIdentifier] {
			return testClient.ApiIntegration.DescribeExternalMcp
		}),
	}
}

func ApiIntegrationExternalMcpDetailsFromObject(t *testing.T, apiIntegrationExternalMcpDetails *sdk.ApiIntegrationExternalMcpDetails) *ApiIntegrationExternalMcpDetailsAssert {
	t.Helper()
	return &ApiIntegrationExternalMcpDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("ApiIntegrationExternalMcpDetails"), apiIntegrationExternalMcpDetails.ID(), apiIntegrationExternalMcpDetails),
	}
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasId(expected sdk.AccountObjectIdentifier) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.Id.Name() != expected.Name() {
			return fmt.Errorf("expected id: %v; got: %v", expected.Name(), o.Id.Name())
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasEnabled(expected bool) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.Enabled != expected {
			return fmt.Errorf("expected enabled: %v; got: %v", expected, o.Enabled)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasApiProvider(expected sdk.ApiIntegrationMcpApiProviderType) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.ApiProvider != strings.ToUpper(string(expected)) {
			return fmt.Errorf("expected api provider: %v; got: %v", expected, o.ApiProvider)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasUserAuthType(expected sdk.ApiIntegrationUserAuthType) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.UserAuthType != string(expected) {
			return fmt.Errorf("expected user auth type: %v; got: %v", expected, o.UserAuthType)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasOauthGrant(expected string) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.OauthGrant != expected {
			return fmt.Errorf("expected oauth grant: %v; got: %v", expected, o.OauthGrant)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasOauthClientId(expected string) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.OauthClientId != expected {
			return fmt.Errorf("expected oauth client id: %v; got: %v", expected, o.OauthClientId)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasOauthClientAuthMethod(expected sdk.ApiIntegrationOauthClientAuthMethod) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.OauthClientAuthMethod != string(expected) {
			return fmt.Errorf("expected oauth client auth method: %v; got: %v", expected, o.OauthClientAuthMethod)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasOauthTokenEndpoint(expected string) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.OauthTokenEndpoint != expected {
			return fmt.Errorf("expected oauth token endpoint: %v; got: %v", expected, o.OauthTokenEndpoint)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasOauthAuthorizationEndpoint(expected string) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.OauthAuthorizationEndpoint != expected {
			return fmt.Errorf("expected oauth authorization endpoint: %v; got: %v", expected, o.OauthAuthorizationEndpoint)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasOauthAccessTokenValidity(expected int) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.OauthAccessTokenValidity != expected {
			return fmt.Errorf("expected oauth access token validity: %v; got: %v", expected, o.OauthAccessTokenValidity)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasOauthRefreshTokenValidity(expected int) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.OauthRefreshTokenValidity != expected {
			return fmt.Errorf("expected oauth refresh token validity: %v; got: %v", expected, o.OauthRefreshTokenValidity)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
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

func (a *ApiIntegrationExternalMcpDetailsAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
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

func (a *ApiIntegrationExternalMcpDetailsAssert) HasComment(expected string) *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasNoBlockedPrefixes() *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if len(o.BlockedPrefixes) != 0 {
			return fmt.Errorf("expected no blocked prefixes; got: %v", o.BlockedPrefixes)
		}
		return nil
	})
	return a
}
