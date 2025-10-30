//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Functions(t *testing.T) {
	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	dataSourceName := "data.snowflake_functions.functions"
	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			{
				Config: functionsConfig(t, schema.ID()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "database", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(dataSourceName, "schema", schema.ID().Name()),
					resource.TestCheckResourceAttr(dataSourceName, "functions.#", "2"),
				),
			},
		},
	})
}

// TODO [SNOW-1348103]: use generated config builder when reworking the datasource
func functionsConfig(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) string {
	t.Helper()

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := testClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchemaNewDataTypes(schemaId, dataType)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchemaNewDataTypes(schemaId, dataType)

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
`, functionsSetup, schemaId.DatabaseName(), schemaId.Name())
}
