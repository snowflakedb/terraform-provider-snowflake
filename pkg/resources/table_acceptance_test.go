//go:build !account_level_tests

package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TableWithSeparateDataRetentionObjectParameterWithoutLifecycle(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfig(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableAndDataRetentionParameterConfigWithoutLifecycle(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_object_parameter.data_retention_in_time", "value", "30"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAcc_TableWithSeparateDataRetentionObjectParameterWithLifecycle(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfig(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableAndDataRetentionParameterConfigWithLifecycle(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_object_parameter.data_retention_in_time", "value", "30"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: updatedTableAndDataRetentionParameterConfigWithLifecycle(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Table with a separate data retention parameter"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_object_parameter.data_retention_in_time", "value", "30"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func TestAcc_Table(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	table1Id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	table2Id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	table3Id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfig(table1Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.schema_evolution_record", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfig2(table1Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfig3(table2Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.0", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
				),
			},
			{
				Config: tableConfig4(table2Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.1", "\"col2\""),
				),
			},
			{
				Config: tableConfig5(table2Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.0", "\"col2\""),
				),
			},
			{
				Config: tableConfig6(table1Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.nullable", "false"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfig7(table1Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.nullable", "false"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.keys.0", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.name", ""),
				),
			},
			{
				Config: tableConfig8(table1Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.nullable", "false"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.keys.0", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.keys.1", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.name", "new_name"),
				),
			},
			{
				Config: tableConfig9CreateTableWithColumnComment(table2Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "COL2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.comment", "some comment"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "primary_key.0"),
				),
			},
			{
				Config: tableConfig10AlterTableColumnComment(table2Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.comment", "other comment"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "COL2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.comment", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "primary_key.0"),
				),
			},
			{
				Config: tableConfig11AlterTableAddColumnWithComment(table2Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.comment", "other comment"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "COL2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.2.name", "COL3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.2.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.2.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.2.comment", "extra"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "primary_key.0"),
				),
			},
			{
				Config: tableConfig12CreateTableWithDataRetention(table3Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "fully_qualified_name", table3Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "data_retention_time_in_days", "10"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.1.name", "column2"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "primary_key.0"),
				),
			},
			{
				Config: tableConfig13AlterTableDataRetention(table3Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "fully_qualified_name", table3Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "data_retention_time_in_days", "0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.1.name", "column2"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "primary_key.0"),
				),
			},
			{
				Config: tableConfig14AlterTableEnableChangeTracking(table3Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "fully_qualified_name", table3Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "data_retention_time_in_days", "0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "change_tracking", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.1.name", "column2"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "primary_key.0"),
				),
			},
			{
				Config: tableConfig15CreateTableWithChangeTracking(table1Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableAndDataRetentionParameterConfigWithoutLifecycle(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}

resource "snowflake_object_parameter" "data_retention_in_time" {
	key = "DATA_RETENTION_TIME_IN_DAYS"
	value = "30"
    object_type = "TABLE"
    object_identifier {
		database = "%[1]s"
		schema   = "%[2]s"
		name     = "%[3]s"
	}
	depends_on = [snowflake_table.test_table]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableAndDataRetentionParameterConfigWithLifecycle(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
	lifecycle {
		ignore_changes = [
			"data_retention_time_in_days"
		]
	}
}

resource "snowflake_object_parameter" "data_retention_in_time" {
	key = "DATA_RETENTION_TIME_IN_DAYS"
	value = "30"
    object_type = "TABLE"
    object_identifier {
		database = "%[1]s"
		schema   = "%[2]s"
		name     = "%[3]s"
	}
	depends_on = [snowflake_table.test_table]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func updatedTableAndDataRetentionParameterConfigWithLifecycle(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Table with a separate data retention parameter"
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
	lifecycle {
		ignore_changes = [
			"data_retention_time_in_days"
		]
	}
}
resource "snowflake_object_parameter" "data_retention_in_time" {
	key = "DATA_RETENTION_TIME_IN_DAYS"
	value = "30"
    object_type = "TABLE"
    object_identifier {
		database = "%[1]s"
		schema   = "%[2]s"
		name     = "%[3]s"
	}
	depends_on = [snowflake_table.test_table]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	data_retention_time_in_days = 1
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig2(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig3(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table2" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	data_retention_time_in_days = 1
	comment  = "Terraform acceptance test"
	cluster_by = ["COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig4(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table2" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	cluster_by = ["COL1","\"col2\""]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig5(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table2" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	cluster_by = ["\"col2\"","COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig6(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
		nullable = false
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig7(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
		nullable = false
	}
	primary_key {
		keys = ["column2"]
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig8(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
		nullable = false
	}
	primary_key {
		name = "new_name"
		keys = ["column2","column3"]
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig9CreateTableWithColumnComment(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table2" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name     = "COL2"
		type     = "FLOAT"
		nullable = false
		comment  = "some comment"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig10AlterTableColumnComment(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table2" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name    = "COL1"
		type    = "VARCHAR(16777216)"
		comment = "other comment"
	}
	column {
		name     = "COL2"
		type     = "FLOAT"
		nullable = false
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig11AlterTableAddColumnWithComment(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table2" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name    = "COL1"
		type    = "VARCHAR(16777216)"
		comment = "other comment"
	}
	column {
		name     = "COL2"
		type     = "FLOAT"
		nullable = false
	}
	column {
		name     = "COL3"
		type     = "NUMBER(38,0)"
		nullable = false
		comment  = "extra"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig12CreateTableWithDataRetention(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table3" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 10
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig13AlterTableDataRetention(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table3" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 0
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig14AlterTableEnableChangeTracking(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table3" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 0
	change_tracking     = true
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfig15CreateTableWithChangeTracking(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 1
	change_tracking     = true
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_TableDefaults(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithDefaults(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARCHAR(16)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.default.0.constant", "hello"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.expression"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.constant"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.default.0.expression", "CURRENT_TIMESTAMP()"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.expression"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.default.0.sequence", fmt.Sprintf(`"%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableId.Name())),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableColumnWithoutDefaults(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARCHAR(16)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.default.0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.expression"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.default.0.sequence", fmt.Sprintf(`"%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableId.Name())),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableColumnWithDefaults(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_seq" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
		default {
			constant = "hello"
		}
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
		default {
			expression = "CURRENT_TIMESTAMP()"
		}
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		default {
			sequence = snowflake_sequence.test_seq.fully_qualified_name
		}
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableColumnWithoutDefaults(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_seq" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		default {
			sequence = snowflake_sequence.test_seq.fully_qualified_name
		}
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_TableTags(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag1, tag1Cleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tag1Cleanup)

	tag2, tag2Cleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tag2Cleanup)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tagValue := random.AlphaN(4)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableWithTags(tableId, tag1.ID(), tag2.ID(), tagValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.0.name", tag1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.0.value", tagValue),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.0.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.0.schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.1.name", tag2.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.1.value", tagValue),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.1.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.1.schema", acc.TestSchemaName),
				),
			},
		},
	})
}

func tableWithTags(tableId sdk.SchemaObjectIdentifier, tagId sdk.SchemaObjectIdentifier, tag2Id sdk.SchemaObjectIdentifier, tagValue string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
	}

	tag {
		database = "%[1]s"
		schema = "%[2]s"
		name = "%[4]s"
		value = "%[6]s"
	}

	tag {
		database = "%[1]s"
		schema = "%[2]s"
		name = "%[5]s"
		value = "%[6]s"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), tagId.Name(), tag2Id.Name(), tagValue)
}

func TestAcc_TableIdentity(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithIdentityDefault(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.expression"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.expression"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.identity.0.start_num", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.identity.0.step_num", "1"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableColumnWithIdentity(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.expression"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.sequence"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.expression"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.identity.0.start_num"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.identity.0.step_num"),
					// we've dropped the previous identity column and making sure that adding a new column as an identity works
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.identity.0.start_num", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.identity.0.step_num", "4"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableColumnWithIdentityDefault(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_seq" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "NUMBER(38,0)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		identity {
		}
	}
	depends_on = [snowflake_sequence.test_seq]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableColumnWithIdentity(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_seq" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "NUMBER(38,0)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}

	column {
		name = "column4"
		type = "NUMBER(38,0)"
		identity {
			start_num = 2
			step_num = 4
		}
	}
	depends_on = [snowflake_sequence.test_seq]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_TableCollate(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithCollate(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.collate", "en"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.collate", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.collate", ""),
				),
			},
			{
				Config: addColumnWithCollate(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "column4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.collate", "utf8"),
				),
			},
			{
				Config:      alterTableColumnWithIncompatibleCollate(tableId),
				ExpectError: regexp.MustCompile("\"VARCHAR\\(100\\) COLLATE 'fr'\" because they have incompatible collations\\."),
			},
		},
	})
}

func tableColumnWithCollate(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "en"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func addColumnWithCollate(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "en"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
	column {
		name = "column4"
		type = "VARCHAR"
		collate = "utf8"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func alterTableColumnWithIncompatibleCollate(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "fr"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
	column {
		name = "column4"
		type = "VARCHAR"
		collate = "utf8"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_TableRename(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	oldId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	oldComment := random.Comment()
	newComment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithName(oldId, oldComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", oldId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", oldId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", oldComment),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfigWithName(newId, newComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_table.test_table", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableConfigWithName(tableId sdk.SchemaObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
    comment  = "%[4]s"
	column {
		name = "column1"
		type = "VARIANT"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), comment)
}

func TestAcc_Table_MaskingPolicy(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	maskingPolicy1, maskingPolicy1Cleanup := acc.TestClient().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarchar)
	t.Cleanup(maskingPolicy1Cleanup)

	maskingPolicy2, maskingPolicy2Cleanup := acc.TestClient().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarchar)
	t.Cleanup(maskingPolicy2Cleanup)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableWithMaskingPolicy(tableId, maskingPolicy1.ID()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.masking_policy", maskingPolicy1.ID().FullyQualifiedName()),
				),
			},
			// this step proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2186
			{
				Config: tableWithMaskingPolicy(tableId, maskingPolicy2.ID()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.masking_policy", maskingPolicy2.ID().FullyQualifiedName()),
				),
			},
		},
	})
}

func tableWithMaskingPolicy(tableId sdk.SchemaObjectIdentifier, maskingPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment  = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
		masking_policy = "\"%s\".\"%s\".\"%s\""
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), maskingPolicyId.DatabaseName(), maskingPolicyId.SchemaName(), maskingPolicyId.Name())
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTime(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acc.TestClient().Database.UpdateDataRetentionTime(t, database.ID(), 5)
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 10)
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 10, 10),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.UpdateDataRetentionTime(t, database.ID(), 10)
					acc.TestClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 3)
				},
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
			{
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 15),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "15"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 15),
				),
			},
			{
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Schema.UnsetDataRetentionTime(t, schema.ID())
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 10, 10),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 5)
				},
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "0"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 5, 0),
				),
			},
			{
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 5, 3),
				),
			},
		},
	})
}

