//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

/*
The following tests are showing the behavior of the provider in cases where objects higher in the hierarchy
like database or schema are renamed when the objects lower in the hierarchy are in the Terraform configuration.
For more information check TODO(SNOW-1672319): link public document.

Shallow hierarchy (database + schema)
- is in config - renamed internally - with implicit dependency (works)

- is in config - renamed internally - without dependency - after rename schema referencing old database name (fails in Read and then it's failing to remove itself in Delete)
- is in config - renamed internally - without dependency - after rename schema referencing new database name (fails in Read and then it's failing to remove itself in Delete)

- is in config - renamed internally - with depends_on - after rename schema referencing old database name (fails in Read and then it's failing to remove itself in Delete)
- is in config - renamed internally - with depends_on - after rename schema referencing new database name (works)

- is in config - renamed externally - with implicit dependency - database holding the same name in config (works; creates a new database with a schema next to the already existing renamed database and schema)
- is in config - renamed externally - with implicit dependency - database holding the new name in config  (fails to create database, because it already exists)

- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing old database name (non-deterministic results depending on the Terraform execution order that seems to be different with every run)
- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing old database name - check impact of the resource order on the plan (seems to fail pretty consistently in Delete because database is dropped before schema)
- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing new database name (fails because schema resource tries to create a new schema that already exists in renamed database)
- is in config - renamed externally - without dependency - after rename database referencing new name and schema referencing old database name (fails because database resource tried to create database that already exists)
- is in config - renamed externally - without dependency - after rename database referencing new name and schema referencing new database name (fails because database resource tried to create database that already exists)

- is in config - renamed externally - with depends_on - after rename database referencing old name and schema referencing old database name (works; creates a new database with a schema next to the already existing renamed database and schema)
- is in config - renamed externally - with depends_on - after rename database referencing old name and schema referencing new database name (fails because schema resource tries to create a new schema that already exists in renamed database)
- is in config - renamed externally - with depends_on - after rename database referencing new name and schema referencing old database name (fails because database resource tried to create database that already exists)
- is in config - renamed externally - with depends_on - after rename database referencing new name and schema referencing new database name (fails because database resource tried to create database that already exists)

- is not in config - renamed externally - referencing old database name (fails because it tries to create a new schema on non-existing database)
- is not in config - renamed externally - referencing new database name (fails because schema resource tries to create a new schema that already exists in renamed database)

Deep hierarchy (database + schema + table)

- are in config - database renamed internally - with database implicit dependency - with no schema dependency 		- with database to schema implicit dependency (fails because table is created before schema)
- are in config - database renamed internally - with database implicit dependency - with implicit schema dependency - with database to schema implicit dependency (works)
- are in config - database renamed internally - with database implicit dependency - with schema depends_on 			- with database to schema implicit dependency (works)

- are in config - database renamed internally - with database implicit dependency - with no schema dependency 		- with database to schema depends_on dependency (fails because table is created before schema)
- are in config - database renamed internally - with database implicit dependency - with implicit schema dependency - with database to schema depends_on dependency (works)
- are in config - database renamed internally - with database implicit dependency - with schema depends_on 			- with database to schema depends_on dependency (works)

- are in config - database renamed internally - with database implicit dependency - with no schema dependency 		- with database to schema no dependency (fails during delete because database is deleted before schema)
- are in config - database renamed internally - with database implicit dependency - with implicit schema dependency - with database to schema no dependency (fails to drop schema after database rename)
- are in config - database renamed internally - with database implicit dependency - with schema depends_on 			- with database to schema no dependency (fails to drop schema after database rename)

- are in config - database renamed internally - with no database dependency - with no schema dependency 	  (fails because table is created before schema)
- are in config - database renamed internally - with no database dependency - with implicit schema dependency (works)
- are in config - database renamed internally - with no database dependency - with schema depends_on 		  (works)

- are in config - database renamed internally - with database depends_on - with no schema dependency 	   (fails because table is created before schema)
- are in config - database renamed internally - with database depends_on - with implicit schema dependency (works)
- are in config - database renamed internally - with database depends_on - with schema depends_on 		   (works)

------------------------------------------------------------------------------------------------------------------------

- are in config - schema renamed internally - with database implicit dependency - with no schema dependency 	  (fails because table is created before schema)
- are in config - schema renamed internally - with database implicit dependency - with implicit schema dependency (works)
- are in config - schema renamed internally - with database implicit dependency - with schema depends_on 		  (works)

- are in config - schema renamed internally - with no database dependency - with no schema dependency 		(fails because table is created before schema)
- are in config - schema renamed internally - with no database dependency - with implicit schema dependency (works)
- are in config - schema renamed internally - with no database dependency - with schema depends_on 			(works)

- are in config - schema renamed internally - with database depends_on - with no schema dependency 		 (fails because table is created before schema)
- are in config - schema renamed internally - with database depends_on - with implicit schema dependency (works)
- are in config - schema renamed internally - with database depends_on - with schema depends_on 		 (works)

------------------------------------------------------------------------------------------------------------------------

- are in config - database renamed externally - with database implicit dependency - with no schema dependency 		(fails because table is created before schema)
- are in config - database renamed externally - with database implicit dependency - with implicit schema dependency (fails because tries to create database when it's already there after rename)
- are in config - database renamed externally - with database implicit dependency - with schema depends_on 			(fails because tries to create database when it's already there after rename)

- are in config - database renamed externally - with no database dependency - with no schema dependency 	  (fails because table is created before schema)
- are in config - database renamed externally - with no database dependency - with implicit schema dependency (fails because tries to create database when it's already there after rename)
- are in config - database renamed externally - with no database dependency - with schema depends_on 		  (fails because tries to create database when it's already there after rename)

- are in config - database renamed externally - with database depends_on - with no schema dependency 	   (fails because table is created before schema)
- are in config - database renamed externally - with database depends_on - with implicit schema dependency (fails because tries to create database when it's already there after rename)
- are in config - database renamed externally - with database depends_on - with schema depends_on 		   (fails because tries to create database when it's already there after rename)

------------------------------------------------------------------------------------------------------------------------

- are in config - schema renamed externally - with database implicit dependency - with no schema dependency 	  (fails because table is created before schema)
- are in config - schema renamed externally - with database implicit dependency - with implicit schema dependency (fails because tries to create database when it's already there after rename)
- are in config - schema renamed externally - with database implicit dependency - with schema depends_on 		  (fails because tries to create database when it's already there after rename)

- are in config - schema renamed externally - with no database dependency - with no schema dependency 		(fails because table is created before schema)
- are in config - schema renamed externally - with no database dependency - with implicit schema dependency (fails because tries to create database when it's already there after rename)
- are in config - schema renamed externally - with no database dependency - with schema depends_on 			(fails because tries to create database when it's already there after rename)

- are in config - schema renamed externally - with database depends_on - with no schema dependency 		 (fails because table is created before schema)
- are in config - schema renamed externally - with database depends_on - with implicit schema dependency (fails because tries to create database when it's already there after rename)
- are in config - schema renamed externally - with database depends_on - with schema depends_on 		 (fails because tries to create database when it's already there after rename)

------------------------------------------------------------------------------------------------------------------------

- are not in config - database renamed externally - referencing old database name (fails because tries to create table on non-existing database)
- are not in config - database renamed externally - referencing new database name (fails because tries to create table that already exists in the renamed database)

- are not in config - schema renamed externally - referencing old schema name (fails because tries to create table on non-existing schema)
- are not in config - schema renamed externally - referencing new schema name (fails because tries to create table that already exists in the renamed schema)

# The list of test cases that were not added:
- (Deep hierarchy) More test cases with varying dependencies between resources
- (Deep hierarchy) Add test cases where old database is referenced to see if hierarchy recreation is possible
- (Deep hierarchy) More test cases could be added when database and schema are renamed at the same time
- (Deep hierarchy) More test cases could be added when either database or schema are in the config
*/

