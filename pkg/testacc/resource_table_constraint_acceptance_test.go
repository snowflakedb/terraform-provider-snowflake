//go:build !account_level_tests

package testacc

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TableConstraint_fk(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintFKConfig(id),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "type", "FOREIGN KEY"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "enforced", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "deferrable", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "comment", "hello fk"),
				),
			},
		},
	})
}

func tableConstraintFKConfig(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table" "fk_t" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "fk_%[3]s"
	column {
		name     = "fk_col1"
		type     = "text"
		nullable = false
	  }
}

resource "snowflake_table_constraint" "fk" {
	name="%[3]s"
	type= "FOREIGN KEY"
	table_id = snowflake_table.t.fully_qualified_name
	columns = ["col1"]
	foreign_key_properties {
	  references {
		table_id = snowflake_table.fk_t.fully_qualified_name
		columns = ["fk_col1"]
	  }
	}
	enforced = false
	deferrable = false
	initially = "IMMEDIATE"
	comment = "hello fk"
}

`, id.DatabaseName(), id.SchemaName(), id.Name())
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2674
// It is connected with https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2629.
// Provider defaults will be reworked during resources redesign.
func TestAcc_TableConstraint_pk(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	constraintName := fmt.Sprintf("pk_%s", tableId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintPKConfig(tableId, constraintName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.pk", "type", "PRIMARY KEY"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.pk", "comment", "hello pk"),
					checkPrimaryKeyExists(tableId, constraintName),
				),
			},
		},
	})
}

func tableConstraintPKConfig(tableId sdk.SchemaObjectIdentifier, constraintName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "pk" {
	name = "%[4]s"
	type = "PRIMARY KEY"
	table_id = snowflake_table.t.fully_qualified_name
	columns = ["col1"]
	enable = false
	deferrable = false
	comment = "hello pk"
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), constraintName)
}

type PrimaryKeys struct {
	ConstraintName string `db:"constraint_name"`
}

func checkPrimaryKeyExists(tableId sdk.SchemaObjectIdentifier, constraintName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		var keys []PrimaryKeys
		err := client.QueryForTests(ctx, &keys, fmt.Sprintf("show primary keys in %s", tableId.FullyQualifiedName()))
		if err != nil {
			return err
		}

		var found bool
		for _, pk := range keys {
			if pk.ConstraintName == strings.ToUpper(constraintName) {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("unable to find primary key %s on table %s, found: %v", constraintName, tableId.FullyQualifiedName(), keys)
		}

		return nil
	}
}

func TestAcc_TableConstraint_unique(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintUniqueConfig(id),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "enforced", "true"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "deferrable", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "comment", "hello unique"),
				),
			},
		},
	})
}

func tableConstraintUniqueConfig(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "unique" {
	name="%[3]s"
	type= "UNIQUE"
	table_id = snowflake_table.t.fully_qualified_name
	columns = ["col1"]
	enforced = true
	deferrable = false
	initially = "IMMEDIATE"
	comment = "hello unique"
}

`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2535
func TestAcc_Table_issue2535_newConstraint(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.86.0"),
				Config:            tableConstraintUniqueConfigUsingTableId(id, "|"),
				ExpectError:       regexp.MustCompile(`.*table id is incorrect.*`),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.89.0"),
				Config:            tableConstraintUniqueConfigUsingTableId(id, "|"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConstraintUniqueConfigUsingTableId(id, "|"),
				ExpectError:              regexp.MustCompile(`.*Expected SchemaObjectIdentifier identifier type, but got:.*`),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConstraintUniqueConfigUsingTableId(id, "."),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
				),
			},
		},
	})
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2535
func TestAcc_Table_issue2535_existingTable(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// reference done by table.id in 0.85.0
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.85.0"),
				Config:            tableConstraintUniqueConfigUsingTableId(id, "|"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
				),
			},
			// switched to qualified_name in 0.86.0
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.86.0"),
				Config:            tableConstraintUniqueConfigUsingQualifiedName(id),
				ExpectError:       regexp.MustCompile(`.*table id is incorrect.*`),
			},
			// fixed in the current version
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConstraintUniqueConfigUsingFullyQualifiedName(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
				),
			},
		},
	})
}

func tableConstraintUniqueConfigUsingTableId(tableId sdk.SchemaObjectIdentifier, idSeparator string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "unique" {
	name     = "%[3]s"
	type     = "UNIQUE"
	table_id = "${snowflake_table.t.database}%[4]s${snowflake_table.t.schema}%[4]s${snowflake_table.t.name}"
	columns  = ["col1"]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), idSeparator)
}

func tableConstraintUniqueConfigUsingQualifiedName(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "unique" {
	name     = "%[3]s"
	type     = "UNIQUE"
	table_id = snowflake_table.t.qualified_name
	columns  = ["col1"]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConstraintUniqueConfigUsingFullyQualifiedName(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "unique" {
	name     = "%[3]s"
	type     = "UNIQUE"
	table_id = snowflake_table.t.fully_qualified_name
	columns  = ["col1"]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

// TODO(issue-2683): Uncomment once the Update operation is ready
// func TestAcc_TableConstraint_Rename(t *testing.T) {
//	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
//	newName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
//		PreCheck:                 func() { TestAccPreCheck(t) },
//		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
//			tfversion.RequireAbove(tfversion.Version1_5_0),
//		},
//		CheckDestroy: nil,
//		Steps: []resource.TestStep{
//			{
//				Config: tableConstraintUniqueConfigUsingFullyQualifiedName(name, TestDatabaseName, TestSchemaName),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "name", name),
//				),
//			},
//			{
//				Config: tableConstraintUniqueConfigUsingFullyQualifiedName(newName, TestDatabaseName, TestSchemaName),
//				ConfigPlanChecks: resource.ConfigPlanChecks{
//					PreApply: []plancheck.PlanCheck{
//						plancheck.ExpectResourceAction("snowflake_table_constraint.unique", plancheck.ResourceActionUpdate),
//					},
//				},
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "name", newName),
//				),
//			},
//		},
//	})
//}

func TestAcc_TableConstraint_ProperlyHandles_EmptyForeignKeyProperties(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      tableConstraintEmptyForeignKeyProperties(id),
				ExpectError: regexp.MustCompile(`At least 1 "references" blocks are required.`),
			},
			{
				Config: tableConstraintForeignKeyProperties(id),
			},
		},
	})
}

func tableConstraintEmptyForeignKeyProperties(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
	resource "snowflake_table" "t" {
		database = "%[1]s"
		schema   = "%[2]s"
		name     = "%[3]s"

		column {
			name = "col1"
			type = "NUMBER(38,0)"
		}
	}

	resource "snowflake_table_constraint" "unique" {
		name = "%[3]s"
		type = "FOREIGN KEY"
		table_id = snowflake_table.t.fully_qualified_name
		columns = ["col1"]
		foreign_key_properties {
		}
	}
	`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConstraintForeignKeyProperties(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
	resource "snowflake_table" "t" {
		database = "%[1]s"
		schema   = "%[2]s"
		name     = "%[3]s"

		column {
			name = "col1"
			type = "NUMBER(38,0)"
		}
	}

	resource "snowflake_table_constraint" "unique" {
		name = "%[3]s"
		type = "FOREIGN KEY"
		table_id = snowflake_table.t.fully_qualified_name
		columns = ["col1"]
		foreign_key_properties {
			references {
				columns = ["col1"]
				table_id = snowflake_table.t.fully_qualified_name
			}
		}
	}
	`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}
