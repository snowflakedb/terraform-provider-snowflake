package gen

import (
	"reflect"
	"slices"

	objectassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SdkObjectShowOutputDetails struct {
	dataSourceDef *dataSourceDef
	genhelpers.SdkObjectDetails
}

var dataSourceMappingNormalized = map[string]dataSourceDef{
	normalized(sdk.Database{}):           {"Databases"},
	normalized(sdk.NetworkRuleDetails{}): {"NetworkRules"},
}

type dataSourceDef struct {
	pluralName string
}

// GetFilteredSdkObjectDetails is currently needed to filter out objects that are not resources because the same underlying list of objects is used.
func GetFilteredSdkObjectDetails() []SdkObjectShowOutputDetails {
	allDetails := objectassertgen.GetSdkObjectDetails()
	filtered := collections.Filter(allDetails, func(d genhelpers.SdkObjectDetails) bool {
		return !slices.Contains(objectNamesNotBeingResources, d.Name)
	})
	return collections.Map(filtered, func(d genhelpers.SdkObjectDetails) SdkObjectShowOutputDetails {
		v, _ := dataSourceMappingNormalized[d.Name]
		return SdkObjectShowOutputDetails{&v, d}
	})
}

var (
	objectsNotBeingResources     = []any{sdk.UserWorkloadIdentityAuthenticationMethod{}}
	objectNamesNotBeingResources = collections.Map(objectsNotBeingResources, func(o any) string {
		return reflect.ValueOf(o).Type().String()
	})
)

func normalized(t any) string {
	return reflect.ValueOf(t).Type().String()
}
