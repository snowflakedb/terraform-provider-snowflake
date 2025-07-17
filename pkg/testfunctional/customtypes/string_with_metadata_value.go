package customtypes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringValuableWithSemanticEquals = StringWithMetadataValue{}

type StringWithMetadataValue struct {
	basetypes.StringValue

	Metadata Metadata // we do not expose it to practitioner in schema
}

func (v StringWithMetadataValue) Type(_ context.Context) attr.Type {
	return StringWithMetadataType{}
}

// Equal returns true if the given value is equivalent.
func (v StringWithMetadataValue) Equal(o attr.Value) bool {
	other, ok := o.(StringWithMetadataValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// ToTerraformValue returns the data contained in the *String as a tftypes.Value.
func (v StringWithMetadataValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	raw := v.ValueString()
	metaBytes, err := json.Marshal(v.Metadata)
	if err != nil {
		return types.StringValue(raw).ToTerraformValue(ctx)
	}
	combined := fmt.Sprintf(`%s|%s`, raw, metaBytes)
	return types.StringValue(combined).ToTerraformValue(ctx)
}

func (v StringWithMetadataValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// The framework should always pass the correct value type, but always check
	newValue, ok := newValuable.(StringWithMetadataValue)

	if !ok {
		diags.AddError("Unexpected type", "Unexpected type was passed to StringWithMetadataValue")
		return false, diags
	}

	return v.ValueString() == newValue.ValueString(), diags
}
