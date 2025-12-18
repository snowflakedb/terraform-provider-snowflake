package main

import (
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// csvUnescape reverses the CSV escaping done by the Terraform HCL generator.
// It converts:
//   - \n -> newline
//   - \r -> carriage return
//   - \\ -> backslash
func csvUnescape(s string) string {
	const placeholder = "\x00BACKSLASH\x00"

	result := s
	result = strings.ReplaceAll(result, "\\\\", placeholder)
	result = strings.ReplaceAll(result, "\\n", "\n")
	result = strings.ReplaceAll(result, "\\r", "\r")
	result = strings.ReplaceAll(result, placeholder, "\\")
	return result
}

func handleOptionalFieldWithBuilder[T any, U any](parameter *T, builder func(T) *U) {
	if parameter != nil {
		builder(*parameter)
	}
}

func handleIfNotEmpty[T any](value string, builder func(string) *T) {
	if value != "" {
		builder(value)
	}
}

func handleIf[T any](condition bool, builder func(string) *T) {
	if condition {
		builder("true")
	}
}

type parameterHandler struct {
	level sdk.ParameterType
}

func newParameterHandler(level sdk.ParameterType) parameterHandler {
	return parameterHandler{
		level: level,
	}
}

func handleParameter[T any](h *parameterHandler, level string, value string, setField **T, parser func(string) (T, error)) error {
	levelParameterType, err := sdk.ToParameterType(level)
	if err != nil {
		return err
	}

	if h.level != levelParameterType {
		return nil
	}

	v, err := parser(value)
	if err != nil {
		return err
	}
	*setField = &v
	return nil
}

func (h *parameterHandler) handleIntegerParameter(level string, value string, setField **int) error {
	return handleParameter(h, level, value, setField, strconv.Atoi)
}

func (h *parameterHandler) handleBooleanParameter(level string, value string, setField **bool) error {
	return handleParameter(h, level, value, setField, strconv.ParseBool)
}

func (h *parameterHandler) handleStringParameter(level string, value string, setField **string) error {
	return handleParameter(h, level, value, setField, func(value string) (string, error) { return value, nil })
}
