package objectassert

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (w *UserAssert) HasDefaults(name string) *UserAssert {
	return w.
		HasName(name).
		HasCreatedOnNotEmpty().
		// login name is always case-insensitive
		HasLoginName(strings.ToUpper(name)).
		HasFirstName("").
		HasLastName("").
		HasEmail("").
		HasMinsToUnlock("").
		HasDaysToExpiry("").
		HasComment("").
		HasDisabled(false).
		HasMustChangePassword(false).
		HasSnowflakeLock(false).
		HasDefaultWarehouse("").
		HasDefaultNamespace("").
		HasDefaultRole("").
		HasDefaultSecondaryRoles(`["ALL"]`).
		HasExtAuthnDuo(false).
		HasExtAuthnUid("").
		HasMinsToBypassMfa("").
		HasLastSuccessLoginEmpty().
		HasExpiresAtTimeEmpty().
		HasLockedUntilTimeEmpty().
		HasHasPassword(false).
		HasHasRsaPublicKey(false)
}

func (w *UserAssert) HasCreatedOnNotEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasOwnerNotEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.Owner == "" {
			return fmt.Errorf("expected owner not empty; got empty")
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasLastSuccessLoginEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.LastSuccessLogin != (time.Time{}) {
			return fmt.Errorf("expected last success login empty; got: %v", o.LastSuccessLogin)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasExpiresAtTimeEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.ExpiresAtTime != (time.Time{}) {
			return fmt.Errorf("expected expires at time empty; got: %v", o.ExpiresAtTime)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasLockedUntilTimeEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.LockedUntilTime != (time.Time{}) {
			return fmt.Errorf("expected locked until time empty; got: %v", o.LockedUntilTime)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasExpiresAtTimeNotEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.ExpiresAtTime == (time.Time{}) {
			return fmt.Errorf("expected expires at time not empty; got: %v", o.ExpiresAtTime)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasLockedUntilTimeNotEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.LockedUntilTime == (time.Time{}) {
			return fmt.Errorf("expected locked until time not empty; got: %v", o.LockedUntilTime)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDefaultNamespaceId(expected sdk.DatabaseObjectIdentifier) *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(o.DefaultNamespace).FullyQualifiedName() != expected.FullyQualifiedName() {
			return fmt.Errorf("expected default namespace: %v; got: %v", expected, o.DefaultNamespace)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDaysToExpiryNotEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.DaysToExpiry == "" {
			return fmt.Errorf("expected days to expiry not empty; got empty")
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDaysToExpiryEmpty() *UserAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.DaysToExpiry != "" {
			return fmt.Errorf("expected days to expiry empty; got: %v", o.DaysToExpiry)
		}
		return nil
	})
	return w
}
