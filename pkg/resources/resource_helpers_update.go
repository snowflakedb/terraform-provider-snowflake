package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func stringAttributeUpdateBuilder[Set any, Unset any](d *schema.ResourceData, key string, setValue func(string) Set, unsetValue func(bool) Unset) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			setValue(v.(string))
		} else {
			unsetValue(true)
		}
	}
	return nil
}

func stringAttributeUpdateSetOnlyNotEmpty(d *schema.ResourceData, key string, setField **string) error {
	if d.HasChange(key) {
		*setField = sdk.String(d.Get(key).(string))
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
			*setField = sdk.Int(v.(int))
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func intAttributeUpdateSetOnly(d *schema.ResourceData, key string, setField **int) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = sdk.Int(v.(int))
		}
	}
	return nil
}

func intAttributeUpdateBuilder[Set any, Unset any](d *schema.ResourceData, key string, setValue func(int) Set, unsetValue func(bool) Unset) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			setValue(v.(int))
		} else {
			unsetValue(true)
		}
	}
	return nil
}

func intAttributeUpdateSetOnlyBuilder[Set any](d *schema.ResourceData, key string, setValue func(int) Set) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			setValue(v.(int))
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

func intAttributeWithSpecialDefaultUpdateBuilder[Set any, Unset any](d *schema.ResourceData, key string, setValue func(int) Set, unsetValue func(bool) Unset) error {
	if d.HasChange(key) {
		if v := d.Get(key).(int); v != IntDefault {
			setValue(v)
		} else {
			unsetValue(true)
		}
	}
	return nil
}

// intAttributeWithSpecialDefaultUnsetFallbackUpdate handles an int attribute whose "not set" sentinel
// is IntDefault (rather than the zero value). Removing it from config re-applies fallbackValue via SET
// instead of unsetting, for attributes where UNSET is unavailable or misbehaves (e.g. auto_suspend).
func intAttributeWithSpecialDefaultUnsetFallbackUpdate(d *schema.ResourceData, key string, setField **int, fallbackValue int) error {
	if d.HasChange(key) {
		if v := d.Get(key).(int); v != IntDefault {
			*setField = sdk.Int(v)
		} else {
			*setField = sdk.Int(fallbackValue)
		}
	}
	return nil
}

func intAttributeUnsetFallbackUpdateWithZeroDefault(d *schema.ResourceData, key string, setField **int, fallbackValue int) error {
	if d.HasChange(key) {
		if v := d.Get(key).(int); v != 0 {
			*setField = sdk.Int(v)
		} else {
			*setField = sdk.Int(fallbackValue)
		}
	}
	return nil
}

