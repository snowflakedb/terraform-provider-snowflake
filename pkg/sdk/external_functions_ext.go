package sdk

import (
	"fmt"
	"strings"
)

func (v *ExternalFunction) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArgumentsNormalized(v.CatalogName, v.SchemaName, v.Name, v.Arguments...)
}

// ArgumentsStrings returns the SQL of each parsed argument type, falling back to the raw SHOW value when parsing failed.
func (v *ExternalFunction) ArgumentsStrings() []string {
	result := make([]string, len(v.ArgumentsRaw))
	for i, raw := range v.ArgumentsRaw {
		if i < len(v.Arguments) && v.Arguments[i] != nil {
			result[i] = v.Arguments[i].ToSql()
		} else {
			result[i] = raw
		}
	}
	return result
}

// ReturnTypeString returns the SQL of the parsed return type, falling back to the raw SHOW value when parsing failed.
func (v *ExternalFunction) ReturnTypeString() string {
	if v.ReturnType != nil {
		return v.ReturnType.ToSql()
	}
	return v.ReturnTypeRaw
}

func (r externalFunctionRow) additionalConvert(result *ExternalFunction) error {
	if r.SchemaName.Valid {
		result.SchemaName = strings.Trim(r.SchemaName.String, `"`)
	}
	if r.CatalogName.Valid {
		result.CatalogName = strings.Trim(r.CatalogName.String, `"`)
	}
	showArguments, err := parseFunctionOrProcedureShowArguments(r.Name, r.Arguments)
	if err != nil {
		return fmt.Errorf("failed to parse external function arguments: %w", err)
	}
	result.Arguments = showArguments.Arguments
	result.ArgumentsRaw = showArguments.ArgumentsRaw
	result.ReturnType = showArguments.ReturnType
	result.ReturnTypeRaw = showArguments.ReturnTypeRaw
	return nil
}
