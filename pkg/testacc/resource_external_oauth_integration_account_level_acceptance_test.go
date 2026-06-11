//go:build account_level_tests

package testacc

import (
	"sort"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalOauthIntegration_completeWithRsaPublicKeysAndBlockedRolesList_paramSet(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	paramCleanup := testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList, "true")
	t.Cleanup(paramCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	expectedRoles := []string{"ACCOUNTADMIN", "SECURITYADMIN", role.ID().Name()}
	sort.Strings(expectedRoles)
	issuer := random.String()
	rsaKey, _ := random.GenerateRSAPublicKey(t)

	integrationModel := model.ExternalOauthSecurityIntegration("test", id.Name(), true, issuer,
		string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress),
		[]string{"foo"},
		string(sdk.ExternalOauthSecurityIntegrationTypeOptionCustom),
	).
		WithComment("foo").
		WithExternalOauthBlockedRoles(role.ID()).
		WithExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable)).
		WithExternalOauthAudiences("foo").
		WithExternalOauthRsaPublicKeyValue(accconfig.MultilineWrapperVariable(rsaKey)).
		WithExternalOauthRsaPublicKey2Value(accconfig.MultilineWrapperVariable(rsaKey)).
		WithExternalOauthScopeDelimiter(".").
		WithExternalOauthScopeMappingAttribute("foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, integrationModel),
				Check: assertThat(t,
					resourceassert.ExternalOauthSecurityIntegrationResource(t, "snowflake_external_oauth_integration.test").
						HasComment("foo").
						HasEnabled(true).
						HasExternalOauthBlockedRolesList(role.ID().Name()).
						HasExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable)).
						HasExternalOauthAudienceList("foo").
						HasExternalOauthIssuer(issuer).
						HasExternalOauthRsaPublicKey(rsaKey).
						HasExternalOauthScopeDelimiter(".").
						HasExternalOauthScopeMappingAttribute("foo").
						HasExternalOauthSnowflakeUserMappingAttribute(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress)).
						HasExternalOauthTokenUserMappingClaim("foo").
						HasName(id.Name()).
						HasExternalOauthType(string(sdk.ExternalOauthSecurityIntegrationTypeOptionCustom)),
					resourceshowoutputassert.SecurityIntegrationShowOutput(t, "snowflake_external_oauth_integration.test").
						HasName(id.Name()).
						HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment("foo").
						HasCreatedOnNotEmpty(),
					resourceshowoutputassert.ExternalOauthSecurityIntegrationDescOutput(t, "snowflake_external_oauth_integration.test").
						HasEnabled("true").
						HasExternalOauthIssuer(issuer).
						HasExternalOauthJwsKeysUrl("").
						HasExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable)).
						HasExternalOauthRsaPublicKey(rsaKey).
						HasExternalOauthRsaPublicKey2(rsaKey).
						HasExternalOauthBlockedRolesList(strings.Join(expectedRoles, ",")).
						HasExternalOauthAllowedRolesList("").
						HasExternalOauthAudienceList("foo").
						HasExternalOauthTokenUserMappingClaim("['foo']").
						HasExternalOauthSnowflakeUserMappingAttribute(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress)).
						HasExternalOauthScopeDelimiter(".").
						HasComment("foo"),
				),
			},
			{
				Config:        accconfig.FromModels(t, integrationModel),
				ResourceName:  "snowflake_external_oauth_integration.test",
				ImportState:   true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedExternalOauthSecurityIntegrationResource(t, resourcehelpers.EncodeResourceIdentifier(id)).
						HasCommentString("foo").
						HasEnabledString("true").
						HasExternalOauthBlockedRolesList(expectedRoles[0], expectedRoles[1], expectedRoles[2]).
						HasExternalOauthAnyRoleModeString(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable)).
						HasExternalOauthAudienceList("foo").
						HasExternalOauthIssuerString(issuer).
						HasExternalOauthRsaPublicKeyString(rsaKey).
						HasExternalOauthScopeDelimiterString(".").
						HasExternalOauthSnowflakeUserMappingAttributeString(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress)).
						HasExternalOauthTokenUserMappingClaim("foo").
						HasNameString(id.Name()).
						HasExternalOauthTypeString(string(sdk.ExternalOauthSecurityIntegrationTypeOptionCustom)),
				),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_completeWithRsaPublicKeysAndBlockedRolesList_paramUnset(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	paramCleanup := testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList, "false")
	t.Cleanup(paramCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := random.String()
	rsaKey, _ := random.GenerateRSAPublicKey(t)

	integrationModel := model.ExternalOauthSecurityIntegration("test", id.Name(), true, issuer,
		string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress),
		[]string{"foo"},
		string(sdk.ExternalOauthSecurityIntegrationTypeOptionCustom),
	).
		WithComment("foo").
		WithExternalOauthBlockedRoles(role.ID()).
		WithExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable)).
		WithExternalOauthAudiences("foo").
		WithExternalOauthRsaPublicKeyValue(accconfig.MultilineWrapperVariable(rsaKey)).
		WithExternalOauthRsaPublicKey2Value(accconfig.MultilineWrapperVariable(rsaKey)).
		WithExternalOauthScopeDelimiter(".").
		WithExternalOauthScopeMappingAttribute("foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, integrationModel),
				Check: assertThat(t,
					resourceassert.ExternalOauthSecurityIntegrationResource(t, "snowflake_external_oauth_integration.test").
						HasComment("foo").
						HasEnabled(true).
						HasExternalOauthBlockedRolesList(role.ID().Name()).
						HasExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable)).
						HasExternalOauthAudienceList("foo").
						HasExternalOauthIssuer(issuer).
						HasExternalOauthRsaPublicKey(rsaKey).
						HasExternalOauthScopeDelimiter(".").
						HasExternalOauthScopeMappingAttribute("foo").
						HasExternalOauthSnowflakeUserMappingAttribute(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress)).
						HasExternalOauthTokenUserMappingClaim("foo").
						HasName(id.Name()).
						HasExternalOauthType(string(sdk.ExternalOauthSecurityIntegrationTypeOptionCustom)),
					resourceshowoutputassert.SecurityIntegrationShowOutput(t, "snowflake_external_oauth_integration.test").
						HasName(id.Name()).
						HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment("foo").
						HasCreatedOnNotEmpty(),
					resourceshowoutputassert.ExternalOauthSecurityIntegrationDescOutput(t, "snowflake_external_oauth_integration.test").
						HasEnabled("true").
						HasExternalOauthIssuer(issuer).
						HasExternalOauthJwsKeysUrl("").
						HasExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeOptionDisable)).
						HasExternalOauthRsaPublicKey(rsaKey).
						HasExternalOauthRsaPublicKey2(rsaKey).
						HasExternalOauthBlockedRolesList(role.ID().Name()).
						HasExternalOauthAllowedRolesList("").
						HasExternalOauthAudienceList("foo").
						HasExternalOauthTokenUserMappingClaim("['foo']").
						HasExternalOauthSnowflakeUserMappingAttribute(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOptionEmailAddress)).
						HasExternalOauthScopeDelimiter(".").
						HasComment("foo"),
				),
			},
			{
				Config:                  accconfig.FromModels(t, integrationModel),
				ResourceName:            "snowflake_external_oauth_integration.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2", "external_oauth_scope_mapping_attribute"},
			},
		},
	})
}
