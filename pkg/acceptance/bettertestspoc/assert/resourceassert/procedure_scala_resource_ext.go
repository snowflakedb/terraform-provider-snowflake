package resourceassert

import (
	"strconv"
)

func (f *ProcedureScalaResourceAssert) HasImportsLength(len int) *ProcedureScalaResourceAssert {
	f.ValueSet("imports.#", strconv.FormatInt(int64(len), 10))
	return f
}
