//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_McpServer_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	hclSpec := model.DefaultSpecAsYamlencodeHCL()
	normalizedSpec, err := sdk.NormalizeMcpServerSpecification(testClient().McpServer.DefaultSpec())
	require.NoError(t, err)

	altHclSpec := model.AltSpecAsYamlencodeHCL()
	altRawSpec := `tools:
  - title: "SQL Execution Tool"
    name: "sql_exec_tool"
    type: "SYSTEM_EXECUTE_SQL"
    description: "Updated description for acceptance tests."
`
	altNormalizedSpec, err := sdk.NormalizeMcpServerSpecification(altRawSpec)
	require.NoError(t, err)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.McpServerWithSpecification("t", id.DatabaseName(), id.SchemaName(), id.Name(), hclSpec)
	complete := model.McpServerWithSpecification("t", id.DatabaseName(), id.SchemaName(), id.Name(), altHclSpec).
		WithComment(comment)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.McpServerResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSpecification(normalizedSpec).
			HasCommentEmpty(),
		resourceshowoutputassert.McpServerShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
		resourceshowoutputassert.McpServerDescribeOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasServerSpec(normalizedSpec),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.McpServerResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSpecification(altNormalizedSpec).
			HasComment(comment),
		resourceshowoutputassert.McpServerShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment),
		resourceshowoutputassert.McpServerDescribeOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasServerSpec(altNormalizedSpec),
	}

	basicAssertionsWithEmptyValues := []assert.TestCheckFuncProvider{
		resourceassert.McpServerResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSpecification(normalizedSpec).
			HasCommentEmpty(),
		resourceshowoutputassert.McpServerShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
		resourceshowoutputassert.McpServerDescribeOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasServerSpec(normalizedSpec),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.McpServer),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Set comment + new specification
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertionsWithEmptyValues...),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, basic),
				Destroy: true,
			},
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change all props externally (via CREATE OR REPLACE — no ALTER available)
			{
				PreConfig: func() {
					// Replace the MCP server with a different spec and comment externally.
					testClient().McpServer.CreateWithRequest(t,
						sdk.NewCreateMcpServerRequest(id, testClient().McpServer.DefaultSpec()).
							WithComment(externalComment).
							WithOrReplace(true),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
		},
	})
}

func TestAcc_McpServer_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	emptySpec := model.McpServer("t", id.DatabaseName(), id.SchemaName(), id.Name(), "")
	specWithDoubleDollar := model.McpServer("t", id.DatabaseName(), id.SchemaName(), id.Name(), "contains $$ sequence")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.McpServer),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, emptySpec),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "specification" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, specWithDoubleDollar),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`cannot contain the \$\$ sequence`),
			},
		},
	})
}
