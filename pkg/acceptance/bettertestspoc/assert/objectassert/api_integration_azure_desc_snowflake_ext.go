package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

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
