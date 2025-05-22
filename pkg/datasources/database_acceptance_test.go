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

func TestAcc_Database(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()

	databaseModel := model.DatabaseWithParametersSet("test", databaseName).
		WithComment(comment)
	databaseDatasourceModel := datasourcemodel.Database("test", databaseName).
		WithDependsOn(databaseModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseModel, databaseDatasourceModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databaseDatasourceModel.DatasourceReference(), "name", databaseName),
					resource.TestCheckResourceAttr(databaseDatasourceModel.DatasourceReference(), "comment", comment),
					resource.TestCheckResourceAttrSet(databaseDatasourceModel.DatasourceReference(), "created_on"),
					resource.TestCheckResourceAttrSet(databaseDatasourceModel.DatasourceReference(), "owner"),
					resource.TestCheckResourceAttrSet(databaseDatasourceModel.DatasourceReference(), "retention_time"),
					resource.TestCheckResourceAttrSet(databaseDatasourceModel.DatasourceReference(), "is_current"),
					resource.TestCheckResourceAttrSet(databaseDatasourceModel.DatasourceReference(), "is_default"),
				),
			},
		},
	})
}
