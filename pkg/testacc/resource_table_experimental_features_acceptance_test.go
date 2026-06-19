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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var tableColumns = []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeNumber}}

// Database A is renamed to B in config. The table resource detects the rename
// and updates its ID without performing any Snowflake modification.
func TestAcc_Experimental_Table_HierarchyRenames_DatabaseRenamed(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModel := model.DatabaseWithParametersSet("db", databaseId.Name())
	databaseModelRenamed := model.DatabaseWithParametersSet("db", newDatabaseId.Name())

	schemaModel := model.SchemaWithImplicitDatabaseDependency("schema", schemaName, databaseModel)
	schemaModelAfterDbRename := model.SchemaWithImplicitDatabaseDependency("schema", schemaName, databaseModelRenamed)

	tableModel := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModel, databaseModel)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelAfterDbRename, databaseModelRenamed)

	expectedNewTableId := sdk.NewSchemaObjectIdentifier(newDatabaseId.Name(), schemaName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create database, schema, and table
			{
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModel, tableModel),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModel.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseId.Name()).
						HasSchemaString(schemaName),
				),
			},
			// Rename database in config — table detects rename (Case A) and updates ID
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(databaseModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelRenamed, schemaModelAfterDbRename, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(newDatabaseId.Name()).
						HasSchemaString(schemaName).
						HasFullyQualifiedNameString(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Both databases exist in config. The table is moved from database A to database B via ALTER TABLE RENAME TO.
func TestAcc_Experimental_Table_HierarchyRenames_TableMoveToAnotherDatabase(t *testing.T) {
	databaseAId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseBId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModelA := model.DatabaseWithParametersSet("dbA", databaseAId.Name())
	databaseModelB := model.DatabaseWithParametersSet("dbB", databaseBId.Name())

	schemaModelA := model.SchemaWithImplicitDatabaseDependency("schemaA", schemaName, databaseModelA)
	schemaModelB := model.SchemaWithImplicitDatabaseDependency("schemaB", schemaName, databaseModelB)

	tableModelBefore := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelA, databaseModelA)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelB, databaseModelB)

	expectedNewTableId := sdk.NewSchemaObjectIdentifier(databaseBId.Name(), schemaName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create both databases, schemas, and table in database A
			{
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelA, schemaModelB, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseAId.Name()).
						HasSchemaString(schemaName),
				),
			},
			// Move table to database B
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelA, schemaModelB, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseBId.Name()).
						HasSchemaString(schemaName).
						HasFullyQualifiedNameString(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Schema X is renamed to Y in config. The table resource detects the rename
// and updates its ID without performing any Snowflake modification.
func TestAcc_Experimental_Table_HierarchyRenames_SchemaRenamed(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()
	newSchemaName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModel := model.DatabaseWithParametersSet("db", databaseId.Name())

	schemaModel := model.SchemaWithImplicitDatabaseDependency("schema", schemaName, databaseModel)
	schemaModelRenamed := model.SchemaWithImplicitDatabaseDependency("schema", newSchemaName, databaseModel)

	tableModel := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModel, databaseModel)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelRenamed, databaseModel)

	expectedNewTableId := sdk.NewSchemaObjectIdentifier(databaseId.Name(), newSchemaName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create database, schema, and table
			{
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModel, tableModel),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModel.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseId.Name()).
						HasSchemaString(schemaName),
				),
			},
			// Rename schema in config — table detects rename (Case B) and updates ID
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(schemaModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModelRenamed, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseId.Name()).
						HasSchemaString(newSchemaName).
						HasFullyQualifiedNameString(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Both schemas exist in config. The table is moved from schema X to schema Y via ALTER TABLE RENAME TO.
func TestAcc_Experimental_Table_HierarchyRenames_TableMoveToAnotherSchema(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaXName := testClient().Ids.Alpha()
	schemaYName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModel := model.DatabaseWithParametersSet("db", databaseId.Name())

	schemaModelX := model.SchemaWithImplicitDatabaseDependency("schemaX", schemaXName, databaseModel)
	schemaModelY := model.SchemaWithImplicitDatabaseDependency("schemaY", schemaYName, databaseModel)

	tableModelBefore := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelX, databaseModel)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelY, databaseModel)

	expectedNewTableId := sdk.NewSchemaObjectIdentifier(databaseId.Name(), schemaYName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create database, both schemas, and table in schema X
			{
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModelX, schemaModelY, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseId.Name()).
						HasSchemaString(schemaXName),
				),
			},
			// Move table to schema Y
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModelX, schemaModelY, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseId.Name()).
						HasSchemaString(schemaYName).
						HasFullyQualifiedNameString(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Database A is renamed to B AND schema X is renamed to Y simultaneously.
// Table detects both renames and updates ID only.
func TestAcc_Experimental_Table_HierarchyRenames_DatabaseRenamed_SchemaRenamed(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()
	newSchemaName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModel := model.DatabaseWithParametersSet("db", databaseId.Name())
	databaseModelRenamed := model.DatabaseWithParametersSet("db", newDatabaseId.Name())

	schemaModel := model.SchemaWithImplicitDatabaseDependency("schema", schemaName, databaseModel)
	schemaModelRenamed := model.SchemaWithImplicitDatabaseDependency("schema", newSchemaName, databaseModelRenamed)

	tableModel := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModel, databaseModel)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelRenamed, databaseModelRenamed)

	expectedNewTableId := sdk.NewSchemaObjectIdentifier(newDatabaseId.Name(), newSchemaName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create database, schema, and table
			{
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModel, tableModel),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModel.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseId.Name()).
						HasSchemaString(schemaName),
				),
			},
			// Rename both database and schema — table detects both renames (Case C Scenario 1)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(databaseModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(schemaModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelRenamed, schemaModelRenamed, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(newDatabaseId.Name()).
						HasSchemaString(newSchemaName).
						HasFullyQualifiedNameString(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Database A is renamed to B. Table is moved from schema X to schema Y within renamed database.
func TestAcc_Experimental_Table_HierarchyRenames_DatabaseRenamed_SchemaMove(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaXName := testClient().Ids.Alpha()
	schemaYName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModel := model.DatabaseWithParametersSet("db", databaseId.Name())
	databaseModelRenamed := model.DatabaseWithParametersSet("db", newDatabaseId.Name())

	schemaModelX := model.SchemaWithImplicitDatabaseDependency("schemaX", schemaXName, databaseModel)
	schemaModelY := model.SchemaWithImplicitDatabaseDependency("schemaY", schemaYName, databaseModel)
	schemaModelXAfter := model.SchemaWithImplicitDatabaseDependency("schemaX", schemaXName, databaseModelRenamed)
	schemaModelYAfter := model.SchemaWithImplicitDatabaseDependency("schemaY", schemaYName, databaseModelRenamed)

	tableModelBefore := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelX, databaseModel)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelYAfter, databaseModelRenamed)

	expectedNewTableId := sdk.NewSchemaObjectIdentifier(newDatabaseId.Name(), schemaYName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create database, both schemas, and table in schema X
			{
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModelX, schemaModelY, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseId.Name()).
						HasSchemaString(schemaXName),
				),
			},
			// Rename database AND move table to different schema (Case C Scenario 2)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(databaseModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelRenamed, schemaModelXAfter, schemaModelYAfter, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(newDatabaseId.Name()).
						HasSchemaString(schemaYName).
						HasFullyQualifiedNameString(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Both databases exist. Schema X is renamed to Y in old database. Table is moved to new database.
func TestAcc_Experimental_Table_HierarchyRenames_DatabaseMove_SchemaRenamed(t *testing.T) {
	databaseAId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseBId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()
	newSchemaName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModelA := model.DatabaseWithParametersSet("dbA", databaseAId.Name())
	databaseModelB := model.DatabaseWithParametersSet("dbB", databaseBId.Name())

	schemaModelBefore := model.SchemaWithImplicitDatabaseDependency("schema", schemaName, databaseModelA)
	schemaModelRenamed := model.SchemaWithImplicitDatabaseDependency("schema", newSchemaName, databaseModelA)
	schemaModelTarget := model.SchemaWithImplicitDatabaseDependency("schemaTarget", newSchemaName, databaseModelB)

	tableModelBefore := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelBefore, databaseModelA)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelTarget, databaseModelB)

	expectedNewTableId := sdk.NewSchemaObjectIdentifier(databaseBId.Name(), newSchemaName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create both databases, schema in A, target schema in B, and table in A.schema
			{
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelBefore, schemaModelTarget, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseAId.Name()).
						HasSchemaString(schemaName),
				),
			},
			// Rename schema in A AND move table to B (Case C Scenario 3)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(schemaModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelRenamed, schemaModelTarget, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseBId.Name()).
						HasSchemaString(newSchemaName).
						HasFullyQualifiedNameString(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Both databases exist. Both schemas exist. Table moved from A.X to B.Y.
func TestAcc_Experimental_Table_HierarchyRenames_DatabaseMove_SchemaMove(t *testing.T) {
	databaseAId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseBId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaXName := testClient().Ids.Alpha()
	schemaYName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModelA := model.DatabaseWithParametersSet("dbA", databaseAId.Name())
	databaseModelB := model.DatabaseWithParametersSet("dbB", databaseBId.Name())

	schemaModelAX := model.SchemaWithImplicitDatabaseDependency("schemaAX", schemaXName, databaseModelA)
	schemaModelBY := model.SchemaWithImplicitDatabaseDependency("schemaBY", schemaYName, databaseModelB)

	tableModelBefore := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelAX, databaseModelA)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelBY, databaseModelB)

	expectedNewTableId := sdk.NewSchemaObjectIdentifier(databaseBId.Name(), schemaYName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create both databases, both schemas, and table in A.X
			{
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelAX, schemaModelBY, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseAId.Name()).
						HasSchemaString(schemaXName),
				),
			},
			// Move table from A.X to B.Y (Case C Scenario 4)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModelA, databaseModelB, schemaModelAX, schemaModelBY, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseBId.Name()).
						HasSchemaString(schemaYName).
						HasFullyQualifiedNameString(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Error case: target database does not exist.
func TestAcc_Experimental_Table_HierarchyRenames_Error_NewDatabaseDoesNotExist(t *testing.T) {
	dbA, cleanupDbA := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDbA)

	schema, cleanupSchema := testClient().Schema.CreateSchemaInDatabase(t, dbA.ID())
	t.Cleanup(cleanupSchema)

	tableName := testClient().Ids.Alpha()
	nonExistentDbName := testClient().Ids.RandomAccountObjectIdentifier().Name()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	tableModelBefore := model.Table("test", dbA.ID().Name(), schema.ID().Name(), tableName, tableColumns)
	tableModelAfter := model.Table("test", nonExistentDbName, schema.ID().Name(), tableName, tableColumns)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create table
			{
				Config: accconfig.FromModels(t, providerModel, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(dbA.ID().Name()),
				),
			},
			// Try to move to non-existent database
			{
				Config:      accconfig.FromModels(t, providerModel, tableModelAfter),
				ExpectError: regexp.MustCompile(`unknown rename use case.*object_renaming_guide`),
			},
		},
	})
}

