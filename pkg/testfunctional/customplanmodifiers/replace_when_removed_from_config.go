package customplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// TODO [mux-PR]: add functional test
func RequiresReplaceIfRemovedFromConfig() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			if !req.StateValue.IsNull() && req.ConfigValue.IsNull() {
				resp.RequiresReplace = true
			}
		},
		"If the value of this attribute is configured and then removed from config, Terraform will destroy and recreate the resource.",
		"If the value of this attribute is configured then removed from config, Terraform will destroy and recreate the resource.",
	)
}
