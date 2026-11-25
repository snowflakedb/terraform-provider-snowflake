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
	"regexp"
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
					assert.Check(resource.TestCheckResourceAttr(fromSource.DatasourceReference(), "id", "replication_groups_read")),
					assert.Check(replicationGroupsInDataSource.Listed(fromSource.DatasourceReference(), id.Name())),
					assert.Check(replicationGroupsInDataSource.Listed(fromSource.DatasourceReference(), replicaID.Name())),
					assert.Check(replicationGroupsInDataSource.HasShape(fromSource.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
					assert.Check(replicationGroupsInDataSource.HasShape(fromSource.DatasourceReference(), replicaID.Name(), false, secondaryTestClient().GetAccountLocator())),
					assert.Check(replicationGroupsInDataSource.Listed(inSource.DatasourceReference(), id.Name())),
					assert.Check(replicationGroupsInDataSource.Listed(inSource.DatasourceReference(), replicaID.Name())),
					assert.Check(replicationGroupsInDataSource.HasShape(inSource.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
					assert.Check(replicationGroupsInDataSource.HasShape(inSource.DatasourceReference(), replicaID.Name(), false, secondaryTestClient().GetAccountLocator())),
				),
			},
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   accconfig.FromModels(t, targetProvider, fromTarget),
				Check: assertThat(t,
					assert.Check(replicationGroupsInDataSource.Listed(fromTarget.DatasourceReference(), replicaID.Name())),
					assert.Check(replicationGroupsInDataSource.HasShape(fromTarget.DatasourceReference(), replicaID.Name(), false, secondaryTestClient().GetAccountLocator())),
					assert.Check(replicationGroupsInDataSource.Listed(fromTarget.DatasourceReference(), id.Name())),
					assert.Check(replicationGroupsInDataSource.HasShape(fromTarget.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
				),
			},
		},
	})
}

var replicationGroupsInDataSource = groupsInDataSource("replication_groups")

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
	return replicationGroupsInDataSource.AtIndex(ref, id.Name(), func(idx int) resource.TestCheckFunc {
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
			HasIsPrimaryReadOnly(false).
			HasAllowedAccounts(sourceAccount.Name(), allowedAccount.Name()).
			ToTerraformTestCheckFunc(t, testClient())
	})
}

func TestAcc_ReplicationGroups_invalidInAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.ReplicationGroups("test").WithInAccount("not an account identifier")),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
		},
	})
}
