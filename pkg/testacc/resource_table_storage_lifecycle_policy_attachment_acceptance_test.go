//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TableStorageLifecyclePolicyAttachment_Table(t *testing.T) {
	table, tableCleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("C1", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("C2", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)
	tableName := table.ID().FullyQualifiedName()

	table2, tableCleanup2 := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("C1", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("C2", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup2)
	tableName2 := table2.ID().FullyQualifiedName()

	policyId := testClient().Ids.RandomSchemaObjectIdentifier()
	policyCleanup := testClient().StorageLifecyclePolicy.CreateWithRequest(t, policyId, sdk.NewCreateStorageLifecyclePolicyRequest(
		policyId,
		[]sdk.CreateStorageLifecyclePolicyArgsRequest{
			*sdk.NewCreateStorageLifecyclePolicyArgsRequest("C1", testdatatypes.DataTypeNumber),
			*sdk.NewCreateStorageLifecyclePolicyArgsRequest("C2", testdatatypes.DataTypeNumber),
		},
		"C1 > 0",
	))
	t.Cleanup(policyCleanup)
	policyName := policyId.FullyQualifiedName()

	policyId2 := testClient().Ids.RandomSchemaObjectIdentifier()
	policyCleanup2 := testClient().StorageLifecyclePolicy.CreateWithRequest(t, policyId2, sdk.NewCreateStorageLifecyclePolicyRequest(
		policyId2,
		[]sdk.CreateStorageLifecyclePolicyArgsRequest{*sdk.NewCreateStorageLifecyclePolicyArgsRequest("C2", testdatatypes.DataTypeNumber)},
		"C2 > 0",
	))
	t.Cleanup(policyCleanup2)
	policyName2 := policyId2.FullyQualifiedName()

	tableType := string(sdk.PolicyEntityDomainTable)

	basic := model.TableStorageLifecyclePolicyAttachment("t", []string{"C1", "C2"}, policyName, tableName, tableType)

	changedColumnsAndPolicy := model.TableStorageLifecyclePolicyAttachment("t", []string{"C2"}, policyName2, tableName, tableType)

	changedTable := model.TableStorageLifecyclePolicyAttachment("t", []string{"C2"}, policyName2, tableName2, tableType)

	ref := basic.ResourceReference()

	basicAssertions := resourceassert.TableStorageLifecyclePolicyAttachmentResource(t, ref).
		HasTableName(tableName).
		HasTableType(tableType).
		HasStorageLifecyclePolicyName(policyName).
		HasOn("C1", "C2")

	changedColumnsAndPolicyAssertions := resourceassert.TableStorageLifecyclePolicyAttachmentResource(t, ref).
		HasTableName(tableName).
		HasTableType(tableType).
		HasStorageLifecyclePolicyName(policyName2).
		HasOn("C2")

	changedTableAssertions := resourceassert.TableStorageLifecyclePolicyAttachmentResource(t, ref).
		HasTableName(tableName2).
		HasTableType(tableType).
		HasStorageLifecyclePolicyName(policyName2).
		HasOn("C2")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: warehouseRequiredProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckTableStorageLifecyclePolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions),
			},
			// Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change columns and policy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, changedColumnsAndPolicy),
				Check:  assertThat(t, changedColumnsAndPolicyAssertions),
			},
			// Change table
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, changedTable),
				Check:  assertThat(t, changedTableAssertions),
			},
			// Drop table externally and trigger destroy - expect empty plan
			{
				PreConfig: func() {
					testClient().Table.DropFunc(t, table2.ID())()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config:  config.FromModels(t, changedTable),
				Destroy: true,
			},
			{
				Config: config.FromModels(t, basic),
			},
			// Unset policy externally - expect attachment to be recreated
			{
				PreConfig: func() {
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(table.ID()).WithDropStorageLifecyclePolicy(new(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions),
			},
		},
	})
}

func TestAcc_TableStorageLifecyclePolicyAttachment_DynamicTable(t *testing.T) {
	sourceTable, sourceTableCleanup := testClient().Table.Create(t)
	t.Cleanup(sourceTableCleanup)

	dynamicTable, dynamicTableCleanup := testClient().DynamicTable.CreateDynamicTable(t, sourceTable.ID())
	t.Cleanup(dynamicTableCleanup)
	dynamicTableName := dynamicTable.ID().FullyQualifiedName()

	policyId := testClient().Ids.RandomSchemaObjectIdentifier()
	policyCleanup := testClient().StorageLifecyclePolicy.CreateWithRequest(t, policyId, sdk.NewCreateStorageLifecyclePolicyRequest(
		policyId,
		[]sdk.CreateStorageLifecyclePolicyArgsRequest{*sdk.NewCreateStorageLifecyclePolicyArgsRequest("ID", testdatatypes.DataTypeNumber)},
		"ID > 0",
	))
	t.Cleanup(policyCleanup)
	policyName := policyId.FullyQualifiedName()

	tableType := string(sdk.PolicyEntityDomainDynamicTable)

	basic := model.TableStorageLifecyclePolicyAttachment("t", []string{"ID"}, policyName, dynamicTableName, tableType)

	ref := basic.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: warehouseRequiredProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckTableStorageLifecyclePolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check: assertThat(t,
					resourceassert.TableStorageLifecyclePolicyAttachmentResource(t, ref).
						HasTableName(dynamicTableName).
						HasTableType(tableType).
						HasStorageLifecyclePolicyName(policyName).
						HasOn("ID"),
				),
			},
			// Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_TableStorageLifecyclePolicyAttachment_Validations(t *testing.T) {
	tableName := testClient().Ids.RandomSchemaObjectIdentifier().FullyQualifiedName()
	policyName := testClient().Ids.RandomSchemaObjectIdentifier().FullyQualifiedName()

	invalidTableType := model.TableStorageLifecyclePolicyAttachment("t", []string{"COLUMN_1"}, policyName, tableName, "INVALID")
	emptyOn := model.TableStorageLifecyclePolicyAttachment("t", []string{}, policyName, tableName, string(sdk.PolicyEntityDomainTable))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: warehouseRequiredProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidTableType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected table_type to be one of .* got INVALID`),
			},
			{
				Config:      config.FromModels(t, emptyOn),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute on requires 1 item minimum`),
			},
		},
	})
}
