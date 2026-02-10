//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
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

func TestAcc_StorageIntegrationAzure_BasicUseCase(t *testing.T) {
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	// TODO [next PRs]: extract allowed location logic and use throughout integration and acceptance tests
	allowedLocations := []sdk.StorageLocation{
		{Path: azureBucketUrl + "allowed-location/"},
	}
	allowedLocations2 := []sdk.StorageLocation{
		{Path: azureBucketUrl + "allowed-location/"},
		{Path: azureBucketUrl + "allowed-location2/"},
	}
	blockedLocations := []sdk.StorageLocation{
		{Path: azureBucketUrl + "blocked-location/"},
	}
	blockedLocations2 := []sdk.StorageLocation{
		{Path: azureBucketUrl + "blocked-location/"},
		{Path: azureBucketUrl + "blocked-location2/"},
	}

	comment := random.Comment()
	newComment := random.Comment()

	storageIntegrationAzureModelNoAttributes := model.StorageIntegrationAzure("w", id.Name(), azureTenantId, false, allowedLocations)
	storageIntegrationAzureModelNoAttributesUsePrivatelinkEndpointExplicit := model.StorageIntegrationAzure("w", id.Name(), azureTenantId, false, allowedLocations).
		WithUsePrivatelinkEndpoint("false")

	storageIntegrationAzureAllAttributes := model.StorageIntegrationAzure("w", id.Name(), azureTenantId, false, allowedLocations).
		WithStorageBlockedLocations(blockedLocations).
		WithComment(comment)

	storageIntegrationAzureAllAttributesChanged := model.StorageIntegrationAzure("w", id.Name(), azureTenantId, true, allowedLocations2).
		WithStorageBlockedLocations(blockedLocations2).
		WithComment(newComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAzure),
		Steps: []resource.TestStep{
			// CREATE WITHOUT ATTRIBUTES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAzureModelNoAttributes.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAzureModelNoAttributes),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAzureResource(t, storageIntegrationAzureModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocations(allowedLocations...).
						HasStorageBlockedLocationsEmpty().
						HasCommentString("").
						HasUsePrivatelinkEndpointString(r.BooleanDefault).
						HasAzureTenantIdString(azureTenantId),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationAzureModelNoAttributes.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(false).
						HasComment("").
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationAzureDescribeOutput(t, storageIntegrationAzureModelNoAttributes.ResourceReference()).
						HasId(id).
						HasEnabled(false).
						HasAllowedLocations(allowedLocations...).
						HasNoBlockedLocations().
						HasProvider("AZURE").
						HasComment("").
						HasUsePrivatelinkEndpoint(false).
						HasTenantId(azureTenantId).
						HasConsentUrlSet().
						HasMultiTenantAppNameSet(),
				),
			},
			// IMPORT
			{
				ResourceName:            storageIntegrationAzureModelNoAttributes.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"use_privatelink_endpoint"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStorageIntegrationAzureResource(t, id.Name()).
						HasUsePrivatelinkEndpointString(r.BooleanFalse),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAzureModelNoAttributesUsePrivatelinkEndpointExplicit.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, storageIntegrationAzureModelNoAttributesUsePrivatelinkEndpointExplicit),
			},
			// DESTROY
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAzureModelNoAttributes.ResourceReference(), plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, storageIntegrationAzureModelNoAttributes),
				Destroy: true,
			},
			// CREATE WITH ALL ATTRIBUTES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAzureAllAttributes.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAzureAllAttributes),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAzureResource(t, storageIntegrationAzureAllAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocations(allowedLocations...).
						HasStorageBlockedLocations(blockedLocations...).
						HasCommentString(comment).
						HasUsePrivatelinkEndpointString(r.BooleanDefault).
						HasAzureTenantIdString(azureTenantId),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationAzureAllAttributes.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(false).
						HasComment(comment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationAzureDescribeOutput(t, storageIntegrationAzureAllAttributes.ResourceReference()).
						HasId(id).
						HasEnabled(false).
						HasAllowedLocations(allowedLocations...).
						HasBlockedLocations(blockedLocations...).
						HasComment(comment).
						HasUsePrivatelinkEndpoint(false).
						HasTenantId(azureTenantId).
						HasConsentUrlSet().
						HasMultiTenantAppNameSet(),
				),
			},
			// CHANGE PROPERTIES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAzureAllAttributesChanged.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAzureAllAttributesChanged),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAzureResource(t, storageIntegrationAzureAllAttributesChanged.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanTrue).
						HasStorageAllowedLocations(allowedLocations2...).
						HasStorageBlockedLocations(blockedLocations2...).
						HasCommentString(newComment).
						HasUsePrivatelinkEndpointString(r.BooleanDefault).
						HasAzureTenantIdString(azureTenantId),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationAzureAllAttributesChanged.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(true).
						HasComment(newComment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationAzureDescribeOutput(t, storageIntegrationAzureAllAttributesChanged.ResourceReference()).
						HasId(id).
						HasEnabled(true).
						HasAllowedLocations(allowedLocations2...).
						HasBlockedLocations(blockedLocations2...).
						HasComment(newComment).
						HasUsePrivatelinkEndpoint(false).
						HasTenantId(azureTenantId).
						HasConsentUrlSet().
						HasMultiTenantAppNameSet(),
				),
			},
			// IMPORT
			{
				ResourceName:            storageIntegrationAzureAllAttributesChanged.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"use_privatelink_endpoint"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStorageIntegrationAzureResource(t, id.Name()).
						HasUsePrivatelinkEndpointString(r.BooleanFalse),
				),
			},
			// UNSET ALL
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAzureModelNoAttributes.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, storageIntegrationAzureModelNoAttributes),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAzureResource(t, storageIntegrationAzureModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocations(allowedLocations...).
						HasStorageBlockedLocationsEmpty().
						HasCommentString("").
						HasUsePrivatelinkEndpointString(r.BooleanDefault).
						HasAzureTenantIdString(azureTenantId),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationAzureModelNoAttributes.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(false).
						HasComment("").
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationAzureDescribeOutput(t, storageIntegrationAzureModelNoAttributes.ResourceReference()).
						HasId(id).
						HasEnabled(false).
						HasAllowedLocations(allowedLocations...).
						HasNoBlockedLocations().
						HasComment("").
						HasUsePrivatelinkEndpoint(false).
						HasTenantId(azureTenantId).
						HasConsentUrlSet().
						HasMultiTenantAppNameSet(),
				),
			},
		},
	})
}

