package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (e *ExecuteResourceAssert) HasQueryResultsLength(len int) *ExecuteResourceAssert {
	e.AddAssertion(assert.ValueSet("query_results.#", strconv.FormatInt(int64(len), 10)))
	return e
}

func (e *ExecuteResourceAssert) HasKeyValueOnIdx(key string, value string, idx int) *ExecuteResourceAssert {
	e.AddAssertion(assert.ValueSet(fmt.Sprintf("query_results.%d.%s", idx, key), value))
	return e
}
