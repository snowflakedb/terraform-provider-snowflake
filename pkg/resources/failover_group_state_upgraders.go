package resources

import (
	"context"
	"errors"
	"fmt"
)

// v2_24_0_FailoverGroupStateUpgrader converts the state written by the hand-written resource:
//   - replication_schedule was a nested block (cron or interval); it is a string now;
//   - from_replica moved to snowflake_failover_group_secondary, so a group that used it cannot be upgraded in place.
func v2_24_0_FailoverGroupStateUpgrader(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	if fromReplica, ok := rawState["from_replica"].([]any); ok && len(fromReplica) > 0 {
		return nil, errors.New("the failover group was created with the from_replica block, which is not a part of snowflake_failover_group anymore; " +
			"remove it from the state (terraform state rm) and import it as snowflake_failover_group_secondary (see the migration guide)")
	}
	delete(rawState, "from_replica")

	schedule, err := failoverGroupScheduleFromState(rawState["replication_schedule"])
	if err != nil {
		return nil, err
	}
	if schedule == "" {
		delete(rawState, "replication_schedule")
	} else {
		rawState["replication_schedule"] = schedule
	}

	return rawState, nil
}

func failoverGroupScheduleFromState(raw any) (string, error) {
	switch value := raw.(type) {
	case nil:
		return "", nil
	case string:
		return value, nil
	case []any:
		if len(value) == 0 {
			return "", nil
		}
		block, ok := value[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("cannot read replication_schedule from the state: %v", value[0])
		}
		if crons, ok := block["cron"].([]any); ok && len(crons) > 0 {
			cron, ok := crons[0].(map[string]any)
			if !ok {
				return "", fmt.Errorf("cannot read replication_schedule.cron from the state: %v", crons[0])
			}
			schedule := "USING CRON " + fmt.Sprint(cron["expression"])
			if timeZone, ok := cron["time_zone"].(string); ok && timeZone != "" {
				schedule += " " + timeZone
			}
			return schedule, nil
		}
		switch interval := block["interval"].(type) {
		case float64:
			if interval > 0 {
				return fmt.Sprintf("%d MINUTE", int(interval)), nil
			}
		case int:
			if interval > 0 {
				return fmt.Sprintf("%d MINUTE", interval), nil
			}
		}
		return "", nil
	default:
		return "", fmt.Errorf("cannot read replication_schedule from the state: %v", raw)
	}
}
