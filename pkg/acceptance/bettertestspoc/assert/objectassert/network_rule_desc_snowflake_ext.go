package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (n *NetworkRuleDetailsAssert) HasCreatedOnNotEmpty() *NetworkRuleDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NetworkRuleDetails) error {
		t.Helper()
		if o.CreatedOn.IsZero() {
			return fmt.Errorf("expected created on to not be empty")
		}
		return nil
	})
	return n
}
