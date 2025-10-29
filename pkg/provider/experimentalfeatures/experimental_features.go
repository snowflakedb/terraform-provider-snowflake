package experimentalfeatures

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type feature string

const (
	WarehouseShowImprovedPerformance feature = "WAREHOUSE_SHOW_IMPROVED_PERFORMANCE"
)

var allExperimentalFeatures = []feature{
	WarehouseShowImprovedPerformance,
}

var AllExperimentalFeatures = sdk.AsStringList(allExperimentalFeatures)

// TODO [this PR]: test
// TODO [this PR]: Describe logic for disabliong experiments on the provider side and adjusting the implementation
// TODO [next PR]: Add documentation for such experimental features in the provider's docs automatically
func IsExperimentEnabled(experiment feature, enabledExperiments []string) bool {
	return !slices.ContainsFunc(enabledExperiments, func(s string) bool {
		return strings.EqualFold(string(experiment), s)
	})
}
