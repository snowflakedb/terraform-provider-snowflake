package main

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func handleOptionalFieldWithBuilder[T any](parameter *T, builder func(T) *model.SchemaModel) {
	if parameter != nil {
		builder(*parameter)
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
