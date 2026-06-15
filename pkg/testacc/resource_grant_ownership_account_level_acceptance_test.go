//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantOwnership_OnTask_Discussion2877(t *testing.T) {
	taskId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleId.Name()),
		"database":          config.StringVariable(taskId.DatabaseName()),
		"schema":            config.StringVariable(taskId.SchemaName()),
		"task":              config.StringVariable(taskId.Name()),
		"child":             config.StringVariable(childId.Name()),
		"warehouse":         config.StringVariable(TestWarehouseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask_Discussion2877/1"),
				ConfigVariables: configVariables,
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask_Discussion2877/2"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("cannot have the given predecessor since they do not share the same owner role"),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask_Discussion2877/3"),
				ConfigVariables: configVariables,
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask_Discussion2877/4"),
				ConfigVariables: configVariables,
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
