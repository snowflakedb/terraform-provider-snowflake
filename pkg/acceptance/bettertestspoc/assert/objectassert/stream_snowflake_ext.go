package objectassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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

func (s *StreamAssert) HasSourceType(expected sdk.StreamSourceType) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.SourceType == nil {
			return fmt.Errorf("expected source type to have value; got: nil")
		}
		if *o.SourceType != expected {
			return fmt.Errorf("expected source type: %v; got: %v", expected, *o.SourceType)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasBaseTables(expected ...sdk.SchemaObjectIdentifier) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if len(o.BaseTables) != len(expected) {
			return fmt.Errorf("expected base tables length: %v; got: %v", len(expected), len(o.BaseTables))
		}
		var errs []error
		for _, wantId := range expected {
			if !slices.ContainsFunc(o.BaseTables, func(gotName string) bool {
				gotId, err := sdk.ParseSchemaObjectIdentifier(gotName)
				if err != nil {
					errs = append(errs, err)
				}
				return wantId.FullyQualifiedName() == gotId.FullyQualifiedName()
			}) {
				errs = append(errs, fmt.Errorf("expected id: %s, to be in the list ids: %v", wantId.FullyQualifiedName(), o.BaseTables))
			}
		}
		return errors.Join(errs...)
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
			if !slices.Contains(o.BaseTables, wantName) {
				errs = append(errs, fmt.Errorf("expected name: %s, to be in the list ids: %v", wantName, o.BaseTables))
			}
		}
		return errors.Join(errs...)
	})
	return s
}

func (s *StreamAssert) HasMode(expected sdk.StreamMode) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.Mode == nil {
			return fmt.Errorf("expected mode to have value; got: nil")
		}
		if *o.Mode != expected {
			return fmt.Errorf("expected mode: %v; got: %v", expected, *o.Mode)
		}
		return nil
	})
	return s
}
