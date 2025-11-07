//go:build non_account_level_tests

package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-1813223): cleanup tests
func TestInt_Tags(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	assertTagHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier, expectedComment string, expectedAllowedValues []string) {
		t.Helper()
		assertThatObject(t, objectassert.Tag(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(expectedComment).
			HasAllowedValuesUnordered(expectedAllowedValues...).
			HasOwnerRoleType("ROLE"),
		)
	}

	t.Run("create tag: comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		request := sdk.NewCreateTagRequest(id).WithComment(&comment)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		assertTagHandle(t, id, comment, nil)
	})

	t.Run("create tag: allowed values", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		values := []string{"value1", "value2"}
		request := sdk.NewCreateTagRequest(id).WithAllowedValues(values)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		assertTagHandle(t, id, "", values)
	})

	t.Run("create tag: comment and allowed values", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		comment := random.Comment()
		values := []string{"value1", "value2"}
		request := sdk.NewCreateTagRequest(id).
			WithOrReplace(true).
			WithComment(&comment).
			WithAllowedValues(values)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		assertTagHandle(t, id, comment, values)
	})

	t.Run("create tag: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Tags.Create(ctx, sdk.NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		assertTagHandle(t, id, "", nil)
	})

	t.Run("drop tag: existing", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()
		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		require.NoError(t, err)

		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop tag: non-existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("undrop tag: existing", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()
		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		require.NoError(t, err)
		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		err = client.Tags.Undrop(ctx, sdk.NewUndropTagRequest(id))
		require.NoError(t, err)
		_, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("alter tag: set and unset comment", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()

		// alter tag with set comment
		comment := random.Comment()
		set := sdk.NewTagSetRequest().WithComment(comment)
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithIfExists(true).WithSet(set))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, comment, tag.Comment)

		// alter tag with unset comment
		unset := sdk.NewTagUnsetRequest().WithComment(true)
		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", tag.Comment)
	})

	t.Run("alter tag: set and unset masking policies", func(t *testing.T) {
		policyTest, policyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(policyCleanup)

		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()

		policies := []sdk.SchemaObjectIdentifier{policyTest.ID()}
		set := sdk.NewTagSetRequest().WithMaskingPolicies(policies)
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(set))
		require.NoError(t, err)

		ref, err := testClientHelper().PolicyReferences.GetPolicyReference(t, tag.ID(), sdk.PolicyEntityDomainTag)
		require.NoError(t, err)
		assert.Equal(t, policyTest.ID().Name(), ref.PolicyName)
		assert.Equal(t, sdk.PolicyKindMaskingPolicy, ref.PolicyKind)

		// assert that setting masking policy does not apply the tag on the masking policy
		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, id, policyTest.ID(), sdk.ObjectTypeMaskingPolicy)
		require.NoError(t, err)
		assert.Nil(t, returnedTagValue)

		unset := sdk.NewTagUnsetRequest().WithMaskingPolicies(policies)
		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)
	})

	t.Run("alter tag: add and drop allowed values", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()

		values := []string{"value1", "value2"}
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd(values))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, values, tag.AllowedValues)

		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithDrop(values))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Empty(t, tag.AllowedValues)
	})

	t.Run("alter tag: rename", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()

		nid := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithRename(nid))
		if err != nil {
			t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))
		} else {
			t.Cleanup(testClientHelper().Tag.DropTagFunc(t, nid))
		}
		require.NoError(t, err)

		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		assertTagHandle(t, nid, "", nil)
	})

	t.Run("alter tag: unset allowed values", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		values := []string{"value1", "value2"}
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd(values))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, values, tag.AllowedValues)

		unset := sdk.NewTagUnsetRequest().WithAllowedValues(true)
		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Empty(t, tag.AllowedValues)
	})

	t.Run("show tag: without like", func(t *testing.T) {
		tag1, tag1Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag1Cleanup)
		tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag2Cleanup)

		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest())
		require.NoError(t, err)

		assert.Len(t, tags, 2)
		assert.Contains(t, tags, *tag1)
		assert.Contains(t, tags, *tag2)
	})

	t.Run("show tag: with like", func(t *testing.T) {
		tag1, tag1Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag1Cleanup)
		tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag2Cleanup)

		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest().WithLike(tag1.Name))
		require.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Contains(t, tags, *tag1)
		assert.NotContains(t, tags, *tag2)
	})

	t.Run("show tag: no matches", func(t *testing.T) {
		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest().WithLike("non-existent"))
		require.NoError(t, err)
		assert.Empty(t, tags)
	})
}

