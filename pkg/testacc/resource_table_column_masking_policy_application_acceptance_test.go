//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TableColumnMaskingPolicyApplication(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	maskingPolicyId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyApplicationTestConfig(tableId, maskingPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_column_masking_policy_application.mpa", "table", tableId.FullyQualifiedName()),
				),
			},
			{
				ResourceName:      "snowflake_table_column_masking_policy_application.mpa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func maskingPolicyApplicationTestConfig(tableId sdk.SchemaObjectIdentifier, maskingPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_masking_policy" "test" {
	database           = "%[1]s"
	schema             = "%[2]s"
	name               = "%[4]s"
	argument {
		name = "val"
		type = "VARCHAR"
	}
	body = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type   = "VARCHAR"
	comment            = "Terraform acceptance test"
}

resource "snowflake_table" "table" {
	database           = "%[1]s"
	schema             = "%[2]s"
	name               = "%[3]s"

	column {
	  name     = "secret"
	  type     = "VARCHAR(16777216)"
	}

	lifecycle {
		ignore_changes = [column[0].masking_policy]
	}
}

resource "snowflake_table_column_masking_policy_application" "mpa" {
	table          = snowflake_table.table.fully_qualified_name
	column         = "secret"
	masking_policy = snowflake_masking_policy.test.fully_qualified_name
}`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), maskingPolicyId.Name())
}

// proves the fix for https://github.com/snowflakedb/terraform-provider-snowflake/issues/4608
func TestAcc_TableColumnMaskingPolicyApplication_BCR2026_02(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	maskingPolicyId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.1"),
				Config:            maskingPolicyApplicationBcr202602TestConfig(tableId, maskingPolicyId),
				ExpectError:       regexp.MustCompile("expected 13 destination arguments in Scan, not 12"),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   maskingPolicyApplicationTestConfig(tableId, maskingPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_column_masking_policy_application.mpa", "table", tableId.FullyQualifiedName()),
				),
			},
		},
	})
}

func maskingPolicyApplicationBcr202602TestConfig(tableId sdk.SchemaObjectIdentifier, maskingPolicyId sdk.SchemaObjectIdentifier) string {
	return `provider "snowflake" {
  preview_features_enabled = ["snowflake_table_resource", "snowflake_table_column_masking_policy_application_resource"]
}
` + maskingPolicyApplicationTestConfig(tableId, maskingPolicyId)
}
