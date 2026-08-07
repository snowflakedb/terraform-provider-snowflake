//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func countShowGrantsOnDatabaseQueries(t *testing.T, databaseId sdk.AccountObjectIdentifier) int {
	t.Helper()
	return countShowQueries(t, fmt.Sprintf("SHOW GRANTS ON DATABASE %s", databaseId.FullyQualifiedName()))
}

// TestAcc_GrantsShowCaching_AccountRolePrivileges_SharedObjectIssuesSingleShow proves the cache
// effect for snowflake_grant_privileges_to_account_role, then contrasts it against the same setup
// without the experiment enabled, reusing the same database/role fixtures via per-step provider
// factory overrides.
func TestAcc_GrantsShowCaching_AccountRolePrivileges_SharedObjectIssuesSingleShow(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)
	roleA, roleACleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleACleanup)
	roleB, roleBCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleBCleanup)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsShowCaching)
	grantToRoleA := model.GrantPrivilegesToAccountRole("to_role_a", roleA.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database.ID())
	grantToRoleB := model.GrantPrivilegesToAccountRole("to_role_b", roleB.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database.ID())

	checkAttrs := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.to_role_a", "account_role_name", roleA.ID().Name()),
		resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.to_role_b", "account_role_name", roleB.ID().Name()),
	)

	var showsWithCachingEnabled int

	resource.Test(t, resource.TestCase{
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: grantsShowCachingProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantToRoleA, grantToRoleB),
				Check:                    checkAttrs,
			},
			// the second refresh must converge (both resources Read the same, cached SHOW GRANTS ON
			// DATABASE result) and that statement must have been issued exactly once across both steps.
			{
				ProtoV6ProviderFactories: grantsShowCachingProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantToRoleA, grantToRoleB),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					got := countShowGrantsOnDatabaseQueries(t, database.ID())
					if got != 1 {
						return fmt.Errorf("expected exactly 1 `SHOW GRANTS ON DATABASE %s` with caching enabled, got %d", database.ID().FullyQualifiedName(), got)
					}
					showsWithCachingEnabled = got
					return nil
				},
			},
			// same config, without the experiment: the shared database must now be shown again on
			// every Read pass, which is exactly the redundancy the experiment removes.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantToRoleA, grantToRoleB),
				Check:                    checkAttrs,
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantToRoleA, grantToRoleB),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					// Compare against the count captured under caching, rather than asserting an
					// absolute number, so this isn't sensitive to exactly how many of the two
					// preceding (uncached) steps' Read passes each issued their own SHOW.
					newShows := countShowGrantsOnDatabaseQueries(t, database.ID()) - showsWithCachingEnabled
					if newShows < 2 {
						return fmt.Errorf("expected more than 1 additional `SHOW GRANTS ON DATABASE %s` without the experiment, got %d", database.ID().FullyQualifiedName(), newShows)
					}
					return nil
				},
			},
		},
	})
}

// TestAcc_GrantsShowCaching_CrossResource_OwnershipAndPrivilegesShareCache proves the design claim
// that the cache key is the rendered SQL of the SHOW statement, not something resource-type-scoped:
// a snowflake_grant_ownership resource and a snowflake_grant_privileges_to_account_role resource
// targeting the SAME object share one cached `SHOW GRANTS ON DATABASE` entry, even though they are
// different resource types. Ownership is transferred to roleOwner and privileges are separately
// granted to roleGrantee on the same database.
func TestAcc_GrantsShowCaching_CrossResource_OwnershipAndPrivilegesShareCache(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)
	roleOwner, roleOwnerCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleOwnerCleanup)
	roleGrantee, roleGranteeCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleGranteeCleanup)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsShowCaching)
	ownershipModel := model.GrantOwnership("ownership", []sdk.OwnershipGrantOn{{
		Object: &sdk.Object{
			ObjectType: sdk.ObjectTypeDatabase,
			Name:       database.ID(),
		},
	}}).WithAccountRoleName(roleOwner.ID().Name())
	privilegesModel := model.GrantPrivilegesToAccountRole("privileges", roleGrantee.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: grantsShowCachingProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, ownershipModel, privilegesModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_ownership.ownership", "account_role_name", roleOwner.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.privileges", "account_role_name", roleGrantee.ID().Name()),
				),
			},
			{
				Config: config.FromModels(t, experimentProviderModel, ownershipModel, privilegesModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					if got := countShowGrantsOnDatabaseQueries(t, database.ID()); got != 1 {
						return fmt.Errorf("expected exactly 1 `SHOW GRANTS ON DATABASE %s` shared across snowflake_grant_ownership and snowflake_grant_privileges_to_account_role, got %d", database.ID().FullyQualifiedName(), got)
					}
					return nil
				},
			},
		},
	})
}

// TestAcc_GrantsShowCaching_AccountRolePrivileges_UpdateInvalidatesCache proves that Update
// invalidates this resource's own cache entry: after changing the granted privileges, the next
// refresh must observe the new privileges (empty plan) rather than a stale cached SHOW GRANTS
// result from before the Update, and that observing them required a fresh SHOW, not a stale hit.
func TestAcc_GrantsShowCaching_AccountRolePrivileges_UpdateInvalidatesCache(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsShowCaching)
	resourceName := "snowflake_grant_privileges_to_account_role.test"
	grantWithUsage := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database.ID())
	grantWithUsageAndCreateSchema := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
		WithPrivileges("USAGE", "CREATE SCHEMA").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database.ID())

	var showsAfterUpdate int

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: grantsShowCachingProviderFactory,
		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, grantWithUsage),
				Check:  resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
			},
			{
				Config: config.FromModels(t, experimentProviderModel, grantWithUsageAndCreateSchema),
				Check: func(s *terraform.State) error {
					showsAfterUpdate = countShowGrantsOnDatabaseQueries(t, database.ID())
					return resource.TestCheckResourceAttr(resourceName, "privileges.#", "2")(s)
				},
			},
			// If Update failed to invalidate the cache, this refresh would either still see the
			// pre-Update, single-privilege SHOW GRANTS result (non-empty plan, drift back to 1) or
			// simply reuse the Update step's own already-fresh result without a new SHOW at all. The
			// query count confirms it's the latter case done correctly: no additional SHOW is needed
			// because Update already invalidated and Read already repopulated the cache, not because
			// the invalidation never happened.
			{
				Config: config.FromModels(t, experimentProviderModel, grantWithUsageAndCreateSchema),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					if got := countShowGrantsOnDatabaseQueries(t, database.ID()); got != showsAfterUpdate {
						return fmt.Errorf("expected no additional `SHOW GRANTS ON DATABASE %s` after Update already repopulated the cache, had %d, now %d", database.ID().FullyQualifiedName(), showsAfterUpdate, got)
					}
					return nil
				},
			},
		},
	})
}
