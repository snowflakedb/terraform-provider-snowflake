//go:build non_account_level_tests

package testacc

import (
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

func TestAcc_ApiIntegrationGitRepositoryToken_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const allowedPrefix = "https://github.com/my-org/"
	const blockedPrefix = "https://github.com/my-org/blocked/"
	const allowedSecrets = "ALL"
	apiProvider := string(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationGitRepositoryToken("t", id.Name(), allowedSecrets, []string{allowedPrefix}, true)
	withOptionals := model.ApiIntegrationGitRepositoryToken("t", id.Name(), allowedSecrets, []string{allowedPrefix}, true).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasAllowedAuthenticationSecretsString(allowedSecrets).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationGitRepositoryTokenDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasAllowedAuthenticationSecrets(allowedSecrets).
			HasComment(""),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasAllowedAuthenticationSecrets(allowedSecrets).
			HasNoUserAuthType().
			HasUsePrivatelinkEndpoint(false).
			HasAllowedPrefixes(allowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasAllowedAuthenticationSecretsString(allowedSecrets).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitRepositoryTokenDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasAllowedAuthenticationSecrets(allowedSecrets).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasAllowedAuthenticationSecrets(allowedSecrets).
			HasNoUserAuthType().
			HasUsePrivatelinkEndpoint(false).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryToken),
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

func TestAcc_ApiIntegrationGitRepositoryToken_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const allowedPrefix = "https://github.com/my-org/"
	const blockedPrefix = "https://github.com/my-org/blocked/"
	const allowedSecrets = "ALL"
	comment := random.Comment()

	testModel := model.ApiIntegrationGitRepositoryToken("t", id.Name(), allowedSecrets, []string{allowedPrefix}, true).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryToken),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: allowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: blockedPrefix}}).
							WithGitHttpsApiTokenBasedProviderParams(*sdk.NewGitHttpsApiTokenBasedParamsRequest().
								WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithAllSecrets(true))),
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

func TestAcc_ApiIntegrationGitRepositoryToken_AllowedSecrets_Update(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const allowedPrefix = "https://github.com/my-org/"

	withAll := model.ApiIntegrationGitRepositoryToken("t", id.Name(), "ALL", []string{allowedPrefix}, true)
	withNone := model.ApiIntegrationGitRepositoryToken("t", id.Name(), "NONE", []string{allowedPrefix}, true)

	ref := withAll.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryToken),
		Steps: []resource.TestStep{
			// Create with ALL secrets
			{
				Config: config.FromModels(t, withAll),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasAllowedAuthenticationSecretsString("ALL"),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets("ALL"),
				),
			},
			// Update to NONE
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withNone),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasAllowedAuthenticationSecretsString("NONE"),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets("NONE"),
				),
			},
			// Update back to ALL
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withAll),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasAllowedAuthenticationSecretsString("ALL"),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets("ALL"),
				),
			},
		},
	})
}