// Error case: target schema does not exist.
func TestAcc_Experimental_Table_HierarchyRenames_Error_NewSchemaDoesNotExist(t *testing.T) {
	dbA, cleanupDbA := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDbA)

	schema, cleanupSchema := testClient().Schema.CreateSchemaInDatabase(t, dbA.ID())
	t.Cleanup(cleanupSchema)

	tableName := testClient().Ids.Alpha()
	nonExistentSchemaName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	tableModelBefore := model.Table("test", dbA.ID().Name(), schema.ID().Name(), tableName, tableColumns)
	tableModelAfter := model.Table("test", dbA.ID().Name(), nonExistentSchemaName, tableName, tableColumns)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create table
			{
				Config: accconfig.FromModels(t, providerModel, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasSchemaString(schema.ID().Name()),
				),
			},
			// Try to move to non-existent schema
			{
				Config:      accconfig.FromModels(t, providerModel, tableModelAfter),
				ExpectError: regexp.MustCompile(`unknown rename use case.*object_renaming_guide`),
			},
		},
	})
}

// When the experiment is NOT enabled, changing the database or schema field still forces recreation (existing behavior preserved).
func TestAcc_Experimental_Table_HierarchyRenames_Disabled_ForceRecreation(t *testing.T) {
	databaseAId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseBId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	databaseModelA := model.DatabaseWithParametersSet("dbA", databaseAId.Name())
	databaseModelB := model.DatabaseWithParametersSet("dbB", databaseBId.Name())

	schemaModelA := model.SchemaWithImplicitDatabaseDependency("schemaA", schemaName, databaseModelA)
	schemaModelB := model.SchemaWithImplicitDatabaseDependency("schemaB", schemaName, databaseModelB)

	// No experiment enabled — use default provider
	tableModelBefore := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelA, databaseModelA)
	tableModelAfter := model.TableWithImplicitDependencies("test", tableName, tableColumns, schemaModelB, databaseModelB)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create both databases, schemas, and table in database A
			{
				Config: accconfig.FromModels(t, databaseModelA, databaseModelB, schemaModelA, schemaModelB, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseAId.Name()),
				),
			},
			// Change database — should force recreation
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, databaseModelA, databaseModelB, schemaModelA, schemaModelB, tableModelAfter),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelAfter.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(databaseBId.Name()),
				),
			},
		},
	})
}

