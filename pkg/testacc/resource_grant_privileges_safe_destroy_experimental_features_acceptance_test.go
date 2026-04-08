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
				ProtoV6ProviderFactories: grantsSafeDestroyProviderFactory,
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
				ProtoV6ProviderFactories: grantsSafeDestroyProviderFactory,
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
				ProtoV6ProviderFactories: grantsSafeDestroyProviderFactory,
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
				ProtoV6ProviderFactories: grantsSafeDestroyProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				Destroy:                  true,
			},
		},
	})
}

// TestAcc_Experimental_GrantPrivilegesToShare_SafeDestroy_MissingShare verifies that destroying
// a grant resource succeeds when the share (grantee) is deleted externally, even without
// the GRANTS_SAFE_DESTROY experiment. Unlike MissingSchema, this case is handled gracefully by
// ReadGrantPrivilegesToShare itself: it always checks share existence via ShowByID first, removes
// the resource from state, and returns before Delete is ever called — so no experiment is needed.
func TestAcc_Experimental_GrantPrivilegesToShare_SafeDestroy_MissingShare(t *testing.T) {
	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, database.ID(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeUsage.String()}, share.ID().Name()).
		WithOnDatabase(database.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with default provider.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.1"),
				Config:            config.FromModels(t, grantModel),
			},
			// Drop the share externally, then try to destroy without experiment — expect failure.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.1"),
				PreConfig:         testClient().Share.DropShareFunc(t, share.ID()),
				Config:            config.FromModels(t, grantModel),
				Destroy:           true,
				ExpectError:       regexp.MustCompile(`revokePrivilegeFromShareOptions fields`),
			},
			// Read detects the missing share via ShowByID and clears the resource from state
			// before Delete is called.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
				Destroy:                  true,
			},
		},
	})
}

// TestAcc_Experimental_GrantPrivilegesToShare_SafeDestroy_MissingSchema verifies that destroying
// a grant resource fails when the target schema is deleted externally (default behavior), and succeeds
// when the GRANTS_SAFE_DESTROY experiment is enabled.
// Uses on_all_tables_in_schema so that Read skips existence checks (returns nil when opts == nil)
// and Delete is actually called even when the schema no longer exists.
func TestAcc_Experimental_GrantPrivilegesToShare_SafeDestroy_MissingSchema(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	// setup grants USAGE on database to share — required by Snowflake before granting schema objects.
	testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, database.ID(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})

	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeSelect.String()}, share.ID().Name()).
		WithOnAllTablesInSchema(schema.ID().FullyQualifiedName())
	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create both grants with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the schema externally, then try to destroy without experiment — expect failure.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().Schema.DropSchemaFunc(t, schema.ID()),
				Config:                   config.FromModels(t, grantModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with GRANTS_SAFE_DESTROY experiment — succeeds.
			{
				ProtoV6ProviderFactories: grantsSafeDestroyProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				Destroy:                  true,
			},
		},
	})
}
