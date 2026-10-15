package defs

import (
	"slices"

	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

var AllParameters = ssotdefs.AllParameters

func ParameterDefsForLevel(level parameterdefs.ParameterLevel) []parameterdefs.ParameterDef {
	return ssotdefs.ParameterDefsForLevel(level)
}

func parameterNamesForLevels(levels ...parameterdefs.ParameterLevel) []string {
	names := make([]string, 0)
	for _, parameter := range AllParameters {
		if slices.ContainsFunc(levels, func(level parameterdefs.ParameterLevel) bool {
			return slices.Contains(parameter.Levels, level)
		}) && !slices.Contains(names, parameter.SqlName) {
			names = append(names, parameter.SqlName)
		}
	}
	return names
}

func withAdditionalParameterNames(names []string, parameters ...parameterdefs.ParameterDef) []string {
	result := slices.Clone(names)
	for _, parameter := range parameters {
		if !slices.Contains(result, parameter.SqlName) {
			result = append(result, parameter.SqlName)
		}
	}
	return result
}
