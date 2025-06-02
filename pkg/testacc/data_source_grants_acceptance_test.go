//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Grants_On_Account(t *testing.T) {
	grantsModel := datasourcemodel.GrantsOnAccount("test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantsModel),
				Check:  checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_AccountObject(t *testing.T) {
	grantsModel := datasourcemodel.GrantsOnAccountObject("test", testClient().Ids.DatabaseId(), sdk.ObjectTypeDatabase)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantsModel),
				Check:  checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_DatabaseObject(t *testing.T) {
	grantsModel := datasourcemodel.GrantsOnDatabaseObject("test", testClient().Ids.SchemaId(), sdk.ObjectTypeSchema)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantsModel),
				Check:  checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_SchemaObject(t *testing.T) {
	viewId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	columnNames := []string{"ROLE_NAME"}

	viewModel := model.View("test", viewId.DatabaseName(), viewId.SchemaName(), viewId.Name(), statement).WithColumnNames(columnNames...)
	grantsModel := datasourcemodel.GrantsOnSchemaObject("test", viewId, sdk.ObjectTypeView).
		WithDependsOn(viewModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, viewModel, grantsModel),
				Check:  checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_SchemaObject_WithArguments(t *testing.T) {
	function := testClient().Function.Create(t, sdk.DataTypeVARCHAR)
	grantsModel := datasourcemodel.GrantsOnSchemaObjectWithArguments("test", function.ID(), sdk.ObjectTypeFunction)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantsModel),
				Check:  checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_Invalid_NoAttribute(t *testing.T) {
	grantsModel := datasourcemodel.GrantsOnEmpty("test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, grantsModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_On_Invalid_MissingObjectType(t *testing.T) {
	grantsModel := datasourcemodel.GrantsOnMissingObjectType("test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, grantsModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}

// TODO [SNOW-1284382]: Implement after snowflake_application and snowflake_application_role resources are introduced.
func TestAcc_Grants_To_Application(t *testing.T) {
	t.Skip("Skipped until snowflake_application and snowflake_application_role resources are introduced. Currently, behavior tested in application_roles_gen_integration_test.go.")
}

// TODO [SNOW-1284382]: Implement after snowflake_application and snowflake_application_role resources are introduced.
func TestAcc_Grants_To_ApplicationRole(t *testing.T) {
	t.Skip("Skipped until snowflake_application and snowflake_application_role resources are introduced. Currently, behavior tested in application_roles_gen_integration_test.go.")
}

func TestAcc_Grants_To_AccountRole(t *testing.T) {
	// TODO [SNOW-1887460]: handle SNOWFLAKE.CORE."AVG(ARG_T TABLE(FLOAT)):FLOAT" and SNOWFLAKE.ACCOUNT_USAGE."TAG_REFERENCES_WITH_LINEAGE(TAG_NAME_INPUT VARCHAR):TABLE:
	t.Skip(`Skipped temporarily because incompatible data types on the current role: SNOWFLAKE.CORE."AVG(ARG_T TABLE(FLOAT)):FLOAT" and SNOWFLAKE.ACCOUNT_USAGE."TAG_REFERENCES_WITH_LINEAGE(TAG_NAME_INPUT VARCHAR):TABLE:`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/To/AccountRole"),
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_To_DatabaseRole(t *testing.T) {
	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifier()
	databaseRoleModel := model.DatabaseRole("test", databaseRoleId.DatabaseName(), databaseRoleId.Name())
	grantsModel := datasourcemodel.GrantsToDatabaseRole("test", databaseRoleId).
		WithDependsOn(databaseRoleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseRoleModel, grantsModel),
				Check:  checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_To_User(t *testing.T) {
	userId := testClient().Context.CurrentUser(t)
	grantsModel := datasourcemodel.GrantsToUser("test", userId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantsModel),
				Check:  checkAtLeastOneGrantPresentLimited(),
			},
		},
	})
}

func TestAcc_Grants_To_Share(t *testing.T) {
	shareId := testClient().Ids.RandomAccountObjectIdentifier()
	configVariables := config.Variables{
		"database": config.StringVariable(TestDatabaseName),
		"share":    config.StringVariable(shareId.Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/To/Share"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

// TODO [SNOW-1284382]: Implement after SHOW GRANTS TO SHARE <share_name> IN APPLICATION PACKAGE <app_package_name> syntax starts working.
func TestAcc_Grants_To_ShareWithApplicationPackage(t *testing.T) {
	t.Skip("Skipped until SHOW GRANTS TO SHARE <share_name> IN APPLICATION PACKAGE <app_package_name> syntax starts working.")
}

func TestAcc_Grants_To_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/To/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_ShareNameMissing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/To/Invalid/ShareNameMissing"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/To/Invalid/DatabaseRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_ApplicationRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/To/Invalid/ApplicationRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_Of_AccountRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/Of/AccountRole"),
				Check:           checkAtLeastOneGrantPresentLimited(),
			},
		},
	})
}

func TestAcc_Grants_Of_DatabaseRole(t *testing.T) {
	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifier()
	configVariables := config.Variables{
		"database":      config.StringVariable(TestDatabaseName),
		"database_role": config.StringVariable(databaseRoleId.Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/Of/DatabaseRole"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresentLimited(),
			},
		},
	})
}

