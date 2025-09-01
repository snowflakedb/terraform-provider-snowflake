package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO(): Use in resource

type GrantAccountRoleId struct {
	AccountRoleName  sdk.AccountObjectIdentifier
	ObjectType       sdk.ObjectType
	TargetIdentifier sdk.AccountObjectIdentifier
}

func (g *GrantAccountRoleId) String() string {
	return helpers.EncodeResourceIdentifier(g.AccountRoleName.FullyQualifiedName(), g.ObjectType.String(), g.TargetIdentifier.FullyQualifiedName())
}

func NewGrantAccountRoleIdToRole(accountRoleName sdk.AccountObjectIdentifier, parentRoleName sdk.AccountObjectIdentifier) GrantAccountRoleId {
	return GrantAccountRoleId{
		AccountRoleName:  accountRoleName,
		ObjectType:       sdk.ObjectTypeRole,
		TargetIdentifier: parentRoleName,
	}
}

func NewGrantAccountRoleIdToUser(accountRoleName sdk.AccountObjectIdentifier, userName sdk.AccountObjectIdentifier) GrantAccountRoleId {
	return GrantAccountRoleId{
		AccountRoleName:  accountRoleName,
		ObjectType:       sdk.ObjectTypeUser,
		TargetIdentifier: userName,
	}
}
