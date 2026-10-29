package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

func ShowOutputAssertionDetails() []SdkObjectShowOutputDetails {
	return mergeShowOutputAssertionDetails(GetFilteredSdkObjectDetails(), ssotShowOutputDetails())
}

func ssotShowOutputDetails() []SdkObjectShowOutputDetails {
	var details []SdkObjectShowOutputDetails
	for _, obj := range defs.All() {
		g := ssot.ObjectGenerators(obj.Generators...)
		if g.ShowOutputAssert.Show {
			details = append(details, showOutputDetailsFromSsot(obj, false))
		}
		if g.ShowOutputAssert.Describe {
			details = append(details, showOutputDetailsFromSsot(obj, true))
		}
	}
	return details
}

// mergeShowOutputAssertionDetails drops reflected names that ssot emits, then appends the ssot models.
// allStructs already repeats some sdk types that are not opted in. Those repeats stay.
// A repeated ssot ObjectName, or an ssot name still present after the skip, panics.
func mergeShowOutputAssertionDetails(reflected, ssotDetails []SdkObjectShowOutputDetails) []SdkObjectShowOutputDetails {
	plugged := make(map[string]struct{}, len(ssotDetails))
	for _, details := range ssotDetails {
		name := details.ObjectName()
		if _, dup := plugged[name]; dup {
			panic(fmt.Sprintf("duplicate show output assertion object name %q", name))
		}
		plugged[name] = struct{}{}
	}
	merged := make([]SdkObjectShowOutputDetails, 0, len(reflected)+len(ssotDetails))
	for _, details := range reflected {
		if _, ok := plugged[details.ObjectName()]; ok {
			continue
		}
		merged = append(merged, details)
	}
	for _, details := range merged {
		if _, ok := plugged[details.ObjectName()]; ok {
			panic(fmt.Sprintf("show output assertion object name %q is still reflected after the ssot skip", details.ObjectName()))
		}
	}
	return append(merged, ssotDetails...)
}
