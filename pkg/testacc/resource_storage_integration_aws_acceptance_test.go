//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

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
	externalId3 := "new-external-id"

	storageIntegrationAwsModelNoAttributes := model.StorageIntegrationAws("w", id.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))
	storageIntegrationAwsModelNoAttributesUsePrivatelinkEndpointExplicit := model.StorageIntegrationAws("w", id.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol)).
		WithUsePrivatelinkEndpoint("false")

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

	storageIntegrationAwsAllAttributesChangedWithDifferentExternalId := model.StorageIntegrationAws("w", id.Name(), true, allowedLocations2, awsRoleArn, string(sdk.RegularS3Protocol)).
		WithStorageBlockedLocations(blockedLocations2).
		WithUsePrivatelinkEndpoint(r.BooleanTrue).
		WithComment(newComment).
		WithStorageAwsExternalId(externalId3).
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
						plancheck.ExpectResourceAction(storageIntegrationAwsModelNoAttributes.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsModelNoAttributes),
				Check: assertThat(
					t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocationsStorageLocation(allowedLocations...).
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
						HasAllowedLocations(allowedLocations...).
						HasNoBlockedLocations().
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
				ResourceName:      storageIntegrationAwsModelNoAttributes.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				// use_privatelink_endpoint is ignored because IMPORT_BOOLEAN_DEFAULT experiment is not enabled
				// in this test
				ImportStateVerifyIgnore: []string{"use_privatelink_endpoint"},
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedStorageIntegrationAwsResource(t, id.Name()).
						HasUsePrivatelinkEndpointString(r.BooleanFalse),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsAllAttributesChanged.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsModelNoAttributesUsePrivatelinkEndpointExplicit),
			},
			// DESTROY
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsModelNoAttributes.ResourceReference(), plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, storageIntegrationAwsModelNoAttributes),
				Destroy: true,
			},
			// CREATE WITH ALL ATTRIBUTES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsAllAttributes.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsAllAttributes),
				Check: assertThat(
					t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsAllAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocationsStorageLocation(allowedLocations...).
						HasStorageBlockedLocationsStorageLocation(blockedLocations...).
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
						HasAllowedLocations(allowedLocations...).
						HasBlockedLocations(blockedLocations...).
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
				Check: assertThat(
					t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsAllAttributesChanged.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanTrue).
						HasStorageAllowedLocationsStorageLocation(allowedLocations2...).
						HasStorageBlockedLocationsStorageLocation(blockedLocations2...).
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
						HasAllowedLocations(allowedLocations2...).
						HasBlockedLocations(blockedLocations2...).
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
				ImportStateVerifyIgnore: []string{"storage_aws_external_id"},
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedStorageIntegrationAwsResource(t, id.Name()).
						HasUsePrivatelinkEndpointString(r.BooleanTrue),
				),
			},
			// CHANGE PROP EXTERNALLY
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterStorageIntegrationRequest(id).WithSet(
						*sdk.NewStorageIntegrationSetRequest().WithS3Params(
							*sdk.NewSetS3StorageParamsRequest().WithStorageAwsExternalId(externalId3),
						),
					)
					testClient().StorageIntegration.Alter(t, alterRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsAllAttributesChanged.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(storageIntegrationAwsAllAttributesChanged.ResourceReference(), "storage_aws_external_id", sdk.String(externalId2), sdk.String(externalId3)),
						planchecks.ExpectChange(storageIntegrationAwsAllAttributesChanged.ResourceReference(), "storage_aws_external_id", tfjson.ActionUpdate, sdk.String(externalId3), sdk.String(externalId2)),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsAllAttributesChanged),
			},
			// CHANGE PROP TO THE CURRENT SNOWFLAKE VALUE
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterStorageIntegrationRequest(id).WithSet(
						*sdk.NewStorageIntegrationSetRequest().WithS3Params(
							*sdk.NewSetS3StorageParamsRequest().WithStorageAwsExternalId(externalId3),
						),
					)
					testClient().StorageIntegration.Alter(t, alterRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsAllAttributesChangedWithDifferentExternalId),
			},
			// UNSET ALL
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsModelNoAttributes.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsModelNoAttributes),
				Check: assertThat(
					t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocationsStorageLocation(allowedLocations...).
						HasStorageBlockedLocationsEmpty().
						HasCommentString("").
						HasStorageAwsRoleArnString(awsRoleArn).
						HasStorageAwsExternalIdEmpty().
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
						HasAllowedLocations(allowedLocations...).
						HasNoBlockedLocations().
						HasProvider(string(sdk.RegularS3Protocol)).
						HasComment("").
						HasUsePrivatelinkEndpoint(false).
						HasIamUserArnSet().
						HasRoleArn(awsRoleArn).
						HasExternalIdSet().
						HasObjectAcl(""),
				),
			},
		},
	})
}

