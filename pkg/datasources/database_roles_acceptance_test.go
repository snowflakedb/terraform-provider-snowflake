//go:build !account_level_tests

package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRoles(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaN(4)
	databaseRoleId1 := acc.TestClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix + "1")
	databaseRoleId2 := acc.TestClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix + "2")
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
