package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type TagReferenceAssert struct {
	*assert.SnowflakeObjectAssert[sdk.TagReference, sdk.AccountObjectIdentifier]
}

func TagReferenceFromObject(t *testing.T, tagReference *sdk.TagReference) *TagReferenceAssert {
	t.Helper()
	return &TagReferenceAssert{
		assert.NewSnowflakeObjectAssertWithObject(
			sdk.ObjectType("TagReference"),
			sdk.NewAccountObjectIdentifier(fmt.Sprintf("%s.%s.%s", tagReference.TagDatabase, tagReference.TagSchema, tagReference.TagName)),
			tagReference,
		),
	}
}

func (a *TagReferenceAssert) HasTagDatabase(expected string) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.TagDatabase != expected {
			return fmt.Errorf("expected tag database: %v; got: %v", expected, o.TagDatabase)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasTagSchema(expected string) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.TagSchema != expected {
			return fmt.Errorf("expected tag schema: %v; got: %v", expected, o.TagSchema)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasTagName(expected string) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.TagName != expected {
			return fmt.Errorf("expected tag name: %v; got: %v", expected, o.TagName)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasTagValue(expected string) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.TagValue != expected {
			return fmt.Errorf("expected tag value: %v; got: %v", expected, o.TagValue)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasLevel(expected sdk.TagReferenceObjectDomain) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.Level != expected {
			return fmt.Errorf("expected level: %v; got: %v", expected, o.Level)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasObjectDatabase(expected string) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.ObjectDatabase == nil {
			return fmt.Errorf("expected object database: %v; got: nil", expected)
		}
		if *o.ObjectDatabase != expected {
			return fmt.Errorf("expected object database: %v; got: %v", expected, *o.ObjectDatabase)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasObjectSchema(expected string) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.ObjectSchema == nil {
			return fmt.Errorf("expected object schema: %v; got: nil", expected)
		}
		if *o.ObjectSchema != expected {
			return fmt.Errorf("expected object schema: %v; got: %v", expected, *o.ObjectSchema)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasObjectName(expected string) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.ObjectName != expected {
			return fmt.Errorf("expected object name: %v; got: %v", expected, o.ObjectName)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasDomain(expected sdk.TagReferenceObjectDomain) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.Domain != expected {
			return fmt.Errorf("expected domain: %v; got: %v", expected, o.Domain)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasColumnNameNil() *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.ColumnName != nil {
			return fmt.Errorf("expected column name to be nil; got: %v", *o.ColumnName)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasColumnName(expected string) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.ColumnName == nil {
			return fmt.Errorf("expected column name: %v; got: nil", expected)
		}
		if *o.ColumnName != expected {
			return fmt.Errorf("expected column name: %v; got: %v", expected, *o.ColumnName)
		}
		return nil
	})
	return a
}

func (a *TagReferenceAssert) HasApplyMethod(expected sdk.TagReferenceApplyMethod) *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.ApplyMethod != expected {
			return fmt.Errorf("expected apply method: %v; got: %v", expected, o.ApplyMethod)
		}
		return nil
	})
	return a
}
