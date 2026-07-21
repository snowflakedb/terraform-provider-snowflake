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
