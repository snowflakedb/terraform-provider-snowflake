package gen

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type SnowflakeObjectParametersAssertionsModel struct {
	Name       string
	IdType     string
	Parameters []ParameterAssertionModel

	genhelpers.PreambleModel
}

type ParameterAssertionModel struct {
	Name             string
	Type             string
	DefaultValue     string
	DefaultLevel     string
	AssertionCreator string
}

func ModelFromSnowflakeObjectParameters(snowflakeObjectParameters SnowflakeObjectParameters, preamble genhelpers.PreambleModel) SnowflakeObjectParametersAssertionsModel {
	parameters := make([]ParameterAssertionModel, len(snowflakeObjectParameters.Parameters))
	for idx, p := range snowflakeObjectParameters.Parameters {
		// TODO [SNOW-1501905]: get a runtime name for the assertion creator
		var assertionCreator string
		switch {
		case p.ParameterType == "bool":
			assertionCreator = "SnowflakeParameterBoolValueSet"
		case p.ParameterType == "int":
			assertionCreator = "SnowflakeParameterIntValueSet"
		case p.ParameterType == "string":
			assertionCreator = "SnowflakeParameterValueSet"
		case strings.HasPrefix(p.ParameterType, "sdk."):
			assertionCreator = "SnowflakeParameterStringUnderlyingValueSet"
		// TODO [SNOW-1501905]: handle other types if needed
		default:
			assertionCreator = "SnowflakeParameterValueSet"
		}

		defaultValue := p.DefaultValue
		// string has to be wrapped in double quotes; all other values are passed explicitly
		if p.ParameterType == "string" {
			defaultValue = fmt.Sprintf(`"%s"`, defaultValue)
		}

		parameters[idx] = ParameterAssertionModel{
			Name:             p.ParameterName,
			Type:             p.ParameterType,
			DefaultValue:     defaultValue,
			DefaultLevel:     p.DefaultLevel,
			AssertionCreator: assertionCreator,
		}
	}

	return SnowflakeObjectParametersAssertionsModel{
		Name:          snowflakeObjectParameters.ObjectName(),
		IdType:        snowflakeObjectParameters.IdType,
		Parameters:    parameters,
		PreambleModel: preamble,
	}
}
