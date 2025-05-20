//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalFunction_basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(id.DatabaseName()),
			"schema":                    config.StringVariable(id.SchemaName()),
			"name":                      config.StringVariable(id.Name()),
			"api_allowed_prefixes":      config.ListVariable(config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/")),
			"url_of_proxy_and_resource": config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
			"comment":                   config.StringVariable("Terraform acceptance test"),
		}
	}

	resourceName := "snowflake_external_function.external_function"
	configVariables := m()
	configVariables2 := m()
	configVariables2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalFunction/basic"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "arg.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "arg.0.name", "ARG1"),
					resource.TestCheckResourceAttr(resourceName, "arg.0.type", "VARCHAR"),
					resource.TestCheckResourceAttr(resourceName, "arg.1.name", "ARG2"),
					resource.TestCheckResourceAttr(resourceName, "arg.1.type", "VARCHAR"),
					resource.TestCheckResourceAttr(resourceName, "null_input_behavior", "CALLED ON NULL INPUT"),
					resource.TestCheckResourceAttr(resourceName, "return_type", "VARIANT"),
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName, "return_behavior", "IMMUTABLE"),
					resource.TestCheckResourceAttrSet(resourceName, "api_integration"),
					resource.TestCheckResourceAttr(resourceName, "compression", "AUTO"),
					resource.TestCheckResourceAttr(resourceName, "url_of_proxy_and_resource", "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalFunction/basic"),
				ConfigVariables: configVariables2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},
			// IMPORT
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalFunction/basic"),
				ConfigVariables:   configVariables2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// these two are not found in either the show or describe command
				ImportStateVerifyIgnore: []string{"return_null_allowed", "api_integration"},
			},
		},
	})
}

func TestAcc_ExternalFunction_no_arguments(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(id.DatabaseName()),
			"schema":                    config.StringVariable(id.SchemaName()),
			"name":                      config.StringVariable(id.Name()),
			"api_allowed_prefixes":      config.ListVariable(config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/")),
			"url_of_proxy_and_resource": config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
			"comment":                   config.StringVariable("Terraform acceptance test"),
		}
	}

	resourceName := "snowflake_external_function.external_function"
	configVariables := m()
	configVariables2 := m()
	configVariables2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalFunction/no_arguments"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "arg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "null_input_behavior", "CALLED ON NULL INPUT"),
					resource.TestCheckResourceAttr(resourceName, "return_type", "VARIANT"),
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName, "return_behavior", "IMMUTABLE"),
					resource.TestCheckResourceAttrSet(resourceName, "api_integration"),
					resource.TestCheckResourceAttr(resourceName, "compression", "AUTO"),
					resource.TestCheckResourceAttr(resourceName, "url_of_proxy_and_resource", "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalFunction/no_arguments"),
				ConfigVariables: configVariables2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},
			// IMPORT
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalFunction/no_arguments"),
				ConfigVariables:   configVariables2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// these two are not found in either the show or describe command
				ImportStateVerifyIgnore: []string{"return_null_allowed", "api_integration"},
			},
		},
	})
}

func TestAcc_ExternalFunction_complete(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifierWithArguments()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(id.DatabaseName()),
			"schema":                    config.StringVariable(id.SchemaName()),
			"name":                      config.StringVariable(id.Name()),
			"api_allowed_prefixes":      config.ListVariable(config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/")),
			"url_of_proxy_and_resource": config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
			"comment":                   config.StringVariable("Terraform acceptance test"),
		}
	}

	resourceName := "snowflake_external_function.external_function"
	configVariables := m()
	configVariables2 := m()
	configVariables2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalFunction/complete"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "arg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "null_input_behavior", "CALLED ON NULL INPUT"),
					resource.TestCheckResourceAttr(resourceName, "return_type", "VARIANT"),
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName, "return_behavior", "IMMUTABLE"),
					resource.TestCheckResourceAttrSet(resourceName, "api_integration"),
					resource.TestCheckResourceAttr(resourceName, "compression", "AUTO"),
					resource.TestCheckResourceAttr(resourceName, "url_of_proxy_and_resource", "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
					resource.TestCheckResourceAttr(resourceName, "header.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "header.0.name", "x-custom-header"),
					resource.TestCheckResourceAttr(resourceName, "header.0.value", "snowflake"),
					resource.TestCheckResourceAttr(resourceName, "max_batch_rows", "500"),
					resource.TestCheckResourceAttr(resourceName, "request_translator", fmt.Sprintf("%s.%s.%s%s", TestDatabaseName, TestSchemaName, id.Name(), "_request_translator")),
					resource.TestCheckResourceAttr(resourceName, "response_translator", fmt.Sprintf("%s.%s.%s%s", TestDatabaseName, TestSchemaName, id.Name(), "_response_translator")),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ExternalFunction/complete"),
				ConfigVariables: configVariables2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},
			// IMPORT
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ExternalFunction/complete"),
				ConfigVariables:   configVariables2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// these four are not found in either the show or describe command
				ImportStateVerifyIgnore: []string{"return_null_allowed", "api_integration", "request_translator", "response_translator"},
			},
		},
	})
}