// Externally dropped schema: table move fails with unknown rename use case.
func TestAcc_Experimental_Table_HierarchyRenames_Error_SchemaDroppedExternally(t *testing.T) {
	dbA, cleanupDbA := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDbA)

	dbB, cleanupDbB := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDbB)

	schema, cleanupSchema := testClient().Schema.CreateSchemaInDatabase(t, dbA.ID())
	t.Cleanup(cleanupSchema)

	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	tableModelBefore := model.Table("test", dbA.ID().Name(), schema.ID().Name(), tableName, tableColumns)
	tableModelAfter := model.Table("test", dbB.ID().Name(), schema.ID().Name(), tableName, tableColumns)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create table in database A
			{
				Config: accconfig.FromModels(t, providerModel, tableModelBefore),
				Check: assertThat(t,
					resourceassert.TableResource(t, tableModelBefore.ResourceReference()).
						HasNameString(tableName).
						HasDatabaseString(dbA.ID().Name()),
				),
			},
			// Drop schema externally, then try to move
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.Execute(func() {
							testClient().Schema.Alter(t, sdk.NewAlterSchemaRequest(schema.ID()).
								WithNewName(sdk.NewDatabaseObjectIdentifier(dbA.ID().Name(), testClient().Ids.Alpha())),
							)
						}),
					},
				},
				Config:      accconfig.FromModels(t, providerModel, tableModelAfter),
				ExpectError: regexp.MustCompile(`unknown rename use case.*object_renaming_guide`),
			},
		},
	})
}
