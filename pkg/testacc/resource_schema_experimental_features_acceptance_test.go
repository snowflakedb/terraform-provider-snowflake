//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Database A is renamed to B in config. The schema resource detects the rename
// and updates its ID without performing any Snowflake modification.
func TestAcc_Experimental_Schema_HierarchyRenames_DatabaseRenamed(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModel := model.DatabaseWithParametersSet("db", databaseId.Name())
	databaseModelRenamed := model.DatabaseWithParametersSet("db", newDatabaseId.Name())

	schemaModel := model.SchemaWithImplicitDatabaseDependency("test", schemaName, databaseModel)

	expectedNewSchemaId := sdk.NewDatabaseObjectIdentifier(newDatabaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create database and schema with implicit dependency
			{
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModel),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModel.ResourceReference()).
						HasNameString(schemaName).
						HasDatabaseString(databaseId.Name()),
				),
			},
			// Rename database in config — schema detects rename (Case 1) and updates ID
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(databaseModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(schemaModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelRenamed, schemaModel),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModel.ResourceReference()).
						HasNameString(schemaName).
						HasDatabaseString(newDatabaseId.Name()).
						HasFullyQualifiedNameString(expectedNewSchemaId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Both databases exist in config. The schema is moved from database A to database B via ALTER SCHEMA RENAME TO.
func TestAcc_Experimental_Schema_HierarchyRenames_SchemaMove(t *testing.T) {
	databaseAId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseBId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModelA := model.DatabaseWithParametersSet("dbA", databaseAId.Name())
	databaseModelB := model.DatabaseWithParametersSet("dbB", databaseBId.Name())

	schemaModelBefore := model.SchemaWithImplicitDatabaseDependency("test", schemaName, databaseModelA)
	schemaModelAfter := model.SchemaWithImplicitDatabaseDependency("test", schemaName, databaseModelB)

	expectedNewSchemaId := sdk.NewDatabaseObjectIdentifier(databaseBId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create both databases and schema in database A
			{
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelBefore),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelBefore.ResourceReference()).
						HasNameString(schemaName).
						HasDatabaseString(databaseAId.Name()),
				),
			},
			// Move schema to database B
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(schemaModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelAfter),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelAfter.ResourceReference()).
						HasNameString(schemaName).
						HasDatabaseString(databaseBId.Name()).
						HasFullyQualifiedNameString(expectedNewSchemaId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Database A is renamed to B in config, and the schema name is also changed from X to Y.
func TestAcc_Experimental_Schema_HierarchyRenames_DatabaseRenamed_WithNameChange(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()
	newSchemaName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModel := model.DatabaseWithParametersSet("db", databaseId.Name())
	databaseModelRenamed := model.DatabaseWithParametersSet("db", newDatabaseId.Name())

	schemaModelBefore := model.SchemaWithImplicitDatabaseDependency("test", schemaName, databaseModel)
	schemaModelAfter := model.SchemaWithImplicitDatabaseDependency("test", newSchemaName, databaseModel)

	expectedNewSchemaId := sdk.NewDatabaseObjectIdentifier(newDatabaseId.Name(), newSchemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create database and schema X with implicit dependency
			{
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModelBefore),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelBefore.ResourceReference()).
						HasNameString(schemaName).
						HasDatabaseString(databaseId.Name()),
				),
			},
			// Rename database in config and change schema name
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(databaseModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(schemaModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelRenamed, schemaModelAfter),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelAfter.ResourceReference()).
						HasNameString(newSchemaName).
						HasDatabaseString(newDatabaseId.Name()).
						HasFullyQualifiedNameString(expectedNewSchemaId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Both databases exist in config. Schema is moved from A to B and renamed from X to Y.
func TestAcc_Experimental_Schema_HierarchyRenames_SchemaMove_WithNameChange(t *testing.T) {
	databaseAId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseBId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()
	newSchemaName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModelA := model.DatabaseWithParametersSet("dbA", databaseAId.Name())
	databaseModelB := model.DatabaseWithParametersSet("dbB", databaseBId.Name())

	schemaModelBefore := model.SchemaWithImplicitDatabaseDependency("test", schemaName, databaseModelA)
	schemaModelAfter := model.SchemaWithImplicitDatabaseDependency("test", newSchemaName, databaseModelB)

	expectedNewSchemaId := sdk.NewDatabaseObjectIdentifier(databaseBId.Name(), newSchemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create both databases and schema X in database A
			{
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelBefore),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelBefore.ResourceReference()).
						HasNameString(schemaName).
						HasDatabaseString(databaseAId.Name()),
				),
			},
			// Move schema to database B and rename to Y
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(schemaModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelAfter),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelAfter.ResourceReference()).
						HasNameString(newSchemaName).
						HasDatabaseString(databaseBId.Name()).
						HasFullyQualifiedNameString(expectedNewSchemaId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Move to the database which does not exist.
func TestAcc_Experimental_Schema_HierarchyRenames_Error_NewDatabaseDoesNotExist(t *testing.T) {
	dbA, cleanupDbA := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDbA)

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(dbA.ID())
	nonExistentDbName := testClient().Ids.RandomAccountObjectIdentifier().Name()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	schemaModelBefore := model.Schema("test", dbA.ID().Name(), schemaId.Name())
	schemaModelAfter := model.Schema("test", nonExistentDbName, schemaId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create schema in database A
			{
				Config: accconfig.FromModels(t, providerModel, schemaModelBefore),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelBefore.ResourceReference()).
						HasNameString(schemaId.Name()).
						HasDatabaseString(dbA.ID().Name()),
				),
			},
			// Try to move to non-existent database
			{
				Config:      accconfig.FromModels(t, providerModel, schemaModelAfter),
				ExpectError: regexp.MustCompile(`unknown rename use case.*object_renaming_guide`),
			},
		},
	})
}

