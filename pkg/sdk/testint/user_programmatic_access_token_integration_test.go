//go:build non_account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_UserProgrammaticAccessToken(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	currentUser, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)

	user, userCleanup := testClientHelper().User.CreateUser(t)
	t.Cleanup(userCleanup)

	t.Run("add - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewAddUserProgrammaticAccessTokenRequest(user.ID(), id)

		token, err := client.Users.AddProgrammaticAccessToken(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.RemoveProgrammaticAccessTokenFunc(t, user.ID(), id))

		// Assert the values returned during token creation.
		assert.NotNil(t, token)
		assert.Equal(t, id.Name(), token.TokenName)
		assert.NotEmpty(t, token.TokenSecret)

		// Assert the token is shown in the SHOW command.
		tokenShowObject := testClientHelper().User.ShowProgrammaticAccessToken(t, user.ID(), id)

		assertThatObject(t, objectassert.ProgrammaticAccessTokenFromObject(t, tokenShowObject).
			HasName(id.Name()).
			HasUserName(user.ID()).
			HasNoRoleRestriction().
			HasExpiresAtNotEmpty().
			HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
			HasNoComment().
			HasCreatedOnNotEmpty().
			HasCreatedBy(currentUser.Name()).
			HasNoMinsToBypassNetworkPolicyRequirement().
			HasRotatedToEmpty(),
		)
	})

	t.Run("add - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		request := sdk.NewAddUserProgrammaticAccessTokenRequest(user.ID(), id).
			WithRoleRestriction(snowflakeroles.Public).
			WithDaysToExpiry(1).
			WithMinsToBypassNetworkPolicyRequirement(10).
			WithComment(comment)

		token, err := client.Users.AddProgrammaticAccessToken(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.RemoveProgrammaticAccessTokenFunc(t, user.ID(), id))

		// Assert the values returned during token creation.
		assert.NotNil(t, token)
		assert.Equal(t, id.Name(), token.TokenName)
		assert.NotEmpty(t, token.TokenSecret)

		// Assert the token is shown in the SHOW command.
		tokenShowObject := testClientHelper().User.ShowProgrammaticAccessToken(t, user.ID(), id)

		assertThatObject(t, objectassert.ProgrammaticAccessTokenFromObject(t, tokenShowObject).
			HasName(id.Name()).
			HasUserName(user.ID()).
			HasRoleRestriction(snowflakeroles.Public).
			// Assert that WithDaysToExpiry(1) takes effect - the expires_at date is before 2 days from now.
			HasExpiresAtBefore(time.Now().Add(time.Hour*24*2)).
			HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
			HasComment(comment).
			HasCreatedOnNotEmpty().
			HasCreatedBy(currentUser.Name()).
			HasMinsToBypassNetworkPolicyRequirementWithTolerance(10).
			HasRotatedToEmpty(),
		)
	})

	t.Run("modify - set and unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewAddUserProgrammaticAccessTokenRequest(user.ID(), id)

		_, err := client.Users.AddProgrammaticAccessToken(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.RemoveProgrammaticAccessTokenFunc(t, user.ID(), id))

		tokenShowObject := testClientHelper().User.ShowProgrammaticAccessToken(t, user.ID(), id)
		assertThatObject(t, objectassert.ProgrammaticAccessTokenFromObject(t, tokenShowObject).
			HasName(id.Name()).
			HasNoComment().
			HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
			HasNoMinsToBypassNetworkPolicyRequirement(),
		)

		comment := random.Comment()
		mins := 15
		setRequest := sdk.NewModifyUserProgrammaticAccessTokenRequest(user.ID(), id).
			WithSet(*sdk.NewModifyProgrammaticAccessTokenSetRequest().
				WithComment(comment).
				WithDisabled(true).
				WithMinsToBypassNetworkPolicyRequirement(mins),
			)

		err = client.Users.ModifyProgrammaticAccessToken(ctx, setRequest)
		require.NoError(t, err)

		tokenShowObject = testClientHelper().User.ShowProgrammaticAccessToken(t, user.ID(), id)
		assertThatObject(t, objectassert.ProgrammaticAccessTokenFromObject(t, tokenShowObject).
			HasName(id.Name()).
			HasComment(comment).
			HasStatus(sdk.ProgrammaticAccessTokenStatusDisabled).
			HasMinsToBypassNetworkPolicyRequirementWithTolerance(mins),
		)

		unsetRequest := sdk.NewModifyUserProgrammaticAccessTokenRequest(user.ID(), id).
			WithUnset(*sdk.NewModifyProgrammaticAccessTokenUnsetRequest().
				WithComment(true).
				WithDisabled(true).
				WithMinsToBypassNetworkPolicyRequirement(true),
			)

		err = client.Users.ModifyProgrammaticAccessToken(ctx, unsetRequest)
		require.NoError(t, err)

		tokenShowObject = testClientHelper().User.ShowProgrammaticAccessToken(t, user.ID(), id)
		assertThatObject(t, objectassert.ProgrammaticAccessTokenFromObject(t, tokenShowObject).
			HasName(id.Name()).
			HasNoComment().
			HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
			HasNoMinsToBypassNetworkPolicyRequirement(),
		)
	})

	t.Run("modify - rename", func(t *testing.T) {
		oldId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewAddUserProgrammaticAccessTokenRequest(user.ID(), oldId)

		_, cleanupToken := testClientHelper().User.AddProgrammaticAccessTokenWithRequest(t, user.ID(), request)
		t.Cleanup(cleanupToken)

		renameRequest := sdk.NewModifyUserProgrammaticAccessTokenRequest(user.ID(), oldId).
			WithRenameTo(newId)

		err := client.Users.ModifyProgrammaticAccessToken(ctx, renameRequest)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.RemoveProgrammaticAccessTokenFunc(t, user.ID(), newId))

		_, err = client.Users.ShowProgrammaticAccessTokenByName(ctx, user.ID(), oldId)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		tokenShowObject, err := client.Users.ShowProgrammaticAccessTokenByName(ctx, user.ID(), newId)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ProgrammaticAccessTokenFromObject(t, tokenShowObject).
			HasName(newId.Name()),
		)
	})

	t.Run("rotate", func(t *testing.T) {
		token, cleanupToken := testClientHelper().User.AddProgrammaticAccessToken(t, user.ID())
		t.Cleanup(cleanupToken)

		rotateRequest := sdk.NewRotateUserProgrammaticAccessTokenRequest(user.ID(), token.ID()).
			WithExpireRotatedTokenAfterHours(0)
		rotateResult, err := client.Users.RotateProgrammaticAccessToken(ctx, rotateRequest)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.RemoveProgrammaticAccessTokenFunc(t, user.ID(), sdk.NewAccountObjectIdentifier(rotateResult.RotatedTokenName)))

		// Assert the values returned during token rotation.
		assert.NotNil(t, rotateResult)
		assert.Equal(t, token.ID().Name(), rotateResult.TokenName)
		assert.NotEmpty(t, rotateResult.TokenSecret)
		prefix := fmt.Sprintf("%s_ROTATED_", token.ID().Name())
		assert.True(t, strings.HasPrefix(rotateResult.RotatedTokenName, prefix))

		// Assert two tokens are shown in the SHOW command.
		showRequest := sdk.NewShowUserProgrammaticAccessTokenRequest().WithUserName(user.ID())
		showResult, err := client.Users.ShowProgrammaticAccessTokens(ctx, showRequest)
		require.NoError(t, err)
		require.NotNil(t, showResult)
		require.Len(t, showResult, 2)

		oldToken, err := collections.FindFirst(showResult, func(t sdk.ProgrammaticAccessToken) bool {
			return t.RotatedTo != nil && *t.RotatedTo == token.ID().Name()
		})
		require.NoError(t, err)

		newToken, err := collections.FindFirst(showResult, func(t sdk.ProgrammaticAccessToken) bool {
			return t.RotatedTo == nil
		})
		require.NoError(t, err)

		assertThatObject(t, objectassert.ProgrammaticAccessTokenFromObject(t, oldToken).
			HasName(rotateResult.RotatedTokenName).
			HasUserName(user.ID()).
			HasNoRoleRestriction().
			HasExpiresAtNotEmpty().
			HasStatus(sdk.ProgrammaticAccessTokenStatusExpired).
			HasNoComment().
			HasCreatedOnNotEmpty().
			HasCreatedBy(currentUser.Name()).
			HasNoMinsToBypassNetworkPolicyRequirement().
			HasRotatedTo(token.ID().Name()),
		)
		assertThatObject(t, objectassert.ProgrammaticAccessTokenFromObject(t, newToken).
			HasName(token.ID().Name()).
			HasUserName(user.ID()).
			HasNoRoleRestriction().
			HasExpiresAtNotEmpty().
			HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
			HasNoComment().
			HasCreatedOnNotEmpty().
			HasCreatedBy(currentUser.Name()).
			HasNoMinsToBypassNetworkPolicyRequirement().
			HasRotatedToEmpty(),
		)
	})

	t.Run("remove", func(t *testing.T) {
		token, cleanupToken := testClientHelper().User.AddProgrammaticAccessToken(t, user.ID())
		t.Cleanup(cleanupToken)

		err := client.Users.RemoveProgrammaticAccessToken(ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(user.ID(), token.ID()))
		require.NoError(t, err)

		_, err = client.Users.ShowProgrammaticAccessTokenByName(ctx, user.ID(), token.ID())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("show", func(t *testing.T) {
		token1, cleanupToken1 := testClientHelper().User.AddProgrammaticAccessToken(t, user.ID())
		t.Cleanup(cleanupToken1)
		token2, cleanupToken2 := testClientHelper().User.AddProgrammaticAccessToken(t, user.ID())
		t.Cleanup(cleanupToken2)

		showRequest := sdk.NewShowUserProgrammaticAccessTokenRequest().WithUserName(user.ID())
		showResult, err := client.Users.ShowProgrammaticAccessTokens(ctx, showRequest)
		require.NoError(t, err)
		require.NotNil(t, showResult)
		require.Len(t, showResult, 2)

		token1ShowObject, err := collections.FindFirst(showResult, func(t sdk.ProgrammaticAccessToken) bool {
			return t.Name == token1.TokenName
		})
		require.NoError(t, err)
		require.NotNil(t, token1ShowObject)

		token2ShowObject, err := collections.FindFirst(showResult, func(t sdk.ProgrammaticAccessToken) bool {
			return t.Name == token2.TokenName
		})
		require.NoError(t, err)
		require.NotNil(t, token2ShowObject)
	})
}
