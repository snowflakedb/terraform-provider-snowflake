package acceptancetests

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
)

// TODO [SNOW-1356199]: add generation tests
// TODO [SNOW-1356199]: use the same fallback suffix for acceptance and integration tests (now two different ones are generated if the env is missing)
var ObjectsSuffix = acceptanceTestsSuffix()

func acceptanceTestsSuffix() string {
	suffix := "AT_" + random.ObjectSuffix()
	log.Printf("[DEBUG] Suffix for the given test run is: %s", suffix)
	return suffix
}
