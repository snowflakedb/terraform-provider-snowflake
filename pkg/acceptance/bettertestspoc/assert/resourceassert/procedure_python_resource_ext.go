package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (f *ProcedurePythonResourceAssert) HasImportsLength(len int) *ProcedurePythonResourceAssert {
	f.AddAssertion(assert.ValueSet("imports.#", strconv.FormatInt(int64(len), 10)))
	return f
}
