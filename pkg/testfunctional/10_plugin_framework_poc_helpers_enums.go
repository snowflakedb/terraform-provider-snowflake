// Content of this file should be moved to production files after proceeding with Terraform Plugin Framework.

package testfunctional

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func stringEnumAttributeCreate[T ~string](stringAttribute types.String, createField **T, mapper func(string) (T, error)) error {
	if !stringAttribute.IsNull() {
		v, err := mapper(stringAttribute.ValueString())
		if err != nil {
			return err
		}
		*createField = sdk.Pointer(v)
	}
	return nil
}

func stringEnumAttributeUpdate[T ~string](planned types.String, inState types.String, setField **T, unsetField **T, mapper func(string) (T, error)) error {
	if !planned.Equal(inState) {
		v, err := mapper(planned.ValueString())
		if err != nil {
			return err
		}
		if planned.IsNull() || planned.IsUnknown() {
			*unsetField = nil
		} else {
			*setField = sdk.Pointer(v)
		}
	}
	return nil
}
