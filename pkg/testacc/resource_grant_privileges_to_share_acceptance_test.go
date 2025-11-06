//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantPrivilegesToShare_OnDatabase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	configVariables := config.Variables{
		"to_share": config.StringVariable(share.ID().Name()),
		"database": config.StringVariable(database.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.ObjectPrivilegeUsage.String()),
		),
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", database.ID().Name()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase_NoGrant"),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnSchema(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(share.ID().Name()),
			"database": config.StringVariable(database.ID().Name()),
			"schema":   config.StringVariable(schemaId.Name()),
		}
		if withGrant {
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeUsage.String()),
			)
		}
		return variables
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnSchema"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_schema", schemaId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnSchema"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnSchema_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnTable(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(share.ID().Name()),
			"database": config.StringVariable(database.ID().Name()),
			"schema":   config.StringVariable(schema.ID().Name()),
			"on_table": config.StringVariable(tableId.Name()),
		}
		if withGrant {
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeSelect.String()),
			)
		}
		return variables
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_table", tableId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnDynamicTable(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	configVariables := config.Variables{
		"to_share":  config.StringVariable(share.ID().Name()),
		"database":  config.StringVariable(id.DatabaseName()),
		"schema":    config.StringVariable(id.SchemaName()),
		"warehouse": config.StringVariable(TestWarehouseName),
		"on_table":  config.StringVariable(id.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.ObjectPrivilegeSelect.String()),
		),
	}

	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDynamicTable"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_table", id.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDynamicTable"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnAllTablesInSchema(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	configVariables := config.Variables{
		"to_share": config.StringVariable(share.ID().Name()),
		"database": config.StringVariable(database.ID().Name()),
		"schema":   config.StringVariable(schema.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.ObjectPrivilegeSelect.String()),
		),
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnAllTablesInSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_all_tables_in_schema", schema.ID().FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnAllTablesInSchema"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnAllTablesInSchema_NoGrant"),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnView(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	viewId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(share.ID().Name()),
			"database": config.StringVariable(database.ID().Name()),
			"schema":   config.StringVariable(schema.ID().Name()),
			"on_table": config.StringVariable(tableId.Name()),
			"on_view":  config.StringVariable(viewId.Name()),
		}
		if withGrant {
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeSelect.String()),
			)
		}
		return variables
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: viewsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnView"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_view", viewId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnView"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnView_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnTag(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	tagId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(share.ID().Name()),
			"database": config.StringVariable(database.ID().Name()),
			"schema":   config.StringVariable(schema.ID().Name()),
			"on_tag":   config.StringVariable(tagId.Name()),
		}
		if withGrant {
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeRead.String()),
			)
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTag"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeRead.String()),
					resource.TestCheckResourceAttr(resourceName, "on_tag", tagId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTag"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTag_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnSchemaObject_OnFunctionWithArguments(t *testing.T) {
	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	function := testClient().Function.CreateSecure(t, sdk.DataTypeFloat)

	configVariables := config.Variables{
		"name":          config.StringVariable(share.ID().Name()),
		"function_name": config.StringVariable(function.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":      config.StringVariable(TestDatabaseName),
		"schema":        config.StringVariable(TestSchemaName),
		"argument_type": config.StringVariable(string(sdk.DataTypeFloat)),
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnFunction"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_function", function.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|USAGE|OnFunction|%s", share.ID().FullyQualifiedName(), function.ID().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnFunction"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnSchemaObject_OnFunctionWithoutArguments(t *testing.T) {
	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	function := testClient().Function.CreateSecure(t)

	configVariables := config.Variables{
		"name":          config.StringVariable(share.ID().Name()),
		"function_name": config.StringVariable(function.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":      config.StringVariable(TestDatabaseName),
		"schema":        config.StringVariable(TestSchemaName),
		"argument_type": config.StringVariable(""),
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnFunction"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_function", function.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|USAGE|OnFunction|%s", share.ID().FullyQualifiedName(), function.ID().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnFunction"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnPrivilegeUpdate(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	configVariables := func(privileges []sdk.ObjectPrivilege) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(share.ID().Name()),
			"database": config.StringVariable(database.ID().Name()),
		}
		if len(privileges) > 0 {
			configPrivileges := make([]config.Variable, len(privileges))
			for i, privilege := range privileges {
				configPrivileges[i] = config.StringVariable(privilege.String())
			}
			variables["privileges"] = config.ListVariable(configPrivileges...)
		}
		return variables
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables([]sdk.ObjectPrivilege{
					sdk.ObjectPrivilegeUsage,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", database.ID().Name()),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables([]sdk.ObjectPrivilege{
					sdk.ObjectPrivilegeUsage,
					sdk.ObjectPrivilegeReferenceUsage,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeReferenceUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", database.ID().Name()),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables([]sdk.ObjectPrivilege{
					sdk.ObjectPrivilegeUsage,
					sdk.ObjectPrivilegeReferenceUsage,
				}),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase_NoGrant"),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnDatabaseWithReferenceUsagePrivilege(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	configVariables := config.Variables{
		"to_share": config.StringVariable(share.ID().Name()),
		"database": config.StringVariable(database.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.ObjectPrivilegeUsage.String()),
			config.StringVariable(sdk.ObjectPrivilegeReferenceUsage.String()),
		),
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeReferenceUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", database.ID().Name()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase_NoGrant"),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_NoPrivileges(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase_NoPrivileges"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`The argument "privileges" is required, but no definition was found.`),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_NoOnOption(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/NoOnOption"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`Invalid combination of arguments`),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToShare_RemoveShareOutsideTerraform(t *testing.T) {
	t.Skip("Should be addressed in SNOW-2048471 - without this task done, the test cannot be properly asserted in the terraform testing framework")

	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	configVariables := config.Variables{
		"to_share": config.StringVariable(share.ID().Name()),
		"database": config.StringVariable(database.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.ObjectPrivilegeUsage.String()),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnCustomShare"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { shareCleanup() },
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnCustomShare"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile(sdk.ErrObjectNotExistOrAuthorized.Error()),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShareWithNameContainingDots_OnTable(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	shareId := testClient().Ids.RandomAccountObjectIdentifierContaining(".foo.bar")
	_, shareCleanup := testClient().Share.CreateShareWithIdentifier(t, shareId)
	t.Cleanup(shareCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareId.Name()),
			"database": config.StringVariable(tableId.DatabaseName()),
			"schema":   config.StringVariable(tableId.SchemaName()),
			"on_table": config.StringVariable(tableId.Name()),
		}
		if withGrant {
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeSelect.String()),
			)
		}
		return variables
	}

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", shareId.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_table", tableId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantPrivilegesToShareBasicConfig(database.ID(), share.ID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "id", fmt.Sprintf(`%s|USAGE|OnDatabase|%s`, share.ID().FullyQualifiedName(), database.ID().FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToShareBasicConfig(database.ID(), share.ID()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_share.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_share.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "id", fmt.Sprintf(`%s|USAGE|OnDatabase|%s`, share.ID().FullyQualifiedName(), database.ID().FullyQualifiedName())),
				),
			},
		},
	})
}

