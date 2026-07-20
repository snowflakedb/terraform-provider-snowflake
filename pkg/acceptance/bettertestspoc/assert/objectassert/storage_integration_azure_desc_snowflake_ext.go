package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAzureDetailsAssert) HasConsentUrlNotEmpty() *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.ConsentUrl == "" {
			return fmt.Errorf("expected consent url not empty; got empty")
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAzureDetailsAssert) HasMultiTenantAppNameNotEmpty() *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.MultiTenantAppName == "" {
			return fmt.Errorf("expected multi tenant app name not empty; got empty")
		}
		return nil
	})
	return s
}