type DependencyType string

const (
	ImplicitDependency  DependencyType = "implicit"
	DependsOnDependency DependencyType = "depends_on"
	NoDependency        DependencyType = "no_dependency"
)

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		Dependency                             DependencyType
		IsRenamedDatabaseReferencedAfterRename bool
	}{
		{Dependency: ImplicitDependency},
		// Error happens during schema's Read operation and then Delete operation (schema cannot be removed).
		// Not able to handle the error produced by Delete operation that results in test always failing
		// {Dependency: NoDependency, IsRenamedDatabaseReferencedAfterRename: false},
		// {Dependency: NoDependency, IsRenamedDatabaseReferencedAfterRename: true},
		// {Dependency: DependsOnDependency, IsRenamedDatabaseReferencedAfterRename: false},
		// {Dependency: DependsOnDependency, IsRenamedDatabaseReferencedAfterRename: true},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestAcc_ dependency: %s", testCase.Dependency), func(t *testing.T) {
			databaseId := testClient().Ids.RandomAccountObjectIdentifier()
			newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
			schemaName := testClient().Ids.Alpha()

			databaseConfigModel := model.DatabaseWithParametersSet("test", databaseId.Name())
			databaseConfigModelWithNewId := model.DatabaseWithParametersSet("test", newDatabaseId.Name())

			var schemaConfigAfterRename string
			var preApplyChecksAfterRename []plancheck.PlanCheck

			switch testCase.Dependency {
			case ImplicitDependency:
				schemaConfigAfterRename = configSchemaWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.Dependency, databaseId.Name(), schemaName)
				preApplyChecksAfterRename = []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
					plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
				}
			case DependsOnDependency, NoDependency:
				if testCase.IsRenamedDatabaseReferencedAfterRename {
					schemaConfigAfterRename = configSchemaWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.Dependency, newDatabaseId.Name(), schemaName)
					preApplyChecksAfterRename = []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					}
				} else {
					schemaConfigAfterRename = configSchemaWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.Dependency, databaseId.Name(), schemaName)
					preApplyChecksAfterRename = []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					}
				}
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Schema),
				Steps: []resource.TestStep{
					{
						Config: config.FromModels(t, databaseConfigModel) + configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.Dependency, databaseId.Name(), schemaName),
					},
					{
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: preApplyChecksAfterRename,
						},
						Config: config.FromModels(t, databaseConfigModelWithNewId) + schemaConfigAfterRename,
					},
				},
			})
		})
	}
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		Dependency                                       DependencyType
		IsRenamedDatabaseReferencedAfterRenameInDatabase bool
		IsRenamedDatabaseReferencedAfterRenameInSchema   bool
		ExpectedError                                    *regexp.Regexp
	}{
		// errors with already exists happen, because we try to create databases or schemas when there's already an existing object in Snowflake with the same name.
		{Dependency: ImplicitDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: false},
		{Dependency: ImplicitDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: true, ExpectedError: regexp.MustCompile(`Object '.*' already exists`)},

		{Dependency: NoDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: false, IsRenamedDatabaseReferencedAfterRenameInSchema: true, ExpectedError: regexp.MustCompile("Failed to create schema")},
		{Dependency: NoDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: true, IsRenamedDatabaseReferencedAfterRenameInSchema: false, ExpectedError: regexp.MustCompile(`Object '.*' already exists`)},
		{Dependency: NoDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: true, IsRenamedDatabaseReferencedAfterRenameInSchema: true, ExpectedError: regexp.MustCompile(`Object '.*' already exists`)},

		{Dependency: DependsOnDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: false, IsRenamedDatabaseReferencedAfterRenameInSchema: false},
		{Dependency: DependsOnDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: false, IsRenamedDatabaseReferencedAfterRenameInSchema: true, ExpectedError: regexp.MustCompile("Failed to create schema")},
		{Dependency: DependsOnDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: true, IsRenamedDatabaseReferencedAfterRenameInSchema: false, ExpectedError: regexp.MustCompile(`Object '.*' already exists`)},
		{Dependency: DependsOnDependency, IsRenamedDatabaseReferencedAfterRenameInDatabase: true, IsRenamedDatabaseReferencedAfterRenameInSchema: true, ExpectedError: regexp.MustCompile(`Object '.*' already exists`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestAcc_ dependency: %s, is using new database name in dataabse resource: %t, in schema resource: %t", testCase.Dependency, testCase.IsRenamedDatabaseReferencedAfterRenameInDatabase, testCase.IsRenamedDatabaseReferencedAfterRenameInSchema), func(t *testing.T) {
			databaseId := testClient().Ids.RandomAccountObjectIdentifier()
			newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
			schemaName := testClient().Ids.Alpha()

			databaseConfigModel := model.DatabaseWithParametersSet("test", databaseId.Name())
			databaseConfigModelWithNewId := model.DatabaseWithParametersSet("test", newDatabaseId.Name())

			var databaseConfigAfterRename string
			var schemaConfigAfterRename string

			if testCase.IsRenamedDatabaseReferencedAfterRenameInDatabase {
				databaseConfigAfterRename = config.FromModels(t, databaseConfigModelWithNewId)
			} else {
				databaseConfigAfterRename = config.FromModels(t, databaseConfigModel)
			}

			switch testCase.Dependency {
			case ImplicitDependency:
				schemaConfigAfterRename = configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.Dependency, "", schemaName)
			case DependsOnDependency, NoDependency:
				if testCase.IsRenamedDatabaseReferencedAfterRenameInSchema {
					schemaConfigAfterRename = configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.Dependency, newDatabaseId.Name(), schemaName)
				} else {
					schemaConfigAfterRename = configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.Dependency, databaseId.Name(), schemaName)
				}
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Schema),
				Steps: []resource.TestStep{
					{
						Config: config.FromModels(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
					},
					{
						PreConfig: func() {
							testClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{
								NewName: &newDatabaseId,
							})
							t.Cleanup(testClient().Database.DropDatabaseFunc(t, newDatabaseId))
						},
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								// Creates are expected, because in refresh Read before apply the database is removing the unknown database from the state using d.SetId("") after failed ShowByID
								plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
								plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
							},
						},
						Config:      databaseConfigAfterRename + schemaConfigAfterRename,
						ExpectError: testCase.ExpectedError,
					},
				},
			})
		})
	}
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	t.Skip("Test results are inconsistent because Terraform execution order is non-deterministic")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()

	databaseConfigModel := model.DatabaseWithParametersSet("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				// This step has inconsistent results, and it depends on the Terraform execution order which seems to be non-deterministic in this case
				PreConfig: func() {
					newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
					testClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(testClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
				// ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

// This test checks if the order of the configuration resources has any impact on the order of resource execution (it seems to have no impact).
func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingOldDatabaseName_ConfigOrderSwap(t *testing.T) {
	t.Skip("Test results are inconsistent because Terraform execution order is non-deterministic")
	// Although the above applies, it seems to be consistently failing on delete operation after the test (because the database is dropped before schema).
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaName := testClient().Ids.Alpha()

	databaseConfigModel := model.DatabaseWithParametersSet("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, schemaModelConfig, databaseConfigModel),
			},
			{
				PreConfig: func() {
					newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
					testClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(testClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, schemaModelConfig, databaseConfigModel),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsNotInConfig_RenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		IsReferencingNewDatabaseName bool
		ExpectedError                *regexp.Regexp
	}{
		{IsReferencingNewDatabaseName: false, ExpectedError: regexp.MustCompile("object does not exist or not authorized")},
		{IsReferencingNewDatabaseName: true, ExpectedError: regexp.MustCompile("Failed to create schema")}, // already exists
	}

	for _, testCase := range testCases {
		database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
		t.Cleanup(databaseCleanup)

		t.Run(fmt.Sprintf("TestAcc_ referencing new database name: %t", testCase.IsReferencingNewDatabaseName), func(t *testing.T) {
			databaseId := database.ID()
			newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
			schemaName := testClient().Ids.Alpha()
			schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

			var schemaConfigAfterRename string
			if testCase.IsReferencingNewDatabaseName {
				schemaConfigAfterRename = configSchemaWithReferences(t, "", NoDependency, newDatabaseId.Name(), schemaName)
			} else {
				schemaConfigAfterRename = configSchemaWithReferences(t, "", NoDependency, databaseId.Name(), schemaName)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Schema),
				Steps: []resource.TestStep{
					{
						Config: config.FromModels(t, schemaModelConfig),
					},
					{
						PreConfig: func() {
							testClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
							t.Cleanup(testClient().Database.DropDatabaseFunc(t, newDatabaseId))
						},
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
							},
						},
						Config:      schemaConfigAfterRename,
						ExpectError: testCase.ExpectedError,
					},
				},
			})
		})
	}
}

