//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_WarehouseInteractive_Basic(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	warehouseModel := model.WarehouseInteractiveWithId(warehouseId)
	warehouseModelWithComment := model.WarehouseInteractiveWithId(warehouseId).
		WithComment(comment).
		WithWarehouseSize(string(sdk.WarehouseSizeXSmall))

	ref := warehouseModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			// create with only required fields
			{
				Config: accconfig.FromModels(t, warehouseModel),
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasNameString(warehouseId.Name()).
						HasWarehouseTypeString(string(sdk.WarehouseTypeInteractive)).
						HasFullyQualifiedNameString(warehouseId.FullyQualifiedName()),
					resourceshowoutputassert.WarehouseShowOutput(t, ref).
						HasName(warehouseId.Name()).
						HasStateNotEmpty(),
				),
			},
			// import after minimal config
			{
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_resume"},
			},
			// set optional fields
			{
				Config: accconfig.FromModels(t, warehouseModelWithComment),
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasCommentString(comment).
						HasWarehouseSizeString(string(sdk.WarehouseSizeXSmall)),
					resourceshowoutputassert.WarehouseShowOutput(t, ref).
						HasComment(comment),
				),
			},
		},
	})
}

func TestAcc_WarehouseInteractive_TablesDelta(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()

	table1, table1Cleanup := testClient().Table.CreateInteractiveTable(t)
	t.Cleanup(table1Cleanup)
	table2, table2Cleanup := testClient().Table.CreateInteractiveTable(t)
	t.Cleanup(table2Cleanup)
	table3, table3Cleanup := testClient().Table.CreateInteractiveTable(t)
	t.Cleanup(table3Cleanup)

	modelWithTwoTables := model.WarehouseInteractiveWithId(warehouseId).
		WithTables(table1.FullyQualifiedName(), table2.FullyQualifiedName())
	// Drop table2, add table3 — should result in a single ADD + single DROP, no full replace.
	modelWithSwappedTable := model.WarehouseInteractiveWithId(warehouseId).
		WithTables(table1.FullyQualifiedName(), table3.FullyQualifiedName())

	ref := modelWithTwoTables.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithTwoTables),
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasTables(table1.FullyQualifiedName(), table2.FullyQualifiedName()),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithSwappedTable),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasTables(table1.FullyQualifiedName(), table3.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_WarehouseInteractive_FallbackWarehouse(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()

	fallback, fallbackCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(fallbackCleanup)

	modelWithFallback := model.WarehouseInteractiveWithId(warehouseId).
		WithFallbackWarehouse(fallback.ID().Name())
	modelWithoutFallback := model.WarehouseInteractiveWithId(warehouseId)

	ref := modelWithFallback.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithFallback),
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasFallbackWarehouseString(fallback.ID().Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithoutFallback),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasFallbackWarehouseEmpty(),
				),
			},
		},
	})
}

func TestAcc_WarehouseInteractive_Import_WrongWarehouseType(t *testing.T) {
	interactiveId := testClient().Ids.RandomAccountObjectIdentifier()
	regularId := testClient().Ids.RandomAccountObjectIdentifier()

	// Create a regular (non-interactive) warehouse outside of Terraform to use as the import target.
	_, regularCleanup := testClient().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(regularId))
	t.Cleanup(regularCleanup)

	interactiveModel := model.WarehouseInteractiveWithId(interactiveId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, interactiveModel),
			},
			{
				ResourceName:  interactiveModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: regularId.Name(),
				ExpectError:   regexp.MustCompile("is not an interactive warehouse"),
			},
		},
	})
}

func TestAcc_WarehouseInteractive_Validations(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()

	modelInvalidAutoSuspend := model.WarehouseInteractiveWithId(warehouseId).
		WithAutoSuspend(100)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidAutoSuspend),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected auto_suspend to be at least \(86400\), got 100`),
			},
		},
	})
}
