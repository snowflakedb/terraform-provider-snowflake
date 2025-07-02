//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Procedures(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

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
				Config: proceduresConfig(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(dataSourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttrSet(dataSourceName, "procedures.#"),
					// resource.TestCheckResourceAttr(dataSourceName, "procedures.#", "3"),
					// Extra 1 in procedure count above due to ASSOCIATE_SEMANTIC_CATEGORY_TAGS appearing in all "SHOW PROCEDURES IN ..." commands
				),
			},
		},
	})
}

// TODO [SNOW-1348103]: use generated config builder when reworking the datasource
func proceduresConfig(t *testing.T) string {
	t.Helper()

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := testClient().Procedure.SampleJavaDefinition(t, className, funcName, argName)

	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

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
`, functionsSetup, TestDatabaseName, TestSchemaName)
}

func TestAcc_Procedures_gh3822_bcr2025_03(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	schema, schemaCleanup := secondaryTestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	_, procedureCleanup := secondaryTestClient().Procedure.CreatePythonInSchema(t, schema.ID())
	t.Cleanup(procedureCleanup)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.ProceduresDatasource))
	proceduresModel := datasourcemodel.Procedures("test", schema.ID().DatabaseName(), schema.ID().Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.2.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.DisableBcrBundle(t, "2025_03")
				},
				Config: config.FromModels(t, providerModel, proceduresModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(proceduresModel.DatasourceReference(), "database", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(proceduresModel.DatasourceReference(), "schema", schema.ID().Name()),
					resource.TestCheckResourceAttrSet(proceduresModel.DatasourceReference(), "procedures.#"),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.2.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_03")
				},
				Config:      config.FromModels(t, providerModel, proceduresModel),
				ExpectError: regexp.MustCompile("could not parse arguments"),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providerModel, proceduresModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(proceduresModel.DatasourceReference(), "database", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(proceduresModel.DatasourceReference(), "schema", schema.ID().Name()),
					resource.TestCheckResourceAttrSet(proceduresModel.DatasourceReference(), "procedures.#"),
				),
			},
		},
	})
}
