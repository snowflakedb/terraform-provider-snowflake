//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	tfjson "github.com/hashicorp/terraform-json"
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

// TestAcc_HybridTable_allDataTypes tests table with comprehensive data types
// Covers scenario HT-A-TF-002
func TestAcc_HybridTable_allDataTypes(t *testing.T) {
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
			// Create with various data types
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "column.#", "15"),
					// Numeric types
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "id"),
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "col_int"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "col_float"),
					resource.TestCheckResourceAttr(resourceName, "column.2.type", "FLOAT"),
					resource.TestCheckResourceAttr(resourceName, "column.3.name", "col_decimal"),
					resource.TestCheckResourceAttr(resourceName, "column.3.type", "NUMBER(10,2)"),
					// String types
					resource.TestCheckResourceAttr(resourceName, "column.4.name", "col_varchar"),
					resource.TestCheckResourceAttr(resourceName, "column.4.type", "VARCHAR(100)"),
					resource.TestCheckResourceAttr(resourceName, "column.5.name", "col_char"),
					resource.TestCheckResourceAttr(resourceName, "column.5.type", "VARCHAR(10)"),
					resource.TestCheckResourceAttr(resourceName, "column.6.name", "col_text"),
					resource.TestCheckResourceAttr(resourceName, "column.6.type", "VARCHAR(134217728)"),
					// Date/Time types
					resource.TestCheckResourceAttr(resourceName, "column.7.name", "col_date"),
					resource.TestCheckResourceAttr(resourceName, "column.7.type", "DATE"),
					resource.TestCheckResourceAttr(resourceName, "column.8.name", "col_time"),
					resource.TestCheckResourceAttr(resourceName, "column.8.type", "TIME(9)"),
					resource.TestCheckResourceAttr(resourceName, "column.9.name", "col_timestamp_ntz"),
					resource.TestCheckResourceAttr(resourceName, "column.9.type", "TIMESTAMP_NTZ"),
					resource.TestCheckResourceAttr(resourceName, "column.10.name", "col_timestamp_ltz"),
					resource.TestCheckResourceAttr(resourceName, "column.10.type", "TIMESTAMP_LTZ"),
					// Semi-structured types
					resource.TestCheckResourceAttr(resourceName, "column.11.name", "col_variant"),
					resource.TestCheckResourceAttr(resourceName, "column.11.type", "VARIANT"),
					resource.TestCheckResourceAttr(resourceName, "column.12.name", "col_object"),
					resource.TestCheckResourceAttr(resourceName, "column.12.type", "OBJECT"),
					resource.TestCheckResourceAttr(resourceName, "column.13.name", "col_array"),
					resource.TestCheckResourceAttr(resourceName, "column.13.type", "ARRAY"),
					// Boolean type
					resource.TestCheckResourceAttr(resourceName, "column.14.name", "col_boolean"),
					resource.TestCheckResourceAttr(resourceName, "column.14.type", "BOOLEAN"),
				),
			},
			// Import
			{
				ConfigDirectory:         config.StaticDirectory("testdata/TestAcc_HybridTable_allDataTypes/1"),
				ConfigVariables:         m(),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint"},
			},
		},
	})
}

// TestAcc_HybridTable_columnLifecycle tests that column changes force table recreation
// Covers scenarios HT-A-TF-003 and HT-A-TF-004
func TestAcc_HybridTable_columnLifecycle(t *testing.T) {
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
			// Step 1: Create with 3 columns (id, name, created_at)
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "column.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "id"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "created_at"),
				),
			},
			// Step 2: Add column (email) - forces recreation
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "column.3.name", "email"),
					resource.TestCheckResourceAttr(resourceName, "column.3.type", "VARCHAR(255)"),
				),
			},
			// Step 3: Drop column (created_at) - forces recreation
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "id"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "email"),
				),
			},
			// Step 4: Update column comment - forces recreation
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "column.1.comment", "User full name"),
				),
			},
			// Step 5: Import
			{
				ConfigDirectory:         config.StaticDirectory("testdata/TestAcc_HybridTable_columnLifecycle/4"),
				ConfigVariables:         m(),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint"},
			},
		},
	})
}

// TestAcc_HybridTable_nameChange tests table rename behavior
// Covers scenario HT-A-TF-007 (partial)
func TestAcc_HybridTable_nameChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(TestDatabaseName),
			"schema":   config.StringVariable(TestSchemaName),
			"comment":  config.StringVariable(comment),
		}
	}

	variableSet2 := m()
	variableSet2["name"] = config.StringVariable(newId.Name())

	resourceName := "snowflake_hybrid_table.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Step 1: Create with original name
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
			// Step 2: Change name - forces recreation (rename not implemented in update function)
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_HybridTable_nameChange/1"),
				ConfigVariables: variableSet2,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", newId.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", newId.FullyQualifiedName()),
				),
			},
		},
	})
}

