package sdk

import (
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func IsValidDataType(v string) bool {
	_, err := datatypes.ParseDataType(v)
	return err == nil
}

func ValidObjectName(name string) bool {
	// https://docs.snowflake.com/en/sql-reference/identifiers-syntax#double-quoted-identifiers
	l := len(name)
	if l == 0 || l > 255 {
		return false
	}
	return true
}

func ValidObjectIdentifier(objectIdentifier ObjectIdentifier) bool {
	return ValidObjectName(objectIdentifier.Name())
}

func anyValueSet(values ...interface{}) bool {
	for _, v := range values {
		if valueSet(v) {
			return true
		}
	}
	return false
}

func exactlyOneValueSet(values ...interface{}) bool {
	var count int
	for _, v := range values {
		if valueSet(v) {
			count++
		}
	}
	return count == 1
}

func moreThanOneValueSet(values ...interface{}) bool {
	var count int
	for _, v := range values {
		if valueSet(v) {
			count++
		}
	}
	return count > 1
}

func everyValueSet(values ...interface{}) bool {
	for _, v := range values {
		if !valueSet(v) {
			return false
		}
	}
	return true
}

func everyValueNil(values ...interface{}) bool {
	for _, v := range values {
		if valueSet(v) {
			return false
		}
	}
	return true
}

func valueSet(value interface{}) bool {
	if value == nil {
		return false
	}
	reflectedValue := reflect.ValueOf(value)
	if reflectedValue.Kind() == reflect.Pointer {
		reflectedValue = reflectedValue.Elem()
	}
	switch reflectedValue.Kind() {
	case reflect.Slice, reflect.String:
		return reflectedValue.Len() > 0
	case reflect.Invalid:
		return false
	case reflect.Struct:
		if _, ok := reflectedValue.Interface().(ObjectIdentifier); ok {
			return ValidObjectIdentifier(reflectedValue.Interface().(ObjectIdentifier))
		}
		return reflectedValue.Interface() != nil
	}
	return true
}

func validateIntInRangeInclusive(value int, min int, max int) bool {
	if value < min || value > max {
		return false
	}
	return true
}

func validateIntGreaterThan(value int, min int) bool {
	return value > min
}

func validateIntGreaterThanOrEqual(value int, min int) bool {
	return value >= min
}

// containsDoubleDollarQuotes reports whether the given value contains the `$$` sequence. It is used to reject user
// input for fields rendered with double dollar quoting, because Snowflake's dollar-quoted string constants are
// interpreted literally (no escaping) and an embedded `$$` would terminate the constant, enabling SQL injection.
func containsDoubleDollarQuotes(value string) bool {
	return strings.Contains(value, "$$")
}
