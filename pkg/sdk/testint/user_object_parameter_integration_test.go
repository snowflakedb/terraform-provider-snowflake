//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_UserSetObjectParameterViaCatalog(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	user, userCleanup := testClientHelper().User.CreateUser(t)
	t.Cleanup(userCleanup)

	object := sdk.Object{
		ObjectType: sdk.ObjectTypeUser,
		Name:       user.ID(),
	}
	const value = "en-ci-upper"
	parameter := sdk.ObjectParameter(sdk.UserParameterQueryTag)

	handled, err := client.Parameters.SetObjectParameterOnObjectViaCatalog(ctx, object, parameter, value)
	require.NoError(t, err)
	require.True(t, handled)
	t.Cleanup(func() {
		handled, err := client.Parameters.UnsetObjectParameterOnObjectViaCatalog(ctx, object, parameter)
		require.NoError(t, err)
		require.True(t, handled)
	})

	param, err := client.Parameters.ShowUserParameter(ctx, sdk.UserParameterQueryTag, user.ID())
	require.NoError(t, err)
	require.Equal(t, value, param.Value)
	require.Equal(t, sdk.ParameterTypeUser, param.Level)
}
