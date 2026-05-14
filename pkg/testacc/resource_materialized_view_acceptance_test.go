//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_MaterializedView(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	viewId := testClient().Ids.RandomSchemaObjectIdentifier()
	viewName := viewId.Name()

	query := fmt.Sprintf(`SELECT ID, DATA FROM "%s"`, tableId.Name())
	otherQuery := fmt.Sprintf(`SELECT ID, DATA FROM "%s" WHERE ID LIKE 'foo%%'`, tableId.Name())

	comment := random.Comment()
	otherComment := random.Comment()

	tableModel := model.TableWithId("test", tableId, []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber},
		{Name: "DATA", Type: testdatatypes.DataTypeVarchar},
	})

	modelBasic := model.MaterializedViewWithId("test", viewId, query, TestWarehouseName).
		WithComment(comment).
		WithIsSecure(true).
		WithOrReplace(false).
		WithDependsOn(tableModel.ResourceReference())

	modelUpdatedParams := model.MaterializedViewWithId("test", viewId, query, TestWarehouseName).
		WithComment(otherComment).
		WithIsSecure(false).
		WithOrReplace(false).
		WithDependsOn(tableModel.ResourceReference())

	modelUpdatedStatement := model.MaterializedViewWithId("test", viewId, otherQuery, TestWarehouseName).
		WithComment(otherComment).
		WithIsSecure(false).
		WithOrReplace(false).
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: viewsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaterializedView),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, modelBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "is_secure", "true"),
				),
			},
			// update parameters
			{
				Config: accconfig.FromModels(t, tableModel, modelUpdatedParams),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", otherComment),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "is_secure", "false"),
				),
			},
			// change statement
			{
				Config: accconfig.FromModels(t, tableModel, modelUpdatedStatement),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", otherComment),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "is_secure", "false"),
				),
			},
			// change statement externally
			{
				PreConfig: func() {
					testClient().MaterializedView.CreateMaterializedViewWithName(t, viewId, query, true)
				},
				Config: accconfig.FromModels(t, tableModel, modelUpdatedStatement),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", otherComment),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "is_secure", "false"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_materialized_view.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"or_replace", "warehouse"},
			},
		},
	})
}

func TestAcc_MaterializedView_Tags(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	viewId := testClient().Ids.RandomSchemaObjectIdentifier()
	tag1Id := testClient().Ids.RandomSchemaObjectIdentifier()
	tag2Id := testClient().Ids.RandomSchemaObjectIdentifier()

	query := fmt.Sprintf(`SELECT ID FROM "%s"`, tableId.Name())

	tableModel := model.TableWithId("test", tableId, []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber},
	})

	tagModel1 := model.TagBase("test_tag", tag1Id)
	tagModel2 := model.TagBase("test_tag_2", tag2Id)

	modelWithTag1 := model.MaterializedViewWithId("test", viewId, query, TestWarehouseName).
		WithTagReference(tagModel1.ResourceReference(), "some_value").
		WithDependsOn(tableModel.ResourceReference())

	modelWithTag2 := model.MaterializedViewWithId("test", viewId, query, TestWarehouseName).
		WithTagReference(tagModel2.ResourceReference(), "some_value").
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: viewsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaterializedView),
		Steps: []resource.TestStep{
			// create tags
			{
				Config: accconfig.FromModels(t, tableModel, tagModel1, tagModel2, modelWithTag1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewId.Name()),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.#", "1"),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.0.name", tag1Id.Name()),
				),
			},
			// update tags
			{
				Config: accconfig.FromModels(t, tableModel, tagModel1, tagModel2, modelWithTag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewId.Name()),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.#", "1"),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.0.name", tag2Id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_materialized_view.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"or_replace", "warehouse", "tag"},
			},
		},
	})
}

func TestAcc_MaterializedView_Rename(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	viewId := testClient().Ids.RandomSchemaObjectIdentifier()
	newViewId := testClient().Ids.RandomSchemaObjectIdentifier()

	query := fmt.Sprintf(`SELECT ID FROM "%s"`, tableId.Name())
	comment := random.Comment()

	tableModel := model.TableWithId("test", tableId, []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber},
	})

	modelBasic := model.MaterializedViewWithId("test", viewId, query, TestWarehouseName).
		WithComment(comment).
		WithIsSecure(true).
		WithOrReplace(false).
		WithDependsOn(tableModel.ResourceReference())

	modelRenamed := model.MaterializedViewWithId("test", newViewId, query, TestWarehouseName).
		WithComment(comment).
		WithIsSecure(false).
		WithOrReplace(false).
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: viewsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaterializedView),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, modelBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewId.Name()),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "fully_qualified_name", viewId.FullyQualifiedName()),
				),
			},
			// rename with one param change
			{
				Config: accconfig.FromModels(t, tableModel, modelRenamed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", newViewId.Name()),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "fully_qualified_name", newViewId.FullyQualifiedName()),
				),
			},
		},
	})
}
