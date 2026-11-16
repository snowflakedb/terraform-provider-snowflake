package resources

import "context"

// v2_23_0_IsAgenticSecurityIntegrationStateUpgrader fills is_agentic on state written before the attribute existed.
// It is shared by snowflake_oauth_integration_for_custom_clients and snowflake_external_oauth_integration.
// The schema default is "default", so an empty prior value would plan an update whose apply runs UNSET IS_AGENTIC,
// resetting a value that may have been set outside Terraform. The upgrader doesn't call Snowflake.
func v2_23_0_IsAgenticSecurityIntegrationStateUpgrader(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	if v, ok := rawState["is_agentic"].(string); !ok || v == "" {
		rawState["is_agentic"] = BooleanDefault
	}
	return rawState, nil
}
