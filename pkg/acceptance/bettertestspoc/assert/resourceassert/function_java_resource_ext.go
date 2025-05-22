package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (f *FunctionJavaResourceAssert) HasImportsLength(len int) *FunctionJavaResourceAssert {
	f.AddAssertion(assert.ValueSet("imports.#", strconv.FormatInt(int64(len), 10)))
	return f
}

func (f *FunctionJavaResourceAssert) HasTargetPathEmpty() *FunctionJavaResourceAssert {
	f.AddAssertion(assert.ValueSet("target_path.#", "0"))
	return f
}
