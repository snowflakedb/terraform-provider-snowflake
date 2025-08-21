package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setFromIntProperty(d *schema.ResourceData, key string, property *sdk.IntProperty) error {
	if property != nil && property.Value != nil {
		if err := d.Set(key, *property.Value); err != nil {
			return err
		}
	}
	return nil
}

func setFromStringProperty(d *schema.ResourceData, key string, property *sdk.StringProperty) error {
	if property != nil {
		if err := d.Set(key, property.Value); err != nil {
			return err
		}
	}
	return nil
}

func setFromStringPropertyIfNotEmpty(d *schema.ResourceData, key string, property *sdk.StringProperty) error {
	if property != nil && property.Value != "" {
		if err := d.Set(key, property.Value); err != nil {
			return err
		}
	}
	return nil
}

func setFromBoolProperty(d *schema.ResourceData, key string, property *sdk.BoolProperty) error {
	if property != nil {
		if err := d.Set(key, property.Value); err != nil {
			return err
		}
	}
	return nil
}

func setBooleanStringFromBoolProperty(d *schema.ResourceData, key string, property *sdk.BoolProperty) error {
	if property != nil {
		if err := d.Set(key, booleanStringFromBool(property.Value)); err != nil {
			return err
		}
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

// TODO [SNOW-1348103]: return error if nil
func setRequiredFromStringPtr(d *schema.ResourceData, key string, ptr *string) error {
	if ptr != nil {
		if err := d.Set(key, *ptr); err != nil {
			return err
		}
	}
	return nil
}
