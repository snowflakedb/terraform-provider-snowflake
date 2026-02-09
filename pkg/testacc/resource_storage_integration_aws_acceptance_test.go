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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
		{Path: awsBucketUrl + "allowed-location/"},
	}
	allowedLocations2 := []sdk.StorageLocation{
		{Path: awsBucketUrl + "allowed-location/"},
		{Path: awsBucketUrl + "allowed-location2/"},
	}
	blockedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "blocked-location/"},
	}
	blockedLocations2 := []sdk.StorageLocation{
		{Path: awsBucketUrl + "blocked-location/"},
		{Path: awsBucketUrl + "blocked-location2/"},
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
		WithUsePrivatelinkEndpoint(r.BooleanTrue).
		WithComment(newComment).
		WithStorageAwsExternalId(externalId2).
		WithStorageAwsObjectAcl("bucket-owner-full-control")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAws),
		Steps: []resource.TestStep{
			// CREATE WITHOUT ATTRIBUTES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsAllAttributesChanged.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
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
						HasStorageAwsObjectAclEmpty(),
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
				ImportStateVerifyIgnore: []string{"use_privatelink_endpoint", "storage_aws_external_id"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStorageIntegrationAwsResource(t, id.Name()).
						HasUsePrivatelinkEndpointString(r.BooleanFalse).
						HasStorageAwsExternalIdNotEmpty(),
				),
			},
			// DESTROY
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsAllAttributesChanged.ResourceReference(), plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, storageIntegrationAwsModelNoAttributes),
				Destroy: true,
			},
			// CREATE WITH ALL ATTRIBUTES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsAllAttributesChanged.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsAllAttributes),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsAllAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocations(allowedLocations...).
						HasStorageBlockedLocations(blockedLocations...).
						HasCommentString(comment).
						HasStorageAwsRoleArnString(awsRoleArn).
						HasStorageAwsExternalIdString(externalId).
						HasStorageAwsObjectAclString("bucket-owner-full-control"),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationAwsAllAttributes.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(false).
						HasComment(comment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationAwsDescribeOutput(t, storageIntegrationAwsAllAttributes.ResourceReference()).
						HasId(id).
						HasEnabled(false).
						HasProvider(string(sdk.RegularS3Protocol)).
						HasComment(comment).
						HasUsePrivatelinkEndpoint(false).
						HasIamUserArnSet().
						HasRoleArn(awsRoleArn).
						HasExternalId(externalId).
						HasObjectAcl("bucket-owner-full-control"),
				),
			},
			// CHANGE PROPERTIES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsAllAttributesChanged.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsAllAttributesChanged),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsAllAttributesChanged.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanTrue).
						HasStorageAllowedLocations(allowedLocations2...).
						HasStorageBlockedLocations(blockedLocations2...).
						HasCommentString(newComment).
						HasStorageAwsRoleArnString(awsRoleArn).
						HasStorageAwsExternalIdString(externalId2).
						HasStorageAwsObjectAclString("bucket-owner-full-control"),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationAwsAllAttributesChanged.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(true).
						HasComment(newComment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationAwsDescribeOutput(t, storageIntegrationAwsAllAttributesChanged.ResourceReference()).
						HasId(id).
						HasEnabled(true).
						HasProvider(string(sdk.RegularS3Protocol)).
						HasComment(newComment).
						HasUsePrivatelinkEndpoint(true).
						HasIamUserArnSet().
						HasRoleArn(awsRoleArn).
						HasExternalId(externalId2).
						HasObjectAcl("bucket-owner-full-control"),
				),
			},
			// IMPORT
			{
				ResourceName:            storageIntegrationAwsAllAttributesChanged.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"use_privatelink_endpoint"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStorageIntegrationAwsResource(t, id.Name()).
						HasUsePrivatelinkEndpointString(r.BooleanTrue).
						HasStorageAwsExternalIdString(externalId2),
				),
			},
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
