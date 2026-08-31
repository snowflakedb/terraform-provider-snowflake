package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowRuntimeDetailsAssert) HasKeyNotEmpty() *OpenflowRuntimeDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowRuntimeDetails) error {
		t.Helper()
		if o.Key == nil || *o.Key == "" {
			return fmt.Errorf("expected key to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowRuntimeDetailsAssert) HasServerUrlNotEmpty() *OpenflowRuntimeDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowRuntimeDetails) error {
		t.Helper()
		if o.ServerUrl == nil || *o.ServerUrl == "" {
			return fmt.Errorf("expected server url to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowRuntimeDetailsAssert) HasNodeTypeTierNotEmpty() *OpenflowRuntimeDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowRuntimeDetails) error {
		t.Helper()
		if o.NodeTypeTier == nil || *o.NodeTypeTier == "" {
			return fmt.Errorf("expected node type tier to be not empty")
		}
		return nil
	})
	return o
}
