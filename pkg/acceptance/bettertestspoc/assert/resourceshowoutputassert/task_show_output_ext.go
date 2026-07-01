package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *TaskShowOutputAssert) HasErrorIntegrationEmpty() *TaskShowOutputAssert {
	t.StringValueSet("error_integration", "")
	return t
}

func (t *TaskShowOutputAssert) HasCreatedOnNotEmpty() *TaskShowOutputAssert {
	t.ValuePresent("created_on")
	return t
}

func (t *TaskShowOutputAssert) HasIdNotEmpty() *TaskShowOutputAssert {
	t.ValuePresent("id")
	return t
}

func (t *TaskShowOutputAssert) HasOwnerNotEmpty() *TaskShowOutputAssert {
	t.ValuePresent("owner")
	return t
}

func (t *TaskShowOutputAssert) HasLastCommittedOnNotEmpty() *TaskShowOutputAssert {
	t.ValuePresent("last_committed_on")
	return t
}

func (t *TaskShowOutputAssert) HasLastSuspendedOnNotEmpty() *TaskShowOutputAssert {
	t.ValuePresent("last_suspended_on")
	return t
}

func (t *TaskShowOutputAssert) HasPredecessors(predecessors ...sdk.SchemaObjectIdentifier) *TaskShowOutputAssert {
	t.StringValueSet("predecessors.#", strconv.Itoa(len(predecessors)))
	for _, predecessor := range predecessors {
		t.SetContainsElem("predecessors", predecessor.FullyQualifiedName())
	}
	return t
}

func (t *TaskShowOutputAssert) HasTaskRelations(expected sdk.TaskRelations) *TaskShowOutputAssert {
	t.StringValueSet("task_relations.#", "1")
	t.StringValueSet("task_relations.0.predecessors.#", strconv.Itoa(len(expected.Predecessors)))
	for i, predecessor := range expected.Predecessors {
		t.StringValueSet(fmt.Sprintf("task_relations.0.predecessors.%d", i), predecessor.FullyQualifiedName())
	}
	if expected.FinalizerTask != nil && len(expected.FinalizerTask.Name()) > 0 {
		t.StringValueSet("task_relations.0.finalizer", expected.FinalizerTask.FullyQualifiedName())
	}
	if expected.FinalizedRootTask != nil && len(expected.FinalizedRootTask.Name()) > 0 {
		t.StringValueSet("task_relations.0.finalized_root_task", expected.FinalizedRootTask.FullyQualifiedName())
	}
	return t
}

func (t *TaskShowOutputAssert) HasScheduleEmpty() *TaskShowOutputAssert {
	t.StringValueSet("schedule", "")
	return t
}

func (t *TaskShowOutputAssert) HasScheduleSeconds(seconds int) *TaskShowOutputAssert {
	t.StringValueSet("schedule", fmt.Sprintf("%d SECOND", seconds))
	return t
}

func (t *TaskShowOutputAssert) HasScheduleMinutes(minutes int) *TaskShowOutputAssert {
	t.StringValueSet("schedule", fmt.Sprintf("%d MINUTE", minutes))
	return t
}

func (t *TaskShowOutputAssert) HasScheduleHours(hours int) *TaskShowOutputAssert {
	t.StringValueSet("schedule", fmt.Sprintf("%d HOUR", hours))
	return t
}

func (t *TaskShowOutputAssert) HasScheduleCron(cron string) *TaskShowOutputAssert {
	t.StringValueSet("schedule", fmt.Sprintf("USING CRON %s", cron))
	return t
}

func (t *TaskShowOutputAssert) HasTargetCompletionIntervalString(expected string) *TaskShowOutputAssert {
	t.StringValueSet("target_completion_interval", expected)
	return t
}

func (t *TaskShowOutputAssert) HasTargetCompletionIntervalEmpty() *TaskShowOutputAssert {
	t.StringValueSet("target_completion_interval.#", "0")
	return t
}

func (t *TaskShowOutputAssert) HasTargetCompletionIntervalSeconds(seconds int) *TaskShowOutputAssert {
	t.StringValueSet("target_completion_interval.#", "1")
	t.StringValueSet("target_completion_interval.0.seconds", strconv.Itoa(seconds))
	t.StringValueSet("target_completion_interval.0.minutes", "0")
	t.StringValueSet("target_completion_interval.0.hours", "0")
	return t
}

func (t *TaskShowOutputAssert) HasTargetCompletionIntervalMinutes(minutes int) *TaskShowOutputAssert {
	t.StringValueSet("target_completion_interval.#", "1")
	t.StringValueSet("target_completion_interval.0.minutes", strconv.Itoa(minutes))
	t.StringValueSet("target_completion_interval.0.seconds", "0")
	t.StringValueSet("target_completion_interval.0.hours", "0")
	return t
}

func (t *TaskShowOutputAssert) HasTargetCompletionIntervalHours(hours int) *TaskShowOutputAssert {
	t.StringValueSet("target_completion_interval.#", "1")
	t.StringValueSet("target_completion_interval.0.hours", strconv.Itoa(hours))
	t.StringValueSet("target_completion_interval.0.seconds", "0")
	t.StringValueSet("target_completion_interval.0.minutes", "0")
	return t
}

func (t *TaskShowOutputAssert) HasWarehouseEmpty() *TaskShowOutputAssert {
	t.StringValueSet("warehouse", "")
	return t
}
