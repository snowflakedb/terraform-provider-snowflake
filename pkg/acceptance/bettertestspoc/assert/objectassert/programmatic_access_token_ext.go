package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ProgrammaticAccessTokenAssert) HasNoRoleRestriction() *ProgrammaticAccessTokenAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ProgrammaticAccessToken) error {
		t.Helper()
		if o.RoleRestriction != nil {
			return fmt.Errorf("expected role_restriction to be empty; got: %s", o.RoleRestriction)
		}
		return nil
	})
	return c
}

func (c *ProgrammaticAccessTokenAssert) HasExpiresAtNotEmpty() *ProgrammaticAccessTokenAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ProgrammaticAccessToken) error {
		t.Helper()
		if o.ExpiresAt == (time.Time{}) {
			return fmt.Errorf("expected expires_at to be not empty")
		}
		return nil
	})
	return c
}

func (c *ProgrammaticAccessTokenAssert) HasExpiresAtBefore(expected time.Time) *ProgrammaticAccessTokenAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ProgrammaticAccessToken) error {
		t.Helper()
		if !o.ExpiresAt.Before(expected) {
			return fmt.Errorf("expected expires_at to be before %s; got: %s", expected, o.ExpiresAt)
		}
		return nil
	})
	return c
}

func (c *ProgrammaticAccessTokenAssert) HasNoComment() *ProgrammaticAccessTokenAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ProgrammaticAccessToken) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be nil; got: %s", *o.Comment)
		}
		return nil
	})
	return c
}

func (c *ProgrammaticAccessTokenAssert) HasCreatedOnNotEmpty() *ProgrammaticAccessTokenAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ProgrammaticAccessToken) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return c
}

func (c *ProgrammaticAccessTokenAssert) HasNoMinsToBypassNetworkPolicyRequirement() *ProgrammaticAccessTokenAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ProgrammaticAccessToken) error {
		t.Helper()
		if o.MinsToBypassNetworkPolicyRequirement != nil {
			return fmt.Errorf("expected mins_to_bypass_network_policy_requirement to be empty; got: %d", *o.MinsToBypassNetworkPolicyRequirement)
		}
		return nil
	})
	return c
}

func (c *ProgrammaticAccessTokenAssert) HasRotatedToEmpty() *ProgrammaticAccessTokenAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ProgrammaticAccessToken) error {
		t.Helper()
		if o.RotatedTo != nil {
			return fmt.Errorf("expected rotated_to to be empty; got: %s", *o.RotatedTo)
		}
		return nil
	})
	return c
}

// HasMinsToBypassNetworkPolicyRequirementWithTolerance can be used to match with an allowed tolerance of 1 minute.
// This is useful because the actual value returned by the API is not always the same as the value set - can be 1 minute off.
func (p *ProgrammaticAccessTokenAssert) HasMinsToBypassNetworkPolicyRequirementWithTolerance(expected int) *ProgrammaticAccessTokenAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.ProgrammaticAccessToken) error {
		t.Helper()
		if o.MinsToBypassNetworkPolicyRequirement == nil {
			return fmt.Errorf("expected mins to bypass network policy requirement to have value; got: nil")
		}
		if *o.MinsToBypassNetworkPolicyRequirement > expected || *o.MinsToBypassNetworkPolicyRequirement < expected-1 {
			return fmt.Errorf("expected mins to bypass network policy requirement: %v; got: %v", expected, *o.MinsToBypassNetworkPolicyRequirement)
		}
		return nil
	})
	return p
}