// TODO [SNOW-1284382]: Implement after snowflake_application and snowflake_application_role resources are introduced.
func TestAcc_Grants_Of_ApplicationRole(t *testing.T) {
	t.Skip("Skipped until snowflake_application and snowflake_application_role resources are introduced. Currently, behavior tested in application_roles_gen_integration_test.go.")
}

// TODO [SNOW-1284394]: Unskip the test
func TestAcc_Grants_Of_Share(t *testing.T) {
	t.Skip("TestAcc_Share are skipped")

	shareId := testClient().Ids.RandomAccountObjectIdentifier()
	accountId := secondaryTestClient().Account.GetAccountIdentifier(t)
	require.NotNil(t, accountId)

	configVariables := config.Variables{
		"database": config.StringVariable(TestDatabaseName),
		"share":    config.StringVariable(shareId.Name()),
		"account":  config.StringVariable(accountId.FullyQualifiedName()),
	}

	datasourceName := "data.snowflake_grants.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/Of/Share"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
					resource.TestCheckNoResourceAttr(datasourceName, "grants.0.created_on"),
				),
			},
		},
	})
}

func TestAcc_Grants_Of_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/Of/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_Of_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/Of/Invalid/DatabaseRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_Of_Invalid_ApplicationRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/Of/Invalid/ApplicationRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Database(t *testing.T) {
	configVariables := config.Variables{
		"database": config.StringVariable(TestDatabaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/FutureIn/Database"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Schema(t *testing.T) {
	configVariables := config.Variables{
		"database": config.StringVariable(TestDatabaseName),
		"schema":   config.StringVariable(TestSchemaName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/FutureIn/Schema"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/FutureIn/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Invalid_SchemaNameNotFullyQualified(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/FutureIn/Invalid/SchemaNameNotFullyQualified"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_AccountRole(t *testing.T) {
	configVariables := config.Variables{
		"database": config.StringVariable(TestDatabaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/FutureTo/AccountRole"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_DatabaseRole(t *testing.T) {
	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifier()
	configVariables := config.Variables{
		"database":      config.StringVariable(TestDatabaseName),
		"database_role": config.StringVariable(databaseRoleId.Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/FutureTo/DatabaseRole"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/FutureTo/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Grants/FutureTo/Invalid/DatabaseRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func checkAtLeastOneGrantPresent() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.created_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.privilege"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_to"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grantee_name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grant_option"),
	)
}

func checkAtLeastOneFutureGrantPresent() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.created_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.privilege"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grantee_name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grant_option"),
	)
}

func checkAtLeastOneGrantPresentLimited() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.created_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_to"),
	)
}
