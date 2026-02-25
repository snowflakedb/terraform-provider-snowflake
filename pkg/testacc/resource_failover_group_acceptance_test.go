//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FailoverGroupBasic(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupBasic(id.Name(), accountName, TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_integration_types.*", "SECURITY INTEGRATIONS"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_integration_types.*", "API INTEGRATIONS"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_integration_types.*", "STORAGE INTEGRATIONS"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_integration_types.*", "EXTERNAL ACCESS INTEGRATIONS"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_integration_types.*", "NOTIFICATION INTEGRATIONS"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.expression", "0 0 10-20 * TUE,THU"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.time_zone", "UTC"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_failover_group.fg",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_edition_check"},
			},
		},
	})
}

func TestAcc_FailoverGroupRemoveObjectTypes(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithInterval(id.Name(), accountName, 20, TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "20"),
				),
			},
			{
				Config: failoverGroupWithNoWarehouse(id.Name(), accountName, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "3"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "20"),
				),
			},
		},
	})
}

func TestAcc_FailoverGroupInterval(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithInterval(id.Name(), accountName, 10, TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "10"),
				),
			},
			// Update Interval
			{
				Config: failoverGroupWithInterval(id.Name(), accountName, 20, TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "20"),
				),
			},
			// Change to Cron Expression
			{
				Config: failoverGroupWithCronExpression(id.Name(), accountName, "0 0 10-20 * TUE,THU", TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.expression", "0 0 10-20 * TUE,THU"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.time_zone", "UTC"),
				),
			},
			// Update Cron Expression
			{
				Config: failoverGroupWithCronExpression(id.Name(), accountName, "0 0 5-20 * TUE,THU", TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.expression", "0 0 5-20 * TUE,THU"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.time_zone", "UTC"),
				),
			},
			// Remove replication schedule
			{
				Config: failoverGroupWithoutReplicationSchedule(id.Name(), accountName, TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "0"),
				),
			},
			// Change to Interval
			{
				Config: failoverGroupWithInterval(id.Name(), accountName, 10, TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "10"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_failover_group.fg",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_edition_check"},
			},
		},
	})
}

func TestAcc_FailoverGroup_issue2517(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithAccountParameters(id.Name(), accountName, TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "5"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.expression", "0 0 10-20 * TUE,THU"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.time_zone", "UTC"),
				),
			},
		},
	})
}

func TestAcc_FailoverGroup_issue2544(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupBasic(id.Name(), accountName, TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
				),
			},
			{
				Config: failoverGroupWithChanges(id.Name(), accountName, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
				),
			},
		},
	})
}

func failoverGroupBasic(randomCharacters, accountName, databaseName string) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES", "DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS", "API INTEGRATIONS", "STORAGE INTEGRATIONS", "EXTERNAL ACCESS INTEGRATIONS", "NOTIFICATION INTEGRATIONS"]
	replication_schedule {
		cron {
			expression = "0 0 10-20 * TUE,THU"
			time_zone = "UTC"
		}
	}
}
`, randomCharacters, accountName, databaseName)
}

func failoverGroupWithInterval(randomCharacters, accountName string, interval int, databaseName string) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		interval = %d
	}
}
`, randomCharacters, accountName, databaseName, interval)
}

func failoverGroupWithoutReplicationSchedule(randomCharacters, accountName string, databaseName string) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
}
`, randomCharacters, accountName, databaseName)
}

func failoverGroupWithNoWarehouse(randomCharacters, accountName string, interval int) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		interval = %d
	}
}
`, randomCharacters, accountName, interval)
}

func failoverGroupWithCronExpression(randomCharacters, accountName, expression, databaseName string) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		cron {
			expression = "%s"
			time_zone  = "UTC"
		}
	}
}
`, randomCharacters, accountName, databaseName, expression)
}

func failoverGroupWithAccountParameters(randomCharacters, accountName, databaseName string) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["ACCOUNT PARAMETERS", "WAREHOUSES", "DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		cron {
			expression = "0 0 10-20 * TUE,THU"
			time_zone = "UTC"
		}
	}
}
`, randomCharacters, accountName, databaseName)
}

func failoverGroupWithChanges(randomCharacters string, accountName string, interval int) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%[1]s"
	object_types = ["DATABASES", "INTEGRATIONS"]
	allowed_accounts= ["%[2]s"]
	allowed_integration_types = ["NOTIFICATION INTEGRATIONS"]
	replication_schedule {
		interval = %d
	}
}
`, randomCharacters, accountName, interval)
}

// TODO: secondary account?
// Error: error removing allowed accounts for failover group BYPRYDAT_4DD10EDF_9310_FE23_A4E3_A969CE8A321E err = 003909 (55000): Disabling replication of the replication group to the account in which the primary currently resides is not allowed.
func TestAcc_FailoverGroup_UpdateAllowedAccounts(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	// providerModel := providermodel.SnowflakeProvider().
	// 	WithPreviewFeaturesEnabled(string(previewfeatures.FailoverGroupResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// // Create with allowed_accounts using old provider version
			// {
			// 	ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
			// 	Config:            config.FromModels(t, providerModel) + failoverGroupBasic(id.Name(), accountName, TestDatabaseName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
			// 		resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
			// 		resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", accountName),
			// 	),
			// },
			// // Remove allowed_accounts with old provider version - proves the bug
			// {
			// 	ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
			// 	Config:            config.FromModels(t, providerModel) + failoverGroupWithRequiredFields(id.Name()),
			// 	ExpectError:       regexp.MustCompile("error removing allowed accounts"),
			// },
			// // Add it back again - proves the bug
			// {
			// 	ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
			// 	Config:            config.FromModels(t, providerModel) + failoverGroupBasic(id.Name(), accountName, TestDatabaseName),
			// 	ExpectError:       regexp.MustCompile("error adding allowed accounts"),
			// },
			// Upgrade to current provider version
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   failoverGroupWithAccount(id.Name(), TestDatabaseName, accountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
				),
			},
			// Remove allowed_accounts with current provider version - proves the fix
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   failoverGroupBare(id.Name(), TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "0"),
				),
			},
			// Add it back again - proves the fix
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   failoverGroupWithAccount(id.Name(), TestDatabaseName, accountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
				),
			},
		},
	})
}

func failoverGroupBare(name, databaseName string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	preview_features_enabled = ["snowflake_failover_group_resource"]
}

resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["DATABASES"]
	allowed_databases = ["%s"]
}
`, name, databaseName)
}

func failoverGroupWithAccount(name, databaseName, accountName string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	preview_features_enabled = ["snowflake_failover_group_resource"]
}

resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["DATABASES"]
	allowed_databases = ["%s"]
	allowed_accounts = ["%s"]
}
`, name, databaseName, accountName)
}
