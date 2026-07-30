package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func UserDefaultWorkloadIdentityAuthenticationMethods(t *testing.T, userId sdk.AccountObjectIdentifier) *UserWorkloadIdentityAuthenticationMethodAssert {
	t.Helper()
	return UserWorkloadIdentityAuthenticationMethod(t, userId, sdk.NewAccountObjectIdentifier("DEFAULT"))
}

func (u *UserWorkloadIdentityAuthenticationMethodAssert) HasNoComment() *UserWorkloadIdentityAuthenticationMethodAssert {
	u.AddAssertion(func(t *testing.T, o *sdk.UserWorkloadIdentityAuthenticationMethod) error {
		t.Helper()
		if o.Comment != "" {
			return fmt.Errorf("expected comment to be empty; got: %s", o.Comment)
		}
		return nil
	})
	return u
}
