package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationExternalMcpDetailsAssert) HasApiProviderType(expected sdk.ApiIntegrationMcpApiProviderType) *ApiIntegrationExternalMcpDetailsAssert {
	return a.HasApiProvider(strings.ToUpper(string(expected)))
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasUserAuthTypeEnum(expected sdk.ApiIntegrationUserAuthType) *ApiIntegrationExternalMcpDetailsAssert {
	return a.HasUserAuthType(string(expected))
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasOauthClientAuthMethodEnum(expected sdk.ApiIntegrationOauthClientAuthMethod) *ApiIntegrationExternalMcpDetailsAssert {
	return a.HasOauthClientAuthMethod(string(expected))
}

func (a *ApiIntegrationExternalMcpDetailsAssert) HasNoUserAuthType() *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.UserAuthType != "" {
			return fmt.Errorf("expected no user auth type; got: %v", o.UserAuthType)
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
