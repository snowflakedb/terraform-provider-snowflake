package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationExternalMcpDetailsAssert) HasApiProviderNotEmpty() *ApiIntegrationExternalMcpDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationExternalMcpDetails) error {
		t.Helper()
		if o.ApiProvider == "" {
			return fmt.Errorf("expected api provider not empty; got empty")
		}
		return nil
	})
	return a
}
