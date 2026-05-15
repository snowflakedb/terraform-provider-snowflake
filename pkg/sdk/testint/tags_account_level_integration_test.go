//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
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

// TestInt_Tags_OnConflict_Bcr2291 verifies that the on_conflict field added to SHOW TAGS by BCR-2291
// (bundle 2026_03) is correctly parsed and surfaced through ShowByID.
func TestInt_Tags_OnConflict_Bcr2291(t *testing.T) {
	t.Run("bundle enabled: create with custom_value is returned by ShowByID", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateWithRequest(t,
			sdk.NewCreateTagRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).
				WithPropagate(*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency).
					WithOnConflict(sdk.TagOnConflict{CustomValue: sdk.String("HIGHLY CONFIDENTIAL")})),
		)
		t.Cleanup(tagCleanup)

		assertThatObject(t, objectassert.Tag(t, tag.ID()).
			HasOnConflict("HIGHLY CONFIDENTIAL"),
		)
	})

	t.Run("bundle enabled: create with allowed_values_sequence is returned by ShowByID", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateWithRequest(t,
			sdk.NewCreateTagRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).
				WithAllowedValues([]string{"confidential", "internal", "public"}).
				WithPropagate(*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency).
					WithOnConflict(sdk.TagOnConflict{AllowedValuesSequence: sdk.Bool(true)})),
		)
		t.Cleanup(tagCleanup)

		assertThatObject(t, objectassert.Tag(t, tag.ID()).
			HasOnConflict("ALLOWED_VALUES_SEQUENCE"),
		)
	})

	t.Run("bundle enabled: alter on_conflict custom value", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateWithRequest(t,
			sdk.NewCreateTagRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).
				WithPropagate(*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency)),
		)
		t.Cleanup(tagCleanup)

		assertThatObject(t, objectassert.Tag(t, tag.ID()).HasOnConflictNil())

		testClientHelper().Tag.Alter(t, sdk.NewAlterTagRequest(tag.ID()).
			WithSet(*sdk.NewTagSetRequest().WithPropagate(*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency).
				WithOnConflict(sdk.TagOnConflict{CustomValue: sdk.String("my_custom_value")}))))

		assertThatObject(t, objectassert.Tag(t, tag.ID()).
			HasOnConflict("my_custom_value"),
		)

		testClientHelper().Tag.Alter(t, sdk.NewAlterTagRequest(tag.ID()).
			WithUnset(*sdk.NewTagUnsetRequest().WithPropagate(true)))

		assertThatObject(t, objectassert.Tag(t, tag.ID()).HasOnConflictNil())
	})

	t.Run("bundle enabled: alter on_conflict allowed values sequence", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateWithRequest(t,
			sdk.NewCreateTagRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).
				WithPropagate(*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency).
					WithOnConflict(sdk.TagOnConflict{AllowedValuesSequence: sdk.Bool(true)})),
		)
		t.Cleanup(tagCleanup)

		assertThatObject(t, objectassert.Tag(t, tag.ID()).HasOnConflict("ALLOWED_VALUES_SEQUENCE"))

		testClientHelper().Tag.Alter(t, sdk.NewAlterTagRequest(tag.ID()).
			WithUnset(*sdk.NewTagUnsetRequest().WithPropagate(true)))

		assertThatObject(t, objectassert.Tag(t, tag.ID()).HasOnConflictNil())
	})

	t.Run("bundle disabled: on_conflict is nil in ShowByID even when set", func(t *testing.T) {
		testClientHelper().BcrBundles.DisableBcrBundle(t, "2026_03")

		tag, tagCleanup := testClientHelper().Tag.CreateWithRequest(t,
			sdk.NewCreateTagRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).
				WithPropagate(*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency).
					WithOnConflict(sdk.TagOnConflict{CustomValue: sdk.String("will_not_appear")})),
		)
		t.Cleanup(tagCleanup)

		assertThatObject(t, objectassert.Tag(t, tag.ID()).HasOnConflictNil())
	})
}
