package resourceshowoutputassert

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TaskDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func TaskDatasourceShowOutput(t *testing.T, name string) *TaskShowOutputAssert {
	t.Helper()

	taskAssert := TaskShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "tasks.0."),
	}
	taskAssert.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &taskAssert
}

func (t *TaskShowOutputAssert) HasErrorIntegrationEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("error_integration", ""))
	return t
}

func (t *TaskShowOutputAssert) HasCreatedOnNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return t
}

func (t *TaskShowOutputAssert) HasIdNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("id"))
	return t
}

func (t *TaskShowOutputAssert) HasOwnerNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return t
}

func (t *TaskShowOutputAssert) HasLastCommittedOnNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("last_committed_on"))
	return t
}

func (t *TaskShowOutputAssert) HasLastSuspendedOnNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("last_suspended_on"))
	return t
}

func (t *TaskShowOutputAssert) HasPredecessors(predecessors ...sdk.SchemaObjectIdentifier) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("predecessors.#", strconv.Itoa(len(predecessors))))
	for i, predecessor := range predecessors {
		t.AddAssertion(assert.ResourceShowOutputValueSet(fmt.Sprintf("predecessors.%d", i), predecessor.FullyQualifiedName()))
	}
	return t
}

func (t *TaskShowOutputAssert) HasTaskRelations(expected sdk.TaskRelations) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("task_relations.#", "1"))
	t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("task_relations.0.predecessors.#", strconv.Itoa(len(expected.Predecessors))))
	for i, predecessor := range expected.Predecessors {
		t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet(fmt.Sprintf("task_relations.0.predecessors.%d", i), predecessor.FullyQualifiedName()))
	}
	if expected.FinalizerTask != nil && len(expected.FinalizerTask.Name()) > 0 {
		t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("task_relations.0.finalizer", expected.FinalizerTask.FullyQualifiedName()))
	}
	if expected.FinalizedRootTask != nil && len(expected.FinalizedRootTask.Name()) > 0 {
		t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("task_relations.0.finalized_root_task", expected.FinalizedRootTask.FullyQualifiedName()))
	}
	return t
}

func (t *TaskShowOutputAssert) HasScheduleMinutes(minutes int) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("schedule", fmt.Sprintf("%d MINUTE", minutes)))
	return t
}

func (t *TaskShowOutputAssert) HasScheduleCron(cron string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("schedule", fmt.Sprintf("USING CRON %s", cron)))
	return t
}
