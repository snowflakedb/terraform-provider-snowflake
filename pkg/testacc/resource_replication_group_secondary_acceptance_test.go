//go:build non_account_level_tests

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

func TestAcc_ReplicationGroupSecondary_BasicUseCase(t *testing.T) {
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
	secondaryTestClient().ReplicationGroup.CreateWithRequest(t, sdk.NewCreateReplicationGroupRequest(
		primaryA,
		[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{targetAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{databaseA.ID()}))
	secondaryTestClient().ReplicationGroup.CreateWithRequest(t, sdk.NewCreateReplicationGroupRequest(
		primaryB,
		[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{targetAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{databaseB.ID()}))
	t.Cleanup(func() {
		testClient().ReplicationGroup.Drop(t, replicaID)
		secondaryTestClient().ReplicationGroup.Drop(t, primaryA)
		secondaryTestClient().ReplicationGroup.Drop(t, primaryB)
	})

	primaryOfA := sdk.NewExternalObjectIdentifier(sourceAccount, primaryA)
	primaryOfB := sdk.NewExternalObjectIdentifier(sourceAccount, primaryB)
	testClient().ReplicationGroup.WaitForPrimaryToBeVisible(t, primaryOfA)
	testClient().ReplicationGroup.WaitForPrimaryToBeVisible(t, primaryOfB)
	basic := model.ReplicationGroupSecondary("t", replicaID.Name()).WithPrimary(primaryOfA.FullyQualifiedName())
	otherPrimary := model.ReplicationGroupSecondary("t", replicaID.Name()).WithPrimary(primaryOfB.FullyQualifiedName())
	ref := basic.ResourceReference()

	assertA := replicationGroupSecondaryAssertions(t, ref, replicaID, primaryOfA, targetAccount)
	assertB := replicationGroupSecondaryAssertions(t, ref, replicaID, primaryOfB, targetAccount)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ReplicationGroupSecondary),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertA...),
			},
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  replicationGroupSecondaryImportChecks(t, replicaID, primaryOfA),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, otherPrimary),
				Check:  assertThat(t, assertB...),
			},
			{
				Config:            config.FromModels(t, otherPrimary),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  replicationGroupSecondaryImportChecks(t, replicaID, primaryOfB),
			},
			{
				PreConfig: func() {
					testClient().ReplicationGroup.Drop(t, replicaID)
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

func TestAcc_ReplicationGroupSecondary_RecreateWhenExternallyPrimary(t *testing.T) {
	targetAccount := testClient().Account.GetAccountIdentifier(t)
	sourceAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	sourceDatabase, cleanupSourceDatabase := secondaryTestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupSourceDatabase)
	localDatabase, cleanupLocalDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupLocalDatabase)

	primaryID := secondaryTestClient().Ids.RandomAccountObjectIdentifier()
	replicaID := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryTestClient().ReplicationGroup.CreateWithRequest(t, sdk.NewCreateReplicationGroupRequest(
		primaryID,
		[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{targetAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{sourceDatabase.ID()}))
	t.Cleanup(func() {
		testClient().ReplicationGroup.Drop(t, replicaID)
		secondaryTestClient().ReplicationGroup.Drop(t, primaryID)
	})

	primaryOf := sdk.NewExternalObjectIdentifier(sourceAccount, primaryID)
	testClient().ReplicationGroup.WaitForPrimaryToBeVisible(t, primaryOf)
	replica := model.ReplicationGroupSecondary("t", replicaID.Name()).WithPrimary(primaryOf.FullyQualifiedName())
	ref := replica.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ReplicationGroupSecondary),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, replica),
				Check:  assertThat(t, replicationGroupSecondaryAssertions(t, ref, replicaID, primaryOf, targetAccount)...),
			},
			{
				// Replace the replica with a primary of the same name. The resource must destroy it and create the replica again.
				PreConfig: func() {
					testClient().ReplicationGroup.Drop(t, replicaID)
					testClient().ReplicationGroup.CreateWithRequest(t, sdk.NewCreateReplicationGroupRequest(
						replicaID,
						[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
						[]sdk.AccountIdentifier{sourceAccount},
					).WithAllowedDatabases([]sdk.AccountObjectIdentifier{localDatabase.ID()}))
				},
				Config: config.FromModels(t, replica),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t, replicationGroupSecondaryAssertions(t, ref, replicaID, primaryOf, targetAccount)...),
			},
		},
	})
}

func replicationGroupSecondaryAssertions(
	t *testing.T,
	ref string,
	id sdk.AccountObjectIdentifier,
	primary sdk.ExternalObjectIdentifier,
	targetAccount sdk.AccountIdentifier,
) []assert.TestCheckFuncProvider {
	t.Helper()

	return []assert.TestCheckFuncProvider{
		resourceassert.ReplicationGroupSecondaryResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasPrimaryString(primary.FullyQualifiedName()).
			HasIsPrimary(false),
		resourceshowoutputassert.ReplicationGroupSecondaryShowOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.ReplicationGroupTypeReplication).
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
		objectassert.ReplicationGroup(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasType(sdk.ReplicationGroupTypeReplication).
			HasIsPrimary(false).
			HasPrimary(primary).
			HasNoComment().
			HasNoObjectTypes().
			HasNoAllowedIntegrationTypes().
			HasNoAllowedAccounts().
			HasNoReplicationSchedule().
			HasNoSecondaryState().
			HasNoReadOnly().
			HasIsListingAutoFulfillmentGroup(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasAccountName(targetAccount.AccountName()).
			HasAccountLocator(testClient().GetAccountLocator()).
			HasOrganizationName(targetAccount.OrganizationName()).
			HasSnowflakeRegionNotEmpty(),
	}
}

func replicationGroupSecondaryImportChecks(t *testing.T, id sdk.AccountObjectIdentifier, primary sdk.ExternalObjectIdentifier) resource.ImportStateCheckFunc {
	t.Helper()
	encoded := helpers.EncodeResourceIdentifier(id)
	return assertThatImport(t,
		resourceassert.ImportedReplicationGroupSecondaryResource(t, encoded).
			HasNameString(id.Name()).
			HasPrimaryString(primary.FullyQualifiedName()).
			HasIsPrimary(false),
		resourceshowoutputassert.ImportedReplicationGroupSecondaryShowOutput(t, encoded).
			HasName(id.Name()).
			HasIsPrimary(false).
			HasPrimary(primary),
	)
}
