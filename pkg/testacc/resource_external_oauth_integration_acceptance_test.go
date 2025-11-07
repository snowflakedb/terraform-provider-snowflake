//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalOauthIntegration_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	issuer := random.String()

	basic := model.ExternalOauthSecurityIntegration("test", id.Name(), true, issuer,
		string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress),
		[]string{"foo"},
		string(sdk.ExternalOauthSecurityIntegrationTypeCustom),
	).
		WithExternalOauthJwsKeysUrlValue(config.SetVariable(config.StringVariable("https://example.com")))

	complete := model.ExternalOauthSecurityIntegration("test", id.Name(), true, issuer,
		string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress),
		[]string{"foo"},
		string(sdk.ExternalOauthSecurityIntegrationTypeCustom),
	).
		WithExternalOauthJwsKeysUrlValue(config.SetVariable(config.StringVariable("https://example.com"))).
		WithExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
		WithExternalOauthAudienceListValue(config.SetVariable(config.StringVariable("bar"))).
		WithExternalOauthScopeDelimiter(".").
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.SecurityIntegration(t, id).
			HasName(id.Name()).
			HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
			HasEnabled(true).
			HasComment(""),

		resourceassert.ExternalOauthSecurityIntegrationResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasEnabledString(r.BooleanTrue).
			HasExternalOauthTypeString(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)).
			HasExternalOauthIssuerString(issuer).
			HasExternalOauthSnowflakeUserMappingAttributeString(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)).
			HasCommentString(""),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, basic.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
			HasEnabled(true).
			HasComment(""),

		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.enabled.0.value", r.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.external_oauth_issuer.0.value", issuer)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress))),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable))),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.external_oauth_scope_delimiter.0.value", ",")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.comment.0.value", "")),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.SecurityIntegration(t, id).
			HasName(id.Name()).
			HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
			HasEnabled(true).
			HasComment(comment),

		resourceassert.ExternalOauthSecurityIntegrationResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasEnabledString(r.BooleanTrue).
			HasExternalOauthTypeString(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)).
			HasExternalOauthIssuerString(issuer).
			HasExternalOauthSnowflakeUserMappingAttributeString(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)).
			HasExternalOauthAnyRoleModeString(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
			HasExternalOauthScopeDelimiterString(".").
			HasCommentString(comment),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
			HasEnabled(true).
			HasComment(comment),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", r.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_issuer.0.value", issuer)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress))),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable))),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_audience_list.0.value", "bar")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_scope_delimiter.0.value", ".")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"external_oauth_jws_keys_url",
					"external_oauth_rsa_public_key",
					"external_oauth_rsa_public_key_2",
					"external_oauth_scope_mapping_attribute",
					"external_oauth_any_role_mode",
					"external_oauth_scope_delimiter",
					"external_oauth_blocked_roles_list",
				},
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"external_oauth_jws_keys_url",
					"external_oauth_rsa_public_key",
					"external_oauth_rsa_public_key_2",
					"external_oauth_scope_mapping_attribute",
					"external_oauth_any_role_mode",
					"external_oauth_scope_delimiter",
					"external_oauth_blocked_roles_list",
				},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().SecurityIntegration.UpdateExternalOauth(t, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).
						WithSet(*sdk.NewExternalOauthIntegrationSetRequest().
							WithEnabled(true).
							WithExternalOauthJwsKeysUrl([]sdk.JwsKeysUrl{{"https://example.com"}}).
							WithExternalOauthAnyRoleMode(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable).
							WithExternalOauthAudienceList(*sdk.NewAudienceListRequest([]sdk.AudienceListItem{{"bar"}})).
							WithExternalOauthSnowflakeUserMappingAttribute(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName).
							WithComment(sdk.StringAllowEmpty{Value: comment}),
						))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy - ensure external oauth integration is deleted
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, basic),
				Check: assertThat(t,
					invokeactionassert.SecurityIntegrationDoesNotExist(t, id),
				),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_CompleteUseCase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	issuer := random.String()

	complete := model.ExternalOauthSecurityIntegration("test", id.Name(), true, issuer,
		string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress),
		[]string{"foo"},
		string(sdk.ExternalOauthSecurityIntegrationTypeCustom),
	).
		WithExternalOauthJwsKeysUrlValue(config.SetVariable(config.StringVariable("https://example.com"))).
		WithExternalOauthScopeMappingAttribute("scp").
		WithExternalOauthAllowedRolesListValue(config.SetVariable(config.StringVariable(role.ID().Name()))).
		WithExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
		WithExternalOauthAudienceListValue(config.SetVariable(config.StringVariable("bar"))).
		WithExternalOauthScopeDelimiter(".").
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			// Create - with all optionals (including optional force-new fields)
			{
				Config: accconfig.FromModels(t, complete),
				Check: assertThat(t,
					objectassert.SecurityIntegration(t, id).
						HasName(id.Name()).
						HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
						HasEnabled(true).
						HasComment(comment),

					resourceassert.ExternalOauthSecurityIntegrationResource(t, complete.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasEnabledString(r.BooleanTrue).
						HasExternalOauthTypeString(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)).
						HasExternalOauthIssuerString(issuer).
						HasExternalOauthSnowflakeUserMappingAttributeString(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)).
						HasExternalOauthScopeMappingAttributeString("scp").
						HasExternalOauthAnyRoleModeString(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
						HasExternalOauthScopeDelimiterString(".").
						HasCommentString(comment),

					resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
						HasName(id.Name()).
						HasIntegrationType("EXTERNAL_OAUTH - CUSTOM").
						HasEnabled(true).
						HasComment(comment),

					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", r.BooleanTrue)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_issuer.0.value", issuer)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress))),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable))),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_allowed_roles_list.0.value", role.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_audience_list.0.value", "bar")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_oauth_scope_delimiter.0.value", ".")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
				),
			},
			// Import - with all optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"external_oauth_jws_keys_url",
					"external_oauth_rsa_public_key",
					"external_oauth_rsa_public_key_2",
					"external_oauth_scope_mapping_attribute",
					"external_oauth_any_role_mode",
					"external_oauth_scope_delimiter",
					"external_oauth_allowed_roles_list",
				},
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_invalidAnyRoleMode(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_allowed_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_any_role_mode":                    config.StringVariable("invalid"),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(random.String()),
			"external_oauth_jws_keys_url":                     config.SetVariable(config.StringVariable("foo")),
			"external_oauth_rsa_public_key":                   config.StringVariable("foo"),
			"external_oauth_rsa_public_key_2":                 config.StringVariable("foo"),
			"external_oauth_scope_delimiter":                  config.StringVariable("foo"),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable("foo"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(),
				ExpectError:     regexp.MustCompile("Error: invalid ExternalOauthSecurityIntegrationAnyRoleModeOption: INVALID"),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_invalidSnowflakeUserMappingAttribute(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_allowed_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(random.String()),
			"external_oauth_jws_keys_url":                     config.SetVariable(config.StringVariable("foo")),
			"external_oauth_rsa_public_key":                   config.StringVariable("foo"),
			"external_oauth_rsa_public_key_2":                 config.StringVariable("foo"),
			"external_oauth_scope_delimiter":                  config.StringVariable("foo"),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable("invalid"),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable("foo"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(),
				ExpectError:     regexp.MustCompile("Error: invalid ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption: INVALID"),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_invalidOauthType(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_allowed_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(random.String()),
			"external_oauth_jws_keys_url":                     config.SetVariable(config.StringVariable("foo")),
			"external_oauth_rsa_public_key":                   config.StringVariable("foo"),
			"external_oauth_rsa_public_key_2":                 config.StringVariable("foo"),
			"external_oauth_scope_delimiter":                  config.StringVariable("foo"),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable("foo"),
			"external_oauth_type": config.StringVariable("invalid"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(),
				ExpectError:     regexp.MustCompile("Error: invalid ExternalOauthSecurityIntegrationTypeOption: INVALID"),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_InvalidIncomplete(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name": config.StringVariable(id.Name()),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			// Some strings are trimmed because of inconsistent '\n' placement from tf error messages.
			`The argument "external_oauth_type" is required, but no definition was found.`,
			`The argument "external_oauth_snowflake_user_mapping_attribute" is required,`,
			`The argument "enabled" is required, but no definition was found.`,
			`The argument "external_oauth_issuer" is required,`,
			`The argument "external_oauth_token_user_mapping_claim" is required,`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalOauthIntegration/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_migrateFromVersion092_withRsaPublicKeysAndBlockedRolesList(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := random.String()
	rsaKey, _ := random.GenerateRSAPublicKey(t)

	resourceName := "snowflake_external_oauth_integration.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.92.0"),
				Config:            externalOauthIntegrationWithRsaPublicKeysAndBlockedRolesListv092(id.Name(), issuer, rsaKey, role.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttr(resourceName, "issuer", issuer),
					resource.TestCheckResourceAttr(resourceName, "token_user_mapping_claims.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_user_mapping_claims.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)),
					resource.TestCheckResourceAttr(resourceName, "scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr(resourceName, "rsa_public_key", rsaKey),
					resource.TestCheckResourceAttr(resourceName, "rsa_public_key_2", rsaKey),
					resource.TestCheckResourceAttr(resourceName, "blocked_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "blocked_roles.0", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "audience_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "audience_urls.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(resourceName, "scope_delimiter", ":"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   externalOauthIntegrationWithRsaPublicKeysAndBlockedRolesListv093(id.Name(), issuer, rsaKey, role.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_rsa_public_key", rsaKey),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_rsa_public_key_2", rsaKey),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_blocked_roles_list.0", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_scope_delimiter", ":"),
				),
			},
		},
	})
}

func externalOauthIntegrationWithRsaPublicKeysAndBlockedRolesListv092(name, issuer, rsaKey, roleName string) string {
	s := `
locals {
  key_raw = <<-EOT
%s
  EOT
  key = trimsuffix(local.key_raw, "\n")
}
resource "snowflake_external_oauth_integration" "test" {
	name                             = "%s"
	enabled = true
	type = "CUSTOM"
	issuer = "%s"
	token_user_mapping_claims = ["foo"]
	snowflake_user_mapping_attribute = "LOGIN_NAME"
	scope_mapping_attribute = "foo"
	rsa_public_key = local.key
	rsa_public_key_2 = local.key
	blocked_roles = ["%s"]
	audience_urls = ["foo"]
	any_role_mode = "DISABLE"
	scope_delimiter = ":"
}`
	return fmt.Sprintf(s, rsaKey, name, issuer, roleName)
}

func externalOauthIntegrationWithRsaPublicKeysAndBlockedRolesListv093(name, issuer, rsaKey, roleName string) string {
	s := `
locals {
  key_raw = <<-EOT
%s
  EOT
  key = trimsuffix(local.key_raw, "\n")
}
resource "snowflake_external_oauth_integration" "test" {
	name                             = "%s"
	enabled = true
	external_oauth_type = "CUSTOM"
	external_oauth_issuer = "%s"
	external_oauth_token_user_mapping_claim = ["foo"]
	external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
	external_oauth_scope_mapping_attribute = "foo"
	external_oauth_rsa_public_key = local.key
	external_oauth_rsa_public_key_2 = local.key
	external_oauth_blocked_roles_list = ["%s"]
	external_oauth_audience_list = ["foo"]
	external_oauth_any_role_mode = "DISABLE"
	external_oauth_scope_delimiter = ":"
}`
	return fmt.Sprintf(s, rsaKey, name, issuer, roleName)
}

func TestAcc_ExternalOauthIntegration_migrateFromVersion092_withJwsKeysUrlAndAllowedRolesList(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := random.String()

	resourceName := "snowflake_external_oauth_integration.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.92.0"),
				Config:            externalOauthIntegrationWithJwsKeysUrlAndAllowedRolesListv092(id.Name(), issuer, role.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttr(resourceName, "issuer", issuer),
					resource.TestCheckResourceAttr(resourceName, "token_user_mapping_claims.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_user_mapping_claims.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)),
					resource.TestCheckResourceAttr(resourceName, "scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr(resourceName, "jws_keys_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "jws_keys_urls.0", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "allowed_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allowed_roles.0", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "audience_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "audience_urls.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(resourceName, "scope_delimiter", ":"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   externalOauthIntegrationWithJwsKeysUrlAndAllowedRolesListv093(id.Name(), issuer, role.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_jws_keys_url.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_jws_keys_url.0", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_allowed_roles_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_allowed_roles_list.0", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_scope_delimiter", ":"),
				),
			},
		},
	})
}

func externalOauthIntegrationWithJwsKeysUrlAndAllowedRolesListv092(name, issuer, roleName string) string {
	s := `
resource "snowflake_external_oauth_integration" "test" {
	name                             = "%s"
	enabled = true
	type = "CUSTOM"
	issuer = "%s"
	token_user_mapping_claims = ["foo"]
	snowflake_user_mapping_attribute = "LOGIN_NAME"
	scope_mapping_attribute = "foo"
	jws_keys_urls = ["https://example.com"]
	allowed_roles = ["%s"]
	audience_urls = ["foo"]
	any_role_mode = "DISABLE"
	scope_delimiter = ":"
}`
	return fmt.Sprintf(s, name, issuer, roleName)
}

func externalOauthIntegrationWithJwsKeysUrlAndAllowedRolesListv093(name, issuer, roleName string) string {
	s := `
resource "snowflake_external_oauth_integration" "test" {
	name                             = "%s"
	enabled = true
	external_oauth_type = "CUSTOM"
	external_oauth_issuer = "%s"
	external_oauth_token_user_mapping_claim = ["foo"]
	external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
	external_oauth_scope_mapping_attribute = "foo"
	external_oauth_jws_keys_url = ["https://example.com"]
	external_oauth_allowed_roles_list = ["%s"]
	external_oauth_audience_list = ["foo"]
	external_oauth_any_role_mode = "DISABLE"
	external_oauth_scope_delimiter = ":"
}`
	return fmt.Sprintf(s, name, issuer, roleName)
}

func TestAcc_ExternalOauthIntegration_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := random.String()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            externalOauthIntegrationBasicConfig(id.Name(), issuer),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   externalOauthIntegrationBasicConfig(id.Name(), issuer),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_WithQuotedName(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())
	issuer := random.String()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             externalOauthIntegrationBasicConfig(quotedId, issuer),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   externalOauthIntegrationBasicConfig(quotedId, issuer),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_oauth_integration.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_oauth_integration.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "id", id.Name()),
				),
			},
		},
	})
}

func externalOauthIntegrationBasicConfig(name string, issuer string) string {
	return fmt.Sprintf(`
resource "snowflake_external_oauth_integration" "test" {
	name               = "%[1]s"
	external_oauth_type                             = "CUSTOM"
	enabled                                         = true
	external_oauth_issuer                           = "%[2]s"
	external_oauth_token_user_mapping_claim         = [ "foo" ]
	external_oauth_snowflake_user_mapping_attribute = "EMAIL_ADDRESS"
	external_oauth_jws_keys_url                     = [ "https://example.com" ]
}
`, name, issuer)
}
