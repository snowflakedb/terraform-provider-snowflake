//go:build account_level_tests

package testacc

import (
	"sort"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
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

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable(role.ID().Name())),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(issuer),
			"external_oauth_rsa_public_key":                   config.StringVariable(rsaKey),
			"external_oauth_rsa_public_key_2":                 config.StringVariable(rsaKey),
			"external_oauth_scope_delimiter":                  config.StringVariable("."),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable(id.Name()),
			"external_oauth_type": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				Check: assertThat(t,
					resourceassert.ExternalOauthSecurityIntegrationResource(t, "snowflake_external_oauth_integration.test").
						HasComment("foo").
						HasEnabled(true).
						HasExternalOauthBlockedRolesList(role.ID().Name()).
						HasExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
						HasExternalOauthAudienceList("foo").
						HasExternalOauthIssuer(issuer).
						HasExternalOauthRsaPublicKey(rsaKey).
						HasExternalOauthScopeDelimiter(".").
						HasExternalOauthScopeMappingAttribute("foo").
						HasExternalOauthSnowflakeUserMappingAttribute(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)).
						HasExternalOauthTokenUserMappingClaim("foo").
						HasName(id.Name()).
						HasExternalOauthType(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
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
						HasExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
						HasExternalOauthRsaPublicKey(rsaKey).
						HasExternalOauthRsaPublicKey2(rsaKey).
						HasExternalOauthBlockedRolesList(strings.Join(expectedRoles, ",")).
						HasExternalOauthAllowedRolesList("").
						HasExternalOauthAudienceList("foo").
						HasExternalOauthTokenUserMappingClaim("['foo']").
						HasExternalOauthSnowflakeUserMappingAttribute(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)).
						HasExternalOauthScopeDelimiter(".").
						HasComment("foo"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				ResourceName:    "snowflake_external_oauth_integration.test",
				ImportState:     true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedExternalOauthSecurityIntegrationResource(t, resourcehelpers.EncodeResourceIdentifier(id)).
						HasCommentString("foo").
						HasEnabledString("true").
						HasExternalOauthBlockedRolesList(expectedRoles[0], expectedRoles[1], expectedRoles[2]).
						HasExternalOauthAnyRoleModeString(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
						HasExternalOauthAudienceList("foo").
						HasExternalOauthIssuerString(issuer).
						HasExternalOauthRsaPublicKeyString(rsaKey).
						HasExternalOauthScopeDelimiterString(".").
						HasExternalOauthSnowflakeUserMappingAttributeString(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)).
						HasExternalOauthTokenUserMappingClaim("foo").
						HasNameString(id.Name()).
						HasExternalOauthTypeString(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
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

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable(role.ID().Name())),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(issuer),
			"external_oauth_rsa_public_key":                   config.StringVariable(rsaKey),
			"external_oauth_rsa_public_key_2":                 config.StringVariable(rsaKey),
			"external_oauth_scope_delimiter":                  config.StringVariable("."),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable(id.Name()),
			"external_oauth_type": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				Check: assertThat(t,
					resourceassert.ExternalOauthSecurityIntegrationResource(t, "snowflake_external_oauth_integration.test").
						HasComment("foo").
						HasEnabled(true).
						HasExternalOauthBlockedRolesList(role.ID().Name()).
						HasExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
						HasExternalOauthAudienceList("foo").
						HasExternalOauthIssuer(issuer).
						HasExternalOauthRsaPublicKey(rsaKey).
						HasExternalOauthScopeDelimiter(".").
						HasExternalOauthScopeMappingAttribute("foo").
						HasExternalOauthSnowflakeUserMappingAttribute(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)).
						HasExternalOauthTokenUserMappingClaim("foo").
						HasName(id.Name()).
						HasExternalOauthType(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
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
						HasExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
						HasExternalOauthRsaPublicKey(rsaKey).
						HasExternalOauthRsaPublicKey2(rsaKey).
						HasExternalOauthBlockedRolesList(role.ID().Name()).
						HasExternalOauthAllowedRolesList("").
						HasExternalOauthAudienceList("foo").
						HasExternalOauthTokenUserMappingClaim("['foo']").
						HasExternalOauthSnowflakeUserMappingAttribute(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)).
						HasExternalOauthScopeDelimiter(".").
						HasComment("foo"),
				),
			},
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables:         m(),
				ResourceName:            "snowflake_external_oauth_integration.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2", "external_oauth_scope_mapping_attribute"},
			},
		},
	})
}
