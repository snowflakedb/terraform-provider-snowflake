package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (t *TaskResourceAssert) HasAfter(ids ...sdk.SchemaObjectIdentifier) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("after.#", strconv.FormatInt(int64(len(ids)), 10)))
	for _, id := range ids {
		t.AddAssertion(assert.SetElem("after.*", id.FullyQualifiedName()))
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
