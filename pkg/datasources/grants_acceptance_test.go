package datasources_test

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Grants_On_Account(t *testing.T) {
	grantsModel := datasourcemodel.GrantsOnAccount("test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	grantsModel := datasourcemodel.GrantsOnAccountObject("test", acc.TestClient().Ids.DatabaseId(), sdk.ObjectTypeDatabase)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	grantsModel := datasourcemodel.GrantsOnDatabaseObject("test", acc.TestClient().Ids.SchemaId(), sdk.ObjectTypeSchema)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	viewId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	columnNames := []string{"ROLE_NAME"}

	viewModel := model.View("test", viewId.DatabaseName(), viewId.Name(), viewId.SchemaName(), statement).WithColumnNames(columnNames...)
	grantsModel := datasourcemodel.GrantsOnSchemaObject("test", viewId, sdk.ObjectTypeView).
		WithDependsOn(viewModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	function := acc.TestClient().Function.Create(t, sdk.DataTypeVARCHAR)
	grantsModel := datasourcemodel.GrantsOnSchemaObjectWithArguments("test", function.ID(), sdk.ObjectTypeFunction)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/AccountRole"),
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_To_DatabaseRole(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	databaseRoleModel := model.DatabaseRole("test", databaseRoleId.DatabaseName(), databaseRoleId.Name())
	grantsModel := datasourcemodel.GrantsToDatabaseRole("test", databaseRoleId).
		WithDependsOn(databaseRoleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
	userId := acc.TestClient().Context.CurrentUser(t)
	configVariables := config.Variables{
		"user": config.StringVariable(userId.Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/User"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresentLimited(),
			},
		},
	})
}

func TestAcc_Grants_To_Share(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	shareName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
		"share":    config.StringVariable(shareName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Share"),
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_ShareNameMissing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Invalid/ShareNameMissing"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Invalid/DatabaseRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_ApplicationRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Invalid/ApplicationRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_Of_AccountRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/AccountRole"),
				Check:           checkAtLeastOneGrantPresentLimited(),
			},
		},
	})
}

func TestAcc_Grants_Of_DatabaseRole(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	databaseRoleName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"database":      config.StringVariable(databaseName),
		"database_role": config.StringVariable(databaseRoleName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/DatabaseRole"),
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
	databaseName := acc.TestClient().Ids.Alpha()
	shareName := acc.TestClient().Ids.Alpha()

	accountId := acc.SecondaryTestClient().Account.GetAccountIdentifier(t)
	require.NotNil(t, accountId)

	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
		"share":    config.StringVariable(shareName),
		"account":  config.StringVariable(accountId.FullyQualifiedName()),
	}
	datasourceName := "data.snowflake_grants.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Share"),
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_Of_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Invalid/DatabaseRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_Of_Invalid_ApplicationRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Invalid/ApplicationRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Database(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureIn/Database"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Schema(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
		"schema":   config.StringVariable(schemaName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureIn/Schema"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureIn/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Invalid_SchemaNameNotFullyQualified(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureIn/Invalid/SchemaNameNotFullyQualified"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_AccountRole(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/AccountRole"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_DatabaseRole(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	databaseRoleName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"database":      config.StringVariable(databaseName),
		"database_role": config.StringVariable(databaseRoleName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/DatabaseRole"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/Invalid/DatabaseRoleIdInvalid"),
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
