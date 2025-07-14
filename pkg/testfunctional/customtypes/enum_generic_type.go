package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*EnumType)(nil)
)

type EnumType struct {
	basetypes.StringType
}

func (t EnumType) String() string {
	return "EnumType"
}

func (t EnumType) ValueType(_ context.Context) attr.Value {
	return EnumValue{}
}

func (t EnumType) Equal(o attr.Type) bool {
	other, ok := o.(EnumType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t EnumType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return EnumValue{
		StringValue: in,
	}, nil
}

func (t EnumType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
