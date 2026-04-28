package objectassert

import (
	"fmt"
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

func (p *PostgresInstanceAssert) HasNoOrigin() *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Origin != nil {
			return fmt.Errorf("expected origin to have nil; got: %s", *o.Origin)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasNoHost() *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Host != nil {
			return fmt.Errorf("expected host to have nil; got: %s", *o.Host)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasNoPostgresSettings() *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.PostgresSettings != nil {
			return fmt.Errorf("expected postgres_settings to have nil; got: %s", *o.PostgresSettings)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasNoComment() *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to have nil; got: %s", *o.Comment)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasNoPrivatelinkServiceIdentifier() *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.PrivatelinkServiceIdentifier != nil {
			return fmt.Errorf("expected privatelink_service_identifier to have nil; got: %s", *o.PrivatelinkServiceIdentifier)
		}
		return nil
	})
	return p
}
