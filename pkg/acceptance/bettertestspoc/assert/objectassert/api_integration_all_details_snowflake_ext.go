package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationAllDetailsAssert) HasApiAwsIamUserArnNotEmpty() *ApiIntegrationAllDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAllDetails) error {
		t.Helper()
		if o.ApiAwsIamUserArn == "" {
			return fmt.Errorf("expected api aws iam user arn not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAllDetailsAssert) HasApiAwsExternalIdNotEmpty() *ApiIntegrationAllDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAllDetails) error {
		t.Helper()
		if o.ApiAwsExternalId == "" {
			return fmt.Errorf("expected api aws external id not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAllDetailsAssert) HasAzureMultiTenantAppNameNotEmpty() *ApiIntegrationAllDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAllDetails) error {
		t.Helper()
		if o.AzureMultiTenantAppName == "" {
			return fmt.Errorf("expected azure multi tenant app name not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAllDetailsAssert) HasAzureConsentUrlNotEmpty() *ApiIntegrationAllDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAllDetails) error {
		t.Helper()
		if o.AzureConsentUrl == "" {
			return fmt.Errorf("expected azure consent url not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAllDetailsAssert) HasGoogleApiServiceAccountNotEmpty() *ApiIntegrationAllDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAllDetails) error {
		t.Helper()
		if o.GoogleApiServiceAccount == "" {
			return fmt.Errorf("expected google api service account not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAllDetailsAssert) HasApiProviderNotEmpty() *ApiIntegrationAllDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAllDetails) error {
		t.Helper()
		if o.ApiProvider == "" {
			return fmt.Errorf("expected api provider not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAllDetailsAssert) HasNoBlockedPrefixes() *ApiIntegrationAllDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAllDetails) error {
		t.Helper()
		if len(o.BlockedPrefixes) != 0 {
			return fmt.Errorf("expected no blocked prefixes; got: %v", o.BlockedPrefixes)
		}
		return nil
	})
	return a
}
