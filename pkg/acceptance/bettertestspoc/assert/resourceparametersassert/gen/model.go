package gen

import (
	"strings"

	objectparametersassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type ResourceParametersAssertionsModel struct {
	Name       string
	Parameters []ResourceParameterAssertionModel

	genhelpers.PreambleModel
}

type ResourceParameterAssertionModel struct {
	Name             string
	Type             string
	AssertionCreator string
}

func ModelFromSnowflakeObjectParameters(snowflakeObjectParameters objectparametersassertgen.SnowflakeObjectParameters, preamble genhelpers.PreambleModel) ResourceParametersAssertionsModel {
	parameters := make([]ResourceParameterAssertionModel, len(snowflakeObjectParameters.Parameters))
	for idx, p := range snowflakeObjectParameters.Parameters {
		// TODO [SNOW-1501905]: get a runtime name for the assertion creator
		var assertionCreator string
		switch {
		case p.ParameterType == "bool":
			assertionCreator = "ResourceParameterBoolValueSet"
		case p.ParameterType == "int":
			assertionCreator = "ResourceParameterIntValueSet"
		case p.ParameterType == "string":
			assertionCreator = "ResourceParameterValueSet"
		case strings.HasPrefix(p.ParameterType, "sdk."):
			assertionCreator = "ResourceParameterStringUnderlyingValueSet"
		default:
			assertionCreator = "ResourceParameterValueSet"
		}

		parameters[idx] = ResourceParameterAssertionModel{
			Name:             p.ParameterName,
			Type:             p.ParameterType,
			AssertionCreator: assertionCreator,
		}
	}

	return ResourceParametersAssertionsModel{
		Name:          snowflakeObjectParameters.ObjectName(),
		Parameters:    parameters,
		PreambleModel: preamble,
	}
}
