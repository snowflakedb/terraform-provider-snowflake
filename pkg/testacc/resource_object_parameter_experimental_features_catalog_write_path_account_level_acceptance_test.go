//go:build account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Experimental_ObjectParameter_CatalogWritePath(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	task, taskCleanup := testClient().Task.Create(t)
	t.Cleanup(taskCleanup)

	hybridTableId, hybridTableCleanup := testClient().HybridTable.Create(t)
	t.Cleanup(hybridTableCleanup)

	icebergTable, icebergTableCleanup := testClient().IcebergTable.Create(t)
	t.Cleanup(icebergTableCleanup)

	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	service, serviceCleanup := testClient().Service.Create(t, computePool.ID())
	t.Cleanup(serviceCleanup)

	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)

	t.Cleanup(func() { testClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterDefaultDdlCollation) })

	accountModel := model.ObjectParameter("account_scope", string(sdk.ObjectParameterDefaultDdlCollation), "en_US").
		WithOnAccount(true)

	objectParameterCatalogWritePathModel := func(name string, parameter sdk.ObjectParameter, value string, objectType sdk.ObjectType, id sdk.ObjectIdentifier) *model.ObjectParameterModel {
		identifier := map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(id.Name()),
		}
		switch id := id.(type) {
		case sdk.DatabaseObjectIdentifier:
			identifier["database"] = tfconfig.StringVariable(id.DatabaseName())
		case sdk.SchemaObjectIdentifier:
			identifier["database"] = tfconfig.StringVariable(id.DatabaseName())
			identifier["schema"] = tfconfig.StringVariable(id.SchemaName())
		}
		return model.ObjectParameter(name, string(parameter), value).
			WithObjectType(objectType.String()).
			WithObjectIdentifierValue(tfconfig.ListVariable(tfconfig.ObjectVariable(identifier)))
	}
	databaseModel := objectParameterCatalogWritePathModel("database_scope", sdk.ObjectParameterDefaultDdlCollation, "en_US", sdk.ObjectTypeDatabase, database.ID())
	schemaModel := objectParameterCatalogWritePathModel("schema_scope", sdk.ObjectParameterDefaultDdlCollation, "en_US", sdk.ObjectTypeSchema, schema.ID())
	warehouseModel := objectParameterCatalogWritePathModel("warehouse", sdk.ObjectParameterMaxConcurrencyLevel, "9", sdk.ObjectTypeWarehouse, warehouse.ID())
	userModel := objectParameterCatalogWritePathModel("user", sdk.ObjectParameterEnableUnredactedQuerySyntaxError, "true", sdk.ObjectTypeUser, user.ID())
	taskModel := objectParameterCatalogWritePathModel("task", sdk.ObjectParameter("QUERY_TAG"), "catalog-path", sdk.ObjectTypeTask, task.ID())
	hybridTableModel := objectParameterCatalogWritePathModel("hybrid_table", sdk.ObjectParameterDataRetentionTimeInDays, "2", sdk.ObjectTypeHybridTable, hybridTableId)
	icebergTableModel := objectParameterCatalogWritePathModel("iceberg_table", sdk.ObjectParameterDefaultDdlCollation, "en_US", sdk.ObjectTypeIcebergTable, icebergTable.ID())
	serviceModel := objectParameterCatalogWritePathModel("service", sdk.ObjectParameterServiceCallerTokenValiditySecs, "120", sdk.ObjectTypeService, service.ID())
	fallbackModel := objectParameterCatalogWritePathModel("generic_fallback", sdk.ObjectParameterDataRetentionTimeInDays, "2", sdk.ObjectTypeTable, table.ID())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config: accconfig.FromModels(t,
					providermodel.SnowflakeProvider().WithPreviewFeaturesEnabled(string(previewfeatures.ObjectParameterResource)),
					accountModel, databaseModel, schemaModel),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(accountModel.ResourceReference(), plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(databaseModel.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(databaseModel.ResourceReference(), tfjsonpath.New("value"), knownvalue.StringExact("en_US")),
						plancheck.ExpectResourceAction(schemaModel.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(schemaModel.ResourceReference(), tfjsonpath.New("value"), knownvalue.StringExact("en_US")),
					},
				},
				Check: assertThat(
					t,
					resourceassert.ObjectParameterResource(t, accountModel.ResourceReference()).HasValueString("en_US"),
					resourceassert.ObjectParameterResource(t, databaseModel.ResourceReference()).HasValueString("EN_US"),
					resourceassert.ObjectParameterResource(t, schemaModel.ResourceReference()).HasValueString("EN_US"),
					objectparametersassert.AccountParameters(t, testClient().Context.CurrentAccountId(t)).HasDefaultDdlCollation("en_US"),
					objectparametersassert.DatabaseParameters(t, database.ID()).HasDefaultDdlCollation("EN_US"),
					objectparametersassert.SchemaParameters(t, schema.ID()).HasDefaultDdlCollation("EN_US"),
				),
			},
			{
				ProtoV6ProviderFactories: objectParameterParseOnSetProviderFactory,
				Config: accconfig.FromModels(t,
					providermodel.SnowflakeProvider().
						WithPreviewFeaturesEnabled(string(previewfeatures.ObjectParameterResource)).
						WithExperimentalFeaturesEnabled(experimentalfeatures.ObjectParameterParseOnSet),
					accountModel,
					databaseModel,
					schemaModel,
					warehouseModel,
					userModel,
					taskModel,
					hybridTableModel,
					icebergTableModel,
					serviceModel,
					fallbackModel,
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(
					t,
					resourceassert.ObjectParameterResource(t, accountModel.ResourceReference()).HasValueString("en_US"),
					resourceassert.ObjectParameterResource(t, databaseModel.ResourceReference()).HasValueString("en_US"),
					resourceassert.ObjectParameterResource(t, schemaModel.ResourceReference()).HasValueString("en_US"),
					resourceassert.ObjectParameterResource(t, warehouseModel.ResourceReference()).HasValueString("9"),
					resourceassert.ObjectParameterResource(t, userModel.ResourceReference()).HasValueString("true"),
					resourceassert.ObjectParameterResource(t, taskModel.ResourceReference()).HasValueString("catalog-path"),
					resourceassert.ObjectParameterResource(t, hybridTableModel.ResourceReference()).HasValueString("2"),
					resourceassert.ObjectParameterResource(t, icebergTableModel.ResourceReference()).HasValueString("en_US"),
					resourceassert.ObjectParameterResource(t, serviceModel.ResourceReference()).HasValueString("120"),
					resourceassert.ObjectParameterResource(t, fallbackModel.ResourceReference()).HasValueString("2"),
					objectparametersassert.AccountParameters(t, testClient().Context.CurrentAccountId(t)).HasDefaultDdlCollation("en_US"),
					objectparametersassert.DatabaseParameters(t, database.ID()).HasDefaultDdlCollation("en_US"),
					objectparametersassert.SchemaParameters(t, schema.ID()).HasDefaultDdlCollation("en_US"),
				),
			},
		},
	})
}
