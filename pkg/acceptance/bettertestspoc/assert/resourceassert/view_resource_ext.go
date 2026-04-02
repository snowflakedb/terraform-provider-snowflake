package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func (v *ViewResourceAssert) HasColumns(columns []sdk.ViewColumn) *ViewResourceAssert {
	v.AddAssertion(assert.ValueSet("column.#", strconv.Itoa(len(columns))))
	for i, col := range columns {
		v.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.column_name", i), col.Name))
		if col.MaskingPolicy != nil {
			v.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.masking_policy.#", i), "1"))
			v.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.masking_policy.0.policy_name", i), col.MaskingPolicy.MaskingPolicy.FullyQualifiedName()))
		}
	}
	return v
}
