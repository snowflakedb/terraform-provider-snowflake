//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkRule(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkRule),
		Steps: []resource.TestStep{
			// basic
			{
				Config: networkRuleIpv4(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "type", "IPV4"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "mode", "INGRESS"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "value_list.#", "2"),
				),
			},
			//// IMPORT - all fields are non-empty
			{
				ResourceName:      "snowflake_network_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			//// CHANGE PROPERTIES - set to empty
			{
				Config: networkRuleIpv4Empty(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "type", "IPV4"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "mode", "INGRESS"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "value_list.#", "0"),
				),
			},
			// IMPORT - incomplete
			{
				ResourceName:      "snowflake_network_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: networkRuleHost(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "type", "HOST_PORT"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "mode", "EGRESS"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "value_list.#", "2"),
				),
			},
		},
	})
}

func networkRuleIpv4(id sdk.SchemaObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_network_rule" "test" {
	database        = "%[1]s"
	schema          = "%[2]s"
	name            = "%[3]s"
	comment         = "%[4]s"
    type            = "IPV4"
    mode			= "INGRESS"
	value_list      = ["192.168.0.100/24", "29.254.123.20"]
}
`, id.DatabaseName(), id.SchemaName(), id.Name(), comment)
}

func networkRuleIpv4Empty(id sdk.SchemaObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_network_rule" "test" {
	database        = "%[1]s"
	schema          = "%[2]s"
	name            = "%[3]s"
	comment         = "%[4]s"
    type            = "IPV4"
    mode			= "INGRESS"
	value_list      = []
}
`, id.DatabaseName(), id.SchemaName(), id.Name(), comment)
}

func networkRuleHost(id sdk.SchemaObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_network_rule" "test" {
	database        = "%[1]s"
	schema          = "%[2]s"
	name            = "%[3]s"
	comment         = "%[4]s"
    type            = "HOST_PORT"
    mode			= "EGRESS"
	value_list      = ["example.com", "company.com:443"]
}
`, id.DatabaseName(), id.SchemaName(), id.Name(), comment)
}

func TestAcc_NetworkRule_migrateFromVersion_0_94_1(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	resourceName := "snowflake_network_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            networkRuleIpv4(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   networkRuleIpv4(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceName, "qualified_name"),
				),
			},
		},
	})
}
