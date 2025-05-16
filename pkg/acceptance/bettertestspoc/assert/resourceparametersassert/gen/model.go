package gen

import (
	"os"
	"strings"

	objectparametersassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type ResourceParametersAssertionsModel struct {
	Name       string
	Parameters []ResourceParameterAssertionModel
	PreambleModel
}

func (m ResourceParametersAssertionsModel) SomeFunc() {
}

type ResourceParameterAssertionModel struct {
	Name             string
	Type             string
	AssertionCreator string
}

func ModelFromSnowflakeObjectParameters(snowflakeObjectParameters objectparametersassertgen.SnowflakeObjectParameters) ResourceParametersAssertionsModel {
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

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceParametersAssertionsModel{
		Name:       snowflakeObjectParameters.ObjectName(),
		Parameters: parameters,
		PreambleModel: PreambleModel{
			PackageName: packageWithGenerateDirective,
		},
	}
}
