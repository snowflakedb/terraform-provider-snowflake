package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type PostgresInstanceAssert struct {
	*assert.SnowflakeObjectAssert[sdk.PostgresInstance, sdk.AccountObjectIdentifier]
}

func PostgresInstance(t *testing.T, id sdk.AccountObjectIdentifier) *PostgresInstanceAssert {
	t.Helper()
	return &PostgresInstanceAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypePostgresInstance, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.PostgresInstance, sdk.AccountObjectIdentifier] {
			return testClient.PostgresInstance.Show
		}),
	}
}

func PostgresInstanceFromObject(t *testing.T, postgresInstance *sdk.PostgresInstance) *PostgresInstanceAssert {
	t.Helper()
	return &PostgresInstanceAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectTypePostgresInstance, postgresInstance.ID(), postgresInstance),
	}
}

func (p *PostgresInstanceAssert) HasName(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasOwner(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Owner != expected {
			return fmt.Errorf("expected owner: %v; got: %v", expected, o.Owner)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasOwnerRoleType(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.OwnerRoleType != expected {
			return fmt.Errorf("expected owner_role_type: %v; got: %v", expected, o.OwnerRoleType)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasCreatedOn(expected time.Time) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created_on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return p
}

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

func (p *PostgresInstanceAssert) HasUpdatedOn(expected time.Time) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.UpdatedOn != expected {
			return fmt.Errorf("expected updated_on: %v; got: %v", expected, o.UpdatedOn)
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

func (p *PostgresInstanceAssert) HasType(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Type != expected {
			return fmt.Errorf("expected type: %v; got: %v", expected, o.Type)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasOrigin(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Origin == nil {
			return fmt.Errorf("expected origin to have value; got: nil")
		}
		if *o.Origin != expected {
			return fmt.Errorf("expected origin: %v; got: %v", expected, *o.Origin)
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

func (p *PostgresInstanceAssert) HasHost(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Host == nil {
			return fmt.Errorf("expected host to have value; got: nil")
		}
		if *o.Host != expected {
			return fmt.Errorf("expected host: %v; got: %v", expected, *o.Host)
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

func (p *PostgresInstanceAssert) HasComputeFamily(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.ComputeFamily != expected {
			return fmt.Errorf("expected compute_family: %v; got: %v", expected, o.ComputeFamily)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasAuthenticationAuthority(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.AuthenticationAuthority != expected {
			return fmt.Errorf("expected authentication_authority: %v; got: %v", expected, o.AuthenticationAuthority)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasStorageSize(expected int) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.StorageSize != expected {
			return fmt.Errorf("expected storage_size: %v; got: %v", expected, o.StorageSize)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasPostgresVersion(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.PostgresVersion != expected {
			return fmt.Errorf("expected postgres_version: %v; got: %v", expected, o.PostgresVersion)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasPostgresSettings(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.PostgresSettings == nil {
			return fmt.Errorf("expected postgres_settings to have value; got: nil")
		}
		if *o.PostgresSettings != expected {
			return fmt.Errorf("expected postgres_settings: %v; got: %v", expected, *o.PostgresSettings)
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

func (p *PostgresInstanceAssert) HasIsHa(expected bool) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.IsHa != expected {
			return fmt.Errorf("expected is_ha: %v; got: %v", expected, o.IsHa)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasRetentionTime(expected int) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.RetentionTime != expected {
			return fmt.Errorf("expected retention_time: %v; got: %v", expected, o.RetentionTime)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasState(expected sdk.PostgresInstanceState) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.State != expected {
			return fmt.Errorf("expected state: %v; got: %v", expected, o.State)
		}
		return nil
	})
	return p
}

func (p *PostgresInstanceAssert) HasComment(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.Comment == nil {
			return fmt.Errorf("expected comment to have value; got: nil")
		}
		if *o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, *o.Comment)
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

func (p *PostgresInstanceAssert) HasPrivatelinkServiceIdentifier(expected string) *PostgresInstanceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstance) error {
		t.Helper()
		if o.PrivatelinkServiceIdentifier == nil {
			return fmt.Errorf("expected privatelink_service_identifier to have value; got: nil")
		}
		if *o.PrivatelinkServiceIdentifier != expected {
			return fmt.Errorf("expected privatelink_service_identifier: %v; got: %v", expected, *o.PrivatelinkServiceIdentifier)
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
