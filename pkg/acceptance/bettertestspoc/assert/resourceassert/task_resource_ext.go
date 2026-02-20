package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (t *TaskResourceAssert) HasAfter(ids ...sdk.SchemaObjectIdentifier) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("after.#", strconv.FormatInt(int64(len(ids)), 10)))
	for _, id := range ids {
		t.AddAssertion(assert.SetElem("after", id.FullyQualifiedName()))
	}
	return t
}

func (t *TaskResourceAssert) HasScheduleSeconds(seconds int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schedule.#", "1"))
	t.AddAssertion(assert.ValueSet("schedule.0.seconds", strconv.Itoa(seconds)))
	t.AddAssertion(assert.ValueSet("schedule.0.minutes", "0"))
	t.AddAssertion(assert.ValueSet("schedule.0.hours", "0"))
	t.AddAssertion(assert.ValueSet("schedule.0.using_cron", ""))
	return t
}

func (t *TaskResourceAssert) HasScheduleMinutes(minutes int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schedule.#", "1"))
	t.AddAssertion(assert.ValueSet("schedule.0.minutes", strconv.Itoa(minutes)))
	t.AddAssertion(assert.ValueSet("schedule.0.seconds", "0"))
	t.AddAssertion(assert.ValueSet("schedule.0.hours", "0"))
	t.AddAssertion(assert.ValueSet("schedule.0.using_cron", ""))
	return t
}

func (t *TaskResourceAssert) HasScheduleHours(hours int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schedule.#", "1"))
	t.AddAssertion(assert.ValueSet("schedule.0.hours", strconv.Itoa(hours)))
	t.AddAssertion(assert.ValueSet("schedule.0.seconds", "0"))
	t.AddAssertion(assert.ValueSet("schedule.0.minutes", "0"))
	t.AddAssertion(assert.ValueSet("schedule.0.using_cron", ""))
	return t
}

func (t *TaskResourceAssert) HasScheduleCron(cron string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schedule.#", "1"))
	t.AddAssertion(assert.ValueSet("schedule.0.using_cron", cron))
	t.AddAssertion(assert.ValueSet("schedule.0.seconds", "0"))
	t.AddAssertion(assert.ValueSet("schedule.0.minutes", "0"))
	t.AddAssertion(assert.ValueSet("schedule.0.hours", "0"))
	return t
}

func (t *TaskResourceAssert) HasNoScheduleSet() *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schedule.#", "0"))
	return t
}

func (t *TaskResourceAssert) HasUserTaskManagedInitialWarehouseSizeEnum(size sdk.WarehouseSize) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("user_task_managed_initial_warehouse_size", string(size)))
	return t
}

func (t *TaskResourceAssert) HasTargetCompletionIntervalSeconds(seconds int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("target_completion_interval.#", "1"))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.seconds", strconv.Itoa(seconds)))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.minutes", "0"))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.hours", "0"))
	return t
}

func (t *TaskResourceAssert) HasTargetCompletionIntervalMinutes(minutes int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("target_completion_interval.#", "1"))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.minutes", strconv.Itoa(minutes)))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.seconds", "0"))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.hours", "0"))
	return t
}

func (t *TaskResourceAssert) HasTargetCompletionIntervalHours(hours int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("target_completion_interval.#", "1"))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.hours", strconv.Itoa(hours)))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.seconds", "0"))
	t.AddAssertion(assert.ValueSet("target_completion_interval.0.minutes", "0"))
	return t
}

func (t *TaskResourceAssert) HasNoTargetCompletionInterval() *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("target_completion_interval.#", "0"))
	return t
}

func (t *TaskResourceAssert) HasServerlessTaskMinStatementSizeEnum(size sdk.WarehouseSize) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("serverless_task_min_statement_size", string(size)))
	return t
}

func (t *TaskResourceAssert) HasServerlessTaskMaxStatementSizeEnum(size sdk.WarehouseSize) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("serverless_task_max_statement_size", string(size)))
	return t
}

func (t *TaskResourceAssert) HasDefaultServerlessTaskMinStatementSize() *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("serverless_task_min_statement_size", string(sdk.WarehouseSizeXSmall)))
	return t
}

func (t *TaskResourceAssert) HasDefaultServerlessTaskMaxStatementSize() *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("serverless_task_max_statement_size", "X2Large"))
	return t
}

func (t *TaskResourceAssert) HasDefaultUserTaskManagedInitialWarehouseSize() *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("user_task_managed_initial_warehouse_size", "Medium"))
	return t
}
