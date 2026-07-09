package resourceassert

import (
	"strconv"
)

func (f *ProcedureJavaResourceAssert) HasImportsLength(len int) *ProcedureJavaResourceAssert {
	f.ValueSet("imports.#", strconv.FormatInt(int64(len), 10))
	return f
}
