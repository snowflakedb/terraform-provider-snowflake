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

func (o *OpenflowRuntimeDetailsAssert) HasDisplayNameNotEmpty() *OpenflowRuntimeDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowRuntimeDetails) error {
		t.Helper()
		if o.DisplayName == nil || *o.DisplayName == "" {
			return fmt.Errorf("expected display name to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowRuntimeDetailsAssert) HasCommentNotEmpty() *OpenflowRuntimeDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowRuntimeDetails) error {
		t.Helper()
		if o.Comment == nil || *o.Comment == "" {
			return fmt.Errorf("expected comment to be not empty")
		}
		return nil
	})
	return o
}

func (o *OpenflowRuntimeDetailsAssert) HasExternalAccessIntegrationsNotEmpty() *OpenflowRuntimeDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowRuntimeDetails) error {
		t.Helper()
		if len(o.ExternalAccessIntegrations) == 0 {
			return fmt.Errorf("expected external access integrations to be not empty")
		}
		return nil
	})
	return o
}
