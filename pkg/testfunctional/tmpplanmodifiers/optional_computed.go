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

	// When the attribute is removed from config and the read is run, then the plan would be empty. We need to react to such a change because we want to unset the value when it is removed from config.
	// However, there is no way to distinguish between the first situation like this one and the subsequent ones, therefore resulting in permadiff.
	// Additional field with previous value is required.
	if request.ConfigValue.IsNull() && !request.StateValue.IsNull() {
		response.PlanValue = types.StringUnknown()
		return
	}
}
