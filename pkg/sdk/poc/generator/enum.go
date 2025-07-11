package generator

import (
	"strings"
)

// Enum defines an enum type with its name and values. For now, only string values are supported.
type Enum struct {
	Name       string
	NamePlural string
	Values     []string
}

func NewEnum(name string, values ...string) *Enum {
	return &Enum{
		Name:   name,
		Values: values,
	}
}

func (e *Enum) WithPlural(plural string) *Enum {
	e.NamePlural = plural
	return e
}

// TypeName returns the Go type name for the enum.
func (e *Enum) TypeName() string {
	return e.Name
}

// valueName returns the constant name for a given enum value.
// E.g. for type ProgrammaticAccessTokenStatus and value "ACTIVE_VALUE" -> "ProgrammaticAccessTokenStatusActiveValue".
func (e *Enum) valueName(value string) string {
	normalized := strings.ToLower(value)
	parts := strings.Split(normalized, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return e.Name + strings.Join(parts, "")
}

// AllValuesSliceName returns the name of the slice containing all enum values
func (e *Enum) AllValuesSliceName() string {
	if e.NamePlural == "" {
		return "all" + e.Name + "s"
	}
	return "all" + e.NamePlural
}

// ConverterFunctionName returns the name of the converter function
func (e *Enum) ConverterFunctionName() string {
	return "to" + e.Name
}

// ValueRepresentations returns all enum values with their Go names and values used in the template.
func (e *Enum) ValueRepresentations() []EnumValueRepresentation {
	valueRepresentations := make([]EnumValueRepresentation, len(e.Values))
	for i, value := range e.Values {
		valueRepresentations[i] = EnumValueRepresentation{
			Name:  e.valueName(value),
			Value: value,
		}
	}
	return valueRepresentations
}

// EnumValueRepresentation represents a single enum value.
type EnumValueRepresentation struct {
	Name  string
	Value string
}
