package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationGoogleDetailsAssert) HasGoogleApiServiceAccountNotEmpty() *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.GoogleApiServiceAccount == "" {
			return fmt.Errorf("expected google api service account not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasNoBlockedPrefixes() *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if len(o.BlockedPrefixes) != 0 {
			return fmt.Errorf("expected no blocked prefixes; got: %v", o.BlockedPrefixes)
		}
		return nil
	})
	return a
}
