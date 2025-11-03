//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecurityIntegrations_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner

	scimModel1 := model.ScimSecurityIntegration("test1", idOne.Name(), false, role.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric))
	scimModel2 := model.ScimSecurityIntegration("test2", idTwo.Name(), false, role.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric))
	scimModel3 := model.ScimSecurityIntegration("test3", idThree.Name(), false, role.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric))

	securityIntegrationsModelLikeFirst := datasourcemodel.SecurityIntegrations("test").
		WithLike(idOne.Name()).
		WithDependsOn(scimModel1.ResourceReference(), scimModel2.ResourceReference(), scimModel3.ResourceReference())

	securityIntegrationsModelLikePrefix := datasourcemodel.SecurityIntegrations("test").
		WithLike(prefix+"%").
		WithDependsOn(scimModel1.ResourceReference(), scimModel2.ResourceReference(), scimModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, scimModel1, scimModel2, scimModel3, securityIntegrationsModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModelLikeFirst.DatasourceReference(), "security_integrations.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, scimModel1, scimModel2, scimModel3, securityIntegrationsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModelLikePrefix.DatasourceReference(), "security_integrations.#", "2"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_CompleteUseCase(t *testing.T) {
	prefix := random.AlphaN(4)
	scimIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	saml2IntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	oauthPartnerIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	oauthCustomIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	externalOauthIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	apiAuthAuthCodeIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	apiAuthClientCredentialsIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	comment := random.Comment()
	role := snowflakeroles.GenericScimProvisioner

	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	externalOauthRole, externalOauthRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(externalOauthRoleCleanup)

	preAuthorizedRole, preAuthorizedRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(preAuthorizedRoleCleanup)

	saml2Cert := random.GenerateX509(t)
	saml2Issuer := testClient().Ids.Alpha()
	saml2ValidUrl := "https://example.com"
	saml2TemporaryVariableName := "saml2_x509_cert"
	saml2TemporaryVariableModel, saml2ConfigVariables := accconfig.SecretStringVariableModelWithConfigVariables(saml2TemporaryVariableName, saml2Cert)

	externalOauthIssuer := random.String()
	externalOauthClaim := random.AlphaN(6)
	externalOauthMappingAttribute := random.AlphaN(6)
	externalOauthAudience := random.AlphaN(6)

	oauthCustomKey, _ := random.GenerateRSAPublicKey(t)

	apiAuthPass1 := random.Password()
	apiAuthPass2 := random.Password()

	scimModel := model.ScimSecurityIntegration("test", scimIntegrationId.Name(), false, role.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric)).
		WithComment(comment)

	saml2Model := model.Saml2SecurityIntegrationVar("test1", saml2IntegrationId.Name(), saml2Issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), saml2ValidUrl, saml2TemporaryVariableName).
		WithComment(comment).
		WithEnabled(datasources.BooleanTrue)

	oauthPartnerModel := model.OauthIntegrationForPartnerApplications("test2", oauthPartnerIntegrationId.Name(), string(sdk.OauthSecurityIntegrationClientTableauServer)).
		WithComment(comment).
		WithEnabled(datasources.BooleanTrue).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN").
		WithOauthIssueRefreshTokens(datasources.BooleanFalse).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))

	oauthCustomModel := model.OauthIntegrationForCustomClients("test3", oauthCustomIntegrationId.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), "https://example.com").
		WithComment(comment).
		WithEnabled(datasources.BooleanTrue).
		WithOauthClientRsaPublicKey(oauthCustomKey).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)).
		WithPreAuthorizedRoles(preAuthorizedRole.ID()).
		WithNetworkPolicy(networkPolicy.ID().Name())

	externalOauthModel := model.ExternalOauthSecurityIntegration("test4", externalOauthIntegrationId.Name(), true, externalOauthIssuer, string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress), []string{externalOauthClaim}, string(sdk.ExternalOauthSecurityIntegrationTypeCustom)).
		WithComment(comment).
		WithExternalOauthAllowedRoles(externalOauthRole.ID()).
		WithExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
		WithExternalOauthAudiences(externalOauthAudience).
		WithExternalOauthJwsKeysUrls("https://example.com").
		WithExternalOauthScopeDelimiter(".").
		WithExternalOauthScopeMappingAttribute(externalOauthMappingAttribute)

	apiAuthAuthCodeModel := model.ApiAuthenticationIntegrationWithAuthorizationCodeGrant("test5", apiAuthAuthCodeIntegrationId.Name(), true, "foo", apiAuthPass1).
		WithComment(comment).
		WithOauthAccessTokenValidity(42).
		WithOauthAuthorizationEndpoint("https://example.com").
		WithOauthClientAuthMethod(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)).
		WithOauthRefreshTokenValidity(12345).
		WithOauthTokenEndpoint("https://example.com").
		WithOauthAllowedScopesValue(config.SetVariable(config.StringVariable("foo")))

	apiAuthClientCredsModel := model.ApiAuthenticationIntegrationWithClientCredentials("test6", apiAuthClientCredentialsIntegrationId.Name(), true, apiAuthPass1, apiAuthPass2).
		WithComment(comment).
		WithOauthAccessTokenValidity(42).
		WithOauthClientAuthMethod(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)).
		WithOauthRefreshTokenValidity(12345).
		WithOauthTokenEndpoint("https://example.com").
		WithOauthAllowedScopesValue(config.SetVariable(config.StringVariable("foo")))

	scimNoDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(scimIntegrationId.Name()).
		WithWithDescribe(false).
		WithDependsOn(scimModel.ResourceReference())

	scimWithDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(scimIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(scimModel.ResourceReference())

	oauthPartnerWithDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(oauthPartnerIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(oauthPartnerModel.ResourceReference())

	oauthCustomWithDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(oauthCustomIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(oauthCustomModel.ResourceReference())

	externalOauthWithDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(externalOauthIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(externalOauthModel.ResourceReference())

	saml2WithDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(saml2IntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(saml2Model.ResourceReference())

	apiAuthAuthCodeWithDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(apiAuthAuthCodeIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(apiAuthAuthCodeModel.ResourceReference())

	apiAuthClientCredentialsWithDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(apiAuthClientCredentialsIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(apiAuthClientCredsModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			// SCIM without describe
			{
				Config: accconfig.FromModels(t, scimModel, scimNoDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecurityIntegrationsDatasourceShowOutput(t, scimNoDescribe.DatasourceReference()).
						HasName(scimIntegrationId.Name()).
						HasIntegrationType("SCIM - GENERIC").
						HasCategory("SECURITY").
						HasEnabled(false).
						HasComment(comment).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(scimNoDescribe.DatasourceReference(), "security_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(scimNoDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "0")),
				),
			},
			// SCIM with describe
			{
				Config: accconfig.FromModels(t, scimModel, scimWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecurityIntegrationsDatasourceShowOutput(t, scimWithDescribe.DatasourceReference()).
						HasName(scimIntegrationId.Name()).
						HasIntegrationType("SCIM - GENERIC").
						HasCategory("SECURITY").
						HasEnabled(false).
						HasComment(comment).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(scimWithDescribe.DatasourceReference(), "security_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(scimWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(scimWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "false")),
					assert.Check(resource.TestCheckResourceAttr(scimWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.run_as_role.0.value", "GENERIC_SCIM_PROVISIONER")),
					assert.Check(resource.TestCheckResourceAttr(scimWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.sync_password.0.value", "false")),
					assert.Check(resource.TestCheckResourceAttr(scimWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
				),
			},
			// OAuth Partner Applications with describe
			{
				Config: accconfig.FromModels(t, oauthPartnerModel, oauthPartnerWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecurityIntegrationsDatasourceShowOutput(t, oauthPartnerWithDescribe.DatasourceReference()).
						HasName(oauthPartnerIntegrationId.Name()).
						HasIntegrationType("OAUTH - TABLEAU_SERVER").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic))),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true")),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN")),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_issue_refresh_tokens.0.value", "false")),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "86400")),
					assert.Check(resource.TestCheckResourceAttr(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
					assert.Check(resource.TestCheckResourceAttrSet(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allowed_authorization_endpoints.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(oauthPartnerWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allowed_token_endpoints.0.value")),
				),
			},
			// OAuth Custom Clients with describe
			{
				Config: accconfig.FromModels(t, oauthCustomModel, oauthCustomWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecurityIntegrationsDatasourceShowOutput(t, oauthCustomWithDescribe.DatasourceReference()).
						HasName(oauthCustomIntegrationId.Name()).
						HasIntegrationType("OAUTH - CUSTOM").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential))),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true")),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone))),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "86400")),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.network_policy.0.value", networkPolicy.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
					assert.Check(resource.TestCheckResourceAttrSet(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_rsa_public_key_fp.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(oauthCustomWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value")),
				),
			},
			// External OAuth with describe
			{
				Config: accconfig.FromModels(t, externalOauthModel, externalOauthWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecurityIntegrationsDatasourceShowOutput(t, externalOauthWithDescribe.DatasourceReference()).
						HasName(externalOauthIntegrationId.Name()).
						HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true")),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_issuer.0.value", externalOauthIssuer)),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable))),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_allowed_roles_list.0.value", externalOauthRole.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_audience_list.0.value", externalOauthAudience)),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_token_user_mapping_claim.0.value", fmt.Sprintf("['%s']", externalOauthClaim))),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress))),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_scope_delimiter.0.value", ".")),
					assert.Check(resource.TestCheckResourceAttr(externalOauthWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
				),
			},
			// SAML2 with describe
			{
				Config:          accconfig.FromModels(t, saml2Model, saml2WithDescribe, saml2TemporaryVariableModel),
				ConfigVariables: saml2ConfigVariables,
				Check: assertThat(t,
					resourceshowoutputassert.SecurityIntegrationsDatasourceShowOutput(t, saml2WithDescribe.DatasourceReference()).
						HasName(saml2IntegrationId.Name()).
						HasIntegrationType("SAML2").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(saml2WithDescribe.DatasourceReference(), "security_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(saml2WithDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(saml2WithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.saml2_issuer.0.value", saml2Issuer)),
					assert.Check(resource.TestCheckResourceAttr(saml2WithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.saml2_sso_url.0.value", saml2ValidUrl)),
					assert.Check(resource.TestCheckResourceAttr(saml2WithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.saml2_provider.0.value", "CUSTOM")),
					assert.Check(resource.TestCheckResourceAttr(saml2WithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
				),
			},
			// API Auth Authorization Code with describe
			{
				Config: accconfig.FromModels(t, apiAuthAuthCodeModel, apiAuthAuthCodeWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecurityIntegrationsDatasourceShowOutput(t, apiAuthAuthCodeWithDescribe.DatasourceReference()).
						HasName(apiAuthAuthCodeIntegrationId.Name()).
						HasIntegrationType("API_AUTHENTICATION").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_access_token_validity.0.value", "42")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "12345")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_auth_method.0.value", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost))),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allowed_scopes.0.value", "[foo]")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode)),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.auth_type.0.value", "OAUTH2")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthAuthCodeWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
				),
			},
			// API Auth Client Credentials with describe
			{
				Config: accconfig.FromModels(t, apiAuthClientCredsModel, apiAuthClientCredentialsWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecurityIntegrationsDatasourceShowOutput(t, apiAuthClientCredentialsWithDescribe.DatasourceReference()).
						HasName(apiAuthClientCredentialsIntegrationId.Name()).
						HasIntegrationType("API_AUTHENTICATION").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_access_token_validity.0.value", "42")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "12345")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_auth_method.0.value", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost))),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allowed_scopes.0.value", "[foo]")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials)),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.auth_type.0.value", "OAUTH2")),
					assert.Check(resource.TestCheckResourceAttr(apiAuthClientCredentialsWithDescribe.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_MultipleTypes(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "1")
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "2")
	issuer := testClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "https://example.com"
	role := snowflakeroles.GenericScimProvisioner

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	saml2Model := model.Saml2SecurityIntegrationVar("test", idOne.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName)
	scimModel := model.ScimSecurityIntegration("test", idTwo.Name(), true, role.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric))
	securityIntegrationsModel := datasourcemodel.SecurityIntegrations("test").
		WithLike(prefix+"%").
		WithDependsOn(saml2Model.ResourceReference(), scimModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, scimModel, saml2Model, securityIntegrationsModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.#", "2"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.name", idTwo.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.run_as_role.0.value", "GENERIC_SCIM_PROVISIONER"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.sync_password.0.value", "false"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.name", idOne.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.0.saml2_provider.0.value", "CUSTOM"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckNoResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.0.saml2_x509_cert.0.value"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_SecurityIntegrationNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_SecurityIntegrations/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one security integration"),
			},
		},
	})
}
