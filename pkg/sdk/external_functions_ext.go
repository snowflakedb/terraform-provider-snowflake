package sdk

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (v *ExternalFunction) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, v.Arguments...)
}

func (r externalFunctionRow) additionalConvert(result *ExternalFunction) error {
	if r.SchemaName.Valid {
		result.SchemaName = strings.Trim(r.SchemaName.String, `"`)
	}
	if r.CatalogName.Valid {
		result.CatalogName = strings.Trim(r.CatalogName.String, `"`)
	}
	arguments := strings.TrimLeft(r.Arguments, r.Name)
	returnIndex := strings.Index(arguments, ") RETURN ")
	parsedArguments, err := ParseFunctionAndProcedureArguments(arguments[:returnIndex+1])
	if err != nil {
		return fmt.Errorf("failed to parse external function arguments: %w", err)
	}
	result.Arguments = collections.Map(parsedArguments, func(a ParsedArgument) DataType {
		return DataType(a.ArgType)
	})
	return nil
}
