package resources

import "context"

func v2_15_0_StreamStateUpgrader(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	if v, ok := rawState["show_initial_rows"]; !ok || v == nil || v == "" {
		rawState["show_initial_rows"] = BooleanDefault
	}
	return rawState, nil
}
