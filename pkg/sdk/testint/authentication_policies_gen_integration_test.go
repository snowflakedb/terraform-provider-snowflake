//go:build non_account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AuthenticationPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertAuthenticationPolicy := func(t *testing.T, id sdk.SchemaObjectIdentifier, expectedComment string) {
		t.Helper()
		assertThatObject(t,
			objectassert.AuthenticationPolicy(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasKind("AUTHENTICATION_POLICY").
				HasOptions("").
				HasOwner("ACCOUNTADMIN").
				HasComment(expectedComment).
				HasOwnerRoleType("ROLE"),
		)
	}

	assertProperty := func(t *testing.T, descriptions []sdk.AuthenticationPolicyDescription, name string, value string) {
		t.Helper()
		description, err := collections.FindFirst(descriptions, func(description sdk.AuthenticationPolicyDescription) bool {
			return description.Property == name
		})
		require.NoError(t, err, "unable to find property %s", name)
		assert.Equal(t, value, description.Value)
	}

	t.Run("Create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.AuthenticationPolicies.Create(ctx, sdk.NewCreateAuthenticationPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().AuthenticationPolicy.DropFunc(t, id))

		assertAuthenticationPolicy(t, id, "")

		desc, err := client.AuthenticationPolicies.Describe(ctx, id)
		require.NoError(t, err)

		assertProperty(t, desc, "COMMENT", "null")
		assertProperty(t, desc, "MFA_ENROLLMENT", string(sdk.MfaEnrollmentRequiredPasswordOnly))
		assertProperty(t, desc, "MFA_AUTHENTICATION_METHODS", "[PASSWORD]")
		assertProperty(t, desc, "SECURITY_INTEGRATIONS", "[ALL]")
		assertProperty(t, desc, "CLIENT_TYPES", "[ALL]")
		assertProperty(t, desc, "AUTHENTICATION_METHODS", "[ALL]")
	})

	t.Run("Create - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		samlIntegration, cleanupSamlIntegration := testClientHelper().SecurityIntegration.CreateSaml2(t)
		t.Cleanup(cleanupSamlIntegration)

		err := client.AuthenticationPolicies.Create(ctx, sdk.NewCreateAuthenticationPolicyRequest(id).
			WithComment(comment).
			WithMfaEnrollment(sdk.MfaEnrollmentOptional).
			WithSecurityIntegrations([]sdk.SecurityIntegrationsOption{
				{Name: samlIntegration.ID()},
			}).
			WithClientTypes([]sdk.ClientTypes{
				{ClientType: sdk.ClientTypesDrivers},
				{ClientType: sdk.ClientTypesSnowSql},
			}).
			WithAuthenticationMethods([]sdk.AuthenticationMethods{
				{Method: sdk.AuthenticationMethodsPassword},
				{Method: sdk.AuthenticationMethodsSaml},
			}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().AuthenticationPolicy.DropFunc(t, id))

		assertAuthenticationPolicy(t, id, comment)

		desc, err := client.AuthenticationPolicies.Describe(ctx, id)
		require.NoError(t, err)

		assertProperty(t, desc, "COMMENT", comment)
		assertProperty(t, desc, "MFA_ENROLLMENT", string(sdk.MfaEnrollmentRequiredPasswordOnly))
		assertProperty(t, desc, "MFA_AUTHENTICATION_METHODS", "[PASSWORD]")
		assertProperty(t, desc, "SECURITY_INTEGRATIONS", fmt.Sprintf("[%s]", samlIntegration.ID().Name()))
		assertProperty(t, desc, "CLIENT_TYPES", "[DRIVERS, SNOWSQL]")
		assertProperty(t, desc, "AUTHENTICATION_METHODS", "[PASSWORD, SAML]")
	})

	t.Run("Alter - set and unset properties", func(t *testing.T) {
		comment := random.Comment()

		authenticationPolicy, cleanupAuthPolicy := testClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(cleanupAuthPolicy)

		samlIntegration, cleanupSamlIntegration := testClientHelper().SecurityIntegration.CreateSaml2(t)
		t.Cleanup(cleanupSamlIntegration)

		err := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(authenticationPolicy.ID()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().
				WithComment(comment).
				WithMfaEnrollment(sdk.MfaEnrollmentRequired).
				WithSecurityIntegrations([]sdk.SecurityIntegrationsOption{
					{Name: samlIntegration.ID()},
				}).
				WithClientTypes([]sdk.ClientTypes{
					{ClientType: sdk.ClientTypesDrivers},
					{ClientType: sdk.ClientTypesSnowSql},
					{ClientType: sdk.ClientTypesSnowflakeUi},
				}).
				WithAuthenticationMethods([]sdk.AuthenticationMethods{
					{Method: sdk.AuthenticationMethodsPassword},
					{Method: sdk.AuthenticationMethodsSaml},
				})))
		require.NoError(t, err)

		desc, err := client.AuthenticationPolicies.Describe(ctx, authenticationPolicy.ID())
		require.NoError(t, err)

		assertProperty(t, desc, "COMMENT", comment)
		assertProperty(t, desc, "MFA_ENROLLMENT", "REQUIRED")
		assertProperty(t, desc, "MFA_AUTHENTICATION_METHODS", "[PASSWORD]")
		assertProperty(t, desc, "SECURITY_INTEGRATIONS", fmt.Sprintf("[%s]", samlIntegration.ID().Name()))
		assertProperty(t, desc, "CLIENT_TYPES", "[DRIVERS, SNOWSQL, SNOWFLAKE_UI]")
		assertProperty(t, desc, "AUTHENTICATION_METHODS", "[PASSWORD, SAML]")

		err = client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(authenticationPolicy.ID()).
			WithUnset(*sdk.NewAuthenticationPolicyUnsetRequest().
				WithComment(true).
				WithMfaEnrollment(true).
				WithSecurityIntegrations(true).
				WithClientTypes(true).
				WithAuthenticationMethods(true)))
		require.NoError(t, err)

		desc, err = client.AuthenticationPolicies.Describe(ctx, authenticationPolicy.ID())
		require.NoError(t, err)

		assertProperty(t, desc, "COMMENT", "null")
		assertProperty(t, desc, "MFA_ENROLLMENT", string(sdk.MfaEnrollmentRequiredPasswordOnly))
		assertProperty(t, desc, "MFA_AUTHENTICATION_METHODS", "[PASSWORD]")
		assertProperty(t, desc, "SECURITY_INTEGRATIONS", "[ALL]")
		assertProperty(t, desc, "CLIENT_TYPES", "[ALL]")
		assertProperty(t, desc, "AUTHENTICATION_METHODS", "[ALL]")
	})

	t.Run("Alter - rename", func(t *testing.T) {
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		authenticationPolicy, cleanupAuthPolicy := testClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(cleanupAuthPolicy)
		t.Cleanup(testClientHelper().AuthenticationPolicy.DropFunc(t, newId))

		err := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(authenticationPolicy.ID()).WithRenameTo(newId))
		require.NoError(t, err)

		_, err = client.AuthenticationPolicies.Describe(ctx, authenticationPolicy.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		_, err = client.AuthenticationPolicies.Describe(ctx, newId)
		assert.NoError(t, err)
	})

	t.Run("Drop - existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.AuthenticationPolicies.Create(ctx, sdk.NewCreateAuthenticationPolicyRequest(id))
		require.NoError(t, err)

		err = client.AuthenticationPolicies.Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.AuthenticationPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("Drop - non-existing", func(t *testing.T) {
		err := client.AuthenticationPolicies.Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_auth_policyzzz")
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_auth_policy_2_")
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_auth_policy_3_")
		id4 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(sdk.NewDatabaseObjectIdentifier(db.Name, "PUBLIC"))

		_, authenticationPolicyCleanup := testClientHelper().AuthenticationPolicy.CreateWithOptions(t, id, sdk.NewCreateAuthenticationPolicyRequest(id))
		t.Cleanup(authenticationPolicyCleanup)

		_, authenticationPolicyCleanup2 := testClientHelper().AuthenticationPolicy.CreateWithOptions(t, id2, sdk.NewCreateAuthenticationPolicyRequest(id2))
		t.Cleanup(authenticationPolicyCleanup2)

		_, authenticationPolicyCleanup3 := testClientHelper().AuthenticationPolicy.CreateWithOptions(t, id3, sdk.NewCreateAuthenticationPolicyRequest(id3))
		t.Cleanup(authenticationPolicyCleanup3)

		_, authenticationPolicyCleanup4 := testClientHelper().AuthenticationPolicy.CreateWithOptions(t, id4, sdk.NewCreateAuthenticationPolicyRequest(id4))
		t.Cleanup(authenticationPolicyCleanup4)

		t.Run("like", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().
				WithLike(sdk.Like{Pattern: sdk.String("test_auth_policy_2_%")}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, authenticationPolicies, 1)
		})

		t.Run("starts_with", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().
				WithStartsWith("test_auth_policy_").
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, authenticationPolicies, 2)
		})

		t.Run("in_account", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}}))
			require.NoError(t, err)
			assert.Greater(t, len(authenticationPolicies), 3)
		})

		t.Run("in_database", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Database: id.DatabaseId()}}))
			require.NoError(t, err)
			assert.Len(t, authenticationPolicies, 3)
		})

		t.Run("in_schema", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, authenticationPolicies, 3)
		})

		t.Run("on_account", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().WithOn(sdk.On{Account: sdk.Pointer(true)}))
			require.NoError(t, err)
			assert.Len(t, authenticationPolicies, 1)
		})

		t.Run("on_user", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().WithOn(sdk.On{User: sdk.NewAccountObjectIdentifier("user_name")}))
			require.NoError(t, err)
			assert.Len(t, authenticationPolicies, 1)
		})

		t.Run("limit", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, authenticationPolicies, 1)
		})

		t.Run("limit from", func(t *testing.T) {
			authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1), From: sdk.String("test_auth_policy_")}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
			require.NoError(t, err)
			require.Len(t, authenticationPolicies, 1)
			require.True(t, strings.HasPrefix(authenticationPolicies[0].Name, "test_auth_policy_2"))
		})
	})

	t.Run("Describe", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		authenticationPolicy, cleanupAuthPolicy := testClientHelper().AuthenticationPolicy.CreateWithOptions(t, id, sdk.NewCreateAuthenticationPolicyRequest(id).WithComment("some_comment"))
		t.Cleanup(cleanupAuthPolicy)

		desc, err := client.AuthenticationPolicies.Describe(ctx, authenticationPolicy.ID())
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(desc), 9)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{
			Property:    "COMMENT",
			Value:       "some_comment",
			Default:     "null",
			Description: "Comment associated with authentication policy.",
		})
	})
}
