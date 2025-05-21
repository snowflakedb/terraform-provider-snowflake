//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantAccountRole_accountRole(t *testing.T) {
	roleId := testClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	parentRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	parentRoleName := parentRoleId.Name()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"role_name":        config.StringVariable(roleName),
			"parent_role_name": config.StringVariable(parentRoleName),
		}
	}

	resourceName := "snowflake_grant_account_role.g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantAccountRole/account_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", parentRoleName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"|ROLE|"%v"`, roleName, parentRoleName)),
				),
			},
			// import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantAccountRole/account_role"),
				ConfigVariables:   m(),
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

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"role_name": config.StringVariable(roleName),
			"user_name": config.StringVariable(userName),
		}
	}

	resourceName := "snowflake_grant_account_role.g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantAccountRole/user"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"|USER|"%v"`, roleName, userName)),
				),
			},
			// import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantAccountRole/user"),
				ConfigVariables:   m(),
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

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantAccountRoleBasicConfig(roleId, parentRoleId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", fmt.Sprintf(`%v|ROLE|%v`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantAccountRoleBasicConfig(roleId, parentRoleId),
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

func grantAccountRoleBasicConfig(roleId sdk.AccountObjectIdentifier, parentRoleId sdk.AccountObjectIdentifier) string {
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
`, roleId.Name(), parentRoleId.Name())
}

func TestAcc_GrantAccountRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	roleId := testClient().Ids.RandomAccountObjectIdentifier()
	parentRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				ExpectError:       regexp.MustCompile("Error: Provider produced inconsistent final plan"),
				Config:            grantAccountRoleConfigWithQuotedIdentifiers(roleId, parentRoleId),
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
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	accountRole, accountRoleCleanup := secondaryTestClient().Role.CreateRole(t)
	t.Cleanup(accountRoleCleanup)

	parentAccountRole, parentAccountRoleCleanup := secondaryTestClient().Role.CreateRole(t)
	t.Cleanup(parentAccountRoleCleanup)

	user, userCleanup := secondaryTestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	testConfig := accconfig.FromModels(t, providerModel) + grantAccountRoleIssue3629Config(accountRole.ID(), parentAccountRole.ID())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					secondaryTestClient().Role.GrantRoleToUser(t, accountRole.ID(), user.ID())
				},
				Config: testConfig,
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", helpers.EncodeResourceIdentifier(accountRole.ID().FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentAccountRole.ID().FullyQualifiedName()))),
				),
			},
		},
	})
}

func grantAccountRoleIssue3629Config(accountRoleId sdk.AccountObjectIdentifier, parentRoleId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_grant_account_role" "test" {
  role_name = "%[1]s"
  parent_role_name = "%[2]s"
}
`, accountRoleId.Name(), parentRoleId.Name())
}
