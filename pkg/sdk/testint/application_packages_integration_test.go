//go:build non_account_level_tests

package testint

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
todo: add tests for:
  - Creates a custom release directive for the specified accounts : https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-release-directive
  - Create application package with insufficient privileges for the following three fields
	-  WithDataRetentionTimeInDays(sdk.Int(1)).
	-  WithMaxDataExtensionTimeInDays(sdk.Int(1)).
	-  WithDefaultDdlCollation(sdk.String("en_US")).
*/

func TestInt_ApplicationPackages(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	tagTest, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	cleanupApplicationPackageHandle := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.ApplicationPackages.Drop(ctx, sdk.NewDropApplicationPackageRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createApplicationPackageHandle := func(t *testing.T) *sdk.ApplicationPackage {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreateApplicationPackageRequest(id).WithDistribution(sdk.DistributionInternal)
		err := client.ApplicationPackages.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationPackageHandle(id))

		e, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	t.Run("create application package", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		request := sdk.NewCreateApplicationPackageRequest(id).
			WithComment(comment).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
			}).
			WithDistribution(sdk.DistributionExternal)
		err := client.ApplicationPackages.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationPackageHandle(id))

		e, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, sdk.DistributionExternal, sdk.Distribution(e.Distribution))
		require.Equal(t, "ACCOUNTADMIN", e.Owner)
		require.Equal(t, comment, e.Comment)

		param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterDataRetentionTimeInDays)
		require.NoError(t, err)

		defaultDataRetentionTimeInDays, err := strconv.Atoi(param.Value)
		require.NoError(t, err)

		require.Equal(t, defaultDataRetentionTimeInDays, e.RetentionTime)
	})

	t.Run("alter application package: set", func(t *testing.T) {
		e := createApplicationPackageHandle(t)
		id := e.ID()

		distribution := sdk.DistributionExternal
		set := sdk.NewApplicationPackageSetRequest().
			WithDistribution(distribution).
			WithComment("test").
			WithDataRetentionTimeInDays(2).
			WithMaxDataExtensionTimeInDays(2).
			WithDefaultDdlCollation("utf8mb4_0900_ai_ci")
		err := client.ApplicationPackages.Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithSet(*set))
		require.NoError(t, err)

		o, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, distribution, sdk.Distribution(o.Distribution))
		assert.Equal(t, 2, o.RetentionTime)
		assert.Equal(t, "test", o.Comment)
	})

	t.Run("alter application package: unset", func(t *testing.T) {
		e := createApplicationPackageHandle(t)
		id := e.ID()

		// unset comment and distribution
		unset := sdk.NewApplicationPackageUnsetRequest().WithComment(true).WithDistribution(true)
		err := client.ApplicationPackages.Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithUnset(*unset))
		require.NoError(t, err)
		o, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Empty(t, o.Comment)
		require.Equal(t, sdk.DistributionInternal, sdk.Distribution(o.Distribution))
	})

	t.Run("show application package for SQL: with like", func(t *testing.T) {
		e := createApplicationPackageHandle(t)

		packages, err := client.ApplicationPackages.Show(ctx, sdk.NewShowApplicationPackageRequest().WithLike(sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Len(t, packages, 1)
		require.Equal(t, *e, packages[0])
	})
}

func TestInt_ApplicationPackagesVersionAndReleaseDirective(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("alter application package: add, patch and drop version", func(t *testing.T) {
		applicationPackage, applicationPackageCleanup := testClientHelper().ApplicationPackage.CreateApplicationPackageWithReleaseChannelsDisabled(t)
		t.Cleanup(applicationPackageCleanup)

		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)
		testClientHelper().Stage.PutOnStage(t, stage.ID(), "manifest.yml")
		testClientHelper().Stage.PutOnStage(t, stage.ID(), "setup.sql")

		version := "V001"
		using := "@" + stage.ID().FullyQualifiedName()
		// add version to application package
		id := applicationPackage.ID()
		vr := sdk.NewAddVersionRequest(using).WithVersionIdentifier(&version).WithLabel("add version V001")
		r1 := sdk.NewAlterApplicationPackageRequest(id).WithAddVersion(*vr)
		err := client.ApplicationPackages.Alter(ctx, r1)
		require.NoError(t, err)
		versions := testClientHelper().ApplicationPackage.ShowVersions(t, applicationPackage.ID())
		require.Len(t, versions, 1)
		require.Equal(t, version, versions[0].Version)
		require.Equal(t, 0, versions[0].Patch)

		// add patch for application package version
		pr := sdk.NewAddPatchForVersionRequest(&version, using).WithLabel("patch version V001")
		r2 := sdk.NewAlterApplicationPackageRequest(id).WithAddPatchForVersion(*pr)
		err = client.ApplicationPackages.Alter(ctx, r2)
		require.NoError(t, err)
		versions = testClientHelper().ApplicationPackage.ShowVersions(t, applicationPackage.ID())
		require.Len(t, versions, 2)
		require.Equal(t, version, versions[0].Version)
		require.Equal(t, 0, versions[0].Patch)
		require.Equal(t, version, versions[1].Version)
		require.Equal(t, 1, versions[1].Patch)

		// drop version from application package
		r3 := sdk.NewAlterApplicationPackageRequest(id).WithDropVersion(*sdk.NewDropVersionRequest(version))
		err = client.ApplicationPackages.Alter(ctx, r3)
		require.NoError(t, err)
		versions = testClientHelper().ApplicationPackage.ShowVersions(t, applicationPackage.ID())
		require.Empty(t, versions)
	})

	t.Run("alter application package: set default release directive", func(t *testing.T) {
		applicationPackage, applicationPackageCleanup := testClientHelper().ApplicationPackage.CreateApplicationPackageWithReleaseChannelsDisabled(t)
		t.Cleanup(applicationPackageCleanup)

		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)
		testClientHelper().Stage.PutOnStage(t, stage.ID(), "manifest.yml")
		testClientHelper().Stage.PutOnStage(t, stage.ID(), "setup.sql")

		version := "V001"
		using := "@" + stage.ID().FullyQualifiedName()
		// add version to application package
		id := applicationPackage.ID()
		vr := sdk.NewAddVersionRequest(using).WithVersionIdentifier(&version).WithLabel("add version V001")
		r1 := sdk.NewAlterApplicationPackageRequest(id).WithAddVersion(*vr)
		err := client.ApplicationPackages.Alter(ctx, r1)
		require.NoError(t, err)
		versions := testClientHelper().ApplicationPackage.ShowVersions(t, applicationPackage.ID())
		require.Len(t, versions, 1)
		require.Equal(t, version, versions[0].Version)
		require.Equal(t, 0, versions[0].Patch)

		// set default release directive
		rr := sdk.NewSetDefaultReleaseDirectiveRequest(version, 0)
		r2 := sdk.NewAlterApplicationPackageRequest(id).WithSetDefaultReleaseDirective(*rr)
		err = client.ApplicationPackages.Alter(ctx, r2)
		require.NoError(t, err)
	})
}
