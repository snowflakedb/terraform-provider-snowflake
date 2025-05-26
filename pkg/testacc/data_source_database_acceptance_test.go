//go:build !account_level_tests

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

func TestAcc_Database(t *testing.T) {
	databaseName := testClient().Ids.Alpha()
	comment := random.Comment()

	databaseModel := model.DatabaseWithParametersSet("test", databaseName).
		WithComment(comment)
	databaseDatasourceModel := datasourcemodel.Database("test", databaseName).
		WithDependsOn(databaseModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
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
