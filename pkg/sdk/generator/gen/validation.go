package gen

import (
	"fmt"
	"log"
	"strings"
)

// ValidationType contains all handled validation types. Below validations are marked to be contained here or not:
// - opts not nil - not present here, handled on template level
// - valid identifier - present here, for now put on level containing given field
// - conflicting fields - present here, put on level containing given fields
// - exactly one value set - present here, put on level containing given fields
// - at least one value set - present here, put on level containing given fields
// - validate nested field - present here, used for common structs which have their own validate() methods specified
// - nested validation conditionally - not present here, handled by putting validations on lower level fields
type ValidationType int64

const (
	ValidIdentifier ValidationType = iota
	ValidIdentifierIfSet
	ConflictingFields
	ExactlyOneValueSet
	AtLeastOneValueSet
	ValidateValue
	ValidateValueSet
)

type Validation struct {
	Type       ValidationType
	FieldNames []string
}

func NewValidation(validationType ValidationType, fieldNames ...string) *Validation {
	return &Validation{
		Type:       validationType,
		FieldNames: fieldNames,
	}
}

func (v *Validation) paramsQuoted() []string {
	params := make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = wrapWith(s, `"`)
	}
	return params
}

func (v *Validation) fieldsWithPath(field *Field) []string {
	params := make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = fmt.Sprintf("opts%s.%s", field.Path(), s)
	}
	return params
}

func (v *Validation) Condition(field *Field) string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("!ValidObjectIdentifier(%s)", strings.Join(v.fieldsWithPath(field), ","))
	case ValidIdentifierIfSet:
		return fmt.Sprintf("%s != nil && !ValidObjectIdentifier(%s)", strings.Join(v.fieldsWithPath(field), ","), strings.Join(v.fieldsWithPath(field), ","))
	case ConflictingFields:
		return fmt.Sprintf("everyValueSet(%s)", strings.Join(v.fieldsWithPath(field), ","))
	case ExactlyOneValueSet:
		return fmt.Sprintf("!exactlyOneValueSet(%s)", strings.Join(v.fieldsWithPath(field), ","))
	case AtLeastOneValueSet:
		return fmt.Sprintf("!anyValueSet(%s)", strings.Join(v.fieldsWithPath(field), ","))
	case ValidateValueSet:
		return fmt.Sprintf("!valueSet(%s)", strings.Join(v.fieldsWithPath(field), ","))
	case ValidateValue:
		if len(v.FieldNames) != 1 {
			log.Panicf("expected ValidateValue to be called exactly one field, got: %v", v.FieldNames)
		}
		return fmt.Sprintf("err := %s.validate(); err != nil", v.fieldsWithPath(field)[0])
	}
	panic("condition for validation unknown")
}

func (v *Validation) ReturnedError(field *Field) string {
	switch v.Type {
	case ValidIdentifier:
		return "ErrInvalidObjectIdentifier"
	case ValidIdentifierIfSet:
		return "ErrInvalidObjectIdentifier"
	case ConflictingFields:
		return fmt.Sprintf(`errOneOf("%s", %s)`, field.PathWithRoot(), strings.Join(v.paramsQuoted(), ","))
	case ExactlyOneValueSet:
		return fmt.Sprintf(`errExactlyOneOf("%s", %s)`, field.PathWithRoot(), strings.Join(v.paramsQuoted(), ","))
	case AtLeastOneValueSet:
		return fmt.Sprintf(`errAtLeastOneOf("%s", %s)`, field.PathWithRoot(), strings.Join(v.paramsQuoted(), ","))
	case ValidateValueSet:
		return fmt.Sprintf(`errNotSet("%s", %s)`, field.PathWithRoot(), strings.Join(v.paramsQuoted(), ","))
	case ValidateValue:
		return "err"
	}
	panic("condition for validation unknown")
}

func (v *Validation) TodoComment(field *Field) string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("validation: valid identifier for %v", v.fieldsWithPath(field))
	case ValidIdentifierIfSet:
		return fmt.Sprintf("validation: valid identifier for %v if set", v.fieldsWithPath(field))
	case ConflictingFields:
		return fmt.Sprintf("validation: conflicting fields for %v", v.fieldsWithPath(field))
	case ExactlyOneValueSet:
		return fmt.Sprintf("validation: exactly one field from %v should be present", v.fieldsWithPath(field))
	case AtLeastOneValueSet:
		return fmt.Sprintf("validation: at least one of the fields %v should be set", v.fieldsWithPath(field))
	case ValidateValueSet:
		return fmt.Sprintf("validation: %v should be set", v.fieldsWithPath(field))
	case ValidateValue:
		return fmt.Sprintf("validation: %v should be valid", v.fieldsWithPath(field)[0])
	}
	panic("condition for validation unknown")
}