func TestInt_TagsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		_, tag1Cleanup := testClientHelper().Tag.CreateTagWithIdentifier(t, id1)
		t.Cleanup(tag1Cleanup)

		_, tag2Cleanup := testClientHelper().Tag.CreateTagWithIdentifier(t, id2)
		t.Cleanup(tag2Cleanup)

		e1, err := client.Tags.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Tags.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

type IDProvider[T sdk.AccountIdentifier | sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier] interface {
	ID() T
}

func TestInt_TagsAssociations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	tagValue := "abc"
	tags := []sdk.TagAssociation{
		{
			Name:  tag.ID(),
			Value: tagValue,
		},
	}
	unsetTags := []sdk.ObjectIdentifier{
		tag.ID(),
	}

	// TODO(SNOW-1813223): Inline these assertions.
	assertTagSet := func(id sdk.ObjectIdentifier, objectType sdk.ObjectType) {
		t.Helper()
		assertTagSet(t, tag.ID(), id, objectType, tagValue)
	}

	assertTagUnset := func(id sdk.ObjectIdentifier, objectType sdk.ObjectType) {
		t.Helper()
		assertTagUnset(t, tag.ID(), id, objectType)
	}

	testTagSet := func(id sdk.ObjectIdentifier, objectType sdk.ObjectType) {
		t.Helper()
		testSetAndUnsetInTagObject(t, tags[0], id, objectType)
	}

	t.Run("TestInt_TagAssociationForAccountLocator", func(t *testing.T) {
		id := testClientHelper().Ids.AccountIdentifierWithLocator()
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			SetTag: tags,
		})
		require.NoError(t, err)

		assertTagSet(id, sdk.ObjectTypeAccount)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			UnsetTag: unsetTags,
		})
		require.NoError(t, err)

		assertTagUnset(id, sdk.ObjectTypeAccount)

		// test tag sdk method
		err = client.Tags.SetOnCurrentAccount(ctx, sdk.NewSetTagOnCurrentAccountRequest().WithSetTags(tags))
		require.NoError(t, err)

		assertTagSet(id, sdk.ObjectTypeAccount)

		err = client.Tags.UnsetOnCurrentAccount(ctx, sdk.NewUnsetTagOnCurrentAccountRequest().WithUnsetTags(unsetTags))
		require.NoError(t, err)

		assertTagUnset(id, sdk.ObjectTypeAccount)
	})

	t.Run("TestInt_TagAssociationForAccount", func(t *testing.T) {
		id := testClientHelper().Context.CurrentAccountIdentifier(t)
		err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeAccount, id).WithSetTags(tags))
		require.NoError(t, err)

		assertTagSet(id, sdk.ObjectTypeAccount)

		err = client.Tags.Unset(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeAccount, id).WithUnsetTags(unsetTags))
		require.NoError(t, err)

		assertTagUnset(id, sdk.ObjectTypeAccount)
	})

	t.Run("for Organization Account with locator", func(t *testing.T) {
		testClientHelper().EnsureValidNonProdOrganizationAccountIsUsed(t)

		id := testClientHelper().Ids.AccountIdentifierWithLocator()
		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSetTags(tags))
		require.NoError(t, err)

		assertTagSet(id, sdk.ObjectTypeAccount)

		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnsetTags(unsetTags))
		require.NoError(t, err)

		assertTagUnset(id, sdk.ObjectTypeAccount)
	})

	// The test is the same as TestInt_TagAssociationForAccount.
	// They are separated as they intend to test different objects and could be later split to different files if needed.
	t.Run("for Organization Account with account identifier", func(t *testing.T) {
		testClientHelper().EnsureValidNonProdOrganizationAccountIsUsed(t)

		id := testClientHelper().Context.CurrentAccountIdentifier(t)
		err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeAccount, id).WithSetTags(tags))
		require.NoError(t, err)

		assertTagSet(id, sdk.ObjectTypeAccount)

		err = client.Tags.Unset(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeAccount, id).WithUnsetTags(unsetTags))
		require.NoError(t, err)

		assertTagUnset(id, sdk.ObjectTypeAccount)
	})

	accountObjectTestCases := []struct {
		name        string
		objectType  sdk.ObjectType
		setupObject func() (IDProvider[sdk.AccountObjectIdentifier], func())
		setTags     func(sdk.AccountObjectIdentifier, []sdk.TagAssociation) error
		unsetTags   func(sdk.AccountObjectIdentifier, []sdk.ObjectIdentifier) error
	}{
		{
			name:       "ApplicationPackage",
			objectType: sdk.ObjectTypeApplicationPackage,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().ApplicationPackage.CreateApplicationPackage(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.ApplicationPackages.Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.ApplicationPackages.Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "NormalDatabase",
			objectType: sdk.ObjectTypeDatabase,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().Database.CreateDatabase(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
					SetTag: tags,
				})
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
					UnsetTag: tags,
				})
			},
		},
		{
			name:       "DatabaseFromShare",
			objectType: sdk.ObjectTypeDatabase,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return createDatabaseFromShare(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
					SetTag: tags,
				})
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
					UnsetTag: tags,
				})
			},
		},
		// TODO [SNOW-1002023]: Add a test for failover groups; Business Critical Snowflake Edition needed
		{
			name:       "ApiIntegration",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().ApiIntegration.CreateApiIntegration(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "NotificationIntegration",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().NotificationIntegration.Create(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "StorageIntegration",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().StorageIntegration.CreateS3(t, awsBucketUrl, awsRoleARN)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.StorageIntegrations.Alter(ctx, sdk.NewAlterStorageIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.StorageIntegrations.Alter(ctx, sdk.NewAlterStorageIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "ApiAuthenticationWithClientCredentialsFlow",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlow(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.SecurityIntegrations.AlterApiAuthenticationWithClientCredentialsFlow(ctx, sdk.NewAlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.SecurityIntegrations.AlterApiAuthenticationWithClientCredentialsFlow(ctx, sdk.NewAlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "ApiAuthenticationWithAuthorizationCodeGrantFlow",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().SecurityIntegration.CreateApiAuthenticationWithAuthorizationCodeGrantFlow(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.SecurityIntegrations.AlterApiAuthenticationWithAuthorizationCodeGrantFlow(ctx, sdk.NewAlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.SecurityIntegrations.AlterApiAuthenticationWithAuthorizationCodeGrantFlow(ctx, sdk.NewAlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "ComputePool",
			objectType: sdk.ObjectTypeComputePool,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().ComputePool.Create(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithUnsetTags(tags))
			},
		},
		// TODO [SNOW-1452191]: add a test for jwt bearer integration
		{
			name:       "ExternalOauth",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().SecurityIntegration.CreateExternalOauth(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.SecurityIntegrations.AlterExternalOauth(ctx, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.SecurityIntegrations.AlterExternalOauth(ctx, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "OauthForCustomClients",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().SecurityIntegration.CreateOauthForCustomClients(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.SecurityIntegrations.AlterOauthForCustomClients(ctx, sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.SecurityIntegrations.AlterOauthForCustomClients(ctx, sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Saml2",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().SecurityIntegration.CreateSaml2(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.SecurityIntegrations.AlterSaml2(ctx, sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.SecurityIntegrations.AlterSaml2(ctx, sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Scim",
			objectType: sdk.ObjectTypeIntegration,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().SecurityIntegration.CreateScim(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Role",
			objectType: sdk.ObjectTypeRole,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().Role.CreateRole(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Share",
			objectType: sdk.ObjectTypeShare,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return createShare(t, ctx, client)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Shares.Alter(ctx, id, &sdk.AlterShareOptions{
					SetTag: tags,
				})
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Shares.Alter(ctx, id, &sdk.AlterShareOptions{
					UnsetTag: tags,
				})
			},
		},
		{
			name:       "User",
			objectType: sdk.ObjectTypeUser,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().User.CreateUser(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
					SetTag: tags,
				})
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
					UnsetTag: tags,
				})
			},
		},
		{
			name:       "Warehouse",
			objectType: sdk.ObjectTypeWarehouse,
			setupObject: func() (IDProvider[sdk.AccountObjectIdentifier], func()) {
				return testClientHelper().Warehouse.CreateWarehouse(t)
			},
			setTags: func(id sdk.AccountObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
					SetTag: tags,
				})
			},
			unsetTags: func(id sdk.AccountObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
					UnsetTag: tags,
				})
			},
		},
	}

	for _, tc := range accountObjectTestCases {
		t.Run(fmt.Sprintf("account object %s", tc.name), func(t *testing.T) {
			idProvider, cleanup := tc.setupObject()
			t.Cleanup(cleanup)
			id := idProvider.ID()
			err := tc.setTags(id, tags)
			require.NoError(t, err)

			assertTagSet(id, tc.objectType)

			err = tc.unsetTags(id, unsetTags)
			require.NoError(t, err)

			assertTagUnset(id, tc.objectType)

			// test object methods
			testTagSet(id, tc.objectType)
		})
	}

	t.Run("account object Application: invalid operation", func(t *testing.T) {
		applicationPackage, applicationPackageCleanup := createApplicationPackage(t)
		t.Cleanup(applicationPackageCleanup)
		application, applicationCleanup := testClientHelper().Application.CreateApplication(t, applicationPackage.ID(), "V01")
		t.Cleanup(applicationCleanup)
		id := application.ID()

		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithSetTags(tags))
		require.NoError(t, err)

		// TODO(SNOW-1746420): adjust after this is fixed on Snowflake side
		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeApplication)
		require.ErrorContains(t, err, "391801 (0A000): SQL compilation error: Object tagging not supported for object type APPLICATION")

		err = client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUnsetTags(unsetTags))
		require.NoError(t, err)
	})

	t.Run("account object database replica: invalid operation", func(t *testing.T) {
		db, dbCleanup := createDatabaseReplica(t)
		t.Cleanup(dbCleanup)
		id := db.ID()

		err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			SetTag: tags,
		})
		require.ErrorContains(t, err, "is a read-only secondary database and cannot be modified.")
	})

	databaseObjectTestCases := []struct {
		name        string
		objectType  sdk.ObjectType
		setupObject func() (IDProvider[sdk.DatabaseObjectIdentifier], func())
		setTags     func(sdk.DatabaseObjectIdentifier, []sdk.TagAssociation) error
		unsetTags   func(sdk.DatabaseObjectIdentifier, []sdk.ObjectIdentifier) error
	}{
		{
			name:       "DatabaseRole",
			objectType: sdk.ObjectTypeDatabaseRole,
			setupObject: func() (IDProvider[sdk.DatabaseObjectIdentifier], func()) {
				return testClientHelper().DatabaseRole.CreateDatabaseRole(t)
			},
			setTags: func(id sdk.DatabaseObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.DatabaseRoles.Alter(ctx, sdk.NewAlterDatabaseRoleRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.DatabaseObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.DatabaseRoles.Alter(ctx, sdk.NewAlterDatabaseRoleRequest(id).WithUnsetTags(unsetTags))
			},
		},
		{
			name:       "Schema",
			objectType: sdk.ObjectTypeSchema,
			setupObject: func() (IDProvider[sdk.DatabaseObjectIdentifier], func()) {
				return testClientHelper().Schema.CreateSchema(t)
			},
			setTags: func(id sdk.DatabaseObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
					SetTag: tags,
				})
			},
			unsetTags: func(id sdk.DatabaseObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
					UnsetTag: tags,
				})
			},
		},
	}

	for _, tc := range databaseObjectTestCases {
		t.Run(fmt.Sprintf("database object %s", tc.name), func(t *testing.T) {
			idProvider, cleanup := tc.setupObject()
			t.Cleanup(cleanup)
			id := idProvider.ID()
			err := tc.setTags(id, tags)
			require.NoError(t, err)

			assertTagSet(id, tc.objectType)

			err = tc.unsetTags(id, unsetTags)
			require.NoError(t, err)

			assertTagUnset(id, tc.objectType)

			// test object methods
			testTagSet(id, tc.objectType)
		})
	}

	schemaObjectTestCases := []struct {
		name        string
		objectType  sdk.ObjectType
		setupObject func() (IDProvider[sdk.SchemaObjectIdentifier], func())
		setTags     func(sdk.SchemaObjectIdentifier, []sdk.TagAssociation) error
		unsetTags   func(sdk.SchemaObjectIdentifier, []sdk.ObjectIdentifier) error
	}{
		{
			name:       "ExternalTable",
			objectType: sdk.ObjectTypeExternalTable,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return createExternalTable(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				setTags := make([]sdk.TagAssociationRequest, len(tags))
				for i, tag := range tags {
					setTags[i] = *sdk.NewTagAssociationRequest(tag.Name, tag.Value)
				}
				return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithSetTags(setTags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "GitRepository",
			objectType: sdk.ObjectTypeGitRepository,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				origin := "https://github.com/octocat/hello-world"

				gitRepositoryId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

				apiIntegrationId, apiIntegrationCleanup := testClientHelper().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)

				repo, repoCleanup := testClientHelper().GitRepository.Create(t, gitRepositoryId, origin, apiIntegrationId)

				cleanup := func() {
					repoCleanup()
					apiIntegrationCleanup()
				}
				return repo, cleanup
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.GitRepositories.Alter(ctx, sdk.NewAlterGitRepositoryRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.GitRepositories.Alter(ctx, sdk.NewAlterGitRepositoryRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "MaterializedView",
			objectType: sdk.ObjectTypeMaterializedView,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return createMaterializedView(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Pipe",
			objectType: sdk.ObjectTypePipe,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return createPipe(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Pipes.Alter(ctx, id, &sdk.AlterPipeOptions{
					SetTag: tags,
				})
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Pipes.Alter(ctx, id, &sdk.AlterPipeOptions{
					UnsetTag: tags,
				})
			},
		},
		{
			name:       "RowAccessPolicy",
			objectType: sdk.ObjectTypeRowAccessPolicy,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "SessionPolicy",
			objectType: sdk.ObjectTypeSessionPolicy,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return testClientHelper().SessionPolicy.CreateSessionPolicy(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.SessionPolicies.Alter(ctx, sdk.NewAlterSessionPolicyRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.SessionPolicies.Alter(ctx, sdk.NewAlterSessionPolicyRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Stage",
			objectType: sdk.ObjectTypeStage,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return testClientHelper().Stage.CreateStage(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Stream",
			objectType: sdk.ObjectTypeStream,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return createStream(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "EventTable",
			objectType: sdk.ObjectTypeEventTable,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return testClientHelper().EventTable.Create(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Table",
			objectType: sdk.ObjectTypeTable,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return testClientHelper().Table.Create(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				setTags := make([]sdk.TagAssociationRequest, len(tags))
				for i, tag := range tags {
					setTags[i] = *sdk.NewTagAssociationRequest(tag.Name, tag.Value)
				}
				return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithSetTags(setTags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Task",
			objectType: sdk.ObjectTypeTask,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return testClientHelper().Task.Create(t)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "View",
			objectType: sdk.ObjectTypeView,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				return testClientHelper().View.CreateView(t, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES")
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "ImageRepository",
			objectType: sdk.ObjectTypeImageRepository,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
				db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
				t.Cleanup(dbCleanup)

				schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
				t.Cleanup(schemaCleanup)
				return testClientHelper().ImageRepository.CreateInSchema(t, schema.ID())
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.ImageRepositories.Alter(ctx, sdk.NewAlterImageRepositoryRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.ImageRepositories.Alter(ctx, sdk.NewAlterImageRepositoryRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Service",
			objectType: sdk.ObjectTypeService,
			setupObject: func() (IDProvider[sdk.SchemaObjectIdentifier], func()) {
				computePool, computePoolCleanup := testClientHelper().ComputePool.Create(t)
				t.Cleanup(computePoolCleanup)

				db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
				t.Cleanup(dbCleanup)

				schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
				t.Cleanup(schemaCleanup)

				id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
				return testClientHelper().Service.CreateWithId(t, computePool.ID(), id)
			},
			setTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.TagAssociation) error {
				return client.Services.Alter(ctx, sdk.NewAlterServiceRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Services.Alter(ctx, sdk.NewAlterServiceRequest(id).WithUnsetTags(tags))
			},
		},
	}

	for _, tc := range schemaObjectTestCases {
		t.Run(fmt.Sprintf("schema object %s", tc.name), func(t *testing.T) {
			idProvider, cleanup := tc.setupObject()
			t.Cleanup(cleanup)
			id := idProvider.ID()
			err := tc.setTags(id, tags)
			require.NoError(t, err)

			assertTagSet(id, tc.objectType)

			err = tc.unsetTags(id, unsetTags)
			require.NoError(t, err)

			assertTagUnset(id, tc.objectType)

			// test object methods
			testTagSet(id, tc.objectType)
		})
	}

	t.Run("schema object MaskingPolicy", func(t *testing.T) {
		maskingPolicy, cleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(cleanup)
		id := maskingPolicy.ID()
		err := client.MaskingPolicies.Alter(ctx, id, &sdk.AlterMaskingPolicyOptions{
			SetTag: tags,
		})
		require.NoError(t, err)

		assertTagSet(id, sdk.ObjectTypeMaskingPolicy)

		// assert that setting masking policy does not apply the tag on the masking policy
		refs, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, tag.ID(), sdk.PolicyEntityDomainTag)
		require.NoError(t, err)
		assert.Empty(t, refs)

		err = client.MaskingPolicies.Alter(ctx, id, &sdk.AlterMaskingPolicyOptions{
			UnsetTag: unsetTags,
		})
		require.NoError(t, err)

		assertTagUnset(id, sdk.ObjectTypeMaskingPolicy)

		// test object methods
		testTagSet(id, sdk.ObjectTypeMaskingPolicy)
	})

	columnTestCases := []struct {
		name        string
		setupObject func() (sdk.TableColumnIdentifier, func())
		setTags     func(sdk.TableColumnIdentifier, []sdk.TagAssociation) error
		unsetTags   func(sdk.TableColumnIdentifier, []sdk.ObjectIdentifier) error
	}{
		{
			name: "Table",
			setupObject: func() (sdk.TableColumnIdentifier, func()) {
				object, objectCleanup := testClientHelper().Table.Create(t)
				columnId := sdk.NewTableColumnIdentifier(object.ID().DatabaseName(), object.ID().SchemaName(), object.ID().Name(), "ID")
				return columnId, objectCleanup
			},
			setTags: func(id sdk.TableColumnIdentifier, tags []sdk.TagAssociation) error {
				return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id.SchemaObjectId()).WithColumnAction(sdk.NewTableColumnActionRequest().
					WithSetTags(sdk.NewTableColumnAlterSetTagsActionRequest(id.Name(), tags))))
			},
			unsetTags: func(id sdk.TableColumnIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id.SchemaObjectId()).WithColumnAction(sdk.NewTableColumnActionRequest().
					WithUnsetTags(sdk.NewTableColumnAlterUnsetTagsActionRequest(id.Name(), tags))))
			},
		},
		{
			name: "View",
			setupObject: func() (sdk.TableColumnIdentifier, func()) {
				object, objectCleanup := testClientHelper().View.CreateView(t, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES")
				t.Cleanup(objectCleanup)
				columnId := sdk.NewTableColumnIdentifier(object.ID().DatabaseName(), object.ID().SchemaName(), object.ID().Name(), "ROLE_NAME")
				return columnId, objectCleanup
			},
			setTags: func(id sdk.TableColumnIdentifier, tags []sdk.TagAssociation) error {
				return client.Views.Alter(ctx, sdk.NewAlterViewRequest(id.SchemaObjectId()).WithSetTagsOnColumn(
					*sdk.NewViewSetColumnTagsRequest("ROLE_NAME", tags),
				))
			},
			unsetTags: func(id sdk.TableColumnIdentifier, tags []sdk.ObjectIdentifier) error {
				return client.Views.Alter(ctx, sdk.NewAlterViewRequest(id.SchemaObjectId()).WithUnsetTagsOnColumn(
					*sdk.NewViewUnsetColumnTagsRequest("ROLE_NAME", tags),
				))
			},
		},
	}

	for _, tc := range columnTestCases {
		t.Run(fmt.Sprintf("column in %s", tc.name), func(t *testing.T) {
			id, cleanup := tc.setupObject()
			t.Cleanup(cleanup)
			err := tc.setTags(id, tags)
			require.NoError(t, err)

			assertTagSet(id, sdk.ObjectTypeColumn)

			err = tc.unsetTags(id, unsetTags)
			require.NoError(t, err)

			assertTagUnset(id, sdk.ObjectTypeColumn)

			// test object methods
			testTagSet(id, sdk.ObjectTypeColumn)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		name        string
		objectType  sdk.ObjectType
		setupObject func() sdk.SchemaObjectIdentifierWithArguments
		setTags     func(sdk.SchemaObjectIdentifierWithArguments, []sdk.TagAssociation) error
		unsetTags   func(sdk.SchemaObjectIdentifierWithArguments, []sdk.ObjectIdentifier) error
	}{
		{
			name:       "Function",
			objectType: sdk.ObjectTypeFunction,
			setupObject: func() sdk.SchemaObjectIdentifierWithArguments {
				// cleanup is set up in the Create function
				function := testClientHelper().Function.Create(t, sdk.DataTypeInt)
				return function.ID()
			},
			setTags: func(id sdk.SchemaObjectIdentifierWithArguments, tags []sdk.TagAssociation) error {
				return client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifierWithArguments, tags []sdk.ObjectIdentifier) error {
				return client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "ExternalFunction",
			objectType: sdk.ObjectTypeExternalFunction,
			setupObject: func() sdk.SchemaObjectIdentifierWithArguments {
				integration, integrationCleanup := testClientHelper().ApiIntegration.CreateApiIntegration(t)
				t.Cleanup(integrationCleanup)
				// cleanup is set up in the Create function
				function := testClientHelper().ExternalFunction.Create(t, integration.ID(), sdk.DataTypeInt)
				return function.ID()
			},
			setTags: func(id sdk.SchemaObjectIdentifierWithArguments, tags []sdk.TagAssociation) error {
				return client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifierWithArguments, tags []sdk.ObjectIdentifier) error {
				return client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetTags(tags))
			},
		},
		{
			name:       "Procedure",
			objectType: sdk.ObjectTypeProcedure,
			setupObject: func() sdk.SchemaObjectIdentifierWithArguments {
				// cleanup is set up in the Create procedure
				procedure, cleanupProcedure := testClientHelper().Procedure.Create(t, sdk.DataTypeInt)
				t.Cleanup(cleanupProcedure)
				return procedure.ID()
			},
			setTags: func(id sdk.SchemaObjectIdentifierWithArguments, tags []sdk.TagAssociation) error {
				return client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithSetTags(tags))
			},
			unsetTags: func(id sdk.SchemaObjectIdentifierWithArguments, tags []sdk.ObjectIdentifier) error {
				return client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithUnsetTags(tags))
			},
		},
	}

	for _, tc := range schemaObjectWithArgumentsTestCases {
		t.Run(fmt.Sprintf("schema object with arguments %s", tc.name), func(t *testing.T) {
			id := tc.setupObject()
			err := tc.setTags(id, tags)
			require.NoError(t, err)

			assertTagSet(id, tc.objectType)

			err = tc.unsetTags(id, unsetTags)
			require.NoError(t, err)

			assertTagUnset(id, tc.objectType)

			// test object methods
			testTagSet(id, tc.objectType)
		})
	}
}
