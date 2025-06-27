package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TestEnsureValidAccountIsUsed(t *testing.T) {
	accountLocator := "ABC123123"
	clientBuilder := sdk.NewClientBuilder().WithAccountLocator(accountLocator)

	testClient := TestClient{
		context: &TestClientContext{
			client: clientBuilder.Build(),
		},
	}

	t.Run("valid account: the test shouldn't be skipped", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "1")
		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), accountLocator)
		defer func() {
			if t.Skipped() {
				t.Errorf("Expected test to run with valid account, but it was skipped")
			}
		}()
		testClient.EnsureValidNonProdAccountIsUsed(t)
	})

	t.Run(fmt.Sprintf("invalid account is used: should skip the tests as %s is not set", testenvs.TestAccountCreate), func(t *testing.T) {
		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), accountLocator)
		t.Setenv(string(testenvs.TestAccountCreate), "")
		defer func() {
			if !t.Skipped() {
				t.Errorf("Expected test to be skipped due to missing %s environment variable", testenvs.TestAccountCreate)
			}
		}()
		testClient.EnsureValidNonProdAccountIsUsed(t)
	})

	t.Run(fmt.Sprintf("invalid account is used: should skip the tests as %s is not set", testenvs.TestNonProdModifiableAccountLocator), func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "1")
		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), "")
		defer func() {
			if !t.Skipped() {
				t.Errorf("Expected test to be skipped due to missing %s environment variable", testenvs.TestNonProdModifiableAccountLocator)
			}
		}()
		testClient.EnsureValidNonProdAccountIsUsed(t)
	})
}
