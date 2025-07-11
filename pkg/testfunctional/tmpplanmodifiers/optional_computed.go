package tmpplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func OptionalComputedString() planmodifier.String {
	return optionalComputedStringModifier{}
}

type optionalComputedStringModifier struct{}

func (m optionalComputedStringModifier) Description(_ context.Context) string {
	return "TODO"
}

func (m optionalComputedStringModifier) MarkdownDescription(_ context.Context) string {
	return "TODO"
}

func (m optionalComputedStringModifier) PlanModifyString(_ context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	// Do nothing if there is no state (resource is being created).
	if request.State.Raw.IsNull() {
		return
	}

	// Do nothing if set in config.
	if !request.ConfigValue.IsNull() {
		return
	}

	if request.ConfigValue.IsNull() && request.StateValue.IsNull() {
		response.PlanValue = types.StringUnknown()
		return
	}
	//
	//// Do nothing if there is a known planned value.
	//if !request.PlanValue.IsUnknown() {
	//	return
	//}
	//
	//// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	//if request.ConfigValue.IsUnknown() {
	//	return
	//}
	//
	//response.PlanValue = request.StateValue
}
