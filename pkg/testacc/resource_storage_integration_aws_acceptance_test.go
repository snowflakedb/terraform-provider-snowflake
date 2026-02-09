//go:build non_account_level_tests

package testacc

import (
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// basic
// complex
// validations test
// external changes/dedicated tests?

func TestAcc_StorageIntegrationAws_BasicUseCase(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	// TODO [next PRs]: extract allowed location logic and use throughout integration and acceptance tests
	allowedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "/allowed-location"},
	}
	allowedLocations2 := []sdk.StorageLocation{
		{Path: awsBucketUrl + "/allowed-location"},
		{Path: awsBucketUrl + "/allowed-location2"},
	}
	blockedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "/blocked-location"},
	}
	blockedLocations2 := []sdk.StorageLocation{
		{Path: awsBucketUrl + "/blocked-location"},
		{Path: awsBucketUrl + "/blocked-location2"},
	}

	comment := random.Comment()
	newComment := random.Comment()

	externalId := "some_external_id"
	externalId2 := "some_external_id_2"

	storageIntegrationAwsModelNoAttributes := model.StorageIntegrationAws("w", id.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))

	storageIntegrationAwsAllAttributes := model.StorageIntegrationAws("w", id.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol)).
		WithStorageBlockedLocations(blockedLocations).
		WithComment(comment).
		WithStorageAwsExternalId(externalId).
		WithStorageAwsObjectAcl("bucket-owner-full-control")

	storageIntegrationAwsAllAttributesChanged := model.StorageIntegrationAws("w", id.Name(), true, allowedLocations2, awsRoleArn, string(sdk.RegularS3Protocol)).
		WithStorageBlockedLocations(blockedLocations2).
		WithComment(newComment).
		WithStorageAwsExternalId(externalId2).
		WithStorageAwsObjectAcl("bucket-owner-full-control")

	_ = storageIntegrationAwsAllAttributes
	_ = storageIntegrationAwsAllAttributesChanged

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAws),
		Steps: []resource.TestStep{
			// CREATE WITHOUT ATTRIBUTES
			{
				Config: config.FromModels(t, storageIntegrationAwsModelNoAttributes),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocations(allowedLocations...).
						HasStorageBlockedLocationsEmpty().
						HasCommentString("").
						HasStorageAwsRoleArnString(awsRoleArn).
						HasNoStorageAwsExternalId().
						HasNoStorageAwsObjectAcl(),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationAwsModelNoAttributes.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(false).
						HasComment("").
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationAwsDescribeOutput(t, storageIntegrationAwsModelNoAttributes.ResourceReference()).
						HasId(id).
						HasEnabled(false).
						HasProvider(string(sdk.RegularS3Protocol)).
						HasComment("").
						HasUsePrivatelinkEndpoint(false).
						HasIamUserArnSet().
						HasRoleArn(awsRoleArn).
						HasExternalIdSet().
						HasObjectAcl(""),
				),
			},
			// IMPORT
			{
				ResourceName:            storageIntegrationAwsModelNoAttributes.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"use_privatelink_endpoint"},
			},
			//// DESTROY
			//{
			//	Config:  config.FromModels(t, userModelNoAttributes),
			//	Destroy: true,
			//},
			//// CREATE WITH ALL ATTRIBUTES
			//{
			//	Config: config.FromModels(t, userModelAllAttributes),
			//	Check: assertThat(t,
			//		resourceassert.UserResource(t, userModelAllAttributes.ResourceReference()).
			//			HasNameString(id.Name()).
			//			HasPasswordString(pass).
			//			HasLoginNameString(loginName).
			//			HasDisplayNameString("Display Name").
			//			HasFirstNameString("Jan").
			//			HasMiddleNameString("Jakub").
			//			HasLastNameString("Testowski").
			//			HasEmailString("fake@email.com").
			//			HasMustChangePassword(true).
			//			HasDisabled(false).
			//			HasDaysToExpiryString("8").
			//			HasMinsToUnlockString("9").
			//			HasDefaultWarehouseString("some_warehouse").
			//			HasDefaultNamespaceString("some.namespace").
			//			HasDefaultRoleString("some_role").
			//			HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
			//			HasMinsToBypassMfaString("10").
			//			HasRsaPublicKeyString(key1).
			//			HasRsaPublicKey2String(key2).
			//			HasCommentString(comment).
			//			HasDisableMfaString(r.BooleanTrue).
			//			HasFullyQualifiedNameString(id.FullyQualifiedName()),
			//	),
			//},
			//// CHANGE PROPERTIES
			//{
			//	Config: config.FromModels(t, userModelAllAttributesChanged(newLoginName)),
			//	Check: assertThat(t,
			//		resourceassert.UserResource(t, userModelAllAttributesChanged(newLoginName).ResourceReference()).
			//			HasNameString(id.Name()).
			//			HasPasswordString(newPass).
			//			HasLoginNameString(newLoginName).
			//			HasDisplayNameString("New Display Name").
			//			HasFirstNameString("Janek").
			//			HasMiddleNameString("Kuba").
			//			HasLastNameString("Terraformowski").
			//			HasEmailString("fake@email.net").
			//			HasMustChangePassword(false).
			//			HasDisabled(true).
			//			HasDaysToExpiryString("12").
			//			HasMinsToUnlockString("13").
			//			HasDefaultWarehouseString("other_warehouse").
			//			HasDefaultNamespaceString("one_part_namespace").
			//			HasDefaultRoleString("other_role").
			//			HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
			//			HasMinsToBypassMfaString("14").
			//			HasRsaPublicKeyString(key2).
			//			HasRsaPublicKey2String(key1).
			//			HasCommentString(newComment).
			//			HasDisableMfaString(r.BooleanFalse).
			//			HasFullyQualifiedNameString(id.FullyQualifiedName()),
			//	),
			//},
			//// IMPORT
			//{
			//	ResourceName:            userModelAllAttributesChanged(newLoginName).ResourceReference(),
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"password", "disable_mfa", "days_to_expiry", "mins_to_unlock", "mins_to_bypass_mfa", "default_namespace", "login_name", "show_output.0.days_to_expiry"},
			//	ImportStateCheck: assertThatImport(t,
			//		resourceassert.ImportedUserResource(t, id.Name()).
			//			HasDefaultNamespaceString("ONE_PART_NAMESPACE").
			//			HasLoginNameString(strings.ToUpper(newLoginName)),
			//	),
			//},
			//// CHANGE PROP TO THE CURRENT SNOWFLAKE VALUE
			//{
			//	PreConfig: func() {
			//		testClient().User.SetLoginName(t, id, loginName)
			//	},
			//	Config: config.FromModels(t, userModelAllAttributesChanged(loginName)),
			//	ConfigPlanChecks: resource.ConfigPlanChecks{
			//		PostApplyPostRefresh: []plancheck.PlanCheck{
			//			plancheck.ExpectEmptyPlan(),
			//		},
			//	},
			//},
			//// UNSET ALL
			//{
			//	Config: config.FromModels(t, userModelNoAttributes),
			//	Check: assertThat(t,
			//		resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
			//			HasNameString(id.Name()).
			//			HasPasswordString("").
			//			HasLoginNameString("").
			//			HasDisplayNameString("").
			//			HasFirstNameString("").
			//			HasMiddleNameString("").
			//			HasLastNameString("").
			//			HasEmailString("").
			//			HasMustChangePasswordString(r.BooleanDefault).
			//			HasDisabledString(r.BooleanDefault).
			//			HasDaysToExpiryString("0").
			//			HasMinsToUnlockString(r.IntDefaultString).
			//			HasDefaultWarehouseString("").
			//			HasDefaultNamespaceString("").
			//			HasDefaultRoleString("").
			//			HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault).
			//			HasMinsToBypassMfaString(r.IntDefaultString).
			//			HasRsaPublicKeyString("").
			//			HasRsaPublicKey2String("").
			//			HasCommentString("").
			//			HasDisableMfaString(r.BooleanDefault).
			//			HasFullyQualifiedNameString(id.FullyQualifiedName()),
			//		resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
			//			HasLoginName(strings.ToUpper(id.Name())).
			//			HasDisplayName(""),
			//	),
			//},
		},
	})
}
