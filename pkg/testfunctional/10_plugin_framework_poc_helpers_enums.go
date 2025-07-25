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

func StringEnumAttributeUpdate[T customtypes.EnumCreator[T]](planValue customtypes.EnumValue[T], stateValue customtypes.EnumValue[T], setField **T, unsetField **T) error {
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
