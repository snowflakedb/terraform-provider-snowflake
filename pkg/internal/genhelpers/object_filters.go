package genhelpers

import (
	"slices"
)

func filterObjectByNameProvider[T ObjectNameProvider](allowedObjectNames []string) func(object T) bool {
	return func(object T) bool {
		if len(allowedObjectNames) == 0 {
			return true
		}
		return slices.Contains(allowedObjectNames, object.ObjectName())
	}
}
