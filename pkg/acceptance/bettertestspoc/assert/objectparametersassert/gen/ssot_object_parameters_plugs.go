package gen

import (
	"fmt"

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
	var out []SnowflakeObjectParameters
	for _, obj := range defs.All() {
		if enabled(ssot.ObjectGenerators(obj.Generators...)) {
			out = append(out, objectParametersFromSsot(obj, true))
		}
	}
	return out
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
