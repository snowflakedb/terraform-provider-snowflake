package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type UserWorkloadIdentityAuthenticationMethodsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.UserWorkloadIdentityAuthenticationMethods, sdk.AccountObjectIdentifier]
}

func UserWorkloadIdentityAuthenticationMethodsFromObject(t *testing.T, userId sdk.AccountObjectIdentifier, userWorkloadIdentityAuthenticationMethods *sdk.UserWorkloadIdentityAuthenticationMethods) *UserWorkloadIdentityAuthenticationMethodsAssert {
	t.Helper()
	return &UserWorkloadIdentityAuthenticationMethodsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("USER WORKLOAD IDENTITY AUTHENTICATION METHODS"), userId, userWorkloadIdentityAuthenticationMethods),
	}
}

func (w *UserWorkloadIdentityAuthenticationMethodsAssert) HasOidcAdditionalInfo(expected *sdk.UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo) *UserWorkloadIdentityAuthenticationMethodsAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.UserWorkloadIdentityAuthenticationMethods) error {
		t.Helper()
		if o.OidcAdditionalInfo == nil {
			return fmt.Errorf("expected oidc additional info not nil")
		}
		if o.AwsAdditionalInfo != nil {
			return fmt.Errorf("expected aws additional info nil; got: %v", o.AwsAdditionalInfo)
		}
		if o.AzureAdditionalInfo != nil {
			return fmt.Errorf("expected azure additional info nil; got: %v", o.AzureAdditionalInfo)
		}
		if o.GcpAdditionalInfo != nil {
			return fmt.Errorf("expected gcp additional info nil; got: %v", o.GcpAdditionalInfo)
		}
		if o.OidcAdditionalInfo.Issuer != expected.Issuer {
			return fmt.Errorf("expected oidc additional info issuer: %v; got: %v", expected.Issuer, o.OidcAdditionalInfo.Issuer)
		}
		if o.OidcAdditionalInfo.Subject != expected.Subject {
			return fmt.Errorf("expected oidc additional info subject: %v; got: %v", expected.Subject, o.OidcAdditionalInfo.Subject)
		}
		if !slices.Equal(o.OidcAdditionalInfo.AudienceList, expected.AudienceList) {
			return fmt.Errorf("expected oidc additional info audience list: %v; got: %v", expected.AudienceList, o.OidcAdditionalInfo.AudienceList)
		}
		return nil
	})
	return w
}
