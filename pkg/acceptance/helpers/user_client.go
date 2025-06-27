package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type UserClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewUserClient(context *TestClientContext, idsGenerator *IdsGenerator) *UserClient {
	return &UserClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *UserClient) client() sdk.Users {
	return c.context.client.Users
}

func (c *UserClient) CreateUser(t *testing.T) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithOptions(t, c.ids.RandomAccountObjectIdentifier(), &sdk.CreateUserOptions{})
}

func (c *UserClient) CreateServiceUser(t *testing.T) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithOptions(t, c.ids.RandomAccountObjectIdentifier(), &sdk.CreateUserOptions{
		ObjectProperties: &sdk.UserObjectProperties{
			Type: sdk.Pointer(sdk.UserTypeService),
		},
	})
}

func (c *UserClient) CreateLegacyServiceUser(t *testing.T) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithOptions(t, c.ids.RandomAccountObjectIdentifier(), &sdk.CreateUserOptions{
		ObjectProperties: &sdk.UserObjectProperties{
			Type: sdk.Pointer(sdk.UserTypeLegacyService),
		},
	})
}

func (c *UserClient) CreateUserWithPrefix(t *testing.T, prefix string) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithOptions(t, c.ids.RandomAccountObjectIdentifierWithPrefix(prefix), &sdk.CreateUserOptions{})
}

func (c *UserClient) CreateUserWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateUserOptions) (*sdk.User, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)
	user, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return user, c.DropUserFunc(t, id)
}

func (c *UserClient) Alter(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.AlterUserOptions) {
	t.Helper()
	err := c.client().Alter(context.Background(), id, opts)
	require.NoError(t, err)
}

func (c *UserClient) AlterCurrentUser(t *testing.T, opts *sdk.AlterUserOptions) {
	t.Helper()
	id, err := c.context.client.ContextFunctions.CurrentUser(context.Background())
	require.NoError(t, err)
	err = c.client().Alter(context.Background(), id, opts)
	require.NoError(t, err)
}

func (c *UserClient) DropUserFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropUserOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}

func (c *UserClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.User, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *UserClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.UserDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, id)
}

func (c *UserClient) Disable(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					Disable: sdk.Bool(true),
				},
			},
		},
	})
	require.NoError(t, err)
}

func (c *UserClient) SetDaysToExpiry(t *testing.T, id sdk.AccountObjectIdentifier, value int) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					DaysToExpiry: sdk.Int(value),
				},
			},
		},
	})
	require.NoError(t, err)
}

func (c *UserClient) SetType(t *testing.T, id sdk.AccountObjectIdentifier, userType sdk.UserType) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					Type: sdk.Pointer(userType),
				},
			},
		},
	})
	require.NoError(t, err)
}

func (c *UserClient) UnsetType(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterUserOptions{
		Unset: &sdk.UserUnset{
			ObjectProperties: &sdk.UserObjectPropertiesUnset{
				Type: sdk.Bool(true),
			},
		},
	})
	require.NoError(t, err)
}

func (c *UserClient) SetLoginName(t *testing.T, id sdk.AccountObjectIdentifier, newLoginName string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					LoginName: sdk.String(newLoginName),
				},
			},
		},
	})
	require.NoError(t, err)
}

func (c *UserClient) UnsetDefaultSecondaryRoles(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterUserOptions{
		Unset: &sdk.UserUnset{
			ObjectProperties: &sdk.UserObjectPropertiesUnset{
				DefaultSecondaryRoles: sdk.Bool(true),
			},
		},
	})
	require.NoError(t, err)
}

func (c *UserClient) AddProgrammaticAccessToken(t *testing.T, id sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier) sdk.AddProgrammaticAccessTokenResult {
	t.Helper()
	ctx := context.Background()
	name := c.ids.Alpha()

	token, err := c.context.client.UserProgrammaticAccessTokens.Add(ctx, sdk.NewAddUserProgrammaticAccessTokenRequest(name).
		WithUser(id).
		WithRoleRestriction(roleId).
		// Expire the token after 1 day to avoid valid leftover tokens.
		WithDaysToExpiry(1),
	)
	require.NoError(t, err)
	require.Len(t, token, 1)
	return token[0]
}
