package random

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
)

var postgresMirrorNameForbiddenCharacters = regexp.MustCompile(`[^A-Za-z0-9_]`)

func PostgresMirrorNameSuffix(suffix string) string {
	return strings.ToLower(postgresMirrorNameForbiddenCharacters.ReplaceAllString(suffix, "_"))
}

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
