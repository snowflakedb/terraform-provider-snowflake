package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationGitHttpsApiDetailsAssert) HasApiProviderNotEmpty() *ApiIntegrationGitHttpsApiDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGitHttpsApiDetails) error {
		t.Helper()
		if o.ApiProvider == "" {
			return fmt.Errorf("expected api provider not empty; got empty")
		}
		return nil
	})
	return a
}
