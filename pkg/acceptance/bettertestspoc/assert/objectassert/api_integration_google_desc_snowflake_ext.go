package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationGoogleDetailsAssert) HasApiProviderType(expected sdk.ApiIntegrationGoogleApiProviderType) *ApiIntegrationGoogleDetailsAssert {
	return a.HasApiProvider(strings.ToUpper(string(expected)))
}

func (a *ApiIntegrationGoogleDetailsAssert) HasApiKeyEmpty() *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.ApiKey != "" {
			return fmt.Errorf("expected api key: %v to be empty", o.ApiKey)
		}
		return nil
	})
	return a
}
