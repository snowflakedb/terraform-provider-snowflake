//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_Experimental_GrantPrivilegesToAccountRole_SafeDestroy_MissingWarehouse verifies that destroying
// a grant resource fails when the target warehouse is deleted externally (default behavior), and succeeds
// when the GRANTS_SAFE_DESTROY experiment is enabled.
// Uses all_privileges = true so that Read skips existence checks and Delete is actually called.
func TestAcc_Experimental_GrantPrivilegesToAccountRole_SafeDestroy_MissingWarehouse(t *testing.T) {
	wh, whCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(whCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	grantModel := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
		WithAllPrivileges(true).
		WithOnAccountObject(sdk.ObjectTypeWarehouse, wh.ID())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantPrivileges_SafeDestroy")

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
			// Drop the warehouse externally, then try to destroy without experiment — expect failure.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().Warehouse.DropWarehouseFunc(t, wh.ID()),
				Config:                   config.FromModels(t, grantModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with GRANTS_SAFE_DESTROY experiment — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				Destroy:                  true,
			},
		},
	})
}

// TestAcc_Experimental_GrantPrivilegesToAccountRole_SafeDestroy_MissingRole verifies that destroying
// a grant resource fails when the grantee role is deleted externally (default behavior), and succeeds
// when the GRANTS_SAFE_DESTROY experiment is enabled.
// Uses all_privileges = true so that Read skips existence checks and Delete is actually called.
func TestAcc_Experimental_GrantPrivilegesToAccountRole_SafeDestroy_MissingRole(t *testing.T) {
	wh, whCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(whCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	grantModel := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
		WithAllPrivileges(true).
		WithOnAccountObject(sdk.ObjectTypeWarehouse, wh.ID())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantPrivileges_SafeDestroy")

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
			// Drop the role externally, then try to destroy without experiment — expect failure.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().Role.DropRoleFunc(t, role.ID()),
				Config:                   config.FromModels(t, grantModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with GRANTS_SAFE_DESTROY experiment — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				Destroy:                  true,
			},
		},
	})
}

// TestAcc_Experimental_GrantPrivilegesToDatabaseRole_SafeDestroy_MissingDatabase verifies that destroying
// a database role grant resource fails when the target database is deleted externally (default behavior),
// and succeeds when the GRANTS_SAFE_DESTROY experiment is enabled.
// Uses all_privileges = true so that Read skips existence checks and Delete is actually called.
func TestAcc_Experimental_GrantPrivilegesToDatabaseRole_SafeDestroy_MissingDatabase(t *testing.T) {
	// Create a separate database to grant on so we can drop it independently of the role.
	db, dbCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(dbCleanup)

	dbRole, dbRoleCleanup := testClient().DatabaseRole.CreateDatabaseRoleInDatabase(t, db.ID())
	t.Cleanup(dbRoleCleanup)

	grantModel := model.GrantPrivilegesToDatabaseRole("test", dbRole.ID().FullyQualifiedName()).
		WithAllPrivileges(true).
		WithOnDatabase(db.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantPrivileges_SafeDestroy")

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
			// Drop the target database externally, then try to destroy without experiment — expect failure.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().Database.DropDatabaseFunc(t, db.ID()),
				Config:                   config.FromModels(t, grantModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with GRANTS_SAFE_DESTROY experiment — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				Destroy:                  true,
			},
		},
	})
}

// TestAcc_Experimental_GrantPrivilegesToDatabaseRole_SafeDestroy_MissingDatabaseRole verifies that destroying
// a database role grant resource fails when the grantee database role is deleted externally (default behavior),
// and succeeds when the GRANTS_SAFE_DESTROY experiment is enabled.
// Uses all_privileges = true so that Read skips existence checks and Delete is actually called.
func TestAcc_Experimental_GrantPrivilegesToDatabaseRole_SafeDestroy_MissingDatabaseRole(t *testing.T) {
	dbRole, dbRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	grantModel := model.GrantPrivilegesToDatabaseRole("test", dbRole.ID().FullyQualifiedName()).
		WithAllPrivileges(true).
		WithOnDatabase(testClient().Ids.DatabaseId().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)
	experimentFactory := providerFactoryUsingCache("TestAcc_Experimental_GrantPrivileges_SafeDestroy")

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
			// Drop the database role externally, then try to destroy without experiment — expect failure.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().DatabaseRole.CleanupDatabaseRoleFunc(t, dbRole.ID()),
				Config:                   config.FromModels(t, grantModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with GRANTS_SAFE_DESTROY experiment — succeeds.
			{
				ProtoV6ProviderFactories: experimentFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				Destroy:                  true,
			},
		},
	})
}
