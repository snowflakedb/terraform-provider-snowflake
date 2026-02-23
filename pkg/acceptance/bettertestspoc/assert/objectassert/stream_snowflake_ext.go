package objectassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StreamAssert) HasTableId(expected sdk.SchemaObjectIdentifier) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.TableName == nil {
			return fmt.Errorf("expected table name to have value; got: nil")
		}
		gotTableId, err := sdk.ParseSchemaObjectIdentifier(*o.TableName)
		if err != nil {
			return err
		}
		if gotTableId.FullyQualifiedName() != expected.FullyQualifiedName() {
			return fmt.Errorf("expected table name: %v; got: %v", expected, *o.TableName)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasStageName(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.TableName == nil {
			return fmt.Errorf("expected table name to have value; got: nil")
		}
		if *o.TableName != expected {
			return fmt.Errorf("expected table name: %v; got: %v", expected, *o.TableName)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasBaseTablesPartiallyQualified(expected ...string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if len(o.BaseTables) != len(expected) {
			return fmt.Errorf("expected base tables length: %v; got: %v", len(expected), len(o.BaseTables))
		}
		var errs []error
		for _, wantName := range expected {
			found := false
			for _, gotName := range o.BaseTables {
				if gotName == wantName {
					found = true
					break
				}
				gotId, err := sdk.ParseSchemaObjectIdentifier(gotName)
				if err == nil {
					wantId, err := sdk.ParseSchemaObjectIdentifier(wantName)
					if err == nil && gotId.FullyQualifiedName() == wantId.FullyQualifiedName() {
						found = true
						break
					}
					if gotId.Name() == wantName {
						found = true
						break
					}
				}
			}
			if !found {
				errs = append(errs, fmt.Errorf("expected name: %s, to be in the list ids: %v", wantName, o.BaseTables))
			}
		}
		return errors.Join(errs...)
	})
	return s
}
