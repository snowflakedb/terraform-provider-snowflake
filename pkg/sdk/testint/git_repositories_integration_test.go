//go:build !account_level_tests

package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_GitRepositories(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	origin := "https://github.com/octocat/hello-world"

	createGitRepositoryWithSecretAndComment := func(t *testing.T) (sdk.SchemaObjectIdentifier, sdk.AccountObjectIdentifier, sdk.SchemaObjectIdentifier) {
		t.Helper()

		gitRepositoryId := testClientHelper().Ids.
			RandomSchemaObjectIdentifierInSchema(schema.ID())

		apiIntegrationId, apiIntegrationCleanup :=
			testClientHelper().ApiIntegration.
				CreateApiIntegrationForGitRepository(t, origin)
		t.Cleanup(apiIntegrationCleanup)

		secretId := testClientHelper().Ids.
			RandomSchemaObjectIdentifierInSchema(schema.ID())
		_, secretCleanup := testClientHelper().Secret.
			CreateWithBasicAuthenticationFlow(t, secretId, "username", "password")
		t.Cleanup(secretCleanup)

		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.
			CreatWithSecretAndComment(t, gitRepositoryId, origin, apiIntegrationId, secretId, "comment")
		t.Cleanup(gitRepositoryCleanup)

		return gitRepositoryId, apiIntegrationId, secretId
	}

	createGitRepository := func(t *testing.T) (sdk.SchemaObjectIdentifier, sdk.AccountObjectIdentifier) {
		t.Helper()

		gitRepositoryId := testClientHelper().Ids.
			RandomSchemaObjectIdentifierInSchema(schema.ID())

		apiIntegrationId, apiIntegrationCleanup :=
			testClientHelper().ApiIntegration.
				CreateApiIntegrationForGitRepository(t, origin)
		t.Cleanup(apiIntegrationCleanup)

		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.
			Create(t, gitRepositoryId, origin, apiIntegrationId)
		t.Cleanup(gitRepositoryCleanup)

		return gitRepositoryId, apiIntegrationId
	}

	t.Run("create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		apiIntegration, apiIntegrationCleanup := testClientHelper().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)
		t.Cleanup(apiIntegrationCleanup)

		request := sdk.NewCreateGitRepositoryRequest(id, origin, apiIntegration)

		err := client.GitRepositories.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().GitRepository.DropGitRepositoryFunc(t, id))

		gitRepository, err := client.GitRepositories.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.GitRepositoryFromObject(t, gitRepository).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOrigin(origin).
			HasApiIntegration(apiIntegration).
			HasGitCredentialsEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		)
	})

	t.Run("create - complete", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		apiIntegration, apiIntegrationCleanup := testClientHelper().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)
		t.Cleanup(apiIntegrationCleanup)

		secretId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		_, secretCleanup := testClientHelper().Secret.CreateWithBasicAuthenticationFlow(t, secretId, "username", "password")
		t.Cleanup(secretCleanup)

		request := sdk.NewCreateGitRepositoryRequest(gitRepositoryId, origin, apiIntegration).WithIfNotExists(true).WithGitCredentials(secretId).WithComment("comment")

		err := client.GitRepositories.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().GitRepository.DropGitRepositoryFunc(t, gitRepositoryId))

		gitRepository, err := client.GitRepositories.ShowByID(ctx, gitRepositoryId)
		require.NoError(t, err)

		assertThatObject(t, objectassert.GitRepositoryFromObject(t, gitRepository).
			HasCreatedOnNotEmpty().
			HasName(gitRepositoryId.Name()).
			HasDatabaseName(gitRepositoryId.DatabaseName()).
			HasSchemaName(gitRepositoryId.SchemaName()).
			HasOrigin(origin).
			HasApiIntegration(apiIntegration).
			HasGitCredentials(secretId).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)
	})

	t.Run("alter: set", func(t *testing.T) {
		gitRepositoryId, apiIntegrationId := createGitRepository(t)

		secretId := testClientHelper().Ids.
			RandomSchemaObjectIdentifierInSchema(schema.ID())
		_, secretCleanup := testClientHelper().Secret.
			CreateWithBasicAuthenticationFlow(t, secretId, "username", "password")
		t.Cleanup(secretCleanup)

		setRequest := sdk.NewGitRepositorySetRequest().
			WithGitCredentials(secretId).
			WithComment("comment")
		alterRequest := sdk.NewAlterGitRepositoryRequest(gitRepositoryId).
			WithSet(*setRequest)

		err := client.GitRepositories.Alter(ctx, alterRequest)
		require.NoError(t, err)

		updatedGitRepository, err := client.GitRepositories.ShowByID(ctx, gitRepositoryId)
		require.NoError(t, err)

		assertThatObject(t, objectassert.GitRepositoryFromObject(t, updatedGitRepository).
			HasName(gitRepositoryId.Name()).
			HasDatabaseName(gitRepositoryId.DatabaseName()).
			HasSchemaName(gitRepositoryId.SchemaName()).
			HasOrigin(origin).
			HasApiIntegration(apiIntegrationId).
			HasGitCredentials(secretId).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)
	})

	t.Run("alter: unset", func(t *testing.T) {
		gitRepositoryId, apiIntegrationId, _ := createGitRepositoryWithSecretAndComment(t)

		unsetRequest := sdk.NewGitRepositoryUnsetRequest().
			WithGitCredentials(true).
			WithComment(true)
		alterRequest := sdk.NewAlterGitRepositoryRequest(gitRepositoryId).
			WithUnset(*unsetRequest)

		err := client.GitRepositories.Alter(ctx, alterRequest)
		require.NoError(t, err)

		updated, err := testClientHelper().GitRepository.Show(t, gitRepositoryId)
		require.NoError(t, err)

		assertThatObject(t, objectassert.GitRepositoryFromObject(t, updated).
			HasName(gitRepositoryId.Name()).
			HasDatabaseName(gitRepositoryId.DatabaseName()).
			HasSchemaName(gitRepositoryId.SchemaName()).
			HasOrigin(origin).
			HasApiIntegration(apiIntegrationId).
			HasGitCredentialsEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		)
	})

	t.Run("drop", func(t *testing.T) {
		gitRepositoryId, _, _ := createGitRepositoryWithSecretAndComment(t)

		err := client.GitRepositories.Drop(ctx, sdk.NewDropGitRepositoryRequest(gitRepositoryId).WithIfExists(true))
		require.NoError(t, err)

		_, err = client.GitRepositories.ShowByID(ctx, gitRepositoryId)
		require.Error(t, err)
	})

	t.Run("show: with like", func(t *testing.T) {
		gitRepositoryId, _, _ := createGitRepositoryWithSecretAndComment(t)
		gitRepository, err := testClientHelper().GitRepository.Show(t, gitRepositoryId)
		require.NoError(t, err)

		pattern := gitRepositoryId.Name()
		gitRepositories, err := client.GitRepositories.Show(ctx, sdk.NewShowGitRepositoryRequest().WithLike(sdk.Like{Pattern: &pattern}))
		require.NoError(t, err)
		require.Equal(t, 1, len(gitRepositories))
		require.Equal(t, *gitRepository, gitRepositories[0])
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(otherSchemaCleanup)

		apiIntegration, apiIntegrationCleanup := testClientHelper().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)
		t.Cleanup(apiIntegrationCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), otherSchema.ID())

		_, gitRepositoryCleanup1 := testClientHelper().GitRepository.Create(t, id1, origin, apiIntegration)
		t.Cleanup(gitRepositoryCleanup1)
		_, gitRepositoryCleanup2 := testClientHelper().GitRepository.Create(t, id2, origin, apiIntegration)
		t.Cleanup(gitRepositoryCleanup2)

		e1, err := client.GitRepositories.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.GitRepositories.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("describe", func(t *testing.T) {
		gitRepositoryId, apiIntegrationId, secretId := createGitRepositoryWithSecretAndComment(t)

		gitRepositories, err := client.GitRepositories.Describe(ctx, gitRepositoryId)
		require.NoError(t, err)
		require.Len(t, gitRepositories, 1)
		gitRepository := gitRepositories[0]

		assertThatObject(t, objectassert.GitRepositoryFromObject(t, &gitRepository).
			HasName(gitRepositoryId.Name()).
			HasDatabaseName(gitRepositoryId.DatabaseName()).
			HasSchemaName(gitRepositoryId.SchemaName()).
			HasOrigin(origin).
			HasApiIntegration(apiIntegrationId).
			HasGitCredentials(secretId).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)
	})

	t.Run("show git branches", func(t *testing.T) {
		gitRepositoryId, _, _ := createGitRepositoryWithSecretAndComment(t)

		branches, err := client.GitRepositories.ShowGitBranches(ctx, sdk.NewShowGitBranchesGitRepositoryRequest(gitRepositoryId))
		require.NoError(t, err)
		require.Len(t, branches, 3)

		expectedNames := []string{"master", "octocat-patch-1", "test"}
		var branchNames []string
		for _, b := range branches {
			branchNames = append(branchNames, strings.ToLower(b.Name))
		}
		require.ElementsMatch(t, expectedNames, branchNames)
	})
}
