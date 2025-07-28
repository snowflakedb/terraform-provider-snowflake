// Content of this file should be moved to production files after proceeding with Terraform Plugin Framework.

package testfunctional

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func BooleanAttributeCreate(boolAttribute types.Bool, createField **bool) error {
	if !boolAttribute.IsNull() {
		*createField = boolAttribute.ValueBoolPointer()
	}
	return nil
}

func booleanAttributeUpdate(planned types.Bool, inState types.Bool, setField **bool, unsetField **bool) {
	if !planned.Equal(inState) {
		if planned.IsNull() {
			*unsetField = nil
		} else {
			*setField = planned.ValueBoolPointer()
		}
	}
}

// TODO [mux-PR]: add functional test for this variant (with unset) to be closer to our implementation
func BooleanAttributeUpdate(planned types.Bool, inState types.Bool, setField **bool, unsetField **bool) error {
	if !planned.Equal(inState) {
		if planned.IsNull() {
			*unsetField = sdk.Bool(true)
		} else {
			*setField = planned.ValueBoolPointer()
		}
	}
	return nil
}

// TODO [mux-PR]: add functional test for this variant
func BooleanAttributeUpdateSetDefaultInsteadOfUnset(planned types.Bool, inState types.Bool, setField **bool, defaultValue bool) error {
	if !planned.Equal(inState) {
		if planned.IsNull() {
			*setField = sdk.Bool(defaultValue)
		} else {
			*setField = planned.ValueBoolPointer()
		}
	}
	return nil
}

func Int64AttributeCreate(int64Attribute types.Int64, createField **int) error {
	if !int64Attribute.IsNull() {
		*createField = sdk.Int(int(int64Attribute.ValueInt64()))
	}
	return nil
}

// For now, we use here two same set/unset pointers as the test server handles a single HTTP call.
// It should be altered when working on the server improvement.
// TODO [mux-PRs]: Handle set/unset instead just single opts
func int64AttributeUpdate(planned types.Int64, inState types.Int64, setField **int, unsetField **int) {
	if !planned.Equal(inState) {
		if planned.IsNull() {
			*unsetField = nil
		} else {
			*setField = sdk.Int(int(planned.ValueInt64()))
		}
	}
}

// TODO [mux-PR]: add functional test for this variant (with unset) to be closer to our implementation
func Int64AttributeUpdate(planned types.Int64, inState types.Int64, setField **int, unsetField **bool) error {
	if !planned.Equal(inState) {
		if planned.IsNull() {
			*unsetField = sdk.Bool(true)
		} else {
			*setField = sdk.Int(int(planned.ValueInt64()))
		}
	}
	return nil
}

// TODO [mux-PR]: add functional test for this variant
func Int64AttributeUpdateSetDefaultInsteadOfUnset(planned types.Int64, inState types.Int64, setField **int, defaultValue int) error {
	if !planned.Equal(inState) {
		if planned.IsNull() {
			*setField = sdk.Int(defaultValue)
		} else {
			*setField = sdk.Int(int(planned.ValueInt64()))
		}
	}
	return nil
}

func StringAttributeCreate(stringAttribute types.String, createField **string) error {
	if !stringAttribute.IsNull() {
		*createField = stringAttribute.ValueStringPointer()
	}
	return nil
}

func stringAttributeUpdate(planned types.String, inState types.String, setField **string, unsetField **string) {
	if !planned.Equal(inState) {
		if planned.IsNull() || planned.IsUnknown() {
			*unsetField = nil
		} else {
			*setField = planned.ValueStringPointer()
		}
	}
}

// TODO [mux-PR]: add functional test for this variant (with unset) to be closer to our implementation
func StringAttributeUpdate(planned types.String, inState types.String, setField **string, unsetField **bool) error {
	if !planned.Equal(inState) {
		if planned.IsNull() || planned.IsUnknown() {
			*unsetField = sdk.Bool(true)
		} else {
			*setField = planned.ValueStringPointer()
		}
	}
	return nil
}
