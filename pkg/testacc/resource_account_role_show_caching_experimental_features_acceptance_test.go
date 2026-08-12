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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func countShowRolesLikeQueries(t *testing.T, roleId sdk.AccountObjectIdentifier) int {
	t.Helper()
	queryHistory := testClient().InformationSchema.GetQueryHistory(t, 1000)
	needle := fmt.Sprintf("SHOW ROLES LIKE '%s'", roleId.Name())
	return len(collections.Filter(queryHistory, func(h helpers.QueryHistory) bool {
		return strings.Contains(h.QueryText, needle)
	}))
}

// TestAcc_AccountRoleShowCaching_SharedRoleIssuesSingleShow proves the cache effect: with the
// experiment enabled, two snowflake_grant_privileges_to_account_role resources granting privileges
// on different objects to the same, pre-existing role result in exactly one `SHOW ROLES LIKE`
// existence check for that role over the provider's lifetime, instead of one per instance.
func TestAcc_AccountRoleShowCaching_SharedRoleIssuesSingleShow(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	database1, database1Cleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(database1Cleanup)
	database2, database2Cleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(database2Cleanup)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.AccountRoleShowCaching)
	grantOnDb1 := model.GrantPrivilegesToAccountRole("on_db1", role.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database1.ID())
	grantOnDb2 := model.GrantPrivilegesToAccountRole("on_db2", role.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database2.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accountRoleShowCachingProviderFactory,
		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, grantOnDb1, grantOnDb2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.on_db1", "account_role_name", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.on_db2", "account_role_name", role.ID().Name()),
				),
			},
			// the second refresh must converge (both instances Read the same, cached role lookup)
			// and that lookup must have been issued exactly once across both steps.
			{
				Config: config.FromModels(t, experimentProviderModel, grantOnDb1, grantOnDb2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					if got := countShowRolesLikeQueries(t, role.ID()); got != 1 {
						return fmt.Errorf("expected exactly 1 `SHOW ROLES LIKE '%s'` with caching enabled, got %d", role.ID().Name(), got)
					}
					return nil
				},
			},
		},
	})
}

// TestAcc_AccountRoleShowCaching_SharedRoleIssuesShowPerInstanceWithoutExperiment is the contrast
// to the test above: without the experiment, the same role referenced by two grant instances is
// shown multiple times (once per Read of each instance), which is exactly the redundancy the
// experiment removes.
func TestAcc_AccountRoleShowCaching_SharedRoleIssuesShowPerInstanceWithoutExperiment(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	database1, database1Cleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(database1Cleanup)
	database2, database2Cleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(database2Cleanup)

	grantOnDb1 := model.GrantPrivilegesToAccountRole("on_db1", role.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database1.ID())
	grantOnDb2 := model.GrantPrivilegesToAccountRole("on_db2", role.ID().Name()).
		WithPrivileges("USAGE").
		WithOnAccountObject(sdk.ObjectTypeDatabase, database2.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, grantOnDb1, grantOnDb2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.on_db1", "account_role_name", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.on_db2", "account_role_name", role.ID().Name()),
				),
			},
			{
				Config: config.FromModels(t, grantOnDb1, grantOnDb2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					// Without the experiment, each instance issues its own SHOW on every Read pass, so
					// the shared role is shown well more than once. We assert "more than one" rather
					// than an exact number to stay robust against Terraform's refresh cadence.
					if got := countShowRolesLikeQueries(t, role.ID()); got <= 1 {
						return fmt.Errorf("expected more than 1 `SHOW ROLES LIKE '%s'` without the experiment, got %d", role.ID().Name(), got)
					}
					return nil
				},
			},
		},
	})
}

// TestAcc_AccountRoleShowCaching_RenameInvalidatesCache proves that UpdateAccountRole invalidates
// its own cached entry on rename: after renaming, the next refresh must observe the role under its
// new identifier (empty plan) rather than a stale cached lookup keyed by the old identifier leaking
// into a later Read for a same-named role, or the new identifier missing a required cache
// invalidation and producing drift.
func TestAcc_AccountRoleShowCaching_RenameInvalidatesCache(t *testing.T) {
	oldId := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.AccountRoleShowCaching)
	resourceName := "snowflake_account_role.test"
	withOldName := model.AccountRole("test", oldId.Name())
	withNewName := model.AccountRole("test", newId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accountRoleShowCachingProviderFactory,
		CheckDestroy:             CheckDestroy(t, resources.AccountRole),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, withOldName),
				Check:  resource.TestCheckResourceAttr(resourceName, "name", oldId.Name()),
			},
			{
				Config: config.FromModels(t, experimentProviderModel, withNewName),
				Check:  resource.TestCheckResourceAttr(resourceName, "name", newId.Name()),
			},
			// If the rename failed to invalidate the cache, this refresh would either still see the
			// pre-rename cached lookup (drift/error) or miss a required re-populate; either way the
			// plan would not be empty.
			{
				Config: config.FromModels(t, experimentProviderModel, withNewName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
