package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/stretchr/testify/require"
)

func Test_mergeSnowflakeObjectAssertionDetails_skipsEveryReflectedName(t *testing.T) {
	reflected := []genhelpers.SdkObjectDetails{
		detailsNamed("sdk.Alpha"),
		detailsNamed("sdk.Beta"),
		detailsNamed("sdk.Alpha"),
		detailsNamed("sdk.Gamma"),
		detailsNamed("sdk.Delta"),
	}
	ssotDetails := []genhelpers.SdkObjectDetails{
		{
			StructDetails: genhelpers.StructDetails{
				Name:   "sdk.Alpha",
				Fields: []genhelpers.Field{{Name: "FromSsot"}},
			},
		},
		detailsNamed("sdk.Gamma"),
	}

	merged := mergeSnowflakeObjectAssertionDetails(reflected, ssotDetails)
	require.Equal(t, []string{"sdk.Beta", "sdk.Delta", "sdk.Alpha", "sdk.Gamma"}, objectNames(merged))
	require.Equal(t, "FromSsot", merged[2].Fields[0].Name)
}

func Test_mergeSnowflakeObjectAssertionDetails_panicsOnDuplicateSsotName(t *testing.T) {
	ssotDetails := []genhelpers.SdkObjectDetails{detailsNamed("sdk.Alpha"), detailsNamed("sdk.Alpha")}
	require.Panics(t, func() {
		mergeSnowflakeObjectAssertionDetails(nil, ssotDetails)
	})
}

func detailsNamed(name string) genhelpers.SdkObjectDetails {
	return genhelpers.SdkObjectDetails{StructDetails: genhelpers.StructDetails{Name: name}}
}

func objectNames(details []genhelpers.SdkObjectDetails) []string {
	return collections.Map(details, func(d genhelpers.SdkObjectDetails) string { return d.ObjectName() })
}