func TestAcc_StorageIntegrationAws_Import(t *testing.T) {
	basicId := testClient().Ids.RandomAccountObjectIdentifier()
	completeId := testClient().Ids.RandomAccountObjectIdentifier()

	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	allowedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "allowed-location/"},
	}
	blockedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "blocked-location/"},
	}
	comment := random.Comment()
	externalId := "some_external_id"

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportBooleanDefault)

	basicStorageIntegrationAwsModel := model.StorageIntegrationAws("w1", basicId.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))

	completeStorageIntegrationAwsModel := model.StorageIntegrationAws("w2", completeId.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol)).
		WithStorageBlockedLocations(blockedLocations).
		WithComment(comment).
		WithStorageAwsExternalId(externalId).
		WithUsePrivatelinkEndpoint("true").
		WithStorageAwsObjectAcl("bucket-owner-full-control")

	basicRef := basicStorageIntegrationAwsModel.ResourceReference()
	completeRef := completeStorageIntegrationAwsModel.ResourceReference()

	_, storageIntegrationCleanup := testClient().StorageIntegration.CreateWithRequest(t, basicId,
		sdk.NewCreateStorageIntegrationRequest(basicId, false, allowedLocations).
			WithS3StorageProviderParams(*sdk.NewS3StorageParamsRequest(sdk.RegularS3Protocol, awsRoleArn)))
	t.Cleanup(storageIntegrationCleanup)

	_, storageIntegrationCleanup = testClient().StorageIntegration.CreateWithRequest(t, completeId,
		sdk.NewCreateStorageIntegrationRequest(completeId, false, allowedLocations).
			WithS3StorageProviderParams(*sdk.NewS3StorageParamsRequest(sdk.RegularS3Protocol, awsRoleArn).
				WithStorageAwsExternalId(externalId).
				WithUsePrivatelinkEndpoint(true).
				WithStorageAwsObjectAcl("bucket-owner-full-control")).
			WithStorageBlockedLocations(blockedLocations).
			WithComment(comment))
	t.Cleanup(storageIntegrationCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: importBooleanDefaultProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAws),
		Steps: []resource.TestStep{
			// Import basic object
			{
				Config:             config.FromModels(t, providerModel, basicStorageIntegrationAwsModel),
				ResourceName:       basicRef,
				ImportState:        true,
				ImportStateId:      basicId.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			// Import complete object
			{
				Config:             config.FromModels(t, providerModel, basicStorageIntegrationAwsModel, completeStorageIntegrationAwsModel),
				ResourceName:       completeRef,
				ImportState:        true,
				ImportStateId:      completeId.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			// Expect empty plan
			{
				Config: config.FromModels(t, providerModel, basicStorageIntegrationAwsModel, completeStorageIntegrationAwsModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_StorageIntegrationAws_Validations(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	allowedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "allowed-location/"},
	}

	storageIntegrationAwsModelNoAllowedLocations := model.StorageIntegrationAws("w", id.Name(), false, []sdk.StorageLocation{}, awsRoleArn, string(sdk.RegularS3Protocol))
	storageIntegrationAwsModelMissingRole := model.StorageIntegrationAws("w", id.Name(), false, allowedLocations, "", string(sdk.RegularS3Protocol))
	storageIntegrationAwsModelIncorrectProtocol := model.StorageIntegrationAws("w", id.Name(), false, allowedLocations, awsRoleArn, "GCS")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAws),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, storageIntegrationAwsModelNoAllowedLocations),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute storage_allowed_locations requires 1 item minimum`),
			},
			{
				Config:      config.FromModels(t, storageIntegrationAwsModelMissingRole),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "storage_aws_role_arn" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, storageIntegrationAwsModelIncorrectProtocol),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid S3 protocol: GCS`),
			},
		},
	})
}

