package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowConnectorAssert) HasLiveVersionLocationUriNotEmpty() *OpenflowConnectorAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnector) error {
		t.Helper()
		if o.LiveVersionLocationUri == nil || *o.LiveVersionLocationUri == "" {
			return fmt.Errorf("expected live version location uri to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowConnectorAssert) HasDefaultVersionNameNotEmpty() *OpenflowConnectorAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnector) error {
		t.Helper()
		if o.DefaultVersionName == nil || *o.DefaultVersionName == "" {
			return fmt.Errorf("expected default version name to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowConnectorAssert) HasConnectorUrlNotEmpty() *OpenflowConnectorAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnector) error {
		t.Helper()
		if o.ConnectorUrl == nil || *o.ConnectorUrl == "" {
			return fmt.Errorf("expected connector url to be not empty")
		}
		return nil
	})
	return o
}
