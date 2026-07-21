package objectassert

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (p *PostgresInstanceAssert) HasCreatedOnNotEmpty() *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasUpdatedOnNotEmpty() *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.UpdatedOn == (time.Time{}) {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasStateOneOf(expected ...sdk.PostgresInstanceState) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if !slices.Contains(expected, o.State) {
			return fmt.Errorf("expected state one of: %v; got: %v", expected, o.State)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasPostgresVersionNotEmpty() *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.PostgresVersion == "" {
			return fmt.Errorf("expected postgres_version to be not empty")
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasOriginContaining(substring string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Origin == nil {
			return fmt.Errorf("expected origin to have value; got: nil")
		}
		if !strings.Contains(*o.Origin, substring) {
			return fmt.Errorf("expected origin to contain %q; got: %s", substring, *o.Origin)
		}
		return nil
	})
	return p
}
