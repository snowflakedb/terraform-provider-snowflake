package testacc

import (
	"fmt"
	"log"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
)

// TODO [next PRs]: contents of this file may be potentially reused for both integration and acceptance tests setups

// timer measures time from invocation point to the end of method.
// It's supposed to be used like:
//
//	defer timer("something to measure name")()
func timer(name string, logger *log.Logger) func() {
	logger.Printf("[INFO] Timer start: %s starting now", name)
	start := time.Now()
	return func() {
		logger.Printf("[INFO] Timer stop: %s took %v", name, time.Since(start))
	}
}

func setUpSdkClient(profile string, tests string) (*gosnowflake.Config, *sdk.Client, error) {
	conf, err := sdk.ProfileConfig(profile)
	if err != nil {
		return nil, nil, err
	}
	if conf == nil {
		return nil, nil, fmt.Errorf("%s config is required to run %s tests", profile, tests)
	}

	c, err := sdk.NewClient(conf)
	if err != nil {
		return nil, nil, err
	}
	return conf, c, nil
}
