package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// UpdateCall is one helper call inside ApplySetUnset.
// Target is the request expression passed after the attribute key.
type UpdateCall struct {
	Helper string
	Target string
}

func attachUpdateMappings(objectName string, src []ssot.Attribute, dst []ResourceAttribute) {
	if len(src) != len(dst) {
		panic(fmt.Sprintf("object %q: %d source attributes and %d resource attributes", objectName, len(src), len(dst)))
	}
	for i := range src {
		if src[i].UpdateMapping == nil {
			continue
		}
		call := src[i].UpdateMapping.UpdateCall(objectName, src[i])
		dst[i].UpdateCall = &UpdateCall{Helper: call.Helper, Target: call.Target}
	}
}
