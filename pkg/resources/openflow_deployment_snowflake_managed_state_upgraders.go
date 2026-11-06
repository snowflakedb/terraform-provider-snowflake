package resources

import "context"

// v2_23_0_OpenflowDeploymentSnowflakeManagedUsePrivateLinkUpgrader fills use_private_link on state
// written before the attribute existed. It is ForceNew, and the schema default is "default", so an
// empty prior value would plan a replacement that only rewrites the same deployment.
func v2_23_0_OpenflowDeploymentSnowflakeManagedUsePrivateLinkUpgrader(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	if v, ok := rawState["use_private_link"].(string); !ok || v == "" {
		rawState["use_private_link"] = BooleanDefault
	}
	return rawState, nil
}
