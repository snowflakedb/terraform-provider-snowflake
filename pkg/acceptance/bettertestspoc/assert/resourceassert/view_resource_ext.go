package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (v *ViewResourceAssert) HasColumnLength(len int) *ViewResourceAssert {
	v.AddAssertion(assert.ValueSet("column.#", strconv.FormatInt(int64(len), 10)))
	return v
}

func (v *ViewResourceAssert) HasAggregationPolicyLength(len int) *ViewResourceAssert {
	v.AddAssertion(assert.ValueSet("aggregation_policy.#", strconv.FormatInt(int64(len), 10)))
	return v
}

func (v *ViewResourceAssert) HasRowAccessPolicyLength(len int) *ViewResourceAssert {
	v.AddAssertion(assert.ValueSet("row_access_policy.#", strconv.FormatInt(int64(len), 10)))
	return v
}

func (v *ViewResourceAssert) HasDataMetricScheduleLength(len int) *ViewResourceAssert {
	v.AddAssertion(assert.ValueSet("data_metric_schedule.#", strconv.FormatInt(int64(len), 10)))
	return v
}

func (v *ViewResourceAssert) HasDataMetricFunctionLength(len int) *ViewResourceAssert {
	v.AddAssertion(assert.ValueSet("data_metric_function.#", strconv.FormatInt(int64(len), 10)))
	return v
}

func (v *ViewResourceAssert) HasNoAggregationPolicyByLength() *ViewResourceAssert {
	v.AddAssertion(assert.ValueNotSet("aggregation_policy.#"))
	return v
}

func (v *ViewResourceAssert) HasNoRowAccessPolicyByLength() *ViewResourceAssert {
	v.AddAssertion(assert.ValueNotSet("row_access_policy.#"))
	return v
}

func (v *ViewResourceAssert) HasNoDataMetricScheduleByLength() *ViewResourceAssert {
	v.AddAssertion(assert.ValueNotSet("data_metric_schedule.#"))
	return v
}

func (v *ViewResourceAssert) HasNoDataMetricFunctionByLength() *ViewResourceAssert {
	v.AddAssertion(assert.ValueNotSet("data_metric_function.#"))
	return v
}
