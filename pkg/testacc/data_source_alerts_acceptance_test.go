//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Alerts_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	warehouseId := testClient().Ids.WarehouseId()
	comment := random.Comment()

	alertModel := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 2", "SELECT 1").
		WithWarehouse(warehouseId.Name()).
		WithAlertScheduleInterval(60).
		WithComment(comment)

	alertsModel := datasourcemodel.Alerts("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(alertModel.ResourceReference())

	alertsModelWithoutDescribe := datasourcemodel.Alerts("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(alertModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.AlertsDatasourceShowOutput(t, alertsModel.DatasourceReference()).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasComment(comment).
		HasWarehouse(warehouseId).
		HasSchedule("60 MINUTE").
		HasState(sdk.AlertStateSuspended).
		HasCondition("SELECT 1").
		HasAction("SELECT 2").
		HasOwnerRoleType("ROLE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, alertModel, alertsModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(alertsModel.DatasourceReference(), "alerts.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.AlertsDatasourceDescribeOutput(t, alertsModel.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment(comment).
						HasWarehouse(warehouseId).
						HasSchedule("60 MINUTE").
						HasState(string(sdk.AlertStateSuspended)).
						HasCondition("SELECT 1").
						HasAction("SELECT 2"),
				),
			},
			{
				Config: accconfig.FromModels(t, alertModel, alertsModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(alertsModelWithoutDescribe.DatasourceReference(), "alerts.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(alertsModelWithoutDescribe.DatasourceReference(), "alerts.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Alerts_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, secondSchema.ID())

	model1 := model.Alert("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), "SELECT 1", "SELECT 1").WithAlertScheduleInterval(60)
	model2 := model.Alert("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), "SELECT 1", "SELECT 1").WithAlertScheduleInterval(60)
	model3 := model.Alert("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), "SELECT 1", "SELECT 1").WithAlertScheduleInterval(60)

	likeFirst := datasourcemodel.Alerts("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	likePrefix := datasourcemodel.Alerts("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inSchema := datasourcemodel.Alerts("test").
		WithInSchema(id1.SchemaId()).
		WithStartsWith(prefix).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inDatabase := datasourcemodel.Alerts("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	withLimit := datasourcemodel.Alerts("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithLimit(1).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likeFirst),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likeFirst.DatasourceReference(), "alerts.#", "1")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likePrefix),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "alerts.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inSchema),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inSchema.DatasourceReference(), "alerts.#", "2")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inDatabase),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inDatabase.DatasourceReference(), "alerts.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, withLimit),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(withLimit.DatasourceReference(), "alerts.#", "1")),
			},
		},
	})
}

func TestAcc_Alerts_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.Alerts("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Alerts_InApplication(t *testing.T) {
	applicationPackage, cleanup := testClient().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(cleanup)

	inApplicationPackage := datasourcemodel.Alerts("test").
		WithInApplicationPackage(applicationPackage.ID())
	inMissingApplication := datasourcemodel.Alerts("test").
		WithInApplication(testClient().Ids.RandomAccountObjectIdentifier())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, inApplicationPackage),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inApplicationPackage.DatasourceReference(), "alerts.#", "0")),
			},
			{
				Config:      accconfig.FromModels(t, inMissingApplication),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_Alerts_EmptyResult(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	alertsModel := datasourcemodel.Alerts("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, alertsModel),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(alertsModel.DatasourceReference(), "alerts.#", "0")),
			},
		},
	})
}