// TestAcc_HybridTable_databaseSchemaChange tests database and schema changes force recreation
// Covers scenario HT-A-TF-007 (partial)
func TestAcc_HybridTable_databaseSchemaChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	secondSchemaId := testClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":          config.StringVariable(id.Name()),
			"database":      config.StringVariable(TestDatabaseName),
			"schema":        config.StringVariable(TestSchemaName),
			"second_schema": config.StringVariable(secondSchemaId.Name()),
			"comment":       config.StringVariable(comment),
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
			// Step 1: Create in first schema
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
				),
			},
			// Step 2: Change schema - should force recreation
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", secondSchemaId.Name()),
				),
			},
		},
	})
}

// TestAcc_HybridTable_compositeConstraints tests multi-column primary key and unique constraints
func TestAcc_HybridTable_compositeConstraints(t *testing.T) {
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
			// Create with composite primary key and unique constraint
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "constraint.#", "2"),
					// Composite primary key
					resource.TestCheckResourceAttr(resourceName, "constraint.0.name", "pk_composite"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.type", "PRIMARY KEY"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.columns.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.columns.0", "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.columns.1", "user_id"),
					// Composite unique constraint
					resource.TestCheckResourceAttr(resourceName, "constraint.1.name", "uq_email_tenant"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.type", "UNIQUE"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.columns.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.columns.0", "email"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.columns.1", "tenant_id"),
				),
			},
			// Import
			{
				ConfigDirectory:         config.StaticDirectory("testdata/TestAcc_HybridTable_compositeConstraints/1"),
				ConfigVariables:         m(),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint"},
			},
		},
	})
}

// TestAcc_HybridTable_columnComments tests that column comment changes force recreation
func TestAcc_HybridTable_columnComments(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	columnComment1 := random.Comment()
	columnComment2 := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":            config.StringVariable(id.Name()),
			"database":        config.StringVariable(TestDatabaseName),
			"schema":          config.StringVariable(TestSchemaName),
			"comment":         config.StringVariable(comment),
			"column_comment":  config.StringVariable(columnComment1),
		}
	}

	variableSet2 := m()
	variableSet2["column_comment"] = config.StringVariable(columnComment2)

	variableSet3 := m()
	variableSet3["column_comment"] = config.StringVariable("")

	resourceName := "snowflake_hybrid_table.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Step 1: Create with column comment
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.1.comment", columnComment1),
				),
			},
			// Step 2: Update column comment - forces recreation
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_HybridTable_columnComments/1"),
				ConfigVariables: variableSet2,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.comment", columnComment2),
				),
			},
			// Step 3: Remove column comment - forces recreation
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_HybridTable_columnComments/1"),
				ConfigVariables: variableSet3,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.comment", ""),
				),
			},
		},
	})
}

// TestAcc_HybridTable_identifierCaseSensitivity tests case handling in table and column names
// Covers scenario HT-A-STATE-001
func TestAcc_HybridTable_identifierCaseSensitivity(t *testing.T) {
	// Use mixed case names
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
			// Create with mixed-case identifiers
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "userId"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "userName"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "createdAt"),
				),
			},
			// Import with case variations
			{
				ConfigDirectory:         config.StaticDirectory("testdata/TestAcc_HybridTable_identifierCaseSensitivity/1"),
				ConfigVariables:         m(),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint"},
			},
		},
	})
}

// TestAcc_HybridTable_specialCharacterIdentifiers tests special characters in identifiers
// Covers scenario HT-A-STATE-002
func TestAcc_HybridTable_specialCharacterIdentifiers(t *testing.T) {
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
			// Create with special characters in column names
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "user-id"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "user name"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "created@time"),
				),
			},
			// Import
			{
				ConfigDirectory:         config.StaticDirectory("testdata/TestAcc_HybridTable_specialCharacterIdentifiers/1"),
				ConfigVariables:         m(),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint"},
			},
		},
	})
}

