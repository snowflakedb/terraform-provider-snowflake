//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalAccessIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	networkRule, networkRuleCleanup := testClientHelper().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalAccessIntegration.DropExternalAccessIntegrationFunc(t, id))

		assertThatObject(
			t, objectassert.ExternalAccessIntegration(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasTypeExternalAccess().
				HasCategoryExternalAccess().
				HasEnabled(true).
				HasComment(""),
		)
		assertThatObject(
			t, objectassert.ExternalAccessIntegrationDetails(t, id).
				HasAllowedNetworkRules(networkRule.ID()).
				HasNoAllowedApiAuthenticationIntegrations().
				HasNoAllowedAuthenticationSecrets().
				HasEnabled(true).
				HasComment(""),
		)
	})

	t.Run("create: without or replace fails", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalAccessIntegration.DropExternalAccessIntegrationFunc(t, id))

		err = client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true),
		)
		require.ErrorContains(t, err, "already exists")
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalAccessIntegration.DropExternalAccessIntegrationFunc(t, id))

		err = client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true).
				WithOrReplace(true).
				WithComment("replaced"),
		)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.ExternalAccessIntegration(t, id).
				HasName(id.Name()).
				HasComment("replaced"),
		)
	})

	t.Run("create: AllowedApiAuthenticationIntegrations = none", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true).
				WithAllowedApiAuthenticationIntegrations(*sdk.NewExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest().WithNone(true)),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalAccessIntegration.DropExternalAccessIntegrationFunc(t, id))

		assertThatObject(
			t, objectassert.ExternalAccessIntegrationDetails(t, id).
				HasNoAllowedApiAuthenticationIntegrations(),
		)
	})

	t.Run("create: AllowedAuthenticationSecrets = all", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true).
				WithAllowedAuthenticationSecrets(*sdk.NewExternalAccessIntegrationAllowedAuthenticationSecretsRequest().WithAll(true)),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalAccessIntegration.DropExternalAccessIntegrationFunc(t, id))

		assertThatObject(
			t, objectassert.ExternalAccessIntegrationDetails(t, id).
				HasAllowedAuthenticationSecrets("ALL"),
		)
	})

	t.Run("create: AllowedAuthenticationSecrets = none", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true).
				WithAllowedAuthenticationSecrets(*sdk.NewExternalAccessIntegrationAllowedAuthenticationSecretsRequest().WithNone(true)),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalAccessIntegration.DropExternalAccessIntegrationFunc(t, id))

		assertThatObject(
			t, objectassert.ExternalAccessIntegrationDetails(t, id).
				HasNoAllowedAuthenticationSecrets(),
		)
	})

	t.Run("create: all options", func(t *testing.T) {
		networkRule2, networkRule2Cleanup := testClientHelper().NetworkRule.Create(t)
		t.Cleanup(networkRule2Cleanup)

		apiAuth1, apiAuth1Cleanup := testClientHelper().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlow(t)
		t.Cleanup(apiAuth1Cleanup)

		apiAuth2, apiAuth2Cleanup := testClientHelper().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlow(t)
		t.Cleanup(apiAuth2Cleanup)

		secretId1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secret1, secret1Cleanup := testClientHelper().Secret.CreateWithGenericString(t, secretId1, "test-secret-1")
		t.Cleanup(secret1Cleanup)

		secretId2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secret2, secret2Cleanup := testClientHelper().Secret.CreateWithGenericString(t, secretId2, "test-secret-2")
		t.Cleanup(secret2Cleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID(), networkRule2.ID()}, false).
				WithAllowedApiAuthenticationIntegrations(*sdk.NewExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest().
					WithIntegrations([]sdk.AccountObjectIdentifier{apiAuth1.ID(), apiAuth2.ID()})).
				WithAllowedAuthenticationSecrets(*sdk.NewExternalAccessIntegrationAllowedAuthenticationSecretsRequest().
					WithSecrets([]sdk.SchemaObjectIdentifier{secret1.ID(), secret2.ID()})).
				WithComment("test comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalAccessIntegration.DropExternalAccessIntegrationFunc(t, id))

		assertThatObject(
			t, objectassert.ExternalAccessIntegration(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasTypeExternalAccess().
				HasCategoryExternalAccess().
				HasEnabled(false).
				HasComment("test comment"),
		)
		assertThatObject(
			t, objectassert.ExternalAccessIntegrationDetails(t, id).
				HasAllowedNetworkRules(networkRule.ID(), networkRule2.ID()).
				HasAllowedApiAuthenticationIntegrations(apiAuth1.ID().Name(), apiAuth2.ID().Name()).
				HasAllowedAuthenticationSecrets(secret1.ID().FullyQualifiedName(), secret2.ID().FullyQualifiedName()).
				HasEnabled(false).
				HasComment("test comment"),
		)
	})

	t.Run("alter: set and unset all", func(t *testing.T) {
		networkRule2, networkRule2Cleanup := testClientHelper().NetworkRule.Create(t)
		t.Cleanup(networkRule2Cleanup)

		apiAuth1, apiAuth1Cleanup := testClientHelper().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlow(t)
		t.Cleanup(apiAuth1Cleanup)

		apiAuth2, apiAuth2Cleanup := testClientHelper().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlow(t)
		t.Cleanup(apiAuth2Cleanup)

		secretId1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secret1, secret1Cleanup := testClientHelper().Secret.CreateWithGenericString(t, secretId1, "test-secret-1")
		t.Cleanup(secret1Cleanup)

		secretId2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secret2, secret2Cleanup := testClientHelper().Secret.CreateWithGenericString(t, secretId2, "test-secret-2")
		t.Cleanup(secret2Cleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalAccessIntegration.DropExternalAccessIntegrationFunc(t, id))

		err = client.ExternalAccessIntegrations.Alter(
			ctx,
			sdk.NewAlterExternalAccessIntegrationRequest(id).
				WithSet(
					*sdk.NewExternalAccessIntegrationSetRequest().
						WithAllowedNetworkRules([]sdk.SchemaObjectIdentifier{networkRule.ID(), networkRule2.ID()}).
						WithAllowedApiAuthenticationIntegrations(*sdk.NewExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest().
							WithIntegrations([]sdk.AccountObjectIdentifier{apiAuth1.ID(), apiAuth2.ID()})).
						WithAllowedAuthenticationSecrets(*sdk.NewExternalAccessIntegrationAllowedAuthenticationSecretsRequest().
							WithSecrets([]sdk.SchemaObjectIdentifier{secret1.ID(), secret2.ID()})).
						WithEnabled(false).
						WithComment("updated comment"),
				),
		)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.ExternalAccessIntegrationDetails(t, id).
				HasAllowedNetworkRules(networkRule.ID(), networkRule2.ID()).
				HasAllowedApiAuthenticationIntegrations(apiAuth1.ID().Name(), apiAuth2.ID().Name()).
				HasAllowedAuthenticationSecrets(secret1.ID().FullyQualifiedName(), secret2.ID().FullyQualifiedName()).
				HasEnabled(false).
				HasComment("updated comment"),
		)

		err = client.ExternalAccessIntegrations.Alter(
			ctx,
			sdk.NewAlterExternalAccessIntegrationRequest(id).
				WithUnset(
					*sdk.NewExternalAccessIntegrationUnsetRequest().
						WithAllowedNetworkRules(true).
						WithAllowedApiAuthenticationIntegrations(true).
						WithAllowedAuthenticationSecrets(true).
						WithEnabled(true).
						WithComment(true),
				),
		)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.ExternalAccessIntegrationDetails(t, id).
				HasNoAllowedNetworkRules().
				HasNoAllowedApiAuthenticationIntegrations().
				HasNoAllowedAuthenticationSecrets().
				HasEnabled(true).
				HasComment(""),
		)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true),
		)
		require.NoError(t, err)

		err = client.ExternalAccessIntegrations.Drop(ctx, sdk.NewDropExternalAccessIntegrationRequest(id))
		require.NoError(t, err)

		_, err = client.ExternalAccessIntegrations.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.ExternalAccessIntegrations.Drop(ctx, sdk.NewDropExternalAccessIntegrationRequest(NonExistingAccountObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists", func(t *testing.T) {
		err := client.ExternalAccessIntegrations.Drop(ctx, sdk.NewDropExternalAccessIntegrationRequest(NonExistingAccountObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalAccessIntegrations.Create(
			ctx,
			sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRule.ID()}, true),
		)
		require.NoError(t, err)

		err = client.ExternalAccessIntegrations.DropSafely(ctx, id)
		require.NoError(t, err)

		_, err = client.ExternalAccessIntegrations.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.ExternalAccessIntegrations.DropSafely(ctx, NonExistingAccountObjectIdentifier)
		require.NoError(t, err)
	})

	t.Run("show: with like filter", func(t *testing.T) {
		id1, id1Cleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
		t.Cleanup(id1Cleanup)

		id2, id2Cleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
		t.Cleanup(id2Cleanup)

		eai1, err := client.ExternalAccessIntegrations.ShowByID(ctx, id1)
		require.NoError(t, err)

		eai2, err := client.ExternalAccessIntegrations.ShowByID(ctx, id2)
		require.NoError(t, err)

		returnedIntegrations, err := client.ExternalAccessIntegrations.Show(ctx, sdk.NewShowExternalAccessIntegrationRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}))
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *eai1)
		assert.NotContains(t, returnedIntegrations, *eai2)
	})

	t.Run("show by id safely: existing", func(t *testing.T) {
		id, idCleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
		t.Cleanup(idCleanup)

		result, err := client.ExternalAccessIntegrations.ShowByIDSafely(ctx, id)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.ExternalAccessIntegrationFromObject(t, result).
				HasName(id.Name()).
				HasEnabled(true),
		)
	})

	t.Run("show by id safely: non-existing", func(t *testing.T) {
		_, err := client.ExternalAccessIntegrations.ShowByIDSafely(ctx, NonExistingAccountObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.ExternalAccessIntegrations.Describe(ctx, NonExistingAccountObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
