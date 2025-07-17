package customtypes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type Metadata struct {
	FieldA string `json:"field_a"`
	FieldB string `json:"field_b,omitempty"`
}

var _ basetypes.StringTypable = StringWithMetadataType{}

type StringWithMetadataType struct {
	basetypes.StringType
}

// String returns a human-readable string of the type name.
func (t StringWithMetadataType) String() string {
	return "customtypes.StringWithMetadataType"
}

func (t StringWithMetadataType) ValueType(_ context.Context) attr.Value {
	return StringWithMetadataValue{}
}

// Equal returns true if the given type is equivalent.
func (t StringWithMetadataType) Equal(o attr.Type) bool {
	_, ok := o.(StringWithMetadataType)
	return ok
}

// ValueFromTerraform returns a Value given a tftypes.Value. This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t StringWithMetadataType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	baseValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := baseValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", baseValue)
	}

	parts := strings.SplitN(stringValue.ValueString(), "|", 2)
	text := parts[0]
	var meta Metadata
	if len(parts) > 1 {
		if err := json.Unmarshal([]byte(parts[1]), &meta); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}
	return StringWithMetadataValue{
		StringValue: types.StringValue(text),
		Metadata:    meta,
	}, nil
}
