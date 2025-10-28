package experimentalfeatures

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type feature string

const (
	ImprovedWarehouseShowQuery feature = "ImprovedWarehouseShowQuery"
)

var allExperimentalFeatures = []feature{
	ImprovedWarehouseShowQuery,
}

var AllExperimentalFeatures = sdk.AsStringList(allExperimentalFeatures)
