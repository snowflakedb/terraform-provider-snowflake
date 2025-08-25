//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Warehouses_Complete(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	warehouseModel := model.Warehouse("test", id.Name()).
		WithComment(comment)
	warehousesModel := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithDependsOn(warehouseModel.ResourceReference())
	warehousesModelOptionalsUnset := datasourcemodel.Warehouses("test").
		WithWithDescribe(false).
		WithWithParameters(false).
		WithLike(id.Name()).
		WithDependsOn(warehouseModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, warehouseModel, warehousesModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.#", "1"),

					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.state"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.type", string(sdk.WarehouseTypeStandard)),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.size", string(sdk.WarehouseSizeXSmall)),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.min_cluster_count", "1"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.max_cluster_count", "1"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.started_clusters"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.running"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.queued"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.is_current", "true"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.auto_suspend", "600"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.auto_resume", "true"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.available"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.provisioning"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.quiescing"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.other"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.created_on"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.resumed_on"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.updated_on"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.enable_query_acceleration", "false"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.query_acceleration_max_scale_factor", "8"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.resource_monitor", ""),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.scaling_policy", string(sdk.ScalingPolicyStandard)),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.owner_role_type"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.show_output.0.resource_constraint"),

					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.#", "1"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.created_on"),
					resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.name"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.kind", "WAREHOUSE"),

					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.parameters.#", "1"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.parameters.0.max_concurrency_level.0.value", "8"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.parameters.0.statement_timeout_in_seconds.0.value", "172800"),
				),
			},
			{
				Config: accconfig.FromModels(t, warehouseModel, warehousesModelOptionalsUnset),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.#", "1"),

					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.state"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.type", string(sdk.WarehouseTypeStandard)),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.size", string(sdk.WarehouseSizeXSmall)),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.min_cluster_count", "1"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.max_cluster_count", "1"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.started_clusters"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.running"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.queued"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.is_current", "true"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.auto_suspend", "600"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.auto_resume", "true"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.available"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.provisioning"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.quiescing"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.other"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.created_on"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.resumed_on"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.updated_on"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.enable_query_acceleration", "false"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.query_acceleration_max_scale_factor", "8"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.resource_monitor", ""),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.scaling_policy", string(sdk.ScalingPolicyStandard)),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.owner_role_type"),
					resource.TestCheckResourceAttrSet(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.show_output.0.resource_constraint"),

					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.describe_output.#", "0"),
					resource.TestCheckResourceAttr(warehousesModelOptionalsUnset.DatasourceReference(), "warehouses.0.parameters.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Warehouses_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModel1 := model.Warehouse("test1", idOne.Name())
	warehouseModel2 := model.Warehouse("test2", idTwo.Name())
	warehouseModel3 := model.Warehouse("test3", idThree.Name())
	warehousesModelLikeFirstOne := datasourcemodel.Warehouses("test").
		WithLike(idOne.Name()).
		WithDependsOn(warehouseModel1.ResourceReference(), warehouseModel2.ResourceReference(), warehouseModel3.ResourceReference())
	warehousesModelLikePrefix := datasourcemodel.Warehouses("test").
		WithLike(prefix+"%").
		WithDependsOn(warehouseModel1.ResourceReference(), warehouseModel2.ResourceReference(), warehouseModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, warehouseModel1, warehouseModel2, warehouseModel3, warehousesModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehousesModelLikeFirstOne.DatasourceReference(), "warehouses.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, warehouseModel1, warehouseModel2, warehouseModel3, warehousesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehousesModelLikePrefix.DatasourceReference(), "warehouses.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Warehouses_WarehouseNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Warehouses/without_warehouse"),
				ExpectError:     regexp.MustCompile("there should be at least one warehouse"),
			},
		},
	})
}
