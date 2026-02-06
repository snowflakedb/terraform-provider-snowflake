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

	flattenLocations := func(locations []sdk.StorageLocation) []string {
		flat := make([]string, len(locations))
		for i, a := range locations {
			flat[i] = a.Path
		}
		return flat
	}

	awsPropertiesAssertions := func(
		t *testing.T,
		totalProps int,
		id sdk.AccountObjectIdentifier,
		props []sdk.StorageIntegrationProperty,
		enabled bool,
		allowedLocations []sdk.StorageLocation,
		blockedLocations []sdk.StorageLocation,
		comment string,
		usePrivateLinkEndpoint bool,
	) *objectassert.StorageIntegrationPropertiesAssert {
		t.Helper()
		return objectassert.StorageIntegrationPropertiesFromObject(t, id, props).
			HasCount(totalProps).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"ENABLED", "Boolean", strconv.FormatBool(enabled), "false"}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_PROVIDER", "String", "S3", ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_ALLOWED_LOCATIONS", "List", strings.Join(flattenLocations(allowedLocations), ","), "[]"}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_BLOCKED_LOCATIONS", "List", strings.Join(flattenLocations(blockedLocations), ","), "[]"}).
			ContainsNotEmptyPropertyWithTypeAndDefault("STORAGE_AWS_IAM_USER_ARN", "String", "").
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_AWS_ROLE_ARN", "String", awsRoleARN, ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"COMMENT", "String", comment, ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"USE_PRIVATELINK_ENDPOINT", "Boolean", strconv.FormatBool(usePrivateLinkEndpoint), "false"})
	}

	gcsPropertiesAssertions := func(
		t *testing.T,
		totalProps int,
		id sdk.AccountObjectIdentifier,
		props []sdk.StorageIntegrationProperty,
		enabled bool,
		allowedLocations []sdk.StorageLocation,
		blockedLocations []sdk.StorageLocation,
		comment string,
	) *objectassert.StorageIntegrationPropertiesAssert {
		t.Helper()
		return objectassert.StorageIntegrationPropertiesFromObject(t, id, props).
			HasCount(totalProps).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"ENABLED", "Boolean", strconv.FormatBool(enabled), "false"}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_PROVIDER", "String", "GCS", ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_ALLOWED_LOCATIONS", "List", strings.Join(flattenLocations(allowedLocations), ","), "[]"}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_BLOCKED_LOCATIONS", "List", strings.Join(flattenLocations(blockedLocations), ","), "[]"}).
			ContainsNotEmptyPropertyWithTypeAndDefault("STORAGE_GCP_SERVICE_ACCOUNT", "String", "").
			// TODO [next PR]: test if use privatelink endpoint can be set on gcs
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"USE_PRIVATELINK_ENDPOINT", "Boolean", "false", "false"}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"COMMENT", "String", comment, ""})
	}

	azurePropertiesAssertions := func(
		t *testing.T,
		totalProps int,
		id sdk.AccountObjectIdentifier,
		props []sdk.StorageIntegrationProperty,
		enabled bool,
		allowedLocations []sdk.StorageLocation,
		blockedLocations []sdk.StorageLocation,
		tenantId string,
		comment string,
		usePrivateLinkEndpoint bool,
	) *objectassert.StorageIntegrationPropertiesAssert {
		t.Helper()

		return objectassert.StorageIntegrationPropertiesFromObject(t, id, props).
			HasCount(totalProps).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"ENABLED", "Boolean", strconv.FormatBool(enabled), "false"}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_PROVIDER", "String", "AZURE", ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_ALLOWED_LOCATIONS", "List", strings.Join(flattenLocations(allowedLocations), ","), "[]"}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_BLOCKED_LOCATIONS", "List", strings.Join(flattenLocations(blockedLocations), ","), "[]"}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"AZURE_TENANT_ID", "String", tenantId, ""}).
			ContainsNotEmptyPropertyWithTypeAndDefault("AZURE_CONSENT_URL", "String", "").
			ContainsNotEmptyPropertyWithTypeAndDefault("AZURE_MULTI_TENANT_APP_NAME", "String", "").
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"COMMENT", "String", comment, ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"USE_PRIVATELINK_ENDPOINT", "Boolean", strconv.FormatBool(usePrivateLinkEndpoint), "false"})
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

		assertThatObject(t, awsPropertiesAssertions(t, 9, id, props, true, s3AllowedLocations, []sdk.StorageLocation{}, "", false).
			DoesNotContainProperty("STORAGE_AWS_OBJECT_ACL").
			ContainsNotEmptyPropertyWithTypeAndDefault("STORAGE_AWS_EXTERNAL_ID", "String", ""),
		)
	})

	t.Run("create: s3", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, awsPropertiesAssertions(t, 10, id, props, true, s3AllowedLocations, s3BlockedLocations, "some comment", true).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_AWS_OBJECT_ACL", "String", "bucket-owner-full-control", ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_AWS_EXTERNAL_ID", "String", "some-external-id", ""}),
		)
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

		assertThatObject(t, gcsPropertiesAssertions(t, 7, id, props, true, gcsAllowedLocations, []sdk.StorageLocation{}, ""))
	})

	t.Run("create: gcs", func(t *testing.T) {
		id := createGCSStorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, gcsPropertiesAssertions(t, 7, id, props, true, gcsAllowedLocations, gcsBlockedLocations, "some comment"))
	})

	t.Run("create: azure basic", func(t *testing.T) {
		id := createAzureStorageIntegrationBasic(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, azurePropertiesAssertions(t, 9, id, props, true, azureAllowedLocations, []sdk.StorageLocation{}, azureTenantId, "", false))
	})

	t.Run("create: azure", func(t *testing.T) {
		id := createAzureStorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, azurePropertiesAssertions(t, 9, id, props, true, azureAllowedLocations, azureBlockedLocations, azureTenantId, "some comment", false))
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

		assertThatObject(t, awsPropertiesAssertions(t, 10, id, props, true, changedS3AllowedLocations, changedS3BlockedLocations, "changed comment", true).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_AWS_OBJECT_ACL", "String", "bucket-owner-full-control", ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_AWS_EXTERNAL_ID", "String", "new-external-id", ""}),
		)

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

		assertThatObject(t, awsPropertiesAssertions(t, 9, id, props, false, changedS3AllowedLocations, []sdk.StorageLocation{}, "", true).
			DoesNotContainProperty("STORAGE_AWS_OBJECT_ACL").
			ContainsPropertyNotEqualToWithTypeAndDefault("STORAGE_AWS_EXTERNAL_ID", "new-external-id", "String", ""),
		)
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

		assertThatObject(t, azurePropertiesAssertions(t, 9, id, props, true, changedAzureAllowedLocations, changedAzureBlockedLocations, azureTenantId, "changed comment", false))

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

		assertThatObject(t, azurePropertiesAssertions(t, 9, id, props, false, changedAzureAllowedLocations, []sdk.StorageLocation{}, azureTenantId, "", false))
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

		assertThatObject(t, gcsPropertiesAssertions(t, 7, id, props, true, changedGcsAllowedLocations, changedGcsBlockedLocations, "changed comment"))

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

		assertThatObject(t, gcsPropertiesAssertions(t, 7, id, props, false, changedGcsAllowedLocations, []sdk.StorageLocation{}, ""))
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

	t.Run("describe: s3", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, awsPropertiesAssertions(t, 10, id, props, true, s3AllowedLocations, s3BlockedLocations, "some comment", true).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_AWS_OBJECT_ACL", "String", "bucket-owner-full-control", ""}).
			ContainsPropertyEqualTo(sdk.StorageIntegrationProperty{"STORAGE_AWS_EXTERNAL_ID", "String", "some-external-id", ""}),
		)
	})

	t.Run("describe: gcs", func(t *testing.T) {
		id := createGCSStorageIntegration(t)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, gcsPropertiesAssertions(t, 7, id, props, true, gcsAllowedLocations, gcsBlockedLocations, "some comment"))
	})

	t.Run("describe: azure", func(t *testing.T) {
		id := createAzureStorageIntegration(t)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, azurePropertiesAssertions(t, 9, id, props, true, azureAllowedLocations, azureBlockedLocations, azureTenantId, "some comment", false))
	})
}
