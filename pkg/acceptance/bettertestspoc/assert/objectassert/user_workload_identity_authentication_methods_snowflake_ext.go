package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func UserWorkloadIdentityAuthenticationMethods(t *testing.T, userId sdk.AccountObjectIdentifier, name string) *UserWorkloadIdentityAuthenticationMethodsAssert {
	t.Helper()
	wifId := helpers.NewUserWorkloadIdentityAuthenticationMethodsObjectIdentifier(userId, name)
	return &UserWorkloadIdentityAuthenticationMethodsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeUserWorkloadIdentityAuthenticationMethod, wifId, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.UserWorkloadIdentityAuthenticationMethod, helpers.UserWorkloadIdentityAuthenticationMethodsObjectIdentifier] {
			return testClient.User.ShowUserWorkloadIdentityAuthenticationMethodOptions
		}),
	}
}

func UserDefaultWorkloadIdentityAuthenticationMethods(t *testing.T, userId sdk.AccountObjectIdentifier) *UserWorkloadIdentityAuthenticationMethodsAssert {
	t.Helper()
	return UserWorkloadIdentityAuthenticationMethods(t, userId, "DEFAULT")
}

func (u *UserWorkloadIdentityAuthenticationMethodsAssert) HasLastUsedNotEmpty() *UserWorkloadIdentityAuthenticationMethodsAssert {
	u.AddAssertion(func(t *testing.T, o *sdk.UserWorkloadIdentityAuthenticationMethod) error {
		t.Helper()
		if o.LastUsed == (time.Time{}) {
			return fmt.Errorf("expected last used not empty; got: %v", o.LastUsed)
		}
		return nil
	})
	return u
}

func (u *UserWorkloadIdentityAuthenticationMethodsAssert) HasCreatedOnNotEmpty() *UserWorkloadIdentityAuthenticationMethodsAssert {
	u.AddAssertion(func(t *testing.T, o *sdk.UserWorkloadIdentityAuthenticationMethod) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return u
}

func (u *UserWorkloadIdentityAuthenticationMethodsAssert) HasNoComment() *UserWorkloadIdentityAuthenticationMethodsAssert {
	u.AddAssertion(func(t *testing.T, o *sdk.UserWorkloadIdentityAuthenticationMethod) error {
		t.Helper()
		if o.Comment != "" {
			return fmt.Errorf("expected comment to be empty; got: %s", o.Comment)
		}
		return nil
	})
	return u
}