func tableConfigWithoutDataRetentionTimeInDays(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
    	name = "id"
    	type = "NUMBER(38,0)"
  	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfigWithDataRetentionTimeInDays(tableId sdk.SchemaObjectIdentifier, dataRetentionTimeInDays int) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	data_retention_time_in_days = %[4]d

	column {
    	name = "id"
    	type = "NUMBER(38,0)"
  	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), dataRetentionTimeInDays)
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acc.TestClient().Database.UpdateDataRetentionTime(t, database.ID(), 5)
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Table.SetDataRetentionTime(t, tableId, 20)
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 10)
				},
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 10, 3),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTimeSettingUnsetting(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acc.TestClient().Database.UpdateDataRetentionTime(t, database.ID(), 10)
					acc.TestClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 3)
				},
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
			{
				Config: tableConfigWithDataRetentionTimeInDays(tableId, -1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				Config: tableConfigWithDataRetentionTimeInDays(tableId, -1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
		},
	})
}

func checkDatabaseSchemaAndTableDataRetentionTime(id sdk.SchemaObjectIdentifier, expectedDatabaseRetentionDays int, expectedSchemaRetentionDays int, expectedTableRetentionsDays int) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		database, err := client.Databases.ShowByID(ctx, id.DatabaseId())
		if err != nil {
			return err
		}

		if database.RetentionTime != expectedDatabaseRetentionDays {
			return fmt.Errorf("invalid database retention time, expected: %d, got: %d", expectedDatabaseRetentionDays, database.RetentionTime)
		}

		s, err := client.Schemas.ShowByID(ctx, id.SchemaId())
		if err != nil {
			return err
		}

		// "retention_time" may sometimes be an empty string instead of an integer
		var schemaRetentionTime int64
		{
			rt := s.RetentionTime
			if rt == "" {
				rt = "0"
			}

			schemaRetentionTime, err = strconv.ParseInt(rt, 10, 64)
			if err != nil {
				return err
			}
		}

		if schemaRetentionTime != int64(expectedSchemaRetentionDays) {
			return fmt.Errorf("invalid schema retention time, expected: %d, got: %d", expectedSchemaRetentionDays, schemaRetentionTime)
		}

		table, err := client.Tables.ShowByID(ctx, id)
		if err != nil {
			return err
		}

		if table.RetentionTime != expectedTableRetentionsDays {
			return fmt.Errorf("invalid table retention time, expected: %d, got: %d", expectedTableRetentionsDays, table.RetentionTime)
		}

		return nil
	}
}

