package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

// objectParametersFromSsot builds the parameter assertion input for one opted-in object.
// Parameters come from the catalog entries for the object's parameter level.
// Constants are sdk.<Name>Parameter.
func objectParametersFromSsot(obj ssot.Object, enabled bool) SnowflakeObjectParameters {
	if !enabled {
		panic(fmt.Sprintf("object %q is not enabled for snowflake parameters assertion generation", obj.Name))
	}
	parameters := defs.ParameterDefsForLevel(obj.ParameterLevel())
	if len(parameters) == 0 {
		panic(fmt.Sprintf("object %q is enabled for snowflake parameters assertions but has no catalog parameters at level %s", obj.Name, obj.ParameterLevel()))
	}
	names := obj.ParameterAssertionNames()
	return SnowflakeObjectParameters{
		Name:                    obj.Name,
		IdType:                  "sdk." + obj.IdType,
		ParameterConstantPrefix: names.ConstantPrefix,
		ObjectTypeName:          names.ObjectTypeName,
		Parameters:              snowflakeParameters(parameters),
	}
}
