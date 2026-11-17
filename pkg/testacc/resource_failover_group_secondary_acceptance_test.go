//go:build account_level_tests

// The primary groups use the shared test databases of the secondary account, and a failover group allows a database in only one
// failover or replication group per account. Failover and replication groups share these restrictions, so the tests of both
// are account-level tests (they cannot run in parallel with each other, other PRs, or the integration/acceptance suites).

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FailoverGroupSecondary_BasicUseCase(t *testing.T) {
	// The replica is created in this account. The primary groups live in the secondary test account.
	// name and primary are both ForceNew, and the resource has no other attribute and no update.
	targetAccount := testClient().Account.GetAccountIdentifier(t)
	sourceAccount := secondaryTestClient().Account.GetAccountIdentifier(t)

	databaseA, cleanupDatabaseA := secondaryTestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabaseA)
	databaseB, cleanupDatabaseB := secondaryTestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabaseB)

	primaryA := secondaryTestClient().Ids.RandomAccountObjectIdentifier()
	primaryB := secondaryTestClient().Ids.RandomAccountObjectIdentifier()
	replicaID := testClient().Ids.RandomAccountObjectIdentifier()
	_, cleanupPrimaryA := secondaryTestClient().FailoverGroup.CreateWithRequest(t, sdk.NewCreateFailoverGroupRequest(
		primaryA,
		[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{targetAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{databaseA.ID()}))
	t.Cleanup(cleanupPrimaryA)
	_, cleanupPrimaryB := secondaryTestClient().FailoverGroup.CreateWithRequest(t, sdk.NewCreateFailoverGroupRequest(
		primaryB,
		[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{targetAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{databaseB.ID()}))
	t.Cleanup(cleanupPrimaryB)
	// A replica must be dropped before its primary.
	t.Cleanup(func() {
		testClient().FailoverGroup.Drop(t, replicaID)
	})

	primaryOfA := sdk.NewExternalObjectIdentifier(sourceAccount, primaryA)
	primaryOfB := sdk.NewExternalObjectIdentifier(sourceAccount, primaryB)
	testClient().FailoverGroup.WaitForPrimaryToBeVisible(t, primaryOfA)
	testClient().FailoverGroup.WaitForPrimaryToBeVisible(t, primaryOfB)
	basic := model.FailoverGroupSecondary("t", replicaID.Name(), primaryOfA.FullyQualifiedName())
	otherPrimary := model.FailoverGroupSecondary("t", replicaID.Name(), primaryOfB.FullyQualifiedName())
	ref := basic.ResourceReference()

	assertA := failoverGroupSecondaryAssertions(t, ref, replicaID, primaryOfA, targetAccount)
	assertB := failoverGroupSecondaryAssertions(t, ref, replicaID, primaryOfB, targetAccount)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroupSecondary),
		Steps: []resource.TestStep{
			// 1. Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertA...),
			},
			// 2. Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  failoverGroupSecondaryImportChecks(t, replicaID, primaryOfA),
			},
			// 3. Update - another primary (ForceNew)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, otherPrimary),
				Check:  assertThat(t, assertB...),
			},
			// 4. Import - another primary
			{
				Config:            config.FromModels(t, otherPrimary),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  failoverGroupSecondaryImportChecks(t, replicaID, primaryOfB),
			},
			// 5. Update - external changes: the replica is dropped and created again
			{
				PreConfig: func() {
					testClient().FailoverGroup.Drop(t, replicaID)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, otherPrimary),
				Check:  assertThat(t, assertB...),
			},
		},
	})
}

func failoverGroupSecondaryAssertions(
	t *testing.T,
	ref string,
	id sdk.AccountObjectIdentifier,
	primary sdk.ExternalObjectIdentifier,
	targetAccount sdk.AccountIdentifier,
) []assert.TestCheckFuncProvider {
	t.Helper()

	return []assert.TestCheckFuncProvider{
		resourceassert.FailoverGroupSecondaryResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasPrimaryString(primary.FullyQualifiedName()).
			HasIsPrimary(false),
		resourceshowoutputassert.FailoverGroupSecondaryShowOutput(t, ref).
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(false).
			HasPrimary(primary).
			HasComment("").
			HasReplicationSchedule("").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasAccountName(targetAccount.AccountName()).
			HasAccountLocator(testClient().GetAccountLocator()).
			HasOrganizationName(targetAccount.OrganizationName()).
			HasIsListingAutoFulfillmentGroup(false).
			HasNoObjectTypes().
			HasNoAllowedIntegrationTypes().
			HasNoAllowedAccounts(),
		objectassert.FailoverGroup(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(false).
			HasPrimary(primary).
			HasNoComment().
			HasNoObjectTypes().
			HasNoAllowedIntegrationTypes().
			HasNoAllowedAccounts().
			HasNoReplicationSchedule().
			HasNoSecondaryState().
			HasIsListingAutoFulfillmentGroup(false).
			HasIsOptimizedRefreshEnabled(false).
			HasRpoAssurance(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasAccountName(targetAccount.AccountName()).
			HasAccountLocator(testClient().GetAccountLocator()).
			HasOrganizationName(targetAccount.OrganizationName()).
			HasSnowflakeRegionNotEmpty(),
	}
}

func failoverGroupSecondaryImportChecks(t *testing.T, id sdk.AccountObjectIdentifier, primary sdk.ExternalObjectIdentifier) resource.ImportStateCheckFunc {
	t.Helper()
	encoded := helpers.EncodeResourceIdentifier(id)
	return assertThatImport(t,
		resourceassert.ImportedFailoverGroupSecondaryResource(t, encoded).
			HasNameString(id.Name()).
			HasPrimaryString(primary.FullyQualifiedName()).
			HasIsPrimary(false),
		resourceshowoutputassert.ImportedFailoverGroupSecondaryShowOutput(t, encoded).
			HasName(id.Name()).
			HasIsPrimary(false).
			HasPrimary(primary),
	)
}
