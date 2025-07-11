package customtypes

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func StringWithBackingFieldAttributeCreate(v StringWithBackingFieldValue, createField **string) {
	if !v.IsNull() {
		*createField = sdk.String(v.ValueString())
	}
}

func StringWithBackingFieldAttributeUpdate(planned StringWithBackingFieldValue, inState StringWithBackingFieldValue, setField **string, unsetField **string) {
	if !planned.Equal(inState) {
		if planned.IsNull() {
			*unsetField = nil
		} else {
			*setField = planned.ValueStringPointer()
		}
	}
}
