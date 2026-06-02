package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// Adjusted manually
func (i *IcebergTableDetailsAssert) HasNoDefault() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.Default != nil {
			return fmt.Errorf("expected default to be nil; got: %s", *o.Default)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableDetailsAssert) HasNoCheck() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.Check != nil {
			return fmt.Errorf("expected check to be nil; got: %s", *o.Check)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableDetailsAssert) HasNoExpression() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.Expression != nil {
			return fmt.Errorf("expected expression to be nil; got: %s", *o.Expression)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableDetailsAssert) HasNoComment() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be nil; got: %s", *o.Comment)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableDetailsAssert) HasNoPolicyName() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.PolicyName != nil {
			return fmt.Errorf("expected policy name to be nil; got: %s", o.PolicyName.FullyQualifiedName())
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableDetailsAssert) HasNoPrivacyDomain() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.PrivacyDomain != nil {
			return fmt.Errorf("expected privacy domain to be nil; got: %s", *o.PrivacyDomain)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableDetailsAssert) HasNoNameMapping() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.NameMapping != nil {
			return fmt.Errorf("expected name mapping to be nil; got: %s", *o.NameMapping)
		}
		return nil
	})
	return i
}

// Adjusted manually
func (i *IcebergTableDetailsAssert) HasNoWriteDefault() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.WriteDefault != nil {
			return fmt.Errorf("expected write default to be nil; got: %s", *o.WriteDefault)
		}
		return nil
	})
	return i
}
