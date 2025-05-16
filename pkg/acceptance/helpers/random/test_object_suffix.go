package random

import (
	"fmt"
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
)

func ObjectSuffix() string {
	var suffix string
	testObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.TestObjectsSuffix))
	if testObjectSuffix != "" {
		suffix = strings.ToUpper(testObjectSuffix)
	} else {
		uuid := UUID()
		formattedUuidSuffix := strings.ToUpper(strings.ReplaceAll(uuid, "-", "_"))
		suffix = formattedUuidSuffix
	}
	return suffix
}
