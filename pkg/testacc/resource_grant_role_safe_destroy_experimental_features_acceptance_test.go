//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// NOTE: In most cases, we call DROP inside PreApply check.
// This is to force destroying the resource after the plan is computed.

// TestAcc_Experimental_GrantAccountRole_SafeDestroy_MissingParentRole verifies that destroying
// a grant_account_role resource succeeds when the grantee (parent) role is deleted externally.
func TestAcc_Experimental_GrantAccountRole_SafeDestroy_MissingParentRole(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantAccountRole("test", role.ID().Name()).
		WithParentRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the parent role in PreApply, then destroy — succeeds.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Config:  config.FromModels(t, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantAccountRole_SafeDestroy_MissingRole verifies that destroying
// a grant_account_role resource succeeds when the granted role itself is deleted externally.
func TestAcc_Experimental_GrantAccountRole_SafeDestroy_MissingRole(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantAccountRole("test", role.ID().Name()).
		WithParentRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the granted role in PreApply, then destroy — succeeds.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, role.ID())),
					},
				},
				Config:  config.FromModels(t, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantDatabaseRole_SafeDestroy_MissingParentRole verifies that destroying
// a grant_database_role resource succeeds when the grantee (parent) account role is deleted externally.
func TestAcc_Experimental_GrantDatabaseRole_SafeDestroy_MissingParentRole(t *testing.T) {
	dbRole, dbRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantDatabaseRole("test", dbRole.ID().FullyQualifiedName()).
		WithParentRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the parent account role in PreApply, then destroy — succeeds.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Config:  config.FromModels(t, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantDatabaseRole_SafeDestroy_MissingDatabaseRole verifies that destroying
// a grant_database_role resource succeeds when the database role itself is deleted externally.
func TestAcc_Experimental_GrantDatabaseRole_SafeDestroy_MissingDatabaseRole(t *testing.T) {
	dbRole, dbRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantDatabaseRole("test", dbRole.ID().FullyQualifiedName()).
		WithParentRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the database role in PreApply, then destroy — succeeds.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().DatabaseRole.CleanupDatabaseRoleFunc(t, dbRole.ID())),
					},
				},
				Config:  config.FromModels(t, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantApplicationRole_SafeDestroy_MissingParentAccountRole verifies that destroying
// a grant_application_role resource (ACCOUNT_ROLE grantee type) succeeds when the grantee account role is
// deleted externally.
//
// Note: unlike grant_privileges_to_* resources, grant_application_role Read already handles missing
// objects gracefully by clearing state. This test verifies the full lifecycle succeeds end-to-end
// via RevokeSafely handling the error in a race condition scenario.
func TestAcc_Experimental_GrantApplicationRole_SafeDestroy_MissingParentAccountRole(t *testing.T) {
	app := createApp(t)
	applicationRoleName := testvars.ApplicationRole1
	appRoleFullName := sdk.NewDatabaseObjectIdentifier(app.ID().Name(), applicationRoleName).FullyQualifiedName()

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantApplicationRole("test", appRoleFullName).
		WithParentAccountRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the parent account role in PreApply, then destroy — succeeds.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Config:  config.FromModels(t, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantApplicationRole_SafeDestroy_MissingApplication verifies that destroying
// a grant_application_role resource (APPLICATION grantee type) succeeds when the grantee application is
// deleted externally.
//
// Note: unlike grant_privileges_to_* resources, grant_application_role Read already handles missing
// objects gracefully by clearing state. This test verifies the full lifecycle succeeds end-to-end
// via RevokeSafely handling the error in a race condition scenario.
func TestAcc_Experimental_GrantApplicationRole_SafeDestroy_MissingApplication(t *testing.T) {
	app, _ := createAppReturnApplicationPackage(t)

	applicationRoleName := testvars.ApplicationRole1
	appRoleFullName := sdk.NewDatabaseObjectIdentifier(app.ID().Name(), applicationRoleName).FullyQualifiedName()

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantApplicationRole("test", appRoleFullName).
		WithParentAccountRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the grantee application in PreApply, then destroy — succeeds.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Application.DropApplicationFunc(t, app.ID())),
					},
				},
				Config:  config.FromModels(t, grantModel),
				Destroy: true,
			},
		},
	})
}