func TestAcc_StorageIntegrationAws_ImportValidation(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	notificationIntegration, notificationIntegrationCleanup := testClient().NotificationIntegration.Create(t)
	t.Cleanup(notificationIntegrationCleanup)

	azureIntegration, azureIntegrationCleanup := testClient().StorageIntegration.CreateAzure(t, azureBucketUrl, azureTenantId)
	t.Cleanup(azureIntegrationCleanup)

	allowedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "allowed-location/"},
	}

	storageIntegrationAwsModel := model.StorageIntegrationAws("w", notificationIntegration.ID().Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))
	storageIntegrationAwsModel2 := model.StorageIntegrationAws("w", azureIntegration.ID().Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAws),
		Steps: []resource.TestStep{
			// import a different integration category
			{
				Config:        config.FromModels(t, storageIntegrationAwsModel),
				ResourceName:  storageIntegrationAwsModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: notificationIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`Integration %s is not a STORAGE integration`, notificationIntegration.ID().Name())),
			},
			// import a different provider type (Azure)
			{
				Config:        config.FromModels(t, storageIntegrationAwsModel2),
				ResourceName:  storageIntegrationAwsModel2.ResourceReference(),
				ImportState:   true,
				ImportStateId: azureIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(`invalid S3 protocol: AZURE`),
			},
		},
	})
}

func TestAcc_StorageIntegrationAws_AllowedLocationsUnordered(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	allowedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "allowed-location/"},
		{Path: awsBucketUrl + "allowed-location2/"},
	}
	allowedLocationsDifferentOrder := []sdk.StorageLocation{
		{Path: awsBucketUrl + "allowed-location2/"},
		{Path: awsBucketUrl + "allowed-location/"},
	}

	storageIntegrationAwsModel := model.StorageIntegrationAws("w", id.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))
	storageIntegrationAwsModel2 := model.StorageIntegrationAws("w", id.Name(), false, allowedLocationsDifferentOrder, awsRoleArn, string(sdk.RegularS3Protocol))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAws),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, storageIntegrationAwsModel),
				Check: assertThat(
					t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsModel.ResourceReference()).
						HasStorageAllowedLocationsStorageLocation(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationAwsDescribeOutput(t, storageIntegrationAwsModel.ResourceReference()).
						HasAllowedLocations(allowedLocations...),
				),
			},
			// change ordering externally
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterStorageIntegrationRequest(id).WithSet(
						*sdk.NewStorageIntegrationSetRequest().WithStorageAllowedLocations(allowedLocationsDifferentOrder),
					)
					testClient().StorageIntegration.Alter(t, alterRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsModel),
				Check: assertThat(
					t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsModel.ResourceReference()).
						HasStorageAllowedLocationsStorageLocation(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationAwsDescribeOutput(t, storageIntegrationAwsModel.ResourceReference()).
						HasAllowedLocations(allowedLocationsDifferentOrder...),
				),
			},
			// change ordering in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAwsModel2.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, storageIntegrationAwsModel2),
				Check: assertThat(
					t,
					resourceassert.StorageIntegrationAwsResource(t, storageIntegrationAwsModel2.ResourceReference()).
						HasStorageAllowedLocationsStorageLocation(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationAwsDescribeOutput(t, storageIntegrationAwsModel2.ResourceReference()).
						HasAllowedLocations(allowedLocationsDifferentOrder...),
				),
			},
		},
	})
}