func configSchemaWithDatabaseReference(databaseReference string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database = %[1]s.name
	name = "%[2]s"
}
`, databaseReference, schemaName)
}

func TestAcc_DeepHierarchy_AreInConfig_DatabaseRenamedInternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		DatabaseDependency         DependencyType
		SchemaDependency           DependencyType
		DatabaseInSchemaDependency DependencyType
		ExpectedFirstStepError     *regexp.Regexp
	}{
		{DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")}, // tries to create table before schema
		{DatabaseDependency: ImplicitDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency},
		{DatabaseDependency: ImplicitDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency},

		{DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: DependsOnDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")}, // tries to create table before schema
		{DatabaseDependency: ImplicitDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: DependsOnDependency},
		{DatabaseDependency: ImplicitDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: DependsOnDependency},

		// {DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: NoDependency}, // fails after incorrect execution order (tries to drop schema after database was dropped); cannot assert
		// {DatabaseDependency: ImplicitDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: NoDependency}, // tries to drop schema after database name was changed; cannot assert
		// {DatabaseDependency: ImplicitDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: NoDependency}, // tries to drop schema after database name was changed; cannot assert

		{DatabaseDependency: DependsOnDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")}, // tries to create table before schema
		{DatabaseDependency: DependsOnDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency},
		{DatabaseDependency: DependsOnDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency},

		{DatabaseDependency: NoDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")}, // tries to create table before schema
		{DatabaseDependency: NoDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency},
		{DatabaseDependency: NoDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestAcc_ database dependency: %s, schema dependency: %s, database in schema dependency: %s", testCase.DatabaseDependency, testCase.SchemaDependency, testCase.DatabaseInSchemaDependency), func(t *testing.T) {
			databaseId := testClient().Ids.RandomAccountObjectIdentifier()
			newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
			schemaName := testClient().Ids.Alpha()
			tableName := testClient().Ids.Alpha()

			databaseConfigModel := model.DatabaseWithParametersSet("test", databaseId.Name())
			databaseConfigModelWithNewId := model.DatabaseWithParametersSet("test", newDatabaseId.Name())

			testSteps := []resource.TestStep{
				{
					Config: config.FromModels(t, databaseConfigModel) +
						configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseInSchemaDependency, databaseId.Name(), schemaName) +
						configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), schemaName, tableName),
					ExpectError: testCase.ExpectedFirstStepError,
				},
			}

			if testCase.ExpectedFirstStepError == nil {
				testSteps = append(testSteps,
					resource.TestStep{
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
								plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
								plancheck.ExpectResourceAction("snowflake_table.test", plancheck.ResourceActionDestroyBeforeCreate),
							},
						},
						Config: config.FromModels(t, databaseConfigModelWithNewId) +
							configSchemaWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.DatabaseInSchemaDependency, newDatabaseId.Name(), schemaName) +
							configTableWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, newDatabaseId.Name(), schemaName, tableName),
					},
				)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Table),
				Steps:        testSteps,
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreInConfig_SchemaRenamedInternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		DatabaseDependency     DependencyType
		SchemaDependency       DependencyType
		ExpectedFirstStepError *regexp.Regexp
	}{
		{DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")}, // tries to create table before schema
		{DatabaseDependency: ImplicitDependency, SchemaDependency: ImplicitDependency},
		{DatabaseDependency: ImplicitDependency, SchemaDependency: DependsOnDependency},

		{DatabaseDependency: DependsOnDependency, SchemaDependency: NoDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")}, // tries to create table before schema
		{DatabaseDependency: DependsOnDependency, SchemaDependency: ImplicitDependency},
		{DatabaseDependency: DependsOnDependency, SchemaDependency: DependsOnDependency},

		{DatabaseDependency: NoDependency, SchemaDependency: NoDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")}, // tries to create table before schema
		{DatabaseDependency: NoDependency, SchemaDependency: ImplicitDependency},
		{DatabaseDependency: NoDependency, SchemaDependency: DependsOnDependency},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestAcc_ database dependency: %s, schema dependency: %s", testCase.DatabaseDependency, testCase.SchemaDependency), func(t *testing.T) {
			databaseId := testClient().Ids.RandomAccountObjectIdentifier()
			schemaName := testClient().Ids.Alpha()
			newSchemaName := testClient().Ids.Alpha()
			tableName := testClient().Ids.Alpha()

			databaseConfigModel := model.DatabaseWithParametersSet("test", databaseId.Name())

			testSteps := []resource.TestStep{
				{
					Config: config.FromModels(t, databaseConfigModel) +
						configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), ImplicitDependency, databaseId.Name(), schemaName) +
						configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), schemaName, tableName),
					ExpectError: testCase.ExpectedFirstStepError,
				},
			}

			if testCase.ExpectedFirstStepError == nil {
				testSteps = append(testSteps,
					resource.TestStep{
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionNoop),
								plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
								plancheck.ExpectResourceAction("snowflake_table.test", plancheck.ResourceActionDestroyBeforeCreate),
							},
						},
						Config: config.FromModels(t, databaseConfigModel) +
							configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), ImplicitDependency, databaseId.Name(), newSchemaName) +
							configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), newSchemaName, tableName),
					},
				)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Table),
				Steps:        testSteps,
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreInConfig_DatabaseRenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		DatabaseDependency         DependencyType
		SchemaDependency           DependencyType
		DatabaseInSchemaDependency DependencyType
		ExpectedFirstStepError     *regexp.Regexp
		ExpectedSecondStepError    *regexp.Regexp
	}{
		{DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: ImplicitDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a database when it's already there
		{DatabaseDependency: ImplicitDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a database when it's already there

		{DatabaseDependency: DependsOnDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: DependsOnDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a database when it's already there
		{DatabaseDependency: DependsOnDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a database when it's already there

		{DatabaseDependency: NoDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: NoDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a database when it's already there
		{DatabaseDependency: NoDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a database when it's already there
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestAcc_ database dependency: %s, schema dependency: %s, database in schema dependency: %s", testCase.DatabaseDependency, testCase.SchemaDependency, testCase.DatabaseInSchemaDependency), func(t *testing.T) {
			databaseId := testClient().Ids.RandomAccountObjectIdentifier()
			newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
			schemaName := testClient().Ids.Alpha()
			tableName := testClient().Ids.Alpha()

			databaseConfigModel := model.DatabaseWithParametersSet("test", databaseId.Name())
			databaseConfigModelWithNewId := model.DatabaseWithParametersSet("test", newDatabaseId.Name())

			testSteps := []resource.TestStep{
				{
					Config: config.FromModels(t, databaseConfigModel) +
						configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseInSchemaDependency, databaseId.Name(), schemaName) +
						configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), schemaName, tableName),
					ExpectError: testCase.ExpectedFirstStepError,
				},
			}

			if testCase.ExpectedFirstStepError == nil {
				testSteps = append(testSteps, resource.TestStep{
					PreConfig: func() {
						testClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{
							NewName: &newDatabaseId,
						})
					},
					Config: config.FromModels(t, databaseConfigModelWithNewId) +
						configSchemaWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.DatabaseInSchemaDependency, newDatabaseId.Name(), schemaName) +
						configTableWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, newDatabaseId.Name(), schemaName, tableName),
					ExpectError: testCase.ExpectedSecondStepError,
				},
				)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Table),
				Steps:        testSteps,
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreInConfig_SchemaRenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		DatabaseDependency         DependencyType
		SchemaDependency           DependencyType
		DatabaseInSchemaDependency DependencyType
		ExpectedFirstStepError     *regexp.Regexp
		ExpectedSecondStepError    *regexp.Regexp
	}{
		{DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: ImplicitDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a database when it's already there
		{DatabaseDependency: ImplicitDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a database when it's already there

		{DatabaseDependency: DependsOnDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: DependsOnDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a database when it's already there
		{DatabaseDependency: DependsOnDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a database when it's already there

		{DatabaseDependency: NoDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: NoDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a database when it's already there
		{DatabaseDependency: NoDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a database when it's already there
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestAcc_ database dependency: %s, schema dependency: %s, database in schema dependency: %s", testCase.DatabaseDependency, testCase.SchemaDependency, testCase.DatabaseInSchemaDependency), func(t *testing.T) {
			databaseId := testClient().Ids.RandomAccountObjectIdentifier()
			schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
			newSchemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
			tableName := testClient().Ids.Alpha()

			databaseConfigModel := model.DatabaseWithParametersSet("test", databaseId.Name())

			testSteps := []resource.TestStep{
				{
					Config: config.FromModels(t, databaseConfigModel) +
						configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseInSchemaDependency, databaseId.Name(), schemaId.Name()) +
						configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), schemaId.Name(), tableName),
					ExpectError: testCase.ExpectedFirstStepError,
				},
			}

			if testCase.ExpectedFirstStepError == nil {
				testSteps = append(testSteps, resource.TestStep{
					PreConfig: func() {
						testClient().Schema.Alter(t, schemaId, &sdk.AlterSchemaOptions{
							NewName: &newSchemaId,
						})
					},
					Config: config.FromModels(t, databaseConfigModel) +
						configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseInSchemaDependency, databaseId.Name(), newSchemaId.Name()) +
						configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), newSchemaId.Name(), tableName),
					ExpectError: testCase.ExpectedSecondStepError,
				},
				)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Table),
				Steps:        testSteps,
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreNotInConfig_DatabaseRenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		UseNewDatabaseNameAfterRename bool
		ExpectedSecondStepError       *regexp.Regexp
	}{
		{UseNewDatabaseNameAfterRename: true, ExpectedSecondStepError: regexp.MustCompile("already exists")},
		{UseNewDatabaseNameAfterRename: false, ExpectedSecondStepError: regexp.MustCompile("object does not exist or not authorized")},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestAcc_ use new database after rename: %t", testCase.UseNewDatabaseNameAfterRename), func(t *testing.T) {
			newDatabaseId := testClient().Ids.RandomAccountObjectIdentifier()
			tableName := testClient().Ids.Alpha()

			database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
			t.Cleanup(databaseCleanup)

			// not cleaning up, because the schema will be dropped with the database anyway
			schema, _ := testClient().Schema.CreateSchemaInDatabase(t, database.ID())

			var secondStepDatabaseName string
			if testCase.UseNewDatabaseNameAfterRename {
				secondStepDatabaseName = newDatabaseId.Name()
			} else {
				secondStepDatabaseName = database.ID().Name()
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Table),
				Steps: []resource.TestStep{
					{
						Config: configTableWithReferences(t, "", NoDependency, "", NoDependency, database.ID().Name(), schema.ID().Name(), tableName),
					},
					{
						PreConfig: func() {
							testClient().Database.Alter(t, database.ID(), &sdk.AlterDatabaseOptions{
								NewName: &newDatabaseId,
							})
							t.Cleanup(testClient().Database.DropDatabaseFunc(t, newDatabaseId))
						},
						Config:      configTableWithReferences(t, "", NoDependency, "", NoDependency, secondStepDatabaseName, schema.ID().Name(), tableName),
						ExpectError: testCase.ExpectedSecondStepError,
					},
				},
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreNotInConfig_SchemaRenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)

	testCases := []struct {
		UseNewSchemaNameAfterRename bool
		ExpectedSecondStepError     *regexp.Regexp
	}{
		{UseNewSchemaNameAfterRename: true, ExpectedSecondStepError: regexp.MustCompile("already exists")},
		{UseNewSchemaNameAfterRename: false, ExpectedSecondStepError: regexp.MustCompile("object does not exist or not authorized")},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestAcc_ use new database after rename: %t", testCase.UseNewSchemaNameAfterRename), func(t *testing.T) {
			database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
			t.Cleanup(databaseCleanup)

			newSchemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
			tableName := testClient().Ids.Alpha()

			// not cleaning up, because the schema will be dropped with the database anyway
			schema, _ := testClient().Schema.CreateSchemaInDatabase(t, database.ID())

			var secondStepSchemaName string
			if testCase.UseNewSchemaNameAfterRename {
				secondStepSchemaName = newSchemaId.Name()
			} else {
				secondStepSchemaName = schema.ID().Name()
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, resources.Table),
				Steps: []resource.TestStep{
					{
						Config: configTableWithReferences(t, "", NoDependency, "", NoDependency, database.ID().Name(), schema.ID().Name(), tableName),
					},
					{
						PreConfig: func() {
							testClient().Schema.Alter(t, schema.ID(), &sdk.AlterSchemaOptions{
								NewName: &newSchemaId,
							})
						},
						Config:      configTableWithReferences(t, "", NoDependency, "", NoDependency, database.ID().Name(), secondStepSchemaName, tableName),
						ExpectError: testCase.ExpectedSecondStepError,
					},
				},
			})
		})
	}
}

func configSchemaWithReferences(t *testing.T, databaseReference string, databaseDependencyType DependencyType, databaseName string, schemaName string) string {
	t.Helper()
	switch databaseDependencyType {
	case ImplicitDependency:
		return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database = %[1]s.name
	name = "%[2]s"
}
`, databaseReference, schemaName)
	case DependsOnDependency:
		return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	depends_on = [%[1]s]
	database = "%[2]s"
	name = "%[3]s"
}
`, databaseReference, databaseName, schemaName)
	case NoDependency:
		return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database = "%[1]s"
	name = "%[2]s"
}
`, databaseName, schemaName)
	default:
		t.Fatalf("configSchemaWithReferences: unknown database reference type: %s", databaseDependencyType)
		return ""
	}
}

