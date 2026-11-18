//go:build account_level_tests

// Account object types other than SHARES and DATABASES (e.g. ROLES, INTEGRATIONS, WAREHOUSES) can be a part of only one
// failover or replication group per account. Snowflake rejects the second group with 500070 "Account object types, other than
// SHARES and DATABASES, can only be within one failover group or replication group" (documented as "Objects other than databases,
// external volumes, and shares must be in the same failover group"). The restriction is on the type, so unique names do not help.
// The tests below use DATABASES, which can be a part of several groups.
// Failover and replication groups share these restrictions, so the tests of both are account-level tests (they cannot run in parallel
// with each other, other PRs, or the integration/acceptance suites).

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FailoverGroups_CompleteUseCase(t *testing.T) {
	// SHOW has no LIKE, so the created group is found by its name.
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	const schedule = "10 MINUTE"

	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)

	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)

	failoverGroupModel := model.FailoverGroup("t", id.Name(), []string{targetAccount.Name()}, []string{"DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()}).
		WithComment(comment).
		WithReplicationSchedule(schedule)
	datasourceModel := datasourcemodel.FailoverGroups("test").
		WithDependsOn(failoverGroupModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, failoverGroupModel, datasourceModel),
				Check: assertThat(t, assert.Check(failoverGroupsShowOutput(
					t,
					datasourceModel.DatasourceReference(),
					id,
					sourceAccount,
					targetAccount,
					comment,
					schedule,
				))),
			},
		},
	})
}

func TestAcc_FailoverGroups_BasicUseCase_DifferentFiltering(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	replicaID := secondaryTestClient().Ids.RandomAccountObjectIdentifier()
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	primaryOf := sdk.NewExternalObjectIdentifier(sourceAccount, id)

	_, cleanupPrimary := testClient().FailoverGroup.CreateWithRequest(t, sdk.NewCreateFailoverGroupRequest(
		id,
		[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{targetAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}))
	t.Cleanup(cleanupPrimary)
	secondaryTestClient().FailoverGroup.WaitForPrimaryToBeVisible(t, primaryOf)
	secondaryTestClient().FailoverGroup.CreateSecondary(t, replicaID, primaryOf)
	t.Cleanup(func() { secondaryTestClient().FailoverGroup.Drop(t, replicaID) })

	fromSource := datasourcemodel.FailoverGroups("source")
	inSource := datasourcemodel.FailoverGroups("in_source").
		WithInAccount(sourceAccount.FullyQualifiedName())
	inTarget := datasourcemodel.FailoverGroups("in_target").
		WithInAccount(targetAccount.FullyQualifiedName())
	fromTarget := datasourcemodel.FailoverGroups("target")
	targetProvider := providermodel.SnowflakeProvider().
		WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.FailoverGroupsDatasource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, fromSource, inSource),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(fromSource.DatasourceReference(), "id", "failover_groups_read")),
					assert.Check(failoverGroupsInDataSource.Listed(fromSource.DatasourceReference(), id.Name())),
					assert.Check(failoverGroupsInDataSource.HasShape(fromSource.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
					assert.Check(failoverGroupsInDataSource.Listed(inSource.DatasourceReference(), id.Name())),
					assert.Check(failoverGroupsInDataSource.HasShape(inSource.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
				),
			},
			{
				// IN ACCOUNT of another account is rejected by SHOW, so the filter has to reach the statement.
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, inTarget),
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   accconfig.FromModels(t, targetProvider, fromTarget),
				Check: assertThat(t,
					assert.Check(failoverGroupsInDataSource.Listed(fromTarget.DatasourceReference(), replicaID.Name())),
					assert.Check(failoverGroupsInDataSource.HasShape(fromTarget.DatasourceReference(), replicaID.Name(), false, secondaryTestClient().GetAccountLocator())),
					assert.Check(failoverGroupsInDataSource.Listed(fromTarget.DatasourceReference(), id.Name())),
					assert.Check(failoverGroupsInDataSource.HasShape(fromTarget.DatasourceReference(), id.Name(), true, testClient().GetAccountLocator())),
				),
			},
		},
	})
}

func TestAcc_FailoverGroups_invalidInAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.FailoverGroups("test").WithInAccount("not an account identifier")),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
		},
	})
}

var failoverGroupsInDataSource = groupsInDataSource("failover_groups")

func failoverGroupsShowOutput(
	t *testing.T,
	ref string,
	id sdk.AccountObjectIdentifier,
	sourceAccount sdk.AccountIdentifier,
	targetAccount sdk.AccountIdentifier,
	comment string,
	schedule string,
) resource.TestCheckFunc {
	t.Helper()
	return failoverGroupsInDataSource.AtIndex(ref, id.Name(), func(idx int) resource.TestCheckFunc {
		return resourceshowoutputassert.FailoverGroupsDatasourceShowOutputOnIdx(t, ref, idx).
			HasCreatedOnNotEmpty().
			HasSnowflakeRegionNotEmpty().
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(true).
			HasPrimary(sdk.NewExternalObjectIdentifier(sourceAccount, id)).
			HasComment(comment).
			HasObjectTypes("DATABASES").
			HasNoAllowedIntegrationTypes().
			HasReplicationSchedule(schedule).
			HasSecondaryStateEmpty().
			HasNextScheduledRefreshNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasAccountName(sourceAccount.AccountName()).
			HasAccountLocator(testClient().GetAccountLocator()).
			HasOrganizationName(sourceAccount.OrganizationName()).
			HasIsListingAutoFulfillmentGroup(false).
			HasIsPrimaryReadOnly(false).
			HasAllowedAccounts(sourceAccount.Name(), targetAccount.Name()).
			ToTerraformTestCheckFunc(t, testClient())
	})
}
