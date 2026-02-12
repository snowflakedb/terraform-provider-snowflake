//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_HybridTable_basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(TestDatabaseName),
			"schema":   config.StringVariable(TestSchemaName),
			"comment":  config.StringVariable(comment),
		}
	}
	variableSet2 := m()
	variableSet2["comment"] = config.StringVariable(newComment)

	resourceName := "snowflake_hybrid_table.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with basic configuration
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "column.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "id"),
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr(resourceName, "column.0.nullable", "false"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "VARCHAR(100)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.nullable", "true"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "created_at"),
					resource.TestCheckResourceAttr(resourceName, "column.2.type", "TIMESTAMP_NTZ"),
					resource.TestCheckResourceAttr(resourceName, "constraint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.name", "pk_id"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.type", "PRIMARY KEY"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.columns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.columns.0", "id"),
				),
			},
			// Update comment
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "comment", newComment),
				),
			},
			// Import
			{
				ConfigDirectory:         ConfigurationSameAsStepN(2),
				ConfigVariables:         variableSet2,
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint"},
			},
		},
	})
}

func TestAcc_HybridTable_withIndexes(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(TestDatabaseName),
			"schema":   config.StringVariable(TestSchemaName),
			"comment":  config.StringVariable(comment),
		}
	}

	resourceName := "snowflake_hybrid_table.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with indexes
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "index.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "index.0.name", "idx_name"),
					resource.TestCheckResourceAttr(resourceName, "index.0.columns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index.0.columns.0", "name"),
					resource.TestCheckResourceAttr(resourceName, "index.1.name", "idx_created"),
					resource.TestCheckResourceAttr(resourceName, "index.1.columns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index.1.columns.0", "created_at"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_withForeignKey(t *testing.T) {
	parentTableId := testClient().Ids.RandomSchemaObjectIdentifier()
	childTableId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"parent_name": config.StringVariable(parentTableId.Name()),
			"child_name":  config.StringVariable(childTableId.Name()),
			"database":    config.StringVariable(TestDatabaseName),
			"schema":      config.StringVariable(TestSchemaName),
			"comment":     config.StringVariable(comment),
		}
	}

	resourceName := "snowflake_hybrid_table.child"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: ComposeCheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", childTableId.Name()),
					resource.TestCheckResourceAttr(resourceName, "constraint.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.type", "FOREIGN KEY"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.columns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.columns.0", "parent_id"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.foreign_key.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.foreign_key.0.table_id", parentTableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.foreign_key.0.columns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.foreign_key.0.columns.0", "id"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_multipleConstraints(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(TestDatabaseName),
			"schema":   config.StringVariable(TestSchemaName),
			"comment":  config.StringVariable(comment),
		}
	}

	resourceName := "snowflake_hybrid_table.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "constraint.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.name", "pk_id"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.type", "PRIMARY KEY"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.name", "uq_email"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.type", "UNIQUE"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.columns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.columns.0", "email"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_constraintForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(TestDatabaseName),
			"schema":   config.StringVariable(TestSchemaName),
			"comment":  config.StringVariable(comment),
		}
	}

	resourceName := "snowflake_hybrid_table.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_HybridTable_constraintForceNew/2"),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAcc_HybridTable_missingPrimaryKey(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(TestDatabaseName),
			"schema":   config.StringVariable(TestSchemaName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ExpectError:     regexp.MustCompile("primary key is required"),
			},
		},
	})
}
