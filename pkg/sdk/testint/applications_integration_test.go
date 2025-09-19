//go:build account_level_tests

package testint

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * todo: (SNOW-1015095) add integration test for `ALTER APPLICATION <name> UPGRADE`
 * todo: (SNOW-1016268) ALTER APPLICATION [ IF EXISTS ] <name> SET [ SHARE_EVENTS_WITH_PROVIDER ]
 *       attention: SHARE_EVENTS_WITH_PROVIDER can only be set/unset if the application is created in a different account from the application package
 */

func TestInt_Applications(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	tagTest, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	cleanupApplicationHandle := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.Applications.Drop(ctx, sdk.NewDropApplicationRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createRegularApplicationPackage := func(t *testing.T, stageId sdk.SchemaObjectIdentifier, version string) *sdk.ApplicationPackage {
		t.Helper()
		applicationPackage, cleanupApplicationPackage := testClientHelper().ApplicationPackage.CreateApplicationPackage(t)
		t.Cleanup(cleanupApplicationPackage)
		testClientHelper().ApplicationPackage.RegisterVersion(t, applicationPackage.ID(), stageId, version)
		return applicationPackage
	}

	createApplicationPackageWithReleaseChannelsDisabled := func(t *testing.T, stageId sdk.SchemaObjectIdentifier, version string, patch int) *sdk.ApplicationPackage {
		t.Helper()
		applicationPackage, cleanupApplicationPackage := testClientHelper().ApplicationPackage.CreateApplicationPackageWithReleaseChannelsDisabled(t)
		t.Cleanup(cleanupApplicationPackage)
		testClientHelper().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stageId, version)
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`ALTER APPLICATION PACKAGE %s SET DEFAULT RELEASE DIRECTIVE VERSION = %s PATCH = %d`, applicationPackage.ID().FullyQualifiedName(), version, patch))
		require.NoError(t, err)
		return applicationPackage
	}

	createApplicationPackageHandle := func(t *testing.T, version string, patch int, defaultReleaseDirective bool) (*sdk.Stage, *sdk.ApplicationPackage) {
		t.Helper()

		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)

		testClientHelper().Stage.PutOnStageWithContent(t, stage.ID(), "manifest.yml", "")
		testClientHelper().Stage.PutOnStageWithContent(t, stage.ID(), "setup.sql", "CREATE APPLICATION ROLE IF NOT EXISTS APP_HELLO_SNOWFLAKE;")

		var applicationPackage *sdk.ApplicationPackage

		// handle release directive for application package
		if defaultReleaseDirective {
			applicationPackage = createApplicationPackageWithReleaseChannelsDisabled(t, stage.ID(), version, patch)
		} else {
			applicationPackage = createRegularApplicationPackage(t, stage.ID(), version)
		}
		return stage, applicationPackage
	}

	createApplicationHandle := func(t *testing.T, version string, patch int, versionDirectory bool, debug bool, addPatch bool) (*sdk.Stage, *sdk.Application, *sdk.ApplicationPackage) {
		t.Helper()

		stage, applicationPackage := createApplicationPackageHandle(t, version, patch, false)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		vr := sdk.NewApplicationVersionRequest().WithVersionAndPatch(sdk.NewVersionAndPatchRequest(version, &patch))
		if versionDirectory {
			vr = sdk.NewApplicationVersionRequest().WithVersionDirectory(sdk.String("@" + stage.ID().FullyQualifiedName()))
		}
		request := sdk.NewCreateApplicationRequest(id, applicationPackage.ID()).WithVersion(vr).WithDebugMode(&debug)
		err := client.Applications.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationHandle(id))

		if addPatch {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`ALTER APPLICATION PACKAGE %s ADD PATCH FOR VERSION %s USING '@%s'`, applicationPackage.ID().FullyQualifiedName(), version, stage.ID().FullyQualifiedName()))
			require.NoError(t, err)
		}

		application, err := client.Applications.ShowByID(ctx, id)
		require.NoError(t, err)
		return stage, application, applicationPackage
	}

	assertApplication := func(t *testing.T, id sdk.AccountObjectIdentifier, applicationPackageName, version string, patch int, comment string) {
		t.Helper()

		param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterDataRetentionTimeInDays)
		require.NoError(t, err)

		defaultDataRetentionTimeInDays, err := strconv.Atoi(param.Value)
		require.NoError(t, err)

		e, err := client.Applications.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.NotEmpty(t, e.CreatedOn)
		assert.Equal(t, id.Name(), e.Name)
		assert.Equal(t, false, e.IsDefault)
		assert.Equal(t, true, e.IsCurrent)
		assert.Equal(t, "APPLICATION PACKAGE", e.SourceType)
		assert.Equal(t, applicationPackageName, e.Source)
		assert.Equal(t, version, e.Version)
		assert.Equal(t, patch, e.Patch)
		assert.Equal(t, "ACCOUNTADMIN", e.Owner)
		assert.Equal(t, comment, e.Comment)
		assert.Equal(t, defaultDataRetentionTimeInDays, e.RetentionTime)
		assert.Empty(t, e.Options)
	}

	t.Run("create application: without version", func(t *testing.T) {
		version, patch := "V001", 0
		_, applicationPackage := createApplicationPackageHandle(t, version, patch, true)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		request := sdk.NewCreateApplicationRequest(id, applicationPackage.ID()).
			WithComment(&comment).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
			})
		err := client.Applications.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationHandle(id))

		e, err := client.Applications.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, "ACCOUNTADMIN", e.Owner)
		require.Equal(t, comment, e.Comment)
		require.Equal(t, "APPLICATION PACKAGE", e.SourceType)
		require.Equal(t, applicationPackage.Name, e.Source)
		require.Equal(t, version, e.Version)
		require.Equal(t, patch, e.Patch)
	})

	t.Run("create application: version and patch", func(t *testing.T) {
		version, patch := "V001", 0
		_, applicationPackage := createApplicationPackageHandle(t, version, patch, false)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		vr := sdk.NewApplicationVersionRequest().WithVersionAndPatch(sdk.NewVersionAndPatchRequest(version, &patch))
		comment := random.Comment()
		request := sdk.NewCreateApplicationRequest(id, applicationPackage.ID()).
			WithDebugMode(sdk.Bool(true)).
			WithComment(&comment).
			WithVersion(vr)
		err := client.Applications.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationHandle(id))

		assertApplication(t, id, applicationPackage.Name, version, patch, comment)
	})

	t.Run("create application: version directory", func(t *testing.T) {
		version, patch := "V001", 0
		stage, applicationPackage := createApplicationPackageHandle(t, version, patch, false)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		vr := sdk.NewApplicationVersionRequest().WithVersionDirectory(sdk.String("@" + stage.ID().FullyQualifiedName()))
		comment := random.Comment()
		request := sdk.NewCreateApplicationRequest(id, applicationPackage.ID()).
			WithDebugMode(sdk.Bool(true)).
			WithComment(&comment).
			WithVersion(vr).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
			})
		err := client.Applications.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationHandle(id))

		assertApplication(t, id, applicationPackage.Name, "UNVERSIONED", patch, comment)
	})

	t.Run("show application: with like", func(t *testing.T) {
		version, patch := "V001", 0
		_, e, _ := createApplicationHandle(t, version, patch, false, true, false)
		packages, err := client.Applications.Show(ctx, sdk.NewShowApplicationRequest().WithLike(sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(packages))
		require.Equal(t, *e, packages[0])
	})

	t.Run("alter application: set", func(t *testing.T) {
		version, patch := "V001", 0
		_, e, _ := createApplicationHandle(t, version, patch, false, true, false)
		id := e.ID()

		comment, mode := random.Comment(), true
		set := sdk.NewApplicationSetRequest().
			WithComment(&comment).
			WithDebugMode(&mode)
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithSet(set))
		require.NoError(t, err)

		details, err := client.Applications.Describe(ctx, id)
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		require.Equal(t, e.SourceType, pairs["source_type"])
		require.Equal(t, e.Source, pairs["source"])
		require.Equal(t, e.Version, pairs["version"])
		require.Equal(t, strconv.Itoa(e.Patch), pairs["patch"])
		require.Equal(t, comment, pairs["comment"])
		require.Equal(t, strconv.FormatBool(mode), pairs["debug_mode"])
	})

	t.Run("alter application: unset", func(t *testing.T) {
		version, patch := "V001", 0
		_, e, _ := createApplicationHandle(t, version, patch, false, true, false)
		id := e.ID()

		unset := sdk.NewApplicationUnsetRequest().WithComment(sdk.Bool(true)).WithDebugMode(sdk.Bool(true))
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUnset(unset))
		require.NoError(t, err)

		o, err := client.Applications.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Empty(t, o.Comment)

		details, err := client.Applications.Describe(ctx, id)
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		require.Equal(t, strconv.FormatBool(false), pairs["debug_mode"])
	})

	t.Run("alter application: upgrade with version and patch", func(t *testing.T) {
		version, patch := "V001", 0
		_, e, applicationPackage := createApplicationHandle(t, version, patch, false, true, true)
		id := e.ID()

		vr := sdk.NewVersionAndPatchRequest(version, sdk.Int(patch+1))
		av := sdk.NewApplicationVersionRequest().WithVersionAndPatch(vr)
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUpgradeVersion(av))
		require.NoError(t, err)
		assertApplication(t, id, applicationPackage.Name, version, patch+1, "")
	})

	t.Run("alter application: upgrade with version directory", func(t *testing.T) {
		version, patch := "V001", 0
		s, e, applicationPackage := createApplicationHandle(t, version, patch, true, true, false)
		id := e.ID()

		av := sdk.NewApplicationVersionRequest().WithVersionDirectory(sdk.String("@" + s.ID().FullyQualifiedName()))
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUpgradeVersion(av))
		require.NoError(t, err)
		assertApplication(t, id, applicationPackage.Name, "UNVERSIONED", patch+1, "")
	})

	t.Run("alter application: unset references", func(t *testing.T) {
		version, patch := "V001", 0
		_, e, applicationPackage := createApplicationHandle(t, version, patch, false, true, false)
		id := e.ID()

		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUnsetReferences(sdk.NewApplicationReferencesRequest()))
		require.NoError(t, err)
		assertApplication(t, id, applicationPackage.Name, version, patch, "")
	})

	t.Run("describe application", func(t *testing.T) {
		version, patch := "V001", 0
		_, e, _ := createApplicationHandle(t, version, patch, false, true, false)
		id := e.ID()

		properties, err := client.Applications.Describe(ctx, id)
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, item := range properties {
			pairs[item.Property] = item.Value
		}
		require.Equal(t, e.SourceType, pairs["source_type"])
		require.Equal(t, e.Source, pairs["source"])
		require.Equal(t, e.Version, pairs["version"])
		require.Equal(t, e.Label, pairs["version_label"])
		require.Equal(t, e.Comment, pairs["comment"])
		require.Equal(t, strconv.Itoa(e.Patch), pairs["patch"])
	})
}
