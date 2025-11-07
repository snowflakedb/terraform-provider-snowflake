//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Listings_BasicUseCase_DifferentFiltering(t *testing.T) {
	// We use a random prefix that does not exist to make assertions deterministic.
	prefix := "tf_acc_no_such_listing_" + random.AlphaN(8)

	listingsLike := datasourcemodel.Listings("test").
		WithLike(prefix)
	listingsStartsWith := datasourcemodel.Listings("test").
		WithStartsWith(prefix)
	listingsLimit := datasourcemodel.Listings("test").
		WithRowsAndFrom(1, "").
		WithStartsWith(prefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, listingsLike),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(listingsLike.DatasourceReference(), "listings.#", "0"),
				),
			},
			{
				Config: accconfig.FromModels(t, listingsStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(listingsStartsWith.DatasourceReference(), "listings.#", "0"),
				),
			},
			{
				Config: accconfig.FromModels(t, listingsLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(listingsLimit.DatasourceReference(), "listings.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Listings_CompleteUseCase(t *testing.T) {
	t.Skip("Skipping: requires existing Marketplace listings to assert fields; enable when environment has listings.")
}
