//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountRoles_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	accountRole1 := model.AccountRole("primary", id1.Name())
	accountRole2 := model.AccountRole("secondary", id2.Name())

	datasourceModelLikeExact := datasourcemodel.AccountRoles("test").
		WithLike(id1.Name()).
		WithDependsOn(accountRole1.ResourceReference(), accountRole2.ResourceReference())

	datasourceModelLikePrefix := datasourcemodel.AccountRoles("test").
		WithLike(prefix+"%").
		WithDependsOn(accountRole1.ResourceReference(), accountRole2.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			// like (prefix)
			{
				Config: config.FromModels(t, accountRole1, accountRole2, datasourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "account_roles.#", "2"),
				),
			},
			// like (exact)
			{
				Config: config.FromModels(t, accountRole1, accountRole2, datasourceModelLikeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "account_roles.#", "1"),
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "account_roles.0.show_output.0.name", id1.Name()),
				),
			},
		},
	})
}

func TestAcc_AccountRoles_CompleteUseCase(t *testing.T) {
	roleName := testClient().Ids.AlphaWithPrefix(random.AlphaN(10))
	comment := random.Comment()

	accountRoleModel := model.AccountRole("role", roleName).WithComment(comment)

	accountRolesModel := datasourcemodel.AccountRoles("test").
		WithLike(roleName).
		WithDependsOn(accountRoleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountRoleModel, accountRolesModel),
				Check: assertThat(t,
					resourceshowoutputassert.AccountRolesDatasourceShowOutput(t, accountRolesModel.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(roleName).
						HasIsDefault(false).
						HasIsCurrent(false).
						HasComment(comment),
				),
			},
		},
	})
}
