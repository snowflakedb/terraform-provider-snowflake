package ssot

import (
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// Enum defines an enum type with its name and values. For now, only string values are supported.
type Enum struct {
	Name       string
	NamePlural string
	Values     []string
	Aliases    map[string][]string
	// explicitNames maps a SQL value to a constant suffix supplied by WithValue.
	// The suffix is concatenated as-is, without SnakeCaseToCamel.
	explicitNames map[string]string
	// SdkType is the sdk Go type when it differs from Name.
	// Empty uses Name. Resource schema validation still calls sdk.To<Name>.
	SdkType string
}

// WithSdkType sets the sdk Go type when Name is only the sdk.To<Name> parser suffix.
func (e *Enum) WithSdkType(sdkType string) *Enum {
	e.SdkType = sdkType
	return e
}

func (e *Enum) sdkTypeName() string {
	if e.SdkType != "" {
		return e.SdkType
	}
	return e.Name
}

func NewEnum(name, namePlural string, values ...string) *Enum {
	return &Enum{
		Name:       name,
		NamePlural: namePlural,
		Values:     values,
		Aliases:    make(map[string][]string),
	}
}

// WithValue appends an enum value whose Go constant suffix is fieldName.
// The suffix is concatenated as-is (Enum name + fieldName), skipping SnakeCaseToCamel.
// Use this when the SQL value cannot produce the desired constant name, including the empty string.
func (e *Enum) WithValue(fieldName, value string) *Enum {
	if fieldName == "" {
		log.Panicf("WithValue: field name must not be empty for enum %s", e.Name)
	}
	if slices.Contains(e.Values, value) {
		log.Panicf("WithValue: value %q is already a declared enum value in %s; declared values: %v", value, e.Name, e.Values)
	}
	constantName := e.Name + fieldName
	for _, existing := range e.Values {
		if e.valueName(existing) == constantName {
			log.Panicf("WithValue: constant name %s is already used by value %q in enum %s", constantName, existing, e.Name)
		}
	}
	if e.explicitNames == nil {
		e.explicitNames = make(map[string]string)
	}
	e.explicitNames[value] = fieldName
	e.Values = append(e.Values, value)
	return e
}

// WithEmptyValue appends an enum value whose SQL form is the empty string.
// fieldName is the constant suffix passed through to WithValue.
func (e *Enum) WithEmptyValue(fieldName string) *Enum {
	return e.WithValue(fieldName, "")
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

// IsLowercase returns true if all enum values are lowercase.
// When true, the generated conversion function uses strings.ToLower instead of strings.ToUpper.
func (e *Enum) IsLowercase() bool {
	for _, v := range e.Values {
		if v != strings.ToLower(v) {
			return false
		}
	}
	return true
}

// IsMixedCase returns true if any enum value contains both upper and lowercase characters.
// When true, the generated conversion function wraps case expressions with the normalization function
// to ensure case-insensitive matching works correctly.
func (e *Enum) IsMixedCase() bool {
	for _, v := range e.Values {
		if v != strings.ToUpper(v) && v != strings.ToLower(v) {
			return true
		}
	}
	return false
}

// Kind should be used in SDK object definitions instead of KindOfT
func (e *Enum) Kind() string {
	return e.Name
}

// KindPtr should be used in SDK object definitions instead of KindOfTPointer
func (e *Enum) KindPtr() string {
	return "*" + e.Kind()
}

// valueName returns the constant name for a given enum value.
// E.g. for type ProgrammaticAccessTokenStatus and value "ACTIVE_VALUE" -> "ProgrammaticAccessTokenStatusActiveValue".
// A lowercase-to-uppercase boundary is split first, so "GeoJSON" becomes "GeoJson" while the SQL value stays "GeoJSON".
func (e *Enum) valueName(value string) string {
	if fieldName, ok := e.explicitNames[value]; ok {
		return e.Name + fieldName
	}
	value = strings.ReplaceAll(value, " ", "_")
	value = strings.ReplaceAll(value, "-", "_")
	value = splitCamelBoundary(value)
	return e.Name + genhelpers.SnakeCaseToCamel(value)
}

// splitCamelBoundary inserts an underscore between a lowercase letter and a following uppercase letter.
func splitCamelBoundary(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	for i := 0; i < len(value); i++ {
		if i > 0 && value[i-1] >= 'a' && value[i-1] <= 'z' && value[i] >= 'A' && value[i] <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteByte(value[i])
	}
	return b.String()
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
