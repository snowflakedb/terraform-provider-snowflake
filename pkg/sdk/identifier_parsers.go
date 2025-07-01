package sdk

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

const IdDelimiter = '.'

// TODO(SNOW-1495053): Temporarily exported, make as private
func ParseIdentifierStringWithOpts(identifier string, opts func(*csv.Reader)) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(identifier))
	if opts != nil {
		opts(reader)
	}
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read identifier: %s, err = %w", identifier, err)
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("incompatible identifier: %s", identifier)
	}
	return lines[0], nil
}

// TODO(SNOW-1495053): Temporarily exported, make as private
func ParseIdentifierString(identifier string) ([]string, error) {
	parts, err := ParseIdentifierStringWithOpts(identifier, func(r *csv.Reader) {
		r.Comma = IdDelimiter
	})
	if err != nil {
		return nil, err
	}
	for _, part := range parts {
		// TODO(SNOW-1571674): Remove the validation
		if strings.Contains(part, `"`) {
			return nil, fmt.Errorf(`unable to parse identifier: %s, currently identifiers containing double quotes are not supported in the provider`, identifier)
		}
	}
	return parts, nil
}

func parseIdentifier[T AccountIdentifier | AccountObjectIdentifier | DatabaseObjectIdentifier | ExternalObjectIdentifier | SchemaObjectIdentifier | TableColumnIdentifier](identifier string, expectedParts int, expectedFormat string, constructFromParts func(parts []string) T) (T, error) {
	var emptyIdentifier T
	parts, err := ParseIdentifierString(identifier)
	if err != nil {
		return emptyIdentifier, err
	}
	if len(parts) != expectedParts {
		return emptyIdentifier, fmt.Errorf(`unexpected number of parts %[1]d in identifier %[2]s, expected %[3]d in a form of "%[4]s"`, len(parts), identifier, expectedParts, expectedFormat)
	}
	return constructFromParts(parts), nil
}

func ParseAccountObjectIdentifier(identifier string) (AccountObjectIdentifier, error) {
	if !(strings.HasPrefix(identifier, `"`) && strings.HasSuffix(identifier, `"`)) {
		identifier = fmt.Sprintf(`"%s"`, identifier)
	}
	return parseIdentifier[AccountObjectIdentifier](
		identifier, 1, "<account_object_name>",
		func(parts []string) AccountObjectIdentifier {
			return NewAccountObjectIdentifier(parts[0])
		},
	)
}

// ParseObjectIdentifierString tries to guess the identifier by the number of parts it contains.
// Because of the overlapping, in some cases, the output ObjectIdentifier can be one of the following implementations:
// - AccountObjectIdentifier for one part
// - DatabaseObjectIdentifier for two parts
// - SchemaObjectIdentifier for three parts (overlaps with ExternalObjectIdentifier)
// - TableColumnIdentifier for four parts
func ParseObjectIdentifierString(identifier string) (ObjectIdentifier, error) {
	parts, err := ParseIdentifierString(identifier)
	if err != nil {
		return nil, err
	}
	switch len(parts) {
	case 1:
		return NewAccountObjectIdentifier(parts[0]), nil
	case 2:
		return NewDatabaseObjectIdentifier(parts[0], parts[1]), nil
	case 3:
		return NewSchemaObjectIdentifier(parts[0], parts[1], parts[2]), nil
	case 4:
		return NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3]), nil
	default:
		return nil, fmt.Errorf("unsupported identifier: %[1]s (number of parts: %[2]d)", identifier, len(parts))
	}
}

func ParseDatabaseObjectIdentifier(identifier string) (DatabaseObjectIdentifier, error) {
	return parseIdentifier[DatabaseObjectIdentifier](
		identifier, 2, "<database_name>.<database_object_name>",
		func(parts []string) DatabaseObjectIdentifier {
			return NewDatabaseObjectIdentifier(parts[0], parts[1])
		},
	)
}

func ParseSchemaObjectIdentifier(identifier string) (SchemaObjectIdentifier, error) {
	return parseIdentifier[SchemaObjectIdentifier](
		identifier, 3, "<database_name>.<schema_name>.<schema_object_name>",
		func(parts []string) SchemaObjectIdentifier {
			return NewSchemaObjectIdentifier(parts[0], parts[1], parts[2])
		},
	)
}