// Move to a schema that has been dropped externally.
func TestAcc_Experimental_Schema_HierarchyRenames_Error_SchemaNotFoundForMove(t *testing.T) {
	dbA, cleanupDbA := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDbA)

	dbB, cleanupDbB := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDbB)

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(dbA.ID())

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	schemaModelBefore := model.Schema("test", dbA.ID().Name(), schemaId.Name())
	schemaModelAfter := model.Schema("test", dbB.ID().Name(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create schema in database A
			{
				Config: accconfig.FromModels(t, providerModel, schemaModelBefore),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelBefore.ResourceReference()).
						HasNameString(schemaId.Name()).
						HasDatabaseString(dbA.ID().Name()),
				),
			},
			// Drop schema externally, then try to move
			{
				// Drop schema externally, after reading
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(func() {
							testClient().Schema.Alter(t, sdk.NewAlterSchemaRequest(schemaId).WithNewName(sdk.NewDatabaseObjectIdentifier(dbA.ID().Name(), testClient().Ids.Alpha())))
						}),
					},
				},
				Config:      accconfig.FromModels(t, providerModel, schemaModelAfter),
				ExpectError: regexp.MustCompile(`unknown rename use case.*object_renaming_guide`),
			},
		},
	})
}

// When the experiment is NOT enabled, changing the database field still forces recreation (existing behavior preserved).
func TestAcc_Experimental_Schema_HierarchyRenames_Disabled_ForceRecreation(t *testing.T) {
	databaseAId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseBId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()

	databaseModelA := model.DatabaseWithParametersSet("dbA", databaseAId.Name())
	databaseModelB := model.DatabaseWithParametersSet("dbB", databaseBId.Name())

	// No experiment enabled — use default provider
	schemaModelBefore := model.SchemaWithImplicitDatabaseDependency("test", schemaName, databaseModelA)
	schemaModelAfter := model.SchemaWithImplicitDatabaseDependency("test", schemaName, databaseModelB)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create both databases and schema in database A
			{
				Config: accconfig.FromModels(t, databaseModelA, databaseModelB, schemaModelBefore),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelBefore.ResourceReference()).
						HasNameString(schemaName).
						HasDatabaseString(databaseAId.Name()),
				),
			},
			// Change database — should force recreation
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(schemaModelAfter.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, databaseModelA, databaseModelB, schemaModelAfter),
				Check: assertThat(
					t,
					resourceassert.SchemaResource(t, schemaModelAfter.ResourceReference()).
						HasNameString(schemaName).
						HasDatabaseString(databaseBId.Name()),
				),
			},
		},
	})
}
