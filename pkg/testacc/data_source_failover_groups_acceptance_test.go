//go:build account_level_tests

// Account object types other than SHARES and DATABASES (e.g. ROLES, INTEGRATIONS, WAREHOUSES) can be a part of only one
// failover or replication group per account. Snowflake rejects the second group with 500070 "Account object types, other than
// SHARES and DATABASES, can only be within one failover group or replication group" (documented as "Objects other than databases,
// external volumes, and shares must be in the same failover group"). The restriction is on the type, so unique names do not help.
// The tests below use such a type (ROLES).
// Failover and replication groups share these restrictions, so the tests of both are account-level tests (they cannot run in parallel
// with each other, other PRs, or the integration/acceptance suites).

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FailoverGroups(t *testing.T) {
	failoverGroupId := testClient().Ids.RandomAccountObjectIdentifier()
	currentAccountId := testClient().Account.GetAccountIdentifier(t)
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: failoverGroupsConfig(failoverGroupId, currentAccountId, secondaryAccountId),
				Check: resource.ComposeTestCheckFunc(
					// Other failover groups may exist (parallel runs), so the group is found by its primary identifier.
					resource.TestCheckTypeSetElemNestedAttrs("data.snowflake_failover_groups.d", "failover_groups.*", map[string]string{
						"primary":            sdk.NewExternalObjectIdentifier(currentAccountId, failoverGroupId).FullyQualifiedName(),
						"object_types.#":     "1",
						"object_types.0":     "ROLES",
						"allowed_accounts.#": "2",
					}),
				),
			},
		},
	})
}

func failoverGroupsConfig(failoverGroupId sdk.AccountObjectIdentifier, currentAccountId sdk.AccountIdentifier, allowedAccountId sdk.AccountIdentifier) string {
	return fmt.Sprintf(`
	resource "snowflake_failover_group" "source_failover_group" {
		name                      = "%s"
		object_types              = ["ROLES"]
		allowed_accounts          = ["%s", "%s"]
	}

	data "snowflake_failover_groups" "d" {
		depends_on = [snowflake_failover_group.source_failover_group]
	}
	`, failoverGroupId.Name(), currentAccountId.Name(), allowedAccountId.Name())
}
