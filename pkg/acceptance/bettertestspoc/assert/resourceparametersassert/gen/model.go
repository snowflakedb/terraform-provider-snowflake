package gen

import (
	"strings"

	objectparametersassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type ResourceParametersAssertionsModel struct {
	Name                  string
	DataSourceName        string
	Parameters            []ResourceParameterAssertionModel
	ParameterConstantName string

	*genhelpers.PreambleModel
}

type ResourceParameterAssertionModel struct {
	Name             string
	Type             string
	AssertionCreator string
	Mapper           string
}

var dataSourceParametersMapping = map[string]string{
	"Database":  "Databases",
	"Task":      "Tasks",
	"User":      "Users",
	"Warehouse": "Warehouses",
}

func ModelFromSnowflakeObjectParameters(snowflakeObjectParameters objectparametersassertgen.SnowflakeObjectParameters, preamble *genhelpers.PreambleModel) ResourceParametersAssertionsModel {
	parameters := make([]ResourceParameterAssertionModel, len(snowflakeObjectParameters.Parameters))
	for idx, p := range snowflakeObjectParameters.Parameters {
		// TODO [SNOW-1501905]: get a runtime name for the assertion creator
		var assertionCreator string
		var mapper string
		switch {
		case p.ParameterType == "bool":
			assertionCreator = "ParameterBoolValueSet"
		case p.ParameterType == "int":
			assertionCreator = "ParameterIntValueSet"
		case p.ParameterType == "string":
			assertionCreator = "ParameterValueSet"
		case strings.HasPrefix(p.ParameterType, "sdk."):
			assertionCreator = "ParameterValueSet"
			mapper = "string"
		default:
			assertionCreator = "ParameterValueSet"
		}

		parameters[idx] = ResourceParameterAssertionModel{
			Name:             p.ParameterName,
			Type:             p.ParameterType,
			AssertionCreator: assertionCreator,
			Mapper:           mapper,
		}
	}

	name := snowflakeObjectParameters.ObjectName()
	dataSourceName := dataSourceParametersMapping[name]

	parameterConstantName := name
	if snowflakeObjectParameters.ParameterConstantPrefix != "" {
		parameterConstantName = snowflakeObjectParameters.ParameterConstantPrefix
	}

	return ResourceParametersAssertionsModel{
		Name:                  name,
		DataSourceName:        dataSourceName,
		Parameters:            parameters,
		ParameterConstantName: parameterConstantName,
		PreambleModel:         preamble,
	}
}
