//go:build account_level_tests

package testacc

import (
	"sort"
	"strings"
	"testing"

	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.0", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_rsa_public_key", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_rsa_public_key.0.value", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_rsa_public_key_2.0.value", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_blocked_roles_list.0.value", strings.Join(expectedRoles, ",")),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				ResourceName:    "snowflake_external_oauth_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_blocked_roles_list.#", "3"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_blocked_roles_list.0", expectedRoles[0]),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_blocked_roles_list.1", expectedRoles[1]),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_blocked_roles_list.2", expectedRoles[2]),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_audience_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_audience_list.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_rsa_public_key", rsaKey),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_scope_delimiter", "."),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_token_user_mapping_claim.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_token_user_mapping_claim.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.0", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_rsa_public_key", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_rsa_public_key.0.value", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_rsa_public_key_2.0.value", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_blocked_roles_list.0.value", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo"),
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
