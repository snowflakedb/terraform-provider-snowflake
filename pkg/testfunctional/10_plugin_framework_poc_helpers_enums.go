// Content of this file should be moved to production files after proceeding with Terraform Plugin Framework.

package testfunctional

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customtypes"
)

func StringEnumAttributeCreate[T customtypes.EnumCreator[T]](attr customtypes.EnumValue[T], createField **T) error {
	if !attr.IsNull() {
		v, err := attr.Normalize()
		if err != nil {
			return err
		}
		*createField = sdk.Pointer(v)
	}
	return nil
}

func stringEnumAttributeUpdate[T customtypes.EnumCreator[T]](planValue customtypes.EnumValue[T], stateValue customtypes.EnumValue[T], setField **T, unsetField **T) error {
	// currently Equal is enough as we have customplanmodifiers.EnumSuppressor which checks normalized equality for planValue and stateValue
	if !planValue.Equal(stateValue) {
		if planValue.IsNull() || planValue.IsUnknown() {
			*unsetField = nil
		} else {
			v, err := planValue.Normalize()
			if err != nil {
				return err
			}
			*setField = sdk.Pointer(v)
		}
	}
	return nil
}

// TODO [SNOW-2296350]: add functional test for this variant (with unset) to be closer to our implementation
func StringEnumAttributeUpdate[T customtypes.EnumCreator[T]](planValue customtypes.EnumValue[T], stateValue customtypes.EnumValue[T], setField **T, unsetField **bool) error {
	// currently Equal is enough as we have customplanmodifiers.EnumSuppressor which checks normalized equality for planValue and stateValue
	if !planValue.Equal(stateValue) {
		if planValue.IsNull() || planValue.IsUnknown() {
			*unsetField = sdk.Bool(true)
		} else {
			v, err := planValue.Normalize()
			if err != nil {
				return err
			}
			*setField = sdk.Pointer(v)
		}
	}
	return nil
}

// TODO [SNOW-2296391]: add functional test for this variant
func StringEnumAttributeUpdateSetDefaultInsteadOfUnset[T customtypes.EnumCreator[T]](planValue customtypes.EnumValue[T], stateValue customtypes.EnumValue[T], setField **T, defaultValue T) error {
	// currently Equal is enough as we have customplanmodifiers.EnumSuppressor which checks normalized equality for planValue and stateValue
	if !planValue.Equal(stateValue) {
		if planValue.IsNull() || planValue.IsUnknown() {
			*setField = sdk.Pointer(defaultValue)
		} else {
			v, err := planValue.Normalize()
			if err != nil {
				return err
			}
			*setField = sdk.Pointer(v)
		}
	}
	return nil
}

// TODO [SNOW-2296391]: add functional test for this variant
func StringEnumAttributeUpdateSetOnly[T customtypes.EnumCreator[T]](planValue customtypes.EnumValue[T], stateValue customtypes.EnumValue[T], setField **T) error {
	// currently Equal is enough as we have customplanmodifiers.EnumSuppressor which checks normalized equality for planValue and stateValue
	if !planValue.Equal(stateValue) {
		if planValue.IsNull() || planValue.IsUnknown() {
			// no-op, should be handled by force new
		} else {
			v, err := planValue.Normalize()
			if err != nil {
				return err
			}
			*setField = sdk.Pointer(v)
		}
	}
	return nil
}
