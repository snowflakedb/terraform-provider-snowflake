package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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
	return c.CreateUserWithRequest(t, sdk.NewCreateUserRequest(c.ids.RandomAccountObjectIdentifier()))
}

func (c *UserClient) CreateServiceUser(t *testing.T) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithRequest(t, sdk.NewCreateUserRequest(c.ids.RandomAccountObjectIdentifier()).
		WithObjectProperties(*sdk.NewUserObjectPropertiesRequest().WithUserType(sdk.UserTypeService)))
}

func (c *UserClient) CreateLegacyServiceUser(t *testing.T) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithRequest(t, sdk.NewCreateUserRequest(c.ids.RandomAccountObjectIdentifier()).
		WithObjectProperties(*sdk.NewUserObjectPropertiesRequest().WithUserType(sdk.UserTypeLegacyService)))
}

func (c *UserClient) CreateUserWithPrefix(t *testing.T, prefix string) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithRequest(t, sdk.NewCreateUserRequest(c.ids.RandomAccountObjectIdentifierWithPrefix(prefix)))
}

func (c *UserClient) CreateUserWithRequest(t *testing.T, request *sdk.CreateUserRequest) (*sdk.User, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Create(ctx, request)
	require.NoError(t, err)
	id := request.ID()
	user, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return user, c.DropUserFunc(t, id)
}

func (c *UserClient) Alter(t *testing.T, request *sdk.AlterUserRequest) {
	t.Helper()
	err := c.client().Alter(context.Background(), request)
	require.NoError(t, err)
}

func (c *UserClient) AlterCurrentUser(t *testing.T, alter func(id sdk.AccountObjectIdentifier) *sdk.AlterUserRequest) {
	t.Helper()
	ctx := context.Background()
	id, err := c.context.client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)
	err = c.client().Alter(ctx, alter(id))
	require.NoError(t, err)
}

func (c *UserClient) DropUserFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropUserRequest(id).WithIfExists(true))
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

	return c.client().DescribeDetails(ctx, id)
}

func (c *UserClient) Disable(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(id).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectProperties(*sdk.NewUserAlterObjectPropertiesRequest().WithDisabled(true))))
	require.NoError(t, err)
}

func (c *UserClient) SetDaysToExpiry(t *testing.T, id sdk.AccountObjectIdentifier, value int) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(id).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectProperties(*sdk.NewUserAlterObjectPropertiesRequest().WithDaysToExpiry(value))))
	require.NoError(t, err)
}

func (c *UserClient) SetType(t *testing.T, id sdk.AccountObjectIdentifier, userType sdk.UserType) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(id).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectProperties(*sdk.NewUserAlterObjectPropertiesRequest().WithUserType(userType))))
	require.NoError(t, err)
}

func (c *UserClient) UnsetType(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(id).
		WithUnset(*sdk.NewUserUnsetRequest().
			WithObjectProperties(*sdk.NewUserObjectPropertiesUnsetRequest().WithUserType(true))))
	require.NoError(t, err)
}

func (c *UserClient) SetLoginName(t *testing.T, id sdk.AccountObjectIdentifier, newLoginName string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(id).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectProperties(*sdk.NewUserAlterObjectPropertiesRequest().WithLoginName(newLoginName))))
	require.NoError(t, err)
}

func (c *UserClient) UnsetDefaultSecondaryRoles(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(id).
		WithUnset(*sdk.NewUserUnsetRequest().
			WithObjectProperties(*sdk.NewUserObjectPropertiesUnsetRequest().WithDefaultSecondaryRoles(true))))
	require.NoError(t, err)
}

func (c *UserClient) AddProgrammaticAccessToken(t *testing.T, userId sdk.AccountObjectIdentifier) (sdk.AddProgrammaticAccessTokenResult, func()) {
	t.Helper()
	name := c.ids.RandomAccountObjectIdentifier()

	return c.AddProgrammaticAccessTokenWithRequest(t, userId, sdk.NewAddUserProgrammaticAccessTokenRequest(userId, name))
}

func (c *UserClient) AddProgrammaticAccessTokenWithRequest(t *testing.T, userId sdk.AccountObjectIdentifier, request *sdk.AddUserProgrammaticAccessTokenRequest) (sdk.AddProgrammaticAccessTokenResult, func()) {
	t.Helper()
	ctx := context.Background()

	// Expire the token after 1 day to avoid valid leftover tokens.
	request.WithDaysToExpiry(1)

	token, err := c.context.client.Users.AddProgrammaticAccessToken(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, token)
	return *token, c.RemoveProgrammaticAccessTokenFunc(t, userId, sdk.NewAccountObjectIdentifier(token.TokenName))
}

func (c *UserClient) ModifyProgrammaticAccessToken(t *testing.T, request *sdk.ModifyUserProgrammaticAccessTokenRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().ModifyProgrammaticAccessToken(ctx, request)
	require.NoError(t, err)
}

