package objectassert

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CortexAgentAssert) HasNoComment() *CortexAgentAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgent) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be empty")
		}
		return nil
	})
	return c
}

func (c *CortexAgentAssert) HasNoProfile() *CortexAgentAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgent) error {
		t.Helper()
		if o.Profile != nil {
			return fmt.Errorf("expected profile to be empty")
		}
		return nil
	})
	return c
}

func (c *CortexAgentAssert) HasCreatedOnNotEmpty() *CortexAgentAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgent) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return c
}

func (c *CortexAgentAssert) HasCortexAgentProfile(expected *sdk.CortexAgentProfile) *CortexAgentAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgent) error {
		t.Helper()
		return assertCortexAgentProfileJsonEqual(o.Profile, expected)
	})
	return c
}

func assertCortexAgentProfileJsonEqual(profile *string, expected *sdk.CortexAgentProfile) error {
	exp := expected
	if exp == nil {
		exp = &sdk.CortexAgentProfile{}
	}

	if profile == nil {
		return fmt.Errorf("expected profile to be non-nil")
	}

	actual, err := sdk.UnmarshalCortexAgentProfile(*profile)
	if err != nil {
		return fmt.Errorf("unmarshal cortex agent profile: %w", err)
	}

	if !reflect.DeepEqual(actual, exp) {
		return fmt.Errorf("expected profile %+v; got %+v", exp, actual)
	}
	return nil
}
