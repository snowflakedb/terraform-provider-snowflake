//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_IntelligenceCortexAgentAttachment_BasicUseCase(t *testing.T) {
	intelligenceId, intelligenceCleanup := testClient().SnowflakeIntelligence.Create(t)
	t.Cleanup(intelligenceCleanup)

	agentId := testClient().Ids.RandomSchemaObjectIdentifier()
	agentCleanup := testClient().CortexAgent.CreateWithId(t, agentId)
	t.Cleanup(agentCleanup)

	newAgentId := testClient().Ids.RandomSchemaObjectIdentifier()
	newAgentCleanup := testClient().CortexAgent.CreateWithId(t, newAgentId)
	t.Cleanup(newAgentCleanup)

	intelligenceName := intelligenceId.FullyQualifiedName()
	agentName := agentId.FullyQualifiedName()
	newAgentName := newAgentId.FullyQualifiedName()

	basicModel := model.IntelligenceCortexAgentAttachment("t", agentName, intelligenceName)
	newAgentModel := model.IntelligenceCortexAgentAttachment("t", newAgentName, intelligenceName)
	ref := basicModel.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IntelligenceCortexAgentAttachmentResource(t, ref).
			HasSnowflakeIntelligenceName(intelligenceName).
			HasCortexAgentName(agentName),
	}
	newAgentAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IntelligenceCortexAgentAttachmentResource(t, ref).
			HasSnowflakeIntelligenceName(intelligenceName).
			HasCortexAgentName(newAgentName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckIntelligenceCortexAgentAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// Create
			{
				Config: config.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, basicAssertions...),
			},
			// Import
			{
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change the attached agent. Both identifiers are ForceNew, so Terraform replaces the attachment.
			{
				Config: config.FromModels(t, newAgentModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, newAgentAssertions...),
			},
			// Drop the referenced agent externally and remove the attachment from config.
			{
				PreConfig: func() {
					testClient().SnowflakeIntelligence.DropAgent(t, intelligenceId, newAgentId)
					newAgentCleanup()
				},
				Config: " ",
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Recreate
			{
				Config: config.FromModels(t, basicModel),
				Check:  assertThat(t, basicAssertions...),
			},
			// Detach the agent externally. The attachment no longer exists, so Terraform creates it again.
			{
				PreConfig: func() {
					testClient().SnowflakeIntelligence.DropAgent(t, intelligenceId, agentId)
				},
				Config: config.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, basicAssertions...),
			},
		},
	})
}
