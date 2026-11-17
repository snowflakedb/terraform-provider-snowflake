//go:build account_level_tests

// Account object types other than SHARES and DATABASES (e.g. ROLES, INTEGRATIONS, WAREHOUSES) can be a part of only one
// failover or replication group per account. Snowflake rejects the second group with 500070 "Account object types, other than
// SHARES and DATABASES, can only be within one failover group or replication group" (documented as "Objects other than databases,
// external volumes, and shares must be in the same failover group"). The restriction is on the type, so unique names do not help.
// The tests below use such types (ROLES, INTEGRATIONS, WAREHOUSES, ACCOUNT PARAMETERS).
// A database can also be added to only one failover or replication group, and these tests put the shared test database into groups.
// Failover and replication groups share these restrictions, so the tests of both are account-level tests (they cannot run in parallel
// with each other, other PRs, or the integration/acceptance suites).

package testacc

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FailoverGroupBasic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupBasic(id, secondaryAccountId, testClient().Ids.DatabaseId()),
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
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "USING CRON 0 0 10-20 * TUE,THU UTC"),
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
	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithInterval(id, secondaryAccountId, 20, testClient().Ids.DatabaseId(), false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "20 MINUTE"),
				),
			},
			{
				Config: failoverGroupWithNoWarehouse(id, secondaryAccountId, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "3"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "20 MINUTE"),
				),
			},
		},
	})
}

func TestAcc_FailoverGroupInterval(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithInterval(id, secondaryAccountId, 10, testClient().Ids.DatabaseId(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "10 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "optimized_refresh", "true"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "rpo_assurance", "false"),
				),
			},
			// Update Interval
			{
				Config: failoverGroupWithInterval(id, secondaryAccountId, 20, testClient().Ids.DatabaseId(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "20 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "optimized_refresh", "true"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "rpo_assurance", "false"),
				),
			},
			// Change to Cron Expression
			{
				Config: failoverGroupWithCronExpression(id, secondaryAccountId, "0 0 10-20 * TUE,THU", testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "USING CRON 0 0 10-20 * TUE,THU UTC"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "optimized_refresh", "false"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "rpo_assurance", "false"),
				),
			},
			// Update Cron Expression
			{
				Config: failoverGroupWithCronExpression(id, secondaryAccountId, "0 0 5-20 * TUE,THU", testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "USING CRON 0 0 5-20 * TUE,THU UTC"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "optimized_refresh", "false"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "rpo_assurance", "false"),
				),
			},
			// Remove replication schedule
			{
				Config: failoverGroupWithoutReplicationSchedule(id, secondaryAccountId, testClient().Ids.DatabaseId(), false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", ""),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "optimized_refresh", "false"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "rpo_assurance", "false"),
				),
			},
			{
				Config: failoverGroupWithoutReplicationSchedule(id, secondaryAccountId, testClient().Ids.DatabaseId(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", ""),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "optimized_refresh", "false"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "rpo_assurance", "true"),
				),
			},
			// Change to Interval
			{
				Config: failoverGroupWithInterval(id, secondaryAccountId, 10, testClient().Ids.DatabaseId(), false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "10 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "optimized_refresh", "false"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "rpo_assurance", "false"),
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
	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithAccountParameters(id, secondaryAccountId, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "5"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule", "USING CRON 0 0 10-20 * TUE,THU UTC"),
				),
			},
		},
	})
}

func TestAcc_FailoverGroup_issue2544(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupBasic(id, secondaryAccountId, testClient().Ids.DatabaseId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
				),
			},
			{
				Config: failoverGroupWithChanges(id, secondaryAccountId, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
				),
			},
		},
	})
}

func failoverGroupBasic(failoverGroupId sdk.AccountObjectIdentifier, allowedAccountId sdk.AccountIdentifier, databaseId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES", "DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts = ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS", "API INTEGRATIONS", "STORAGE INTEGRATIONS", "EXTERNAL ACCESS INTEGRATIONS", "NOTIFICATION INTEGRATIONS"]
	replication_schedule = "USING CRON 0 0 10-20 * TUE,THU UTC"
}
`, failoverGroupId.Name(), allowedAccountId.Name(), databaseId.Name())
}

func failoverGroupWithInterval(id sdk.AccountObjectIdentifier, allowedAccountId sdk.AccountIdentifier, interval int, databaseId sdk.AccountObjectIdentifier, optimizedRefresh bool) string {
	optimizedRefreshConfig := ""
	if optimizedRefresh {
		optimizedRefreshConfig = "\n\toptimized_refresh = true"
	}
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts = ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule = "%d MINUTE"%s
}
`, id.Name(), allowedAccountId.Name(), databaseId.Name(), interval, optimizedRefreshConfig)
}

