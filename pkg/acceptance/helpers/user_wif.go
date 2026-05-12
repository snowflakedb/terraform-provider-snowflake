package helpers

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// UserWorkloadIdentityAuthenticationMethodsObjectIdentifier is a identifier for a user workload identity authentication method.
// It is a pseudo-identifier to satisfy the ObjectIdentifier interface. It can be used only in tests.
type UserWorkloadIdentityAuthenticationMethodsObjectIdentifier struct {
	userId sdk.AccountObjectIdentifier
	name   string
}

func NewUserWorkloadIdentityAuthenticationMethodsObjectIdentifier(userId sdk.AccountObjectIdentifier, name string) UserWorkloadIdentityAuthenticationMethodsObjectIdentifier {
	return UserWorkloadIdentityAuthenticationMethodsObjectIdentifier{
		userId: userId,
		name:   name,
	}
}

func (i UserWorkloadIdentityAuthenticationMethodsObjectIdentifier) FullyQualifiedName() string {
	return fmt.Sprintf("%s.%s", i.userId.FullyQualifiedName(), i.name)
}

func (i UserWorkloadIdentityAuthenticationMethodsObjectIdentifier) Name() string {
	return i.name
}
