//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ObjectParameter(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: objectParameterConfigBasic(database.ID(), sdk.DatabaseParameterUserTaskTimeoutMs, "1000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.DatabaseParameterUserTaskTimeoutMs)),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "1000"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "false"),
				),
			},
		},
	})
}

func objectParameterConfigBasic(databaseId sdk.AccountObjectIdentifier, key sdk.DatabaseParameter, value string) string {
	return fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[2]s"
	value = "%[3]s"
	object_type = "DATABASE"
	object_identifier {
		name = "%[1]s"
	}
}
`, databaseId.Name(), key, value)
}

func TestAcc_ObjectParameterAccount(t *testing.T) {
	// TODO [SNOW-2010844]: unskip
	t.Skip("Skipping temporarily as it messes with the account level setting.")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: objectParameterConfigOnAccount(sdk.AccountParameterDataRetentionTimeInDays, "5"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.AccountParameterDataRetentionTimeInDays)),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "5"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "true"),
				),
			},
		},
	})
}

func objectParameterConfigOnAccount(key sdk.AccountParameter, value string) string {
	return fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[1]s"
	value = "%[2]s"
	on_account = true
}
`, key, value)
}

func TestAcc_UserParameter(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: userParameterConfigBasic(user.ID(), sdk.UserParameterEnableUnredactedQuerySyntaxError, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.UserParameterEnableUnredactedQuerySyntaxError)),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "true"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "false"),
				),
			},
		},
	})
}

func userParameterConfigBasic(userId sdk.AccountObjectIdentifier, key sdk.UserParameter, value string) string {
	return fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[2]s"
	value = "%[3]s"
	object_type = "USER"
	object_identifier {
		name = "%[1]s"
	}
}
`, userId.Name(), key, value)
}
