//go:build non_account_level_tests

package testint

import (
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_StorageIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	assertStorageIntegrationShowResult := func(t *testing.T, s *sdk.StorageIntegration, id sdk.AccountObjectIdentifier, comment string) {
		t.Helper()

		assertThatObject(t, objectassert.StorageIntegrationFromObject(t, s).
			HasName(id.Name()).
			HasEnabled(true).
			HasStorageTypeExternal().
			HasCategoryStorage().
			HasComment(comment),
		)
	}

	findProp := func(t *testing.T, props []sdk.StorageIntegrationProperty, name string) *sdk.StorageIntegrationProperty {
		t.Helper()
		prop, err := collections.FindFirst(props, func(property sdk.StorageIntegrationProperty) bool { return property.Name == name })
		require.NoError(t, err)
		return prop
	}

	// TODO [next PR]: replace with fluent assertions like ContainsDetail for semantic views
	assertS3StorageIntegrationDescResult := func(
		t *testing.T,
		props []sdk.StorageIntegrationProperty,
		enabled bool,
		allowedLocations []sdk.StorageLocation,
		blockedLocations []sdk.StorageLocation,
		comment string,
		usePrivateLinkEndpoint bool,
	) {
		t.Helper()
		allowed := make([]string, len(allowedLocations))
		for i, a := range allowedLocations {
			allowed[i] = a.Path
		}
		blocked := make([]string, len(blockedLocations))
		for i, b := range blockedLocations {
			blocked[i] = b.Path
		}
		assert.Equal(t, "Boolean", findProp(t, props, "ENABLED").Type)
		assert.Equal(t, strconv.FormatBool(enabled), findProp(t, props, "ENABLED").Value)
		assert.Equal(t, "false", findProp(t, props, "ENABLED").Default)
		assert.Equal(t, "S3", findProp(t, props, "STORAGE_PROVIDER").Value)
		assert.Equal(t, strings.Join(allowed, ","), findProp(t, props, "STORAGE_ALLOWED_LOCATIONS").Value)
		assert.Equal(t, strings.Join(blocked, ","), findProp(t, props, "STORAGE_BLOCKED_LOCATIONS").Value)
		assert.NotEmpty(t, findProp(t, props, "STORAGE_AWS_IAM_USER_ARN").Value)
		assert.NotEmpty(t, findProp(t, props, "STORAGE_AWS_ROLE_ARN").Value)
		assert.NotEmpty(t, findProp(t, props, "STORAGE_AWS_EXTERNAL_ID").Value)
		assert.Equal(t, comment, findProp(t, props, "COMMENT").Value)
		assert.Equal(t, strconv.FormatBool(usePrivateLinkEndpoint), findProp(t, props, "USE_PRIVATELINK_ENDPOINT").Value)
	}

	// TODO [next PR]: replace with fluent assertions like ContainsDetail for semantic views
	assertGCSStorageIntegrationDescResult := func(
		t *testing.T,
		props []sdk.StorageIntegrationProperty,
		enabled bool,
		allowedLocations []sdk.StorageLocation,
		blockedLocations []sdk.StorageLocation,
		comment string,
	) {
		t.Helper()
		allowed := make([]string, len(allowedLocations))
		for i, a := range allowedLocations {
			allowed[i] = a.Path
		}
		blocked := make([]string, len(blockedLocations))
		for i, b := range blockedLocations {
			blocked[i] = b.Path
		}
		assert.Equal(t, "Boolean", findProp(t, props, "ENABLED").Type)
		assert.Equal(t, strconv.FormatBool(enabled), findProp(t, props, "ENABLED").Value)
		assert.Equal(t, "false", findProp(t, props, "ENABLED").Default)
		assert.Equal(t, "GCS", findProp(t, props, "STORAGE_PROVIDER").Value)
		assert.Equal(t, strings.Join(allowed, ","), findProp(t, props, "STORAGE_ALLOWED_LOCATIONS").Value)
		assert.Equal(t, strings.Join(blocked, ","), findProp(t, props, "STORAGE_BLOCKED_LOCATIONS").Value)
		assert.NotEmpty(t, findProp(t, props, "STORAGE_GCP_SERVICE_ACCOUNT").Value)
		assert.Equal(t, comment, findProp(t, props, "COMMENT").Value)
	}

	// TODO [next PR]: replace with fluent assertions like ContainsDetail for semantic views
	assertAzureStorageIntegrationDescResult := func(
		t *testing.T,
		props []sdk.StorageIntegrationProperty,
		enabled bool,
		allowedLocations []sdk.StorageLocation,
		blockedLocations []sdk.StorageLocation,
		comment string,
	) {
		t.Helper()
		allowed := make([]string, len(allowedLocations))
		for i, a := range allowedLocations {
			allowed[i] = a.Path
		}
		blocked := make([]string, len(blockedLocations))
		for i, b := range blockedLocations {
			blocked[i] = b.Path
		}
		assert.Equal(t, "Boolean", findProp(t, props, "ENABLED").Type)
		assert.Equal(t, strconv.FormatBool(enabled), findProp(t, props, "ENABLED").Value)
		assert.Equal(t, "false", findProp(t, props, "ENABLED").Default)
		assert.Equal(t, "AZURE", findProp(t, props, "STORAGE_PROVIDER").Value)
		assert.Equal(t, strings.Join(allowed, ","), findProp(t, props, "STORAGE_ALLOWED_LOCATIONS").Value)
		assert.Equal(t, strings.Join(blocked, ","), findProp(t, props, "STORAGE_BLOCKED_LOCATIONS").Value)
		assert.NotEmpty(t, findProp(t, props, "AZURE_TENANT_ID").Value)
		assert.NotEmpty(t, findProp(t, props, "AZURE_CONSENT_URL").Value)
		assert.NotEmpty(t, findProp(t, props, "AZURE_MULTI_TENANT_APP_NAME").Value)
		assert.Equal(t, comment, findProp(t, props, "COMMENT").Value)
	}

	allowedLocations := func(prefix string) []sdk.StorageLocation {
		return []sdk.StorageLocation{
			{
				Path: prefix + "/allowed-location",
			},
			{
				Path: prefix + "/allowed-location2",
			},
		}
	}
	s3AllowedLocations := allowedLocations(awsBucketUrl)
	gcsAllowedLocations := allowedLocations(gcsBucketUrl)
	azureAllowedLocations := allowedLocations(azureBucketUrl)

	blockedLocations := func(prefix string) []sdk.StorageLocation {
		return []sdk.StorageLocation{
			{
				Path: prefix + "/blocked-location",
			},
			{
				Path: prefix + "/blocked-location2",
			},
		}
	}
	s3BlockedLocations := blockedLocations(awsBucketUrl)
	gcsBlockedLocations := blockedLocations(gcsBucketUrl)
	azureBlockedLocations := blockedLocations(azureBucketUrl)

	createS3StorageIntegrationBasic := func(t *testing.T, protocol sdk.S3Protocol) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, s3AllowedLocations).
			WithS3StorageProviderParams(*sdk.NewS3StorageParamsRequest(protocol, awsRoleARN))

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(testClientHelper().StorageIntegration.DropFunc(t, id))
		return id
	}

	createS3StorageIntegration := func(t *testing.T, protocol sdk.S3Protocol) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, s3AllowedLocations).
			WithIfNotExists(true).
			WithS3StorageProviderParams(*sdk.NewS3StorageParamsRequest(protocol, awsRoleARN).
				WithStorageAwsExternalId("some-external-id").
				WithStorageAwsObjectAcl("bucket-owner-full-control").
				WithUsePrivatelinkEndpoint(true),
			).
			WithStorageBlockedLocations(s3BlockedLocations).
			WithComment("some comment")

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(testClientHelper().StorageIntegration.DropFunc(t, id))
		return id
	}

	createGCSStorageIntegrationBasic := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, gcsAllowedLocations).
			WithGCSStorageProviderParams(*sdk.NewGCSStorageParamsRequest())

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(testClientHelper().StorageIntegration.DropFunc(t, id))
		return id
	}

	createGCSStorageIntegration := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, gcsAllowedLocations).
			WithIfNotExists(true).
			WithGCSStorageProviderParams(*sdk.NewGCSStorageParamsRequest()).
			WithStorageBlockedLocations(gcsBlockedLocations).
			WithComment("some comment")

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(testClientHelper().StorageIntegration.DropFunc(t, id))
		return id
	}

	createAzureStorageIntegrationBasic := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, azureAllowedLocations).
			WithAzureStorageProviderParams(*sdk.NewAzureStorageParamsRequest(azureTenantId))

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(testClientHelper().StorageIntegration.DropFunc(t, id))
		return id
	}

	// TODO [SNOW-2356128]: Add test for use_privatelink_endpoint on Azure deployment (preprod?)
	createAzureStorageIntegration := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, azureAllowedLocations).
			WithIfNotExists(true).
			WithAzureStorageProviderParams(*sdk.NewAzureStorageParamsRequest(azureTenantId)).
			WithStorageBlockedLocations(azureBlockedLocations).
			WithComment("some comment")

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(testClientHelper().StorageIntegration.DropFunc(t, id))
		return id
	}

	// Enabled is required even though it can be UNSET. Not using it in create results in:
	// 002029 (42601): SQL compilation error: Missing option(s): ENABLED
	t.Run("create: without enabled, even though it can be unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := testClientHelper().StorageIntegration.CreateWithoutEnabled(t, id, awsRoleARN, s3AllowedLocations[0])

		require.Error(t, err)
		require.ErrorContains(t, err, "Missing option(s): ENABLED")
	})

	t.Run("create: s3 basic", func(t *testing.T) {
		id := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		// TODO [next PR]: assert all properties from describe
		assert.NotEmpty(t, findProp(t, props, "STORAGE_AWS_EXTERNAL_ID").Value)
	})

	t.Run("create: s3", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		// TODO [next PR]: assert all properties from describe
		assert.Equal(t, "some-external-id", findProp(t, props, "STORAGE_AWS_EXTERNAL_ID").Value)
	})

	// TODO [SNOW-1820099]: Run integration tests on gov preprod
	t.Run("create: s3gov", func(t *testing.T) {
		t.Skip("TODO [SNOW-1820099]: Run integration tests on gov preprod")
		id := createS3StorageIntegration(t, sdk.GovS3Protocol)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")

		// TODO [SNOW-1820099]: assert describe
	})

	t.Run("create: gcs basic", func(t *testing.T) {
		id := createGCSStorageIntegrationBasic(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "GCS", findProp(t, props, "STORAGE_PROVIDER").Value)
		// TODO [next PR]: assert all properties from describe
	})

	t.Run("create: gcs", func(t *testing.T) {
		id := createGCSStorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "GCS", findProp(t, props, "STORAGE_PROVIDER").Value)
		// TODO [next PR]: assert all properties from describe
	})

	t.Run("create: azure basic", func(t *testing.T) {
		id := createAzureStorageIntegrationBasic(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "AZURE", findProp(t, props, "STORAGE_PROVIDER").Value)
		// TODO [next PR]: assert all properties from describe
	})

	t.Run("create: azure", func(t *testing.T) {
		id := createAzureStorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "AZURE", findProp(t, props, "STORAGE_PROVIDER").Value)
		// TODO [next PR]: assert all properties from describe
	})

	t.Run("alter: s3, set and unset", func(t *testing.T) {
		id := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)

		changedS3AllowedLocations := append([]sdk.StorageLocation{{Path: awsBucketUrl + "/allowed-location3"}}, s3AllowedLocations...)
		changedS3BlockedLocations := append([]sdk.StorageLocation{{Path: awsBucketUrl + "/blocked-location3"}}, s3BlockedLocations...)
		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithSet(
				*sdk.NewStorageIntegrationSetRequest().
					WithS3Params(
						*sdk.NewSetS3StorageParamsRequest().
							WithStorageAwsRoleArn(awsRoleARN).
							WithStorageAwsObjectAcl("bucket-owner-full-control").
							WithStorageAwsExternalId("new-external-id").
							WithUsePrivatelinkEndpoint(true),
					).
					WithEnabled(true).
					WithStorageAllowedLocations(changedS3AllowedLocations).
					WithStorageBlockedLocations(changedS3BlockedLocations).
					WithComment("changed comment"),
			)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertS3StorageIntegrationDescResult(t, props, true, changedS3AllowedLocations, changedS3BlockedLocations, "changed comment", true)
		assert.Equal(t, "new-external-id", findProp(t, props, "STORAGE_AWS_EXTERNAL_ID").Value)

		unset := sdk.NewAlterStorageIntegrationRequest(id).
			WithUnset(
				*sdk.NewStorageIntegrationUnsetRequest().
					WithS3Params(*sdk.NewUnsetS3StorageParamsRequest().
						// unsetting private link omitted on purpose - check "alter: unset privatelink endpoint does not work" test
						WithStorageAwsObjectAcl(true).
						WithStorageAwsExternalId(true),
					).
					WithEnabled(true).
					WithStorageBlockedLocations(true).
					WithComment(true),
			)
		err = client.StorageIntegrations.Alter(ctx, unset)
		require.NoError(t, err)

		props, err = client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertS3StorageIntegrationDescResult(t, props, false, changedS3AllowedLocations, []sdk.StorageLocation{}, "", true)
		assert.NotEqual(t, "new-external-id", findProp(t, props, "STORAGE_AWS_EXTERNAL_ID").Value)
	})

	// TODO [SNOW-2356049]: Adjust this test when UNSET starts working correctly
	t.Run("alter: unset privatelink endpoint does not work", func(t *testing.T) {
		id := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)

		req := sdk.NewAlterStorageIntegrationRequest(id).WithUnset(
			*sdk.NewStorageIntegrationUnsetRequest().WithS3Params(
				*sdk.NewUnsetS3StorageParamsRequest().WithUsePrivatelinkEndpoint(true),
			),
		)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.ErrorContains(t, err, "Cannot unset property 'USE_PRIVATELINK_ENDPOINT' on integration")
	})

	// The docs currently list STORAGE_AWS_ROLE_ARN as a required parameter in ALTER. It seems the alter can be successfully run without it.
	t.Run("alter: S3, without STORAGE_AWS_ROLE_ARN", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		req := sdk.NewAlterStorageIntegrationRequest(id).WithSet(
			*sdk.NewStorageIntegrationSetRequest().WithS3Params(
				*sdk.NewSetS3StorageParamsRequest().
					WithStorageAwsExternalId("new-external-id"),
			),
		)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.NoError(t, err)
	})

	// TODO [SNOW-2356128]: Add test for use_privatelink_endpoint on Azure deployment (preprod?)
	t.Run("alter: azure, set and unset", func(t *testing.T) {
		id := createAzureStorageIntegration(t)

		changedAzureAllowedLocations := append([]sdk.StorageLocation{{Path: azureBucketUrl + "/allowed-location3"}}, azureAllowedLocations...)
		changedAzureBlockedLocations := append([]sdk.StorageLocation{{Path: azureBucketUrl + "/blocked-location3"}}, azureBlockedLocations...)
		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithSet(
				*sdk.NewStorageIntegrationSetRequest().
					WithAzureParams(*sdk.NewSetAzureStorageParamsRequest().WithAzureTenantId(azureTenantId)).
					WithEnabled(true).
					WithStorageAllowedLocations(changedAzureAllowedLocations).
					WithStorageBlockedLocations(changedAzureBlockedLocations).
					WithComment("changed comment"),
			)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertAzureStorageIntegrationDescResult(t, props, true, changedAzureAllowedLocations, changedAzureBlockedLocations, "changed comment")

		unset := sdk.NewAlterStorageIntegrationRequest(id).WithUnset(
			*sdk.NewStorageIntegrationUnsetRequest().
				// TODO [SNOW-2356128]: Add test for unsetting use_privatelink_endpoint on Azure deployment (preprod?)
				WithEnabled(true).
				WithStorageBlockedLocations(true).
				WithComment(true),
		)
		err = client.StorageIntegrations.Alter(ctx, unset)
		require.NoError(t, err)

		props, err = client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertAzureStorageIntegrationDescResult(t, props, false, changedAzureAllowedLocations, []sdk.StorageLocation{}, "")
	})

	// TODO [SNOW-2356128]: Unskip when can be run on Azure deployment (preprod?)
	t.Run("alter: azure, set without AZURE_TENANT_ID", func(t *testing.T) {
		t.Skip("Unskip when can be run on the azure env. Current error: 511300 (0A000): SQL compilation error: Privatelink endpoints are not supported with 'Azure Storage' locations on 'aws' platform.")
		id := createAzureStorageIntegration(t)

		req := sdk.NewAlterStorageIntegrationRequest(id).WithSet(
			*sdk.NewStorageIntegrationSetRequest().WithAzureParams(
				*sdk.NewSetAzureStorageParamsRequest().WithUsePrivatelinkEndpoint(true),
			),
		)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.NoError(t, err)
	})

	t.Run("alter: gcs, set and unset", func(t *testing.T) {
		id := createGCSStorageIntegrationBasic(t)

		changedGcsAllowedLocations := append([]sdk.StorageLocation{{Path: gcsBucketUrl + "/allowed-location3"}}, gcsAllowedLocations...)
		changedGcsBlockedLocations := append([]sdk.StorageLocation{{Path: gcsBucketUrl + "/blocked-location3"}}, gcsBlockedLocations...)
		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithSet(
				*sdk.NewStorageIntegrationSetRequest().
					WithEnabled(true).
					WithStorageAllowedLocations(changedGcsAllowedLocations).
					WithStorageBlockedLocations(changedGcsBlockedLocations).
					WithComment("changed comment"),
			)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertGCSStorageIntegrationDescResult(t, props, true, changedGcsAllowedLocations, changedGcsBlockedLocations, "changed comment")

		unset := sdk.NewAlterStorageIntegrationRequest(id).WithUnset(
			*sdk.NewStorageIntegrationUnsetRequest().
				WithEnabled(true).
				WithStorageBlockedLocations(true).
				WithComment(true),
		)
		err = client.StorageIntegrations.Alter(ctx, unset)
		require.NoError(t, err)

		props, err = client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertGCSStorageIntegrationDescResult(t, props, false, changedGcsAllowedLocations, []sdk.StorageLocation{}, "")
	})

	t.Run("show: without like", func(t *testing.T) {
		id := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)
		id2 := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)

		integrations, err := client.StorageIntegrations.Show(ctx, sdk.NewShowStorageIntegrationRequest())
		require.NoError(t, err)
		ids := collections.Map(integrations, func(i sdk.StorageIntegration) sdk.AccountObjectIdentifier { return i.ID() })

		require.GreaterOrEqual(t, len(ids), 2)
		require.Contains(t, ids, id)
		require.Contains(t, ids, id2)
	})

	t.Run("show: with like", func(t *testing.T) {
		id := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)
		id2 := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)

		integrations, err := client.StorageIntegrations.Show(ctx, sdk.NewShowStorageIntegrationRequest().WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
		require.NoError(t, err)
		ids := collections.Map(integrations, func(i sdk.StorageIntegration) sdk.AccountObjectIdentifier { return i.ID() })

		require.Len(t, ids, 1)
		require.Contains(t, ids, id)
		require.NotContains(t, ids, id2)
	})

	t.Run("show: s3, no matches", func(t *testing.T) {
		integrations, err := client.StorageIntegrations.Show(ctx, sdk.NewShowStorageIntegrationRequest().
			WithLike(sdk.Like{Pattern: sdk.String(NonExistingSchemaObjectIdentifier.Name())}),
		)
		require.NoError(t, err)
		require.Empty(t, integrations)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)

		err := client.StorageIntegrations.Drop(ctx, sdk.NewDropStorageIntegrationRequest(id))
		require.NoError(t, err)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.StorageIntegrations.Drop(ctx, sdk.NewDropStorageIntegrationRequest(id))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show by id safely: existing", func(t *testing.T) {
		id := createS3StorageIntegrationBasic(t, sdk.RegularS3Protocol)

		storageIntegration, err := client.StorageIntegrations.ShowByIDSafely(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "")
	})

	t.Run("show by id safely: non-existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		_, err := client.StorageIntegrations.ShowByIDSafely(ctx, id)
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	// TODO [next PR]: adjust describe tests when introducing details and dedicated assertions
	t.Run("describe: s3", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		desc, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertS3StorageIntegrationDescResult(t, desc, true, s3AllowedLocations, s3BlockedLocations, "some comment", true)
	})

	t.Run("describe: gcs", func(t *testing.T) {
		id := createGCSStorageIntegration(t)

		desc, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertGCSStorageIntegrationDescResult(t, desc, true, gcsAllowedLocations, gcsBlockedLocations, "some comment")
	})

	t.Run("describe: azure", func(t *testing.T) {
		id := createAzureStorageIntegration(t)

		desc, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertAzureStorageIntegrationDescResult(t, desc, true, azureAllowedLocations, azureBlockedLocations, "some comment")
	})
}
