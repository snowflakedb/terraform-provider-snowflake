// Code generated by assertions generator; DO NOT EDIT.

package objectassert

import (
	"fmt"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type StreamAssert struct {
	*assert.SnowflakeObjectAssert[sdk.Stream, sdk.SchemaObjectIdentifier]
}

func Stream(t *testing.T, id sdk.SchemaObjectIdentifier) *StreamAssert {
	t.Helper()
	return &StreamAssert{
		assert.NewSnowflakeObjectAssertWithProvider(sdk.ObjectTypeStream, id, acc.TestClient().Stream.Show),
	}
}

func StreamFromObject(t *testing.T, stream *sdk.Stream) *StreamAssert {
	t.Helper()
	return &StreamAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectTypeStream, stream.ID(), stream),
	}
}

func (s *StreamAssert) HasCreatedOn(expected time.Time) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasName(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasDatabaseName(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.DatabaseName != expected {
			return fmt.Errorf("expected database name: %v; got: %v", expected, o.DatabaseName)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasSchemaName(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.SchemaName != expected {
			return fmt.Errorf("expected schema name: %v; got: %v", expected, o.SchemaName)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasOwner(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.Owner == nil {
			return fmt.Errorf("expected owner to have value; got: nil")
		}
		if *o.Owner != expected {
			return fmt.Errorf("expected owner: %v; got: %v", expected, *o.Owner)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasComment(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.Comment == nil {
			return fmt.Errorf("expected comment to have value; got: nil")
		}
		if *o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, *o.Comment)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasTableName(expected string) *StreamAssert {
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

func (s *StreamAssert) HasSourceType(expected string) *StreamAssert {
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

func (s *StreamAssert) HasBaseTables(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.BaseTables == nil {
			return fmt.Errorf("expected base tables to have value; got: nil")
		}
		if *o.BaseTables != expected {
			return fmt.Errorf("expected base tables: %v; got: %v", expected, *o.BaseTables)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasType(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.Type == nil {
			return fmt.Errorf("expected type to have value; got: nil")
		}
		if *o.Type != expected {
			return fmt.Errorf("expected type: %v; got: %v", expected, *o.Type)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasStale(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.Stale == nil {
			return fmt.Errorf("expected stale to have value; got: nil")
		}
		if *o.Stale != expected {
			return fmt.Errorf("expected stale: %v; got: %v", expected, *o.Stale)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasMode(expected string) *StreamAssert {
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

func (s *StreamAssert) HasStaleAfter(expected time.Time) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.StaleAfter == nil {
			return fmt.Errorf("expected stale after to have value; got: nil")
		}
		if *o.StaleAfter != expected {
			return fmt.Errorf("expected stale after: %v; got: %v", expected, *o.StaleAfter)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasInvalidReason(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.InvalidReason == nil {
			return fmt.Errorf("expected invalid reason to have value; got: nil")
		}
		if *o.InvalidReason != expected {
			return fmt.Errorf("expected invalid reason: %v; got: %v", expected, *o.InvalidReason)
		}
		return nil
	})
	return s
}

func (s *StreamAssert) HasOwnerRoleType(expected string) *StreamAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stream) error {
		t.Helper()
		if o.OwnerRoleType == nil {
			return fmt.Errorf("expected owner role type to have value; got: nil")
		}
		if *o.OwnerRoleType != expected {
			return fmt.Errorf("expected owner role type: %v; got: %v", expected, *o.OwnerRoleType)
		}
		return nil
	})
	return s
}
