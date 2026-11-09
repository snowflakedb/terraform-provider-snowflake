package ssot

import (
	"fmt"
	"strconv"

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
	BoolAttributeCreateBuilder                    HelperMapping = "boolAttributeCreateBuilder"
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

// SetMappedValueCreate is setMappedValueCreateBuilder with the name of the per-element mapper function.
type SetMappedValueCreateMapping struct {
	Mapper string
}

func (m SetMappedValueCreateMapping) CreateCall(objectName string, attr Attribute) CreateCall {
	if m.Mapper == "" {
		panic(fmt.Sprintf("object %q: attribute %q SetMappedValueCreate needs a mapper function", objectName, attr.Name))
	}
	return CreateCall{
		Helper: "setMappedValueCreateBuilder",
		Target: fmt.Sprintf("request.With%s, %s", genhelpers.SnakeCaseToCamel(attr.Name), m.Mapper),
	}
}

func SetMappedValueCreate(mapper string) SetMappedValueCreateMapping {
	return SetMappedValueCreateMapping{Mapper: mapper}
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
// Set and Unset say which setters the helper receives from the SET and UNSET request that holds the attribute: the
// With<Field> methods, or with Fields the request fields (&set.<Field>). Mapper is the extra argument after the setters,
// if any, such as a literal. The generator writes the final arguments, because only it knows which request each
// attribute belongs to.
type UpdateCall struct {
	Helper string
	Set    bool
	Unset  bool
	Fields bool
	// Field is the name of the request field or setter when it differs from the attribute, such as SetBody for body.
	Field  string
	Mapper string
}

// SetOnlyMapping is an update helper called with the SET setter.
// The value is the helper function name.
type SetOnlyMapping string

func (m SetOnlyMapping) UpdateCall(_ string, attr Attribute) UpdateCall {
	return UpdateCall{
		Helper: string(m),
		Set:    true,
	}
}

const (
	IntAttributeUpdateSetOnlyBuilder     SetOnlyMapping = "intAttributeUpdateSetOnlyBuilder"
	BooleanAttributeUpdateSetOnlyBuilder SetOnlyMapping = "booleanAttributeUpdateSetOnlyBuilder"
)

// SetUnsetMapping is an update helper called with the SET and UNSET setters.
// The value is the helper function name.
type SetUnsetMapping string

func (m SetUnsetMapping) UpdateCall(_ string, attr Attribute) UpdateCall {
	return UpdateCall{
		Helper: string(m),
		Set:    true,
		Unset:  true,
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
	return UpdateCall{
		Helper: "listValueUpdateBuilder",
		Set:    true,
		Unset:  true,
		Mapper: m.Mapper,
	}
}

func ListValueUpdate(mapper string) ListValueUpdateMapping {
	return ListValueUpdateMapping{Mapper: mapper}
}

// ListValueUpdateSetOnlyMapping is listValueUpdateSetOnlyBuilder. The list can be set and not unset.
type ListValueUpdateSetOnlyMapping struct {
	Mapper string
}

func (m ListValueUpdateSetOnlyMapping) UpdateCall(objectName string, attr Attribute) UpdateCall {
	if m.Mapper == "" {
		panic(fmt.Sprintf("object %q: attribute %q ListValueUpdateSetOnly needs a mapper function", objectName, attr.Name))
	}
	return UpdateCall{
		Helper: "listValueUpdateSetOnlyBuilder",
		Set:    true,
		Mapper: m.Mapper,
	}
}

func ListValueUpdateSetOnly(mapper string) ListValueUpdateSetOnlyMapping {
	return ListValueUpdateSetOnlyMapping{Mapper: mapper}
}

// SetUnsetFieldsMapping is an update helper called with the fields of the SET and UNSET requests,
// &set.<Field> and &unset.<Field>. The value is the helper function name.
type SetUnsetFieldsMapping string

func (m SetUnsetFieldsMapping) UpdateCall(_ string, _ Attribute) UpdateCall {
	return UpdateCall{
		Helper: string(m),
		Set:    true,
		Unset:  true,
		Fields: true,
	}
}

const (
	AccountObjectIdentifierAttributeUpdate SetUnsetFieldsMapping = "accountObjectIdentifierAttributeUpdate"
)

// SetOnlyFieldsMapping is an update helper called with the field of the SET request, &set.<Field>. It never unsets the
// attribute. The value is the helper function name.
type SetOnlyFieldsMapping string

func (m SetOnlyFieldsMapping) UpdateCall(_ string, _ Attribute) UpdateCall {
	return UpdateCall{
		Helper: string(m),
		Set:    true,
		Fields: true,
	}
}

const (
	BooleanStringAttributeUpdateSetOnly SetOnlyFieldsMapping = "booleanStringAttributeUpdateSetOnly"
)

// AttributeMappedValueUpdateMapping is attributeMappedValueUpdate with the name of the function that turns the Terraform
// value into the SDK value. An attribute that is removed from the configuration is unset. It also serves a nested
// attribute that is updated as one value, such as a block that the SDK takes as a string.
type AttributeMappedValueUpdateMapping struct {
	Mapper string
}

func (m AttributeMappedValueUpdateMapping) UpdateCall(objectName string, attr Attribute) UpdateCall {
	if m.Mapper == "" {
		panic(fmt.Sprintf("object %q: attribute %q AttributeMappedValueUpdate needs a mapper function", objectName, attr.Name))
	}
	return UpdateCall{
		Helper: "attributeMappedValueUpdate",
		Set:    true,
		Unset:  true,
		Fields: true,
		Mapper: m.Mapper,
	}
}

func AttributeMappedValueUpdate(mapper string) AttributeMappedValueUpdateMapping {
	return AttributeMappedValueUpdateMapping{Mapper: mapper}
}

// WholeValueUpdateMapping is implemented by an UpdateMapping that updates a nested attribute as one value. The generator
// skips every other nested attribute, because their fields are not updated one by one yet.
type WholeValueUpdateMapping interface {
	UpdateMapping
	updatesWholeValue()
}

func (AttributeMappedValueUpdateMapping) updatesWholeValue() {}

// UpdatesWholeValue reports whether the mapping updates a nested attribute as one value.
func UpdatesWholeValue(mapping UpdateMapping) bool {
	_, ok := mapping.(WholeValueUpdateMapping)
	return ok
}

// StringAttributeUpdateSetOnlyNotEmptyMapping is stringAttributeUpdateSetOnlyNotEmpty: the value is set when it changes,
// and it is never unset. Field is the request field when it differs from the attribute.
type StringAttributeUpdateSetOnlyNotEmptyMapping struct {
	Field string
}

func (m StringAttributeUpdateSetOnlyNotEmptyMapping) UpdateCall(_ string, _ Attribute) UpdateCall {
	return UpdateCall{
		Helper: "stringAttributeUpdateSetOnlyNotEmpty",
		Set:    true,
		Fields: true,
		Field:  m.Field,
	}
}

func StringAttributeUpdateSetOnlyNotEmpty(field string) StringAttributeUpdateSetOnlyNotEmptyMapping {
	return StringAttributeUpdateSetOnlyNotEmptyMapping{Field: field}
}

// BooleanStringAttributeUnsetFallbackUpdateMapping is booleanStringAttributeUnsetFallbackUpdate. A boolean string that is
// removed from the configuration is set to Fallback, because the object has no UNSET for it.
type BooleanStringAttributeUnsetFallbackUpdateMapping struct {
	Fallback bool
}

func (m BooleanStringAttributeUnsetFallbackUpdateMapping) UpdateCall(_ string, _ Attribute) UpdateCall {
	return UpdateCall{
		Helper: "booleanStringAttributeUnsetFallbackUpdate",
		Set:    true,
		Fields: true,
		Mapper: strconv.FormatBool(m.Fallback),
	}
}

func BooleanStringAttributeUnsetFallbackUpdate(fallback bool) BooleanStringAttributeUnsetFallbackUpdateMapping {
	return BooleanStringAttributeUnsetFallbackUpdateMapping{Fallback: fallback}
}
