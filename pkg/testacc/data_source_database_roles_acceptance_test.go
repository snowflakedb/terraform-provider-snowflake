//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRoles(t *testing.T) {
	prefix := random.AlphaN(4)
	databaseRoleId1 := testClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix + "1")
	databaseRoleId2 := testClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix + "2")
	comment := random.Comment()

	databaseRoleModel1 := model.DatabaseRole("test", databaseRoleId1.DatabaseName(), databaseRoleId1.Name()).
		WithComment(comment)
	databaseRoleModel2 := model.DatabaseRole("test2", databaseRoleId2.DatabaseName(), databaseRoleId2.Name()).
		WithComment(comment)
	databaseRolesDatasourceModel := datasourcemodel.DatabaseRoles("test", databaseRoleId1.DatabaseName()).
		WithLike(prefix+"%").
		WithRowsAndFrom(1, databaseRoleId1.Name()).
		WithDependsOn(databaseRoleModel1.ResourceReference(), databaseRoleModel2.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseRoleModel1, databaseRoleModel2, databaseRolesDatasourceModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.#", "1"),
					resource.TestCheckResourceAttrSet(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.name", databaseRoleId2.Name()),
					resource.TestCheckResourceAttr(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttrSet(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.is_current"),
					resource.TestCheckResourceAttr(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.granted_to_database_roles", "0"),
					resource.TestCheckResourceAttr(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.granted_database_roles", "0"),
					resource.TestCheckResourceAttrSet(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(databaseRolesDatasourceModel.DatasourceReference(), "database_roles.0.show_output.0.owner_role_type"),
				),
			},
		},
	})
}
