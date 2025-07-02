package helpers

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TestEnsureValidAccountIsUsed(t *testing.T) {
	client := sdk.Client{}

	testClient := TestClient{
		context: &TestClientContext{
			client: &client,
		},
	}

	t.Run("valid account: the test shouldn't be skipped", func(t *testing.T) {
		accountLocator := "ABC123123"
		client.SetAccountLocatorForTests(accountLocator)

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
		accountLocator := "ABC123123"
		client.SetAccountLocatorForTests(accountLocator)
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

type mockOrganizationAccounts struct {
	ShowResult    []sdk.OrganizationAccount
	ShowResultErr error
}

func (m *mockOrganizationAccounts) Show(ctx context.Context, request *sdk.ShowOrganizationAccountRequest) ([]sdk.OrganizationAccount, error) {
	return m.ShowResult, m.ShowResultErr
}

func (m *mockOrganizationAccounts) Create(ctx context.Context, request *sdk.CreateOrganizationAccountRequest) error {
	return nil
}

func (m *mockOrganizationAccounts) Alter(ctx context.Context, request *sdk.AlterOrganizationAccountRequest) error {
	return nil
}

func (m *mockOrganizationAccounts) ShowParameters(ctx context.Context) ([]*sdk.Parameter, error) {
	return nil, nil
}

func (m *mockOrganizationAccounts) UnsetAllParameters(ctx context.Context) error {
	return nil
}

func TestEnsureValidOrganizationAccountIsUsed(t *testing.T) {
	accountLocator := "ABC123123"
	anotherAccountLocator := "DEF456456"

	organizationAccounts := &mockOrganizationAccounts{
		ShowResult: []sdk.OrganizationAccount{
			{
				AccountLocator: accountLocator,
			},
		},
	}

	client := sdk.Client{
		OrganizationAccounts: organizationAccounts,
	}
	client.SetAccountLocatorForTests(accountLocator)

	testClient := TestClient{
		context: &TestClientContext{
			client: &client,
		},
	}

	t.Run("valid organization account: the test shouldn't be skipped", func(t *testing.T) {
		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), accountLocator)
		defer func() {
			if t.Skipped() {
				t.Errorf("Expected test to run with valid account, but it was skipped")
			}
		}()
		testClient.EnsureValidNonProdOrganizationAccountIsUsed(t)
	})

	t.Run(fmt.Sprintf("invalid account is used: should skip the tests as %s is not set", testenvs.TestNonProdModifiableAccountLocator), func(t *testing.T) {
		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), "")
		defer func() {
			if !t.Skipped() {
				t.Errorf("Expected test to be skipped due to missing %s environment variable", testenvs.TestNonProdModifiableAccountLocator)
			}
		}()
		testClient.EnsureValidNonProdOrganizationAccountIsUsed(t)
	})

	t.Run("invalid organization account: SHOW ORGANIZATION ACCOUNTS error", func(t *testing.T) {
		organizationAccounts.ShowResultErr = errors.New("SHOW ORGANIZATION ACCOUNTS error")
		t.Cleanup(func() {
			organizationAccounts.ShowResultErr = nil
		})

		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), accountLocator)
		defer func() {
			if !t.Skipped() {
				t.Errorf("Expected test to be skipped due to SHOW ORGANIZATION ACCOUNTS error")
			}
		}()

		testClient.EnsureValidNonProdOrganizationAccountIsUsed(t)
	})

	t.Run("invalid organization account: account locator not matching", func(t *testing.T) {
		t.Setenv(string(testenvs.TestNonProdModifiableAccountLocator), anotherAccountLocator)
		defer func() {
			if !t.Skipped() {
				t.Errorf("Expected test to be skipped due to SHOW ORGANIZATION ACCOUNTS error")
			}
		}()
		testClient.EnsureValidNonProdOrganizationAccountIsUsed(t)
	})
}
