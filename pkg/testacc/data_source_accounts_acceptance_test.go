//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Accounts_BasicUseCase_DifferentFiltering(t *testing.T) {
	currentAccount := testClient().Context.CurrentAccountId(t)

	accountsWithAccountName := datasourcemodel.Accounts("test").
		WithLike(currentAccount.AccountName()).
		WithWithParameters(false)
	accountsWithNoMatch := datasourcemodel.Accounts("test").
		WithLike(random.AlphaUpperN(12)).
		WithWithParameters(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountsWithAccountName),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(accountsWithAccountName.DatasourceReference(), "accounts.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(accountsWithAccountName.DatasourceReference(), "accounts.0.parameters.#", "0")),
					resourceshowoutputassert.AccountsDatasourceShowOutput(t, accountsWithAccountName.DatasourceReference()).
						HasOrganizationName(currentAccount.OrganizationName()).
						HasAccountName(currentAccount.AccountName()).
						HasAccountLocator(testClient().GetAccountLocator()),
				),
			},
			{
				Config: config.FromModels(t, accountsWithNoMatch),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(accountsWithNoMatch.DatasourceReference(), "accounts.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Accounts_CompleteUseCase_Parameters(t *testing.T) {
	currentAccountName := testClient().Context.CurrentAccountName(t)

	accountsWithParameters := datasourcemodel.Accounts("test").
		WithLike(currentAccountName).
		WithWithParameters(true)
	accountsWithoutParameters := datasourcemodel.Accounts("test").
		WithLike(currentAccountName).
		WithWithParameters(false)

	parametersAssert := resourceparametersassert.AccountsDatasourceParameters(t, accountsWithParameters.DatasourceReference()).
		HasAllParametersPresent()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountsWithParameters),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(accountsWithParameters.DatasourceReference(), "accounts.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(accountsWithParameters.DatasourceReference(), "accounts.0.parameters.#", "1")),
					parametersAssert,
				),
			},
			{
				Config: config.FromModels(t, accountsWithoutParameters),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(accountsWithoutParameters.DatasourceReference(), "accounts.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(accountsWithoutParameters.DatasourceReference(), "accounts.0.parameters.#", "0")),
				),
			},
		},
	})
}
