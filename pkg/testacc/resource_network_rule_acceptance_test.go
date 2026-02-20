//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkRule_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	changedComment := random.Comment()

	values := []string{"192.168.0.100/24", "29.254.123.20"}
	changedValues := []string{"192.168.0.100/24", "29.254.123.20", "172.16.0.0/12"}
	hostPortValues := []string{"example.com", "snowflake.com"}

	modelBasic := model.NetworkRuleFromId(id, sdk.NetworkRuleModeIngress, sdk.NetworkRuleTypeIpv4, values)
	modelAfterUnset := model.NetworkRuleFromId(id, sdk.NetworkRuleModeIngress, sdk.NetworkRuleTypeIpv4, []string{})

	modelComplete := model.NetworkRuleFromId(id, sdk.NetworkRuleModeIngress, sdk.NetworkRuleTypeIpv4, values).
		WithComment(changedComment)

	modelWithAlteredValues := model.NetworkRuleFromId(id, sdk.NetworkRuleModeIngress, sdk.NetworkRuleTypeIpv4, changedValues).
		WithComment(changedComment)

	modelWithChangedForceNewValues := model.NetworkRuleFromId(id, sdk.NetworkRuleModeEgress, sdk.NetworkRuleTypeHostPort, hostPortValues).
		WithComment(changedComment)

	ref := modelBasic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.NetworkRuleResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTypeEnum(sdk.NetworkRuleTypeIpv4).
			HasModeEnum(sdk.NetworkRuleModeIngress).
			HasCommentEmpty().
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasValueList(values),
		resourceshowoutputassert.NetworkRuleShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasCommentEmpty().
			HasEntriesInValueList(2).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
		resourceshowoutputassert.NetworkRuleDescOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasCommentEmpty().
			HasValueList(values).
			HasOwner(snowflakeroles.Accountadmin.Name()),
	}

	unsetAssertions := []assert.TestCheckFuncProvider{
		resourceassert.NetworkRuleResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTypeEnum(sdk.NetworkRuleTypeIpv4).
			HasModeEnum(sdk.NetworkRuleModeIngress).
			HasCommentEmpty().
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasValueList([]string{}),
		resourceshowoutputassert.NetworkRuleShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasCommentEmpty().
			HasEntriesInValueList(0).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
		resourceshowoutputassert.NetworkRuleDescOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasCommentEmpty().
			HasValueList([]string{}).
			HasOwner(snowflakeroles.Accountadmin.Name()),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.NetworkRuleResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTypeEnum(sdk.NetworkRuleTypeIpv4).
			HasModeEnum(sdk.NetworkRuleModeIngress).
			HasCommentString(changedComment).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasValueList(values),
		resourceshowoutputassert.NetworkRuleShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasComment(changedComment).
			HasEntriesInValueList(2).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
		resourceshowoutputassert.NetworkRuleDescOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasComment(changedComment).
			HasValueList(values).
			HasOwner(snowflakeroles.Accountadmin.Name()),
	}

	alteredValuesAssertions := []assert.TestCheckFuncProvider{
		resourceassert.NetworkRuleResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTypeEnum(sdk.NetworkRuleTypeIpv4).
			HasModeEnum(sdk.NetworkRuleModeIngress).
			HasCommentString(changedComment).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasValueList(changedValues),
		resourceshowoutputassert.NetworkRuleShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasComment(changedComment).
			HasEntriesInValueList(3).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
		resourceshowoutputassert.NetworkRuleDescOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeIpv4).
			HasMode(sdk.NetworkRuleModeIngress).
			HasComment(changedComment).
			HasValueList(changedValues).
			HasOwner(snowflakeroles.Accountadmin.Name()),
	}

	forceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.NetworkRuleResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTypeEnum(sdk.NetworkRuleTypeHostPort).
			HasModeEnum(sdk.NetworkRuleModeEgress).
			HasCommentString(changedComment).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasValueList(hostPortValues),
		resourceshowoutputassert.NetworkRuleShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeHostPort).
			HasMode(sdk.NetworkRuleModeEgress).
			HasComment(changedComment).
			HasEntriesInValueList(len(hostPortValues)).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"),
		resourceshowoutputassert.NetworkRuleDescOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.NetworkRuleTypeHostPort).
			HasMode(sdk.NetworkRuleModeEgress).
			HasComment(changedComment).
			HasValueList(hostPortValues).
			HasOwner(snowflakeroles.Accountadmin.Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkRule),
		Steps: []resource.TestStep{
			// Create
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:            accconfig.FromModels(t, modelBasic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Complete
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Altered values
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWithAlteredValues),
				Check:  assertThat(t, alteredValuesAssertions...),
			},
			// Changed force new values
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelWithChangedForceNewValues),
				Check:  assertThat(t, forceNewAssertions...),
			},
			// External Changes
			{
				PreConfig: func() {
					testClient().NetworkRule.Alter(t, sdk.NewAlterNetworkRuleRequest(id).WithSet(
						*sdk.NewNetworkRuleSetRequest().WithValueList([]sdk.NetworkRuleValue{
							{Value: "example.com:8080"},
						}).WithComment("external comment"),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWithChangedForceNewValues),
				Check:  assertThat(t, forceNewAssertions...),
			},
			// Bring back the original values
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelAfterUnset),
				Check:  assertThat(t, unsetAssertions...),
			},
		},
	})
}

func TestAcc_NetworkRule_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	values := []string{"192.168.0.100/24", "29.254.123.20"}
	modelComplete := model.NetworkRuleFromId(id, sdk.NetworkRuleModeIngress, sdk.NetworkRuleTypeIpv4, values).WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkRule),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.NetworkRuleResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasTypeEnum(sdk.NetworkRuleTypeIpv4).
						HasModeEnum(sdk.NetworkRuleModeIngress).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasValueList(values),
					resourceshowoutputassert.NetworkRuleShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.NetworkRuleTypeIpv4).
						HasMode(sdk.NetworkRuleModeIngress).
						HasComment(comment).
						HasEntriesInValueList(2).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
					resourceshowoutputassert.NetworkRuleDescOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.NetworkRuleTypeIpv4).
						HasMode(sdk.NetworkRuleModeIngress).
						HasComment(comment).
						HasValueList(values).
						HasOwner(snowflakeroles.Accountadmin.Name()),
				),
			},
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_NetworkRule_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelInvalidType := model.NetworkRule("test", id.DatabaseName(), id.SchemaName(), id.Name(), string(sdk.NetworkRuleModeIngress), "invalid", []string{})
	modelInvalidMode := model.NetworkRule("test", id.DatabaseName(), id.SchemaName(), id.Name(), "invalid", string(sdk.NetworkRuleTypeIpv4), []string{})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkRule),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid network rule type: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidMode),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid network rule mode: INVALID`),
			},
		},
	})
}

func TestAcc_NetworkRule_migrateFromVersion_2_13_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	modelBasic := model.NetworkRuleFromId(id, sdk.NetworkRuleModeIngress, sdk.NetworkRuleTypeIpv4, []string{})
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.NetworkRuleResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            accconfig.FromModels(t, providerModel, modelBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "id", helpers.EncodeSnowflakeID(id)),
					resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
		},
	})
}
