//go:build non_account_level_tests

package testacc

import (
	"context"
	"fmt"
	"testing"

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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDbtProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbtProjectConfig_basic(databaseName, schemaName, stageName, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "name", projectName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "from.0.stage", stageName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "default_version", "LAST"),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDbtProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbtProjectConfig_gitSource(databaseName, schemaName, stageName, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "name", projectName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.repository_url", "https://github.com/Snowflake-Labs/getting-started-with-dbt-on-snowflake.git"),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.branch", "main"),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.path", "tasty_bytes_dbt_demo"),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.stage", stageName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.stage_path", "dbt_demo"),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "default_version", "LAST"),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDbtProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbtProjectConfig_gitSourceWithTag(databaseName, schemaName, stageName, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "name", projectName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.repository_url", "https://github.com/Snowflake-Labs/getting-started-with-dbt-on-snowflake.git"),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.tag", "v0"),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.path", "tasty_bytes_dbt_demo"),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "git_source.0.stage", stageName),
					resource.TestCheckResourceAttr("snowflake_dbt_project.test", "default_version", "LAST"),
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

func testAccDbtProjectConfig_basic(databaseName, schemaName, stageName, projectName string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%s"
}

resource "snowflake_schema" "test" {
  database = snowflake_database.test.name
  name     = "%s"
}

resource "snowflake_stage" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "%s"
}

resource "snowflake_dbt_project" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "%s"

  from {
    stage = snowflake_stage.test.fully_qualified_name
  }

  default_version = "LAST"
  comment         = "Test DBT project"
}
`, databaseName, schemaName, stageName, projectName)
}

func testAccDbtProjectConfig_gitSource(databaseName, schemaName, stageName, projectName string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%s"
}

resource "snowflake_schema" "test" {
  database = snowflake_database.test.name
  name     = "%s"
}

resource "snowflake_stage" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "%s"
}

resource "snowflake_dbt_project" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "%s"

  git_source {
    repository_url = "https://github.com/Snowflake-Labs/getting-started-with-dbt-on-snowflake.git"
    branch         = "main"
    path           = "tasty_bytes_dbt_demo"
    stage          = snowflake_stage.test.fully_qualified_name
    stage_path     = "dbt_demo"
  }

  default_version = "LAST"
  comment         = "Test DBT project with Git integration"
}
`, databaseName, schemaName, stageName, projectName)
}

func testAccDbtProjectConfig_gitSourceWithTag(databaseName, schemaName, stageName, projectName string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%s"
}

resource "snowflake_schema" "test" {
  database = snowflake_database.test.name
  name     = "%s"
}

resource "snowflake_stage" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "%s"
}

resource "snowflake_dbt_project" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "%s"

  git_source {
    repository_url = "https://github.com/Snowflake-Labs/getting-started-with-dbt-on-snowflake.git"
    tag            = "v0"
    path           = "tasty_bytes_dbt_demo"
    stage          = snowflake_stage.test.fully_qualified_name
  }

  default_version = "LAST"
  comment         = "Test DBT project with Git tag"
}
`, databaseName, schemaName, stageName, projectName)
}
