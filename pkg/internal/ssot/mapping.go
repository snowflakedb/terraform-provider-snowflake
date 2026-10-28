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

// AttributeMappedValueCreate is attributeMappedValueCreateBuilder.
// The value is the mapper function name.
type AttributeMappedValueCreate string

func (m AttributeMappedValueCreate) CreateCall(objectName string, attr Attribute) CreateCall {
	if m == "" {
		panic(fmt.Sprintf("object %q: attribute %q attributeMappedValueCreateBuilder needs a mapper function", objectName, attr.Name))
	}
	return CreateCall{
		Helper: "attributeMappedValueCreateBuilder",
		Target: fmt.Sprintf("request.With%s, %s", genhelpers.SnakeCaseToCamel(attr.Name), string(m)),
	}
}

// UpdateMapping names the resource update helper for one attribute.
// An attribute with a nil UpdateMapping is left out of set/unset.
type UpdateMapping interface {
	UpdateCall(objectName string, attr Attribute) UpdateCall
}

// UpdateCall is one generated ApplySetUnset call.
// Target is the argument after the attribute key.
type UpdateCall struct {
	Helper string
	Target string
}

// SetOnlyMapping is an update helper called with set.With<Field>.
// The value is the helper function name.
type SetOnlyMapping string

func (m SetOnlyMapping) UpdateCall(_ string, attr Attribute) UpdateCall {
	return UpdateCall{
		Helper: string(m),
		Target: "set.With" + genhelpers.SnakeCaseToCamel(attr.Name),
	}
}

const IntAttributeUpdateSetOnlyBuilder SetOnlyMapping = "intAttributeUpdateSetOnlyBuilder"

// SetUnsetMapping is an update helper called with set.With<Field>, unset.With<Field>.
// The value is the helper function name.
type SetUnsetMapping string

func (m SetUnsetMapping) UpdateCall(_ string, attr Attribute) UpdateCall {
	goName := genhelpers.SnakeCaseToCamel(attr.Name)
	return UpdateCall{
		Helper: string(m),
		Target: fmt.Sprintf("set.With%s, unset.With%s", goName, goName),
	}
}

const (
	StringAttributeUpdateBuilder                SetUnsetMapping = "stringAttributeUpdateBuilder"
	IntAttributeUpdateBuilder                   SetUnsetMapping = "intAttributeUpdateBuilder"
	IntAttributeWithSpecialDefaultUpdateBuilder SetUnsetMapping = "intAttributeWithSpecialDefaultUpdateBuilder"
	BooleanStringAttributeUpdateBuilder         SetUnsetMapping = "booleanStringAttributeUpdateBuilder"
)

// ListValueUpdateMapping is listValueUpdateBuilder with the name of the per-element mapper function.
type ListValueUpdateMapping struct {
	Mapper string
}

func (m ListValueUpdateMapping) UpdateCall(objectName string, attr Attribute) UpdateCall {
	if m.Mapper == "" {
		panic(fmt.Sprintf("object %q: attribute %q ListValueUpdate needs a mapper function", objectName, attr.Name))
	}
	goName := genhelpers.SnakeCaseToCamel(attr.Name)
	return UpdateCall{
		Helper: "listValueUpdateBuilder",
		Target: fmt.Sprintf("set.With%s, unset.With%s, %s", goName, goName, m.Mapper),
	}
}

func ListValueUpdate(mapper string) ListValueUpdateMapping {
	return ListValueUpdateMapping{Mapper: mapper}
}
