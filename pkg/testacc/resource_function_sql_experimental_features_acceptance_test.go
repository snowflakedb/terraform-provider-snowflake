//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FunctionSql_ParametersIgnoreValueChangesIfNotOnObjectLevel(t *testing.T) {
	database, databaseCleanup := secondaryTestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)
	secondaryTestClient().Database.UpdateLogLevel(t, database.ID(), sdk.LogLevelError)

	schema, schemaCleanup := secondaryTestClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)
	secondaryTestClient().Schema.UpdateLogLevel(t, schema.ID(), sdk.LogLevelWarn)

	argName := "abc"
	dataType := testdatatypes.DataTypeFloat
	id := secondaryTestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchemaNewDataTypes(schema.ID(), dataType)

	definition := secondaryTestClient().Function.SampleSqlDefinitionWithArgument(t, argName)

	functionModel := model.FunctionSqlBasicInline("test", id, definition, dataType.ToLegacyDataTypeSql()).
		WithArgument(argName, dataType)
	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithExperimentalFeaturesEnabled(experimentalfeatures.ParametersIgnoreValueChangesIfNotOnObjectLevel)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: functionsAndProceduresProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Parameter value taken from schema because it's set
			{
				Config: config.FromModels(t, providerModel, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionSqlResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasLogLevelString(string(sdk.LogLevelWarn)),
					// modifying the parameter as assertion, as this is the moment in-between apply and refresh
					invokeactionassert.Invoke(t, func() error {
						secondaryTestClient().Schema.UpdateLogLevel(t, schema.ID(), sdk.LogLevelInfo)
						return nil
					}),
				),
			},
			// Value changed on schema level
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, providerModel, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionSqlResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasLogLevelString(string(sdk.LogLevelInfo)),
					// modifying the parameter as assertion, as this is the moment in-between apply and refresh
					invokeactionassert.Invoke(t, func() error {
						secondaryTestClient().Schema.UnsetLogLevel(t, schema.ID())
						return nil
					}),
				),
			},
			// Unset on schema level, fallback to database
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, providerModel, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionSqlResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasLogLevelString(string(sdk.LogLevelError)),
				),
			},
		},
	})
}
