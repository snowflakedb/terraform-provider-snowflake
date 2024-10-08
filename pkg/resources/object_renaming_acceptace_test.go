package resources_test

import (
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

/*
The following tests are showing the behavior of the provider in cases where objects higher in the hierarchy
like database or schema are renamed when the objects lower in the hierarchy are in the Terraform configuration.
For more information check TODO(SNOW-1672319): link public document.

Shallow hierarchy (database + schema)
- is in config - renamed internally - with implicit dependency

- is in config - renamed internally - without dependency - after rename schema referencing old database name
- is in config - renamed internally - without dependency - after rename schema referencing new database name

- is in config - renamed internally - with depends_on - after rename schema referencing old database name
- is in config - renamed internally - with depends_on - after rename schema referencing new database name

- is in config - renamed externally - with implicit dependency - database holding the same name in config
- is in config - renamed externally - with implicit dependency - database holding the new name in config

- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing old database name
- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing old database name - check impact of the resource order on the plan
- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing new database name
- is in config - renamed externally - without dependency - after rename database referencing new name and schema referencing old database name
- is in config - renamed externally - without dependency - after rename database referencing new name and schema referencing new database name

- is in config - renamed externally - with depends_on - after rename database referencing old name and schema referencing old database name
- is in config - renamed externally - with depends_on - after rename database referencing old name and schema referencing new database name
- is in config - renamed externally - with depends_on - after rename database referencing new name and schema referencing old database name
- is in config - renamed externally - with depends_on - after rename database referencing new name and schema referencing new database name

- is not in config - renamed externally - referencing old database name
- is not in config - renamed externally - referencing new database name

Deep hierarchy (database + schema + table)
- TODO: More test cases could be added where there's no dependency, because there could be either old or new database/schema name referenced in the table resource config

- are in config - database renamed internally - with database implicit dependency - with no schema dependency 		- with database to schema implicit dependency
- are in config - database renamed internally - with database implicit dependency - with implicit schema dependency - with database to schema implicit dependency
- are in config - database renamed internally - with database implicit dependency - with schema depends_on 			- with database to schema implicit dependency

- are in config - database renamed internally - with database implicit dependency - with no schema dependency 		- with database to schema depends_on dependency
- are in config - database renamed internally - with database implicit dependency - with implicit schema dependency - with database to schema depends_on dependency
- are in config - database renamed internally - with database implicit dependency - with schema depends_on 			- with database to schema depends_on dependency

- are in config - database renamed internally - with database implicit dependency - with no schema dependency 		- with database to schema no dependency
- are in config - database renamed internally - with database implicit dependency - with implicit schema dependency - with database to schema no dependency
- are in config - database renamed internally - with database implicit dependency - with schema depends_on 			- with database to schema no dependency

- are in config - database renamed internally - with no database dependency - with no schema dependency
- are in config - database renamed internally - with no database dependency - with implicit schema dependency
- are in config - database renamed internally - with no database dependency - with schema depends_on

- are in config - database renamed internally - with database depends_on - with no schema dependency
- are in config - database renamed internally - with database depends_on - with implicit schema dependency
- are in config - database renamed internally - with database depends_on - with schema depends_on

------------------------------------------------------------------------------------------------------------------------

- are in config - schema renamed internally - with database implicit dependency - with no schema dependency
- are in config - schema renamed internally - with database implicit dependency - with implicit schema dependency
- are in config - schema renamed internally - with database implicit dependency - with schema depends_on

- are in config - schema renamed internally - with no database dependency - with no schema dependency
- are in config - schema renamed internally - with no database dependency - with implicit schema dependency
- are in config - schema renamed internally - with no database dependency - with schema depends_on

- are in config - schema renamed internally - with database depends_on - with no schema dependency
- are in config - schema renamed internally - with database depends_on - with implicit schema dependency
- are in config - schema renamed internally - with database depends_on - with schema depends_on

------------------------------------------------------------------------------------------------------------------------

- are in config - database renamed externally - with database implicit dependency - with no schema dependency
- are in config - database renamed externally - with database implicit dependency - with implicit schema dependency
- are in config - database renamed externally - with database implicit dependency - with schema depends_on

- are in config - database renamed externally - with no database dependency - with no schema dependency
- are in config - database renamed externally - with no database dependency - with implicit schema dependency
- are in config - database renamed externally - with no database dependency - with schema depends_on

- are in config - database renamed externally - with database depends_on - with no schema dependency
- are in config - database renamed externally - with database depends_on - with implicit schema dependency
- are in config - database renamed externally - with database depends_on - with schema depends_on

------------------------------------------------------------------------------------------------------------------------

- are in config - schema renamed externally - with database implicit dependency - with no schema dependency
- are in config - schema renamed externally - with database implicit dependency - with implicit schema dependency
- are in config - schema renamed externally - with database implicit dependency - with schema depends_on

- are in config - schema renamed externally - with no database dependency - with no schema dependency
- are in config - schema renamed externally - with no database dependency - with implicit schema dependency
- are in config - schema renamed externally - with no database dependency - with schema depends_on

- are in config - schema renamed externally - with database depends_on - with no schema dependency
- are in config - schema renamed externally - with database depends_on - with implicit schema dependency
- are in config - schema renamed externally - with database depends_on - with schema depends_on

- TODO: More test cases could be added when database and schema are renamed at the same time

------------------------------------------------------------------------------------------------------------------------

- are not in config - database renamed externally - referencing old database name
- are not in config - database renamed externally - referencing new database name

- are not in config - schema renamed externally - referencing old schema name
- are not in config - schema renamed externally - referencing new schema name

- TODO: More test cases could be added when either database or schema are in the config
*/

type DependencyType string

const (
	ImplicitDependency  DependencyType = "implicit"
	DependsOnDependency DependencyType = "depends_on"
	NoDependency        DependencyType = "no_dependency"
)

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithImplicitDependency(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModel(t, databaseConfigModelWithNewId) + configSchemaWithDatabaseReference(databaseConfigModelWithNewId.ResourceReference(), schemaName),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithoutDependency_AfterRenameSchemaReferencingOldDatabaseName(t *testing.T) {
	// Error happens during schema's Read operation and then Delete operation (schema cannot be removed).
	t.Skip("Not able to handle the error produced by Delete operation that results in test always failing")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfig),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithoutDependency_AfterRenameSchemaReferencingNewDatabaseName(t *testing.T) {
	// Error happens during schema's Read operation and then Delete operation (schema cannot be removed).
	t.Skip("Not able to handle the error produced by Delete operation that results in test always failing")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithDependsOn_AfterRenameSchemaReferencingOldDatabaseName(t *testing.T) {
	// Error happens during schema's Read operation and then Delete operation (schema cannot be removed).
	t.Skip("Not able to handle the error produced by Delete operation that results in test always failing")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfig),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithDependsOn_AfterRenameSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId.SetDependsOn(databaseConfigModelWithNewId.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithImplicitDependency_DatabaseHoldingTheOldNameInConfig(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{
						NewName: &newDatabaseId,
					})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithImplicitDependency_DatabaseHoldingTheNewNameInConfig(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{
						NewName: &newDatabaseId,
					})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// TODO: Why Create? This case could be handled in a way that it shouldn't require any action (implicit import)
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModel(t, databaseConfigModelWithNewId) + configSchemaWithDatabaseReference(databaseConfigModelWithNewId.ResourceReference(), schemaName),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	t.Skip("Test results are inconsistent because Terraform execution order is non-deterministic")
	// Although the above applies, it seems to be consistently failing on delete operation after the test (because the database is dropped before schema).
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				PreConfig: func() {
					newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				// TODO: This test may have inconsistent result (depending on Terraform execution order which is not deterministic).
				Config:      config.FromModels(t, databaseConfigModel, schemaModelConfig),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

// This test checks if the order of the configuration resources has any impact on the order of resource execution (it seems to have no impact).
func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingOldDatabaseName_ConfigOrderSwap(t *testing.T) {
	t.Skip("Test results are inconsistent because Terraform execution order is non-deterministic")
	// Although the above applies, it seems to be consistently failing on delete operation after the test (because the database is dropped before schema).
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, schemaModelConfig, databaseConfigModel),
			},
			{
				PreConfig: func() {
					newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				// TODO: This test may have inconsistent result (depending on Terraform execution order which is not deterministic).
				Config:      config.FromModels(t, schemaModelConfig, databaseConfigModel),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModel, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile("Failed to create schema"), // already exists
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingNewNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfig),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingNewNameAndSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithDependsOn_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				PreConfig: func() {
					newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithDependsOn_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId.SetDependsOn(databaseConfigModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModel, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile("Failed to create schema"), // already exists
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithDependsOn_AfterRenameDatabaseReferencingNewNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfig),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithDependsOn_AfterRenameDatabaseReferencingNewNameAndSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId.SetDependsOn(databaseConfigModelWithNewId.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsNotInConfig_RenamedExternally_ReferencingOldName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithIdentifier(t, databaseId)
					t.Cleanup(databaseCleanup)
				},
				Config: config.FromModel(t, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModel(t, schemaModelConfig),
				ExpectError: regexp.MustCompile("object does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsNotInConfig_RenamedExternally_ReferencingNewName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithIdentifier(t, databaseId)
					t.Cleanup(databaseCleanup)
				},
				Config: config.FromModel(t, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModel(t, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile("Failed to create schema"), // already exists
			},
		},
	})
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

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

		//{DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: NoDependency}, // fails after incorrect execution order (tries to drop schema after database was dropped); cannot assert
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
		t.Run(fmt.Sprintf("database dependency: %s, schema dependency: %s, database in schema dependency: %s", testCase.DatabaseDependency, testCase.SchemaDependency, testCase.DatabaseInSchemaDependency), func(t *testing.T) {
			databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
			newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
			schemaName := acc.TestClient().Ids.Alpha()
			tableName := acc.TestClient().Ids.Alpha()

			databaseConfigModel := model.Database("test", databaseId.Name())
			databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

			testSteps := []resource.TestStep{
				{
					Config: config.FromModel(t, databaseConfigModel) +
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
						Config: config.FromModel(t, databaseConfigModelWithNewId) +
							configSchemaWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.DatabaseInSchemaDependency, newDatabaseId.Name(), schemaName) +
							configTableWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, newDatabaseId.Name(), schemaName, tableName),
					},
				)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: acc.CheckDestroy(t, resources.Table),
				Steps:        testSteps,
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreInConfig_SchemaRenamedInternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

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
		t.Run(fmt.Sprintf("database dependency: %s, schema dependency: %s", testCase.DatabaseDependency, testCase.SchemaDependency), func(t *testing.T) {
			databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
			schemaName := acc.TestClient().Ids.Alpha()
			newSchemaName := acc.TestClient().Ids.Alpha()
			tableName := acc.TestClient().Ids.Alpha()

			databaseConfigModel := model.Database("test", databaseId.Name())

			testSteps := []resource.TestStep{
				{
					Config: config.FromModel(t, databaseConfigModel) +
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
						Config: config.FromModel(t, databaseConfigModel) +
							configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), ImplicitDependency, databaseId.Name(), newSchemaName) +
							configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), newSchemaName, tableName),
					},
				)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: acc.CheckDestroy(t, resources.Table),
				Steps:        testSteps,
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreInConfig_DatabaseRenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	testCases := []struct {
		DatabaseDependency         DependencyType
		SchemaDependency           DependencyType
		DatabaseInSchemaDependency DependencyType
		ExpectedFirstStepError     *regexp.Regexp
		ExpectedSecondStepError    *regexp.Regexp
	}{
		{DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: ImplicitDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a table when it's already there
		{DatabaseDependency: ImplicitDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a table when it's already there

		{DatabaseDependency: DependsOnDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: DependsOnDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a table when it's already there
		{DatabaseDependency: DependsOnDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a table when it's already there

		{DatabaseDependency: NoDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: NoDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a table when it's already there
		{DatabaseDependency: NoDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a table when it's already there
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("database dependency: %s, schema dependency: %s, database in schema dependency: %s", testCase.DatabaseDependency, testCase.SchemaDependency, testCase.DatabaseInSchemaDependency), func(t *testing.T) {
			databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
			newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
			schemaName := acc.TestClient().Ids.Alpha()
			tableName := acc.TestClient().Ids.Alpha()

			databaseConfigModel := model.Database("test", databaseId.Name())
			databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

			testSteps := []resource.TestStep{
				{
					Config: config.FromModel(t, databaseConfigModel) +
						configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseInSchemaDependency, databaseId.Name(), schemaName) +
						configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), schemaName, tableName),
					ExpectError: testCase.ExpectedFirstStepError,
				},
			}

			if testCase.ExpectedFirstStepError == nil {
				testSteps = append(testSteps, resource.TestStep{
					PreConfig: func() {
						acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{
							NewName: &newDatabaseId,
						})
					},
					Config: config.FromModel(t, databaseConfigModelWithNewId) +
						configSchemaWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.DatabaseInSchemaDependency, newDatabaseId.Name(), schemaName) +
						configTableWithReferences(t, databaseConfigModelWithNewId.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, newDatabaseId.Name(), schemaName, tableName),
					ExpectError: testCase.ExpectedSecondStepError,
				},
				)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: acc.CheckDestroy(t, resources.Table),
				Steps:        testSteps,
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreInConfig_SchemaRenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	testCases := []struct {
		DatabaseDependency         DependencyType
		SchemaDependency           DependencyType
		DatabaseInSchemaDependency DependencyType
		ExpectedFirstStepError     *regexp.Regexp
		ExpectedSecondStepError    *regexp.Regexp
	}{
		{DatabaseDependency: ImplicitDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: ImplicitDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a table when it's already there
		{DatabaseDependency: ImplicitDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a table when it's already there

		{DatabaseDependency: DependsOnDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: DependsOnDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a table when it's already there
		{DatabaseDependency: DependsOnDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a table when it's already there

		{DatabaseDependency: NoDependency, SchemaDependency: NoDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedFirstStepError: regexp.MustCompile("error creating table")},   // tries to create table before schema
		{DatabaseDependency: NoDependency, SchemaDependency: ImplicitDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")},  // tries to create a table when it's already there
		{DatabaseDependency: NoDependency, SchemaDependency: DependsOnDependency, DatabaseInSchemaDependency: ImplicitDependency, ExpectedSecondStepError: regexp.MustCompile("already exists")}, // tries to create a table when it's already there
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("database dependency: %s, schema dependency: %s, database in schema dependency: %s", testCase.DatabaseDependency, testCase.SchemaDependency, testCase.DatabaseInSchemaDependency), func(t *testing.T) {
			databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
			schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
			newSchemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
			tableName := acc.TestClient().Ids.Alpha()

			databaseConfigModel := model.Database("test", databaseId.Name())

			testSteps := []resource.TestStep{
				{
					Config: config.FromModel(t, databaseConfigModel) +
						configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseInSchemaDependency, databaseId.Name(), schemaId.Name()) +
						configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), schemaId.Name(), tableName),
					ExpectError: testCase.ExpectedFirstStepError,
				},
			}

			if testCase.ExpectedFirstStepError == nil {
				testSteps = append(testSteps, resource.TestStep{
					PreConfig: func() {
						acc.TestClient().Schema.Alter(t, schemaId, &sdk.AlterSchemaOptions{
							NewName: &newSchemaId,
						})
					},
					Config: config.FromModel(t, databaseConfigModel) +
						configSchemaWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseInSchemaDependency, databaseId.Name(), newSchemaId.Name()) +
						configTableWithReferences(t, databaseConfigModel.ResourceReference(), testCase.DatabaseDependency, "snowflake_schema.test", testCase.SchemaDependency, databaseId.Name(), newSchemaId.Name(), tableName),
					ExpectError: testCase.ExpectedSecondStepError,
				},
				)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: acc.CheckDestroy(t, resources.Table),
				Steps:        testSteps,
			})
		})
	}
}

func TestAcc_DeepHierarchy_AreNotInConfig_DatabaseRenamedExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	testCases := []struct {
		UseNewDatabaseNameAfterRename bool
		ExpectedSecondStepError       *regexp.Regexp
	}{
		{UseNewDatabaseNameAfterRename: true, ExpectedSecondStepError: regexp.MustCompile("already exists")},
		{UseNewDatabaseNameAfterRename: false, ExpectedSecondStepError: regexp.MustCompile("object does not exist or not authorized")},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("use new database after rename: %t", testCase.UseNewDatabaseNameAfterRename), func(t *testing.T) {
			newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
			tableName := acc.TestClient().Ids.Alpha()

			database, databaseCleanup := acc.TestClient().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			// not cleaning up, because the schema will be dropped with the database anyway
			schema, _ := acc.TestClient().Schema.CreateSchemaInDatabase(t, database.ID())

			var secondStepDatabaseName string
			if testCase.UseNewDatabaseNameAfterRename {
				secondStepDatabaseName = newDatabaseId.Name()
			} else {
				secondStepDatabaseName = database.ID().Name()
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: acc.CheckDestroy(t, resources.Table),
				Steps: []resource.TestStep{
					{
						Config: configTableWithReferences(t, "", NoDependency, "", NoDependency, database.ID().Name(), schema.ID().Name(), tableName),
					},
					{
						PreConfig: func() {
							acc.TestClient().Database.Alter(t, database.ID(), &sdk.AlterDatabaseOptions{
								NewName: &newDatabaseId,
							})
							t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	testCases := []struct {
		UseNewSchemaNameAfterRename bool
		ExpectedSecondStepError     *regexp.Regexp
	}{
		{UseNewSchemaNameAfterRename: true, ExpectedSecondStepError: regexp.MustCompile("already exists")},
		{UseNewSchemaNameAfterRename: false, ExpectedSecondStepError: regexp.MustCompile("object does not exist or not authorized")},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("use new database after rename: %t", testCase.UseNewSchemaNameAfterRename), func(t *testing.T) {
			database, databaseCleanup := acc.TestClient().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			newSchemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
			tableName := acc.TestClient().Ids.Alpha()

			// not cleaning up, because the schema will be dropped with the database anyway
			schema, _ := acc.TestClient().Schema.CreateSchemaInDatabase(t, database.ID())

			var secondStepSchemaName string
			if testCase.UseNewSchemaNameAfterRename {
				secondStepSchemaName = newSchemaId.Name()
			} else {
				secondStepSchemaName = schema.ID().Name()
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: acc.CheckDestroy(t, resources.Table),
				Steps: []resource.TestStep{
					{
						Config: configTableWithReferences(t, "", NoDependency, "", NoDependency, database.ID().Name(), schema.ID().Name(), tableName),
					},
					{
						PreConfig: func() {
							acc.TestClient().Schema.Alter(t, schema.ID(), &sdk.AlterSchemaOptions{
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
