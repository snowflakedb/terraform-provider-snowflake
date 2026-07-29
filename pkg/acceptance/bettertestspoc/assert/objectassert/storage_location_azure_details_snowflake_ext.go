package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLocationAzureDetailsAssert) HasAzureMultiTenantAppNameNotEmpty() *StorageLocationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLocationAzureDetails) error {
		t.Helper()
		if o.AzureMultiTenantAppName == "" {
			return fmt.Errorf("expected azure multi tenant app name not empty; got empty")
		}
		return nil
	})
	return s
}

func (s *StorageLocationAzureDetailsAssert) HasAzureConsentUrlNotEmpty() *StorageLocationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLocationAzureDetails) error {
		t.Helper()
		if o.AzureConsentUrl == "" {
			return fmt.Errorf("expected azure consent url not empty; got empty")
		}
		return nil
	})
	return s
}

func (s *StorageLocationAzureDetailsAssert) HasUsePrivatelinkEndpoint(expected bool) *StorageLocationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLocationAzureDetails) error {
		t.Helper()
		if o.UsePrivatelinkEndpoint == nil {
			return fmt.Errorf("expected use privatelink endpoint: %v; got: nil", expected)
		}
		if *o.UsePrivatelinkEndpoint != expected {
			return fmt.Errorf("expected use privatelink endpoint: %v; got: %v", expected, *o.UsePrivatelinkEndpoint)
		}
		return nil
	})
	return s
}

func (s *StorageLocationAzureDetailsAssert) HasUsePrivatelinkEndpointEmpty() *StorageLocationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLocationAzureDetails) error {
		t.Helper()
		if o.UsePrivatelinkEndpoint != nil {
			return fmt.Errorf("expected use privatelink endpoint to be nil; got: %v", *o.UsePrivatelinkEndpoint)
		}
		return nil
	})
	return s
}
