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

func TestAcc_StorageIntegrationGcs_BasicUseCase(t *testing.T) {
	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	// TODO [next PRs]: extract allowed location logic and use throughout integration and acceptance tests
	allowedLocations := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "allowed-location/"},
	}
	allowedLocations2 := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "allowed-location/"},
		{Path: gcsBucketUrl + "allowed-location2/"},
	}
	blockedLocations := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "blocked-location/"},
	}
	blockedLocations2 := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "blocked-location/"},
		{Path: gcsBucketUrl + "blocked-location2/"},
	}

	comment := random.Comment()
	newComment := random.Comment()

	storageIntegrationGcsModelNoAttributes := model.StorageIntegrationGcs("w", id.Name(), false, allowedLocations)

	storageIntegrationGcsAllAttributes := model.StorageIntegrationGcs("w", id.Name(), false, allowedLocations).
		WithStorageBlockedLocations(blockedLocations).
		WithComment(comment)

	storageIntegrationGcsAllAttributesChanged := model.StorageIntegrationGcs("w", id.Name(), true, allowedLocations2).
		WithStorageBlockedLocations(blockedLocations2).
		WithComment(newComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationGcs),
		Steps: []resource.TestStep{
			// CREATE WITHOUT ATTRIBUTES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationGcsModelNoAttributes.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, storageIntegrationGcsModelNoAttributes),
				Check: assertThat(t,
					resourceassert.StorageIntegrationGcsResource(t, storageIntegrationGcsModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocations(allowedLocations...).
						HasStorageBlockedLocationsEmpty().
						HasCommentString(""),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationGcsModelNoAttributes.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(false).
						HasComment("").
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationGcsDescribeOutput(t, storageIntegrationGcsModelNoAttributes.ResourceReference()).
						HasId(id).
						HasEnabled(false).
						HasAllowedLocations(allowedLocations...).
						HasNoBlockedLocations().
						HasProvider("GCS").
						HasComment("").
						HasUsePrivatelinkEndpoint(false).
						HasServiceAccountSet(),
				),
			},
			// IMPORT
			{
				ResourceName:      storageIntegrationGcsModelNoAttributes.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// DESTROY
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationGcsModelNoAttributes.ResourceReference(), plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, storageIntegrationGcsModelNoAttributes),
				Destroy: true,
			},
			// CREATE WITH ALL ATTRIBUTES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationGcsAllAttributes.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, storageIntegrationGcsAllAttributes),
				Check: assertThat(t,
					resourceassert.StorageIntegrationGcsResource(t, storageIntegrationGcsAllAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocations(allowedLocations...).
						HasStorageBlockedLocations(blockedLocations...).
						HasCommentString(comment),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationGcsAllAttributes.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(false).
						HasComment(comment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationGcsDescribeOutput(t, storageIntegrationGcsAllAttributes.ResourceReference()).
						HasId(id).
						HasEnabled(false).
						HasAllowedLocations(allowedLocations...).
						HasBlockedLocations(blockedLocations...).
						HasComment(comment).
						HasUsePrivatelinkEndpoint(false).
						HasServiceAccountSet(),
				),
			},
			// CHANGE PROPERTIES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationGcsAllAttributesChanged.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, storageIntegrationGcsAllAttributesChanged),
				Check: assertThat(t,
					resourceassert.StorageIntegrationGcsResource(t, storageIntegrationGcsAllAttributesChanged.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanTrue).
						HasStorageAllowedLocations(allowedLocations2...).
						HasStorageBlockedLocations(blockedLocations2...).
						HasCommentString(newComment),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationGcsAllAttributesChanged.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(true).
						HasComment(newComment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationGcsDescribeOutput(t, storageIntegrationGcsAllAttributesChanged.ResourceReference()).
						HasId(id).
						HasEnabled(true).
						HasAllowedLocations(allowedLocations2...).
						HasBlockedLocations(blockedLocations2...).
						HasComment(newComment).
						HasUsePrivatelinkEndpoint(false).
						HasServiceAccountSet(),
				),
			},
			// IMPORT
			{
				ResourceName:      storageIntegrationGcsAllAttributesChanged.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// UNSET ALL
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationGcsModelNoAttributes.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, storageIntegrationGcsModelNoAttributes),
				Check: assertThat(t,
					resourceassert.StorageIntegrationGcsResource(t, storageIntegrationGcsModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasStorageAllowedLocations(allowedLocations...).
						HasStorageBlockedLocationsEmpty().
						HasCommentString(""),
					resourceshowoutputassert.StorageIntegrationShowOutput(t, storageIntegrationGcsModelNoAttributes.ResourceReference()).
						HasName(id.Name()).
						HasEnabled(false).
						HasComment("").
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					resourceshowoutputassert.StorageIntegrationGcsDescribeOutput(t, storageIntegrationGcsModelNoAttributes.ResourceReference()).
						HasId(id).
						HasEnabled(false).
						HasAllowedLocations(allowedLocations...).
						HasNoBlockedLocations().
						HasComment("").
						HasUsePrivatelinkEndpoint(false).
						HasServiceAccountSet(),
				),
			},
		},
	})
}

