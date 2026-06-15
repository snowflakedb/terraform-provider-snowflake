//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_ApiIntegrationExternalMcpOAuth2_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const oauthClientId = "oauth-client-id-123"
	const oauthClientSecret = "oauth-client-secret-456"
	const oauthTokenEndpoint = "https://auth.example.com/token"
	const oauthAuthorizationEndpoint = "https://auth.example.com/authorize"
	const allowedPrefix = "https://mcp.example.com/api/"
	const blockedPrefix = "https://mcp.example.com/api/blocked/"
	apiProvider := string(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationExternalMcpOAuth2("t", id.Name(), []string{allowedPrefix}, true, oauthAuthorizationEndpoint, oauthClientId, oauthClientSecret, oauthTokenEndpoint)
	withOptionals := model.ApiIntegrationExternalMcpOAuth2("t", id.Name(), []string{allowedPrefix}, true, oauthAuthorizationEndpoint, oauthClientId, oauthClientSecret, oauthTokenEndpoint).
		WithOauthClientAuthMethod(string(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost)).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationExternalMcpOAuth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasOauthClientIdString(oauthClientId).
			HasOauthClientSecretString(oauthClientSecret).
			HasOauthTokenEndpointString(oauthTokenEndpoint).
			HasOauthAuthorizationEndpointString(oauthAuthorizationEndpoint).
			HasOauthClientAuthMethodEmpty().
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationExternalMcpOAuth2DescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasNoOauthClientAuthMethod().
			HasNoBlockedPrefixes().
			HasNoComment(),
		objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasAllowedPrefixes(allowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationExternalMcpOAuth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasOauthClientIdString(oauthClientId).
			HasOauthClientSecretString(oauthClientSecret).
			HasOauthTokenEndpointString(oauthTokenEndpoint).
			HasOauthAuthorizationEndpointString(oauthAuthorizationEndpoint).
			HasOauthClientAuthMethodString(string(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost)).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationExternalMcpOAuth2DescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasOauthClientAuthMethod(string(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost)).
			HasComment(comment),
		objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasOauthClientAuthMethod(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpOAuth2),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:                  config.FromModels(t, basic),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
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
				Config:                  config.FromModels(t, withOptionals),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
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

// TestAcc_ApiIntegrationExternalMcpOAuth2_Import verifies that importing a resource created outside Terraform
// produces an in-place update (not destroy-recreate) to sync oauth_client_secret into state, then a noop.
func TestAcc_ApiIntegrationExternalMcpOAuth2_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const oauthClientId = "oauth-client-id-123"
	const oauthClientSecret = "oauth-client-secret-456"
	const oauthTokenEndpoint = "https://auth.example.com/token"
	const oauthAuthorizationEndpoint = "https://auth.example.com/authorize"
	const allowedPrefix = "https://mcp.example.com/api/"

	testModel := model.ApiIntegrationExternalMcpOAuth2("t", id.Name(), []string{allowedPrefix}, true, oauthAuthorizationEndpoint, oauthClientId, oauthClientSecret, oauthTokenEndpoint)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpOAuth2),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					auth := sdk.NewOAuth2McpUserAuthenticationRequest(oauthClientId, oauthClientSecret, oauthTokenEndpoint, oauthAuthorizationEndpoint)
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: allowedPrefix}}, true).
							WithExternalMcpOAuth2ProviderParams(*sdk.NewExternalMcpOAuth2ParamsRequest().WithApiUserAuthentication(*auth)),
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
						plancheck.ExpectResourceAction(testModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(testModel.ResourceReference(), "oauth_client_secret", tfjson.ActionUpdate, nil, sdk.String(oauthClientSecret)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_ApiIntegrationExternalMcpOAuth2_Import_WrongProviderType(t *testing.T) {
	// Create a DynamicClient MCP integration outside Terraform to use as the import target.
	dynamicClientIntegration, dynamicClientCleanup := testClient().ApiIntegration.CreateMcpDynamicClient(t)
	t.Cleanup(dynamicClientCleanup)

	mcpOAuth2Id := testClient().Ids.RandomAccountObjectIdentifier()
	mcpOAuth2Model := model.ApiIntegrationExternalMcpOAuth2("t", mcpOAuth2Id.Name(),
		[]string{"https://mcp.example.com/api/"},
		true,
		"https://auth.example.com/authorize",
		"oauth-client-id-123",
		"oauth-client-secret-456",
		"https://auth.example.com/token",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpOAuth2),
		Steps: []resource.TestStep{
			// Attempt to import a DynamicClient MCP integration via the OAuth2 resource — expects a user auth type mismatch error.
			{
				Config:        config.FromModels(t, mcpOAuth2Model),
				ResourceName:  mcpOAuth2Model.ResourceReference(),
				ImportState:   true,
				ImportStateId: dynamicClientIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_external_mcp_oauth2"),
			},
		},
	})
}
