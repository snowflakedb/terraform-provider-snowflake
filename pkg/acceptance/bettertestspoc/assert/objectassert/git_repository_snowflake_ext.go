package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *GitRepositoryAssert) HasCreatedOnNotEmpty() *GitRepositoryAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.GitRepository) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *GitRepositoryAssert) HasGitCredentialsEmpty() *GitRepositoryAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.GitRepository) error {
		t.Helper()
		if o.GitCredentials != nil {
			return fmt.Errorf("expected git_credentials to be empty")
		}
		return nil
	})
	return a
}