func failoverGroupWithoutReplicationSchedule(id sdk.AccountObjectIdentifier, allowedAccountId sdk.AccountIdentifier, databaseId sdk.AccountObjectIdentifier, rpoAssurance bool) string {
	rpoAssuranceConfig := ""
	if rpoAssurance {
		rpoAssuranceConfig = "\n\trpo_assurance = true"
	}
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts = ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]%s
}
`, id.Name(), allowedAccountId.Name(), databaseId.Name(), rpoAssuranceConfig)
}

func failoverGroupWithNoWarehouse(id sdk.AccountObjectIdentifier, allowedAccountId sdk.AccountIdentifier, interval int) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule = "%d MINUTE"
}
`, id.Name(), allowedAccountId.Name(), interval)
}

func failoverGroupWithCronExpression(id sdk.AccountObjectIdentifier, allowedAccountId sdk.AccountIdentifier, expression string, databaseId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts = ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule = "USING CRON %s UTC"
}
`, id.Name(), allowedAccountId.Name(), databaseId.Name(), expression)
}

func failoverGroupWithAccountParameters(id sdk.AccountObjectIdentifier, allowedAccountId sdk.AccountIdentifier, databaseId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["ACCOUNT PARAMETERS", "WAREHOUSES", "DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts = ["%s"]
	allowed_databases = ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule = "USING CRON 0 0 10-20 * TUE,THU UTC"
}
`, id.Name(), allowedAccountId.Name(), databaseId.Name())
}

func failoverGroupWithChanges(id sdk.AccountObjectIdentifier, allowedAccountId sdk.AccountIdentifier, interval int) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name = "%[1]s"
	object_types = ["DATABASES", "INTEGRATIONS"]
	allowed_accounts = ["%[2]s"]
	allowed_integration_types = ["NOTIFICATION INTEGRATIONS"]
	replication_schedule = "%d MINUTE"
}
`, id.Name(), allowedAccountId.Name(), interval)
}

func TestAcc_FailoverGroup_UpdateAllowedAccounts(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	// We need to use the secondary account to test the failover group. Otherwise, we will get an error:
	// Error: error removing allowed accounts for failover group TEST err = 003909 (55000): Disabling replication of the replication group to the account in which the primary currently resides is not allowed.
	accountIdentifier := testClient().Context.CurrentAccountIdentifier(t)

	secondaryAccountIdentifier := secondaryTestClient().Account.GetAccountIdentifier(t)
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
			// Upgrade to the current provider version. The source account is no longer a part of allowed_accounts, so the config lists only the target account.
			// The secondary account was removed externally above, so this step adds it back with the current provider - proves the fix.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   failoverGroupWithAccounts(id.Name(), TestDatabaseName, secondaryAccountIdentifier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", secondaryAccountIdentifier.Name()),
				),
			},
			// Removing an allowed account needs a second target account, because allowed_accounts requires at least one.
			// TestAcc_FailoverGroup_AllowedAccounts_ThreeAccounts covers it where a third account is available.
		},
	})
}

// TestAcc_FailoverGroup_AllowedAccounts_ThreeAccounts removes a target account from allowed_accounts and adds it back.
func TestAcc_FailoverGroup_AllowedAccounts_ThreeAccounts(t *testing.T) {
	// The prod environment has no setup with three accounts on the same deployment, which is needed to have two target accounts.
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeProdEnvironment {
		t.Skip("Skipping for the prod Snowflake environment: three accounts on the same deployment are required")
	}

	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)
	thirdAccountId := thirdTestClient().Account.GetAccountIdentifier(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithAccounts(id.Name(), TestDatabaseName, secondaryAccountId, thirdAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", secondaryAccountId.Name()),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", thirdAccountId.Name()),
				),
			},
			// Remove an allowed account
			{
				Config: failoverGroupWithAccounts(id.Name(), TestDatabaseName, secondaryAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", secondaryAccountId.Name()),
				),
			},
			// Add it back
			{
				Config: failoverGroupWithAccounts(id.Name(), TestDatabaseName, secondaryAccountId, thirdAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "2"),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", secondaryAccountId.Name()),
					resource.TestCheckTypeSetElemAttr("snowflake_failover_group.fg", "allowed_accounts.*", thirdAccountId.Name()),
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
			// An object type outside the supported set is rejected during validation.
			{
				Config:      failoverGroupWithObjectTypes(id.Name(), `"DATABASES", "INVALID OBJECT TYPE;"`),
				ExpectError: regexp.MustCompile("invalid failover group object type: INVALID OBJECT TYPE;"),
			},
		},
	})
}

func failoverGroupWithObjectTypes(name, objectTypes string) string {
	return fmt.Sprintf(`
