//go:build non_account_level_tests

package testacc

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DbtProject_basic(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	databaseName := random.AlphaN(12)
	schemaName := random.AlphaN(12)
	stageName := random.AlphaN(12)
	projectName := random.AlphaN(12)

	configModels := []config.ConfigModel{
		model.Database("test_db", databaseName),
		model.Schema("test_schema", databaseName, schemaName).
			WithDependsOn("snowflake_database.test_db"),
		model.Stage("test_stage", databaseName, schemaName, stageName).
			WithDependsOn("snowflake_schema.test_schema"),
		model.BasicDbtProjectModel("test", databaseName, schemaName, projectName).
			WithFromStage("${snowflake_stage.test_stage.fully_qualified_name}").
			WithDependsOn("snowflake_stage.test_stage"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDbtProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModels...),
				Check: resource.ComposeTestCheckFunc(
					resourceassert.DbtProjectResource(t, "snowflake_dbt_project.test").
						HasNameString(projectName).
						HasDatabaseString(databaseName).
						HasSchemaString(schemaName).
						HasFromStage("${snowflake_stage.test_stage.fully_qualified_name}").
						HasDefaultVersionString("LAST").
						HasCommentString("Test DBT project"),
				),
			},
		},
	})
}

func TestAcc_DbtProject_gitSource(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	databaseName := random.AlphaN(12)
	schemaName := random.AlphaN(12)
	stageName := random.AlphaN(12)
	projectName := random.AlphaN(12)

	configModels := []config.ConfigModel{
		model.Database("test_db", databaseName),
		model.Schema("test_schema", databaseName, schemaName).
			WithDependsOn("snowflake_database.test_db"),
		model.Stage("test_stage", databaseName, schemaName, stageName).
			WithDependsOn("snowflake_schema.test_schema"),
		model.BasicDbtProjectModel("test", databaseName, schemaName, projectName).
			WithGitSource(
				"https://github.com/Snowflake-Labs/getting-started-with-dbt-on-snowflake.git",
				"${snowflake_stage.test_stage.fully_qualified_name}",
				"main",
				"tasty_bytes_dbt_demo",
				"dbt_demo",
			).
			WithDependsOn("snowflake_stage.test_stage"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDbtProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModels...),
				Check: resource.ComposeTestCheckFunc(
					resourceassert.DbtProjectResource(t, "snowflake_dbt_project.test").
						HasNameString(projectName).
						HasDatabaseString(databaseName).
						HasSchemaString(schemaName).
						HasGitSourceWithBranch(
							"https://github.com/Snowflake-Labs/getting-started-with-dbt-on-snowflake.git",
							"main",
							"${snowflake_stage.test_stage.fully_qualified_name}",
						).
						HasGitSourcePath("tasty_bytes_dbt_demo").
						HasGitSourceStagePath("dbt_demo").
						HasDefaultVersionString("LAST").
						HasCommentString("Test DBT project"),
				),
			},
		},
	})
}

func TestAcc_DbtProject_gitSourceWithTag(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	databaseName := random.AlphaN(12)
	schemaName := random.AlphaN(12)
	stageName := random.AlphaN(12)
	projectName := random.AlphaN(12)

	configModels := []config.ConfigModel{
		model.Database("test_db", databaseName),
		model.Schema("test_schema", databaseName, schemaName).
			WithDependsOn("snowflake_database.test_db"),
		model.Stage("test_stage", databaseName, schemaName, stageName).
			WithDependsOn("snowflake_schema.test_schema"),
		model.BasicDbtProjectModel("test", databaseName, schemaName, projectName).
			WithGitSourceTag(
				"https://github.com/Snowflake-Labs/getting-started-with-dbt-on-snowflake.git",
				"${snowflake_stage.test_stage.fully_qualified_name}",
				"v0",
				"tasty_bytes_dbt_demo",
				"",
			).
			WithDependsOn("snowflake_stage.test_stage"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDbtProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModels...),
				Check: resource.ComposeTestCheckFunc(
					resourceassert.DbtProjectResource(t, "snowflake_dbt_project.test").
						HasNameString(projectName).
						HasDatabaseString(databaseName).
						HasSchemaString(schemaName).
						HasGitSourceWithTag(
							"https://github.com/Snowflake-Labs/getting-started-with-dbt-on-snowflake.git",
							"v0",
							"${snowflake_stage.test_stage.fully_qualified_name}",
						).
						HasGitSourcePath("tasty_bytes_dbt_demo").
						HasDefaultVersionString("LAST").
						HasCommentString("Test DBT project"),
				),
			},
		},
	})
}

func testAccCheckDbtProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_dbt_project" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(rs.Primary.ID)
		existingDbtProject, err := client.DbtProjects.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("dbt project %v still exists", existingDbtProject.Name)
		}
	}
	return nil
}
