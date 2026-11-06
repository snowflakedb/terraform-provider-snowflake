package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// SelectObjects returns, in All order, the objects a generator takes as input. An object is selected when the generator's
// flag is on (enabled) and the generator's own filter accepts it (selects). The filter is one rule per generator, so an
// object without the generator's input is skipped instead of failing, and the flag can be on for every object.
func SelectObjects(enabled func(ssot.Generators) bool, selects func(ssot.Object) bool) []ssot.Object {
	var out []ssot.Object
	for _, obj := range All() {
		if enabled(ssot.ObjectGenerators(obj.Generators...)) && selects(obj) {
			out = append(out, obj)
		}
	}
	return out
}
