//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalAccessIntegration_BasicUseCase(t *testing.T) {
	networkRule1Id := testClient().Ids.RandomSchemaObjectIdentifier()
	_, networkRule1Cleanup := testClient().NetworkRule.CreateEgressWithIdentifier(t, networkRule1Id)
	t.Cleanup(networkRule1Cleanup)

	networkRule2Id := testClient().Ids.RandomSchemaObjectIdentifier()
	_, networkRule2Cleanup := testClient().NetworkRule.CreateEgressWithIdentifier(t, networkRule2Id)
	t.Cleanup(networkRule2Cleanup)

	secret1Id, secret1Cleanup := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(secret1Cleanup)

	secret2Id, secret2Cleanup := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(secret2Cleanup)

	apiAuthIntegration, apiAuthCleanup := testClient().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlow(t)
	t.Cleanup(apiAuthCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id}, true)

	complete := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id}, true).
		WithComment(comment).
		WithAllowedAuthenticationSecretsSecrets([]string{secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()}).
		WithAllowedApiAuthenticationIntegrationsIntegrations([]string{apiAuthIntegration.ID().Name()})

	disabled := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id}, false).
		WithComment(comment).
		WithAllowedAuthenticationSecretsSecrets([]string{secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()}).
		WithAllowedApiAuthenticationIntegrationsIntegrations([]string{apiAuthIntegration.ID().Name()})

	withTwoRules := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id, networkRule2Id}, true).
		WithComment(comment).
		WithAllowedAuthenticationSecretsSecrets([]string{secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()}).
		WithAllowedApiAuthenticationIntegrationsIntegrations([]string{apiAuthIntegration.ID().Name()})

	allSecrets := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id, networkRule2Id}, true).
		WithComment(comment).
		WithAllowedAuthenticationSecretsAll().
		WithAllowedApiAuthenticationIntegrationsIntegrations([]string{apiAuthIntegration.ID().Name()})

	noSecrets := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id, networkRule2Id}, true).
		WithComment(comment).
		WithAllowedAuthenticationSecretsNone().
		WithAllowedApiAuthenticationIntegrationsIntegrations([]string{apiAuthIntegration.ID().Name()})

	unsetSecrets := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id, networkRule2Id}, true).
		WithComment(comment).
		WithAllowedApiAuthenticationIntegrationsIntegrations([]string{apiAuthIntegration.ID().Name()})

	noApiAuth := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id, networkRule2Id}, true).
		WithComment(comment).
		WithAllowedApiAuthenticationIntegrationsNone()

	unsetApiAuth := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRule1Id, networkRule2Id}, true).
		WithComment(comment)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasCommentEmpty().
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsNotSet().
			HasAllowedApiAuthenticationIntegrationsNotSet().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasId(id).
			HasEnabled(true).
			HasComment("").
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName()).
			HasNoAllowedAuthenticationSecrets().
			HasNoAllowedApiAuthenticationIntegrations(),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsSecrets(secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()).
			HasAllowedApiAuthenticationIntegrationsIntegrations(apiAuthIntegration.ID().Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasId(id).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecrets(secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()).
			HasAllowedApiAuthenticationIntegrations(apiAuthIntegration.ID().Name()),
	}

	disabledAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasEnabled(false).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsSecrets(secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()).
			HasAllowedApiAuthenticationIntegrationsIntegrations(apiAuthIntegration.ID().Name()),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasEnabled(false).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecrets(secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()).
			HasAllowedApiAuthenticationIntegrations(apiAuthIntegration.ID().Name()),
	}

	withTwoRulesAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsSecrets(secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()).
			HasAllowedApiAuthenticationIntegrationsIntegrations(apiAuthIntegration.ID().Name()),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecrets(secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()).
			HasAllowedApiAuthenticationIntegrations(apiAuthIntegration.ID().Name()),
	}

	allSecretsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsAll().
			HasAllowedApiAuthenticationIntegrationsIntegrations(apiAuthIntegration.ID().Name()),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecrets("ALL").
			HasAllowedApiAuthenticationIntegrations(apiAuthIntegration.ID().Name()),
	}

	noSecretsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsNone().
			HasAllowedApiAuthenticationIntegrationsIntegrations(apiAuthIntegration.ID().Name()),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasNoAllowedAuthenticationSecrets().
			HasAllowedApiAuthenticationIntegrations(apiAuthIntegration.ID().Name()),
	}

	unsetSecretsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsNotSet().
			HasAllowedApiAuthenticationIntegrationsIntegrations(apiAuthIntegration.ID().Name()),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasNoAllowedAuthenticationSecrets().
			HasAllowedApiAuthenticationIntegrations(apiAuthIntegration.ID().Name()),
	}

	noApiAuthAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsNotSet().
			HasAllowedApiAuthenticationIntegrationsNone(),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasNoAllowedAuthenticationSecrets().
			HasNoAllowedApiAuthenticationIntegrations(),
	}

	unsetApiAuthAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsNotSet().
			HasAllowedApiAuthenticationIntegrationsNotSet(),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRule1Id.FullyQualifiedName(), networkRule2Id.FullyQualifiedName()).
			HasNoAllowedAuthenticationSecrets().
			HasNoAllowedApiAuthenticationIntegrations(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalAccessIntegration),
		Steps: []resource.TestStep{
			// Create with required fields only
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import with required fields only
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Set all optional fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Disable
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, disabled),
				Check:  assertThat(t, disabledAssertions...),
			},
			// External change: alter all fields in Snowflake, expect drift detected and corrected
			{
				PreConfig: func() {
					testClient().ExternalAccessIntegration.Alter(
						t,
						sdk.NewAlterExternalAccessIntegrationRequest(id).WithSet(
							*sdk.NewExternalAccessIntegrationSetRequest().
								WithEnabled(true).
								WithComment(externalComment).
								WithAllowedNetworkRules([]sdk.SchemaObjectIdentifier{networkRule2Id}).
								WithAllowedAuthenticationSecrets(
									*sdk.NewExternalAccessIntegrationAllowedAuthenticationSecretsRequest().WithAll(true),
								).
								WithAllowedApiAuthenticationIntegrations(
									*sdk.NewExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest().WithNone(true),
								),
						),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "enabled", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(ref, "enabled", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						planchecks.ExpectDrift(ref, "allowed_authentication_secrets.0.all", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(ref, "allowed_authentication_secrets.0.all", tfjson.ActionUpdate, sdk.String("true"), nil),
						planchecks.ExpectDrift(ref, "allowed_api_authentication_integrations.0.none", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(ref, "allowed_api_authentication_integrations.0.none", tfjson.ActionUpdate, sdk.String("true"), nil),
					},
				},
				Config: accconfig.FromModels(t, disabled),
				Check:  assertThat(t, disabledAssertions...),
			},
			// Add second network rule, re-enable
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, withTwoRules),
				Check:  assertThat(t, withTwoRulesAssertions...),
			},
			// Auth secrets: list → all
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, allSecrets),
				Check:  assertThat(t, allSecretsAssertions...),
			},
			// Auth secrets: all → none
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, noSecrets),
				Check:  assertThat(t, noSecretsAssertions...),
			},
			// Auth secrets: none → unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, unsetSecrets),
				Check:  assertThat(t, unsetSecretsAssertions...),
			},
			// API auth: list → none
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, noApiAuth),
				Check:  assertThat(t, noApiAuthAssertions...),
			},
			// API auth: none → unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, unsetApiAuth),
				Check:  assertThat(t, unsetApiAuthAssertions...),
			},
			// Import with fields unset
			{
				Config:            accconfig.FromModels(t, unsetApiAuth),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ExternalAccessIntegration_CompleteUseCase(t *testing.T) {
	networkRuleId := testClient().Ids.RandomSchemaObjectIdentifier()
	_, networkRuleCleanup := testClient().NetworkRule.CreateEgressWithIdentifier(t, networkRuleId)
	t.Cleanup(networkRuleCleanup)

	secret1Id, secret1Cleanup := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(secret1Cleanup)

	secret2Id, secret2Cleanup := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(secret2Cleanup)

	apiAuthIntegration, apiAuthCleanup := testClient().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlow(t)
	t.Cleanup(apiAuthCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	allAttributes := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRuleId}, true).
		WithComment(comment).
		WithAllowedAuthenticationSecretsSecrets([]string{secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()}).
		WithAllowedApiAuthenticationIntegrationsIntegrations([]string{apiAuthIntegration.ID().Name()})

	ref := allAttributes.ResourceReference()

	allAttributesAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ExternalAccessIntegrationResource(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRuleId.FullyQualifiedName()).
			HasAllowedAuthenticationSecretsSecrets(secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()).
			HasAllowedApiAuthenticationIntegrationsIntegrations(apiAuthIntegration.ID().Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.ExternalAccessIntegrationShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ExternalAccessIntegrationDescribeOutput(t, ref).
			HasId(id).
			HasEnabled(true).
			HasComment(comment).
			HasAllowedNetworkRules(networkRuleId.FullyQualifiedName()).
			HasAllowedAuthenticationSecrets(secret1Id.FullyQualifiedName(), secret2Id.FullyQualifiedName()).
			HasAllowedApiAuthenticationIntegrations(apiAuthIntegration.ID().Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalAccessIntegration),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, allAttributes),
				Check:  assertThat(t, allAttributesAssertions...),
			},
			{
				Config:            accconfig.FromModels(t, allAttributes),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"allowed_authentication_secrets",
					"allowed_api_authentication_integrations",
				},
			},
		},
	})
}

func TestAcc_ExternalAccessIntegration_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	networkRuleId := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidNetworkRule := model.ExternalAccessIntegration("t", id.Name(), []string{"invalid_network_rule"}, true)

	invalidSecret := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRuleId}, true).
		WithAllowedAuthenticationSecretsSecrets([]string{"invalid_secret"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalAccessIntegration),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, invalidNetworkRule),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid identifier type`),
			},
			{
				Config:      accconfig.FromModels(t, invalidSecret),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid identifier type`),
			},
		},
	})
}
