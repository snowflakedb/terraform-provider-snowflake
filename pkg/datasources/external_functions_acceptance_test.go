//go:build !account_level_tests

package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalFunctions_basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"schema":                    config.StringVariable(acc.TestSchemaName),
			"name":                      config.StringVariable(id.Name()),
			"api_allowed_prefixes":      config.ListVariable(config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/")),
			"url_of_proxy_and_resource": config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
			"comment":                   config.StringVariable(comment),
		}
	}
	configVariables := m()

	dataSourceName := "data.snowflake_external_functions.external_functions"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalFunctions/basic"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(dataSourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(dataSourceName, "external_functions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "external_functions.0.name", id.Name()),
					resource.TestCheckResourceAttr(dataSourceName, "external_functions.0.database", id.DatabaseName()),
					resource.TestCheckResourceAttr(dataSourceName, "external_functions.0.schema", id.SchemaName()),
					resource.TestCheckResourceAttr(dataSourceName, "external_functions.0.comment", comment),
					resource.TestCheckResourceAttrSet(dataSourceName, "external_functions.0.language"),
				),
			},
		},
	})
}

func TestAcc_ExternalFunctions_no_database(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"schema":                    config.StringVariable(acc.TestSchemaName),
			"name":                      config.StringVariable(id.Name()),
			"api_allowed_prefixes":      config.ListVariable(config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/")),
			"url_of_proxy_and_resource": config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
			"comment":                   config.StringVariable(comment),
		}
	}

	dataSourceName := "data.snowflake_external_functions.external_functions"
	configVariables := m()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalFunctions/no_filter"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "external_functions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "external_functions.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "external_functions.0.comment"),
					resource.TestCheckResourceAttrSet(dataSourceName, "external_functions.0.language"),
				),
			},
		},
	})
}