// TestAcc_HybridTable_driftDetection tests detection of out-of-band changes
// Covers scenario HT-A-TF-010
func TestAcc_HybridTable_driftDetection(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	manualComment := random.Comment()

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
			// Step 1: Create table
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			// Step 2: Detect drift after manual change
			{
				PreConfig: func() {
					testClient().HybridTable.UpdateComment(t, id, manualComment)
				},
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_HybridTable_driftDetection/1"),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift(resourceName, "comment", sdk.String(comment), sdk.String(manualComment)),
						planchecks.ExpectChange(resourceName, "comment", tfjson.ActionUpdate, sdk.String(manualComment), sdk.String(comment)),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

// TestAcc_HybridTable_indexLifecycle tests adding and removing indexes
// Covers scenario HT-A-TF-006
func TestAcc_HybridTable_indexLifecycle(t *testing.T) {
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
			// Step 1: Create without indexes
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "index.#", "0"),
				),
			},
			// Step 2: Add first index
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index.0.name", "idx_name"),
					resource.TestCheckResourceAttr(resourceName, "index.0.columns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index.0.columns.0", "name"),
				),
			},
			// Step 3: Add second index
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "index.0.name", "idx_name"),
					resource.TestCheckResourceAttr(resourceName, "index.1.name", "idx_email"),
					resource.TestCheckResourceAttr(resourceName, "index.1.columns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index.1.columns.0", "email"),
				),
			},
			// Step 4: Remove first index
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index.0.name", "idx_email"),
				),
			},
			// Step 5: Remove all indexes
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index.#", "0"),
				),
			},
		},
	})
}

// TestAcc_HybridTable_emptyCommentHandling tests comment empty string vs null handling
func TestAcc_HybridTable_emptyCommentHandling(t *testing.T) {
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

	variableSet2 := m()
	variableSet2["comment"] = config.StringVariable("")

	resourceName := "snowflake_hybrid_table.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Step 1: Create with comment
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			// Step 2: Set comment to empty string
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_HybridTable_emptyCommentHandling/1"),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
				),
			},
			// Step 3: Verify empty comment persists
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_HybridTable_emptyCommentHandling/1"),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
				),
			},
		},
	})
}

// TestAcc_HybridTable_importWithConstraints tests import state handling of constraints
func TestAcc_HybridTable_importWithConstraints(t *testing.T) {
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
			// Create with multiple constraints
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "constraint.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "constraint.0.type", "PRIMARY KEY"),
					resource.TestCheckResourceAttr(resourceName, "constraint.1.type", "UNIQUE"),
					resource.TestCheckResourceAttr(resourceName, "constraint.2.type", "UNIQUE"),
				),
			},
			// Import - must ignore constraint due to ForceNew
			{
				ConfigDirectory:         config.StaticDirectory("testdata/TestAcc_HybridTable_importWithConstraints/1"),
				ConfigVariables:         m(),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
		},
	})
}

// TestAcc_HybridTable_columnOrderPreservation tests that column order matches config order
// Covers scenario HT-A-STATE-004
func TestAcc_HybridTable_columnOrderPreservation(t *testing.T) {
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
			// Create with specific column order
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "column.#", "5"),
					// Verify column order matches config
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "id"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "created_at"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "column.3.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.4.name", "email"),
				),
			},
			// Import and verify order preserved
			{
				ConfigDirectory:         config.StaticDirectory("testdata/TestAcc_HybridTable_columnOrderPreservation/1"),
				ConfigVariables:         m(),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint"},
				Check: resource.ComposeTestCheckFunc(
					// Verify order maintained after import
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "id"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "created_at"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "column.3.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.4.name", "email"),
				),
			},
		},
	})
}

// TestAcc_HybridTable_multipleIndexOrder tests that index order is consistent
// Covers scenario HT-A-STATE-005
func TestAcc_HybridTable_multipleIndexOrder(t *testing.T) {
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
			// Create with 3 indexes in specific order
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "index.#", "3"),
					// Verify index order matches config
					resource.TestCheckResourceAttr(resourceName, "index.0.name", "idx_name"),
					resource.TestCheckResourceAttr(resourceName, "index.0.columns.0", "name"),
					resource.TestCheckResourceAttr(resourceName, "index.1.name", "idx_email"),
					resource.TestCheckResourceAttr(resourceName, "index.1.columns.0", "email"),
					resource.TestCheckResourceAttr(resourceName, "index.2.name", "idx_created"),
					resource.TestCheckResourceAttr(resourceName, "index.2.columns.0", "created_at"),
				),
			},
			// Import - indexes not read on import (limitation)
			{
				ConfigDirectory:         config.StaticDirectory("testdata/TestAcc_HybridTable_multipleIndexOrder/1"),
				ConfigVariables:         m(),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"constraint", "index"},
			},
		},
	})
}