func ParseTableColumnIdentifier(identifier string) (TableColumnIdentifier, error) {
	return parseIdentifier[TableColumnIdentifier](
		identifier, 4, "<database_name>.<schema_name>.<table_name>.<table_column_name>",
		func(parts []string) TableColumnIdentifier {
			return NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3])
		},
	)
}

// ParseAccountIdentifier is implemented with an assumption that the recommended format is used that contains two parts,
// organization name and account name.
func ParseAccountIdentifier(identifier string) (AccountIdentifier, error) {
	return parseIdentifier[AccountIdentifier](
		identifier, 2, "<organization_name>.<account_name>",
		func(parts []string) AccountIdentifier {
			return NewAccountIdentifier(parts[0], parts[1])
		},
	)
}

// ParseExternalObjectIdentifier is implemented with an assumption that the identifier consists of three parts, because:
//   - After identifier rework, we expect account identifiers to always have two parts "<organization_name>.<account_name>".
//   - So far, the only external things that we referred to with external identifiers had only one part (not including the account identifier),
//     meaning it will always be represented as sdk.AccountObjectIdentifier. Documentation also doesn't describe any case where
//     account identifier would be used as part of the identifier that would refer to the "lower level" object.
//     Reference: https://docs.snowflake.com/en/user-guide/admin-account-identifier#where-are-account-identifiers-used.
func ParseExternalObjectIdentifier(identifier string) (ExternalObjectIdentifier, error) {
	return parseIdentifier[ExternalObjectIdentifier](
		identifier, 3, "<organization_name>.<account_name>.<external_object_name>",
		func(parts []string) ExternalObjectIdentifier {
			return NewExternalObjectIdentifier(NewAccountIdentifier(parts[0], parts[1]), NewAccountObjectIdentifier(parts[2]))
		},
	)
}

func ParseSchemaObjectIdentifierWithArguments(fullyQualifiedName string) (SchemaObjectIdentifierWithArguments, error) {
	splitIdIndex := strings.IndexRune(fullyQualifiedName, '(')
	if splitIdIndex == -1 {
		return SchemaObjectIdentifierWithArguments{}, errors.New("unable to parse identifier: '(' not present")
	}
	parts, err := ParseIdentifierString(fullyQualifiedName[:splitIdIndex])
	if err != nil {
		return SchemaObjectIdentifierWithArguments{}, err
	}
	parsedArguments, err := ParseFunctionAndProcedureArguments(fullyQualifiedName[splitIdIndex:])
	if err != nil {
		return SchemaObjectIdentifierWithArguments{}, err
	}
	return NewSchemaObjectIdentifierWithArguments(
		parts[0],
		parts[1],
		parts[2],
		collections.Map(parsedArguments, func(a ParsedArgument) DataType {
			return DataType(a.ArgType)
		})...,
	), nil
}

// ParseSchemaObjectIdentifierWithArgumentsAndReturnType parses names in the following format: <database>.<schema>."<function>(<argname> <argtype>...):<returntype>"
// Return type is not part of an identifier.
// TODO(SNOW-1625030): Remove and use ParseSchemaObjectIdentifierWithArguments instead
func ParseSchemaObjectIdentifierWithArgumentsAndReturnType(fullyQualifiedName string) (SchemaObjectIdentifierWithArguments, error) {
	parts, err := ParseIdentifierStringWithOpts(fullyQualifiedName, func(r *csv.Reader) {
		r.Comma = IdDelimiter
	})
	if err != nil {
		return SchemaObjectIdentifierWithArguments{}, err
	}
	functionHeader := parts[2]
	leftParenthesisIndex := strings.IndexRune(functionHeader, '(')
	functionName := functionHeader[:leftParenthesisIndex]
	argsRaw := functionHeader[leftParenthesisIndex:]
	returnTypeIndex := strings.LastIndex(argsRaw, ":")
	if returnTypeIndex != -1 {
		argsRaw = argsRaw[:returnTypeIndex]
	}
	dataTypes, err := ParseFunctionArgumentsFromString(argsRaw)
	if err != nil {
		return SchemaObjectIdentifierWithArguments{}, err
	}
	return NewSchemaObjectIdentifierWithArguments(
		parts[0],
		parts[1],
		functionName,
		dataTypes...,
	), nil
}

