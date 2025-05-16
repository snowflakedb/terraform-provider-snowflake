//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func connectionsData() string {
	return `
    data "snowflake_connections" "test" {
        depends_on = [snowflake_primary_connection.test]
    }`
}

func TestAcc_Connections_Minimal(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	accountId := acc.TestClient().Account.GetAccountIdentifier(t)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	connectionModel := model.PrimaryConnection("test", id.Name())

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	dataConnections := accconfig.FromModels(t, connectionModel) + connectionsData()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.PrimaryConnection),
		Steps: []resource.TestStep{
			{
				Config: dataConnections,
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "1")),
					resourceshowoutputassert.ConnectionShowOutput(t, "snowflake_primary_connection.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
						HasAccountLocator(acc.TestClient().GetAccountLocator()).
						HasAccountName(accountId.AccountName()).
						HasOrganizationName(accountId.OrganizationName()).
						HasComment("").
						HasIsPrimary(true).
						HasPrimaryIdentifier(primaryConnectionAsExternalId).
						HasFailoverAllowedToAccounts(accountId).
						HasConnectionUrl(
							acc.TestClient().Connection.GetConnectionUrl(accountId.OrganizationName(), id.Name()),
						),
				),
			},
		},
	})
}

func TestAcc_Connections_Complete(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	accountId := acc.TestClient().Account.GetAccountIdentifier(t)
	secondaryAccountId := acc.SecondaryTestClient().Account.GetAccountIdentifier(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	connectionModel := model.PrimaryConnection("test", id.Name()).
		WithEnableFailover(secondaryAccountId).
		WithComment("test comment")

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	dataConnections := accconfig.FromModels(t, connectionModel) + connectionsData()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.PrimaryConnection),
		Steps: []resource.TestStep{
			{
				Config: dataConnections,
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "1")),
					resourceshowoutputassert.ConnectionShowOutput(t, "snowflake_connection.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
						HasAccountLocator(acc.TestClient().GetAccountLocator()).
						HasAccountName(accountId.AccountName()).
						HasOrganizationName(accountId.OrganizationName()).
						HasComment("test comment").
						HasIsPrimary(true).
						HasPrimaryIdentifier(primaryConnectionAsExternalId).
						HasFailoverAllowedToAccounts(accountId, secondaryAccountId).
						HasConnectionUrl(
							acc.TestClient().Connection.GetConnectionUrl(accountId.OrganizationName(), id.Name()),
						),
				),
			},
		},
	})
}

func TestAcc_Connections_Filtering(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	// TODO: [SNOW-1788041] - need to uppercase as connection name in snowflake is returned in uppercase
	prefix := random.AlphaN(4)
	prefix = strings.ToUpper(prefix)

	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	connectionModelOne := model.PrimaryConnection("c1", idOne.Name())
	connectionModelTwo := model.PrimaryConnection("c2", idTwo.Name())
	connectionModelThree := model.PrimaryConnection("c3", idThree.Name())

	configWithLike := accconfig.FromModels(t, connectionModelOne) +
		accconfig.FromModels(t, connectionModelTwo) +
		accconfig.FromModels(t, connectionModelThree)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.PrimaryConnection),
		Steps: []resource.TestStep{
			// with like
			{
				Config: configWithLike + connectionDatasourceWithLike(prefix+"%"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Connections_FilteringWithReplica(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	// TODO: [SNOW-1788041] - need to uppercase as connection name in snowflake is returned in uppercase
	prefix := random.AlphaN(4)
	prefix = strings.ToUpper(prefix)

	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.SecondaryTestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	accountId := acc.TestClient().Account.GetAccountIdentifier(t)

	_, cleanup := acc.SecondaryTestClient().Connection.CreateWithIdentifier(t, idTwo)
	t.Cleanup(cleanup)

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, idTwo)
	acc.SecondaryTestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(idTwo).
		WithEnableConnectionFailover(*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{accountId})))

	connectionModelOne := model.PrimaryConnection("c1", idOne.Name())
	connectionModelTwo := model.SecondaryConnection("c2", idTwo.Name(), primaryConnectionAsExternalId.FullyQualifiedName())

	configWithLike := accconfig.FromModels(t, connectionModelOne) +
		accconfig.FromModels(t, connectionModelTwo)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.ComposeCheckDestroy(t, resources.PrimaryConnection, resources.SecondaryConnection),
		Steps: []resource.TestStep{
			// with like
			{
				Config: configWithLike + connectionAndSecondaryConnectionDatasourceWithLike(prefix+"%"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "2"),
				),
			},
		},
	})
}

func connectionDatasourceWithLike(like string) string {
	return fmt.Sprintf(`
    data "snowflake_connections" "test" {
        depends_on = [snowflake_primary_connection.c1, snowflake_primary_connection.c2, snowflake_primary_connection.c3]

        like = "%s"
    }
`, like)
}

func connectionAndSecondaryConnectionDatasourceWithLike(like string) string {
	return fmt.Sprintf(`
    data "snowflake_connections" "test" {
        depends_on = [snowflake_primary_connection.c1, snowflake_secondary_connection.c2]

        like = "%s"
    }
`, like)
}

func TestAcc_Connections_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      connectionNonExisting(),
				ExpectError: regexp.MustCompile("there should be at least one connection"),
			},
		},
	})
}

func connectionNonExisting() string {
	return `
data "snowflake_connections" "test" {
  like = "non-existing-connection"

  lifecycle {
    postcondition {
      condition     = length(self.connections) > 0
      error_message = "there should be at least one connection"
    }
  }
}
`
}
