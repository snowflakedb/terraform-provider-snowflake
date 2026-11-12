//go:build account_level_tests

// Account object types other than SHARES and DATABASES (e.g. ROLES, INTEGRATIONS, WAREHOUSES) can be a part of only one
// failover or replication group per account. Snowflake rejects the second group with 500070 "Account object types, other than
// SHARES and DATABASES, can only be within one failover group or replication group" (documented as "Objects other than databases,
// external volumes, and shares must be in the same failover group"). The restriction is on the type, so unique names do not help.
// The tests below use such types (ROLES, INTEGRATIONS).
// Failover and replication groups share these restrictions, so the tests of both are account-level tests (they cannot run in parallel
// with each other, other PRs, or the integration/acceptance suites).

package testacc

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
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
	share, cleanupShare := testClient().Share.Create(t)
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
	replicaID := secondaryTestClient().Ids.RandomAccountObjectIdentifier()
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	primaryOf := sdk.NewExternalObjectIdentifier(sourceAccount, id)

	testClient().ReplicationGroup.CreateWithRequest(t, sdk.NewCreateReplicationGroupRequest(
		id,
		[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{targetAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}))
	t.Cleanup(func() { testClient().ReplicationGroup.Drop(t, id) })
	secondaryTestClient().ReplicationGroup.WaitForPrimaryToBeVisible(t, primaryOf)
	secondaryTestClient().ReplicationGroup.CreateSecondary(t, replicaID, primaryOf)
	t.Cleanup(func() { secondaryTestClient().ReplicationGroup.Drop(t, replicaID) })

	fromSource := datasourcemodel.ReplicationGroups("source")
	inSource := datasourcemodel.ReplicationGroups("in_source").
		WithInAccount(sourceAccount.FullyQualifiedName())
	fromTarget := datasourcemodel.ReplicationGroups("target")

	targetProvider := providermodel.SnowflakeProvider().
		WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.ReplicationGroupsDatasource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, fromSource, inSource),
				Check: assertThat(t,
					assert.Check(replicationGroupListed(t, fromSource.DatasourceReference(), id.Name(), true)),
					assert.Check(replicationGroupListed(t, fromSource.DatasourceReference(), replicaID.Name(), true)),
					assert.Check(replicationGroupHasShape(t, fromSource.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
					assert.Check(replicationGroupHasShape(t, fromSource.DatasourceReference(), replicaID.Name(), false, secondaryTestClient().GetAccountLocator())),
					assert.Check(replicationGroupListed(t, inSource.DatasourceReference(), id.Name(), true)),
					assert.Check(replicationGroupListed(t, inSource.DatasourceReference(), replicaID.Name(), true)),
					assert.Check(replicationGroupHasShape(t, inSource.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
					assert.Check(replicationGroupHasShape(t, inSource.DatasourceReference(), replicaID.Name(), false, secondaryTestClient().GetAccountLocator())),
				),
			},
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   accconfig.FromModels(t, targetProvider, fromTarget),
				Check: assertThat(t,
					assert.Check(replicationGroupListed(t, fromTarget.DatasourceReference(), replicaID.Name(), true)),
					assert.Check(replicationGroupHasShape(t, fromTarget.DatasourceReference(), replicaID.Name(), false, secondaryTestClient().GetAccountLocator())),
					assert.Check(replicationGroupListed(t, fromTarget.DatasourceReference(), id.Name(), true)),
					assert.Check(replicationGroupHasShape(t, fromTarget.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
				),
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
			return fmt.Errorf("replication group %s listed=%t, want %t\n%s", name, found, want, replicationGroupShowListing(rs))
		}
		return nil
	}
}

// IN ACCOUNT still returns every replication group in that account, so the created row is not at index 0.
// Other data source tests filter with LIKE until the row is at a known index, then use TestCheckResourceAttr.
func replicationGroupHasShape(t *testing.T, ref, name string, isPrimary bool, locator string) resource.TestCheckFunc {
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
		for i := range count {
			prefix := fmt.Sprintf("replication_groups.%d.show_output.0.", i)
			if rs.Primary.Attributes[prefix+"name"] != name {
				continue
			}
			gotPrimary := rs.Primary.Attributes[prefix+"is_primary"]
			gotLocator := rs.Primary.Attributes[prefix+"account_locator"]
			if gotPrimary != strconv.FormatBool(isPrimary) || gotLocator != locator {
				return fmt.Errorf("replication group %s is_primary=%s locator=%s, want is_primary=%t locator=%s\n%s",
					name, gotPrimary, gotLocator, isPrimary, locator, replicationGroupShowListing(rs))
			}
			return nil
		}
		return fmt.Errorf("replication group %s not found in %s\n%s", name, ref, replicationGroupShowListing(rs))
	}
}

// IN ACCOUNT still returns every replication group in that account, so the created row is not at index 0.
// Other data source tests filter with LIKE until the row is at a known index, then use TestCheckResourceAttr.
func replicationGroupShowListing(rs *terraform.ResourceState) string {
	count, err := strconv.Atoi(rs.Primary.Attributes["replication_groups.#"])
	if err != nil {
		return rs.Primary.ID
	}
	var b strings.Builder
	for i := range count {
		prefix := fmt.Sprintf("replication_groups.%d.show_output.0.", i)
		fmt.Fprintf(&b, "  name=%s is_primary=%s locator=%s account=%s.%s\n",
			rs.Primary.Attributes[prefix+"name"],
			rs.Primary.Attributes[prefix+"is_primary"],
			rs.Primary.Attributes[prefix+"account_locator"],
			rs.Primary.Attributes[prefix+"organization_name"],
			rs.Primary.Attributes[prefix+"account_name"],
		)
	}
	return b.String()
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
	// Region group is reported only on the production deployment.
	expectedRegionGroup := ""
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeProdEnvironment {
		expectedRegionGroup = "PUBLIC"
	}
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
			HasRegionGroup(expectedRegionGroup).
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
