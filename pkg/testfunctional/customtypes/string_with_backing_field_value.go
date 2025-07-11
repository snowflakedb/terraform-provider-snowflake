package customtypes

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ attr.Value = StringWithBackingFieldValue{}

type StringWithBackingFieldValue struct {
	basetypes.StringValue

	InternalValue string // we do not expose it to practitioner in schema
}

func (v StringWithBackingFieldValue) Type(_ context.Context) attr.Type {
	return StringWithBackingFieldType{}
}

// Equal returns true if the given value is equivalent.
// TODO [mux-PRs]: adjust
func (v StringWithBackingFieldValue) Equal(o attr.Value) bool {
	other, ok := o.(StringWithBackingFieldValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// TODO [mux-PRs]: adjust
func (v StringWithBackingFieldValue) Unmarshal(target any) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("String with backing field Unmarshal Error", "string value is null"))
		return diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("String with backing field Unmarshal Error", "string value is unknown"))
		return diags
	}

	err := json.Unmarshal([]byte(v.ValueString()), target)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("String with backing field Unmarshal Error", err.Error()))
	}

	return diags
}

// NewStringWithBackingFieldValueNull creates a StringWithBackingFieldValue with a null value. Determine whether the value is null via IsNull method.
func NewStringWithBackingFieldValueNull() StringWithBackingFieldValue {
	return StringWithBackingFieldValue{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewStringWithBackingFieldValueUnknown creates a StringWithBackingFieldValue with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewStringWithBackingFieldValueUnknown() StringWithBackingFieldValue {
	return StringWithBackingFieldValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewStringWithBackingFieldValueValue creates a StringWithBackingFieldValue with a known value. Access the value via ValueString method.
func NewStringWithBackingFieldValueValue(value string) StringWithBackingFieldValue {
	return StringWithBackingFieldValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewStringWithBackingFieldValuePointerValue creates a StringWithBackingFieldValue with a null value if nil or a known value. Access the value via ValueStringPointer method.
func NewStringWithBackingFieldValuePointerValue(value *string) StringWithBackingFieldValue {
	return StringWithBackingFieldValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
