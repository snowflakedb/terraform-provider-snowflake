//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ScimIntegration_basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	role, role2 := snowflakeroles.GenericScimProvisioner, snowflakeroles.OktaProvisioner
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	scimModelBasic := model.ScimSecurityIntegration("test", id.Name(), false, role.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric))
	scimModelOktaFull := model.ScimSecurityIntegration("test", id.Name(), true, role2.Name(), string(sdk.ScimSecurityIntegrationScimClientOkta)).
		WithSyncPassword(r.BooleanFalse).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithComment(comment)
	scimModelOkta := model.ScimSecurityIntegration("test", id.Name(), true, role2.Name(), string(sdk.ScimSecurityIntegrationScimClientOkta))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, scimModelBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "scim_client", "GENERIC"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "run_as_role", role.Name()),
					resource.TestCheckNoResourceAttr(scimModelBasic.ResourceReference(), "network_policy"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "sync_password", r.BooleanDefault),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(scimModelBasic.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "describe_output.0.network_policy.0.value", ""),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "describe_output.0.run_as_role.0.value", role.Name()),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "describe_output.0.sync_password.0.value", "false"),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "describe_output.0.comment.0.value", ""),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, scimModelBasic),
				ResourceName: scimModelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "scim_client", "GENERIC"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "run_as_role", role.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "sync_password", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, scimModelOktaFull),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "scim_client", "OKTA"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "run_as_role", role2.Name()),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "network_policy", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "sync_password", "false"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "show_output.0.integration_type", "SCIM - OKTA"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(scimModelOktaFull.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "describe_output.0.network_policy.0.value", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "describe_output.0.run_as_role.0.value", role2.Name()),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "describe_output.0.sync_password.0.value", "false"),
					resource.TestCheckResourceAttr(scimModelOktaFull.ResourceReference(), "describe_output.0.comment.0.value", comment),
				),
			},
			// import - complete
			{
				Config:       accconfig.FromModels(t, scimModelOktaFull),
				ResourceName: scimModelOktaFull.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "scim_client", "OKTA"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "run_as_role", role2.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", networkPolicy.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "sync_password", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, scimModelOkta),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(scimModelOkta.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(scimModelOkta.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(scimModelOkta.ResourceReference(), "scim_client", "OKTA"),
					resource.TestCheckResourceAttr(scimModelOkta.ResourceReference(), "run_as_role", role2.Name()),
					resource.TestCheckResourceAttr(scimModelOkta.ResourceReference(), "network_policy", ""),
					resource.TestCheckResourceAttr(scimModelOkta.ResourceReference(), "sync_password", r.BooleanDefault),
					resource.TestCheckResourceAttr(scimModelOkta.ResourceReference(), "comment", ""),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner
	comment := random.Comment()

	scimCompleteModel := model.ScimSecurityIntegration("test", id.Name(), false, role.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric)).
		WithSyncPassword(r.BooleanFalse).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, scimCompleteModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(scimCompleteModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(scimCompleteModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(scimCompleteModel.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(scimCompleteModel.ResourceReference(), "scim_client", "GENERIC"),
					resource.TestCheckResourceAttr(scimCompleteModel.ResourceReference(), "run_as_role", role.Name()),
					resource.TestCheckResourceAttr(scimCompleteModel.ResourceReference(), "network_policy", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr(scimCompleteModel.ResourceReference(), "sync_password", "false"),
					resource.TestCheckResourceAttr(scimCompleteModel.ResourceReference(), "comment", comment),
				),
			},
			{
				Config:            accconfig.FromModels(t, scimCompleteModel),
				ResourceName:      scimCompleteModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ScimIntegration_completeAzure(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner
	comment := random.Comment()

	scimCompleteAzureModel := model.ScimSecurityIntegration("test", id.Name(), false, role.Name(), string(sdk.ScimSecurityIntegrationScimClientAzure)).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, scimCompleteAzureModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(scimCompleteAzureModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(scimCompleteAzureModel.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(scimCompleteAzureModel.ResourceReference(), "scim_client", string(sdk.ScimSecurityIntegrationScimClientAzure)),
					resource.TestCheckResourceAttr(scimCompleteAzureModel.ResourceReference(), "run_as_role", role.Name()),
					resource.TestCheckResourceAttr(scimCompleteAzureModel.ResourceReference(), "network_policy", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr(scimCompleteAzureModel.ResourceReference(), "sync_password", r.BooleanDefault),
					resource.TestCheckResourceAttr(scimCompleteAzureModel.ResourceReference(), "comment", comment),
				),
			},
			{
				Config:            accconfig.FromModels(t, scimCompleteAzureModel),
				ResourceName:      scimCompleteAzureModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidScimClient(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	scimModelBasic := model.ScimSecurityIntegration("test", id.Name(), false, snowflakeroles.GenericScimProvisioner.Name(), "invalid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, scimModelBasic),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid ScimSecurityIntegrationScimClientOption: INVALID`),
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidRunAsRole(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	scimModelBasic := model.ScimSecurityIntegration("test", id.Name(), false, "invalid", string(sdk.ScimSecurityIntegrationScimClientGeneric))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, scimModelBasic),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid ScimSecurityIntegrationRunAsRoleOption: INVALID`),
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidIncomplete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			`The argument "scim_client" is required, but no definition was found.`,
			`The argument "run_as_role" is required, but no definition was found.`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/invalid"),
				PlanOnly:        true,
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidCreateWithSyncPasswordWithAzure(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	scimCompleteAzureModel := model.ScimSecurityIntegration("test", id.Name(), false, snowflakeroles.GenericScimProvisioner.Name(), string(sdk.ScimSecurityIntegrationScimClientAzure)).
		WithComment(comment).
		WithSyncPassword(r.BooleanFalse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			"can not CREATE scim integration with field `sync_password`",
		}),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, scimCompleteAzureModel),
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidUpdateWithSyncPasswordWithAzure(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	scimBasicAzureModel := model.ScimSecurityIntegration("test", id.Name(), false, snowflakeroles.GenericScimProvisioner.Name(), string(sdk.ScimSecurityIntegrationScimClientAzure))
	scimCompleteAzureModel := model.ScimSecurityIntegration("test", id.Name(), false, snowflakeroles.GenericScimProvisioner.Name(), string(sdk.ScimSecurityIntegrationScimClientAzure)).
		WithComment(comment).
		WithSyncPassword(r.BooleanFalse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			"can not SET and UNSET field `sync_password`",
		}),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, scimBasicAzureModel),
			},
			{
				Config: accconfig.FromModels(t, scimCompleteAzureModel),
			},
		},
	})
}

func TestAcc_ScimIntegration_migrateFromVersion092EnabledTrue(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner

	resourceName := "snowflake_scim_integration.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.92.0"),
				Config:            scimIntegrationV092(id, role, sdk.ScimSecurityIntegrationScimClientGeneric),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "provisioner_role", role.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationV093(id, role, true, sdk.ScimSecurityIntegrationScimClientGeneric),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "name", tfjson.ActionUpdate, sdk.String(id.Name()), sdk.String(id.Name())),
						planchecks.ExpectChange(resourceName, "enabled", tfjson.ActionUpdate, sdk.String("true"), sdk.String("true")),
						planchecks.ExpectChange(resourceName, "scim_client", tfjson.ActionUpdate, sdk.String("GENERIC"), sdk.String("GENERIC")),
						planchecks.ExpectChange(resourceName, "run_as_role", tfjson.ActionUpdate, sdk.String(role.Name()), sdk.String(role.Name())),
						planchecks.ExpectChange(resourceName, "network_policy", tfjson.ActionUpdate, sdk.String(""), sdk.String("")),
						planchecks.ExpectChange(resourceName, "sync_password", tfjson.ActionUpdate, nil, sdk.String(r.BooleanDefault)),
						planchecks.ExpectChange(resourceName, "comment", tfjson.ActionUpdate, nil, nil),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "run_as_role", role.Name()),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_migrateFromVersion092EnabledFalse(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner

	resourceName := "snowflake_scim_integration.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.92.0"),
				Config:            scimIntegrationV092(id, role, sdk.ScimSecurityIntegrationScimClientGeneric),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "provisioner_role", role.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationV093(id, role, false, sdk.ScimSecurityIntegrationScimClientGeneric),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "run_as_role", role.Name()),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_migrateFromVersion093HandleSyncPassword(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner

	resourceName := "snowflake_scim_integration.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			// create resource with v0.92
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.92.0"),
				Config:            scimIntegrationV092(id, role, sdk.ScimSecurityIntegrationScimClientAzure),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
				),
			},
			// migrate to v0.93 - there is a diff due to new field sync_password in state
			{
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.93.0"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.ExpectChange(resourceName, "sync_password", tfjson.ActionUpdate, nil, sdk.String(r.BooleanDefault)),
					},
				},
				Config: scimIntegrationV093(id, role, true, sdk.ScimSecurityIntegrationScimClientAzure),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
				),
				ExpectError: regexp.MustCompile("invalid property 'SYNC_PASSWORD' for 'INTEGRATION"),
			},
			// check with newest version - the value in state was set to boolean default, so there should be no diff
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationV093(id, role, true, sdk.ScimSecurityIntegrationScimClientAzure),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "sync_password", r.BooleanDefault),
				),
			},
		},
	})
}

func scimIntegrationV092(scimId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, scimClient sdk.ScimSecurityIntegrationScimClientOption) string {
	return fmt.Sprintf(`
resource "snowflake_scim_integration" "test" {
	name             = "%[1]s"
	provisioner_role = "%[2]s"
	scim_client      = "%[3]s"
}
`, scimId.Name(), roleId.Name(), scimClient)
}

func scimIntegrationV093(scimId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, enabled bool, scimClient sdk.ScimSecurityIntegrationScimClientOption) string {
	return fmt.Sprintf(`
resource "snowflake_scim_integration" "test" {
	name             = "%[1]s"
	run_as_role		 = "%[2]s"
	scim_client      = "%[3]s"
	enabled          = %[4]t
}
`, scimId.Name(), roleId.Name(), scimClient, enabled)
}

func TestAcc_ScimIntegration_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	scimModelBasic := model.ScimSecurityIntegration("test", id.Name(), false, snowflakeroles.GenericScimProvisioner.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, scimModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, scimModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_IdentifierQuotingDiffSuppression(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`"%s"`, id.Name())

	scimModelBasic := model.ScimSecurityIntegration("test", quotedId, false, snowflakeroles.GenericScimProvisioner.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, scimModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, scimModelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(scimModelBasic.ResourceReference(), plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(scimModelBasic.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(scimModelBasic.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}
