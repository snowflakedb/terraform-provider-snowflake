package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationAzureDetailsAssert) HasApiProviderType(expected sdk.ApiIntegrationAzureApiProviderType) *ApiIntegrationAzureDetailsAssert {
	return a.HasApiProvider(strings.ToUpper(string(expected)))
}

func (a *ApiIntegrationAzureDetailsAssert) HasApiKeyNotEmpty() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.ApiKey == "" {
			return fmt.Errorf("expected api key not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasNoAllowedPrefixes() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if len(o.AllowedPrefixes) != 0 {
			return fmt.Errorf("expected no allowed prefixes; got: %v", o.AllowedPrefixes)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAzureMultiTenantAppNameNotEmpty() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.AzureMultiTenantAppName == "" {
			return fmt.Errorf("expected azure multi tenant app name not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAzureConsentUrlNotEmpty() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.AzureConsentUrl == "" {
			return fmt.Errorf("expected azure consent url not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasNoBlockedPrefixes() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if len(o.BlockedPrefixes) != 0 {
			return fmt.Errorf("expected no blocked prefixes; got: %v", o.BlockedPrefixes)
		}
		return nil
	})
	return a
}
