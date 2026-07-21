package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasApiProviderType(expected sdk.ApiIntegrationGitApiProviderType) *ApiIntegrationGitHttpsApiDetailsAssert {
	return a.HasApiProvider(strings.ToUpper(string(expected)))
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasUserAuthTypeEnum(expected sdk.ApiIntegrationUserAuthType) *ApiIntegrationGitHttpsApiDetailsAssert {
	return a.HasUserAuthType(string(expected))
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasOauthAllowedScopesEnum(expected ...sdk.ApiIntegrationOauthAllowedScope) *ApiIntegrationGitHttpsApiDetailsAssert {
	strs := make([]string, len(expected))
	for i, s := range expected {
		strs[i] = string(s)
	}
	return a.HasOauthAllowedScopes(strs...)
}

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasNoUserAuthType() *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.UserAuthType != "" {
			return fmt.Errorf("expected no user auth type; got: %v", o.UserAuthType)
		}
		return nil
	})
	return a
}
