//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"strings"
	"testing"
	"time"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Connections_CompleteUseCase(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	accountId := testClient().Account.GetAccountIdentifier(t)
	id := testClient().Ids.RandomAccountObjectIdentifier()
	connectionModel := model.PrimaryConnection("test", id.Name())

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	connectionsModel := datasourcemodel.Connections("test").
		WithDependsOn(connectionModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PrimaryConnection),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, connectionModel, connectionsModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "1")),
					resourceshowoutputassert.ConnectionShowOutput(t, "snowflake_primary_connection.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasSnowflakeRegion(strings.TrimPrefix(testClient().Context.CurrentRegion(t), "PUBLIC.")).
						HasAccountLocator(testClient().GetAccountLocator()).
						HasAccountName(accountId.AccountName()).
						HasOrganizationName(accountId.OrganizationName()).
						HasComment("").
						HasIsPrimary(true).
						HasPrimary(primaryConnectionAsExternalId).
						HasFailoverAllowedToAccounts(accountId).
						HasConnectionUrl(
							testClient().Connection.GetConnectionUrl(accountId.OrganizationName(), id.Name()),
						),
				),
			},
		},
	})
}

func TestAcc_Connections_Complete(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	accountId := testClient().Account.GetAccountIdentifier(t)
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	connectionModel := model.PrimaryConnection("test", id.Name()).
		WithEnableFailover(secondaryAccountId).
		WithComment("test comment")

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	connectionsModel := datasourcemodel.Connections("test").
		WithDependsOn(connectionModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PrimaryConnection),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, connectionModel, connectionsModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "1")),
					resourceshowoutputassert.ConnectionShowOutput(t, "snowflake_primary_connection.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasSnowflakeRegion(strings.TrimPrefix(testClient().Context.CurrentRegion(t), "PUBLIC.")).
						HasAccountLocator(testClient().GetAccountLocator()).
						HasAccountName(accountId.AccountName()).
						HasOrganizationName(accountId.OrganizationName()).
						HasComment("test comment").
						HasIsPrimary(true).
						HasPrimary(primaryConnectionAsExternalId).
						HasFailoverAllowedToAccounts(accountId, secondaryAccountId).
						HasConnectionUrl(
							testClient().Connection.GetConnectionUrl(accountId.OrganizationName(), id.Name()),
						),
				),
			},
		},
	})
}

func TestAcc_Connections_BasicUseCase_DifferentFiltering(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	// TODO: [SNOW-1788041] - need to uppercase as connection name in snowflake is returned in uppercase
	prefix := random.AlphaN(4)
	prefix = strings.ToUpper(prefix)

	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	connectionModelOne := model.PrimaryConnection("c1", idOne.Name())
	connectionModelTwo := model.PrimaryConnection("c2", idTwo.Name())
	connectionModelThree := model.PrimaryConnection("c3", idThree.Name())

	connectionsDatasourceModel := datasourcemodel.Connections("test").
		WithLike(prefix+"%").
		WithDependsOn(connectionModelOne.ResourceReference(), connectionModelTwo.ResourceReference(), connectionModelThree.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PrimaryConnection),
		Steps: []resource.TestStep{
			// with like
			{
				Config: accconfig.FromModels(t, connectionModelOne, connectionModelTwo, connectionModelThree, connectionsDatasourceModel),
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

	// TODO: [SNOW-1788041] - need to uppercase as connection name in snowflake is returned in uppercase
	prefix := random.AlphaN(4)
	prefix = strings.ToUpper(prefix)

	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := secondaryTestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	accountId := testClient().Account.GetAccountIdentifier(t)
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)

	_, cleanup := secondaryTestClient().Connection.CreateWithIdentifier(t, idTwo)
	t.Cleanup(cleanup)

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(secondaryAccountId, idTwo)
	secondaryTestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(idTwo).
		WithEnableConnectionFailover(*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{accountId})))

	require.Eventually(t, func() bool {
		err := testClient().Connection.CreateAsReplicaOf(t, idTwo, primaryConnectionAsExternalId)
		if err != nil {
			t.Logf("waiting for connection failover to propagate: %s", err)
			return false
		}
		require.NoError(t, testClient().Connection.Drop(t, idTwo))
		return true
	}, 3*time.Minute, 10*time.Second)

	connectionModelOne := model.PrimaryConnection("c1", idOne.Name())
	connectionModelTwo := model.SecondaryConnection("c2", idTwo.Name(), primaryConnectionAsExternalId.FullyQualifiedName())

	connectionsDatasourceModel := datasourcemodel.Connections("test").
		WithLike(prefix+"%").
		WithDependsOn(connectionModelOne.ResourceReference(), connectionModelTwo.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: ComposeCheckDestroy(t, resources.PrimaryConnection, resources.SecondaryConnection),
		Steps: []resource.TestStep{
			// with like
			{
				Config: accconfig.FromModels(t, connectionModelOne, connectionModelTwo, connectionsDatasourceModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "3"),
				),
			},
		},
	})
}

func TestAcc_Connections_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