// proves issues https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2110 and https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2495
func TestAcc_Table_ClusterBy(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithComplexClusterBy(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.0", "date_trunc('month', LAST_LOAD_TIME)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.1", "COL1"),
				),
			},
		},
	})
}

func tableConfigWithComplexClusterBy(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database            = "%[1]s"
	schema              = "%[2]s"
	name                = "%[3]s"
	cluster_by = ["date_trunc('month', LAST_LOAD_TIME)", "COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
    column {
        name     = "LAST_LOAD_TIME"
        type     = "TIMESTAMP_LTZ(6)"
        nullable = true
    }
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

// TODO [SNOW-1348114]: do not trim the data type (e.g. NUMBER(38,0) -> NUMBER(36,0) diff is ignored); finish the test
// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2588 is fixed
func TestAcc_ColumnTypeChangeWithNonTextType(t *testing.T) {
	t.Skipf("Will be fixed with tables redesign in SNOW-1348114")

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithNumberColumnType(tableId, "NUMBER(38,0)"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "id"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
				),
			},
			{
				Config: tableConfigWithNumberColumnType(tableId, "NUMBER(36,0)"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "id"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(36,0)"),
				),
			},
		},
	})
}

func tableConfigWithNumberColumnType(tableId sdk.SchemaObjectIdentifier, columnType string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database            = "%[1]s"
	schema              = "%[2]s"
	name                = "%[3]s"

	column {
		name = "id"
		type = "%[4]s"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), columnType)
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2733 is fixed
func TestAcc_Table_gh2733(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigGh2733(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "MY_INT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "MY_STRING"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "VARCHAR(16777216)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "MY_DATE"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "MY_DATE2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.type", "TIMESTAMP_NTZ(9)"),
				),
			},
		},
	})
}

