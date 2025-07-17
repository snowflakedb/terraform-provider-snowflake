package testfunctional_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
)

var functionalTestLog = log.New(os.Stdout, "", log.LstdFlags)

func TestMain(m *testing.M) {
	exitVal := execute(m)
	os.Exit(exitVal)
}

func execute(m *testing.M) int {
	defer cleanup()
	functionalTestLog.Printf("[INFO] Setting up functional tests")
	err := initialize()
	if err != nil {
		functionalTestLog.Printf("[ERROR] Functional tests setup failed with: `%s`", err)
		return 1
	}
	exitVal := m.Run()
	return exitVal
}

func initialize() error {
	enableAcceptance := os.Getenv(fmt.Sprintf("%v", testenvs.EnableAcceptance))
	if enableAcceptance == "" {
		return fmt.Errorf("acceptance tests cannot be run; set %s env to run them", testenvs.EnableAcceptance)
	}

	generatedRandomValue := os.Getenv(fmt.Sprintf("%v", testenvs.GeneratedRandomValue))
	if generatedRandomValue == "" {
		return fmt.Errorf("generated random value is required for this test run; set %s env", testenvs.GeneratedRandomValue)
	}

	if err := setUpProvidersForFunctionalTests(); err != nil {
		return fmt.Errorf("cannot set up test providers for the functional tests, err: %w", err)
	}

	if err := setUpTestHttpServer(); err != nil {
		return fmt.Errorf("cannot set up test http server for the functional tests, err: %w", err)
	}

	return nil
}

func cleanup() {
	functionalTestLog.Printf("[INFO] Running functional tests cleanup")

	if serverCleanup != nil {
		defer serverCleanup()
	}
}
