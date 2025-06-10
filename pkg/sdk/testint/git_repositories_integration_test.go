//go:build !account_level_tests

package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_GitRepositories(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	gitRepositoryOrigin := testvars.ExampleGitRepositoryOrigin

	apiIntegrationId, apiIntegrationCleanup := testClientHelper().ApiIntegration.
		CreateApiIntegrationForGitRepository(t, gitRepositoryOrigin)
	t.Cleanup(apiIntegrationCleanup)

	secretId, secretCleanup := testClientHelper().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(secretCleanup)

	t.Run("create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		request := sdk.NewCreateGitRepositoryRequest(id, gitRepositoryOrigin, apiIntegrationId)

		err := client.GitRepositories.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().GitRepository.DropFunc(t, id))

		gitRepository, err := client.GitRepositories.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.GitRepositoryFromObject(t, gitRepository).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOrigin(gitRepositoryOrigin).
			HasApiIntegration(apiIntegrationId).
			HasGitCredentialsEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		)
	})

	t.Run("create - complete", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		request := sdk.NewCreateGitRepositoryRequest(gitRepositoryId, gitRepositoryOrigin, apiIntegrationId).WithIfNotExists(true).WithGitCredentials(secretId).WithComment("comment")

		err := client.GitRepositories.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().GitRepository.DropFunc(t, gitRepositoryId))

		gitRepository, err := client.GitRepositories.ShowByID(ctx, gitRepositoryId)
		require.NoError(t, err)

		assertThatObject(t, objectassert.GitRepositoryFromObject(t, gitRepository).
			HasCreatedOnNotEmpty().
			HasName(gitRepositoryId.Name()).
			HasDatabaseName(gitRepositoryId.DatabaseName()).
			HasSchemaName(gitRepositoryId.SchemaName()).
			HasOrigin(gitRepositoryOrigin).
			HasApiIntegration(apiIntegrationId).
			HasGitCredentials(secretId).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)
	})

	t.Run("alter: set", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.
			Create(t, gitRepositoryId, gitRepositoryOrigin, apiIntegrationId)
		t.Cleanup(gitRepositoryCleanup)

		newApiIntegrationId, newApiIntegrationCleanup := testClientHelper().ApiIntegration.
			CreateApiIntegrationForGitRepository(t, gitRepositoryOrigin)
		t.Cleanup(newApiIntegrationCleanup)

		setRequest := sdk.NewGitRepositorySetRequest().
			WithApiIntegration(newApiIntegrationId).
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
			HasOrigin(gitRepositoryOrigin).
			HasApiIntegration(newApiIntegrationId).
			HasGitCredentials(secretId).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)
	})

	t.Run("alter: unset", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createRequest := sdk.NewCreateGitRepositoryRequest(gitRepositoryId, gitRepositoryOrigin, apiIntegrationId).WithGitCredentials(secretId).WithComment("comment")
		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.CreateWithRequest(t, createRequest)
		t.Cleanup(gitRepositoryCleanup)

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
			HasOrigin(gitRepositoryOrigin).
			HasApiIntegration(apiIntegrationId).
			HasGitCredentialsEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		)
	})

	t.Run("drop", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.
			Create(t, gitRepositoryId, gitRepositoryOrigin, apiIntegrationId)
		t.Cleanup(gitRepositoryCleanup)

		err := client.GitRepositories.Drop(ctx, sdk.NewDropGitRepositoryRequest(gitRepositoryId).WithIfExists(true))
		require.NoError(t, err)

		_, err = client.GitRepositories.ShowByID(ctx, gitRepositoryId)
		require.Error(t, err)
	})

	t.Run("show: with like", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createRequest := sdk.NewCreateGitRepositoryRequest(gitRepositoryId, gitRepositoryOrigin, apiIntegrationId).WithGitCredentials(secretId).WithComment("comment")
		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.CreateWithRequest(t, createRequest)
		t.Cleanup(gitRepositoryCleanup)

		gitRepository, err := testClientHelper().GitRepository.Show(t, gitRepositoryId)
		require.NoError(t, err)

		pattern := gitRepositoryId.Name()
		gitRepositories, err := client.GitRepositories.Show(ctx, sdk.NewShowGitRepositoryRequest().WithLike(sdk.Like{Pattern: &pattern}))
		require.NoError(t, err)
		require.Equal(t, 1, len(gitRepositories))
		require.Equal(t, *gitRepository, gitRepositories[0])
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(otherSchemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), otherSchema.ID())

		_, gitRepositoryCleanup1 := testClientHelper().GitRepository.Create(t, id1, gitRepositoryOrigin, apiIntegrationId)
		t.Cleanup(gitRepositoryCleanup1)
		_, gitRepositoryCleanup2 := testClientHelper().GitRepository.Create(t, id2, gitRepositoryOrigin, apiIntegrationId)
		t.Cleanup(gitRepositoryCleanup2)

		e1, err := client.GitRepositories.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.GitRepositories.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("show git tags", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.
			Create(t, gitRepositoryId, gitRepositoryOrigin, apiIntegrationId)
		t.Cleanup(gitRepositoryCleanup)

		tags, err := client.GitRepositories.ShowGitTags(ctx, sdk.NewShowGitTagsGitRepositoryRequest(gitRepositoryId))
		require.NoError(t, err)
		require.Zero(t, len(tags))
	})

	t.Run("describe", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createRequest := sdk.NewCreateGitRepositoryRequest(gitRepositoryId, gitRepositoryOrigin, apiIntegrationId).WithGitCredentials(secretId).WithComment("comment")
		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.CreateWithRequest(t, createRequest)
		t.Cleanup(gitRepositoryCleanup)

		gitRepository, err := client.GitRepositories.Describe(ctx, gitRepositoryId)
		require.NoError(t, err)

		assertThatObject(t, objectassert.GitRepositoryFromObject(t, gitRepository).
			HasName(gitRepositoryId.Name()).
			HasDatabaseName(gitRepositoryId.DatabaseName()).
			HasSchemaName(gitRepositoryId.SchemaName()).
			HasOrigin(gitRepositoryOrigin).
			HasApiIntegration(apiIntegrationId).
			HasGitCredentials(secretId).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)
	})

	t.Run("show git branches", func(t *testing.T) {
		gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, gitRepositoryCleanup := testClientHelper().
			GitRepository.
			Create(t, gitRepositoryId, gitRepositoryOrigin, apiIntegrationId)
		t.Cleanup(gitRepositoryCleanup)

		branches, err := client.GitRepositories.ShowGitBranches(ctx, sdk.NewShowGitBranchesGitRepositoryRequest(gitRepositoryId))
		require.NoError(t, err)
		require.NotZero(t, len(branches))

		var branchNames []string
		for _, b := range branches {
			branchNames = append(branchNames, strings.ToLower(b.Name))
		}
		require.Contains(t, branchNames, "master")
	})
}
