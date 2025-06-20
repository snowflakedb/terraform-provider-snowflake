package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
)

func TestEnsureValidAccountIsUsed(t *testing.T) {
	t.Run("valid account", func(t *testing.T) {
		accountLocator := testClient().Context.CurrentAccount(t)
		t.Setenv(string(testenvs.TestAccountCreate), "1")
		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), accountLocator)
		defer func() {
			if t.Skipped() {
				t.Errorf("Expected test to run with valid account, but it was skipped")
			}
		}()
		testClient().EnsureValidNonProdAccountIsUsed(t)
	})

	t.Run(fmt.Sprintf("invalid account: %s not set", testenvs.TestAccountCreate), func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "")
		defer func() {
			if !t.Skipped() {
				t.Errorf("Expected test to be skipped due to missing %s environment variable", testenvs.TestAccountCreate)
			}
		}()
		testClient().EnsureValidNonProdAccountIsUsed(t)
	})

	t.Run(fmt.Sprintf("invalid account: %s not set", testenvs.TestNonProdModifiableAccountLocator), func(t *testing.T) {
		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), "")
		defer func() {
			if !t.Skipped() {
				t.Errorf("Expected test to be skipped due to missing %s environment variable", testenvs.TestNonProdModifiableAccountLocator)
			}
		}()
		testClient().EnsureValidNonProdAccountIsUsed(t)
	})
}
