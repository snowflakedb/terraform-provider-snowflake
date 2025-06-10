package helpers

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"testing"
)

func TestEnsureValidAccountIsUsed(t *testing.T) {
	t.Run("valid account", func(t *testing.T) {})

	t.Run(fmt.Sprintf("invalid account: %s not set", testenvs.TestAccountCreate), func(t *testing.T) {})

	t.Run(fmt.Sprintf("invalid account: %s not set", testenvs.TestNonProdModifiableAccountLocator), func(t *testing.T) {})
}
