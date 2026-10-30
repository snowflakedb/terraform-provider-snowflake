//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ReplicationGroups_CompleteUseCase(t *testing.T) {
	// DESC REPLICATION GROUP is not published. SHOW has no LIKE, so this step finds the created group by name.
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	const schedule = "10 MINUTE"

	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	database2, cleanupDatabase2 := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase2)
	share, cleanupShare := testClient().Share.CreateShare(t)
	t.Cleanup(cleanupShare)
	externalVolume, cleanupExternalVolume := testClient().ExternalVolume.Create(t)
	t.Cleanup(cleanupExternalVolume)
	errorIntegration, cleanupErrorIntegration := testClient().NotificationIntegration.Create(t)
	t.Cleanup(cleanupErrorIntegration)

	allowedAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	// SHOW returns these in its own order. The resource stores that order.
	objectTypes := []string{
		string(sdk.ReplicationGroupObjectTypeDatabases),
		string(sdk.ReplicationGroupObjectTypeRoles),
		string(sdk.ReplicationGroupObjectTypeExternalVolumes),
		string(sdk.ReplicationGroupObjectTypeShares),
		string(sdk.ReplicationGroupObjectTypeIntegrations),
	}
	integrationTypes := []string{
		string(sdk.ReplicationGroupIntegrationTypeApiIntegrations),
		string(sdk.ReplicationGroupIntegrationTypeExternalAccessIntegrations),
		string(sdk.ReplicationGroupIntegrationTypeNotificationIntegrations),
		string(sdk.ReplicationGroupIntegrationTypeSecurityIntegrations),
		string(sdk.ReplicationGroupIntegrationTypeStorageIntegrations),
	}

	// CREATE has no COMMENT clause. The second config is what sends the ALTER.
	// WithComment mutates the model, so the two configs need separate builders.
	group := func(withComment bool) *model.ReplicationGroupModel {
		m := model.ReplicationGroup("t", id.Name(), []string{allowedAccount.FullyQualifiedName()}, objectTypes).
			WithAllowedDatabases([]string{database.ID().FullyQualifiedName(), database2.ID().FullyQualifiedName()}).
			WithAllowedExternalVolumes([]string{externalVolume.FullyQualifiedName()}).
			WithAllowedShares([]string{share.ID().FullyQualifiedName()}).
			WithAllowedIntegrationTypes(integrationTypes).
			WithReplicationSchedule(schedule).
			WithErrorIntegration(errorIntegration.ID().FullyQualifiedName())
		if withComment {
			m.WithComment(comment)
		}
		return m
	}
	withComment := group(true)
	datasourceModel := datasourcemodel.ReplicationGroups("test").
		WithDependsOn(withComment.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, group(false)),
			},
			{
				Config: accconfig.FromModels(t, withComment, datasourceModel),
				Check: assertThat(t, assert.Check(replicationGroupsShowOutput(
					t,
					datasourceModel.DatasourceReference(),
					id,
					sourceAccount,
					allowedAccount,
					comment,
					schedule,
				))),
			},
		},
	})
}

func TestAcc_ReplicationGroups_BasicUseCase_DifferentFiltering(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	otherAccount := secondaryTestClient().Account.GetAccountIdentifier(t)

	group := model.ReplicationGroup("t", id.Name(), []string{otherAccount.FullyQualifiedName()}, []string{string(sdk.ReplicationGroupObjectTypeDatabases)}).
		WithAllowedDatabases([]string{database.ID().FullyQualifiedName()})
	// IN ACCOUNT of another account is not authorized for this role, so the filter step uses this account.
	inSource := datasourcemodel.ReplicationGroups("test").
		WithInAccount(sourceAccount.FullyQualifiedName()).
		WithDependsOn(group.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, group, inSource),
				Check:  assertThat(t, assert.Check(replicationGroupListed(t, inSource.DatasourceReference(), id.Name(), true))),
			},
		},
	})
}

// IN ACCOUNT still returns every replication group in that account, so the created row is not at index 0.
// Other data source tests filter with LIKE until the row is at a known index, then use TestCheckResourceAttr.
func replicationGroupListed(t *testing.T, ref string, name string, want bool) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[ref]
		if !ok {
			return fmt.Errorf("data source %s not found", ref)
		}
		count, err := strconv.Atoi(rs.Primary.Attributes["replication_groups.#"])
		if err != nil {
			return err
		}
		found := false
		for i := range count {
			if rs.Primary.Attributes[fmt.Sprintf("replication_groups.%d.show_output.0.name", i)] == name {
				found = true
				break
			}
		}
		if found != want {
			return fmt.Errorf("replication group %s listed=%t, want %t", name, found, want)
		}
		return nil
	}
}

// The show-output builder needs an index. This finds the created group for the same reason as replicationGroupListed.
func replicationGroupsShowOutput(
	t *testing.T,
	ref string,
	id sdk.AccountObjectIdentifier,
	sourceAccount sdk.AccountIdentifier,
	allowedAccount sdk.AccountIdentifier,
	comment string,
	schedule string,
) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[ref]
		if !ok {
			return fmt.Errorf("data source %s not found", ref)
		}
		count, err := strconv.Atoi(rs.Primary.Attributes["replication_groups.#"])
		if err != nil {
			return err
		}
		idx := -1
		for i := range count {
			if rs.Primary.Attributes[fmt.Sprintf("replication_groups.%d.show_output.0.name", i)] == id.Name() {
				idx = i
				break
			}
		}
		if idx < 0 {
			return fmt.Errorf("replication group %s not found in %s", id.Name(), ref)
		}
		return resourceshowoutputassert.ReplicationGroupsDatasourceShowOutputOnIdx(t, ref, idx).
			HasCreatedOnNotEmpty().
			HasSnowflakeRegionNotEmpty().
			HasRegionGroup("").
			HasName(id.Name()).
			HasType(sdk.ReplicationGroupTypeReplication).
			HasIsPrimary(true).
			HasPrimary(sdk.NewExternalObjectIdentifier(sourceAccount, id)).
			HasComment(comment).
			HasReplicationSchedule(schedule).
			HasSecondaryStateEmpty().
			HasNextScheduledRefreshNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasAccountName(sourceAccount.AccountName()).
			HasAccountLocator(testClient().GetAccountLocator()).
			HasOrganizationName(sourceAccount.OrganizationName()).
			HasIsListingAutoFulfillmentGroup(false).
			HasReadOnly(false).
			HasAllowedAccounts(sourceAccount.Name(), allowedAccount.Name()).
			ToTerraformTestCheckFunc(t, testClient())(s)
	}
}
