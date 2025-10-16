//go:build non_account_level_tests

package testint

import (
	"strconv"
	"strings"
	"testing"

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

	assertStorageIntegrationShowResult := func(t *testing.T, s *sdk.StorageIntegration, name sdk.AccountObjectIdentifier, comment string) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.True(t, s.Enabled)
		assert.Equal(t, "EXTERNAL_STAGE", s.StorageType)
		assert.Equal(t, "STORAGE", s.Category)
		assert.Equal(t, comment, s.Comment)
	}

	findProp := func(t *testing.T, props []sdk.StorageIntegrationProperty, name string) *sdk.StorageIntegrationProperty {
		t.Helper()
		prop, err := collections.FindFirst(props, func(property sdk.StorageIntegrationProperty) bool { return property.Name == name })
		require.NoError(t, err)
		return prop
	}

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

	createS3StorageIntegration := func(t *testing.T, protocol sdk.S3Protocol) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, s3AllowedLocations).
			WithIfNotExists(true).
			WithS3StorageProviderParams(*sdk.NewS3StorageParamsRequest(protocol, awsRoleARN).WithStorageAwsExternalId("some-external-id").WithUsePrivateLinkEndpoint(true)).
			WithStorageBlockedLocations(s3BlockedLocations).
			WithComment("some comment")

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.StorageIntegrations.Drop(ctx, sdk.NewDropStorageIntegrationRequest(id))
			require.NoError(t, err)
		})

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

		t.Cleanup(func() {
			err := client.StorageIntegrations.Drop(ctx, sdk.NewDropStorageIntegrationRequest(id))
			require.NoError(t, err)
		})

		return id
	}

	// TODO(SNOW-2356128): Test use_privatelink_endpoint
	createAzureStorageIntegration := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, azureAllowedLocations).
			WithIfNotExists(true).
			WithAzureStorageProviderParams(*sdk.NewAzureStorageParamsRequest(sdk.String(azureTenantId))).
			WithStorageBlockedLocations(azureBlockedLocations).
			WithComment("some comment")

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.StorageIntegrations.Drop(ctx, sdk.NewDropStorageIntegrationRequest(id))
			require.NoError(t, err)
		})

		return id
	}

	t.Run("Create - S3", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "some-external-id", findProp(t, props, "STORAGE_AWS_EXTERNAL_ID").Value)
	})

	t.Run("Create - S3GOV", func(t *testing.T) {
		t.Skip("TODO(SNOW-1820099): Setup GOV accounts to be able to run this test on CI")
		id := createS3StorageIntegration(t, sdk.GovS3Protocol)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")
	})

	t.Run("Create - GCS", func(t *testing.T) {
		id := createGCSStorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")
	})

	t.Run("Create - Azure", func(t *testing.T) {
		id := createAzureStorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")
	})

	t.Run("Alter - set - S3", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		changedS3AllowedLocations := append([]sdk.StorageLocation{{Path: awsBucketUrl + "/allowed-location3"}}, s3AllowedLocations...)
		changedS3BlockedLocations := append([]sdk.StorageLocation{{Path: awsBucketUrl + "/blocked-location3"}}, s3BlockedLocations...)
		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithSet(
				*sdk.NewStorageIntegrationSetRequest().
					WithS3Params(
						*sdk.NewSetS3StorageParamsRequest(awsRoleARN).
							WithStorageAwsExternalId("new-external-id").
							WithUsePrivateLinkEndpoint(true),
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
	})

	// TODO(SNOW-2356128): Add test for use_privatelink_endpoint
	t.Run("Alter - set - Azure", func(t *testing.T) {
		id := createAzureStorageIntegration(t)

		changedAzureAllowedLocations := append([]sdk.StorageLocation{{Path: azureBucketUrl + "/allowed-location3"}}, azureAllowedLocations...)
		changedAzureBlockedLocations := append([]sdk.StorageLocation{{Path: azureBucketUrl + "/blocked-location3"}}, azureBlockedLocations...)
		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithSet(
				*sdk.NewStorageIntegrationSetRequest().
					WithAzureParams(*sdk.NewSetAzureStorageParamsRequest(azureTenantId)).
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
	})

	t.Run("Alter - unset", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithUnset(
				*sdk.NewStorageIntegrationUnsetRequest().
					WithStorageAwsObjectAcl(true).
					WithEnabled(true).
					WithStorageBlockedLocations(true).
					WithStorageAwsExternalId(true).
					WithComment(true),
			)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertS3StorageIntegrationDescResult(t, props, false, s3AllowedLocations, []sdk.StorageLocation{}, "", true)
		assert.NotEqual(t, "some-external-id", findProp(t, props, "STORAGE_AWS_EXTERNAL_ID").Value)
	})

	// TODO(SNOW-2356049): Adjust this test when UNSET starts working correctly
	t.Run("Alter - unset privatelink endpoint does not work", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithUnset(
				*sdk.NewStorageIntegrationUnsetRequest().
					WithUsePrivateLinkEndpoint(true),
			)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.ErrorContains(t, err, "Cannot unset property 'USE_PRIVATELINK_ENDPOINT' on integration")
	})

	t.Run("Describe - S3", func(t *testing.T) {
		id := createS3StorageIntegration(t, sdk.RegularS3Protocol)

		desc, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertS3StorageIntegrationDescResult(t, desc, true, s3AllowedLocations, s3BlockedLocations, "some comment", true)
	})

	t.Run("Describe - GCS", func(t *testing.T) {
		id := createGCSStorageIntegration(t)

		desc, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertGCSStorageIntegrationDescResult(t, desc, true, gcsAllowedLocations, gcsBlockedLocations, "some comment")
	})

	t.Run("Describe - Azure", func(t *testing.T) {
		id := createAzureStorageIntegration(t)

		desc, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertAzureStorageIntegrationDescResult(t, desc, true, azureAllowedLocations, azureBlockedLocations, "some comment")
	})
}
