package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (p *PostgresInstanceDetailsAssert) HasCreatedOnNotEmpty() *PostgresInstanceDetailsAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstanceDetails) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceDetailsAssert) HasUpdatedOnNotEmpty() *PostgresInstanceDetailsAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstanceDetails) error {
		t.Helper()
		if o.UpdatedOn == "" {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceDetailsAssert) HasHostNotEmpty() *PostgresInstanceDetailsAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstanceDetails) error {
		t.Helper()
		if o.Host == "" {
			return fmt.Errorf("expected host to be not empty")
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceDetailsAssert) HasStateNotEmpty() *PostgresInstanceDetailsAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstanceDetails) error {
		t.Helper()
		if o.State == "" {
			return fmt.Errorf("expected state to be not empty")
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceDetailsAssert) HasPostgresVersionNotEmpty() *PostgresInstanceDetailsAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstanceDetails) error {
		t.Helper()
		if o.PostgresVersion == "" {
			return fmt.Errorf("expected postgres_version to be not empty")
		}
		return nil
	})
	return p
}
