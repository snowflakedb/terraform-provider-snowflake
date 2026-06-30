package resourceassert

import (
	"strconv"
)

func (f *FunctionJavaResourceAssert) HasImportsLength(len int) *FunctionJavaResourceAssert {
	f.ValueSet("imports.#", strconv.FormatInt(int64(len), 10))
	return f
}
