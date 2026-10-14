package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

func SnowflakeObjectAssertionDetails() []genhelpers.SdkObjectDetails {
	return mergeSnowflakeObjectAssertionDetails(GetSdkObjectDetails(), ssotObjectAssertionDetails())
}

func ssotObjectAssertionDetails() []genhelpers.SdkObjectDetails {
	var details []genhelpers.SdkObjectDetails
	for _, obj := range defs.All() {
		g := ssot.ObjectGenerators(obj.Generators...)
		if g.ObjectAssert.Show {
			details = append(details, objectAssertionDetailsFromSsot(obj, false))
		}
		if g.ObjectAssert.Describe {
			details = append(details, objectAssertionDetailsFromSsot(obj, true))
		}
	}
	return details
}

// mergeSnowflakeObjectAssertionDetails drops reflected names that ssot emits, then appends the ssot models.
// A repeated ssot ObjectName, or an ssot name still present after the skip, panics.
func mergeSnowflakeObjectAssertionDetails(reflected, ssotDetails []genhelpers.SdkObjectDetails) []genhelpers.SdkObjectDetails {
	plugged := make(map[string]struct{}, len(ssotDetails))
	for _, details := range ssotDetails {
		name := details.ObjectName()
		if _, dup := plugged[name]; dup {
			panic(fmt.Sprintf("duplicate snowflake object assertion object name %q", name))
		}
		plugged[name] = struct{}{}
	}
	merged := make([]genhelpers.SdkObjectDetails, 0, len(reflected)+len(ssotDetails))
	for _, details := range reflected {
		if _, ok := plugged[details.ObjectName()]; ok {
			continue
		}
		merged = append(merged, details)
	}
	for _, details := range merged {
		if _, ok := plugged[details.ObjectName()]; ok {
			panic(fmt.Sprintf("snowflake object assertion object name %q is still reflected after the ssot skip", details.ObjectName()))
		}
	}
	return append(merged, ssotDetails...)
}
