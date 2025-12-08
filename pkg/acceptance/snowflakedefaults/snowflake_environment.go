package snowflakedefaults

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
)

type SnowflakeEnvironment string

const (
	SnowflakeProdEnvironment    SnowflakeEnvironment = "PROD"
	SnowflakeNonProdEnvironment SnowflakeEnvironment = "NON_PROD"
)

var allSnowflakeEnvironments = []SnowflakeEnvironment{
	SnowflakeProdEnvironment,
	SnowflakeNonProdEnvironment,
}

func parseSnowflakeEnvironment(environment string) (SnowflakeEnvironment, error) {
	snowflakeEnvironment := SnowflakeEnvironment(strings.ToUpper(environment))
	if slices.Contains(allSnowflakeEnvironments, snowflakeEnvironment) {
		return snowflakeEnvironment, nil
	}
	return "", fmt.Errorf("invalid Snowflake environment: %s, valid values are: %v", environment, allSnowflakeEnvironments)
}

func getSnowflakeEnvironmentWithProdDefault() SnowflakeEnvironment {
	env := os.Getenv(string(testenvs.SnowflakeTestingEnvironment))
	if env == "" {
		log.Printf("[DEBUG] Snowflake environment variable %s not set, returning default PROD environment", testenvs.SnowflakeTestingEnvironment)
		return SnowflakeProdEnvironment
	}
	snowflakeEnvironment, err := parseSnowflakeEnvironment(env)
	if err != nil {
		log.Printf("[DEBUG] Failed to parse Snowflake environment variable (%s), returning default PROD environment, err = %s", testenvs.SnowflakeTestingEnvironment, err)
		return SnowflakeProdEnvironment
	}
	return snowflakeEnvironment
}