func booleanAttributeUpdate(d *schema.ResourceData, key string, setField **bool, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = sdk.Bool(v.(bool))
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func booleanAttributeUpdateSetOnly(d *schema.ResourceData, key string, setField **bool) error {
	if d.HasChange(key) {
		*setField = sdk.Bool(d.Get(key).(bool))
	}
	return nil
}

func booleanAttributeUpdateSetOnlyBuilder[Set any](d *schema.ResourceData, key string, setValue func(bool) Set) error {
	if d.HasChange(key) {
		setValue(d.Get(key).(bool))
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

func booleanStringAttributeUpdateBuilder[Set any, Unset any](d *schema.ResourceData, key string, setValue func(bool) Set, unsetValue func(bool) Unset) error {
	if d.HasChange(key) {
		if v := d.Get(key).(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return err
			}
			setValue(parsed)
		} else {
			unsetValue(true)
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
			*setField = sdk.Bool(parsed)
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
			*setField = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
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

func schemaObjectIdentifierAttributeUpdate(d *schema.ResourceData, key string, setField **sdk.SchemaObjectIdentifier, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			id, err := sdk.ParseSchemaObjectIdentifier(v.(string))
			if err != nil {
				return err
			}
			*setField = sdk.Pointer(id)
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

// singleBlockAttributeUpdate unsets a MaxItems=1 list block when it is removed.
// While the block is present, fill applies the existing attribute update helpers to nested fields at "<key>.0.".
// The block is set only when fill changes the request.
func singleBlockAttributeUpdate[T any](d *schema.ResourceData, key string, setField **T, unsetField **bool, fill func(*schema.ResourceData, string, *T) error) error {
	if !d.HasChange(key) {
		return nil
	}
	if _, ok := d.GetOk(key); !ok {
		*unsetField = new(true)
		return nil
	}
	req := new(T)
	if err := fill(d, key+".0.", req); err != nil {
		return err
	}
	if !reflect.DeepEqual(req, new(T)) {
		*setField = req
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
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

// setValueUpdateFallback reuses setValueUpdate. Clearing a set assigns fallback
// (for example ALL) instead of unsetting the attribute.
func setValueUpdateFallback[T any](d *schema.ResourceData, key string, setField *[]T, mapper func(any) (T, error), fallback []T) error {
	var unset *bool
	if err := setValueUpdate(d, key, setField, &unset, mapper); err != nil {
		return err
	}
	if unset != nil {
		*setField = fallback
	}
	return nil
}

func listValueUpdateSetOnlyBuilder[T any, Set any](d *schema.ResourceData, key string, setValue func([]T) Set, mapper func(any) (T, error)) error {
	if !d.HasChange(key) {
		return nil
	}
	mappedValue, err := collections.MapErr(d.Get(key).([]any), mapper)
	if err != nil {
		return err
	}
	if len(mappedValue) == 0 {
		return fmt.Errorf("cannot unset %s", key)
	}
	setValue(mappedValue)
	return nil
}

func listValueUpdateBuilder[T any, Set any, Unset any](d *schema.ResourceData, key string, setValue func([]T) Set, unsetValue func(bool) Unset, mapper func(any) (T, error)) error {
	if d.HasChange(key) {
		v := d.Get(key)
		mappedValue, err := collections.MapErr(v.([]any), mapper)
		if err != nil {
			return err
		}

		if len(mappedValue) > 0 {
			setValue(mappedValue)
		} else {
			unsetValue(true)
		}
	}
	return nil
}

// attributeMappedValueUpdateSetOnlyAlways maps the current value whenever the attribute changed,
// including values GetOk treats as empty. There is no unset.
func attributeMappedValueUpdateSetOnlyAlways[T, R any](d *schema.ResourceData, key string, setField **R, mapper func(T) (R, error)) error {
	if d.HasChange(key) {
		mappedValue, err := mapper(d.Get(key).(T))
		if err != nil {
			return err
		}
		*setField = sdk.Pointer(mappedValue)
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
			*setField = sdk.Pointer(mappedValue)
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
			*setField = sdk.Pointer(mappedValue)
		} else {
			*setField = sdk.Pointer(fallbackValue)
		}
	}
	return nil
}

// attributeMappedValueUpdateSetOnlyRawConfigFallback is attributeMappedValueUpdateSetOnlyFallback that checks
// the raw config instead of GetOk, because GetOk cannot distinguish an absent list attribute from one
// explicitly set to an empty list (e.g. `null_if = []`).
func attributeMappedValueUpdateSetOnlyRawConfigFallback[T, R any](d *schema.ResourceData, key string, setField **R, mapper func(T) (R, error), fallbackValue R) error {
	if d.HasChange(key) {
		if !d.GetRawConfig().AsValueMap()[key].IsNull() {
			mappedValue, err := mapper(d.Get(key).(T))
			if err != nil {
				return err
			}
			*setField = sdk.Pointer(mappedValue)
		} else {
			*setField = sdk.Pointer(fallbackValue)
		}
	}
	return nil
}

// attributeMappedValueUpdateSetOnlySliceFallback is attributeMappedValueUpdateSetOnlyFallback for set
// fields typed as a plain slice (e.g. []sdk.NullString) rather than a pointer to a value.
// Unlike attributeMappedValueUpdateSetOnlyFallback, it checks the raw config instead of GetOk,
// because GetOk cannot distinguish an absent list attribute from one explicitly set to an empty list.
func attributeMappedValueUpdateSetOnlySliceFallback[T, E any](d *schema.ResourceData, key string, setField *[]E, mapper func(T) ([]E, error), fallbackValue []E) error {
	if d.HasChange(key) {
		if !d.GetRawConfig().AsValueMap()[key].IsNull() {
			mappedValue, err := mapper(d.Get(key).(T))
			if err != nil {
				return err
			}
			*setField = mappedValue
		} else {
			*setField = fallbackValue
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
			*setField = sdk.Pointer(mappedValue)
		} else {
			*setField = sdk.Pointer(fallbackValue)
		}
	}
	return nil
}

// renameAccountLevelResource renames an account-level object when name changed.
// The new identifier is the name attribute.
// On success the Terraform ID is set with EncodeResourceIdentifier and the new identifier is returned.
// When name is unchanged, id is returned unchanged.
func renameAccountLevelResource(
	ctx context.Context,
	d *schema.ResourceData,
	id sdk.AccountObjectIdentifier,
	rename func(context.Context, sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) error,
) (sdk.AccountObjectIdentifier, error) {
	if !d.HasChange("name") {
		return id, nil
	}
	target := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
	if err := rename(ctx, id, target); err != nil {
		return id, err
	}
	d.SetId(helpers.EncodeResourceIdentifier(target))
	return target, nil
}

func attributeMappedValueUpdateIf[T, R any](d *schema.ResourceData, key string, setField **R, unsetField **bool, condition func(T) bool, mapper func(T) (R, error)) error {
	if d.HasChange(key) {
		v := d.Get(key)
		if condition(v.(T)) {
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

// removeThenAddWithPlaceholder adds a temporary item, applies removeThenAdd, and removes the
// temporary item even when the update fails. Use this when every existing item is replaced
// and the object requires at least one item at all times.
func removeThenAddWithPlaceholder[T any](placeholder T, removed, added []T, remove, add func(T) error) (err error) {
	if err := add(placeholder); err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, remove(placeholder))
	}()

	return removeThenAdd(removed, added, remove, add)
}

// removeThenAdd removes items before adding replacements, so a replacement can reuse a name
// that must stay unique. len(removed) must leave at least one existing item when the object
// requires a minimum of one; otherwise use removeThenAddWithPlaceholder.
func removeThenAdd[T any](removed, added []T, remove, add func(T) error) error {
	for _, item := range removed {
		if err := remove(item); err != nil {
			return err
		}
	}
	for _, item := range added {
		if err := add(item); err != nil {
			return err
		}
	}
	return nil
}
