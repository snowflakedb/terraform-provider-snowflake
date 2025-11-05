//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRoles_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	databaseRoleId1 := testClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix)
	databaseRoleId2 := testClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix)
	databaseRoleId3 := testClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	databaseRoleModel1 := model.DatabaseRole("test", databaseRoleId1.DatabaseName(), databaseRoleId1.Name()).WithComment(comment)
	databaseRoleModel2 := model.DatabaseRole("test1", databaseRoleId2.DatabaseName(), databaseRoleId2.Name()).WithComment(comment)
	databaseRoleModel3 := model.DatabaseRole("test2", databaseRoleId3.DatabaseName(), databaseRoleId3.Name()).WithComment(comment)

	datasourceModelLikeExact := datasourcemodel.DatabaseRoles("test", databaseRoleId1.DatabaseName()).
		WithLike(databaseRoleId1.Name()).
		WithDependsOn(databaseRoleModel1.ResourceReference(), databaseRoleModel2.ResourceReference(), databaseRoleModel3.ResourceReference())

	datasourceModelLikePrefix := datasourcemodel.DatabaseRoles("test", databaseRoleId1.DatabaseName()).
		WithLike(prefix+"%").
		WithDependsOn(databaseRoleModel1.ResourceReference(), databaseRoleModel2.ResourceReference(), databaseRoleModel3.ResourceReference())

	datasourceModelInDatabase := datasourcemodel.DatabaseRoles("test", databaseRoleId3.DatabaseName()).
		WithInDatabase(databaseRoleId3.DatabaseName()).
		WithLike(databaseRoleId3.Name()).
		WithDependsOn(databaseRoleModel1.ResourceReference(), databaseRoleModel2.ResourceReference(), databaseRoleModel3.ResourceReference())

	datasourceModelLimitRows := datasourcemodel.DatabaseRoles("test", databaseRoleId1.DatabaseName()).
		WithRows(1).
		WithDependsOn(databaseRoleModel1.ResourceReference(), databaseRoleModel2.ResourceReference(), databaseRoleModel3.ResourceReference())

	datasourceModelLikePrefixWithLimit := datasourcemodel.DatabaseRoles("test", databaseRoleId1.DatabaseName()).
		WithLike(prefix+"%").
		WithRows(1).
		WithDependsOn(databaseRoleModel1.ResourceReference(), databaseRoleModel2.ResourceReference(), databaseRoleModel3.ResourceReference())

	datasourceModelLimitRowsFrom := datasourcemodel.DatabaseRoles("test", databaseRoleId1.DatabaseName()).
		WithRowsAndFrom(1, databaseRoleId1.Name()).
		WithDependsOn(databaseRoleModel1.ResourceReference(), databaseRoleModel2.ResourceReference(), databaseRoleModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DatabaseRole),
		Steps: []resource.TestStep{
			// like (exact)
			{
				Config: accconfig.FromModels(t, databaseRoleModel1, databaseRoleModel2, databaseRoleModel3, datasourceModelLikeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "database_roles.#", "1"),
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "database_roles.0.show_output.0.name", databaseRoleId1.Name()),
				),
			},
			// like (prefix)
			{
				Config: accconfig.FromModels(t, databaseRoleModel1, databaseRoleModel2, databaseRoleModel3, datasourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "database_roles.#", "2"),
				),
			},
			// explicit in_database filtering for role3's database
			{
				Config: accconfig.FromModels(t, databaseRoleModel1, databaseRoleModel2, databaseRoleModel3, datasourceModelInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInDatabase.DatasourceReference(), "database_roles.#", "1"),
					resource.TestCheckResourceAttr(datasourceModelInDatabase.DatasourceReference(), "database_roles.0.show_output.0.name", databaseRoleId3.Name()),
				),
			},
			// limit rows only (no from)
			{
				Config: accconfig.FromModels(t, databaseRoleModel1, databaseRoleModel2, databaseRoleModel3, datasourceModelLimitRows),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLimitRows.DatasourceReference(), "database_roles.#", "1"),
				),
			},
			// like + limit rows only
			{
				Config: accconfig.FromModels(t, databaseRoleModel1, databaseRoleModel2, databaseRoleModel3, datasourceModelLikePrefixWithLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefixWithLimit.DatasourceReference(), "database_roles.#", "1"),
				),
			},
			// limit rows with from
			{
				Config: accconfig.FromModels(t, databaseRoleModel1, databaseRoleModel2, databaseRoleModel3, datasourceModelLimitRowsFrom),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLimitRowsFrom.DatasourceReference(), "database_roles.#", "1"),
					resource.TestCheckResourceAttr(datasourceModelLimitRowsFrom.DatasourceReference(), "database_roles.0.show_output.0.name", databaseRoleId2.Name()),
				),
			},
		},
	})
}

func TestAcc_DatabaseRoles_CompleteUseCase(t *testing.T) {
	prefix := random.AlphaN(4)
	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix)
	comment := random.Comment()

	databaseRoleModel := model.DatabaseRole("test", databaseRoleId.DatabaseName(), databaseRoleId.Name()).WithComment(comment)

	databaseRolesModel := datasourcemodel.DatabaseRoles("test", databaseRoleId.DatabaseName()).
		WithInDatabase(databaseRoleId.DatabaseName()).
		WithLike(databaseRoleId.Name()).
		WithDependsOn(databaseRoleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DatabaseRole),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseRoleModel, databaseRolesModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resourceshowoutputassert.DatabaseRolesDatasourceShowOutput(t, databaseRolesModel.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(databaseRoleId.Name()).
						HasDatabaseName(databaseRoleId.DatabaseName()).
						HasIsDefault(false).
						HasIsCurrent(false).
						HasIsInherited(false).
						HasGrantedToRoles(0).
						HasGrantedToDatabaseRoles(0).
						HasGrantedDatabaseRoles(0).
						HasOwnerNotEmpty().
						HasComment(comment).
						HasOwnerRoleTypeNotEmpty().
						ToTerraformTestCheckFunc(t, testClient()),
				),
			},
		},
	})
}
