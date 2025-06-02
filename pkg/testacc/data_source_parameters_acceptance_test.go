//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ParametersOnAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnAccount(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "pattern", "AUTOCOMMIT"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "parameters.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "parameters.0.key", "AUTOCOMMIT"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.value"),
				),
			},
		},
	})
}

func TestAcc_ParametersOnSession(t *testing.T) {
	userId := testClient().Context.CurrentUser(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnSession(userId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.key"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "user", userId.Name()),
				),
			},
		},
	})
}

func TestAcc_ParametersOnObject(t *testing.T) {
	dbId := testClient().Ids.DatabaseId()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnObject(dbId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.key"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "object_type", "DATABASE"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "object_name", dbId.Name()),
				),
			},
		},
	})
}

func parametersConfigOnAccount() string {
	return `data "snowflake_parameters" "p" {
		parameter_type = "ACCOUNT"
		pattern = "AUTOCOMMIT"
	}`
}

func parametersConfigOnSession(userId sdk.AccountObjectIdentifier) string {
	s := `
	data "snowflake_parameters" "p" {
		parameter_type = "SESSION"
		user = "%s"
	}`
	return fmt.Sprintf(s, userId.Name())
}

func parametersConfigOnObject(databaseId sdk.AccountObjectIdentifier) string {
	stmt := `
	data "snowflake_parameters" "p" {
		parameter_type = "OBJECT"
		object_type = "DATABASE"
		object_name = "%s"
	}`
	return fmt.Sprintf(stmt, databaseId.Name())
}
