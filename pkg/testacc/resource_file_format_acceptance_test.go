//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FileFormatCSV(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigCSV(id, ";", "'", "Terraform acceptance test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "CSV"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "record_delimiter", "\r"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "field_delimiter", ";"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "file_extension", ".ssv"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "parse_header", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_blank_lines", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "date_format", "YYY-MM-DD"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "time_format", "HH24:MI"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "timestamp_format", "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "binary_format", "UTF8"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "escape", "\\"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "escape_unenclosed_field", "!"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "field_optionally_enclosed_by", "'"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "2"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.1", ""),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "error_on_column_count_mismatch", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "empty_field_as_null", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_byte_order_mark", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "encoding", "UTF-16"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
			// UPDATE
			{
				Config: fileFormatConfigCSV(id, ",", "'", "Terraform acceptance test 2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "field_delimiter", ","),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test 2"),
				),
			},
			// UPDATE: field_optionally_enclosed_by can take the value NONE
			{
				Config: fileFormatConfigCSV(id, ",", "NONE", "Terraform acceptance test 2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "field_optionally_enclosed_by", "NONE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatJSON(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigJSON(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "JSON"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "date_format", "YYY-MM-DD"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "time_format", "HH24:MI"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "timestamp_format", "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "binary_format", "UTF8"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "file_extension", ".jsn"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "enable_octal", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "allow_duplicate", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "strip_outer_array", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "strip_null_values", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "ignore_utf8_errors", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_byte_order_mark", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_FileFormatAvro(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigAvro(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "AVRO"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_FileFormatORC(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigORC(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "ORC"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_FileFormatParquet(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigParquet(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "PARQUET"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "SNAPPY"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "binary_as_text", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_FileFormatXML(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigXML(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "XML"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "preserve_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "strip_outer_element", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "disable_snowflake_data", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "disable_auto_convert", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_byte_order_mark", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

// The following tests check that Terraform will accept the default values generated at creation and not drift.
// See https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1706
func TestAcc_FileFormatCSVDefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(id, "CSV"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "CSV"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatJSONDefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(id, "JSON"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "JSON"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatAVRODefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(id, "AVRO"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "AVRO"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatORCDefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(id, "ORC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "ORC"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatPARQUETDefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(id, "PARQUET"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "PARQUET"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatXMLDefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(id, "XML"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "XML"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAcc_FileFormat_issue1947 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1947 issue.
func TestAcc_FileFormat_issue1947(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(id, "XML"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
				),
			},
			/*
			 * Before the fix this step resulted in
			 *     Error: only one of Rename or Set can be set at once.
			 * which matches the issue description exactly
			 */
			{
				// we set param which is not right for XML but this allows us to run update on terraform apply
				Config: fileFormatConfigFullDefaultsWithAdditionalParam(id, "XML"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
				),
			},
		},
	})
}

func TestAcc_FileFormat_Rename(t *testing.T) {
	oldId := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()
	resourceName := "snowflake_file_format.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigWithComment(oldId, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", oldId.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", oldId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: fileFormatConfigWithComment(newId, newComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", newId.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "comment", newComment),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_FileFormat_CSV_EmptyNullIf_NoPlanDiff(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigEmptyNullIf(id, "CSV"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "CSV"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "0"),
				),
			},
		},
	})
}

func TestAcc_FileFormat_JSON_EmptyNullIf_NoPlanDiff(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigEmptyNullIf(id, "JSON"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "JSON"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "0"),
				),
			},
		},
	})
}

func TestAcc_FileFormat_AVRO_EmptyNullIf_NoPlanDiff(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigEmptyNullIf(id, "AVRO"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "AVRO"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "0"),
				),
			},
		},
	})
}

func TestAcc_FileFormat_ORC_EmptyNullIf_NoPlanDiff(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigEmptyNullIf(id, "ORC"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "ORC"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "0"),
				),
			},
		},
	})
}

