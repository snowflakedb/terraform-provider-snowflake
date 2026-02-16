//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CatalogIntegration_ObjectStore_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()

	catalogIntegrationModelBasic := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "ICEBERG", false)

	catalogIntegrationModelWithComment := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "ICEBERG", false).
		WithComment(comment)

	catalogIntegrationModelEnabled := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "ICEBERG", true).
		WithComment(newComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegration),
		Steps: []resource.TestStep{
			// CREATE WITHOUT OPTIONAL ATTRIBUTES
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(catalogIntegrationModelBasic.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(catalogIntegrationModelBasic.ResourceReference(), "catalog_source", "OBJECT_STORE"),
					resource.TestCheckResourceAttr(catalogIntegrationModelBasic.ResourceReference(), "table_format", "ICEBERG"),
					resource.TestCheckResourceAttr(catalogIntegrationModelBasic.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(catalogIntegrationModelBasic.ResourceReference(), "comment", ""),
				),
			},
			// IMPORT
			{
				ResourceName:      catalogIntegrationModelBasic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// UPDATE - ADD COMMENT
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(catalogIntegrationModelWithComment.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationModelWithComment.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(catalogIntegrationModelWithComment.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(catalogIntegrationModelWithComment.ResourceReference(), "comment", comment),
				),
			},
			// UPDATE - ENABLE AND CHANGE COMMENT
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(catalogIntegrationModelEnabled.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationModelEnabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationModelEnabled.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(catalogIntegrationModelEnabled.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(catalogIntegrationModelEnabled.ResourceReference(), "comment", newComment),
				),
			},
			// DESTROY
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(catalogIntegrationModelEnabled.ResourceReference(), plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, catalogIntegrationModelEnabled),
				Destroy: true,
			},
		},
	})
}

func TestAcc_CatalogIntegration_Glue_Complete(t *testing.T) {
	// Note: This test requires AWS environment variables to be set
	t.Skip("Skipping GLUE catalog integration test - requires AWS Glue setup")

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	glueRoleArn := "arn:aws:iam::123456789012:role/SnowflakeGlueRole"
	glueCatalogId := "123456789012"
	glueRegion := "us-west-2"

	catalogIntegrationModelGlue := model.CatalogIntegration("w", id.Name(), "GLUE", "ICEBERG", true).
		WithGlueAwsRoleArn(glueRoleArn).
		WithGlueCatalogId(glueCatalogId).
		WithGlueRegion(glueRegion).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegration),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, catalogIntegrationModelGlue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationModelGlue.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(catalogIntegrationModelGlue.ResourceReference(), "catalog_source", "GLUE"),
					resource.TestCheckResourceAttr(catalogIntegrationModelGlue.ResourceReference(), "table_format", "ICEBERG"),
					resource.TestCheckResourceAttr(catalogIntegrationModelGlue.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(catalogIntegrationModelGlue.ResourceReference(), "glue_aws_role_arn", glueRoleArn),
					resource.TestCheckResourceAttr(catalogIntegrationModelGlue.ResourceReference(), "glue_catalog_id", glueCatalogId),
					resource.TestCheckResourceAttr(catalogIntegrationModelGlue.ResourceReference(), "glue_region", glueRegion),
					resource.TestCheckResourceAttr(catalogIntegrationModelGlue.ResourceReference(), "comment", comment),
				),
			},
			// IMPORT
			{
				ResourceName:      catalogIntegrationModelGlue.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAcc_CatalogIntegration_ForceNew verifies that changing table_format (a ForceNew field)
// causes Terraform to destroy and recreate the resource.
func TestAcc_CatalogIntegration_ForceNew(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogIntegrationIceberg := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "ICEBERG", true)
	catalogIntegrationDelta := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "DELTA", true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegration),
		Steps: []resource.TestStep{
			// CREATE with ICEBERG
			{
				Config: config.FromModels(t, catalogIntegrationIceberg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationIceberg.ResourceReference(), "table_format", "ICEBERG"),
				),
			},
			// CHANGE table_format to DELTA - should force recreate
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(catalogIntegrationDelta.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationDelta),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationDelta.ResourceReference(), "table_format", "DELTA"),
				),
			},
		},
	})
}

func TestAcc_CatalogIntegration_Unset(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	catalogIntegrationWithComment := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "ICEBERG", true).
		WithComment(comment)

	catalogIntegrationWithoutComment := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "ICEBERG", true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegration),
		Steps: []resource.TestStep{
			// CREATE WITH COMMENT
			{
				Config: config.FromModels(t, catalogIntegrationWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationWithComment.ResourceReference(), "comment", comment),
				),
			},
			// UNSET COMMENT
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(catalogIntegrationWithoutComment.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationWithoutComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationWithoutComment.ResourceReference(), "comment", ""),
				),
			},
		},
	})
}

func TestAcc_CatalogIntegration_Invalid(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      catalogIntegrationConfigInvalidSource(id.Name()),
				ExpectError: regexp.MustCompile(`expected catalog_source to be one of`),
			},
			{
				Config:      catalogIntegrationConfigInvalidTableFormat(id.Name()),
				ExpectError: regexp.MustCompile(`expected table_format to be one of`),
			},
		},
	})
}

func catalogIntegrationConfigInvalidSource(name string) string {
	return fmt.Sprintf(`
resource "snowflake_catalog_integration" "test" {
  name           = "%s"
  catalog_source = "INVALID_SOURCE"
  table_format   = "ICEBERG"
  enabled        = true
}
`, name)
}

func catalogIntegrationConfigInvalidTableFormat(name string) string {
	return fmt.Sprintf(`
resource "snowflake_catalog_integration" "test" {
  name           = "%s"
  catalog_source = "OBJECT_STORE"
  table_format   = "INVALID_FORMAT"
  enabled        = true
}
`, name)
}
