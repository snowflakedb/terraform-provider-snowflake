package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO(SNOW-2314062): Use in resource

type GrantDatabaseRoleId struct {
	DatabaseRoleName sdk.DatabaseObjectIdentifier
	ObjectType       sdk.ObjectType
	TargetIdentifier sdk.ObjectIdentifier
}

func (g *GrantDatabaseRoleId) String() string {
	return helpers.EncodeResourceIdentifier(g.DatabaseRoleName.FullyQualifiedName(), g.ObjectType.String(), g.TargetIdentifier.FullyQualifiedName())
}

func NewGrantDatabaseRoleIdToDatabaseRole(databaseRoleName sdk.DatabaseObjectIdentifier, parentRoleName sdk.DatabaseObjectIdentifier) GrantDatabaseRoleId {
	return GrantDatabaseRoleId{
		DatabaseRoleName: databaseRoleName,
		ObjectType:       sdk.ObjectTypeDatabaseRole,
		TargetIdentifier: parentRoleName,
	}
}

func NewGrantDatabaseRoleIdToRole(databaseRoleName sdk.DatabaseObjectIdentifier, parentRoleName sdk.AccountObjectIdentifier) GrantDatabaseRoleId {
	return GrantDatabaseRoleId{
		DatabaseRoleName: databaseRoleName,
		ObjectType:       sdk.ObjectTypeRole,
		TargetIdentifier: parentRoleName,
	}
}

func NewGrantDatabaseRoleIdToUser(databaseRoleName sdk.DatabaseObjectIdentifier, userName sdk.AccountObjectIdentifier) GrantDatabaseRoleId {
	return GrantDatabaseRoleId{
		DatabaseRoleName: databaseRoleName,
		ObjectType:       sdk.ObjectTypeRole,
		TargetIdentifier: userName,
	}
}

func NewGrantDatabaseRoleIdToShare(databaseRoleName sdk.DatabaseObjectIdentifier, shareName sdk.AccountObjectIdentifier) GrantDatabaseRoleId {
	return GrantDatabaseRoleId{
		DatabaseRoleName: databaseRoleName,
		ObjectType:       sdk.ObjectTypeRole,
		TargetIdentifier: shareName,
	}
}
