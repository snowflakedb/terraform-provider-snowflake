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

func TestAcc_ApiIntegrationExternalMcpDynamicClient_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiProvider := string(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationExternalMcpDynamicClient("t", id.Name(), []string{mcpAllowedPrefix}, true, mcpOauthResourceUrl)
	withOptionals := model.ApiIntegrationExternalMcpDynamicClient("t", id.Name(), []string{mcpAllowedPrefix}, true, mcpOauthResourceUrl).
		WithApiBlockedPrefixes([]string{mcpBlockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationExternalMcpDynamicClientResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasOauthResourceUrlString(mcpOauthResourceUrl).
			HasApiAllowedPrefixes(mcpAllowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationExternalMcpDynamicClientDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauthDynamicClient)).
			HasOauthResourceUrl(mcpOauthResourceUrl).
			HasNoBlockedPrefixes().
			HasComment(""),
		objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauthDynamicClient).
			HasOauthResourceUrl(mcpOauthResourceUrl).
			HasAllowedPrefixes(mcpAllowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationExternalMcpDynamicClientResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasOauthResourceUrlString(mcpOauthResourceUrl).
			HasApiAllowedPrefixes(mcpAllowedPrefix).
			HasApiBlockedPrefixes(mcpBlockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationExternalMcpDynamicClientDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauthDynamicClient)).
			HasOauthResourceUrl(mcpOauthResourceUrl).
			HasComment(comment),
		objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauthDynamicClient).
			HasOauthResourceUrl(mcpOauthResourceUrl).
			HasAllowedPrefixes(mcpAllowedPrefix).
			HasBlockedPrefixes(mcpBlockedPrefix).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpDynamicClient),
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

// TestAcc_ApiIntegrationExternalMcpDynamicClient_Import verifies that importing a resource created outside
// Terraform produces no destroy-before-create plan.
func TestAcc_ApiIntegrationExternalMcpDynamicClient_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()

	testModel := model.ApiIntegrationExternalMcpDynamicClient("t", id.Name(), []string{mcpAllowedPrefix}, true, mcpOauthResourceUrl).
		WithApiBlockedPrefixes([]string{mcpBlockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpDynamicClient),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					auth := sdk.NewDynamicClientMcpUserAuthenticationRequest(mcpOauthResourceUrl)
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: mcpAllowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: mcpBlockedPrefix}}).
							WithExternalMcpDynamicClientProviderParams(*sdk.NewExternalMcpDynamicClientParamsRequest().WithApiUserAuthentication(*auth)),
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

func TestAcc_ApiIntegrationExternalMcpDynamicClient_Import_WrongProviderType(t *testing.T) {
	// Create an OAuth2 MCP integration outside Terraform to use as the import target.
	oauth2Integration, oauth2Cleanup := testClient().ApiIntegration.CreateMcpOAuth2(t)
	t.Cleanup(oauth2Cleanup)

	dynamicClientId := testClient().Ids.RandomAccountObjectIdentifier()
	dynamicClientModel := model.ApiIntegrationExternalMcpDynamicClient("t", dynamicClientId.Name(),
		[]string{mcpAllowedPrefix},
		true,
		mcpOauthResourceUrl,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpDynamicClient),
		Steps: []resource.TestStep{
			// Attempt to import an OAuth2 MCP integration via the DynamicClient resource — expects a user auth type mismatch error.
			{
				Config:        config.FromModels(t, dynamicClientModel),
				ResourceName:  dynamicClientModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: oauth2Integration.ID().Name(),
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_external_mcp_dynamic_client"),
			},
		},
	})
}
