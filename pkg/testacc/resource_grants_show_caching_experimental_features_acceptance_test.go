//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// countShowGrantsOnDatabaseQueries returns how many `SHOW GRANTS ON DATABASE <id>` statements were
// issued (across all sessions of the test user) for the given, test-unique database. Because the
// database name is randomly generated per test, this count is isolated to the queries produced by
// the test under inspection. Used to prove the GRANTS_SHOW_CACHING experiment collapses N identical
// SHOW calls into one, including across different resource types (see the mirror-key comment on
// showGrantsCached in pkg/resources/grant_helpers.go).
func countShowGrantsOnDatabaseQueries(t *testing.T, databaseId sdk.AccountObjectIdentifier) int {
	t.Helper()
	queryHistory := testClient().InformationSchema.GetQueryHistory(t, 1000)
	needle := fmt.Sprintf("SHOW GRANTS ON DATABASE %s", databaseId.FullyQualifiedName())
	return len(collections.Filter(queryHistory, func(h helpers.QueryHistory) bool {
		return strings.Contains(h.QueryText, needle)
	}))
}

// TestAcc_GrantsShowCaching_AccountRolePrivileges_SharedObjectIssuesSingleShow proves the cache
// effect for snowflake_grant_privileges_to_account_role: with the experiment enabled, two grants on
// the same database to two different roles result in exactly one `SHOW GRANTS ON DATABASE` call for
// that database over the provider's lifetime, instead of one per instance.
func TestAcc_GrantsShowCaching_AccountRolePrivileges_SharedObjectIssuesSingleShow(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: grantsShowCachingProviderFactory,
		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, grantToRoleA, grantToRoleB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.to_role_a", "account_role_name", roleA.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.to_role_b", "account_role_name", roleB.ID().Name()),
				),
			},
			// the second refresh must converge (both resources Read the same, cached SHOW GRANTS ON
			// DATABASE result) and that statement must have been issued exactly once across both steps.
			{
				Config: config.FromModels(t, experimentProviderModel, grantToRoleA, grantToRoleB),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					if got := countShowGrantsOnDatabaseQueries(t, database.ID()); got != 1 {
						return fmt.Errorf("expected exactly 1 `SHOW GRANTS ON DATABASE %s` with caching enabled, got %d", database.ID().FullyQualifiedName(), got)
					}
					return nil
				},
			},
		},
	})
}

// TestAcc_GrantsShowCaching_AccountRolePrivileges_SharedObjectIssuesShowPerInstanceWithoutExperiment
// is the contrast to the test above: without the experiment, the same database targeted by two
// separate grants is shown multiple times (once per Read of each instance), which is exactly the
// redundancy the experiment removes.
func TestAcc_GrantsShowCaching_AccountRolePrivileges_SharedObjectIssuesShowPerInstanceWithoutExperiment(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)
	roleA, roleACleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleACleanup)
	roleB, roleBCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleBCleanup)

	grantToRoleA := model.GrantPrivilegesToAccountRole("to_role_a", roleA.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database.ID())
	grantToRoleB := model.GrantPrivilegesToAccountRole("to_role_b", roleB.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, grantToRoleA, grantToRoleB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.to_role_a", "account_role_name", roleA.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.to_role_b", "account_role_name", roleB.ID().Name()),
				),
			},
			{
				Config: config.FromModels(t, grantToRoleA, grantToRoleB),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					// Without the experiment, each instance issues its own SHOW on every Read pass, so the
					// shared database is shown well more than once. We assert "more than one" rather than
					// an exact number to stay robust against Terraform's refresh cadence.
					if got := countShowGrantsOnDatabaseQueries(t, database.ID()); got <= 1 {
						return fmt.Errorf("expected more than 1 `SHOW GRANTS ON DATABASE %s` without the experiment, got %d", database.ID().FullyQualifiedName(), got)
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
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
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
			// Both resources' Read must converge on a second refresh, and — the point of this test —
			// the SHOW they share must have been issued exactly once across the two resource TYPES,
			// not once per type.
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
// result from before the Update.
func TestAcc_GrantsShowCaching_AccountRolePrivileges_UpdateInvalidatesCache(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
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
				Check:  resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
			},
			// If Update failed to invalidate the cache, this refresh would still see the pre-Update,
			// single-privilege SHOW GRANTS result and produce a non-empty plan (drift back to 1).
			{
				Config: config.FromModels(t, experimentProviderModel, grantWithUsageAndCreateSchema),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
