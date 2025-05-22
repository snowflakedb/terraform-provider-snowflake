package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringAttributeUpdate(d *schema.ResourceData, key string, setField **string, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = sdk.String(v.(string))
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func intAttributeUpdate(d *schema.ResourceData, key string, setField **int, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = sdk.Int(v.(int))
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func intAttributeWithSpecialDefaultUpdate(d *schema.ResourceData, key string, setField **int, unsetField **bool) error {
	if d.HasChange(key) {
		if v := d.Get(key).(int); v != IntDefault {
			*setField = sdk.Int(v)
		} else {
			*unsetField = sdk.Bool(true)
		}
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
			*setField = sdk.Bool(parsed)
		} else {
			*unsetField = sdk.Bool(true)
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
			*setField = sdk.Bool(parsed)
		} else {
			*setField = sdk.Bool(fallbackValue)
		}
	}
	return nil
}

func accountObjectIdentifierAttributeUpdate(d *schema.ResourceData, key string, setField **sdk.AccountObjectIdentifier, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
		} else {
			*unsetField = sdk.Bool(true)
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
			*setField = sdk.Pointer(objectIdentifier)
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func attributeDirectValueUpdate[T any](d *schema.ResourceData, key string, setField **T, value *T, unsetField **bool) error {
	if d.HasChange(key) {
		if _, ok := d.GetOk(key); ok {
			*setField = value
		} else {
			*unsetField = sdk.Bool(true)
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
			*setField = sdk.Pointer(mappedValue)
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}