func TestAcc_StorageIntegrationGcs_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	storageIntegrationGcsModelNoAllowedLocations := model.StorageIntegrationGcs("w", id.Name(), false, []sdk.StorageLocation{})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationGcs),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, storageIntegrationGcsModelNoAllowedLocations),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute storage_allowed_locations requires 1 item minimum`),
			},
		},
	})
}

func TestAcc_StorageIntegrationGcs_ImportValidation(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)

	notificationIntegration, notificationIntegrationCleanup := testClient().NotificationIntegration.Create(t)
	t.Cleanup(notificationIntegrationCleanup)

	awsIntegration, awsIntegrationCleanup := testClient().StorageIntegration.CreateS3(t, awsBucketUrl, awsRoleArn)
	t.Cleanup(awsIntegrationCleanup)

	allowedLocations := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "allowed-location/"},
	}

	storageIntegrationGcsModel := model.StorageIntegrationGcs("w", notificationIntegration.ID().Name(), false, allowedLocations)
	storageIntegrationGcsModel2 := model.StorageIntegrationGcs("w", awsIntegration.ID().Name(), false, allowedLocations)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationGcs),
		Steps: []resource.TestStep{
			// import a different integration category
			{
				Config:        config.FromModels(t, storageIntegrationGcsModel),
				ResourceName:  storageIntegrationGcsModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: notificationIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`Integration %s is not a STORAGE integration`, notificationIntegration.ID().Name())),
			},
			// import a different provider type (AWS)
			{
				Config:        config.FromModels(t, storageIntegrationGcsModel2),
				ResourceName:  storageIntegrationGcsModel2.ResourceReference(),
				ImportState:   true,
				ImportStateId: awsIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(`.*expected GCS storage provider got S3`),
			},
		},
	})
}

func TestAcc_StorageIntegrationGcs_AllowedLocationsUnordered(t *testing.T) {
	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	allowedLocations := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "allowed-location/"},
		{Path: gcsBucketUrl + "allowed-location2/"},
	}
	allowedLocationsDifferentOrder := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "allowed-location2/"},
		{Path: gcsBucketUrl + "allowed-location/"},
	}

	storageIntegrationGcsModel := model.StorageIntegrationGcs("w", id.Name(), false, allowedLocations)
	storageIntegrationGcsModel2 := model.StorageIntegrationGcs("w", id.Name(), false, allowedLocationsDifferentOrder)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationGcs),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, storageIntegrationGcsModel),
				Check: assertThat(t,
					resourceassert.StorageIntegrationGcsResource(t, storageIntegrationGcsModel.ResourceReference()).
						HasStorageAllowedLocations(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationGcsDescribeOutput(t, storageIntegrationGcsModel.ResourceReference()).
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
						plancheck.ExpectResourceAction(storageIntegrationGcsModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, storageIntegrationGcsModel),
				Check: assertThat(t,
					resourceassert.StorageIntegrationGcsResource(t, storageIntegrationGcsModel.ResourceReference()).
						HasStorageAllowedLocations(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationGcsDescribeOutput(t, storageIntegrationGcsModel.ResourceReference()).
						HasAllowedLocations(allowedLocationsDifferentOrder...),
				),
			},
			// change ordering in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(storageIntegrationGcsModel2.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, storageIntegrationGcsModel2),
				Check: assertThat(t,
					resourceassert.StorageIntegrationGcsResource(t, storageIntegrationGcsModel2.ResourceReference()).
						HasStorageAllowedLocations(allowedLocations...),
					resourceshowoutputassert.StorageIntegrationGcsDescribeOutput(t, storageIntegrationGcsModel2.ResourceReference()).
						HasAllowedLocations(allowedLocationsDifferentOrder...),
				),
			},
		},
	})
}
