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

func TestAcc_ApiIntegrationGitRepositoryGithubApp_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationGitRepositoryGithubApp("t", id.Name(), []string{gitAllowedPrefix}, true)
	withOptionals := model.ApiIntegrationGitRepositoryGithubApp("t", id.Name(), []string{gitAllowedPrefix}, true).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryGithubAppResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasApiAllowedPrefixes(gitAllowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasApiProvider(string(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi)).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeSnowflakeGithubApp)).
			HasComment(""),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeSnowflakeGithubApp).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryGithubAppResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasApiAllowedPrefixes(gitAllowedPrefix).
			HasApiBlockedPrefixes(gitBlockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasApiProvider(string(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi)).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeSnowflakeGithubApp)).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeSnowflakeGithubApp).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryGithubApp),
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

func TestAcc_ApiIntegrationGitRepositoryGithubApp_Import_WrongAuthType(t *testing.T) {
	// Create a token-based git integration to use as the import target.
	tokenIntegration, tokenCleanup := testClient().ApiIntegration.CreateGitToken(t)
	t.Cleanup(tokenCleanup)

	githubAppId := testClient().Ids.RandomAccountObjectIdentifier()
	githubAppModel := model.ApiIntegrationGitRepositoryGithubApp("t", githubAppId.Name(),
		[]string{gitAllowedPrefix},
		true,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryGithubApp),
		Steps: []resource.TestStep{
			// Attempt to import a token-based integration via the GitHub App resource — expects auth type mismatch error.
			{
				Config:        config.FromModels(t, githubAppModel),
				ResourceName:  githubAppModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: tokenIntegration.Name,
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_git_repository_github_app"),
			},
		},
	})
}

func TestAcc_ApiIntegrationGitRepositoryGithubApp_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()

	testModel := model.ApiIntegrationGitRepositoryGithubApp("t", id.Name(), []string{gitAllowedPrefix}, true).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryGithubApp),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitBlockedPrefix}}).
							WithGitHttpsApiGithubAppProviderParams(*sdk.NewGitHttpsApiGithubAppParamsRequest()),
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
