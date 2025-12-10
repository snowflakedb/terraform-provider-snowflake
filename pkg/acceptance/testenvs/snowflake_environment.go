package testenvs

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
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

func GetSnowflakeEnvironmentWithProdDefault() SnowflakeEnvironment {
	env := os.Getenv(string(SnowflakeTestingEnvironment))
	if env == "" {
		log.Printf("[DEBUG] Snowflake environment variable %s not set, returning default PROD environment", SnowflakeTestingEnvironment)
		return SnowflakeProdEnvironment
	}
	snowflakeEnvironment, err := parseSnowflakeEnvironment(env)
	if err != nil {
		log.Printf("[DEBUG] Failed to parse Snowflake environment variable (%s), returning default PROD environment, err = %s", SnowflakeTestingEnvironment, err)
		return SnowflakeProdEnvironment
	}
	return snowflakeEnvironment
}
