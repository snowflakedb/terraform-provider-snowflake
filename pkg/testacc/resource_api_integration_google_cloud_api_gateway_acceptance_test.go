//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_ApiIntegrationGoogleCloudApiGateway_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const googleAudience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
	const allowedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod"
	const blockedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod/blocked/"
	apiProvider := string(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationGoogleCloudApiGateway("t", id.Name(), []string{allowedPrefix}, true, googleAudience)
	withOptionals := model.ApiIntegrationGoogleCloudApiGateway("t", id.Name(), []string{allowedPrefix}, true, googleAudience).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGoogleCloudApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasGoogleAudienceString(googleAudience).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationGoogleDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasGoogleAudience(googleAudience).
			HasApiKey("").
			HasComment(""),
		objectassert.ApiIntegrationGoogleDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway).
			HasGoogleAudience(googleAudience).
			HasAllowedPrefixes(allowedPrefix).
			HasNoBlockedPrefixes().
			HasApiKey("").
			HasComment("").
			HasGoogleApiServiceAccountNotEmpty(),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGoogleCloudApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasGoogleAudienceString(googleAudience).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGoogleDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasGoogleAudience(googleAudience).
			HasApiKey("").
			HasComment(comment),
		objectassert.ApiIntegrationGoogleDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway).
			HasGoogleAudience(googleAudience).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasApiKey("").
			HasComment(comment).
			HasGoogleApiServiceAccountNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGoogleCloudApiGateway),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// Import - with optionals
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - external changes
			{
				PreConfig: func() {
					testClient().ApiIntegration.Alter(t, sdk.NewAlterApiIntegrationRequest(id).WithSet(
						*sdk.NewApiIntegrationSetRequest().WithComment(externalComment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
			},
			// Create - with optionals
			{
				PreConfig: func() {
					_, err := testClient().ApiIntegration.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
		},
	})
}

func TestAcc_ApiIntegrationGoogleCloudApiGateway_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const googleAudience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
	const allowedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod"
	const blockedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod/blocked/"
	apiProvider := string(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway)

	comment := random.Comment()

	allAttributes := model.ApiIntegrationGoogleCloudApiGateway("t", id.Name(), []string{allowedPrefix}, true, googleAudience).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	ref := allAttributes.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGoogleCloudApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasGoogleAudienceString(googleAudience).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGoogleDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasGoogleAudience(googleAudience).
			HasApiKey("").
			HasComment(comment).
			HasGoogleApiServiceAccountNotEmpty(),
		objectassert.ApiIntegrationGoogleDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway).
			HasGoogleAudience(googleAudience).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasApiKey("").
			HasComment(comment).
			HasGoogleApiServiceAccountNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGoogleCloudApiGateway),
		Steps: []resource.TestStep{
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import with all attributes
			{
				Config:            config.FromModels(t, allAttributes),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ApiIntegrationGoogleCloudApiGateway_Import_WrongProviderType(t *testing.T) {
	// Create an AWS API integration outside of Terraform to use as the import target.
	awsIntegration, awsCleanup := testClient().ApiIntegration.CreateAws(t)
	t.Cleanup(awsCleanup)

	googleId := testClient().Ids.RandomAccountObjectIdentifier()
	googleModel := model.ApiIntegrationGoogleCloudApiGateway("t", googleId.Name(),
		[]string{"https://gateway-id-123456.uc.gateway.dev/prod"},
		true,
		"api-gateway-id-123456.apigateway.gcp-project.cloud.goog",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGoogleCloudApiGateway),
		Steps: []resource.TestStep{
			// Attempt to import an AWS API integration via the Google Cloud API Gateway resource — expects a provider type mismatch error.
			{
				Config:        config.FromModels(t, googleModel),
				ResourceName:  googleModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: awsIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_google_cloud_api_gateway"),
			},
		},
	})
}

// TestAcc_ApiIntegrationGoogleCloudApiGateway_Import verifies that importing a resource created outside Terraform
// populates state correctly so that no destroy-before-create plan is produced.
func TestAcc_ApiIntegrationGoogleCloudApiGateway_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const googleAudience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
	const allowedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod"
	const blockedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod/blocked/"

	comment := random.Comment()

	testModel := model.ApiIntegrationGoogleCloudApiGateway("t", id.Name(), []string{allowedPrefix}, true, googleAudience).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGoogleCloudApiGateway),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: allowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: blockedPrefix}}).
							WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest(googleAudience)),
					)
					t.Cleanup(cleanup)
				},
				Config:             config.FromModels(t, testModel),
				ResourceName:       testModel.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			{
				Config: config.FromModels(t, testModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAcc_ApiIntegrationGoogleCloudApiGateway_ExternalProviderTypeMismatch(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const allowedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod"

	googleModel := model.ApiIntegrationGoogleCloudApiGateway("t", id.Name(), []string{allowedPrefix}, true, "api-gateway-id-123456.apigateway.gcp-project.cloud.goog")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGoogleCloudApiGateway),
		Steps: []resource.TestStep{
			// Create Google resource.
			{
				Config: config.FromModels(t, googleModel),
			},
			// External change: drop the Google integration and recreate with the same name as Azure API Management.
			// The next read should detect the provider type mismatch and return an error.
			{
				PreConfig: func() {
					testClient().ApiIntegration.DropApiIntegrationFunc(t, id)()
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://apim-hello-world.azure-api.net/dev"}}, true).
							WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest("00000000-0000-0000-0000-000000000000", "11111111-1111-1111-1111-111111111111")),
					)
					t.Cleanup(cleanup)
				},
				Config:      config.FromModels(t, googleModel),
				ExpectError: regexp.MustCompile("could not normalize api_provider value"),
			},
		},
	})
}