func grantPrivilegesToShareBasicConfig(databaseId sdk.AccountObjectIdentifier, shareId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = "%[2]s"
  privileges  = ["USAGE"]
  on_database = "%[1]s"
}
`, databaseId.Name(), shareId.Name())
}

func TestAcc_GrantPrivilegesToShare_IdentifierQuotingDiffSuppression(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	shareId := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				ExpectError:       regexp.MustCompile("Error: Provider produced inconsistent final plan"),
				Config:            grantPrivilegesToShareQuotedIdentifiers(database.ID(), shareId),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToShareQuotedIdentifiers(database.ID(), shareId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_share.test", plancheck.ResourceActionCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_share.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "to_share", shareId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "on_database", database.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "id", fmt.Sprintf(`%s|USAGE|OnDatabase|%s`, shareId.FullyQualifiedName(), database.ID().FullyQualifiedName())),
				),
			},
		},
	})
}

func grantPrivilegesToShareQuotedIdentifiers(databaseId sdk.AccountObjectIdentifier, shareId sdk.AccountObjectIdentifier) string {
	quotedShareId := fmt.Sprintf(`\"%s\"`, shareId.Name())

	return fmt.Sprintf(`
resource "snowflake_share" "test" {
  name       = "%[2]s"
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = snowflake_share.test.name
  privileges  = ["USAGE"]
  on_database = "%[1]s"
}
`, databaseId.Name(), quotedShareId)
}
