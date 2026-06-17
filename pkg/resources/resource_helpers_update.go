package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringAttributeUpdate(d *schema.ResourceData, key string, setField **string, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = new(v.(string))
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func stringAttributeUpdateSetOnlyNotEmpty(d *schema.ResourceData, key string, setField **string) error {
	if d.HasChange(key) {
		*setField = new(d.Get(key).(string))
	}
	return nil
}

func stringAttributeUpdateSetOnly(d *schema.ResourceData, key string, setField **sdk.StringAllowEmpty) error {
	if d.HasChange(key) {
		*setField = &sdk.StringAllowEmpty{Value: d.Get(key).(string)}
	}
	return nil
}

func intAttributeUpdate(d *schema.ResourceData, key string, setField **int, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = new(v.(int))
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func intAttributeUpdateSetOnly(d *schema.ResourceData, key string, setField **int) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = new(v.(int))
		}
	}
	return nil
}

func intAttributeWithSpecialDefaultUpdate(d *schema.ResourceData, key string, setField **int, unsetField **bool) error {
	if d.HasChange(key) {
		if v := d.Get(key).(int); v != IntDefault {
			*setField = new(v)
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func intAttributeUnsetFallbackUpdateWithZeroDefault(d *schema.ResourceData, key string, setField **int, fallbackValue int) error {
	if d.HasChange(key) {
		if v := d.Get(key).(int); v != 0 {
			*setField = new(v)
		} else {
			*setField = new(fallbackValue)
		}
	}
	return nil
}

func booleanAttributeUpdate(d *schema.ResourceData, key string, setField **bool, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = new(v.(bool))
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func booleanAttributeUpdateSetOnly(d *schema.ResourceData, key string, setField **bool) error {
	if d.HasChange(key) {
		*setField = new(d.Get(key).(bool))
	}
	return nil
}

func booleanStringAttributeUpdate(d *schema.ResourceData, key string, setField **bool, unsetField **bool) error {
	if d.HasChange(key) {
		if v := d.Get(key).(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return err
			}
			*setField = new(parsed)
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func booleanStringAttributeUpdateSetOnly(d *schema.ResourceData, key string, setField **bool) error {
	if d.HasChange(key) {
		if v := d.Get(key).(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return err
			}
			*setField = new(parsed)
		}
	}
	return nil
}

func booleanStringAttributeUnsetFallbackUpdate(d *schema.ResourceData, key string, setField **bool, fallbackValue bool) error {
	if d.HasChange(key) {
		if v := d.Get(key).(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return err
			}
			*setField = new(parsed)
		} else {
			*setField = new(fallbackValue)
		}
	}
	return nil
}

func booleanStringAttributeUnsetFallbackUpdateBuilder[T any](d *schema.ResourceData, key string, setValue func(bool) T, fallbackValue bool) error {
	if d.HasChange(key) {
		if v := d.Get(key).(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return err
			}
			setValue(parsed)
		} else {
			setValue(fallbackValue)
		}
	}
	return nil
}

func accountObjectIdentifierAttributeSetOnly(d *schema.ResourceData, key string, setField **sdk.AccountObjectIdentifier) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = new(sdk.NewAccountObjectIdentifier(v.(string)))
		}
	}
	return nil
}

func accountObjectIdentifierAttributeUpdate(d *schema.ResourceData, key string, setField **sdk.AccountObjectIdentifier, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = new(sdk.NewAccountObjectIdentifier(v.(string)))
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func schemaObjectIdentifierAttributeUpdate(d *schema.ResourceData, key string, setField **sdk.SchemaObjectIdentifier, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			id, err := sdk.ParseSchemaObjectIdentifier(v.(string))
			if err != nil {
				return err
			}
			*setField = new(id)
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func objectIdentifierAttributeUpdate(d *schema.ResourceData, key string, setField **sdk.ObjectIdentifier, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			objectIdentifier, err := sdk.ParseObjectIdentifierString(v.(string))
			if err != nil {
				return err
			}
			*setField = new(objectIdentifier)
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func attributeDirectValueUpdate[T any](d *schema.ResourceData, key string, setField **T, value *T, unsetField **bool) error {
	if d.HasChange(key) {
		if _, ok := d.GetOk(key); ok {
			*setField = value
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func attributeMappedValueUpdate[T, R any](d *schema.ResourceData, key string, setField **R, unsetField **bool, mapper func(T) (R, error)) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			mappedValue, err := mapper(v.(T))
			if err != nil {
				return err
			}
			*setField = new(mappedValue)
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func setValueUpdate[T any](d *schema.ResourceData, key string, setField *[]T, unsetField **bool, mapper func(any) (T, error)) error {
	if d.HasChange(key) {
		v := d.Get(key)
		mappedValue, err := collections.MapErr(v.(*schema.Set).List(), mapper)
		if err != nil {
			return err
		}

		if len(mappedValue) > 0 {
			*setField = mappedValue
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}

func attributeMappedValueUpdateSetOnly[T, R any](d *schema.ResourceData, key string, setField **R, mapper func(T) (R, error)) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			mappedValue, err := mapper(v.(T))
			if err != nil {
				return err
			}
			*setField = new(mappedValue)
		}
	}
	return nil
}

func attributeMappedValueUpdateSetOnlyFallback[T, R any](d *schema.ResourceData, key string, setField **R, mapper func(T) (R, error), fallbackValue R) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			mappedValue, err := mapper(v.(T))
			if err != nil {
				return err
			}
			*setField = new(mappedValue)
		} else {
			*setField = new(fallbackValue)
		}
	}
	return nil
}

func attributeMappedValueUpdateSetOnlyFallbackNested[R any](d *schema.ResourceData, key string, setField **R, mapper func(*schema.ResourceData) (R, error), fallbackValue R) error {
	if d.HasChange(key) {
		if _, ok := d.GetOk(key); ok {
			mappedValue, err := mapper(d)
			if err != nil {
				return err
			}
			*setField = new(mappedValue)
		} else {
			*setField = new(fallbackValue)
		}
	}
	return nil
}

func attributeMappedValueUpdateIf[T, R any](d *schema.ResourceData, key string, setField **R, unsetField **bool, condition func(T) bool, mapper func(T) (R, error)) error {
	if d.HasChange(key) {
		v := d.Get(key)
		if condition(v.(T)) {
			mappedValue, err := mapper(v.(T))
			if err != nil {
				return err
			}
			*setField = new(mappedValue)
		} else {
			*unsetField = new(true)
		}
	}
	return nil
}
