package generator

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/stringhelpers"

// Enum defines an enum type with its name and values. For now, only string values are supported.
// Limitations (also added to the README.md):
// Generate unit tests for the enum type converters
// Handle synonyms (e.g. [`ToWarehouseSize`](https://github.com/snowflakedb/terraform-provider-snowflake/blob/5bdcd127d9288212b10ea7b138bebc0cb770c5b9/pkg/sdk/warehouses.go#L77))
type Enum struct {
	Name       string
	NamePlural string
	Values     []string
}

func NewEnum(name, namePlural string, values ...string) *Enum {
	return &Enum{
		Name:       name,
		NamePlural: namePlural,
		Values:     values,
	}
}

// valueName returns the constant name for a given enum value.
// E.g. for type ProgrammaticAccessTokenStatus and value "ACTIVE_VALUE" -> "ProgrammaticAccessTokenStatusActiveValue".
func (e *Enum) valueName(value string) string {
	return e.Name + stringhelpers.SnakeCaseToCamel(value)
}

// AllValuesSliceName returns the name of the slice containing all enum values
func (e *Enum) AllValuesSliceName() string {
	return "All" + e.NamePlural
}

// ConverterFunctionName returns the name of the converter function
func (e *Enum) ConverterFunctionName() string {
	return "To" + e.Name
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
