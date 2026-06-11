//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantAccountRole_accountRole(t *testing.T) {
	roleId := testClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	parentRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	parentRoleName := parentRoleId.Name()

	roleModel := model.AccountRole("role", roleName)
	parentRoleModel := model.AccountRole("parent_role", parentRoleName)
	grantModel := model.GrantAccountRole("g", roleName).
		WithParentRoleName(parentRoleName).
		WithDependsOn(roleModel.ResourceReference(), parentRoleModel.ResourceReference())

	resourceName := "snowflake_grant_account_role.g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, roleModel, parentRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", parentRoleName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"|ROLE|"%v"`, roleName, parentRoleName)),
				),
			},
			// import
			{
				Config:            accconfig.FromModels(t, roleModel, parentRoleModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantAccountRole_user(t *testing.T) {
	roleId := testClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	userName := userId.Name()

	roleModel := model.AccountRole("role", roleName)
	userModel := model.User("user", userName)
	grantModel := model.GrantAccountRole("g", roleName).
		WithUserName(userName).
		WithDependsOn(roleModel.ResourceReference(), userModel.ResourceReference())

	resourceName := "snowflake_grant_account_role.g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, roleModel, userModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"|USER|"%v"`, roleName, userName)),
				),
			},
			// import
			{
				Config:            accconfig.FromModels(t, roleModel, userModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantAccountRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	roleId := testClient().Ids.RandomAccountObjectIdentifier()
	parentRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	roleModel := model.AccountRole("role", roleId.Name())
	parentRoleModel := model.AccountRole("parent_role", parentRoleId.Name())
	grantModel := model.GrantAccountRole("test", roleId.Name()).
		WithParentRoleName(parentRoleId.Name()).
		WithDependsOn(roleModel.ResourceReference(), parentRoleModel.ResourceReference())

	resourceConfig := accconfig.FromModels(t, roleModel, parentRoleModel, grantModel)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            providerConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", fmt.Sprintf(`%v|ROLE|%v`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   resourceConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_account_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_account_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", fmt.Sprintf(`%v|ROLE|%v`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
		},
	})
}

func TestAcc_GrantAccountRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	roleId := testClient().Ids.RandomAccountObjectIdentifier()
	parentRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				ExpectError:       regexp.MustCompile("Error: Provider produced inconsistent final plan"),
				Config:            providerConfig + grantAccountRoleConfigWithQuotedIdentifiers(roleId, parentRoleId),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantAccountRoleConfigWithQuotedIdentifiers(roleId, parentRoleId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_account_role.test", plancheck.ResourceActionCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_account_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "role_name", roleId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "parent_role_name", parentRoleId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", fmt.Sprintf(`%v|ROLE|%v`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantAccountRoleConfigWithQuotedIdentifiers(roleId sdk.AccountObjectIdentifier, parentRoleId sdk.AccountObjectIdentifier) string {
	quotedRoleId := fmt.Sprintf(`\"%s\"`, roleId.Name())
	quotedParentRoleId := fmt.Sprintf(`\"%s\"`, parentRoleId.Name())

	return fmt.Sprintf(`
resource "snowflake_account_role" "role" {
  name = "%s"
}

resource "snowflake_account_role" "parent_role" {
  name = "%s"
}

resource "snowflake_grant_account_role" "test" {
  role_name        = snowflake_account_role.role.name
  parent_role_name = snowflake_account_role.parent_role.name
}
`, quotedRoleId, quotedParentRoleId)
}

// proves that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3629 (UBAC) doesn't affect the grant account role resource
func TestAcc_GrantAccountRole_Issue_3629(t *testing.T) {
	accountRole, accountRoleCleanup := secondaryTestClient().Role.CreateRole(t)
	t.Cleanup(accountRoleCleanup)

	parentAccountRole, parentAccountRoleCleanup := secondaryTestClient().Role.CreateRole(t)
	t.Cleanup(parentAccountRoleCleanup)

	user, userCleanup := secondaryTestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	grantModel := model.GrantAccountRole("test", accountRole.ID().Name()).
		WithParentRoleName(parentAccountRole.ID().Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: secondaryAccountProviderFactory,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					secondaryTestClient().Role.GrantRoleToUser(t, accountRole.ID(), user.ID())
				},
				Config: accconfig.FromModels(t, providerModel, grantModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", helpers.EncodeResourceIdentifier(accountRole.ID().FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentAccountRole.ID().FullyQualifiedName()))),
				),
			},
		},
	})
}
