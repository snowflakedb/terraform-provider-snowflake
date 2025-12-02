package main

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

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

func (h *parameterHandler) handleIntegerParameter(level string, value string, setField **int) error {
	levelParameterType, err := sdk.ToParameterType(level)
	if err != nil {
		return err
	}

	if h.level != levelParameterType {
		return nil
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*setField = &v
	return nil
}

func (h *parameterHandler) handleBooleanParameter(level string, value string, setField **bool) error {
	levelParameterType, err := sdk.ToParameterType(level)
	if err != nil {
		return err
	}

	if h.level != levelParameterType {
		return nil
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	*setField = &b
	return nil
}

func (h *parameterHandler) handleStringParameter(level string, value string, setField **string) error {
	levelParameterType, err := sdk.ToParameterType(level)
	if err != nil {
		return err
	}

	if h.level != levelParameterType {
		return nil
	}

	*setField = &value
	return nil
}
