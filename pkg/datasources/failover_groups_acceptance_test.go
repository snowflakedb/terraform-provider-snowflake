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

func TestAcc_FailoverGroups(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	failoverGroupId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: failoverGroupsConfig(failoverGroupId, accountName),
				Check: resource.ComposeTestCheckFunc(
					// TODO [SNOW-1348343]: fix these assertions - there might be multiple failover groups if we run tests in parallel
					resource.TestCheckResourceAttr("data.snowflake_failover_groups.d", "failover_groups.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_failover_groups.d", "failover_groups.0.object_types.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_failover_groups.d", "failover_groups.0.object_types.0", "ROLES"),
					resource.TestCheckResourceAttr("data.snowflake_failover_groups.d", "failover_groups.0.allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_failover_groups.d", "failover_groups.0.allowed_accounts.0", accountName),
				),
			},
		},
	})
}

func failoverGroupsConfig(failoverGroupId sdk.AccountObjectIdentifier, allowedAccount string) string {
	return fmt.Sprintf(`
	resource "snowflake_failover_group" "source_failover_group" {
		name                      = "%s"
		object_types              = ["ROLES"]
		allowed_accounts          = ["%s"]
	}

	data "snowflake_failover_groups" "d" {
		depends_on = [snowflake_failover_group.source_failover_group]
	}
	`, failoverGroupId.Name(), allowedAccount)
}