func TestAcc_FileFormat_PARQUET_EmptyNullIf_NoPlanDiff(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FileFormat),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigEmptyNullIf(id, "PARQUET"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "PARQUET"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "0"),
				),
			},
		},
	})
}

func fileFormatConfigCSV(id sdk.SchemaObjectIdentifier, fieldDelimiter string, fieldOptionallyEnclosedBy string, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "CSV"
	compression = "GZIP"
	record_delimiter = "\r"
	field_delimiter = "%s"
	file_extension = ".ssv"
	parse_header = true
	skip_blank_lines = true
	date_format = "YYY-MM-DD"
	time_format = "HH24:MI"
	timestamp_format = "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"
	binary_format = "UTF8"
	escape = "\\"
	escape_unenclosed_field = "!"
	trim_space = true
	field_optionally_enclosed_by = "%s"
	null_if = ["NULL", ""]
	error_on_column_count_mismatch = true
	replace_invalid_characters = true
	empty_field_as_null = false
	skip_byte_order_mark = false
	encoding = "UTF-16"
	comment = "%s"
}
`, id.Name(), id.DatabaseName(), id.SchemaName(), fieldDelimiter, fieldOptionallyEnclosedBy, comment)
}

func fileFormatConfigJSON(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "JSON"
	compression = "GZIP"
	date_format = "YYY-MM-DD"
	time_format = "HH24:MI"
	timestamp_format = "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"
	binary_format = "UTF8"
	trim_space = true
	null_if = ["NULL"]
	file_extension = ".jsn"
	enable_octal = true
	allow_duplicate = true
	strip_outer_array = true
	strip_null_values = true
	ignore_utf8_errors = true
	skip_byte_order_mark = false
	comment = "Terraform acceptance test"
}
`, id.Name(), id.DatabaseName(), id.SchemaName())
}

func fileFormatConfigAvro(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "AVRO"
	compression = "GZIP"
	trim_space = true
	null_if = ["NULL"]
	comment = "Terraform acceptance test"
}
`, id.Name(), id.DatabaseName(), id.SchemaName())
}

func fileFormatConfigORC(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "ORC"
	trim_space = true
	null_if = ["NULL"]
	comment = "Terraform acceptance test"
}
`, id.Name(), id.DatabaseName(), id.SchemaName())
}

func fileFormatConfigParquet(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "PARQUET"
	compression = "SNAPPY"
	binary_as_text = true
	trim_space = true
	null_if = ["NULL"]
	comment = "Terraform acceptance test"
}
`, id.Name(), id.DatabaseName(), id.SchemaName())
}

func fileFormatConfigXML(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "XML"
	compression = "GZIP"
	ignore_utf8_errors = true
	preserve_space = true
	strip_outer_element = true
	disable_snowflake_data =  true
	disable_auto_convert =  true
	skip_byte_order_mark = false
	comment = "Terraform acceptance test"
}
`, id.Name(), id.DatabaseName(), id.SchemaName())
}

func fileFormatConfigFullDefaults(id sdk.SchemaObjectIdentifier, formatType string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "%s"
}
`, id.Name(), id.DatabaseName(), id.SchemaName(), formatType)
}

func fileFormatConfigWithComment(id sdk.SchemaObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "XML"
	comment = "%s"
}
`, id.Name(), id.DatabaseName(), id.SchemaName(), comment)
}

func fileFormatConfigFullDefaultsWithAdditionalParam(id sdk.SchemaObjectIdentifier, formatType string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	format_type = "%s"
    encoding = "UTF-16"
}
`, id.Name(), id.DatabaseName(), id.SchemaName(), formatType)
}

func fileFormatConfigEmptyNullIf(id sdk.SchemaObjectIdentifier, formatType string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
  name        = "%s"
  database    = "%s"
  schema      = "%s"
  format_type = "%s"
  null_if     = []
}
`, id.Name(), id.DatabaseName(), id.SchemaName(), formatType)
}