func TestAcc_StorageIntegrationAzure_Validations(t *testing.T) {
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	allowedLocations := []sdk.StorageLocation{
		{Path: azureBucketUrl + "allowed-location/"},
	}

	storageIntegrationAzureModelNoAllowedLocations := model.StorageIntegrationAzure("w", id.Name(), azureTenantId, false, []sdk.StorageLocation{})
	storageIntegrationAzureModelMissingTenantId := model.StorageIntegrationAzure("w", id.Name(), "", false, allowedLocations)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAzure),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, storageIntegrationAzureModelNoAllowedLocations),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute storage_allowed_locations requires 1 item minimum`),
			},
			{
				Config:      config.FromModels(t, storageIntegrationAzureModelMissingTenantId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "azure_tenant_id" to not be an empty string`),
			},
		},
	})
}

func TestAcc_StorageIntegrationAzure_ImportValidation(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	notificationIntegration, notificationIntegrationCleanup := testClient().NotificationIntegration.Create(t)
	t.Cleanup(notificationIntegrationCleanup)

	awsIntegration, awsIntegrationCleanup := testClient().StorageIntegration.CreateS3(t, awsBucketUrl, awsRoleArn)
	t.Cleanup(awsIntegrationCleanup)

	allowedLocations := []sdk.StorageLocation{
		{Path: azureBucketUrl + "allowed-location/"},
	}

	storageIntegrationAzureModel := model.StorageIntegrationAzure("w", notificationIntegration.ID().Name(), azureTenantId, false, allowedLocations)
	storageIntegrationAzureModel2 := model.StorageIntegrationAzure("w", awsIntegration.ID().Name(), azureTenantId, false, allowedLocations)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAzure),
		Steps: []resource.TestStep{
			// import a different integration category
			{
				Config:        config.FromModels(t, storageIntegrationAzureModel),
				ResourceName:  storageIntegrationAzureModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: notificationIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`Integration %s is not a STORAGE integration`, notificationIntegration.ID().Name())),
			},
			// import a different provider type (AWS)
			{
				Config:        config.FromModels(t, storageIntegrationAzureModel2),
				ResourceName:  storageIntegrationAzureModel2.ResourceReference(),
				ImportState:   true,
				ImportStateId: awsIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(`.*expected AZURE storage provider got S3`),
			},
		},
	})
}

func TestAcc_StorageIntegrationAzure_AllowedLocationsUnordered(t *testing.T) {
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	allowedLocations := []sdk.StorageLocation{
		{Path: azureBucketUrl + "allowed-location/"},
		{Path: azureBucketUrl + "allowed-location2/"},
	}
	allowedLocationsDifferentOrder := []sdk.StorageLocation{
		{Path: azureBucketUrl + "allowed-location2/"},
		{Path: azureBucketUrl + "allowed-location/"},
	}

	storageIntegrationAzureModel := model.StorageIntegrationAzure("w", id.Name(), azureTenantId, false, allowedLocations)
	storageIntegrationAzureModel2 := model.StorageIntegrationAzure("w", id.Name(), azureTenantId, false, allowedLocationsDifferentOrder)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAzure),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, storageIntegrationAzureModel),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAzureResource(t, storageIntegrationAzureModel.ResourceReference()).
						HasStorageAllowedLocations(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationAzureDescribeOutput(t, storageIntegrationAzureModel.ResourceReference()).
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
						plancheck.ExpectResourceAction(storageIntegrationAzureModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, storageIntegrationAzureModel),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAzureResource(t, storageIntegrationAzureModel.ResourceReference()).
						HasStorageAllowedLocations(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationAzureDescribeOutput(t, storageIntegrationAzureModel.ResourceReference()).
						HasAllowedLocations(allowedLocationsDifferentOrder...),
				),
			},
			// change ordering in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationAzureModel2.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, storageIntegrationAzureModel2),
				Check: assertThat(t,
					resourceassert.StorageIntegrationAzureResource(t, storageIntegrationAzureModel2.ResourceReference()).
						HasStorageAllowedLocations(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationAzureDescribeOutput(t, storageIntegrationAzureModel2.ResourceReference()).
						HasAllowedLocations(allowedLocationsDifferentOrder...),
				),
			},
		},
	})
}
