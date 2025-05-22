//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FileFormats(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	fileFormatId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatsInSchema(fileFormatId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "database", fileFormatId.DatabaseName()),
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "schema", fileFormatId.SchemaName()),
					resource.TestCheckResourceAttrSet("data.snowflake_file_formats.t", "file_formats.#"),
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "file_formats.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "file_formats.0.name", fileFormatId.Name()),
				),
			},
		},
	})
}

func TestAcc_FileFormatsEmpty(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatsInSchemaWithoutCreation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_file_formats.t", "file_formats.#"),
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "file_formats.#", "0"),
				),
			},
		},
	})
}

func fileFormatsInSchema(fileFormatId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
	resource snowflake_file_format "t"{
		name 	 	= "%[3]s"
		database 	= "%[1]s"
		schema 	 	= "%[2]s"
		format_type = "CSV"
		compression = "GZIP"
		record_delimiter = "\r"
		field_delimiter = ";"
		file_extension = ".ssv"
		skip_header = 1
		skip_blank_lines = true
		date_format = "YYY-MM-DD"
		time_format = "HH24:MI"
		timestamp_format = "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"
		binary_format = "UTF8"
		escape = "\\"
		escape_unenclosed_field = "!"
		trim_space = true
		field_optionally_enclosed_by = "'"
		null_if = ["NULL"]
		error_on_column_count_mismatch = true
		replace_invalid_characters = true
		empty_field_as_null = false
		skip_byte_order_mark = false
		encoding = "UTF-16"
		comment = "Terraform acceptance test"
	}

	data snowflake_file_formats "t" {
		database = "%[1]s"
		schema = "%[2]s"
		depends_on = [snowflake_file_format.t]
	}
	`, fileFormatId.DatabaseName(), fileFormatId.SchemaName(), fileFormatId.Name())
}

func fileFormatsInSchemaWithoutCreation() string {
	return fmt.Sprintf(`
	data snowflake_file_formats "t" {
		database = "%[1]s"
		schema 	 = "%[2]s"
	}
	`, acc.TestDatabaseName, acc.TestSchemaName)
}
