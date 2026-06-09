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

func TestAcc_ApiIntegrationAmazonApiGateway_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const awsRoleArn = "arn:aws:iam::000000000001:role/test"
	const allowedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	const blockedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/blocked/"
	apiProvider := string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{allowedPrefix}, awsRoleArn, apiProvider, true)
	withOptionals := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{allowedPrefix}, awsRoleArn, apiProvider, true).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAmazonApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasApiProviderString(apiProvider).
			HasApiAwsRoleArnString(awsRoleArn).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty().
			HasApiKeyEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationAwsDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasApiAwsRoleArn(awsRoleArn).
			HasApiKey("").
			HasComment(""),
		objectassert.ApiIntegrationAwsDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
			HasApiAwsRoleArn(awsRoleArn).
			HasAllowedPrefixes(allowedPrefix).
			HasNoBlockedPrefixes().
			HasApiKey("").
			HasComment("").
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty(),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAmazonApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasApiProviderString(apiProvider).
			HasApiAwsRoleArnString(awsRoleArn).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment).
			HasApiKeyEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationAwsDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasApiAwsRoleArn(awsRoleArn).
			HasApiKey("").
			HasComment(comment),
		objectassert.ApiIntegrationAwsDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
			HasApiAwsRoleArn(awsRoleArn).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasApiKey("").
			HasComment(comment).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
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

func TestAcc_ApiIntegrationAmazonApiGateway_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const awsRoleArn = "arn:aws:iam::000000000001:role/test"
	const allowedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	const blockedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/blocked/"
	apiProvider := string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway)

	comment := random.Comment()
	apiKey := random.AlphanumericN(10)
	externalApiKey := random.AlphanumericN(10)

	allAttributes := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{allowedPrefix}, awsRoleArn, apiProvider, true).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithApiKey(apiKey).
		WithComment(comment)

	ref := allAttributes.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAmazonApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasApiProviderString(apiProvider).
			HasApiAwsRoleArnString(awsRoleArn).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment).
			HasApiKey(apiKey),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationAwsDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasApiAwsRoleArn(awsRoleArn).
			HasApiKey(apiKey).
			HasComment(comment),
		objectassert.ApiIntegrationAwsDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
			HasApiAwsRoleArn(awsRoleArn).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasApiKey(apiKey).
			HasComment(comment).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
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
			// Import with all attributes (api_key is not readable after import, set to "")
			{
				Config:                  config.FromModels(t, allAttributes),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedApiIntegrationAmazonApiGatewayResource(t, id.Name()).
						HasNameString(id.Name()).
						HasApiKeyEmpty(),
				),
			},
			// Change api_key externally — plan detects drift and updates
			{
				PreConfig: func() {
					testClient().ApiIntegration.Alter(t, sdk.NewAlterApiIntegrationRequest(id).WithSet(
						*sdk.NewApiIntegrationSetRequest().WithAwsParams(
							*sdk.NewSetAwsApiParamsRequest().WithApiKey(externalApiKey),
						),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// NOTE: planchecks.ExpectDrift cannot check api_key value — Sensitive field is redacted in plan JSON
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Change api_key to current Snowflake value — expect no-op (DiffSuppressFunc suppresses the diff)
			{
				PreConfig: func() {
					testClient().ApiIntegration.Alter(t, sdk.NewAlterApiIntegrationRequest(id).WithSet(
						*sdk.NewApiIntegrationSetRequest().WithAwsParams(
							*sdk.NewSetAwsApiParamsRequest().WithApiKey(apiKey),
						),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
		},
	})
}

func TestAcc_ApiIntegrationAmazonApiGateway_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const awsRoleArn = "arn:aws:iam::000000000001:role/test"
	const allowedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"

	invalidProvider := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{allowedPrefix}, awsRoleArn, "INVALID_PROVIDER", true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidProvider),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid ApiIntegrationAwsApiProviderType: INVALID_PROVIDER`),
			},
		},
	})
}

func TestAcc_ApiIntegrationAmazonApiGateway_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const awsRoleArn = "arn:aws:iam::000000000001:role/test"
	const allowedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	const blockedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/blocked/"
	apiProvider := string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway)

	comment := random.Comment()
	apiKey := random.AlphanumericN(10)

	allAttributes := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{allowedPrefix}, awsRoleArn, apiProvider, true).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithApiKey(apiKey).
		WithComment(comment)

	ref := allAttributes.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, allAttributes),
			},
			// Import verifies that api_key is empty in state after import (external changes marking)
			{
				Config:                  config.FromModels(t, allAttributes),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedApiIntegrationAmazonApiGatewayResource(t, id.Name()).
						HasNameString(id.Name()).
						HasApiKeyEmpty(),
				),
			},
		},
	})
}
