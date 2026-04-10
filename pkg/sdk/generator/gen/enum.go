package gen

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// Enum defines an enum type with its name and values. For now, only string values are supported.
type Enum struct {
	Name       string
	NamePlural string
	Values     []string
	Aliases    map[string][]string
}

func NewEnum(name, namePlural string, values ...string) *Enum {
	return &Enum{
		Name:       name,
		NamePlural: namePlural,
		Values:     values,
		Aliases:    make(map[string][]string),
	}
}

// WithAliases adds aliases for a canonical enum value.
// E.g. for canonical value "XSMALL" with aliases "X-SMALL".
func (e *Enum) WithAliases(value string, aliases ...string) *Enum {
	if !slices.Contains(e.Values, value) {
		log.Panicf("WithAliases: value %q is not a declared enum value in %s; declared values: %v", value, e.Name, e.Values)
	}
	e.Aliases[value] = append(e.Aliases[value], aliases...)
	return e
}

// HasAliases returns true if any enum value has at least one alias.
func (e *Enum) HasAliases() bool {
	for _, aliases := range e.Aliases {
		if len(aliases) > 0 {
			return true
		}
	}
	return false
}

// valueName returns the constant name for a given enum value.
// E.g. for type ProgrammaticAccessTokenStatus and value "ACTIVE_VALUE" -> "ProgrammaticAccessTokenStatusActiveValue".
func (e *Enum) valueName(value string) string {
	return e.Name + genhelpers.SnakeCaseToCamel(value)
}

// ValueRepresentations returns all enum values with their Go names, values, and aliases used in the template.
func (e *Enum) ValueRepresentations() []EnumValueRepresentation {
	valueRepresentations := make([]EnumValueRepresentation, len(e.Values))
	for i, value := range e.Values {
		valueRepresentations[i] = EnumValueRepresentation{
			Name:    e.valueName(value),
			Value:   value,
			Aliases: e.Aliases[value],
		}
	}
	return valueRepresentations
}

// EnumValueRepresentation represents a single enum value with its Go constant name, SQL value, and optional aliases.
type EnumValueRepresentation struct {
	Name    string
	Value   string
	Aliases []string
}
