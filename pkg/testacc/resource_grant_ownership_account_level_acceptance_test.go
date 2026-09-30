//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantOwnership_OnTask_Discussion2877(t *testing.T) {
	taskId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	parentTaskModel := model.TaskWithId("test", taskId, false, "SELECT CURRENT_TIMESTAMP").
		WithWarehouse(TestWarehouseName)
	childTaskModel := model.TaskWithId("child", childId, false, "SELECT CURRENT_TIMESTAMP").
		WithWarehouse(TestWarehouseName).
		WithAfterValue(tfconfig.SetVariable(tfconfig.StringVariable(taskId.FullyQualifiedName())))

	grantOnParentModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeTask, taskId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), parentTaskModel.ResourceReference())
	grantOnParentWithChildDependencyModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeTask, taskId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), childTaskModel.ResourceReference())
	grantOnChildModel := model.GrantOwnershipWithRawOn("child").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeTask, childId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), childTaskModel.ResourceReference())
	grantOnAllTasksModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnAllInSchema(sdk.PluralObjectTypeTasks, taskId.SchemaId().FullyQualifiedName()).
		WithDependsOn(parentTaskModel.ResourceReference(), childTaskModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, parentTaskModel, grantOnParentModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test", "name", taskId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TASK|%s", accountRoleId.FullyQualifiedName(), taskId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
						ObjectType: sdk.ObjectTypeTask,
						Name:       taskId,
					})), sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName()),
				),
			},
			{
				Config:      accconfig.FromModels(t, accountRoleModel, parentTaskModel, childTaskModel, grantOnParentWithChildDependencyModel, grantOnChildModel),
				ExpectError: regexp.MustCompile("cannot have the given predecessor since they do not share the same owner role"),
			},
			{
				Config: accconfig.FromModels(t, accountRoleModel, parentTaskModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test", "name", taskId.Name()),
					checkResourceOwnershipIsGranted(sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
						ObjectType: sdk.ObjectTypeTask,
						Name:       taskId,
					})), sdk.ObjectTypeTask, testClient().Context.CurrentRole(t).Name(), taskId.FullyQualifiedName()),
				),
			},
			{
				Config: accconfig.FromModels(t, accountRoleModel, parentTaskModel, childTaskModel, grantOnAllTasksModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test", "name", taskId.Name()),
					resource.TestCheckResourceAttr("snowflake_task.child", "name", childId.Name()),
					resource.TestCheckResourceAttr("snowflake_task.child", "after.0", taskId.FullyQualifiedName()),
					checkResourceOwnershipIsGranted(sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
						ObjectType: sdk.ObjectTypeTask,
						Name:       taskId,
					})), sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName()),
					checkResourceOwnershipIsGranted(sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
						ObjectType: sdk.ObjectTypeTask,
						Name:       childId,
					})), sdk.ObjectTypeTask, accountRoleId.Name(), childId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_bcr2026_06_OnDatabaseRole(t *testing.T) {
	// Regression test for SNOW-4075337: the 2026_06 bundle (BCR-2371) makes SHOW GRANTS ON DATABASE ROLE
	// report granted_on = DATABASE_ROLE instead of ROLE. Before the fix, ReadGrantOwnership matched the
	// OWNERSHIP grant only when granted_on was ROLE, so with the bundle enabled it cleared the id and apply
	// failed with an "inconsistent result after apply" error. The fixed provider accepts both values.
	// TODO(SNOW-4075337): skip or remove this test once the 2026_06 bundle is enforced and can no longer be disabled.

	// Toggling a behavior change bundle is account-wide, so run on the secondary account.
	database, databaseCleanup := secondaryTestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseRoleId := secondaryTestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	accountRoleId := secondaryTestClient().Ids.RandomAccountObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	dbRoleModel := model.DatabaseRole("test", databaseId.Name(), databaseRoleId.Name())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeDatabaseRole, databaseRoleId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), dbRoleModel.ResourceReference())
	expectedId := fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE ROLE|%s", accountRoleId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Latest released provider, bundle disabled: granted_on is ROLE and the grant matches on read.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.20.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.DisableBcrBundle(t, "2026_06")
				},
				Config: accconfig.FromModels(t, providerModel, accountRoleModel, dbRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", expectedId),
				),
			},
			// Same released provider, bundle enabled: granted_on becomes DATABASE_ROLE, read no longer matches,
			// the id is cleared and apply fails with an inconsistent-result error.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.20.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2026_06")
				},
				Config:      accconfig.FromModels(t, providerModel, accountRoleModel, dbRoleModel, grantModel),
				ExpectError: regexp.MustCompile("Provider produced inconsistent result after apply"),
			},
			// Fixed provider, bundle still enabled: DATABASE_ROLE is accepted, the grant matches on read and the
			// resource is stable. secondaryAccountProviderFactory is required because TestAccProtoV6ProviderFactories
			// is cached against the primary account.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   accconfig.FromModels(t, providerModel, accountRoleModel, dbRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", expectedId),
				),
			},
		},
	})
}
