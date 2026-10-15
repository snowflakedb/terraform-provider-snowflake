package ssot

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// CreateMapping names the resource create helper for one attribute.
// A required attribute with a nil CreateMapping is a create constructor argument.
type CreateMapping interface {
	CreateCall(objectName string, attr Attribute) CreateCall
}

// CreateCall is one generated ApplyCreateOptionals call.
// Target is the argument after the attribute key.
type CreateCall struct {
	Helper string
	Target string
}

// HelperMapping is a create helper called with the attribute key and request.With<Field>.
// The value is the helper function name.
type HelperMapping string

func (m HelperMapping) CreateCall(_ string, attr Attribute) CreateCall {
	return CreateCall{
		Helper: string(m),
		Target: "request.With" + genhelpers.SnakeCaseToCamel(attr.Name),
	}
}

const (
	StringAttributeCreateBuilder                  HelperMapping = "stringAttributeCreateBuilder"
	IntAttributeCreateBuilder                     HelperMapping = "intAttributeCreateBuilder"
	IntAttributeWithSpecialDefaultCreateBuilder   HelperMapping = "intAttributeWithSpecialDefaultCreateBuilder"
	BooleanStringAttributeCreateBuilder           HelperMapping = "booleanStringAttributeCreateBuilder"
	AccountObjectIdentifierAttributeCreateBuilder HelperMapping = "accountObjectIdentifierAttributeCreateBuilder"
)

// ListValueCreateMapping is listValueCreate with the name of the per-element mapper function.
type ListValueCreateMapping struct {
	Mapper string
}

func (m ListValueCreateMapping) CreateCall(objectName string, attr Attribute) CreateCall {
	if m.Mapper == "" {
		panic(fmt.Sprintf("object %q: attribute %q ListValueCreate needs a mapper function", objectName, attr.Name))
	}
	return CreateCall{
		Helper: "listValueCreate",
		Target: fmt.Sprintf("&request.%s, %s", genhelpers.SnakeCaseToCamel(attr.Name), m.Mapper),
	}
}

func ListValueCreate(mapper string) ListValueCreateMapping {
	return ListValueCreateMapping{Mapper: mapper}
}
