//go:build !account_level_tests

package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_ImageRepositories(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	t.Run("create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateImageRepositoryRequest(id)

		err := client.ImageRepositories.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ImageRepository.DropImageRepositoryFunc(t, id))

		imageRepository, err := client.ImageRepositories.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ImageRepositoryFromObject(t, imageRepository).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasRepositoryUrlNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("").
			HasPrivatelinkRepositoryUrl(""),
		)
	})

	t.Run("create - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		comment := random.Comment()
		request := sdk.NewCreateImageRepositoryRequest(id).WithIfNotExists(true).WithComment(comment)

		err := client.ImageRepositories.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ImageRepository.DropImageRepositoryFunc(t, id))

		imageRepository, err := client.ImageRepositories.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ImageRepositoryFromObject(t, imageRepository).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasRepositoryUrlNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment).
			HasPrivatelinkRepositoryUrl(""),
		)
	})

	// TODO(SNOW-2070746): Using symbols that require quoting the name fails - in this case, lowercase letters. Remove this after fix on snowflake side.
	t.Run("create with an ID with lowercase letters fails", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateImageRepositoryRequest(id)

		err := client.ImageRepositories.Create(ctx, req)
		require.ErrorContains(t, err, "db, schema and repo names for an image repo must be unquoted identifiers")
	})

	t.Run("alter: set", func(t *testing.T) {
		imageRepository, funcCleanup := testClientHelper().ImageRepository.CreateInSchema(t, schema.ID())
		t.Cleanup(funcCleanup)

		comment := random.Comment()

		req := sdk.NewAlterImageRepositoryRequest(imageRepository.ID()).WithSet(
			*sdk.NewImageRepositorySetRequest().
				WithComment(sdk.StringAllowEmpty{Value: comment}),
		)
		err := client.ImageRepositories.Alter(ctx, req)
		require.NoError(t, err)

		imageRepository, err = client.ImageRepositories.ShowByID(ctx, imageRepository.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ImageRepositoryFromObject(t, imageRepository).
			HasComment(comment),
		)

		// Set comment to an empty string.
		// TODO(SNOW-2070753): After UNSET COMMENT is fixed in Snowflake, remove this check and fallback to UNSET COMMENT.
		req = sdk.NewAlterImageRepositoryRequest(imageRepository.ID()).WithSet(
			*sdk.NewImageRepositorySetRequest().
				WithComment(sdk.StringAllowEmpty{Value: ""}),
		)
		err = client.ImageRepositories.Alter(ctx, req)
		require.NoError(t, err)

		imageRepository, err = client.ImageRepositories.ShowByID(ctx, imageRepository.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ImageRepositoryFromObject(t, imageRepository).
			HasComment(""),
		)
	})

	// TODO(SNOW-2070753): Adjust this test after it's fixed in Snowflake.
	t.Run("alter: unset comment fails", func(t *testing.T) {
		imageRepository, funcCleanup := testClientHelper().ImageRepository.CreateInSchema(t, schema.ID())
		t.Cleanup(funcCleanup)

		_, err := client.ExecForTests(context.Background(), fmt.Sprintf("ALTER IMAGE REPOSITORY %s UNSET COMMENT", imageRepository.ID().FullyQualifiedName()))
		require.ErrorContains(t, err, "000002 (0A000): Unsupported feature 'UNSET'")
	})

	t.Run("show: with like", func(t *testing.T) {
		imageRepository, funcCleanup := testClientHelper().ImageRepository.CreateInSchema(t, schema.ID())
		t.Cleanup(funcCleanup)

		imageRepositories, err := client.ImageRepositories.Show(ctx, sdk.NewShowImageRepositoryRequest().WithLike(sdk.Like{Pattern: &imageRepository.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(imageRepositories))
		require.Equal(t, *imageRepository, imageRepositories[0])
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(otherSchemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), otherSchema.ID())

		req1 := sdk.NewCreateImageRepositoryRequest(id1)
		req2 := sdk.NewCreateImageRepositoryRequest(id2)
		_, cleanup1 := testClientHelper().ImageRepository.CreateWithRequest(t, req1)
		t.Cleanup(cleanup1)
		_, cleanup2 := testClientHelper().ImageRepository.CreateWithRequest(t, req2)
		t.Cleanup(cleanup2)

		e1, err := client.ImageRepositories.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.ImageRepositories.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
