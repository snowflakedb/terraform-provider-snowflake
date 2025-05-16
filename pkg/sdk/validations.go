package sdk

import (
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

func IsValidDataType(v string) bool {
	_, err := datatypes.ParseDataType(v)
	return err == nil
}

func ValidObjectIdentifier(objectIdentifier ObjectIdentifier) bool {
	// https://docs.snowflake.com/en/sql-reference/identifiers-syntax#double-quoted-identifiers
	l := len(objectIdentifier.Name())
	if l == 0 || l > 255 {
		return false
	}
	return true
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
	if reflectedValue.Kind() == reflect.Ptr {
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
