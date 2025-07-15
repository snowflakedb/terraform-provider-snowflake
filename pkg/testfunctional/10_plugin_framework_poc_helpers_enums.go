// Content of this file should be moved to production files after proceeding with Terraform Plugin Framework.

package testfunctional

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customtypes"
)

func sameAfterNormalization[T ~string](oldValue string, newValue string, normalize func(string) (T, error)) (bool, error) {
	oldNormalized, err := normalize(oldValue)
	if err != nil {
		return false, err
	}
	newNormalized, err := normalize(newValue)
	if err != nil {
		return false, err
	}

	return oldNormalized == newNormalized, nil
}

func stringEnumAttributeCreate[T customtypes.EnumCreator[T]](attr customtypes.EnumValue[T], createField **T, mapper func(string) (T, error)) error {
	if !attr.IsNull() {
		v, err := mapper(attr.ValueString())
		if err != nil {
			return err
		}
		*createField = sdk.Pointer(v)
	}
	return nil
}

func stringEnumAttributeUpdate[T customtypes.EnumCreator[T]](planned customtypes.EnumValue[T], inState customtypes.EnumValue[T], setField **T, unsetField **T, mapper func(string) (T, error)) error {
	if !planned.Equal(inState) {
		if planned.IsNull() || planned.IsUnknown() {
			*unsetField = nil
		} else {
			v, err := mapper(planned.ValueString())
			if err != nil {
				return err
			}
			*setField = sdk.Pointer(v)
		}
	}
	return nil
}
