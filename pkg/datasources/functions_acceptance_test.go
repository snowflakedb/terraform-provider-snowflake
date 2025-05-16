//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Functions(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	dataSourceName := "data.snowflake_functions.functions"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			{
				Config: functionsConfig(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(dataSourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttrSet(dataSourceName, "functions.#"),
				),
			},
		},
	})
}

// TODO [SNOW-1348103]: use generated config builder when reworking the datasource
func functionsConfig(t *testing.T) string {
	t.Helper()

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	id1 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	id2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	functionsSetup := config.FromModels(t,
		model.FunctionJavaBasicInline("f1", id1, dataType, handler, definition).WithArgument(argName, dataType),
		model.FunctionJavaBasicInline("f2", id2, dataType, handler, definition).WithArgument(argName, dataType),
	)

	return fmt.Sprintf(`
%s
data "snowflake_functions" "functions" {
  database   = "%s"
  schema     = "%s"
  depends_on = [snowflake_function_java.f1, snowflake_function_java.f2]
}
`, functionsSetup, acc.TestDatabaseName, acc.TestSchemaName)
}
