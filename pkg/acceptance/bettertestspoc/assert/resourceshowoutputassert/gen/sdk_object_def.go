package gen

import (
	"reflect"
	"slices"

	objectassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// GetFilteredSdkObjectDetails is currently needed to filter out objects that are not resources because the same underlying list of objects is used.
func GetFilteredSdkObjectDetails() []genhelpers.SdkObjectDetails {
	allDetails := objectassertgen.GetSdkObjectDetails()
	return collections.Filter(allDetails, func(d genhelpers.SdkObjectDetails) bool {
		return !slices.Contains(objectNamesNotBeingResources, d.Name)
	})
}

var (
	objectsNotBeingResources     = []any{sdk.UserWorkloadIdentityAuthenticationMethod{}}
	objectNamesNotBeingResources = collections.Map(objectsNotBeingResources, func(o any) string {
		return reflect.ValueOf(o).Type().String()
	})
)
