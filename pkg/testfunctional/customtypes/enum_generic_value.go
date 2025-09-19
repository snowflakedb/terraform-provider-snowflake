package customtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable                   = (*EnumValue[dummyEnumType])(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*EnumValue[dummyEnumType])(nil)
	_ xattr.ValidateableAttribute                = (*EnumValue[dummyEnumType])(nil)
)

type EnumValue[T EnumCreator[T]] struct {
	basetypes.StringValue
	et T
}

func (v EnumValue[T]) Type(_ context.Context) attr.Type {
	return EnumType[T]{}
}

func (v EnumValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(EnumValue[T])

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v EnumValue[T]) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	newValue, ok := newValuable.(EnumValue[T])
	if !ok {
		// TODO [SNOW-2298131]: better diags later
		diags.AddError("Incomparable types", newValuable.String())
		return false, diags
	}

	result, err := v.sameAfterNormalization(newValue.ValueString(), v.ValueString())
	if err != nil {
		// TODO [SNOW-2298131]: better diags later
		diags.AddError("Normalization failed", err.Error())
		return false, diags
	}

	return result, diags
}

func (v EnumValue[T]) sameAfterNormalization(oldValue string, newValue string) (bool, error) {
	oldNormalized, err := v.et.FromString(oldValue)
	if err != nil {
		return false, err
	}
	newNormalized, err := v.et.FromString(newValue)
	if err != nil {
		return false, err
	}

	return oldNormalized == newNormalized, nil
}

func (v EnumValue[T]) ValidateAttribute(_ context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	_, err := v.et.FromString(v.ValueString())
	if err != nil {
		// TODO [SNOW-2298131]: better diags later
		resp.Diagnostics.AddAttributeError(req.Path, "Incorrect value for attribute", err.Error())
		return
	}
}

func NewEnumValue[T EnumCreator[T]](value T) EnumValue[T] {
	return EnumValue[T]{
		StringValue: types.StringValue(string(value)),
	}
}

func NewEnumValueFromStringValue[T EnumCreator[T]](value types.String) EnumValue[T] {
	return EnumValue[T]{
		StringValue: value,
	}
}

func (v EnumValue[T]) Normalize() (T, error) {
	return v.et.FromString(v.ValueString())
}
