package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (w *WarehouseShowOutputAssert) HasResourceConstraintEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValueSet("resource_constraint", ""))
	return w
}
