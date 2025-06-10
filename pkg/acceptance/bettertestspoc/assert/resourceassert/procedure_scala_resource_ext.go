package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (f *ProcedureScalaResourceAssert) HasImportsLength(len int) *ProcedureScalaResourceAssert {
	f.AddAssertion(assert.ValueSet("imports.#", strconv.FormatInt(int64(len), 10)))
	return f
}