resource "snowflake_failover_group" "fg" {
	name         = "%[1]s"
	object_types = [%[2]s]
	allowed_accounts = ["MYORG.MYACCOUNT"]
}
`, name, objectTypes)
}

// The state written by the last release with the hand-written resource (nested replication_schedule, boolean ignore_edition_check,
// the source account in allowed_accounts) is upgraded without a diff once the configuration is adjusted.
func TestAcc_FailoverGroup_UpgradeFromV2_23_0(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	sourceAccount := testClient().Context.CurrentAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.FailoverGroupResource))

	// The nested replication_schedule block cannot be written with the current model builders.
	configInV2_23_0 := fmt.Sprintf(`
resource "snowflake_failover_group" "t" {
	name                 = "%s"
	object_types         = ["DATABASES"]
	allowed_accounts     = ["%s", "%s"]
	allowed_databases    = ["%s"]
	ignore_edition_check = true
	optimized_refresh    = true
	replication_schedule {
		interval = 10
	}
}
`, id.Name(), sourceAccount.Name(), targetAccount.Name(), database.ID().Name())
	current := model.FailoverGroup("t", id.Name(), []string{targetAccount.Name()}, []string{"DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()}).
		WithIgnoreEditionCheck(true).
		WithOptimizedRefresh(true).
		WithReplicationSchedule("10 MINUTE")
	ref := current.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// Create with the last release
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.23.0"),
				Config:            config.FromModels(t, providerModel) + configInV2_23_0,
				Check: assertThat(t,
					resourceassert.FailoverGroupResource(t, ref).
						HasNameString(id.Name()).
						HasAllowedAccounts(sourceAccount.Name(), targetAccount.Name()),
					objectassert.FailoverGroup(t, id).
						HasName(id.Name()).
						HasReplicationSchedule("10 MINUTE").
						HasIsOptimizedRefreshEnabled(true),
				),
			},
			// Upgrade: the state is converted, the configuration is rewritten (string schedule, no source account), and nothing is left to apply
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, current),
				Check: assertThat(t, []assert.TestCheckFuncProvider{
					resourceassert.FailoverGroupResource(t, ref).
						HasNameString(id.Name()).
						HasObjectTypes("DATABASES").
						HasAllowedAccounts(targetAccount.Name()).
						HasAllowedDatabases(database.ID().Name()).
						HasReplicationSchedule("10 MINUTE").
						HasOptimizedRefresh(true).
						HasIgnoreEditionCheck(true),
					objectassert.FailoverGroup(t, id).
						HasName(id.Name()).
						HasReplicationSchedule("10 MINUTE").
						HasIsOptimizedRefreshEnabled(true),
				}...),
			},
		},
	})
}

// The error integration is read from SHOW, so changing it or unsetting it outside of Terraform is drift that the next apply restores.
func TestAcc_FailoverGroup_ErrorIntegration_ExternalChanges(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	allowedAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	errorIntegration, cleanupErrorIntegration := testClient().NotificationIntegration.Create(t)
	t.Cleanup(cleanupErrorIntegration)
	otherIntegration, cleanupOtherIntegration := testClient().NotificationIntegration.Create(t)
	t.Cleanup(cleanupOtherIntegration)

	configured := model.FailoverGroup("t", id.Name(), []string{allowedAccount.FullyQualifiedName()}, []string{"DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()}).
		WithErrorIntegration(errorIntegration.ID().Name())
	ref := configured.ResourceReference()

	restored := []assert.TestCheckFuncProvider{
		resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasErrorIntegrationString(errorIntegration.ID().Name()),
		objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasErrorIntegration(errorIntegration.ID()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// 1. Create with the error integration
			{
				Config: config.FromModels(t, configured),
				Check:  assertThat(t, restored...),
			},
			// 2. The integration is changed outside of Terraform
			{
				PreConfig: func() {
					testClient().FailoverGroup.Alter(t, sdk.NewAlterFailoverGroupRequest(id).
						WithSet(*sdk.NewFailoverGroupSetRequest().WithErrorIntegration(otherIntegration.ID())))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "error_integration", sdk.String(errorIntegration.ID().Name()), sdk.String(otherIntegration.ID().Name())),
						planchecks.ExpectChange(ref, "error_integration", tfjson.ActionUpdate, sdk.String(otherIntegration.ID().Name()), sdk.String(errorIntegration.ID().Name())),
					},
				},
				Config: config.FromModels(t, configured),
				Check:  assertThat(t, restored...),
			},
			// 3. The integration is unset outside of Terraform
			{
				PreConfig: func() {
					testClient().FailoverGroup.Alter(t, sdk.NewAlterFailoverGroupRequest(id).
						WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithErrorIntegration(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "error_integration", sdk.String(errorIntegration.ID().Name()), sdk.String("")),
						planchecks.ExpectChange(ref, "error_integration", tfjson.ActionUpdate, sdk.String(""), sdk.String(errorIntegration.ID().Name())),
					},
				},
				Config: config.FromModels(t, configured),
				Check:  assertThat(t, restored...),
			},
		},
	})
}
