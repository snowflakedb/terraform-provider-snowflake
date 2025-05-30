// Code generated by assertions generator; DO NOT EDIT.

package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// to ensure sdk package is used
var _ = sdk.Object{}

type TaskShowOutputAssert struct {
	*assert.ResourceAssert
}

func TaskShowOutput(t *testing.T, name string) *TaskShowOutputAssert {
	t.Helper()

	task := TaskShowOutputAssert{
		ResourceAssert: assert.NewResourceAssert(name, "show_output"),
	}
	task.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &task
}

func ImportedTaskShowOutput(t *testing.T, id string) *TaskShowOutputAssert {
	t.Helper()

	task := TaskShowOutputAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "show_output"),
	}
	task.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &task
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (t *TaskShowOutputAssert) HasCreatedOn(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("created_on", expected))
	return t
}

func (t *TaskShowOutputAssert) HasName(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("name", expected))
	return t
}

func (t *TaskShowOutputAssert) HasId(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("id", expected))
	return t
}

func (t *TaskShowOutputAssert) HasDatabaseName(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("database_name", expected))
	return t
}

func (t *TaskShowOutputAssert) HasSchemaName(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("schema_name", expected))
	return t
}

func (t *TaskShowOutputAssert) HasOwner(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("owner", expected))
	return t
}

func (t *TaskShowOutputAssert) HasComment(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("comment", expected))
	return t
}

func (t *TaskShowOutputAssert) HasWarehouse(expected sdk.AccountObjectIdentifier) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("warehouse", expected.Name()))
	return t
}

func (t *TaskShowOutputAssert) HasSchedule(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("schedule", expected))
	return t
}

func (t *TaskShowOutputAssert) HasState(expected sdk.TaskState) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("state", expected))
	return t
}

func (t *TaskShowOutputAssert) HasDefinition(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("definition", expected))
	return t
}

func (t *TaskShowOutputAssert) HasCondition(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("condition", expected))
	return t
}

func (t *TaskShowOutputAssert) HasAllowOverlappingExecution(expected bool) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputBoolValueSet("allow_overlapping_execution", expected))
	return t
}

func (t *TaskShowOutputAssert) HasErrorIntegration(expected sdk.AccountObjectIdentifier) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("error_integration", expected.Name()))
	return t
}

func (t *TaskShowOutputAssert) HasLastCommittedOn(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("last_committed_on", expected))
	return t
}

func (t *TaskShowOutputAssert) HasLastSuspendedOn(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("last_suspended_on", expected))
	return t
}

func (t *TaskShowOutputAssert) HasOwnerRoleType(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("owner_role_type", expected))
	return t
}

func (t *TaskShowOutputAssert) HasConfig(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("config", expected))
	return t
}

func (t *TaskShowOutputAssert) HasBudget(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("budget", expected))
	return t
}

func (t *TaskShowOutputAssert) HasLastSuspendedReason(expected string) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("last_suspended_reason", expected))
	return t
}
