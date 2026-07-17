//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FailoverGroupBasic(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	currentAccountId := testClient().Account.GetAccountIdentifier(t)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupBasic(id, currentAccountId, businessCriticalAccountId, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
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
	currentAccountId := testClient().Account.GetAccountIdentifier(t)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithInterval(id, currentAccountId, businessCriticalAccountId, 20, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "20"),
				),
			},
			{
				Config: failoverGroupWithNoWarehouse(id, currentAccountId, businessCriticalAccountId, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "3"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
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
	currentAccountId := testClient().Account.GetAccountIdentifier(t)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithInterval(id, currentAccountId, businessCriticalAccountId, 10, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "10"),
				),
			},
			// Update Interval
			{
				Config: failoverGroupWithInterval(id, currentAccountId, businessCriticalAccountId, 20, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "20"),
				),
			},
			// Change to Cron Expression
			{
				Config: failoverGroupWithCronExpression(id, currentAccountId, businessCriticalAccountId, "0 0 10-20 * TUE,THU", testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
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
				Config: failoverGroupWithCronExpression(id, currentAccountId, businessCriticalAccountId, "0 0 5-20 * TUE,THU", testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
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
				Config: failoverGroupWithoutReplicationSchedule(id, currentAccountId, businessCriticalAccountId, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "0"),
				),
			},
			// Change to Interval
			{
				Config: failoverGroupWithInterval(id, currentAccountId, businessCriticalAccountId, 10, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
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
	currentAccountId := testClient().Account.GetAccountIdentifier(t)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithAccountParameters(id, currentAccountId, businessCriticalAccountId, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "5"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
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
	currentAccountId := testClient().Account.GetAccountIdentifier(t)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupBasic(id, currentAccountId, businessCriticalAccountId, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
				),
			},
			{
				Config: failoverGroupWithChanges(id, currentAccountId, businessCriticalAccountId, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
				),
			},
		},
	})
}

// TODO [SNOW-1348343]: current account id passed on purpose; handle properly during resource rework
func failoverGroupBasic(failoverGroupId sdk.AccountObjectIdentifier, currentAccountId sdk.AccountIdentifier, allowedAccountId sdk.AccountIdentifier, databaseId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES", "DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s", "%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS", "API INTEGRATIONS", "STORAGE INTEGRATIONS", "EXTERNAL ACCESS INTEGRATIONS", "NOTIFICATION INTEGRATIONS"]
	replication_schedule {
		cron {
			expression = "0 0 10-20 * TUE,THU"
			time_zone = "UTC"
		}
	}
}
`, failoverGroupId.Name(), currentAccountId.Name(), allowedAccountId.Name(), databaseId.Name())
}

func failoverGroupWithInterval(id sdk.AccountObjectIdentifier, currentAccountId sdk.AccountIdentifier, allowedAccountId sdk.AccountIdentifier, interval int, databaseId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s", "%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		interval = %d
	}
}
`, id.Name(), currentAccountId.Name(), allowedAccountId.Name(), databaseId.Name(), interval)
}

func failoverGroupWithoutReplicationSchedule(id sdk.AccountObjectIdentifier, currentAccountId sdk.AccountIdentifier, allowedAccountId sdk.AccountIdentifier, databaseId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s", "%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
}
`, id.Name(), currentAccountId.Name(), allowedAccountId.Name(), databaseId.Name())
}

func failoverGroupWithNoWarehouse(id sdk.AccountObjectIdentifier, currentAccountId sdk.AccountIdentifier, allowedAccountId sdk.AccountIdentifier, interval int) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s", "%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		interval = %d
	}
}
`, id.Name(), currentAccountId.Name(), allowedAccountId.Name(), interval)
}

func failoverGroupWithCronExpression(id sdk.AccountObjectIdentifier, currentAccountId sdk.AccountIdentifier, allowedAccountId sdk.AccountIdentifier, expression string, databaseId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s", "%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		cron {
			expression = "%s"
			time_zone  = "UTC"
		}
	}
}
`, id.Name(), currentAccountId.Name(), allowedAccountId.Name(), databaseId.Name(), expression)
}

func failoverGroupWithAccountParameters(id sdk.AccountObjectIdentifier, currentAccountId sdk.AccountIdentifier, allowedAccountId sdk.AccountIdentifier, databaseId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["ACCOUNT PARAMETERS", "WAREHOUSES", "DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s", "%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		cron {
			expression = "0 0 10-20 * TUE,THU"
			time_zone = "UTC"
		}
	}
}
`, id.Name(), currentAccountId.Name(), allowedAccountId.Name(), databaseId.Name())
}

