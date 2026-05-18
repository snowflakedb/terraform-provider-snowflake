package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
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

func (c *CortexAgentDetailsAssert) HasCortexAgentSpec(expected string) *CortexAgentDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgentDetails) error {
		t.Helper()
		normalizedActual, err := sdk.NormalizeCortexAgentSpecification(o.AgentSpec)
		require.NoError(t, err)
		normalizedExpected, err := sdk.NormalizeCortexAgentSpecification(expected)
		require.NoError(t, err)

		if normalizedActual != normalizedExpected {
			return fmt.Errorf("expected agent spec: %v; got: %v", normalizedExpected, normalizedActual)
		}
		return nil
	})
	return c
}
