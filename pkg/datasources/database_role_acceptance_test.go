//go:build !account_level_tests

package datasources_test

import (
	"regexp"
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

func TestAcc_DatabaseRole(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	databaseRoleModel := model.DatabaseRole("test", databaseRoleId.DatabaseName(), databaseRoleId.Name()).
		WithComment(comment)
	databaseRoleDatasourceModel := datasourcemodel.DatabaseRole("test", databaseRoleId.DatabaseName(), databaseRoleId.Name()).
		WithDependsOn(databaseRoleModel.ResourceReference())
	databaseRoleNotExistingDatasourceModel := datasourcemodel.DatabaseRole("test", acc.NonExistingDatabaseObjectIdentifier.DatabaseName(), acc.NonExistingDatabaseObjectIdentifier.Name()).
		WithDependsOn(databaseRoleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseRoleModel, databaseRoleDatasourceModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(databaseRoleDatasourceModel.DatasourceReference(), "name"),
					resource.TestCheckResourceAttrSet(databaseRoleDatasourceModel.DatasourceReference(), "comment"),
					resource.TestCheckResourceAttrSet(databaseRoleDatasourceModel.DatasourceReference(), "owner"),
				),
			},
			{
				Config:      accconfig.FromModels(t, databaseRoleModel, databaseRoleNotExistingDatasourceModel),
				ExpectError: regexp.MustCompile("Error: object does not exist"),
			},
		},
	})
}
