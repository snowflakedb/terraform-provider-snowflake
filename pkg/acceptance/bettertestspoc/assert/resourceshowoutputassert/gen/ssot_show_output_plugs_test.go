package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/stretchr/testify/require"
)

func Test_mergeShowOutputAssertionDetails_skipsReflectedAndKeepsRepeats(t *testing.T) {
	reflected := []SdkObjectShowOutputDetails{
		detailsNamed("sdk.Alpha"),
		detailsNamed("sdk.Beta"),
		detailsNamed("sdk.Alpha"),
		detailsNamed("sdk.Gamma"),
	}
	ssotDetails := []SdkObjectShowOutputDetails{{
		SdkObjectDetails: genhelpers.SdkObjectDetails{StructDetails: genhelpers.StructDetails{
			Name:   "sdk.Beta",
			Fields: []genhelpers.Field{{Name: "FromSsot"}},
		}},
	}}

	merged := mergeShowOutputAssertionDetails(reflected, ssotDetails)
	require.Equal(t, []string{"sdk.Alpha", "sdk.Alpha", "sdk.Gamma", "sdk.Beta"}, objectNames(merged))
	require.Equal(t, "FromSsot", merged[3].Fields[0].Name)
}

func Test_mergeShowOutputAssertionDetails_panicsOnDuplicateSsotName(t *testing.T) {
	ssotDetails := []SdkObjectShowOutputDetails{detailsNamed("sdk.Beta"), detailsNamed("sdk.Beta")}
	require.Panics(t, func() {
		mergeShowOutputAssertionDetails(nil, ssotDetails)
	})
}

func detailsNamed(name string) SdkObjectShowOutputDetails {
	return SdkObjectShowOutputDetails{
		SdkObjectDetails: genhelpers.SdkObjectDetails{StructDetails: genhelpers.StructDetails{Name: name}},
	}
}

func objectNames(details []SdkObjectShowOutputDetails) []string {
	return collections.Map(details, func(d SdkObjectShowOutputDetails) string { return d.ObjectName() })
}
