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
// - additional validations invocation - present here, emits a call to additionalValidations() method on the struct
type ValidationType int64

const (
	ValidIdentifier ValidationType = iota
	ValidIdentifierIfSet
	ConflictingFields
	MoreThanOneValueSet
	ExactlyOneValueSet
	AtLeastOneValueSet
	ValidateValue
	ValidateValueSet
	AdditionalValidations
	NoDoubleDollarQuotes
	NoDoubleDollarQuotesIfSet
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

func (v *Validation) IsAdditionalValidations() bool {
	return v.Type == AdditionalValidations
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

func (v *Validation) fieldsInSlicePath(elemVar string) []string {
	params := make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = fmt.Sprintf("%s.%s", elemVar, s)
	}
	return params
}

func (v *Validation) Condition(field *Field) string {
	var fieldNamesProvider func(*Field) []string
	if field.IsSlice() {
		fieldNamesProvider = func(f *Field) []string {
			return v.fieldsInSlicePath(f.SliceElemVar())
		}
	} else {
		fieldNamesProvider = func(f *Field) []string {
			return v.fieldsWithPath(f)
		}
	}

	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("!ValidObjectIdentifier(%s)", strings.Join(fieldNamesProvider(field), ","))
	case ValidIdentifierIfSet:
		return fmt.Sprintf("%s != nil && !ValidObjectIdentifier(%s)", strings.Join(fieldNamesProvider(field), ","), strings.Join(fieldNamesProvider(field), ","))
	case ConflictingFields:
		return fmt.Sprintf("everyValueSet(%s)", strings.Join(fieldNamesProvider(field), ","))
	case MoreThanOneValueSet:
		return fmt.Sprintf("moreThanOneValueSet(%s)", strings.Join(fieldNamesProvider(field), ","))
	case ExactlyOneValueSet:
		return fmt.Sprintf("!exactlyOneValueSet(%s)", strings.Join(fieldNamesProvider(field), ","))
	case AtLeastOneValueSet:
		return fmt.Sprintf("!anyValueSet(%s)", strings.Join(fieldNamesProvider(field), ","))
	case ValidateValueSet:
		return fmt.Sprintf("!valueSet(%s)", strings.Join(fieldNamesProvider(field), ","))
	case ValidateValue:
		if len(v.FieldNames) != 1 {
			log.Panicf("expected ValidateValue to be called exactly one field, got: %v", v.FieldNames)
		}
		return fmt.Sprintf("err := %s.validate(); err != nil", fieldNamesProvider(field)[0])
	case AdditionalValidations:
		log.Panicf("Condition() must not be called for AdditionalValidations type")
	case NoDoubleDollarQuotes:
		if len(v.FieldNames) != 1 {
			log.Panicf("expected NoDoubleDollarQuotes to be called with exactly one field, got: %v", v.FieldNames)
		}
		return fmt.Sprintf("containsDoubleDollarQuotes(%s)", fieldNamesProvider(field)[0])
	case NoDoubleDollarQuotesIfSet:
		if len(v.FieldNames) != 1 {
			log.Panicf("expected NoDoubleDollarQuotesIfSet to be called with exactly one field, got: %v", v.FieldNames)
		}
		fieldPath := fieldNamesProvider(field)[0]
		return fmt.Sprintf("%s != nil && containsDoubleDollarQuotes(*%s)", fieldPath, fieldPath)
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
	case MoreThanOneValueSet:
		return fmt.Sprintf(`errMoreThanOneOf("%s", %s)`, field.PathWithRoot(), strings.Join(v.paramsQuoted(), ","))
	case ExactlyOneValueSet:
		return fmt.Sprintf(`errExactlyOneOf("%s", %s)`, field.PathWithRoot(), strings.Join(v.paramsQuoted(), ","))
	case AtLeastOneValueSet:
		return fmt.Sprintf(`errAtLeastOneOf("%s", %s)`, field.PathWithRoot(), strings.Join(v.paramsQuoted(), ","))
	case ValidateValueSet:
		return fmt.Sprintf(`errNotSet("%s", %s)`, field.PathWithRoot(), strings.Join(v.paramsQuoted(), ","))
	case ValidateValue:
		return "err"
	case AdditionalValidations:
		log.Panicf("ReturnedError() must not be called for AdditionalValidations type")
	case NoDoubleDollarQuotes, NoDoubleDollarQuotesIfSet:
		return fmt.Sprintf(`errDoubleDollarQuotesNotAllowed("%s", "%s")`, field.PathWithRoot(), v.FieldNames[0])
	}
	panic("condition for validation unknown")
}

func (v *Validation) TodoComment(field *Field) string {
	switch v.Type {
	case AdditionalValidations:
		return fmt.Sprintf("validation: additional validations for opts%s", field.Path())
	case ValidIdentifier:
		return fmt.Sprintf("validation: valid identifier for %v", v.fieldsWithPath(field))
	case ValidIdentifierIfSet:
		return fmt.Sprintf("validation: valid identifier for %v if set", v.fieldsWithPath(field))
	case ConflictingFields:
		return fmt.Sprintf("validation: conflicting fields for %v", v.fieldsWithPath(field))
	case MoreThanOneValueSet:
		return fmt.Sprintf("validation: more than one field from %v cannot be set at the same time", v.fieldsWithPath(field))
	case ExactlyOneValueSet:
		return fmt.Sprintf("validation: exactly one field from %v should be present", v.fieldsWithPath(field))
	case AtLeastOneValueSet:
		return fmt.Sprintf("validation: at least one of the fields %v should be set", v.fieldsWithPath(field))
	case ValidateValueSet:
		return fmt.Sprintf("validation: %v should be set", v.fieldsWithPath(field))
	case ValidateValue:
		return fmt.Sprintf("validation: %v should be valid", v.fieldsWithPath(field)[0])
	case NoDoubleDollarQuotes:
		return fmt.Sprintf("validation: %v must not contain $$", v.fieldsWithPath(field)[0])
	case NoDoubleDollarQuotesIfSet:
		return fmt.Sprintf("validation: %v must not contain $$ if set", v.fieldsWithPath(field)[0])
	}
	panic("condition for validation unknown")
}
