package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// Adjusted manually
func (i *IcebergTableAssert) HasNoOwner() *IcebergTableAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTable) error {
		t.Helper()
		if o.Owner != nil {
			return fmt.Errorf("expected owner to be nil; got: %s", *o.Owner)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableAssert) HasNoExternalVolumeName() *IcebergTableAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTable) error {
		t.Helper()
		if o.ExternalVolumeName != nil {
			return fmt.Errorf("expected external volume name to be nil; got: %s", (*o.ExternalVolumeName).Name())
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableAssert) HasNoCatalogName() *IcebergTableAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTable) error {
		t.Helper()
		if o.CatalogName != nil {
			return fmt.Errorf("expected catalog name to be nil; got: %s", (*o.CatalogName).Name())
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableAssert) HasNoCatalogTableName() *IcebergTableAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTable) error {
		t.Helper()
		if o.CatalogTableName != nil {
			return fmt.Errorf("expected catalog table name to be nil; got: %s", *o.CatalogTableName)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableAssert) HasNoCatalogNamespace() *IcebergTableAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTable) error {
		t.Helper()
		if o.CatalogNamespace != nil {
			return fmt.Errorf("expected catalog namespace to be nil; got: %s", *o.CatalogNamespace)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableAssert) HasNoComment() *IcebergTableAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTable) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be nil; got: %s", *o.Comment)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableAssert) HasNoNameMapping() *IcebergTableAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTable) error {
		t.Helper()
		if o.NameMapping != nil {
			return fmt.Errorf("expected name mapping to be nil; got: %s", *o.NameMapping)
		}
		return nil
	})
	return i
}
