package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setStringIfNotEmpty(d *schema.ResourceData, key, value string) error {
	if value != "" {
		return d.Set(key, value)
	}
	return nil
}

func attributeMappedValueReadOrDefault[T, R any](d *schema.ResourceData, key string, value *T, mapper func(*T) (R, error), defaultValue *R) error {
	if value != nil {
		mappedValue, err := mapper(value)
		if err != nil {
			return err
		}
		return d.Set(key, mappedValue)
	}
	if defaultValue != nil {
		return d.Set(key, *defaultValue)
	}
	return d.Set(key, nil)
}

func setOptionalValueWithMapping[T, R any](d *schema.ResourceData, key string, value *T, mapper func(*T) R) error {
	if value != nil {
		return d.Set(key, mapper(value))
	}
	return d.Set(key, nil)
}

func setOptionalFromStringPtr(d *schema.ResourceData, key string, ptr *string) error {
	if ptr != nil {
		if err := d.Set(key, *ptr); err != nil {
			return err
		}
	}
	return nil
}

func setOptionalFromPtr[T any](d *schema.ResourceData, key string, ptr *T) error {
	if ptr != nil {
		if err := d.Set(key, *ptr); err != nil {
			return err
		}
	}
	return nil
}

// TODO [SNOW-1348103]: return error if nil
func setRequiredFromStringPtr(d *schema.ResourceData, key string, ptr *string) error {
	if ptr != nil {
		if err := d.Set(key, *ptr); err != nil {
			return err
		}
	}
	return nil
}

func optionalStringOutputMapping[T ~string](value *T) (any, string) {
	if value != nil {
		return *value, string(*value)
	}
	return nil, ""
}

func optionalBooleanStringOutputMapping(value *bool) (any, string) {
	if value != nil {
		return *value, booleanStringFromBool(*value)
	}
	return nil, BooleanDefault
}

func optionalIntOutputMapping[T ~int](value *T) any {
	if value != nil {
		return int(*value)
	}
	return 0
}

func optionalIntOutputMappingIntDefault[T ~int](value *T) any {
	if value != nil {
		return int(*value)
	}
	return IntDefault
}
