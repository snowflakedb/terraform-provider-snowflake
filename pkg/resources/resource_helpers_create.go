package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringAttributeCreate(d *schema.ResourceData, key string, createField **string) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.String(v.(string))
	}
	return nil
}

func stringAttributeCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(string) T) error {
	if v, ok := d.GetOk(key); ok {
		setValue(v.(string))
	}
	return nil
}

func intAttributeCreate(d *schema.ResourceData, key string, createField **int) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Int(v.(int))
	}
	return nil
}

func intAttributeCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(int) T) error {
	if v, ok := d.GetOk(key); ok {
		setValue(v.(int))
	}
	return nil
}

func intAttributeWithSpecialDefaultCreate(d *schema.ResourceData, key string, createField **int) error {
	if v := d.Get(key).(int); v != IntDefault {
		*createField = sdk.Int(v)
	}
	return nil
}

func intAttributeWithSpecialDefaultCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(int) T) error {
	if v := d.Get(key).(int); v != IntDefault {
		setValue(v)
	}
	return nil
}

func boolAttributeCreate(d *schema.ResourceData, key string, createField **bool) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Bool(v.(bool))
	}
	return nil
}

func boolAttributeCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(bool) T) error {
	if v, ok := d.GetOk(key); ok {
		setValue(v.(bool))
	}
	return nil
}

func booleanStringAttributeCreate(d *schema.ResourceData, key string, createField **bool) error {
	if v := d.Get(key).(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return err
		}
		*createField = sdk.Bool(parsed)
	}
	return nil
}

func booleanStringAttributeCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(bool) T) error {
	if v := d.Get(key).(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return err
		}
		setValue(parsed)
	}
	return nil
}

// booleanStringPairAttributeCreate reads a tri-state boolean string field and sets exactly one of the
// two given pointer fields: positiveField when the value is "true", negativeField when it is "false".
// Both are left untouched when the value is BooleanDefault. Useful for squashing a pair of mutually
// exclusive SQL keywords (e.g. ENFORCED / NOT ENFORCED) into a single schema field.
func booleanStringPairAttributeCreate(d *schema.ResourceData, key string, positiveField, negativeField **bool) error {
	v := d.Get(key).(string)
	if v == BooleanDefault {
		return nil
	}
	parsed, err := booleanStringToBool(v)
	if err != nil {
		return err
	}
	if parsed {
		*positiveField = new(true)
	} else {
		*negativeField = new(true)
	}
	return nil
}

func schemaObjectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.SchemaObjectIdentifier) error {
	if v, ok := d.GetOk(key); ok {
		objectIdentifier, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return err
		}
		*createField = sdk.Pointer(objectIdentifier)
	}
	return nil
}

func accountObjectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.AccountObjectIdentifier) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	return nil
}

func accountObjectIdentifierAttributeCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(sdk.AccountObjectIdentifier) T) error {
	if v, ok := d.GetOk(key); ok {
		setValue(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	return nil
}

func objectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.ObjectIdentifier) error {
	if v, ok := d.GetOk(key); ok {
		objectIdentifier, err := sdk.ParseObjectIdentifierString(v.(string))
		if err != nil {
			return err
		}
		*createField = sdk.Pointer(objectIdentifier)
	}
	return nil
}

func attributeDirectValueCreate[T any](d *schema.ResourceData, key string, createField **T, value *T) error {
	if _, ok := d.GetOk(key); ok {
		*createField = value
	}
	return nil
}

func attributeMappedValueCreate[T any](d *schema.ResourceData, key string, createField **T, mapper func(value any) (*T, error)) error {
	if v, ok := d.GetOk(key); ok {
		value, err := mapper(v)
		if err != nil {
			return err
		}
		*createField = value
	}
	return nil
}

func attributeMappedValueCreateBuilder[InputType any, MappedType any, RequestBuilder any](d *schema.ResourceData, key string, setValue func(MappedType) RequestBuilder, mapper func(value InputType) (MappedType, error)) error {
	if v, ok := d.GetOk(key); ok {
		value, err := mapper(v.(InputType))
		if err != nil {
			return err
		}
		setValue(value)
	}
	return nil
}

// attributeMappedValueCreateBuilderRawConfig is attributeMappedValueCreateBuilder that checks the raw config
// instead of GetOk, because GetOk cannot distinguish an absent list attribute from one explicitly set to an
// empty list (e.g. `null_if = []`).
func attributeMappedValueCreateBuilderRawConfig[InputType any, MappedType any, RequestBuilder any](d *schema.ResourceData, key string, setValue func(MappedType) RequestBuilder, mapper func(value InputType) (MappedType, error)) error {
	if !d.GetRawConfig().AsValueMap()[key].IsNull() {
		value, err := mapper(d.Get(key).(InputType))
		if err != nil {
			return err
		}
		setValue(value)
	}
	return nil
}

func attributeMappedValueCreateBuilderNested[MappedType any, RequestBuilder any](d *schema.ResourceData, key string, setValue func(MappedType) RequestBuilder, mapper func(*schema.ResourceData) (MappedType, error)) error {
	if _, ok := d.GetOk(key); ok {
		value, err := mapper(d)
		if err != nil {
			return err
		}
		setValue(value)
	}
	return nil
}

// singleBlockAttributeCreate parses a MaxItems=1 list block when it is set.
// parse receives the nested prefix "<key>.0." and reads fields with the existing attribute helpers.
func singleBlockAttributeCreate[Value any, RequestBuilder any](d *schema.ResourceData, key string, setValue func(Value) RequestBuilder, parse func(*schema.ResourceData, string) (Value, error)) error {
	if _, ok := d.GetOk(key); !ok {
		return nil
	}
	value, err := parse(d, key+".0.")
	if err != nil {
		return err
	}
	setValue(value)
	return nil
}

// listBlockAttributeCreate parses a TypeList of nested blocks into createField.
// parse receives the nested prefix "<key>.<index>." and reads fields with the existing attribute helpers.
// An empty list leaves createField untouched.
func listBlockAttributeCreate[T any](d *schema.ResourceData, key string, createField *[]T, mapper func(*schema.ResourceData, string) (T, error)) error {
	count := d.Get(key + ".#").(int)
	if count == 0 {
		return nil
	}
	values := make([]T, count)
	for i := range count {
		value, err := mapper(d, fmt.Sprintf("%s.%d.", key, i))
		if err != nil {
			return err
		}
		values[i] = value
	}
	*createField = values
	return nil
}

// setValueCreate maps each element of a schema.Set attribute and assigns the mapped slice when the
// attribute is set. An absent or empty set leaves createField untouched.
func setValueCreate[T any](d *schema.ResourceData, key string, createField *[]T, mapper func(any) (T, error)) error {
	if v, ok := d.GetOk(key); ok {
		mappedValue, err := collections.MapErr(v.(*schema.Set).List(), mapper)
		if err != nil {
			return err
		}
		*createField = mappedValue
	}
	return nil
}

// listValueCreate maps each element of a TypeList attribute and assigns the mapped slice when the
// attribute is set. An absent list leaves createField untouched.
func listValueCreate[T any](d *schema.ResourceData, key string, createField *[]T, mapper func(any) (T, error)) error {
	if v, ok := d.GetOk(key); ok {
		mappedValue, err := collections.MapErr(v.([]any), mapper)
		if err != nil {
			return err
		}
		*createField = mappedValue
	}
	return nil
}

// setMappedValueCreateBuilder maps each element of a schema.Set attribute and passes the mapped
// slice to setValue when the attribute is set. An absent or empty set leaves the request unchanged.
func setMappedValueCreateBuilder[T any, RequestBuilder any](d *schema.ResourceData, key string, setValue func([]T) RequestBuilder, mapper func(any) (T, error)) error {
	if v, ok := d.GetOk(key); ok {
		mappedValue, err := collections.MapErr(v.(*schema.Set).List(), mapper)
		if err != nil {
			return err
		}
		setValue(mappedValue)
	}
	return nil
}

// listObjects maps a list of nested objects. raw must be a []any whose elements are map[string]any.
func listObjects[T any](raw any, mapper func(map[string]any) (T, error)) ([]T, error) {
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("expected list of objects, got %T", raw)
	}
	return collections.MapErr(list, func(item any) (T, error) {
		block, ok := item.(map[string]any)
		if !ok {
			var zero T
			return zero, fmt.Errorf("expected object, got %T", item)
		}
		return mapper(block)
	})
}

func copyGrantsAttributeCreate(d *schema.ResourceData, isOrReplace bool, orReplaceField, copyGrantsField **bool) error {
	if isOrReplace {
		*orReplaceField = sdk.Bool(true)
		if d.GetRawConfig().AsValueMap()["copy_grants"].True() {
			*copyGrantsField = sdk.Bool(true)
		}
	}
	return nil
}
