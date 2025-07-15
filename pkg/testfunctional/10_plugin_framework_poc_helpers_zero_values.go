// Content of this file should be moved to production files after proceeding with Terraform Plugin Framework.

package testfunctional

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func booleanAttributeCreate(boolAttribute types.Bool, createField **bool) {
	if !boolAttribute.IsNull() {
		*createField = boolAttribute.ValueBoolPointer()
	}
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

func int64AttributeCreate(int64Attribute types.Int64, createField **int) {
	if !int64Attribute.IsNull() {
		*createField = sdk.Int(int(int64Attribute.ValueInt64()))
	}
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

func stringAttributeCreate(stringAttribute types.String, createField **string) {
	if !stringAttribute.IsNull() {
		*createField = stringAttribute.ValueStringPointer()
	}
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
