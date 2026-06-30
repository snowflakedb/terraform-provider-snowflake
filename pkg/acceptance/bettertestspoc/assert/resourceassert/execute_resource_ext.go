package resourceassert

import (
	"fmt"
	"strconv"
)

func (e *ExecuteResourceAssert) HasQueryResultsLength(len int) *ExecuteResourceAssert {
	e.ValueSet("query_results.#", strconv.FormatInt(int64(len), 10))
	return e
}

func (e *ExecuteResourceAssert) HasKeyValueOnIdx(key string, value string, idx int) *ExecuteResourceAssert {
	e.ValueSet(fmt.Sprintf("query_results.%d.%s", idx, key), value)
	return e
}
