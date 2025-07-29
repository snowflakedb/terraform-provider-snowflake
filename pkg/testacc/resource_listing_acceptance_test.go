//go:build !account_level_tests

package testacc

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
)

func TestAcc_Listing_Basic_Inlined(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	basicManifest, _ := testClient().Listing.BasicManifest(t)
	modelBasic := model.ListingWithInlineManifest("test", id.Name(), basicManifest)
	modelComplete := model.ListingWithInlineManifest("test", id.Name(), basicManifest)
	modelCompleteWithDifferentValues := model.ListingWithInlineManifest("test", id.Name(), basicManifest)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			// create without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelBasic.ResourceReference()),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()),
				),
			},
			// import without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedListingResource(t, helpers.EncodeResourceIdentifier(id)),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelComplete.ResourceReference()),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()),
				),
			},
			// import complete object
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedListingResource(t, helpers.EncodeResourceIdentifier(id)),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)),
				),
			},
			// alter
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelCompleteWithDifferentValues.ResourceReference()),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()),
				),
			},
			// change externally
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelCompleteWithDifferentValues.ResourceReference()),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelBasic.ResourceReference()),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()),
				),
			},
		},
	})
}

func TestAcc_Listing_Basic_FromStage(t *testing.T) {}

func TestAcc_Listing_Complete_Inlined(t *testing.T)   {}
func TestAcc_Listing_Complete_FromStage(t *testing.T) {}

func TestAcc_Listing_NewVersions_Inlined(t *testing.T)   {}
func TestAcc_Listing_NewVersions_FromStage(t *testing.T) {}

func TestAcc_Listing_Updates_Inlined(t *testing.T)   {}
func TestAcc_Listing_Updates_FromStage(t *testing.T) {}

func TestAcc_Listing_UpdateManifestSource(t *testing.T) {}
func TestAcc_Listing_Validations(t *testing.T)          {}
