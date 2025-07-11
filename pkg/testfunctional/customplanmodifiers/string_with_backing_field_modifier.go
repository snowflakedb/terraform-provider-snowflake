package customplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func NewStringWithBackingFieldModifier() planmodifier.String {
	return stringWithBackingFieldModifier{}
}

type stringWithBackingFieldModifier struct{}

func (m stringWithBackingFieldModifier) Description(_ context.Context) string {
	return "TODO"
}

func (m stringWithBackingFieldModifier) MarkdownDescription(_ context.Context) string {
	return "TODO"
}

func (m stringWithBackingFieldModifier) PlanModifyString(_ context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	// TODO
	_ = request.StateValue
}
