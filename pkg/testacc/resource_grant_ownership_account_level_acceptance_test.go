//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
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
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       taskId,
							},
						},
					}, sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName()),
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
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       taskId,
							},
						},
					}, sdk.ObjectTypeTask, testClient().Context.CurrentRole(t).Name(), taskId.FullyQualifiedName()),
				),
			},
			{
				Config: accconfig.FromModels(t, accountRoleModel, parentTaskModel, childTaskModel, grantOnAllTasksModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test", "name", taskId.Name()),
					resource.TestCheckResourceAttr("snowflake_task.child", "name", childId.Name()),
					resource.TestCheckResourceAttr("snowflake_task.child", "after.0", taskId.FullyQualifiedName()),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       taskId,
							},
						},
					}, sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName()),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       childId,
							},
						},
					}, sdk.ObjectTypeTask, accountRoleId.Name(), childId.FullyQualifiedName()),
				),
			},
		},
	})
}
