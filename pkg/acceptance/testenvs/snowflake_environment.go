package testenvs

import (
	"fmt"
	"os"
	"slices"
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
	snowflakeEnvironment := SnowflakeEnvironment(environment)
	if slices.Contains(allSnowflakeEnvironments, snowflakeEnvironment) {
		return snowflakeEnvironment, nil
	}
	return "", fmt.Errorf("invalid Snowflake environment: %s", environment)
}

func GetSnowflakeEnvironmentWithProdDefault() SnowflakeEnvironment {
	env := os.Getenv(string(SnowflakeTestingEnvironment))
	if env == "" {
		return SnowflakeProdEnvironment
	}
	snowflakeEnvironment, err := parseSnowflakeEnvironment(env)
	if err != nil {
		return SnowflakeProdEnvironment
	}
	return snowflakeEnvironment
}
