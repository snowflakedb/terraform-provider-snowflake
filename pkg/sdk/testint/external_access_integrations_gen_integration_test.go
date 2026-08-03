//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
				HasNoAllowedApiAuthenticationIntegrationsList().
				HasNoAllowedAuthenticationSecretsList().
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
				HasNoAllowedApiAuthenticationIntegrationsList(),
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
				HasAllowedAuthenticationSecretsAll(true),
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
				HasNoAllowedAuthenticationSecretsList(),
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
				HasAllowedApiAuthenticationIntegrationsList(apiAuth1.ID(), apiAuth2.ID()).
				HasAllowedAuthenticationSecretsList(secret1.ID(), secret2.ID()).
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
				HasAllowedApiAuthenticationIntegrationsList(apiAuth1.ID(), apiAuth2.ID()).
				HasAllowedAuthenticationSecretsList(secret1.ID(), secret2.ID()).
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
				HasNoAllowedApiAuthenticationIntegrationsList().
				HasAllowedAuthenticationSecretsAll(false).
				HasNoAllowedAuthenticationSecretsList().
				HasEnabled(true).
				HasComment(""),
		)
	})
}
