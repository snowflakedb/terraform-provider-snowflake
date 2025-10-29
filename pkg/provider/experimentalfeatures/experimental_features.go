package experimentalfeatures

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO [SNOW-2398035]: allow discontinuing the experiment
// TODO [SNOW-2398035]: generate docs with the experiment description
// TODO [SNOW-2398035]: automatically fix the description for the experimental feature when experiment is discontinued

type ExperimentalFeature string

const (
	WarehouseShowImprovedPerformance ExperimentalFeature = "WAREHOUSE_SHOW_IMPROVED_PERFORMANCE"
)

var allExperimentalFeatures = []ExperimentalFeature{
	WarehouseShowImprovedPerformance,
}

var AllExperimentalFeatures = sdk.AsStringList(allExperimentalFeatures)

func IsExperimentEnabled(experiment ExperimentalFeature, enabledExperiments []string) bool {
	return slices.ContainsFunc(enabledExperiments, func(s string) bool {
		return strings.EqualFold(string(experiment), s)
	})
}
