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
		assert.Equal(t, ok, true)
	})

	t.Run("with invalid object identifier", func(t *testing.T) {
		ok := ValidObjectIdentifier(emptyAccountObjectIdentifier)
		assert.Equal(t, ok, false)
	})

	t.Run("over 255 characters", func(t *testing.T) {
		ok := ValidObjectIdentifier(invalidAccountObjectIdentifier)
		assert.Equal(t, ok, false)
	})

	t.Run("with 255 characters in each of db, schema and name", func(t *testing.T) {
		ok := ValidObjectIdentifier(longSchemaObjectIdentifier)
		assert.Equal(t, ok, true)
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
		assert.Equal(t, ok, true)
	})

	t.Run("with no values", func(t *testing.T) {
		ok := anyValueSet()
		assert.Equal(t, ok, false)
	})

	t.Run("with multiple values set", func(t *testing.T) {
		ok := anyValueSet(String("foo"), String("bar"))
		assert.Equal(t, ok, true)
	})

	t.Run("with multiple values set and nil", func(t *testing.T) {
		ok := anyValueSet(String("foo"), nil, String("bar"))
		assert.Equal(t, ok, true)
	})
}

func TestExactlyOneValueSet(t *testing.T) {
	t.Run("with one value set", func(t *testing.T) {
		ok := exactlyOneValueSet(String("foo"))
		assert.Equal(t, ok, true)
	})

	t.Run("with no values", func(t *testing.T) {
		ok := exactlyOneValueSet()
		assert.Equal(t, ok, false)
	})

	t.Run("with multiple values set", func(t *testing.T) {
		ok := exactlyOneValueSet(String("foo"), String("bar"))
		assert.Equal(t, ok, false)
	})

	t.Run("with multiple values set and nil", func(t *testing.T) {
		ok := exactlyOneValueSet(String("foo"), nil, String("bar"))
		assert.Equal(t, ok, false)
	})
}

func TestEveryValueSet(t *testing.T) {
	t.Run("with one value set", func(t *testing.T) {
		ok := everyValueSet(String("foo"))
		assert.Equal(t, ok, true)
	})

	t.Run("with no values", func(t *testing.T) {
		ok := everyValueSet()
		assert.Equal(t, ok, true)
	})

	t.Run("with multiple values set", func(t *testing.T) {
		ok := everyValueSet(String("foo"), String("bar"))
		assert.Equal(t, ok, true)
	})

	t.Run("with multiple values set and nil", func(t *testing.T) {
		ok := everyValueSet(String("foo"), nil, String("bar"))
		assert.Equal(t, ok, false)
	})
}

func TestEveryValueNil(t *testing.T) {
	t.Run("with one value set", func(t *testing.T) {
		ok := everyValueNil(String("foo"))
		assert.Equal(t, ok, false)
	})

	t.Run("with no values", func(t *testing.T) {
		ok := everyValueNil()
		assert.Equal(t, ok, true)
	})

	t.Run("with multiple values set", func(t *testing.T) {
		ok := everyValueNil(String("foo"), String("bar"))
		assert.Equal(t, ok, false)
	})

	t.Run("with multiple values set and nil", func(t *testing.T) {
		ok := everyValueNil(String("foo"), nil, String("bar"))
		assert.Equal(t, ok, false)
	})
}

func TestValueSet(t *testing.T) {
	t.Run("with value set", func(t *testing.T) {
		ok := valueSet(String("foo"))
		assert.Equal(t, ok, true)
	})

	t.Run("with no value", func(t *testing.T) {
		ok := valueSet(nil)
		assert.Equal(t, ok, false)
	})

	t.Run("with valid identifier", func(t *testing.T) {
		ok := valueSet(NewAccountObjectIdentifier("foo"))
		assert.Equal(t, ok, true)
	})

	t.Run("with invalid identifier", func(t *testing.T) {
		ok := valueSet(emptyAccountObjectIdentifier)
		assert.Equal(t, ok, false)
	})

	t.Run("with zero ObjectType", func(t *testing.T) {
		s := struct {
			ot *ObjectType
		}{}
		ok := valueSet(s.ot)
		assert.Equal(t, ok, false)
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
		assert.Equal(t, ok, true)
	})

	t.Run("with value out of range", func(t *testing.T) {
		ok := validateIntInRangeInclusive(5, 10, 20)
		assert.Equal(t, ok, false)
	})
}

func TestValidateIntGreaterThanOrEqual(t *testing.T) {
	t.Run("with value in range", func(t *testing.T) {
		ok := validateIntGreaterThanOrEqual(5, 0)
		assert.Equal(t, ok, true)
	})

	t.Run("with value out of range", func(t *testing.T) {
		ok := validateIntGreaterThanOrEqual(5, 10)
		assert.Equal(t, ok, false)
	})
}
