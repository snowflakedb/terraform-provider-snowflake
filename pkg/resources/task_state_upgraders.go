package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func v098TaskStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["condition"] = rawState["when"]
	rawState["started"] = rawState["enabled"].(bool)
	rawState["allow_overlapping_execution"] = booleanStringFromBool(rawState["allow_overlapping_execution"].(bool))
	if rawState["after"] != nil {
		if afterSlice, okType := rawState["after"].([]any); okType {
			newAfter := make([]string, len(afterSlice))
			for i, name := range afterSlice {
				newAfter[i] = sdk.NewSchemaObjectIdentifier(rawState["database"].(string), rawState["schema"].(string), name.(string)).FullyQualifiedName()
			}
			rawState["after"] = newAfter
		}
	}
	if rawState["session_parameters"] != nil {
		if sessionParamsMap, okType := rawState["session_parameters"].(map[string]any); okType {
			for k, v := range sessionParamsMap {
				rawState[k] = v
			}
		}
	}
	delete(rawState, "session_parameters")

	if rawState["schedule"] != nil && len(rawState["schedule"].(string)) > 0 {
		taskSchedule, err := sdk.ParseTaskSchedule(rawState["schedule"].(string))
		scheduleMap := make(map[string]any)
		if err != nil {
			return nil, err
		}
		switch {
		case len(taskSchedule.Cron) > 0:
			scheduleMap["using_cron"] = taskSchedule.Cron
		case taskSchedule.Minutes > 0:
			scheduleMap["minutes"] = taskSchedule.Minutes
		}
		rawState["schedule"] = []any{scheduleMap}
	} else {
		delete(rawState, "schedule")
	}

	return migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx, rawState, meta)
}
