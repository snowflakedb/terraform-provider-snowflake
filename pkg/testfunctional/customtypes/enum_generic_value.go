package customtypes

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable                   = (*EnumValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*EnumValue)(nil)
	_ xattr.ValidateableAttribute                = (*EnumValue)(nil)
)

type EnumValue struct {
	basetypes.StringValue
}

func (v EnumValue) Type(_ context.Context) attr.Type {
	return EnumType{}
}

func (v EnumValue) Equal(o attr.Value) bool {
	other, ok := o.(EnumValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v EnumValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	newValue, ok := newValuable.(EnumValue)
	if !ok {
		diags.AddError("TODO", "TODO")
		return false, diags
	}

	// TODO: parameterize func
	result, err := sameAfterNormalization(newValue.ValueString(), v.ValueString(), sdk.ToWarehouseType)
	if err != nil {
		diags.AddError("TODO", "TODO")
		return false, diags
	}

	return result, diags
}

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

func (v EnumValue) ValidateAttribute(_ context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	// TODO: parameterize func
	_, err := sdk.ToWarehouseType(v.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "TODO", "TODO")
		return
	}
}
