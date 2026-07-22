package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (g *GitRepositoryAssert) HasGitCredentialsEmpty() *GitRepositoryAssert {
	g.AddAssertion(func(t *testing.T, o *sdk.GitRepository) error {
		t.Helper()
		if o.GitCredentials != nil {
			return fmt.Errorf("expected git_credentials to be empty")
		}
		return nil
	})
	return g
}
