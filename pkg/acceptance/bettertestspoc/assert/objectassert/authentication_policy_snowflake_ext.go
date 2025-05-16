package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (a *AuthenticationPolicyAssert) HasCreatedOnNotEmpty() *AuthenticationPolicyAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.AuthenticationPolicy) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}