// ParseFunctionArgumentsFromString parses function argument from arguments string with optional argument names.
// Varying types are not supported (e.g. VARCHAR(200)), because Snowflake outputs them in a shortened version
// (VARCHAR in this case). The only exception is newly added type VECTOR that has the following structure
// VECTOR(<type>, n) where <type> right now can be either INT or FLOAT and n is the number of elements in the VECTOR.
// Snowflake returns vectors with their exact type, and this function supports it.
func ParseFunctionArgumentsFromString(arguments string) ([]DataType, error) {
	dataTypes := make([]DataType, 0)

	if len(arguments) > 0 && arguments[0] == '(' && arguments[len(arguments)-1] == ')' {
		arguments = arguments[1 : len(arguments)-1]
	}
	stringBuffer := bytes.NewBufferString(arguments)

	for stringBuffer.Len() > 0 {
		stringBuffer = bytes.NewBufferString(strings.TrimSpace(stringBuffer.String()))

		// When a function is created with a default value for an argument, in the SHOW output ("arguments" column)
		// the argument's data type is prefixed with "DEFAULT ", e.g. "(DEFAULT INT, DEFAULT VARCHAR)".
		if strings.HasPrefix(stringBuffer.String(), "DEFAULT") {
			if _, err := stringBuffer.ReadString(' '); err != nil {
				return nil, fmt.Errorf("failed to skip default keyword, err = %w", err)
			}
		}

		// We use another buffer to peek into next data type (needed for vector parsing)
		peekDataType, _ := bytes.NewBufferString(stringBuffer.String()).ReadString(',')
		// Skip argument name, if present
		if strings.ContainsRune(peekDataType, ' ') && !strings.HasPrefix(peekDataType, "VECTOR(") {
			// Ignore err, because stringBuffer always contains ' ' here
			_, _ = stringBuffer.ReadString(' ')
			peekDataType, _ = bytes.NewBufferString(stringBuffer.String()).ReadString(',')
		}

		switch {
		// For now, only vectors need special parsing behavior
		case strings.HasPrefix(peekDataType, "VECTOR"):
			vectorDataType, err := stringBuffer.ReadString(')')
			if err != nil {
				return nil, fmt.Errorf("failed to parse vector type, couldn't find the closing bracket, err = %w", err)
			}

			vectorDataType = strings.TrimSpace(vectorDataType)
			vectorTypeBuffer := bytes.NewBufferString(vectorDataType)
			if _, err := vectorTypeBuffer.ReadString('('); err != nil {
				return nil, fmt.Errorf("failed to parse vector type, couldn't find the opening bracket, err = %w", err)
			}

			vectorInnerType, err := vectorTypeBuffer.ReadString(',')
			if err != nil {
				return nil, fmt.Errorf("failed to parse vector inner type: %w", err)
			}

			vectorInnerType = vectorInnerType[:len(vectorInnerType)-1]
			if !slices.Contains(allowedVectorInnerTypes, DataType(vectorInnerType)) {
				return nil, fmt.Errorf("invalid vector inner type: %s, allowed vector types are: %v", vectorInnerType, allowedVectorInnerTypes)
			}

			vectorSize, err := vectorTypeBuffer.ReadString(')')
			if err != nil {
				return nil, fmt.Errorf("failed to parse vector size: %w", err)
			}

			vectorSize = strings.TrimSpace(vectorSize[:len(vectorSize)-1])
			_, err = strconv.ParseInt(vectorSize, 0, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid vector size: %s (not a number): %w", vectorSize, err)
			}

			if stringBuffer.Len() > 0 {
				commaByte, err := stringBuffer.ReadByte()
				if commaByte != ',' {
					return nil, fmt.Errorf("expected a comma delimited string but found %s", string(commaByte))
				}
				if err != nil {
					return nil, err
				}
			}
			dataTypes = append(dataTypes, DataType(vectorDataType))
		default:
			argument, err := stringBuffer.ReadString(',')
			if err == nil {
				argument = argument[:len(argument)-1]
			}
			dataTypes = append(dataTypes, DataType(argument))
		}
	}

	return dataTypes, nil
}
