//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_Experimental_GrantAccountRole_SafeDestroy_MissingParentRole verifies that destroying
// a grant_account_role resource fails when the grantee (parent) role is deleted externally (default behavior),
// and succeeds when the GRANTS_SAFE_DESTROY experiment is enabled.
func TestAcc_Experimental_GrantAccountRole_SafeDestroy_MissingParentRole(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantAccountRole("test", role.ID().Name()).
		WithParentRoleName(parentRole.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantRole_SafeDestroy")

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
			// Drop the parent role externally and destroy WITHOUT experiment — must error.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Destroy:     true,
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
			// Recreate the parent role, drop it again in PreApply, then destroy with GRANTS_SAFE_DESTROY — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				PreConfig: func() {
					testClient().Role.CreateRoleWithIdentifier(t, parentRole.ID())
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Config:  config.FromModels(t, experimentProviderModel, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantAccountRole_SafeDestroy_MissingRole verifies that destroying
// a grant_account_role resource fails when the granted role itself is deleted externally (default behavior),
// and succeeds when the GRANTS_SAFE_DESTROY experiment is enabled.
func TestAcc_Experimental_GrantAccountRole_SafeDestroy_MissingRole(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantAccountRole("test", role.ID().Name()).
		WithParentRoleName(parentRole.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantRole_SafeDestroy")

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
			// Drop the granted role externally WITHOUT experiment — must error.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, role.ID())),
					},
				},
				Destroy:     true,
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
			// Recreate the granted role, drop it again in PreApply, then destroy with GRANTS_SAFE_DESTROY — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				PreConfig: func() {
					testClient().Role.CreateRoleWithIdentifier(t, role.ID())
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, role.ID())),
					},
				},
				Config:  config.FromModels(t, experimentProviderModel, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantDatabaseRole_SafeDestroy_MissingParentRole verifies that destroying
// a grant_database_role resource fails when the grantee (parent) account role is deleted externally
// (default behavior), and succeeds when the GRANTS_SAFE_DESTROY experiment is enabled.
func TestAcc_Experimental_GrantDatabaseRole_SafeDestroy_MissingParentRole(t *testing.T) {
	dbRole, dbRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantDatabaseRole("test", dbRole.ID().FullyQualifiedName()).
		WithParentRoleName(parentRole.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantRole_SafeDestroy")

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
			// Drop the parent account role externally WITHOUT experiment — must error.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Destroy:     true,
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
			// Recreate the parent account role, drop it again in PreApply, then destroy with GRANTS_SAFE_DESTROY — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				PreConfig: func() {
					testClient().Role.CreateRoleWithIdentifier(t, parentRole.ID())
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Config:  config.FromModels(t, experimentProviderModel, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantDatabaseRole_SafeDestroy_MissingDatabaseRole verifies that destroying
// a grant_database_role resource fails when the database role itself is deleted externally
// (default behavior), and succeeds when the GRANTS_SAFE_DESTROY experiment is enabled.
func TestAcc_Experimental_GrantDatabaseRole_SafeDestroy_MissingDatabaseRole(t *testing.T) {
	dbRole, dbRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantDatabaseRole("test", dbRole.ID().FullyQualifiedName()).
		WithParentRoleName(parentRole.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantRole_SafeDestroy")

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
			// Drop the database role externally WITHOUT experiment — must error.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().DatabaseRole.CleanupDatabaseRoleFunc(t, dbRole.ID())),
					},
				},
				Destroy:     true,
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
			// Recreate the database role, drop it again in PreApply, then destroy with GRANTS_SAFE_DESTROY — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				PreConfig: func() {
					testClient().DatabaseRole.CreateDatabaseRoleInDatabaseWithName(t, dbRole.ID().DatabaseId(), dbRole.ID().Name())
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().DatabaseRole.CleanupDatabaseRoleFunc(t, dbRole.ID())),
					},
				},
				Config:  config.FromModels(t, experimentProviderModel, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantApplicationRole_SafeDestroy_MissingParentAccountRole verifies that destroying
// a grant_application_role resource (ACCOUNT_ROLE grantee type) succeeds when the grantee account role is
// deleted externally and the GRANTS_SAFE_DESTROY experiment is enabled.
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

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantApplicationRole_SafeDestroy")

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
			// Drop the parent account role externally and destroy with GRANTS_SAFE_DESTROY — succeeds.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Destroy:     true,
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
			// Recreate the parent account role, drop it again in PreApply, then destroy with GRANTS_SAFE_DESTROY — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				PreConfig: func() {
					testClient().Role.CreateRoleWithIdentifier(t, parentRole.ID())
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Role.DropRoleFunc(t, parentRole.ID())),
					},
				},
				Config:  config.FromModels(t, experimentProviderModel, grantModel),
				Destroy: true,
			},
		},
	})
}

// TestAcc_Experimental_GrantApplicationRole_SafeDestroy_MissingApplication verifies that destroying
// a grant_application_role resource (APPLICATION grantee type) succeeds when the grantee application is
// deleted externally and the GRANTS_SAFE_DESTROY experiment is enabled.
//
// Note: unlike grant_privileges_to_* resources, grant_application_role Read already handles missing
// objects gracefully by clearing state. This test verifies the full lifecycle succeeds end-to-end
// via RevokeSafely handling the error in a race condition scenario.
func TestAcc_Experimental_GrantApplicationRole_SafeDestroy_MissingApplication(t *testing.T) {
	app, appPackage := createAppReturnApplicationPackage(t)

	applicationRoleName := testvars.ApplicationRole1
	appRoleFullName := sdk.NewDatabaseObjectIdentifier(app.ID().Name(), applicationRoleName).FullyQualifiedName()

	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantApplicationRole("test", appRoleFullName).
		WithParentAccountRoleName(parentRole.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantApplicationRole_SafeDestroy")

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
			// Drop the grantee application externally WITHOUT experiment — must error.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Application.DropApplicationFunc(t, app.ID())),
					},
				},
				Destroy:     true,
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
			// Recreate the grantee application, drop it again in PreApply, then destroy with GRANTS_SAFE_DESTROY — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				PreConfig: func() {
					testClient().Application.CreateApplicationWithIdentifier(t, app.ID(), appPackage.ID(), "v1")
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(testClient().Application.DropApplicationFunc(t, app.ID())),
					},
				},
				Config:  config.FromModels(t, experimentProviderModel, grantModel),
				Destroy: true,
			},
		},
	})
}
