//go:build !account_level_tests

package testacc

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Procedures(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	dataSourceName := "data.snowflake_procedures.procedures"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ProcedureJava),
		Steps: []resource.TestStep{
			{
				Config: proceduresConfig(t, schema.ID()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "database", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(dataSourceName, "schema", schema.ID().Name()),
					// Extra 1 in procedure count above due to ASSOCIATE_SEMANTIC_CATEGORY_TAGS appearing in all "SHOW PROCEDURES IN ..." commands
					resource.TestCheckResourceAttr(dataSourceName, "procedures.#", "3"),
				),
			},
		},
	})
}

// TODO [SNOW-1348103]: use generated config builder when reworking the datasource
func proceduresConfig(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) string {
	t.Helper()

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := testClient().Procedure.SampleJavaDefinition(t, className, funcName, argName)

	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchemaNewDataTypes(schemaId, dataType)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchemaNewDataTypes(schemaId, dataType)

	functionsSetup := config.FromModels(t,
		model.ProcedureJavaBasicInline("p1", id1, dataType, handler, definition).WithArgument(argName, dataType),
		model.ProcedureJavaBasicInline("p2", id2, dataType, handler, definition).WithArgument(argName, dataType),
	)

	return fmt.Sprintf(`
%s
data "snowflake_procedures" "procedures" {
  database   = "%s"
  schema     = "%s"
  depends_on = [snowflake_procedure_java.p1, snowflake_procedure_java.p2]
}
`, functionsSetup, schemaId.DatabaseName(), schemaId.Name())
}
