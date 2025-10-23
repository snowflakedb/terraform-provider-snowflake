package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SchemaAssert) HasCreatedOnNotEmpty() *SchemaAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Schema) error {
		t.Helper()
		if o.CreatedOn.IsZero() {
			return fmt.Errorf("expected created on to have value; got: empty")
		}
		return nil
	})
	return s
}

func (s *SchemaAssert) HasOwnerNotEmpty() *SchemaAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Schema) error {
		t.Helper()
		if o.Owner == "" {
			return fmt.Errorf("expected owner to have value; got: empty")
		}
		return nil
	})
	return s
}

func (s *SchemaAssert) HasOwnerRoleTypeNotEmpty() *SchemaAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Schema) error {
		t.Helper()
		if o.OwnerRoleType == "" {
			return fmt.Errorf("expected owner role type to have value; got: empty")
		}
		return nil
	})
	return s
}
