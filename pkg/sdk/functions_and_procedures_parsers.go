package sdk

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type ParsedArgument struct {
	IsDefault bool
	ArgName   string
	ArgType   string
}

// ParseFunctionAndProcedureArguments parses function argument from arguments string with optional argument names.
// Varying types are supported (e.g. VARCHAR(200)), because Snowflake started to output them with 2025_03 Bundle.
// https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03/bcr-1944
// The input format for this function is: (DEFAULT argName argType(x, ...), ...), where:
// - enclosing parentheses are optional
// - DEFAULT string is optional
// - argName is optional and should not contain commas, parentheses, or spaces
// - argType has optional attributes specifying the type; they are not empty and comma-separated
// - various spaces around the whole string, DEFAULT, argName, argType, and attributes are allowed to some extent
func ParseFunctionAndProcedureArguments(arguments string) ([]ParsedArgument, error) {
	log.Printf("[DEBUG] Parsing arguments string: `%s`", arguments)
	arguments = strings.TrimSpace(arguments)
	if len(arguments) > 0 && arguments[0] == '(' && arguments[len(arguments)-1] == ')' {
		arguments = arguments[1 : len(arguments)-1]
	}
	arguments = strings.TrimSpace(arguments)

	args, err := splitArgs(arguments)
	if err != nil {
		return nil, fmt.Errorf("arguments string %s could not be parsed; %w", arguments, err)
	}
	parsedArguments := make([]ParsedArgument, 0)
	for _, arg := range args {
		log.Printf("[DEBUG] Processing arg: `%s`", arg)
		arg = strings.TrimSpace(arg)
		parsedArgument := ParsedArgument{}

		// When a function is created with a default value for an argument, in the SHOW output ("arguments" column)
		// the argument's data type is prefixed with "DEFAULT ", e.g. "(DEFAULT INT, DEFAULT VARCHAR)".
		arg, parsedArgument.IsDefault = strings.CutPrefix(arg, "DEFAULT ")
		arg = strings.TrimSpace(arg)

		// argName is optional
		firstSpaceIdx := strings.Index(arg, " ")
		firstParenthesesIdx := strings.Index(arg, "(")
		spacePresent := func() bool {
			return firstSpaceIdx != -1
		}
		noParen := func() bool {
			return firstParenthesesIdx == -1
		}
		firstSpaceBeforeParen := func() bool {
			return firstParenthesesIdx != -1 && firstSpaceIdx < firstParenthesesIdx
		}
		if spacePresent() && (firstSpaceBeforeParen() || noParen()) {
			parsedArgument.ArgName = arg[:firstSpaceIdx]
			arg = arg[firstSpaceIdx+1:]
		}

		// the rest is just argType
		parsedArgument.ArgType = strings.TrimSpace(arg)
		parsedArguments = append(parsedArguments, parsedArgument)
	}

	log.Printf("[DEBUG] Parsing arguments string resulted in: %v", parsedArguments)
	return parsedArguments, nil
}

func splitArgs(arguments string) ([]string, error) {
	var args []string
	start := 0
	level := 0

	if arguments == "" {
		return args, nil
	}

	for i, c := range arguments {
		switch c {
		case '(':
			level++
		case ')':
			level--
		case ',':
			if level == 0 {
				args = append(args, arguments[start:i])
				start = i + 1
			}
		}
	}
	if level != 0 {
		return nil, errors.New("opening and closing parentheses do not match")
	}
	if start >= len(arguments) {
		return nil, errors.New("can't end arguments list with a comma")
	}
	args = append(args, arguments[start:])
	return args, nil
}