func TestAcc_ExternalFunction_migrateFromVersion085(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsOld(sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)
	name := id.Name()
	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.85.0"),
				Config:            externalFunctionConfig(TestDatabaseName, TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|%s|%s|VARCHAR-VARCHAR", TestDatabaseName, TestSchemaName, name)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", "\""+TestDatabaseName+"\""),
					resource.TestCheckResourceAttr(resourceName, "schema", "\""+TestSchemaName+"\""),
					resource.TestCheckNoResourceAttr(resourceName, "return_null_allowed"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            externalFunctionConfig(TestDatabaseName, TestSchemaName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
			},
		},
	})
}

func TestAcc_ExternalFunction_migrateFromVersion085_issue2694_previousValuePresent(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.85.0"),
				Config:            externalFunctionConfigWithReturnNullAllowed(TestDatabaseName, TestSchemaName, name, sdk.Bool(true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            externalFunctionConfig(TestDatabaseName, TestSchemaName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
			},
		},
	})
}

func TestAcc_ExternalFunction_migrateFromVersion085_issue2694_previousValueRemoved(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.85.0"),
				Config:            externalFunctionConfigWithReturnNullAllowed(TestDatabaseName, TestSchemaName, name, sdk.Bool(true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.85.0"),
				Config:            externalFunctionConfig(TestDatabaseName, TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "return_null_allowed"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   externalFunctionConfig(TestDatabaseName, TestSchemaName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
			},
		},
	})
}

// Proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2528.
// The problem originated from ShowById without IN clause. There was no IN clause in the docs at the time.
// It was raised with the appropriate team in Snowflake.
func TestAcc_ExternalFunction_issue2528(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	accName := id.Name()
	secondSchemaId := testClient().Ids.RandomDatabaseObjectIdentifier()
	secondSchema := secondSchemaId.Name()

	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				Config: externalFunctionConfigIssue2528(TestDatabaseName, TestSchemaName, accName, secondSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", accName),
				),
			},
		},
	})
}

func TestAcc_ExternalFunction_issue3392_returnVarchar(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	accName := id.Name()

	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				Config: externalFunctionConfigWithReturnType(TestDatabaseName, TestSchemaName, accName, "VARCHAR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_type", "VARCHAR"),
				),
			},
		},
	})
}

func TestAcc_ExternalFunction_issue3392_returnVarcharWithSize(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	accName := id.Name()

	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				Config: externalFunctionConfigWithReturnType(TestDatabaseName, TestSchemaName, accName, "VARCHAR(10)"),
				Check: resource.ComposeTestCheckFunc(
					// Snowflake drops the size from VARCHAR when it's specified
					resource.TestCheckResourceAttr(resourceName, "return_type", "VARCHAR"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Proves that header parsing handles values wrapped in curly braces, e.g. `value = "{1}"`
func TestAcc_ExternalFunction_HeaderParsing(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.93.0"),
				Config:            externalFunctionConfigIssueCurlyHeader(id),
				// Previous implementation produces a plan with the following changes
				//
				// - header { # forces replacement
				//   - name  = "name" -> null
				//   - value = "0" -> null
				// }
				//
				// + header { # forces replacement
				//   + name  = "name"
				//   + value = "{0}"
				// }
				ExpectNonEmptyPlan: true,
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   externalFunctionConfigIssueCurlyHeader(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "header.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "header.0.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "header.0.value", "{0}"),
				),
			},
		},
	})
}

func externalFunctionConfig(database string, schema string, name string) string {
	return externalFunctionConfigWithReturnNullAllowed(database, schema, name, nil)
}

func externalFunctionConfigWithReturnNullAllowed(database string, schema string, name string, returnNullAllowed *bool) string {
	returnNullAllowedText := ""
	if returnNullAllowed != nil {
		returnNullAllowedText = fmt.Sprintf("return_null_allowed = \"%t\"", *returnNullAllowed)
	}

	return fmt.Sprintf(`
resource "snowflake_api_integration" "test_api_int" {
 name                 = "%[3]s"
 api_provider         = "aws_api_gateway"
 api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
 api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
 enabled              = true
}

resource "snowflake_external_function" "f" {
 name     = "%[3]s"
 database = "%[1]s"
 schema   = "%[2]s"
 arg {
   name = "ARG1"
   type = "VARCHAR"
 }
 arg {
   name = "ARG2"
   type = "VARCHAR"
 }
 return_type               = "VARIANT"
 return_behavior           = "IMMUTABLE"
 api_integration           = snowflake_api_integration.test_api_int.name
 url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
 %[4]s
}

`, database, schema, name, returnNullAllowedText)
}