func failoverGroupWithChanges(id sdk.AccountObjectIdentifier, currentAccountId sdk.AccountIdentifier, allowedAccountId sdk.AccountIdentifier, interval int) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%[1]s"
	object_types = ["DATABASES", "INTEGRATIONS"]
	allowed_accounts= ["%[2]s", "%[3]s"]
	allowed_integration_types = ["NOTIFICATION INTEGRATIONS"]
	replication_schedule {
		interval = %d
	}
}
`, id.Name(), currentAccountId.Name(), allowedAccountId.Name(), interval)
}

func TestAcc_FailoverGroup_UpdateAllowedAccounts(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	// We need to use the secondary account to test the failover group. Otherwise, we will get an error:
	// Error: error removing allowed accounts for failover group TEST err = 003909 (55000): Disabling replication of the replication group to the account in which the primary currently resides is not allowed.
	accountIdentifier := testClient().Context.CurrentAccountIdentifier(t)

	secondaryAccountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	secondaryAccountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(secondaryAccountName)
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.FailoverGroupResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// Create with allowed_accounts using old provider version
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            config.FromModels(t, providerModel) + failoverGroupWithAccounts(id.Name(), TestDatabaseName, accountIdentifier, secondaryAccountIdentifier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", accountIdentifier.Name()),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", secondaryAccountIdentifier.Name()),
				),
			},
			// Remove an allowed account with old provider version - proves the bug
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            config.FromModels(t, providerModel) + failoverGroupWithAccounts(id.Name(), TestDatabaseName, accountIdentifier),
				ExpectError:       regexp.MustCompile("error removing allowed accounts"),
			},
			// Remove it externally
			{
				PreConfig: func() {
					testClient().FailoverGroup.RemoveAllowedAccounts(t, id, secondaryAccountIdentifier)
				},
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            config.FromModels(t, providerModel) + failoverGroupWithAccounts(id.Name(), TestDatabaseName, accountIdentifier),
			},
			// Add it back again - proves the bug
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            config.FromModels(t, providerModel) + failoverGroupWithAccounts(id.Name(), TestDatabaseName, accountIdentifier, secondaryAccountIdentifier),
				ExpectError:       regexp.MustCompile("error adding allowed accounts"),
			},
			// Upgrade to current provider version
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   failoverGroupWithAccounts(id.Name(), TestDatabaseName, accountIdentifier, secondaryAccountIdentifier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", accountIdentifier.Name()),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", secondaryAccountIdentifier.Name()),
				),
			},
			// Remove an allowed account with current provider version - proves the fix
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   failoverGroupWithAccounts(id.Name(), TestDatabaseName, accountIdentifier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", accountIdentifier.Name()),
				),
			},
			// Add it back again - proves the fix
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   failoverGroupWithAccounts(id.Name(), TestDatabaseName, accountIdentifier, secondaryAccountIdentifier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", accountIdentifier.Name()),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", secondaryAccountIdentifier.Name()),
				),
			},
		},
	})
}

func failoverGroupWithAccounts(name, databaseName string, accountIds ...sdk.AccountIdentifier) string {
	accountNames := collections.Map(accountIds, func(account sdk.AccountIdentifier) string {
		return fmt.Sprintf(`"%s"`, account.Name())
	})
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["DATABASES"]
	allowed_databases = ["%s"]
	allowed_accounts = [ %s ]
}
`, name, databaseName, strings.Join(accountNames, ", "))
}

func TestAcc_FailoverGroup_InvalidObjectType(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// An object type outside the supported plural object type set is rejected during create.
			{
				Config:      failoverGroupWithObjectTypes(id.Name(), `"DATABASES", "INVALID OBJECT TYPE;"`),
				ExpectError: regexp.MustCompile("invalid plural object type: INVALID OBJECT TYPE;"),
			},
		},
	})
}

func failoverGroupWithObjectTypes(name, objectTypes string) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name         = "%[1]s"
	object_types = [%[2]s]
}
`, name, objectTypes)
}
