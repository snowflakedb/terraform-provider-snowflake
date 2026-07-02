package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO [SNOW-1827324]: add TestClient ref to each specific client, so that we enhance specific client and not the base one
// TODO [SNOW-1827324]: consider using these in other places where user is set up

func (c *TestClient) SetUpTemporaryLegacyServiceUser(t *testing.T) *TmpLegacyServiceUser {
	t.Helper()

	pass := random.Password()
	tmpUser := c.setUpTmpUserWithBasicAccess(t, func(userId sdk.AccountObjectIdentifier) (*sdk.User, func()) {
		return c.User.CreateUserWithRequest(t, sdk.NewCreateUserRequest(userId).
			WithObjectProperties(*sdk.NewUserObjectPropertiesRequest().
				WithUserType(sdk.UserTypeLegacyService).
				WithPassword(pass)))
	})

	return &TmpLegacyServiceUser{
		Pass:    pass,
		TmpUser: tmpUser,
	}
}

func (c *TestClient) SetUpTemporaryServiceUser(t *testing.T) *TmpServiceUser {
	t.Helper()

	pass := random.Password()
	privateKey, encryptedKey, publicKey, _ := random.GenerateRSAKeyPair(t, pass)
	tmpUser := c.setUpTmpUserWithBasicAccess(t, func(userId sdk.AccountObjectIdentifier) (*sdk.User, func()) {
		return c.User.CreateUserWithRequest(t, sdk.NewCreateUserRequest(userId).
			WithObjectProperties(*sdk.NewUserObjectPropertiesRequest().
				WithUserType(sdk.UserTypeLegacyService).
				WithRsaPublicKey(publicKey)))
	})

	return &TmpServiceUser{
		PublicKey:           publicKey,
		PrivateKey:          privateKey,
		EncryptedPrivateKey: encryptedKey,
		Pass:                pass,
		TmpUser:             tmpUser,
	}
}

func (c *TestClient) SetUpTemporaryLegacyServiceUserWithPat(t *testing.T) *TmpServiceUserWithPat {
	t.Helper()

	tmpUser := c.setUpTmpUserWithBasicAccess(t, func(userId sdk.AccountObjectIdentifier) (*sdk.User, func()) {
		return c.User.CreateUserWithRequest(t, sdk.NewCreateUserRequest(userId).
			WithObjectProperties(*sdk.NewUserObjectPropertiesRequest().
				WithUserType(sdk.UserTypeLegacyService)))
	})
	req := sdk.NewAddUserProgrammaticAccessTokenRequest(tmpUser.UserId, c.Ids.RandomAccountObjectIdentifier()).WithRoleRestriction(tmpUser.RoleId)
	pat, cleanupPat := c.User.AddProgrammaticAccessTokenWithRequest(t, tmpUser.UserId, req)
	t.Cleanup(cleanupPat)

	return &TmpServiceUserWithPat{
		Pat:     pat.TokenSecret,
		TmpUser: tmpUser,
	}
}

func (c *TestClient) SetUpTemporaryUserForOauthClientCredentials(t *testing.T, loginName string) *TmpUser {
	t.Helper()
	userId := c.Ids.RandomAccountObjectIdentifier()
	user, userCleanup := c.User.CreateUserWithRequest(t, sdk.NewCreateUserRequest(userId).
		WithObjectProperties(*sdk.NewUserObjectPropertiesRequest().
			WithMustChangePassword(false).
			WithLoginName(loginName)))
	t.Cleanup(userCleanup)

	return &TmpUser{
		UserId: user.ID(),
		RoleId: snowflakeroles.Public,
	}
}

func (c *TestClient) setUpTmpUserWithBasicAccess(t *testing.T, userCreator func(userId sdk.AccountObjectIdentifier) (*sdk.User, func())) TmpUser {
	t.Helper()

	warehouseId := c.Ids.SnowflakeWarehouseId()
	accountId := c.Context.CurrentAccountId(t)

	tmpUserId := c.Ids.RandomAccountObjectIdentifier()
	_, userCleanup := userCreator(tmpUserId)
	t.Cleanup(userCleanup)

	tmpRole, roleCleanup := c.Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	tmpRoleId := tmpRole.ID()

	c.Grant.GrantPrivilegesOnDatabaseToAccountRole(t, tmpRoleId, c.Ids.DatabaseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Grant.GrantPrivilegesOnWarehouseToAccountRole(t, tmpRoleId, warehouseId, []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Role.GrantRoleToUser(t, tmpRoleId, tmpUserId)

	return TmpUser{
		UserId:      tmpUserId,
		RoleId:      tmpRoleId,
		WarehouseId: warehouseId,
		AccountId:   accountId,
	}
}

type TmpUser struct {
	UserId      sdk.AccountObjectIdentifier
	RoleId      sdk.AccountObjectIdentifier
	WarehouseId sdk.AccountObjectIdentifier
	AccountId   sdk.AccountIdentifier
}

func (u *TmpUser) OrgAndAccount() string {
	return fmt.Sprintf("%s-%s", u.AccountId.OrganizationName(), u.AccountId.AccountName())
}

type TmpServiceUser struct {
	PublicKey           string
	PrivateKey          string
	EncryptedPrivateKey string
	Pass                string
	TmpUser
}

type TmpLegacyServiceUser struct {
	Pass string
	TmpUser
}

type TmpServiceUserWithPat struct {
	Pat string
	TmpUser
}
