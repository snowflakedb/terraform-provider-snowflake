package objectassert

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CortexAgentDetailsAssert) HasNoComment() *CortexAgentDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgentDetails) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be empty")
		}
		return nil
	})
	return c
}

func (c *CortexAgentDetailsAssert) HasNoProfile() *CortexAgentDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgentDetails) error {
		t.Helper()
		if o.Profile != nil {
			return fmt.Errorf("expected profile to be empty")
		}
		return nil
	})
	return c
}

func (c *CortexAgentDetailsAssert) HasCreatedOnNotEmpty() *CortexAgentDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgentDetails) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return c
}

func (c *CortexAgentDetailsAssert) HasCortexAgentProfile(expected sdk.CortexAgentProfile) *CortexAgentDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgentDetails) error {
		t.Helper()
		return assertCortexAgentProfileJsonEqual(o.Profile, expected)
	})
	return c
}

func (c *CortexAgentDetailsAssert) HasCortexAgentSpec(expected map[string]any) *CortexAgentDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgentDetails) error {
		t.Helper()
		return assertCortexAgentSpecMapEqual(o.AgentSpec, expected)
	})
	return c
}

func assertCortexAgentSpecMapEqual(agentSpec string, expected map[string]any) error {
	exp := expected
	if exp == nil {
		exp = map[string]any{}
	}

	actual, err := sdk.UnmarshalCortexAgentSpec(agentSpec)
	if err != nil {
		return fmt.Errorf("unmarshal cortex agent spec: %w", err)
	}

	if !reflect.DeepEqual(actual, exp) {
		return fmt.Errorf("expected agent spec %+v; got %+v", exp, actual)
	}
	return nil
}