func tableConfigGh2733(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
  database            = "%[1]s"
  schema              = "%[2]s"
  name                = "%[3]s"

  column {
    name = "MY_INT"
    type = "int"
    # type  = "NUMBER(38,0)" # Should be equivalent
  }

  column {
    name = "MY_STRING"
    type = "VARCHAR(16777216)"
    # type = "STRING" # Should be equivalent
  }

  column {
    name = "MY_DATE"
    type = "TIMESTAMP_NTZ"
    # type = "TIMESTAMP_NTZ(9)" # Should be equivalent
  }

  column {
    name = "MY_DATE2"
    type = "DATETIME"
    # type = "TIMESTAMP_NTZ" # Equivalent to TIMESTAMP_NTZ
  }
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_Table_migrateFromVersion_0_94_1(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resourceName := "snowflake_table.test_table"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            tableConfig(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", tableId.Name()),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", tableId.FullyQualifiedName()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfig(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", tableId.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", tableId.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceName, "qualified_name"),
				),
			},
		},
	})
}

func TestAcc_Table_SuppressQuotingOnDefaultSequence_issue2644(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resourceName := "snowflake_table.test_table"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             tableConfigWithSequence(tableId),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigWithSequence(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.default.0.sequence", tableId.FullyQualifiedName()),
				),
			},
		},
	})
}

