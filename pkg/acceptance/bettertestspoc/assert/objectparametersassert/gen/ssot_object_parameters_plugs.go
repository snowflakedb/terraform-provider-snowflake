package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

// SnowflakeObjectParametersForAssertions returns the hand-written definitions merged with objects
// opted in through WithObjectParametersAssert.
func SnowflakeObjectParametersForAssertions() []SnowflakeObjectParameters {
	return parametersForAssertions(func(g ssot.Generators) bool {
		return g.ObjectParametersAssert.Enabled
	})
}

// SnowflakeResourceParametersForAssertions returns the hand-written definitions merged with objects
// opted in through WithResourceParametersAssert.
func SnowflakeResourceParametersForAssertions() []SnowflakeObjectParameters {
	return parametersForAssertions(func(g ssot.Generators) bool {
		return g.ResourceParametersAssert.Enabled
	})
}

func parametersForAssertions(enabled func(ssot.Generators) bool) []SnowflakeObjectParameters {
	return mergeSnowflakeObjectParameters(GetAllSnowflakeObjectParameters(), ssotObjectParameters(enabled))
}

func ssotObjectParameters(enabled func(ssot.Generators) bool) []SnowflakeObjectParameters {
	return collections.Map(defs.SelectObjects(enabled, selectsForParameters), func(obj ssot.Object) SnowflakeObjectParameters {
		return objectParametersFromSsot(obj, true)
	})
}

// selectsForParameters is the one filter of both parameter assertion generators: the catalog has entries at the object's
// parameter level. An enabled object that fails it is skipped, so the generators can be enabled for every object.
func selectsForParameters(obj ssot.Object) bool {
	return len(defs.ParameterDefsForLevel(obj.ParameterLevel())) > 0
}

// mergeSnowflakeObjectParameters drops hand-written names that ssot emits, then appends the ssot definitions.
// A repeated ssot name panics.
func mergeSnowflakeObjectParameters(handWritten, ssotParameters []SnowflakeObjectParameters) []SnowflakeObjectParameters {
	plugged := make(map[string]struct{}, len(ssotParameters))
	for _, parameters := range ssotParameters {
		name := parameters.ObjectName()
		if _, dup := plugged[name]; dup {
			panic(fmt.Sprintf("duplicate snowflake object parameters assertion object name %q", name))
		}
		plugged[name] = struct{}{}
	}
	merged := make([]SnowflakeObjectParameters, 0, len(handWritten)+len(ssotParameters))
	for _, parameters := range handWritten {
		if _, ok := plugged[parameters.ObjectName()]; ok {
			continue
		}
		merged = append(merged, parameters)
	}
	return append(merged, ssotParameters...)
}
