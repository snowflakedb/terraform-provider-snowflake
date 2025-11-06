package resources

import (
	"context"
)

func v2_10_0_CurrentAccountStateUpgrader(_ context.Context, state map[string]any, _ any) (map[string]any, error) {
	if state == nil {
		return state, nil
	}

	if v, ok := state["saml_identity_provider"]; ok && v != nil && v.(string) == "" {
		delete(state, "saml_identity_provider")
	}

	return state, nil
}
