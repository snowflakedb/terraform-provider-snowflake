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
}
