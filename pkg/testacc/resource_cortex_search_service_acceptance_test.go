//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CortexSearchService_basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	newWarehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(id.Name()),
			"on":         config.StringVariable("SOME_TEXT"),
			"database":   config.StringVariable(TestDatabaseName),
			"schema":     config.StringVariable(TestSchemaName),
			"warehouse":  config.StringVariable(TestWarehouseName),
			"query":      config.StringVariable(fmt.Sprintf("select SOME_TEXT from %s", tableId.FullyQualifiedName())),
			"comment":    config.StringVariable("Terraform acceptance test"),
			"table_name": config.StringVariable(tableId.Name()),
		}
	}
	variableSet2 := m()
	variableSet2["attributes"] = config.SetVariable(config.StringVariable("SOME_OTHER_TEXT"))
	variableSet2["warehouse"] = config.StringVariable(newWarehouseId.Name())
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")
	variableSet2["query"] = config.StringVariable(fmt.Sprintf("select SOME_TEXT, SOME_OTHER_TEXT from %s", tableId.FullyQualifiedName()))
	variableSet2["embedding_model"] = config.StringVariable("snowflake-arctic-embed-m-v1.5")

	resourceName := "snowflake_cortex_search_service.css"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CortexSearchService),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "on", "SOME_TEXT"),
					resource.TestCheckNoResourceAttr(resourceName, "attributes"),
					resource.TestCheckResourceAttr(resourceName, "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr(resourceName, "target_lag", "2 minutes"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "query", fmt.Sprintf("select SOME_TEXT from %s", tableId.FullyQualifiedName())),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet2,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "on", "SOME_TEXT"),
					resource.TestCheckResourceAttr(resourceName, "attributes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0", "SOME_OTHER_TEXT"),
					resource.TestCheckResourceAttr(resourceName, "warehouse", newWarehouseId.Name()),
					resource.TestCheckResourceAttr(resourceName, "target_lag", "2 minutes"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
					resource.TestCheckResourceAttr(resourceName, "query", fmt.Sprintf("select SOME_TEXT, SOME_OTHER_TEXT from %s", tableId.FullyQualifiedName())),
					resource.TestCheckResourceAttr(resourceName, "embedding_model", "snowflake-arctic-embed-m-v1.5"),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
				),
			},
			// test import
			{
				ConfigDirectory:   config.TestStepDirectory(),
				ConfigVariables:   variableSet2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// currently not set in read because the early implementation on Snowflake side did not return these values on SHOW/DESCRIBE
				ImportStateVerifyIgnore: []string{"attributes", "on", "query", "target_lag", "warehouse"},
			},
		},
	})
}