func (c *UserClient) RemoveProgrammaticAccessTokenFunc(t *testing.T, userId sdk.AccountObjectIdentifier, tokenName sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.context.client.Users.RemoveProgrammaticAccessTokenSafely(ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(userId, tokenName))
		if err != nil && !errors.Is(err, sdk.ErrPatNotFound) {
			t.Errorf("failed to remove programmatic access token: %v", err)
		}
	}
}

func (c *UserClient) ShowProgrammaticAccessToken(t *testing.T, userId sdk.AccountObjectIdentifier, tokenName sdk.AccountObjectIdentifier) *sdk.ProgrammaticAccessToken {
	t.Helper()
	ctx := context.Background()

	token, err := c.context.client.Users.ShowProgrammaticAccessTokenByName(ctx, userId, tokenName)
	require.NoError(t, err)
	require.NotNil(t, token)
	return token
}

func (c *UserClient) ShowUserWorkloadIdentityAuthenticationMethodOptions(t *testing.T, id UserWorkloadIdentityAuthenticationMethodsObjectIdentifier) (*sdk.UserWorkloadIdentityAuthenticationMethod, error) {
	t.Helper()
	ctx := context.Background()

	methods, err := c.context.client.Users.ShowUserWorkloadIdentityAuthenticationMethodOptions(ctx, sdk.NewShowUserWorkloadIdentityAuthenticationMethodOptionsUserRequest(id.userId))
	if err != nil {
		return nil, err
	}
	wif, err := collections.FindFirst(methods, func(method sdk.UserWorkloadIdentityAuthenticationMethod) bool {
		return method.Name == id.name
	})
	if err != nil {
		return nil, err
	}
	return wif, nil
}

// SetOidcWorkloadIdentity sets the OIDC workload identity configuration for a user.
func (c *UserClient) SetOidcWorkloadIdentity(t *testing.T, userId sdk.AccountObjectIdentifier, issuer, subject string, audienceList ...string) {
	t.Helper()
	ctx := context.Background()

	audiences := make([]sdk.StringListItemWrapper, len(audienceList))
	for i, v := range audienceList {
		audiences[i] = sdk.StringListItemWrapper{Value: v}
	}

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(userId).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectProperties(*sdk.NewUserAlterObjectPropertiesRequest().
				WithWorkloadIdentity(*sdk.NewUserObjectWorkloadIdentityPropertiesRequest().
					WithOidcType(*sdk.NewUserObjectWorkloadIdentityOidcRequest().
						WithIssuer(issuer).
						WithSubject(subject).
						WithOidcAudienceList(audiences))))))
	require.NoError(t, err)
}

// SetGcpWorkloadIdentity sets the GCP workload identity configuration for a user.
func (c *UserClient) SetGcpWorkloadIdentity(t *testing.T, userId sdk.AccountObjectIdentifier, subject string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(userId).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectProperties(*sdk.NewUserAlterObjectPropertiesRequest().
				WithWorkloadIdentity(*sdk.NewUserObjectWorkloadIdentityPropertiesRequest().
					WithGcpType(*sdk.NewUserObjectWorkloadIdentityGcpRequest().
						WithSubject(subject))))))
	require.NoError(t, err)
}

// SetAzureWorkloadIdentity sets the Azure workload identity configuration for a user.
func (c *UserClient) SetAzureWorkloadIdentity(t *testing.T, userId sdk.AccountObjectIdentifier, issuer, subject string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(userId).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectProperties(*sdk.NewUserAlterObjectPropertiesRequest().
				WithWorkloadIdentity(*sdk.NewUserObjectWorkloadIdentityPropertiesRequest().
					WithAzureType(*sdk.NewUserObjectWorkloadIdentityAzureRequest().
						WithIssuer(issuer).
						WithSubject(subject))))))
	require.NoError(t, err)
}

// SetAwsWorkloadIdentity sets the AWS workload identity configuration for a user.
func (c *UserClient) SetAwsWorkloadIdentity(t *testing.T, userId sdk.AccountObjectIdentifier, arn string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(userId).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectProperties(*sdk.NewUserAlterObjectPropertiesRequest().
				WithWorkloadIdentity(*sdk.NewUserObjectWorkloadIdentityPropertiesRequest().
					WithAwsType(*sdk.NewUserObjectWorkloadIdentityAwsRequest().
						WithArn(arn))))))
	require.NoError(t, err)
}

// UnsetWorkloadIdentity removes the workload identity configuration for a user.
func (c *UserClient) UnsetWorkloadIdentity(t *testing.T, userId sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(userId).
		WithUnset(*sdk.NewUserUnsetRequest().
			WithObjectProperties(*sdk.NewUserObjectPropertiesUnsetRequest().WithWorkloadIdentity(true))))
	require.NoError(t, err)
}

func (c *UserClient) UpdateEnableUnredactedQuerySyntaxError(t *testing.T, userId sdk.AccountObjectIdentifier, newValue bool) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterUserRequest(userId).
		WithSet(*sdk.NewUserSetRequest().
			WithObjectParameters(*sdk.NewUserObjectParametersRequest().WithEnableUnredactedQuerySyntaxError(newValue))))
	require.NoError(t, err)
}
