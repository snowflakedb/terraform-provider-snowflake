package resourceassert

import (
	"strconv"
)

func (f *ProcedurePythonResourceAssert) HasImportsLength(len int) *ProcedurePythonResourceAssert {
	f.ValueSet("imports.#", strconv.FormatInt(int64(len), 10))
	return f
}
