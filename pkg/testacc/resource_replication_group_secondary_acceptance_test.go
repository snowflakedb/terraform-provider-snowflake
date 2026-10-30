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
			HasPrimaryString(primary.FullyQualifiedName()),
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
			HasPrimaryString(primary.FullyQualifiedName()),
		resourceshowoutputassert.ImportedReplicationGroupSecondaryShowOutput(t, encoded).
			HasName(id.Name()).
			HasIsPrimary(false).
			HasPrimary(primary),
	)
}
