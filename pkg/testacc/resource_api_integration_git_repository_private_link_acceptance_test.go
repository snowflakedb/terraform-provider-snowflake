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

func TestAcc_ApiIntegrationGitRepositoryPrivateLink_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const allowedPrefix = "https://github.com/my-org/"
	const blockedPrefix = "https://github.com/my-org/blocked/"
	apiProvider := string(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{allowedPrefix}, true, true)
	withOptionals := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{allowedPrefix}, true, true).
		WithAllowedAuthenticationSecrets("ALL").
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasUsePrivatelinkEndpointString(r.BooleanTrue).
			HasAllowedAuthenticationSecretsEmpty().
			HasApiBlockedPrefixesEmpty().
			HasTlsTrustedCertificatesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasAllowedAuthenticationSecrets("").
			HasUsePrivatelinkEndpoint(true).
			HasNoTlsTrustedCertificates().
			HasNoBlockedPrefixes().
			HasComment(""),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets("").
			HasUsePrivatelinkEndpoint(true).
			HasAllowedPrefixes(allowedPrefix).
			HasNoBlockedPrefixes().
			HasNoTlsTrustedCertificates().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasUsePrivatelinkEndpointString(r.BooleanTrue).
			HasAllowedAuthenticationSecrets("ALL").
			HasApiBlockedPrefixes(blockedPrefix).
			HasTlsTrustedCertificatesEmpty().
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasAllowedAuthenticationSecrets("ALL").
			HasUsePrivatelinkEndpoint(true).
			HasNoTlsTrustedCertificates().
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets("ALL").
			HasUsePrivatelinkEndpoint(true).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasNoTlsTrustedCertificates().
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryPrivateLink),
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

// TestAcc_ApiIntegrationGitRepositoryPrivateLink_Import verifies that importing a resource created outside Terraform
// produces no destroy-before-create plan.
func TestAcc_ApiIntegrationGitRepositoryPrivateLink_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const allowedPrefix = "https://github.com/my-org/"
	const blockedPrefix = "https://github.com/my-org/blocked/"
	comment := random.Comment()

	testModel := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{allowedPrefix}, true, true).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryPrivateLink),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: allowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: blockedPrefix}}).
							WithGitHttpsApiPrivateLinkProviderParams(*sdk.NewGitHttpsApiPrivateLinkParamsRequest(true)),
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

// TestAcc_ApiIntegrationGitRepositoryPrivateLink_Import_WrongProviderType verifies that importing a non-private-link
// git integration via this resource returns a descriptive error.
func TestAcc_ApiIntegrationGitRepositoryPrivateLink_Import_WrongProviderType(t *testing.T) {
	gitTokenIntegration, gitTokenCleanup := testClient().ApiIntegration.CreateGitToken(t)
	t.Cleanup(gitTokenCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	dummyModel := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(),
		[]string{"https://github.com/my-org/"}, true, true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryPrivateLink),
		Steps: []resource.TestStep{
			// Attempt to import a token-based git integration via the private link resource — expects a mismatch error.
			{
				Config:        config.FromModels(t, dummyModel),
				ResourceName:  dummyModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: gitTokenIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_git_repository_private_link"),
			},
		},
	})
}
