//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_TagsAccountLevel(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

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
	t.Run("account object OauthForPartnerApplications", func(t *testing.T) {
		securityIntegration, securityIntegrationCleanup := testClientHelper().SecurityIntegration.CreateOauthForPartnerApplications(t)
		t.Cleanup(securityIntegrationCleanup)
		id := securityIntegration.ID()

		err := client.SecurityIntegrations.AlterOauthForPartnerApplications(ctx, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSetTags(tags))
		require.NoError(t, err)
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeIntegration, tagValue)

		err = client.SecurityIntegrations.AlterOauthForPartnerApplications(ctx, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithUnsetTags(unsetTags))
		require.NoError(t, err)
		assertTagUnset(t, tag.ID(), id, sdk.ObjectTypeIntegration)

		testSetAndUnsetInTagObject(t, tags[0], id, sdk.ObjectTypeIntegration)
	})

	// The failover group helper creates a group with ROLES, an account object type that can be a part of only one failover/replication group per account.
	// The failover group SDK has no dedicated SET TAG / UNSET TAG in ALTER, so the generic tag methods are used.
	t.Run("account object FailoverGroup", func(t *testing.T) {
		failoverGroup, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(failoverGroupCleanup)
		id := failoverGroup.ID()

		err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeFailoverGroup, id).WithSetTags(tags))
		require.NoError(t, err)
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeFailoverGroup, tagValue)

		err = client.Tags.Unset(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeFailoverGroup, id).WithUnsetTags(unsetTags))
		require.NoError(t, err)
		assertTagUnset(t, tag.ID(), id, sdk.ObjectTypeFailoverGroup)

		testSetAndUnsetInTagObject(t, tags[0], id, sdk.ObjectTypeFailoverGroup)
	})
}
