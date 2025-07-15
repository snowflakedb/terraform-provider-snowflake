package customplanmodifiers

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customtypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func EnumSuppressor[T customtypes.EnumCreator[T]]() planmodifier.String {
	return enumSuppressorModifier[T]{}
}

type enumSuppressorModifier[T customtypes.EnumCreator[T]] struct {
	et T
}

func (m enumSuppressorModifier[T]) Description(_ context.Context) string {
	return "enum suppressor"
}

func (m enumSuppressorModifier[T]) MarkdownDescription(_ context.Context) string {
	return "enum suppressor"
}

// PlanModifyString implements the plan modification logic.
func (m enumSuppressorModifier[T]) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	// Do nothing if there is unknown planned value.
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}

	plan := customtypes.NewEnumValueFromStringValue[T](req.PlanValue)
	state := customtypes.NewEnumValueFromStringValue[T](req.StateValue)
	result, d := state.StringSemanticEquals(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	if result {
		resp.PlanValue = req.StateValue
	}
}
