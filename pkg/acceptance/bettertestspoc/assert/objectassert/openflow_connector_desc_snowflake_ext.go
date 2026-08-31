package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowConnectorDetailsAssert) HasConnectorUrlNotEmpty() *OpenflowConnectorDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnectorDetails) error {
		t.Helper()
		if o.ConnectorUrl == nil || *o.ConnectorUrl == "" {
			return fmt.Errorf("expected connector url to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowConnectorDetailsAssert) HasLiveVersionLocationUriNotEmpty() *OpenflowConnectorDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnectorDetails) error {
		t.Helper()
		if o.LiveVersionLocationUri == nil || *o.LiveVersionLocationUri == "" {
			return fmt.Errorf("expected live version location uri to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowConnectorDetailsAssert) HasDefaultVersionNameNotEmpty() *OpenflowConnectorDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnectorDetails) error {
		t.Helper()
		if o.DefaultVersionName == nil || *o.DefaultVersionName == "" {
			return fmt.Errorf("expected default version name to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowConnectorDetailsAssert) HasDefaultVersionLocationUriNotEmpty() *OpenflowConnectorDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnectorDetails) error {
		t.Helper()
		if o.DefaultVersionLocationUri == nil || *o.DefaultVersionLocationUri == "" {
			return fmt.Errorf("expected default version location uri to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowConnectorDetailsAssert) HasLastVersionNameNotEmpty() *OpenflowConnectorDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnectorDetails) error {
		t.Helper()
		if o.LastVersionName == nil || *o.LastVersionName == "" {
			return fmt.Errorf("expected last version name to be not empty")
		}
		return nil
	})
	return o
}
