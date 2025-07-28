package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/stretchr/testify/assert"
)

func TestIsValidDataType(t *testing.T) {
	t.Run("with valid data type", func(t *testing.T) {
		ok := IsValidDataType("VARCHAR")
		assert.True(t, ok)
	})

	t.Run("with invalid data type", func(t *testing.T) {
		ok := IsValidDataType("foo")
		assert.False(t, ok)
	})
}

func TestValidObjectIdentifier(t *testing.T) {
	t.Run("with valid object identifier", func(t *testing.T) {
		ok := ValidObjectIdentifier(randomAccountObjectIdentifier())
		assert.True(t, ok)
	})

	t.Run("with invalid object identifier", func(t *testing.T) {
		ok := ValidObjectIdentifier(emptyAccountObjectIdentifier)
		assert.False(t, ok)
	})

	t.Run("over 255 characters", func(t *testing.T) {
		ok := ValidObjectIdentifier(invalidAccountObjectIdentifier)
		assert.False(t, ok)
	})

	t.Run("with 255 characters in each of db, schema and name", func(t *testing.T) {
		ok := ValidObjectIdentifier(longSchemaObjectIdentifier)
		assert.True(t, ok)
	})
}

func TestValidObjectName(t *testing.T) {
	t.Run("with valid object name", func(t *testing.T) {
		ok := ValidObjectName("test_name")
		assert.True(t, ok)
	})

	t.Run("with empty name", func(t *testing.T) {
		ok := ValidObjectName("")
		assert.False(t, ok)
	})

	t.Run("with name over 255 characters", func(t *testing.T) {
		longName := random.AlphaN(256)
		ok := ValidObjectName(longName)
		assert.False(t, ok)
	})

	t.Run("with name exactly 255 characters", func(t *testing.T) {
		name := random.AlphaN(255)
		ok := ValidObjectName(name)
		assert.True(t, ok)
	})
}

func TestAnyValueSet(t *testing.T) {
	t.Run("with one value set", func(t *testing.T) {
		ok := anyValueSet(String("foo"))
		assert.True(t, ok)
	})

	t.Run("with no values", func(t *testing.T) {
		ok := anyValueSet()
		assert.False(t, ok)
	})

	t.Run("with multiple values set", func(t *testing.T) {
		ok := anyValueSet(String("foo"), String("bar"))
		assert.True(t, ok)
	})

	t.Run("with multiple values set and nil", func(t *testing.T) {
		ok := anyValueSet(String("foo"), nil, String("bar"))
		assert.True(t, ok)
	})
}

func TestExactlyOneValueSet(t *testing.T) {
	t.Run("with one value set", func(t *testing.T) {
		ok := exactlyOneValueSet(String("foo"))
		assert.True(t, ok)
	})

	t.Run("with no values", func(t *testing.T) {
		ok := exactlyOneValueSet()
		assert.False(t, ok)
	})

	t.Run("with multiple values set", func(t *testing.T) {
		ok := exactlyOneValueSet(String("foo"), String("bar"))
		assert.False(t, ok)
	})

	t.Run("with multiple values set and nil", func(t *testing.T) {
		ok := exactlyOneValueSet(String("foo"), nil, String("bar"))
		assert.False(t, ok)
	})
}

func TestEveryValueSet(t *testing.T) {
	t.Run("with one value set", func(t *testing.T) {
		ok := everyValueSet(String("foo"))
		assert.True(t, ok)
	})

	t.Run("with no values", func(t *testing.T) {
		ok := everyValueSet()
		assert.True(t, ok)
	})

	t.Run("with multiple values set", func(t *testing.T) {
		ok := everyValueSet(String("foo"), String("bar"))
		assert.True(t, ok)
	})

	t.Run("with multiple values set and nil", func(t *testing.T) {
		ok := everyValueSet(String("foo"), nil, String("bar"))
		assert.False(t, ok)
	})
}

func TestEveryValueNil(t *testing.T) {
	t.Run("with one value set", func(t *testing.T) {
		ok := everyValueNil(String("foo"))
		assert.False(t, ok)
	})

	t.Run("with no values", func(t *testing.T) {
		ok := everyValueNil()
		assert.True(t, ok)
	})

	t.Run("with multiple values set", func(t *testing.T) {
		ok := everyValueNil(String("foo"), String("bar"))
		assert.False(t, ok)
	})

	t.Run("with multiple values set and nil", func(t *testing.T) {
		ok := everyValueNil(String("foo"), nil, String("bar"))
		assert.False(t, ok)
	})
}

func TestValueSet(t *testing.T) {
	t.Run("with value set", func(t *testing.T) {
		ok := valueSet(String("foo"))
		assert.True(t, ok)
	})

	t.Run("with no value", func(t *testing.T) {
		ok := valueSet(nil)
		assert.False(t, ok)
	})

	t.Run("with valid identifier", func(t *testing.T) {
		ok := valueSet(NewAccountObjectIdentifier("foo"))
		assert.True(t, ok)
	})

	t.Run("with invalid identifier", func(t *testing.T) {
		ok := valueSet(emptyAccountObjectIdentifier)
		assert.False(t, ok)
	})

	t.Run("with zero ObjectType", func(t *testing.T) {
		s := struct {
			ot *ObjectType
		}{}
		ok := valueSet(s.ot)
		assert.False(t, ok)
	})

	t.Run("with invalid empty string", func(t *testing.T) {
		invalid := ""
		assert.False(t, valueSet(invalid))
	})

	t.Run("with valid non-empty string", func(t *testing.T) {
		valid := "non-empty string"
		assert.True(t, valueSet(valid))
	})
}

func TestValidateIntInRangeInclusive(t *testing.T) {
	t.Run("with value in range", func(t *testing.T) {
		ok := validateIntInRangeInclusive(5, 0, 10)
		assert.True(t, ok)
	})

	t.Run("with value out of range", func(t *testing.T) {
		ok := validateIntInRangeInclusive(5, 10, 20)
		assert.False(t, ok)
	})
}

func TestValidateIntGreaterThanOrEqual(t *testing.T) {
	t.Run("with value in range", func(t *testing.T) {
		ok := validateIntGreaterThanOrEqual(5, 0)
		assert.True(t, ok)
	})

	t.Run("with value out of range", func(t *testing.T) {
		ok := validateIntGreaterThanOrEqual(5, 10)
		assert.False(t, ok)
	})
}