func tableConfigWithSequence(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_sequence" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	data_retention_time_in_days = 1
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "NUMBER"
		default {
			sequence = "%[1]s.%[2]s.%[3]s"
		}
	}
	depends_on = [snowflake_sequence.test_sequence]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_Table_issue3007_textColumn(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	defaultVarchar := fmt.Sprintf("VARCHAR(%d)", datatypes.DefaultVarcharLength)

	resourceName := "snowflake_table.test_table"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR(3)"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "VARCHAR(3)"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR(256)"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("VARCHAR(3)"), sdk.String("VARCHAR(256)")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "VARCHAR(256)"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("VARCHAR(256)"), sdk.String("VARCHAR")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, defaultVarchar),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "text"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
		},
	})
}

// TODO [SNOW-1348114]: visit with table rework (e.g. changing scale is not supported: err 040052 (22000): SQL compilation error: cannot change column SOME_COLUMN from type NUMBER(38,0) to NUMBER(11,2) because changing the scale of a number is not supported.)
func TestAcc_Table_issue3007_numberColumn(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	defaultNumber := fmt.Sprintf("NUMBER(%d,%d)", datatypes.DefaultNumberPrecision, datatypes.DefaultNumberScale)

	resourceName := "snowflake_table.test_table"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "NUMBER(38,0)"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER(11)"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("NUMBER(38,0)"), sdk.String("NUMBER(11)")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "NUMBER(11,0)"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("NUMBER(11,0)"), sdk.String("NUMBER")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, defaultNumber),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "decimal"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
		},
	})
}

func tableConfigIssue3007(tableId sdk.SchemaObjectIdentifier, dataType string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
    database = "%[1]s"
    schema   = "%[2]s"
    name     = "%[3]s"
    comment  = "Issue 3007 confirmation"

    column {
        name = "ID"
        type = "NUMBER(11,2)"
    }

    column {
        name = "SOME_COLUMN"
        type = "%[4]s"
    }
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), dataType)
}

func TestAcc_Table_SchemaRemovedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetLegacyConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("1.0.5"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_table.test_table", plancheck.ResourceActionCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_table.test_table", plancheck.ResourceActionNoop),
					},
				},
				Config: tableConfig(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "id", helpers.EncodeSnowflakeID(tableId)),
				),
			},
			{
				ExternalProviders: acc.ExternalProviderWithExactVersion("1.0.5"),
				PreConfig:         func() { schemaCleanup() },
				Config:            tableConfig(tableId),
				ExpectError:       regexp.MustCompile("object does not exist or not authorized"),
			},
			// The New version removes table from the state which ends up with create operation
			// and the "error creating table" error (because of the missing underlying schema).
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfig(tableId),
				ExpectError:              regexp.MustCompile("error creating table"),
			},
		},
	})
}
