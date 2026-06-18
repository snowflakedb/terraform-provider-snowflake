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

	apiProvider := string(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationGitRepositoryToken("t", id.Name(), []string{gitAllowedPrefix}, true).
		WithAllAllowedAuthenticationSecrets(true)
	withOptionals := model.ApiIntegrationGitRepositoryToken("t", id.Name(), []string{gitAllowedPrefix}, true).
		WithAllAllowedAuthenticationSecrets(true).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasAllAllowedAuthenticationSecrets(true).
			HasApiAllowedPrefixes(gitAllowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationGitRepositoryTokenDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasAllowedAuthenticationSecrets("ALL").
			HasComment(""),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasAllowedAuthenticationSecrets("ALL").
			HasNoUserAuthType().
			HasUsePrivatelinkEndpoint(false).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasAllAllowedAuthenticationSecrets(true).
			HasApiAllowedPrefixes(gitAllowedPrefix).
			HasApiBlockedPrefixes(gitBlockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitRepositoryTokenDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasAllowedAuthenticationSecrets("ALL").
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasAllowedAuthenticationSecrets("ALL").
			HasNoUserAuthType().
			HasUsePrivatelinkEndpoint(false).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
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

	comment := random.Comment()

	testModel := model.ApiIntegrationGitRepositoryToken("t", id.Name(), []string{gitAllowedPrefix}, true).
		WithAllAllowedAuthenticationSecrets(true).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
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
							[]sdk.ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitBlockedPrefix}}).
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

	secretId, cleanupSecret := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(cleanupSecret)

	withAll := model.ApiIntegrationGitRepositoryToken("t", id.Name(), []string{gitAllowedPrefix}, true).
		WithAllAllowedAuthenticationSecrets(true)
	withNone := model.ApiIntegrationGitRepositoryToken("t", id.Name(), []string{gitAllowedPrefix}, true).
		WithNoAllowedAuthenticationSecrets(true)
	withList := model.ApiIntegrationGitRepositoryToken("t", id.Name(), []string{gitAllowedPrefix}, true).
		WithAllowedAuthenticationSecrets([]string{secretId.FullyQualifiedName()})

	ref := withAll.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryToken),
		Steps: []resource.TestStep{
			// Create with ALL
			{
				Config: config.FromModels(t, withAll),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasAllAllowedAuthenticationSecrets(true),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets("ALL"),
				),
			},
			// ALL → NONE
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate)},
				},
				Config: config.FromModels(t, withNone),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasNoAllowedAuthenticationSecrets(true),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets("NONE"),
				),
			},
			// NONE → specific list
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate)},
				},
				Config: config.FromModels(t, withList),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasAllowedAuthenticationSecrets(secretId.FullyQualifiedName()),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets(secretId.FullyQualifiedName()),
				),
			},
			// specific list → ALL
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate)},
				},
				Config: config.FromModels(t, withAll),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasAllAllowedAuthenticationSecrets(true),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets("ALL"),
				),
			},
			// ALL → specific list
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate)},
				},
				Config: config.FromModels(t, withList),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasAllowedAuthenticationSecrets(secretId.FullyQualifiedName()),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets(secretId.FullyQualifiedName()),
				),
			},
			// specific list → NONE
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate)},
				},
				Config: config.FromModels(t, withNone),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasNoAllowedAuthenticationSecrets(true),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets("NONE"),
				),
			},
			// NONE → ALL
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate)},
				},
				Config: config.FromModels(t, withAll),
				Check: assertThat(t,
					resourceassert.ApiIntegrationGitRepositoryTokenResource(t, ref).
						HasAllAllowedAuthenticationSecrets(true),
					objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
						HasAllowedAuthenticationSecrets("ALL"),
				),
			},
		},
	})
}
