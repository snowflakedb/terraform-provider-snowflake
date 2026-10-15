package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// CreateConstructorArg is one argument of NewCreate<Object>Request after the id.
type CreateConstructorArg struct {
	Get string
}

// CreateOptionalCall is one helper call inside ApplyCreateOptionals.
// Target is the request expression passed after the attribute key.
type CreateOptionalCall struct {
	Helper string
	Target string
}

func attachCreateMappings(objectName string, src []ssot.Attribute, dst []ResourceAttribute) {
	if len(src) != len(dst) {
		panic(fmt.Sprintf("object %q: %d source attributes and %d resource attributes", objectName, len(src), len(dst)))
	}
	for i := range src {
		if get, ok := src[i].ConstructorGet(); ok {
			dst[i].CreateConstructorArg = &CreateConstructorArg{Get: get}
			continue
		}
		if src[i].CreateMapping == nil {
			continue
		}
		call := src[i].CreateMapping.CreateCall(objectName, src[i])
		dst[i].CreateOptionalCall = &CreateOptionalCall{Helper: call.Helper, Target: call.Target}
	}
}
