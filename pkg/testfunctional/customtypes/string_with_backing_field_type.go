package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ attr.Type = StringWithBackingFieldType{}

type StringWithBackingFieldType struct {
	basetypes.StringType
}

// String returns a human-readable string of the type name.
func (t StringWithBackingFieldType) String() string {
	return "customtypes.StringWithBackingFieldType"
}

func (t StringWithBackingFieldType) ValueType(_ context.Context) attr.Value {
	return StringWithBackingFieldValue{}
}

// Equal returns true if the given type is equivalent.
// TODO [mux-PRs]: adjust
func (t StringWithBackingFieldType) Equal(o attr.Type) bool {
	other, ok := o.(StringWithBackingFieldType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// ValueFromString returns a StringValuable type given a StringValue.
// TODO [mux-PRs]: adjust
func (t StringWithBackingFieldType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return StringWithBackingFieldValue{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value. This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
// TODO [mux-PRs]: adjust
func (t StringWithBackingFieldType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