func configTableWithReferences(t *testing.T, databaseReference string, databaseDependencyType DependencyType, schemaReference string, schemaDependencyType DependencyType, databaseName string, schemaName string, tableName string) string {
	t.Helper()
	builder := new(strings.Builder)
	builder.WriteString("resource \"snowflake_table\" \"test\" {\n")

	dependsOn := make([]string, 0)
	database := ""
	schema := ""

	switch databaseDependencyType {
	case ImplicitDependency:
		database = fmt.Sprintf("%s.name", databaseReference)
	case DependsOnDependency:
		dependsOn = append(dependsOn, databaseReference)
		database = strconv.Quote(databaseName)
	case NoDependency:
		database = strconv.Quote(databaseName)
	}

	switch schemaDependencyType {
	case ImplicitDependency:
		schema = fmt.Sprintf("%s.name", schemaReference)
	case DependsOnDependency:
		dependsOn = append(dependsOn, schemaReference)
		schema = strconv.Quote(schemaName)
	case NoDependency:
		schema = strconv.Quote(schemaName)
	}

	if len(dependsOn) > 0 {
		builder.WriteString(fmt.Sprintf("depends_on = [%s]\n", strings.Join(dependsOn, ", ")))
	}
	builder.WriteString(fmt.Sprintf("database = %s\n", database))
	builder.WriteString(fmt.Sprintf("schema = %s\n", schema))
	builder.WriteString(fmt.Sprintf("name = \"%s\"\n", tableName))
	builder.WriteString(`
column {
	type = "NUMBER(38,0)"
	name = "N"
}
`)
	builder.WriteString(`}`)
	return builder.String()
}
