package main

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func handleOptionalFieldWithBuilder[T any, U any](parameter *T, builder func(T) *U) {
	if parameter != nil {
		builder(*parameter)
	}
}

func handleOptionalFieldWithStringBuilder[T any, U any](parameter *T, builder func(string) *U) {
	if parameter != nil {
		builder(fmt.Sprintf("%v", *parameter))
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

func (h *parameterHandler) handleIntegerParameter(level sdk.ParameterType, value string, setField **int) error {
	if h.level != level {
		return nil
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*setField = &v
	return nil
}

func (h *parameterHandler) handleBooleanParameter(level sdk.ParameterType, value string, setField **bool) error {
	if h.level != level {
		return nil
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	*setField = &b
	return nil
}

func (h *parameterHandler) handleStringParameter(level sdk.ParameterType, value string, setField **string) error {
	if h.level != level {
		return nil
	}
	*setField = &value
	return nil
}