func externalFunctionConfigWithReturnType(database string, schema string, name string, returnType string) string {
	return fmt.Sprintf(`
resource "snowflake_api_integration" "test_api_int" {
 name                 = "%[3]s"
 api_provider         = "aws_api_gateway"
 api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
 api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
 enabled              = true
}

resource "snowflake_external_function" "f" {
 name     = "%[3]s"
 database = "%[1]s"
 schema   = "%[2]s"
 arg {
   name = "ARG1"
   type = "VARCHAR"
 }
 arg {
   name = "ARG2"
   type = "VARCHAR"
 }
 return_type               = "%[4]s"
 return_behavior           = "IMMUTABLE"
 api_integration           = snowflake_api_integration.test_api_int.name
 url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}

`, database, schema, name, returnType)
}

func externalFunctionConfigIssue2528(database string, schema string, name string, schema2 string) string {
	return fmt.Sprintf(`
resource "snowflake_api_integration" "test_api_int" {
 name                 = "%[3]s"
 api_provider         = "aws_api_gateway"
 api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
 api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
 enabled              = true
}

resource "snowflake_schema" "s2" {
 database            = "%[1]s"
 name                = "%[4]s"
}

resource "snowflake_external_function" "f" {
 name     = "%[3]s"
 database = "%[1]s"
 schema   = "%[2]s"
 arg {
   name = "SNS_NOTIF"
   type = "OBJECT"
 }
 return_type = "VARIANT"
 return_behavior = "VOLATILE"
 api_integration = snowflake_api_integration.test_api_int.name
 url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}

resource "snowflake_external_function" "f2" {
 depends_on = [snowflake_schema.s2]

 name     = "%[3]s"
 database = "%[1]s"
 schema   = "%[4]s"
 arg {
   name = "SNS_NOTIF"
   type = "OBJECT"
 }
 return_type = "VARIANT"
 return_behavior = "VOLATILE"
 api_integration = snowflake_api_integration.test_api_int.name
 url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}
`, database, schema, name, schema2)
}

func externalFunctionConfigIssueCurlyHeader(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_api_integration" "test_api_int" {
 name                 = "%[3]s"
 api_provider         = "aws_api_gateway"
 api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
 api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
 enabled              = true
}

resource "snowflake_external_function" "f" {
 name     = "%[3]s"
 database = "%[1]s"
 schema   = "%[2]s"
 arg {
   name = "ARG1"
   type = "VARCHAR"
 }
 arg {
   name = "ARG2"
   type = "VARCHAR"
 }
 header {
	name = "name"
	value = "{0}"
 }
 return_type               = "VARIANT"
 return_behavior           = "IMMUTABLE"
 api_integration           = snowflake_api_integration.test_api_int.name
 url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}

`, id.DatabaseName(), id.SchemaName(), id.Name())
}

func TestAcc_ExternalFunction_EnsureSmoothResourceIdMigrationToV0950(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            externalFunctionConfigWithMoreArguments(TestDatabaseName, TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%s"."%s"."%s"(VARCHAR, FLOAT, NUMBER)`, TestDatabaseName, TestSchemaName, name)),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   externalFunctionConfigWithMoreArguments(TestDatabaseName, TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%s"."%s"."%s"(VARCHAR, FLOAT, NUMBER)`, TestDatabaseName, TestSchemaName, name)),
				),
			},
		},
	})
}

func externalFunctionConfigWithMoreArguments(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_api_integration" "test_api_int" {
 name                 = "%[3]s"
 api_provider         = "aws_api_gateway"
 api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
 api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
 enabled              = true
}

resource "snowflake_external_function" "f" {
 database = "%[1]s"
 schema   = "%[2]s"
 name     = "%[3]s"

 arg {
   name = "ARG1"
   type = "VARCHAR"
 }

 arg {
   name = "ARG2"
   type = "FLOAT"
 }

 arg {
   name = "ARG3"
   type = "NUMBER"
 }

 return_type               = "VARIANT"
 return_behavior           = "IMMUTABLE"
 api_integration           = snowflake_api_integration.test_api_int.name
 url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}
`, database, schema, name)
}

func TestAcc_ExternalFunction_EnsureSmoothResourceIdMigrationToV0950_WithoutArguments(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            externalFunctionConfigWithoutArguments(TestDatabaseName, TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%s"."%s"."%s"`, TestDatabaseName, TestSchemaName, name)),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   externalFunctionConfigWithoutArguments(TestDatabaseName, TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%s"."%s"."%s"()`, TestDatabaseName, TestSchemaName, name)),
				),
			},
		},
	})
}

func externalFunctionConfigWithoutArguments(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_api_integration" "test_api_int" {
 name                 = "%[3]s"
 api_provider         = "aws_api_gateway"
 api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
 api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
 enabled              = true
}

resource "snowflake_external_function" "f" {
 database = "%[1]s"
 schema   = "%[2]s"
 name     = "%[3]s"

 return_type               = "VARIANT"
 return_behavior           = "IMMUTABLE"
 api_integration           = snowflake_api_integration.test_api_int.name
 url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}

`, database, schema, name)
}
